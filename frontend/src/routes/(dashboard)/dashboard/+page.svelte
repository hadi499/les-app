<script lang="ts">
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
      class="text-2xl font-extrabold text-orange-950 sm:text-3xl tracking-tight drop-shadow-sm"
    >
      Selamat datang kembali! 👋
    </h1>
  </div>

  <!-- Main Content -->
  <div class="max-w-4xl mx-auto space-y-6">
    <!-- Typing Progress Section -->
    <div
      class="bg-white/60 backdrop-blur-md rounded-3xl border border-orange-200 shadow-xl shadow-orange-900/10 overflow-hidden"
    >
      <div
        class="px-6 py-5 border-b border-orange-200 flex justify-between items-center bg-white/40"
      >
        <h3 class="text-lg font-bold text-orange-950 drop-shadow-sm m-0">
          Progress Mengetik 10 Jari
        </h3>
        <a
          href="/mengetik"
          class="text-sm font-medium text-orange-800 hover:text-orange-950 bg-white/50 border border-orange-300 hover:bg-white/80 px-3 py-1.5 rounded-lg transition-colors no-underline"
          >Lanjutkan Latihan</a
        >
      </div>

      <div class="p-6">
        <h4 class="text-sm font-medium text-orange-800 uppercase tracking-wider mb-4 m-0">Statistik Pelajaran</h4>
        <div class="grid grid-cols-2 gap-4 mb-8">
          <div
            class="bg-amber-100/50 rounded-2xl p-4 flex flex-col items-center justify-center border border-amber-300 hover:-translate-y-1 hover:shadow-lg hover:shadow-amber-900/10 transition-all"
          >
            <span class="text-3xl mb-2 drop-shadow-sm">⭐</span>
            <span class="text-2xl font-bold text-amber-600 drop-shadow-sm">{totalStars}</span>
            <span
              class="text-xs font-medium text-amber-800 uppercase tracking-wide mt-1"
              >Total Bintang</span
            >
          </div>
          <div
            class="bg-emerald-100/50 rounded-2xl p-4 flex flex-col items-center justify-center border border-emerald-300 hover:-translate-y-1 hover:shadow-lg hover:shadow-emerald-900/10 transition-all"
          >
            <span class="text-3xl mb-2 drop-shadow-sm">📚</span>
            <span class="text-2xl font-bold text-emerald-700 drop-shadow-sm"
              >{completedLessons}</span
            >
            <span
              class="text-xs font-medium text-emerald-800 uppercase tracking-wide text-center mt-1"
              >Pelajaran Selesai</span
            >
          </div>
        </div>

        <div class="mt-4 pt-4 border-t border-orange-200 flex justify-end">
          <a
            href="/dashboard/card-memory"
            class="text-sm font-medium text-orange-900 hover:text-orange-950 bg-white border border-orange-300 hover:bg-orange-50 px-4 py-2 rounded-lg transition-colors no-underline flex items-center gap-2 shadow-sm"
          >
            <span>Buka Card Memory</span>
            <span class="text-lg drop-shadow-sm">🎴</span>
          </a>
        </div>

        <h4 class="text-sm font-medium text-orange-800 uppercase tracking-wider mb-4 mt-8 m-0">High Score Game Mengetik</h4>
        <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-4">
          <div
            class="bg-rose-100/50 rounded-2xl p-4 flex flex-col items-center justify-center border border-rose-300 hover:-translate-y-1 hover:shadow-lg hover:shadow-rose-900/10 transition-all"
          >
            <span class="text-3xl mb-2 drop-shadow-sm">👈</span>
            <span class="text-2xl font-bold text-rose-600 drop-shadow-sm"
              >{$gameScores.left || 0}</span
            >
            <span
              class="text-[0.65rem] font-medium text-rose-800 uppercase tracking-wide text-center mt-1"
              >Tangan Kiri</span
            >
          </div>
          <div
            class="bg-teal-100/50 rounded-2xl p-4 flex flex-col items-center justify-center border border-teal-300 hover:-translate-y-1 hover:shadow-lg hover:shadow-teal-900/10 transition-all"
          >
            <span class="text-3xl mb-2 drop-shadow-sm">👉</span>
            <span class="text-2xl font-bold text-teal-600 drop-shadow-sm"
              >{$gameScores.right || 0}</span
            >
            <span
              class="text-[0.65rem] font-medium text-teal-800 uppercase tracking-wide text-center mt-1"
              >Tangan Kanan</span
            >
          </div>
          <div
            class="bg-blue-100/50 rounded-2xl p-4 flex flex-col items-center justify-center border border-blue-300 hover:-translate-y-1 hover:shadow-lg hover:shadow-blue-900/10 transition-all"
          >
            <span class="text-3xl mb-2 drop-shadow-sm">🎯</span>
            <span class="text-2xl font-bold text-blue-600 drop-shadow-sm"
              >{$gameScores.both || 0}</span
            >
            <span
              class="text-[0.65rem] font-medium text-blue-800 uppercase tracking-wide text-center mt-1"
              >Home Row</span
            >
          </div>
          <div
            class="bg-orange-100/50 rounded-2xl p-4 flex flex-col items-center justify-center border border-orange-300 hover:-translate-y-1 hover:shadow-lg hover:shadow-orange-900/10 transition-all"
          >
            <span class="text-3xl mb-2 drop-shadow-sm">🔤</span>
            <span class="text-2xl font-bold text-orange-600 drop-shadow-sm"
              >{$gameScores.letters || 0}</span
            >
            <span
              class="text-[0.65rem] font-medium text-orange-800 uppercase tracking-wide text-center mt-1"
              >Semua Huruf</span
            >
          </div>
          <div
            class="bg-purple-100/50 rounded-2xl p-4 flex flex-col items-center justify-center border border-purple-300 hover:-translate-y-1 hover:shadow-lg hover:shadow-purple-900/10 transition-all"
          >
            <span class="text-3xl mb-2 drop-shadow-sm">🚀</span>
            <span class="text-2xl font-bold text-purple-600 drop-shadow-sm"
              >{$gameScores.all || 0}</span
            >
            <span
              class="text-[0.65rem] font-medium text-purple-800 uppercase tracking-wide text-center mt-1"
              >Semua Tombol</span
            >
          </div>
        </div>
      </div>
    </div>

    <!-- History Section -->
    <div
      class="bg-white/60 backdrop-blur-md rounded-3xl border border-orange-200 shadow-xl shadow-orange-900/10 overflow-hidden"
    >
      <div class="px-6 py-5 border-b border-orange-200 bg-white/40">
        <h3 class="text-lg font-bold text-orange-950 drop-shadow-sm m-0">
          Riwayat Aktivitas Mengetik
        </h3>
      </div>

      <div class="divide-y divide-orange-200 max-h-[400px] overflow-y-auto">
        {#if $gameHistoryStore.length === 0 && $lessonHistoryStore.length === 0}
          <div class="p-8 text-center text-orange-700 font-medium">
            Belum ada riwayat aktivitas. Ayo mulai mengetik!
          </div>
        {:else}
          {@const allHistory = [
            ...$gameHistoryStore.map((h) => ({ type: "game" as const, ...h })),
            ...$lessonHistoryStore.map((h) => ({ type: "lesson" as const, ...h })),
          ]
            .sort(
              (a, b) =>
                new Date(b.created_at).getTime() -
                new Date(a.created_at).getTime(),
            )
            .slice(0, 20)}
          {#each allHistory as item}
            <div
              class="p-4 sm:p-5 flex items-center gap-4 hover:bg-white/50 transition-colors"
            >
              <div
                class="w-10 h-10 rounded-full flex items-center justify-center border border-orange-300 shadow-sm {item.type ===
                'game'
                  ? 'bg-amber-100 text-amber-600'
                  : 'bg-emerald-100 text-emerald-700'} flex-shrink-0 drop-shadow-sm"
              >
                {item.type === "game" ? "🎮" : "📚"}
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm font-bold text-orange-950 m-0">
                  {item.type === "game"
                    ? `Game Mengetik: ${item.mode.toUpperCase()}`
                    : `Pelajaran ${item.lessonId}`}
                </p>
                <p class="text-xs text-orange-800 mt-0.5 font-medium m-0">
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
                  <span class="text-sm font-bold text-amber-600 drop-shadow-sm"
                    >Skor: {item.score}</span
                  >
                {:else}
                  <div class="flex flex-col items-end gap-1">
                    <span
                      class="text-xs font-bold px-2 py-0.5 rounded-md bg-emerald-100 text-emerald-700 border border-emerald-300"
                      >{item.wpm} WPM</span
                    >
                    <span class="text-[0.65rem] font-medium text-orange-800 uppercase tracking-wider"
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
