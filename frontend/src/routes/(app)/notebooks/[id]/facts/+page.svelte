<script lang="ts">
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<div id="notebook-main" class="flex flex-1 flex-col overflow-auto p-6">
	<h1 class="text-2xl font-bold">Facts</h1>
	<p class="text-muted-foreground text-sm">{data.notebook.name}</p>
</div>
