<script lang="ts">
  import type { Card } from "$lib/types";
  import { printSettings } from "$lib/stores/print-settings.svelte";
  import { renderMathContent } from "$lib/extensions/math";

  let {
    card,
    index,
    oncopy,
    onremove,
    onedit,
    ondelete,
    ontoggleprint,
    showqueue,
    selectable,
    selected,
    onselect,
  }: {
    card: Card;
    index: number;
    oncopy?: (card: Card) => void;
    onremove?: (card: Card) => void;
    onedit?: (card: Card) => void;
    ondelete?: (card: Card) => void;
    ontoggleprint?: (card: Card) => void;
    showqueue?: boolean;
    selectable?: boolean;
    selected?: boolean;
    onselect?: (card: Card) => void;
  } = $props();

  let zoom = $derived(printSettings.imageZoom);

  let headerScale = $derived(printSettings.headerSize);
  let titleScale = $derived(printSettings.titleSize);
  let contentScale = $derived(printSettings.contentSize);

  let perPage = $derived(parseInt(card.size || "6") || 6);
  let isImageCard = $derived(card.cardType === "image");
  let layoutScale = $derived(1.0);
  let imageBasis = $derived(isImageCard ? "80%" : "50%");
  let contentClamp = $derived(
    perPage === 6 ? 6 : perPage === 4 ? 12 : isImageCard ? 0 : 24,
  );
  let titleClamp = $derived(
    isImageCard ? 2 : perPage === 2 ? 3 : perPage === 4 ? 2 : 2,
  );

  let headerFont = $derived(`${10 * headerScale * layoutScale}pt`);
  let titleFont = $derived(`${14 * titleScale * layoutScale}pt`);
  let contentFont = $derived(`${12 * contentScale * layoutScale}pt`);
  let numFont = $derived(`${8 * contentScale * layoutScale}pt`);

  let renderedContent = $derived(renderMathContent(card.content || ""));
</script>

<div
  class="card-item flex flex-col border-2 {selected
    ? 'border-indigo-500 bg-indigo-50/30'
    : 'border-gray-400'} bg-white overflow-hidden h-full relative group"
>
  <!-- Selection checkbox (top-left) -->
  {#if selectable}
    <div class="absolute top-1 left-1 z-10 print:hidden">
      <button
        onclick={() => onselect?.(card)}
        class="w-6 h-6 flex items-center justify-center rounded border-2 cursor-pointer {selected
          ? 'bg-indigo-600 border-indigo-600 text-orange-950'
          : 'bg-white/90 border-gray-300 hover:border-indigo-400'}"
      >
        {#if selected}
          <svg
            class="w-4 h-4"
            fill="none"
            stroke="currentColor"
            stroke-width="3"
            viewBox="0 0 24 24"
            ><path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M5 13l4 4L19 7"
            /></svg
          >
        {/if}
      </button>
    </div>
  {/if}
  <!-- Hover action buttons -->
  {#if oncopy || onremove || onedit || ondelete || ontoggleprint}
    <div
      class="absolute top-1 right-1 hidden group-hover:flex gap-1 z-10 print:hidden"
    >
      {#if ontoggleprint}
        <button
          onclick={() => ontoggleprint(card)}
          class="w-10 h-10 flex items-center justify-center rounded border {showqueue
            ? 'bg-indigo-50 border-indigo-300 text-indigo-600 hover:bg-indigo-100'
            : 'bg-white/90 border-gray-300 text-gray-500 hover:bg-indigo-50 hover:border-indigo-300 hover:text-indigo-600'} cursor-pointer"
          title={showqueue ? "Hapus dari Print" : "Tambah ke Print"}
        >
          <svg
            class="w-5 h-5"
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
      {/if}
      {#if onedit}
        <button
          onclick={() => onedit(card)}
          class="w-10 h-10 flex items-center justify-center rounded bg-white/90 border border-gray-300 hover:bg-green-50 hover:border-green-300 text-gray-500 hover:text-green-600 cursor-pointer"
          title="Edit"
        >
          <svg
            class="w-6 h-6"
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
      {/if}
      {#if ondelete}
        <button
          onclick={() => ondelete(card)}
          class="w-10 h-10 flex items-center justify-center rounded bg-white/90 border border-gray-300 hover:bg-red-50 hover:border-red-300 text-gray-500 hover:text-red-600 cursor-pointer"
          title="Hapus"
        >
          <svg
            class="w-6 h-6"
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
      {#if oncopy}
        <button
          onclick={() => oncopy(card)}
          class="w-10 h-10 flex items-center justify-center rounded bg-white/90 border border-gray-300 hover:bg-blue-50 hover:border-blue-300 text-gray-500 hover:text-blue-600 cursor-pointer"
          title="Salin"
        >
          <svg
            class="w-6 h-6"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            viewBox="0 0 24 24"
          >
            <rect x="9" y="9" width="13" height="13" rx="2" ry="2" />
            <path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1" />
          </svg>
        </button>
      {/if}
      {#if onremove}
        <button
          onclick={() => onremove(card)}
          class="w-10 h-10 flex items-center justify-center rounded bg-white/90 border border-gray-300 hover:bg-red-50 hover:border-red-300 text-gray-500 hover:text-red-600 cursor-pointer"
          title="Hapus dari antrian"
        >
          <svg
            class="w-6 h-6"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            viewBox="0 0 24 24"
          >
            <line x1="18" y1="6" x2="6" y2="18" /><line
              x1="6"
              y1="6"
              x2="18"
              y2="18"
            />
          </svg>
        </button>
      {/if}
    </div>
  {/if}

  <!-- Header: Kelas & Kategori (hidden for image card) -->
  {#if !isImageCard}
    <div class="text-center shrink-0 py-1.5 px-3">
      <div
        class="font-semibold uppercase tracking-wider truncate text-black"
        style="font-size: {headerFont}; "
      >
        {card.category || "Umum"}
      </div>
    </div>
  {/if}

  {#if isImageCard}
    <!-- Image Card: 20% title + 80% image -->
    <div
      class="shrink-0 px-3 py-2 text-center border-b border-gray-200 flex items-center justify-center"
      style="flex-basis: 20%;"
    >
      <h3
        class="font-bold text-black leading-tight text-center w-full"
        style="font-size: {titleFont}; display: -webkit-box; -webkit-line-clamp: {titleClamp}; -webkit-box-orient: vertical; overflow: hidden;"
      >
        {card.title || "Judul"}
      </h3>
    </div>
    <div class="flex-1 overflow-hidden" style="flex-basis: 80%;">
      {#if card.image}
        <img
          src={card.image}
          alt={card.title}
          class="w-full h-full object-contain"
          style="transform: scale({zoom}); transform-origin: center center;"
        />
      {:else}
        <div
          class="w-full h-full flex items-center justify-center text-gray-300"
        >
          <svg
            class="w-16 h-16"
            fill="none"
            stroke="currentColor"
            stroke-width="1"
            viewBox="0 0 24 24"
          >
            <rect x="3" y="3" width="18" height="18" rx="2" ry="2" /><circle
              cx="8.5"
              cy="8.5"
              r="1.5"
            /><polyline points="21 15 16 10 5 21" />
          </svg>
        </div>
      {/if}
    </div>
  {:else}
    <!-- Section 1: Image (only when present, takes 50% height) -->
    {#if card.image}
      <div class="shrink-0 overflow-hidden" style="flex-basis: {imageBasis};">
        <img
          src={card.image}
          alt={card.title}
          class="w-full h-full object-contain"
          style="transform: scale({zoom}); transform-origin: center center;"
        />
      </div>
    {/if}

    <!-- Section 2 + 3: Title + Content (fills remaining height) -->
    <div class="flex-1 flex flex-col min-h-0">
      <div class="px-3 py-1.5 text-center shrink-0">
        <h3
          class="font-bold text-black leading-tight"
          style="font-size: {titleFont}; display: -webkit-box; -webkit-line-clamp: {titleClamp}; -webkit-box-orient: vertical; overflow: hidden;"
        >
          {card.title || "Judul"}
        </h3>
      </div>

      <div
        class="flex-1 px-3 py-2 text-black leading-relaxed overflow-hidden"
        style="font-size: {contentFont};"
      >
        <div
          class="rich-content"
          style="display: -webkit-box; -webkit-line-clamp: {contentClamp}; -webkit-box-orient: vertical; overflow: hidden;"
        >
          {@html renderedContent || "Konten ringkasan..."}
        </div>
      </div>
    </div>
  {/if}
</div>
