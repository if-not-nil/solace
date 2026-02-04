import { SvelteDate, SvelteSet } from 'svelte/reactivity';
import { api } from './api';
import { logger } from './log';
const log = logger('caches');
import type {
	ChannelResponse,
	MessageHistoryResponse,
	ServerResponse,
	UserResponse
} from './orm_types';

export const server_cache = $state<Record<string, ServerResponse>>({});

export const get_server = async (id: string): Promise<ServerResponse> => {
	if (!id) {
		log.err('get_server: invalid server id');
		throw new Error('server id is required');
	}
	if (server_cache[id]) return server_cache[id];

	try {
		const server: ServerResponse = await api.server.get(id);
		server_cache[id] = server;
		server.channels?.forEach((channel) => {
			channel_cache[channel.id] = { meta: channel };
		});
		return server;
	} catch (err) {
		log.err(`get_server: failed to fetch server ${id}:`, err);
		throw new Error(`failed to load server: ${err}`);
	}
};

export type ChannelRecord = { meta: ChannelResponse; msgs?: MessageHistoryResponse[] };
export const channel_cache = $state<Record<string, ChannelRecord>>({});

export const get_channel = async (id: string, server_id?: string): Promise<ChannelRecord> => {
	if (!id) {
		log.err('get_channel: invalid channel id');
		throw new Error('channel id is required');
	}
	if (channel_cache[id]) return channel_cache[id];

	if (!server_id) {
		log.err(`get_channel: channel not found and no server_id provided for ${id}`);
		throw new Error('channel not found');
	}

	try {
		await get_server(server_id);
		const channel = channel_cache[id];
		if (!channel) {
			log.err(`get_channel: channel ${id} not found in server ${server_id}`);
			throw new Error('channel not found in server');
		}
		return channel;
	} catch (err) {
		log.err(`get_channel: failed to load channel ${id}:`, err);
		throw err;
	}
};

export async function message_history(
	channel_id: string,
	until: number,
	count: number = 30
): Promise<MessageHistoryResponse[]> {
	if (!channel_id) {
		log.err('message_history: invalid channel id');
		throw new Error('channel id is required');
	}
	if (count < 1) {
		log.err('message_history: invalid count');
		throw new Error('count must be at least 1');
	}

	try {
		const res = await api.channel.history(channel_id, until, count);

		if (!res.users || !Array.isArray(res.users)) {
			log.err('message_history: invalid users response');
			res.users = [];
		}
		if (!res.messages || !Array.isArray(res.messages)) {
			log.err('message_history: invalid messages response');
			res.messages = [];
		}

		res.users.forEach((u) => {
			if (u?.id) {
				user_cache[u.id] = u;
			}
		});
		messages_add(channel_id, res.messages);
		return res.messages;
	} catch (err) {
		log.err(`message_history: failed to fetch history for ${channel_id}:`, err);
		throw new Error(`failed to load message history: ${err}`);
	}
}

export function messages_add(channel_id: string, new_messages: MessageHistoryResponse[]) {
	if (!channel_id) {
		log.err('messages_add: invalid channel id');
		return;
	}
	if (!Array.isArray(new_messages)) {
		log.err('messages_add: new_messages is not an array');
		return;
	}

	const channelRecord = channel_cache[channel_id];
	if (!channelRecord) {
		log.err(`messages_add: channel ${channel_id} not found in cache`);
		return;
	}

	const existing = channelRecord.msgs || [];
	const existingIds = new SvelteSet(existing.map((m) => m.id));

	const filtered = new_messages.filter((m) => m?.id && !existingIds.has(m.id));
	const sorted = [...existing, ...filtered].sort((a, b) => {
		try {
			return new SvelteDate(a.created_at).getTime() - new SvelteDate(b.created_at).getTime();
		} catch (err) {
			log.err('messages_add: failed to parse dates:', err);
			return 0;
		}
	});
	channelRecord.msgs = sorted;
}

export function messages_remove(channel_id: string, message_id: string) {
	if (!channel_id) {
		log.err('messages_remove: invalid channel id');
		return;
	}
	if (!message_id) {
		log.err('messages_remove: invalid message id');
		return;
	}

	const channelRecord = channel_cache[channel_id];
	if (!channelRecord) {
		log.err(`messages_remove: channel ${channel_id} not found in cache`);
		return;
	}

	if (!channelRecord.msgs) {
		log.warn(`messages_remove: no messages in channel ${channel_id}`);
		return;
	}

	const filtered = channelRecord.msgs.filter((m) => m.id !== message_id);
	channelRecord.msgs = filtered;
}

export const user_cache = $state<Record<string, UserResponse>>({});

export const get_user = async (id: string): Promise<UserResponse> => {
	if (!id) {
		log.err('get_user: invalid user id');
		return Promise.reject(new Error('user id is required'));
	}

	if (user_cache[id]) return Promise.resolve(user_cache[id]);

	try {
		const user = await api.user.by_id(id);
		if (!user?.id) {
			log.err('get_user: invalid user response');
			throw new Error('invalid user response');
		}
		user_cache[id] = user;
		return user;
	} catch (err) {
		log.err(`get_user: failed to fetch user ${id}:`, err);
		throw new Error(`failed to load user: ${err}`);
	}
};
