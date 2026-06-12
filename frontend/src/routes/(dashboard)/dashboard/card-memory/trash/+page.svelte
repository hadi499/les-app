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

<div class="min-h-screen bg-transparent text-white">
  <header class="fixed top-16 md:top-0 left-0 md:left-64 right-0 bg-zinc-900/90 backdrop-blur-md border-b border-zinc-800 px-6 md:px-10 py-3 z-40 transition-all">
    <div
      class="max-w-5xl mx-auto flex items-center justify-between gap-3 flex-wrap"
    >
      <div class="flex items-center gap-4">
        <nav class="flex items-center gap-1 bg-zinc-800/80 p-1 rounded-xl">
          <a
            href="/dashboard/card-memory"
            class="px-4 py-1.5 text-sm rounded-lg text-blue-200 hover:text-white hover:bg-zinc-900/60 font-medium transition-all"
            >Arsip</a
          >
          <a
            href="/dashboard/card-memory/print"
            class="px-4 py-1.5 text-sm rounded-lg text-blue-200 hover:text-white hover:bg-zinc-900/60 font-medium flex items-center gap-1.5 transition-all"
            >Print</a
          >
          <span
            class="px-4 py-1.5 text-sm rounded-lg bg-zinc-900 text-white shadow-sm font-medium transition-all"
          >
            Sampah ({trashCards.length})
          </span>
        </nav>
      </div>

      <div class="flex items-center gap-2 flex-wrap">
      {#if trashCards.length > 0}
        <button
          onclick={() => (showEmptyConfirm = true)}
          class="px-4 py-1.5 text-sm rounded-xl bg-zinc-900 border border-zinc-800 text-blue-200 hover:text-red-400 hover:border-red-800/50 hover:bg-red-900/30 font-medium transition-all shadow-sm cursor-pointer"
        >
          Kosongkan Sampah
        </button>
      {/if}
      </div>
    </div>
  </header>

  <main class="max-w-5xl mx-auto p-4 pt-24 md:pt-20">
    {#if loading}
      <div class="text-center py-12 text-zinc-500">Memuat...</div>
    {:else if error}
      <div class="text-center py-12 bg-zinc-900 rounded-xl border border-red-800/50">
        <p class="text-red-400 mb-3">{error}</p>
        <button
          onclick={loadTrash}
          class="px-4 py-2 text-sm rounded-lg bg-zinc-800 text-blue-200 hover:bg-zinc-700 cursor-pointer"
          >Coba Lagi</button
        >
      </div>
    {:else if trashCards.length === 0}
      <div class="text-center py-12 bg-zinc-900 rounded-xl border border-zinc-800">
        <p class="text-zinc-500">Tempat sampah kosong</p>
      </div>
    {:else}
      <div
        class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-3"
      >
        {#each trashCards as card}
          <div
            class="bg-zinc-900 border border-red-800/50 rounded-lg p-4 text-base group"
          >
            <div class="flex items-center justify-between mb-2">
              <div class="flex items-center gap-2 min-w-0">
                {#if card.cardType === "image"}
                  <span
                    class="text-xs font-medium bg-purple-900/30 text-purple-500 px-2 py-0.5 rounded shrink-0"
                    >Gambar</span
                  >
                {/if}

                <span class="font-semibold text-red-400 text-sm truncate"
                  >{card.category || "Umum"}</span
                >
              </div>
            </div>
            <p class="font-bold text-blue-100 text-base truncate mb-1">
              {card.title}
            </p>
            <div class="text-sm text-zinc-400 line-clamp-2 mb-3 rich-content">
              {@html renderMathContent(card.content)}
            </div>
            <div class="flex gap-1">
              <button
                onclick={() => handleRestore(card.id)}
                class="flex-1 text-center text-sm py-1 rounded bg-green-900/30 text-green-300 hover:bg-green-900/50 cursor-pointer"
              >
                Pulihkan
              </button>
              <button
                onclick={() => handleForceDeleteClick(card.id)}
                class="flex-1 text-center text-sm py-1 rounded bg-red-900/30 text-red-300 hover:bg-red-900/50 cursor-pointer"
              >
                Hapus
              </button>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </main>

  <!-- Force Delete Confirmation -->
  <Modal show={!!forceTarget} onclose={() => (forceTarget = null)}>
    <div class="space-y-4 text-center">
      <div
        class="w-12 h-12 mx-auto rounded-full bg-red-900/50 flex items-center justify-center"
      >
        <svg
          class="w-6 h-6 text-red-400"
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
        <h3 class="text-lg font-semibold text-white">Hapus Permanen</h3>
        {#if forceCard}
          <p class="text-sm font-medium text-blue-100 mt-1">
            "{forceCard.title}"
          </p>
        {/if}
        <p class="text-sm text-zinc-400 mt-1">Kartu tidak bisa dikembalikan.</p>
      </div>
      <div class="flex gap-2 justify-center pt-2">
        <button
          onclick={() => (forceTarget = null)}
          class="px-4 py-2 text-sm rounded-lg border border-zinc-700 text-blue-200 hover:bg-transparent text-white cursor-pointer"
          >Batal</button
        >
        <button
          onclick={confirmForceDelete}
          class="px-4 py-2 text-sm rounded-lg bg-red-600 text-white hover:bg-red-700 cursor-pointer"
          >Hapus</button
        >
      </div>
    </div>
  </Modal>

  <!-- Empty Trash Confirmation -->
  <Modal show={showEmptyConfirm} onclose={() => (showEmptyConfirm = false)}>
    <div class="space-y-4 text-center">
      <div
        class="w-12 h-12 mx-auto rounded-full bg-red-900/50 flex items-center justify-center"
      >
        <svg
          class="w-6 h-6 text-red-400"
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
        <h3 class="text-lg font-semibold text-white">Kosongkan Sampah</h3>
        <p class="text-sm text-zinc-400 mt-1">
          Semua kartu di sampah akan dihapus permanen.
        </p>
      </div>
      <div class="flex gap-2 justify-center pt-2">
        <button
          onclick={() => (showEmptyConfirm = false)}
          class="px-4 py-2 text-sm rounded-lg border border-zinc-700 text-blue-200 hover:bg-transparent text-white cursor-pointer"
          >Batal</button
        >
        <button
          onclick={confirmEmptyTrash}
          class="px-4 py-2 text-sm rounded-lg bg-red-600 text-white hover:bg-red-700 cursor-pointer"
          >Hapus Semua</button
        >
      </div>
    </div>
  </Modal>
</div>

