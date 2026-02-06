import type { Handle } from '@sveltejs/kit';
import { apiRequest } from '$lib/server/api';
import type { User } from '$lib/api/types';

function shouldValidateSession(pathname: string): boolean {
	if (pathname.startsWith('/api/')) return false;
	if (pathname.startsWith('/_app/')) return false;
	if (pathname === '/favicon.ico') return false;
	if (pathname === '/robots.txt') return false;
	if (pathname === '/manifest.webmanifest') return false;
	return true;
}

function shouldClearSessionCookie(status: number): boolean {
	return status === 401 || status === 403;
}

export const handle: Handle = async ({ event, resolve }) => {
	const sessionToken = event.cookies.get('session_token');

	event.locals.user = null;
	event.locals.sessionToken = sessionToken ?? null;

	const validateSession =
		!!sessionToken &&
		event.request.method === 'GET' &&
		shouldValidateSession(event.url.pathname);

	if (validateSession && sessionToken) {
		const result = await apiRequest<User>('/v1/auth/me', {
			token: sessionToken,
		});

		if (result.ok) {
			event.locals.user = result.data;
		} else if (shouldClearSessionCookie(result.status)) {
			// Token is invalid, clear the cookie.
			event.cookies.delete('session_token', { path: '/' });
			event.locals.sessionToken = null;
		}
	}

	return resolve(event);
};
