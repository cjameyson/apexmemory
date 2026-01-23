import { fail, redirect } from '@sveltejs/kit';
import { dev } from '$app/environment';
import type { Actions } from './$types';
import { apiRequest } from '$lib/server/api';
import type { AuthResponse } from '$lib/api/types';
import { validateEmail, validatePassword } from '$lib/validation/auth';
import { setSessionCookie } from '$lib/auth/session';

interface FieldErrors {
	email?: string;
	password?: string;
}

interface FormValues {
	email: string;
}

export const actions: Actions = {
	default: async ({ request, cookies }) => {
		const formData = await request.formData();
		const email = formData.get('email')?.toString().trim() ?? '';
		const password = formData.get('password')?.toString() ?? '';
		const values: FormValues = { email };

		const fieldErrors: FieldErrors = {};
		const emailError = validateEmail(email);
		const passwordError = validatePassword(password);

		if (emailError) fieldErrors.email = emailError;
		if (passwordError) fieldErrors.password = passwordError;

		if (Object.keys(fieldErrors).length > 0) {
			return fail(400, { fieldErrors, values });
		}

		const result = await apiRequest<AuthResponse>('/v1/auth/login', {
			method: 'POST',
			body: { email, password },
		});

		if (!result.ok) {
			let error = 'Invalid email or password';
			if (result.status === 0) {
				error = 'Unable to connect to server. Please try again.';
			}
			return fail(result.status || 400, { error, values });
		}

		setSessionCookie(cookies, result.data.session_token, !dev);

		redirect(302, '/dashboard');
	},
};
