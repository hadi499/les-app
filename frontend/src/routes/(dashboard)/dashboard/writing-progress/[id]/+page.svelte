<script lang="ts">
  import { onMount } from "svelte";
  import { page } from "$app/state";

  let progressId = $derived(page.params.id);

  type WritingProgress = {
    id: number;
    user_id: number;
    date: string;
    image: string;
    user?: { username: string };
  };

  let progress = $state<WritingProgress | null>(null);
  let isLoading = $state(true);
  let errorMsg = $state("");

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
</script>

<svelte:head>
  <title>Detail Perkembangan Menulis - Portal PAUD/TK</title>
</svelte:head>

<div
  class="animate-in fade-in zoom-in-95 duration-500 max-w-4xl mx-auto py-8 sm:px-6"
>
  <div class="mb-8">
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

  {#if isLoading}
    <div class="flex justify-center p-20">
      <div
        class="w-12 h-12 border-4 border-slate-200 border-t-indigo-500 rounded-full animate-spin"
      ></div>
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
