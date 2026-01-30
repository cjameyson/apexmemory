import type { ApiCard, ApiFact, ApiFactDetail } from '$lib/api/types';
import type { Card, Fact, FactDetail } from '$lib/types/fact';

export function toCard(api: ApiCard): Card {
	return {
		id: api.id,
		factId: api.fact_id,
		notebookId: api.notebook_id,
		elementId: api.element_id,
		state: api.state,
		stability: api.stability,
		difficulty: api.difficulty,
		due: api.due,
		reps: api.reps,
		lapses: api.lapses,
		suspendedAt: api.suspended_at,
		buriedUntil: api.buried_until,
		createdAt: api.created_at,
		updatedAt: api.updated_at
	};
}

export function toFact(api: ApiFact): Fact {
	return {
		id: api.id,
		notebookId: api.notebook_id,
		factType: api.fact_type,
		content: api.content,
		sourceId: api.source_id,
		cardCount: api.card_count,
		tags: api.tags ?? [],
		dueCount: api.due_count ?? 0,
		createdAt: api.created_at,
		updatedAt: api.updated_at
	};
}

export function toFactDetail(api: ApiFactDetail): FactDetail {
	return {
		...toFact(api),
		cards: api.cards.map(toCard)
	};
}
