//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import "errors"

type ConversationAssignmentStatus string

const (
	ConversationAssignmentStatus_Assigned   ConversationAssignmentStatus = "Assigned"
	ConversationAssignmentStatus_Unassigned ConversationAssignmentStatus = "Unassigned"
)

func (e *ConversationAssignmentStatus) Scan(value interface{}) error {
	var enumValue string
	switch val := value.(type) {
	case string:
		enumValue = val
	case []byte:
		enumValue = string(val)
	default:
		return errors.New("jet: Invalid scan value for AllTypesEnum enum. Enum value has to be of type string or []byte")
	}

	switch enumValue {
	case "Assigned":
		*e = ConversationAssignmentStatus_Assigned
	case "Unassigned":
		*e = ConversationAssignmentStatus_Unassigned
	default:
		return errors.New("jet: Invalid scan value '" + enumValue + "' for ConversationAssignmentStatus enum")
	}

	return nil
}

func (e ConversationAssignmentStatus) String() string {
	return string(e)
}
