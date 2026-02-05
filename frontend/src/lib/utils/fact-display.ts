import type { FactType } from '$lib/types/fact';

interface FactField {
	name: string;
	type: string;
	value: string | Record<string, unknown>;
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

/** Extract plain text from a TipTap JSON document by walking text nodes. */
function extractTextFromDoc(doc: Record<string, unknown>): string {
	const parts: string[] = [];
	function walk(node: Record<string, unknown>) {
		if (node.type === 'text' && typeof node.text === 'string') {
			parts.push(node.text);
		}
		if (Array.isArray(node.content)) {
			for (const child of node.content) {
				walk(child as Record<string, unknown>);
			}
		}
	}
	walk(doc);
	return parts.join(' ').replace(/\s+/g, ' ').trim();
}

function getFieldText(content: FactContent, name: string): string {
	const field = content.fields.find((f) => f.name === name);
	if (!field) return '';
	if (field.type === 'rich_text' && typeof field.value === 'object' && field.value !== null) {
		return extractTextFromDoc(field.value as Record<string, unknown>);
	}
	return stripHtml(String(field.value ?? ''));
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
			const front = getFieldText(content, 'front');
			const back = getFieldText(content, 'back');
			return {
				primary: truncate(front || 'No front text', 120),
				secondary: back ? truncate(back, 80) : null
			};
		}
		case 'cloze': {
			const text = getFieldText(content, 'text');
			const masked = text.replace(/\{\{c\d+::([^}]*?)(?:::[^}]*)?\}\}/g, '[...]');
			return {
				primary: truncate(masked || 'No cloze text', 120),
				secondary: null
			};
		}
		case 'image_occlusion': {
			const title = getFieldText(content, 'title');
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
