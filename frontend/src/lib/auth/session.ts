import type { Cookies } from '@sveltejs/kit';

export function setSessionCookie(cookies: Cookies, token: string, isSecure: boolean): void {
	cookies.set('session_token', token, {
		path: '/',
		httpOnly: true,
		sameSite: 'lax',
		secure: isSecure,
		maxAge: 60 * 60 * 24 * 30, // 30 days
	});
}
