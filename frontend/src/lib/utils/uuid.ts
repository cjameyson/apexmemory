import { v7 as uuidv7 } from 'uuid';

/**
 * Generate a UUIDv7 (time-sortable) string.
 * Consistent with backend primary key generation.
 */
export function generateUUID(): string {
	return uuidv7();
}
