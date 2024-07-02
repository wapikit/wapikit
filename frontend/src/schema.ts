import { ContactStatusEnum, UserRoleEnum } from 'root/.generated'
import { z } from 'zod'

export const UserTokenPayloadSchema = z.object({
	unique_id: z.string(),
	username: z.string(),
	email: z.string(),
	role: z.nativeEnum(UserRoleEnum),
	organization_id: z.string(),
	name: z.string()
})

export const NewTeamMemberFormScheam = z.object({
	email: z.string().email({ message: 'Enter a valid email address' }),
	role: z.nativeEnum(UserRoleEnum)
})

export const NewContactFormSchema = z.object({
	name: z.string().min(3, { message: 'Name must be at least 3 characters' }),
	description: z.string().min(3, { message: 'Description must be at least 3 characters' }),
	phone: z.string().min(10, { message: 'Phone number must be at least 10 characters' }),
	lists: z.string().array(),
	status: z.nativeEnum(ContactStatusEnum),
	attributes: z.any()
})

export const NewContactListFormSchema = z.object({
	name: z.string().min(3, { message: 'Name must be at least 3 characters' }),
	description: z.string().min(3, { message: 'Description must be at least 3 characters' }),
	tagIds: z.string().array()
})

export const NewCampaignSchema = z.object({
	name: z.string().min(3, { message: 'Name must be at least 3 characters' }),
	description: z.string().min(3, { message: 'Description must be at least 3 characters' }),
	tags: z.string().array(),
	lists: z.string().array(),
	templateId: z.string(),
	isLinkTrackingEnabled: z.boolean(),
	templateParameter: z.object({
		parameter: z.string(),
		parameterIndex: z.string(),
		parameterType: z.string(),
		value: z.string()
	}),
	schedule: z.object({
		date: z.string(),
		time: z.string()
	})
})
