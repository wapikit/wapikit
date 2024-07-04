'use client'
import { Checkbox } from '~/components/ui/checkbox'
import { type ColumnDef } from '@tanstack/react-table'
import { type Contact } from '~/types'
import {
	type OrganizationMemberSchema,
	type CampaignSchema,
	type ContactListSchema,
	type ContactSchema
} from 'root/.generated'

export const ContactTableColumns: ColumnDef<ContactSchema>[] = [
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
		header: 'Created At'
	},
	{
		accessorKey: 'phone',
		header: 'PHONE'
	}
	// {
	// 	id: 'actions',
	// 	cell: ({ row }) => <CellAction data={row.original} />
	// }
]

export const CampaignTableColumns: ColumnDef<CampaignSchema>[] = [
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
		header: 'Created At'
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

export const ContactListTableColumns: ColumnDef<ContactListSchema>[] = []

export const OrganizationMembersTableColumns: ColumnDef<OrganizationMemberSchema>[] = [
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
	}
]

export const RolesTableColumns: ColumnDef<OrganizationMemberSchema>[] = [
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
	}
]

export const columns: ColumnDef<Contact>[] = [
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
