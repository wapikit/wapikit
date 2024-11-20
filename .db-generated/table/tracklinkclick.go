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

var TrackLinkClick = newTrackLinkClickTable("public", "TrackLinkClick", "")

type trackLinkClickTable struct {
	postgres.Table

	// Columns
	UniqueId    postgres.ColumnString
	CreatedAt   postgres.ColumnTimestampz
	UpdatedAt   postgres.ColumnTimestampz
	TrackLinkId postgres.ColumnString
	ContactId   postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type TrackLinkClickTable struct {
	trackLinkClickTable

	EXCLUDED trackLinkClickTable
}

// AS creates new TrackLinkClickTable with assigned alias
func (a TrackLinkClickTable) AS(alias string) *TrackLinkClickTable {
	return newTrackLinkClickTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new TrackLinkClickTable with assigned schema name
func (a TrackLinkClickTable) FromSchema(schemaName string) *TrackLinkClickTable {
	return newTrackLinkClickTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new TrackLinkClickTable with assigned table prefix
func (a TrackLinkClickTable) WithPrefix(prefix string) *TrackLinkClickTable {
	return newTrackLinkClickTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new TrackLinkClickTable with assigned table suffix
func (a TrackLinkClickTable) WithSuffix(suffix string) *TrackLinkClickTable {
	return newTrackLinkClickTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newTrackLinkClickTable(schemaName, tableName, alias string) *TrackLinkClickTable {
	return &TrackLinkClickTable{
		trackLinkClickTable: newTrackLinkClickTableImpl(schemaName, tableName, alias),
		EXCLUDED:            newTrackLinkClickTableImpl("", "excluded", ""),
	}
}

func newTrackLinkClickTableImpl(schemaName, tableName, alias string) trackLinkClickTable {
	var (
		UniqueIdColumn    = postgres.StringColumn("UniqueId")
		CreatedAtColumn   = postgres.TimestampzColumn("CreatedAt")
		UpdatedAtColumn   = postgres.TimestampzColumn("UpdatedAt")
		TrackLinkIdColumn = postgres.StringColumn("TrackLinkId")
		ContactIdColumn   = postgres.StringColumn("ContactId")
		allColumns        = postgres.ColumnList{UniqueIdColumn, CreatedAtColumn, UpdatedAtColumn, TrackLinkIdColumn, ContactIdColumn}
		mutableColumns    = postgres.ColumnList{CreatedAtColumn, UpdatedAtColumn, TrackLinkIdColumn, ContactIdColumn}
	)

	return trackLinkClickTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		UniqueId:    UniqueIdColumn,
		CreatedAt:   CreatedAtColumn,
		UpdatedAt:   UpdatedAtColumn,
		TrackLinkId: TrackLinkIdColumn,
		ContactId:   ContactIdColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}