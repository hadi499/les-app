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

  let systemInfo: any = $state(null);
  let isRefreshing = $state(false);
  
  let isClassOpenPaud = $state(true);
  let isClassOpenSd = $state(true);
  let isLoadingTogglePaud = $state(false);
  let isLoadingToggleSd = $state(false);
  let userRole = $state("");
  let isSettingsLoaded = $state(false);
  let isRoleLoaded = $state(false);

  async function fetchSettings() {
    try {
      const res = await fetch("/api/settings", { credentials: "include" });
      if (res.ok) {
        const json = await res.json();
        isClassOpenPaud = json.is_class_open_paud === "true";
        isClassOpenSd = json.is_class_open_sd === "true";
      }
    } catch (e) {
      console.error("Gagal memuat pengaturan", e);
    } finally {
      isSettingsLoaded = true;
    }
  }

  async function toggleClassOpenPaud() {
    isLoadingTogglePaud = true;
    const newValue = !isClassOpenPaud;
    try {
      const res = await fetch("/api/settings/is_class_open_paud", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ value: newValue ? "true" : "false" })
      });
      if (res.ok) {
        isClassOpenPaud = newValue;
      } else {
        alert("Gagal memperbarui status kelas PAUD/TK.");
      }
    } catch (e) {
      console.error(e);
      alert("Terjadi kesalahan koneksi.");
    } finally {
      isLoadingTogglePaud = false;
    }
  }

  async function toggleClassOpenSd() {
    isLoadingToggleSd = true;
    const newValue = !isClassOpenSd;
    try {
      const res = await fetch("/api/settings/is_class_open_sd", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ value: newValue ? "true" : "false" })
      });
      if (res.ok) {
        isClassOpenSd = newValue;
      } else {
        alert("Gagal memperbarui status kelas SD/SMP.");
      }
    } catch (e) {
      console.error(e);
      alert("Terjadi kesalahan koneksi.");
    } finally {
      isLoadingToggleSd = false;
    }
  }

  async function fetchSystemInfo() {
    isRefreshing = true;
    try {
      const meRes = await fetch("/me", { credentials: "include" });
      if (meRes.ok) {
        const meData = await meRes.json();
        if (meData.authenticated) {
          userRole = meData.user.role;
          if (userRole === "teacher") {
            const res = await fetch("/api/system/info", { credentials: "include" });
            if (res.ok) {
              const json = await res.json();
              systemInfo = json.data;
            }
          }
        }
      }
    } catch (e) {
      console.error(e);
    } finally {
      isRefreshing = false;
      isRoleLoaded = true;
    }
  }

  onMount(async () => {
    await fetchAllProgress();
    await fetchAllGameScores();
    await fetchHistory();
    await fetchSystemInfo();
    await fetchSettings();
  });

  let completedLessons = $derived(
    Array.from($progressMap.values()).filter((p) => p.completed).length,
  );
  let totalStars = $derived(
    Array.from($progressMap.values()).reduce((sum, p) => sum + p.stars, 0),
  );

  function formatBytes(bytes: number) {
    if (!bytes) return "0 B";
    const k = 1024;
    const sizes = ["B", "KB", "MB", "GB", "TB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
  }
</script>

<svelte:head>
  <title>Dashboard - Portal Les Balongarut</title>
</svelte:head>

<div class="animate-in fade-in duration-500">
  <!-- Header -->
  <div class="mb-8">
    <h1
      class="text-2xl font-extrabold text-slate-900 sm:text-3xl tracking-tight drop-shadow-sm"
    >
      Selamat datang kembali! 👋
    </h1>
  </div>

  <!-- Main Content -->
  <div class="max-w-4xl mx-auto space-y-6">
    <!-- Control Panel -->
    <div class="bg-white/60 backdrop-blur-md rounded-3xl border border-slate-200 shadow-xl shadow-slate-800/10 overflow-hidden">
      <div class="p-6 flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h4 class="font-bold text-slate-800 m-0">Status Kelas Hari Ini</h4>
        </div>
        
        <div class="flex flex-col sm:flex-row items-center gap-4">
          <!-- PAUD/TK -->
          <div class="flex items-center gap-3 bg-slate-50 px-4 py-2 rounded-2xl border border-slate-100">
            <span class="text-sm font-semibold text-slate-700">PAUD/TK:</span>
            {#if !isSettingsLoaded || !isRoleLoaded}
              <div class="w-24 h-8 bg-slate-200/70 rounded-full animate-pulse"></div>
            {:else if userRole === 'teacher'}
              <span class="text-sm font-bold {isClassOpenPaud ? 'text-emerald-600' : 'text-rose-600'}">{isClassOpenPaud ? 'BUKA' : 'LIBUR'}</span>
              <button
                onclick={toggleClassOpenPaud}
                disabled={isLoadingTogglePaud}
                class="relative inline-flex h-8 w-14 items-center rounded-full transition-colors focus:outline-none disabled:opacity-50 {isClassOpenPaud ? 'bg-emerald-500' : 'bg-rose-500'} border-none cursor-pointer"
              >
                <span class="sr-only">Toggle Kelas PAUD</span>
                <span
                  class="inline-block h-6 w-6 transform rounded-full bg-white transition-transform shadow-sm {isClassOpenPaud ? 'translate-x-7' : 'translate-x-1'}"
                ></span>
              </button>
            {:else}
              <div class="px-3 py-1 rounded-full font-bold text-xs flex items-center gap-1.5 shadow-sm {isClassOpenPaud ? 'bg-emerald-100 text-emerald-700 border border-emerald-200' : 'bg-rose-100 text-rose-700 border border-rose-200'}">
                {#if isClassOpenPaud}
                  <span class="relative flex h-2 w-2">
                    <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                    <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
                  </span>
                  BUKA
                {:else}
                  <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  LIBUR
                {/if}
              </div>
            {/if}
          </div>

          <!-- SD/SMP -->
          <div class="flex items-center gap-3 bg-slate-50 px-4 py-2 rounded-2xl border border-slate-100">
            <span class="text-sm font-semibold text-slate-700">SD/SMP:</span>
            {#if !isSettingsLoaded || !isRoleLoaded}
              <div class="w-24 h-8 bg-slate-200/70 rounded-full animate-pulse"></div>
            {:else if userRole === 'teacher'}
              <span class="text-sm font-bold {isClassOpenSd ? 'text-emerald-600' : 'text-rose-600'}">{isClassOpenSd ? 'BUKA' : 'LIBUR'}</span>
              <button
                onclick={toggleClassOpenSd}
                disabled={isLoadingToggleSd}
                class="relative inline-flex h-8 w-14 items-center rounded-full transition-colors focus:outline-none disabled:opacity-50 {isClassOpenSd ? 'bg-emerald-500' : 'bg-rose-500'} border-none cursor-pointer"
              >
                <span class="sr-only">Toggle Kelas SD</span>
                <span
                  class="inline-block h-6 w-6 transform rounded-full bg-white transition-transform shadow-sm {isClassOpenSd ? 'translate-x-7' : 'translate-x-1'}"
                ></span>
              </button>
            {:else}
              <div class="px-3 py-1 rounded-full font-bold text-xs flex items-center gap-1.5 shadow-sm {isClassOpenSd ? 'bg-emerald-100 text-emerald-700 border border-emerald-200' : 'bg-rose-100 text-rose-700 border border-rose-200'}">
                {#if isClassOpenSd}
                  <span class="relative flex h-2 w-2">
                    <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                    <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
                  </span>
                  BUKA
                {:else}
                  <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  LIBUR
                {/if}
              </div>
            {/if}
          </div>
        </div>
      </div>
    </div>

    {#if systemInfo}
      <!-- System Info Section (VPS) -->
      <div
        class="bg-white/60 backdrop-blur-md rounded-3xl border border-slate-200 shadow-xl shadow-slate-800/10 overflow-hidden"
      >
        <div class="px-6 py-5 border-b border-slate-200 bg-white/40 flex justify-between items-center">
          <h3 class="text-lg font-bold text-slate-900 drop-shadow-sm m-0 flex items-center gap-2">
            <span>🖥️</span> Status Server (VPS)
          </h3>
          <button 
            onclick={fetchSystemInfo}
            disabled={isRefreshing}
            class="text-sm font-medium text-slate-600 hover:text-slate-900 bg-white/50 border border-slate-300 hover:bg-white/80 px-3 py-1.5 rounded-lg transition-colors flex items-center gap-2 shadow-sm disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span class={isRefreshing ? "animate-spin inline-block" : "inline-block"}>🔄</span>
            {isRefreshing ? "Memperbarui..." : "Perbarui"}
          </button>
        </div>
        <div class="p-6">
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <!-- Memory -->
            <div class="bg-indigo-50/50 rounded-2xl p-5 border border-indigo-200">
              <div class="flex items-center gap-3 mb-3">
                <div class="w-8 h-8 rounded-full bg-indigo-100 flex items-center justify-center text-indigo-600">🧠</div>
                <h4 class="font-bold text-slate-800 m-0">Memori (RAM)</h4>
              </div>
              <div class="space-y-2">
                <div class="flex justify-between text-sm">
                  <span class="text-slate-600">Terpakai</span>
                  <span class="font-bold text-slate-900">{formatBytes(systemInfo.memory.used)} / {formatBytes(systemInfo.memory.total)}</span>
                </div>
                <div class="w-full bg-slate-200 rounded-full h-2.5 overflow-hidden">
                  <div class="bg-indigo-500 h-2.5 rounded-full" style="width: {(systemInfo.memory.used / systemInfo.memory.total) * 100}%"></div>
                </div>
              </div>
            </div>

            <!-- Storage -->
            <div class="bg-fuchsia-50/50 rounded-2xl p-5 border border-fuchsia-200">
              <div class="flex items-center gap-3 mb-3">
                <div class="w-8 h-8 rounded-full bg-fuchsia-100 flex items-center justify-center text-fuchsia-600">💾</div>
                <h4 class="font-bold text-slate-800 m-0">Penyimpanan (Disk)</h4>
              </div>
              <div class="space-y-2">
                <div class="flex justify-between text-sm">
                  <span class="text-slate-600">Terpakai</span>
                  <span class="font-bold text-slate-900">{formatBytes(systemInfo.storage.used)} / {formatBytes(systemInfo.storage.total)}</span>
                </div>
                <div class="w-full bg-slate-200 rounded-full h-2.5 overflow-hidden">
                  <div class="bg-fuchsia-500 h-2.5 rounded-full" style="width: {(systemInfo.storage.used / systemInfo.storage.total) * 100}%"></div>
                </div>
              </div>
            </div>

            <!-- Uploads Folder -->
            <div class="bg-amber-50/50 rounded-2xl p-5 border border-amber-200">
              <div class="flex items-center gap-3 mb-3">
                <div class="w-8 h-8 rounded-full bg-amber-100 flex items-center justify-center text-amber-600">📁</div>
                <h4 class="font-bold text-slate-800 m-0">Folder Uploads</h4>
              </div>
              <div class="space-y-2">
                <div class="flex justify-between text-sm mt-6">
                  <span class="text-slate-600">Total Ukuran</span>
                  {#if systemInfo.uploads && systemInfo.uploads.error}
                    <span class="font-bold text-red-600 cursor-help" title={systemInfo.uploads.error}>Error (Arahkan Mouse)</span>
                  {:else}
                    <span class="font-bold text-slate-900">{systemInfo.uploads ? formatBytes(systemInfo.uploads.size) : '0 B'}</span>
                  {/if}
                </div>
                <div class="w-full bg-slate-200 rounded-full h-2.5 overflow-hidden">
                  <div class="bg-amber-500 h-2.5 rounded-full w-full"></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    {/if}

    <!-- Typing Progress Section -->
    <div
      class="bg-white/60 backdrop-blur-md rounded-3xl border border-slate-200 shadow-xl shadow-slate-800/10 overflow-hidden"
    >
      <div
        class="px-6 py-5 border-b border-slate-200 flex justify-between items-center bg-white/40"
      >
        <h3 class="text-lg font-bold text-slate-900 drop-shadow-sm m-0">
          Progress Mengetik 10 Jari
        </h3>
        <a
          href="/mengetik"
          class="text-sm font-medium text-slate-600 hover:text-slate-900 bg-white/50 border border-slate-300 hover:bg-white/80 px-3 py-1.5 rounded-lg transition-colors no-underline"
          >Lanjutkan Latihan</a
        >
      </div>

      <div class="p-6">
        <h4 class="text-sm font-medium text-slate-600 uppercase tracking-wider mb-4 m-0">Statistik Pelajaran</h4>
        <div class="grid grid-cols-2 gap-4 mb-8">
          <div
            class="bg-sky-100/50 rounded-2xl p-4 flex flex-col items-center justify-center border border-sky-300 hover:-translate-y-1 hover:shadow-lg hover:shadow-sky-900/10 transition-all"
          >
            <span class="text-3xl mb-2 drop-shadow-sm">⭐</span>
            <span class="text-2xl font-bold text-sky-600 drop-shadow-sm">{totalStars}</span>
            <span
              class="text-xs font-medium text-sky-800 uppercase tracking-wide mt-1"
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

        <div class="mt-4 pt-4 border-t border-slate-200 flex justify-end">
          <a
            href="/dashboard/card-memory"
            class="text-sm font-medium text-slate-800 hover:text-slate-900 bg-white border border-slate-300 hover:bg-slate-50 px-4 py-2 rounded-lg transition-colors no-underline flex items-center gap-2 shadow-sm"
          >
            <span>Buka Card Memory</span>
            <span class="text-lg drop-shadow-sm">🎴</span>
          </a>
        </div>

        <h4 class="text-sm font-medium text-slate-600 uppercase tracking-wider mb-4 mt-8 m-0">High Score Game Mengetik</h4>
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
            class="bg-slate-100/50 rounded-2xl p-4 flex flex-col items-center justify-center border border-slate-300 hover:-translate-y-1 hover:shadow-lg hover:shadow-slate-800/10 transition-all"
          >
            <span class="text-3xl mb-2 drop-shadow-sm">🔤</span>
            <span class="text-2xl font-bold text-blue-600 drop-shadow-sm"
              >{$gameScores.letters || 0}</span
            >
            <span
              class="text-[0.65rem] font-medium text-slate-600 uppercase tracking-wide text-center mt-1"
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
      class="bg-white/60 backdrop-blur-md rounded-3xl border border-slate-200 shadow-xl shadow-slate-800/10 overflow-hidden"
    >
      <div class="px-6 py-5 border-b border-slate-200 bg-white/40">
        <h3 class="text-lg font-bold text-slate-900 drop-shadow-sm m-0">
          Riwayat Aktivitas Mengetik
        </h3>
      </div>

      <div class="divide-y divide-slate-200 max-h-[400px] overflow-y-auto">
        {#if $gameHistoryStore.length === 0 && $lessonHistoryStore.length === 0}
          <div class="p-8 text-center text-slate-500 font-medium">
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
                class="w-10 h-10 rounded-full flex items-center justify-center border border-slate-300 shadow-sm {item.type ===
                'game'
                  ? 'bg-sky-100 text-sky-600'
                  : 'bg-emerald-100 text-emerald-700'} flex-shrink-0 drop-shadow-sm"
              >
                {item.type === "game" ? "🎮" : "📚"}
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm font-bold text-slate-900 m-0">
                  {item.type === "game"
                    ? `Game Mengetik: ${item.mode.toUpperCase()}`
                    : `Pelajaran ${item.lessonId}`}
                </p>
                <p class="text-xs text-slate-600 mt-0.5 font-medium m-0">
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
                  <span class="text-sm font-bold text-sky-600 drop-shadow-sm"
                    >Skor: {item.score}</span
                  >
                {:else}
                  <div class="flex flex-col items-end gap-1">
                    <span
                      class="text-xs font-bold px-2 py-0.5 rounded-md bg-emerald-100 text-emerald-700 border border-emerald-300"
                      >{item.wpm} WPM</span
                    >
                    <span class="text-[0.65rem] font-medium text-slate-600 uppercase tracking-wider"
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
