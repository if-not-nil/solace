<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
	import { me, try_update_auth } from '$lib/user_store.svelte';
	import './app.css';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import Alerts from '$lib/components/Alerts.svelte';

	onMount(() => {
		const auth_path = resolve('/auth');
		if (!me.isLoggedIn) {
			try_update_auth().then((e) => {
				if (!e) {
					console.log('off u go');
					goto(auth_path);
				}
			});
		} else goto(auth_path);
	});

	let { children } = $props();
</script>

<Alerts />
<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

{@render children()}
