/**
 * Undo/Redo history management using Svelte 5 runes.
 *
 * Uses the command pattern - each undoable action is wrapped in a command object
 * that knows how to execute and undo itself.
 */

import type { EditorCommand } from './types';

/** Maximum number of commands to keep in history */
const MAX_HISTORY_SIZE = 100;

/** Time window (ms) for merging consecutive commands of the same type */
const MERGE_WINDOW_MS = 500;

export interface HistoryManager {
	/** Whether undo is available */
	readonly canUndo: boolean;
	/** Whether redo is available */
	readonly canRedo: boolean;
	/** Number of commands in undo stack */
	readonly undoCount: number;
	/** Number of commands in redo stack */
	readonly redoCount: number;
	/** Execute a command and add to history */
	execute(command: EditorCommand): void;
	/** Undo the last command */
	undo(): void;
	/** Redo the last undone command */
	redo(): void;
	/** Clear all history */
	clear(): void;
}

/**
 * Create a history manager with undo/redo support.
 */
export function createHistoryManager(): HistoryManager {
	// svelte-ignore state_referenced_locally
	let undoStack = $state<EditorCommand[]>([]);
	// svelte-ignore state_referenced_locally
	let redoStack = $state<EditorCommand[]>([]);
	let lastExecuteTime = 0;

	const canUndo = $derived(undoStack.length > 0);
	const canRedo = $derived(redoStack.length > 0);
	const undoCount = $derived(undoStack.length);
	const redoCount = $derived(redoStack.length);

	function execute(command: EditorCommand): void {
		const now = Date.now();

		// Try to merge with the last command if within the time window
		if (
			undoStack.length > 0 &&
			now - lastExecuteTime < MERGE_WINDOW_MS
		) {
			const lastCommand = undoStack[undoStack.length - 1];
			if (
				lastCommand.type === command.type &&
				lastCommand.targetId === command.targetId &&
				lastCommand.merge
			) {
				// Execute the new command
				command.execute();
				// Try to merge - if successful, we don't add a new entry
				if (lastCommand.merge(command)) {
					lastExecuteTime = now;
					return;
				}
			}
		}

		// Execute the command
		command.execute();

		// Add to undo stack
		undoStack = [...undoStack, command];

		// Trim history if it exceeds max size
		if (undoStack.length > MAX_HISTORY_SIZE) {
			undoStack = undoStack.slice(-MAX_HISTORY_SIZE);
		}

		// Clear redo stack (new action invalidates redo history)
		redoStack = [];

		lastExecuteTime = now;
	}

	function undo(): void {
		if (undoStack.length === 0) return;

		const command = undoStack[undoStack.length - 1];
		command.undo();

		undoStack = undoStack.slice(0, -1);
		redoStack = [...redoStack, command];
	}

	function redo(): void {
		if (redoStack.length === 0) return;

		const command = redoStack[redoStack.length - 1];
		command.execute();

		redoStack = redoStack.slice(0, -1);
		undoStack = [...undoStack, command];
	}

	function clear(): void {
		undoStack = [];
		redoStack = [];
		lastExecuteTime = 0;
	}

	return {
		get canUndo() {
			return canUndo;
		},
		get canRedo() {
			return canRedo;
		},
		get undoCount() {
			return undoCount;
		},
		get redoCount() {
			return redoCount;
		},
		execute,
		undo,
		redo,
		clear
	};
}
