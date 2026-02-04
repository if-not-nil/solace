<script lang="ts">
	// TODO: make them render above everything
	import { alerts } from '$lib/alerts';
	import { fly } from 'svelte/transition';

	const typeClasses = {
		error: 'alert-error',
		success: 'alert-success',
		info: 'alert-info',
		warning: 'alert-warning'
	};
</script>

<div class="alert-container">
	{#each $alerts as alert (alert.id)}
		<div
			in:fly={{ x: 200, duration: 200 }}
			out:fly={{ x: 200, duration: 200 }}
			class="alert {typeClasses[alert.type ?? 'error']}"
		>
			{alert.message}
		</div>
	{/each}
</div>

<style>
	.alert-container {
		position: fixed;
		top: 1rem;
		right: 1rem;
		width: 300px;
		z-index: 1000;
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.alert {
		padding: 1rem;
		border-radius: 6px;
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
		font-size: 0.95rem;
		font-weight: 500;
		transition: opacity 0.3s ease;
	}

	.alert-error {
		background-color: var(--base08);
		color: var(--base05);
	}

	.alert-success {
		background-color: var(--base0B);
		color: var(--base00);
	}

	.alert-info {
		background-color: var(--base0D);
		color: var(--base00);
	}

	.alert-warning {
		background-color: var(--base0A);
		color: var(--base01);
	}
</style>
