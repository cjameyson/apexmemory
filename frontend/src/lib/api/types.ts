// User returned from /v1/auth/me and in auth responses
export interface User {
	id: string;
	email: string;
	username: string;
}

// Response from /v1/auth/login and /v1/auth/register
export interface AuthResponse {
	session_token: string;
	expires_at: string;
	user: User;
}

// API error response structure
export interface ApiError {
	error: string;
	code?: string;
	details?: Record<string, string>;
	fieldErrors?: Record<string, string>;
}

// ============================================================================
// Notebook API Types
// ============================================================================

// Mirror backend FSRSSettings struct
export interface ApiFSRSSettings {
	desired_retention: number;
	version: number;
	params: number[];
	learning_steps: number[];
	relearning_steps: number[];
	maximum_interval: number;
	enable_fuzzing: boolean;
}

// Mirror backend NotebookResponse exactly
export interface ApiNotebook {
	id: string;
	name: string;
	description: string | null;
	emoji: string | null;
	color: string | null;
	fsrs_settings: ApiFSRSSettings;
	position: number;
	created_at: string;
	updated_at: string;
}

// POST /v1/notebooks request body
export interface CreateNotebookRequest {
	name: string;
	description?: string;
	emoji?: string | null;
	color?: string | null;
}

// PATCH /v1/notebooks/{id} request body
export interface UpdateNotebookRequest {
	name?: string;
	description?: string | null; // null = clear, undefined = no change
	desired_retention?: number;
	position?: number;
	emoji?: string | null;
	color?: string | null;
}

// ============================================================================
// Fact API Types
// ============================================================================

export interface ApiFact {
	id: string;
	notebook_id: string;
	fact_type: 'basic' | 'cloze' | 'image_occlusion';
	content: Record<string, unknown>;
	source_id: string | null;
	card_count: number;
	tags?: string[];
	due_count?: number;
	created_at: string;
	updated_at: string;
}

export interface ApiFactDetail extends ApiFact {
	cards: ApiCard[];
}

export interface ApiCard {
	id: string;
	fact_id: string;
	notebook_id: string;
	element_id: string;
	state: 'new' | 'learning' | 'review' | 'relearning';
	stability: number | null;
	difficulty: number | null;
	due: string | null;
	reps: number;
	lapses: number;
	suspended_at: string | null;
	buried_until: string | null;
	created_at: string;
	updated_at: string;
}

export interface CreateFactRequest {
	fact_type?: 'basic' | 'cloze' | 'image_occlusion';
	content: { version: number; fields: unknown[] };
}

export interface UpdateFactRequest {
	content: { version: number; fields: unknown[] };
}

export interface UpdateFactResponse {
	fact_id: string;
	created: number;
	deleted: number;
	unchanged: number;
}

export interface PaginatedResponse<T> {
	data: T[];
	total: number;
	has_more: boolean;
}

// Stats returned when ?stats=true on facts list
export interface ApiFactStats {
	total_facts: number;
	total_cards: number;
	total_due: number;
	by_type: {
		basic: number;
		cloze: number;
		image_occlusion: number;
	};
}

// Extended facts list response when stats=true
export interface ApiFactsListWithStats extends PaginatedResponse<ApiFact> {
	stats: ApiFactStats;
}

// ============================================================================
// Review API Types
// ============================================================================

// Study card returned by GET /v1/reviews/study and /v1/reviews/practice
export interface ApiStudyCard extends ApiCard {
	fact_type: 'basic' | 'cloze' | 'image_occlusion';
	fact_content: Record<string, unknown>;
	intervals: { again: string; hard: string; good: string; easy: string };
}

export interface ApiReviewRequest {
	id: string;
	card_id: string;
	rating: 'again' | 'hard' | 'good' | 'easy';
	duration_ms?: number;
	mode: 'scheduled' | 'practice';
}

export interface ApiReviewResponse {
	review: {
		id: string;
		card_id: string;
		rating: string;
		mode: string;
		reviewed_at: string;
	};
	card: ApiCard;
}

// Study counts response from GET /v1/reviews/study-counts
export interface ApiStudyCountsResponse {
	counts: Record<string, { due: number; new: number; total: number }>;
	total_due: number;
	total_new: number;
}

// Undo review response from DELETE /v1/reviews/{id}
export interface ApiUndoReviewResponse {
	card: ApiCard | null; // null for practice mode
}

// Discriminated union for API results
export type ApiResult<T> =
	| { ok: true; data: T }
	| { ok: false; error: ApiError; status: number };
