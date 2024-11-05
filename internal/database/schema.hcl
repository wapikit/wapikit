schema "public" {
}

// ===== ENUMS ====

enum "UserAccountStatusEnum" {
  schema = schema.public
  values = ["Active", "Deleted", "Suspended"]
}

enum "OauthProviderEnum" {
  schema = schema.public
  values = ["Google"]
}

enum "OrganizationInviteStatusEnum" {
  schema = schema.public
  values = ["Pending", "Redeemed"]
}

enum "ContactStatus" {
  schema = schema.public
  values = ["Active", "Inactive", "Blocked", "Deleted"]
}

enum "ConversationStatus" {
  schema = schema.public
  values = ["Active", "Closed", "Deleted"]
}

enum "MessageDirection" {
  schema = schema.public
  values = ["InBound", "OutBound"]
}

enum "MessageStatus" {
  schema = schema.public
  values = ["Sent", "Delivered", "Read", "Failed", "UnDelivered"]
}

enum "ConversationInitiatedEnum" {
  schema = schema.public
  values = ["Cotact", "Campaign"]
}

enum "CampaignStatus" {
  schema = schema.public
  values = ["Draft", "Running", "Finished", "Paused", "Cancelled", "Scheduled"]
}

enum "AccessLogType" {
  schema = schema.public
  values = ["WebInterface", "ApiAccess"]
}

enum "UserPermissionLevel" {
  schema = schema.public
  values = ["Owner", "Member"]
}

enum "OrganizaRolePermissionEnum" {
  schema = schema.public
  values = [
    "GetTeam",
    "UpdateTeam",
    "GetCampaign",
    "CreateCampaign",
    "UpdateCampaign",
    "DeleteCampaign",
    "GetConversations",
    "GetConversation",
    "UpdateConversation",
    "DeleteConversation",
    "AssignConversation",
    "UnassignConversation",
    "GetMessages",
    "GetList",
    "CreateList",
    "UpdateList",
    "DeleteList",
    "GetApiKey",
    "RegenerateApiKey",
    "GetAppSettings",
    "UpdateAppSettings",
    "GetContacts",
    "GetContact",
    "CreateContact",
    "UpdateContact",
    "DeleteContact",
    "BulkImportContacts",
    "GetPrimaryAnalytics",
    "GetSecondaryAnalytics",
    "GetCampaignAnalytics",
    "GetCampaignsAnalytics",
    "GetMetadata",
    "GetOrganizations",
    "CreateOrganization",
    "GetOrganization",
    "UpdateOrganization",
    "TransferOwnership",
    "ManageOrganizationSettings",
    "ManageOrganizationTags",
    "ManageOrganizationInvites",
    "GetMembers",
    "ManageMember",
    "AssignRoleToMember",
    "ManageRoles",
    "ManageRole",
    "ManageIntegrations",
    "SwitchOrganization",
    "JoinOrganization"
  ]
}


// ===== PRIMARY TABLES ====

table "User" {
  schema = schema.public
  column "UniqueId" {
    type    = uuid
    null    = false
    default = sql("gen_random_uuid()")
  }
  column "CreatedAt" {
    type    = timestamptz
    null    = false
    default = sql("now()")
  }
  column "UpdatedAt" {
    type = timestamptz
    null = false
  }

  column "Name" {
    type = text
    null = false
  }
  column "Email" {
    type = text
    null = false
  }
  column "PhoneNumber" {
    type = text
    null = true
  }

  column "Username" {
    type = text
    null = false
  }
  column "Password" {
    type = text
    null = true
  }

  column "OauthProvider" {
    type = enum.OauthProviderEnum
    null = true
  }

  column "ProfilePictureUrl" {
    type = text
    null = true
  }

  column "Status" {
    type = enum.UserAccountStatusEnum
    null = false
  }

  primary_key {
    columns = [column.UniqueId]
  }

  index "UserEmailIndex" {
    columns = [column.Email]
    unique  = true
  }

  index "UserUsernameIndex" {
    columns = [column.Username]
    unique  = true
  }
}

table "Organization" {
  schema = schema.public
  column "UniqueId" {
    type    = uuid
    null    = false
    default = sql("gen_random_uuid()")
  }
  column "CreatedAt" {
    type    = timestamptz
    null    = false
    default = sql("now()")
  }
  column "UpdatedAt" {
    type = timestamptz
    null = false
  }
  column "Name" {
    type = text
    null = false
  }

  column "Description" {
    type = text
    null = true
  }

  column "WebsiteUrl" {
    type = text
    null = true
  }
  column "LogoUrl" {
    type = text
    null = true
  }
  column "FaviconUrl" {
    type = text
  }

  # adding this so that we can notify the organization members
  column "SlackWebhookUrl" {
    type = text
    null = true
  }

  column "DiscordWebhookUrl" {
    type = text
    null = true
  }

  # below details so that self hosted system can be sent email notifications, if they want
  column "SmtpClientHost" {
    type = text
    null = true
  }

  column "SmtpClientUsername" {
    type = text
    null = true
  }

  column "SmtpClientPassword" {
    type = text
    null = true
  }

  primary_key {
    columns = [column.UniqueId]
  }
}

table "OrganizationMember" {
  schema = schema.public
  column "UniqueId" {
    type    = uuid
    null    = false
    default = sql("gen_random_uuid()")
  }
  column "CreatedAt" {
    type    = timestamptz
    null    = false
    default = sql("now()")
  }
  column "UpdatedAt" {
    type = timestamptz
    null = false
  }

  column "AccessLevel" {
    type = enum.UserPermissionLevel
    null = false
  }

  column "OrganizationId" {
    type = uuid
    null = false
  }

  column "UserId" {
    type = uuid
    null = false
  }

  column "InviteId" {
    type = uuid
    null = true
  }

  primary_key {
    columns = [column.UniqueId]
  }

  foreign_key "OrganizationToOrganizationMemberForeignKey" {
    columns     = [column.OrganizationId]
    ref_columns = [table.Organization.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  foreign_key "OrganizationInviteToOrganizationMemberForeignKey" {
    columns     = [column.InviteId]
    ref_columns = [table.OrganizationMemberInvite.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }


  foreign_key "OrganizationMemberToUserForeignKey" {
    columns     = [column.UserId]
    ref_columns = [table.User.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  index "OrganizationMemberOrganizationIdIndex" {
    columns = [column.OrganizationId]
  }

  index "OrganizationMemberUserIdIndex" {
    columns = [column.UserId]
  }
}

table "OrganizationMemberInvite" {
  schema = schema.public

  column "UniqueId" {
    type    = uuid
    null    = false
    default = sql("gen_random_uuid()")
  }

  column "CreatedAt" {
    type    = timestamptz
    null    = false
    default = sql("now()")
  }

  column "UpdatedAt" {
    type = timestamptz
    null = false
  }

  column "Slug" {
    type = text
    null = false
  }

  column "email" {
    type = text
    null = false
  }

  column "AccessLevel" {
    type = enum.UserPermissionLevel
    null = false
  }

  column "OrganizationId" {
    type = uuid
    null = false
  }

  column "Status" {
    type    = enum.OrganizationInviteStatusEnum
    null    = false
    default = "Pending"
  }

  column "InvitedByUserId" {
    type = uuid
    null = false
  }

  primary_key {
    columns = [column.UniqueId]
  }

  foreign_key "UserToOrganizationInviteForeignKey" {
    columns     = [column.InvitedByUserId]
    ref_columns = [table.User.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }


  foreign_key "OrganizationToOrganizationInviteForeignKey" {
    columns     = [column.OrganizationId]
    ref_columns = [table.Organization.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  index "OrganizationInviteOrganizationIdIndex" {
    columns = [column.OrganizationId]
  }

  index "OrganizationInviteInvitedByUserIdIndex" {
    columns = [column.InvitedByUserId]
  }

}

table "OrganizationRole" {
  schema = schema.public

  column "UniqueId" {
    type    = uuid
    null    = false
    default = sql("gen_random_uuid()")
  }

  column "CreatedAt" {
    type    = timestamptz
    null    = false
    default = sql("now()")
  }

  column "UpdatedAt" {
    type = timestamptz
    null = false
  }

  column "Name" {
    type = text
    null = false
  }

  column "Description" {
    type = text
    null = true
  }

  column "Permissions" {
    type    = text
    null    = false
    default = ""
  }

  column "OrganizationId" {
    type = uuid
    null = false
  }

  primary_key {
    columns = [column.UniqueId]
  }

  foreign_key "OrganizationToOrganizationRoleForeignKey" {
    columns     = [column.OrganizationId]
    ref_columns = [table.Organization.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  index "OrganizationRoleOrganizationIdIndex" {
    columns = [column.OrganizationId]
  }
}

table "RoleAssignment" {
  schema = schema.public

  column "UniqueId" {
    type    = uuid
    null    = false
    default = sql("gen_random_uuid()")
  }

  column "CreatedAt" {
    type    = timestamptz
    null    = false
    default = sql("now()")
  }
  column "UpdatedAt" {
    type = timestamptz
    null = false
  }

  column "OrganizationRoleId" {
    type = uuid
    null = false
  }

  column "OrganizationMemberId" {
    type = uuid
    null = false
  }

  primary_key {
    columns = [column.UniqueId]
  }

  foreign_key "OrganizationRoleToRoleAssignmentForeignKey" {
    columns     = [column.OrganizationRoleId]
    ref_columns = [table.OrganizationRole.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  index "RoleAssignmentOrganizationRoleIdIndex" {
    columns = [column.OrganizationRoleId]
  }

  foreign_key "OrganizationMemberToRoleAssignmentForeignKey" {
    columns     = [column.OrganizationMemberId]
    ref_columns = [table.OrganizationMember.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  index "RoleAssignmentOrganizationMemberIdIndex" {
    columns = [column.OrganizationMemberId]
  }

  index "RoleAssignmentUniqueIndex" {
    columns = [column.OrganizationRoleId, column.OrganizationMemberId]
    unique  = true
  }
}

table "ApiKey" {
  schema = schema.public
  column "UniqueId" {
    type    = uuid
    null    = false
    default = sql("gen_random_uuid()")
  }
  column "CreatedAt" {
    type    = timestamptz
    null    = false
    default = sql("now()")
  }
  column "UpdatedAt" {
    type = timestamptz
    null = false
  }

  column "MemberId" {
    type = uuid
    null = false
  }

  column "Key" {
    type = text
    null = false
  }

  column "OrganizationId" {
    type = uuid
    null = false
  }

  primary_key {
    columns = [column.UniqueId]
  }

  foreign_key "ApiKeyToOrganizationForeignKey" {
    columns     = [column.OrganizationId]
    ref_columns = [table.Organization.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  foreign_key "ApiKeyToOrganizationMemberForeignKey" {
    columns     = [column.MemberId]
    ref_columns = [table.OrganizationMember.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  index "ApiKeyOrganizationIdIndex" {
    columns = [column.OrganizationId]
  }

  index "ApiKeyIndex" {
    columns = [column.Key]
    unique  = true
  }

  index "ApiKeyOrganizationMemberIdIndex" {
    columns = [column.MemberId]
    unique  = true
  }
}


table "WhatsappBusinessAccount" {
  schema = schema.public
  column "UniqueId" {
    type    = uuid
    null    = false
    default = sql("gen_random_uuid()")
  }
  column "CreatedAt" {
    type    = timestamptz
    null    = false
    default = sql("now()")
  }
  column "UpdatedAt" {
    type = timestamptz
    null = false
  }

  column "AccountId" {
    type = text
    null = false
  }

  column "AccessToken" {
    type = text
    null = false
  }

  column "WebhookSecret" {
    type = text
    null = false
  }


  column "OrganizationId" {
    type = uuid
    null = false
  }

  primary_key {
    columns = [column.UniqueId]
  }

  foreign_key "WhatsappBusinessAccountToOrganizationForeignKey" {
    columns     = [column.OrganizationId]
    ref_columns = [table.Organization.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  index "WhatsappBusinessAccountOrganizationIdIndex" {
    columns = [column.OrganizationId]
  }

  index "WhatsappBusinessAccountAccountIdIndex" {
    columns = [column.AccountId]
    unique  = true
  }
}

table "Contact" {
  schema = schema.public
  column "UniqueId" {
    type    = uuid
    null    = false
    default = sql("gen_random_uuid()")
  }
  column "CreatedAt" {
    type    = timestamptz
    null    = false
    default = sql("now()")
  }
  column "UpdatedAt" {
    type = timestamptz
    null = false
  }
  column "OrganizationId" {
    type = uuid
    null = false
  }
  column "Status" {
    type = enum.ContactStatus
    null = false
  }
  column "Name" {
    type = text
    null = false
  }
  column "PhoneNumber" {
    type = text
    null = false
  }

  column "Attributes" {
    type = jsonb
    null = true
  }

  primary_key {
    columns = [column.UniqueId]
  }

  foreign_key "OrganizationToContactForeignKey" {
    columns     = [column.OrganizationId]
    ref_columns = [table.Organization.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }
  index "ContactOrganizationIdIndex" {
    columns = [column.OrganizationId]
  }
  index "ContactPhoneNumberIndex" {
    columns = [column.PhoneNumber]
    unique  = true
  }

  index "ContactNumberOrganizationIdUniqueIndex" {
    columns = [column.OrganizationId, column.PhoneNumber]
    unique  = true
  }
}

table "ContactList" {
  schema = schema.public
  column "UniqueId" {
    type    = uuid
    null    = false
    default = sql("gen_random_uuid()")
  }
  column "CreatedAt" {
    type    = timestamptz
    null    = false
    default = sql("now()")
  }
  column "UpdatedAt" {
    type = timestamptz
    null = false
  }
  column "OrganizationId" {
    type = uuid
    null = false
  }
  column "Name" {
    type = text
    null = false
  }
  primary_key {
    columns = [column.UniqueId]
  }
  foreign_key "OrganizationToContactListForeignKey" {
    columns     = [column.OrganizationId]
    ref_columns = [table.Organization.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }
  index "ContactListOrganizationIdIndex" {
    columns = [column.OrganizationId]

  }
}


table "Campaign" {
  schema = schema.public
  column "UniqueId" {
    type    = uuid
    null    = false
    default = sql("gen_random_uuid()")
  }
  column "CreatedAt" {
    type    = timestamptz
    null    = false
    default = sql("now()")
  }
  column "UpdatedAt" {
    type = timestamptz
    null = false
  }

  column "Name" {
    type = text
    null = false
  }

  column "Status" {
    type    = enum.CampaignStatus
    null    = false
    default = "Draft"
  }

  column "IsLinkTrackingEnabled" {
    type    = boolean
    default = false
    null    = false
  }

  column "CreatedByOrganizationMemberId" {
    type = uuid
    null = false
  }

  column "OrganizationId" {
    type = uuid
    null = false
  }

  // this would be the template Id provided by whatsapp business platform only
  column "MessageTemplateId" {
    type = text
    null = true
  }

  column "PhoneNumber" {
    type = text
    null = false
  }

  primary_key {
    columns = [column.UniqueId]
  }

  foreign_key "CampaignToOrganizationMemberForeignKey" {
    columns     = [column.CreatedByOrganizationMemberId]
    ref_columns = [table.OrganizationMember.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  foreign_key "OrganizationToCampaignForeignKey" {
    columns     = [column.OrganizationId]
    ref_columns = [table.Organization.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  index "CampaignCreatedByOrganizationMemberIdIndex" {
    columns = [column.CreatedByOrganizationMemberId]
  }

  index "CampaignMessageTemplateIndex" {
    columns = [column.MessageTemplateId]
  }
}

table "Conversation" {
  schema = schema.public
  column "UniqueId" {
    type    = uuid
    null    = false
    default = sql("gen_random_uuid()")
  }
  column "CreatedAt" {
    type    = timestamptz
    null    = false
    default = sql("now()")
  }
  column "UpdatedAt" {
    type = timestamptz
    null = false
  }

  column "ContactId" {
    type = uuid
    null = false
  }

  column "OrganizationId" {
    type = uuid
    null = false
  }

  column "Status" {
    type = enum.ConversationStatus
    null = false
  }

  column "WhatsappBusinessAccountPhoneNumberId" {
    type = uuid
    null = false
  }

  column "InitiatedBy" {
    type = enum.ConversationInitiatedEnum
    null = false
  }

  primary_key {
    columns = [column.UniqueId]
  }

  foreign_key "ConversationToContactForeignKey" {
    columns     = [column.ContactId]
    ref_columns = [table.Contact.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  foreign_key "ConversationToWhatsappBusinessAccountPhoneNumberForeignKey" {
    columns     = [column.WhatsappBusinessAccountPhoneNumberId]
    ref_columns = [table.WhatsappBusinessAccountPhoneNumber.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  foreign_key "ConversationToOrganizationForeignKey" {
    columns     = [column.OrganizationId]
    ref_columns = [table.Organization.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  index "ConversationContactIdIndex" {
    columns = [column.ContactId]
  }

  index "ConversationWhatsappBusinessAccountPhoneNumberIdIndex" {
    columns = [column.WhatsappBusinessAccountPhoneNumberId]
  }
}

table "Message" {
  schema = schema.public
  column "UniqueId" {
    type    = uuid
    null    = false
    default = sql("gen_random_uuid()")
  }
  column "CreatedAt" {
    type    = timestamptz
    null    = false
    default = sql("now()")
  }
  column "UpdatedAt" {
    type = timestamptz
    null = false
  }

  column "ConversationId" {
    type = uuid
    null = true
  }

  column "CampaignId" {
    type = uuid
    null = true
  }

  column "ContactId" {
    type = uuid
    null = false
  }

  column "WhatsappBusinessAccountPhoneNumberId" {
    type = uuid
    null = false
  }

  column "Direction" {
    type = enum.MessageDirection
    null = false
  }

  column "Content" {
    type = text
    null = true
  }

  column "OrganizationId" {
    type = uuid
    null = false
  }

  column "Status" {
    type = enum.MessageStatus
    null = false
  }

  primary_key {
    columns = [column.UniqueId]
  }

  foreign_key "MessageToCampaignForeignKey" {
    columns     = [column.CampaignId]
    ref_columns = [table.Campaign.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  foreign_key "MessageToWhatsappBusinessAccountPhoneNumberForeignKey" {
    columns     = [column.WhatsappBusinessAccountPhoneNumberId]
    ref_columns = [table.WhatsappBusinessAccountPhoneNumber.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  foreign_key "MessageToContactForeignKey" {
    columns     = [column.ContactId]
    ref_columns = [table.Contact.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  foreign_key "MessageToConversationForeignKey" {
    columns     = [column.ConversationId]
    ref_columns = [table.Conversation.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  foreign_key "MessageToOrganizationForeignKey" {
    columns     = [column.OrganizationId]
    ref_columns = [table.Organization.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  index "MessageCampaignIdIndex" {
    columns = [column.CampaignId]
  }

  index "MessageContactIdIndex" {
    columns = [column.ContactId]
  }

  index "MessageWhatsappBusinessAccountPhoneNumberIdIndex" {
    columns = [column.WhatsappBusinessAccountPhoneNumberId]
  }
}

table "TrackLink" {
  schema = schema.public
  column "UniqueId" {
    type    = uuid
    null    = false
    default = sql("gen_random_uuid()")
  }
  column "CreatedAt" {
    type    = timestamptz
    null    = false
    default = sql("now()")
  }
  column "UpdatedAt" {
    type = timestamptz
    null = false
  }

  column "OrganizationId" {
    type = uuid
    null = false
  }

  column "CampaignId" {
    type = uuid
    null = false
  }

  column "Slug" {
    type = text
    null = false
  }

  column "DestinationUrl" {
    type = text
    null = true
  }

  primary_key {
    columns = [column.UniqueId]
  }

  foreign_key "TrackLinkToCampaignForeignKey" {
    columns     = [column.CampaignId]
    ref_columns = [table.Campaign.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  foreign_key "TrackLinkToOrganizationForeignKey" {
    columns     = [column.OrganizationId]
    ref_columns = [table.Organization.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  index "TrackLinkCampaignIdIndex" {
    columns = [column.CampaignId]

  }
}

table "TrackLinkClick" {
  schema = schema.public
  column "UniqueId" {
    type    = uuid
    null    = false
    default = sql("gen_random_uuid()")
  }
  column "CreatedAt" {
    type    = timestamptz
    null    = false
    default = sql("now()")
  }
  column "UpdatedAt" {
    type = timestamptz
    null = false
  }

  column "TrackLinkId" {
    type = uuid
    null = false
  }

  column "ContactId" {
    type = uuid
    null = false
  }

  primary_key {
    columns = [column.UniqueId]
  }

  foreign_key "TrackLinkClickToTrackLinkForeignKey" {
    columns     = [column.TrackLinkId]
    ref_columns = [table.TrackLink.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  foreign_key "TrackLinkClickToContactForeignKey" {
    columns     = [column.ContactId]
    ref_columns = [table.Contact.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  index "TrackLinkClickTrackLinkIdIndex" {
    columns = [column.TrackLinkId]
  }

  index "TrackLinkClickContactIdIndex" {
    columns = [column.ContactId]
  }
}

table "Tag" {
  schema = schema.public
  column "UniqueId" {
    type    = uuid
    null    = false
    default = sql("gen_random_uuid()")
  }
  column "CreatedAt" {
    type    = timestamptz
    null    = false
    default = sql("now()")
  }
  column "UpdatedAt" {
    type = timestamptz
    null = false
  }

  column "Label" {
    type = text
    null = false
  }

  column "Slug" {
    type = text
    null = false
  }

  column "OrganizationId" {
    type = uuid
    null = false
  }

  primary_key {
    columns = [column.UniqueId]
  }

  unique "UniqueSlug" {
    columns = [column.Slug]
  }

  index "slugIndex" {
    columns = [column.Slug]
  }

  foreign_key "TagToOrganizationForeignKey" {
    columns     = [column.OrganizationId]
    ref_columns = [table.Organization.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  index "TagOrganizationIdIndex" {
    columns = [column.OrganizationId]
  }

  index "TagLabelOragnizationIdUniqueIndex" {
    columns = [column.Label, column.OrganizationId]
    unique  = true
  }
}

table "Integration" {
  schema = schema.public

  column "UniqueId" {
    type    = uuid
    null    = false
    default = sql("gen_random_uuid()")
  }
  column "CreatedAt" {
    type    = timestamptz
    null    = false
    default = sql("now()")
  }
  column "UpdatedAt" {
    type = timestamptz
    null = false
  }

  primary_key {
    columns = [column.UniqueId]
  }
}

// this stores the installed integration for a Organization
table "OrganizationIntegration" {
  schema = schema.public

  column "UniqueId" {
    type    = uuid
    null    = false
    default = sql("gen_random_uuid()")
  }
  column "CreatedAt" {
    type    = timestamptz
    null    = false
    default = sql("now()")
  }
  column "UpdatedAt" {
    type = timestamptz
    null = false
  }

  primary_key {
    columns = [column.UniqueId]
  }

}

table "Notification" {
  schema = schema.public

  column "UniqueId" {
    type    = uuid
    null    = false
    default = sql("gen_random_uuid()")
  }
  column "CreatedAt" {
    type    = timestamptz
    null    = false
    default = sql("now()")
  }
  column "UpdatedAt" {
    type = timestamptz
    null = false
  }

  column "ctaUrl" {
    type = text
    null = true
  }

  column "title" {
    type = text
    null = false
  }

  column "description" {
    type = text
    null = false
  }

  column "type" {
    type = text
    null = true
  }

  column "isBroadcast" {
    type    = boolean
    default = false
    null    = false
  }

  // if the above broadcast is true then the user id can be null, because the notification has been sent to all platform users
  column "UserId" {
    type = uuid
    null = true
  }

  primary_key {
    columns = [column.UniqueId]
  }
}

table "NotificationReadLog" {
  schema = schema.public

  column "UniqueId" {
    type    = uuid
    null    = false
    default = sql("gen_random_uuid()")
  }
  column "CreatedAt" {
    type    = timestamptz
    null    = false
    default = sql("now()")
  }
  column "UpdatedAt" {
    type = timestamptz
    null = false
  }

  column "ReadByUserId" {
    type = uuid
    null = false
  }

  column "NotificationId" {
    type = uuid
    null = false
  }

  primary_key {
    columns = [column.UniqueId]
  }

  foreign_key "NotificationReadLogToNotificationForeignKey" {
    columns     = [column.NotificationId]
    ref_columns = [table.Notification.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  foreign_key "NotificationReadLogToUserForeignKey" {
    columns     = [column.ReadByUserId]
    ref_columns = [table.User.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  index "NotificationReadLogNotificationIdIndex" {
    columns = [column.NotificationId]
  }

  index "NotificationReadLogReadByUserIdIndex" {
    columns = [column.ReadByUserId]
  }
}

// ==== JOIN TABLES ======

table "ContactListContact" {
  schema = schema.public
  column "CreatedAt" {
    type    = timestamptz
    null    = false
    default = sql("now()")
  }
  column "UpdatedAt" {
    type = timestamptz
    null = false
  }

  column "ContactListId" {
    type = uuid
    null = false
  }

  column "ContactId" {
    type = uuid
    null = false
  }

  primary_key {
    columns = [column.ContactListId, column.ContactId]
  }

  foreign_key "ContactListContactToContactListForeignKey" {
    columns     = [column.ContactListId]
    ref_columns = [table.ContactList.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  foreign_key "ContactListContactToContactForeignKey" {
    columns     = [column.ContactId]
    ref_columns = [table.Contact.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }
}

table "ContactListTag" {
  schema = schema.public
  column "CreatedAt" {
    type    = timestamptz
    null    = false
    default = sql("now()")
  }
  column "UpdatedAt" {
    type = timestamptz
    null = false
  }

  column "ContactListId" {
    type = uuid
    null = false
  }

  column "TagId" {
    type = uuid
    null = false
  }

  primary_key {
    columns = [column.ContactListId, column.TagId]
  }

  foreign_key "ContactListTagToContactListForeignKey" {
    columns     = [column.ContactListId]
    ref_columns = [table.ContactList.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  foreign_key "ContactListTagToTagForeignKey" {
    columns     = [column.TagId]
    ref_columns = [table.Tag.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }
}

table "CampaignList" {

  schema = schema.public

  column "CreatedAt" {
    type    = timestamptz
    null    = false
    default = sql("now()")
  }
  column "UpdatedAt" {
    type = timestamptz
    null = false
  }

  column "ContactListId" {
    type = uuid
    null = false
  }

  column "CampaignId" {
    type = uuid
    null = false
  }

  primary_key {
    columns = [column.ContactListId, column.CampaignId]
  }

  foreign_key "CampaignListToContactListForeignKey" {
    columns     = [column.ContactListId]
    ref_columns = [table.ContactList.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  foreign_key "CampaignListToCampaignForeignKey" {
    columns     = [column.CampaignId]
    ref_columns = [table.Campaign.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

}

table "ConversationTag" {
  schema = schema.public
  column "CreatedAt" {
    type    = timestamptz
    null    = false
    default = sql("now()")
  }
  column "UpdatedAt" {
    type = timestamptz
    null = false
  }

  column "ConversationId" {
    type = uuid
    null = false
  }

  column "TagId" {
    type = uuid
    null = false
  }

  primary_key {
    columns = [column.ConversationId, column.TagId]
  }

  foreign_key "ConversationTagToConversationForeignKey" {
    columns     = [column.ConversationId]
    ref_columns = [table.Conversation.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  foreign_key "ConversationTagToTagForeignKey" {
    columns     = [column.TagId]
    ref_columns = [table.Tag.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }
}

table "CampaignTag" {
  schema = schema.public
  column "CreatedAt" {
    type    = timestamptz
    null    = false
    default = sql("now()")
  }
  column "UpdatedAt" {
    type = timestamptz
    null = false
  }

  column "CampaignId" {
    type = uuid
    null = false
  }

  column "TagId" {
    type = uuid
    null = false
  }

  primary_key {
    columns = [column.CampaignId, column.TagId]
  }

  foreign_key "CampaignTagToCampaignForeignKey" {
    columns     = [column.CampaignId]
    ref_columns = [table.Campaign.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  foreign_key "CampaignTagToTagForeignKey" {
    columns     = [column.TagId]
    ref_columns = [table.Tag.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  index "CampaignTagIdCampaignIdUniqueIndex" {
    columns = [column.CampaignId, column.TagId]
    unique  = true
  }
}

table "MessageReply" {
  schema = schema.public
  column "CreatedAt" {
    type    = timestamptz
    null    = false
    default = sql("now()")
  }
  column "UpdatedAt" {
    type = timestamptz
    null = false
  }

  column "MessageId" {
    type = uuid
    null = false
  }

  column "ReplyMessageId" {
    type = uuid
    null = false
  }

  primary_key {
    columns = [column.MessageId, column.ReplyMessageId]
  }

  foreign_key "MessageReplyToMessageForeignKey" {
    columns     = [column.MessageId]
    ref_columns = [table.Message.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  foreign_key "MessageReplyToReplyMessageForeignKey" {
    columns     = [column.ReplyMessageId]
    ref_columns = [table.Message.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }
}
