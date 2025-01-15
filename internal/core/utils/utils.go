package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	mathRandom "math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/nyaruka/phonenumbers"
	binder "github.com/oapi-codegen/runtime"
	"github.com/oklog/ulid"
)

func GenerateUniqueId() string {
	newUlid, err := ulid.New(ulid.Now(), nil)
	if err != nil {
		panic(err)
	}
	return newUlid.String()
}

func ParseUlid(id string) uint64 {
	parsedUlid, err := ulid.Parse(id)
	if err != nil {
		panic(err)
	}

	return parsedUlid.Time()
}

func GenerateOtp() string {
	mathRandom.Seed(time.Now().UnixNano())
	min := 100000
	max := 999999
	otp := mathRandom.Intn(max-min+1) + min
	return strconv.Itoa(otp)
}

func BindQueryParams(context echo.Context, dest interface{}) error {
	// Iterate through the fields of the destination struct
	typeOfDest := reflect.TypeOf(dest).Elem()

	valueOfDest := reflect.ValueOf(dest).Elem()
	for i := 0; i < typeOfDest.NumField(); i++ {
		field := typeOfDest.Field(i)
		fieldTag := field.Tag
		structFieldVal := valueOfDest.Field(i)
		// Check if the field has a 'query' tag
		paramName := fieldTag.Get("form")

		if paramName == "" {
			continue
		}

		contactsOmitempty := strings.Contains(paramName, "omitempty")
		// Determine if the parameter is required or optional
		required := !contactsOmitempty
		if contactsOmitempty {
			paramName = strings.Split(paramName, ",")[0]
		}

		// Bind the query parameter to the field
		err := binder.BindQueryParameter("form", true, required, paramName, context.QueryParams(), structFieldVal.Addr().Interface())
		if err != nil {
			return err
		}
	}

	return nil
}

func IsValidEmail(email string) bool {
	pattern := `^(([^<>()[\].,;:\s@"]+(\.[^<>()[\].,;:\s@"]+)*)|(".+"))@(([^<>()[\].,;:\s@"]+\.)+[^<>()[\].,;:\s@"]{2,})$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func ParsePhoneNumber(phoneNumber string) (*phonenumbers.PhoneNumber, error) {
	parsedPhoneNumber := phonenumbers.PhoneNumber{}
	err := phonenumbers.ParseAndKeepRawInputToNumber(phoneNumber, "IN", &parsedPhoneNumber)

	if err != nil {
		return nil, err
	}

	return &parsedPhoneNumber, err
}

func EnumExpression(value string) StringExpression {
	return RawString(strings.Join([]string{"'", value, "'"}, ""))
}

func GenerateWebsocketEventId() string {
	return uuid.NewString()
}

type WebhookSecretData struct {
	WhatsappBusinessAccountId string `json:"whatsapp_business_account_id"`
	OrganizationId            string `json:"organization_id"`
}

// generateUniqueWebhookSecret returns an encrypted token that includes the
// WhatsAppBusinessAccountId and the organizationId. The token is opaque to
// external parties, but you can decrypt it internally to recover the data.
func GenerateUniqueWebhookSecret(whatsappBusinessAccountId, organizationId, encryptionKey string) (string, error) {
	// 1. Construct the data struct
	secretData := WebhookSecretData{
		WhatsappBusinessAccountId: whatsappBusinessAccountId,
		OrganizationId:            organizationId,
	}

	// 2. Serialize to JSON
	plaintextJSON, err := json.Marshal(secretData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal secret data to JSON: %w", err)
	}

	// 3. Encrypt (AES-256 in GCM mode)
	block, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM cipher: %w", err)
	}

	// Generate a random nonce of the correct length
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the plaintext JSON
	ciphertext := aesGCM.Seal(nil, nonce, plaintextJSON, nil)

	// 4. Combine nonce + ciphertext, then base64-encode
	//    A typical approach is: final = nonce || ciphertext
	finalBytes := append(nonce, ciphertext...)
	token := base64.URLEncoding.EncodeToString(finalBytes)

	return token, nil
}

// decryptWebhookSecret does the reverse of generateUniqueWebhookSecret.
// It takes the token (nonce + ciphertext in base64), decrypts it, and returns
// the WABA ID and Organization ID.
func DecryptWebhookSecret(token, encryptionKey string) (*WebhookSecretData, error) {
	// 1. Base64 decode
	raw, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("failed to base64 decode token: %w", err)
	}

	// 2. Decrypt using AES-GCM
	block, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM cipher: %w", err)
	}

	nonceSize := aesGCM.NonceSize()
	if len(raw) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := raw[:nonceSize], raw[nonceSize:]

	// 3. Decrypt
	plaintextJSON, err := aesGCM.Open(nil, nonce, ciphertext, nil)

	fmt.Println("plaintextJSON", string(plaintextJSON))

	if err != nil {
		return nil, fmt.Errorf("failed to decrypt token: %w", err)
	}

	// 4. Unmarshal JSON
	var secretData WebhookSecretData
	if err = json.Unmarshal(plaintextJSON, &secretData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal secret data: %w", err)
	}

	return &secretData, nil
}
