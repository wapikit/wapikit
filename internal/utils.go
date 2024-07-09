package internal

import (
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
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
