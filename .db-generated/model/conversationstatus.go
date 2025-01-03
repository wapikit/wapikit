//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import "errors"

type ConversationStatus string

const (
	ConversationStatus_Active  ConversationStatus = "Active"
	ConversationStatus_Closed  ConversationStatus = "Closed"
	ConversationStatus_Deleted ConversationStatus = "Deleted"
)

func (e *ConversationStatus) Scan(value interface{}) error {
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
	case "Active":
		*e = ConversationStatus_Active
	case "Closed":
		*e = ConversationStatus_Closed
	case "Deleted":
		*e = ConversationStatus_Deleted
	default:
		return errors.New("jet: Invalid scan value '" + enumValue + "' for ConversationStatus enum")
	}

	return nil
}

func (e ConversationStatus) String() string {
	return string(e)
}
