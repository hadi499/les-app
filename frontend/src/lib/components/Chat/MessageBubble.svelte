<script lang="ts">
	let { id = 0, isMine = false, content = "", time = "", onDelete } = $props<{
		id?: number;
		isMine?: boolean;
		content?: string;
		time?: string;
		onDelete?: (id: number) => void;
	}>();

	let pressTimer: ReturnType<typeof setTimeout>;

	function startPress() {
		if (!isMine || !onDelete || !id) return;
		pressTimer = setTimeout(() => {
			onDelete(id!);
		}, 600);
	}

	function cancelPress() {
		if (pressTimer) clearTimeout(pressTimer);
	}
</script>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div 
	class={`flex w-full ${isMine ? 'justify-end select-none' : 'justify-start'} mb-4 group`}
	style={isMine ? '-webkit-touch-callout: none;' : ''}
	onmousedown={startPress}
	onmouseup={cancelPress}
	onmouseleave={cancelPress}
	ontouchstart={startPress}
	ontouchend={cancelPress}
	ontouchmove={cancelPress}
	oncontextmenu={(e) => { if (isMine) e.preventDefault(); }}
>
	{#if isMine && onDelete && id}
		<div class="hidden md:group-hover:flex items-center justify-center mr-2 opacity-0 group-hover:opacity-100 transition-opacity">
			<button 
				class="text-gray-400 hover:text-red-500 p-1 rounded-full hover:bg-red-50 transition-colors"
				title="Hapus pesan"
				onclick={() => onDelete(id)}
			>
				<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>
			</button>
		</div>
	{/if}
	<div class={`max-w-[85%] md:max-w-[75%] rounded-2xl px-4 py-2 text-sm shadow-sm transition-transform active:scale-[0.98] ${isMine ? 'bg-blue-600 text-white rounded-br-none' : 'bg-white text-gray-800 border border-gray-100 rounded-bl-none'}`}>
		<p class="whitespace-pre-wrap break-words">{content}</p>
		{#if time}
			<span class={`text-[10px] mt-1 block text-right ${isMine ? 'text-blue-100' : 'text-gray-400'}`}>
				{time}
			</span>
		{/if}
	</div>
</div>
