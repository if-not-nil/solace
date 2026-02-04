<script lang="ts">
	import { type InviteLink } from '$lib/orm_types';
	import { api } from '$lib/api';
	import { can } from '$lib/permissions';
	let { server } = $props();
	let maxUsers = $state(10);
	let inviteLinks = $state<InviteLink[]>([]);
	$effect(() => {
		if (!server?.id) return;
		api.server.linksList(server.id).then((res) => (inviteLinks = res));
	});
</script>

{#if server?.id && can('manage_invites', server.id)}
	<div class="section">
		<h3>invite links</h3>
		<div class="form-group">
			<input class="input" type="number" bind:value={maxUsers} min="1" />
			<button
				class="btn btn-primary"
				onclick={async () => {
					await api.server.linksNew(server!.id, maxUsers);
					inviteLinks = await api.server.linksList(server!.id);
				}}>create</button
			>
		</div>
		<div class="list">
			{#if inviteLinks.length < 1}
				<p class="info-text">no invite links</p>
			{/if}
			{#each inviteLinks as invite}
				<div class="list-item">
					<div class="list-item-info">
						<span class="list-item-title">{invite.id}</span>
						<span class="list-item-subtitle">left: {invite.joins_left}</span>
					</div>
					<button
						class="btn btn-danger"
						onclick={async () => {
							await api.server.linksRevoke(server!.id, invite.id);
							inviteLinks = await api.server.linksList(server!.id);
						}}>revoke</button
					>
				</div>
			{/each}
		</div>
	</div>
{/if}
