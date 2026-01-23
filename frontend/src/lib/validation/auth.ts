export const USERNAME_PATTERN = /^[a-zA-Z][a-zA-Z0-9_-]*$/;

export function validateEmail(email: string): string | undefined {
	if (!email) return 'Email is required';
	if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) return 'Invalid email format';
	return undefined;
}

export function validateUsername(username: string): string | undefined {
	if (!username) return 'Username is required';
	if (username.length < 3) return 'Username must be at least 3 characters';
	if (username.length > 30) return 'Username must be at most 30 characters';
	if (!USERNAME_PATTERN.test(username)) {
		return 'Username must start with a letter and contain only letters, numbers, underscores, and hyphens';
	}
	return undefined;
}

export function validatePassword(password: string): string | undefined {
	if (!password) return 'Password is required';
	if (password.length < 8) return 'Password must be at least 8 characters';
	if (password.length > 128) return 'Password must be at most 128 characters';
	return undefined;
}
