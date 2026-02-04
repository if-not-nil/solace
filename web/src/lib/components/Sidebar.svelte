<script lang="ts">
	import { me } from '$lib/user_store.svelte';
	import { type ServerResponse } from '$lib/orm_types';
	import { location } from '$lib/navigation.svelte';
	import { can } from '$lib/permissions';
	import { goto } from '$app/navigation';
	import ServerModal from './ServerModal.svelte';

	// props
	let {
		onServerSelect,
		onChannelSelect
	}: { onServerSelect: (serverId: string) => void; onChannelSelect: (channelId: string) => void } =
		$props();

	// form state
	let showCreateServer = $state(false);
</script>

{#snippet server_icon(server: ServerResponse, active: boolean)}
	<img
		class={`server-icon ${active ? 'active' : ''}`}
		onclick={() => {
			onServerSelect?.(server.id);
		}}
		src={server.avatar || '/server-placeholder.jpg'}
		alt={server.name}
	/>
{/snippet}

{#snippet server_action(icon: string, action: () => void)}
	<button class="server-action" onclick={action}>
		{icon}
	</button>
{/snippet}

<ServerModal bind:show={showCreateServer} />

<aside class="sidebar">
	<div class="servers">
		{@render server_action('+', () => (showCreateServer = true))}
		{#if me.servers}
			{#each me.servers as server (server.id)}
				{@render server_icon(server, location.server?.id === server.id)}
			{/each}
		{/if}
	</div>

	<div class="channels">
		<div class="server-banner">
			<h2>{location.server?.name || 'server not selected!'}</h2>
		</div>

		<div class="channel-list">
			{#if location.server?.id && can('manage_server', location.server.id)}
				<p style="text-align: center" onclick={() => goto(`/admin?id=${location.server!.id}`)}>
					things for admins
				</p>
			{/if}

			{#if location.server?.channels}
				{#each location.server.channels as channel (channel.id)}
					<p
						class={location.channel?.meta.id === channel.id ? 'active' : ''}
						onclick={() => {
							onChannelSelect?.(channel.id);
						}}
					>
						# {channel.name}
					</p>
				{/each}
			{:else}
				<p class="info-text">no channels available!</p>
			{/if}
		</div>
	</div>

	<div class="user-info">
		<img src={me.user?.avatar || '/flight.jpg'} alt="pfp" />
		<span>
			<a> {me.user?.name || 'huh'} </a>
			<span>online</span>
		</span>
	</div>
</aside>

<style>
	.sidebar {
		position: relative;
		display: flex;
		flex-direction: row;
		height: 100vh;
		width: 18em;
		background-color: var(--base01);
		color: white;
		font-family: var(--nc-font-mono);
	}

	/* servers */
	.servers {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.3em;
		padding: 0.3em 0;
		width: 3em;
		background: var(--base01);
		overflow-y: auto;
	}

	.server-icon {
		width: 85%;
		aspect-ratio: 1/1;
		padding: 0.1em;
		border-radius: 0.8em;
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--base01);
		cursor: pointer;
		font-size: 1.1em;
		transition: all 0.2s;
		flex-shrink: 0;
	}
	.server-icon:hover {
		background: var(--base03);
		border-radius: 0.7em;
		transform: scale(1.1) rotate(5deg);
	}
	.server-icon.active {
		background: var(--base03);
		border-radius: 0;
		border-radius: 0.5em;
		animation: pulse 2s infinite;
	}
	@keyframes pulse {
		0%,
		100% {
			transform: scale(1);
		}
		50% {
			transform: scale(1.05);
		}
	}

	/* channels */
	.channels {
		display: flex;
		flex-direction: column;
		flex: 1;
		background: var(--base01);
	}

	.server-banner {
		padding: 0.3em 0.3em 0.3em 0.5em;
		font-weight: bold;
	}

	.server-banner h2 {
		margin: 0;
		font-size: 1em;
		font-family: var(--nc-font-mono);
		color: var(--base05);
	}

	.channel-list {
		display: flex;
		flex-direction: column;
		overflow-y: auto;
		flex: 1;
		padding: 0.3em 0;
	}
	.channel-list p {
		margin: 0;
		font-family: var(--nc-font-mono);
		padding: 0.3em 0.5em;
		border-radius: 2px;
		cursor: pointer;
		transition: background 0.2s;
		font-size: 0.85em;
	}
	.channel-list p:hover {
		background: var(--base02);
		transform: translateX(3px);
	}
	.channel-list p.active {
		background: var(--base05);
		color: var(--base00);
	}

	.user-info {
		position: absolute;
		bottom: 0;
		left: 0;
		display: flex;
		align-items: center;
		gap: 0.5em;
		padding: 0.5em 0.8em;
		width: 100%;
		background: var(--base02);
		z-index: 5;
		box-sizing: border-box;
	}

	.user-info span {
		background: var(--base02);
		font-size: 0.9em;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
		display: flex;
		flex-direction: column;
	}

	.user-info img {
		width: 2em;
		height: 2em;
		border-radius: 0.5em;
		object-fit: cover;
	}

	/* server action button */
	.server-action {
		width: 85%;
		aspect-ratio: 1/1;
		border-radius: 0.5em;
		padding: 0.1em;
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--base03);
		color: var(--base05);
		border: none;
		cursor: pointer;
		font-size: 1.3em;
		transition: all 0.2s;
		flex-shrink: 0;
		font-family: var(--nc-font-mono);
	}
	.server-action:hover {
		background: var(--base0D);
		color: var(--base00);
		border-radius: 0.2em;
		transform: rotate(90deg) scale(1.1);
	}

	/* modal header styling */
	.modal-header {
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 0.4em;
		border-bottom: 1px solid var(--base02);
		font-family: var(--nc-font-mono);
		font-size: 1em;
		color: var(--base05);
	}

	.modal-header span {
		cursor: pointer;
		transition: all 0.2s;
		padding: 0.2em 0.3em;
		border-radius: 2px;
	}

	.modal-header span:hover {
		background: var(--base02);
		color: var(--base06);
	}

	.modal-header span.active {
		color: var(--base0D);
		font-weight: bold;
	}

	.modal-header span.active {
		color: var(--base0D);
		font-weight: bold;
	}

	.modal .section {
		margin-bottom: 1em;
	}

	.modal .section h3 {
		margin: 0 0 0.8em 0;
		color: var(--base05);
		font-size: 1.1em;
		font-family: var(--nc-font-mono);
	}

	@media (max-width: 768px) {
		.sidebar {
			width: 18em;
		}

		.servers {
			width: 3em;
		}

		.form-group {
			flex-direction: column;
			align-items: stretch;
		}
	}
</style>
