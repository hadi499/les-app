<script lang="ts">
  import type { Card } from "$lib/types";
  import { printQueue } from "$lib/stores/print-queue.svelte";

  let {
    cards,
    activeCategory,
    isTeacher,
    selectedIds,
    onToggleSelect,
    onEdit,
    onDelete,
    onDetailClick,
  }: {
    cards: Card[];
    activeCategory: string | null;
    isTeacher: boolean;
    selectedIds: Set<string>;
    onToggleSelect: (card: Card) => void;
    onEdit: (card: Card) => void;
    onDelete: (card: Card) => void;
    onDetailClick: (card: Card) => void;
  } = $props();

  function layoutLabel(card: Card): string {
    if (card.cardType === "image") return "Gambar";
    return `${card.size || "6"}/A4`;
  }
</script>

<div class="bg-white/80 rounded-xl border border-slate-200 overflow-hidden">
  <div class="overflow-x-auto">
    <table class="w-full text-sm min-w-[600px]">
      <thead>
        <tr
          class="bg-transparent border-b border-slate-200 text-left text-xs text-slate-800 uppercase tracking-wider"
        >
          {#if isTeacher}<th class="px-4 py-3 w-8"></th>{/if}
          <th class="px-4 py-3">Title</th>
          {#if activeCategory === null}
            <th class="px-4 py-3">Folder</th>
          {/if}
          <th class="px-4 py-3">Layout Print</th>
          <th class="px-4 py-3 text-right">Aksi</th>
        </tr>
      </thead>
      <tbody>
        {#each cards as card (card.id)}
          <tr
            class="border-b border-slate-200 hover:bg-slate-50/50 text-slate-900 transition-colors"
          >
            {#if isTeacher}
              <td class="px-4 py-3">
                <button
                  onclick={() => onToggleSelect(card)}
                  class="w-5 h-5 flex items-center justify-center rounded border-2 cursor-pointer {selectedIds.has(
                    card.id,
                  )
                    ? 'bg-blue-500 border-blue-500'
                    : 'border-slate-300'}"
                >
                  {#if selectedIds.has(card.id)}
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
                </button>
              </td>
            {/if}
            <td class="px-4 py-3 font-semibold text-slate-800 text-md"
              >{card.title}</td
            >
            {#if activeCategory === null}
              <td class="px-4 py-3 text-slate-600"
                >{card.card_folder?.name || "-"}</td
              >
            {/if}
            <td class="px-4 py-3">
              <span
                class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium {card.cardType ===
                'image'
                  ? 'bg-purple-100 text-purple-800'
                  : 'bg-slate-100 text-slate-600'}"
              >
                {layoutLabel(card)}
              </span>
            </td>
            <td class="px-4 py-3 text-right">
              <div class="flex items-center justify-end gap-1">
                <button
                  onclick={() => onDetailClick(card)}
                  class="p-1.5 rounded-lg text-blue-500 hover:text-blue-700 hover:bg-blue-50 cursor-pointer"
                  title="Detail"
                >
                  <svg
                    class="w-4 h-4"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    viewBox="0 0 24 24"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                    />
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                    />
                  </svg>
                </button>
                <button
                  onclick={() => printQueue.toggle(card)}
                  class="p-1.5 rounded-lg cursor-pointer {printQueue.has(
                    card.id,
                  )
                    ? 'bg-slate-200 text-slate-600'
                    : 'text-emerald-500 hover:text-emerald-700 hover:bg-emerald-50'}"
                  title={printQueue.has(card.id)
                    ? "Hapus dari antrian print"
                    : "Tambah ke antrian print"}
                >
                  <svg
                    class="w-4 h-4"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    viewBox="0 0 24 24"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z"
                    />
                  </svg>
                </button>
                {#if isTeacher}
                  <button
                    onclick={() => onEdit(card)}
                    class="p-1.5 rounded-lg text-amber-500 hover:text-amber-700 hover:bg-amber-50 cursor-pointer"
                    title="Edit"
                  >
                    <svg
                      class="w-4 h-4"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                      />
                    </svg>
                  </button>
                  <button
                    onclick={() => onDelete(card)}
                    class="p-1.5 rounded-lg text-red-600 hover:text-red-700 hover:bg-red-100 cursor-pointer"
                    title="Hapus"
                  >
                    <svg
                      class="w-4 h-4"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                      />
                    </svg>
                  </button>
                {/if}
              </div>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>
