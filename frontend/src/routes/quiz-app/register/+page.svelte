<script lang="ts">
  import { goto } from "$app/navigation";
  
  let username = $state("");
  let password = $state("");
  let confirmPassword = $state("");
  let status = $state("umum");
  let errorMsg = $state("");
  let isLoading = $state(false);

  let showPassword = $state(false);
  let showConfirmPassword = $state(false);

  async function handleRegister(e: Event) {
    e.preventDefault();
    if (password !== confirmPassword) {
      errorMsg = "Password dan Konfirmasi Password tidak sama";
      return;
    }
    isLoading = true;
    errorMsg = "";
    
    try {
      const res = await fetch("/api/kuisapp/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password, status }),
      });
      
      const data = await res.json();
      
      if (res.ok) {
        // Auto login
        const loginRes = await fetch("/api/kuisapp/login", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ username, password }),
        });
        if (loginRes.ok) {
          // Fix for iOS Safari race condition where Set-Cookie is dropped if we navigate too fast
          await new Promise(resolve => setTimeout(resolve, 300));
          window.location.href = "/quiz-app/dashboard";
        } else {
          goto("/quiz-app/login");
        }
      } else {
        errorMsg = data.error || "Gagal mendaftar";
      }
    } catch (err) {
      errorMsg = "Terjadi kesalahan koneksi";
    } finally {
      isLoading = false;
    }
  }
</script>

<div class="max-w-md mx-auto mt-16 sm:mt-20">
  <div class="bg-white py-8 px-4 shadow-xl shadow-indigo-100/50 sm:rounded-2xl sm:px-10 border border-slate-100">
    <div class="sm:mx-auto sm:w-full sm:max-w-md mb-8">
      <h2 class="text-center text-3xl font-black bg-gradient-to-r from-indigo-600 to-violet-600 bg-clip-text text-transparent">
        Daftar QuizApp
      </h2>
      <p class="mt-2 text-center text-sm text-slate-500">
        Sudah punya akun? <a href="/quiz-app/login" class="font-medium text-indigo-600 hover:text-indigo-500">Masuk di sini</a>
      </p>
    </div>

    {#if errorMsg}
      <div class="mb-4 bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-xl text-sm font-medium flex items-center gap-2">
        <svg class="w-5 h-5 shrink-0" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd"></path></svg>
        {errorMsg}
      </div>
    {/if}

    <form class="space-y-5" onsubmit={handleRegister}>
      <div>
        <label for="username" class="block text-sm font-semibold text-slate-700">Username</label>
        <div class="mt-2">
          <input id="username" type="text" required bind:value={username} class="block w-full appearance-none rounded-xl border border-slate-200 px-4 py-3 placeholder-slate-400 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm transition-all" placeholder="Pilih username">
        </div>
      </div>

      <div>
        <label for="password" class="block text-sm font-semibold text-slate-700">Password</label>
        <div class="mt-2 relative">
          <input id="password" type={showPassword ? "text" : "password"} required bind:value={password} class="block w-full appearance-none rounded-xl border border-slate-200 px-4 py-3 pr-12 placeholder-slate-400 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm transition-all" placeholder="Minimal 6 karakter">
          <button type="button" class="absolute inset-y-0 right-0 flex items-center px-4 text-slate-400 hover:text-slate-600" onclick={() => showPassword = !showPassword}>
            {#if showPassword}
              <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.29 3.29m0 0a10.05 10.05 0 015.51-2.122c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0l-3.29-3.29"/></svg>
            {:else}
              <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/></svg>
            {/if}
          </button>
        </div>
      </div>

      <div>
        <label for="confirmPassword" class="block text-sm font-semibold text-slate-700">Konfirmasi Password</label>
        <div class="mt-2 relative">
          <input id="confirmPassword" type={showConfirmPassword ? "text" : "password"} required bind:value={confirmPassword} class="block w-full appearance-none rounded-xl border border-slate-200 px-4 py-3 pr-12 placeholder-slate-400 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm transition-all" placeholder="Ketik ulang password">
          <button type="button" class="absolute inset-y-0 right-0 flex items-center px-4 text-slate-400 hover:text-slate-600" onclick={() => showConfirmPassword = !showConfirmPassword}>
            {#if showConfirmPassword}
              <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.29 3.29m0 0a10.05 10.05 0 015.51-2.122c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0l-3.29-3.29"/></svg>
            {:else}
              <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/></svg>
            {/if}
          </button>
        </div>
      </div>

      <div>
        <label for="status" class="block text-sm font-semibold text-slate-700 mb-2">Status Pendaftar</label>
        <div class="grid grid-cols-2 gap-3">
          <label class="relative flex cursor-pointer rounded-xl border bg-white p-3 shadow-sm focus:outline-none {status === 'pelajar' ? 'border-indigo-500 ring-1 ring-indigo-500' : 'border-slate-200'}">
            <input type="radio" name="status" value="pelajar" bind:group={status} class="sr-only">
            <span class="flex flex-1">
              <span class="flex flex-col">
                <span class="block text-sm font-medium text-slate-900">Pelajar</span>
                <span class="mt-1 flex items-center text-xs text-slate-500">Siswa/Mahasiswa</span>
              </span>
            </span>
            <svg class="h-5 w-5 {status === 'pelajar' ? 'text-indigo-600' : 'text-transparent'}" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
            </svg>
          </label>
          <label class="relative flex cursor-pointer rounded-xl border bg-white p-3 shadow-sm focus:outline-none {status === 'umum' ? 'border-indigo-500 ring-1 ring-indigo-500' : 'border-slate-200'}">
            <input type="radio" name="status" value="umum" bind:group={status} class="sr-only">
            <span class="flex flex-1">
              <span class="flex flex-col">
                <span class="block text-sm font-medium text-slate-900">Umum</span>
                <span class="mt-1 flex items-center text-xs text-slate-500">Masyarakat Umum</span>
              </span>
            </span>
            <svg class="h-5 w-5 {status === 'umum' ? 'text-indigo-600' : 'text-transparent'}" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
            </svg>
          </label>
        </div>
      </div>

      <div>
        <button type="submit" disabled={isLoading} class="mt-2 w-full flex justify-center py-3 px-4 border border-transparent rounded-xl shadow-sm text-sm font-bold text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-all disabled:opacity-70">
          {#if isLoading}
            <div class="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
          {:else}
            Daftar Sekarang
          {/if}
        </button>
      </div>
    </form>
  </div>
  <div class="mt-6 text-center">
    <a href="/quiz-app" class="inline-flex items-center gap-2 text-sm font-semibold text-slate-500 hover:text-indigo-600 transition-colors">
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path></svg>
      Kembali ke Beranda
    </a>
  </div>
</div>
