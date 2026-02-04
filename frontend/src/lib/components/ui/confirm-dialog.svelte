<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog';
	import Button from '$lib/components/ui/button/button.svelte';

	let {
		open = $bindable(false),
		title,
		description,
		confirmLabel = 'Delete',
		cancelLabel = 'Cancel',
		loading = false,
		onconfirm,
		oncancel
	}: {
		open: boolean;
		title: string;
		description: string;
		confirmLabel?: string;
		cancelLabel?: string;
		loading?: boolean;
		onconfirm: () => void;
		oncancel: () => void;
	} = $props();
</script>

<Dialog.Root bind:open onOpenChange={(v) => !v && oncancel()}>
	<Dialog.Content class="sm:max-w-md">
		<Dialog.Header>
			<Dialog.Title>{title}</Dialog.Title>
			<Dialog.Description>{description}</Dialog.Description>
		</Dialog.Header>
		<Dialog.Footer>
			<Button variant="outline" onclick={oncancel} disabled={loading}>
				{cancelLabel}
			</Button>
			<Button variant="destructive" onclick={onconfirm} disabled={loading}>
				{loading ? 'Deleting...' : confirmLabel}
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
