<script lang="ts">
  import { onMount } from "svelte";

  type Quiz = {
    id: number;
    title: string;
    category: string;
    timeLimit: number;
  };

  type ScoreQuiz = {
    id: number;
    username: string;
    quiz_id: number;
    quiz?: Quiz;
    score: number;
    created_at: string;
  };

  let isTeacher = $state(false);
  let activeTab = $state<"quizzes" | "scores">("quizzes");
  
  let quizzes: Quiz[] = $state([]);
  let scores: ScoreQuiz[] = $state([]);
  let isLoading = $state(true);
  let isLoadingScores = $state(false);

  // Paginasi State
  let currentQuizPage = $state(1);
  let totalQuizPages = $state(1);
  
  let currentScorePage = $state(1);
  let totalScorePages = $state(1);
  
  let itemsPerPage = 24;
  let isInitialLoad = true;

  async function checkRole() {
    try {
      const res = await fetch(`/me`, {
        credentials: "include",
      });
      const data = await res.json();
      if (data.authenticated && data.user && data.user.role === "teacher") {
        isTeacher = true;
      } else {
        if (!localStorage.getItem("quizzesActiveTab")) {
          activeTab = "scores"; 
        }
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function fetchQuizzes() {
    try {
      const res = await fetch(`/api/quizzes?page=${currentQuizPage}&limit=${itemsPerPage}`, {
        credentials: "include",
      });
      if (res.ok) {
        const json = await res.json();
        quizzes = json.data || [];
        totalQuizPages = json.total_pages || 1;
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function fetchScores() {
    isLoadingScores = true;
    try {
      const url = isTeacher 
        ? `/api/quizzes/scores?page=${currentScorePage}&limit=${itemsPerPage}`
        : `/api/scores/quizzes?page=${currentScorePage}&limit=${itemsPerPage}`;
      
      const res = await fetch(url, { credentials: "include" });
      if (res.ok) {
        const json = await res.json();
        scores = json.data || [];
        totalScorePages = json.total_pages || 1;
      }
    } catch (e) {
      console.error(e);
    } finally {
      isLoadingScores = false;
    }
  }

  async function deleteQuiz(id: number) {
    if (!confirm("Yakin ingin menghapus kuis ini?")) return;
    try {
      const res = await fetch(`/api/quizzes/${id}`, {
        method: "DELETE",
        credentials: "include",
      });
      if (res.ok) {
        quizzes = quizzes.filter(q => q.id !== id);
      } else {
        alert("Gagal menghapus kuis");
      }
    } catch (e) {
      console.error(e);
      alert("Terjadi kesalahan");
    }
  }

  $effect(() => {
    if (!isInitialLoad) {
      fetchQuizzes();
    }
  });

  $effect(() => {
    if (!isInitialLoad) {
      fetchScores();
    }
  });

  onMount(async () => {
    const savedTab = localStorage.getItem("quizzesActiveTab");
    if (savedTab === "quizzes" || savedTab === "scores") {
      activeTab = savedTab;
    }

    await checkRole();
    await fetchQuizzes();
    await fetchScores();
    isInitialLoad = false;
    isLoading = false;
  });

  $effect(() => {
    if (!isLoading) {
      localStorage.setItem("quizzesActiveTab", activeTab);
    }
  });
</script>

<svelte:head>
  <title>{isTeacher ? 'Manajemen Kuis' : 'Nilai Kuis'} | Les Balongarut</title>
</svelte:head>

<div class="mb-8 flex flex-col md:flex-row items-start md:justify-between md:items-end gap-4">
  <div>
    <h1 class="text-3xl font-bold text-slate-900 mb-2">Kuis</h1>
    <p class="text-slate-600">
      {isTeacher ? 'Kelola kuis dan pantau nilai siswa.' : 'Lihat daftar kuis dan riwayat nilaimu.'}
    </p>
  </div>
  {#if isTeacher && activeTab === "quizzes"}
    <a
      href="/dashboard/quizzes/create"
      class="inline-flex items-center gap-2 px-4 py-2 text-sm font-bold text-blue-700 bg-blue-50 hover:bg-blue-100 border border-blue-200 rounded-xl transition-all shadow-sm no-underline"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
      </svg>
      Tambah Kuis
    </a>
  {/if}
</div>

<!-- Tabs Navigation -->
<div class="mb-6 flex">
  <div class="grid grid-cols-2 w-full md:w-max md:flex p-1.5 gap-1.5 bg-slate-100/80 backdrop-blur-sm rounded-2xl">
    <button
      class="px-4 py-2.5 rounded-xl text-sm font-bold transition-all duration-200 {activeTab === 'quizzes' ? 'bg-white text-blue-600 shadow-sm ring-1 ring-slate-900/5' : 'text-slate-500 hover:text-slate-700 hover:bg-slate-200/50'}"
      onclick={() => activeTab = 'quizzes'}
    >
      Daftar Kuis
    </button>
    <button
      class="px-4 py-2.5 rounded-xl text-sm font-bold transition-all duration-200 {activeTab === 'scores' ? 'bg-white text-blue-600 shadow-sm ring-1 ring-slate-900/5' : 'text-slate-500 hover:text-slate-700 hover:bg-slate-200/50'}"
      onclick={() => activeTab = 'scores'}
    >
      Riwayat Nilai
    </button>
  </div>
</div>

{#if isLoading}
  <div class="flex justify-center p-12">
    <div class="w-8 h-8 border-4 border-slate-200 border-t-blue-600 rounded-full animate-spin"></div>
  </div>
{:else}

  {#if activeTab === "quizzes"}
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-6">
      {#each quizzes as quiz}
        <div class="bg-white border border-slate-200 rounded-2xl p-5 shadow-sm hover:shadow-md transition-shadow relative flex flex-col h-full">
          <div class="flex justify-between items-start mb-2">
            <div>
              <span class="text-xs font-bold text-slate-400">#{quiz.id}</span>
              <h3 class="text-lg font-bold text-slate-900 mt-1 leading-tight">{quiz.title}</h3>
            </div>
            <span class="px-2 py-1 bg-slate-100 text-slate-600 rounded-md text-[10px] font-bold border border-slate-200 whitespace-nowrap uppercase tracking-wide">
              {quiz.category}
            </span>
          </div>
          
          <div class="flex items-center text-xs font-medium text-slate-500 mb-6 bg-slate-50 w-max px-2.5 py-1.5 rounded-md border border-slate-100">
            <svg class="w-4 h-4 mr-1.5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
            {quiz.timeLimit} Menit
          </div>

          <div class="mt-auto flex gap-2 pt-4 border-t border-slate-100">
            {#if isTeacher}
              <a
                href="/dashboard/quizzes/{quiz.id}/edit"
                class="flex-1 inline-flex items-center justify-center p-2.5 text-sm font-semibold text-blue-700 bg-blue-50 rounded-xl border border-blue-100 hover:bg-blue-100 transition-colors no-underline"
              >
                <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                </svg>
                Edit
              </a>
              <button
                onclick={() => deleteQuiz(quiz.id)}
                class="flex-1 inline-flex items-center justify-center p-2.5 text-sm font-semibold text-red-600 bg-red-50 rounded-xl border border-red-100 hover:bg-red-100 transition-colors cursor-pointer"
              >
                <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
                Hapus
              </button>
            {:else}
              <a
                href="/quiz/{quiz.id}"
                class="w-full inline-flex items-center justify-center px-4 py-2.5 text-sm font-bold text-blue-700 bg-blue-50 rounded-xl border border-blue-100 hover:bg-blue-100 transition-colors no-underline"
              >
                Kerjakan Kuis
              </a>
            {/if}
          </div>
        </div>
      {:else}
        <div class="col-span-full text-center p-12 bg-white rounded-2xl border border-slate-200 shadow-sm text-sm text-slate-500">Belum ada kuis yang tersedia.</div>
      {/each}
    </div>

    <!-- Pagination Controls Kuis -->
    {#if totalQuizPages > 1}
      <div class="flex items-center justify-between bg-white px-4 py-3 border border-slate-200 rounded-xl shadow-sm">
        <div class="text-sm text-slate-600">
          Halaman <span class="font-bold text-slate-900">{currentQuizPage}</span> dari <span class="font-bold text-slate-900">{totalQuizPages}</span>
        </div>
        <div class="flex items-center gap-2">
          <button
            class="p-2 rounded-lg border border-slate-200 text-slate-600 hover:bg-slate-100 disabled:opacity-50 disabled:cursor-not-allowed transition-colors cursor-pointer"
            onclick={() => currentQuizPage = Math.max(1, currentQuizPage - 1)}
            disabled={currentQuizPage === 1}
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" /></svg>
          </button>
          
          <button
            class="p-2 rounded-lg border border-slate-200 text-slate-600 hover:bg-slate-100 disabled:opacity-50 disabled:cursor-not-allowed transition-colors cursor-pointer"
            onclick={() => currentQuizPage = Math.min(totalQuizPages, currentQuizPage + 1)}
            disabled={currentQuizPage === totalQuizPages}
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" /></svg>
          </button>
        </div>
      </div>
    {/if}
  {/if}

  {#if activeTab === "scores"}
    {#if isLoadingScores}
      <div class="flex justify-center p-8">
        <div class="w-6 h-6 border-4 border-slate-200 border-t-blue-600 rounded-full animate-spin"></div>
      </div>
    {:else}
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-6">
        {#each scores as s}
          <div class="bg-white border border-slate-200 rounded-2xl p-5 shadow-sm hover:shadow-md transition-shadow relative flex flex-col h-full">
            <div class="flex justify-between items-start mb-4">
              <div class="pr-4">
                <h3 class="text-lg font-bold text-slate-900 leading-tight">
                  {s.quiz ? s.quiz.title : `Kuis #${s.quiz_id}`}
                </h3>
                {#if isTeacher}
                  <div class="flex items-center text-sm text-slate-500 mt-2 font-medium bg-slate-50 px-2.5 py-1 rounded-md w-max border border-slate-100">
                    <svg class="w-4 h-4 mr-1.5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" /></svg>
                    @{s.username}
                  </div>
                {/if}
              </div>
              <div class="shrink-0 text-right mt-0.5">
                <span class={`inline-block px-3 py-1.5 rounded-xl text-sm font-black border shadow-sm tracking-wide ${s.score >= 80 ? 'bg-green-50 text-green-700 border-green-200' : s.score >= 60 ? 'bg-yellow-50 text-yellow-700 border-yellow-200' : 'bg-red-50 text-red-600 border-red-200'}`}>
                  SKOR {s.score}
                </span>
              </div>
            </div>
            
            <div class="mt-auto flex items-center justify-between pt-4 border-t border-slate-100">
              <div class="flex items-center text-xs text-slate-600 font-semibold bg-slate-100 px-3 py-1.5 rounded-lg border border-slate-200/60">
                <svg class="w-4 h-4 mr-1.5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg>
                {new Date(s.created_at).toLocaleDateString("id-ID", { day: "numeric", month: "short", year: "numeric" })}
              </div>
              <div class="text-xs text-slate-500 font-medium bg-slate-50 px-2.5 py-1.5 rounded-lg border border-slate-100">
                {new Date(s.created_at).toLocaleTimeString("id-ID", { hour: "2-digit", minute: "2-digit" })}
              </div>
            </div>
          </div>
        {:else}
          <div class="col-span-full text-center p-12 bg-white rounded-2xl border border-slate-200 shadow-sm text-sm text-slate-500">Belum ada riwayat pengerjaan kuis.</div>
        {/each}
      </div>

      <!-- Pagination Controls Skor -->
      {#if totalScorePages > 1}
        <div class="flex items-center justify-between bg-white px-4 py-3 border border-slate-200 rounded-xl shadow-sm">
          <div class="text-sm text-slate-600">
            Halaman <span class="font-bold text-slate-900">{currentScorePage}</span> dari <span class="font-bold text-slate-900">{totalScorePages}</span>
          </div>
          <div class="flex items-center gap-2">
            <button
              class="p-2 rounded-lg border border-slate-200 text-slate-600 hover:bg-slate-100 disabled:opacity-50 disabled:cursor-not-allowed transition-colors cursor-pointer"
              onclick={() => currentScorePage = Math.max(1, currentScorePage - 1)}
              disabled={currentScorePage === 1}
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" /></svg>
            </button>
            
            <button
              class="p-2 rounded-lg border border-slate-200 text-slate-600 hover:bg-slate-100 disabled:opacity-50 disabled:cursor-not-allowed transition-colors cursor-pointer"
              onclick={() => currentScorePage = Math.min(totalScorePages, currentScorePage + 1)}
              disabled={currentScorePage === totalScorePages}
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" /></svg>
            </button>
          </div>
        </div>
      {/if}
    {/if}
  {/if}

{/if}
