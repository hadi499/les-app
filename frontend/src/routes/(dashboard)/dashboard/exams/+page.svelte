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

  // Form State
  let formUserId: number | string = $state("");
  let formExamDate = $state("");
  let formExamName = $state("");
  let formSubjectId: number | string = $state("");
  let formScore = $state(0);

  async function fetchExams() {
    try {
      const res = await fetch("http://localhost:8080/api/exams", {
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
      const res = await fetch("http://localhost:8080/api/users", {
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
      const res = await fetch("http://localhost:8080/api/subjects", {
        credentials: "include",
      });
      if (!res.ok) throw new Error("Gagal mengambil data mapel");
      subjects = (await res.json()) || [];
    } catch (e) {
      console.error(e);
    }
  }

  async function loadData() {
    isLoading = true;
    errorMsg = "";
    await Promise.all([fetchExams(), fetchUsers(), fetchSubjects()]);
    isLoading = false;
  }

  onMount(() => {
    loadData();
  });

  function openAddModal() {
    isEditing = false;
    currentExamId = null;
    formUserId = "";
    formExamDate = new Date().toISOString().split("T")[0];
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
      let url = "http://localhost:8080/api/exams";
      let method = "POST";

      if (isEditing) {
        url = `http://localhost:8080/api/exams/${currentExamId}`;
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

  async function deleteExam(id: number) {
    if (!confirm("Apakah Anda yakin ingin menghapus nilai ini?")) return;

    try {
      const res = await fetch(`http://localhost:8080/api/exams/${id}`, {
        method: "DELETE",
        credentials: "include",
      });

      if (!res.ok) throw new Error("Gagal menghapus data");
      await fetchExams();
    } catch (err) {
      alert(err instanceof Error ? err.message : String(err));
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
    class="flex flex-col sm:flex-row sm:items-end justify-between gap-4 mb-8"
  >
    <div>
      <h1
        class="text-2xl font-bold text-orange-950 sm:text-3xl tracking-tight drop-shadow-sm"
      >
        Nilai Harian
      </h1>
      <p class="mt-2 text-orange-800 text-sm sm:text-base tracking-wide">
        Kelola nilai ujian harian siswa.
      </p>
    </div>
    <button
      onclick={openAddModal}
      class="inline-flex items-center gap-2 px-4 py-2.5 bg-indigo-300 hover:bg-indigo-200 text-orange-950 font-medium rounded-xl transition-all shadow-md shadow-indigo-900/20"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"
        ><path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M12 4v16m8-8H4"
        ></path></svg
      >
      Nilai Ujian
    </button>
  </div>

  {#if isLoading}
    <div class="flex justify-center p-12">
      <div
        class="w-10 h-10 border-4 border-orange-200 border-t-orange-500 rounded-full animate-spin"
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
    <div class="flex gap-6 border-b border-orange-200 pb-0">
      <button
        onclick={() => (activeTab = "table")}
        class="pb-3 px-1 font-medium transition-colors {activeTab === 'table'
          ? 'text-indigo-400 border-b-2 border-indigo-400'
          : 'text-orange-700 hover:text-orange-900'}">Tabel Nilai</button
      >
      <button
        onclick={() => (activeTab = "chart")}
        class="pb-3 px-1 font-medium transition-colors {activeTab === 'chart'
          ? 'text-indigo-400 border-b-2 border-indigo-400'
          : 'text-orange-700 hover:text-orange-900'}"
        >Grafik Perkembangan</button
      >
    </div>

    {#if activeTab === "table"}
      <div
        class="bg-white/60 backdrop-blur-md rounded-3xl border border-orange-200 shadow-lg shadow-orange-900/10 overflow-hidden"
      >
        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="bg-white/40 border-b border-orange-200">
                <th class="py-4 px-6 align-bottom">
                  <div class="font-bold text-orange-950 text-sm mb-2">
                    Tanggal
                  </div>
                  <div class="relative w-full">
                    <input
                      type="date"
                      bind:value={filterDate}
                      class="w-full bg-white/80 border border-orange-300 rounded-lg pl-2 pr-8 py-1 text-xs text-orange-950 focus:ring-1 focus:ring-indigo-500 outline-none transition-all"
                    />
                    {#if filterDate}
                      <button
                        onclick={() => (filterDate = "")}
                        class="absolute right-2 top-1/2 -translate-y-1/2 text-orange-800 hover:text-red-600 transition-colors"
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
                <th class="py-4 px-6 align-bottom">
                  <div class="font-bold text-orange-950 text-sm mb-2">
                    Murid
                  </div>
                  <select
                    bind:value={filterUser}
                    class="w-full bg-white/80 border border-orange-300 rounded-lg px-2 py-1 text-xs text-orange-950 focus:ring-1 focus:ring-indigo-500 outline-none transition-all"
                  >
                    <option value="">-- Semua Murid --</option>
                    {#each users as u}
                      <option value={u.username}>{u.username}</option>
                    {/each}
                  </select>
                </th>
                <th
                  class="py-4 px-6 align-bottom font-bold text-orange-950 text-sm pb-5"
                  >Nama Ujian</th
                >
                <th class="py-4 px-6 align-bottom">
                  <div class="font-bold text-orange-950 text-sm mb-2">
                    Mata Pelajaran
                  </div>
                  <select
                    bind:value={filterSubject}
                    class="w-full bg-white/80 border border-orange-300 rounded-lg px-2 py-1 text-xs text-orange-950 focus:ring-1 focus:ring-indigo-500 outline-none transition-all"
                  >
                    <option value="">-- Semua Pelajaran --</option>
                    {#each subjects as s}
                      <option value={s.name}>{s.name}</option>
                    {/each}
                  </select>
                </th>
                <th
                  class="py-4 px-6 align-bottom font-bold text-orange-950 text-sm text-center pb-5"
                  >Nilai</th
                >
                <th
                  class="py-4 px-6 align-bottom font-bold text-orange-950 text-sm text-right pb-5"
                  >Aksi</th
                >
              </tr>
            </thead>
            <tbody class="divide-y divide-orange-200">
              {#each paginatedExams as exam}
                <tr class="hover:bg-white/40 transition-colors">
                  <td class="py-4 px-6 text-sm text-orange-800"
                    >{formatDate(exam.exam_date)}</td
                  >
                  <td
                    class="py-4 px-6 text-sm font-medium text-orange-950 drop-shadow-sm"
                    >{exam.user?.username || "Unknown"}</td
                  >
                  <td class="py-4 px-6 text-sm text-orange-900"
                    >{exam.exam_name}</td
                  >
                  <td class="py-4 px-6 text-sm text-orange-900"
                    >{exam.subject?.name || "Unknown"}</td
                  >
                  <td class="py-4 px-6 text-sm text-center">
                    <span
                      class={`px-3 py-1 rounded-full font-bold text-xs border ${exam.score >= 80 ? "bg-green-100 text-green-700 border-green-500" : exam.score >= 60 ? "bg-yellow-100 text-yellow-700 border-yellow-800/50" : "bg-red-100 text-red-600 border-red-300"}`}
                    >
                      {exam.score}
                    </span>
                  </td>
                  <td class="py-4 px-6 text-right space-x-2">
                    <button
                      onclick={() => openEditModal(exam)}
                      class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-orange-600 bg-blue-100 hover:bg-blue-800/50 rounded-lg transition-colors border border-blue-300"
                    >
                      Edit
                    </button>
                    <button
                      onclick={() => deleteExam(exam.id)}
                      class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-red-600 bg-red-100 hover:bg-red-200 rounded-lg transition-colors border border-red-300"
                    >
                      Hapus
                    </button>
                  </td>
                </tr>
              {/each}
              {#if filteredExams.length === 0}
                <tr>
                  <td
                    colspan="6"
                    class="py-12 text-center text-orange-700 font-light"
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
            class="px-6 py-4 border-t border-orange-200 flex items-center justify-between"
          >
            <div class="text-sm text-orange-800">
              Menampilkan <span class="font-medium text-orange-950"
                >{filteredExams.length > 0
                  ? (currentPage - 1) * itemsPerPage + 1
                  : 0}</span
              >
              sampai
              <span class="font-medium text-orange-950"
                >{Math.min(
                  currentPage * itemsPerPage,
                  filteredExams.length,
                )}</span
              >
              dari
              <span class="font-medium text-orange-950"
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
                  : 'bg-white text-orange-900 hover:bg-orange-50 hover:text-orange-950'}"
              >
                Sebelumnya
              </button>
              <div class="flex items-center gap-1 px-2">
                {#each Array(totalPages) as _, i}
                  <button
                    onclick={() => goToPage(i + 1)}
                    class="w-8 h-8 flex items-center justify-center text-sm font-medium rounded-lg transition-colors {currentPage ===
                    i + 1
                      ? 'bg-indigo-600 text-orange-950'
                      : 'text-orange-800 hover:bg-white hover:text-orange-950'}"
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
                  : 'bg-white text-orange-900 hover:bg-orange-50 hover:text-orange-950'}"
              >
                Selanjutnya
              </button>
            </div>
          </div>
        {/if}
      </div>
    {:else if activeTab === "chart"}
      <div
        class="bg-white/60 backdrop-blur-md rounded-3xl border border-orange-200 shadow-lg shadow-orange-900/10 p-6 space-y-6"
      >
        <div class="flex flex-col sm:flex-row gap-4">
          <div class="w-full sm:max-w-xs">
            <label
              class="block text-sm font-medium text-orange-800 mb-2"
              for="chartUser">Pilih Murid</label
            >
            <select
              id="chartUser"
              bind:value={chartSelectedUserId}
              class="w-full bg-white border border-orange-300 rounded-xl px-4 py-2.5 text-orange-950 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all"
            >
              <option value="" disabled>-- Pilih Murid --</option>
              {#each users as u}
                <option value={u.id}>{u.username}</option>
              {/each}
            </select>
          </div>

          {#if chartSelectedUserId}
            <div class="w-full sm:max-w-[12rem]">
              <label
                class="block text-sm font-medium text-orange-800 mb-2"
                for="chartTimeframe">Rentang Waktu</label
              >
              <select
                id="chartTimeframe"
                bind:value={chartTimeframe}
                class="w-full bg-white border border-orange-300 rounded-xl px-4 py-2.5 text-orange-950 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all"
              >
                <option value="harian">Harian</option>
                <option value="mingguan">Mingguan</option>
                <option value="bulanan">Bulanan</option>
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
            class="py-16 text-center text-orange-700 font-light border-2 border-dashed border-orange-200 rounded-2xl"
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
      class="relative bg-[#EAE4BD] border border-orange-300 rounded-2xl shadow-2xl w-full max-w-md overflow-hidden animate-in zoom-in-95 duration-200"
    >
      <div
        class="p-6 border-b border-orange-200 flex justify-between items-center"
      >
        <h3 class="text-xl font-bold text-orange-950">
          {isEditing ? "Edit Nilai Ujian" : "Tambah Nilai Ujian"}
        </h3>
        <button
          onclick={closeModal}
          class="text-orange-700 hover:text-orange-950 transition-colors"
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

      <form onsubmit={saveExam} class="p-6 space-y-4">
        <div>
          <label
            class="block text-sm font-medium text-orange-800 mb-1.5"
            for="user">Pilih Murid</label
          >
          <select
            id="user"
            bind:value={formUserId}
            class="w-full bg-white border border-orange-300 rounded-xl px-4 py-2.5 text-orange-950 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all"
            required
          >
            <option value="" disabled>-- Pilih Murid --</option>
            {#each users as u}
              <option value={u.id}>{u.username}</option>
            {/each}
          </select>
        </div>

        <div>
          <label
            class="block text-sm font-medium text-orange-800 mb-1.5"
            for="date">Tanggal Ujian</label
          >
          <input
            type="date"
            id="date"
            bind:value={formExamDate}
            class="w-full bg-white border border-orange-300 rounded-xl px-4 py-2.5 text-orange-950 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all"
            required
          />
        </div>

        <div>
          <label
            class="block text-sm font-medium text-orange-800 mb-1.5"
            for="examName">Nama Ujian</label
          >
          <input
            type="text"
            id="examName"
            bind:value={formExamName}
            placeholder="Misal: Ujian Tengah Semester"
            class="w-full bg-white border border-orange-300 rounded-xl px-4 py-2.5 text-orange-950 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all"
            required
          />
        </div>

        <div>
          <label
            class="block text-sm font-medium text-orange-800 mb-1.5"
            for="subject">Mata Pelajaran</label
          >
          <select
            id="subject"
            bind:value={formSubjectId}
            class="w-full bg-white border border-orange-300 rounded-xl px-4 py-2.5 text-orange-950 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all"
            required
          >
            <option value="" disabled>-- Pilih Mata Pelajaran --</option>
            {#each subjects as s}
              <option value={s.id}>{s.name}</option>
            {/each}
          </select>
        </div>

        <div>
          <label
            class="block text-sm font-medium text-orange-800 mb-1.5"
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
            class="w-full bg-white border border-orange-300 rounded-xl px-4 py-2.5 text-orange-950 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all"
            required
          />
        </div>

        <div class="pt-4 flex justify-end gap-3">
          <button
            type="button"
            onclick={closeModal}
            class="px-4 py-2.5 text-sm font-medium text-orange-900 hover:text-orange-950 bg-white hover:bg-orange-50 rounded-xl transition-colors"
          >
            Batal
          </button>
          <button
            type="submit"
            class="px-4 py-2.5 text-sm font-medium text-orange-950 bg-indigo-600 hover:bg-indigo-500 rounded-xl transition-all shadow-md shadow-indigo-900/20"
          >
            {isEditing ? "Simpan Perubahan" : "Tambahkan"}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}

<style>
  /* Override global dark color-scheme agar sesuai dengan tema terang di halaman ini */
  select,
  input[type="date"] {
    color-scheme: light;
  }

  option {
    background-color: #ffffff;
    color: #431407; /* text-orange-950 */
  }
</style>
