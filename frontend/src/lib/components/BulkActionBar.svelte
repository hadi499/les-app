<script lang="ts">
  let {
    allSelected,
    selectedCount,
    totalCards,
    onSelectAll,
    onDeleteBulk,
  }: {
    allSelected: boolean;
    selectedCount: number;
    totalCards: number;
    onSelectAll: () => void;
    onDeleteBulk: () => void;
  } = $props();
</script>

<div class="flex items-center gap-3 mb-3">
  <button
    onclick={onSelectAll}
    class="flex items-center gap-1.5 text-xs text-slate-600 hover:text-slate-800 cursor-pointer"
  >
    <span
      class="w-5 h-5 flex items-center justify-center rounded border-2 {allSelected
        ? 'bg-blue-500 border-blue-500'
        : 'border-slate-300'}"
    >
      {#if allSelected}
        <svg
          class="w-3.5 h-3.5 text-slate-900"
          fill="none"
          stroke="currentColor"
          stroke-width="3"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M5 13l4 4L19 7"
          />
        </svg>
      {/if}
    </span>
    {allSelected ? "Batal Pilih Semua" : "Pilih Semua"} ({totalCards})
  </button>
  {#if selectedCount > 0}
    <span class="text-xs text-slate-500">|</span>
    <span class="text-xs text-slate-600">{selectedCount} dipilih</span>
    <button
      onclick={onDeleteBulk}
      class="px-3 py-1 text-xs rounded-lg bg-red-500 text-white hover:bg-red-700 cursor-pointer"
    >
      Hapus ({selectedCount})
    </button>
  {/if}
</div>
