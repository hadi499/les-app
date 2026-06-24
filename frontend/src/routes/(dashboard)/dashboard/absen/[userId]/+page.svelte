<script lang="ts">
  import { page } from "$app/state";
  import { onMount } from "svelte";

  let userId = $state(page.params.userId);
  let user = $state<{ id: number; username: string } | null>(null);
  
  type AbsenceData = {
    id: number;
    date: string;
    reason_type: string;
    note: string;
  };
  
  let historyData: AbsenceData[] = $state([]);
  let isLoading = $state(true);
  let errorMsg = $state("");

  // Edit Modal State
  let showEditModal = $state(false);
  let editId = $state<number | null>(null);
  let editDate = $state("");
  let editReasonType = $state("Sakit");
  let editNote = $state("");
  let isSaving = $state(false);

  // Delete Modal State
  let showDeleteModal = $state(false);
  let deleteId = $state<number | null>(null);
  let isDeleting = $state(false);

  async function fetchHistory() {
    isLoading = true;
    errorMsg = "";
    try {
      const res = await fetch(`/api/absences/user/${userId}`, { credentials: "include" });
      if (!res.ok) throw new Error("Gagal mengambil detail historis");
      const data = await res.json();
      historyData = data.history || [];
      if (historyData.length > 0 && data.history[0].user) {
         user = data.history[0].user;
      }
    } catch (e) {
      errorMsg = e instanceof Error ? e.message : String(e);
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    fetchHistory();
  });

  // Edit actions
  function openEdit(item: AbsenceData) {
    editId = item.id;
    editDate = item.date.split("T")[0];
    editReasonType = item.reason_type;
    editNote = item.note || "";
    showEditModal = true;
  }

  function closeEdit() {
    showEditModal = false;
    editId = null;
  }

  async function handleSaveEdit() {
    if (!editId) return;
    isSaving = true;
    try {
      const res = await fetch(`/api/absences/${editId}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({
          date: editDate,
          reason_type: editReasonType,
          note: editNote
        })
      });

      if (!res.ok) {
        const errData = await res.json();
        throw new Error(errData.error || "Gagal mengupdate absensi");
      }

      closeEdit();
      await fetchHistory();
    } catch (e) {
      alert("Terjadi kesalahan: " + (e instanceof Error ? e.message : String(e)));
    } finally {
      isSaving = false;
    }
  }

  // Delete actions
  function promptDelete(id: number) {
    deleteId = id;
    showDeleteModal = true;
  }

  function cancelDelete() {
    showDeleteModal = false;
    deleteId = null;
  }

  async function confirmDelete() {
    if (!deleteId) return;
    isDeleting = true;
    try {
      const res = await fetch(`/api/absences/${deleteId}`, {
        method: "DELETE",
        credentials: "include",
      });

      if (!res.ok) {
        const data = await res.json();
        throw new Error(data.error || "Gagal menghapus data");
      }

      showDeleteModal = false;
      deleteId = null;
      await fetchHistory();
    } catch (e) {
      alert("Terjadi kesalahan: " + (e instanceof Error ? e.message : String(e)));
    } finally {
      isDeleting = false;
    }
  }
</script>

<svelte:head>
  <title>Detail Historis Absensi - Portal Guru</title>
</svelte:head>

<div class="animate-in fade-in duration-500 pb-12">
  <div class="mb-8 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <div>
      <h1 class="text-2xl font-bold text-slate-900 sm:text-3xl tracking-tight drop-shadow-sm flex items-center gap-3">
        <a href="/dashboard/absen" aria-label="Kembali ke daftar absensi" class="text-slate-400 hover:text-blue-600 transition-colors">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path></svg>
        </a>
        Detail Historis Siswa
      </h1>
      <p class="mt-2 text-slate-600 text-sm sm:text-base font-light tracking-wide pl-9">
        Riwayat ketidakhadiran {user?.username ? `untuk ${user.username}` : ''}.
      </p>
    </div>
  </div>

  {#if isLoading}
    <div class="flex justify-center p-12">
      <div class="w-10 h-10 border-4 border-slate-200 border-t-blue-500 rounded-full animate-spin"></div>
    </div>
  {:else if errorMsg}
    <div class="bg-red-100 text-red-800 p-6 rounded-2xl border border-red-300 font-medium mb-8">
      {errorMsg}
    </div>
  {:else}
    <div class="bg-white/60 backdrop-blur-md rounded-3xl border border-slate-200 shadow-lg shadow-slate-800/10 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-white/40 border-b border-slate-200">
              <th class="py-4 px-6 font-bold text-slate-900 text-sm w-16">No</th>
              <th class="py-4 px-6 font-bold text-slate-900 text-sm">Tanggal</th>
              <th class="py-4 px-6 font-bold text-slate-900 text-sm">Alasan</th>
              <th class="py-4 px-6 font-bold text-slate-900 text-sm">Catatan</th>
              <th class="py-4 px-6 font-bold text-slate-900 text-sm text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-200">
            {#each historyData as item, i}
              <tr class="hover:bg-white/40 transition-colors">
                <td class="py-4 px-6 text-sm text-slate-600 font-medium">{i + 1}</td>
                <td class="py-4 px-6 text-sm font-medium text-slate-900">
                  {new Date(item.date).toLocaleDateString("id-ID", { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })}
                </td>
                <td class="py-4 px-6">
                  <span class={`px-2.5 py-1 rounded-md text-xs font-bold border ${
                    item.reason_type === 'Sakit' ? 'bg-amber-100 text-amber-700 border-amber-200' :
                    item.reason_type === 'Izin' ? 'bg-blue-100 text-blue-700 border-blue-200' :
                    'bg-red-100 text-red-700 border-red-200'
                  }`}>
                    {item.reason_type}
                  </span>
                </td>
                <td class="py-4 px-6 text-sm text-slate-600">
                  {#if item.note}
                    <span class="italic text-slate-500">"{item.note}"</span>
                  {:else}
                    <span class="text-slate-400">-</span>
                  {/if}
                </td>
                <td class="py-4 px-6 text-right">
                  <div class="flex items-center justify-end gap-2">
                    <button 
                      onclick={() => openEdit(item)}
                      class="p-1.5 text-slate-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors border border-transparent hover:border-blue-200 cursor-pointer"
                      title="Edit">
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path></svg>
                    </button>
                    <button 
                      onclick={() => promptDelete(item.id)}
                      class="p-1.5 text-slate-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors border border-transparent hover:border-red-200 cursor-pointer"
                      title="Hapus">
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
                    </button>
                  </div>
                </td>
              </tr>
            {/each}
            {#if historyData.length === 0}
              <tr>
                <td colspan="5" class="py-8 text-center text-slate-500 font-light">Tidak ada riwayat ketidakhadiran.</td>
              </tr>
            {/if}
          </tbody>
        </table>
      </div>
    </div>
  {/if}

  <!-- Edit Modal -->
  {#if showEditModal}
    <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-sm">
      <div class="bg-slate-50 backdrop-blur-md rounded-3xl p-6 w-full max-w-md shadow-2xl border border-slate-200 transform transition-all">
        <h3 class="text-xl font-bold text-slate-900 mb-4">Edit Ketidakhadiran</h3>
        
        <div class="space-y-4 mb-6">
          <div class="flex flex-col gap-1.5">
            <label for="editDate" class="text-sm font-semibold text-slate-700">Tanggal</label>
            <input id="editDate" type="date" bind:value={editDate} class="px-3 py-2 text-sm bg-white border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500/50 transition-shadow" />
          </div>

          <div class="flex flex-col gap-1.5">
            <label for="editReason" class="text-sm font-semibold text-slate-700">Alasan</label>
            <select id="editReason" bind:value={editReasonType} class="px-3 py-2 text-sm bg-white border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500/50 transition-shadow">
              <option value="Sakit">Sakit</option>
              <option value="Izin">Izin</option>
              <option value="Alpa">Alpa (Tanpa Keterangan)</option>
            </select>
          </div>

          <div class="flex flex-col gap-1.5">
            <label for="editNote" class="text-sm font-semibold text-slate-700">Catatan</label>
            <input id="editNote" type="text" bind:value={editNote} class="px-3 py-2 text-sm bg-white border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500/50 transition-shadow" />
          </div>
        </div>

        <div class="flex justify-end gap-3">
          <button onclick={closeEdit} disabled={isSaving} class="px-4 py-2 text-sm font-medium text-slate-700 bg-white border border-slate-200 hover:bg-slate-100 rounded-xl transition-all disabled:opacity-50 cursor-pointer">
            Batal
          </button>
          <button onclick={handleSaveEdit} disabled={isSaving} class="px-4 py-2 text-sm font-semibold text-white bg-blue-600 hover:bg-blue-700 rounded-xl transition-all shadow-md flex items-center gap-2 disabled:opacity-50 cursor-pointer">
            {#if isSaving}
              <div class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
            {/if}
            Simpan
          </button>
        </div>
      </div>
    </div>
  {/if}

  <!-- Delete Modal -->
  {#if showDeleteModal}
    <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-sm">
      <div class="bg-slate-50 backdrop-blur-md rounded-3xl p-6 w-full max-w-sm shadow-2xl border border-slate-200 transform transition-all text-center">
        <div class="w-12 h-12 rounded-full bg-red-100 text-red-600 flex items-center justify-center mx-auto mb-4">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path></svg>
        </div>
        <h3 class="text-xl font-bold text-slate-900 mb-2">Hapus Data?</h3>
        <p class="text-slate-600 text-sm mb-6">Data absensi ini akan dihapus secara permanen dan tidak dapat dikembalikan.</p>
        
        <div class="flex justify-center gap-3">
          <button onclick={cancelDelete} disabled={isDeleting} class="px-4 py-2 text-sm font-medium text-slate-700 bg-white border border-slate-200 hover:bg-slate-100 rounded-xl transition-all disabled:opacity-50 cursor-pointer">
            Batal
          </button>
          <button onclick={confirmDelete} disabled={isDeleting} class="px-4 py-2 text-sm font-semibold text-white bg-red-600 hover:bg-red-700 rounded-xl transition-all shadow-md flex items-center gap-2 disabled:opacity-50 cursor-pointer">
            {#if isDeleting}
              <div class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
            {/if}
            Ya, Hapus
          </button>
        </div>
      </div>
    </div>
  {/if}
</div>
