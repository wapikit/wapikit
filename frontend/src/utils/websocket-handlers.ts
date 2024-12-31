import type { z } from 'zod';
import type { WebsocketEventDataMap, WebsocketEventEnum } from '../websocket-events'

export async function messageEventHandler(message: z.infer<typeof WebsocketEventDataMap[WebsocketEventEnum.MessageEvent]['shape']['data']>): Promise<boolean> {
    try {

        const { conversationId } = message

        // ! get the conversation from the store
        // ! append the above message to the conversation
        // ! return true if the message was appended successfully

        return true

    } catch (error) {
        console.error(error)
        return false
    }
}

export async function conversationAssignedEventHandler(message: z.infer<typeof WebsocketEventDataMap[WebsocketEventEnum.ConversationAssignmentEvent]['shape']['data']>): Promise<boolean> {
    try {

        const { conversationId } = message

        // ! get the conversation from the store
        // ! update the conversation with the new assignee

        // ! show a notification that the conversation has been assigned if the conversation is assigned to the current user


        return true

        // ! append the above message to the conversation
    } catch (error) {
        console.error(error)
        return false
    }
}

export async function conversationUnassignedEventHandler(message: z.infer<typeof WebsocketEventDataMap[WebsocketEventEnum.ConversationClosedEvent]['shape']['data']>) {
    try {

        const { conversationId } = message

        return true
    } catch (error) {
        console.error(error)
        return false
    }
}




