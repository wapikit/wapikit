import { OrganizationMemberSchemaRole } from 'root/.generated'
import { z } from 'zod'

export const UserTokenPayloadSchema = z.object({
	unique_id: z.string(),
	username: z.string(),
	email: z.string(),
	role: z.nativeEnum(OrganizationMemberSchemaRole),
	organization_id: z.string(),
	name: z.string()
})
