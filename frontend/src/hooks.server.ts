import type { Handle } from '@sveltejs/kit';
import { apiRequest } from '$lib/server/api';
import type { User } from '$lib/api/types';

export const handle: Handle = async ({ event, resolve }) => {
	const sessionToken = event.cookies.get('session_token');

	event.locals.user = null;
	event.locals.sessionToken = sessionToken ?? null;

	if (sessionToken) {
		const result = await apiRequest<User>('/v1/auth/me', {
			token: sessionToken,
		});

		if (result.ok) {
			event.locals.user = result.data;
		} else {
			// Token is invalid, clear the cookie
			event.cookies.delete('session_token', { path: '/' });
			event.locals.sessionToken = null;
		}
	}

	return resolve(event);
};
