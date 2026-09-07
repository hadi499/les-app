<script lang="ts">
  import { onMount } from "svelte";

  type Quote = { id: number; quote: string; arti: string; author: string; is_published: boolean };

  let quotes: Quote[] = $state([]);
  let isLoading = $state(true);
  let showLoadingSpinner = $state(false);
  let errorMsg = $state("");

  // Modal State
  let showModal = $state(false);
  let isEditing = $state(false);
  let currentQuoteId: number | null = $state(null);

  let showDeleteModal = $state(false);
  let quoteToDelete: number | null = $state(null);

  // Form State
  let formQuote = $state("");
  let formArti = $state("");
  let formAuthor = $state("");
  let formIsPublished = $state(true);

  // Computed state for counting published quotes
  let publishedCount = $derived(quotes.filter(q => q.is_published).length);

  async function fetchQuotes() {
    isLoading = true;
    showLoadingSpinner = false;
    setTimeout(() => { showLoadingSpinner = true; }, 150);
    errorMsg = "";
    try {
      const res = await fetch(`/api/quotes/all`, {
        credentials: "include",
      });
      if (!res.ok) throw new Error("Gagal mengambil data quotes");
      quotes = ((await res.json()) as Quote[]) || [];
    } catch (e) {
      errorMsg = e instanceof Error ? e.message : String(e);
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    fetchQuotes();
  });

  function openAddModal() {
    isEditing = false;
    currentQuoteId = null;
    formQuote = "";
    formArti = "";
    formAuthor = "";
    formIsPublished = false;
    showModal = true;
  }

  function openEditModal(quote: Quote) {
    isEditing = true;
    currentQuoteId = quote.id;
    formQuote = quote.quote;
    formArti = quote.arti;
    formAuthor = quote.author;
    formIsPublished = quote.is_published;
    showModal = true;
  }

  function closeModal() {
    showModal = false;
  }

  async function saveQuote(e: Event) {
    e.preventDefault();
    if (!formQuote.trim() || !formArti.trim() || !formAuthor.trim()) {
      alert("Semua kolom harus diisi!");
      return;
    }

    const payload = {
      quote: formQuote.trim(),
      arti: formArti.trim(),
      author: formAuthor.trim(),
      is_published: formIsPublished,
    };

    try {
      let url = `/api/quotes`;
      let method = "POST";

      if (isEditing) {
        url = `/api/quotes/${currentQuoteId}`;
        method = "PUT";
      }

      const res = await fetch(url, {
        method,
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
        credentials: "include",
      });

      if (!res.ok) {
        const errData = await res.json();
        throw new Error(errData.error || "Gagal menyimpan data");
      }

      closeModal();
      await fetchQuotes();
    } catch (err) {
      alert(err instanceof Error ? err.message : String(err));
    }
  }

  function openDeleteModal(id: number) {
    quoteToDelete = id;
    showDeleteModal = true;
  }

  function closeDeleteModal() {
    showDeleteModal = false;
    quoteToDelete = null;
  }

  async function executeDelete() {
    if (quoteToDelete === null) return;

    try {
      const res = await fetch(`/api/quotes/${quoteToDelete}`, {
        method: "DELETE",
        credentials: "include",
      });

      if (!res.ok) {
        const errData = await res.json();
        throw new Error(errData.error || "Gagal menghapus data");
      }
      await fetchQuotes();
      closeDeleteModal();
    } catch (err) {
      alert(err instanceof Error ? err.message : String(err));
      closeDeleteModal();
    }
  }
</script>

<svelte:head>
  <title>Manajemen Quotes - Portal Guru</title>
</svelte:head>

<div class="animate-in fade-in duration-500 max-w-5xl mx-auto space-y-6">
  <div
    class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 mb-8"
  >
    <div>
      <h1
        class="text-2xl font-bold text-slate-900 sm:text-3xl tracking-tight drop-shadow-sm"
      >
        Quotes Inspirasi
      </h1>
      <p
        class="mt-2 text-slate-600 text-sm sm:text-base font-light tracking-wide"
      >
        Kelola quotes yang akan ditampilkan di halaman utama.
      </p>
    </div>
    <button
      onclick={openAddModal}
      class="inline-flex items-center self-start sm:self-auto gap-2 px-4 py-2.5 text-sm font-bold text-indigo-700 bg-indigo-50 hover:bg-indigo-100 border border-indigo-200 rounded-xl transition-all shadow-sm cursor-pointer"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"
        ><path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M12 4v16m8-8H4"
        ></path></svg
      >
      Quote Baru
    </button>
  </div>

  {#if isLoading}
    <div class="fixed inset-0 z-[100] flex flex-col items-center justify-center bg-slate-50/50 backdrop-blur-sm {showLoadingSpinner ? 'opacity-100' : 'opacity-0'} transition-opacity duration-300">
      <div class="w-12 h-12 border-4 border-slate-200 border-t-indigo-600 rounded-full animate-spin shadow-sm"></div>
    </div>
  {:else if errorMsg}
    <div
      class="bg-red-100 text-red-800 p-6 rounded-2xl border border-red-300 font-medium flex items-center gap-3"
    >
      <svg
        class="w-6 h-6 flex-shrink-0"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
        ><path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
        ></path></svg
      >
      {errorMsg}
    </div>
  {:else}
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      {#each quotes as quote}
        <div
          class="bg-white/60 backdrop-blur-md rounded-2xl border border-slate-200 p-5 shadow-lg shadow-slate-800/5 flex flex-col justify-between hover:border-slate-300 transition-colors gap-4"
        >
          <div>
            <div class="flex justify-between items-start mb-2 gap-2">
              <p class="text-slate-900 font-semibold text-base leading-relaxed italic">"{quote.quote}"</p>
              {#if quote.is_published}
                <span class="inline-flex shrink-0 items-center px-2 py-0.5 rounded text-[10px] font-bold bg-green-100 text-green-700">Publik</span>
              {:else}
                <span class="inline-flex shrink-0 items-center px-2 py-0.5 rounded text-[10px] font-bold bg-slate-100 text-slate-500">Draft</span>
              {/if}
            </div>
            <p class="text-slate-500 text-sm leading-relaxed border-t border-slate-200 pt-3 mt-1 mb-2">{quote.arti}</p>
            <div class="flex items-center gap-2">
              <div class="w-4 h-0.5 rounded-full bg-indigo-400"></div>
              <p class="text-slate-700 font-bold text-xs tracking-wider uppercase">{quote.author}</p>
            </div>
          </div>

          <div class="flex justify-end gap-2 pt-3 border-t border-slate-200">
            <button
              onclick={() => openEditModal(quote)}
              class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-blue-600 bg-blue-100 hover:bg-blue-500/10 rounded-lg transition-colors border border-blue-300"
            >
              Edit
            </button>
            <button
              onclick={() => openDeleteModal(quote.id)}
              class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-red-600 bg-red-100 hover:bg-red-200 rounded-lg transition-colors border border-red-300"
            >
              Hapus
            </button>
          </div>
        </div>
      {/each}
      {#if quotes.length === 0}
        <div
          class="col-span-full py-12 text-center text-slate-500 font-light bg-[#EAE4BD]/30 rounded-3xl border border-slate-200"
        >
          Belum ada quotes yang ditambahkan.
        </div>
      {/if}
    </div>
  {/if}
</div>

<!-- Modal Dialog -->
{#if showModal}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
    <div
      class="absolute inset-0 bg-black/60 backdrop-blur-sm"
      onclick={closeModal}
    ></div>
    <div
      class="relative bg-slate-50 border border-slate-300 rounded-2xl shadow-2xl w-full max-w-lg overflow-hidden animate-in zoom-in-95 duration-200 max-h-[90vh] flex flex-col"
    >
      <div
        class="p-6 border-b border-slate-200 flex justify-between items-center"
      >
        <h3 class="text-xl font-bold text-slate-900">
          {isEditing ? "Edit Quote" : "Tambah Quote"}
        </h3>
        <button
          onclick={closeModal}
          class="text-slate-500 hover:text-slate-900 transition-colors"
        >
          <svg
            class="w-6 h-6"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            ><path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M6 18L18 6M6 6l12 12"
            ></path></svg
          >
        </button>
      </div>

      <div class="overflow-y-auto">
        <form onsubmit={saveQuote} class="p-6 space-y-4">
          <div>
            <label
              class="block text-sm font-medium text-slate-600 mb-1.5"
              for="quote">Quote (Inggris / Asli)</label
            >
            <textarea
              id="quote"
              bind:value={formQuote}
              placeholder="Masukkan quote..."
              class="w-full bg-white border border-slate-300 rounded-xl px-4 py-2.5 text-slate-900 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all min-h-[100px]"
              required
            ></textarea>
          </div>

          <div>
            <label
              class="block text-sm font-medium text-slate-600 mb-1.5"
              for="arti">Arti (Indonesia)</label
            >
            <textarea
              id="arti"
              bind:value={formArti}
              placeholder="Masukkan arti dari quote..."
              class="w-full bg-white border border-slate-300 rounded-xl px-4 py-2.5 text-slate-900 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all min-h-[100px]"
              required
            ></textarea>
          </div>

          <div>
            <label
              class="block text-sm font-medium text-slate-600 mb-1.5"
              for="author">Penulis / Tokoh</label
            >
            <input
              type="text"
              id="author"
              bind:value={formAuthor}
              placeholder="Misal: Albert Einstein"
              class="w-full bg-white border border-slate-300 rounded-xl px-4 py-2.5 text-slate-900 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all"
              required
            />
          </div>

          <div class="flex items-center mt-2">
            <input
              type="checkbox"
              id="isPublished"
              bind:checked={formIsPublished}
              disabled={(!isEditing && publishedCount >= 1) || (isEditing && !quotes.find(q => q.id === currentQuoteId)?.is_published && publishedCount >= 1)}
              class="w-4 h-4 text-indigo-600 bg-white border-slate-300 rounded focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
            />
            <label for="isPublished" class="ml-2 text-sm font-medium text-slate-700 {((!isEditing && publishedCount >= 1) || (isEditing && !quotes.find(q => q.id === currentQuoteId)?.is_published && publishedCount >= 1)) ? 'opacity-50' : ''}">
              Tampilkan di Beranda (Publish)
            </label>
          </div>
          {#if (!isEditing && publishedCount >= 1) || (isEditing && !quotes.find(q => q.id === currentQuoteId)?.is_published && publishedCount >= 1)}
            <p class="text-xs text-amber-600 mt-1">Batas maksimal 1 quote publish telah tercapai. Quote ini akan disimpan sebagai Draft.</p>
          {/if}

          <div class="pt-4 flex justify-end gap-3">
            <button
              type="button"
              onclick={closeModal}
              class="px-4 py-2.5 text-sm font-medium text-slate-800 hover:text-slate-900 bg-white hover:bg-slate-100 rounded-xl transition-colors shadow-md"
            >
              Batal
            </button>
            <button
              type="submit"
              class="px-4 py-2.5 text-sm font-medium text-slate-900 bg-indigo-100 hover:bg-indigo-200 rounded-xl transition-all shadow-md shadow-indigo-900/20"
            >
              {isEditing ? "Simpan" : "Tambahkan"}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
{/if}

<!-- Delete Confirmation Modal -->
{#if showDeleteModal}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
    <div
      class="absolute inset-0 bg-black/60 backdrop-blur-sm"
      onclick={closeDeleteModal}
    ></div>
    <div
      class="relative bg-slate-50 border border-red-300 rounded-2xl shadow-2xl w-full max-w-md overflow-hidden animate-in zoom-in-95 duration-200 p-6 text-center"
    >
      <div
        class="w-12 h-12 mx-auto rounded-full bg-red-100 flex items-center justify-center mb-4"
      >
        <svg
          class="w-6 h-6 text-red-600"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
          />
        </svg>
      </div>
      <h3 class="text-xl font-bold text-slate-900 mb-2">
        Hapus Quote?
      </h3>
      <p class="text-sm text-slate-600 mb-6">
        Apakah Anda yakin ingin menghapus quote ini? Tindakan ini tidak dapat dibatalkan.
      </p>
      <div class="flex justify-center gap-3">
        <button
          onclick={closeDeleteModal}
          class="px-4 py-2.5 text-sm font-medium text-slate-800 bg-white border border-slate-300 rounded-xl hover:bg-slate-50 transition-colors"
        >
          Batal
        </button>
        <button
          onclick={executeDelete}
          class="px-4 py-2.5 text-sm font-medium text-white bg-red-600 hover:bg-red-700 rounded-xl transition-all shadow-md shadow-red-900/20"
        >
          Ya, Hapus
        </button>
      </div>
    </div>
  </div>
{/if}
