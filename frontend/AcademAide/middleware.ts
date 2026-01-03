// import { NextResponse } from 'next/server'
// import type { NextRequest } from 'next/server'

// export function middleware(request: NextRequest) {
//     const token = request.cookies.get('token')?.value

//     // Define protected routes
//     const protectedRoutes = ['/dashboard', '/chat', '/timetable', '/recommendations', '/profile']

//     const isProtectedRoute = protectedRoutes.some(route =>
//         request.nextUrl.pathname.startsWith(route)
//     )

//     if (isProtectedRoute && !token) {
//         return NextResponse.redirect(new URL('/login', request.url))
//     }

//     // Redirect to dashboard if logged in and trying to access login/signup
//     if (token && (request.nextUrl.pathname === '/login' || request.nextUrl.pathname === '/signup')) {
//         return NextResponse.redirect(new URL('/dashboard', request.url))
//     }

//     return NextResponse.next()
// }

// export const config = {
//     matcher: ['/((?!api|_next/static|_next/image|favicon.ico).*)'],
// }

import { NextResponse } from "next/server"
import type { NextRequest } from "next/server"

export function middleware(request: NextRequest) {
    // TEMPORARILY DISABLE AUTH FOR DEVELOPMENT
    return NextResponse.next()
}
