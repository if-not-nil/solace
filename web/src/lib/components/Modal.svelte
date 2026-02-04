<script>
	let { showModal = $bindable(), header = null, children } = $props();

	let dialog = $state(); // HTMLDialogElement

	$effect(() => {
		if (showModal) dialog.showModal();
	});
</script>

<dialog
	class="modal"
	bind:this={dialog}
	onclose={() => (showModal = false)}
	onclick={(e) => {
		if (e.target === dialog) dialog.close();
	}}
>
	<div>
		{#if header}
			<div class="modal-header">{@render header?.()}</div>
		{/if}
		<div class="modal-content">{@render children?.()}</div>
	</div>
</dialog>

<style>
	/* modal container */
	.modal {
		background: var(--base00);
		border: 1px solid var(--base02);
		border-radius: 6px;
		width: 90%;
		max-width: 600px;
		max-height: 90vh;
		overflow-y: auto;
	}

	/* modal header */
	.modal-header {
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 0.4em;
		/* border-bottom: 1px solid var(--base02); */
		font-family: var(--nc-font-mono);
		font-size: 1em;
		color: var(--base05);
	}

	/* modal content */
	.modal-content {
		padding: 0.8em;
	}

	/* responsive */
	@media (max-width: 768px) {
		.modal {
			width: 95%;
			margin: 1em;
		}
	}

	/* legacy dialog support */
	dialog {
		border-radius: 0.2em;
		border: none;
		padding: 0;
		background: var(--base00);
		color: var(--base07);
	}
	dialog::backdrop {
		background: rgba(0, 0, 0, 0.3);
		backdrop-filter: blur(4px);
		transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
	}
	dialog > div {
		padding: 1em;
	}
	dialog[open] {
		animation: zoom 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
	}
	@keyframes zoom {
		0% {
			transform: scale(0.8) rotate(-2deg);
		}
		50% {
			transform: scale(1.02) rotate(1deg);
		}
		100% {
			transform: scale(1) rotate(0deg);
		}
	}
	dialog[open]::backdrop {
		animation: fade 0.2s ease-out;
	}
	@keyframes fade {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}
</style>
