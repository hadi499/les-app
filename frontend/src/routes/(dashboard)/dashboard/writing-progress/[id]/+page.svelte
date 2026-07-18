<script lang="ts">
  import { onMount } from "svelte";
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import { compressImageFile } from "$lib/utils";

  let progressId = $derived(page.params.id);

  type WritingProgress = {
    id: number;
    user_id: number;
    date: string;
    image: string;
    user?: { username: string };
  };
  type User = { id: number; username: string; role?: string };

  let progress = $state<WritingProgress | null>(null);
  let isLoading = $state(true);
  let showLoadingSpinner = $state(false);
  let errorMsg = $state("");
  let isTeacher = $state(false);
  let users = $state<User[]>([]);

  // Modals state
  let showEditModal = $state(false);
  let showDeleteModal = $state(false);
  
  // Form state
  let formUserId: number | string = $state("");
  let formDate = $state("");
  let formImage = $state("");
  let formFile = $state<File | null>(null);
  let localImagePreview = $state("");
  let isUploadingImage = $state(false);

  async function checkRole() {
    try {
      const res = await fetch(`/me`, { credentials: "include" });
      const data = await res.json();
      if (data.authenticated && data.user) {
        if (data.user.role === "teacher") {
          isTeacher = true;
          fetchUsers();
        }
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function fetchUsers() {
    try {
      const res = await fetch(`/api/users`, { credentials: "include" });
      if (res.ok) {
        const data = await res.json();
        users = data.users.filter((u: any) => u.role === "student") || [];
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function fetchDetail() {
    try {
      const res = await fetch(`/api/writing-progress/${progressId}`, {
        credentials: "include",
      });
      if (!res.ok) throw new Error("Gagal memuat detail perkembangan menulis");
      progress = await res.json();
    } catch (e) {
      errorMsg = e instanceof Error ? e.message : String(e);
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    checkRole();
    fetchDetail();
  });

  function formatDate(dateStr: string | undefined) {
    if (!dateStr) return "-";
    return new Date(dateStr).toLocaleDateString("id-ID", {
      weekday: "long",
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  }

  function openEditModal() {
    if (!progress) return;
    formUserId = progress.user_id;
    formDate = progress.date ? progress.date.split("T")[0] : "";
    formImage = progress.image || "";
    formFile = null;
    localImagePreview = "";
    showEditModal = true;
  }

  async function handleImageUpload(e: Event) {
    const target = e.target as HTMLInputElement;
    if (!target.files || target.files.length === 0) return;
    isUploadingImage = true;
    try {
      let file = target.files[0];
      file = await compressImageFile(file);
      if (file.size > 1 * 1024 * 1024) {
        alert("Ukuran file setelah kompresi masih lebih dari 1MB. Silakan pilih gambar lain.");
        return;
      }
      formFile = file;
      localImagePreview = URL.createObjectURL(file);
    } finally {
      isUploadingImage = false;
    }
  }

  async function saveProgress(e: Event) {
    e.preventDefault();
    if (!formUserId || !formDate || (!formImage && !formFile)) {
      alert("Harap lengkapi semua data, termasuk gambar!");
      return;
    }

    let finalImageUrl = formImage;
    if (formFile) {
      isUploadingImage = true;
      const formData = new FormData();
      formData.append("image", formFile);
      try {
        const res = await fetch("/api/upload?type=writing_progress", {
          method: "POST",
          body: formData,
          credentials: "include",
        });
        if (!res.ok) throw new Error("Gagal mengunggah gambar");
        const data = await res.json();
        finalImageUrl = data.url;
      } catch (err) {
        alert(err instanceof Error ? err.message : String(err));
        isUploadingImage = false;
        return;
      }
    }

    const payload = {
      user_id: Number(formUserId),
      date: formDate,
      image: finalImageUrl,
    };

    try {
      const res = await fetch(`/api/writing-progress/${progressId}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
        credentials: "include",
      });
      if (!res.ok) throw new Error("Gagal menyimpan data");
      showEditModal = false;
      await fetchDetail();
    } catch (err) {
      alert(err instanceof Error ? err.message : String(err));
    } finally {
      isUploadingImage = false;
    }
  }

  async function executeDelete() {
    try {
      const res = await fetch(`/api/writing-progress/${progressId}`, {
        method: "DELETE",
        credentials: "include",
      });
      if (!res.ok) throw new Error("Gagal menghapus data");
      goto("/dashboard/writing-progress");
    } catch (err) {
      alert(err instanceof Error ? err.message : String(err));
      showDeleteModal = false;
    }
  }
</script>

<svelte:head>
  <title>Detail Perkembangan Menulis - Portal PAUD/TK</title>
</svelte:head>

<div
  class="animate-in fade-in zoom-in-95 duration-500 max-w-4xl mx-auto py-8 sm:px-6"
>
  <div class="mb-8 flex flex-col sm:flex-row sm:items-end justify-between gap-4">
    <div>
      <a
        href="/dashboard/writing-progress"
        class="inline-flex items-center gap-2 text-sm font-medium text-slate-500 hover:text-slate-800 transition-colors mb-4"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"
          ><path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M10 19l-7-7m0 0l7-7m-7 7h18"
          ></path></svg
        >
        Kembali
      </a>
      <div>
        <h1
          class="text-2xl font-bold text-slate-900 sm:text-3xl tracking-tight drop-shadow-sm"
        >
          Detail Perkembangan Menulis
        </h1>
        <p class="mt-1 text-slate-500 text-sm tracking-wide">
          Informasi lengkap hasil karya tulisan murid PAUD/TK.
        </p>
      </div>
    </div>
    {#if isTeacher && progress}
      <div class="flex items-center gap-3">
        <button
          onclick={openEditModal}
          class="inline-flex items-center gap-2 px-4 py-2 text-sm font-bold text-blue-700 bg-blue-50 hover:bg-blue-100 border border-blue-200 rounded-xl transition-all shadow-sm cursor-pointer"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path></svg>
          Edit
        </button>
        <button
          onclick={() => showDeleteModal = true}
          class="inline-flex items-center gap-2 px-4 py-2 text-sm font-bold text-red-700 bg-red-50 hover:bg-red-100 border border-red-200 rounded-xl transition-all shadow-sm cursor-pointer"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
          Hapus
        </button>
      </div>
    {/if}
  </div>

    {#if isLoading}
    <div class="fixed inset-0 z-[100] flex flex-col items-center justify-center bg-slate-50/50 backdrop-blur-sm {showLoadingSpinner ? 'opacity-100' : 'opacity-0'} transition-opacity duration-300">
      <div class="w-12 h-12 border-4 border-slate-200 border-t-indigo-600 rounded-full animate-spin shadow-sm"></div>
    </div>
  {:else if errorMsg}
    <div
      class="bg-red-50 text-red-700 p-6 rounded-3xl border border-red-200 font-medium flex items-center justify-center gap-3 shadow-sm"
    >
      <svg
        class="w-6 h-6 shrink-0"
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
  {:else if progress}
    <div
      class="bg-white/60 backdrop-blur-xl rounded-xl border border-slate-200/60 shadow-xl shadow-slate-800/5 overflow-hidden"
    >
      <!-- Header Area -->
      <div
        class="relative bg-gradient-to-br from-indigo-500 via-purple-500 to-indigo-600 p-5 sm:p-12 text-white overflow-hidden"
      >
        <div
          class="absolute top-0 right-0 -mt-16 -mr-16 w-64 h-64 bg-white opacity-10 rounded-full blur-3xl"
        ></div>
        <div
          class="absolute bottom-0 left-0 -mb-16 -ml-16 w-48 h-48 bg-indigo-900 opacity-20 rounded-full blur-2xl"
        ></div>

        <div
          class="relative z-10 flex flex-col sm:flex-row sm:items-center justify-between gap-6"
        >
          <div>
            <div
              class="inline-flex items-center gap-2 px-3 py-1 bg-white/20 backdrop-blur-md rounded-xl text-sm font-medium mb-4 shadow-sm border border-white/20"
            >
              <svg
                class="w-4 h-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"
                ></path></svg
              >
              Hasil Tulisan
            </div>
            <h2
              class="text-2xl sm:text-3xl font-semibold tracking-tight drop-shadow-md mb-2"
            >
              {progress.user?.username || "Tidak diketahui"}
            </h2>
            <p class="text-indigo-100 font-medium flex items-center gap-2">
              <svg
                class="w-5 h-5 opacity-80"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                ></path></svg
              >
              {formatDate(progress.date)}
            </p>
          </div>
        </div>
      </div>

      <!-- Content Area -->
      <div class="p-5 sm:p-12 space-y-10">
        <!-- Image Evidence Area -->
        <div>
          <h3
            class="text-xl font-bold text-slate-800 mb-6 flex items-center gap-2"
          >
            <svg
              class="w-6 h-6 text-indigo-500"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              ><path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
              ></path></svg
            >
            Foto Tulisan Anak
          </h3>

          {#if progress.image}
            <div
              class="mb-6 bg-amber-50 border border-amber-200 text-amber-700 rounded-2xl p-4 flex items-start gap-3 shadow-sm"
            >
              <svg
                class="w-5 h-5 shrink-0 mt-0.5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                ></path></svg
              >
              <div>
                <strong class="block font-semibold mb-0.5 text-sm"
                  >Informasi Penyimpanan</strong
                >
                <p class="text-sm opacity-90 leading-relaxed">
                  File gambar bukti karya tulis ini akan dihapus secara otomatis oleh
                  sistem dalam waktu 30 hari.
                </p>
              </div>
            </div>
            <div>
              <div
                class="rounded-3xl overflow-hidden border-4 border-white shadow-xl bg-slate-100"
              >
                <img
                  src={progress.image}
                  alt="Tulisan {progress.user?.username}"
                  class="w-full h-auto object-cover"
                />
              </div>
              <div
                class="mt-4 flex flex-col sm:flex-row sm:items-center justify-between gap-4 px-2"
              >
                <span
                  class="text-slate-500 font-medium text-sm text-center sm:text-left"
                  >File Foto</span
                >
                <a
                  href={progress.image}
                  download
                  target="_blank"
                  class="px-5 py-2.5 bg-indigo-50 hover:bg-indigo-100 text-indigo-700 font-medium rounded-xl transition-colors shadow-sm border border-indigo-100 flex items-center justify-center w-full sm:w-auto gap-2"
                >
                  <svg
                    class="w-4 h-4"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    ><path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"
                    ></path></svg
                  >
                  Buka Ukuran Penuh
                </a>
              </div>
            </div>
          {:else}
            <div
              class="bg-slate-50 border-2 border-dashed border-slate-200 rounded-3xl p-12 text-center"
            >
              <div
                class="w-20 h-20 bg-white rounded-full flex items-center justify-center mx-auto mb-4 shadow-sm border border-slate-100"
              >
                <svg
                  class="w-10 h-10 text-slate-300"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  ><path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
                  ></path></svg
                >
              </div>
              <h4 class="text-lg font-bold text-slate-700 mb-1">
                Tidak Ada Foto
              </h4>
              <p class="text-slate-500 max-w-sm mx-auto">
                Progress ini belum memiliki lampiran foto tulisan.
              </p>
            </div>
          {/if}
        </div>
      </div>
    </div>
  {/if}
</div>

<!-- Edit Modal -->
{#if showEditModal}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
    <div
      class="absolute inset-0 bg-black/60 backdrop-blur-sm"
      onclick={() => showEditModal = false}
      role="button"
      tabindex="0"
      onkeydown={(e) => {
        if (e.key === "Escape") showEditModal = false;
      }}
    ></div>
    <div class="relative bg-white border border-slate-200 rounded-3xl shadow-2xl w-[95%] sm:w-full max-w-md overflow-hidden animate-in zoom-in-95 duration-200 flex flex-col max-h-[90vh]">
      <div class="p-4 sm:p-6 border-b border-slate-100 flex justify-between items-center bg-slate-50/50">
        <h3 class="text-xl font-bold text-slate-900">Edit Data</h3>
        <button
          onclick={() => showEditModal = false}
          class="text-slate-400 hover:text-slate-700 hover:bg-slate-100 p-2 rounded-full transition-colors cursor-pointer"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
        </button>
      </div>

      <form onsubmit={saveProgress} class="p-4 sm:p-6 space-y-5 overflow-y-auto">
        <div>
          <label class="block text-sm font-bold text-slate-700 mb-1.5" for="user">Pilih Murid</label>
          <select
            id="user"
            bind:value={formUserId}
            class="w-full bg-slate-50 border border-slate-200 focus:bg-white focus:border-indigo-400 focus:ring-2 focus:ring-indigo-400/20 rounded-xl pl-4 pr-10 py-2.5 text-sm text-slate-700 outline-none transition-all shadow-sm appearance-none bg-[url('data:image/svg+xml;charset=utf-8,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20fill%3D%22none%22%20viewBox%3D%220%200%2024%2024%22%20stroke%3D%22%2364748b%22%3E%3Cpath%20stroke-linecap%3D%22round%22%20stroke-linejoin%3D%22round%22%20stroke-width%3D%222%22%20d%3D%22M19%209l-7%207-7-7%22%2F%3E%3C%2Fsvg%3E')] bg-[length:1.25rem_1.25rem] bg-[position:right_1rem_center] bg-no-repeat cursor-pointer"
            required
          >
            <option value="" disabled hidden>-- Pilih Murid --</option>
            {#each users as u}
              <option value={u.id}>{u.username}</option>
            {/each}
          </select>
        </div>

        <div>
          <label class="block text-sm font-bold text-slate-700 mb-1.5" for="date">Tanggal</label>
          <input
            id="date"
            type="date"
            bind:value={formDate}
            required
            class="w-full bg-slate-50 border border-slate-200 focus:bg-white focus:border-indigo-400 focus:ring-2 focus:ring-indigo-400/20 rounded-xl px-4 py-2.5 text-sm text-slate-700 outline-none transition-all shadow-sm"
          />
        </div>

        <div>
          <label class="block text-sm font-bold text-slate-700 mb-1.5" for="image">Upload Foto (Opsional / Maks 1MB)</label>
          <input
            id="image"
            type="file"
            accept="image/*"
            onchange={handleImageUpload}
            disabled={isUploadingImage}
            class="w-full text-sm text-slate-600 file:mr-4 file:py-2.5 file:px-4 file:rounded-xl file:border-0 file:text-sm file:font-semibold file:bg-indigo-50 file:text-indigo-700 hover:file:bg-indigo-100 transition-all border border-slate-200 rounded-xl"
          />
          {#if isUploadingImage}
            <p class="mt-2 text-sm text-indigo-600 font-medium animate-pulse">Mengunggah gambar...</p>
          {/if}
          {#if localImagePreview || formImage}
            <div class="mt-3 rounded-xl overflow-hidden border border-slate-200 shadow-sm relative group">
              <img src={localImagePreview || formImage} alt="Preview" class="w-full h-40 object-cover" />
              <div class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                <span class="text-white text-xs font-medium px-3 py-1 bg-black/60 rounded-full">Preview</span>
              </div>
            </div>
          {/if}
        </div>

        <div class="pt-4 border-t border-slate-100 flex justify-end gap-3">
          <button
            type="button"
            onclick={() => showEditModal = false}
            class="px-5 py-2.5 text-sm font-semibold text-slate-600 hover:text-slate-900 hover:bg-slate-100 rounded-xl transition-colors cursor-pointer"
          >
            Batal
          </button>
          <button
            type="submit"
            disabled={isUploadingImage}
            class="px-5 py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white text-sm font-bold rounded-xl transition-all shadow-md shadow-indigo-600/20 disabled:opacity-50 cursor-pointer"
          >
            Simpan Perubahan
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}

<!-- Delete Modal -->
{#if showDeleteModal}
  <div class="fixed inset-0 z-[60] flex items-center justify-center p-4">
    <div
      class="absolute inset-0 bg-black/60 backdrop-blur-sm"
      onclick={() => showDeleteModal = false}
      role="button"
      tabindex="0"
    ></div>
    <div class="relative bg-white rounded-3xl p-6 sm:p-8 max-w-sm w-full text-center shadow-2xl animate-in zoom-in-95 duration-200 border border-slate-100">
      <div class="w-16 h-16 bg-red-100 text-red-600 rounded-full flex items-center justify-center mx-auto mb-4 shadow-inner border border-red-200">
        <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path></svg>
      </div>
      <h3 class="text-xl font-bold text-slate-900 mb-2">Hapus Data?</h3>
      <p class="text-slate-500 text-sm mb-8">
        Apakah Anda yakin ingin menghapus data perkembangan ini? Tindakan ini tidak dapat dibatalkan.
      </p>
      <div class="flex gap-3 justify-center">
        <button
          onclick={() => showDeleteModal = false}
          class="px-5 py-2.5 bg-slate-100 hover:bg-slate-200 text-slate-700 text-sm font-semibold rounded-xl transition-colors w-full cursor-pointer"
        >
          Batal
        </button>
        <button
          onclick={executeDelete}
          class="px-5 py-2.5 bg-red-600 hover:bg-red-700 text-white text-sm font-bold rounded-xl transition-all shadow-md shadow-red-600/20 w-full cursor-pointer"
        >
          Ya, Hapus
        </button>
      </div>
    </div>
  </div>
{/if}
