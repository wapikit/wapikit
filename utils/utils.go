package utils

import (
	mathRandom "math/rand"
	"net"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	. "github.com/go-jet/jet/v2/postgres"
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

func GenerateOtp(isProduction bool) string {
	if !isProduction {
		return "123456"
	}
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

type WebhookSecretData struct {
	WhatsappBusinessAccountId string `json:"whatsapp_business_account_id"`
	OrganizationId            string `json:"organization_id"`
}

// GetUserIpFromRequest extracts the user's IP address from an HTTP request.
func GetUserIpFromRequest(r *http.Request) string {
	// Check X-Forwarded-For header (common in reverse proxies and load balancers)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		// Take the first IP (original client IP)
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header (another common header for client IP)
	if xRealIP := r.Header.Get("X-Real-IP"); xRealIP != "" {
		return xRealIP
	}

	// Fallback to RemoteAddr (from the TCP connection)
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // Use the raw address if splitting fails
	}
	return ip
}

// GetUserCountryFromRequest determines the user's country based on their IP address.
// This function assumes the existence of a GeoIP database/service (e.g., MaxMind or IP2Location).
func GetUserCountryFromRequest(r *http.Request) string {
	userIP := GetUserIpFromRequest(r)
	// Example: Using a fictional `GetCountryFromIP` function that uses a GeoIP database
	country, err := GetCountryFromIP(userIP)
	if err != nil {
		return "Unknown" // Return "Unknown" if the country cannot be determined
	}
	return country
}

// GetCountryFromIP is a placeholder for a GeoIP lookup function.
// Replace this with an actual GeoIP database query (e.g., MaxMind, IP2Location).
func GetCountryFromIP(ip string) (string, error) {
	// Here you can integrate an external library or API for IP to country mapping
	// Example: Using MaxMind GeoIP2 reader
	// Replace the below logic with actual GeoIP implementation
	if ip == "127.0.0.1" || strings.HasPrefix(ip, "192.168.") {
		return "Local", nil // Local addresses are not mapped to a country
	}

	// Simulated response for demonstration purposes
	return "United States", nil
}

func GetCurrentTimeAndDateInUTCString() string {
	return time.Now().UTC().String()
}
