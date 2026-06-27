<script lang="ts">
  import type { Card } from "$lib/types";

  let {
    groupedCards,
    onSelectCategory,
  }: {
    groupedCards: Record<string, Card[]>;
    onSelectCategory: (category: string) => void;
  } = $props();
</script>

<div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
  {#each Object.entries(groupedCards) as [category, categoryCards]}
    <button
      onclick={() => onSelectCategory(category)}
      class="bg-white/80 border border-slate-200 rounded-2xl p-5 flex flex-col items-center justify-center gap-3 hover:bg-slate-50 hover:border-blue-200 hover:shadow-md transition-all cursor-pointer group"
    >
      <div class="relative">
        <svg
          class="w-16 h-16 text-blue-400 group-hover:text-blue-500 transition-colors"
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
        <h3 class="font-semibold text-slate-700 text-sm truncate w-full px-1">
          {category}
        </h3>
      </div>
    </button>
  {/each}
</div>
