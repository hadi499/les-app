<script lang="ts">
  import Modal from "$lib/components/Modal.svelte";
  import { fade } from "svelte/transition";
  import { onMount } from "svelte";

  let categories: any[] = $state([]);
  let quizzes: any[] = $state([]);
  let myResults: any[] = $state([]);
  let userProfile: any = $state(null);
  let isLoading = $state(true);
  let showLoading = $state(false);

  // Change Password State
  let showChangePasswordModal = $state(false);
  let oldPassword = $state("");
  let newPassword = $state("");
  let confirmPassword = $state("");
  let passwordErrorMsg = $state("");
  let passwordSuccessMsg = $state("");
  let isChangingPassword = $state(false);
  let showOldPassword = $state(false);
  let showNewPassword = $state(false);
  let showConfirmPassword = $state(false);

  async function handleChangePassword(e: Event) {
    e.preventDefault();
    if (newPassword !== confirmPassword) {
      passwordErrorMsg = "Password Baru dan Konfirmasi Password tidak sama";
      return;
    }
    isChangingPassword = true;
    passwordErrorMsg = "";
    passwordSuccessMsg = "";

    try {
      const res = await fetch("/api/kuisapp/change-password", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ old_password: oldPassword, new_password: newPassword }),
        credentials: "include"
      });

      const data = await res.json();
      if (res.ok) {
        passwordSuccessMsg = "Password berhasil diubah!";
        oldPassword = "";
        newPassword = "";
        confirmPassword = "";
        setTimeout(() => {
          showChangePasswordModal = false;
          passwordSuccessMsg = "";
        }, 800);
      } else {
        passwordErrorMsg = data.error || "Gagal mengganti password";
      }
    } catch (err) {
      passwordErrorMsg = "Terjadi kesalahan koneksi";
    } finally {
      isChangingPassword = false;
    }
  }

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
    const loadingTimer = setTimeout(() => {
      showLoading = true;
    }, 250);

    try {
      const [catRes, quizRes, resRes, userRes] = await Promise.all([
        fetch("/api/kuisapp/categories", { credentials: "include" }),
        fetch("/api/kuisapp/quizzes", { credentials: "include" }),
        fetch("/api/kuisapp/my-results", { credentials: "include" }),
        fetch(`/api/kuisapp/me?t=${Date.now()}`, { credentials: "include", cache: "no-store" }),
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
      console.error(e);
    } finally {
      clearTimeout(loadingTimer);
      isLoading = false;
    }
  });
</script>

<svelte:head>
  <title>Quiz | Dashboard</title>
</svelte:head>

{#if isLoading}
  {#if showLoading}
    <div in:fade={{ duration: 300 }} class="w-full min-h-[50vh] flex flex-col items-center justify-center gap-4">
      <div class="w-10 h-10 border-4 border-slate-200 border-t-indigo-600 rounded-full animate-spin"></div>
      <div class="text-sm font-medium text-slate-500 animate-pulse">Memuat data kuis...</div>
    </div>
  {/if}
{:else}
  <div in:fade={{ duration: 300 }} class="space-y-8 pb-12">
    <!-- Welcome Banner -->
    <div
      class="relative overflow-hidden rounded-2xl bg-gradient-to-br from-indigo-600 via-violet-600 to-fuchsia-600 p-5 sm:p-6 text-white shadow-lg flex flex-col gap-4 sm:gap-5"
    >
      <!-- Background element -->
      <div class="absolute top-0 right-0 -mt-8 -mr-8 text-white/10 hidden sm:block">
        <svg class="w-32 h-32 transform rotate-12" fill="currentColor" viewBox="0 0 24 24"><path d="M12 2L2 22h20L12 2z" /></svg>
      </div>
      
      <!-- App Signature -->
      <div class="relative z-10 flex items-center justify-between">
        <div class="flex items-center">
          <span class="font-black tracking-wider uppercase text-xs sm:text-sm drop-shadow-sm">LB Quiz</span>
        </div>
      </div>

      <div class="relative z-10 flex flex-row justify-between items-center gap-3">
        <!-- Profile Info -->
        <div class="flex flex-row items-center gap-3 sm:gap-4 overflow-hidden">
          <!-- Texts -->
          <div class="min-w-0">
            <h1 class="text-lg sm:text-xl font-bold mb-0.5 tracking-tight truncate">
              {userProfile?.username || "Peserta"}
            </h1>
            <div class="flex items-center gap-2">
              <span class="inline-flex items-center px-2 py-0.5 rounded text-[10px] sm:text-xs font-medium bg-white/20 backdrop-blur-sm text-indigo-50 border border-white/10 capitalize">
                {userProfile?.status || "Umum"}
              </span>
              <button 
                onclick={() => showChangePasswordModal = true}
                class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded text-[10px] sm:text-xs font-medium bg-indigo-500/30 hover:bg-indigo-500/50 backdrop-blur-sm text-indigo-50 border border-indigo-300/30 transition-colors cursor-pointer"
                title="Ganti Password"
              >
                <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path></svg>
                <span class="hidden sm:inline">Ganti Password</span>
              </button>
            </div>
          </div>
        </div>

        <!-- Points Badge -->
        {#if userProfile}
          <div class="flex flex-col items-center shrink-0">
            <div class="text-[9px] sm:text-[10px] text-indigo-100/90 font-medium uppercase tracking-widest mb-1">Total Poin</div>
            <div class="flex items-center gap-1.5 bg-white/10 backdrop-blur-md px-3 py-1.5 rounded-xl border border-white/20 shadow-sm">
              <span class="text-xl sm:text-2xl font-bold text-amber-300 drop-shadow-sm">{userProfile.points || 0}</span>
            </div>
          </div>
        {/if}
      </div>
    </div>

    <!-- Quizzes -->
    <section>
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
{/if}

<!-- Change Password Modal -->
<Modal show={showChangePasswordModal} onclose={() => { showChangePasswordModal = false; passwordErrorMsg = ''; passwordSuccessMsg = ''; }}>
  <div class="p-2">
    <div class="flex items-center gap-3 mb-6">
      <div class="w-10 h-10 rounded-full bg-indigo-100 flex items-center justify-center text-indigo-600">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path></svg>
      </div>
      <h3 class="text-xl font-bold text-slate-800">Ganti Password</h3>
    </div>

    {#if passwordErrorMsg}
      <div class="mb-4 bg-red-50 text-red-600 px-4 py-3 rounded-xl text-sm font-medium border border-red-100 flex items-start gap-2">
        <svg class="w-5 h-5 shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"></path></svg>
        {passwordErrorMsg}
      </div>
    {/if}

    {#if passwordSuccessMsg}
      <div class="mb-4 bg-emerald-50 text-emerald-600 px-4 py-3 rounded-xl text-sm font-medium border border-emerald-100 flex items-start gap-2">
        <svg class="w-5 h-5 shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"></path></svg>
        {passwordSuccessMsg}
      </div>
    {/if}

    <form class="space-y-4" onsubmit={handleChangePassword}>
      <div>
        <label for="oldPassword" class="block text-sm font-semibold text-slate-700 mb-1.5">Password Lama</label>
        <div class="relative">
          <input id="oldPassword" type={showOldPassword ? "text" : "password"} required bind:value={oldPassword} class="block w-full rounded-xl border border-slate-200 px-4 py-2.5 pr-12 text-sm focus:border-indigo-500 focus:ring-indigo-500 transition-colors" placeholder="Masukkan password lama">
          <button type="button" class="absolute inset-y-0 right-0 flex items-center px-4 text-slate-400 hover:text-slate-600" onclick={() => showOldPassword = !showOldPassword}>
            {#if showOldPassword}
              <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.29 3.29m0 0a10.05 10.05 0 015.51-2.122c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0l-3.29-3.29"/></svg>
            {:else}
              <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/></svg>
            {/if}
          </button>
        </div>
      </div>

      <div>
        <label for="newPassword" class="block text-sm font-semibold text-slate-700 mb-1.5">Password Baru</label>
        <div class="relative">
          <input id="newPassword" type={showNewPassword ? "text" : "password"} required bind:value={newPassword} class="block w-full rounded-xl border border-slate-200 px-4 py-2.5 pr-12 text-sm focus:border-indigo-500 focus:ring-indigo-500 transition-colors" placeholder="Minimal 6 karakter">
          <button type="button" class="absolute inset-y-0 right-0 flex items-center px-4 text-slate-400 hover:text-slate-600" onclick={() => showNewPassword = !showNewPassword}>
            {#if showNewPassword}
              <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.29 3.29m0 0a10.05 10.05 0 015.51-2.122c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0l-3.29-3.29"/></svg>
            {:else}
              <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/></svg>
            {/if}
          </button>
        </div>
      </div>

      <div>
        <label for="confirmPassword" class="block text-sm font-semibold text-slate-700 mb-1.5">Konfirmasi Password Baru</label>
        <div class="relative">
          <input id="confirmPassword" type={showConfirmPassword ? "text" : "password"} required bind:value={confirmPassword} class="block w-full rounded-xl border border-slate-200 px-4 py-2.5 pr-12 text-sm focus:border-indigo-500 focus:ring-indigo-500 transition-colors" placeholder="Ketik ulang password baru">
          <button type="button" class="absolute inset-y-0 right-0 flex items-center px-4 text-slate-400 hover:text-slate-600" onclick={() => showConfirmPassword = !showConfirmPassword}>
            {#if showConfirmPassword}
              <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.29 3.29m0 0a10.05 10.05 0 015.51-2.122c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0l-3.29-3.29"/></svg>
            {:else}
              <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/></svg>
            {/if}
          </button>
        </div>
      </div>

      <div class="pt-2 flex justify-end gap-3">
        <button type="button" onclick={() => { showChangePasswordModal = false; passwordErrorMsg = ''; passwordSuccessMsg = ''; }} class="px-5 py-2.5 text-sm font-semibold text-slate-600 hover:bg-slate-100 rounded-xl transition-colors">Batal</button>
        <button type="submit" disabled={isChangingPassword || passwordSuccessMsg !== ''} class="px-5 py-2.5 text-sm font-bold text-white bg-indigo-600 hover:bg-indigo-700 rounded-xl shadow-sm hover:shadow-indigo-500/25 transition-all disabled:opacity-70 disabled:cursor-not-allowed flex items-center gap-2">
          {#if isChangingPassword}
            <div class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
            Menyimpan...
          {:else}
            Simpan Password
          {/if}
        </button>
      </div>
    </form>
  </div>
</Modal>
