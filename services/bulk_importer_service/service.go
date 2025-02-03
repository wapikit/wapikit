package bulk_importer_service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/wapikit/wapikit/.db-generated/model"
	"github.com/wapikit/wapikit/.db-generated/table"
	"github.com/wapikit/wapikit/interfaces"
)

var (
	ImporterServiceInstance *ImporterService
)

type ImporterService struct {
	Logger *slog.Logger
	Db     *sql.DB
}

func NewImporterService(logger *slog.Logger, db *sql.DB) *ImporterService {
	if ImporterServiceInstance == nil {
		ImporterServiceInstance = &ImporterService{Logger: logger}
	}
	return ImporterServiceInstance
}

type ImportEvent struct {
	Type    string `json:"type"` // "progress", "error", "complete"
	Message string `json:"message"`
	Current int    `json:"current"`
	Total   int    `json:"total"`
}

func (importer *ImporterService) SendEvent(enc *json.Encoder, context *interfaces.ContextWithSession, event ImportEvent) error {
	if err := enc.Encode(event); err != nil {
		return err
	}
	context.Response().Flush()
	return nil
}

func (importer *ImporterService) ProcessRecord(record []string, orgUuid uuid.UUID) (model.Contact, error) {
	// Validate record length
	if len(record) < 3 {
		return model.Contact{}, fmt.Errorf("invalid record length, expected at least 3 columns")
	}

	// Extract fields
	name := strings.TrimSpace(record[0])
	phone := strings.TrimSpace(record[1])
	attributes := strings.TrimSpace(record[2])

	// Validate required fields
	if phone == "" {
		return model.Contact{}, fmt.Errorf("phone number is required")
	}

	// Validate name length
	if utf8.RuneCountInString(name) > 255 {
		return model.Contact{}, fmt.Errorf("name exceeds maximum length of 255 characters")
	}

	// Parse attributes JSON
	var attrMap map[string]interface{}
	if attributes != "" {
		if err := json.Unmarshal([]byte(attributes), &attrMap); err != nil {
			return model.Contact{}, fmt.Errorf("invalid attributes JSON: %v", err)
		}
	}

	// Prepare contact model
	jsonAttributes, _ := json.Marshal(attrMap)
	stringAttributes := string(jsonAttributes)

	return model.Contact{
		OrganizationId: orgUuid,
		Name:           name,
		PhoneNumber:    phone,
		Attributes:     &stringAttributes,
		Status:         model.ContactStatusEnum_Active,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

// Batch insert function
func (importer *ImporterService) InsertBatch(ctx context.Context, batch []model.Contact, listUuids []uuid.UUID, db *sql.DB) error {
	// Start transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	// Insert contacts
	insertQuery := table.Contact.
		INSERT(table.Contact.MutableColumns).
		MODELS(batch).
		ON_CONFLICT(table.Contact.PhoneNumber, table.Contact.OrganizationId).
		DO_NOTHING()

	var insertedContacts []model.Contact
	if err := insertQuery.QueryContext(ctx, db, &insertedContacts); err != nil {
		return fmt.Errorf("batch insert failed: %v", err)
	}

	// Insert into contact lists
	if len(listUuids) > 0 && len(insertedContacts) > 0 {
		var listContacts []model.ContactListContact
		now := time.Now()

		for _, listId := range listUuids {
			for _, contact := range insertedContacts {
				listContacts = append(listContacts, model.ContactListContact{
					ContactId:     contact.UniqueId,
					ContactListId: listId,
					CreatedAt:     now,
					UpdatedAt:     now,
				})
			}
		}

		listQuery := table.ContactListContact.
			INSERT(table.ContactListContact.AllColumns).
			MODELS(listContacts).
			ON_CONFLICT(table.ContactListContact.ContactId, table.ContactListContact.ContactListId).
			DO_NOTHING()

		if _, err := listQuery.ExecContext(ctx, db); err != nil {
			return fmt.Errorf("failed to insert list associations: %v", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit failed: %v", err)
	}

	return nil
}
