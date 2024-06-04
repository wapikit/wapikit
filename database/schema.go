package database

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Organization struct {
	gorm.Model
	UniqueId  uint      `json:"uniqueId" gorm:"primarykey,unique"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Name       sql.NullString `json:"name"`
	WebsiteUrl sql.NullString `json:"websiteUrl"`
	LegalName  sql.NullString `json:"legalName"`

	// MetaData `json:"metaData" gorm:"type:jsonb"` // JSONB column definition

	OrganizationMembers []OrganizationMember    `json:"organizationMembers" gorm:"foreignKey:OrganizationId;references:UniqueId"`
	Contacts            []Contact               `json:"contacts" gorm:"foreignKey:OrganizationId;references:UniqueId"`
	ContactLists        []ContactList           `json:"contactLists" gorm:"foreignKey:OrganizationId;references:UniqueId"`
	BusinessAccount     WhatsappBusinessAccount `json:"businessAccount" gorm:"foreignKey:OrganizationId;references:UniqueId"`
}

type OrganizationMember struct {
	gorm.Model

	UniqueId  uint      `json:"uniqueId" gorm:"primarykey,unique"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Name        string         `json:"name"`
	Email       string         `json:"email" gorm:"unique"`
	PhoneNumber sql.NullString `json:"phoneNumber" gorm:"unique"`

	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`

	OrganizationId uint `json:"organizationId" gorm:"index"`

	// this can be a empty array in case of member, because by default members do not have any permissions
	Permissions []string `json:"permissions"`

	Role string `json:"role" gorm:"type:enum('super_admin', 'admin', 'member');default:'member';index"`
}

type WhatsappBusinessAccount struct {
	gorm.Model
	UniqueId  uint      `json:"uniqueId" gorm:"primarykey,unique"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	AccountId      string `json:"accountId" gorm:"unique"` // this is the account Id provided by WhatsApp itself
	OrganizationId uint   `json:"organizationId" gorm:"index"`

	PhoneNumbers []WhatsappBusinessAccountPhoneNumber `json:"phoneNumbers" gorm:"foreignKey:BusinessAccountId;references:UniqueId"`
}

type WhatsappBusinessAccountPhoneNumber struct {
	gorm.Model
	UniqueId  uint      `json:"uniqueId" gorm:"primarykey,unique"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	PhoneNumber string         `json:"phoneNumber"`
	Description sql.NullString `json:"description"` // this is for if in case user want to add some meta description for internal purpose

	BusinessAccountId uint       `json:"businessAccountId" gorm:"index"`
	Messages          []Message  `json:"messages" gorm:"foreignKey:WhatsappBusinessAccountPhoneNumberId;references:UniqueId"`
	Campaigns         []Campaign `json:"campaigns" gorm:"foreignKey:PhoneNumberId;references:UniqueId"`
}

type Contact struct {
	gorm.Model
	UniqueId  uint      `json:"uniqueId" gorm:"primarykey,unique"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	ContactNumber string `json:"contactNumber" gorm:"unique"`
	Name          string `json:"name"`

	Status string `json:"status" gorm:"type:enum('active', 'blocked');default:'active'"` // Active, Blocked

	OrganizationId uint `json:"organizationId" gorm:"index"`

	// MetaData `json:"metaData" gorm:"type:jsonb"` // JSONB column definition

	Lists        []ContactList  `json:"lists" gorm:"many2many:contact_lists;foreignKey:ContactId;References:UniqueId;"`
	Messages     []Message      `json:"messages" gorm:"foreignKey:ContactId;references:UniqueId"`
	Conversation []Conversation `json:"conversation" gorm:"foreignKey:ContactId;references:UniqueId"`
}

type ContactList struct {
	gorm.Model
	UniqueId  uint      `json:"uniqueId" gorm:"primarykey,unique"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Name        string `json:"name"`
	Description string `json:"description"`

	OrganizationId uint `json:"organizationId" gorm:"index"`

	Contacts []Contact `json:"contacts" gorm:"many2many:contact_lists;foreignKey:ListId;references:UniqueId"`
	Tags     []Tag     `json:"tags" gorm:"many2many:contactLists_tags;foreignKey:ContactListId;references:UniqueId"`

	Campaigns []Campaign `json:"campaigns" gorm:"foreignKey:ContactListId;references:UniqueId"`
}

type Campaign struct {
	gorm.Model
	UniqueId  uint      `json:"uniqueId" gorm:"primarykey,unique"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Name string `json:"name"`

	PhoneNumberId uint `json:"phoneNumberId" gorm:"index"` // relate to the phoneNumber table
	ContactListId uint `json:"contactListId" gorm:"index"`

	Messages []Message   `json:"messages" gorm:"foreignKey:CampaignId"`
	Links    []TrackLink `json:"links" gorm:"foreignKey:CampaignId;references:UniqueId"`

	Tags []Tag `json:"tags" gorm:"foreignKey:CampaignId"` // many to many relation need here, a tag can have multiple campaigns and a campaign can have multiple tags
}

type Conversation struct {
	gorm.Model
	UniqueId  uint      `json:"uniqueId" gorm:"primarykey,unique"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	InitiatedBy string `json:"initiatedBy" gorm:"type:enum('user', 'campaign')"` // ! TODO: load enum User, Campaign

	ContactId uint `json:"contactId" gorm:"index"`

	Messages []Message `json:"messages" gorm:"foreignKey:ConversationId"`

	Tags []Tag `json:"tags" gorm:"many2many:conversation_tags;foreignKey:ConversationId;references:UniqueId"`
}

// this table will store the outgoing message and incoming both the messages
// for 1.Campaigns 2.Customer Support Integration
type Message struct {
	gorm.Model

	UniqueId  uint      `json:"uniqueId" gorm:"primarykey,unique"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Status string `json:"status" gorm:"type:enum('sent', 'delivered', 'read', 'failed', 'undelivered');default:'sent';index"`

	Type      string `json:"type" gorm:"type:enum('audio', 'image');default:'text';index"` // ! TODO: load enum Audio, Video, Image, Text, Document, Location, Contact, Template, Product, Product List etc etc
	Direction string `json:"direction" gorm:"type:enum('incoming', 'outgoing');default:'incoming';index"`

	ContactId uint `json:"contactId" gorm:"index"`

	// can be a null string in case if no incoming message, has been initiated from the user end
	ConversationId sql.NullString `json:"conversationId"`

	// populated in case of this is a message sent via a campaign broadcasts
	CampaignId sql.NullString `json:"campaignId"`

	WhatsappBusinessAccountPhoneNumberId uint `json:"whatsappBusinessAccountPhoneNumberId" gorm:"index"`

	Content string `json:"content"` // should be encrypted message

	// this is for a message can have multiple replies, so its a one to many self referential relation
	ReplyToMessageId sql.NullString
	ReplyMessages    []Message `json:"replyMessages" gorm:"foreignKey:ReplyToMessageId;references:UniqueId"`
}

type TrackLink struct {
	gorm.Model

	UniqueId  uint      `json:"uniqueId" gorm:"primarykey,unique"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Destination string `json:"destination"`
	Slug        string `json:"slug"`

	Clicks []TrackLinkClick `json:"clicks" gorm:"foreignKey:TrackLinkId;references:UniqueId"`

	CampaignId string `json:"campaignId" gorm:"index"`
}

type TrackLinkClick struct {
	gorm.Model

	UniqueId  uint      `json:"uniqueId" gorm:"primarykey,unique"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	ClickedByPhoneNumber string `json:"clickedByPhoneNumber"`

	TrackLinkId uint `json:"trackLinkId" gorm:"index"`
}

type Tag struct {
	gorm.Model
	UniqueId  uint      `json:"uniqueId" gorm:"primarykey,unique"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Name string `json:"name" gorm:"unique,index"`

	Conversations []Conversation `json:"conversations" gorm:"many2many:conversation_tags;foreignKey:TagId;references:UniqueId"`
	ContactLists  []ContactList  `json:"contactLists" gorm:"many2many:contactLists_tags;foreignKey:TagId;references:UniqueId"`
}
