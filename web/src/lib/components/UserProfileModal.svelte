<script lang="ts">
	import Modal from './Modal.svelte';
	import type { UserResponse } from '$lib/orm_types';

	let { user, showModal = $bindable(false) }: { user: UserResponse | null; showModal: boolean } =
		$props();
</script>

<Modal bind:showModal>
	<!-- {#snippet header()} -->
	<!-- 	{user?.name || 'user profile'} -->
	<!-- {/snippet} -->

	{#if user}
		<div class="user-profile-modal">
			<div class="user-avatar-large">
				<img src={user.avatar || '/default_avatar.png'} alt="user avatar" />
			</div>
			<div class="user-details-large">
				<h2>{user.name}</h2>
				<div class="user-meta-large">
					<div class="meta-item">
						<span class="meta-label">id</span>
						<span class="meta-value">{user.id}</span>
					</div>
					<div class="meta-item">
						<span class="meta-label">status</span>
						<span class="meta-value status-online">online</span>
					</div>
				</div>
			</div>
		</div>
	{:else}
		<div class="loading-profile">
			<i><p>loading user profile...</p></i>
		</div>
	{/if}
</Modal>

<style>
	.user-profile-modal {
		display: flex;
		flex-direction: row;
		gap: 2em;
		padding: 1em;
		min-width: 500px;
	}

	.user-avatar-large img {
		width: 120px;
		height: 120px;
		/* border-radius: 1em; */
		object-fit: cover;
		/* border: 1px solid var(--base02); */
	}

	.user-details-large h2 {
		margin: 0 0 1em 0;
		color: var(--base05);
		font-family: var(--nc-font-mono);
		font-size: 1.5em;
		font-weight: bold;
	}

	.user-meta-large {
		display: flex;
		flex-direction: column;
		gap: 0.5em;
	}

	.meta-item {
		display: flex;
		align-items: center;
		gap: 0.5em;
		padding: 0.5em;
		/* background: var(--base01); */
		/* border-radius: 3px; */
	}

	.meta-label {
		color: var(--base04);
		font-family: var(--nc-font-mono);
		font-size: 0.9em;
		min-width: 60px;
	}

	.meta-value {
		color: var(--base05);
		font-family: var(--nc-font-mono);
		font-size: 0.9em;
		font-weight: bold;
	}

	.status-online {
		color: var(--base0B) !important;
	}

	.loading-profile {
		text-align: center;
		color: var(--base04);
		padding: 3em;
		font-family: var(--nc-font-mono);
		font-size: 1.1em;
	}

	@media (max-width: 768px) {
		.user-profile-modal {
			min-width: auto;
			padding: 0.5em;
		}

		.user-avatar-large img {
			width: 100px;
			height: 100px;
		}

		.user-details-large h2 {
			font-size: 1.5em;
		}
	}
</style>
