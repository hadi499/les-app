<script lang="ts">
  import { onMount } from "svelte";
  import { toast } from "$lib/stores/toast.svelte";

  let user = $state<{ id?: number; username?: string; role?: string } | null>(null);

  // State for absence form
  let selectedUserId = $state("");

  function getLocalDate() {
    const d = new Date();
    d.setMinutes(d.getMinutes() - d.getTimezoneOffset());
    return d.toISOString().split("T")[0];
  }
  let selectedDate = $state(getLocalDate());

  let selectedReasonType = $state("Sakit");
  let reasonNote = $state("");
  let isSubmitting = $state(false);

  // Real data
  type User = { id: number; username: string; role: string };
  let users: User[] = $state([]);

  type RecapData = {
    user_id: number;
    name: string;
    total: number;
    sakit: number;
    izin: number;
    alpa: number;
    total_all: number;
  };
  let recapData: RecapData[] = $state([]);
  let isLoading = $state(true);
  let errorMsg = $state("");

  async function fetchUsers() {
    try {
      const res = await fetch(`/api/users`, { credentials: "include" });
      if (!res.ok) throw new Error("Gagal mengambil data siswa");
      const data = await res.json();
      // Filter only students
      users = (data.users || []).filter((u: User) => u.role === "student");
    } catch (e) {
      errorMsg = e instanceof Error ? e.message : String(e);
    }
  }

  async function fetchRecap() {
    try {
      const res = await fetch(`/api/absences/recap`, {
        credentials: "include",
      });
      if (!res.ok) throw new Error("Gagal mengambil rekap absensi");
      const data = await res.json();
      recapData = data.recap || [];
    } catch (e) {
      errorMsg = e instanceof Error ? e.message : String(e);
    }
  }

  onMount(async () => {
    isLoading = true;
    try {
      const res = await fetch("/me", { credentials: "include" });
      if (res.ok) {
        const data = await res.json();
        if (data.authenticated) {
          user = data.user;
        }
      }
    } catch (e) {
      console.error(e);
    }

    if (user?.role === "teacher") {
      await Promise.all([fetchUsers(), fetchRecap()]);
    } else {
      await fetchRecap();
    }
    isLoading = false;
  });

  async function handleSubmit() {
    if (!selectedUserId) {
      toast.error("Pilih siswa terlebih dahulu.");
      return;
    }
    isSubmitting = true;
    try {
      const res = await fetch(`/api/absences`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({
          user_id: Number(selectedUserId),
          date: selectedDate,
          reason_type: selectedReasonType,
          note: reasonNote,
        }),
      });

      if (!res.ok) {
        const errData = await res.json();
        throw new Error(errData.error || "Gagal menyimpan absensi");
      }

      toast.success("Ketidakhadiran berhasil dicatat!");
      // Reset form
      selectedUserId = "";
      reasonNote = "";

      // Refresh recap
      await fetchRecap();
    } catch (e) {
      toast.error(
        "Terjadi kesalahan: " + (e instanceof Error ? e.message : String(e)),
      );
    } finally {
      isSubmitting = false;
    }
  }

</script>

<svelte:head>
  <title>Laporan Ketidakhadiran - Portal Guru</title>
</svelte:head>

<div class="animate-in fade-in duration-500 pb-12">
  <div
    class="mb-8 flex flex-col sm:flex-row sm:items-center justify-between gap-4"
  >
    <div>
      <h1
        class="text-2xl font-bold text-slate-900 sm:text-3xl tracking-tight drop-shadow-sm"
      >
        Laporan Ketidakhadiran
      </h1>
      <p
        class="mt-2 text-slate-600 text-sm sm:text-base font-light tracking-wide"
      >
        Catat dan pantau riwayat ketidakhadiran siswa.
      </p>
    </div>

  </div>

  {#if isLoading}
    <div class="flex justify-center p-12">
      <div
        class="w-10 h-10 border-4 border-slate-200 border-t-blue-500 rounded-full animate-spin"
      ></div>
    </div>
  {:else if errorMsg}
    <div
      class="bg-red-100 text-red-800 p-6 rounded-2xl border border-red-300 font-medium mb-8"
    >
      {errorMsg}
    </div>
  {:else}
    {#if user?.role === "teacher"}
      <!-- Form Input Ketidakhadiran -->
      <div
        class="max-w-4xl bg-white/60 backdrop-blur-md rounded-3xl border border-slate-200 shadow-lg shadow-slate-800/10 p-6 mb-8"
      >
      <h2 class="text-lg font-bold text-slate-900 mb-4 flex items-center gap-2">
        <svg
          class="w-5 h-5 text-blue-500"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          ><path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
          ></path></svg
        >
        Catat Ketidakhadiran Baru
      </h2>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-4">
        <!-- Select User -->
        <div class="flex flex-col gap-1.5">
          <label for="userSelect" class="text-sm font-semibold text-slate-700"
            >Nama Siswa</label
          >
          <select
            id="userSelect"
            bind:value={selectedUserId}
            class="px-3 py-2.5 text-sm bg-white border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500/50 transition-shadow"
          >
            <option value="" disabled>Pilih Siswa...</option>
            {#each users as user}
              <option value={user.id}>{user.username}</option>
            {/each}
          </select>
        </div>

        <!-- Select Date -->
        <div class="flex flex-col gap-1.5">
          <label for="dateSelect" class="text-sm font-semibold text-slate-700"
            >Tanggal</label
          >
          <input
            id="dateSelect"
            type="date"
            bind:value={selectedDate}
            class="px-3 py-2.5 text-sm bg-white border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500/50 transition-shadow"
          />
        </div>

        <!-- Tipe Alasan -->
        <div class="flex flex-col gap-1.5">
          <label for="reasonSelect" class="text-sm font-semibold text-slate-700"
            >Alasan</label
          >
          <select
            id="reasonSelect"
            bind:value={selectedReasonType}
            class="px-3 py-2.5 text-sm bg-white border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500/50 transition-shadow"
          >
            <option value="Sakit">Sakit</option>
            <option value="Izin">Izin</option>
            <option value="Alpa">Alpa (Tanpa Keterangan)</option>
          </select>
        </div>

        <!-- Keterangan -->
        <div class="flex flex-col gap-1.5 lg:col-span-1 md:col-span-2">
          <label for="noteInput" class="text-sm font-semibold text-slate-700"
            >Keterangan Detail</label
          >
          <input
            id="noteInput"
            type="text"
            bind:value={reasonNote}
            placeholder="Contoh: Sakit demam, acara keluarga..."
            class="px-3 py-2.5 text-sm bg-white border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500/50 transition-shadow placeholder:text-slate-400"
          />
        </div>
      </div>

      <div class="flex justify-end">
        <button
          onclick={handleSubmit}
          disabled={isSubmitting || !selectedUserId}
          class="px-5 py-2 text-sm font-semibold text-white bg-blue-600 shadow-lg shadow-blue-600/30 hover:bg-blue-700 rounded-xl transition-all flex items-center gap-2 disabled:opacity-70 disabled:cursor-not-allowed cursor-pointer"
        >
          {#if isSubmitting}
            <div
              class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"
            ></div>
            Menyimpan...
          {:else}
            <svg
              class="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              ><path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M12 4v16m8-8H4"
              ></path></svg
            >
            Tambahkan Data
          {/if}
        </button>
      </div>
    </div>
    {/if}

    <!-- Rekap Heading -->
    <div class="mb-6 mt-8 flex items-center justify-between">
      <h2 class="text-xl font-bold text-slate-900 flex items-center gap-2 drop-shadow-sm">
        <svg class="w-6 h-6 text-purple-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
        </svg>
        Rekap Ketidakhadiran
      </h2>
    </div>

    <!-- Grid Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      {#each recapData as data}
        <div class="bg-white/90 backdrop-blur-sm border border-slate-200 rounded-2xl p-5 shadow-sm hover:shadow-md transition-all flex flex-col h-full group">
          <div class="flex justify-between items-start mb-5">
            <h3 class="text-lg font-bold text-slate-900 group-hover:text-blue-700 transition-colors">{data.name}</h3>
          </div>

          <div class="grid grid-cols-2 gap-3 mb-5">
            <div class="{data.total > 0 ? 'bg-red-50 border-red-100 text-red-600' : 'bg-green-50 border-green-100 text-green-600'} border rounded-xl p-3 flex flex-col items-center justify-center">
              <span class="text-3xl font-black mb-0.5">{data.total}</span>
              <span class="text-[10px] font-bold {data.total > 0 ? 'text-red-400' : 'text-green-500'} uppercase tracking-wider">Bulan Ini</span>
            </div>
            <div class="{data.total_all > 0 ? 'bg-orange-50 border-orange-100 text-orange-600' : 'bg-green-50 border-green-100 text-green-600'} border rounded-xl p-3 flex flex-col items-center justify-center">
              <span class="text-3xl font-black mb-0.5">{data.total_all}</span>
              <span class="text-[10px] font-bold {data.total_all > 0 ? 'text-orange-400' : 'text-green-500'} uppercase tracking-wider">Total Hari</span>
            </div>
          </div>

          <div class="flex items-center justify-between mb-5 px-4 py-3 bg-slate-50/80 rounded-xl border border-slate-100">
            <div class="flex flex-col items-center w-1/3">
              <span class="text-[10px] text-slate-500 font-bold uppercase tracking-wider mb-1">Sakit</span>
              <span class="text-base font-bold text-slate-700">{data.sakit}</span>
            </div>
            <div class="w-px h-8 bg-slate-200"></div>
            <div class="flex flex-col items-center w-1/3">
              <span class="text-[10px] text-slate-500 font-bold uppercase tracking-wider mb-1">Izin</span>
              <span class="text-base font-bold text-slate-700">{data.izin}</span>
            </div>
            <div class="w-px h-8 bg-slate-200"></div>
            <div class="flex flex-col items-center w-1/3">
              <span class="text-[10px] text-slate-500 font-bold uppercase tracking-wider mb-1">Alpa</span>
              <span class="text-base font-bold text-slate-700">{data.alpa}</span>
            </div>
          </div>

          <div class="mt-auto pt-4 border-t border-slate-100">
            <a
              href={`/dashboard/absen/${data.user_id}`}
              class="w-full inline-flex items-center justify-center gap-1.5 px-4 py-2.5 text-sm font-semibold text-blue-700 bg-blue-50 hover:bg-blue-100 rounded-xl transition-colors border border-blue-200 cursor-pointer no-underline"
            >
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" /></svg>
              Detail Historis
            </a>
          </div>
        </div>
      {/each}
      
      {#if recapData.length === 0}
        <div class="col-span-full py-16 flex flex-col items-center justify-center bg-white/40 backdrop-blur-sm rounded-3xl border border-slate-200">
          <svg class="w-16 h-16 text-slate-300 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
          </svg>
          <span class="text-slate-500 font-medium text-center text-lg">Belum ada data siswa.</span>
        </div>
      {/if}
    </div>
  {/if}
</div>
