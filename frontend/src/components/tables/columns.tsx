'use client'
import { Checkbox } from '~/components/ui/checkbox'
import { type ColumnDef } from '@tanstack/react-table'
import { type Contact } from '~/types'
import {
	type OrganizationMemberSchema,
	type CampaignSchema,
	type ContactListSchema,
	type ContactSchema,
	type OrganizationRoleSchema
} from 'root/.generated'

export const ContactTableColumns: ColumnDef<ContactSchema>[] = [
	{
		id: 'uniqueId',
		accessorKey: 'uniqueId',
		enableHiding: true,
		size: 0
	},
	{
		id: 'select',
		header: ({ table }) => (
			<Checkbox
				checked={table.getIsAllPageRowsSelected()}
				onCheckedChange={(value: any) => table.toggleAllPageRowsSelected(!!value)}
				aria-label="Select all"
			/>
		),
		cell: ({ row }) => (
			<Checkbox
				checked={row.getIsSelected()}
				onCheckedChange={(value: any) => row.toggleSelected(!!value)}
				aria-label="Select row"
			/>
		),
		enableSorting: false,
		enableHiding: false
	},
	{
		accessorKey: 'name',
		header: 'NAME'
	},
	{
		accessorKey: 'created_at',
		header: 'Created At',
		accessorFn: (originalRow: ContactSchema) => {
			return new Date(originalRow.createdAt).toDateString()
		}
	},
	{
		accessorKey: 'phone',
		header: 'PHONE'
	}
]

export const CampaignTableColumns: ColumnDef<CampaignSchema>[] = [
	{
		id: 'uniqueId',
		accessorKey: 'uniqueId',
		enableHiding: true,
		size: 0
	},
	{
		id: 'select',
		header: ({ table }) => (
			<Checkbox
				checked={table.getIsAllPageRowsSelected()}
				onCheckedChange={(value: any) => table.toggleAllPageRowsSelected(!!value)}
				aria-label="Select all"
			/>
		),
		cell: ({ row }) => (
			<Checkbox
				checked={row.getIsSelected()}
				onCheckedChange={(value: any) => row.toggleSelected(!!value)}
				aria-label="Select row"
			/>
		),
		enableSorting: false,
		enableHiding: false
	},
	{
		accessorKey: 'name',
		header: 'NAME'
	},
	{
		accessorKey: 'created_at',
		header: 'Created At',
		accessorFn: (originalRow: CampaignSchema) => {
			return new Date(originalRow.createdAt).toDateString()
		}
	},
	{
		accessorKey: 'status',
		header: 'STATUS'
	},
	{
		accessorKey: 'listId',
		header: 'LISTS'
	},
	{
		accessorKey: 'tags',
		header: 'TAGS'
	}
]

export const ContactListTableColumns: ColumnDef<ContactListSchema>[] = [
	{
		id: 'uniqueId',
		accessorKey: 'uniqueId',
		enableHiding: true,
		size: 0
	},
	{
		id: 'select',
		header: ({ table }) => (
			<Checkbox
				checked={table.getIsAllPageRowsSelected()}
				onCheckedChange={(value: any) => table.toggleAllPageRowsSelected(!!value)}
				aria-label="Select all"
			/>
		),
		cell: ({ row }) => (
			<Checkbox
				checked={row.getIsSelected()}
				onCheckedChange={(value: any) => row.toggleSelected(!!value)}
				aria-label="Select row"
			/>
		),
		enableSorting: false,
		enableHiding: false
	},
	{
		accessorKey: 'name',
		header: 'NAME'
	},
	{
		accessorKey: 'created_at',
		header: 'Created At',
		accessorFn: (originalRow: ContactListSchema) => {
			return new Date(originalRow.createdAt).toDateString()
		}
	},
	{
		accessorKey: 'numberOfCampaignsSent',
		header: 'Campaigns Sent'
	},
	{
		accessorKey: 'numberOfContacts',
		header: 'Contacts'
	},
	{
		accessorKey: 'tags',
		header: 'TAGS'
	}
]

export const OrganizationMembersTableColumns: ColumnDef<OrganizationMemberSchema>[] = [
	{
		id: 'uniqueId',
		accessorKey: 'uniqueId',
		enableHiding: true,
		size: 0
	},
	{
		id: 'select',
		header: ({ table }) => (
			<Checkbox
				checked={table.getIsAllPageRowsSelected()}
				onCheckedChange={(value: any) => table.toggleAllPageRowsSelected(!!value)}
				aria-label="Select all"
			/>
		),
		cell: ({ row }) => (
			<Checkbox
				checked={row.getIsSelected()}
				onCheckedChange={(value: any) => row.toggleSelected(!!value)}
				aria-label="Select row"
			/>
		),
		enableSorting: false,
		enableHiding: false
	},
	{
		accessorKey: 'name',
		header: 'NAME'
	},
	{
		accessorKey: 'email',
		header: 'EMAIL'
	},
	{
		accessorKey: 'accessLevel',
		header: 'ACCESS LEVEL'
	},
	{
		accessorKey: 'roles',
		header: 'ROLES'
	},
	{
		accessorKey: 'createdAt',
		header: 'Joined At',
		accessorFn: (originalRow: OrganizationMemberSchema) => {
			return new Date(originalRow.createdAt).toDateString()
		}
	}
]

export const RolesTableColumns: ColumnDef<OrganizationRoleSchema>[] = [
	{
		id: 'uniqueId',
		accessorKey: 'uniqueId',
		enableHiding: true,
		size: 0
	},
	{
		id: 'select',
		header: ({ table }) => (
			<Checkbox
				checked={table.getIsAllPageRowsSelected()}
				onCheckedChange={(value: any) => table.toggleAllPageRowsSelected(!!value)}
				aria-label="Select all"
			/>
		),
		cell: ({ row }) => (
			<Checkbox
				checked={row.getIsSelected()}
				onCheckedChange={(value: any) => row.toggleSelected(!!value)}
				aria-label="Select row"
			/>
		),
		enableSorting: false,
		enableHiding: false
	},
	{
		accessorKey: 'name',
		header: 'NAME'
	},
	{
		accessorKey: 'permissions',
		header: 'PERMS',
		enablePinning: true
	}
]

export const columns: ColumnDef<Contact>[] = [
	{
		id: 'uniqueId',
		accessorKey: 'uniqueId',
		enableHiding: true,
		size: 0
	},
	{
		id: 'select',
		header: ({ table }) => (
			<Checkbox
				checked={table.getIsAllPageRowsSelected()}
				onCheckedChange={(value: any) => table.toggleAllPageRowsSelected(!!value)}
				aria-label="Select all"
			/>
		),
		cell: ({ row }) => (
			<Checkbox
				checked={row.getIsSelected()}
				onCheckedChange={(value: any) => row.toggleSelected(!!value)}
				aria-label="Select row"
			/>
		),
		enableSorting: false,
		enableHiding: false
	},
	{
		accessorKey: 'name',
		header: 'NAME'
	},
	{
		accessorKey: 'phone',
		header: 'PHONE'
	},
	{
		accessorKey: 'list',
		header: 'EMAIL'
	}
]
