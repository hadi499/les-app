<script>
  import { onMount } from 'svelte';

  let subjects = $state([]);
  let isLoading = $state(true);
  let errorMsg = $state('');

  // Modal State
  let showModal = $state(false);
  let isEditing = $state(false);
  let currentSubjectId = $state(null);
  
  // Form State
  let formName = $state('');

  async function fetchSubjects() {
    isLoading = true;
    errorMsg = '';
    try {
      const res = await fetch('http://localhost:8080/api/subjects', { credentials: 'include' });
      if (!res.ok) throw new Error('Gagal mengambil data mata pelajaran');
      subjects = await res.json() || [];
    } catch(e) {
      errorMsg = e.message;
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
    formName = '';
    showModal = true;
  }

  function openEditModal(subject) {
    isEditing = true;
    currentSubjectId = subject.id;
    formName = subject.name;
    showModal = true;
  }

  function closeModal() {
    showModal = false;
  }

  async function saveSubject(e) {
    e.preventDefault();
    if (!formName.trim()) {
      alert("Nama mata pelajaran tidak boleh kosong!");
      return;
    }

    const payload = {
      name: formName.trim()
    };

    try {
      let url = 'http://localhost:8080/api/subjects';
      let method = 'POST';

      if (isEditing) {
        url = `http://localhost:8080/api/subjects/${currentSubjectId}`;
        method = 'PUT';
      }

      const res = await fetch(url, {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
        credentials: 'include'
      });

      if (!res.ok) {
        const errData = await res.json();
        throw new Error(errData.error || 'Gagal menyimpan data');
      }

      closeModal();
      await fetchSubjects();
    } catch (err) {
      alert(err.message);
    }
  }

  async function deleteSubject(id) {
    if (!confirm('Apakah Anda yakin ingin menghapus mata pelajaran ini? Jika ada nilai yang menggunakan mata pelajaran ini, penghapusan mungkin akan gagal.')) return;

    try {
      const res = await fetch(`http://localhost:8080/api/subjects/${id}`, {
        method: 'DELETE',
        credentials: 'include'
      });

      if (!res.ok) {
        const errData = await res.json();
        throw new Error(errData.error || 'Gagal menghapus data');
      }
      await fetchSubjects();
    } catch (err) {
      alert(err.message);
    }
  }
</script>

<svelte:head>
  <title>Manajemen Mata Pelajaran - Portal Guru</title>
</svelte:head>

<div class="animate-in fade-in duration-500 max-w-5xl mx-auto space-y-6">
  <div class="flex flex-col sm:flex-row sm:items-end justify-between gap-4 mb-8">
    <div>
      <h1 class="text-2xl font-bold text-orange-950 sm:text-3xl tracking-tight drop-shadow-sm">Mata Pelajaran</h1>
      <p class="mt-2 text-orange-800 text-sm sm:text-base font-light tracking-wide">Kelola daftar mata pelajaran untuk nilai ujian harian.</p>
    </div>
    <button
      onclick={openAddModal}
      class="inline-flex items-center gap-2 px-4 py-2.5 bg-indigo-600 hover:bg-indigo-500 text-orange-950 font-medium rounded-xl transition-all shadow-md shadow-indigo-900/20"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
      Tambah Mata Pelajaran
    </button>
  </div>

  {#if isLoading}
    <div class="flex justify-center p-12">
      <div class="w-10 h-10 border-4 border-orange-200 border-t-orange-500 rounded-full animate-spin"></div>
    </div>
  {:else if errorMsg}
    <div class="bg-red-100 text-red-800 p-6 rounded-2xl border border-red-300 font-medium flex items-center gap-3">
      <svg class="w-6 h-6 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
      {errorMsg}
    </div>
  {:else}
    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
      {#each subjects as subject}
        <div class="bg-white/60 backdrop-blur-md rounded-2xl border border-orange-200 p-5 shadow-lg shadow-orange-900/5 flex flex-col justify-between hover:border-orange-300 transition-colors">
          <div class="flex items-center gap-3 mb-4">
            <div class="w-10 h-10 rounded-full bg-indigo-900/30 text-indigo-400 flex items-center justify-center font-bold text-lg border border-indigo-800/50">
              {subject.name.charAt(0).toUpperCase()}
            </div>
            <h3 class="text-lg font-bold text-orange-950 truncate">{subject.name}</h3>
          </div>
          
          <div class="flex justify-end gap-2 pt-3 border-t border-orange-200">
            <button 
              onclick={() => openEditModal(subject)}
              class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-orange-600 bg-blue-100 hover:bg-blue-800/50 rounded-lg transition-colors border border-blue-300"
            >
              Edit
            </button>
            <button 
              onclick={() => deleteSubject(subject.id)}
              class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-red-600 bg-red-100 hover:bg-red-200 rounded-lg transition-colors border border-red-300"
            >
              Hapus
            </button>
          </div>
        </div>
      {/each}
      {#if subjects.length === 0}
        <div class="col-span-full py-12 text-center text-orange-700 font-light bg-[#EAE4BD]/30 rounded-3xl border border-orange-200">
          Belum ada mata pelajaran.
        </div>
      {/if}
    </div>
  {/if}
</div>

<!-- Modal Dialog -->
{#if showModal}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
    <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" onclick={closeModal}></div>
    <div class="relative bg-[#EAE4BD] border border-orange-300 rounded-2xl shadow-2xl w-full max-w-md overflow-hidden animate-in zoom-in-95 duration-200">
      <div class="p-6 border-b border-orange-200 flex justify-between items-center">
        <h3 class="text-xl font-bold text-orange-950">{isEditing ? 'Edit Mata Pelajaran' : 'Tambah Mata Pelajaran'}</h3>
        <button onclick={closeModal} class="text-orange-700 hover:text-orange-950 transition-colors">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
        </button>
      </div>
      
      <form onsubmit={saveSubject} class="p-6 space-y-4">
        <div>
          <label class="block text-sm font-medium text-orange-800 mb-1.5" for="name">Nama Mata Pelajaran</label>
          <input type="text" id="name" bind:value={formName} placeholder="Misal: Matematika, Bahasa Inggris" class="w-full bg-white border border-orange-300 rounded-xl px-4 py-2.5 text-orange-950 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all" required>
        </div>

        <div class="pt-4 flex justify-end gap-3">
          <button type="button" onclick={closeModal} class="px-4 py-2.5 text-sm font-medium text-orange-900 hover:text-orange-950 bg-white hover:bg-orange-50 rounded-xl transition-colors">
            Batal
          </button>
          <button type="submit" class="px-4 py-2.5 text-sm font-medium text-orange-950 bg-indigo-600 hover:bg-indigo-500 rounded-xl transition-all shadow-md shadow-indigo-900/20">
            {isEditing ? 'Simpan' : 'Tambahkan'}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}
