-- Add new schema named "public"
CREATE SCHEMA IF NOT EXISTS "public";
-- Create "OrganizationIntegration" table
CREATE TABLE "public"."OrganizationIntegration" (
  "UniqueId" uuid NOT NULL DEFAULT gen_random_uuid(),
  "CreatedAt" timestamptz NOT NULL DEFAULT now(),
  "UpdatedAt" timestamptz NOT NULL,
  PRIMARY KEY ("UniqueId")
);
-- Create enum type "OrganizaRolePermissionEnum"
CREATE TYPE "public"."OrganizaRolePermissionEnum" AS ENUM ('Get:OrganizationMember', 'Create:OrganizationMember', 'Update:OrganizationMember', 'Delete:OrganizationMember', 'Get:Campaign', 'Create:Campaign', 'Update:Campaign', 'Delete:Campaign', 'Get:Conversation', 'Update:Conversation', 'Delete:Conversation', 'Assign:Conversation', 'Unassign:Conversation', 'Get:List', 'Create:List', 'Update:List', 'Delete:List', 'Get:Tag', 'Create:Tag', 'Update:Tag', 'Delete:Tag', 'Get:ApiKey', 'Regenerate:ApiKey', 'Get:AppSettings', 'Update:AppSettings', 'Get:Contact', 'Create:Contact', 'Update:Contact', 'Delete:Contact', 'BulkImport:Contacts', 'Get:PrimaryAnalytics', 'Get:SecondaryAnalytics', 'Get:CampaignAnalytics', 'Update:Organization', 'Get:OrganizationRole', 'Create:OrganizationRole', 'Update:OrganizationRole', 'Delete:OrganizationRole', 'Update:IntegrationSettings', 'Get:MessageTemplates', 'Get:PhoneNumbers');
-- Create enum type "OrganizationInviteStatusEnum"
CREATE TYPE "public"."OrganizationInviteStatusEnum" AS ENUM ('Pending', 'Redeemed');
-- Create enum type "ContactStatus"
CREATE TYPE "public"."ContactStatus" AS ENUM ('Active', 'Inactive', 'Blocked', 'Deleted');
-- Create enum type "UserPermissionLevel"
CREATE TYPE "public"."UserPermissionLevel" AS ENUM ('Owner', 'Member');
-- Create enum type "UserAccountStatusEnum"
CREATE TYPE "public"."UserAccountStatusEnum" AS ENUM ('Active', 'Deleted', 'Suspended');
-- Create enum type "OauthProviderEnum"
CREATE TYPE "public"."OauthProviderEnum" AS ENUM ('Google');
-- Create enum type "AccessLogType"
CREATE TYPE "public"."AccessLogType" AS ENUM ('WebInterface', 'ApiAccess');
-- Create enum type "ConversationInitiatedEnum"
CREATE TYPE "public"."ConversationInitiatedEnum" AS ENUM ('Contact', 'Campaign');
-- Create enum type "MessageStatus"
CREATE TYPE "public"."MessageStatus" AS ENUM ('Sent', 'Delivered', 'Read', 'Failed', 'UnDelivered');
-- Create enum type "MessageDirection"
CREATE TYPE "public"."MessageDirection" AS ENUM ('InBound', 'OutBound');
-- Create enum type "ConversationStatus"
CREATE TYPE "public"."ConversationStatus" AS ENUM ('Active', 'Closed', 'Deleted');
-- Create "Organization" table
CREATE TABLE "public"."Organization" (
  "UniqueId" uuid NOT NULL DEFAULT gen_random_uuid(),
  "CreatedAt" timestamptz NOT NULL DEFAULT now(),
  "UpdatedAt" timestamptz NOT NULL,
  "Name" text NOT NULL,
  "Description" text NULL,
  "WebsiteUrl" text NULL,
  "LogoUrl" text NULL,
  "FaviconUrl" text NOT NULL,
  "SlackWebhookUrl" text NULL,
  "DiscordWebhookUrl" text NULL,
  "SmtpClientHost" text NULL,
  "SmtpClientUsername" text NULL,
  "SmtpClientPassword" text NULL,
  PRIMARY KEY ("UniqueId")
);
-- Create "Integration" table
CREATE TABLE "public"."Integration" (
  "UniqueId" uuid NOT NULL DEFAULT gen_random_uuid(),
  "CreatedAt" timestamptz NOT NULL DEFAULT now(),
  "UpdatedAt" timestamptz NOT NULL,
  PRIMARY KEY ("UniqueId")
);
-- Create enum type "CampaignStatus"
CREATE TYPE "public"."CampaignStatus" AS ENUM ('Draft', 'Running', 'Finished', 'Paused', 'Cancelled', 'Scheduled');
-- Create "User" table
CREATE TABLE "public"."User" (
  "UniqueId" uuid NOT NULL DEFAULT gen_random_uuid(),
  "CreatedAt" timestamptz NOT NULL DEFAULT now(),
  "UpdatedAt" timestamptz NOT NULL,
  "Name" text NOT NULL,
  "Email" text NOT NULL,
  "PhoneNumber" text NULL,
  "Username" text NOT NULL,
  "Password" text NULL,
  "OauthProvider" "public"."OauthProviderEnum" NULL,
  "ProfilePictureUrl" text NULL,
  "Status" "public"."UserAccountStatusEnum" NOT NULL,
  PRIMARY KEY ("UniqueId")
);
-- Create index "UserEmailIndex" to table: "User"
CREATE UNIQUE INDEX "UserEmailIndex" ON "public"."User" ("Email");
-- Create index "UserUsernameIndex" to table: "User"
CREATE UNIQUE INDEX "UserUsernameIndex" ON "public"."User" ("Username");
-- Create "OrganizationMemberInvite" table
CREATE TABLE "public"."OrganizationMemberInvite" (
  "UniqueId" uuid NOT NULL DEFAULT gen_random_uuid(),
  "CreatedAt" timestamptz NOT NULL DEFAULT now(),
  "UpdatedAt" timestamptz NOT NULL,
  "Slug" text NOT NULL,
  "email" text NOT NULL,
  "AccessLevel" "public"."UserPermissionLevel" NOT NULL,
  "OrganizationId" uuid NOT NULL,
  "Status" "public"."OrganizationInviteStatusEnum" NOT NULL DEFAULT 'Pending',
  "InvitedByUserId" uuid NOT NULL,
  PRIMARY KEY ("UniqueId"),
  CONSTRAINT "OrganizationToOrganizationInviteForeignKey" FOREIGN KEY ("OrganizationId") REFERENCES "public"."Organization" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "UserToOrganizationInviteForeignKey" FOREIGN KEY ("InvitedByUserId") REFERENCES "public"."User" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "OrganizationInviteInvitedByUserIdIndex" to table: "OrganizationMemberInvite"
CREATE INDEX "OrganizationInviteInvitedByUserIdIndex" ON "public"."OrganizationMemberInvite" ("InvitedByUserId");
-- Create index "OrganizationInviteOrganizationIdIndex" to table: "OrganizationMemberInvite"
CREATE INDEX "OrganizationInviteOrganizationIdIndex" ON "public"."OrganizationMemberInvite" ("OrganizationId");
-- Create "OrganizationMember" table
CREATE TABLE "public"."OrganizationMember" (
  "UniqueId" uuid NOT NULL DEFAULT gen_random_uuid(),
  "CreatedAt" timestamptz NOT NULL DEFAULT now(),
  "UpdatedAt" timestamptz NOT NULL,
  "AccessLevel" "public"."UserPermissionLevel" NOT NULL,
  "OrganizationId" uuid NOT NULL,
  "UserId" uuid NOT NULL,
  "InviteId" uuid NULL,
  PRIMARY KEY ("UniqueId"),
  CONSTRAINT "OrganizationInviteToOrganizationMemberForeignKey" FOREIGN KEY ("InviteId") REFERENCES "public"."OrganizationMemberInvite" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "OrganizationMemberToUserForeignKey" FOREIGN KEY ("UserId") REFERENCES "public"."User" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "OrganizationToOrganizationMemberForeignKey" FOREIGN KEY ("OrganizationId") REFERENCES "public"."Organization" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "OrganizationMemberOrganizationIdIndex" to table: "OrganizationMember"
CREATE INDEX "OrganizationMemberOrganizationIdIndex" ON "public"."OrganizationMember" ("OrganizationId");
-- Create index "OrganizationMemberUserIdIndex" to table: "OrganizationMember"
CREATE INDEX "OrganizationMemberUserIdIndex" ON "public"."OrganizationMember" ("UserId");
-- Create "ApiKey" table
CREATE TABLE "public"."ApiKey" (
  "UniqueId" uuid NOT NULL DEFAULT gen_random_uuid(),
  "CreatedAt" timestamptz NOT NULL DEFAULT now(),
  "UpdatedAt" timestamptz NOT NULL,
  "MemberId" uuid NOT NULL,
  "Key" text NOT NULL,
  "OrganizationId" uuid NOT NULL,
  PRIMARY KEY ("UniqueId"),
  CONSTRAINT "ApiKeyToOrganizationForeignKey" FOREIGN KEY ("OrganizationId") REFERENCES "public"."Organization" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "ApiKeyToOrganizationMemberForeignKey" FOREIGN KEY ("MemberId") REFERENCES "public"."OrganizationMember" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "ApiKeyIndex" to table: "ApiKey"
CREATE UNIQUE INDEX "ApiKeyIndex" ON "public"."ApiKey" ("Key");
-- Create index "ApiKeyOrganizationIdIndex" to table: "ApiKey"
CREATE INDEX "ApiKeyOrganizationIdIndex" ON "public"."ApiKey" ("OrganizationId");
-- Create index "ApiKeyOrganizationMemberIdIndex" to table: "ApiKey"
CREATE UNIQUE INDEX "ApiKeyOrganizationMemberIdIndex" ON "public"."ApiKey" ("MemberId");
-- Create "Campaign" table
CREATE TABLE "public"."Campaign" (
  "UniqueId" uuid NOT NULL DEFAULT gen_random_uuid(),
  "CreatedAt" timestamptz NOT NULL DEFAULT now(),
  "UpdatedAt" timestamptz NOT NULL,
  "Description" text NULL,
  "Name" text NOT NULL,
  "Status" "public"."CampaignStatus" NOT NULL DEFAULT 'Draft',
  "LastContactSent" uuid NULL,
  "IsLinkTrackingEnabled" boolean NOT NULL DEFAULT false,
  "CreatedByOrganizationMemberId" uuid NOT NULL,
  "OrganizationId" uuid NOT NULL,
  "MessageTemplateId" text NULL,
  "PhoneNumber" text NOT NULL,
  "TemplateMessageComponentParameters" jsonb NULL,
  PRIMARY KEY ("UniqueId"),
  CONSTRAINT "CampaignToOrganizationMemberForeignKey" FOREIGN KEY ("CreatedByOrganizationMemberId") REFERENCES "public"."OrganizationMember" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "OrganizationToCampaignForeignKey" FOREIGN KEY ("OrganizationId") REFERENCES "public"."Organization" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "CampaignCreatedByOrganizationMemberIdIndex" to table: "Campaign"
CREATE INDEX "CampaignCreatedByOrganizationMemberIdIndex" ON "public"."Campaign" ("CreatedByOrganizationMemberId");
-- Create index "CampaignMessageTemplateIndex" to table: "Campaign"
CREATE INDEX "CampaignMessageTemplateIndex" ON "public"."Campaign" ("MessageTemplateId");
-- Create "ContactList" table
CREATE TABLE "public"."ContactList" (
  "UniqueId" uuid NOT NULL DEFAULT gen_random_uuid(),
  "CreatedAt" timestamptz NOT NULL DEFAULT now(),
  "UpdatedAt" timestamptz NOT NULL,
  "OrganizationId" uuid NOT NULL,
  "Name" text NOT NULL,
  PRIMARY KEY ("UniqueId"),
  CONSTRAINT "OrganizationToContactListForeignKey" FOREIGN KEY ("OrganizationId") REFERENCES "public"."Organization" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "ContactListOrganizationIdIndex" to table: "ContactList"
CREATE INDEX "ContactListOrganizationIdIndex" ON "public"."ContactList" ("OrganizationId");
-- Create "CampaignList" table
CREATE TABLE "public"."CampaignList" (
  "CreatedAt" timestamptz NOT NULL DEFAULT now(),
  "UpdatedAt" timestamptz NOT NULL,
  "ContactListId" uuid NOT NULL,
  "CampaignId" uuid NOT NULL,
  PRIMARY KEY ("ContactListId", "CampaignId"),
  CONSTRAINT "CampaignListToCampaignForeignKey" FOREIGN KEY ("CampaignId") REFERENCES "public"."Campaign" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "CampaignListToContactListForeignKey" FOREIGN KEY ("ContactListId") REFERENCES "public"."ContactList" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "Tag" table
CREATE TABLE "public"."Tag" (
  "UniqueId" uuid NOT NULL DEFAULT gen_random_uuid(),
  "CreatedAt" timestamptz NOT NULL DEFAULT now(),
  "UpdatedAt" timestamptz NOT NULL,
  "Label" text NOT NULL,
  "Slug" text NOT NULL,
  "OrganizationId" uuid NOT NULL,
  PRIMARY KEY ("UniqueId"),
  CONSTRAINT "UniqueSlug" UNIQUE ("Slug"),
  CONSTRAINT "TagToOrganizationForeignKey" FOREIGN KEY ("OrganizationId") REFERENCES "public"."Organization" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "TagLabelOragnizationIdUniqueIndex" to table: "Tag"
CREATE UNIQUE INDEX "TagLabelOragnizationIdUniqueIndex" ON "public"."Tag" ("Label", "OrganizationId");
-- Create index "TagOrganizationIdIndex" to table: "Tag"
CREATE INDEX "TagOrganizationIdIndex" ON "public"."Tag" ("OrganizationId");
-- Create index "slugIndex" to table: "Tag"
CREATE INDEX "slugIndex" ON "public"."Tag" ("Slug");
-- Create "CampaignTag" table
CREATE TABLE "public"."CampaignTag" (
  "CreatedAt" timestamptz NOT NULL DEFAULT now(),
  "UpdatedAt" timestamptz NOT NULL,
  "CampaignId" uuid NOT NULL,
  "TagId" uuid NOT NULL,
  PRIMARY KEY ("CampaignId", "TagId"),
  CONSTRAINT "CampaignTagToCampaignForeignKey" FOREIGN KEY ("CampaignId") REFERENCES "public"."Campaign" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "CampaignTagToTagForeignKey" FOREIGN KEY ("TagId") REFERENCES "public"."Tag" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "CampaignTagIdCampaignIdUniqueIndex" to table: "CampaignTag"
CREATE UNIQUE INDEX "CampaignTagIdCampaignIdUniqueIndex" ON "public"."CampaignTag" ("CampaignId", "TagId");
-- Create "Contact" table
CREATE TABLE "public"."Contact" (
  "UniqueId" uuid NOT NULL DEFAULT gen_random_uuid(),
  "CreatedAt" timestamptz NOT NULL DEFAULT now(),
  "UpdatedAt" timestamptz NOT NULL,
  "OrganizationId" uuid NOT NULL,
  "Status" "public"."ContactStatus" NOT NULL,
  "Name" text NOT NULL,
  "PhoneNumber" text NOT NULL,
  "Attributes" jsonb NULL,
  PRIMARY KEY ("UniqueId"),
  CONSTRAINT "OrganizationToContactForeignKey" FOREIGN KEY ("OrganizationId") REFERENCES "public"."Organization" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "ContactNumberOrganizationIdUniqueIndex" to table: "Contact"
CREATE UNIQUE INDEX "ContactNumberOrganizationIdUniqueIndex" ON "public"."Contact" ("OrganizationId", "PhoneNumber");
-- Create index "ContactOrganizationIdIndex" to table: "Contact"
CREATE INDEX "ContactOrganizationIdIndex" ON "public"."Contact" ("OrganizationId");
-- Create index "ContactPhoneNumberIndex" to table: "Contact"
CREATE UNIQUE INDEX "ContactPhoneNumberIndex" ON "public"."Contact" ("PhoneNumber");
-- Create "ContactListContact" table
CREATE TABLE "public"."ContactListContact" (
  "CreatedAt" timestamptz NOT NULL DEFAULT now(),
  "UpdatedAt" timestamptz NOT NULL,
  "ContactListId" uuid NOT NULL,
  "ContactId" uuid NOT NULL,
  PRIMARY KEY ("ContactListId", "ContactId"),
  CONSTRAINT "ContactListContactToContactForeignKey" FOREIGN KEY ("ContactId") REFERENCES "public"."Contact" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "ContactListContactToContactListForeignKey" FOREIGN KEY ("ContactListId") REFERENCES "public"."ContactList" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "ContactListTag" table
CREATE TABLE "public"."ContactListTag" (
  "CreatedAt" timestamptz NOT NULL DEFAULT now(),
  "UpdatedAt" timestamptz NOT NULL,
  "ContactListId" uuid NOT NULL,
  "TagId" uuid NOT NULL,
  PRIMARY KEY ("ContactListId", "TagId"),
  CONSTRAINT "ContactListTagToContactListForeignKey" FOREIGN KEY ("ContactListId") REFERENCES "public"."ContactList" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "ContactListTagToTagForeignKey" FOREIGN KEY ("TagId") REFERENCES "public"."Tag" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "Conversation" table
CREATE TABLE "public"."Conversation" (
  "UniqueId" uuid NOT NULL DEFAULT gen_random_uuid(),
  "CreatedAt" timestamptz NOT NULL DEFAULT now(),
  "UpdatedAt" timestamptz NOT NULL,
  "ContactId" uuid NOT NULL,
  "OrganizationId" uuid NOT NULL,
  "Status" "public"."ConversationStatus" NOT NULL,
  "PhoneNumberUsed" text NOT NULL,
  "InitiatedBy" "public"."ConversationInitiatedEnum" NOT NULL,
  "InitiatedByCampaignId" uuid NULL,
  PRIMARY KEY ("UniqueId"),
  CONSTRAINT "ConversationToContactForeignKey" FOREIGN KEY ("ContactId") REFERENCES "public"."Contact" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "ConversationToOrganizationForeignKey" FOREIGN KEY ("OrganizationId") REFERENCES "public"."Organization" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "ConversationContactIdIndex" to table: "Conversation"
CREATE INDEX "ConversationContactIdIndex" ON "public"."Conversation" ("ContactId");
-- Create "ConversationTag" table
CREATE TABLE "public"."ConversationTag" (
  "CreatedAt" timestamptz NOT NULL DEFAULT now(),
  "UpdatedAt" timestamptz NOT NULL,
  "ConversationId" uuid NOT NULL,
  "TagId" uuid NOT NULL,
  PRIMARY KEY ("ConversationId", "TagId"),
  CONSTRAINT "ConversationTagToConversationForeignKey" FOREIGN KEY ("ConversationId") REFERENCES "public"."Conversation" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "ConversationTagToTagForeignKey" FOREIGN KEY ("TagId") REFERENCES "public"."Tag" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "Message" table
CREATE TABLE "public"."Message" (
  "UniqueId" uuid NOT NULL DEFAULT gen_random_uuid(),
  "CreatedAt" timestamptz NOT NULL DEFAULT now(),
  "UpdatedAt" timestamptz NOT NULL,
  "ConversationId" uuid NULL,
  "CampaignId" uuid NULL,
  "ContactId" uuid NOT NULL,
  "PhoneNumberUsed" text NOT NULL,
  "Direction" "public"."MessageDirection" NOT NULL,
  "Content" text NULL,
  "OrganizationId" uuid NOT NULL,
  "Status" "public"."MessageStatus" NOT NULL,
  "RepliedTo" uuid NULL,
  PRIMARY KEY ("UniqueId"),
  CONSTRAINT "MessageToCampaignForeignKey" FOREIGN KEY ("CampaignId") REFERENCES "public"."Campaign" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "MessageToContactForeignKey" FOREIGN KEY ("ContactId") REFERENCES "public"."Contact" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "MessageToConversationForeignKey" FOREIGN KEY ("ConversationId") REFERENCES "public"."Conversation" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "MessageToOrganizationForeignKey" FOREIGN KEY ("OrganizationId") REFERENCES "public"."Organization" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "MessageCampaignIdIndex" to table: "Message"
CREATE INDEX "MessageCampaignIdIndex" ON "public"."Message" ("CampaignId");
-- Create index "MessageContactIdIndex" to table: "Message"
CREATE INDEX "MessageContactIdIndex" ON "public"."Message" ("ContactId");
-- Create "Notification" table
CREATE TABLE "public"."Notification" (
  "UniqueId" uuid NOT NULL DEFAULT gen_random_uuid(),
  "CreatedAt" timestamptz NOT NULL DEFAULT now(),
  "UpdatedAt" timestamptz NOT NULL,
  "ctaUrl" text NULL,
  "title" text NOT NULL,
  "description" text NOT NULL,
  "type" text NULL,
  "isBroadcast" boolean NOT NULL DEFAULT false,
  "UserId" uuid NULL,
  PRIMARY KEY ("UniqueId")
);
-- Create "NotificationReadLog" table
CREATE TABLE "public"."NotificationReadLog" (
  "UniqueId" uuid NOT NULL DEFAULT gen_random_uuid(),
  "CreatedAt" timestamptz NOT NULL DEFAULT now(),
  "UpdatedAt" timestamptz NOT NULL,
  "ReadByUserId" uuid NOT NULL,
  "NotificationId" uuid NOT NULL,
  PRIMARY KEY ("UniqueId"),
  CONSTRAINT "NotificationReadLogToNotificationForeignKey" FOREIGN KEY ("NotificationId") REFERENCES "public"."Notification" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "NotificationReadLogToUserForeignKey" FOREIGN KEY ("ReadByUserId") REFERENCES "public"."User" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "NotificationReadLogNotificationIdIndex" to table: "NotificationReadLog"
CREATE INDEX "NotificationReadLogNotificationIdIndex" ON "public"."NotificationReadLog" ("NotificationId");
-- Create index "NotificationReadLogReadByUserIdIndex" to table: "NotificationReadLog"
CREATE INDEX "NotificationReadLogReadByUserIdIndex" ON "public"."NotificationReadLog" ("ReadByUserId");
-- Create "OrganizationRole" table
CREATE TABLE "public"."OrganizationRole" (
  "UniqueId" uuid NOT NULL DEFAULT gen_random_uuid(),
  "CreatedAt" timestamptz NOT NULL DEFAULT now(),
  "UpdatedAt" timestamptz NOT NULL,
  "Name" text NOT NULL,
  "Description" text NULL,
  "Permissions" text NOT NULL DEFAULT '',
  "OrganizationId" uuid NOT NULL,
  PRIMARY KEY ("UniqueId"),
  CONSTRAINT "OrganizationToOrganizationRoleForeignKey" FOREIGN KEY ("OrganizationId") REFERENCES "public"."Organization" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "OrganizationRoleOrganizationIdIndex" to table: "OrganizationRole"
CREATE INDEX "OrganizationRoleOrganizationIdIndex" ON "public"."OrganizationRole" ("OrganizationId");
-- Create "RoleAssignment" table
CREATE TABLE "public"."RoleAssignment" (
  "UniqueId" uuid NOT NULL DEFAULT gen_random_uuid(),
  "CreatedAt" timestamptz NOT NULL DEFAULT now(),
  "UpdatedAt" timestamptz NOT NULL,
  "OrganizationRoleId" uuid NOT NULL,
  "OrganizationMemberId" uuid NOT NULL,
  PRIMARY KEY ("UniqueId"),
  CONSTRAINT "OrganizationMemberToRoleAssignmentForeignKey" FOREIGN KEY ("OrganizationMemberId") REFERENCES "public"."OrganizationMember" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "OrganizationRoleToRoleAssignmentForeignKey" FOREIGN KEY ("OrganizationRoleId") REFERENCES "public"."OrganizationRole" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "RoleAssignmentOrganizationMemberIdIndex" to table: "RoleAssignment"
CREATE INDEX "RoleAssignmentOrganizationMemberIdIndex" ON "public"."RoleAssignment" ("OrganizationMemberId");
-- Create index "RoleAssignmentOrganizationRoleIdIndex" to table: "RoleAssignment"
CREATE INDEX "RoleAssignmentOrganizationRoleIdIndex" ON "public"."RoleAssignment" ("OrganizationRoleId");
-- Create index "RoleAssignmentUniqueIndex" to table: "RoleAssignment"
CREATE UNIQUE INDEX "RoleAssignmentUniqueIndex" ON "public"."RoleAssignment" ("OrganizationRoleId", "OrganizationMemberId");
-- Create "TrackLink" table
CREATE TABLE "public"."TrackLink" (
  "UniqueId" uuid NOT NULL DEFAULT gen_random_uuid(),
  "CreatedAt" timestamptz NOT NULL DEFAULT now(),
  "UpdatedAt" timestamptz NOT NULL,
  "OrganizationId" uuid NOT NULL,
  "CampaignId" uuid NOT NULL,
  "Slug" text NOT NULL,
  "DestinationUrl" text NULL,
  PRIMARY KEY ("UniqueId"),
  CONSTRAINT "TrackLinkToCampaignForeignKey" FOREIGN KEY ("CampaignId") REFERENCES "public"."Campaign" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "TrackLinkToOrganizationForeignKey" FOREIGN KEY ("OrganizationId") REFERENCES "public"."Organization" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "TrackLinkCampaignIdIndex" to table: "TrackLink"
CREATE INDEX "TrackLinkCampaignIdIndex" ON "public"."TrackLink" ("CampaignId");
-- Create "TrackLinkClick" table
CREATE TABLE "public"."TrackLinkClick" (
  "UniqueId" uuid NOT NULL DEFAULT gen_random_uuid(),
  "CreatedAt" timestamptz NOT NULL DEFAULT now(),
  "UpdatedAt" timestamptz NOT NULL,
  "TrackLinkId" uuid NOT NULL,
  "ContactId" uuid NOT NULL,
  PRIMARY KEY ("UniqueId"),
  CONSTRAINT "TrackLinkClickToContactForeignKey" FOREIGN KEY ("ContactId") REFERENCES "public"."Contact" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "TrackLinkClickToTrackLinkForeignKey" FOREIGN KEY ("TrackLinkId") REFERENCES "public"."TrackLink" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "TrackLinkClickContactIdIndex" to table: "TrackLinkClick"
CREATE INDEX "TrackLinkClickContactIdIndex" ON "public"."TrackLinkClick" ("ContactId");
-- Create index "TrackLinkClickTrackLinkIdIndex" to table: "TrackLinkClick"
CREATE INDEX "TrackLinkClickTrackLinkIdIndex" ON "public"."TrackLinkClick" ("TrackLinkId");
-- Create "WhatsappBusinessAccount" table
CREATE TABLE "public"."WhatsappBusinessAccount" (
  "UniqueId" uuid NOT NULL DEFAULT gen_random_uuid(),
  "CreatedAt" timestamptz NOT NULL DEFAULT now(),
  "UpdatedAt" timestamptz NOT NULL,
  "AccountId" text NOT NULL,
  "AccessToken" text NOT NULL,
  "WebhookSecret" text NOT NULL,
  "OrganizationId" uuid NOT NULL,
  PRIMARY KEY ("UniqueId"),
  CONSTRAINT "WhatsappBusinessAccountToOrganizationForeignKey" FOREIGN KEY ("OrganizationId") REFERENCES "public"."Organization" ("UniqueId") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "WhatsappBusinessAccountAccountIdIndex" to table: "WhatsappBusinessAccount"
CREATE UNIQUE INDEX "WhatsappBusinessAccountAccountIdIndex" ON "public"."WhatsappBusinessAccount" ("AccountId");
-- Create index "WhatsappBusinessAccountOrganizationIdIndex" to table: "WhatsappBusinessAccount"
CREATE INDEX "WhatsappBusinessAccountOrganizationIdIndex" ON "public"."WhatsappBusinessAccount" ("OrganizationId");
