<script lang="ts">
	import { onMount, onDestroy, tick } from 'svelte';
	import { chatStore } from '$lib/stores/chatStore';
	import ContactList from '$lib/components/Chat/ContactList.svelte';
	import MessageBubble from '$lib/components/Chat/MessageBubble.svelte';
	import { fade } from 'svelte/transition';
	
	let contacts: any[] = $state([]);
	let activeUser: any = $state(null);
	let messageInput = $state('');
	let chatContainer: HTMLElement | undefined = $state();
	let myUserId: number | null = $state(null);
	let loadingContacts = $state(true);
	let loadingHistory = $state(false);
	
	let showDeleteModal = $state(false);
	let messageToDelete: number | null = $state(null);
	
	let currentPage = $state(1);
	let hasMore = $state(false);
	let loadingMore = $state(false);
	let isPrepending = false;
	
	$effect(() => {
		// Svelte 5 tracks this dependency automatically.
		// Whenever unreadCount changes (goes up or down), refresh the contact list
		// to ensure the unread badges are perfectly in sync with the server.
		const currentCount = $chatStore.unreadCount;
		fetchContacts();
	});

	async function fetchContacts() {
		try {
			const res = await fetch(`/api/chat/contacts`, { credentials: "include" });
			if (res.ok) {
				contacts = await res.json();
			}
		} catch (error) {
			console.error("Failed to load contacts", error);
		}
	}

	onMount(async () => {
		try {
			// Fetch current user info to know our ID
			const meRes = await fetch(`/me`, { credentials: "include" });
			if (meRes.ok) {
				const meData = await meRes.json();
				myUserId = meData.user?.id;
			}

			// Fetch contacts
			await fetchContacts();

			// Restore active user from localStorage
			const savedUserId = localStorage.getItem('lastActiveChatUserId');
			if (savedUserId && contacts.length > 0) {
				const userToSelect = contacts.find(c => c.id.toString() === savedUserId);
				if (userToSelect) {
					selectContact(userToSelect);
				}
			}
		} catch (error) {
			console.error("Failed to initialize chat:", error);
		} finally {
			loadingContacts = false;
		}
	});

	onDestroy(() => {
		chatStore.setActiveUser(null);
	});

	async function selectContact(user: any) {
		activeUser = user;
		chatStore.setActiveUser(user.id);
		localStorage.setItem('lastActiveChatUserId', user.id.toString());
		loadingHistory = true;

		// Locally clear unread count for this contact to make it snappy
		const contactIndex = contacts.findIndex(c => c.id === user.id);
		if (contactIndex !== -1 && contacts[contactIndex].unread_count > 0) {
			contacts[contactIndex].unread_count = 0;
		}

		// Fetch history and mark as read
		try {
			currentPage = 1;
			const res = await fetch(`/api/chat/history/${user.id}?page=1`, { credentials: "include" });
			if (res.ok) {
				const data = await res.json();
				chatStore.setMessages(data.messages || []);
				hasMore = data.has_more || false;
				scrollToBottom();
				
				// Reset unread count since we just opened this chat
				// Optimistic update: we assume any unread messages from this user are now read.
				// For a simple approach, we can just refetch the total unread count from the server
				const unreadRes = await fetch("/api/chat/unread-count", { credentials: "include" });
				if (unreadRes.ok) {
					const unreadData = await unreadRes.json();
					chatStore.setUnreadCount(unreadData.count || 0);
				}
			}
		} catch (error) {
			console.error("Failed to load history", error);
		} finally {
			loadingHistory = false;
		}
	}

	function handleSendMessage() {
		if (!messageInput.trim() || !activeUser) return;
		
		chatStore.sendMessage(activeUser.id, messageInput.trim());
		messageInput = '';
		scrollToBottom();
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter' && !e.shiftKey) {
			e.preventDefault();
			handleSendMessage();
		}
	}

	// Reactive statement to scroll when messages change
	$effect(() => {
		if ($chatStore.messages && !isPrepending) {
			scrollToBottom();
		}
	});

	async function handleScroll() {
		if (!chatContainer || loadingMore || !hasMore || !activeUser) return;
		if (chatContainer.scrollTop < 50) {
			loadingMore = true;
			const oldScrollHeight = chatContainer.scrollHeight;
			const oldScrollTop = chatContainer.scrollTop;
			
			currentPage += 1;
			try {
				const res = await fetch(`/api/chat/history/${activeUser.id}?page=${currentPage}`, { credentials: "include" });
				if (res.ok) {
					const data = await res.json();
					if (data.messages && data.messages.length > 0) {
						isPrepending = true;
						chatStore.prependMessages(data.messages);
					}
					hasMore = data.has_more || false;
					
					await tick();
					if (chatContainer) {
						chatContainer.scrollTop = chatContainer.scrollHeight - oldScrollHeight + oldScrollTop;
					}
					setTimeout(() => { isPrepending = false; }, 50);
				}
			} catch(e) {
				console.error(e);
			} finally {
				loadingMore = false;
			}
		}
	}

	function scrollToBottom() {
		setTimeout(() => {
			if (chatContainer) {
				chatContainer.scrollTop = chatContainer.scrollHeight;
			}
		}, 50);
	}
	
	function formatJustTime(dateString: string) {
		if (!dateString) return '';
		const date = new Date(dateString);
		return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', hour12: false });
	}

	function getDateLabel(dateString: string) {
		if (!dateString) return '';
		const date = new Date(dateString);
		const now = new Date();
		
		const isToday = date.getDate() === now.getDate() && date.getMonth() === now.getMonth() && date.getFullYear() === now.getFullYear();
		
		const yesterday = new Date(now);
		yesterday.setDate(now.getDate() - 1);
		const isYesterday = date.getDate() === yesterday.getDate() && date.getMonth() === yesterday.getMonth() && date.getFullYear() === yesterday.getFullYear();

		if (isToday) {
			return 'Hari Ini';
		} else if (isYesterday) {
			return 'Kemarin';
		} else {
			return date.toLocaleDateString('id-ID', { day: '2-digit', month: 'long', year: 'numeric' });
		}
	}

	function confirmDeleteMessage(id: number) {
		messageToDelete = id;
		showDeleteModal = true;
	}

	function cancelDelete() {
		showDeleteModal = false;
		messageToDelete = null;
	}

	async function executeDeleteMessage() {
		if (!messageToDelete) return;
		const id = messageToDelete;
		
		try {
			const res = await fetch(`/api/chat/messages/${id}`, {
				method: 'DELETE',
				credentials: "include"
			});
			if (res.ok) {
				chatStore.removeMessage(id);
			} else {
				console.error("Failed to delete message");
			}
		} catch (error) {
			console.error("Error deleting message", error);
		} finally {
			cancelDelete();
		}
	}
</script>

<div class="fixed top-16 inset-x-0 bottom-0 z-40 md:relative md:top-auto md:bottom-auto md:inset-x-auto md:z-auto md:h-[calc(100dvh-6rem)] flex flex-col md:flex-row bg-white md:rounded-2xl shadow-sm md:border md:border-gray-100 overflow-hidden">
	
	<!-- Sidebar: Contact List -->
	<div class={`w-full md:w-80 border-r border-gray-100 flex flex-col ${activeUser ? 'hidden md:flex' : 'flex'}`}>
		<div class="p-4 border-b border-gray-100 bg-gray-50/50">
			<h2 class="text-lg font-bold text-gray-800 flex items-center gap-2">
				<svg xmlns="http://www.w3.org/2000/svg" class="text-blue-500 w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"></path></svg>
				Pesan
			</h2>
		</div>
		
		{#if loadingContacts}
			<div class="p-8 text-center text-gray-400">
				<span class="loading loading-spinner loading-md"></span>
			</div>
		{:else}
			<ContactList 
				{contacts} 
				activeUserId={$chatStore.activeUserId}
				onSelect={selectContact} 
			/>
		{/if}
	</div>

	<!-- Main Chat Area -->
	<div class={`flex-1 flex flex-col min-h-0 bg-gray-50/30 ${!activeUser ? 'hidden md:flex' : 'flex'}`}>
		{#if loadingContacts}
			<div class="flex-1 flex items-center justify-center text-gray-400">
				<span class="loading loading-spinner loading-lg text-blue-500"></span>
			</div>
		{:else if !activeUser}
			<div class="flex-1 flex flex-col items-center justify-center text-gray-400" transition:fade>
				<div class="w-24 h-24 bg-blue-50 rounded-full flex items-center justify-center mb-4 text-blue-500">
					<svg xmlns="http://www.w3.org/2000/svg" class="w-12 h-12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"></path></svg>
				</div>
				<p class="text-lg font-medium text-gray-600">Pilih kontak untuk mulai chat</p>
				<p class="text-sm mt-1">Chat langsung secara real-time</p>
			</div>
		{:else}
			<!-- Chat Header -->
			<div class="p-4 border-b border-gray-100 bg-white flex items-center gap-3">
				<button class="md:hidden text-gray-500 mr-1 cursor-pointer" onclick={() => {
					activeUser = null;
					localStorage.removeItem('lastActiveChatUserId');
				}}>
					<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
					</svg>
				</button>
				<div class="h-10 w-10 rounded-full bg-gradient-to-tr from-blue-400 to-indigo-500 flex items-center justify-center text-white font-bold flex-shrink-0">
					{activeUser.username.charAt(0).toUpperCase()}
				</div>
				<div>
					<h3 class="font-bold text-gray-800">{activeUser.username}</h3>
					<p class="text-xs text-green-500 capitalize flex items-center gap-1">
						<span class="w-2 h-2 rounded-full bg-green-500 inline-block"></span>
						{activeUser.role}
					</p>
				</div>
			</div>

			<!-- Messages -->
			<div class="flex-1 overflow-y-auto min-h-0 p-4 md:p-6" bind:this={chatContainer} onscroll={handleScroll}>
				{#if loadingHistory}
					<div class="flex justify-center p-4">
						<span class="loading loading-dots loading-md text-blue-500"></span>
					</div>
				{:else}
					{#if loadingMore}
						<div class="flex justify-center p-2 mb-2">
							<span class="loading loading-spinner loading-xs text-blue-500"></span>
						</div>
					{/if}
					{#each $chatStore.messages as msg, i (msg.id)}
						{@const msgDateStr = msg.created_at || new Date().toISOString()}
						{@const currentDateLabel = getDateLabel(msgDateStr)}
						{@const prevMsgDateStr = i > 0 ? ($chatStore.messages[i - 1].created_at || new Date().toISOString()) : null}
						{@const prevDateLabel = prevMsgDateStr ? getDateLabel(prevMsgDateStr) : null}
						
						{#if currentDateLabel !== prevDateLabel}
							<div class="flex justify-center my-4">
								<span class="bg-gray-200/60 text-gray-600 text-[11px] px-3 py-1 rounded-md font-medium shadow-sm backdrop-blur-sm">
									{currentDateLabel}
								</span>
							</div>
						{/if}
						<MessageBubble 
							id={msg.id}
							isMine={msg.sender_id === myUserId} 
							content={msg.content} 
							time={formatJustTime(msgDateStr)}
							isRead={msg.is_read}
							onDelete={confirmDeleteMessage}
						/>
					{/each}
				{/if}
			</div>

			<!-- Message Input -->
			<div class="p-4 bg-white border-t border-gray-100">
				<div class="flex items-end gap-2 bg-gray-50 p-2 rounded-2xl border border-gray-200 focus-within:border-blue-300 focus-within:bg-white transition-colors">
					<textarea
						bind:value={messageInput}
						onkeydown={handleKeydown}
						placeholder="Tulis pesan..."
						class="flex-1 bg-transparent border-none focus:ring-0 resize-none max-h-32 min-h-[40px] px-3 py-2 text-sm outline-none"
						rows="1"
					></textarea>
					<button 
						class="bg-blue-600 hover:bg-blue-700 text-white rounded-xl p-3 transition-colors disabled:opacity-50 flex-shrink-0 cursor-pointer"
						disabled={!messageInput.trim() || !$chatStore.connected}
						onclick={handleSendMessage}
					>
						<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="22" y1="2" x2="11" y2="13"></line><polygon points="22 2 15 22 11 13 2 9 22 2"></polygon></svg>
					</button>
				</div>
				{#if !$chatStore.connected}
					<p class="text-xs text-red-500 mt-2 text-center">Terputus dari server. Mencoba menghubungkan kembali...</p>
				{/if}
			</div>
		{/if}
	</div>
</div>

<!-- Delete Confirmation Modal -->
{#if showDeleteModal}
	<div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm" transition:fade>
		<div class="bg-white rounded-2xl shadow-xl w-full max-w-sm overflow-hidden" onclick={(e) => e.stopPropagation()}>
			<div class="p-6 text-center">
				<div class="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4 text-red-500">
					<svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 6h18"></path><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path><line x1="10" y1="11" x2="10" y2="17"></line><line x1="14" y1="11" x2="14" y2="17"></line></svg>
				</div>
				<h3 class="text-xl font-bold text-gray-900 mb-2">Hapus Pesan</h3>
				<p class="text-gray-500 mb-6">Apakah Anda yakin ingin menghapus pesan ini?</p>
				<div class="flex gap-3">
					<button 
						class="flex-1 py-2.5 px-4 bg-gray-100 hover:bg-gray-200 text-gray-700 font-medium rounded-xl transition-colors cursor-pointer"
						onclick={cancelDelete}
					>
						Batal
					</button>
					<button 
						class="flex-1 py-2.5 px-4 bg-red-500 hover:bg-red-600 text-white font-medium rounded-xl transition-colors shadow-sm shadow-red-200 cursor-pointer"
						onclick={executeDeleteMessage}
					>
						Hapus
					</button>
				</div>
			</div>
		</div>
	</div>
{/if}
