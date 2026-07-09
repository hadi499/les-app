<script lang="ts">
  import { onMount } from "svelte";
  import { page } from "$app/state";
  import { compressImageFile } from "$lib/utils";

  type WritingProgress = {
    id: number;
    user_id: number;
    date: string;
    image: string;
    user?: { username: string };
  };
  type User = { id: number; username: string; role?: string };

  let progresses: WritingProgress[] = $state([]);
  let users: User[] = $state([]);
  let isLoading = $state(true);
  let errorMsg = $state("");
  let isTeacher = $state(false);

  // Pagination State
  let currentPage = $state(1);
  const itemsPerPage = 50;
  let totalPages = $state(1);
  let totalRecords = $state(0);

  function nextPage() {
    if (currentPage < totalPages) {
      currentPage++;
      fetchProgresses();
    }
  }

  function prevPage() {
    if (currentPage > 1) {
      currentPage--;
      fetchProgresses();
    }
  }

  function goToPage(page: number) {
    currentPage = page;
    fetchProgresses();
  }

  // Modal State
  let showModal = $state(false);
  let isEditing = $state(false);
  let currentProgressId: number | null = $state(null);

  let showDeleteModal = $state(false);
  let progressToDelete: number | null = $state(null);

  // Form State
  let formUserId: number | string = $state("");
  let formDate = $state("");
  let formImage = $state("");
  let formFile = $state<File | null>(null);
  let localImagePreview = $state("");
  let isUploadingImage = $state(false);

  async function handleImageUpload(e: Event) {
    const target = e.target as HTMLInputElement;
    if (!target.files || target.files.length === 0) return;

    isUploadingImage = true;
    try {
      let file = target.files[0];

      // Compress image
      file = await compressImageFile(file);

      if (file.size > 1 * 1024 * 1024) {
        alert(
          "Ukuran file setelah kompresi masih lebih dari 1MB. Silakan pilih gambar lain.",
        );
        return;
      }

      formFile = file;
      localImagePreview = URL.createObjectURL(file);
    } finally {
      isUploadingImage = false;
    }
  }

  async function fetchProgresses() {
    try {
      const res = await fetch(`/api/writing-progress?page=${currentPage}&limit=${itemsPerPage}`, {
        credentials: "include",
      });
      if (!res.ok) throw new Error("Gagal mengambil data perkembangan menulis");
      const data = await res.json();
      progresses = data.data || [];
      totalPages = data.total_pages || 1;
      totalRecords = data.total || 0;
    } catch (e) {
      errorMsg = e instanceof Error ? e.message : String(e);
    }
  }

  async function fetchUsers() {
    try {
      const res = await fetch(`/api/users`, {
        credentials: "include",
      });
      if (!res.ok) throw new Error("Gagal mengambil data murid");
      const data = (await res.json()) as { users: User[] };
      users = data.users.filter((u) => u.role === "student") || [];
    } catch (e) {
      console.error(e);
    }
  }

  async function checkRole() {
    try {
      const res = await fetch(`/me`, { credentials: "include" });
      const data = await res.json();
      if (data.authenticated && data.user) {
        if (data.user.role === "teacher") {
          isTeacher = true;
        }
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function loadData() {
    isLoading = true;
    errorMsg = "";
    await checkRole();
    const promises: Promise<any>[] = [fetchProgresses()];
    if (isTeacher) {
      promises.push(fetchUsers());
    }
    await Promise.all(promises);
    isLoading = false;
  }

  onMount(() => {
    loadData();
  });

  function openAddModal() {
    isEditing = false;
    currentProgressId = null;
    formUserId = "";
    formDate = new Date().toLocaleDateString("en-CA");
    formImage = "";
    formFile = null;
    localImagePreview = "";
    showModal = true;
  }

  function openEditModal(progress: WritingProgress) {
    isEditing = true;
    currentProgressId = progress.id;
    formUserId = progress.user_id;
    formDate = progress.date ? progress.date.split("T")[0] : "";
    formImage = progress.image || "";
    formFile = null;
    localImagePreview = "";
    showModal = true;
  }

  function closeModal() {
    showModal = false;
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
        if (!res.ok) {
          const err = await res.json();
          throw new Error(err.error || "Gagal mengunggah gambar");
        }
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
      let url = `/api/writing-progress`;
      let method = "POST";

      if (isEditing) {
        url = `/api/writing-progress/${currentProgressId}`;
        method = "PUT";
      }

      const res = await fetch(url, {
        method,
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
        credentials: "include",
      });

      if (!res.ok) {
        const errData = await res.json();
        throw new Error(errData.error || "Gagal menyimpan data");
      }

      closeModal();
      await fetchProgresses();
    } catch (err) {
      alert(err instanceof Error ? err.message : String(err));
    } finally {
      isUploadingImage = false;
    }
  }

  function openDeleteModal(id: number) {
    progressToDelete = id;
    showDeleteModal = true;
  }

  function closeDeleteModal() {
    showDeleteModal = false;
    progressToDelete = null;
  }

  async function executeDelete() {
    if (progressToDelete === null) return;

    try {
      const res = await fetch(`/api/writing-progress/${progressToDelete}`, {
        method: "DELETE",
        credentials: "include",
      });

      if (!res.ok) throw new Error("Gagal menghapus data");
      await fetchProgresses();
      closeDeleteModal();
    } catch (err) {
      alert(err instanceof Error ? err.message : String(err));
      closeDeleteModal();
    }
  }

  function formatDate(dateStr: string | null) {
    if (!dateStr) return "-";
    return new Date(dateStr).toLocaleDateString("id-ID", {
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  }
</script>

<svelte:head>
  <title>Perkembangan Menulis - Portal {isTeacher ? "Guru" : "Siswa"}</title>
</svelte:head>

<div
  class="animate-in fade-in duration-500 max-w-4xl mx-auto space-y-6 p-4 sm:p-0"
>
  <div
    class="flex flex-col sm:flex-row items-start sm:items-end justify-between gap-4 mb-8"
  >
    <div>
      <h1
        class="text-2xl font-bold text-slate-900 sm:text-3xl tracking-tight drop-shadow-sm"
      >
        Perkembangan Menulis
      </h1>
      <p class="mt-2 text-slate-600 text-sm sm:text-base tracking-wide">
        Pantau progres menulis anak PAUD dan TK.
      </p>
    </div>
    {#if isTeacher}
      <button
        onclick={openAddModal}
        class="inline-flex items-center gap-2 px-4 py-2.5 bg-indigo-300 hover:bg-indigo-200 text-slate-900 font-medium rounded-xl transition-all shadow-md shadow-indigo-900/20 cursor-pointer"
      >
        <svg
          class="w-5 h-5"
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
        Tambah Data
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
    {#if progresses.length === 0}
      <div
        class="py-16 text-center bg-white/60 backdrop-blur-md rounded-xl border border-slate-200 shadow-sm flex flex-col items-center justify-center"
      >
        <div
          class="w-16 h-16 bg-slate-100 rounded-full flex items-center justify-center text-slate-400 mb-4"
        >
          <svg
            class="w-8 h-8"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            ><path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
            ></path></svg
          >
        </div>
        <p class="text-slate-500 font-medium text-lg">
          Belum ada data perkembangan menulis.
        </p>
        {#if isTeacher}
          <p class="text-slate-400 text-sm mt-1">
            Silakan tambahkan data baru melalui tombol di atas.
          </p>
        {/if}
      </div>
    {:else}
      <div
        class="bg-white/80 backdrop-blur-md rounded-2xl shadow-sm border border-slate-200 overflow-hidden"
      >
        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse">
            <thead>
              <tr
                class="bg-slate-50/50 border-b border-slate-200 text-slate-500 text-sm"
              >
                {#if isTeacher}
                  <th class="px-6 py-4 font-semibold">Siswa</th>
                {/if}
                <th class="px-6 py-4 font-semibold">Tanggal</th>
                <th class="px-6 py-4 font-semibold text-center w-32">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              {#each progresses as progress}
                <tr class="hover:bg-slate-50/50 transition-colors group">
                  {#if isTeacher}
                    <td class="px-6 py-4">
                      <div class="flex items-center">
                        <span class="font-bold text-slate-800">
                          {progress.user?.username || "Siswa Tidak Diketahui"}
                        </span>
                      </div>
                    </td>
                  {/if}
                  <td class="px-6 py-4">
                    <div
                      class="flex items-center gap-2 text-sm text-slate-600 font-medium"
                    >
                      {formatDate(progress.date)}
                    </div>
                  </td>
                  <td class="px-6 py-4">
                    <div class="flex items-center justify-center gap-2">
                      <a
                        href={`/dashboard/writing-progress/${progress.id}`}
                        title="Lihat Detail"
                        class="inline-flex items-center justify-center w-8 h-8 text-green-700 bg-green-50 hover:bg-green-100 border border-green-100 rounded-xl transition-colors cursor-pointer"
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
                            d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                          ></path><path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                          ></path></svg
                        >
                      </a>
                      {#if isTeacher}
                        <button
                          onclick={() => openEditModal(progress)}
                          title="Edit"
                          class="inline-flex items-center justify-center w-8 h-8 text-blue-700 bg-blue-50 hover:bg-blue-100 rounded-xl transition-colors border border-blue-100 cursor-pointer"
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
                              d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                            ></path></svg
                          >
                        </button>
                        <button
                          onclick={() => openDeleteModal(progress.id)}
                          title="Hapus"
                          class="inline-flex items-center justify-center w-8 h-8 text-red-700 bg-red-50 hover:bg-red-100 rounded-xl transition-colors border border-red-100 cursor-pointer"
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
                              d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                            ></path></svg
                          >
                        </button>
                      {/if}
                    </div>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>
    {/if}

    <!-- Pagination Controls -->
    {#if totalPages > 1}
      <div
        class="mt-8 px-6 py-4 bg-white/60 backdrop-blur-md rounded-2xl border border-slate-200 shadow-sm flex flex-col sm:flex-row items-center justify-between gap-4"
      >
        <div class="text-sm text-slate-600 text-center sm:text-left">
          Menampilkan <span class="font-medium text-slate-900"
            >{totalRecords > 0
              ? (currentPage - 1) * itemsPerPage + 1
              : 0}</span
          >
          sampai
          <span class="font-medium text-slate-900"
            >{Math.min(currentPage * itemsPerPage, totalRecords)}</span
          >
          dari
          <span class="font-medium text-slate-900">{totalRecords}</span> data
        </div>
        <div class="flex gap-2">
          <button
            onclick={prevPage}
            disabled={currentPage === 1}
            class="px-3 py-1.5 text-sm font-medium rounded-lg transition-colors cursor-pointer {currentPage ===
            1
              ? 'bg-white/40 text-zinc-600 cursor-not-allowed'
              : 'bg-white text-slate-800 hover:bg-slate-50 hover:text-slate-900'}"
          >
            Sebelumnya
          </button>
          <div class="flex items-center gap-1 px-2 flex-wrap justify-center">
            {#each Array(totalPages) as _, i}
              <button
                onclick={() => goToPage(i + 1)}
                class="w-8 h-8 flex items-center justify-center text-sm font-medium rounded-lg transition-colors cursor-pointer {currentPage ===
                i + 1
                  ? 'bg-indigo-600 text-slate-100'
                  : 'text-slate-600 hover:bg-white hover:text-slate-900'}"
              >
                {i + 1}
              </button>
            {/each}
          </div>
          <button
            onclick={nextPage}
            disabled={currentPage === totalPages}
            class="px-3 py-1.5 text-sm font-medium rounded-lg transition-colors cursor-pointer {currentPage ===
            totalPages
              ? 'bg-white/40 text-zinc-600 cursor-not-allowed'
              : 'bg-white text-slate-800 hover:bg-slate-50 hover:text-slate-900'}"
          >
            Selanjutnya
          </button>
        </div>
      </div>
    {/if}
  {/if}
</div>

<!-- Modal Dialog -->
{#if showModal}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
    <div
      class="absolute inset-0 bg-black/60 backdrop-blur-sm"
      onclick={closeModal}
      role="button"
      tabindex="0"
      onkeydown={(e) => {
        if (e.key === "Escape") closeModal();
      }}
    ></div>
    <div
      class="relative bg-white border border-slate-200 rounded-3xl shadow-2xl w-[95%] sm:w-full max-w-md overflow-hidden animate-in zoom-in-95 duration-200 flex flex-col max-h-[90vh]"
    >
      <div
        class="p-4 sm:p-6 border-b border-slate-100 flex justify-between items-center bg-slate-50/50"
      >
        <h3 class="text-xl font-bold text-slate-900">
          {isEditing ? "Edit Data" : "Tambah Data"}
        </h3>
        <button
          onclick={closeModal}
          class="text-slate-400 hover:text-slate-700 hover:bg-slate-100 p-2 rounded-full transition-colors cursor-pointer"
        >
          <svg
            class="w-5 h-5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            ><path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M6 18L18 6M6 6l12 12"
            ></path></svg
          >
        </button>
      </div>

      <form
        onsubmit={saveProgress}
        class="p-4 sm:p-6 space-y-5 overflow-y-auto"
      >
        <div>
          <label
            class="block text-sm font-bold text-slate-700 mb-1.5"
            for="user">Pilih Murid</label
          >
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
          <label
            class="block text-sm font-bold text-slate-700 mb-1.5"
            for="date">Tanggal</label
          >
          <input
            id="date"
            type="date"
            bind:value={formDate}
            required
            class="w-full bg-slate-50 border border-slate-200 focus:bg-white focus:border-indigo-400 focus:ring-2 focus:ring-indigo-400/20 rounded-xl px-4 py-2.5 text-sm text-slate-700 outline-none transition-all shadow-sm"
          />
        </div>

        <div>
          <label
            class="block text-sm font-bold text-slate-700 mb-1.5"
            for="image">Upload Foto (Opsional / Maks 1MB)</label
          >
          <input
            id="image"
            type="file"
            accept="image/*"
            onchange={handleImageUpload}
            disabled={isUploadingImage}
            class="w-full text-sm text-slate-600 file:mr-4 file:py-2.5 file:px-4 file:rounded-xl file:border-0 file:text-sm file:font-semibold file:bg-indigo-50 file:text-indigo-700 hover:file:bg-indigo-100 transition-all border border-slate-200 rounded-xl"
          />
          {#if isUploadingImage}
            <p class="mt-2 text-sm text-indigo-600 font-medium animate-pulse">
              Mengunggah gambar...
            </p>
          {/if}
          {#if localImagePreview || formImage}
            <div
              class="mt-3 rounded-xl overflow-hidden border border-slate-200 shadow-sm relative group"
            >
              <img
                src={localImagePreview || formImage}
                alt="Preview"
                class="w-full h-40 object-cover"
              />
              <div
                class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center"
              >
                <span
                  class="text-white text-xs font-medium px-3 py-1 bg-black/60 rounded-full"
                  >Preview</span
                >
              </div>
            </div>
          {/if}
        </div>

        <div class="pt-4 border-t border-slate-100 flex justify-end gap-3">
          <button
            type="button"
            onclick={closeModal}
            class="px-5 py-2.5 text-sm font-semibold text-slate-600 hover:text-slate-900 hover:bg-slate-100 rounded-xl transition-colors cursor-pointer"
          >
            Batal
          </button>
          <button
            type="submit"
            disabled={isUploadingImage}
            class="px-5 py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white text-sm font-bold rounded-xl transition-all shadow-md shadow-indigo-600/20 disabled:opacity-50 cursor-pointer"
          >
            {isEditing ? "Simpan Perubahan" : "Simpan Data"}
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
      onclick={closeDeleteModal}
      role="button"
      tabindex="0"
      onkeydown={(e) => {
        if (e.key === "Escape") closeDeleteModal();
      }}
    ></div>
    <div
      class="relative bg-white border border-slate-200 rounded-3xl shadow-2xl w-full max-w-sm p-6 text-center animate-in zoom-in-95 duration-200"
    >
      <div
        class="w-16 h-16 bg-red-100 text-red-600 rounded-full flex items-center justify-center mx-auto mb-4 shadow-sm"
      >
        <svg
          class="w-8 h-8"
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
      </div>
      <h3 class="text-xl font-bold text-slate-900 mb-2">Hapus Data?</h3>
      <p class="text-slate-600 mb-6 font-medium">
        Tindakan ini tidak dapat dibatalkan.
      </p>
      <div class="flex justify-center gap-3">
        <button
          onclick={closeDeleteModal}
          class="px-5 py-2.5 text-sm font-bold text-slate-700 bg-slate-100 hover:bg-slate-200 rounded-xl transition-colors cursor-pointer w-full"
        >
          Batal
        </button>
        <button
          onclick={executeDelete}
          class="px-5 py-2.5 text-sm font-bold text-white bg-red-600 hover:bg-red-700 shadow-md shadow-red-600/20 rounded-xl transition-colors cursor-pointer w-full"
        >
          Ya, Hapus
        </button>
      </div>
    </div>
  </div>
{/if}
