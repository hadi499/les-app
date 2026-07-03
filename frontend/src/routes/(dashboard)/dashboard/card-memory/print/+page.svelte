<script lang="ts">
  import { printQueue } from "$lib/stores/print-queue.svelte";
  import { printSettings } from "$lib/stores/print-settings.svelte";
  import { trashCount } from "$lib/stores/trash-count.svelte";
  import type { Card } from "$lib/types";
  import CardGrid from "$lib/components/CardGrid.svelte";
  import Modal from "$lib/components/Modal.svelte";

  let queueCards = $derived(printQueue.items);

  // Group by card size
  let groups = $derived(() => {
    const result: { size: number; cards: Card[] }[] = [];
    const s6 = queueCards.filter((c) => !c.size || String(c.size) === "6");
    const s4 = queueCards.filter((c) => String(c.size) === "4");
    const s2 = queueCards.filter((c) => String(c.size) === "2");
    if (s6.length) result.push({ size: 6, cards: s6 });
    if (s4.length) result.push({ size: 4, cards: s4 });
    if (s2.length) result.push({ size: 2, cards: s2 });
    return result;
  });

  let zoomPct = $derived(Math.round(printSettings.imageZoom * 100));
  let headerPct = $derived(Math.round(printSettings.headerSize * 100));
  let titlePct = $derived(Math.round(printSettings.titleSize * 100));
  let contentPct = $derived(Math.round(printSettings.contentSize * 100));

  function handlePrint() {
    window.print();
  }

  let showClearConfirm = $state(false);

  function handleClear() {
    printQueue.clear();
    showClearConfirm = false;
  }

  function handleCopy(card: Card) {
    printQueue.pushCopy(card);
  }
</script>

<svelte:head>
  <title>Print</title>
</svelte:head>

<div
  class="min-h-screen bg-transparent text-slate-900 print:bg-white print:text-black print:min-h-0"
>
  <header
    class="fixed top-16 md:top-0 left-0 md:left-[var(--sidebar-width)] right-0 bg-white/80/90 backdrop-blur-md px-6 md:px-10 py-3 z-40 transition-all print:hidden"
  >
    <div
      class="max-w-5xl mx-auto flex items-center justify-between gap-3 flex-wrap"
    >
      <div class="flex items-center gap-4">
        <nav class="flex items-center gap-1 bg-white/80 p-1 rounded-xl">
          <a
            href="/dashboard/card-memory"
            class="px-4 py-1.5 text-sm rounded-lg text-slate-600 hover:text-slate-900 hover:bg-white/80/60 font-medium transition-all"
            >Arsip</a
          >
          <a
            href="/dashboard/card-memory/recent"
            class="px-4 py-1.5 text-sm rounded-lg text-slate-600 hover:text-slate-900 hover:bg-white/80/60 font-medium transition-all"
            >Recent</a
          >
          <span
            class="px-4 py-1.5 text-sm rounded-lg bg-slate-800 text-white shadow-md font-medium transition-all"
          >
            Print ({queueCards.length})
          </span>
          <a
            href="/dashboard/card-memory/trash"
            class="px-4 py-1.5 text-sm rounded-lg text-slate-600 hover:text-slate-900 hover:bg-white/80/60 font-medium flex items-center gap-1.5 transition-all"
          >
            Sampah
            {#if trashCount.value > 0}
              <span
                class="bg-red-500 text-white text-[10px] rounded-full px-1.5 py-0.5 leading-none"
                >{trashCount.value}</span
              >
            {/if}
          </a>
        </nav>
      </div>

      <div class="flex items-center gap-2 flex-wrap">
        {#if queueCards.length > 0}
          <button
            onclick={handlePrint}
            class="flex items-center gap-1.5 px-4 py-1.5 text-sm rounded-xl bg-blue-500 text-white hover:bg-blue-600 font-medium transition-all shadow-sm cursor-pointer"
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
            Cetak ({queueCards.length})
          </button>
        {/if}
        <button
          onclick={() => (showClearConfirm = true)}
          class="px-4 py-1.5 text-sm rounded-xl bg-white/80 border border-slate-200 text-slate-600 hover:text-red-600 hover:border-red-300 hover:bg-red-100 font-medium transition-all shadow-sm cursor-pointer"
        >
          Kosongkan
        </button>
      </div>
    </div>
  </header>

  <main class="max-w-5xl mx-auto p-4 print:p-0 pt-24 md:pt-20 print:pt-0">
    <!-- Settings toolbar -->
    <div
      class="mb-4 p-3 w-fit bg-white/80 rounded-xl border border-slate-200 flex items-center gap-4 flex-wrap print:hidden"
    >
      <div class="flex items-center gap-1.5">
        <span class="text-xs text-slate-600 w-16">Gambar</span>
        <button
          onclick={() => printSettings.zoomOut()}
          class="w-7 h-7 flex items-center justify-center rounded border border-slate-300 text-slate-600 hover:bg-white text-sm cursor-pointer"
          >-</button
        >
        <span class="text-xs font-medium text-slate-800 w-10 text-center"
          >{zoomPct}%</span
        >
        <button
          onclick={() => printSettings.zoomIn()}
          class="w-7 h-7 flex items-center justify-center rounded border border-slate-300 text-slate-600 hover:bg-white text-sm cursor-pointer"
          >+</button
        >
      </div>
      <div class="flex items-center gap-1.5">
        <span class="text-xs text-slate-600 w-16">Judul</span>
        <button
          onclick={() => printSettings.titleDown()}
          class="w-7 h-7 flex items-center justify-center rounded border border-slate-300 text-slate-600 hover:bg-white text-sm cursor-pointer"
          >A</button
        >
        <span class="text-xs font-medium text-slate-800 w-10 text-center"
          >{titlePct}%</span
        >
        <button
          onclick={() => printSettings.titleUp()}
          class="w-7 h-7 flex items-center justify-center rounded border border-slate-300 text-slate-600 hover:bg-white text-base font-bold cursor-pointer"
          >A</button
        >
      </div>
      <div class="w-px h-6 bg-gray-200"></div>

      <div class="flex items-center gap-1.5">
        <span class="text-xs text-slate-600 w-16">Konten</span>
        <button
          onclick={() => printSettings.contentDown()}
          class="w-7 h-7 flex items-center justify-center rounded border border-slate-300 text-slate-600 hover:bg-white text-sm cursor-pointer"
          >A</button
        >
        <span class="text-xs font-medium text-slate-800 w-10 text-center"
          >{contentPct}%</span
        >
        <button
          onclick={() => printSettings.contentUp()}
          class="w-7 h-7 flex items-center justify-center rounded border border-slate-300 text-slate-600 hover:bg-white text-base font-bold cursor-pointer"
          >A</button
        >
      </div>
      <div class="w-px h-6 bg-gray-200"></div>

      <button
        onclick={() => printSettings.reset()}
        class="px-2 py-1 text-xs rounded border border-slate-200 text-slate-500 hover:text-slate-600 cursor-pointer"
      >
        Reset
      </button>
    </div>

    <div class="mb-4 print:hidden">
      <p class="text-sm text-slate-600">
        {queueCards.length} kartu dalam antrian
      </p>
    </div>

    {#if queueCards.length === 0}
      <div
        class="text-center py-12 bg-white/80 rounded-xl border border-slate-200"
      >
        <p class="text-slate-500 mb-2">Belum ada kartu di antrian print.</p>
        <a
          href="/dashboard/card-memory"
          class="text-indigo-400 hover:underline text-sm"
          >Kembali ke Arsip untuk menambah kartu</a
        >
      </div>
    {:else}
      {#each groups() as group, gi}
        {@const label =
          group.size === 6 ? "Kecil" : group.size === 4 ? "Sedang" : "Besar"}
        <div class="mb-6 print:mb-0">
          {#if groups().length > 1}
            <p class="text-xs text-slate-600 mb-2 print:hidden">
              {label} · {group.size}/A4 · {group.cards.length} kartu
            </p>
          {/if}
          <div class="flex justify-center print:contents">
            <div class="w-full max-w-[210mm] print:contents">
              <CardGrid
                cards={group.cards}
                perPage={group.size}
                oncopy={handleCopy}
                onremove={(c) => printQueue.removeOne(c.id)}
                contentOnly={true}
                nogap={true}
              />
            </div>
          </div>
        </div>
      {/each}
    {/if}
  </main>

  <footer class="text-center text-xs text-slate-500 py-6 print:hidden">
    Kartu Rangkuman Pelajaran SD & SMP — Cetak di kertas A4 Portrait (210mm ×
    297mm)
  </footer>

  <Modal show={showClearConfirm} onclose={() => (showClearConfirm = false)} maxWidth="max-w-sm">
    <div class="space-y-4 text-center">
      <div
        class="w-12 h-12 mx-auto rounded-full bg-red-100 flex items-center justify-center"
      >
        <svg
          class="w-6 h-6 text-red-600"
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
      </div>
      <div>
        <h3 class="text-lg font-semibold text-slate-900">
          Kosongkan Antrian Print
        </h3>
        <p class="text-sm text-slate-600 mt-1">
          Semua {queueCards.length} kartu akan dihapus dari antrian.
        </p>
      </div>
      <div class="flex gap-2 justify-center pt-2">
        <button
          onclick={() => (showClearConfirm = false)}
          class="px-4 py-2 text-sm rounded-lg border border-slate-300 text-slate-600 hover:bg-transparent text-slate-900 cursor-pointer"
          >Batal</button
        >
        <button
          onclick={handleClear}
          class="px-4 py-2 text-sm rounded-lg bg-red-500 text-white hover:bg-red-700 cursor-pointer"
          >Hapus Semua</button
        >
      </div>
    </div>
  </Modal>
</div>
