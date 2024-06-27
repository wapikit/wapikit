package internal

import (
	"reflect"
	"strings"

	"github.com/labstack/echo/v4"
	binder "github.com/oapi-codegen/runtime"
	"github.com/oklog/ulid"
)

// // Create a new ULID (using current time and default entropy source)
// newUlid := ulid.Make()
// fmt.Println("New ULID:", newUlid)

// // Create a ULID from a specific time
// timestamp := time.Now()
// ulidFromTime := ulid.MustNew(ulid.Timestamp(timestamp), ulid.DefaultEntropy())
// fmt.Println("ULID from time:", ulidFromTime)

// // Parse a ULID string
// ulidString := "01ARZ3NDEKTSV4RRFFQ69G5FAV"
// parsedUlid, err := ulid.Parse(ulidString)
// if err != nil {
// 	panic(err)
// }
// fmt.Println("Parsed ULID:", parsedUlid)

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
