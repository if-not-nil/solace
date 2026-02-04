import { writable } from 'svelte/store';

type Alert = {
	message: string;
	type?: 'error' | 'success' | 'info' | 'warning';
	id: number;
};

function createAlertStore() {
	const { subscribe, update } = writable<Alert[]>([]);
	let idCounter = 0;

	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	function add(message: any, type: Alert['type'] = 'error') {
		const id = idCounter++;
		update((alerts) => [...alerts, { message, type, id }]);
		// Auto-remove after 5 seconds
		setTimeout(() => {
			remove(id);
		}, 5000);
	}

	function remove(id: number) {
		update((alerts) => alerts.filter((alert) => alert.id !== id));
	}

	return { subscribe, add, remove };
}

export const alerts = createAlertStore();
