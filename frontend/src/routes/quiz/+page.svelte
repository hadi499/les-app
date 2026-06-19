<script lang="ts">
  import { onMount } from "svelte";

  type Quiz = {
    id: number;
    title: string;
    category: string;
    timeLimit: number;
  };

  let quizzes: Quiz[] = $state([]);
  let isLoading = $state(true);

  import { goto } from "$app/navigation";

  onMount(async () => {
    try {
      // Cek apakah user sudah login
      const authRes = await fetch(`/me`, { credentials: "include" });
      const authData = await authRes.json();
      if (!authData.authenticated) {
        goto("/login");
        return;
      }

      const res = await fetch(`/api/quizzes`, {
        credentials: "include",
      });
      if (res.ok) {
        const json = await res.json();
        quizzes = json.data || [];
      } else {
        console.error("Gagal memuat kuis:", res.status);
      }
    } catch (e) {
      console.error(e);
    } finally {
      isLoading = false;
    }
  });
</script>

<svelte:head>
  <title>Kuis Tersedia | Les Balongarut</title>
</svelte:head>

<div
  class="min-h-screen bg-orange-50 font-sans flex flex-col relative overflow-x-hidden pt-24 px-4"
>
  <!-- Background Ambient -->
  <div class="absolute inset-0 z-0 pointer-events-none fixed">
    <div
      class="absolute top-1/4 left-1/4 w-[600px] h-[600px] bg-white/40 rounded-full blur-[120px]"
    ></div>
    <div
      class="absolute bottom-1/4 right-1/4 w-[500px] h-[500px] bg-amber-100/50 rounded-full blur-[120px]"
    ></div>
    <div
      class="absolute inset-0 bg-[linear-gradient(to_right,#fbbf24_1px,transparent_1px),linear-gradient(to_bottom,#fbbf24_1px,transparent_1px)] bg-[size:4rem_4rem] [mask-image:radial-gradient(ellipse_80%_80%_at_50%_50%,#000_20%,transparent_100%)] opacity-20"
    ></div>
  </div>

  <div class="relative z-10 w-full max-w-4xl mx-auto pb-12">
    <div class="text-center mb-12">
      <h1
        class="text-4xl sm:text-5xl font-bold tracking-[0.1em] text-orange-950 uppercase mb-4"
      >
        Daftar Kuis
      </h1>
      <p
        class="text-orange-800 tracking-[0.1em] text-sm font-medium max-w-lg mx-auto"
      >
        Uji kemampuanmu dengan memilih salah satu kuis yang tersedia di bawah
        ini.
      </p>
    </div>

    {#if isLoading}
      <div class="flex justify-center p-12">
        <div
          class="w-10 h-10 border-4 border-orange-200 border-t-orange-600 rounded-full animate-spin"
        ></div>
      </div>
    {:else if quizzes.length === 0}
      <div
        class="text-center bg-white/60 backdrop-blur-md p-8 rounded-2xl border border-orange-200 shadow-sm"
      >
        <p class="text-orange-900 font-medium">
          Belum ada kuis yang tersedia saat ini.
        </p>
      </div>
    {:else}
      <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
        {#each quizzes as quiz}
          <div
            class="bg-white/80 backdrop-blur-md rounded-2xl p-6 border border-orange-200 shadow-sm hover:shadow-md hover:border-orange-400 transition-all group flex flex-col"
          >
            <div class="mb-4">
              <span
                class="inline-block px-3 py-1 bg-orange-100 text-orange-800 text-xs font-bold uppercase tracking-wider rounded-lg mb-3"
              >
                {quiz.category}
              </span>
              <h2
                class="text-xl font-bold text-orange-950 line-clamp-2 leading-tight"
              >
                {quiz.title}
              </h2>
            </div>
            <div
              class="mt-auto pt-4 flex items-center justify-between border-t border-orange-100"
            >
              <div
                class="flex items-center gap-1.5 text-xs font-semibold text-orange-800"
              >
                <svg
                  class="w-4 h-4"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                  />
                </svg>
                {quiz.timeLimit}s / soal
              </div>
              <a
                href="/quiz/{quiz.id}"
                class="text-sm font-bold text-orange-600 group-hover:text-orange-700 uppercase tracking-wider transition-colors no-underline flex items-center gap-1"
              >
                Mulai
                <svg
                  class="w-4 h-4 transform group-hover:translate-x-1 transition-transform"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M9 5l7 7-7 7"
                  />
                </svg>
              </a>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>
