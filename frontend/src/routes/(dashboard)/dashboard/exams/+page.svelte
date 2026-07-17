<script lang="ts">
  import { onMount } from "svelte";
  import { compressImageFile } from "$lib/utils";

  type Exam = {
    id: number;
    user_id: number;
    exam_date: string;
    exam_name: string;
    subject_id: number;
    score: number;
    image?: string;
    user?: { username: string };
    subject?: { name: string };
  };
  type User = { id: number; username: string; role?: string };
  type Subject = { id: number; name: string };

  let exams: Exam[] = $state([]);
  let users: User[] = $state([]);
  let subjects: Subject[] = $state([]);
  let isLoading = $state(true);
  let errorMsg = $state("");
  let isTeacher = $state(false);
  let activeTab: "card" | "table" = $state("card");
  let isLoadingMore = $state(false);

  // Pagination State
  let currentPage = $state(1);
  let itemsPerPage = $state(20);
  let totalPages = $state(1);
  let totalRecords = $state(0);

  function switchTab(tab: "card" | "table") {
    if (activeTab === tab) return;
    activeTab = tab;
    currentPage = 1;
    localStorage.setItem("examsActiveTab", tab);
    fetchExams();
  }

  function infiniteScroll(node: HTMLElement) {
    const observer = new IntersectionObserver((entries) => {
      if (
        entries[0].isIntersecting &&
        activeTab === "card" &&
        currentPage < totalPages &&
        !isLoadingMore &&
        !isLoading
      ) {
        isLoadingMore = true;
        currentPage++;
        fetchExams().finally(() => { isLoadingMore = false; });
      }
    });
    observer.observe(node);
    return {
      destroy() {
        observer.disconnect();
      },
    };
  }

  function nextPage() {
    if (currentPage < totalPages) {
      currentPage++;
      fetchExams();
    }
  }

  function prevPage() {
    if (currentPage > 1) {
      currentPage--;
      fetchExams();
    }
  }

  function goToPage(page: number) {
    currentPage = page;
    fetchExams();
  }

  // Modal State
  let showModal = $state(false);
  let isEditing = $state(false);
  let currentExamId: number | null = $state(null);

  let showDeleteModal = $state(false);
  let examToDelete: number | null = $state(null);
  let isBulkDelete = $state(false);

  // Form State
  let formUserId: number | string = $state("");
  let formExamDate = $state("");
  let formExamName = $state("");
  let formSubjectId: number | string = $state("");
  let formScore = $state(0);
  let formImage = $state("");
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

      const formData = new FormData();
      formData.append("image", file);

      const res = await fetch("/api/upload?type=exam", {
        method: "POST",
        body: formData,
        credentials: "include",
      });
      if (!res.ok) {
        const err = await res.json();
        throw new Error(err.error || "Gagal mengunggah gambar");
      }
      const data = await res.json();
      formImage = data.url;
    } catch (err: any) {
      alert(err.message || "Terjadi kesalahan");
    } finally {
      isUploadingImage = false;
    }
  }

  async function fetchExams() {
    try {
      const res = await fetch(`/api/exams?page=${currentPage}&limit=${itemsPerPage}`, {
        credentials: "include",
      });
      if (!res.ok) throw new Error("Gagal mengambil data ujian");
      const data = await res.json();
      
      if (activeTab === "card" && currentPage > 1) {
        exams = [...exams, ...(data.data || [])];
      } else {
        exams = data.data || [];
      }
      
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

  async function fetchSubjects() {
    try {
      const res = await fetch(`/api/subjects`, {
        credentials: "include",
      });
      if (!res.ok) throw new Error("Gagal mengambil data mapel");
      subjects = (await res.json()) || [];
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
    const promises: Promise<any>[] = [fetchExams(), fetchSubjects()];
    if (isTeacher) {
      promises.push(fetchUsers());
    }
    await Promise.all(promises);
    isLoading = false;
  }

  onMount(() => {
    if (window.innerWidth < 768) {
      itemsPerPage = 10;
    } else {
      itemsPerPage = 20;
    }

    const savedTab = localStorage.getItem("examsActiveTab");
    if (savedTab === "card" || savedTab === "table") {
      activeTab = savedTab;
    }
    loadData();
  });

  function openAddModal() {
    isEditing = false;
    currentExamId = null;
    formUserId = "";
    formExamDate = new Date().toLocaleDateString("en-CA");
    formExamName = "";
    formSubjectId = "";
    formScore = 0;
    formImage = "";
    showModal = true;
  }

  function openEditModal(exam: Exam) {
    isEditing = true;
    currentExamId = exam.id;
    formUserId = exam.user_id;
    formExamDate = exam.exam_date ? exam.exam_date.split("T")[0] : "";
    formExamName = exam.exam_name;
    formSubjectId = exam.subject_id;
    formScore = exam.score;
    formImage = exam.image || "";
    showModal = true;
  }

  function closeModal() {
    showModal = false;
  }

  async function saveExam(e: Event) {
    e.preventDefault();
    if (!formUserId || !formExamDate || !formExamName || !formSubjectId) {
      alert("Harap lengkapi semua data!");
      return;
    }

    const payload = {
      user_id: Number(formUserId),
      exam_date: formExamDate,
      exam_name: formExamName,
      subject_id: Number(formSubjectId),
      score: formScore,
      image: formImage,
    };

    try {
      let url = `/api/exams`;
      let method = "POST";

      if (isEditing) {
        url = `/api/exams/${currentExamId}`;
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
      await fetchExams();
    } catch (err) {
      alert(err instanceof Error ? err.message : String(err));
    }
  }

  function openDeleteModal(id: number) {
    isBulkDelete = false;
    examToDelete = id;
    showDeleteModal = true;
  }

  function openBulkDeleteModal() {
    isBulkDelete = true;
    showDeleteModal = true;
  }

  function closeDeleteModal() {
    showDeleteModal = false;
    examToDelete = null;
    isBulkDelete = false;
  }

  async function executeDelete() {
    if (isBulkDelete) {
      await executeBulkDelete();
      closeDeleteModal();
    } else {
      if (examToDelete === null) return;

      try {
        const res = await fetch(`/api/exams/${examToDelete}`, {
          method: "DELETE",
          credentials: "include",
        });

        if (!res.ok) throw new Error("Gagal menghapus data");
        await fetchExams();
        closeDeleteModal();
      } catch (err) {
        alert(err instanceof Error ? err.message : String(err));
        closeDeleteModal();
      }
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

  // Bulk Delete State
  let selectedExams = $state<Set<number>>(new Set());
  let allSelectedOnPage = $derived(
    exams.length > 0 && exams.every(e => selectedExams.has(e.id))
  );

  function toggleAllOnPage() {
    const newSet = new Set(selectedExams);
    if (allSelectedOnPage) {
      exams.forEach(e => newSet.delete(e.id));
    } else {
      exams.forEach(e => newSet.add(e.id));
    }
    selectedExams = newSet;
  }

  function toggleExamSelection(id: number) {
    const newSet = new Set(selectedExams);
    if (newSet.has(id)) {
      newSet.delete(id);
    } else {
      newSet.add(id);
    }
    selectedExams = newSet;
  }

  async function executeBulkDelete() {
    if (selectedExams.size === 0) return;

    try {
      const res = await fetch(`/api/exams/bulk-delete`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ids: Array.from(selectedExams) }),
        credentials: "include",
      });

      if (!res.ok) throw new Error("Gagal menghapus data secara massal");
      selectedExams = new Set();
      await fetchExams();
    } catch (err) {
      alert(err instanceof Error ? err.message : String(err));
    }
  }
</script>

<svelte:head>
  <title>Nilai Harian Ujian - Portal {isTeacher ? "Guru" : "Siswa"}</title>
</svelte:head>

<div
  class="animate-in fade-in duration-500 max-w-7xl mx-auto space-y-6 p-4 sm:p-0"
>
  <div
    class="flex flex-col sm:flex-row items-start sm:items-end justify-between gap-4 mb-8"
  >
    <div>
      <h1
        class="text-2xl font-bold text-slate-900 sm:text-3xl tracking-tight drop-shadow-sm"
      >
        Nilai Harian
      </h1>
      <p class="mt-2 text-slate-600 text-sm sm:text-base tracking-wide">
        Kelola nilai ujian harian siswa.
      </p>
    </div>
    {#if isTeacher}
      <button
        onclick={openAddModal}
        class="inline-flex items-center gap-2 px-4 py-2.5 text-sm font-bold text-indigo-700 bg-indigo-50 hover:bg-indigo-100 border border-indigo-200 rounded-xl transition-all shadow-sm cursor-pointer"
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
    {#if exams.length === 0}
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
        <p class="text-slate-500 font-medium text-lg">Belum ada data ujian.</p>
        {#if isTeacher}
          <p class="text-slate-400 text-sm mt-1">
            Silakan tambahkan data baru melalui tombol di atas.
          </p>
        {/if}
      </div>
    {:else}
      {#if isTeacher}
        <div class="flex items-center gap-2 mb-6">
          <button
            onclick={() => switchTab("card")}
            class="px-4 py-2 text-sm font-medium rounded-xl transition-all {activeTab ===
            'card'
              ? 'bg-indigo-600 text-white shadow-md'
              : 'bg-white text-slate-600 hover:bg-slate-50 border border-slate-200 cursor-pointer'}"
          >
            Tampilan Kartu
          </button>
          <button
            onclick={() => switchTab("table")}
            class="px-4 py-2 text-sm font-medium rounded-xl transition-all {activeTab ===
            'table'
              ? 'bg-indigo-600 text-white shadow-md'
              : 'bg-white text-slate-600 hover:bg-slate-50 border border-slate-200 cursor-pointer'}"
          >
            Tampilan Tabel
          </button>
        </div>
      {/if}

      {#if activeTab === "table" && isTeacher}
        {#if selectedExams.size > 0}
          <div class="mb-4 p-4 bg-indigo-50 border border-indigo-100 rounded-xl flex items-center justify-between">
            <span class="text-sm font-semibold text-indigo-800">{selectedExams.size} data terpilih</span>
            <button
              onclick={openBulkDeleteModal}
              class="px-4 py-2 bg-red-600 text-white text-sm font-bold rounded-lg hover:bg-red-700 transition-colors shadow-sm cursor-pointer inline-flex items-center gap-2"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
              Hapus Terpilih
            </button>
          </div>
        {/if}
        <div
          class="bg-white/80 backdrop-blur-md rounded-2xl shadow-sm border border-slate-200 overflow-hidden"
        >
          <div class="overflow-x-auto">
            <table class="w-full text-left border-collapse min-w-[600px]">
              <thead>
                <tr
                  class="bg-slate-50/50 border-b border-slate-200 text-slate-500 text-sm"
                >
                  {#if isTeacher}
                    <th class="px-6 py-4 w-12 text-center">
                      <input
                        type="checkbox"
                        class="w-4 h-4 rounded border-slate-300 text-indigo-600 focus:ring-indigo-500 cursor-pointer"
                        checked={allSelectedOnPage}
                        onchange={toggleAllOnPage}
                        aria-label="Pilih Semua"
                      />
                    </th>
                    <th class="px-6 py-4 font-semibold">Murid</th>
                  {/if}
                  <th class="px-6 py-4 font-semibold">Tanggal</th>
                  <th class="px-6 py-4 font-semibold">Ujian</th>
                  <th class="px-6 py-4 font-semibold">Pelajaran</th>
                  <th class="px-6 py-4 font-semibold text-center">Nilai</th>
                  <th class="px-6 py-4 font-semibold text-center w-32"
                    >Detail</th
                  >
                  {#if isTeacher}
                    <th class="px-6 py-4 font-semibold text-center w-32"
                      >Aksi</th
                    >
                  {/if}
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100">
                {#each exams as exam}
                  <tr class="hover:bg-slate-50/50 transition-colors group">
                    {#if isTeacher}
                      <td class="px-6 py-4 text-center">
                        <input
                          type="checkbox"
                          class="w-4 h-4 rounded border-slate-300 text-indigo-600 focus:ring-indigo-500 cursor-pointer"
                          checked={selectedExams.has(exam.id)}
                          onchange={() => toggleExamSelection(exam.id)}
                          aria-label="Pilih Data"
                        />
                      </td>
                      <td class="px-6 py-4">
                        <span class="font-bold text-slate-800"
                          >{exam.user?.username || "Unknown"}</span
                        >
                      </td>
                    {/if}
                    <td class="px-6 py-4">
                      <span class="text-sm text-slate-600 font-medium"
                        >{formatDate(exam.exam_date)}</span
                      >
                    </td>
                    <td class="px-6 py-4 text-sm text-slate-800"
                      >{exam.exam_name}</td
                    >
                    <td class="px-6 py-4 text-sm text-slate-800"
                      >{exam.subject?.name || "Unknown"}</td
                    >
                    <td class="px-6 py-4 text-center">
                      <span
                        class={`px-3.5 py-1.5 rounded-full font-bold text-sm border ${exam.score >= 80 ? "bg-green-100 text-green-700 border-green-500" : exam.score >= 60 ? "bg-yellow-100 text-yellow-700 border-yellow-800/50" : "bg-red-100 text-red-600 border-red-300"}`}
                      >
                        {exam.score}
                      </span>
                    </td>
                    <td class="px-6 py-4 text-center">
                      <a
                        href={`/dashboard/exams/${exam.id}`}
                        title="Lihat Detail"
                        class="inline-flex items-center justify-center p-1.5 text-indigo-600 bg-indigo-100 hover:bg-indigo-200 rounded-lg transition-colors border border-indigo-300 cursor-pointer"
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
                    </td>
                    {#if isTeacher}
                      <td class="px-6 py-4">
                        <div class="flex items-center justify-center gap-2">
                          <button
                            onclick={() => openEditModal(exam)}
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
                            onclick={() => openDeleteModal(exam.id)}
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
                        </div>
                      </td>
                    {/if}
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        </div>
      {:else}
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {#each exams as exam}
            <a
              href={`/dashboard/exams/${exam.id}`}
              class="group bg-white/80 backdrop-blur-md border border-slate-200 rounded-2xl overflow-hidden shadow-sm hover:shadow-md transition-all flex flex-col cursor-pointer"
            >
              <div class="h-40 w-full bg-slate-100 overflow-hidden relative">
                {#if exam.image}
                  <img
                    src={exam.image}
                    alt="Exam"
                    class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                  />
                {:else}
                  <div
                    class="w-full h-full flex items-center justify-center text-slate-400 bg-indigo-50/50"
                  >
                    <svg
                      class="w-12 h-12 text-indigo-200"
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
                {/if}
                <div class="absolute top-2 right-2">
                  <span
                    class={`px-3.5 py-1.5 rounded-full font-bold text-sm border shadow-sm ${exam.score >= 80 ? "bg-green-100 text-green-700 border-green-500" : exam.score >= 60 ? "bg-yellow-100 text-yellow-700 border-yellow-800/50" : "bg-red-100 text-red-600 border-red-300"}`}
                  >
                    {exam.score}
                  </span>
                </div>
              </div>
              <div class="p-4 flex-1 flex flex-col justify-between gap-3">
                <div>
                  <h3
                    class="font-semibold text-slate-900 line-clamp-2 min-h-12"
                  >
                    {exam.exam_name}
                  </h3>
                  {#if isTeacher}
                    <p class="text-sm font-medium text-slate-600 mt-1">
                      {exam.user?.username || "Unknown"}
                    </p>
                  {/if}
                  <p class="text-sm text-slate-500 mt-1 flex justify-between">
                    <span>{exam.subject?.name}</span>
                    <span>{formatDate(exam.exam_date)}</span>
                  </p>
                </div>
                <div class="pt-3 border-t border-slate-100 flex justify-end">
                  <span
                    class="text-xs font-semibold text-indigo-600 group-hover:text-indigo-700 bg-indigo-50 px-2.5 py-1 rounded-lg transition-colors"
                    >Lihat Detail &rarr;</span
                  >
                </div>
              </div>
            </a>
          {/each}
        </div>
        {#if activeTab === "card" && currentPage < totalPages}
          <div use:infiniteScroll class="flex justify-center p-6">
            <div
              class="w-8 h-8 border-4 border-slate-200 border-t-indigo-500 rounded-full animate-spin"
            ></div>
          </div>
        {/if}
      {/if}
    {/if}

    <!-- Pagination Controls (Table Only) -->
    {#if activeTab === "table" && totalPages > 1}
      <div
        class="mt-8 px-4 py-3 sm:px-6 sm:py-4 bg-white/60 backdrop-blur-md rounded-2xl border border-slate-200 shadow-sm flex flex-row items-center justify-between gap-4"
      >
        <div class="text-sm text-slate-600 text-left">
          <span class="font-medium text-slate-900"
            >{Math.min(currentPage * itemsPerPage, totalRecords)}</span
          >
          dari
          <span class="font-medium text-slate-900">{totalRecords}</span> data
        </div>
        <div class="flex gap-2">
          <button
            aria-label="Sebelumnya"
            onclick={prevPage}
            disabled={currentPage === 1}
            class="p-2 text-sm font-medium rounded-lg transition-colors cursor-pointer flex items-center justify-center {currentPage ===
            1
              ? 'bg-white/40 text-zinc-600 cursor-not-allowed'
              : 'bg-white text-slate-800 hover:bg-slate-50 hover:text-slate-900'}"
          >
            <svg
              class="w-5 h-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              ><path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M15 19l-7-7 7-7"
              /></svg
            >
          </button>
          <div
            class="hidden sm:flex items-center gap-1 px-2 flex-wrap justify-center"
          >
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
            aria-label="Selanjutnya"
            onclick={nextPage}
            disabled={currentPage === totalPages}
            class="p-2 text-sm font-medium rounded-lg transition-colors cursor-pointer flex items-center justify-center {currentPage ===
            totalPages
              ? 'bg-white/40 text-zinc-600 cursor-not-allowed'
              : 'bg-white text-slate-800 hover:bg-slate-50 hover:text-slate-900'}"
          >
            <svg
              class="w-5 h-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              ><path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M9 5l7 7-7 7"
              /></svg
            >
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
    ></div>
    <div
      class="relative bg-slate-50 border border-slate-300 rounded-2xl shadow-2xl w-[95%] sm:w-full max-w-md overflow-hidden animate-in zoom-in-95 duration-200 flex flex-col max-h-[90vh]"
    >
      <div
        class="p-4 sm:p-6 border-b border-slate-200 flex justify-between items-center bg-white"
      >
        <h3 class="text-xl font-bold text-slate-900">
          {isEditing ? "Edit Nilai Ujian" : "Tambah Nilai Ujian"}
        </h3>
        <button
          onclick={closeModal}
          class="text-slate-500 hover:text-slate-900 transition-colors cursor-pointer"
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

      <form onsubmit={saveExam} class="p-4 sm:p-6 space-y-4 overflow-y-auto">
        <div>
          <label
            class="block text-sm font-medium text-slate-600 mb-1.5"
            for="user">Pilih Murid</label
          >
          <select
            id="user"
            bind:value={formUserId}
            class="w-full bg-slate-50 hover:bg-white border border-slate-200 focus:bg-white focus:border-indigo-400 focus:ring-4 focus:ring-indigo-400/20 rounded-2xl pl-4 pr-10 py-3 text-sm text-slate-700 outline-none transition-all shadow-sm appearance-none bg-[url('data:image/svg+xml;charset=utf-8,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20fill%3D%22none%22%20viewBox%3D%220%200%2024%2024%22%20stroke%3D%22%2364748b%22%3E%3Cpath%20stroke-linecap%3D%22round%22%20stroke-linejoin%3D%22round%22%20stroke-width%3D%222%22%20d%3D%22M19%209l-7%207-7-7%22%2F%3E%3C%2Fsvg%3E')] bg-[length:1.25rem_1.25rem] bg-[position:right_1rem_center] bg-no-repeat [color-scheme:light]"
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
            class="block text-sm font-medium text-slate-600 mb-1.5"
            for="date">Tanggal Ujian</label
          >
          <input
            type="date"
            id="date"
            bind:value={formExamDate}
            class="w-full bg-slate-50 hover:bg-white border border-slate-200 focus:bg-white focus:border-indigo-400 focus:ring-4 focus:ring-indigo-400/20 rounded-2xl px-4 py-3 text-sm text-slate-700 outline-none transition-all shadow-sm [color-scheme:light]"
            required
          />
        </div>

        <div>
          <label
            class="block text-sm font-medium text-slate-600 mb-1.5"
            for="examName">Nama Ujian</label
          >
          <input
            type="text"
            id="examName"
            bind:value={formExamName}
            placeholder="Misal: Ujian Tengah Semester"
            class="w-full bg-slate-50 hover:bg-white border border-slate-200 focus:bg-white focus:border-indigo-400 focus:ring-4 focus:ring-indigo-400/20 rounded-2xl px-4 py-3 text-sm text-slate-700 outline-none transition-all shadow-sm"
            required
          />
        </div>

        <div>
          <label
            class="block text-sm font-medium text-slate-600 mb-1.5"
            for="subject">Mata Pelajaran</label
          >
          <select
            id="subject"
            bind:value={formSubjectId}
            class="w-full bg-slate-50 hover:bg-white border border-slate-200 focus:bg-white focus:border-indigo-400 focus:ring-4 focus:ring-indigo-400/20 rounded-2xl pl-4 pr-10 py-3 text-sm text-slate-700 outline-none transition-all shadow-sm appearance-none bg-[url('data:image/svg+xml;charset=utf-8,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20fill%3D%22none%22%20viewBox%3D%220%200%2024%2024%22%20stroke%3D%22%2364748b%22%3E%3Cpath%20stroke-linecap%3D%22round%22%20stroke-linejoin%3D%22round%22%20stroke-width%3D%222%22%20d%3D%22M19%209l-7%207-7-7%22%2F%3E%3C%2Fsvg%3E')] bg-[length:1.25rem_1.25rem] bg-[position:right_1rem_center] bg-no-repeat [color-scheme:light]"
            required
          >
            <option value="" disabled hidden>-- Pilih Mata Pelajaran --</option>
            {#each subjects as s}
              <option value={s.id}>{s.name}</option>
            {/each}
          </select>
        </div>

        <div>
          <label
            class="block text-sm font-medium text-slate-600 mb-1.5"
            for="score">Nilai (0 - 100)</label
          >
          <input
            type="number"
            id="score"
            value={formScore}
            oninput={(e) => {
              let val = parseInt(e.currentTarget.value, 10);
              if (isNaN(val)) val = 0;
              if (val > 100) val = 100;
              if (val < 0) val = 0;
              formScore = val;
              e.currentTarget.value = String(val);
            }}
            min="0"
            max="100"
            class="w-full bg-slate-50 hover:bg-white border border-slate-200 focus:bg-white focus:border-indigo-400 focus:ring-4 focus:ring-indigo-400/20 rounded-2xl px-4 py-3 text-sm text-slate-700 outline-none transition-all shadow-sm"
            required
          />
        </div>

        <div>
          <label
            class="block text-sm font-medium text-slate-600 mb-1.5"
            for="image">Bukti Ujian (Opsional, Maks 1MB)</label
          >
          <input
            type="file"
            id="image"
            accept="image/*"
            onchange={handleImageUpload}
            disabled={isUploadingImage}
            class="w-full bg-slate-50 hover:bg-white border border-slate-200 focus:bg-white focus:border-indigo-400 focus:ring-4 focus:ring-indigo-400/20 rounded-2xl px-4 py-3 text-sm text-slate-700 outline-none transition-all shadow-sm file:mr-4 file:py-2 file:px-4 file:rounded-xl file:border-0 file:text-sm file:font-semibold file:bg-indigo-50 file:text-indigo-700 hover:file:bg-indigo-100"
          />
          {#if isUploadingImage}
            <p
              class="text-xs text-indigo-600 mt-2 flex items-center gap-2 animate-pulse"
            >
              Mengunggah...
            </p>
          {/if}
          {#if formImage}
            <div class="mt-3 relative inline-block">
              <img
                src={formImage}
                alt="Bukti Ujian"
                class="h-24 object-cover rounded-xl border border-slate-200 shadow-sm"
              />
              <button
                type="button"
                onclick={() => (formImage = "")}
                class="absolute -top-2 -right-2 bg-red-100 text-red-600 hover:bg-red-200 rounded-full p-1 border border-red-300 shadow-sm transition-colors cursor-pointer"
                title="Hapus Bukti"
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
                    d="M6 18L18 6M6 6l12 12"
                  ></path></svg
                >
              </button>
            </div>
          {/if}
        </div>

        <div class="pt-4 flex justify-end gap-3 border-t border-slate-200 mt-4">
          <button
            type="button"
            onclick={closeModal}
            class="px-4 py-2.5 text-sm font-medium text-slate-800 bg-white border border-slate-200 hover:bg-slate-100 shadow-md rounded-xl transition-colors cursor-pointer"
          >
            Batal
          </button>
          <button
            type="submit"
            class="px-4 py-2.5 text-sm font-medium text-slate-900 bg-indigo-50 border border-indigo-200 hover:bg-indigo-100 rounded-xl transition-all shadow-md shadow-indigo-900/20 cursor-pointer"
          >
            {isEditing ? "Simpan Perubahan" : "Tambahkan"}
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
          ><path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
          /></svg
        >
      </div>
      <h3 class="text-xl font-bold text-slate-900 mb-2">Hapus Nilai Ujian?</h3>
      <p class="text-sm text-slate-600 mb-6">
        {#if isBulkDelete}
          Apakah Anda yakin ingin menghapus {selectedExams.size} data yang terpilih? Data yang dihapus tidak dapat dikembalikan.
        {:else}
          Apakah Anda yakin ingin menghapus nilai ini? Data yang dihapus tidak dapat dikembalikan.
        {/if}
      </p>
      <div class="flex justify-center gap-3">
        <button
          onclick={closeDeleteModal}
          class="px-4 py-2.5 text-sm font-medium text-slate-800 bg-white border border-slate-300 rounded-xl hover:bg-slate-50 transition-colors cursor-pointer"
          >Batal</button
        >
        <button
          onclick={executeDelete}
          class="px-4 py-2.5 text-sm font-medium text-white bg-red-600 hover:bg-red-700 rounded-xl transition-all shadow-md shadow-red-900/20 cursor-pointer"
          >Ya, Hapus</button
        >
      </div>
    </div>
  </div>
{/if}
