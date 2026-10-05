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
  
  let isLoading = $state(true);
  
  // Modals state
  let showCardModal = $state(false);
  let isEditingCard = $state(false);
  
  let showCategoryModal = $state(false);
  let isEditingCategory = $state(false);

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
      categories = await fetchFlashcardCategories();
      if (selectedCategoryId) {
        flashcards = await fetchFlashcards(selectedCategoryId as number);
      } else {
        flashcards = await fetchFlashcards();
      }
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
    currentCategory = { ...cat };
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
    currentCard = { ...card, options: [...card.options] };
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
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-slate-50 border-b border-slate-200 text-slate-600 text-sm">
              <th class="py-4 px-6 w-16">ID</th>
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
              {#each categories as cat}
                <tr class="border-b border-slate-100 hover:bg-slate-50 transition-colors">
                  <td class="py-4 px-6 text-slate-500 text-sm font-mono">{cat.id}</td>
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
  {:else}
    <!-- TAB: FLASHCARDS -->
    <div class="animate-fade-in">
      <div class="bg-white p-5 rounded-2xl shadow-sm border border-slate-200 mb-6 flex flex-col sm:flex-row gap-4 items-start sm:items-center justify-between">
        <div class="flex items-center gap-3 w-full sm:w-auto">
          <label for="category-filter" class="font-bold text-slate-700 whitespace-nowrap">Filter Kategori:</label>
          <select 
            id="category-filter" 
            bind:value={selectedCategoryId} 
            onchange={loadData}
            class="border-2 border-slate-200 rounded-xl px-4 py-2 bg-slate-50 min-w-[220px] font-semibold text-slate-700 focus:border-teal-500 outline-none transition-colors"
          >
            <option value="">Semua Kategori</option>
            {#each categories as cat}
              <option value={cat.id}>{cat.name}</option>
            {/each}
          </select>
        </div>
        
        <button onclick={openCreateCardModal} class="bg-teal-600 hover:bg-teal-700 text-white px-5 py-2.5 rounded-xl font-bold shadow-sm transition-all active:scale-95 flex items-center gap-2 whitespace-nowrap">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
          Tambah Kartu
        </button>
      </div>

      <!-- Table Flashcards -->
      <div class="bg-white rounded-2xl shadow-sm overflow-hidden border border-slate-200">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-slate-50 border-b border-slate-200 text-slate-600 text-sm">
              <th class="py-4 px-6 w-16">ID</th>
              <th class="py-4 px-6">Kategori</th>
              <th class="py-4 px-6">Pertanyaan</th>
              <th class="py-4 px-6">Jawaban Benar</th>
              <th class="py-4 px-6 w-32 text-center">Aksi</th>
            </tr>
          </thead>
          <tbody>
            {#if isLoading}
              <tr><td colspan="5" class="py-12 text-center text-slate-500 font-medium">Memuat flashcard...</td></tr>
            {:else if flashcards.length === 0}
              <tr><td colspan="5" class="py-12 text-center text-slate-500 font-medium">Belum ada kartu di kategori ini.</td></tr>
            {:else}
              {#each flashcards as card}
                <tr class="border-b border-slate-100 hover:bg-slate-50 transition-colors">
                  <td class="py-4 px-6 text-slate-500 text-sm font-mono">{card.id}</td>
                  <td class="py-4 px-6">
                    <span class="bg-teal-50 text-teal-800 px-3 py-1.5 rounded-lg border border-teal-100 text-xs font-bold">
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
