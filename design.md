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


