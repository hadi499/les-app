<script lang="ts">
  import { onMount } from "svelte";

  type UserRef = { username: string };
  type Progress = { user?: UserRef; lessonId: number; bestWPM: number; bestAccuracy: number; stars: number; completed: boolean; attempts: number };
  type GameScore = { user?: UserRef; mode: string; score: number };
  type LessonHistory = { created_at: string; user?: UserRef; lessonId: number; wpm: number; accuracy: number; stars: number };
  type GameHistory = { created_at: string; user?: UserRef; mode: string; score: number };

  let activeTab = $state("progress");
  let isLoading = $state(true);
  let errorMsg = $state("");

  let progressData: Progress[] = $state([]);
  let gameScoresData: GameScore[] = $state([]);
  let lessonHistoryData: LessonHistory[] = $state([]);
  let gameHistoryData: GameHistory[] = $state([]);

  async function fetchData() {
    isLoading = true;
    errorMsg = "";
    try {
      const [resProgress, resGameScores, resLessonHistory, resGameHistory] =
        await Promise.all([
          fetch(`/api/typing/admin/progress`, {
            credentials: "include",
          }),
          fetch(`/api/typing/admin/game-scores`, {
            credentials: "include",
          }),
          fetch(`/api/typing/admin/history/lesson`, {
            credentials: "include",
          }),
          fetch(`/api/typing/admin/history/game`, {
            credentials: "include",
          }),
        ]);

      if (resProgress.status === 403) {
        errorMsg = "Akses ditolak. Anda bukan Guru.";
        return;
      }

      if (
        !resProgress.ok ||
        !resGameScores.ok ||
        !resLessonHistory.ok ||
        !resGameHistory.ok
      ) {
        throw new Error("Gagal mengambil data dari server");
      }

      progressData = (await resProgress.json() as Progress[]) || [];
      gameScoresData = (await resGameScores.json() as GameScore[]) || [];
      lessonHistoryData = (await resLessonHistory.json() as LessonHistory[]) || [];
      gameHistoryData = (await resGameHistory.json() as GameHistory[]) || [];
    } catch (e) {
      errorMsg = e instanceof Error ? e.message : String(e);
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    fetchData();
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
      class="text-2xl font-bold text-orange-950 sm:text-3xl tracking-tight drop-shadow-sm"
    >
      Pemantauan Mengetik
    </h1>
    <p class="mt-2 text-orange-800 text-md tracking-wide">
      Pantau perkembangan kemampuan mengetik seluruh murid.
    </p>
  </div>

  {#if isLoading}
    <div class="flex justify-center p-12">
      <div
        class="w-10 h-10 border-4 border-orange-200 border-t-orange-500 rounded-full animate-spin"
      ></div>
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
    <div class="flex flex-wrap gap-2 border-b border-orange-200 pb-2">
      {#each tabs as tab}
        <button
          class="px-4 py-2 rounded-t-lg text-sm font-medium transition-all duration-200 border-b-2 {activeTab ===
          tab.id
            ? 'border-indigo-500 text-indigo-600 bg-indigo-800/10'
            : 'border-transparent text-orange-800  hover:bg-white/50'}"
          onclick={() => (activeTab = tab.id)}
        >
          {tab.label}
        </button>
      {/each}
    </div>

    <!-- Content -->
    <div
      class="bg-white/60 backdrop-blur-md rounded-b-3xl rounded-tr-3xl border border-orange-200 shadow-lg shadow-orange-900/10 overflow-hidden"
    >
      <div class="overflow-x-auto">
        {#if activeTab === "progress"}
          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="bg-white/40 border-b border-orange-200">
                <th class="py-4 px-6 font-bold text-orange-950 text-sm"
                  >Murid</th
                >
                <th
                  class="py-4 px-6 font-bold text-orange-950 text-sm text-center"
                  >Lesson ID</th
                >
                <th
                  class="py-4 px-6 font-bold text-orange-950 text-sm text-center"
                  >Best WPM</th
                >
                <th
                  class="py-4 px-6 font-bold text-orange-950 text-sm text-center"
                  >Akurasi</th
                >
                <th
                  class="py-4 px-6 font-bold text-orange-950 text-sm text-center"
                  >Bintang</th
                >
                <th
                  class="py-4 px-6 font-bold text-orange-950 text-sm text-center"
                  >Selesai</th
                >
                <th
                  class="py-4 px-6 font-bold text-orange-950 text-sm text-center"
                  >Percobaan</th
                >
              </tr>
            </thead>
            <tbody class="divide-y divide-orange-200">
              {#each progressData as p}
                <tr class="hover:bg-white/40 transition-colors">
                  <td class="py-3 px-6 text-sm font-medium text-orange-950"
                    >{p.user?.username || "Unknown"}</td
                  >
                  <td class="py-3 px-6 text-sm text-orange-900 text-center"
                    >{p.lessonId}</td
                  >
                  <td
                    class="py-3 px-6 text-sm text-indigo-500 font-bold text-center"
                    >{p.bestWPM}</td
                  >
                  <td class="py-3 px-6 text-sm text-orange-900 text-center"
                    >{p.bestAccuracy}%</td
                  >
                  <td class="py-3 px-6 text-sm text-yellow-500 text-center"
                    >{"★".repeat(p.stars)}{"☆".repeat(3 - p.stars)}</td
                  >
                  <td class="py-3 px-6 text-sm text-center">
                    <span
                      class={`px-2 py-1 rounded-full text-xs font-medium border ${p.completed ? "bg-green-100 text-green-600 border-green-800/50" : "bg-white text-orange-800 border-orange-300"}`}
                    >
                      {p.completed ? "Ya" : "Belum"}
                    </span>
                  </td>
                  <td class="py-3 px-6 text-sm text-orange-800 text-center"
                    >{p.attempts}</td
                  >
                </tr>
              {/each}
              {#if progressData.length === 0}
                <tr
                  ><td colspan="7" class="py-8 text-center text-orange-700"
                    >Belum ada data progress.</td
                  ></tr
                >
              {/if}
            </tbody>
          </table>
        {/if}

        {#if activeTab === "gameScores"}
          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="bg-white/40 border-b border-orange-200">
                <th class="py-4 px-6 font-bold text-orange-950 text-sm"
                  >Murid</th
                >
                <th
                  class="py-4 px-6 font-bold text-orange-950 text-sm text-center"
                  >Mode</th
                >
                <th
                  class="py-4 px-6 font-bold text-orange-950 text-sm text-center"
                  >Skor Tertinggi</th
                >
              </tr>
            </thead>
            <tbody class="divide-y divide-orange-200">
              {#each gameScoresData as gs}
                <tr class="hover:bg-white/40 transition-colors">
                  <td class="py-3 px-6 text-sm font-medium text-orange-950"
                    >{gs.user?.username || "Unknown"}</td
                  >
                  <td
                    class="py-3 px-6 text-sm text-orange-900 text-center capitalize"
                    >{gs.mode}</td
                  >
                  <td
                    class="py-3 px-6 text-sm text-green-400 font-bold text-center"
                    >{gs.score}</td
                  >
                </tr>
              {/each}
              {#if gameScoresData.length === 0}
                <tr
                  ><td colspan="3" class="py-8 text-center text-orange-700"
                    >Belum ada skor game.</td
                  ></tr
                >
              {/if}
            </tbody>
          </table>
        {/if}

        {#if activeTab === "lessonHistory"}
          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="bg-white/40 border-b border-orange-200">
                <th class="py-4 px-6 font-bold text-orange-950 text-sm"
                  >Tanggal</th
                >
                <th class="py-4 px-6 font-bold text-orange-950 text-sm"
                  >Murid</th
                >
                <th
                  class="py-4 px-6 font-bold text-orange-950 text-sm text-center"
                  >Lesson ID</th
                >
                <th
                  class="py-4 px-6 font-bold text-orange-950 text-sm text-center"
                  >WPM</th
                >
                <th
                  class="py-4 px-6 font-bold text-orange-950 text-sm text-center"
                  >Akurasi</th
                >
                <th
                  class="py-4 px-6 font-bold text-orange-950 text-sm text-center"
                  >Bintang</th
                >
              </tr>
            </thead>
            <tbody class="divide-y divide-orange-200">
              {#each lessonHistoryData as lh}
                <tr class="hover:bg-white/40 transition-colors">
                  <td class="py-3 px-6 text-xs text-orange-800"
                    >{formatDate(lh.created_at)}</td
                  >
                  <td class="py-3 px-6 text-sm font-medium text-orange-950"
                    >{lh.user?.username || "Unknown"}</td
                  >
                  <td class="py-3 px-6 text-sm text-orange-900 text-center"
                    >{lh.lessonId}</td
                  >
                  <td class="py-3 px-6 text-sm text-orange-900 text-center"
                    >{lh.wpm}</td
                  >
                  <td class="py-3 px-6 text-sm text-orange-900 text-center"
                    >{lh.accuracy}%</td
                  >
                  <td class="py-3 px-6 text-sm text-yellow-500 text-center"
                    >{"★".repeat(lh.stars)}{"☆".repeat(3 - lh.stars)}</td
                  >
                </tr>
              {/each}
              {#if lessonHistoryData.length === 0}
                <tr
                  ><td colspan="6" class="py-8 text-center text-orange-700"
                    >Belum ada riwayat lesson.</td
                  ></tr
                >
              {/if}
            </tbody>
          </table>
        {/if}

        {#if activeTab === "gameHistory"}
          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="bg-white/40 border-b border-orange-200">
                <th class="py-4 px-6 font-bold text-orange-950 text-sm"
                  >Tanggal</th
                >
                <th class="py-4 px-6 font-bold text-orange-950 text-sm"
                  >Murid</th
                >
                <th
                  class="py-4 px-6 font-bold text-orange-950 text-sm text-center"
                  >Mode</th
                >
                <th
                  class="py-4 px-6 font-bold text-orange-950 text-sm text-center"
                  >Skor</th
                >
              </tr>
            </thead>
            <tbody class="divide-y divide-orange-200">
              {#each gameHistoryData as gh}
                <tr class="hover:bg-white/40 transition-colors">
                  <td class="py-3 px-6 text-xs text-orange-800"
                    >{formatDate(gh.created_at)}</td
                  >
                  <td class="py-3 px-6 text-sm font-medium text-orange-950"
                    >{gh.user?.username || "Unknown"}</td
                  >
                  <td
                    class="py-3 px-6 text-sm text-orange-900 text-center capitalize"
                    >{gh.mode}</td
                  >
                  <td
                    class="py-3 px-6 text-sm text-orange-900 font-bold text-center"
                    >{gh.score}</td
                  >
                </tr>
              {/each}
              {#if gameHistoryData.length === 0}
                <tr
                  ><td colspan="4" class="py-8 text-center text-orange-700"
                    >Belum ada riwayat game.</td
                  ></tr
                >
              {/if}
            </tbody>
          </table>
        {/if}
      </div>
    </div>
  {/if}
</div>
