<script lang="ts">
  import { onMount } from "svelte";
  import Chart from "chart.js/auto";

  type Exam = {
    id: number;
    user_id: number;
    exam_date: string;
    exam_name: string;
    subject_id: number;
    score: number;
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
  let currentUserId = $state<number | null>(null);

  // Tabs & Chart State
  let activeTab = $state("table");
  let chartCanvas: HTMLCanvasElement | null = $state(null);
  let chartInstance: Chart | null = null;
  let chartSelectedUserId = $state("");
  let chartTimeframe = $state("harian");

  // Pagination & Filter State
  let currentPage = $state(1);
  const itemsPerPage = 20;

  let filterUser = $state("");
  let filterSubject = $state("");
  let filterDate = $state("");

  let filteredExams = $derived(
    exams.filter((exam) => {
      const matchUser =
        filterUser === "" ||
        (exam.user?.username || "")
          .toLowerCase()
          .includes(filterUser.toLowerCase());
      const matchSubject =
        filterSubject === "" ||
        (exam.subject?.name || "")
          .toLowerCase()
          .includes(filterSubject.toLowerCase());
      const matchDate =
        filterDate === "" || (exam.exam_date || "").startsWith(filterDate);
      return matchUser && matchSubject && matchDate;
    }),
  );

  let paginatedExams = $derived(
    filteredExams.slice(
      (currentPage - 1) * itemsPerPage,
      currentPage * itemsPerPage,
    ),
  );
  let totalPages = $derived(
    Math.ceil(filteredExams.length / itemsPerPage) || 1,
  );

  $effect(() => {
    // Reset page to 1 when filters change
    filterUser;
    filterSubject;
    filterDate;
    currentPage = 1;
  });

  function nextPage() {
    if (currentPage < totalPages) currentPage++;
  }

  function prevPage() {
    if (currentPage > 1) currentPage--;
  }

  function goToPage(page: number) {
    currentPage = page;
  }

  // Modal State
  let showModal = $state(false);
  let isEditing = $state(false);
  let currentExamId: number | null = $state(null);

  let showDeleteModal = $state(false);
  let examToDelete: number | null = $state(null);

  // Form State
  let formUserId: number | string = $state("");
  let formExamDate = $state("");
  let formExamName = $state("");
  let formSubjectId: number | string = $state("");
  let formScore = $state(0);

  async function fetchExams() {
    try {
      const res = await fetch(`/api/exams`, {
        credentials: "include",
      });
      if (!res.ok) throw new Error("Gagal mengambil data ujian");
      exams = (await res.json()) || [];
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
        currentUserId = data.user.id;
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
    if (!isTeacher && currentUserId) {
      chartSelectedUserId = currentUserId.toString();
    }
    isLoading = false;
  }

  onMount(() => {
    loadData();
  });

  function openAddModal() {
    isEditing = false;
    currentExamId = null;
    formUserId = "";
    formExamDate = new Date().toLocaleDateString("en-CA"); // "en-CA" formats as YYYY-MM-DD
    formExamName = "";
    formSubjectId = "";
    formScore = 0;
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
    examToDelete = id;
    showDeleteModal = true;
  }

  function closeDeleteModal() {
    showDeleteModal = false;
    examToDelete = null;
  }

  async function executeDelete() {
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

  function formatDate(dateStr: string | null) {
    if (!dateStr) return "-";
    return new Date(dateStr).toLocaleDateString("id-ID", {
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  }

  function getWeekNumber(d: Date) {
    const date = new Date(Date.UTC(d.getFullYear(), d.getMonth(), d.getDate()));
    const dayNum = date.getUTCDay() || 7;
    date.setUTCDate(date.getUTCDate() + 4 - dayNum);
    const yearStart = new Date(Date.UTC(date.getUTCFullYear(), 0, 1));
    return Math.ceil(
      ((date.getTime() - yearStart.getTime()) / 86400000 + 1) / 7,
    );
  }

  function getTimeframeLabel(dateStr: string, timeframe: string) {
    const d = new Date(dateStr);
    if (timeframe === "harian") {
      return formatDate(dateStr);
    } else if (timeframe === "mingguan") {
      const w = getWeekNumber(d);
      return `Minggu ${w}, ${d.getFullYear()}`;
    } else if (timeframe === "bulanan") {
      return d.toLocaleDateString("id-ID", { month: "long", year: "numeric" });
    }
    return formatDate(dateStr);
  }

  function drawChart() {
    if (activeTab !== "chart" || !chartCanvas) return;
    if (chartInstance) chartInstance.destroy();

    if (!chartSelectedUserId) return;

    // Filter exams for selected user
    const userExams = exams.filter(
      (e) => e.user_id === parseInt(chartSelectedUserId),
    );

    // Sort by date ascending (base sort)
    userExams.sort(
      (a, b) =>
        new Date(a.exam_date).getTime() - new Date(b.exam_date).getTime(),
    );

    // Group by subject.name
    const subjectsMap: Record<string, { x: string; y: number }[]> = {};
    const datesSet = new Set<string>();
    const labelSortValue: Record<string, number> = {};

    userExams.forEach((e) => {
      const subjectName = e.subject?.name || "Unknown";
      const d = new Date(e.exam_date);
      const label = getTimeframeLabel(e.exam_date, chartTimeframe);

      datesSet.add(label);
      if (!labelSortValue[label]) {
        labelSortValue[label] = d.getTime();
      }

      if (!subjectsMap[subjectName]) subjectsMap[subjectName] = [];
      subjectsMap[subjectName].push({ x: label, y: e.score });
    });

    const labels = Array.from(datesSet).sort(
      (a, b) => labelSortValue[a] - labelSortValue[b],
    );

    const datasets = [];
    const colors = [
      "#60a5fa",
      "#34d399",
      "#fbbf24",
      "#f87171",
      "#a78bfa",
      "#f472b6",
    ];
    let colorIndex = 0;

    for (const [subjectName, dataPoints] of Object.entries(subjectsMap)) {
      // Map dataPoints to labels and calculate average if multiple scores in the same timeframe
      const data = labels.map((label) => {
        const points = dataPoints.filter((dp) => dp.x === label);
        if (points.length === 0) return null;
        const sum = points.reduce((acc, curr) => acc + curr.y, 0);
        return sum / points.length;
      });

      datasets.push({
        label: subjectName,
        data: data,
        borderColor: colors[colorIndex % colors.length],
        backgroundColor: colors[colorIndex % colors.length],
        pointBackgroundColor: colors[colorIndex % colors.length],
        tension: 0.3,
        fill: false,
        spanGaps: true,
      });
      colorIndex++;
    }

    chartInstance = new Chart(chartCanvas, {
      type: "line",
      data: { labels, datasets },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: {
            labels: {
              color: "#431407",
              font: { weight: "bold" },
              usePointStyle: true,
              pointStyle: "circle",
              padding: 20,
            },
          },
        },
        scales: {
          y: {
            beginAtZero: true,
            max: 100,
            grid: { color: "#fed7aa" },
            ticks: { color: "#9a3412", font: { weight: 500 } },
          },
          x: {
            grid: { color: "#fed7aa" },
            ticks: { color: "#9a3412", font: { weight: 500 } },
          },
        },
      },
    });
  }

  $effect(() => {
    if (activeTab === "chart") {
      // Small delay to ensure canvas is mounted
      setTimeout(drawChart, 0);
    }
  });

  $effect(() => {
    // Redraw if selected user or timeframe changes
    chartSelectedUserId;
    chartTimeframe;
    setTimeout(drawChart, 0);
  });
</script>

<svelte:head>
  <title>Nilai Harian Ujian - Portal Guru</title>
</svelte:head>

<div class="animate-in fade-in duration-500 max-w-7xl mx-auto space-y-6">
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
        class="inline-flex items-center gap-2 px-4 py-2.5 bg-indigo-300 hover:bg-indigo-200 text-slate-900 font-medium rounded-xl transition-all shadow-md shadow-indigo-900/20"
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
        Nilai Ujian
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
    <!-- Tabs Header -->
    <div
      class="flex gap-4 sm:gap-6 border-b border-slate-200 pb-0 overflow-x-auto whitespace-nowrap px-1"
    >
      <button
        onclick={() => (activeTab = "table")}
        class="pb-3 px-1 font-medium transition-colors {activeTab === 'table'
          ? 'text-indigo-400 border-b-2 border-indigo-400'
          : 'text-slate-500 hover:text-slate-800'}">Tabel Nilai</button
      >
      <button
        onclick={() => (activeTab = "chart")}
        class="pb-3 px-1 font-medium transition-colors {activeTab === 'chart'
          ? 'text-indigo-400 border-b-2 border-indigo-400'
          : 'text-slate-500 hover:text-slate-800'}">Grafik Perkembangan</button
      >
    </div>

    {#if activeTab === "table"}
      <div
        class="bg-white/60 backdrop-blur-md rounded-3xl border border-slate-200 shadow-lg shadow-slate-800/10 overflow-hidden"
      >
        <div class="overflow-x-auto w-full">
          <table
            class="w-full text-left border-collapse whitespace-nowrap min-w-[800px]"
          >
            <thead>
              <tr class="bg-white/40 border-b border-slate-200">
                <th class="py-4 px-6 align-bottom min-w-[200px]">
                  <div class="font-bold text-slate-900 text-sm mb-2">
                    Tanggal
                  </div>
                  <div class="relative w-full">
                    <input
                      type="date"
                      bind:value={filterDate}
                      class="w-full bg-white border border-slate-300 rounded-lg pl-3 pr-8 py-2 text-sm text-slate-900 focus:ring-1 focus:ring-indigo-500 outline-none transition-all scheme-light"
                    />
                    {#if filterDate}
                      <button
                        onclick={() => (filterDate = "")}
                        class="absolute right-2 top-1/2 -translate-y-1/2 text-slate-600 hover:text-red-600 transition-colors"
                        title="Hapus Filter Tanggal"
                      >
                        <svg
                          class="w-3.5 h-3.5"
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
                    {/if}
                  </div>
                </th>
                {#if isTeacher}
                  <th class="py-4 px-6 align-bottom min-w-[200px]">
                    <div class="font-bold text-slate-900 text-sm mb-2">
                      Murid
                    </div>
                    <select
                      bind:value={filterUser}
                      class="w-full bg-white border border-slate-300 rounded-lg pl-3 pr-8 py-2 text-sm text-slate-900 focus:ring-1 focus:ring-indigo-500 outline-none transition-all appearance-none bg-[url('data:image/svg+xml;charset=utf-8,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20fill%3D%22none%22%20viewBox%3D%220%200%2024%2024%22%20stroke%3D%22%2364748b%22%3E%3Cpath%20stroke-linecap%3D%22round%22%20stroke-linejoin%3D%22round%22%20stroke-width%3D%222%22%20d%3D%22M19%209l-7%207-7-7%22%2F%3E%3C%2Fsvg%3E')] bg-[length:1.25rem_1.25rem] bg-[position:right_0.5rem_center] bg-no-repeat scheme-light"
                    >
                      <option value="" class="bg-white text-slate-900">-- Semua Murid --</option>
                      {#each users as u}
                        <option value={u.username} class="bg-white text-slate-900">{u.username}</option>
                      {/each}
                    </select>
                  </th>
                {/if}
                <th
                  class="py-4 px-6 align-bottom font-bold text-slate-900 text-sm pb-5"
                  >Nama Ujian</th
                >
                <th class="py-4 px-6 align-bottom min-w-[200px]">
                  <div class="font-bold text-slate-900 text-sm mb-2">
                    Mata Pelajaran
                  </div>
                  <select
                    bind:value={filterSubject}
                    class="w-full bg-white border border-slate-300 rounded-lg pl-3 pr-8 py-2 text-sm text-slate-900 focus:ring-1 focus:ring-indigo-500 outline-none transition-all appearance-none bg-[url('data:image/svg+xml;charset=utf-8,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20fill%3D%22none%22%20viewBox%3D%220%200%2024%2024%22%20stroke%3D%22%2364748b%22%3E%3Cpath%20stroke-linecap%3D%22round%22%20stroke-linejoin%3D%22round%22%20stroke-width%3D%222%22%20d%3D%22M19%209l-7%207-7-7%22%2F%3E%3C%2Fsvg%3E')] bg-[length:1.25rem_1.25rem] bg-[position:right_0.5rem_center] bg-no-repeat scheme-light"
                  >
                    <option value="" class="bg-white text-slate-900">-- Semua Pelajaran --</option>
                    {#each subjects as s}
                      <option value={s.name} class="bg-white text-slate-900">{s.name}</option>
                    {/each}
                  </select>
                </th>
                <th
                  class="py-4 px-6 align-bottom font-bold text-slate-900 text-sm text-center pb-5"
                  >Nilai</th
                >
                {#if isTeacher}
                  <th
                    class="py-4 px-6 align-bottom font-bold text-slate-900 text-sm text-center pb-5"
                    >Aksi</th
                  >
                {/if}
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-200">
              {#each paginatedExams as exam}
                <tr class="hover:bg-white/40 transition-colors">
                  <td class="py-4 px-6 text-sm text-slate-600"
                    >{formatDate(exam.exam_date)}</td
                  >
                  {#if isTeacher}
                    <td
                      class="py-4 px-6 text-sm font-medium text-slate-900 drop-shadow-sm"
                      >{exam.user?.username || "Unknown"}</td
                    >
                  {/if}
                  <td class="py-4 px-6 text-sm text-slate-800"
                    >{exam.exam_name}</td
                  >
                  <td class="py-4 px-6 text-sm text-slate-800"
                    >{exam.subject?.name || "Unknown"}</td
                  >
                  <td class="py-4 px-6 text-sm text-center">
                    <span
                      class={`px-3 py-1 rounded-full font-bold text-xs border ${exam.score >= 80 ? "bg-green-100 text-green-700 border-green-500" : exam.score >= 60 ? "bg-yellow-100 text-yellow-700 border-yellow-800/50" : "bg-red-100 text-red-600 border-red-300"}`}
                    >
                      {exam.score}
                    </span>
                  </td>
                  {#if isTeacher}
                    <td class="py-4 px-6 text-center space-x-2">
                      <button
                        onclick={() => openEditModal(exam)}
                        class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-blue-600 bg-blue-100 hover:bg-blue-800/50 rounded-lg transition-colors border border-blue-300"
                      >
                        Edit
                      </button>
                      <button
                        onclick={() => openDeleteModal(exam.id)}
                        class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-red-600 bg-red-100 hover:bg-red-200 rounded-lg transition-colors border border-red-300"
                      >
                        Hapus
                      </button>
                    </td>
                  {/if}
                </tr>
              {/each}
              {#if filteredExams.length === 0}
                <tr>
                  <td
                    colspan={isTeacher ? 6 : 4}
                    class="py-12 text-center text-slate-500 font-light"
                    >Belum ada data nilai ujian.</td
                  >
                </tr>
              {/if}
            </tbody>
          </table>
        </div>

        <!-- Pagination Controls -->
        {#if totalPages > 1}
          <div
            class="px-4 sm:px-6 py-4 border-t border-slate-200 flex flex-col sm:flex-row items-center justify-between gap-4"
          >
            <div class="text-sm text-slate-600 text-center sm:text-left">
              Menampilkan <span class="font-medium text-slate-900"
                >{filteredExams.length > 0
                  ? (currentPage - 1) * itemsPerPage + 1
                  : 0}</span
              >
              sampai
              <span class="font-medium text-slate-900"
                >{Math.min(
                  currentPage * itemsPerPage,
                  filteredExams.length,
                )}</span
              >
              dari
              <span class="font-medium text-slate-900"
                >{filteredExams.length}</span
              > data
            </div>
            <div class="flex gap-2">
              <button
                onclick={prevPage}
                disabled={currentPage === 1}
                class="px-3 py-1.5 text-sm font-medium rounded-lg transition-colors {currentPage ===
                1
                  ? 'bg-white/40 text-zinc-600 cursor-not-allowed'
                  : 'bg-white text-slate-800 hover:bg-slate-50 hover:text-slate-900'}"
              >
                Sebelumnya
              </button>
              <div
                class="flex items-center gap-1 px-2 flex-wrap justify-center"
              >
                {#each Array(totalPages) as _, i}
                  <button
                    onclick={() => goToPage(i + 1)}
                    class="w-8 h-8 flex items-center justify-center text-sm font-medium rounded-lg transition-colors {currentPage ===
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
                class="px-3 py-1.5 text-sm font-medium rounded-lg transition-colors {currentPage ===
                totalPages
                  ? 'bg-white/40 text-zinc-600 cursor-not-allowed'
                  : 'bg-white text-slate-800 hover:bg-slate-50 hover:text-slate-900'}"
              >
                Selanjutnya
              </button>
            </div>
          </div>
        {/if}
      </div>
    {:else if activeTab === "chart"}
      <div
        class="bg-white/60 backdrop-blur-md rounded-3xl border border-slate-200 shadow-lg shadow-slate-800/10 p-6 space-y-6"
      >
        <div class="flex flex-col sm:flex-row gap-4">
          {#if isTeacher}
            <div class="w-full sm:max-w-xs">
              <label
                class="block text-sm font-medium text-slate-600 mb-2"
                for="chartUser">Pilih Murid</label
              >
              <select
                id="chartUser"
                bind:value={chartSelectedUserId}
                class="w-full bg-slate-50 hover:bg-white border border-slate-200 focus:bg-white focus:border-indigo-400 focus:ring-4 focus:ring-indigo-400/20 rounded-2xl pl-4 pr-10 py-3 text-sm text-slate-700 outline-none transition-all shadow-sm cursor-pointer appearance-none bg-[url('data:image/svg+xml;charset=utf-8,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20fill%3D%22none%22%20viewBox%3D%220%200%2024%2024%22%20stroke%3D%22%2364748b%22%3E%3Cpath%20stroke-linecap%3D%22round%22%20stroke-linejoin%3D%22round%22%20stroke-width%3D%222%22%20d%3D%22M19%209l-7%207-7-7%22%2F%3E%3C%2Fsvg%3E')] bg-[length:1.25rem_1.25rem] bg-[position:right_1rem_center] bg-no-repeat"
              >
                <option value="" disabled hidden>-- Pilih Murid --</option>
                {#each users as u}
                  <option value={u.id} class="bg-slate-700 text-white"
                    >{u.username}</option
                  >
                {/each}
              </select>
            </div>
          {/if}

          {#if chartSelectedUserId}
            <div class="w-full sm:max-w-[12rem]">
              <label
                class="block text-sm font-medium text-slate-600 mb-2"
                for="chartTimeframe">Rentang Waktu</label
              >
              <select
                id="chartTimeframe"
                bind:value={chartTimeframe}
                class="w-full bg-slate-50 hover:bg-white border border-slate-200 focus:bg-white focus:border-indigo-400 focus:ring-4 focus:ring-indigo-400/20 rounded-2xl pl-4 pr-10 py-3 text-sm text-slate-700 outline-none transition-all shadow-sm cursor-pointer appearance-none bg-[url('data:image/svg+xml;charset=utf-8,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20fill%3D%22none%22%20viewBox%3D%220%200%2024%2024%22%20stroke%3D%22%2364748b%22%3E%3Cpath%20stroke-linecap%3D%22round%22%20stroke-linejoin%3D%22round%22%20stroke-width%3D%222%22%20d%3D%22M19%209l-7%207-7-7%22%2F%3E%3C%2Fsvg%3E')] bg-[length:1.25rem_1.25rem] bg-[position:right_1rem_center] bg-no-repeat"
              >
                <option value="harian" class="bg-slate-700 text-white"
                  >Harian</option
                >
                <option value="mingguan" class="bg-slate-700 text-white"
                  >Mingguan</option
                >
                <option value="bulanan" class="bg-slate-700 text-white"
                  >Bulanan</option
                >
              </select>
            </div>
          {/if}
        </div>

        {#if chartSelectedUserId}
          <div class="h-96 w-full">
            <canvas bind:this={chartCanvas}></canvas>
          </div>
        {:else}
          <div
            class="py-16 text-center text-slate-500 font-light border-2 border-dashed border-slate-200 rounded-2xl"
          >
            Silakan pilih nama murid terlebih dahulu untuk melihat grafik
            perkembangan.
          </div>
        {/if}
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
        class="p-4 sm:p-6 border-b border-slate-200 flex justify-between items-center"
      >
        <h3 class="text-xl font-bold text-slate-900">
          {isEditing ? "Edit Nilai Ujian" : "Tambah Nilai Ujian"}
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
              <option value={u.id} class="bg-slate-700 text-white"
                >{u.username}</option
              >
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
              <option value={s.id} class="bg-slate-700 text-white"
                >{s.name}</option
              >
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

        <div class="pt-4 flex justify-end gap-3">
          <button
            type="button"
            onclick={closeModal}
            class="px-4 py-2.5 text-sm font-medium text-slate-800 bg-white border border-slate-200 hover:bg-slate-100 shadow-md rounded-xl transition-colors"
          >
            Batal
          </button>
          <button
            type="submit"
            class="px-4 py-2.5 text-sm font-medium text-slate-900 bg-indigo-50 border border-indigo-200 hover:bg-indigo-100 rounded-xl transition-all shadow-md shadow-indigo-900/20"
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
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
          />
        </svg>
      </div>
      <h3 class="text-xl font-bold text-slate-900 mb-2">Hapus Nilai Ujian?</h3>
      <p class="text-sm text-slate-600 mb-6">
        Apakah Anda yakin ingin menghapus nilai ini? Data yang dihapus tidak
        dapat dikembalikan.
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

