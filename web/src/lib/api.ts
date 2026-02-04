import type { AxiosPromise } from 'axios';
import type {
	UserResponse,
	ServerResponse,
	ChannelResponse,
	InviteLink,
	ChannelHistoryResponse
} from './orm_types';
import axios from 'axios';
import { token } from './user_store.svelte';
import type { Permission } from './permissions';

export const BASE = import.meta.env.VITE_API_URL || 'http://localhost:1323/api';
const axe = axios.create({
	baseURL: BASE,
	headers: {
		'Content-Type': 'application/json'
	}
});

axe.interceptors.request.use(
	(config) => {
		const t = token.get();
		if (t) {
			config.headers.Authorization = `Bearer ${t}`;
		}
		return config;
	},
	(error) => {
		return Promise.reject(error);
	}
);

export async function safeRequest<T>(request: AxiosPromise<T>): Promise<T> {
	try {
		const response = await request;
		return response.data;
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
	} catch (error: any) {
		console.log(error);
		if (error.response?.data?.message) throw error.response.data.message;
		throw error.message || 'unknown error';
	}
}

export interface LoginResponse {
	message: string;
	token: string;
	user: UserResponse;
}

export interface RegisterResponse {
	message: string;
	token: string;
	user: UserResponse;
}

export const api = {
	auth: {
		login: (name: string, password: string): Promise<LoginResponse> =>
			safeRequest(axe.post('/auth/login', { name, password })),

		register: (name: string, password: string): Promise<RegisterResponse> =>
			safeRequest(axe.post('/auth/register', { name, password })),

		refresh: (): Promise<{ token: string }> => safeRequest(axe.get('/auth/refresh'))
	},

	user: {
		me: async (): Promise<{ user: UserResponse; servers: ServerResponse[] }> =>
			safeRequest(axe.get('/user/me')),
		by_id: async (id: string): Promise<UserResponse> => (await api.user.me()).user // temp
		// safeRequest(axe.get("/user/me")),
	},

	server: {
		get: (id: string): Promise<ServerResponse> => safeRequest(axe.get(`/server/${id}`)),

		create: (name: string): Promise<ServerResponse> => safeRequest(axe.post('/server', { name })),

		kick: (id: string, userId: string): Promise<{ success: boolean }> =>
			safeRequest(axe.post(`/server/${id}/kick`, { id: userId })),

		newChannel: (id: string, name: string): Promise<ChannelResponse> =>
			safeRequest(axe.post(`/server/${id}/new_channel`, { name })),

		linksNew: (id: string, maxUsers: number): Promise<InviteLink> =>
			safeRequest(axe.post(`/server/${id}/links`, { max_users: maxUsers })),

		linksList: (id: string): Promise<InviteLink[]> => safeRequest(axe.get(`/server/${id}/links`)),

		linksRevoke: (id: string, inviteId: string): Promise<{ message: string }> =>
			safeRequest(axe.delete(`/server/${id}/links/${inviteId}`)),

		linksJoin: (inviteId: string): Promise<string> =>
			safeRequest(axe.post(`/server/join/${inviteId}`)),

		setAvatar: (serverId: string, fileId: string): Promise<{ id: string }> =>
			safeRequest(axe.post(`/server/${serverId}/avatar?id=${encodeURIComponent(fileId)}`)),

		assignRole: (serverId: string, userId: string, roleId: string): Promise<{ message: string }> =>
			safeRequest(axe.put(`/server/${serverId}/roles/assign/${userId}/${roleId}`)),

		rolesList: (serverId: string): Promise<import('./orm_types').Role[]> =>
			safeRequest(axe.get(`/server/${serverId}/roles`)),

		userRoles: (serverId: string): Promise<import('./orm_types').Role[]> =>
			safeRequest(axe.get(`/server/${serverId}/user/roles`)),

		roleCreate: (
			serverId: string,
			name: string,
			permissions: Permission[]
		): Promise<import('./orm_types').Role> =>
			safeRequest(axe.post(`/server/${serverId}/roles`, { name, permissions })),

		roleUpdate: (
			serverId: string,
			roleId: string,
			permissions: Record<Permission, boolean>
		): Promise<import('./orm_types').Role> =>
			safeRequest(axe.put(`/server/${serverId}/roles/${roleId}`, permissions)),

		roleDelete: (serverId: string, roleId: string): Promise<{ message: string }> =>
			safeRequest(axe.delete(`/server/${serverId}/roles/${roleId}`))
	},

	channel: {
		send: (id: string, content: string): Promise<{ message: string }> =>
			safeRequest(axe.post(`/channel/${id}/send`, { content })),

		history: (id: string, until?: number, count = 20): Promise<ChannelHistoryResponse> =>
			safeRequest(axe.post(`/channel/${id}/history`, { until, count })),

		delete: (id: string): Promise<{ success: boolean }> => safeRequest(axe.delete(`/channel/${id}`))
	},

	message: {
		delete: (id: string): Promise<{ message: string }> => safeRequest(axe.delete(`/message/${id}`)),
		edit: (id: string, content: string) => safeRequest(axe.patch(`/message/${id}`, { content }))
	},

	file: {
		get: (id: string): Promise<Blob> =>
			safeRequest(axe.get(`/file/${id}`, { responseType: 'blob' })),

		upload: (formData: FormData): Promise<{ id: string; type: string }> =>
			safeRequest(
				axe.post('/file/upload', formData, {
					headers: { 'Content-Type': 'multipart/form-data' }
				})
			)
	},

	ws: {
		connect: (url: string) => new WebSocket(url)
	}
};
