package utils

import (
	"fmt"
	"math/rand"
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
	wapi "github.com/wapikit/wapi.go/pkg/client"
	"github.com/wapikit/wapikit/internal/api_types"
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
	rand.Seed(time.Now().UnixNano())
	min := 100000
	max := 999999
	otp := rand.Intn(max-min+1) + min
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

		// Determine if the parameter is required or optional
		required := !strings.Contains(paramName, "omitempty")

		// Bind the query parameter to the field
		err := binder.BindQueryParameter("form", true, required, paramName, context.QueryParams(), structFieldVal.Addr().Interface())
		if err != nil {
			return err
		}
	}

	return nil
}

func GetFeatureFlags(userId, organizationId uuid.UUID) (api_types.FeatureFlags, error) {
	// Fetch the feature flags from the database

	response := api_types.FeatureFlags{
		SystemFeatureFlags: &api_types.SystemFeatureFlags{
			IsApiAccessEnabled:              true,
			IsMultiOrganizationEnabled:      true,
			IsRoleBasedAccessControlEnabled: true,
		},
		IntegrationFeatureFlags: &api_types.IntegrationFeatureFlags{
			IsCustomChatBoxIntegrationEnabled: true,
			IsOpenAiIntegrationEnabled:        true,
			IsSlackIntegrationEnabled:         true,
		},
	}

	return response, nil

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

func FetchBusinessPhoneNumberId(wapiClient *wapi.Client, phoneNumber string) (string, error) {
	// Fetch the business phone number id from the database

	// add caching here
	phoneNumberDetails, err := wapiClient.Business.PhoneNumber.Fetch(phoneNumber)

	fmt.Println("Phone number details: ", phoneNumberDetails)

	if err != nil {
		fmt.Println("Error fetching phone number details: ", err)
		return "", err
	}

	return phoneNumberDetails.Id, nil

}
