<script lang="ts">
  import { onMount } from "svelte";

  let categories: any[] = $state([]);
  let quizzes: any[] = $state([]);
  let myResults: any[] = $state([]);
  let userProfile: any = $state(null);
  let isLoading = $state(true);

  let filteredQuizzes = $derived(
    quizzes
      .filter((q) => q.is_published)
      .map((q) => {
        // Find if user has completed this quiz
        const result = myResults.find((r) => r.quiz_id === q.id);
        return {
          ...q,
          status: result ? "completed" : "new",
          score: result ? result.score : null,
          questionsCount: q.questions ? q.questions.length : 0,
        };
      }),
  );

  onMount(async () => {
    try {
      const [catRes, quizRes, resRes, userRes] = await Promise.all([
        fetch("/api/kuisapp/categories", { credentials: "include" }),
        fetch("/api/kuisapp/quizzes", { credentials: "include" }),
        fetch("/api/kuisapp/my-results", { credentials: "include" }),
        fetch("/api/kuisapp/me", { credentials: "include" }),
      ]);

      if (catRes.ok) {
        const d = await catRes.json();
        categories = d.data || [];
      }
      if (quizRes.ok) {
        const d = await quizRes.json();
        quizzes = d.data || [];
      }
      if (resRes.ok) {
        const d = await resRes.json();
        myResults = d.data || [];
      }
      if (userRes.ok) {
        const d = await userRes.json();
        userProfile = d.user;
      }
    } catch (e) {
      console.error("Failed to load dashboard data", e);
    } finally {
      isLoading = false;
    }
  });
</script>

<svelte:head>
  <title>Quiz | Dashboard</title>
</svelte:head>

<div class="space-y-8 pb-12">
  <!-- Welcome Banner -->
  <div
    class="relative overflow-hidden rounded-3xl bg-gradient-to-br from-indigo-600 via-violet-600 to-fuchsia-600 p-6 sm:p-10 text-white shadow-xl"
  >
    <div
      class="absolute top-0 right-0 -mt-16 -mr-16 text-white/10 hidden sm:block"
    >
      <svg
        class="w-64 h-64 transform rotate-12"
        fill="currentColor"
        viewBox="0 0 24 24"><path d="M12 2L2 22h20L12 2z" /></svg
      >
    </div>
    <div
      class="relative z-10 flex flex-col sm:flex-row justify-between items-center gap-5 sm:gap-6 text-center sm:text-left"
    >
      <div
        class="flex flex-col sm:flex-row items-center gap-4 sm:gap-5 w-full sm:w-auto"
      >
        <div
          class="w-16 h-16 sm:w-24 sm:h-24 rounded-full bg-white/20 backdrop-blur-md border-4 border-white/30 flex items-center justify-center text-white shadow-lg shrink-0"
        >
          <svg
            class="w-8 h-8 sm:w-14 sm:h-14 opacity-90"
            fill="currentColor"
            viewBox="0 0 24 24"
            ><path
              d="M12 12a5 5 0 100-10 5 5 0 000 10zm0 2c-5.33 0-8 2.67-8 8v1h16v-1c0-5.33-2.67-8-8-8z"
            /></svg
          >
        </div>
        <div>
          <div
            class="inline-flex items-center justify-center sm:justify-start gap-1 sm:gap-2 px-2.5 sm:px-3 py-1 rounded-full bg-white/20 backdrop-blur-sm border border-white/20 text-[10px] sm:text-xs font-bold uppercase tracking-widest text-indigo-100 mb-1 sm:mb-2"
          >
            Profil Peserta
          </div>
          <h1
            class="text-2xl sm:text-4xl font-black mb-0.5 sm:mb-1 tracking-tight"
          >
            {userProfile?.username || "Peserta"}
          </h1>
          <p class="text-indigo-100 text-xs sm:text-base opacity-90 capitalize">
            Status: <strong class="text-white"
              >{userProfile?.status || "Umum"}</strong
            >
          </p>
        </div>
      </div>
      {#if userProfile}
        <div
          class="bg-white/10 backdrop-blur-md px-4 sm:px-6 py-3 sm:py-4 rounded-2xl border border-white/20 text-center w-full sm:w-auto min-w-37.5"
        >
          <div
            class="text-indigo-100 text-xs sm:text-sm font-semibold uppercase tracking-wider mb-0.5 sm:mb-1"
          >
            Total Poin
          </div>
          <div
            class="text-3xl sm:text-4xl font-black text-amber-300 drop-shadow-sm"
          >
            {userProfile.points || 0}
          </div>
        </div>
      {/if}
    </div>
  </div>

  <!-- Quizzes -->
  <section>
    <h2 class="text-xl font-bold text-slate-800 mb-6">Kuis Tersedia</h2>
    {#if filteredQuizzes.length === 0}
      <div
        class="text-center py-12 bg-white rounded-2xl border border-slate-100 border-dashed"
      >
        <div class="text-4xl mb-3">📭</div>
        <h3 class="text-lg font-medium text-slate-700">
          Belum ada kuis yang tersedia saat ini.
        </h3>
      </div>
    {:else}
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
        {#each filteredQuizzes as quiz}
          <div
            class="bg-white rounded-2xl p-6 shadow-sm border border-slate-100 hover:shadow-md transition-shadow flex flex-col h-full gap-5 items-start justify-between group"
          >
            <div class="flex-1 w-full">
              <div class="flex items-center gap-3 mb-2">
                <span
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold bg-indigo-50 text-indigo-600"
                >
                  {categories.find((c) => c.id === quiz.category_id)?.name ||
                    "Tanpa Kategori"}
                </span>
                {#if quiz.status === "completed"}
                  <span
                    class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold bg-emerald-50 text-emerald-600"
                  >
                    Selesai
                  </span>
                {/if}
              </div>
              <h3
                class="text-lg font-bold text-slate-800 group-hover:text-indigo-600 transition-colors"
              >
                {quiz.title}
              </h3>
              <div
                class="flex items-center gap-4 mt-3 text-sm text-slate-500 font-medium"
              >
                <div class="flex items-center gap-1.5">
                  <svg
                    class="w-4 h-4 text-slate-400"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    ><path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                    /></svg
                  >
                  {quiz.timeLimit} Detik / Soal
                </div>
                <div class="flex items-center gap-1.5">
                  <svg
                    class="w-4 h-4 text-slate-400"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    ><path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                    /></svg
                  >
                  {quiz.questionsCount} Soal
                </div>
              </div>
            </div>

            <a
              href="/quiz-app/dashboard/quizzes/{quiz.id}"
              class="w-full mt-auto inline-flex items-center justify-center gap-2 px-6 py-3 bg-slate-900 hover:bg-indigo-600 text-white text-sm font-bold rounded-xl transition-colors shadow-sm hover:shadow-indigo-500/25 no-underline"
            >
              {#if quiz.status === "completed"}
                Ulangi Kuis
              {:else}
                Mulai Kuis
              {/if}
              <svg
                class="w-4 h-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M14 5l7 7m0 0l-7 7m7-7H3"
                /></svg
              >
            </a>
          </div>
        {/each}
      </div>
    {/if}
  </section>
</div>
