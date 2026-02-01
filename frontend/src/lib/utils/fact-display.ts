import type { FactType } from '$lib/types/fact';

interface FactField {
	name: string;
	type: string;
	value: string;
}

interface FactContent {
	version: number;
	fields: FactField[];
}

export interface FactDisplayText {
	primary: string;
	secondary: string | null;
}

function truncate(text: string, maxLength: number): string {
	if (text.length <= maxLength) return text;
	return text.slice(0, maxLength).trimEnd() + '...';
}

function getFieldValue(content: FactContent, name: string): string {
	return content.fields.find((f) => f.name === name)?.value ?? '';
}

function stripHtml(html: string): string {
	return html.replace(/<[^>]*>/g, '').trim();
}

export function getFactDisplayText(
	rawContent: Record<string, unknown>,
	factType: FactType
): FactDisplayText {
	const content = rawContent as unknown as FactContent;

	if (!content?.fields?.length) {
		return { primary: 'Empty fact', secondary: null };
	}

	switch (factType) {
		case 'basic': {
			const front = stripHtml(getFieldValue(content, 'front'));
			const back = stripHtml(getFieldValue(content, 'back'));
			return {
				primary: truncate(front || 'No front text', 120),
				secondary: back ? truncate(back, 80) : null
			};
		}
		case 'cloze': {
			const text = stripHtml(getFieldValue(content, 'text'));
			const masked = text.replace(/\{\{c\d+::([^}]*?)(?:::[^}]*)?\}\}/g, '[...]');
			return {
				primary: truncate(masked || 'No cloze text', 120),
				secondary: null
			};
		}
		case 'image_occlusion': {
			const title = stripHtml(getFieldValue(content, 'title'));
			const regions = content.fields.filter((f) => f.type === 'image_occlusion');
			return {
				primary: truncate(title || 'Image occlusion', 120),
				secondary: `${regions.length} region${regions.length !== 1 ? 's' : ''}`
			};
		}
		default:
			return { primary: 'Unknown fact type', secondary: null };
	}
}
