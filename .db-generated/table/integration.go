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

var Integration = newIntegrationTable("public", "Integration", "")

type integrationTable struct {
	postgres.Table

	// Columns
	UniqueId  postgres.ColumnString
	CreatedAt postgres.ColumnTimestampz
	UpdatedAt postgres.ColumnTimestampz

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type IntegrationTable struct {
	integrationTable

	EXCLUDED integrationTable
}

// AS creates new IntegrationTable with assigned alias
func (a IntegrationTable) AS(alias string) *IntegrationTable {
	return newIntegrationTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new IntegrationTable with assigned schema name
func (a IntegrationTable) FromSchema(schemaName string) *IntegrationTable {
	return newIntegrationTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new IntegrationTable with assigned table prefix
func (a IntegrationTable) WithPrefix(prefix string) *IntegrationTable {
	return newIntegrationTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new IntegrationTable with assigned table suffix
func (a IntegrationTable) WithSuffix(suffix string) *IntegrationTable {
	return newIntegrationTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newIntegrationTable(schemaName, tableName, alias string) *IntegrationTable {
	return &IntegrationTable{
		integrationTable: newIntegrationTableImpl(schemaName, tableName, alias),
		EXCLUDED:         newIntegrationTableImpl("", "excluded", ""),
	}
}

func newIntegrationTableImpl(schemaName, tableName, alias string) integrationTable {
	var (
		UniqueIdColumn  = postgres.StringColumn("UniqueId")
		CreatedAtColumn = postgres.TimestampzColumn("CreatedAt")
		UpdatedAtColumn = postgres.TimestampzColumn("UpdatedAt")
		allColumns      = postgres.ColumnList{UniqueIdColumn, CreatedAtColumn, UpdatedAtColumn}
		mutableColumns  = postgres.ColumnList{CreatedAtColumn, UpdatedAtColumn}
	)

	return integrationTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		UniqueId:  UniqueIdColumn,
		CreatedAt: CreatedAtColumn,
		UpdatedAt: UpdatedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
