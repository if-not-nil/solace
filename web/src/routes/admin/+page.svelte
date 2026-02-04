<script lang="ts">
	import { type ServerResponse } from '$lib/orm_types';
	import { api } from '$lib/api';
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { logger } from '$lib/log';
	import Channels from './channels.svelte';
	import Invites from './invites.svelte';
	import Users from './users.svelte';

	let log = logger('admin page');

	let loading = $state(true);
	let server = $state<ServerResponse>();
	let activeTab = $state('channels');

	onMount(async () => {
		const id = page.url.searchParams.get('id');
		if (!id) return goto('/');
		try {
			server = await api.server.get(id);
		} catch (err) {
			log.err(`Couldnt find server ${id}: ${err}`);
			goto('/');
		} finally {
			loading = false;
		}
	});
</script>

{#if !loading}
	<div class="page">
		<header class="page-header">
			<h1 class="page-title">server admin: {server?.name}</h1>
		</header>

		<nav class="tabs">
			{#each ['channels', 'invites', 'users'] as tab}
				<button class="tab" class:active={activeTab === tab} onclick={() => (activeTab = tab)}>
					{tab}
				</button>
			{/each}
		</nav>

		<main class="content">
			{#if activeTab === 'channels'}
				<Channels {server} />
			{:else if activeTab == 'invites'}
				<Invites {server} />
			{:else if activeTab == 'users'}
				<Users {server} />
			{/if}
		</main>
	</div>
{/if}
