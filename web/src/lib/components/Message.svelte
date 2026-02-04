<script lang="ts">
	import { get_user, channel_cache } from '$lib/caches.svelte';
	import UserProfileModal from '$lib/components/UserProfileModal.svelte';
	import { api } from '$lib/api';
	import { messaging } from '$lib/messaging.svelte';
	import { alerts } from '$lib/alerts';
	import { can } from '$lib/permissions';
	import { me } from '$lib/user_store.svelte';
	import type { MessageHistoryResponse, UserResponse } from '$lib/orm_types';

	const { message, serverId }: { message: MessageHistoryResponse; serverId?: string } = $props();

	const linkify = (text: string) =>
		text.replace(/(https?:\/\/[^\s]+)/g, '<a href="$1" target="_blank">$1</a>');

	const detectImage = (text: string) => {
		const m = text.match(
			/((?:https?:\/\/)?(?:localhost|[\d.]+|[\w.-]+)\S*\.(?:jpg|jpeg|png|gif|webp))/i
		);
		return m ? m[0] : null;
	};

	const nameColor = (name: string) => {
		const baseSuffixes = ['08', '09', '0A', '0B', '0C', '0D', '0E', '0F'];
		return name.length ? baseSuffixes[name[0].charCodeAt(0) % 8] : '08';
	};

	const attachment = detectImage(message.content);
	let showModal = $state(false);
	let showMenu = $state(false);
	let selectedUser = $state<UserResponse | null>(null);
	let isEditing = $state(false);
	let editContent = $state(message.content);
	let editInputElement = $state<HTMLInputElement>();

	$effect(() => {
		editContent = message.content;
	});

	async function handleDeleteMessage(messageId: string) {
		const messageToDelete = message;

		const currentChannel = window.location.search.match(/c=([^&]*)/)?.[1];
		if (!currentChannel) {
			alerts.add('unable to determine current channel', 'error');
			return;
		}

		const channelRecord = channel_cache[currentChannel];
		if (channelRecord?.msgs) {
			channelRecord.msgs = channelRecord.msgs.filter((msg) => msg.id !== messageId);
		}

		showMenu = false;

		try {
			if (!messageId.startsWith('optimistic-')) {
				await api.message.delete(messageId);
			}
			alerts.add('message deleted', 'success');
		} catch (error) {
			if (channelRecord?.msgs) {
				channelRecord.msgs = [...channelRecord.msgs, messageToDelete].sort(
					(a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
				);
			}
			alerts.add(`failed to delete message: ${error}`, 'error');
		}
	}

	function startEditing() {
		isEditing = true;
		editContent = message.content;
		showMenu = false;
		setTimeout(() => editInputElement?.focus(), 0);
	}

	function cancelEdit() {
		isEditing = false;
		editContent = message.content;
	}

	async function saveEdit(e: Event) {
		e.preventDefault();
		e.stopPropagation();

		const currentChannel = window.location.search.match(/c=([^&]*)/)?.[1];
		if (!currentChannel) {
			alerts.add('unable to determine current channel', 'error');
			return;
		}

		if (!editContent.trim()) {
			cancelEdit();
			return;
		}

		await messaging.editMessage(message.id, editContent.trim(), currentChannel);
		isEditing = false;
	}

	function handleUsernameClick(user: UserResponse) {
		selectedUser = user;
		showModal = true;
	}

	function handleMenuClick(user: UserResponse) {
		selectedUser = user;
		showMenu = !showMenu;
	}

	function handleClickOutside(event: MouseEvent) {
		if (showMenu && !(event.target as Element).closest('.message-menu')) {
			showMenu = false;
		}
	}

	$effect(() => {
		if (showMenu) {
			document.addEventListener('click', handleClickOutside);
			return () => document.removeEventListener('click', handleClickOutside);
		}
	});
</script>

<div class="message">
	<div class="message-content">
		<div class="message-header">
			{#await get_user(message.user_id) then user}
				<span
					onclick={() => handleUsernameClick(user)}
					class="username"
					style="color: var(--base{nameColor(user.name)})"
				>
					{user.name}
				</span>
			{/await}
			{@render message_menu()}
		</div>
		<div class="message-body">
			{#if isEditing}
				<form onsubmit={saveEdit} class="edit-form">
					<input
						bind:value={editContent}
						bind:this={editInputElement}
						class="edit-input"
						placeholder="edit message"
						onkeydown={(e) => {
							if (e.key === 'Escape') cancelEdit();
						}}
					/>
					<div class="edit-actions">
						<button type="submit" class="edit-save">save</button>
						<span class="t-slash">slash</span>
						<button type="button" onclick={cancelEdit} class="edit-cancel">cancel</button>
					</div>
				</form>
			{:else}
				{@html linkify(message.content)}
				{#if attachment}
					<img alt="attachment" class="attachment" src={attachment} />
				{/if}
			{/if}
		</div>
	</div>
</div>

<UserProfileModal bind:showModal user={selectedUser} />

{#snippet message_menu()}
	<div class="message-menu">
		{#await get_user(message.user_id) then user}
			{#if serverId && (can('manage_messages', serverId) || message.user_id === me.user?.id)}
				<button
					class="menu-trigger"
					onclick={() => handleMenuClick(user)}
					aria-label="message actions"
				>
					⋮
				</button>
			{/if}
			{#if showMenu && selectedUser?.id === user.id}
				<div class="menu-dropdown">
					<button
						onclick={() => {
							showModal = true;
							showMenu = false;
						}}
						class="menu-item"
					>
						view profile
					</button>

					{#if serverId && (can('manage_messages', serverId) || message.user_id === me.user?.id)}
						{#if message.user_id === me.user?.id}
							<button onclick={() => startEditing()} class="menu-item"> edit message </button>
						{/if}
						<button onclick={() => handleDeleteMessage(message.id)} class="menu-item delete">
							delete message
						</button>
					{/if}
				</div>
			{/if}
		{/await}
	</div>
{/snippet}

<style>
	.t-slash:hover {
		transform: trasnalateY(10em);
	}
	.message {
		padding: 0.1em 0.8em;
		position: relative;
		border-radius: 3px;
		transition: transform 0.2s;
	}

	.message:hover {
		background-color: var(--base01);
		/* transform: translateX(-0.2em); */
	}

	.message-content {
		display: flex;
		flex-direction: column;
		gap: 0.2em;
	}

	.message-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		min-height: 1.5em;
	}

	.username {
		cursor: pointer;
		font-weight: bold;
		transition: scale 0.2s;
		font-family: var(--nc-font-mono);
		font-size: 0.9em;
		color: var(--base05);
	}

	.username:hover {
		transform: scale(1.2);
	}

	.message-menu {
		position: relative;
	}

	.menu-trigger {
		background: none;
		border: none;
		color: var(--base04);
		cursor: pointer;
		font-size: 1.2em;
		padding: 0.2em;
		border-radius: 3px;
		transition: background-color 0.2s;
		opacity: 0;
		transition: opacity 0.2s;
	}

	.message:hover .menu-trigger {
		opacity: 1;
	}

	.menu-trigger:hover {
		background-color: var(--base02);
	}

	.menu-dropdown {
		position: absolute;
		right: 0;
		top: 100%;
		background: var(--base00);
		border: 1px solid var(--base02);
		border-radius: 4px;
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
		min-width: 150px;
		z-index: 1000;
		padding: 0.5em 0;
	}

	.menu-item {
		display: block;
		width: 100%;
		padding: 0.5em 1em;
		background: none;
		border: none;
		color: var(--base05);
		cursor: pointer;
		text-align: left;
		transition: background-color 0.2s;
		font-size: 0.9em;
		font-family: var(--nc-font-mono);
		border-radius: 3px;
	}

	.menu-item:hover {
		background-color: var(--base01);
	}

	.menu-item.delete {
		color: var(--base08);
	}

	.menu-item.delete:hover {
		background-color: var(--base08);
		color: var(--base07);
	}

	.menu-divider {
		height: 1px;
		background-color: var(--base02);
		margin: 0.25em 0;
	}

	.menu-section {
		padding: 0;
	}

	.menu-label {
		padding: 0.25em 1em;
		font-size: 0.8em;
		color: var(--base04);
	}

	.message-body {
		color: var(--base05);
		line-height: 1.4;
		word-wrap: break-word;
	}

	.attachment {
		max-width: 300px;
		max-height: 200px;
		border-radius: 4px;
		margin-top: 0.5em;
		display: block;
	}

	.edit-form {
		display: flex;
		flex-direction: column;
		gap: 0.5em;
	}

	.edit-input {
		background: var(--base02);
		border: 1px solid var(--base03);
		color: var(--base05);
		padding: 0.5em;
		border-radius: 3px;
		font-family: var(--nc-font-mono);
		font-size: 0.9em;
		outline: none;
	}

	.edit-input:focus {
		border-color: var(--base0B);
	}

	.edit-actions {
		display: flex;
		gap: 0.5em;
		justify-content: flex-start;
	}

	.edit-save,
	.edit-cancel {
		border: none;
		cursor: pointer;
		font-family: var(--nc-font-mono);
		font-size: 0.8em;
		transition: background-color 0.2s;
		color: var(--base0B);
		background: var(--base00);
		border: 1px solid var(--base0B);
	}

	.edit-cancel:hover {
		transform: scale(1.1) rotate(10deg);
	}
	.edit-save:hover {
		transform: scale(1.3) rotate(-10deg);
	}
</style>
