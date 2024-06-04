package database

import "fmt"

type OrganizationMemberRole string

const (
	OrganizationMemberRoleSuperAdmin OrganizationMemberRole = "SuperAdmin"
	OrganizationMemberRoleAdmin      OrganizationMemberRole = "Admin"
	OrganizationMemberRoleMember     OrganizationMemberRole = "Member"
)

func (e *OrganizationMemberRole) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid str")
	}
	*e = OrganizationMemberRole(str)
	return nil
}

func (e OrganizationMemberRole) Value() (interface{}, error) {
	return string(e), nil
}

type ContactStatus string

const (
	ContactStatusActive   ContactStatus = "Active"
	ContactStatusInActive ContactStatus = "Blocked"
)

func (e *ContactStatus) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid str")
	}
	*e = ContactStatus(str)
	return nil
}

func (e ContactStatus) Value() (interface{}, error) {
	return string(e), nil
}

type OrganizationMemberPermission string

const (
	OrganizationMemberPermissionReadMember         OrganizationMemberPermission = "ReadMember"
	OrganizationMemberPermissionWriteMember        OrganizationMemberPermission = "WriteMember"
	OrganizationMemberPermissionReadCampaign       OrganizationMemberPermission = "ReadCampaign"
	OrganizationMemberPermissionWriteCampaign      OrganizationMemberPermission = "WriteCampaign"
	OrganizationMemberPermissionReadContact        OrganizationMemberPermission = "ReadContact"
	OrganizationMemberPermissionWriteContact       OrganizationMemberPermission = "WriteContact"
	OrganizationMemberPermissionReadContactList    OrganizationMemberPermission = "ReadContactList"
	OrganizationMemberPermissionWriteContactList   OrganizationMemberPermission = "WriteContactList"
	OrganizationMemberPermissionReadSettings       OrganizationMemberPermission = "ReadSettings"
	OrganizationMemberPermissionWriteSettings      OrganizationMemberPermission = "WriteSettings"
	OrganizationMemberPermissionReadConversations  OrganizationMemberPermission = "ReadConversations"
	OrganizationMemberPermissionWriteConversations OrganizationMemberPermission = "WriteConversations"
)

func (e *OrganizationMemberPermission) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid str")
	}
	*e = OrganizationMemberPermission(str)
	return nil
}

func (e OrganizationMemberPermission) Value() (interface{}, error) {
	return string(e), nil
}

type MessageStatus string

const (
	MessageStatusSent        MessageStatus = "Sent"
	MessageStatusDelivered   MessageStatus = "Delivered"
	MessageStatusRead        MessageStatus = "Read"
	MessageStatusFailed      MessageStatus = "Failed"
	MessageStatusUnDelivered MessageStatus = "UnDelivered"
)

func (e *MessageStatus) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid str")
	}
	*e = MessageStatus(str)
	return nil
}

func (e MessageStatus) Value() (interface{}, error) {
	return string(e), nil
}

type MessageDirection string

const (
	MessageDirectionIncoming MessageDirection = "Incoming"
	MessageDirectionOutgoing MessageDirection = "Outgoing"
)

func (e *MessageDirection) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid str")
	}
	*e = MessageDirection(str)
	return nil
}

func (e MessageDirection) Value() (interface{}, error) {
	return string(e), nil
}
