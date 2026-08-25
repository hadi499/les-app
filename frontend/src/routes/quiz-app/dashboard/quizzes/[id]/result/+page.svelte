<script lang="ts">
  import { onMount } from "svelte";
  import { fade, slide } from "svelte/transition";
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import katex from "katex";
  import "katex/dist/katex.min.css";

  let quizId = $derived(page.params.id);

  let finalScore = $state(0);
  let userAnswers = $state<any[]>([]);
  let quizTitle = $state("");
  let quizQuestionsLength = $state(0);
  let pointsEarned = $state(0);
  let pointsAlreadyClaimed = $state(false);
  let isLoaded = $state(false);
  let duration = $state(0);

  onMount(() => {
    // Ambil data hasil kuis dari sessionStorage
    const storedResult = sessionStorage.getItem(`kuisapp_result_${quizId}`);
    if (storedResult) {
      try {
        const parsed = JSON.parse(storedResult);
        finalScore = parsed.finalScore || 0;
        userAnswers = parsed.userAnswers || [];
        quizTitle = parsed.quizTitle || "Kuis";
        quizQuestionsLength = parsed.quizQuestionsLength || 0;
        pointsEarned = parsed.pointsEarned || 0;
        pointsAlreadyClaimed = parsed.pointsAlreadyClaimed || false;
        duration = parsed.duration || 0;
        isLoaded = true;
      } catch (e) {
        console.error("Gagal membaca hasil kuis", e);
        goto("/quiz-app/dashboard");
      }
    } else {
      // Jika tidak ada data hasil kuis, kembalikan ke dashboard
      goto("/quiz-app/dashboard");
    }
  });

  function renderText(text: string | null) {
    if (!text) return "";

    // Replace block math $$...$$
    let rendered = text.replace(/\$\$([\s\S]*?)\$\$/g, (match, math) => {
      try {
        return katex.renderToString(math, {
          displayMode: true,
          throwOnError: false,
        });
      } catch (e) {
        return match;
      }
    });

    // Replace inline math $...$
    rendered = rendered.replace(/\$([^$]*?)\$/g, (match, math) => {
      try {
        return katex.renderToString(math, {
          displayMode: false,
          throwOnError: false,
        });
      } catch (e) {
        return match;
      }
    });

    return rendered;
  }

  function formatTime(seconds: number) {
    if (!seconds || seconds < 0) return "0 Detik";
    return `${seconds} Detik`;
  }
</script>

<svelte:head>
  <title>Hasil {quizTitle} | KuisApp</title>
</svelte:head>

<div
  class="w-full flex flex-col items-center py-6 sm:py-12 relative overflow-visible"
>
  <div class="relative z-10 w-full max-w-2xl mx-auto flex flex-col gap-6">
    {#if isLoaded}
      <div
        in:fade={{ duration: 300 }}
        class="bg-white rounded-3xl p-8 shadow-xl border border-slate-200 flex flex-col gap-8 relative overflow-hidden"
      >
        <div
          class="absolute top-0 left-0 w-full h-2 bg-gradient-to-r from-emerald-400 to-teal-500"
        ></div>
        <!-- Score Section -->
        <div
          class="flex flex-col items-center gap-2 text-center relative z-10 border-b border-slate-200 pb-8"
        >
          <div
            class="text-sm font-bold tracking-[0.2em] uppercase text-slate-600 mb-2"
          >
            Skor Akhir
          </div>
          <div
            class="text-[5rem] leading-none font-black text-blue-600 drop-shadow-sm"
          >
            {Math.round(finalScore)}
          </div>
          <p class="text-slate-600 font-medium text-sm">
            Total Benar: {userAnswers.filter((a) => a.isCorrect).length} dari {quizQuestionsLength}
            Soal
          </p>
          {#if duration > 0}
            <div class="inline-flex items-center gap-1.5 px-3 py-1 bg-slate-100 text-slate-600 rounded-full text-xs font-bold mt-2">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
              {formatTime(duration)}
            </div>
          {/if}
          {#if pointsEarned > 0}
            <div
              class="mt-4 px-6 py-3 bg-emerald-50 border-2 border-emerald-200 rounded-2xl flex flex-col items-center justify-center gap-1"
            >
              <span class="text-sm font-bold text-emerald-700"
                >{pointsEarned} poin berhasil didapatkan</span
              >
            </div>
          {:else if pointsAlreadyClaimed}
            <div
              class="mt-4 px-4 py-2 bg-slate-100 rounded-xl text-xs font-semibold text-slate-500"
            >
              Poin sudah diklaim pada percobaan sebelumnya.
            </div>
          {/if}
        </div>

        <!-- Review Section -->
        <div class="flex flex-col gap-6 relative z-10">
          <h3
            class="text-lg font-bold tracking-[0.1em] text-slate-800 uppercase border-l-4 border-blue-400 pl-3"
          >
            Ulasan Jawaban
          </h3>

          <div
            class="flex flex-col gap-4 max-h-[50vh] overflow-y-auto pr-2 custom-scrollbar"
          >
            {#each userAnswers as ans, index}
              <div
                class="p-4 rounded-xl border border-slate-200 bg-white/50 space-y-2"
              >
                <div class="flex justify-between items-start gap-4">
                  <div class="flex flex-col gap-2">
                    {#if ans.image}
                      <img
                        src={ans.image}
                        alt="Gambar Pertanyaan"
                        class="w-full max-w-50 md:max-w-48 rounded-lg border border-slate-200 mb-1"
                      />
                    {/if}
                    <p class="font-semibold text-slate-800 text-sm m-0">
                      {index + 1}. {@html renderText(ans.question)}
                    </p>
                  </div>
                  {#if ans.isCorrect}
                    <span
                      class="inline-flex items-center justify-center px-2 py-1 rounded text-xs font-bold bg-green-100 text-green-700"
                      >Benar</span
                    >
                  {:else}
                    <span
                      class="inline-flex items-center justify-center px-2 py-1 rounded text-xs font-bold bg-red-100 text-red-700"
                      >Salah</span
                    >
                  {/if}
                </div>

                <div class="text-sm flex flex-col gap-1 mt-2">
                  <div class="flex gap-2 items-center">
                    <span class="text-slate-500 font-medium">Jawaban Anda:</span
                    >
                    <span
                      class={ans.isCorrect
                        ? "text-green-700 font-semibold"
                        : "text-red-600 font-semibold"}
                    >
                      {@html renderText(
                        ans.answer === null ? "Tidak Menjawab" : ans.answer,
                      )}
                    </span>
                  </div>
                  {#if !ans.isCorrect}
                    <div class="flex gap-2 items-center">
                      <span class="text-slate-500 font-medium"
                        >Jawaban Benar:</span
                      >
                      <span class="text-green-700 font-semibold"
                        >{@html renderText(ans.correct)}</span
                      >
                    </div>
                  {/if}
                </div>
              </div>
            {/each}
          </div>
        </div>

        <div class="flex justify-center pt-4 relative z-10 gap-4">
          <a
            href="/quiz-app/dashboard"
            class="inline-flex items-center justify-center px-6 py-3 text-xs tracking-[0.1em] font-bold uppercase text-slate-700 bg-white border border-slate-300 hover:border-slate-500 hover:text-slate-900 transition-all duration-300 rounded-lg shadow-sm cursor-pointer no-underline"
          >
            Selesai
          </a>
          <a
            href="/quiz-app/dashboard/quizzes/{quizId}"
            class="inline-flex items-center justify-center px-8 py-3 text-xs tracking-[0.2em] font-bold uppercase text-white bg-blue-600 border border-transparent hover:bg-blue-700 transition-all duration-300 rounded-lg overflow-hidden cursor-pointer shadow-md hover:shadow-lg"
          >
            Ulangi
          </a>
        </div>
      </div>
    {/if}
  </div>
</div>

<style>
  .custom-scrollbar::-webkit-scrollbar {
    width: 6px;
  }
  .custom-scrollbar::-webkit-scrollbar-track {
    background: rgba(241, 245, 249, 0.5); /* slate-100 */
    border-radius: 4px;
  }
  .custom-scrollbar::-webkit-scrollbar-thumb {
    background: rgba(148, 163, 184, 0.5); /* slate-400 */
    border-radius: 4px;
  }
  .custom-scrollbar::-webkit-scrollbar-thumb:hover {
    background: rgba(71, 85, 105, 0.8); /* slate-600 */
  }
</style>
