import { fail, redirect } from '@sveltejs/kit';
import { dev } from '$app/environment';
import type { Actions } from './$types';
import { apiRequest } from '$lib/server/api';
import type { AuthResponse } from '$lib/api/types';
import { validateEmail, validateUsername, validatePassword } from '$lib/validation/auth';
import { setSessionCookie } from '$lib/auth/session';

interface FieldErrors {
	email?: string;
	username?: string;
	password?: string;
	confirmPassword?: string;
}

interface FormValues {
	email: string;
	username: string;
}

export const actions: Actions = {
	default: async ({ request, cookies }) => {
		const formData = await request.formData();
		const email = formData.get('email')?.toString().trim() ?? '';
		const username = formData.get('username')?.toString().trim() ?? '';
		const password = formData.get('password')?.toString() ?? '';
		const confirmPassword = formData.get('confirmPassword')?.toString() ?? '';
		const values: FormValues = { email, username };

		const fieldErrors: FieldErrors = {};

		const emailError = validateEmail(email);
		const usernameError = validateUsername(username);
		const passwordError = validatePassword(password);

		if (emailError) fieldErrors.email = emailError;
		if (usernameError) fieldErrors.username = usernameError;
		if (passwordError) fieldErrors.password = passwordError;

		if (password && confirmPassword && password !== confirmPassword) {
			fieldErrors.confirmPassword = 'Passwords do not match';
		} else if (!confirmPassword) {
			fieldErrors.confirmPassword = 'Please confirm your password';
		}

		if (Object.keys(fieldErrors).length > 0) {
			return fail(400, { fieldErrors, values });
		}

		const result = await apiRequest<AuthResponse>('/v1/auth/register', {
			method: 'POST',
			body: { email, username, password },
		});

		if (!result.ok) {
			// Check for error codes from backend
			if (result.error.code === 'EMAIL_EXISTS') {
				return fail(400, {
					fieldErrors: { email: 'This email is already registered' } as FieldErrors,
					values,
				});
			}
			if (result.error.code === 'USERNAME_EXISTS') {
				return fail(400, {
					fieldErrors: { username: 'This username is already taken' } as FieldErrors,
					values,
				});
			}

			let error = 'Registration failed. Please try again.';
			if (result.status === 0) {
				error = 'Unable to connect to server. Please try again.';
			} else if (result.error.error) {
				error = result.error.error;
			}

			return fail(result.status || 400, { error, values });
		}

		setSessionCookie(cookies, result.data.session_token, !dev);

		redirect(302, '/dashboard');
	},
};
