<script lang="ts">
  import { onMount } from "svelte";
  import { page } from "$app/state";

  let examId = $derived(page.params.id);

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

  let exam = $state<Exam | null>(null);
  let isLoading = $state(true);
  let showLoadingSpinner = $state(false);
  let errorMsg = $state("");

  async function fetchExamDetail() {
    try {
      const res = await fetch(`/api/exams/${examId}`, {
        credentials: "include",
      });
      if (!res.ok) throw new Error("Gagal memuat detail nilai ujian");
      exam = await res.json();
    } catch (e) {
      errorMsg = e instanceof Error ? e.message : String(e);
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    fetchExamDetail();
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
</script>

<svelte:head>
  <title>Detail Nilai Ujian - Portal Siswa</title>
</svelte:head>

<div
  class="animate-in fade-in zoom-in-95 duration-500 max-w-4xl mx-auto py-8 sm:px-6"
>
  <div class="mb-8">
    <a
      href="/dashboard/exams"
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
        Detail Nilai Ujian
      </h1>
      <p class="mt-1 text-slate-500 text-sm tracking-wide">
        Informasi lengkap hasil ujian dan bukti unggahan.
      </p>
    </div>
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
  {:else if exam}
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
                  d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"
                ></path></svg
              >
              {exam.subject?.name || "Mata Pelajaran"}
            </div>
            <h2
              class="text-2xl sm:text-3xl font-semibold tracking-tight drop-shadow-md mb-2"
            >
              {exam.exam_name}
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
                  d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
                ></path></svg
              >
              {exam.user?.username || "Tidak diketahui"}
            </p>
          </div>

          <div
            class="bg-white/10 backdrop-blur-md border border-white/20 rounded-3xl p-6 text-center shadow-lg transform hover:scale-105 transition-transform duration-300"
          >
            <div
              class="text-sm text-indigo-100 font-medium uppercase tracking-wider mb-1"
            >
              Skor Akhir
            </div>
            <div class="text-4xl font-black drop-shadow-lg">{exam.score}</div>
          </div>
        </div>
      </div>

      <!-- Content Area -->
      <div class="p-5 sm:p-12 space-y-10">
        <!-- Info Grid -->
        <div class="grid grid-cols-1 gap-8">
          <div class="flex items-start gap-4">
            <div
              class="w-12 h-12 rounded-2xl bg-blue-50 text-blue-600 flex items-center justify-center flex-shrink-0 shadow-sm border border-blue-100"
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
                  d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                ></path></svg
              >
            </div>
            <div>
              <h4
                class="text-sm font-semibold text-slate-500 uppercase tracking-wide"
              >
                Tanggal Ujian
              </h4>
              <p class="mt-1 text-lg font-medium text-slate-900">
                {formatDate(exam.exam_date)}
              </p>
            </div>
          </div>
        </div>

        <hr class="border-slate-100" />

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
            Bukti Ujian
          </h3>

          {#if exam.image}
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
                  File gambar bukti ujian ini akan dihapus secara otomatis oleh
                  sistem dalam waktu 30 hari.
                </p>
              </div>
            </div>
            <div>
              <div
                class="rounded-3xl overflow-hidden border-4 border-white shadow-xl bg-slate-100"
              >
                <img
                  src={exam.image}
                  alt="Bukti {exam.exam_name}"
                  class="w-full h-auto object-cover max-h-[600px]"
                />
              </div>
              <div
                class="mt-4 flex flex-col sm:flex-row sm:items-center justify-between gap-4 px-2"
              >
                <span
                  class="text-slate-500 font-medium text-sm text-center sm:text-left"
                  >Dokumen Terlampir</span
                >
                <a
                  href={exam.image}
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
                  Buka Gambar
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
                Tidak Ada Bukti Lampiran
              </h4>
              <p class="text-slate-500 max-w-sm mx-auto">
                Ujian ini disubmit tanpa menyertakan dokumen gambar bukti
                tambahan.
              </p>
            </div>
          {/if}
        </div>
      </div>
    </div>
  {/if}
</div>
