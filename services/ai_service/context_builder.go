package ai_service

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/google/uuid"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/wapikit/wapikit/.db-generated/model"
	"github.com/wapikit/wapikit/.db-generated/table"
)

type ContextBuilder struct {
	Db     *sql.DB
	Logger *slog.Logger
}

func NewContextBuilder(db *sql.DB, logger *slog.Logger) *ContextBuilder {
	return &ContextBuilder{
		Db: db,
	}
}

type ContextCampaign struct {
	model.Campaign
	MessagesSent           int    `json:"messagesSent"`
	MessagesRead           int    `json:"messagesRead"`
	MessagesReplied        int    `json:"messagesReplied"`
	TemplateUsed           string `json:"templateUsed"`
	MessagesFailedToBeSent int    `json:"messagesFailedToBeSent"`
	Tags                   []struct {
		model.Tag
	}
	Lists []struct {
		model.ContactList
		NumberOfContacts int `json:"numberOfContacts"`
	}
}

type ContextConversation struct {
	model.Conversation
	Contact    model.Contact `json:"contact"`
	AssignedTo struct {
		Member model.OrganizationMember `json:"member"`
	}
	Messages []model.Message `json:"messages"`
}

type UnifiedContext struct {
	Campaigns     *[]ContextCampaign     `json:"campaigns"`
	Conversations *[]ContextConversation `json:"conversations"`
}

func (cb *ContextBuilder) containsIntent(intent *DetectIntentResponse, target UserQueryIntent) bool {
	if intent == nil {
		return false
	}

	primary := strings.ToLower(string(intent.PrimaryIntent))
	t := strings.ToLower(string(target))

	// Check for exact equality.
	if primary == t {
		return true
	}

	if strings.HasSuffix(primary, "s") && primary[:len(primary)-1] == t {
		return true
	}

	if strings.HasSuffix(t, "s") && t[:len(t)-1] == primary {
		return true
	}

	return false
}

func (cb *ContextBuilder) fetchCampaignData(ctx context.Context, organizationId uuid.UUID, temporalRange *TemporalRange) (*[]ContextCampaign, error) {
	var dest []ContextCampaign

	whereCondition := table.Campaign.OrganizationId.EQ(UUID(organizationId))
	if !temporalRange.Start.IsZero() {
		whereCondition = whereCondition.AND(table.Campaign.CreatedAt.GT_EQ(TimestampzT(*temporalRange.Start)))
	}

	if !temporalRange.End.IsZero() {
		whereCondition = whereCondition.AND(table.Campaign.CreatedAt.LT_EQ(TimestampzT(*temporalRange.End)))
	}

	campaignQuery := SELECT(
		table.Campaign.AllColumns,
		table.Tag.AllColumns,
		table.CampaignList.AllColumns,
		table.ContactList.AllColumns,
		table.CampaignTag.AllColumns,
		COUNT(table.Campaign.UniqueId).OVER().AS("totalCampaigns"),
	).
		FROM(table.Campaign.
			LEFT_JOIN(table.CampaignTag, table.CampaignTag.CampaignId.EQ(table.Campaign.UniqueId)).
			LEFT_JOIN(table.Tag, table.Tag.UniqueId.EQ(table.CampaignTag.TagId)).
			LEFT_JOIN(table.CampaignList, table.CampaignList.CampaignId.EQ(table.Campaign.UniqueId)).
			LEFT_JOIN(table.ContactList, table.ContactList.UniqueId.EQ(table.CampaignList.ContactListId)),
		).
		WHERE(whereCondition)

	err := campaignQuery.QueryContext(ctx, cb.Db, &dest)

	if err != nil {
		return nil, err
	}

	return &dest, nil
}

func (cb *ContextBuilder) fetchConversationData(organizationId uuid.UUID) (*[]ContextConversation, error) {
	var dest []ContextConversation

	// ! TODO: implement this

	_ = dest

	return &dest, nil
}

// fetchRelevantContext builds a unified context by fetching both campaign and conversation data,
// and then generating proactive insights (if needed).
func (cb *ContextBuilder) fetchRelevantContext(organizationId uuid.UUID, intent DetectIntentResponse) string {
	unifiedCtx := UnifiedContext{}
	if cb.containsIntent(&intent, UserIntentCampaigns) {
		campaigns, err := cb.fetchCampaignData(context.Background(), organizationId, &intent.TemporalContext)
		if err != nil {
			unifiedCtx.Campaigns = &[]ContextCampaign{}
		} else {
			unifiedCtx.Campaigns = campaigns
		}
	}

	if cb.containsIntent(&intent, UserIntentConversation) {
		conversations, err := cb.fetchConversationData(organizationId)
		if err != nil {
			unifiedCtx.Conversations = &[]ContextConversation{}
		} else {
			unifiedCtx.Conversations = conversations
		}
	}

	data, err := json.Marshal(unifiedCtx)
	if err != nil {
		return ""
	}
	return string(data)
}
