import type { ServerResponse } from './orm_types';
import { get_channel, get_server, type ChannelRecord } from './caches.svelte';
import { fetchUserPermissions } from './permissions';

export const location = $state<{ server?: ServerResponse; channel?: ChannelRecord }>({
	server: undefined,
	channel: undefined
});

export const location_methods = {
	set_both_id: async (server_id: string, channel_id: string) => {
		location.server = await get_server(server_id);
		location.channel = await get_channel(channel_id, server_id);
		// preload permissions for the server
		fetchUserPermissions(server_id);
	},

	set_server_id: async (server_id: string) => {
		const server = await get_server(server_id);
		// pick first channel of server if exists
		const channel = server.channels.length
			? await get_channel(server.channels[0].id, server_id)
			: undefined;
		location.server = server;
		location.channel = channel;
		// preload permissions for the server
		fetchUserPermissions(server_id);
	},

	set_channel_id: async (channel_id: string) => {
		const channel = await get_channel(channel_id);
		const server = await get_server(channel.meta.server_id);
		location.server = server;
		location.channel = channel;
		// preload permissions for the server
		fetchUserPermissions(channel.meta.server_id);
	}
};
