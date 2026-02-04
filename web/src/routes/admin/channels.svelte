<script lang="ts">
	import { alerts } from '$lib/alerts';
	import { api } from '$lib/api';
	import { can } from '$lib/permissions';

	let newChannelName = $state('');
	let showConfirmDelete = $state<string | null>(null);
	let { server } = $props();
	async function refreshServer() {
		if (server?.id) server = await api.server.get(server.id);
	}
	const wrapAction = (fn: () => Promise<any>, msg: string) => async () => {
		try {
			await fn();
			alerts.add(msg, 'success');
		} catch (e) {
			alerts.add(`error: ${e}`, 'error');
		}
	};

	const handleCreateChannel = wrapAction(async () => {
		await api.server.newChannel(server!.id, newChannelName.trim());
		newChannelName = '';
		await refreshServer();
	}, 'channel created');

	const handleDeleteChannel = async (id: string) => {
		await api.channel.delete(id);
		showConfirmDelete = null;
		await refreshServer();
	};
</script>

{#if server?.id && can('manage_channels', server.id)}
	<section class="section">
		<h3>channel management</h3>
		<div class="form-group">
			<input
				class="input"
				bind:value={newChannelName}
				placeholder="channel name"
				onkeydown={(e) => e.key === 'Enter' && handleCreateChannel()}
			/>
			<button
				class="btn btn-primary"
				onclick={handleCreateChannel}
				disabled={!newChannelName.trim()}>create</button
			>
		</div>

		<h4>existing channels</h4>
		<div class="list">
			{#each server.channels as channel (channel.id)}
				<div class="list-item">
					<span># {channel.name}</span>
					{#if showConfirmDelete === channel.id}
						<div class="list-item-actions">
							<button class="btn btn-danger" onclick={() => handleDeleteChannel(channel.id)}
								>confirm</button
							>
							<button class="btn btn-secondary" onclick={() => (showConfirmDelete = null)}
								>cancel</button
							>
						</div>
					{:else}
						<button
							class="btn btn-danger"
							onclick={() => (showConfirmDelete = channel.id)}
							disabled={server.channels.length <= 1}>delete</button
						>
					{/if}
				</div>
			{/each}
			{#if server.channels.length < 1}
				<p class="info-text">no channels</p>
			{/if}
		</div>
	</section>
{:else}
	<p class="error-text">insufficient permissions</p>
{/if}
