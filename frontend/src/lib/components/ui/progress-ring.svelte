<script lang="ts">
	import { cn } from '$lib/utils';

	interface Props {
		progress: number; // 0-100
		size?: number;
		stroke?: number;
		trackColor?: string;
		progressColor?: string;
		class?: string;
	}

	let {
		progress,
		size = 80,
		stroke = 3,
		trackColor = 'currentColor',
		progressColor = 'currentColor',
		class: className
	}: Props = $props();

	// SVG calculations
	let radius = $derived((36 - stroke) / 2);
	let circumference = $derived(2 * Math.PI * radius);
	let dashOffset = $derived(circumference - (progress / 100) * circumference);

	// Path for the circle
	let circlePath = $derived(
		`M18 ${stroke / 2 + (36 - stroke) / 2 - radius} a ${radius} ${radius} 0 0 1 0 ${radius * 2} a ${radius} ${radius} 0 0 1 0 ${-radius * 2}`
	);
</script>

<svg
	class={cn('-rotate-90', className)}
	width={size}
	height={size}
	viewBox="0 0 36 36"
	aria-valuenow={progress}
	aria-valuemin={0}
	aria-valuemax={100}
	role="progressbar"
>
	<!-- Track circle -->
	<circle
		class="opacity-20"
		cx="18"
		cy="18"
		r={radius}
		stroke={trackColor}
		stroke-width={stroke}
		fill="none"
	/>

	<!-- Progress circle -->
	<circle
		cx="18"
		cy="18"
		r={radius}
		stroke={progressColor}
		stroke-width={stroke}
		fill="none"
		stroke-dasharray={circumference}
		stroke-dashoffset={dashOffset}
		stroke-linecap="round"
		class="transition-[stroke-dashoffset] duration-500 ease-out"
	/>
</svg>
