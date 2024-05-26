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

func loadEnums(sb *strings.Builder) {
	enums := []string{
		`CREATE TYPE gender AS ENUM (
			'MALE',
			'FEMALE'
		);`,
	}
	for _, enum := range enums {
		sb.WriteString(enum)
		sb.WriteString(";\n")
	}
}

func loadModels(sb *strings.Builder) {
	models := []interface{}{
		&database.User{},
		&database.Subscriber{},
		&database.Admin{},
		&database.SubscriberList{},
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
	loadEnums(sb)
	loadModels(sb)

	io.WriteString(os.Stdout, sb.String())
}
