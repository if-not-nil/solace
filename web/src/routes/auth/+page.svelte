<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { alerts } from '$lib/alerts';
	import { api } from '$lib/api';
	import { logger } from '$lib/log';
	import { me, token } from '$lib/user_store.svelte';

	let signing_up = $state(false);
	let name = $state('');
	let password = $state('');
	let password_confirm = $state('');
	let signing_in_progress = $state(false);
	let log = logger('sign up page');

	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	function handle_error(what: any) {
		alerts.add(what, 'error');
	}

	async function handle_click() {
		const trimmedName = name.trim();
		const trimmedPassword = password.trim();

		if (trimmedName.length < 3 || trimmedName.length > 50) {
			handle_error('name must be between 3 and 50 characters');
			return;
		}
		if (trimmedPassword.length < 8 || trimmedPassword.length > 128) {
			handle_error('password must be between 8 and 128 characters');
			return;
		}
		if (signing_up && trimmedPassword !== password_confirm.trim()) {
			handle_error('passwords do not match');
			return;
		}

		if (signing_in_progress) return;
		signing_in_progress = true;

		try {
			if (signing_up) {
				log.info('signing up');
				const res = await api.auth.register(trimmedName, trimmedPassword);
				token.set(res.token);
				const self = await api.user.me();
				me.isLoggedIn = true;
				me.user = self.user;
				me.servers = self.servers;
				goto(resolve('/'));
			} else {
				const res = await api.auth.login(trimmedName, trimmedPassword);
				token.set(res.token);
				const self = await api.user.me();
				me.isLoggedIn = true;
				me.user = self.user;
				me.servers = self.servers;
				goto(resolve('/'));
			}
		} catch (err) {
			alerts.add(err, 'error');
			log.err(err);
		} finally {
			signing_in_progress = false;
		}
	}

	function handle_form_submit(e: SubmitEvent) {
		e.preventDefault();
		handle_click();
	}
</script>

<div class="main-box">
	<div class="subbox left">
		<h1>solace kinda</h1>
		<form class="input-fields" onsubmit={handle_form_submit}>
			<input
				type="text"
				bind:value={name}
				placeholder="username or email"
				disabled={signing_in_progress}
				required
			/>
			<input
				type="password"
				bind:value={password}
				placeholder="password"
				disabled={signing_in_progress}
				required
			/>
			{#if signing_up}
				<input
					type="password"
					bind:value={password_confirm}
					placeholder="password again"
					disabled={signing_in_progress}
					required
				/>
			{/if}
			<div class="reglogbuttons">
				<button
					type="button"
					onclick={() => {
						signing_up = !signing_up;
						password_confirm = '';
					}}
					disabled={signing_in_progress}
				>
					{signing_up ? 'log in' : 'sign up'} instead
				</button>
				<span class="separator">or</span>
				<button type="submit" disabled={signing_in_progress}>
					{signing_in_progress ? 'please wait...' : signing_up ? 'sign up' : 'log in'}
				</button>
			</div>
		</form>
	</div>
	<div class="subbox right"></div>
</div>

<style>
	.reglogbuttons {
		display: flex;
		justify-content: space-between;
		align-items: center;
		background: var(--base01);
		margin-top: auto;
		gap: 1em;
	}
	.reglogbuttons > button {
		width: 50%;
	}
	.separator {
		font-weight: bold;
	}

	input {
		padding: 0.5em;
		border-radius: 4px;
		border: 1px solid var(--base02);
		transition: border-color 0.2s;
	}

	input:focus {
		outline: none;
		border-color: var(--accent);
	}

	input:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	button {
		padding: 0.5em 1em;
		border-radius: 4px;
		border: none;
		cursor: pointer;
		transition: opacity 0.2s;
	}

	button:hover:not(:disabled) {
		opacity: 0.9;
	}

	button:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.subbox {
		flex-shrink: 1;
		flex-basis: 0;
	}
	.input-fields {
		display: flex;
		flex-direction: column;
		background: var(--base01);
		gap: 1em;
		width: 20em;
		padding: 2em 1em 2em 1em;
		border-radius: 8px;
	}
	.left {
		flex: 1.2;
		display: flex;
		justify-content: center;
		align-items: center;
		flex-direction: column;
	}

	@media (max-width: 768px) {
		.right {
			display: none;
		}
	}

	.right {
		flex: 0.8;
		background-image: url('/flight.jpg');
		background-size: cover;
		object-fit: fill;
		background-size: 100% 100%;
		background-position: center;
	}

	.main-box {
		margin: 0px;
		padding: 0px;
		display: flex;
		width: 100vw;
		height: 100vh;
		flex-direction: row;
	}
</style>
