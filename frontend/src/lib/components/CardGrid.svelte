<script lang="ts">
  import type { Card } from '$lib/types';
  import { printSettings } from '$lib/stores/print-settings.svelte';
  import CardItem from './CardItem.svelte';

  let { cards, oncopy, onremove, onedit, ondelete, ontoggleprint, showqueue, perPage: perPageProp, nopaper, selectable, selectedIds, onselect, contentOnly, nogap, isPrint = true }: {
    cards: Card[];
    oncopy?: (card: Card) => void;
    onremove?: (card: Card) => void;
    onedit?: (card: Card) => void;
    ondelete?: (card: Card) => void;
    ontoggleprint?: (card: Card) => void;
    showqueue?: (id: string) => boolean;
    perPage?: number;
    nopaper?: boolean;
    selectable?: boolean;
    selectedIds?: Set<string>;
    onselect?: (card: Card) => void;
    contentOnly?: boolean;
    nogap?: boolean;
    isPrint?: boolean;
  } = $props();

  let perPage = $derived(perPageProp ?? printSettings.cardsPerPage);
  let cols = $derived(perPage === 2 ? 1 : 2);
  let rows = $derived(perPage === 6 ? 3 : perPage === 4 ? 2 : 2);

  let chunks = $derived(() => {
    const result: Card[][] = [];
    for (let i = 0; i < cards.length; i += perPage) {
      result.push(cards.slice(i, i + perPage));
    }
    return result;
  });
</script>

<div class="print-card-grid" class:no-paper={nopaper}>
  {#each chunks() as page, pageIdx}
    <div class="card-page" class:no-paper={nopaper}>
      <div class="card-page-grid" class:no-gap={nogap} style="grid-template-columns: repeat({cols}, 1fr); grid-template-rows: repeat({rows}, 1fr);">
        {#each page as card, idx}
          <CardItem {card} index={pageIdx * perPage + idx} {oncopy} {onremove} {onedit} {ondelete} {ontoggleprint} showqueue={showqueue?.(card.id)} {selectable} selected={selectedIds?.has(card.id)} {onselect} {contentOnly} {isPrint} />
        {/each}
        {#each Array(perPage - page.length) as _}
          <div class="border border-dashed border-gray-200"></div>
        {/each}
      </div>
    </div>
  {/each}
</div>

{#if chunks().length === 0}
  <div class="text-center py-12 text-gray-400">
    <p class="text-lg mb-2">Belum ada kartu</p>
    <p class="text-sm">Tambah kartu baru menggunakan form di atas</p>
  </div>
{/if}

<style>
  .print-card-grid {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .print-card-grid.no-paper {
    gap: 1rem;
  }

  .card-page {
    background: white;
    box-shadow: 0 1px 3px rgba(0,0,0,0.12);
    border-radius: 4px;
    width: 210mm;
    height: 297mm;
    padding: 10mm;
    box-sizing: border-box;
    overflow: hidden;
  }

  .card-page.no-paper {
    background: transparent;
    box-shadow: none;
    border-radius: 0;
    width: 210mm;
    height: auto;
    padding: 0;
    aspect-ratio: unset;
  }

  .card-page.no-paper .card-page-grid {
    aspect-ratio: 210 / 297;
    height: auto;
    gap: 3mm;
  }

  .card-page-grid {
    display: grid;
    gap: 3mm;
    height: 100%;
  }

  .card-page-grid.no-gap {
    gap: 0;
  }

  /* Print styles */
  @media print {
    .print-card-grid {
      display: block;
    }

    .card-page {
      box-shadow: none;
      border-radius: 0;
    }

    .card-page + .card-page {
      page-break-before: always;
    }
  }

  @page {
    size: A4;
    margin: 0;
  }
</style>
