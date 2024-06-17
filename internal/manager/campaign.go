package manager

// - it checks for the running status campaigns in db
// - it queues the message to be sent to user for the campaigns that are running
// - it checks if the campaign has ended and updates the status in db
// - on every message send it updates the last user id for the campaign in the db
// - it fetches the next batch of users to send the message to
// - it updates the campaign status to completed in db
// - it runs all this in memory
// - it must be executed in a go routine because it a long running blocking function, which continuously check for campaign and messages to be sent.

type CampaignManager struct {
}

func NewCampaignManager() *CampaignManager {
	return &CampaignManager{}
}

// Run starts the campaign manager
// main blocking function must be executed in a go routine
func (cm *CampaignManager) Run() {
	// scan campaign
}
