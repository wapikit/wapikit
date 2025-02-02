//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var WhatsappBusinessAccount = newWhatsappBusinessAccountTable("public", "WhatsappBusinessAccount", "")

type whatsappBusinessAccountTable struct {
	postgres.Table

	// Columns
	UniqueId       postgres.ColumnString
	CreatedAt      postgres.ColumnTimestampz
	UpdatedAt      postgres.ColumnTimestampz
	AccountId      postgres.ColumnString
	AccessToken    postgres.ColumnString
	WebhookSecret  postgres.ColumnString
	OrganizationId postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type WhatsappBusinessAccountTable struct {
	whatsappBusinessAccountTable

	EXCLUDED whatsappBusinessAccountTable
}

// AS creates new WhatsappBusinessAccountTable with assigned alias
func (a WhatsappBusinessAccountTable) AS(alias string) *WhatsappBusinessAccountTable {
	return newWhatsappBusinessAccountTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new WhatsappBusinessAccountTable with assigned schema name
func (a WhatsappBusinessAccountTable) FromSchema(schemaName string) *WhatsappBusinessAccountTable {
	return newWhatsappBusinessAccountTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new WhatsappBusinessAccountTable with assigned table prefix
func (a WhatsappBusinessAccountTable) WithPrefix(prefix string) *WhatsappBusinessAccountTable {
	return newWhatsappBusinessAccountTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new WhatsappBusinessAccountTable with assigned table suffix
func (a WhatsappBusinessAccountTable) WithSuffix(suffix string) *WhatsappBusinessAccountTable {
	return newWhatsappBusinessAccountTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newWhatsappBusinessAccountTable(schemaName, tableName, alias string) *WhatsappBusinessAccountTable {
	return &WhatsappBusinessAccountTable{
		whatsappBusinessAccountTable: newWhatsappBusinessAccountTableImpl(schemaName, tableName, alias),
		EXCLUDED:                     newWhatsappBusinessAccountTableImpl("", "excluded", ""),
	}
}

func newWhatsappBusinessAccountTableImpl(schemaName, tableName, alias string) whatsappBusinessAccountTable {
	var (
		UniqueIdColumn       = postgres.StringColumn("UniqueId")
		CreatedAtColumn      = postgres.TimestampzColumn("CreatedAt")
		UpdatedAtColumn      = postgres.TimestampzColumn("UpdatedAt")
		AccountIdColumn      = postgres.StringColumn("AccountId")
		AccessTokenColumn    = postgres.StringColumn("AccessToken")
		WebhookSecretColumn  = postgres.StringColumn("WebhookSecret")
		OrganizationIdColumn = postgres.StringColumn("OrganizationId")
		allColumns           = postgres.ColumnList{UniqueIdColumn, CreatedAtColumn, UpdatedAtColumn, AccountIdColumn, AccessTokenColumn, WebhookSecretColumn, OrganizationIdColumn}
		mutableColumns       = postgres.ColumnList{CreatedAtColumn, UpdatedAtColumn, AccountIdColumn, AccessTokenColumn, WebhookSecretColumn, OrganizationIdColumn}
	)

	return whatsappBusinessAccountTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		UniqueId:       UniqueIdColumn,
		CreatedAt:      CreatedAtColumn,
		UpdatedAt:      UpdatedAtColumn,
		AccountId:      AccountIdColumn,
		AccessToken:    AccessTokenColumn,
		WebhookSecret:  WebhookSecretColumn,
		OrganizationId: OrganizationIdColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
