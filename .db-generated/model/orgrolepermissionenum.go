//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import "errors"

type OrgRolePermissionEnum string

const (
	OrgRolePermissionEnum_GetColonOrganizationmember     OrgRolePermissionEnum = "Get:OrganizationMember"
	OrgRolePermissionEnum_CreateColonOrganizationmember  OrgRolePermissionEnum = "Create:OrganizationMember"
	OrgRolePermissionEnum_UpdateColonOrganizationmember  OrgRolePermissionEnum = "Update:OrganizationMember"
	OrgRolePermissionEnum_DeleteColonOrganizationmember  OrgRolePermissionEnum = "Delete:OrganizationMember"
	OrgRolePermissionEnum_GetColonCampaign               OrgRolePermissionEnum = "Get:Campaign"
	OrgRolePermissionEnum_CreateColonCampaign            OrgRolePermissionEnum = "Create:Campaign"
	OrgRolePermissionEnum_UpdateColonCampaign            OrgRolePermissionEnum = "Update:Campaign"
	OrgRolePermissionEnum_DeleteColonCampaign            OrgRolePermissionEnum = "Delete:Campaign"
	OrgRolePermissionEnum_GetColonConversation           OrgRolePermissionEnum = "Get:Conversation"
	OrgRolePermissionEnum_UpdateColonConversation        OrgRolePermissionEnum = "Update:Conversation"
	OrgRolePermissionEnum_DeleteColonConversation        OrgRolePermissionEnum = "Delete:Conversation"
	OrgRolePermissionEnum_AssignColonConversation        OrgRolePermissionEnum = "Assign:Conversation"
	OrgRolePermissionEnum_UnassignColonConversation      OrgRolePermissionEnum = "Unassign:Conversation"
	OrgRolePermissionEnum_GetColonList                   OrgRolePermissionEnum = "Get:List"
	OrgRolePermissionEnum_CreateColonList                OrgRolePermissionEnum = "Create:List"
	OrgRolePermissionEnum_UpdateColonList                OrgRolePermissionEnum = "Update:List"
	OrgRolePermissionEnum_DeleteColonList                OrgRolePermissionEnum = "Delete:List"
	OrgRolePermissionEnum_GetColonTag                    OrgRolePermissionEnum = "Get:Tag"
	OrgRolePermissionEnum_CreateColonTag                 OrgRolePermissionEnum = "Create:Tag"
	OrgRolePermissionEnum_UpdateColonTag                 OrgRolePermissionEnum = "Update:Tag"
	OrgRolePermissionEnum_DeleteColonTag                 OrgRolePermissionEnum = "Delete:Tag"
	OrgRolePermissionEnum_GetColonApikey                 OrgRolePermissionEnum = "Get:ApiKey"
	OrgRolePermissionEnum_RegenerateColonApikey          OrgRolePermissionEnum = "Regenerate:ApiKey"
	OrgRolePermissionEnum_GetColonAppsettings            OrgRolePermissionEnum = "Get:AppSettings"
	OrgRolePermissionEnum_UpdateColonAppsettings         OrgRolePermissionEnum = "Update:AppSettings"
	OrgRolePermissionEnum_GetColonContact                OrgRolePermissionEnum = "Get:Contact"
	OrgRolePermissionEnum_CreateColonContact             OrgRolePermissionEnum = "Create:Contact"
	OrgRolePermissionEnum_UpdateColonContact             OrgRolePermissionEnum = "Update:Contact"
	OrgRolePermissionEnum_DeleteColonContact             OrgRolePermissionEnum = "Delete:Contact"
	OrgRolePermissionEnum_BulkimportColonContacts        OrgRolePermissionEnum = "BulkImport:Contacts"
	OrgRolePermissionEnum_GetColonPrimaryanalytics       OrgRolePermissionEnum = "Get:PrimaryAnalytics"
	OrgRolePermissionEnum_GetColonSecondaryanalytics     OrgRolePermissionEnum = "Get:SecondaryAnalytics"
	OrgRolePermissionEnum_GetColonCampaignanalytics      OrgRolePermissionEnum = "Get:CampaignAnalytics"
	OrgRolePermissionEnum_UpdateColonOrganization        OrgRolePermissionEnum = "Update:Organization"
	OrgRolePermissionEnum_GetColonOrganizationrole       OrgRolePermissionEnum = "Get:OrganizationRole"
	OrgRolePermissionEnum_CreateColonOrganizationrole    OrgRolePermissionEnum = "Create:OrganizationRole"
	OrgRolePermissionEnum_UpdateColonOrganizationrole    OrgRolePermissionEnum = "Update:OrganizationRole"
	OrgRolePermissionEnum_DeleteColonOrganizationrole    OrgRolePermissionEnum = "Delete:OrganizationRole"
	OrgRolePermissionEnum_UpdateColonIntegrationsettings OrgRolePermissionEnum = "Update:IntegrationSettings"
	OrgRolePermissionEnum_GetColonMessagetemplates       OrgRolePermissionEnum = "Get:MessageTemplates"
	OrgRolePermissionEnum_GetColonPhonenumbers           OrgRolePermissionEnum = "Get:PhoneNumbers"
)

func (e *OrgRolePermissionEnum) Scan(value interface{}) error {
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
	case "Get:OrganizationMember":
		*e = OrgRolePermissionEnum_GetColonOrganizationmember
	case "Create:OrganizationMember":
		*e = OrgRolePermissionEnum_CreateColonOrganizationmember
	case "Update:OrganizationMember":
		*e = OrgRolePermissionEnum_UpdateColonOrganizationmember
	case "Delete:OrganizationMember":
		*e = OrgRolePermissionEnum_DeleteColonOrganizationmember
	case "Get:Campaign":
		*e = OrgRolePermissionEnum_GetColonCampaign
	case "Create:Campaign":
		*e = OrgRolePermissionEnum_CreateColonCampaign
	case "Update:Campaign":
		*e = OrgRolePermissionEnum_UpdateColonCampaign
	case "Delete:Campaign":
		*e = OrgRolePermissionEnum_DeleteColonCampaign
	case "Get:Conversation":
		*e = OrgRolePermissionEnum_GetColonConversation
	case "Update:Conversation":
		*e = OrgRolePermissionEnum_UpdateColonConversation
	case "Delete:Conversation":
		*e = OrgRolePermissionEnum_DeleteColonConversation
	case "Assign:Conversation":
		*e = OrgRolePermissionEnum_AssignColonConversation
	case "Unassign:Conversation":
		*e = OrgRolePermissionEnum_UnassignColonConversation
	case "Get:List":
		*e = OrgRolePermissionEnum_GetColonList
	case "Create:List":
		*e = OrgRolePermissionEnum_CreateColonList
	case "Update:List":
		*e = OrgRolePermissionEnum_UpdateColonList
	case "Delete:List":
		*e = OrgRolePermissionEnum_DeleteColonList
	case "Get:Tag":
		*e = OrgRolePermissionEnum_GetColonTag
	case "Create:Tag":
		*e = OrgRolePermissionEnum_CreateColonTag
	case "Update:Tag":
		*e = OrgRolePermissionEnum_UpdateColonTag
	case "Delete:Tag":
		*e = OrgRolePermissionEnum_DeleteColonTag
	case "Get:ApiKey":
		*e = OrgRolePermissionEnum_GetColonApikey
	case "Regenerate:ApiKey":
		*e = OrgRolePermissionEnum_RegenerateColonApikey
	case "Get:AppSettings":
		*e = OrgRolePermissionEnum_GetColonAppsettings
	case "Update:AppSettings":
		*e = OrgRolePermissionEnum_UpdateColonAppsettings
	case "Get:Contact":
		*e = OrgRolePermissionEnum_GetColonContact
	case "Create:Contact":
		*e = OrgRolePermissionEnum_CreateColonContact
	case "Update:Contact":
		*e = OrgRolePermissionEnum_UpdateColonContact
	case "Delete:Contact":
		*e = OrgRolePermissionEnum_DeleteColonContact
	case "BulkImport:Contacts":
		*e = OrgRolePermissionEnum_BulkimportColonContacts
	case "Get:PrimaryAnalytics":
		*e = OrgRolePermissionEnum_GetColonPrimaryanalytics
	case "Get:SecondaryAnalytics":
		*e = OrgRolePermissionEnum_GetColonSecondaryanalytics
	case "Get:CampaignAnalytics":
		*e = OrgRolePermissionEnum_GetColonCampaignanalytics
	case "Update:Organization":
		*e = OrgRolePermissionEnum_UpdateColonOrganization
	case "Get:OrganizationRole":
		*e = OrgRolePermissionEnum_GetColonOrganizationrole
	case "Create:OrganizationRole":
		*e = OrgRolePermissionEnum_CreateColonOrganizationrole
	case "Update:OrganizationRole":
		*e = OrgRolePermissionEnum_UpdateColonOrganizationrole
	case "Delete:OrganizationRole":
		*e = OrgRolePermissionEnum_DeleteColonOrganizationrole
	case "Update:IntegrationSettings":
		*e = OrgRolePermissionEnum_UpdateColonIntegrationsettings
	case "Get:MessageTemplates":
		*e = OrgRolePermissionEnum_GetColonMessagetemplates
	case "Get:PhoneNumbers":
		*e = OrgRolePermissionEnum_GetColonPhonenumbers
	default:
		return errors.New("jet: Invalid scan value '" + enumValue + "' for OrgRolePermissionEnum enum")
	}

	return nil
}

func (e OrgRolePermissionEnum) String() string {
	return string(e)
}
