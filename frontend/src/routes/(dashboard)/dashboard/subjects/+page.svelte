<script lang="ts">
  import { onMount } from "svelte";

  type Subject = { id: number; name: string };

  let subjects: Subject[] = $state([]);
  let isLoading = $state(true);
  let errorMsg = $state("");

  // Modal State
  let showModal = $state(false);
  let isEditing = $state(false);
  let currentSubjectId: number | null = $state(null);

  let showDeleteModal = $state(false);
  let subjectToDelete: number | null = $state(null);

  // Form State
  let formName = $state("");

  async function fetchSubjects() {
    isLoading = true;
    errorMsg = "";
    try {
      const res = await fetch(`/api/subjects`, {
        credentials: "include",
      });
      if (!res.ok) throw new Error("Gagal mengambil data mata pelajaran");
      subjects = ((await res.json()) as Subject[]) || [];
    } catch (e) {
      errorMsg = e instanceof Error ? e.message : String(e);
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    fetchSubjects();
  });

  function openAddModal() {
    isEditing = false;
    currentSubjectId = null;
    formName = "";
    showModal = true;
  }

  function openEditModal(subject: Subject) {
    isEditing = true;
    currentSubjectId = subject.id;
    formName = subject.name;
    showModal = true;
  }

  function closeModal() {
    showModal = false;
  }

  async function saveSubject(e: Event) {
    e.preventDefault();
    if (!formName.trim()) {
      alert("Nama mata pelajaran tidak boleh kosong!");
      return;
    }

    const payload = {
      name: formName.trim(),
    };

    try {
      let url = `/api/subjects`;
      let method = "POST";

      if (isEditing) {
        url = `/api/subjects/${currentSubjectId}`;
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
      await fetchSubjects();
    } catch (err) {
      alert(err instanceof Error ? err.message : String(err));
    }
  }

  function openDeleteModal(id: number) {
    subjectToDelete = id;
    showDeleteModal = true;
  }

  function closeDeleteModal() {
    showDeleteModal = false;
    subjectToDelete = null;
  }

  async function executeDelete() {
    if (subjectToDelete === null) return;

    try {
      const res = await fetch(`/api/subjects/${subjectToDelete}`, {
        method: "DELETE",
        credentials: "include",
      });

      if (!res.ok) {
        const errData = await res.json();
        throw new Error(errData.error || "Gagal menghapus data");
      }
      await fetchSubjects();
      closeDeleteModal();
    } catch (err) {
      alert(err instanceof Error ? err.message : String(err));
      closeDeleteModal();
    }
  }
</script>

<svelte:head>
  <title>Manajemen Mata Pelajaran - Portal Guru</title>
</svelte:head>

<div class="animate-in fade-in duration-500 max-w-5xl mx-auto space-y-6">
  <div
    class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 mb-8"
  >
    <div>
      <h1
        class="text-2xl font-bold text-slate-900 sm:text-3xl tracking-tight drop-shadow-sm"
      >
        Mata Pelajaran
      </h1>
      <p
        class="mt-2 text-slate-600 text-sm sm:text-base font-light tracking-wide"
      >
        Kelola daftar mata pelajaran untuk nilai ujian harian.
      </p>
    </div>
    <button
      onclick={openAddModal}
      class="inline-flex items-center self-start sm:self-auto gap-2 px-4 py-2.5 bg-indigo-300 hover:bg-indigo-200 text-slate-900 font-medium rounded-xl transition-all shadow-md shadow-indigo-900/20"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"
        ><path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M12 4v16m8-8H4"
        ></path></svg
      >
      Pelajaran
    </button>
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
    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
      {#each subjects as subject}
        <div
          class="bg-white/60 backdrop-blur-md rounded-2xl border border-slate-200 p-5 shadow-lg shadow-slate-800/5 flex flex-col justify-between hover:border-slate-300 transition-colors"
        >
          <div class="flex items-center gap-3 mb-4">
            <div
              class="w-10 h-10 rounded-full bg-indigo-100 text-indigo-600 flex items-center justify-center font-bold text-lg border border-indigo-800/50"
            >
              {subject.name.charAt(0).toUpperCase()}
            </div>
            <h3 class="text-lg font-bold text-slate-900 truncate">
              {subject.name}
            </h3>
          </div>

          <div class="flex justify-end gap-2 pt-3 border-t border-slate-200">
            <button
              onclick={() => openEditModal(subject)}
              class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-blue-600 bg-blue-100 hover:bg-blue-500/10 rounded-lg transition-colors border border-blue-300"
            >
              Edit
            </button>
            <button
              onclick={() => openDeleteModal(subject.id)}
              class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-red-600 bg-red-100 hover:bg-red-200 rounded-lg transition-colors border border-red-300"
            >
              Hapus
            </button>
          </div>
        </div>
      {/each}
      {#if subjects.length === 0}
        <div
          class="col-span-full py-12 text-center text-slate-500 font-light bg-[#EAE4BD]/30 rounded-3xl border border-slate-200"
        >
          Belum ada mata pelajaran.
        </div>
      {/if}
    </div>
  {/if}
</div>

<!-- Modal Dialog -->
{#if showModal}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
    <div
      class="absolute inset-0 bg-black/60 backdrop-blur-sm"
      onclick={closeModal}
    ></div>
    <div
      class="relative bg-slate-50 border border-slate-300 rounded-2xl shadow-2xl w-full max-w-md overflow-hidden animate-in zoom-in-95 duration-200"
    >
      <div
        class="p-6 border-b border-slate-200 flex justify-between items-center"
      >
        <h3 class="text-xl font-bold text-slate-900">
          {isEditing ? "Edit Mata Pelajaran" : "Tambah Mata Pelajaran"}
        </h3>
        <button
          onclick={closeModal}
          class="text-slate-500 hover:text-slate-900 transition-colors"
        >
          <svg
            class="w-6 h-6"
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

      <form onsubmit={saveSubject} class="p-6 space-y-4">
        <div>
          <label
            class="block text-sm font-medium text-slate-600 mb-1.5"
            for="name">Nama Mata Pelajaran</label
          >
          <input
            type="text"
            id="name"
            bind:value={formName}
            placeholder="Misal: Matematika, Bahasa Inggris"
            class="w-full bg-white border border-slate-300 rounded-xl px-4 py-2.5 text-slate-900 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all"
            required
          />
        </div>

        <div class="pt-4 flex justify-end gap-3">
          <button
            type="button"
            onclick={closeModal}
            class="px-4 py-2.5 text-sm font-medium text-slate-800 hover:text-slate-900 bg-white hover:bg-slate-100 rounded-xl transition-colors shadow-md"
          >
            Batal
          </button>
          <button
            type="submit"
            class="px-4 py-2.5 text-sm font-medium text-slate-900 bg-indigo-100 hover:bg-indigo-200 rounded-xl transition-all shadow-md shadow-indigo-900/20"
          >
            {isEditing ? "Simpan" : "Tambahkan"}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}

<!-- Delete Confirmation Modal -->
{#if showDeleteModal}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
    <div
      class="absolute inset-0 bg-black/60 backdrop-blur-sm"
      onclick={closeDeleteModal}
    ></div>
    <div
      class="relative bg-slate-50 border border-red-300 rounded-2xl shadow-2xl w-full max-w-md overflow-hidden animate-in zoom-in-95 duration-200 p-6 text-center"
    >
      <div
        class="w-12 h-12 mx-auto rounded-full bg-red-100 flex items-center justify-center mb-4"
      >
        <svg
          class="w-6 h-6 text-red-600"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
          />
        </svg>
      </div>
      <h3 class="text-xl font-bold text-slate-900 mb-2">
        Hapus Mata Pelajaran?
      </h3>
      <p class="text-sm text-slate-600 mb-6">
        Apakah Anda yakin ingin menghapus mata pelajaran ini? Jika ada nilai
        yang menggunakan mata pelajaran ini, penghapusan mungkin akan gagal.
      </p>
      <div class="flex justify-center gap-3">
        <button
          onclick={closeDeleteModal}
          class="px-4 py-2.5 text-sm font-medium text-slate-800 bg-white border border-slate-300 rounded-xl hover:bg-slate-50 transition-colors"
        >
          Batal
        </button>
        <button
          onclick={executeDelete}
          class="px-4 py-2.5 text-sm font-medium text-white bg-red-600 hover:bg-red-700 rounded-xl transition-all shadow-md shadow-red-900/20"
        >
          Ya, Hapus
        </button>
      </div>
    </div>
  </div>
{/if}
