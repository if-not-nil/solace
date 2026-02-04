import { BASE } from './api';
import { messages_add, messages_remove, channel_cache } from './caches.svelte';
import { logger } from './log';
import type { MessageHistoryResponse, MessageResponse } from './orm_types';
import { token as tok } from './user_store.svelte';

const log = logger('ws');

interface WSMessage {
	type: string;
	payload: any;
}

class WebSocketStore {
	private ws: WebSocket | null = null;
	private reconnectTimeout: ReturnType<typeof setTimeout> | null = null;
	private onMessageCallback: (() => void) | null = null;
	private onAuthFailure: (() => void) | null = null;
	private currentChannel: string | null = null;
	private _isConnected = false;
	private _isConnecting = false;

	private async _connect() {
		if (this._isConnecting || this._isConnected) return;

		const token = tok.get();
		if (!token) {
			log.warn('no token available for websocket connection');
			return;
		}

		this._isConnecting = true;
		const wsBase = BASE.replace('http://', 'ws://').replace('https://', 'wss://');

		try {
			this.ws = new WebSocket(`${wsBase}/ws/connect`);
			log.debug('initiating websocket connection');

			this.ws.onopen = () => {
				log.info('websocket connected');
				this._isConnected = true;
				this._isConnecting = false;

				// send auth immediately
				this.sendAuth();
			};

			this.ws.onmessage = (event) => {
				this.handleMessage(event);
			};

			this.ws.onclose = (event) => {
				log.debug(`websocket closed: ${event.code} ${event.reason || 'no reason'}`);
				this._isConnected = false;
				this._isConnecting = false;
				this.ws = null;

				// auto-reconnect for unexpected closures (not manual disconnects)
				if (event.code !== 1000 && tok.get()) {
					this.scheduleReconnect();
				}
			};

			this.ws.onerror = (error) => {
				log.err('websocket error:', error);
				this._isConnecting = false;
			};
		} catch (error) {
			log.err('failed to create websocket:', error);
			this._isConnecting = false;
		}
	}

	private scheduleReconnect() {
		if (this.reconnectTimeout) return;

		const delay = Math.min(1000 * Math.pow(2, Math.floor(Math.random() * 5)), 30000);
		log.info(`reconnecting in ${delay}ms`);

		this.reconnectTimeout = setTimeout(() => {
			this.reconnectTimeout = null;
			this._connect();
		}, delay);
	}

	private sendAuth() {
		if (!this.ws || this.ws.readyState !== WebSocket.OPEN) return;

		const message = {
			type: 'auth',
			payload: {
				token: tok.get(),
				channel_id: this.currentChannel || ''
			}
		};

		this.ws.send(JSON.stringify(message));
		log.debug('sent auth message');
	}

	private _sendChannelSwitch() {
		if (!this.ws || this.ws.readyState !== WebSocket.OPEN || !this.currentChannel) return;

		const message = {
			type: 'swap_channel',
			payload: { channel_id: this.currentChannel }
		};

		this.ws.send(JSON.stringify(message));
		log.debug(`switched to channel: ${this.currentChannel}`);
	}

	private handleMessage(event: MessageEvent) {
		try {
			const wsMsg: WSMessage = JSON.parse(event.data);

			switch (wsMsg.type) {
				case 'message':
					this.handleIncomingMessage(wsMsg.payload);
					break;
				case 'message_edited':
					this.handleMessageEdited(wsMsg.payload);
					break;
				case 'message_deleted':
					this.handleMessageDeleted(wsMsg.payload);
					break;
				case 'auth_failed':
					log.warn('websocket auth failed');
					if (this.onAuthFailure) {
						this.onAuthFailure();
					}
					break;
				case 'pong':
					// keep-alive response, ignore
					break;
				default:
					log.debug(`unknown message type: ${wsMsg.type}`);
			}
		} catch (err) {
			log.err('error parsing websocket message:', err);
		}
	}

	private handleIncomingMessage(msg: MessageResponse) {
		if (!this.currentChannel) {
			log.warn('received message but no current channel set');
			return;
		}

		const convertedMsg: MessageHistoryResponse = {
			id: msg.id,
			content: msg.content,
			user_id: msg.user_id,
			created_at: msg.created_at,
			updated_at: msg.created_at
		};

		// Remove any optimistic message that matches this real message
		const channelRecord = channel_cache[this.currentChannel];
		if (channelRecord?.msgs) {
			channelRecord.msgs = channelRecord.msgs.filter(
				(existingMsg) =>
					!(
						existingMsg.id.startsWith('optimistic-') &&
						existingMsg.content === msg.content &&
						existingMsg.user_id === msg.user_id
					)
			);
		}

		messages_add(this.currentChannel, [convertedMsg]);

		if (this.onMessageCallback) {
			this.onMessageCallback();
		}
	}

	private handleMessageEdited(payload: { message_id: string; new_content: string }) {
		if (!this.currentChannel) {
			log.warn('received message edit but no current channel set');
			return;
		}

		const channelRecord = channel_cache[this.currentChannel];
		if (channelRecord?.msgs) {
			channelRecord.msgs = channelRecord.msgs.map((msg) =>
				msg.id === payload.message_id
					? { ...msg, content: payload.new_content, updated_at: new Date().toISOString() }
					: msg
			);
		}

		if (this.onMessageCallback) {
			this.onMessageCallback();
		}
	}

	private handleMessageDeleted(payload: { message_id: string }) {
		if (!this.currentChannel) {
			log.warn('received message deletion but no current channel set');
			return;
		}

		messages_remove(this.currentChannel, payload.message_id);

		if (this.onMessageCallback) {
			this.onMessageCallback();
		}
	}

	get isConnected(): boolean {
		return this._isConnected;
	}

	get isConnecting(): boolean {
		return this._isConnecting;
	}

	get currentChannelId(): string | null {
		return this.currentChannel;
	}

	setChannel(channelId: string | null) {
		this.currentChannel = channelId;
	}

	setOnMessageCallback(callback: (() => void) | null) {
		this.onMessageCallback = callback;
	}

	setOnAuthFailure(callback: (() => void) | null) {
		this.onAuthFailure = callback;
	}

	// public methods for component control
	connect() {
		this._connect();
	}

	sendChannelSwitch() {
		this._sendChannelSwitch();
	}

	disconnect() {
		if (this.reconnectTimeout) {
			clearTimeout(this.reconnectTimeout);
			this.reconnectTimeout = null;
		}

		if (this.ws) {
			this.ws.close(1000, 'manual disconnect');
			this.ws = null;
		}

		this._isConnected = false;
		this._isConnecting = false;
		this.currentChannel = null;
	}
}

export const wsStore = new WebSocketStore();
