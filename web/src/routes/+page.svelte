<script lang="ts">
	import { SvelteURLSearchParams } from 'svelte/reactivity';
	import Sidebar from '$lib/components/Sidebar.svelte';
	import { location, location_methods } from '$lib/navigation.svelte';
	import { goto } from '$app/navigation';
	import { onMount, onDestroy } from 'svelte';
	import { get_channel, channel_cache } from '$lib/caches.svelte';
	import Message from '$lib/components/Message.svelte';
	import { messaging } from '$lib/messaging.svelte';

	let inputValue = $state('');
	let inputElement = $state<HTMLInputElement>();
	let msgListElement = $state<HTMLElement | null>(null);

	// sync url params with location store
	onMount(() => {
		const params = new SvelteURLSearchParams(window.location.search);
		const serverId = params.get('s');
		const channelId = params.get('c');

		if (channelId) {
			try {
				get_channel(channelId, serverId || undefined).then((channel) => {
					if (channel) {
						location_methods.set_both_id(channel.meta.server_id, channel.meta.id);
					}
				});
			} catch (e) {
				console.error('failed to load channel from url:', e);
			}
		} else if (serverId) {
			try {
				location_methods.set_server_id(serverId);
			} catch (e) {
				console.error('failed to load server from url:', e);
			}
		}

		// setup messaging
		messaging.setup();
	});

	onDestroy(() => {
		messaging.cleanup();
	});

	// add/remove scroll event listener when msgListElement changes
	$effect(() => {
		if (msgListElement) {
			msgListElement.addEventListener(
				'scroll',
				() => messaging.debouncedScrollHandler(msgListElement),
				{ passive: true }
			);
			return () => {
				msgListElement!.removeEventListener('scroll', () =>
					messaging.debouncedScrollHandler(msgListElement)
				);
			};
		}
	});

	// update url when location changes
	$effect(() => {
		if (location && location.channel) {
			let params = new SvelteURLSearchParams(window.location.search);
			params.set('s', location.channel.meta.server_id);
			params.set('c', location.channel.meta.id);
			goto(`?${params.toString()}`);
		}
	});

	// handle channel changes
	$effect(() => {
		const channelID = location.channel?.meta.id || null;
		messaging.initializeForChannel(channelID);
	});

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		const content = inputValue.trim();
		inputValue = '';
		if (!content?.length || !location.channel) return;

		await messaging.sendMessage(content, location.channel.meta.id);
		inputElement?.focus();
	}

	let messages = $derived(location.channel?.msgs);

	// auto-connect websocket when logged in
	$effect(() => {
		messaging.autoConnect();
	});

	// auto-switch channels when connected
	$effect(() => {
		messaging.autoSwitchChannel();
	});
</script>

<main>
	<Sidebar
		onChannelSelect={(e: string) => {
			location_methods.set_channel_id(e);
		}}
		onServerSelect={(e: string) => {
			location_methods.set_server_id(e);
		}}
	/>

	<section class="content">
		<div class="content-wrapper">
			{#if location.channel}
				<div class="msg-list" bind:this={msgListElement}>
					{#if messaging.state.hasError}
						<div class="error-state">
							<p>failed to load messages :(</p>
							<button onclick={() => messaging.loadInitialMessages(location.channel!.meta.id)}>
								retry
							</button>
						</div>
					{:else if messaging.state.isLoading}
						<div class="info-text">
							<p>crickets...</p>
						</div>
					{:else if channel_cache[location.channel.meta.id]?.msgs}
						{#each messages || [] as msg (msg.id)}
							<Message message={msg} serverId={location.channel.meta.server_id} />
							<!-- <Message {msg} user={get_user(msg.user_id)} /> -->
						{/each}
					{:else}
						<p class="no-messages">crickets</p>
					{/if}
				</div>

				{#if location.channel}
					<!-- <div class="input"> -->
					<form onsubmit={handleSubmit} class="input-wrapper">
						<!-- svelte-ignore a11y_autofocus -->
						<input
							autofocus
							bind:value={inputValue}
							bind:this={inputElement}
							placeholder="say something"
						/>
						<button>=></button>
					</form>
					<!-- </div> -->
				{/if}
			{:else}
				<div class="info-text empty-state">
					<h2>welcome!</h2>
					<p>select a channel please</p>
				</div>
			{/if}
		</div>
	</section>
</main>

<style>
	:root {
		--input-height: 3em;
	}
	main {
		display: flex;
		height: 100vh;
		width: 100vw;
		box-sizing: border-box;
		overflow: hidden;
	}

	.content {
		flex: 1;
		display: flex;
		flex-direction: column;
		background-color: var(--base00);
		color: var(--base05);
		position: relative;
		overflow: hidden;
	}

	.content-wrapper {
		display: flex;
		flex-direction: column;
		height: 100%;
		width: 100%;
		position: relative;
	}

	.msg-list {
		display: flex;
		flex-direction: column;
		margin-top: 0;
		margin-bottom: calc(var(--input-height) + 2em);
		overflow-y: auto;
		/* padding: 1em; */
		gap: 0.5em;
		scroll-behavior: smooth;
		font-family: var(--nc-font-mono);
		background: var(--base00);
		min-height: 90vh;
	}

	/* Terminal-like styling for messages */
	.msg-list::-webkit-scrollbar {
		width: 8px;
	}

	.msg-list::-webkit-scrollbar-track {
		background: var(--base01);
	}

	.msg-list::-webkit-scrollbar-thumb {
		background: var(--base03);
		border-radius: 0;
	}

	/* Loading and error states */
	.loading-state,
	.error-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 2em;
		color: var(--base04);
		text-align: center;
	}

	.loading-state p,
	.error-state p {
		margin: 0 0 1em 0;
		font-family: var(--nc-font-mono);
	}

	.error-state button {
		background: var(--base08);
		border: 1px solid var(--base09);
	}

	.error-state button:hover {
		background: var(--base09);
		color: var(--base00);
	}

	.input-wrapper {
		position: fixed;
		bottom: 1em;
		width: calc(100vw - 20em); /* 1em margins */
		left: 19em;
		height: var(--input-height);
		box-sizing: border-box;
		z-index: 2;
		border-radius: 0;
		display: flex;
		flex-direction: row;
	}
	.input-wrapper > input {
		border-radius: 0px;
		flex: 1;
		padding: 0.75em 1em;
		font-size: 1em;
		border: none;
		outline: none;
		border-radius: 0px;
		background: var(--base02);
		color: var(--base05);
		height: 100%;
	}

	/* .input-wrapper input:focus { */
	/* border: 1px solid var(--base03); */
	/* } */
	.input-wrapper button {
		padding: 0.75em 0.75em;
		background-color: var(--base02);
		color: var(--base05);
		border: none;
		border-radius: 0px;
		cursor: pointer;
		font-size: 1em;
		transition: background-color 0.2s;
	}

	.input-wrapper button:hover {
		background-color: var(--base03);
	}

	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		height: 100%;
		opacity: 0.5;
	}
</style>
