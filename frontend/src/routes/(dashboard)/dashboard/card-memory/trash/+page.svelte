<script lang="ts">
  import { browser } from "$app/environment";
  import type { Card } from "$lib/types";
  import { cardsStore } from "$lib/stores/cards.svelte";
  import * as api from "$lib/api";
  import { renderMathContent } from "$lib/extensions/math";
  import Modal from "$lib/components/Modal.svelte";
  import { toast } from "$lib/stores/toast.svelte";
  import { trashCount } from "$lib/stores/trash-count.svelte";

  let trashCards = $state<Card[]>([]);
  let loading = $state(true);
  let error = $state("");
  let forceTarget = $state<string | null>(null);
  let showEmptyConfirm = $state(false);
  let forceCard = $derived(trashCards.find((c) => c.id === forceTarget));

  async function loadTrash() {
    loading = true;
    error = "";
    try {
      const res = await api.fetchTrashCards();
      trashCards = res.data;
      trashCount.set(res.total);
    } catch (e: any) {
      error = e.message || "Gagal memuat data";
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    if (browser) loadTrash();
  });

  async function handleRestore(id: string) {
    try {
      await api.restoreCard(id);
      trashCards = trashCards.filter((c) => c.id !== id);
      trashCount.dec();
      await cardsStore.fetch();
      toast.success("Kartu berhasil dipulihkan");
    } catch {
      toast.error("Gagal memulihkan kartu");
    }
  }

  function handleForceDeleteClick(id: string) {
    forceTarget = id;
  }

  async function confirmForceDelete() {
    if (!forceTarget) return;
    await api.forceDeleteCard(forceTarget);
    trashCards = trashCards.filter((c) => c.id !== forceTarget);
    trashCount.dec();
    toast.success("Kartu dihapus permanen");
    forceTarget = null;
  }

  async function confirmEmptyTrash() {
    await api.emptyTrash();
    trashCards = [];
    trashCount.set(0);
    toast.success("Sampah dikosongkan");
    showEmptyConfirm = false;
  }
</script>

<svelte:head>
  <title>Trash</title>
</svelte:head>

<div class="min-h-screen bg-transparent text-slate-900">
  <header
    class="fixed top-16 md:top-0 left-0 md:left-[var(--sidebar-width)] right-0 bg-white/80/90 backdrop-blur-md px-6 md:px-10 py-3 z-40 transition-all"
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
            href="/dashboard/card-memory/print"
            class="px-4 py-1.5 text-sm rounded-lg text-slate-600 hover:text-slate-900 hover:bg-white/80/60 font-medium flex items-center gap-1.5 transition-all"
            >Print</a
          >
          <span
            class="px-4 py-1.5 text-sm rounded-lg bg-white/80 text-slate-900 shadow-sm font-medium transition-all"
          >
            Sampah ({trashCards.length})
          </span>
        </nav>
      </div>

      <div class="flex items-center gap-2 flex-wrap">
        {#if trashCards.length > 0}
          <button
            onclick={() => (showEmptyConfirm = true)}
            class="px-4 py-1.5 text-sm rounded-xl bg-white/80 border border-slate-200 text-slate-600 hover:text-red-600 hover:border-red-300 hover:bg-red-100 font-medium transition-all shadow-sm cursor-pointer"
          >
            Kosongkan Sampah
          </button>
        {/if}
      </div>
    </div>
  </header>

  <main class="max-w-5xl mx-auto p-4 pt-24 md:pt-20">
    {#if loading}
      <div class="text-center py-12 text-slate-500">Memuat...</div>
    {:else if error}
      <div
        class="text-center py-12 bg-white/80 rounded-xl border border-red-300"
      >
        <p class="text-red-600 mb-3">{error}</p>
        <button
          onclick={loadTrash}
          class="px-4 py-2 text-sm rounded-lg bg-white text-slate-600 hover:bg-slate-50 cursor-pointer"
          >Coba Lagi</button
        >
      </div>
    {:else if trashCards.length === 0}
      <div
        class="text-center py-12 bg-white/80 rounded-xl border border-slate-200"
      >
        <p class="text-slate-500">Tempat sampah kosong</p>
      </div>
    {:else}
      <div class="flex flex-col gap-2">
        {#each trashCards as card}
          <div
            class="flex items-center justify-between bg-white border border-slate-200 rounded-lg p-3 hover:border-slate-300 transition-colors gap-4"
          >
            <div class="flex items-center gap-3 overflow-hidden">
              <div class="flex items-center gap-2 shrink-0">
                {#if card.cardType === "image"}
                  <span
                    class="text-[10px] font-bold uppercase tracking-wider bg-purple-100 text-purple-600 px-2 py-0.5 rounded"
                    >Img</span
                  >
                {/if}
                {#if card.card_folder?.name}
                  <span class="text-xs font-medium text-red-500 bg-red-50 px-2 py-0.5 rounded"
                    >{card.card_folder.name}</span
                  >
                {/if}
              </div>
              <p class="font-medium text-slate-700 text-sm truncate">
                {card.title}
              </p>
            </div>
            <div class="flex gap-2 shrink-0">
              <button
                onclick={() => handleRestore(card.id)}
                class="p-2 text-green-600 bg-green-50 hover:bg-green-100 rounded-md transition-colors border border-green-200 cursor-pointer"
                title="Pulihkan"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-rotate-ccw"><path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"/><path d="M3 3v5h5"/></svg>
              </button>
              <button
                onclick={() => handleForceDeleteClick(card.id)}
                class="p-2 text-red-600 bg-red-50 hover:bg-red-100 rounded-md transition-colors border border-red-200 cursor-pointer"
                title="Hapus Permanen"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-trash-2"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/><line x1="10" x2="10" y1="11" y2="17"/><line x1="14" x2="14" y1="11" y2="17"/></svg>
              </button>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </main>

  <!-- Force Delete Confirmation -->
  <Modal show={!!forceTarget} onclose={() => (forceTarget = null)} maxWidth="max-w-sm">
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
        <h3 class="text-lg font-semibold text-slate-900">Hapus Permanen</h3>
        {#if forceCard}
          <p class="text-sm font-medium text-slate-800 mt-1">
            "{forceCard.title}"
          </p>
        {/if}
        <p class="text-sm text-slate-600 mt-1">
          Kartu tidak bisa dikembalikan.
        </p>
      </div>
      <div class="flex gap-2 justify-center pt-2">
        <button
          onclick={() => (forceTarget = null)}
          class="px-4 py-2 text-sm rounded-lg border border-slate-300 text-slate-600 hover:bg-transparent text-slate-900 cursor-pointer"
          >Batal</button
        >
        <button
          onclick={confirmForceDelete}
          class="px-4 py-2 text-sm rounded-lg bg-red-500 text-white hover:bg-red-700 cursor-pointer"
          >Hapus</button
        >
      </div>
    </div>
  </Modal>

  <!-- Empty Trash Confirmation -->
  <Modal show={showEmptyConfirm} onclose={() => (showEmptyConfirm = false)} maxWidth="max-w-sm">
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
        <h3 class="text-lg font-semibold text-slate-900">Kosongkan Sampah</h3>
        <p class="text-sm text-slate-600 mt-1">
          Semua kartu di sampah akan dihapus permanen.
        </p>
      </div>
      <div class="flex gap-2 justify-center pt-2">
        <button
          onclick={() => (showEmptyConfirm = false)}
          class="px-4 py-2 text-sm rounded-lg border border-slate-300 text-slate-600 hover:bg-transparent text-slate-900 cursor-pointer"
          >Batal</button
        >
        <button
          onclick={confirmEmptyTrash}
          class="px-4 py-2 text-sm rounded-lg bg-red-500 text-white hover:bg-red-700 cursor-pointer"
          >Hapus Semua</button
        >
      </div>
    </div>
  </Modal>
</div>
