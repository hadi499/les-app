<script lang="ts">
  import { onMount } from "svelte";
  import { 
    fetchFlashcards, createFlashcard, updateFlashcard, deleteFlashcard,
    fetchFlashcardCategories, createFlashcardCategory, updateFlashcardCategory, deleteFlashcardCategory
  } from "$lib/api_english_flashcard";
  import type { EnglishFlashcard, EnglishFlashcardCategory } from "$lib/api_english_flashcard";
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";

  let activeTab = $derived($page.url.searchParams.get("tab") || "categories");

  function switchTab(tab: string) {
    const newUrl = new URL($page.url);
    newUrl.searchParams.set("tab", tab);
    goto(newUrl.toString(), { keepFocus: true, noScroll: true });
  }

  let flashcards = $state<EnglishFlashcard[]>([]);
  let categories = $state<EnglishFlashcardCategory[]>([]);
  let selectedCategoryId = $state<number | "">("");
  
  let currentPage = $state(1);
  const itemsPerPage = 30;
  let paginatedFlashcards = $derived(flashcards.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage));
  let totalPages = $derived(Math.ceil(flashcards.length / itemsPerPage) || 1);
  
  let searchQuery = $state("");
  let searchTimeout: any;

  function handleSearchInput() {
    clearTimeout(searchTimeout);
    searchTimeout = setTimeout(() => {
      loadData();
    }, 400);
  }

  function clearSearch() {
    searchQuery = "";
    loadData();
  }
  
  let isLoading = $state(true);
  
  // Modals state
  let showCardModal = $state(false);
  let isEditingCard = $state(false);
  
  let showCategoryModal = $state(false);
  let isEditingCategory = $state(false);

  let showJsonModal = $state(false);
  let jsonInput = $state("");
  let isImporting = $state(false);
  let jsonCategoryId = $state<number | "">("");

  // Forms
  let currentCard = $state<EnglishFlashcard>({
    category_id: 0,
    question: "",
    options: ["", "", "", ""],
    answer: "",
    explanation: ""
  });
  
  let currentCategory = $state<EnglishFlashcardCategory>({
    name: "",
    description: ""
  });

  async function loadData() {
    isLoading = true;
    try {
      currentPage = 1;
      const [cats, cards] = await Promise.all([
        fetchFlashcardCategories(),
        fetchFlashcards(selectedCategoryId ? (selectedCategoryId as number) : undefined, searchQuery)
      ]);
      categories = cats;
      flashcards = cards;
    } catch (e) {
      console.error(e);
      alert("Gagal memuat data");
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    loadData();
  });

  // --- CATEGORY ACTIONS ---
  
  function openCreateCategoryModal() {
    isEditingCategory = false;
    currentCategory = { name: "", description: "" };
    showCategoryModal = true;
  }
  
  function openEditCategoryModal(cat: EnglishFlashcardCategory) {
    isEditingCategory = true;
    currentCategory.id = cat.id;
    currentCategory.name = cat.name;
    currentCategory.description = cat.description;
    showCategoryModal = true;
  }
  
  async function saveCategory() {
    if (!currentCategory.name) return alert("Nama Kategori wajib diisi");
    
    try {
      if (isEditingCategory && currentCategory.id) {
        await updateFlashcardCategory(currentCategory.id, currentCategory);
      } else {
        await createFlashcardCategory(currentCategory);
      }
      showCategoryModal = false;
      loadData();
    } catch (e) {
      console.error(e);
      alert("Gagal menyimpan kategori");
    }
  }

  async function handleDeleteCategory(id: number | undefined) {
    if (!id) return;
    if (confirm("Yakin menghapus kategori ini? Semua flashcard di dalamnya akan terhapus!")) {
      try {
        await deleteFlashcardCategory(id);
        if (selectedCategoryId === id) selectedCategoryId = "";
        loadData();
      } catch (e) {
        console.error(e);
        alert("Gagal menghapus kategori");
      }
    }
  }

  // --- FLASHCARD ACTIONS ---
  
  function openCreateCardModal() {
    if (!categories.length) {
      switchTab("categories");
      return alert("Silakan buat kategori terlebih dahulu!");
    }
    isEditingCard = false;
    currentCard = {
      category_id: (selectedCategoryId as number) || categories[0].id || 0,
      question: "",
      options: ["", "", "", ""],
      answer: "",
      explanation: ""
    };
    showCardModal = true;
  }

  function openEditCardModal(card: EnglishFlashcard) {
    isEditingCard = true;
    currentCard.id = card.id;
    currentCard.category_id = card.category_id;
    currentCard.question = card.question;
    currentCard.options = [...card.options];
    currentCard.answer = card.answer;
    currentCard.explanation = card.explanation;
    showCardModal = true;
  }

  async function saveCard() {
    if (!currentCard.question || !currentCard.answer || !currentCard.category_id) {
      alert("Mohon lengkapi data wajib (Pertanyaan, Jawaban, Kategori)");
      return;
    }

    try {
      if (isEditingCard && currentCard.id) {
        await updateFlashcard(currentCard.id, currentCard);
      } else {
        await createFlashcard(currentCard);
      }
      showCardModal = false;
      loadData();
    } catch (e) {
      console.error(e);
      alert("Gagal menyimpan kartu");
    }
  }

  async function handleDeleteCard(id: number | undefined) {
    if (!id) return;
    if (confirm("Yakin ingin menghapus kartu ini?")) {
      try {
        await deleteFlashcard(id);
        loadData();
      } catch (e) {
        console.error(e);
        alert("Gagal menghapus data");
      }
    }
  }

  function openJsonModal() {
    if (!categories.length) {
      switchTab("categories");
      return alert("Silakan buat kategori terlebih dahulu!");
    }
    jsonCategoryId = (selectedCategoryId as number) || categories[0].id || "";
    jsonInput = ""; // Dikosongkan agar user bisa langsung paste
    showJsonModal = true;
  }

  async function handleJsonImport() {
    if (!jsonCategoryId) {
      alert("Pilih kategori terlebih dahulu!");
      return;
    }
    try {
      const parsedData = JSON.parse(jsonInput);
      if (!Array.isArray(parsedData)) {
        alert("Format JSON harus berupa array of objects!");
        return;
      }
      
      isImporting = true;
      for (const item of parsedData) {
        if (!item.question || !item.answer || !item.options || item.options.length === 0) {
          console.warn("Skipping invalid item:", item);
          continue;
        }
        
        const newCard: EnglishFlashcard = {
          category_id: jsonCategoryId as number,
          question: item.question,
          options: item.options,
          answer: item.answer,
          explanation: item.explanation || ""
        };
        
        await createFlashcard(newCard);
      }
      
      alert(`Berhasil import flashcards!`);
      showJsonModal = false;
      loadData();
    } catch (e) {
      console.error(e);
      alert("Format JSON tidak valid atau terjadi kesalahan saat import!");
    } finally {
      isImporting = false;
    }
  }
</script>

<div class="p-6 max-w-6xl mx-auto">
  <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-6 gap-4">
    <h1 class="text-3xl font-black text-slate-800 tracking-tight">Manajemen Flashcard B.Inggris</h1>
  </div>

  <!-- Tabs Navigation -->
  <div class="flex border-b-2 border-slate-200 mb-6 w-full">
    <button 
      class="px-8 py-4 font-bold text-lg transition-colors border-b-4 -mb-[2px] {activeTab === 'categories' ? 'text-indigo-600 border-indigo-600' : 'text-slate-500 border-transparent hover:text-slate-700'}" 
      onclick={() => switchTab('categories')}
    >
      Daftar Kategori
    </button>
    <button 
      class="px-8 py-4 font-bold text-lg transition-colors border-b-4 -mb-[2px] {activeTab === 'flashcards' ? 'text-teal-600 border-teal-600' : 'text-slate-500 border-transparent hover:text-slate-700'}" 
      onclick={() => switchTab('flashcards')}
    >
      Daftar Kartu (Flashcard)
    </button>
  </div>

  {#if activeTab === 'categories'}
    <!-- TAB: KATEGORI -->
    <div class="animate-fade-in">
      <div class="flex justify-end mb-4">
        <button onclick={openCreateCategoryModal} class="bg-indigo-600 hover:bg-indigo-700 text-white px-5 py-2.5 rounded-xl font-bold shadow-sm transition-all active:scale-95 flex items-center gap-2">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
          Buat Kategori Baru
        </button>
      </div>

      <div class="bg-white rounded-2xl shadow-sm overflow-hidden border border-slate-200">
        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse min-w-[600px]">
            <thead>
              <tr class="bg-slate-50 border-b border-slate-200 text-slate-600 text-sm whitespace-nowrap">
                <th class="py-4 px-6 w-16">No</th>
                <th class="py-4 px-6">Nama Kategori</th>
                <th class="py-4 px-6">Deskripsi</th>
                <th class="py-4 px-6 w-32 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {#if isLoading}
                <tr><td colspan="4" class="py-12 text-center text-slate-500 font-medium">Memuat data kategori...</td></tr>
              {:else if categories.length === 0}
                <tr><td colspan="4" class="py-12 text-center text-slate-500 font-medium">Belum ada kategori. Silakan buat kategori pertama Anda.</td></tr>
              {:else}
                {#each categories as cat, i}
                  <tr class="border-b border-slate-100 hover:bg-slate-50 transition-colors">
                    <td class="py-4 px-6 text-slate-500 text-sm font-mono">{i + 1}</td>
                    <td class="py-4 px-6 font-bold text-slate-800">{cat.name}</td>
                    <td class="py-4 px-6 text-slate-600 text-sm">{cat.description || '-'}</td>
                    <td class="py-4 px-6">
                      <div class="flex justify-center gap-2">
                        <button onclick={() => openEditCategoryModal(cat)} class="text-blue-600 hover:bg-blue-100 bg-blue-50 p-2 rounded-lg transition-colors" title="Edit">
                          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path></svg>
                        </button>
                        <button onclick={() => handleDeleteCategory(cat.id)} class="text-red-600 hover:bg-red-100 bg-red-50 p-2 rounded-lg transition-colors" title="Hapus">
                          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
                        </button>
                      </div>
                    </td>
                  </tr>
                {/each}
              {/if}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  {:else}
    <!-- TAB: FLASHCARDS -->
    <div class="animate-fade-in">
      <div class="bg-white p-5 rounded-2xl shadow-sm border border-slate-200 mb-6 flex flex-col sm:flex-row gap-4 items-start sm:items-center justify-between">
        <div class="flex items-center gap-3 w-full sm:w-auto">
          <label for="category-filter" class="font-bold text-slate-700 whitespace-nowrap hidden sm:block">Kategori:</label>
          <select 
            id="category-filter" 
            bind:value={selectedCategoryId} 
            onchange={loadData}
            class="border-2 border-slate-200 rounded-xl px-3 py-2 bg-slate-50 w-full sm:w-48 font-semibold text-slate-700 focus:border-teal-500 outline-none transition-colors"
          >
            <option value="">Semua Kategori</option>
            {#each categories as cat}
              <option value={cat.id}>{cat.name}</option>
            {/each}
          </select>
          <div class="relative w-full sm:w-64">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg>
            </div>
            <input 
              type="text" 
              bind:value={searchQuery}
              oninput={handleSearchInput}
              placeholder="Cari kartu..." 
              class="w-full border-2 border-slate-200 rounded-xl pl-10 pr-10 py-2 bg-slate-50 font-semibold text-slate-700 focus:border-teal-500 outline-none transition-colors"
            />
            {#if searchQuery}
              <button 
                onclick={clearSearch}
                class="absolute inset-y-0 right-0 pr-3 flex items-center text-slate-400 hover:text-slate-600 transition-colors"
                title="Hapus pencarian"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
              </button>
            {/if}
          </div>
        </div>
        
        <div class="flex gap-2 w-full sm:w-auto flex-col sm:flex-row mt-4 sm:mt-0">
          <button onclick={openJsonModal} class="bg-indigo-600 hover:bg-indigo-700 text-white px-5 py-2.5 rounded-xl font-bold shadow-sm transition-all active:scale-95 flex items-center gap-2 whitespace-nowrap justify-center">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"></path></svg>
            Import JSON
          </button>
          <button onclick={openCreateCardModal} class="bg-teal-600 hover:bg-teal-700 text-white px-5 py-2.5 rounded-xl font-bold shadow-sm transition-all active:scale-95 flex items-center gap-2 whitespace-nowrap justify-center">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
            Tambah Kartu
          </button>
        </div>
      </div>

      <!-- Table Flashcards -->
      <div class="bg-white rounded-2xl shadow-sm overflow-hidden border border-slate-200">
        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse min-w-[900px]">
            <thead>
              <tr class="bg-slate-50 border-b border-slate-200 text-slate-600 text-sm whitespace-nowrap">
                <th class="py-4 px-6 w-16">No</th>
                <th class="py-4 px-6 w-48 min-w-[150px]">Kategori</th>
                <th class="py-4 px-6">Pertanyaan</th>
                <th class="py-4 px-6 w-40 min-w-[150px]">Jawaban Benar</th>
                <th class="py-4 px-6 w-32 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {#if isLoading}
                <tr><td colspan="5" class="py-12 text-center text-slate-500 font-medium">Memuat flashcard...</td></tr>
              {:else if flashcards.length === 0 && searchQuery}
                <tr><td colspan="5" class="py-12 text-center text-slate-500 font-medium">Tidak ada kartu yang cocok dengan pencarian Anda.</td></tr>
              {:else if flashcards.length === 0}
                <tr><td colspan="5" class="py-12 text-center text-slate-500 font-medium">Belum ada kartu di kategori ini.</td></tr>
              {:else}
                {#each paginatedFlashcards as card, i}
                  <tr class="border-b border-slate-100 hover:bg-slate-50 transition-colors">
                    <td class="py-4 px-6 text-slate-500 text-sm font-mono">{(currentPage - 1) * itemsPerPage + i + 1}</td>
                    <td class="py-4 px-6">
                      <span class="bg-teal-50 text-teal-800 px-3 py-1.5 rounded-lg border border-teal-100 text-xs font-bold inline-block max-w-full truncate align-middle">
                        {card.category?.name || 'Unknown'}
                      </span>
                    </td>
                    <td class="py-4 px-6 text-slate-800 font-bold">{card.question}</td>
                    <td class="py-4 px-6 text-green-600 font-bold">{card.answer}</td>
                    <td class="py-4 px-6">
                      <div class="flex justify-center gap-2">
                        <button onclick={() => openEditCardModal(card)} class="text-blue-600 hover:bg-blue-100 bg-blue-50 p-2 rounded-lg transition-colors" title="Edit">
                          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path></svg>
                        </button>
                        <button onclick={() => handleDeleteCard(card.id)} class="text-red-600 hover:bg-red-100 bg-red-50 p-2 rounded-lg transition-colors" title="Hapus">
                          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
                        </button>
                      </div>
                    </td>
                  </tr>
                {/each}
              {/if}
            </tbody>
          </table>
        </div>
        
        {#if flashcards.length > itemsPerPage}
          <div class="p-4 border-t border-slate-200 bg-slate-50 flex flex-col sm:flex-row items-center justify-between gap-4">
            <span class="text-sm text-slate-500 font-medium text-center sm:text-left">
              Menampilkan {(currentPage - 1) * itemsPerPage + 1} - {Math.min(currentPage * itemsPerPage, flashcards.length)} dari {flashcards.length} kartu
            </span>
            <div class="flex gap-2 w-full sm:w-auto">
              <button 
                onclick={() => currentPage > 1 && currentPage--} 
                disabled={currentPage === 1}
                class="flex-1 sm:flex-none px-2 sm:px-4 py-2 rounded-lg border border-slate-300 bg-white text-slate-600 font-semibold text-sm sm:text-base disabled:opacity-50 disabled:cursor-not-allowed hover:bg-slate-100 active:scale-95 transition-all shadow-sm text-center"
              >
                Prev
              </button>
              <div class="flex items-center justify-center px-3 sm:px-4 font-bold text-slate-700 bg-white border border-slate-200 rounded-lg shadow-sm text-sm sm:text-base">
                {currentPage} / {totalPages}
              </div>
              <button 
                onclick={() => currentPage < totalPages && currentPage++} 
                disabled={currentPage === totalPages}
                class="flex-1 sm:flex-none px-2 sm:px-4 py-2 rounded-lg border border-slate-300 bg-white text-slate-600 font-semibold text-sm sm:text-base disabled:opacity-50 disabled:cursor-not-allowed hover:bg-slate-100 active:scale-95 transition-all shadow-sm text-center"
              >
                Next
              </button>
            </div>
          </div>
        {/if}
      </div>
    </div>
  {/if}
</div>

<!-- Modal Kategori -->
{#if showCategoryModal}
  <div class="fixed inset-0 bg-slate-900/60 backdrop-blur-sm z-50 flex items-center justify-center p-4 animate-fade-in">
    <div class="bg-white rounded-3xl shadow-2xl w-full max-w-lg flex flex-col overflow-hidden animate-pop">
      <div class="p-6 border-b border-slate-200 flex justify-between items-center bg-slate-50">
        <h2 class="text-xl font-bold text-slate-800">{isEditingCategory ? 'Edit Kategori' : 'Buat Kategori Baru'}</h2>
        <button onclick={() => showCategoryModal = false} class="text-slate-400 hover:text-slate-600 bg-slate-200 hover:bg-slate-300 p-1.5 rounded-full transition-colors">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
        </button>
      </div>
      
      <div class="p-6 space-y-5">
        <div>
          <label class="block text-sm font-bold text-slate-700 mb-2" for="catTitle">Nama Kategori</label>
          <input type="text" id="catTitle" bind:value={currentCategory.name} class="w-full border-2 border-slate-200 rounded-xl px-4 py-3 font-semibold focus:border-indigo-500 outline-none transition-colors" placeholder="Cth: SD Kelas 1, Grammar..." />
        </div>
        <div>
          <label class="block text-sm font-bold text-slate-700 mb-2" for="catDesc">Deskripsi (Opsional)</label>
          <textarea id="catDesc" bind:value={currentCategory.description} rows="2" class="w-full border-2 border-slate-200 rounded-xl px-4 py-3 focus:border-indigo-500 outline-none transition-colors" placeholder="Deskripsi materi kategori ini..."></textarea>
        </div>
      </div>
      
      <div class="p-5 border-t border-slate-200 bg-slate-50 flex justify-end gap-3">
        <button onclick={() => showCategoryModal = false} class="px-6 py-2.5 rounded-xl font-bold text-slate-600 hover:bg-slate-200 transition-colors">Batal</button>
        <button onclick={saveCategory} class="px-6 py-2.5 rounded-xl font-bold text-white bg-indigo-600 hover:bg-indigo-700 transition-colors shadow-md active:scale-95">Simpan Kategori</button>
      </div>
    </div>
  </div>
{/if}

<!-- Modal Card Form -->
{#if showCardModal}
  <div class="fixed inset-0 bg-slate-900/60 backdrop-blur-sm z-50 flex items-center justify-center p-4 animate-fade-in">
    <div class="bg-white rounded-3xl shadow-2xl w-full max-w-2xl max-h-[90vh] flex flex-col overflow-hidden animate-pop">
      <div class="p-6 border-b border-slate-200 flex justify-between items-center bg-slate-50">
        <h2 class="text-xl font-bold text-slate-800">{isEditingCard ? 'Edit Kartu' : 'Tambah Kartu Baru'}</h2>
        <button onclick={() => showCardModal = false} class="text-slate-400 hover:text-slate-600 bg-slate-200 hover:bg-slate-300 p-1.5 rounded-full transition-colors">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
        </button>
      </div>
      
      <div class="p-6 overflow-y-auto flex-1 space-y-5">
        <div>
          <label class="block text-sm font-bold text-slate-700 mb-2" for="cardCategory">Kategori</label>
          <select id="cardCategory" bind:value={currentCard.category_id} class="w-full border-2 border-slate-200 rounded-xl px-4 py-3 bg-white font-semibold focus:border-teal-500 outline-none transition-colors">
            {#each categories as cat}
              <option value={cat.id}>{cat.name}</option>
            {/each}
          </select>
        </div>
        
        <div>
          <label class="block text-sm font-bold text-slate-700 mb-2" for="question">Pertanyaan</label>
          <textarea id="question" bind:value={currentCard.question} rows="2" class="w-full border-2 border-slate-200 rounded-xl px-4 py-3 font-semibold focus:border-teal-500 outline-none transition-colors" placeholder="Cth: Apa bahasa inggrisnya Kucing?"></textarea>
        </div>
        
        <div class="bg-teal-50/50 p-5 rounded-2xl border border-teal-100">
          <label class="block text-sm font-bold text-teal-800 mb-3">Pilihan Jawaban (4 Pilihan)</label>
          <div class="grid grid-cols-2 gap-4">
            {#each currentCard.options as _, i}
              <input type="text" bind:value={currentCard.options[i]} class="border-2 border-teal-200 rounded-xl px-4 py-2 font-medium focus:border-teal-500 outline-none" placeholder={`Pilihan ${i+1}`} />
            {/each}
          </div>
        </div>
        
        <div>
          <label class="block text-sm font-bold text-slate-700 mb-2" for="answer">Jawaban Benar</label>
          <select id="answer" bind:value={currentCard.answer} class="w-full border-2 border-slate-200 rounded-xl px-4 py-3 bg-white font-bold text-teal-700 focus:border-teal-500 outline-none transition-colors">
            <option value="" disabled>-- Pilih jawaban benar --</option>
            {#each currentCard.options as opt}
              {#if opt}
                <option value={opt}>{opt}</option>
              {/if}
            {/each}
          </select>
          <p class="text-xs text-slate-500 mt-2">*Pastikan opsi di atas terisi dahulu lalu pilih jawaban benar.</p>
        </div>
        
        <div>
          <label class="block text-sm font-bold text-slate-700 mb-2" for="explanation">Penjelasan (Tampil saat kartu dibalik)</label>
          <textarea id="explanation" bind:value={currentCard.explanation} rows="2" class="w-full border-2 border-slate-200 rounded-xl px-4 py-3 focus:border-teal-500 outline-none transition-colors" placeholder="Cth: Cat artinya Kucing. Dog artinya Anjing..."></textarea>
        </div>
      </div>
      
      <div class="p-5 border-t border-slate-200 bg-slate-50 flex justify-end gap-3">
        <button onclick={() => showCardModal = false} class="px-6 py-2.5 rounded-xl font-bold text-slate-600 hover:bg-slate-200 transition-colors">Batal</button>
        <button onclick={saveCard} class="px-6 py-2.5 rounded-xl font-bold text-white bg-teal-600 hover:bg-teal-700 transition-colors shadow-md active:scale-95">Simpan Kartu</button>
      </div>
    </div>
  </div>
{/if}

<!-- Modal Import JSON -->
{#if showJsonModal}
  <div class="fixed inset-0 bg-slate-900/60 backdrop-blur-sm z-50 flex items-center justify-center p-4 animate-fade-in">
    <div class="bg-white rounded-3xl shadow-2xl w-full max-w-3xl max-h-[90vh] flex flex-col overflow-hidden animate-pop">
      <div class="p-6 border-b border-slate-200 flex justify-between items-center bg-slate-50">
        <h2 class="text-xl font-bold text-slate-800">Import Flashcard via JSON</h2>
        <button onclick={() => showJsonModal = false} class="text-slate-400 hover:text-slate-600 bg-slate-200 hover:bg-slate-300 p-1.5 rounded-full transition-colors">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
        </button>
      </div>
      
      <div class="p-6 overflow-y-auto flex-1 space-y-5">
        <div class="bg-blue-50 text-blue-800 p-4 rounded-xl text-sm font-medium border border-blue-200">
          <p class="mb-2">Format JSON harus berupa <strong>array of objects</strong>. Contoh format:</p>
          <pre class="bg-white/60 p-2 rounded-lg text-xs overflow-x-auto text-blue-900 border border-blue-100">
[
  &#123;
    "question": "Pertanyaan 1",
    "options": ["Opsi A", "Opsi B", "Opsi C", "Opsi D"],
    "answer": "Opsi A",
    "explanation": "Penjelasan jawaban A"
  &#125;
]
          </pre>
        </div>
        
        <div>
          <label class="block text-sm font-bold text-slate-700 mb-2" for="jsonCategory">Pilih Kategori Tujuan</label>
          <select id="jsonCategory" bind:value={jsonCategoryId} class="w-full border-2 border-slate-200 rounded-xl px-4 py-3 bg-white font-semibold focus:border-indigo-500 outline-none transition-colors">
            {#each categories as cat}
              <option value={cat.id}>{cat.name}</option>
            {/each}
          </select>
        </div>
        
        <div>
          <label class="block text-sm font-bold text-slate-700 mb-2" for="jsonInput">Data JSON (seperti raw body Postman)</label>
          <textarea id="jsonInput" bind:value={jsonInput} rows="10" class="w-full border-2 border-slate-200 rounded-xl px-4 py-3 font-mono text-sm focus:border-indigo-500 outline-none transition-colors" placeholder="Paste JSON array di sini..."></textarea>
        </div>
      </div>
      
      <div class="p-5 border-t border-slate-200 bg-slate-50 flex justify-end gap-3">
        <button onclick={() => showJsonModal = false} class="px-6 py-2.5 rounded-xl font-bold text-slate-600 hover:bg-slate-200 transition-colors">Batal</button>
        <button onclick={handleJsonImport} disabled={isImporting || !jsonCategoryId} class="px-6 py-2.5 rounded-xl font-bold text-white bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors shadow-md active:scale-95 flex items-center gap-2">
          {#if isImporting}
            <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Mengimpor...
          {:else}
            Mulai Import
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }
  @keyframes pop {
    0% { transform: scale(0.95); opacity: 0; }
    100% { transform: scale(1); opacity: 1; }
  }
  .animate-fade-in {
    animation: fadeIn 0.2s ease-out forwards;
  }
  .animate-pop {
    animation: pop 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275) forwards;
  }
</style>
