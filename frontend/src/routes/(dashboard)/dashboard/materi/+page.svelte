<script lang="ts">
  import "katex/dist/katex.min.css";
  import Modal from "$lib/components/Modal.svelte";
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import { renderMathContent } from "$lib/extensions/math";
  import { toast } from "$lib/stores/toast.svelte";

  type Subject = {
    id: number;
    name: string;
  };

  type User = {
    id: number;
    username: string;
    role: string;
  };

  type Materi = {
    id: number;
    title: string;
    subject_id: number;
    subject: Subject;
    users: User[];
    content: string;
    created_at: string;
  };

  let materis = $state<Materi[]>([]);
  let subjects = $state<Subject[]>([]);
  let users = $state<User[]>([]);
  let searchQuery = $state("");
  let viewingMateri = $state<Materi | null>(null);
  let currentUser = $state<User | null>(null);

  let currentPage = $state(1);
  let limit = $state(12);
  let totalPages = $state(1);
  let totalItems = $state(0);
  let selectedSubjectId = $state<number | "all">("all");
  let subjectsList = $state<Subject[]>([]);
  let searchTimeout: any;

  async function fetchMateris() {
    try {
      const url = new URL("/api/materis", window.location.origin);
      url.searchParams.set("page", currentPage.toString());
      url.searchParams.set("limit", limit.toString());
      if (searchQuery.trim()) url.searchParams.set("search", searchQuery.trim());
      if (selectedSubjectId !== "all") url.searchParams.set("subject_id", selectedSubjectId.toString());

      const resMateris = await fetch(url.toString(), { credentials: "include" });
      if (resMateris.ok) {
        const json = await resMateris.json();
        if (json.data) {
          materis = json.data;
          currentPage = json.page;
          totalPages = json.total_pages;
          totalItems = json.total;
        } else {
          materis = json;
          totalPages = 1;
        }
      }
    } catch (e) {
      console.error(e);
    }
  }

  function handleSearchInput() {
    clearTimeout(searchTimeout);
    searchTimeout = setTimeout(() => {
      currentPage = 1;
      fetchMateris();
    }, 400);
  }

  function handleSubjectChange() {
    currentPage = 1;
    fetchMateris();
  }

  function changePage(pageStr: number) {
    if (pageStr >= 1 && pageStr <= totalPages) {
      currentPage = pageStr;
      fetchMateris();
    }
  }

  onMount(async () => {
    try {
      const resMe = await fetch("/me", { credentials: "include" });
      if (resMe.ok) {
        const data = await resMe.json();
        currentUser = data.user;
      }

      const resSubjects = await fetch("/api/subjects", { credentials: "include" });
      if (resSubjects.ok) {
        subjectsList = await resSubjects.json();
      }

      await fetchMateris();
    } catch (e) {
      console.error(e);
    }
  });

  $effect(() => {
    const materiId = page.url.searchParams.get("id");
    if (materiId && materis.length > 0) {
      const found = materis.find((m) => m.id === Number(materiId));
      if (found && viewingMateri?.id !== found.id) viewingMateri = found;
    } else if (materis.length > 0 && !materiId) {
      if (viewingMateri) viewingMateri = null;
    }
  });

  function openMateri(materi: Materi) {
    viewingMateri = materi;
    const url = new URL(page.url);
    url.searchParams.set("id", String(materi.id));
    goto(url, { keepFocus: true, noScroll: true });
  }

  function closeMateri() {
    viewingMateri = null;
    const url = new URL(page.url);
    url.searchParams.delete("id");
    goto(url, { keepFocus: true, noScroll: true });
  }

  // Modals state
  let showDeleteModal = $state(false);
  let materiToDelete = $state<Materi | null>(null);
  let openMenuId = $state<number | null>(null);

  function openEdit(materi?: Materi) {
    if (materi) {
      goto(`/dashboard/materi/form?id=${materi.id}`);
    } else {
      goto("/dashboard/materi/form");
    }
  }

  function deleteMateri(materi: Materi) {
    materiToDelete = materi;
    showDeleteModal = true;
  }

  async function confirmDelete() {
    if (!materiToDelete) return;
    const id = materiToDelete.id;
    showDeleteModal = false;
    try {
      const res = await fetch(`/api/materis/${id}`, {
        method: "DELETE",
        credentials: "include",
      });
      if (res.ok) {
        materis = materis.filter((m) => m.id !== id);
        if (viewingMateri?.id === id) closeMateri();
        toast.success("Materi berhasil dihapus");
      } else {
        toast.error("Gagal menghapus materi");
      }
    } catch (e) {
      console.error(e);
      toast.error("Terjadi kesalahan");
    } finally {
      materiToDelete = null;
    }
  }

  function printMateri() {
    window.print();
  }

  let printFontSize = $state(16);
</script>

<svelte:head>
  <title>Materi Pelajaran - Les Balongarut</title>
  <style>
    @media print {
      @page {
        size: A4 portrait;
        margin: 20mm;
      }
      body {
        background: white !important;
      }
    }
  </style>
</svelte:head>

<svelte:window onclick={() => (openMenuId = null)} />

<!-- Delete Modal -->
<Modal show={showDeleteModal} onclose={() => (showDeleteModal = false)} maxWidth="max-w-sm">
  <div class="space-y-4 text-center">
    <div class="w-12 h-12 mx-auto rounded-full bg-red-100 flex items-center justify-center">
      <svg class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
      </svg>
    </div>
    <div>
      <h3 class="text-lg font-semibold text-slate-900">Hapus Materi</h3>
      <p class="text-sm text-slate-600 mt-1">Materi "{materiToDelete?.title}" akan dihapus permanen.</p>
    </div>
    <div class="flex gap-2 justify-center pt-2">
      <button onclick={() => (showDeleteModal = false)} class="px-4 py-2 text-sm rounded-lg border border-slate-300 hover:bg-slate-50 text-slate-900 cursor-pointer">Batal</button>
      <button onclick={confirmDelete} class="px-4 py-2 text-sm rounded-lg bg-red-500 text-white hover:bg-red-700 cursor-pointer">Hapus</button>
    </div>
  </div>
</Modal>


{#if viewingMateri}
  <!-- VIEW: DETAIL MATERI -->
  <div class="w-full mx-auto bg-transparent md:bg-white p-4 sm:p-6 md:p-[20mm] md:max-w-[210mm] md:rounded-sm md:shadow-xl md:border border-slate-200 print:max-w-none print:border-none print:shadow-none print:p-0 print:m-0 print:min-h-0 animate-in fade-in zoom-in-95 duration-200 relative min-h-[50vh] md:min-h-[297mm]">
    <div class="flex flex-col sm:flex-row justify-between items-center gap-5 sm:gap-4 print:hidden mb-8 border-b border-slate-100 pb-5">
      <button onclick={() => closeMateri()} class="flex items-center gap-2 text-slate-500 hover:text-slate-900 font-medium transition-colors cursor-pointer w-full justify-center sm:justify-start sm:w-auto">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path></svg>
        Kembali
      </button>
      <div class="flex flex-wrap items-center justify-center sm:justify-end gap-4 sm:gap-6 w-full sm:w-auto mt-2 sm:mt-0">
        {#if currentUser?.role === "teacher"}
          <button onclick={() => deleteMateri(viewingMateri!)} class="p-2.5 text-red-600 bg-red-50 hover:bg-red-100 border border-red-200 rounded-xl transition-colors cursor-pointer flex items-center justify-center shrink-0" title="Hapus">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
          </button>
          <button onclick={() => openEdit(viewingMateri!)} class="p-2.5 text-blue-600 bg-blue-50 hover:bg-blue-100 border border-blue-200 rounded-xl transition-colors cursor-pointer flex items-center justify-center shrink-0" title="Edit">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path></svg>
          </button>
          <button onclick={printMateri} class="p-2.5 text-slate-700 bg-slate-100 hover:bg-slate-200 border border-slate-200 rounded-xl transition-colors cursor-pointer flex items-center justify-center shrink-0" title="Cetak">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z"></path></svg>
          </button>
          <div class="flex items-center bg-slate-100 rounded-xl overflow-hidden print:hidden border border-slate-200 shrink-0">
            <button onclick={() => (printFontSize = Math.max(10, printFontSize - 1))} class="px-3 py-2 text-sm font-bold text-slate-600 hover:bg-slate-200 hover:text-slate-900 transition-colors cursor-pointer" title="Perkecil Font">A-</button>
            <div class="px-2 py-2 text-xs font-medium text-slate-600 bg-white border-x border-slate-200 text-center min-w-[3rem]">{printFontSize}px</div>
            <button onclick={() => (printFontSize = Math.min(36, printFontSize + 1))} class="px-3 py-2 text-base font-bold text-slate-600 hover:bg-slate-200 hover:text-slate-900 transition-colors cursor-pointer" title="Perbesar Font">A+</button>
          </div>
        {/if}
      </div>
    </div>

    <!-- Print Area -->
    <div class="print:block">
      <h1 class="text-xl sm:text-2xl font-semibold text-slate-900 mb-4 tracking-tight leading-tight text-center sm:text-left">
        {viewingMateri.title}
      </h1>
      <div class="flex flex-col sm:flex-row items-center sm:justify-start gap-3 text-sm text-slate-500 mb-8 pb-8 border-b border-slate-100">
        <span class="text-slate-700 font-medium flex items-center gap-1.5">
          <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"></path></svg>
          {viewingMateri.subject?.name || "Mata Pelajaran"}
        </span>
        {#if currentUser?.role === "teacher"}
          <span class="flex items-center gap-1.5 text-slate-700 font-medium max-w-[200px] sm:max-w-[300px] truncate" title={viewingMateri.users?.map(u => u.username).join(', ') || "Unknown"}>
            <svg class="w-4 h-4 text-slate-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path></svg>
            <span class="truncate">Siswa: {viewingMateri.users?.map(u => u.username).join(', ') || "Unknown"}</span>
          </span>
        {/if}
        <span class="flex items-center gap-1.5">
          <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg>
          {new Date(viewingMateri.created_at).toLocaleDateString("id-ID", { year: "numeric", month: "short", day: "numeric" })}
        </span>
      </div>
      <div class="prose prose-slate max-w-none text-slate-800 leading-loose whitespace-pre-wrap text-justify wrap-break-word text-lg sm:text-(length:--base-size) print:text-(length:--base-size)" style="--base-size: {printFontSize}px; tab-size: 4;">
        {@html renderMathContent(viewingMateri.content)}
      </div>
    </div>
  </div>
{:else}
  <!-- VIEW: LIST -->
  <div class="max-w-4xl mx-auto space-y-6 animate-in fade-in duration-300 print:hidden">
    <!-- Header & Actions -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold text-slate-900 tracking-tight">Materi Pelajaran</h1>
        <p class="text-sm text-slate-500 mt-1">Daftar materi yang tersedia.</p>
      </div>

      <div class="flex items-center gap-3">
        {#if currentUser?.role === "teacher"}
          <button onclick={() => openEdit()} class="inline-flex items-center gap-2 px-5 py-2.5 text-sm font-bold rounded-full bg-blue-50 text-blue-700 border border-blue-200 hover:bg-blue-100 transition-colors shadow-sm cursor-pointer">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
            Materi Baru
          </button>
        {/if}
      </div>
    </div>

    <!-- Search Bar & Filter -->
    <div class="flex flex-col sm:flex-row gap-3 max-w-2xl">
      <div class="relative flex-1">
        <div class="absolute inset-y-0 left-0 pl-5 flex items-center pointer-events-none">
          <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg>
        </div>
        <input type="text" bind:value={searchQuery} oninput={handleSearchInput} placeholder="Cari materi..." class="block w-full pl-14 pr-12 py-2.5 border border-slate-200 rounded-full bg-white/60 focus:bg-white text-sm text-slate-800 placeholder-slate-400 focus:outline-none focus:border-slate-300 focus:ring-4 focus:ring-slate-100 transition-all" />
        {#if searchQuery}
          <button onclick={() => { searchQuery = ""; handleSearchInput(); }} class="absolute right-4 top-1/2 -translate-y-1/2 p-1 text-slate-400 hover:text-slate-600 rounded-full hover:bg-slate-100 transition-colors cursor-pointer" title="Clear">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
          </button>
        {/if}
      </div>
      
      <div class="relative shrink-0 sm:w-56">
        <select 
          bind:value={selectedSubjectId}
          onchange={handleSubjectChange}
          class="block w-full pl-4 pr-10 py-2.5 border border-slate-200 rounded-full bg-white/60 focus:bg-white text-sm text-slate-800 focus:outline-none focus:border-slate-300 focus:ring-4 focus:ring-slate-100 transition-all appearance-none cursor-pointer truncate"
        >
          <option value="all">Semua Mata Pelajaran</option>
          {#each subjectsList as subject}
            <option value={subject.id}>{subject.name}</option>
          {/each}
        </select>
        <div class="absolute inset-y-0 right-0 pr-4 flex items-center pointer-events-none">
          <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
        </div>
      </div>
    </div>

    <!-- Content Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      {#each materis as materi}
        <div
          class="bg-white border border-slate-200 rounded-2xl p-5 shadow-sm hover:shadow-md transition-shadow relative flex flex-col h-full {openMenuId === materi.id ? 'z-50' : 'z-0'}"
        >
          <div class="flex justify-between items-center mb-2.5 gap-2">
            <div class="flex items-center gap-2 flex-wrap">
              <span class="text-xs font-bold text-slate-400 bg-slate-50 px-2 py-1 rounded-md border border-slate-100">#{materi.id}</span>
              {#if currentUser?.role === "teacher"}
                <span class="px-2 py-1 bg-indigo-50 text-indigo-700 rounded-md text-[10px] font-bold border border-indigo-100 whitespace-nowrap uppercase tracking-wide flex items-center gap-1 max-w-[120px] sm:max-w-[150px] truncate" title={materi.users?.map(u => u.username).join(', ') || "Siswa"}>
                  <svg class="w-3 h-3 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path></svg>
                  <span class="truncate">{materi.users?.length ? (materi.users.length > 2 ? `${materi.users[0].username}, ${materi.users[1].username} +${materi.users.length - 2}` : materi.users.map(u => u.username).join(', ')) : "Siswa"}</span>
                </span>
              {/if}
            </div>
            <span class="px-2.5 py-1 bg-blue-50 text-blue-700 rounded-md text-[10px] font-bold border border-blue-100 whitespace-nowrap uppercase tracking-wide">
              {materi.subject?.name || "Mata Pelajaran"}
            </span>
          </div>
          
          <div class="grow">
            <h3 class="text-[17px] font-bold text-slate-900 leading-snug line-clamp-3 mb-3" title={materi.title}>{materi.title}</h3>
          </div>
          
          <div class="flex items-center text-xs font-medium text-slate-500 mb-4 bg-slate-50 w-max px-2.5 py-1.5 rounded-md border border-slate-100">
            <svg class="w-4 h-4 mr-1.5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg>
            {new Date(materi.created_at).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })}
          </div>

          <div class="mt-auto flex gap-2 pt-4 border-t border-slate-100">
            <button
              onclick={() => openMateri(materi)}
              class="flex-1 inline-flex items-center justify-center p-2.5 text-sm font-semibold text-blue-700 bg-blue-50 rounded-xl border border-blue-100 hover:bg-blue-100 transition-colors cursor-pointer"
              title="Buka Materi"
            >
              <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg>
              Buka
            </button>
            
            {#if currentUser?.role === "teacher"}
              <div class="relative shrink-0">
                <button
                  onclick={(e) => {
                    e.stopPropagation();
                    openMenuId = openMenuId === materi.id ? null : materi.id;
                  }}
                  class="p-2.5 text-slate-500 bg-slate-50 border border-slate-200 hover:bg-slate-100 rounded-xl transition-colors cursor-pointer flex items-center justify-center h-full"
                  title="Menu"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"></path></svg>
                </button>

                {#if openMenuId === materi.id}
                  <div class="absolute bottom-full right-0 mb-2 w-48 bg-white rounded-xl shadow-xl border border-slate-100 overflow-hidden animate-in fade-in zoom-in-95 duration-150 z-50 origin-bottom-right">
                    <button
                      onclick={(e) => {
                        e.stopPropagation();
                        openMenuId = null;
                        openEdit(materi);
                      }}
                      class="w-full text-left px-4 py-2.5 text-sm font-medium text-slate-700 hover:bg-slate-50 hover:text-blue-600 transition-colors flex items-center gap-2 cursor-pointer"
                    >
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path></svg>
                      Edit Materi
                    </button>
                    <div class="h-px bg-slate-100"></div>
                    <button
                      onclick={(e) => {
                        e.stopPropagation();
                        openMenuId = null;
                        materiToDelete = materi;
                        showDeleteModal = true;
                      }}
                      class="w-full text-left px-4 py-2.5 text-sm font-medium text-red-600 hover:bg-red-50 transition-colors flex items-center gap-2 cursor-pointer"
                    >
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
                      Hapus
                    </button>
                  </div>
                {/if}
              </div>
            {/if}
          </div>
        </div>
      {:else}
        <div class="col-span-full py-16 text-center text-slate-500 bg-white/50 border border-slate-200 border-dashed rounded-2xl">
          <svg class="w-16 h-16 mx-auto mb-4 text-slate-300" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"></path></svg>
          <p class="text-base font-medium text-slate-700">Belum ada materi</p>
          <p class="text-sm mt-1">Buat materi baru untuk memulai.</p>
        </div>
      {/each}
    </div>

    <!-- Pagination Controls -->
    {#if totalPages > 1}
      <div class="flex items-center justify-center gap-2 mt-8 mb-4">
        <button 
          onclick={() => changePage(currentPage - 1)}
          disabled={currentPage === 1}
          class="flex items-center justify-center w-10 h-10 md:w-auto md:px-4 md:py-2 border border-slate-200 rounded-xl text-sm font-medium text-slate-600 bg-white hover:bg-slate-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors cursor-pointer shadow-sm"
          title="Sebelumnya"
        >
          <svg class="w-5 h-5 md:hidden" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path></svg>
          <span class="hidden md:inline">Sebelumnya</span>
        </button>
        
        <div class="flex items-center text-sm font-medium text-slate-700 px-4 h-10 bg-slate-50 border border-slate-100 rounded-xl">
          <span class="md:hidden">{currentPage} <span class="text-slate-400 mx-1">/</span> {totalPages}</span>
          <span class="hidden md:inline">Halaman {currentPage} dari {totalPages}</span>
        </div>
        
        <button 
          onclick={() => changePage(currentPage + 1)}
          disabled={currentPage === totalPages}
          class="flex items-center justify-center w-10 h-10 md:w-auto md:px-4 md:py-2 border border-slate-200 rounded-xl text-sm font-medium text-slate-600 bg-white hover:bg-slate-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors cursor-pointer shadow-sm"
          title="Selanjutnya"
        >
          <svg class="w-5 h-5 md:hidden" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path></svg>
          <span class="hidden md:inline">Selanjutnya</span>
        </button>
      </div>
    {/if}
  </div>
{/if}
