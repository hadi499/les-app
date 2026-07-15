<script lang="ts">
	let { contacts = [], activeUserId = null, onSelect } = $props<{
		contacts?: any[];
		activeUserId?: number | null;
		onSelect: (user: any) => void;
	}>();
</script>

<div class="flex-1 overflow-y-auto">
	{#if contacts.length === 0}
		<div class="p-4 text-center text-sm text-gray-500">Belum ada kontak.</div>
	{:else}
		<ul class="divide-y divide-gray-50">
			{#each contacts as contact}
				<li>
					<button
						class={`w-full flex items-center gap-3 p-4 hover:bg-blue-50 transition-colors text-left ${activeUserId === contact.id ? 'bg-blue-50 border-r-4 border-blue-500' : ''}`}
						onclick={() => onSelect(contact)}
					>
						<div class="h-10 w-10 rounded-full bg-gradient-to-tr from-blue-400 to-indigo-500 flex items-center justify-center text-white font-bold flex-shrink-0 relative">
							{contact.username.charAt(0).toUpperCase()}
							{#if contact.unread_count > 0}
								<div class="absolute -top-1 -right-1 bg-red-500 text-white text-[10px] font-bold h-5 min-w-[20px] px-1 rounded-full flex items-center justify-center border-2 border-white">
									{contact.unread_count}
								</div>
							{/if}
							{#if contact.last_active_at && Date.now() - new Date(contact.last_active_at).getTime() < 5 * 60 * 1000}
								<div class="absolute -bottom-0.5 -right-0.5 w-3.5 h-3.5 bg-green-500 border-2 border-white rounded-full" title="Online"></div>
							{/if}
						</div>
						<div class="flex-1 text-left min-w-0">
							<p class="font-semibold text-gray-800 truncate">{contact.username}</p>
							<p class="text-xs text-gray-500 capitalize truncate">{contact.role}</p>
						</div>
					</button>
				</li>
			{/each}
		</ul>
	{/if}
</div>
