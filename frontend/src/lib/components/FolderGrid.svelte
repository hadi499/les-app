<script lang="ts">
  import type { Card } from "$lib/types";

  let {
    groupedCards,
    onSelectCategory,
    isTeacher = false,
    onEditFolder = () => {},
    onDeleteFolder = () => {},
  }: {
    groupedCards: Record<string, Card[]>;
    onSelectCategory: (category: string) => void;
    isTeacher?: boolean;
    onEditFolder?: (category: string) => void;
    onDeleteFolder?: (category: string) => void;
  } = $props();

  let openMenuId = $state<string | null>(null);
</script>

<div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-3 md:gap-4">
  {#each Object.entries(groupedCards) as [category, categoryCards]}
    <div
      role="button"
      tabindex="0"
      onclick={() => onSelectCategory(category)}
      onkeydown={(e) => e.key === "Enter" && onSelectCategory(category)}
      class="bg-white/80 border border-slate-200 rounded-2xl p-3 md:p-5 flex flex-col items-center justify-center gap-3 hover:bg-slate-50 hover:border-blue-200 hover:shadow-md transition-all cursor-pointer group relative {openMenuId === category ? 'z-50' : 'z-0'}"
    >
      <div class="relative">
        <svg
          class="w-12 h-12 md:w-16 md:h-16 text-blue-400 group-hover:text-blue-500 transition-colors"
          viewBox="0 0 24 24"
          fill="currentColor"
        >
          <path
            d="M4 4c-1.11 0-2 .89-2 2v12c0 1.11.89 2 2 2h16c1.11 0 2-.89 2-2V8c0-1.11-.89-2-2-2h-8l-2-2H4z"
          />
        </svg>
        <div class="absolute inset-0 flex items-center justify-center mt-2">
          <span
            class="bg-white/90 text-blue-700 text-[10px] font-bold px-1.5 py-0.5 rounded-md shadow-sm"
          >
            {categoryCards.length}
          </span>
        </div>
      </div>
      <div class="text-center w-full">
        <h3 class="font-semibold text-slate-700 text-xs md:text-sm break-words leading-tight w-full px-1">
          {category}
        </h3>
      </div>
      {#if isTeacher && category !== 'Tidak Berkategori'}
        <div class="absolute top-2 right-2 flex items-center gap-1 transition-all z-10">
          <button
            onclick={(e) => {
              e.stopPropagation();
              openMenuId = openMenuId === category ? null : category;
            }}
            class="p-1.5 text-slate-400 hover:bg-slate-100 rounded-xl transition-all cursor-pointer {openMenuId === category ? 'bg-slate-100' : ''}"
            aria-label="Menu folder"
          >
            <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24"><path d="M12 8a2 2 0 110-4 2 2 0 010 4zm0 6a2 2 0 110-4 2 2 0 010 4zm0 6a2 2 0 110-4 2 2 0 010 4z"/></svg>
          </button>

          <div class="absolute right-0 top-10 {openMenuId === category ? 'flex flex-col' : 'hidden'} p-2 bg-white border border-slate-100 shadow-xl rounded-xl items-center gap-1 transition-all z-50 w-max">
            <button
              onclick={(e) => {
                e.stopPropagation();
                openMenuId = null;
                onEditFolder(category);
              }}
              class="w-full flex items-center gap-2 p-2 text-slate-600 hover:text-blue-600 cursor-pointer rounded-lg hover:bg-blue-50 transition-colors"
              title="Edit folder"
            >
              <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path></svg>
              <span class="text-sm font-medium pr-2">Edit</span>
            </button>
            <button
              onclick={(e) => {
                e.stopPropagation();
                openMenuId = null;
                onDeleteFolder(category);
              }}
              class="w-full flex items-center gap-2 p-2 text-red-600 hover:text-red-500 cursor-pointer rounded-lg hover:bg-red-50 transition-colors"
              title="Hapus folder"
            >
              <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
              <span class="text-sm font-medium pr-2">Hapus</span>
            </button>
          </div>
        </div>
      {/if}
    </div>
  {/each}
</div>
