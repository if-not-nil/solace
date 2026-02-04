import { api } from './api';
import { logger } from './log';
import { alerts } from './alerts';
import { wsStore } from './socket.svelte';
import { message_history, messages_add, channel_cache } from './caches.svelte';
import { me } from './user_store.svelte';
import { location } from './navigation.svelte';
import type { MessageHistoryResponse } from './orm_types';

const log = logger('messaging');

class MessagingStore {
	state = $state({
		inputValue: '',
		isLoading: false,
		hasError: false,
		isInputFocused: true,
		scrollTimeout: null as ReturnType<typeof setTimeout> | null
	});

	sendMessage = async (content: string, channelId: string) => {
		if (!content?.trim() || !channelId) return;

		if (content.length > 2000) {
			alerts.add('message too long (max 2000 characters)', 'error');
			return;
		}

		const optimisticMessage: MessageHistoryResponse = {
			id: `optimistic-${Date.now()}-${Math.random()}`,
			content: content.trim(),
			user_id: me.user?.id || '',
			created_at: new Date().toISOString(),
			updated_at: new Date().toISOString()
		};

		messages_add(channelId, [optimisticMessage]);

		try {
			await api.channel.send(channelId, content.trim());
		} catch (e) {
			// remove optimistic message on failure
			this.removeOptimisticMessage(channelId, optimisticMessage.id);
			alerts.add(`failed to send message: ${e}`, 'error');
			log.err('failed to send message:', e);
		}
	};

	editMessage = async (messageId: string, newContent: string, channelId: string) => {
		if (!newContent?.trim() || !messageId || !channelId) return;

		if (newContent.length > 2000) {
			alerts.add('message too long (max 2000 characters)', 'error');
			return;
		}

		const channelRecord = channel_cache[channelId];
		if (!channelRecord?.msgs) return;

		const originalMessage = channelRecord.msgs.find((msg) => msg.id === messageId);
		if (!originalMessage) return;

		const originalContent = originalMessage.content;

		// optimistic update
		const updatedMessage = {
			...originalMessage,
			content: newContent.trim(),
			updated_at: new Date().toISOString()
		};

		const updatedMessages = channelRecord.msgs.map((msg) =>
			msg.id === messageId ? updatedMessage : msg
		);

		channelRecord.msgs = updatedMessages;

		try {
			await api.message.edit(messageId, newContent.trim());
		} catch (e) {
			// revert optimistic update on failure
			const revertedMessages = channelRecord.msgs.map((msg) =>
				msg.id === messageId
					? { ...msg, content: originalContent, updated_at: originalMessage.updated_at }
					: msg
			);
			channelRecord.msgs = revertedMessages;
			alerts.add(`failed to edit message: ${e}`, 'error');
			log.err('failed to edit message:', e);
		}
	};

	removeOptimisticMessage = (channelId: string, messageId: string) => {
		const channelRecord = channel_cache[channelId];
		if (channelRecord?.msgs) {
			channelRecord.msgs = channelRecord.msgs.filter((msg) => msg.id !== messageId);
		}
	};

	loadInitialMessages = async (channelId: string) => {
		try {
			this.state.isLoading = true;
			await message_history(channelId, Date.now(), 20);
		} catch (error) {
			log.err('failed to load channel messages:', error);
			this.state.hasError = true;
		} finally {
			this.state.isLoading = false;
		}
	};

	// for infinite scroll
	loadMoreMessages = async (channelId: string) => {
		try {
			const currentMessages = channel_cache[channelId]?.msgs || [];
			if (currentMessages.length === 0) return;

			const oldestMessage = currentMessages[0];
			const olderMessages = await message_history(
				channelId,
				new Date(oldestMessage.created_at).getTime(),
				20
			);

			// prepend older messages
			if (olderMessages.length > 0) {
				const existingMessages = channel_cache[channelId]?.msgs || [];
				const newMessages = olderMessages.filter(
					(msg) => !existingMessages.some((existing) => existing.id === msg.id)
				);

				if (newMessages.length > 0) {
					channel_cache[channelId].msgs = [...newMessages, ...existingMessages];
				}
			}
		} catch (error) {
			log.err('failed to load more messages:', error);
		}
	};

	// for infinite loading
	handleScroll = (msgListElement: HTMLElement | null) => {
		if (!msgListElement || !location.channel) return;

		const { scrollTop } = msgListElement;
		if (scrollTop < 100) {
			this.loadMoreMessages(location.channel.meta.id);
		}
	};

	debouncedScrollHandler = (msgListElement: HTMLElement | null) => {
		if (this.state.scrollTimeout) {
			clearTimeout(this.state.scrollTimeout);
		}
		this.state.scrollTimeout = setTimeout(() => {
			this.handleScroll(msgListElement);
		}, 100);
	};

	scrollToBottom = (msgListElement: HTMLElement | null, behavior: ScrollBehavior = 'smooth') => {
		if (msgListElement) {
			requestAnimationFrame(() => {
				msgListElement.scroll({
					top: msgListElement.scrollHeight,
					behavior
				});
			});
		}
	};

	initializeForChannel = (channelId: string | null) => {
		if (channelId !== wsStore.currentChannelId) {
			wsStore.setChannel(channelId);
			if (channelId) {
				this.state.hasError = false;
				this.loadInitialMessages(channelId);
			}
		}
	};

	setup = () => {
		wsStore.setOnMessageCallback(() => {
			this.scrollToBottom(document.querySelector('.msg-list') as HTMLElement);
		});

		wsStore.setOnAuthFailure(() => {
			window.location.href = '/auth';
		});
	};

	cleanup = () => {
		wsStore.disconnect();
		if (this.state.scrollTimeout) {
			clearTimeout(this.state.scrollTimeout);
		}
	};

	autoConnect = () => {
		const hasToken = !!me.isLoggedIn;
		if (hasToken && !wsStore.isConnected && !wsStore.isConnecting) {
			wsStore.connect();
		} else if (!hasToken && wsStore.isConnected) {
			wsStore.disconnect();
		}
	};

	autoSwitchChannel = () => {
		if (wsStore.isConnected && location.channel) {
			wsStore.sendChannelSwitch();
		}
	};
}

export const messaging = new MessagingStore();
