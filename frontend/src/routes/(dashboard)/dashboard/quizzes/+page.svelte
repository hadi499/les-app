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

  async function checkRole() {
    try {
      const res = await fetch(`/me`, {
        credentials: "include",
      });
      const data = await res.json();
      if (data.authenticated && data.user && data.user.role === "teacher") {
        isTeacher = true;
      } else {
        // Jika murid, mungkin tab default langsung ke skor jika tidak mau akses manajemen kuis.
        // Tapi kita biarkan mereka bisa lihat daftar kuis juga (optional).
        // Kita juga bisa sembunyikan Manajemen Kuis dari murid kalau memang hanya untuk guru.
        // Untuk sekarang, murid diarahkan default ke "scores" atau daftar kuis.
        activeTab = "scores"; // Asumsi student lebih peduli nilainya
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function fetchQuizzes() {
    try {
      const res = await fetch(`/api/quizzes`, {
        credentials: "include",
      });
      if (res.ok) {
        const json = await res.json();
        quizzes = json.data || [];
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function fetchScores() {
    isLoadingScores = true;
    try {
      const url = isTeacher 
        ? `/api/quizzes/scores`
        : `/api/scores/quizzes`;
      
      const res = await fetch(url, { credentials: "include" });
      if (res.ok) {
        const json = await res.json();
        scores = json.data || [];
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

  onMount(async () => {
    await checkRole();
    await fetchQuizzes();
    await fetchScores();
    isLoading = false;
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
      class="inline-flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg font-medium hover:bg-slate-500 transition-colors shadow-sm no-underline"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
      </svg>
      Tambah Kuis
    </a>
  {/if}
</div>

<!-- Tabs Navigation -->
<div class="flex gap-2 mb-6 border-b border-slate-200 overflow-x-auto whitespace-nowrap">
  <button
    class="px-6 py-3 font-semibold text-sm border-b-2 transition-colors cursor-pointer {activeTab === 'quizzes' ? 'border-blue-600 text-slate-800' : 'border-transparent text-slate-500 hover:text-slate-800'}"
    onclick={() => activeTab = 'quizzes'}
  >
    Daftar Kuis
  </button>
  <button
    class="px-6 py-3 font-semibold text-sm border-b-2 transition-colors cursor-pointer {activeTab === 'scores' ? 'border-blue-600 text-slate-800' : 'border-transparent text-slate-500 hover:text-slate-800'}"
    onclick={() => activeTab = 'scores'}
  >
    Riwayat Nilai
  </button>
</div>

{#if isLoading}
  <div class="flex justify-center p-12">
    <div class="w-8 h-8 border-4 border-slate-200 border-t-blue-600 rounded-full animate-spin"></div>
  </div>
{:else}

  {#if activeTab === "quizzes"}
    <div class="bg-white rounded-2xl shadow-sm border border-slate-200 overflow-hidden">
      <div class="overflow-x-auto w-full">
        <table class="w-full text-left border-collapse whitespace-nowrap min-w-[700px]">
          <thead>
            <tr class="bg-slate-50/50 border-b border-slate-200">
              <th class="px-6 py-4 text-xs font-bold text-slate-800 uppercase tracking-wider">ID</th>
              <th class="px-6 py-4 text-xs font-bold text-slate-800 uppercase tracking-wider">Judul Kuis</th>
              <th class="px-6 py-4 text-xs font-bold text-slate-800 uppercase tracking-wider">Kategori</th>
              <th class="px-6 py-4 text-xs font-bold text-slate-800 uppercase tracking-wider">Batas Waktu</th>
              <th class="px-6 py-4 text-xs font-bold text-slate-800 uppercase tracking-wider text-center">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            {#each quizzes as quiz}
              <tr class="hover:bg-slate-50/50 transition-colors">
                <td class="px-6 py-4 text-sm font-medium text-slate-900">#{quiz.id}</td>
                <td class="px-6 py-4 text-sm font-semibold text-slate-800">{quiz.title}</td>
                <td class="px-6 py-4 text-sm text-slate-600">
                  <span class="px-2 py-1 bg-slate-100 text-slate-600 rounded-md text-xs font-medium border border-slate-200">
                    {quiz.category}
                  </span>
                </td>
                <td class="px-6 py-4 text-sm text-slate-600">{quiz.timeLimit} Menit</td>
                <td class="px-6 py-4 text-center space-x-2">
                  {#if isTeacher}
                    <a
                      href="/dashboard/quizzes/{quiz.id}/edit"
                      class="inline-flex items-center justify-center p-2 text-blue-600 bg-blue-50 rounded-lg hover:bg-blue-100 transition-colors"
                      title="Edit Kuis"
                    >
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                      </svg>
                    </a>
                    <button
                      onclick={() => deleteQuiz(quiz.id)}
                      class="inline-flex items-center justify-center p-2 text-red-600 bg-red-50 rounded-lg hover:bg-red-100 transition-colors cursor-pointer border-none"
                      title="Hapus Kuis"
                    >
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                      </svg>
                    </button>
                  {:else}
                    <a
                      href="/quiz/{quiz.id}"
                      class="inline-flex items-center justify-center px-4 py-2 text-xs font-bold text-white bg-blue-500 rounded-lg hover:bg-blue-600 transition-colors no-underline"
                    >
                      Kerjakan Kuis
                    </a>
                  {/if}
                </td>
              </tr>
            {:else}
              <tr>
                <td colspan="5" class="px-6 py-8 text-center text-sm text-slate-600">
                  Belum ada kuis yang tersedia.
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>
  {/if}

  {#if activeTab === "scores"}
    {#if isLoadingScores}
      <div class="flex justify-center p-8">
        <div class="w-6 h-6 border-4 border-slate-200 border-t-blue-600 rounded-full animate-spin"></div>
      </div>
    {:else}
      <div class="bg-white rounded-2xl shadow-sm border border-slate-200 overflow-hidden">
        <div class="overflow-x-auto w-full">
          <table class="w-full text-left border-collapse whitespace-nowrap min-w-[500px]">
            <thead>
              <tr class="bg-slate-50/50 border-b border-slate-200">
                <th class="px-6 py-4 text-xs font-bold text-slate-800 uppercase tracking-wider">Kuis</th>
                {#if isTeacher}
                  <th class="px-6 py-4 text-xs font-bold text-slate-800 uppercase tracking-wider">Siswa</th>
                {/if}
                <th class="px-6 py-4 text-xs font-bold text-slate-800 uppercase tracking-wider">Skor</th>
                <th class="px-6 py-4 text-xs font-bold text-slate-800 uppercase tracking-wider">Waktu Selesai</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              {#each scores as s}
                <tr class="hover:bg-slate-50/50 transition-colors">
                  <td class="px-6 py-4 text-sm font-semibold text-slate-800">
                    {s.quiz ? s.quiz.title : `Kuis #${s.quiz_id}`}
                  </td>
                  {#if isTeacher}
                    <td class="px-6 py-4 text-sm text-slate-600 font-medium">@{s.username}</td>
                  {/if}
                  <td class="px-6 py-4">
                    <span class={`px-3 py-1 rounded-full text-xs font-bold border ${s.score >= 80 ? 'bg-green-100 text-green-700 border-green-300' : s.score >= 60 ? 'bg-yellow-100 text-yellow-700 border-yellow-300' : 'bg-red-100 text-red-600 border-red-300'}`}>
                      {s.score}
                    </span>
                  </td>
                  <td class="px-6 py-4 text-sm text-slate-600">
                    {new Date(s.created_at).toLocaleString("id-ID", {
                      day: "numeric", month: "short", year: "numeric", hour: "2-digit", minute: "2-digit"
                    })}
                  </td>
                </tr>
              {:else}
                <tr>
                  <td colspan={isTeacher ? 4 : 3} class="px-6 py-8 text-center text-sm text-slate-600">
                    Belum ada riwayat pengerjaan kuis.
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>
    {/if}
  {/if}

{/if}
