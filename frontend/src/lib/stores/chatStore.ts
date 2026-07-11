import { writable } from 'svelte/store';

export interface ChatMessage {
	id?: number;
	sender_id: number;
	receiver_id: number;
	content: string;
	created_at?: string;
	is_read?: boolean;
	is_deleted?: boolean;
}

function createChatStore() {
	const { subscribe, set, update } = writable<{
		messages: ChatMessage[];
		connected: boolean;
		activeUserId: number | null;
		unreadCount: number;
		myUserId: number | null;
	}>({
		messages: [],
		connected: false,
		activeUserId: null,
		unreadCount: 0,
		myUserId: null
	});

	let socket: WebSocket | null = null;
	let pingInterval: ReturnType<typeof setInterval> | null = null;

	return {
		subscribe,
		setMyUserId: (id: number) => update(s => ({ ...s, myUserId: id })),
		setUnreadCount: (count: number) => update(s => ({ ...s, unreadCount: count })),
		connect: () => {
			const wsProtocol = window.location.protocol === 'https:' ? 'wss://' : 'ws://';
			// Now that vite.config.ts proxies /ws to 8080, we can safely use the current host for all environments!
			const host = window.location.host;
			const wsUrl = wsProtocol + host + '/ws/chat';
			
			if (socket && socket.readyState === WebSocket.OPEN) return;

			socket = new WebSocket(wsUrl);

			socket.onopen = () => {
				update(s => ({ ...s, connected: true }));
				
				if (pingInterval) clearInterval(pingInterval);
				pingInterval = setInterval(() => {
					if (socket && socket.readyState === WebSocket.OPEN) {
						socket.send(JSON.stringify({ receiver_id: 0, content: "PING" }));
					}
				}, 30000);
			};

			socket.onmessage = (event) => {
				const message = JSON.parse(event.data) as ChatMessage;
				update(s => {
					// Check if it's a delete signal
					if (message.is_deleted) {
						return { ...s, messages: s.messages.filter(m => m.id !== message.id) };
					}

					let newUnreadCount = s.unreadCount;
					// Check if message is addressed to the current logged in user (and not from me)
					if (s.myUserId && message.receiver_id === s.myUserId) {
						if (s.activeUserId !== message.sender_id) {
							// Increment unread count if we're not actively chatting with them
							newUnreadCount += 1;
						}
					}

					let newMessages = s.messages;
					// Add message if it's from/to the active user
					if (s.activeUserId === message.sender_id || s.activeUserId === message.receiver_id) {
						// avoid duplicates
						if (!message.id || !s.messages.some(m => m.id === message.id)) {
							newMessages = [...s.messages, message];
						}
					}
					
					return { ...s, unreadCount: newUnreadCount, messages: newMessages };
				});
			};

			socket.onclose = () => {
				update(s => ({ ...s, connected: false }));
				socket = null;
				if (pingInterval) clearInterval(pingInterval);
				
				// Auto-reconnect after 3 seconds
				setTimeout(() => {
					chatStore.connect();
				}, 3000);
			};
		},
		disconnect: () => {
			if (socket) {
				socket.close();
				socket = null;
			}
		},
		setMessages: (messages: ChatMessage[]) => update(s => ({ ...s, messages })),
		removeMessage: (messageId: number) => update(s => ({ ...s, messages: s.messages.filter(m => m.id !== messageId) })),
		setActiveUser: (userId: number | null) => update(s => ({ ...s, activeUserId: userId })),
		sendMessage: (receiverId: number, content: string) => {
			if (socket && socket.readyState === WebSocket.OPEN) {
				socket.send(JSON.stringify({ receiver_id: receiverId, content }));
			} else {
				console.error("Socket is not connected");
			}
		}
	};
}

export const chatStore = createChatStore();
