import { env } from '$env/dynamic/private';
import type { Cookies } from '@sveltejs/kit';

export function setSessionCookie(cookies: Cookies, token: string, isSecure: boolean): void {
	// Allow disabling secure cookies for HTTP deployments (e.g., behind Tailscale)
	const secure = env.COOKIES_SECURE === 'false' ? false : isSecure;

	cookies.set('session_token', token, {
		path: '/',
		httpOnly: true,
		sameSite: 'lax',
		secure,
		maxAge: 60 * 60 * 24 * 30, // 30 days
	});
}
