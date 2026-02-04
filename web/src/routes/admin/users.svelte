<script lang="ts">
	import { type Role } from '$lib/orm_types';
	import { api } from '$lib/api';
	import { alerts } from '$lib/alerts';
	import { can, Permissions, type Permission } from '$lib/permissions';
	import Modal from '$lib/components/Modal.svelte';
	let { server } = $props();
	let selectedUserId = $state('');
	const default_perms: Permission[] = [
		'attach_files',
		'send_messages',
		'view_channels',
		'embed_links'
	];
	$effect(() => {
		if (!server?.id) return;
		api.server.rolesList(server.id).then((res) => (roles = res));
	});
	function openRoleModal(role?: Role) {
		editingRoleId = role?.id ?? null;
		newRoleName = role?.name ?? '';
		Permissions.forEach((p) => {
			rolePermissions[p] = role?.permissions.includes(p) || default_perms.includes(p);
		});
		showRoleModal = true;
	}
	async function handleSaveRole() {
		if (!server?.id || !newRoleName.trim()) return;
		const selected = Permissions.filter((p) => rolePermissions[p]);

		try {
			if (editingRoleId) {
				await api.server.roleDelete(server.id, editingRoleId);
			}
			await api.server.roleCreate(server.id, newRoleName.trim(), selected);
			alerts.add(`role "${newRoleName}" saved`, 'success');
			showRoleModal = false;
			roles = await api.server.rolesList(server.id);
		} catch (e) {
			alerts.add(`failed: ${e}`, 'error');
		}
	}
	let roles = $state<Role[]>([]);

	let showRoleModal = $state(false);
	let editingRoleId = $state<string | null>(null);
	let newRoleName = $state('');
	let rolePermissions = $state<Record<Permission, boolean>>(
		Object.fromEntries(Permissions.map((p) => [p, false])) as Record<Permission, boolean>
	);
</script>

<Modal header={null} bind:showModal={showRoleModal}>
	<div class="modal-header">
		<h2>{editingRoleId ? 'edit' : 'create'} role</h2>
	</div>
	<div class="form-group">
		<label>name</label>
		<input class="input" bind:value={newRoleName} />
	</div>
	<h4>permissions</h4>
	<div class="permissions-grid">
		{#each Permissions as perm}
			<label class="perm-item">
				<input type="checkbox" bind:checked={rolePermissions[perm]} />
				<span>{perm.replace('_', ' ')}</span>
			</label>
		{/each}
	</div>
	<div class="modal-actions">
		<button class="btn btn-secondary" onclick={() => (showRoleModal = false)}>cancel</button>
		<button class="btn btn-primary" onclick={handleSaveRole}>save</button>
	</div>
</Modal>

{#if server?.id && can('kick_users', server.id)}
	<section class="section">
		<h4>kick user</h4>
		<div class="form-group">
			<input class="input" bind:value={selectedUserId} placeholder="user id" />
			<button
				class="btn btn-danger"
				onclick={async () => {
					await api.server.kick(server!.id, selectedUserId);
					selectedUserId = '';
				}}>kick</button
			>
		</div>
	</section>
{/if}

{#if server?.id && can('manage_roles', server.id)}
	<section class="section">
		<h4>roles</h4>
		<button class="btn" onclick={() => openRoleModal()}>create new role</button>
		<div class="list">
			{#each roles as role}
				<div class="list-item">
					<div class="list-item-info">
						<span class="list-item-title">{role.name}</span>
						<span class="list-item-subtitle">{role.permissions.join(', ')}</span>
					</div>
					<div class="list-item-actions">
						<button class="btn btn-primary" onclick={() => openRoleModal(role)}>edit</button>
						<button
							class="btn btn-danger"
							onclick={async () => {
								await api.server.roleDelete(server!.id, role.id);
								roles = await api.server.rolesList(server!.id);
							}}>delete</button
						>
					</div>
				</div>
			{/each}

			{#if roles.length < 1}
				<p class="info-text">no roles</p>
			{/if}
		</div>
	</section>
{/if}

<style>
	.permissions-grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 0.5rem;
		margin: 1rem 0;
	}
	.perm-item {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 0.9rem;
		cursor: pointer;
	}

	@media (max-width: 400px) {
		.permissions-grid {
			grid-template-columns: 1fr;
		}
	}
</style>
