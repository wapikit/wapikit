//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"github.com/google/uuid"
	"time"
)

type ConversationAssignment struct {
	CreatedAt                      time.Time
	UpdatedAt                      time.Time
	ConversationId                 uuid.UUID `sql:"primary_key"`
	AssignedToOrganizationMemberId uuid.UUID `sql:"primary_key"`
	Status                         ConversationAssignmentStatus
}