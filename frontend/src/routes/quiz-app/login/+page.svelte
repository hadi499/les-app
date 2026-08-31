<script lang="ts">
  import { goto } from "$app/navigation";

  let username = $state("");
  let password = $state("");
  let showPassword = $state(false);
  let errorMsg = $state("");
  let isLoading = $state(false);

  async function handleLogin(e: Event) {
    e.preventDefault();
    isLoading = true;
    errorMsg = "";

    try {
      const res = await fetch("/api/kuisapp/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
      });

      const data = await res.json();

      if (res.ok) {
        // Fix for iOS Safari race condition where Set-Cookie is dropped if we navigate too fast
        await new Promise((resolve) => setTimeout(resolve, 300));
        // Reload to let the layout fetch the user and redirect properly
        window.location.href = "/quiz-app/dashboard";
      } else {
        errorMsg = data.error || "Login gagal";
      }
    } catch (err) {
      errorMsg = "Terjadi kesalahan koneksi";
    } finally {
      isLoading = false;
    }
  }
</script>

<div class="max-w-md mx-auto mt-16 sm:mt-20">
  <div
    class="bg-white py-8 px-4 shadow-xl shadow-indigo-100/50 sm:rounded-2xl sm:px-10 border border-slate-100"
  >
    <div class="sm:mx-auto sm:w-full sm:max-w-md mb-8">
      <h2
        class="text-center text-3xl font-black bg-gradient-to-r from-indigo-600 to-violet-600 bg-clip-text text-transparent"
      >
        Masuk QuizApp
      </h2>
      <p class="mt-2 text-center text-sm text-slate-500">
        Atau <a
          href="/quiz-app/register"
          class="font-medium text-indigo-600 hover:text-indigo-500"
          >daftar akun baru</a
        >
      </p>
    </div>

    {#if errorMsg}
      <div
        class="mb-4 bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-xl text-sm font-medium flex items-center gap-2"
      >
        <svg class="w-5 h-5 shrink-0" fill="currentColor" viewBox="0 0 20 20"
          ><path
            fill-rule="evenodd"
            d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z"
            clip-rule="evenodd"
          ></path></svg
        >
        {errorMsg}
      </div>
    {/if}

    <form class="space-y-6" onsubmit={handleLogin}>
      <div>
        <label for="username" class="block text-sm font-semibold text-slate-700"
          >Username</label
        >
        <div class="mt-2">
          <input
            id="username"
            type="text"
            required
            bind:value={username}
            class="block w-full appearance-none rounded-xl border border-slate-200 px-4 py-3 placeholder-slate-400 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm transition-all"
            placeholder="Masukkan username"
          />
        </div>
      </div>

      <div>
        <label for="password" class="block text-sm font-semibold text-slate-700"
          >Password</label
        >
        <div class="mt-2 relative">
          <input
            id="password"
            type={showPassword ? "text" : "password"}
            required
            bind:value={password}
            class="block w-full appearance-none rounded-xl border border-slate-200 pl-4 pr-10 py-3 placeholder-slate-400 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm transition-all"
            placeholder="••••••••"
          />
          <button
            type="button"
            onclick={() => (showPassword = !showPassword)}
            class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 focus:outline-none"
            title={showPassword ? "Sembunyikan password" : "Tampilkan password"}
          >
            {#if showPassword}
              <svg
                class="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"
                /></svg
              >
            {:else}
              <svg
                class="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                /><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                /></svg
              >
            {/if}
          </button>
        </div>
      </div>

      <div>
        <button
          type="submit"
          disabled={isLoading}
          class="w-full flex justify-center py-3 px-4 border border-transparent rounded-xl shadow-sm text-sm font-bold text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-all disabled:opacity-70 cursor-pointer disabled:cursor-not-allowed"
        >
          {#if isLoading}
            <div
              class="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"
            ></div>
          {:else}
            Masuk
          {/if}
        </button>
      </div>
    </form>
  </div>
  <div class="mt-6 text-center">
    <a
      href="/quiz-app"
      class="inline-flex items-center gap-2 text-sm font-semibold text-slate-500 hover:text-indigo-600 transition-colors"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"
        ><path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M10 19l-7-7m0 0l7-7m-7 7h18"
        ></path></svg
      >
      Kembali ke Beranda
    </a>
  </div>
</div>
