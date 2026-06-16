<script lang="ts">
  import { onMount } from "svelte";
  import type { Card } from "$lib/types";
  import { printQueue } from "$lib/stores/print-queue.svelte";
  import * as api from "$lib/api";
  import CardForm from "$lib/components/CardForm.svelte";
  import CardImageForm from "$lib/components/CardImageForm.svelte";
  import Modal from "$lib/components/Modal.svelte";
  import CardItem from "$lib/components/CardItem.svelte";
  import { toast } from "$lib/stores/toast.svelte";
  import { trashCount } from "$lib/stores/trash-count.svelte";

  let showForm = $state(false);
  let isTeacher = $state(false);
  let showImageForm = $state(false);
  let editingCard = $state<Card | null>(null);
  let deleteCardId = $state<string | null>(null);
  let detailCard = $state<Card | null>(null);
  let searchQuery = $state("");
  let cardWrapperWidth = $state(0);

  let cards = $state<Card[]>([]);
  let total = $state(0);
  let totalPages = $state(1);
  let currentPage = $state(1);
  let loading = $state(true);
  let error = $state("");

  async function loadCards(params: Record<string, any> = {}) {
    loading = true;
    error = "";
    try {
      const q: Record<string, any> = {
        page: params.page || 1,
        limit: 20,
        all: true,
        ...params,
      };
      const res = await api.fetchCards(q);
      cards = res.data;
      total = res.total;
      totalPages = res.totalPages;
      currentPage = res.page;
    } catch (e: any) {
      error = e.message || "Gagal memuat data";
    } finally {
      loading = false;
    }
  }

  let deleteTarget = $derived(cards.find((c) => c.id === deleteCardId));
  let queueCount = $derived(printQueue.count);
  let trashCountVal = $derived(trashCount.value);

  // Bulk selection
  let selectedIds = $state(new Set<string>());
  let showBulkDeleteConfirm = $state(false);

  function toggleSelect(card: Card) {
    const next = new Set(selectedIds);
    if (next.has(card.id)) {
      next.delete(card.id);
    } else {
      next.add(card.id);
    }
    selectedIds = next;
  }

  function selectAll() {
    if (selectedIds.size === cards.length) {
      selectedIds = new Set();
    } else {
      selectedIds = new Set(cards.map((c) => c.id));
    }
  }

  async function bulkDelete() {
    const ids = [...selectedIds];
    for (const id of ids) {
      await api.deleteCard(id);
      printQueue.removeAll(id);
    }
    trashCount.set(trashCount.value + ids.length);
    toast.success(`${ids.length} kartu dipindahkan ke sampah`);
    selectedIds = new Set();
    showBulkDeleteConfirm = false;
    loadCards();
  }

  let selectedCount = $derived(selectedIds.size);
  let allSelected = $derived(
    cards.length > 0 && selectedCount === cards.length,
  );

  onMount(async () => {
    try {
      const res = await fetch("http://localhost:8080/me", {
        credentials: "include",
      });
      const data = await res.json();
      if (data.authenticated && data.user && data.user.role === "teacher") {
        isTeacher = true;
      }
    } catch (e) {
      console.error(e);
    }
    loadCards();
    trashCount.init();
  });

  function changePage(n: number) {
    loadCards({
      page: n,
      ...(searchQuery.trim() ? { search: searchQuery.trim() } : {}),
    });
  }

  function applySearch() {
    loadCards({
      ...(searchQuery.trim() ? { search: searchQuery.trim() } : {}),
    });
  }

  function layoutLabel(card: Card): string {
    if (card.cardType === "image") return "Gambar";
    return `${card.size || "6"}/A4`;
  }

  function getCardAspect(size: string | undefined): string {
    const s = parseInt(size || "6") || 6;
    if (s === 2) return "190 / 137";
    if (s === 4) return "93.5 / 137";
    return "93.5 / 90.33";
  }

  function getCardPrintWidth(size: string | undefined): number {
    const s = parseInt(size || "6") || 6;
    return s === 2 ? 190 : 93.5;
  }

  function getCardPrintHeight(size: string | undefined): number {
    const s = parseInt(size || "6") || 6;
    if (s === 2 || s === 4) return 137;
    return 90.33;
  }

  function getModalMaxWidth(size: string | undefined): string {
    const s = parseInt(size || "6") || 6;
    if (s === 2) return "max-w-[190mm]";
    return "max-w-[93.5mm]";
  }

  function handleEdit(card: Card) {
    editingCard = card;
    if (card.cardType === "image") {
      showImageForm = true;
    } else {
      showForm = true;
    }
  }

  async function handleAdd(card: Omit<Card, "id">) {
    await api.createCard(card);
    showForm = false;
    toast.success("Kartu berhasil dibuat");
    loadCards();
  }

  async function handleUpdate(card: Omit<Card, "id">) {
    if (editingCard) {
      await api.updateCard(editingCard.id, card);
      editingCard = null;
      showForm = false;
      showImageForm = false;
      toast.success("Kartu berhasil diupdate");
      loadCards();
    }
  }

  async function handleAddImage(card: Omit<Card, "id">) {
    await api.createCard(card);
    showImageForm = false;
    toast.success("Kartu gambar berhasil dibuat");
    loadCards();
  }

  function handleDeleteClick(card: Card) {
    deleteCardId = card.id;
  }

  function handleDetailClick(card: Card) {
    detailCard = card;
  }

  async function confirmDelete() {
    if (!deleteCardId) return;
    await api.deleteCard(deleteCardId);
    printQueue.removeAll(deleteCardId);
    trashCount.inc();
    if (editingCard?.id === deleteCardId) {
      editingCard = null;
      showForm = false;
    }
    cards = cards.filter((c) => c.id !== deleteCardId);
    total = total - 1;
    toast.success("Kartu dipindahkan ke sampah");
    deleteCardId = null;
  }
</script>

<svelte:head>
  <title>Arsip</title>
</svelte:head>

<div class="min-h-screen bg-transparent text-orange-950">
  <header
    class="fixed top-16 md:top-0 left-0 md:left-64 right-0 bg-white/80/90 backdrop-blur-md px-6 md:px-10 py-3 z-40 transition-all"
  >
    <div
      class="max-w-5xl mx-auto flex items-center justify-between gap-3 flex-wrap"
    >
      <div class="flex items-center gap-4">
        <nav class="flex items-center gap-1 bg-white/80 p-1 rounded-xl">
          <span
            class="px-4 py-1.5 text-sm rounded-lg bg-white/80 text-orange-950 shadow-sm font-base transition-all"
          >
            Arsip
          </span>
          <a
            href="/dashboard/card-memory/print"
            class="px-4 py-1.5 text-sm rounded-lg text-orange-800 hover:text-orange-950 hover:bg-white/80/60 font-base flex items-center gap-1.5 transition-all"
          >
            Print
            {#if queueCount > 0}
              <span
                class="bg-orange-500 text-white text-[10px] rounded-full px-1.5 py-0.5 leading-none"
                >{queueCount}</span
              >
            {/if}
          </a>
          {#if isTeacher}
            <a
              href="/dashboard/card-memory/trash"
              class="px-4 py-1.5 text-sm rounded-lg text-orange-800 hover:text-orange-950 hover:bg-white/80/60 font-base flex items-center gap-1.5 transition-all"
            >
              Sampah
              {#if trashCountVal > 0}
                <span
                  class="bg-red-500 text-white text-[10px] rounded-full px-1.5 py-0.5 leading-none"
                  >{trashCountVal}</span
                >
              {/if}
            </a>
          {/if}
        </nav>
      </div>

      <div class="flex items-center gap-3 flex-wrap">
        <div class="relative">
          <input
            type="text"
            bind:value={searchQuery}
            placeholder="Cari kartu..."
            class="px-4 py-1.5 pr-8 text-sm bg-white/80 text-orange-950 border border-orange-200 rounded-xl focus:ring-2 focus:ring-gray-200 focus:bg-white/80 focus:border-transparent outline-none w-56 transition-all"
            onkeydown={(e) => {
              if (e.key === "Enter") applySearch();
            }}
          />
          {#if searchQuery}
            <button
              onclick={() => {
                searchQuery = "";
                applySearch();
              }}
              class="absolute right-8 top-1/2 -translate-y-1/2 text-orange-700 hover:text-orange-800 cursor-pointer text-base"
              >&times;</button
            >
          {/if}
          <button
            onclick={applySearch}
            class="absolute right-2.5 top-1/2 -translate-y-1/2 text-orange-700 hover:text-orange-950 cursor-pointer"
            title="Cari"
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
                d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
              />
            </svg>
          </button>
        </div>
        {#if isTeacher}
          <div class="flex items-center gap-2">
            <button
              onclick={() => {
                editingCard = null;
                showForm = !showForm;
              }}
              class="flex items-center gap-1.5 px-4 py-1.5 text-sm rounded-xl bg-orange-500 text-white hover:bg-orange-600 font-medium transition-all shadow-sm cursor-pointer"
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
                  d="M12 4v16m8-8H4"
                />
              </svg>
              {showForm ? "Tutup Form" : "Kartu"}
            </button>
            <button
              onclick={() => {
                editingCard = null;
                showImageForm = !showImageForm;
              }}
              class="flex items-center gap-1.5 px-4 py-1.5 text-sm rounded-xl bg-indigo-500 border border-indigo-500 hover:bg-indigo-600 text-white font-medium transition-all shadow-sm cursor-pointer"
            >
              <svg
                class="w-4 h-4"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                viewBox="0 0 24 24"
              >
                <rect x="3" y="3" width="18" height="18" rx="2" ry="2" /><circle
                  cx="8.5"
                  cy="8.5"
                  r="1.5"
                /><polyline points="21 15 16 10 5 21" />
              </svg>
              {showImageForm ? "Tutup" : "Kartu Gambar"}
            </button>
          </div>
        {/if}
      </div>
    </div>
  </header>

  <main class="max-w-5xl mx-auto p-4 pt-24 md:pt-20">
    <!-- Bulk selection bar -->
    {#if cards.length > 0 && isTeacher}
      <div class="flex items-center gap-3 mb-3">
        <button
          onclick={selectAll}
          class="flex items-center gap-1.5 text-xs text-orange-800 hover:text-orange-900 cursor-pointer"
        >
          <span
            class="w-5 h-5 flex items-center justify-center rounded border-2 {allSelected
              ? 'bg-orange-500 border-orange-500'
              : 'border-orange-300'}"
          >
            {#if allSelected}
              <svg
                class="w-3.5 h-3.5 text-orange-950"
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
          </span>
          {allSelected ? "Batal Pilih Semua" : "Pilih Semua"} ({cards.length})
        </button>
        {#if selectedCount > 0}
          <span class="text-xs text-orange-700">|</span>
          <span class="text-xs text-orange-800">{selectedCount} dipilih</span>
          <button
            onclick={() => (showBulkDeleteConfirm = true)}
            class="px-3 py-1 text-xs rounded-lg bg-red-500 text-white hover:bg-red-700 cursor-pointer"
          >
            Hapus ({selectedCount})
          </button>
        {/if}
      </div>
    {/if}

    <Modal
      show={showForm}
      onclose={() => {
        showForm = false;
        editingCard = null;
      }}
      title={editingCard ? "Edit Kartu" : "Tambah Kartu Baru"}
    >
      <CardForm
        onsave={editingCard ? handleUpdate : handleAdd}
        oncancel={() => {
          showForm = false;
          editingCard = null;
        }}
        edit={editingCard}
      />
    </Modal>

    <Modal
      show={showImageForm}
      onclose={() => {
        showImageForm = false;
        editingCard = null;
      }}
      title={editingCard ? "Edit Kartu Gambar" : "Tambah Kartu Gambar"}
    >
      <CardImageForm
        onsave={editingCard ? handleUpdate : handleAddImage}
        oncancel={() => {
          showImageForm = false;
          editingCard = null;
        }}
        edit={editingCard}
      />
    </Modal>

    <div>
      <h2 class="text-md font-semibold text-orange-900 mb-3">
        Arsip ({total} kartu)
        {#if totalPages > 1}
          <span class="text-orange-700"> · hlm {currentPage}/{totalPages}</span>
        {/if}
      </h2>

      {#if loading}
        <div
          class="text-center py-8 bg-white/80 rounded-xl border border-orange-200"
        >
          <p class="text-orange-700">Memuat data dari server...</p>
        </div>
      {:else if error}
        <div
          class="text-center py-8 bg-white/80 rounded-xl border border-red-300"
        >
          <p class="text-red-600 mb-3">{error}</p>
          <button
            onclick={() => loadCards()}
            class="px-4 py-2 text-sm rounded-lg bg-white text-orange-800 hover:bg-orange-50 cursor-pointer"
            >Coba Lagi</button
          >
        </div>
      {:else if cards.length === 0}
        <div
          class="text-center py-8 bg-white/80 rounded-xl border border-orange-200"
        >
          <p class="text-orange-700 mb-3">
            Belum ada kartu. Tambah kartu baru via form di atas.
          </p>
        </div>
      {:else}
        <div
          class="bg-white/80 rounded-xl border border-orange-200 overflow-hidden"
        >
          <table class="w-full text-sm">
            <thead>
              <tr
                class="bg-transparent border-b border-orange-200 text-left text-xs text-orange-900 uppercase tracking-wider"
              >
                {#if isTeacher}<th class="px-4 py-3 w-8"></th>{/if}
                <th class="px-4 py-3">Title</th>
                <th class="px-4 py-3">Kategori</th>
                <th class="px-4 py-3">Layout Print</th>
                <th class="px-4 py-3 text-right">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {#each cards as card (card.id)}
                <tr
                  class="border-b border-orange-200 hover:bg-transparent text-orange-950 transition-colors"
                >
                  {#if isTeacher}
                    <td class="px-4 py-3">
                      <button
                        onclick={() => toggleSelect(card)}
                        class="w-5 h-5 flex items-center justify-center rounded border-2 cursor-pointer {selectedIds.has(
                          card.id,
                        )
                          ? 'bg-orange-500 border-orange-500'
                          : 'border-orange-300'}"
                      >
                        {#if selectedIds.has(card.id)}
                          <svg
                            class="w-3.5 h-3.5 text-orange-950"
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
                  <td class="px-4 py-3 font-semibold text-orange-900 text-md"
                    >{card.title}</td
                  >
                  <td class="px-4 py-3 text-orange-800"
                    >{card.category || "-"}</td
                  >
                  <td class="px-4 py-3">
                    <span
                      class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium {card.cardType ===
                      'image'
                        ? 'bg-purple-100 text-purple-800'
                        : 'bg-orange-100 text-orange-800'}"
                    >
                      {layoutLabel(card)}
                    </span>
                  </td>
                  <td class="px-4 py-3 text-right">
                    <div class="flex items-center justify-end gap-1">
                      <button
                        onclick={() => handleDetailClick(card)}
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
                          ? 'bg-orange-200 text-orange-800'
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
                          onclick={() => handleEdit(card)}
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
                          onclick={() => handleDeleteClick(card)}
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

        {#if totalPages > 1}
          <div class="flex items-center justify-center gap-1 mt-6">
            <button
              onclick={() => changePage(1)}
              disabled={currentPage === 1}
              class="px-2 py-1 text-xs rounded border border-orange-300 text-orange-800 hover:bg-white disabled:opacity-30 disabled:cursor-default cursor-pointer"
              >&laquo;</button
            >
            <button
              onclick={() => changePage(currentPage - 1)}
              disabled={currentPage === 1}
              class="px-2 py-1 text-xs rounded border border-orange-300 text-orange-800 hover:bg-white disabled:opacity-30 disabled:cursor-default cursor-pointer"
              >&lsaquo;</button
            >

            {#each Array(totalPages) as _, i}
              {@const p = i + 1}
              {#if p === 1 || p === totalPages || (p >= currentPage - 2 && p <= currentPage + 2)}
                <button
                  onclick={() => changePage(p)}
                  class="px-2 py-1 text-xs rounded border cursor-pointer {p ===
                  currentPage
                    ? 'bg-orange-500 text-white border-indigo-600'
                    : 'border-orange-300 text-orange-800 hover:bg-white'}"
                  >{p}</button
                >
              {:else if p === currentPage - 3 || p === currentPage + 3}
                <span class="px-1 text-xs text-orange-700">...</span>
              {/if}
            {/each}

            <button
              onclick={() => changePage(currentPage + 1)}
              disabled={currentPage === totalPages}
              class="px-2 py-1 text-xs rounded border border-orange-300 text-orange-800 hover:bg-white disabled:opacity-30 disabled:cursor-default cursor-pointer"
              >&rsaquo;</button
            >
            <button
              onclick={() => changePage(totalPages)}
              disabled={currentPage === totalPages}
              class="px-2 py-1 text-xs rounded border border-orange-300 text-orange-800 hover:bg-white disabled:opacity-30 disabled:cursor-default cursor-pointer"
              >&raquo;</button
            >
          </div>
        {/if}
      {/if}
    </div>
  </main>

  <Modal show={!!deleteCardId} onclose={() => (deleteCardId = null)}>
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
        <h3 class="text-lg font-semibold text-orange-950">Hapus Kartu</h3>
        {#if deleteTarget}
          <p class="text-sm font-medium text-orange-900 mt-1">
            "{deleteTarget.title}"
          </p>
        {/if}
        <p class="text-sm text-orange-800 mt-1">
          Kartu akan dipindahkan ke tempat sampah. Kamu bisa memulihkannya
          nanti.
        </p>
      </div>
      <div class="flex gap-2 justify-center pt-2">
        <button
          onclick={() => (deleteCardId = null)}
          class="px-4 py-2 text-sm rounded-lg border border-orange-300 hover:bg-transparent text-orange-950 cursor-pointer"
          >Batal</button
        >
        <button
          onclick={confirmDelete}
          class="px-4 py-2 text-sm rounded-lg bg-red-500 text-white hover:bg-red-700 cursor-pointer"
          >Hapus</button
        >
      </div>
    </div>
  </Modal>

  <Modal
    show={showBulkDeleteConfirm}
    onclose={() => (showBulkDeleteConfirm = false)}
  >
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
        <h3 class="text-lg font-semibold text-orange-950">
          Hapus {selectedCount} Kartu
        </h3>
        <p class="text-sm text-orange-800 mt-1">
          {selectedCount} kartu akan dipindahkan ke tempat sampah.
        </p>
      </div>
      <div class="flex gap-2 justify-center pt-2">
        <button
          onclick={() => (showBulkDeleteConfirm = false)}
          class="px-4 py-2 text-sm rounded-lg border border-orange-300 hover:bg-transparent text-orange-950 cursor-pointer"
          >Batal</button
        >
        <button
          onclick={bulkDelete}
          class="px-4 py-2 text-sm rounded-lg bg-red-500 text-white hover:bg-red-700 cursor-pointer"
          >Hapus {selectedCount} Kartu</button
        >
      </div>
    </div>
  </Modal>

  <Modal
    show={!!detailCard}
    onclose={() => {
      detailCard = null;
      cardWrapperWidth = 0;
    }}
  >
    {#if detailCard}
      <div
        class="space-y-5 w-full {getModalMaxWidth(
          detailCard.size,
        )} mx-auto pb-2 transition-all duration-300"
      >
        <div
          class="w-full mx-auto shadow-lg rounded-xl overflow-hidden border border-orange-200 transition-all duration-300 bg-white/80"
          style="aspect-ratio: {getCardAspect(detailCard.size)};"
          bind:clientWidth={cardWrapperWidth}
        >
          {#if cardWrapperWidth > 0}
            <div
              style="width: {getCardPrintWidth(
                detailCard.size,
              )}mm; height: {getCardPrintHeight(
                detailCard.size,
              )}mm; transform: scale({cardWrapperWidth /
                (getCardPrintWidth(detailCard.size) *
                  3.7795)}); transform-origin: top left;"
            >
              <CardItem card={detailCard} index={0} />
            </div>
          {/if}
        </div>

        <div class="flex items-center justify-center gap-4">
          <button
            onclick={() => printQueue.toggle(detailCard!)}
            class="w-12 h-12 flex items-center justify-center rounded-full transition-all cursor-pointer shadow-sm border {printQueue.has(
              detailCard.id,
            )
              ? 'bg-indigo-900/30 border-indigo-800/50 text-orange-700 hover:bg-indigo-900/50 shadow-indigo-100'
              : 'bg-white/80 border-orange-200 text-orange-800 hover:text-orange-700 hover:border-indigo-800/50 hover:bg-indigo-900/30'}"
            title={printQueue.has(detailCard.id)
              ? "Hapus dari Print"
              : "Tambah ke Print"}
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

          {#if isTeacher}
            <button
              onclick={() => {
                const card = detailCard;
                detailCard = null;
                cardWrapperWidth = 0;
                handleEdit(card!);
              }}
              class="w-12 h-12 flex items-center justify-center rounded-full bg-white/80 border border-orange-200 text-orange-800 hover:text-amber-600 hover:border-amber-300 hover:bg-amber-50 transition-all cursor-pointer shadow-sm"
              title="Edit Kartu"
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
                  d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                />
              </svg>
            </button>

            <button
              onclick={() => {
                const card = detailCard;
                detailCard = null;
                cardWrapperWidth = 0;
                handleDeleteClick(card!);
              }}
              class="w-12 h-12 flex items-center justify-center rounded-full bg-white/80 border border-orange-200 text-orange-800 hover:text-red-600 hover:border-red-300 hover:bg-red-100 transition-all cursor-pointer shadow-sm"
              title="Hapus Kartu"
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
                  d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                />
              </svg>
            </button>
          {/if}

          <button
            onclick={() => {
              detailCard = null;
              cardWrapperWidth = 0;
            }}
            class="w-12 h-12 flex items-center justify-center rounded-full bg-gray-800 border border-gray-800 text-orange-950 hover:bg-indigo-600 hover:border-gray-900 transition-all cursor-pointer shadow-md hover:shadow-lg"
            title="Tutup"
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
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>
      </div>
    {/if}
  </Modal>
</div>
