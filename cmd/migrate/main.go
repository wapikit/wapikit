package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	_ "ariga.io/atlas-go-sdk/recordriver"
	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/sarthakjdev/wapikit/database"
)

// func loadEnums(sb *strings.Builder) (string, error) {
// 	enums := []interface{}{
// 		database.ContactStatus(0),
// 		database.OrganizationMemberPermission(0),
// 		database.OrganizationMemberRole(0),
// 	}

// 	for _, enum := range enums {
// 		// Get the type of the enum value
// 		enumType := reflect.TypeOf(enum)

// 		// Check if it's a valid enum type
// 		if enumType.Kind() != reflect.Int {
// 			return "", fmt.Errorf("invalid enum type: %s", enumType.Name())
// 		}

// 		// Get the name of the enum type
// 		enumName := enumType.Name()

// 		fmt.Println("enumName", enumName)

// 		// Start building the SQL statement
// 		sb.WriteString(fmt.Sprintf("CREATE TYPE %s AS ENUM (\n", enumName))

// 		// Get the values of the enum
// 		enumValues := reflect.ValueOf(enum)
// 		for i := 0; i < enumValues.NumField(); i++ {
// 			value := enumValues.Type().Field(i).Name
// 			sb.WriteString(fmt.Sprintf("    '%s'", value))
// 			if i < enumValues.NumField()-1 {
// 				sb.WriteString(",\n")
// 			} else {
// 				sb.WriteString("\n")
// 			}
// 		}

// 		sb.WriteString(");\n\n")
// 	}

// 	return sb.String(), nil
// }

func loadModels(sb *strings.Builder) {
	models := []interface{}{
		&database.Organization{},
		&database.OrganizationMember{},
		&database.WhatsappBusinessAccount{},
		&database.WhatsappBusinessAccountPhoneNumber{},
		&database.Contact{},
		&database.ContactList{},
		&database.Campaign{},
		&database.Conversation{},
		&database.Message{},
		&database.TrackLink{},
		&database.TrackLinkClick{},
		&database.Tag{},
	}

	stmts, err := gormschema.New("postgres").Load(models...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	sb.WriteString(stmts)
	sb.WriteString(";\n")
}

func main() {
	sb := &strings.Builder{}
	// loadEnums(sb)
	loadModels(sb)
	io.WriteString(os.Stdout, sb.String())
}
