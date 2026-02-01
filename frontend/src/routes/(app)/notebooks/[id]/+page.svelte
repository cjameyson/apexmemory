<script lang="ts">
	import type { PageData } from './$types';
	import NotebookDashboard from '$lib/components/notebook/notebook-dashboard.svelte';

	let { data }: { data: PageData } = $props();
</script>

<div id="notebook-main" class="flex flex-1 flex-col overflow-auto">
	<NotebookDashboard notebook={data.notebook} sources={data.sources} cards={data.cards} />
</div>
