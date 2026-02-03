import type { ApiStudyCard } from '$lib/api/types';
import type { StudyCard, CardDisplay } from '$lib/types/review';
import { toast } from 'svelte-sonner';

const RATING_MAP: Record<number, 'again' | 'hard' | 'good' | 'easy'> = {
	1: 'again',
	2: 'hard',
	3: 'good',
	4: 'easy'
};

export function ratingToString(rating: 1 | 2 | 3 | 4): 'again' | 'hard' | 'good' | 'easy' {
	return RATING_MAP[rating];
}

export function toStudyCard(api: ApiStudyCard): StudyCard {
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
		factType: api.fact_type,
		factContent: api.fact_content,
		intervals: api.intervals
	};
}

export function toStudyCards(apiList: ApiStudyCard[]): StudyCard[] {
	return apiList.map(toStudyCard);
}

export async function fetchStudyCards(notebookId?: string): Promise<StudyCard[]> {
	const params = new URLSearchParams({ limit: '50' });
	if (notebookId) params.set('notebook_id', notebookId);

	const res = await fetch(`/api/reviews/study?${params}`);
	if (!res.ok) {
		toast.error('Failed to load study cards. Please try again.');
		return [];
	}
	const data: ApiStudyCard[] = await res.json();
	return toStudyCards(data);
}

export async function fetchPracticeCards(notebookId?: string): Promise<StudyCard[]> {
	const params = new URLSearchParams({ limit: '50' });
	if (notebookId) params.set('notebook_id', notebookId);

	const res = await fetch(`/api/reviews/practice?${params}`);
	if (!res.ok) {
		toast.error('Failed to load practice cards. Please try again.');
		return [];
	}
	const response = await res.json();
	// Practice endpoint returns paginated response
	const data: ApiStudyCard[] = response.data ?? response;
	return toStudyCards(data);
}

/**
 * Extract front/back display text from a study card's fact content.
 */
export function extractCardDisplay(card: StudyCard): CardDisplay {
	const content = card.factContent as { version?: number; fields?: Array<{ name: string; type: string; value: string }> };
	const fields = content?.fields ?? [];

	if (card.factType === 'basic') {
		return {
			front: fields[0]?.value ?? '',
			back: fields[1]?.value ?? ''
		};
	}

	if (card.factType === 'cloze') {
		const textField = fields.find((f) => f.type === 'cloze_text') ?? fields[0];
		const text = textField?.value ?? '';
		const elementId = card.elementId; // e.g. "c1"

		// Extract the answer for the target cloze
		let clozeAnswer = '';
		text.replace(/\{\{(c\d+)::([^}]+)\}\}/g, (_, id, content) => {
			if (id === elementId) clozeAnswer = content;
			return '';
		});

		// Front: replace the target cloze with [...], show others filled
		const front = text.replace(/\{\{(c\d+)::([^}]+)\}\}/g, (_, id, content) => {
			return id === elementId ? '[...]' : content;
		});

		// Back: show all cloze content filled
		const back = text.replace(/\{\{c\d+::([^}]+)\}\}/g, (_, content) => content);

		return { front, back, isCloze: true, clozeAnswer };
	}

	// image_occlusion: placeholder for now
	if (card.factType === 'image_occlusion') {
		const title = fields.find((f) => f.name === 'title')?.value ?? 'Image Occlusion';
		return { front: title, back: title };
	}

	return { front: 'Unknown card type', back: '' };
}
