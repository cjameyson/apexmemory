// Mock sources and cards for development
// Returns sample data for any notebook ID until APIs exist

import type { Source, Card } from '$lib/types';

/**
 * Returns mock sources for any notebook ID.
 * Uses the notebook ID to associate sources with the correct notebook.
 */
export function getMockSourcesForNotebook(notebookId: string): Source[] {
	return [
		{
			id: 'src-textbook-ch1',
			notebookId,
			name: 'Introduction to the Subject',
			type: 'pdf',
			cards: 45,
			excerpt: 'This chapter covers the fundamental concepts and terminology...',
			pages: 87,
			addedAt: '2024-01-15T10:30:00Z'
		},
		{
			id: 'src-lecture-video',
			notebookId,
			name: 'Lecture 1 - Core Concepts',
			type: 'youtube',
			cards: 28,
			excerpt: 'Professor explains the key principles with visual examples...',
			duration: '14:32',
			addedAt: '2024-01-18T14:00:00Z'
		},
		{
			id: 'src-class-notes',
			notebookId,
			name: 'Class Notes - Week 1',
			type: 'notes',
			cards: 33,
			excerpt: 'Key takeaways from the first week of lectures and readings',
			addedAt: '2024-01-20T09:15:00Z'
		},
		{
			id: 'src-reference-article',
			notebookId,
			name: 'Supplementary Reading',
			type: 'url',
			cards: 12,
			excerpt: 'Additional context and real-world applications...',
			addedAt: '2024-01-22T16:45:00Z'
		}
	];
}

/**
 * Returns mock cards for any notebook ID.
 * Generates cards linked to the mock sources above.
 */
export function getMockCardsForNotebook(notebookId: string): Card[] {
	const sources = getMockSourcesForNotebook(notebookId);

	return [
		// Cards from textbook
		{
			id: 'card-1',
			notebookId,
			sourceId: sources[0].id,
			front: 'What is the fundamental unit of study in this subject?',
			back: 'The basic building block that serves as the foundation for all further concepts.',
			due: true,
			interval: '1d',
			tags: ['fundamentals', 'chapter-1']
		},
		{
			id: 'card-2',
			notebookId,
			sourceId: sources[0].id,
			front: 'Define the key terminology introduced in Chapter 1.',
			back: 'A comprehensive definition covering scope, application, and historical context.',
			due: true,
			interval: '2d',
			tags: ['definitions', 'chapter-1']
		},
		{
			id: 'card-3',
			notebookId,
			sourceId: sources[0].id,
			front: 'What are the three main categories discussed?',
			back: '1. Category A - description\n2. Category B - description\n3. Category C - description',
			due: false,
			interval: '7d',
			tags: ['categories', 'chapter-1']
		},
		// Cards from video
		{
			id: 'card-4',
			notebookId,
			sourceId: sources[1].id,
			front: 'According to the lecture, what is the most important concept to understand first?',
			back: 'The foundational principle that everything else builds upon.',
			due: true,
			interval: '1d',
			tags: ['lecture', 'core-concepts']
		},
		{
			id: 'card-5',
			notebookId,
			sourceId: sources[1].id,
			front: 'What example did the professor use to illustrate the main point?',
			back: 'The real-world scenario demonstrating practical application.',
			due: false,
			interval: '14d',
			tags: ['lecture', 'examples']
		},
		// Cards from notes
		{
			id: 'card-6',
			notebookId,
			sourceId: sources[2].id,
			front: 'What were the key takeaways from Week 1?',
			back: '- Point 1: Understanding basics\n- Point 2: Applying concepts\n- Point 3: Common mistakes to avoid',
			due: true,
			interval: '3d',
			tags: ['notes', 'week-1']
		},
		{
			id: 'card-7',
			notebookId,
			sourceId: sources[2].id,
			front: 'What question was raised during the class discussion?',
			back: 'How does this concept apply in edge cases?',
			due: false,
			interval: '21d',
			tags: ['notes', 'discussion']
		},
		// Cards from article
		{
			id: 'card-8',
			notebookId,
			sourceId: sources[3].id,
			front: 'What real-world application was described in the article?',
			back: 'Industry use case showing practical implementation and results.',
			due: true,
			interval: '1d',
			tags: ['article', 'applications']
		}
	];
}

/**
 * Get a specific mock source by ID.
 */
export function getMockSource(notebookId: string, sourceId: string): Source | undefined {
	const sources = getMockSourcesForNotebook(notebookId);
	return sources.find((s) => s.id === sourceId);
}
