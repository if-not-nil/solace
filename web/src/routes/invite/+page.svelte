<script lang="ts">
	import { api } from '$lib/api';
	import { alerts } from '$lib/alerts';
	import { goto } from '$app/navigation';
	import { try_update_auth } from '$lib/user_store.svelte';

	let { url }: { url: URL } = $props();
	let inviteCode = $state(url.searchParams.get('code'));
	let joining = $state(false);
	let error = $state('');

	async function handleJoinInvite() {
		if (!inviteCode) {
			error = 'no invite code provided';
			return;
		}

		joining = true;
		error = '';

		try {
			await api.server.linksJoin(inviteCode);
			alerts.add('joined server successfully', 'success');
			await try_update_auth();
			goto('/');
		} catch (e) {
			error = `failed to join: ${e}`;
		} finally {
			joining = false;
		}
	}

	if (inviteCode) {
		handleJoinInvite();
	}
</script>

<div class="page">
	<div class="page-header">
		<h1 class="page-title">join server</h1>
	</div>

	<div class="section">
		{#if joining}
			<p class="info-text">joining server...</p>
		{:else if error}
			<p class="error-text">{error}</p>
		{:else if !inviteCode}
			<div class="form-group">
				<input
					class="input"
					bind:value={inviteCode}
					placeholder="invite code"
					onkeydown={(e) => e.key === 'Enter' && handleJoinInvite()}
				/>
				<button class="btn btn-primary" onclick={handleJoinInvite} disabled={!inviteCode?.trim()}>
					join
				</button>
			</div>
		{:else}
			<p class="info-text">joined successfully! redirecting...</p>
		{/if}
	</div>
</div>
