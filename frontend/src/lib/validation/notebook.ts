// Notebook validation functions
// Following existing function-based pattern from auth.ts

/**
 * Validate notebook name.
 * @returns Error message if invalid, undefined if valid.
 */
export function validateNotebookName(name: string): string | undefined {
	const trimmed = name.trim();
	if (!trimmed) return 'Name is required';
	if (trimmed.length > 255) return 'Name must not exceed 255 characters';
	return undefined;
}

/**
 * Validate notebook description.
 * @returns Error message if invalid, undefined if valid.
 */
export function validateNotebookDescription(
	description: string | null | undefined
): string | undefined {
	if (!description) return undefined;
	if (description.length > 10000) return 'Description must not exceed 10,000 characters';
	return undefined;
}

/**
 * Validate desired retention value.
 * Frontend constraint: 0.70-0.99 (backend allows 0-1 exclusive).
 * @returns Error message if invalid, undefined if valid.
 */
export function validateDesiredRetention(value: number | undefined): string | undefined {
	if (value === undefined) return undefined;
	if (value < 0.7 || value > 0.99) {
		return 'Desired retention must be between 0.70 and 0.99';
	}
	return undefined;
}

/**
 * Validate all notebook fields for create/update.
 * @returns Object with field errors, or null if all valid.
 */
export function validateNotebook(data: {
	name?: string;
	description?: string | null;
	desiredRetention?: number;
}): Record<string, string> | null {
	const errors: Record<string, string> = {};

	if (data.name !== undefined) {
		const nameError = validateNotebookName(data.name);
		if (nameError) errors.name = nameError;
	}

	if (data.description !== undefined) {
		const descError = validateNotebookDescription(data.description);
		if (descError) errors.description = descError;
	}

	if (data.desiredRetention !== undefined) {
		const retentionError = validateDesiredRetention(data.desiredRetention);
		if (retentionError) errors.desiredRetention = retentionError;
	}

	return Object.keys(errors).length > 0 ? errors : null;
}
