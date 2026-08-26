<script lang="ts">
  import { onMount } from "svelte";
  import { fade } from "svelte/transition";

  let myResults: any[] = $state([]);
  let isLoading = $state(true);
  let showLoading = $state(false);
  let errorMsg = $state("");

  onMount(async () => {
    const loadingTimer = setTimeout(() => {
      showLoading = true;
    }, 250);

    try {
      const res = await fetch("/api/kuisapp/my-results", {
        credentials: "include",
      });
      if (!res.ok) {
        throw new Error("Gagal mengambil data riwayat kuis.");
      }
      const data = await res.json();
      // Sort by finished_at descending (newest first)
      myResults = (data.data || []).sort((a: any, b: any) => {
        return (
          new Date(b.finished_at).getTime() - new Date(a.finished_at).getTime()
        );
      });
    } catch (e) {
      errorMsg = e instanceof Error ? e.message : String(e);
    } finally {
      clearTimeout(loadingTimer);
      isLoading = false;
    }
  });

  function formatDate(dateString: string) {
    const d = new Date(dateString);
    const datePart = new Intl.DateTimeFormat("id-ID", {
      day: "numeric",
      month: "long",
      year: "numeric",
    }).format(d);

    const timePart = [
      d.getHours().toString().padStart(2, "0"),
      d.getMinutes().toString().padStart(2, "0"),
      d.getSeconds().toString().padStart(2, "0"),
    ].join(":");

    return `${datePart}, ${timePart}`;
  }

  function formatDuration(seconds: number) {
    if (!seconds || seconds < 0) return "0 Detik";
    return `${seconds} Detik`;
  }

  function getScoreColor(score: number) {
    if (score >= 80) return "text-emerald-600 bg-emerald-50 border-emerald-100";
    if (score >= 60) return "text-yellow-600 bg-yellow-50 border-yellow-100";
    return "text-red-600 bg-red-50 border-red-100";
  }
</script>

<svelte:head>
  <title>Quiz | Riwayat Kuis</title>
</svelte:head>

<div class="space-y-8 pb-12 animate-in fade-in duration-500">
  <!-- Header Banner -->
  <div
    class="relative overflow-hidden rounded-3xl bg-gradient-to-r from-slate-800 to-slate-900 p-8 sm:p-10 text-white shadow-xl"
  >
    <div class="absolute top-0 right-0 -mt-16 -mr-16 text-white/5">
      <svg
        class="w-64 h-64 transform -rotate-12"
        fill="currentColor"
        viewBox="0 0 24 24"><path d="M12 2L2 22h20L12 2z" /></svg
      >
    </div>
    <div
      class="relative z-10 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-6"
    >
      <div>
        <h1 class="text-2xl sm:text-3xl font-black mb-2 tracking-tight">
          Riwayat Kuis 📜
        </h1>
        <p class="text-slate-300 max-w-xl text-md">
          Lihat rekam jejak nilai dari semua kuis yang telah Anda kerjakan.
        </p>
      </div>
    </div>
  </div>

  <!-- Content -->
  <section>
    {#if isLoading}
      {#if showLoading}
        <div in:fade={{ duration: 300 }} class="flex justify-center py-20">
          <div
            class="w-10 h-10 border-4 border-slate-200 border-t-indigo-600 rounded-full animate-spin"
          ></div>
        </div>
      {/if}
    {:else if errorMsg}
      <div
        class="bg-red-50 text-red-700 p-6 rounded-2xl border border-red-100 font-medium flex items-center gap-3 shadow-sm"
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
    {:else if myResults.length === 0}
      <div
        class="text-center py-16 bg-white rounded-3xl border-2 border-slate-100 border-dashed"
      >
        <div class="text-5xl mb-4">📭</div>
        <h3 class="text-xl font-bold text-slate-800 mb-2">
          Belum ada riwayat kuis
        </h3>
        <p class="text-slate-500 max-w-md mx-auto mb-6">
          Anda belum mengerjakan kuis apa pun. Mulai kerjakan kuis sekarang
          untuk melihat nilainya di sini!
        </p>
        <a
          href="/quiz-app/dashboard"
          class="inline-flex items-center justify-center gap-2 px-6 py-3 bg-indigo-600 hover:bg-indigo-700 text-white font-bold rounded-xl transition-colors shadow-sm"
        >
          Pergi ke Dashboard
        </a>
      </div>
    {:else}
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
        {#each myResults as result}
          <a
            href="/quiz-app/dashboard/quizzes/{result.quiz_id}"
            class="bg-white rounded-2xl p-6 border border-slate-100 shadow-sm hover:shadow-md hover:border-indigo-200 transition-all relative block no-underline group"
          >
            <div class="flex justify-between items-start mb-4">
              <div>
                <div
                  class="font-bold text-slate-900 text-lg group-hover:text-indigo-600 transition-colors"
                >
                  {result.quiz?.title || "Kuis Tidak Diketahui"}
                </div>
              </div>
              <span
                class="inline-flex items-center justify-center px-3 py-1 rounded-xl text-sm font-bold border shrink-0 ml-4 {getScoreColor(
                  result.score,
                )}"
              >
                {Math.round(result.score)}
              </span>
            </div>

            <div
              class="flex items-center justify-between mt-4 pt-4 border-t border-slate-100"
            >
              <div class="flex flex-col gap-1">
                <div
                  class="text-xs font-medium text-slate-500 flex items-center gap-1.5"
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
                      d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                    /></svg
                  >
                  {formatDate(result.finished_at)}
                </div>
                {#if result.duration > 0}
                <div
                  class="text-xs font-medium text-slate-500 flex items-center gap-1.5"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                  {formatDuration(result.duration)}
                </div>
                {/if}
              </div>
              <div>
                {#if result.points_earned > 0}
                  <span
                    class="inline-flex items-center justify-center px-2.5 py-1 rounded-lg text-xs font-bold bg-emerald-100 text-emerald-700 border border-emerald-200"
                  >
                    +{result.points_earned} Poin
                  </span>
                {:else}
                  <span class="text-slate-400 text-xs font-medium">0 Poin</span>
                {/if}
              </div>
            </div>
          </a>
        {/each}
      </div>
    {/if}
  </section>
</div>
