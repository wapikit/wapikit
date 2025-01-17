package utils

import (
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
