import type { CardState, FactType } from './fact';
import type { JSONContent } from '@tiptap/core';
import type { Region, OcclusionMode, RevealStyle } from '$lib/components/image-occlusion/types';

export type ReviewMode = 'scheduled' | 'practice';

export interface StudyCard {
	id: string;
	factId: string;
	notebookId: string;
	elementId: string;
	state: CardState;
	stability: number | null;
	difficulty: number | null;
	due: string | null;
	reps: number;
	lapses: number;
	factType: FactType;
	factContent: Record<string, unknown>;
	intervals: { again: string; hard: string; good: string; easy: string };
}

export interface BasicCardDisplay {
	type: 'basic';
	front: JSONContent | string;
	back: JSONContent | string;
}

export interface ClozeCardDisplay {
	type: 'cloze';
	front: string;
	clozeAnswer: string;
}

export interface ImageOcclusionCardDisplay {
	type: 'image_occlusion';
	title: string;
	imageUrl: string;
	imageWidth: number;
	imageHeight: number;
	imageRotation: 0 | 90 | 180 | 270;
	regions: Region[];
	targetRegionId: string;
	mode: OcclusionMode;
	revealStyle: RevealStyle;
}

export type CardDisplay = BasicCardDisplay | ClozeCardDisplay | ImageOcclusionCardDisplay;
