<script lang="ts">
	import { try_update_auth } from '$lib/user_store.svelte';
	import { api } from '$lib/api';
	import { alerts } from '$lib/alerts';
	import Modal from './Modal.svelte';
	let { show = $bindable() } = $props();
	let activeTab = $state<'join' | 'create'>('join');

	// handlers
	async function handleCreateServer() {
		if (!newServerName.trim()) return;
		isCreating = true;
		try {
			const server = await api.server.create(newServerName.trim());
			alerts.add(`server "${server.name}" created`, 'success');
			show = false;
			newServerName = '';
			await try_update_auth();
		} catch (e) {
			alerts.add(`failed: ${e}`, 'error');
		} finally {
			isCreating = false;
		}
	}

	async function handleJoinServer() {
		if (!inviteLink.trim()) return;
		isCreating = true;
		try {
			await api.server.linksJoin(inviteLink.trim());
			alerts.add('joined server', 'success');
			show = false;
			inviteLink = '';
			await try_update_auth();
		} catch (e) {
			alerts.add(`failed: ${e}`, 'error');
		} finally {
			isCreating = false;
		}
	}
	let newServerName = $state('');
	let inviteLink = $state('');
	let isCreating = $state(false);
</script>

<Modal bind:showModal={show}>
	{#snippet header()}
		[<span
			class="tab-button"
			class:active={activeTab === 'join'}
			onclick={() => (activeTab = 'join')}>join</span
		>/<span
			class="tab-button"
			class:active={activeTab === 'create'}
			onclick={() => (activeTab = 'create')}>create</span
		>] server
	{/snippet}

	{#if activeTab === 'join'}
		<div class="section">
			<h3>join server</h3>
			<div class="form-group">
				<input
					class="input"
					bind:value={inviteLink}
					placeholder="invite link or code"
					onkeydown={(e) => e.key === 'Enter' && handleJoinServer()}
					disabled={isCreating}
				/>
				<button
					class="btn btn-primary"
					onclick={handleJoinServer}
					disabled={!inviteLink.trim() || isCreating}
				>
					{isCreating ? 'joining...' : 'join'}
				</button>
			</div>
		</div>
	{:else if activeTab === 'create'}
		<div class="section">
			<h3>create server</h3>
			<div class="form-group">
				<input
					class="input"
					bind:value={newServerName}
					placeholder="server name"
					onkeydown={(e) => e.key === 'Enter' && handleCreateServer()}
					disabled={isCreating}
				/>
				<button
					class="btn btn-primary"
					onclick={handleCreateServer}
					disabled={!newServerName.trim() || isCreating}
				>
					{isCreating ? 'creating...' : 'create'}
				</button>
			</div>
		</div>
	{/if}
</Modal>

<style>
	@media (max-width: 768px) {
		.form-group {
			flex-direction: column;
			align-items: stretch;
		}
	}
	.tab-button {
		color: var(--base03);
		text-decoration: underline;
	}
	.tab-button.active {
		color: var(--base07);
	}
</style>
