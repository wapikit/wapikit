import { NextResponse } from 'next/server'
import { type NextRequest } from 'next/server'
import { decode } from 'jsonwebtoken'

// This function can be marked `async` if using `await` inside
export function middleware(request: NextRequest) {
	const authCookie = request.cookies.get('auth_token')

	if (!authCookie) {
		// NO AUTH COOKIE, MEANS THIS IS A UN-AUTHENTICATED REQ REDIRECT TO LOGIN
		return NextResponse.redirect(new URL('/login', request.url))
	} else {
		// paths which can not be accessed by authed user
		const decoded = decode(authCookie.value)

		if (!decoded || typeof decoded === 'string') {
			return NextResponse.redirect(new URL('/login', request.url))
		}

		let destinationPath = '/dashboard'

		if (decoded['isStudent']) {
			destinationPath = '/student/dashboard'
		} else if (decoded['isInstructor']) {
			destinationPath = '/instructor/dashboard'
		} else if (decoded['isCompanyMember']) {
			destinationPath = '/company/dashboard'
		} else if (decoded['isEducationalInstitutionMember']) {
			destinationPath = '/educational-institution/dashboard'
		}

		return NextResponse.redirect(new URL(destinationPath, request.url))
	}
}

export const config = {
	matcher: [
		{
			source: '/login',
			has: [
				{
					type: 'cookie',
					key: 'auth_token',
					value: '(^[A-Za-z0-9-_]+.[A-Za-z0-9-_]+.[A-Za-z0-9-_.+/=]+$)'
				}
			]
		},
		{
			source: '/dashboard',
			has: [
				{
					type: 'cookie',
					key: 'auth_token',
					value: '(^[A-Za-z0-9-_]+.[A-Za-z0-9-_]+.[A-Za-z0-9-_.+/=]+$)'
				}
			]
		},
		{
			source: '/register',
			has: [
				{
					type: 'cookie',
					key: 'auth_token',
					value: '(^[A-Za-z0-9-_]+.[A-Za-z0-9-_]+.[A-Za-z0-9-_.+/=]+$)'
				}
			]
		},
		{
			source: '/register/oauth',
			has: [
				{
					type: 'cookie',
					key: 'auth_token',
					value: '(^[A-Za-z0-9-_]+.[A-Za-z0-9-_]+.[A-Za-z0-9-_.+/=]+$)'
				}
			]
		},
		{
			source: '/',
			has: [
				{
					type: 'cookie',
					key: 'auth_token',
					value: '(^[A-Za-z0-9-_]+.[A-Za-z0-9-_]+.[A-Za-z0-9-_.+/=]+$)'
				}
			]
		},
		{
			source: '/company',
			has: [
				{
					type: 'cookie',
					key: 'auth_token',
					value: '(^[A-Za-z0-9-_]+.[A-Za-z0-9-_]+.[A-Za-z0-9-_.+/=]+$)'
				}
			]
		},
		{
			source: '/instructor',
			has: [
				{
					type: 'cookie',
					key: 'auth_token',
					value: '(^[A-Za-z0-9-_]+.[A-Za-z0-9-_]+.[A-Za-z0-9-_.+/=]+$)'
				}
			]
		},
		{
			source: '/educational-institution',
			has: [
				{
					type: 'cookie',
					key: 'auth_token',
					value: '(^[A-Za-z0-9-_]+.[A-Za-z0-9-_]+.[A-Za-z0-9-_.+/=]+$)'
				}
			]
		}
	]
}
