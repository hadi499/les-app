<script module>
  // Cache sederhana agar tidak ada kedipan (flicker) saat kembali ke halaman ini
  let cachedCards: any[] = [];
  let cachedTotal = 0;
</script>

<script lang="ts">
  import { onMount } from "svelte";
  import type { Card } from "$lib/types";
  import { printQueue } from "$lib/stores/print-queue.svelte";
  import * as api from "$lib/api";
  import CardForm from "$lib/components/CardForm.svelte";
  import CardImageForm from "$lib/components/CardImageForm.svelte";
  import Modal from "$lib/components/Modal.svelte";
  import CardItem from "$lib/components/CardItem.svelte";
  import FolderGrid from "$lib/components/FolderGrid.svelte";
  import CardTable from "$lib/components/CardTable.svelte";
  import BulkActionBar from "$lib/components/BulkActionBar.svelte";
  import { toast } from "$lib/stores/toast.svelte";
  import { trashCount } from "$lib/stores/trash-count.svelte";
  import { cardFolders, updateCardFolder, deleteCardFolder, fetchCardFolders, createCardFolder, type CardFolder } from "$lib/stores/cardFolders";

  let showForm = $state(false);
  let isTeacher = $state(true);
  let showImageForm = $state(false);
  let editingCard = $state<Card | null>(null);
  let deleteCardId = $state<string | null>(null);
  let detailCard = $state<Card | null>(null);
  let activeCategory = $state<string | null>(null);
  let searchQuery = $state("");
  let cardWrapperWidth = $state(0);

  let cards = $state<Card[]>(cachedCards);
  let total = $state(cachedTotal);
  let loading = $state(cachedCards.length === 0);
  let error = $state("");

  let showFolderForm = $state(false);
  let editingFolder = $state<CardFolder | null>(null);
  let newFolderName = $state("");
  let showDeleteFolderModal = $state(false);
  let folderToDelete = $state<CardFolder | null>(null);

  function openNewFolder() {
    editingFolder = null;
    newFolderName = "";
    showFolderForm = true;
  }

  function handleEditFolder(folderName: string) {
    const cf = $cardFolders.find(f => f.name === folderName);
    if (cf) {
      editingFolder = cf;
      newFolderName = cf.name;
      showFolderForm = true;
    }
  }

  function handleDeleteFolderClick(folderName: string) {
    const cf = $cardFolders.find(f => f.name === folderName);
    if (cf) {
      folderToDelete = cf;
      showDeleteFolderModal = true;
    }
  }

  async function saveFolder(e: Event) {
    e.preventDefault();
    const name = newFolderName.trim();
    if (!name) return;

    if (editingFolder) {
      const ok = await updateCardFolder(editingFolder.id, name);
      if (ok) {
        toast.success("Folder berhasil diupdate");
        loadCards();
      } else {
        toast.error("Gagal mengupdate folder");
      }
    } else {
      const ok = await createCardFolder(name);
      if (ok) {
        toast.success("Folder berhasil dibuat");
        loadCards();
      } else {
        toast.error("Gagal membuat folder");
      }
    }
    showFolderForm = false;
  }

  async function confirmDeleteFolder() {
    if (!folderToDelete) return;
    const ok = await deleteCardFolder(folderToDelete.id);
    if (ok) {
      toast.success(`Folder "${folderToDelete.name}" berhasil dihapus`);
      loadCards();
    } else {
      toast.error("Gagal menghapus folder");
    }
    showDeleteFolderModal = false;
    folderToDelete = null;
  }

  async function loadCards(params: Record<string, any> = {}) {
    if (cachedCards.length === 0) loading = true;
    error = "";
    try {
      const q: Record<string, any> = {
        page: params.page || 1,
        limit: 10000,
        all: true,
        ...params,
      };
      const res = await api.fetchCards(q);
      cards = res.data;
      total = res.total;
      
      // Update cache
      if (!params.search && !params.category) {
        cachedCards = res.data;
        cachedTotal = res.total;
      }
    } catch (e: any) {
      error = e.message || "Gagal memuat data";
    } finally {
      loading = false;
    }
  }

  let deleteTarget = $derived(cards.find((c) => c.id === deleteCardId));
  let queueCount = $derived(printQueue.count);
  let trashCountVal = $derived(trashCount.value);

  let groupedCards = $derived.by(() => {
    const acc: Record<string, Card[]> = {};
    // Inisialisasi folder kosong
    for (const f of $cardFolders) {
      acc[f.name] = [];
    }
    // Kelompokkan kartu berdasarkan folder
    for (const card of cards) {
      const cat = card.card_folder?.name || 'Tidak Berkategori';
      if (!acc[cat]) acc[cat] = [];
      acc[cat].push(card);
    }
    return acc;
  });

  // Bulk selection
  let selectedIds = $state(new Set<string>());
  let showBulkDeleteConfirm = $state(false);

  let activeFolderId = $derived.by(() => {
    if (!activeCategory) return null;
    const folder = $cardFolders.find(f => f.name === activeCategory);
    return folder ? folder.id : null;
  });

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
    const currentCards = activeCategory !== null ? (groupedCards[activeCategory] || []) : (searchQuery.trim() ? cards : []);
    if (selectedIds.size === currentCards.length) {
      selectedIds = new Set();
    } else {
      selectedIds = new Set(currentCards.map((c) => c.id));
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
  let activeCategoryCards = $derived.by(() => {
    let list = activeCategory !== null ? (groupedCards[activeCategory] || []) : cards;
    const q = searchQuery.trim().toLowerCase();
    if (q) {
      list = list.filter(c => 
        (c.title && c.title.toLowerCase().includes(q)) || 
        (c.card_folder?.name && c.card_folder.name.toLowerCase().includes(q))
      );
    }
    return list;
  });
  let allSelected = $derived(
    activeCategoryCards.length > 0 && selectedCount === activeCategoryCards.length,
  );
  let detailCardIndex = $derived(detailCard ? activeCategoryCards.findIndex((c) => c.id === detailCard!.id) : -1);

  onMount(async () => {
    if (typeof window !== 'undefined') {
      const params = new URLSearchParams(window.location.search);
      const cat = params.get("category");
      if (cat) activeCategory = cat;
    }

    try {
      const res = await fetch(`/me`, {
        credentials: "include",
      });
      const data = await res.json();
      if (data.authenticated && data.user && data.user.role === "teacher") {
        isTeacher = true;
      }
    } catch (e) {
      console.error(e);
    }
    fetchCardFolders();
    loadCards();
    trashCount.init();
  });

  function selectCategory(cat: string | null) {
    activeCategory = cat;
    selectedIds = new Set();
    if (typeof window !== 'undefined') {
      const url = new URL(window.location.href);
      if (cat) {
        url.searchParams.set("category", cat);
      } else {
        url.searchParams.delete("category");
      }
      window.history.replaceState({}, '', url);
    }
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
    // Note: Server limits or folder-based limits should be handled in backend
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

<div class="min-h-screen bg-transparent text-slate-900">
  <header
    class="fixed top-16 md:top-0 left-0 md:left-(--sidebar-width) right-0 bg-white/80/90 backdrop-blur-md px-6 md:px-10 py-3 z-40 transition-all"
  >
    <div
      class="max-w-5xl mx-auto flex flex-col md:flex-row items-stretch md:items-center justify-between gap-4"
    >
      <div class="flex items-center overflow-x-auto hide-scrollbar -mx-2 px-2 md:mx-0 md:px-0">
        <nav class="flex items-center gap-1 bg-white/80 p-1 rounded-xl w-max">
          <span
            class="px-4 py-1.5 text-sm rounded-lg bg-slate-800 text-white shadow-md font-base transition-all"
          >
            Arsip
          </span>
          <a
            href="/dashboard/card-memory/recent"
            class="px-4 py-1.5 text-sm rounded-lg text-slate-600 hover:text-slate-900 hover:bg-white/80/60 font-base transition-all"
          >
            Recent
          </a>
          <a
            href="/dashboard/card-memory/print"
            class="px-4 py-1.5 text-sm rounded-lg text-slate-600 hover:text-slate-900 hover:bg-white/80/60 font-base flex items-center gap-1.5 transition-all"
          >
            Print
            {#if queueCount > 0}
              <span
                class="bg-blue-500 text-white text-[10px] rounded-full px-1.5 py-0.5 leading-none"
                >{queueCount}</span
              >
            {/if}
          </a>
          {#if isTeacher}
            <a
              href="/dashboard/card-memory/trash"
              class="px-4 py-1.5 text-sm rounded-lg text-slate-600 hover:text-slate-900 hover:bg-white/80/60 font-base flex items-center gap-1.5 transition-all"
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

      <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-3">
        <div class="relative w-full sm:w-auto">
          <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none">
            <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
            </svg>
          </div>
          <input
            type="text"
            bind:value={searchQuery}
            placeholder="Cari kartu..."
            class="block w-full pl-10 pr-10 py-1.5 border border-slate-200 rounded-full bg-white/60 focus:bg-white text-sm text-slate-800 placeholder-slate-400 focus:outline-none focus:border-slate-300 focus:ring-4 focus:ring-slate-100 transition-all md:w-56"
          />
          {#if searchQuery}
            <button
              onclick={() => (searchQuery = "")}
              aria-label="Bersihkan pencarian"
              class="absolute right-2 top-1/2 -translate-y-1/2 p-1 text-slate-400 hover:text-slate-600 rounded-full hover:bg-slate-100 transition-colors cursor-pointer"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
              </svg>
            </button>
          {/if}
        </div>
        {#if isTeacher && activeCategory !== null}
          <div class="flex flex-row items-center gap-2 w-full md:w-auto mt-1 md:mt-0">
            <button
              onclick={() => {
                editingCard = null;
                showForm = !showForm;
              }}
              class="flex-1 md:flex-none flex items-center justify-center gap-1.5 px-3 py-2 md:px-4 md:py-1.5 text-xs sm:text-sm rounded-xl bg-blue-500 text-white hover:bg-blue-600 font-medium transition-all shadow-sm cursor-pointer whitespace-nowrap"
            >
              <svg
                class="w-4 h-4 shrink-0"
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
              class="flex-1 md:flex-none flex items-center justify-center gap-1.5 px-3 py-2 md:px-4 md:py-1.5 text-xs sm:text-sm rounded-xl bg-indigo-500 border border-indigo-500 hover:bg-indigo-600 text-white font-medium transition-all shadow-sm cursor-pointer whitespace-nowrap"
            >
              <svg
                class="w-4 h-4 shrink-0"
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
              {showImageForm ? "Tutup Form" : "Kartu Gambar"}
            </button>
          </div>
        {/if}
      </div>
    </div>
  </header>

  <main class="max-w-5xl mx-auto p-2 sm:p-4 {activeCategory !== null ? 'pt-44' : 'pt-32'} md:pt-20">
    <!-- Bulk selection bar -->
    {#if cards.length > 0 && isTeacher && (activeCategory !== null || searchQuery.trim())}
      <BulkActionBar
        allSelected={allSelected}
        selectedCount={selectedCount}
        totalCards={activeCategoryCards.length}
        onSelectAll={selectAll}
        onDeleteBulk={() => (showBulkDeleteConfirm = true)}
      />
    {/if}

    <Modal
      show={showForm}
      onclose={() => {
        showForm = false;
        editingCard = null;
      }}
      title={editingCard ? "Edit Kartu" : "Tambah Kartu Baru"}
    >
      {#if showForm}
        <CardForm
          onsave={editingCard ? handleUpdate : handleAdd}
          oncancel={() => {
            showForm = false;
            editingCard = null;
          }}
          edit={editingCard}
          initialFolderId={activeFolderId}
        />
      {/if}
    </Modal>

    <Modal
      show={showImageForm}
      onclose={() => {
        showImageForm = false;
        editingCard = null;
      }}
      title={editingCard ? "Edit Kartu Gambar" : "Tambah Kartu Gambar"}
    >
      {#if showImageForm}
        <CardImageForm
          onsave={editingCard ? handleUpdate : handleAddImage}
          oncancel={() => {
            showImageForm = false;
            editingCard = null;
          }}
          edit={editingCard}
          initialFolderId={activeFolderId}
        />
      {/if}
    </Modal>

    <div>
      <div class="flex items-center justify-between mb-3 min-h-[40px]">
        <h2 class="text-md font-semibold text-slate-800">
          Arsip ({total} kartu)
        </h2>
        {#if isTeacher && activeCategory === null && !searchQuery.trim()}
          <button
            onclick={openNewFolder}
            class="inline-flex items-center gap-1.5 px-3 md:px-4 py-2 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-xl transition-all shadow-sm cursor-pointer"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 13h6m-3-3v6m-9 1V7a2 2 0 012-2h6l2 2h6a2 2 0 012 2v8a2 2 0 01-2 2H5a2 2 0 01-2-2z"></path></svg>
            Folder Baru
          </button>
        {/if}
      </div>

      {#if loading}
        <div class="text-center py-12 bg-white/80 rounded-xl border border-slate-200">
          <p class="text-slate-500">Memuat kartu...</p>
        </div>
      {:else if error}
        <div
          class="text-center py-8 bg-white/80 rounded-xl border border-red-300"
        >
          <p class="text-red-600 mb-3">{error}</p>
          <button
            onclick={() => loadCards()}
            class="px-4 py-2 text-sm rounded-lg bg-white text-slate-600 hover:bg-slate-50 cursor-pointer"
            >Coba Lagi</button
          >
        </div>
      {:else if cards.length === 0 && !loading}
        <div
          class="text-center py-8 bg-white/80 rounded-xl border border-slate-200"
        >
          <p class="text-slate-500 mb-3">
            Belum ada kartu. Tambah kartu baru via form di atas.
          </p>
        </div>
      {:else}
        {#if activeCategory === null && !searchQuery.trim()}
          <FolderGrid
            groupedCards={groupedCards}
            onSelectCategory={selectCategory}
            isTeacher={isTeacher}
            onEditFolder={handleEditFolder}
            onDeleteFolder={handleDeleteFolderClick}
          />
        {:else}
          {#if activeCategory !== null}
            <div class="mb-4 flex items-center justify-between">
              <button
                onclick={() => selectCategory(null)}
                class="p-1.5 text-slate-500 bg-white/50 hover:text-slate-900 hover:bg-white rounded-lg border border-slate-200 shadow-sm transition-all cursor-pointer"
                title="Kembali"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
                </svg>
              </button>
              <h3 class="font-semibold text-slate-800 text-lg flex items-center gap-2">
                <svg class="w-5 h-5 text-blue-400" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M4 4c-1.11 0-2 .89-2 2v12c0 1.11.89 2 2 2h16c1.11 0 2-.89 2-2V8c0-1.11-.89-2-2-2h-8l-2-2H4z" />
                </svg>
                {activeCategory}
              </h3>
            </div>
          {:else}
            <div class="mb-4">
              <h3 class="font-semibold text-slate-800 text-lg">Hasil Pencarian</h3>
            </div>
          {/if}
          <CardTable
            cards={activeCategoryCards}
            activeCategory={activeCategory}
            isTeacher={isTeacher}
            selectedIds={selectedIds}
            onToggleSelect={toggleSelect}
            onEdit={handleEdit}
            onDelete={handleDeleteClick}
            onDetailClick={handleDetailClick}
          />
        {/if}
      {/if}
    </div>
  </main>

  <Modal show={!!deleteCardId} onclose={() => (deleteCardId = null)} maxWidth="max-w-sm">
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
        <h3 class="text-lg font-semibold text-slate-900">Hapus Kartu</h3>
        {#if deleteTarget}
          <p class="text-sm font-medium text-slate-800 mt-1">
            "{deleteTarget.title}"
          </p>
        {/if}
        <p class="text-sm text-slate-600 mt-1">
          Kartu akan dipindahkan ke tempat sampah. Kamu bisa memulihkannya
          nanti.
        </p>
      </div>
      <div class="flex gap-2 justify-center pt-2">
        <button
          onclick={() => (deleteCardId = null)}
          class="px-4 py-2 text-sm rounded-lg border border-slate-300 hover:bg-transparent text-slate-900 cursor-pointer"
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
    maxWidth="max-w-sm"
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
        <h3 class="text-lg font-semibold text-slate-900">
          Hapus {selectedCount} Kartu
        </h3>
        <p class="text-sm text-slate-600 mt-1">
          {selectedCount} kartu akan dipindahkan ke tempat sampah.
        </p>
      </div>
      <div class="flex gap-2 justify-center pt-2">
        <button
          onclick={() => (showBulkDeleteConfirm = false)}
          class="px-4 py-2 text-sm rounded-lg border border-slate-300 hover:bg-transparent text-slate-900 cursor-pointer"
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
    maxWidth="max-w-[190mm]"
  >
    {#if detailCard}
      <div
        class="space-y-5 w-full max-w-[190mm] mx-auto pb-2 transition-all duration-300 relative group"
      >
        {#if detailCardIndex > 0}
          <button
            onclick={() => {
              detailCard = activeCategoryCards[detailCardIndex - 1];
            }}
            class="absolute top-[40%] -left-3 md:-left-4 -translate-y-1/2 w-10 h-10 md:w-12 md:h-12 flex items-center justify-center rounded-full bg-white/90 border border-slate-200 text-slate-600 hover:text-blue-600 hover:border-blue-300 hover:bg-blue-50 transition-all cursor-pointer shadow-md z-10 opacity-80 group-hover:opacity-100"
            title="Sebelumnya"
          >
            <svg class="w-5 h-5 md:w-6 md:h-6" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" />
            </svg>
          </button>
        {/if}

        {#if detailCardIndex >= 0 && detailCardIndex < activeCategoryCards.length - 1}
          <button
            onclick={() => {
              detailCard = activeCategoryCards[detailCardIndex + 1];
            }}
            class="absolute top-[40%] -right-3 md:-right-4 -translate-y-1/2 w-10 h-10 md:w-12 md:h-12 flex items-center justify-center rounded-full bg-white/90 border border-slate-200 text-slate-600 hover:text-blue-600 hover:border-blue-300 hover:bg-blue-50 transition-all cursor-pointer shadow-md z-10 opacity-80 group-hover:opacity-100"
            title="Selanjutnya"
          >
            <svg class="w-5 h-5 md:w-6 md:h-6" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
            </svg>
          </button>
        {/if}

        <div 
          class="w-full flex items-center justify-center bg-transparent rounded-xl"
          style="aspect-ratio: 190 / 137;"
        >
          <div
            class="w-full {getModalMaxWidth(
              detailCard.size,
            )} mx-auto shadow-lg rounded-xl overflow-hidden border border-slate-200 transition-all duration-300 bg-white"
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
        </div>

        <div class="flex items-center justify-center gap-4 flex-wrap">
          <button
            onclick={() => printQueue.toggle(detailCard!)}
            class="w-12 h-12 flex items-center justify-center rounded-full transition-all cursor-pointer shadow-sm border {printQueue.has(
              detailCard.id,
            )
              ? 'bg-indigo-900/30 border-indigo-800/50 text-slate-500 hover:bg-indigo-900/50 shadow-indigo-100'
              : 'bg-white/80 border-slate-200 text-slate-600 hover:text-slate-500 hover:border-indigo-800/50 hover:bg-indigo-900/30'}"
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
              class="w-12 h-12 flex items-center justify-center rounded-full bg-white/80 border border-slate-200 text-slate-600 hover:text-sky-600 hover:border-sky-300 hover:bg-amber-50 transition-all cursor-pointer shadow-sm"
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
              class="w-12 h-12 flex items-center justify-center rounded-full bg-white/80 border border-slate-200 text-slate-600 hover:text-red-600 hover:border-red-300 hover:bg-red-100 transition-all cursor-pointer shadow-sm"
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
            class="w-12 h-12 flex items-center justify-center rounded-full bg-white/80 border border-slate-200 text-slate-600 hover:text-red-600 hover:border-red-300 hover:bg-red-100 transition-all cursor-pointer shadow-sm"
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

  <!-- Form Folder -->
  <Modal
    show={showFolderForm}
    onclose={() => (showFolderForm = false)}
    title={editingFolder ? "Edit Folder" : "Tambah Folder Baru"}
    maxWidth="max-w-md"
  >
    <form onsubmit={saveFolder} class="space-y-4">
      <div>
        <label
          for="folderName"
          class="block text-sm font-medium text-slate-700 mb-1"
        >
          Nama Folder
        </label>
        <input
          id="folderName"
          type="text"
          bind:value={newFolderName}
          placeholder="Contoh: Matematika"
          class="block w-full px-4 py-2 border border-slate-200 rounded-xl focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-colors"
          required
        />
      </div>
      <div class="flex justify-end gap-3 pt-4 border-t border-slate-100">
        <button
          type="button"
          onclick={() => (showFolderForm = false)}
          class="px-4 py-2 text-sm font-medium text-slate-600 bg-slate-100 hover:bg-slate-200 rounded-xl transition-colors cursor-pointer"
        >
          Batal
        </button>
        <button
          type="submit"
          class="px-6 py-2 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-xl transition-colors shadow-sm cursor-pointer"
          disabled={!newFolderName.trim()}
        >
          Simpan
        </button>
      </div>
    </form>
  </Modal>

  <!-- Konfirmasi Hapus Folder -->
  <Modal
    show={showDeleteFolderModal}
    onclose={() => (showDeleteFolderModal = false)}
    maxWidth="max-w-sm"
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
        <h3 class="text-lg font-semibold text-slate-900">Hapus Folder</h3>
        <p class="text-sm text-slate-600 mt-1">
          Folder "{folderToDelete?.name}" akan dihapus. Kartu di dalamnya mungkin akan kehilangan foldernya.
        </p>
      </div>
      <div class="flex gap-2 justify-center pt-2">
        <button
          onclick={() => (showDeleteFolderModal = false)}
          class="px-4 py-2 text-sm rounded-lg border border-slate-300 hover:bg-transparent text-slate-900 cursor-pointer"
        >
          Batal
        </button>
        <button
          onclick={confirmDeleteFolder}
          class="px-4 py-2 text-sm rounded-lg bg-red-500 text-white hover:bg-red-700 cursor-pointer"
        >
          Hapus
        </button>
      </div>
    </div>
  </Modal>
</div>
