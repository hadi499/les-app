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

  let isResetting = $state(false);
  let showResetModal = $state(false);
  async function handleReset() {
    isResetting = true;
    try {
      const res = await fetch(`/api/absences/reset`, {
        method: "POST",
        credentials: "include",
      });

      if (!res.ok) {
        throw new Error("Gagal mereset absensi");
      }

      toast.success("Seluruh data absensi berhasil direset!");
      showResetModal = false;
      await fetchRecap();
    } catch (e) {
      toast.error(
        "Terjadi kesalahan: " + (e instanceof Error ? e.message : String(e)),
      );
    } finally {
      isResetting = false;
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

    {#if user?.role === "teacher"}
      <button
        onclick={() => (showResetModal = true)}
      disabled={isResetting}
      class="px-4 py-2 text-sm font-bold text-white bg-red-600 shadow-md shadow-red-600/30 hover:bg-red-700 rounded-xl transition-all flex items-center gap-2 disabled:opacity-70 disabled:cursor-not-allowed w-fit"
    >
      {#if isResetting}
        <div
          class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"
        ></div>
        Mereset...
      {:else}
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="16"
          height="16"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="lucide lucide-trash-2"
          ><path d="M3 6h18" /><path
            d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"
          /><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" /><line
            x1="10"
            x2="10"
            y1="11"
            y2="17"
          /><line x1="14" x2="14" y1="11" y2="17" /></svg
        >
        Reset Data Absensi
      {/if}
      </button>
    {/if}
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
      class="bg-white/60 backdrop-blur-md rounded-3xl border border-slate-200 shadow-lg shadow-slate-800/10 p-6 mb-8"
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

    <!-- Table Recap -->
    <div
      class="bg-white/60 backdrop-blur-md rounded-3xl border border-slate-200 shadow-lg shadow-slate-800/10 overflow-hidden"
    >
      <div class="p-5 border-b border-slate-200 bg-white/40">
        <h2 class="text-lg font-bold text-slate-900 flex items-center gap-2">
          <svg
            class="w-5 h-5 text-purple-500"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            ><path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
            ></path></svg
          >
          Rekap Ketidakhadiran
        </h2>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse whitespace-nowrap min-w-[700px]">
          <thead>
            <tr class="bg-white/40 border-b border-slate-200">
              <th class="py-4 px-6 font-bold text-slate-900 text-sm"
                >Nama Siswa</th
              >
              <th class="py-4 px-6 font-bold text-slate-900 text-sm text-center"
                >Bulan Ini</th
              >
              <th class="py-4 px-6 font-bold text-slate-900 text-sm text-center"
                >Total</th
              >
              <th class="py-4 px-6 font-bold text-slate-900 text-sm text-center"
                >Rincian</th
              >
              <th class="py-4 px-6 font-bold text-slate-900 text-sm text-center"
                >Aksi</th
              >
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-200">
            {#each recapData as data}
              <tr class="hover:bg-white/40 transition-colors">
                <td class="py-4 px-6 text-sm font-medium text-slate-900">
                  {data.name}
                </td>
                <td class="py-4 px-6 text-center">
                  {#if data.total > 0}
                    <span
                      class="inline-flex items-center justify-center px-2.5 py-1 rounded-full text-xs font-bold bg-red-100 text-red-700 border border-red-200"
                    >
                      {data.total} Hari
                    </span>
                  {:else}
                    <span
                      class="inline-flex items-center justify-center px-2.5 py-1 rounded-full text-xs font-bold bg-green-100 text-green-700 border border-green-200"
                    >
                      0 Hari
                    </span>
                  {/if}
                </td>
                <td class="py-4 px-6 text-center">
                  {#if data.total_all > 0}
                    <span
                      class="inline-flex items-center justify-center px-2.5 py-1 rounded-full text-xs font-bold bg-orange-100 text-orange-700 border border-orange-200"
                    >
                      {data.total_all} Hari
                    </span>
                  {:else}
                    <span
                      class="inline-flex items-center justify-center px-2.5 py-1 rounded-full text-xs font-bold bg-green-100 text-green-700 border border-green-200"
                    >
                      0 Hari
                    </span>
                  {/if}
                </td>
                <td class="py-4 px-6 text-center">
                  <div class="flex items-center justify-center gap-2">
                    <div class="flex flex-col items-center">
                      <span
                        class="text-[10px] text-slate-500 font-semibold uppercase"
                        >Sakit</span
                      >
                      <span class="text-sm font-bold text-slate-700"
                        >{data.sakit}</span
                      >
                    </div>
                    <div class="w-px h-6 bg-slate-200"></div>
                    <div class="flex flex-col items-center">
                      <span
                        class="text-[10px] text-slate-500 font-semibold uppercase"
                        >Izin</span
                      >
                      <span class="text-sm font-bold text-slate-700"
                        >{data.izin}</span
                      >
                    </div>
                    <div class="w-px h-6 bg-slate-200"></div>
                    <div class="flex flex-col items-center">
                      <span
                        class="text-[10px] text-slate-500 font-semibold uppercase"
                        >Alpa</span
                      >
                      <span class="text-sm font-bold text-slate-700"
                        >{data.alpa}</span
                      >
                    </div>
                  </div>
                </td>
                <td class="py-4 px-6 text-center">
                  <a
                    href={`/dashboard/absen/${data.user_id}`}
                    class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-blue-600 bg-blue-50 hover:bg-blue-100 rounded-lg transition-colors border border-blue-200 cursor-pointer no-underline"
                  >
                    Detail Historis
                  </a>
                </td>
              </tr>
            {/each}
            {#if recapData.length === 0}
              <tr>
                <td
                  colspan="5"
                  class="py-8 text-center text-slate-500 font-light"
                  >Belum ada data siswa.</td
                >
              </tr>
            {/if}
          </tbody>
        </table>
      </div>
    </div>
  {/if}
</div>

{#if showResetModal}
  <div
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/50 backdrop-blur-sm animate-in fade-in duration-200"
  >
    <div
      class="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden animate-in zoom-in-95 duration-200"
    >
      <div class="p-6">
        <div
          class="w-12 h-12 rounded-full bg-red-100 flex items-center justify-center mb-4 text-red-600"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="lucide lucide-alert-triangle"
            ><path
              d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"
            /><path d="M12 9v4" /><path d="M12 17h.01" /></svg
          >
        </div>
        <h3 class="text-xl font-bold text-slate-900 mb-2">
          Reset Seluruh Data Absensi?
        </h3>
        <p class="text-slate-600 text-sm">
          PERINGATAN: Tindakan ini akan menghapus <strong
            >seluruh data riwayat absensi murid</strong
          > secara permanen. Data yang telah dihapus tidak dapat dikembalikan lagi.
          Apakah Anda yakin ingin melanjutkan?
        </p>
      </div>
      <div
        class="px-6 py-4 bg-slate-50 border-t border-slate-100 flex justify-end gap-3"
      >
        <button
          onclick={() => (showResetModal = false)}
          disabled={isResetting}
          class="px-4 py-2 text-sm font-semibold text-slate-700 hover:bg-slate-200 rounded-xl transition-colors disabled:opacity-70"
        >
          Batal
        </button>
        <button
          onclick={handleReset}
          disabled={isResetting}
          class="px-4 py-2 text-sm font-semibold text-white bg-red-600 hover:bg-red-700 rounded-xl transition-colors disabled:opacity-70 flex items-center gap-2"
        >
          {#if isResetting}
            <div
              class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"
            ></div>
            Mereset...
          {:else}
            Ya, Reset Data
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}
