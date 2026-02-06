import { beforeEach, describe, expect, it, vi } from 'vitest';
import { handle } from '../hooks.server';
import { apiRequest } from '$lib/server/api';

vi.mock('$lib/server/api', () => ({
	apiRequest: vi.fn()
}));

const mockedApiRequest = vi.mocked(apiRequest);

interface TestContextOptions {
	method?: string;
	pathname?: string;
	sessionToken?: string;
}

function createContext(options: TestContextOptions = {}) {
	const { method = 'GET', pathname = '/', sessionToken = 'test-session-token' } = options;

	const cookies = {
		get: vi.fn((name: string) => (name === 'session_token' ? sessionToken : undefined)),
		delete: vi.fn()
	};

	const event = {
		request: new Request(`http://localhost${pathname}`, { method }),
		url: new URL(`http://localhost${pathname}`),
		cookies,
		locals: {
			user: null,
			sessionToken: null
		}
	};

	const resolve = vi.fn(async () => new Response('ok'));

	return { event, cookies, resolve };
}

function asHandleArgs(
	event: ReturnType<typeof createContext>['event'],
	resolve: ReturnType<typeof createContext>['resolve']
) {
	return { event, resolve } as unknown as Parameters<typeof handle>[0];
}

describe('hooks.server handle', () => {
	beforeEach(() => {
		mockedApiRequest.mockReset();
	});

	it('does not clear the session cookie on transient backend failures', async () => {
		mockedApiRequest.mockResolvedValue({
			ok: false,
			status: 0,
			error: { error: 'Network error' }
		});

		const { event, cookies, resolve } = createContext({ pathname: '/notebooks' });

		await handle(asHandleArgs(event, resolve));

		expect(mockedApiRequest).toHaveBeenCalledWith('/v1/auth/me', {
			token: 'test-session-token'
		});
		expect(cookies.delete).not.toHaveBeenCalled();
		expect(event.locals.sessionToken).toBe('test-session-token');
		expect(resolve).toHaveBeenCalledTimes(1);
	});

	it('clears the session cookie when the backend returns 401', async () => {
		mockedApiRequest.mockResolvedValue({
			ok: false,
			status: 401,
			error: { error: 'Unauthorized' }
		});

		const { event, cookies, resolve } = createContext({ pathname: '/notebooks' });

		await handle(asHandleArgs(event, resolve));

		expect(cookies.delete).toHaveBeenCalledWith('session_token', { path: '/' });
		expect(event.locals.sessionToken).toBeNull();
		expect(resolve).toHaveBeenCalledTimes(1);
	});

	it('skips auth validation for internal /api routes', async () => {
		const { event, resolve } = createContext({ pathname: '/api/reviews/study' });

		await handle(asHandleArgs(event, resolve));

		expect(mockedApiRequest).not.toHaveBeenCalled();
		expect(resolve).toHaveBeenCalledTimes(1);
	});

	it('skips auth validation for non-GET requests', async () => {
		const { event, resolve } = createContext({ method: 'POST', pathname: '/notebooks' });

		await handle(asHandleArgs(event, resolve));

		expect(mockedApiRequest).not.toHaveBeenCalled();
		expect(resolve).toHaveBeenCalledTimes(1);
	});
});
