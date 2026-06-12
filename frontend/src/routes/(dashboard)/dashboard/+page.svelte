<script>
  import { onMount } from "svelte";
  import {
    fetchAllProgress,
    fetchAllGameScores,
    fetchHistory,
    progressMap,
    gameScores,
    gameHistoryStore,
    lessonHistoryStore,
  } from "$lib/stores/progress";

  onMount(async () => {
    await fetchAllProgress();
    await fetchAllGameScores();
    await fetchHistory();
  });

  let completedLessons = $derived(
    Array.from($progressMap.values()).filter((p) => p.completed).length,
  );
  let totalStars = $derived(
    Array.from($progressMap.values()).reduce((sum, p) => sum + p.stars, 0),
  );
</script>

<svelte:head>
  <title>Dashboard - Portal Les Balongarut</title>
</svelte:head>

<div class="animate-in fade-in duration-500">
  <!-- Header -->
  <div class="mb-8">
    <h1
      class="text-2xl font-bold text-white sm:text-3xl tracking-tight drop-shadow-sm"
    >
      Selamat datang kembali! 👋
    </h1>
  </div>

  <!-- Main Content -->
  <div class="max-w-4xl mx-auto space-y-6">
    <!-- Typing Progress Section -->
    <div
      class="bg-zinc-900/50 backdrop-blur-md rounded-3xl border border-zinc-800 shadow-lg shadow-black/30 overflow-hidden"
    >
      <div
        class="px-6 py-5 border-b border-zinc-800 flex justify-between items-center bg-zinc-800/50"
      >
        <h3 class="text-lg font-medium text-white drop-shadow-sm m-0">
          Progress Mengetik 10 Jari
        </h3>
        <a
          href="/mengetik"
          class="text-sm font-medium text-blue-300 hover:text-white bg-blue-900/30 border border-blue-800/50 hover:bg-blue-800/50 px-3 py-1.5 rounded-lg transition-colors no-underline"
          >Lanjutkan Latihan</a
        >
      </div>

      <div class="p-6">
        <h4 class="text-sm font-light text-zinc-400 uppercase tracking-wider mb-4 m-0">Statistik Pelajaran</h4>
        <div class="grid grid-cols-2 gap-4 mb-8">
          <div
            class="bg-amber-900/20 rounded-2xl p-4 flex flex-col items-center justify-center border border-amber-800/50 hover:-translate-y-1 hover:shadow-lg hover:shadow-amber-900/20 transition-all"
          >
            <span class="text-3xl mb-2 drop-shadow-sm">⭐</span>
            <span class="text-2xl font-bold text-amber-400 drop-shadow-sm">{totalStars}</span>
            <span
              class="text-xs font-light text-amber-200/80 uppercase tracking-wide mt-1"
              >Total Bintang</span
            >
          </div>
          <div
            class="bg-emerald-900/20 rounded-2xl p-4 flex flex-col items-center justify-center border border-emerald-800/50 hover:-translate-y-1 hover:shadow-lg hover:shadow-emerald-900/20 transition-all"
          >
            <span class="text-3xl mb-2 drop-shadow-sm">📚</span>
            <span class="text-2xl font-bold text-emerald-400 drop-shadow-sm"
              >{completedLessons}</span
            >
            <span
              class="text-xs font-light text-emerald-200/80 uppercase tracking-wide text-center mt-1"
              >Pelajaran Selesai</span
            >
          </div>
        </div>

        <div class="mt-4 pt-4 border-t border-zinc-800/50 flex justify-end">
          <a
            href="/dashboard/card-memory"
            class="text-sm font-medium text-blue-200 hover:text-white bg-zinc-800 border border-zinc-700 hover:bg-zinc-700 px-4 py-2 rounded-lg transition-colors no-underline flex items-center gap-2 shadow-sm"
          >
            <span>Buka Card Memory</span>
            <span class="text-lg drop-shadow-sm">🎴</span>
          </a>
        </div>

        <h4 class="text-sm font-light text-zinc-400 uppercase tracking-wider mb-4 mt-8 m-0">High Score Game Mengetik</h4>
        <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-4">
          <div
            class="bg-rose-900/20 rounded-2xl p-4 flex flex-col items-center justify-center border border-rose-800/50 hover:-translate-y-1 hover:shadow-lg hover:shadow-rose-900/20 transition-all"
          >
            <span class="text-3xl mb-2 drop-shadow-sm">👈</span>
            <span class="text-2xl font-bold text-rose-400 drop-shadow-sm"
              >{$gameScores.left || 0}</span
            >
            <span
              class="text-[0.65rem] font-light text-rose-200/80 uppercase tracking-wide text-center mt-1"
              >Tangan Kiri</span
            >
          </div>
          <div
            class="bg-teal-900/20 rounded-2xl p-4 flex flex-col items-center justify-center border border-teal-800/50 hover:-translate-y-1 hover:shadow-lg hover:shadow-teal-900/20 transition-all"
          >
            <span class="text-3xl mb-2 drop-shadow-sm">👉</span>
            <span class="text-2xl font-bold text-teal-400 drop-shadow-sm"
              >{$gameScores.right || 0}</span
            >
            <span
              class="text-[0.65rem] font-light text-teal-200/80 uppercase tracking-wide text-center mt-1"
              >Tangan Kanan</span
            >
          </div>
          <div
            class="bg-blue-900/20 rounded-2xl p-4 flex flex-col items-center justify-center border border-blue-800/50 hover:-translate-y-1 hover:shadow-lg hover:shadow-blue-900/20 transition-all"
          >
            <span class="text-3xl mb-2 drop-shadow-sm">🎯</span>
            <span class="text-2xl font-bold text-blue-400 drop-shadow-sm"
              >{$gameScores.both || 0}</span
            >
            <span
              class="text-[0.65rem] font-light text-blue-200/80 uppercase tracking-wide text-center mt-1"
              >Home Row</span
            >
          </div>
          <div
            class="bg-orange-900/20 rounded-2xl p-4 flex flex-col items-center justify-center border border-orange-800/50 hover:-translate-y-1 hover:shadow-lg hover:shadow-orange-900/20 transition-all"
          >
            <span class="text-3xl mb-2 drop-shadow-sm">🔤</span>
            <span class="text-2xl font-bold text-orange-400 drop-shadow-sm"
              >{$gameScores.letters || 0}</span
            >
            <span
              class="text-[0.65rem] font-light text-orange-200/80 uppercase tracking-wide text-center mt-1"
              >Semua Huruf</span
            >
          </div>
          <div
            class="bg-purple-900/20 rounded-2xl p-4 flex flex-col items-center justify-center border border-purple-800/50 hover:-translate-y-1 hover:shadow-lg hover:shadow-purple-900/20 transition-all"
          >
            <span class="text-3xl mb-2 drop-shadow-sm">🚀</span>
            <span class="text-2xl font-bold text-purple-400 drop-shadow-sm"
              >{$gameScores.all || 0}</span
            >
            <span
              class="text-[0.65rem] font-light text-purple-200/80 uppercase tracking-wide text-center mt-1"
              >Semua Tombol</span
            >
          </div>
        </div>
      </div>
    </div>

    <!-- History Section -->
    <div
      class="bg-zinc-900/50 backdrop-blur-md rounded-3xl border border-zinc-800 shadow-lg shadow-black/30 overflow-hidden"
    >
      <div class="px-6 py-5 border-b border-zinc-800 bg-zinc-800/50">
        <h3 class="text-lg font-medium text-white drop-shadow-sm m-0">
          Riwayat Aktivitas Mengetik
        </h3>
      </div>

      <div class="divide-y divide-zinc-800/50 max-h-[400px] overflow-y-auto">
        {#if $gameHistoryStore.length === 0 && $lessonHistoryStore.length === 0}
          <div class="p-8 text-center text-zinc-500 font-light">
            Belum ada riwayat aktivitas. Ayo mulai mengetik!
          </div>
        {:else}
          {@const allHistory = [
            ...$gameHistoryStore.map((h) => ({ type: "game", ...h })),
            ...$lessonHistoryStore.map((h) => ({ type: "lesson", ...h })),
          ]
            .sort(
              (a, b) =>
                new Date(b.created_at).getTime() -
                new Date(a.created_at).getTime(),
            )
            .slice(0, 20)}
          {#each allHistory as item}
            <div
              class="p-4 sm:p-5 flex items-center gap-4 hover:bg-zinc-800/50 transition-colors"
            >
              <div
                class="w-10 h-10 rounded-full flex items-center justify-center border border-zinc-700/50 shadow-sm {item.type ===
                'game'
                  ? 'bg-amber-900/30 text-amber-400'
                  : 'bg-emerald-900/30 text-emerald-400'} flex-shrink-0 drop-shadow-sm"
              >
                {item.type === "game" ? "🎮" : "📚"}
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium text-white m-0">
                  {item.type === "game"
                    ? `Game Mengetik: ${item.mode.toUpperCase()}`
                    : `Pelajaran ${item.lessonId}`}
                </p>
                <p class="text-xs text-blue-200 mt-0.5 font-light m-0">
                  {new Date(item.created_at).toLocaleDateString("id-ID", {
                    day: "numeric",
                    month: "short",
                    hour: "2-digit",
                    minute: "2-digit",
                  })}
                </p>
              </div>
              <div class="text-right">
                {#if item.type === "game"}
                  <span class="text-sm font-bold text-amber-400 drop-shadow-sm"
                    >Skor: {item.score}</span
                  >
                {:else}
                  <div class="flex flex-col items-end gap-1">
                    <span
                      class="text-xs font-medium px-2 py-0.5 rounded-md bg-emerald-900/40 text-emerald-300 border border-emerald-800/50"
                      >{item.wpm} WPM</span
                    >
                    <span class="text-[0.65rem] font-light text-zinc-400 uppercase tracking-wider"
                      >Akurasi {item.accuracy}%</span
                    >
                  </div>
                {/if}
              </div>
            </div>
          {/each}
        {/if}
      </div>
    </div>
  </div>
</div>
