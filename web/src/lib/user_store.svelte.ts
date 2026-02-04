import type { ServerResponse, UserResponse } from '$lib/orm_types';
import { api } from './api';

export const me = $state({
	servers: undefined as ServerResponse[] | undefined,
	user: undefined as UserResponse | undefined,
	isLoggedIn: false
});

export const logout = () => {
	localStorage.removeItem('token');
	me.servers = undefined;
	me.user = undefined;
	me.isLoggedIn = false;
};

export const login = () => {
	me.servers = undefined;
	me.user = undefined;
	me.isLoggedIn = false;
};

export const try_update_auth = async (): Promise<boolean> => {
	try {
		await api.auth.refresh();
		const res = await api.user.me();
		me.isLoggedIn = true;
		me.user = res.user;
		me.servers = res.servers;
		return true;
	} catch (e) {
		console.error(e);
		return false;
	}
};

export const token = {
	set: (val: string) => {
		localStorage.setItem('token', val);
	},
	get: function (): string | null {
		return localStorage.getItem('token');
	}
};
