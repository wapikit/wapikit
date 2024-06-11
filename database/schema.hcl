schema "public" {
}


// ===== ENUMS ====


enum "ContactStatus" {
  schema = schema.public
  values = ["active", "inactive"]
}

enum "MessageDirection" {
  schema = schema.public
  values = ["inbound", "outbound"]
}

enum "MessageStatus" {
  schema = schema.public
  values = ["sent", "delivered", "read", "failed", "undelivered"]
}

enum "ConversationInitiatedEnum" {
  schema = schema.public
  values = ["contact", "campaign"]
}

enum "CampaignStatus" {
  schema = schema.public
  values = ["draft", "running", "finished", "paused", "cancelled"]
}

// ===== PRIMARY TABLES ====



table "Organisation" {
  schema = schema.public
  column "UniqueId" {
    type = uuid
  }
  column "CreatedAt" {
    type = timestamp
  }
  column "UpdatedAt" {
    type = timestamp
  }
  column "Name" {
    type = text
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
  primary_key {
    columns = [column.UniqueId]
  }

}

enum "OrganisationMemberRole" {
  schema = schema.public
  values = ["super_admin", "admin", "member"]
}

table "OrganisationMember" {
  schema = schema.public
  column "UniqueId" {
    type = uuid
  }
  column "CreatedAt" {
    type = timestamp
  }
  column "UpdatedAt" {
    type = timestamp
  }

  column "Role" {
    type = enum.OrganisationMemberRole
  }

  column "Name" {
    type = text
  }
  column "Email" {
    type = text
  }
  column "PhoneNumber" {
    type = text
  }
  column "Username" {
    type = text
  }
  column "Password" {
    type = text
  }

  column "OrganisationId" {
    type = uuid
  }

  primary_key {
    columns = [column.UniqueId]
  }

  foreign_key "OrganisationToOrganisationMemberForeignKey" {
    columns     = [column.OrganisationId]
    ref_columns = [table.Organisation.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  index "OrganisationMemberOrganisationIdIndex" {
    columns = [column.OrganisationId]
  }

  index "OrganisationMemberEmailIndex" {
    columns = [column.Email]
    unique  = true
  }

  index "OrganisationMemberUsernameIndex" {
    columns = [column.Username]
    unique  = true
  }
}

table "WhatsappBusinessAccount" {
  schema = schema.public
  column "UniqueId" {
    type = uuid
  }
  column "CreatedAt" {
    type = timestamp
  }
  column "UpdatedAt" {
    type = timestamp
  }

  column "AccountId" {
    type = text
  }

  column "OrganisationId" {
    type = uuid
  }

  primary_key {
    columns = [column.UniqueId]
  }

  foreign_key "WhatsappBusinessAccountToOrganisationForeignKey" {
    columns     = [column.OrganisationId]
    ref_columns = [table.Organisation.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  index "WhatsappBusinessAccountOrganisationIdIndex" {
    columns = [column.OrganisationId]
  }

  index "WhatsappBusinessAccountAccountIdIndex" {
    columns = [column.AccountId]
    unique  = true
  }
}


table "WhatsappBusinessAccountPhoneNumber" {
  schema = schema.public
  column "UniqueId" {
    type = uuid
  }
  column "CreatedAt" {
    type = timestamp
  }
  column "UpdatedAt" {
    type = timestamp
  }

  column "WhatsappBusinessAccountId" {
    type = uuid
  }

  column "MetaTitle" {
    type = text
  }

  column "MetaDescription" {
    type = text
  }

  column "PhoneNumber" {
    type = text
  }

  primary_key {
    columns = [column.UniqueId]
  }

  foreign_key "PhoneNumberToWhatsappBusinessAccountForeignKey" {
    columns     = [column.WhatsappBusinessAccountId]
    ref_columns = [table.WhatsappBusinessAccount.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  index "PhoneNumberWhatsappBusinessAccountIdIndex" {
    columns = [column.WhatsappBusinessAccountId]
  }

  index "PhoneNumberPhoneNumberIndex" {
    columns = [column.PhoneNumber]
    unique  = true
  }

}


table "Contact" {
  schema = schema.public
  column "UniqueId" {
    type = uuid
  }
  column "CreatedAt" {
    type = timestamp
  }
  column "UpdatedAt" {
    type = timestamp
  }
  column "OrganisationId" {
    type = uuid
  }
  column "Status" {
    type = enum.ContactStatus
  }
  column "Name" {
    type = text
  }
  column "PhoneNumber" {
    type = text
  }

  column "Attributes" {
    type = jsonb
    null = true
  }

  primary_key {
    columns = [column.UniqueId]
  }

  foreign_key "OrganisationToContactForeignKey" {
    columns     = [column.OrganisationId]
    ref_columns = [table.Organisation.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }
  index "ContactOrganisationIdIndex" {
    columns = [column.OrganisationId]
  }
  index "ContactPhoneNumberIndex" {
    columns = [column.PhoneNumber]
    unique  = true
  }
}

table "ContactList" {
  schema = schema.public
  column "UniqueId" {
    type = uuid
  }
  column "CreatedAt" {
    type = timestamp
  }
  column "UpdatedAt" {
    type = timestamp
  }
  column "OrganisationId" {
    type = uuid
  }
  column "Name" {
    type = text
  }
  primary_key {
    columns = [column.UniqueId]
  }
  foreign_key "OrganisationToContactListForeignKey" {
    columns     = [column.OrganisationId]
    ref_columns = [table.Organisation.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }
  index "ContactListOrganisationIdIndex" {
    columns = [column.OrganisationId]

  }
}


table "Campaign" {
  schema = schema.public
  column "UniqueId" {
    type = uuid
  }
  column "CreatedAt" {
    type = timestamp
  }
  column "UpdatedAt" {
    type = timestamp
  }

  column "Name" {
    type = text
  }

  column "Status" {
    type    = enum.CampaignStatus
    default = "draft"
  }

  column "CreatedByOrganisationMemberId" {
    type = uuid
  }

  column "OrganisationId" {
    type = uuid
  }

  // this would be the template Id provided by whatsapp business platform only
  column "MessageTemplateId" {
    type = text
  }

  primary_key {
    columns = [column.UniqueId]
  }

  foreign_key "CampaignToOrganisationMemberForeignKey" {
    columns     = [column.CreatedByOrganisationMemberId]
    ref_columns = [table.OrganisationMember.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  foreign_key "OrganisationToCampaignForeignKey" {
    columns     = [column.OrganisationId]
    ref_columns = [table.Organisation.column.UniqueId]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  index "CampaignCreatedByOrganisationMemberIdIndex" {
    columns = [column.CreatedByOrganisationMemberId]
  }

  index "CampaignMessageTemplateIndex" {
    columns = [column.MessageTemplateId]
  }
}

table "Conversation" {
  schema = schema.public
  column "UniqueId" {
    type = uuid
  }
  column "CreatedAt" {
    type = timestamp
  }
  column "UpdatedAt" {
    type = timestamp
  }

  column "ContactId" {
    type = uuid
  }

  column "WhatsappBusinessAccountPhoneNumberId" {
    type = uuid
  }

  column "InitiatedBy" {
    type = enum.ConversationInitiatedEnum
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
    type = uuid
  }
  column "CreatedAt" {
    type = timestamp
  }
  column "UpdatedAt" {
    type = timestamp
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
  }

  column "WhatsappBusinessAccountPhoneNumberId" {
    type = uuid
  }

  column "Direction" {
    type = enum.MessageDirection
  }

  column "Content" {
    type = text
  }

  column "Status" {
    type = enum.MessageStatus
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
    type = uuid
  }
  column "CreatedAt" {
    type = timestamp
  }
  column "UpdatedAt" {
    type = timestamp
  }

  column "CampaignId" {
    type = uuid
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

  index "TrackLinkCampaignIdIndex" {
    columns = [column.CampaignId]

  }
}

table "TrackLinkClick" {
  schema = schema.public
  column "UniqueId" {
    type = uuid
  }
  column "CreatedAt" {
    type = timestamp
  }
  column "UpdatedAt" {
    type = timestamp
  }

  column "TrackLinkId" {
    type = uuid
  }

  column "ContactId" {
    type = uuid
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
    type = uuid
  }
  column "CreatedAt" {
    type = timestamp
  }
  column "UpdatedAt" {
    type = timestamp
  }

  primary_key {
    columns = [column.UniqueId]
  }
}



// ==== JOINT TABLES ======


table "ContactListContact" {
  schema = schema.public
  column "CreatedAt" {
    type = timestamp
  }
  column "UpdatedAt" {
    type = timestamp
  }

  column "ContactListId" {
    type = uuid
  }

  column "ContactId" {
    type = uuid
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
    type = timestamp
  }
  column "UpdatedAt" {
    type = timestamp
  }

  column "ContactListId" {
    type = uuid
  }

  column "TagId" {
    type = uuid
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
    type = timestamp
  }
  column "UpdatedAt" {
    type = timestamp
  }

  column "ContactListId" {
    type = uuid
  }

  column "CampaignId" {
    type = uuid
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
    type = timestamp
  }
  column "UpdatedAt" {
    type = timestamp
  }

  column "ConversationId" {
    type = uuid
  }

  column "TagId" {
    type = uuid
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
    type = timestamp
  }
  column "UpdatedAt" {
    type = timestamp
  }

  column "CampaignId" {
    type = uuid
  }

  column "TagId" {
    type = uuid
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
}

table "MessageReply" {
  schema = schema.public
  column "CreatedAt" {
    type = timestamp
  }
  column "UpdatedAt" {
    type = timestamp
  }

  column "MessageId" {
    type = uuid
  }

  column "ReplyMessageId" {
    type = uuid
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
