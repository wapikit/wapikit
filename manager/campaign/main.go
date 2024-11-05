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

	// ! we need to adhere the rate limits of whatsapp api

	// ! call the whatsapp api for the account to check fo the current limit

	// ! check for the running status campaigns in db

	

	//

	// for p := range m.nextPipes {
	// 	has, err := p.NextSubscribers()
	// 	if err != nil {
	// 		m.log.Printf("error processing campaign batch (%s): %v", p.camp.Name, err)
	// 		continue
	// 	}

	// 	if has {
	// 		// There are more subscribers to fetch. Queue again.
	// 		select {
	// 		case m.nextPipes <- p:
	// 		default:
	// 		}
	// 	} else {
	// 		// Mark the pseudo counter that's added in makePipe() that is used
	// 		// to force a wait on a pipe.
	// 		p.wg.Done()
	// 	}
	// }

	// scan for new running campaigns

	// scan campaign
}

func (cm *CampaignManager) Stop() {
	// stop the campaign manager
}

// ! TODO: possibly we will need a worker for every campaign
