<script lang="ts">
	let { children } = $props();
</script>

<div class="flex flex-1 items-center justify-center px-4 py-8">
	<div class="w-full max-w-md">
		{@render children()}
	</div>
</div>
