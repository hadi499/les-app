<script lang="ts">
  import { onMount } from "svelte";

  type UserRef = { username: string };
  type Progress = { user?: UserRef; lessonId: number; bestWPM: number; bestAccuracy: number; stars: number; completed: boolean; attempts: number };
  type GameScore = { user?: UserRef; mode: string; score: number };
  type LessonHistory = { created_at: string; user?: UserRef; lessonId: number; wpm: number; accuracy: number; stars: number };
  type GameHistory = { created_at: string; user?: UserRef; mode: string; score: number };

  let activeTab = $state("progress");

  onMount(() => {
    const savedTab = localStorage.getItem("typingAdminTab");
    if (savedTab) {
      activeTab = savedTab;
    }
  });

  $effect(() => {
    localStorage.setItem("typingAdminTab", activeTab);
  });
  let isLoading = $state(true);
  let showLoadingSpinner = $state(false);
  let errorMsg = $state("");

  let progressData: Progress[] = $state([]);
  let gameScoresData: GameScore[] = $state([]);
  let lessonHistoryData: LessonHistory[] = $state([]);
  let gameHistoryData: GameHistory[] = $state([]);

  let pageState = $state({ progress: 1, gameScores: 1, lessonHistory: 1, gameHistory: 1 });
  let totalPages = $state({ progress: 1, gameScores: 1, lessonHistory: 1, gameHistory: 1 });
  const limit = 24;

  async function fetchTabData(tab: string, page: number) {
    isLoading = true;
    showLoadingSpinner = false;
    setTimeout(() => { showLoadingSpinner = true; }, 150);
    errorMsg = "";
    
    let endpoint = "";
    if (tab === "progress") endpoint = `/api/typing/admin/progress?page=${page}&limit=${limit}`;
    else if (tab === "gameScores") endpoint = `/api/typing/admin/game-scores?page=${page}&limit=${limit}`;
    else if (tab === "lessonHistory") endpoint = `/api/typing/admin/history/lesson?page=${page}&limit=${limit}`;
    else if (tab === "gameHistory") endpoint = `/api/typing/admin/history/game?page=${page}&limit=${limit}`;

    try {
      const res = await fetch(endpoint, { credentials: "include" });
      if (res.status === 403) {
        errorMsg = "Akses ditolak. Anda bukan Guru.";
        return;
      }
      if (!res.ok) throw new Error("Gagal mengambil data dari server");
      
      const resJson = await res.json();
      const data = resJson.data || [];
      const total = resJson.total_pages || 1;
      
      if (tab === "progress") {
        progressData = data;
        totalPages.progress = total;
      } else if (tab === "gameScores") {
        gameScoresData = data;
        totalPages.gameScores = total;
      } else if (tab === "lessonHistory") {
        lessonHistoryData = data;
        totalPages.lessonHistory = total;
      } else if (tab === "gameHistory") {
        gameHistoryData = data;
        totalPages.gameHistory = total;
      }
    } catch (e) {
      errorMsg = e instanceof Error ? e.message : String(e);
    } finally {
      isLoading = false;
    }
  }

  $effect(() => {
    fetchTabData(activeTab, pageState[activeTab as keyof typeof pageState]);
  });

  const tabs = [
    { id: "progress", label: "Lesson Progress" },
    { id: "gameScores", label: "Game High Scores" },
    { id: "lessonHistory", label: "Lesson History" },
    { id: "gameHistory", label: "Game History" },
  ];

  function formatDate(dateStr: string | null | undefined) {
    if (!dateStr) return "-";
    return new Date(dateStr).toLocaleString("id-ID", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  }
</script>

<svelte:head>
  <title>Typing Monitoring - Portal Guru</title>
</svelte:head>

<div class="animate-in fade-in duration-500 max-w-7xl mx-auto space-y-6">
  <div class="mb-8">
    <h1
      class="text-2xl font-bold text-slate-900 sm:text-3xl tracking-tight drop-shadow-sm"
    >
      Pemantauan Mengetik
    </h1>
    <p class="mt-2 text-slate-600 text-md tracking-wide">
      Pantau perkembangan kemampuan mengetik seluruh murid.
    </p>
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
    <!-- Tabs -->
    <div class="mb-6 flex">
      <div class="grid grid-cols-2 w-full md:w-max md:flex p-1.5 gap-1.5 bg-slate-100/80 backdrop-blur-sm rounded-2xl">
        {#each tabs as tab}
          <button
            class="px-4 py-2.5 rounded-xl text-sm font-bold transition-all duration-200 {activeTab === tab.id ? 'bg-white text-indigo-600 shadow-sm ring-1 ring-slate-900/5' : 'text-slate-500 hover:text-slate-700 hover:bg-slate-200/50'}"
            onclick={() => (activeTab = tab.id)}
          >
            {tab.label}
          </button>
        {/each}
      </div>
    </div>

    <!-- Content -->
    <div class="mt-6 pt-2">
      <div>
        {#if activeTab === "progress"}
          {#if progressData.length === 0}
            <div class="py-12 text-center">
              <h3 class="text-lg font-bold text-slate-700 mb-1">Belum ada data progress.</h3>
            </div>
          {:else}
            <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
              {#each progressData as p}
                <div class="bg-white border border-slate-200 rounded-2xl p-5 shadow-sm hover:shadow-md transition-shadow flex flex-col h-full">
                  <div class="flex justify-between items-start mb-4">
                    <div class="flex flex-col">
                      <span class="text-xs font-bold text-slate-400">Lesson {p.lessonId}</span>
                      <h3 class="text-lg font-bold text-slate-900 mt-0.5">@{p.user?.username || "Unknown"}</h3>
                    </div>
                    <span class={`px-2 py-1 rounded-md text-[10px] font-bold border whitespace-nowrap uppercase tracking-wide ${p.completed ? "bg-emerald-50 text-emerald-600 border-emerald-100" : "bg-slate-50 text-slate-500 border-slate-200"}`}>
                      {p.completed ? "Selesai" : "Belum Selesai"}
                    </span>
                  </div>
                  
                  <div class="grid grid-cols-2 gap-3 mb-4">
                    <div class="bg-indigo-50 border border-indigo-100 rounded-xl p-3 flex flex-col items-center justify-center">
                      <span class="text-2xl font-black text-indigo-600">{p.bestWPM}</span>
                      <span class="text-[10px] font-bold text-indigo-400 uppercase tracking-wider">Best WPM</span>
                    </div>
                    <div class="bg-blue-50 border border-blue-100 rounded-xl p-3 flex flex-col items-center justify-center">
                      <span class="text-2xl font-black text-blue-600">{p.bestAccuracy}%</span>
                      <span class="text-[10px] font-bold text-blue-400 uppercase tracking-wider">Akurasi</span>
                    </div>
                  </div>

                  <div class="mt-auto flex items-center justify-between pt-4 border-t border-slate-100">
                    <div class="flex gap-1 text-lg">
                      <span class="text-amber-400">{"★".repeat(p.stars)}</span><span class="text-slate-200">{"★".repeat(3 - p.stars)}</span>
                    </div>
                    <span class="text-xs font-medium text-slate-500 bg-slate-50 px-2.5 py-1 rounded-md border border-slate-100">
                      {p.attempts} Percobaan
                    </span>
                  </div>
                </div>
              {/each}
            </div>

            <!-- Pagination Controls -->
            {#if totalPages.progress > 1}
            <div class="mt-8 px-4 py-3 sm:px-6 sm:py-4 bg-white/60 backdrop-blur-md rounded-2xl border border-slate-200 shadow-sm flex flex-row items-center justify-between gap-4 mb-4">
              <div class="text-sm text-slate-600 text-left">
                Halaman <span class="font-medium text-slate-900">{pageState.progress}</span> dari <span class="font-medium text-slate-900">{totalPages.progress}</span>
              </div>
              <div class="flex gap-2">
                <button
                  aria-label="Sebelumnya"
                  onclick={() => pageState.progress--}
                  disabled={pageState.progress === 1}
                  class="p-2 text-sm font-medium rounded-lg transition-colors cursor-pointer flex items-center justify-center {pageState.progress === 1 ? 'bg-white/40 text-slate-400 cursor-not-allowed' : 'bg-white border border-slate-200 text-slate-800 hover:bg-slate-50 hover:text-slate-900'}"
                >
                  <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" /></svg>
                </button>
                <div class="hidden sm:flex items-center gap-1 px-2 flex-wrap justify-center">
                  {#each Array(totalPages.progress) as _, i}
                    <button
                      onclick={() => pageState.progress = i + 1}
                      class="w-8 h-8 flex items-center justify-center text-sm font-medium rounded-lg transition-colors cursor-pointer {pageState.progress === i + 1 ? 'bg-indigo-600 text-slate-100' : 'text-slate-600 hover:bg-white hover:text-slate-900'}"
                    >
                      {i + 1}
                    </button>
                  {/each}
                </div>
                <button
                  aria-label="Selanjutnya"
                  onclick={() => pageState.progress++}
                  disabled={pageState.progress === totalPages.progress}
                  class="p-2 text-sm font-medium rounded-lg transition-colors cursor-pointer flex items-center justify-center {pageState.progress === totalPages.progress ? 'bg-white/40 text-slate-400 cursor-not-allowed' : 'bg-white border border-slate-200 text-slate-800 hover:bg-slate-50 hover:text-slate-900'}"
                >
                  <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" /></svg>
                </button>
              </div>
            </div>
            {/if}
          {/if}
        {/if}

        {#if activeTab === "gameScores"}
          {#if gameScoresData.length === 0}
            <div class="py-12 text-center">
              <h3 class="text-lg font-bold text-slate-700 mb-1">Belum ada skor game.</h3>
            </div>
          {:else}
            <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
              {#each gameScoresData as gs}
                <div class="bg-white border border-slate-200 rounded-2xl p-5 shadow-sm hover:shadow-md transition-shadow flex items-center justify-between">
                  <div>
                    <h3 class="text-lg font-bold text-slate-900">@{gs.user?.username || "Unknown"}</h3>
                    <div class="inline-flex items-center mt-1 px-2 py-0.5 bg-blue-50 text-blue-600 rounded-md text-[10px] font-bold border border-blue-100 uppercase tracking-wide">
                      Mode: {gs.mode}
                    </div>
                  </div>
                  <div class="bg-emerald-50 border border-emerald-100 rounded-xl px-4 py-2 flex flex-col items-center justify-center min-w-[80px]">
                    <span class="text-xl font-black text-emerald-600">{gs.score}</span>
                    <span class="text-[9px] font-bold text-emerald-400 uppercase tracking-wider">Tertinggi</span>
                  </div>
                </div>
              {/each}
            </div>

            <!-- Pagination Controls -->
            {#if totalPages.gameScores > 1}
            <div class="mt-8 px-4 py-3 sm:px-6 sm:py-4 bg-white/60 backdrop-blur-md rounded-2xl border border-slate-200 shadow-sm flex flex-row items-center justify-between gap-4 mb-4">
              <div class="text-sm text-slate-600 text-left">
                Halaman <span class="font-medium text-slate-900">{pageState.gameScores}</span> dari <span class="font-medium text-slate-900">{totalPages.gameScores}</span>
              </div>
              <div class="flex gap-2">
                <button
                  aria-label="Sebelumnya"
                  onclick={() => pageState.gameScores--}
                  disabled={pageState.gameScores === 1}
                  class="p-2 text-sm font-medium rounded-lg transition-colors cursor-pointer flex items-center justify-center {pageState.gameScores === 1 ? 'bg-white/40 text-slate-400 cursor-not-allowed' : 'bg-white border border-slate-200 text-slate-800 hover:bg-slate-50 hover:text-slate-900'}"
                >
                  <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" /></svg>
                </button>
                <div class="hidden sm:flex items-center gap-1 px-2 flex-wrap justify-center">
                  {#each Array(totalPages.gameScores) as _, i}
                    <button
                      onclick={() => pageState.gameScores = i + 1}
                      class="w-8 h-8 flex items-center justify-center text-sm font-medium rounded-lg transition-colors cursor-pointer {pageState.gameScores === i + 1 ? 'bg-indigo-600 text-slate-100' : 'text-slate-600 hover:bg-white hover:text-slate-900'}"
                    >
                      {i + 1}
                    </button>
                  {/each}
                </div>
                <button
                  aria-label="Selanjutnya"
                  onclick={() => pageState.gameScores++}
                  disabled={pageState.gameScores === totalPages.gameScores}
                  class="p-2 text-sm font-medium rounded-lg transition-colors cursor-pointer flex items-center justify-center {pageState.gameScores === totalPages.gameScores ? 'bg-white/40 text-slate-400 cursor-not-allowed' : 'bg-white border border-slate-200 text-slate-800 hover:bg-slate-50 hover:text-slate-900'}"
                >
                  <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" /></svg>
                </button>
              </div>
            </div>
            {/if}
          {/if}
        {/if}

        {#if activeTab === "lessonHistory"}
          {#if lessonHistoryData.length === 0}
            <div class="py-12 text-center">
              <h3 class="text-lg font-bold text-slate-700 mb-1">Belum ada riwayat lesson.</h3>
            </div>
          {:else}
            <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
              {#each lessonHistoryData as lh}
                <div class="bg-white border border-slate-200 rounded-2xl p-5 shadow-sm hover:shadow-md transition-shadow flex flex-col h-full">
                  <div class="flex justify-between items-start mb-3">
                    <div class="flex flex-col">
                      <span class="text-sm font-bold text-slate-400">Lesson {lh.lessonId}</span>
                      <h3 class="text-lg font-bold text-slate-900 mt-0.5">@{lh.user?.username || "Unknown"}</h3>
                    </div>
                    <span class="text-xs font-medium text-slate-500 bg-slate-50 px-2 py-1 rounded-md border border-slate-100">
                      {formatDate(lh.created_at)}
                    </span>
                  </div>
                  
                  <div class="grid grid-cols-2 gap-2 mb-4">
                    <div class="bg-indigo-50/50 border border-indigo-100 rounded-lg p-2 text-center">
                      <span class="text-lg font-bold text-indigo-600 block">{lh.wpm}</span>
                      <span class="text-[9px] font-bold text-indigo-400 uppercase tracking-wider">WPM</span>
                    </div>
                    <div class="bg-blue-50/50 border border-blue-100 rounded-lg p-2 text-center">
                      <span class="text-lg font-bold text-blue-600 block">{lh.accuracy}%</span>
                      <span class="text-[9px] font-bold text-blue-400 uppercase tracking-wider">Akurasi</span>
                    </div>
                  </div>

                  <div class="mt-auto flex justify-center pt-3 border-t border-slate-100">
                    <div class="flex gap-1 text-base">
                      <span class="text-amber-400">{"★".repeat(lh.stars)}</span><span class="text-slate-200">{"★".repeat(3 - lh.stars)}</span>
                    </div>
                  </div>
                </div>
              {/each}
            </div>

            <!-- Pagination Controls -->
            {#if totalPages.lessonHistory > 1}
            <div class="mt-8 px-4 py-3 sm:px-6 sm:py-4 bg-white/60 backdrop-blur-md rounded-2xl border border-slate-200 shadow-sm flex flex-row items-center justify-between gap-4 mb-4">
              <div class="text-sm text-slate-600 text-left">
                Halaman <span class="font-medium text-slate-900">{pageState.lessonHistory}</span> dari <span class="font-medium text-slate-900">{totalPages.lessonHistory}</span>
              </div>
              <div class="flex gap-2">
                <button
                  aria-label="Sebelumnya"
                  onclick={() => pageState.lessonHistory--}
                  disabled={pageState.lessonHistory === 1}
                  class="p-2 text-sm font-medium rounded-lg transition-colors cursor-pointer flex items-center justify-center {pageState.lessonHistory === 1 ? 'bg-white/40 text-slate-400 cursor-not-allowed' : 'bg-white border border-slate-200 text-slate-800 hover:bg-slate-50 hover:text-slate-900'}"
                >
                  <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" /></svg>
                </button>
                <div class="hidden sm:flex items-center gap-1 px-2 flex-wrap justify-center">
                  {#each Array(totalPages.lessonHistory) as _, i}
                    <button
                      onclick={() => pageState.lessonHistory = i + 1}
                      class="w-8 h-8 flex items-center justify-center text-sm font-medium rounded-lg transition-colors cursor-pointer {pageState.lessonHistory === i + 1 ? 'bg-indigo-600 text-slate-100' : 'text-slate-600 hover:bg-white hover:text-slate-900'}"
                    >
                      {i + 1}
                    </button>
                  {/each}
                </div>
                <button
                  aria-label="Selanjutnya"
                  onclick={() => pageState.lessonHistory++}
                  disabled={pageState.lessonHistory === totalPages.lessonHistory}
                  class="p-2 text-sm font-medium rounded-lg transition-colors cursor-pointer flex items-center justify-center {pageState.lessonHistory === totalPages.lessonHistory ? 'bg-white/40 text-slate-400 cursor-not-allowed' : 'bg-white border border-slate-200 text-slate-800 hover:bg-slate-50 hover:text-slate-900'}"
                >
                  <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" /></svg>
                </button>
              </div>
            </div>
            {/if}
          {/if}
        {/if}

        {#if activeTab === "gameHistory"}
          {#if gameHistoryData.length === 0}
            <div class="py-12 text-center">
              <h3 class="text-lg font-bold text-slate-700 mb-1">Belum ada riwayat game.</h3>
            </div>
          {:else}
            <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
              {#each gameHistoryData as gh}
                <div class="bg-white border border-slate-200 rounded-2xl p-5 shadow-sm hover:shadow-md transition-shadow flex items-center justify-between">
                  <div class="flex flex-col">
                    <span class="text-xs font-medium text-slate-500 mb-1.5">{formatDate(gh.created_at)}</span>
                    <h3 class="text-lg font-bold text-slate-900 leading-none">@{gh.user?.username || "Unknown"}</h3>
                    <div class="mt-2.5 w-max px-2.5 py-1 bg-slate-100 text-slate-600 rounded-md text-xs font-bold border border-slate-200 uppercase tracking-wide">
                      Mode: {gh.mode}
                    </div>
                  </div>
                  <div class="bg-blue-50 border border-blue-100 rounded-xl px-4 py-2 flex flex-col items-center justify-center min-w-[70px]">
                    <span class="text-xl font-black text-blue-600">{gh.score}</span>
                    <span class="text-[9px] font-bold text-blue-400 uppercase tracking-wider">Skor</span>
                  </div>
                </div>
              {/each}
            </div>

            <!-- Pagination Controls -->
            {#if totalPages.gameHistory > 1}
            <div class="mt-8 px-4 py-3 sm:px-6 sm:py-4 bg-white/60 backdrop-blur-md rounded-2xl border border-slate-200 shadow-sm flex flex-row items-center justify-between gap-4 mb-4">
              <div class="text-sm text-slate-600 text-left">
                Halaman <span class="font-medium text-slate-900">{pageState.gameHistory}</span> dari <span class="font-medium text-slate-900">{totalPages.gameHistory}</span>
              </div>
              <div class="flex gap-2">
                <button
                  aria-label="Sebelumnya"
                  onclick={() => pageState.gameHistory--}
                  disabled={pageState.gameHistory === 1}
                  class="p-2 text-sm font-medium rounded-lg transition-colors cursor-pointer flex items-center justify-center {pageState.gameHistory === 1 ? 'bg-white/40 text-slate-400 cursor-not-allowed' : 'bg-white border border-slate-200 text-slate-800 hover:bg-slate-50 hover:text-slate-900'}"
                >
                  <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" /></svg>
                </button>
                <div class="hidden sm:flex items-center gap-1 px-2 flex-wrap justify-center">
                  {#each Array(totalPages.gameHistory) as _, i}
                    <button
                      onclick={() => pageState.gameHistory = i + 1}
                      class="w-8 h-8 flex items-center justify-center text-sm font-medium rounded-lg transition-colors cursor-pointer {pageState.gameHistory === i + 1 ? 'bg-indigo-600 text-slate-100' : 'text-slate-600 hover:bg-white hover:text-slate-900'}"
                    >
                      {i + 1}
                    </button>
                  {/each}
                </div>
                <button
                  aria-label="Selanjutnya"
                  onclick={() => pageState.gameHistory++}
                  disabled={pageState.gameHistory === totalPages.gameHistory}
                  class="p-2 text-sm font-medium rounded-lg transition-colors cursor-pointer flex items-center justify-center {pageState.gameHistory === totalPages.gameHistory ? 'bg-white/40 text-slate-400 cursor-not-allowed' : 'bg-white border border-slate-200 text-slate-800 hover:bg-slate-50 hover:text-slate-900'}"
                >
                  <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" /></svg>
                </button>
              </div>
            </div>
            {/if}
          {/if}
        {/if}
      </div>
    </div>
  {/if}
</div>
