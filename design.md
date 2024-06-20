## Application Design

## User flows:

#### Organization Configs

- Organization config and super admin can be configured at the time of application setup in the config.toml file.
- Super admin can admins, by default admin has all the permissions, but the super admin can limit every individual permission for a admin.
- By default organization member of type member have no permission, super admin or admin with `MemberWrite` permission has to grant the suitable permission required.
- A single hosting of the application can only have one organization and only one business accounts but may have multiple phone number.
-
- Super Admin or other member with suitable permission can add contacts.
- Super Admin or other member with suitable permission can create contact list and do the other write and read operations.

##### Permission System:

- Two entities only Admin and team member
- Admin can make other people admin
- rest admin can create roles
- while creating roles admin can select granular permission
- user can be alloted with roles

### Settings

- Can Add phone numbers
- Can update whatsapp profile
- API tokens, if the organization want to add the contact or something similar from their technical product on say user sign up.

### Dashboard

- Quality Rating
- Current Quota and Remaining Quota
- Campaigns Meta
- Contacts Meta
- Aggregate Metric Graphs

#### Messaging, Campaigns and Chat Configs

- Each time a contact replies to a message, or initiates a new conversation would be converted to a conversation, a conversation is assignable to members in the organization and by default who ever from the organization member opens/reads the message for the first time, would be assigned the same conversation.
- Contact who only have been sent outgoing message won't be considered as an conversation and just be a contact in the list and the part of the campaign as contact list subscriber.
- User can add auto reply for particular messages.
- User can create a campaign with a selected template message, and at the time of template message selection, they would be asked to add the parameter if applicable to be used for the template message.
- User can test the message before starting the real campaign.
- While creating campaign user can add a opt-out keyword. which will unsubscribe the user from the contact list.

### Media Library

- Each Phone number has its own media library.
- Media can be uploaded using a public link or uploading a file directly.
- Each Media has a media Id provided by whatsapp only.
- Media can be drag and dropped in the conversations to send.
- Media would be stored on the same hosting server mounting a static file path, accessible via the same central http server , on path `/media/*`

## Technical Details and Feasibility:

### Campaign:

- To get a campaign details, read all the message from the message table which is linked to a campaign and on the basis of status calculate the delivery rate, read rate, fail rate etc.
- After a campaign has been started or finished, user can retarget via selecting the user, which will create another campaign but with only those users, a list will also be created with some tags such as "retarget" , and a linked campaign id.


### Service access and communication: 

- suppose if a admin updated 

### FEATURE FLAGS:

- IS_ROLE_BASED_ACCESS_CONTROL_ENABLED
  communityVersion: true
  default: true
  free: false
- IS_SINGLE_BINARY_MODE_ENABLED
  communityVersion: true
  default: true
  hostedVersion: false

- IS_INTEGRATIONS_ENABLED

- IS_QUICK_KEYWORD_REPLIES_ENABLED

- IS_SELF_HOSTED
  default: true
  managedHosting: false

### Campaign Manager

- When user starts a campaign it should start sending messages.
- When user schedules a message it should on time change the status of the campaign to running and then start sending messages.
- When user pauses the campaign, it should stop sending the messages to the users and when user reruns the campaign it should consider the message which have already been sent, and start sending messages again to the remaining user.
- It should keep the whatsapp API rate limits in mind.
- It should retry if a whatsapp message does not get delivered due to network issue or any other non critical error code.
- At the end of all messages has been sent, it should change the status of the campaign to finished.
