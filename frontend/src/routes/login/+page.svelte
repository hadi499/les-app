<script lang="ts">
  import { goto } from "$app/navigation";
  import { onMount, getContext } from "svelte";
  import { page } from "$app/state";

  let authState = getContext<{isAuthenticated: boolean, authChecked: boolean}>("authState");

  let username = $state("");
  let password = $state("");
  let errorMsg = $state("");
  let successMsg = $state("");
  let isLoading = $state(false);
  let showPassword = $state(false);

  onMount(async () => {
    if (page.url.searchParams.get("registered") === "true") {
      successMsg = "Pendaftaran berhasil! Silakan masuk.";
      // Hapus query param dari URL tanpa mereload halaman
      const url = new URL(window.location.href);
      url.searchParams.delete("registered");
      window.history.replaceState({}, "", url);
    }

    try {
      const res = await fetch(`/me`, {
        credentials: "include",
      });
      if (res.ok) {
        const data = await res.json();
        if (data.authenticated) {
          goto("/");
        }
      }
    } catch (e) {
      // Abaikan error jaringan jika tidak bisa mengecek sesi
    }
  });

  async function handleLogin(e: SubmitEvent) {
    e.preventDefault();
    isLoading = true;
    errorMsg = "";

    try {
      const res = await fetch(`/api/auth/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
        credentials: "include",
      });

      const data = await res.json();

      if (!res.ok) {
        throw new Error(data.error || "Login gagal");
      }

      if (authState) {
        authState.isAuthenticated = true;
      }

      goto("/");
    } catch (err) {
      errorMsg = err instanceof Error ? err.message : String(err);
    } finally {
      isLoading = false;
    }
  }
</script>

<svelte:head>
  <title>Masuk — Portal Les</title>
  <meta name="description" content="Masuk ke portal belajar Anda." />
  <meta name="robots" content="noindex, nofollow" />
</svelte:head>

<!-- Page wrapper -->
<div
  class="relative min-h-screen bg-slate-50 flex items-center justify-center px-4 py-12 overflow-hidden font-sans selection:bg-blue-200 selection:text-blue-900"
>

  <!-- Content wrapper -->
  <main class="relative z-10 w-full max-w-sm flex flex-col items-center gap-6">
    <!-- Header -->
    <div class="flex flex-col items-center gap-2 text-center">
      <div
        class="w-13 h-13 bg-white border border-slate-200 text-blue-600 rounded-2xl flex items-center justify-center shadow-sm mb-1"
        aria-hidden="true"
      >
        <svg
          width="26"
          height="26"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d="M12 2L2 7l10 5 10-5-10-5z" />
          <path d="M2 17l10 5 10-5" />
          <path d="M2 12l10 5 10-5" />
        </svg>
      </div>
      <h1
        class="text-2xl font-bold tracking-tight text-slate-900 drop-shadow-sm"
      >
        Selamat datang
      </h1>
      <p class="text-sm text-slate-600 font-medium">
        Masuk untuk melanjutkan ke portal belajar
      </p>
    </div>

    <!-- Card -->
    <div
      class="w-full bg-white rounded-2xl border border-slate-200 shadow-xl shadow-slate-200/50 p-8"
    >
      <form onsubmit={handleLogin} novalidate class="flex flex-col gap-4">
        <!-- Error -->
        {#if errorMsg}
          <div
            class="flex items-start gap-2.5 bg-red-50 border border-red-200 text-red-800 rounded-xl px-3.5 py-3 text-sm leading-snug"
            role="alert"
          >
            <svg
              class="w-4 h-4 flex-shrink-0 mt-0.5"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
              viewBox="0 0 24 24"
              aria-hidden="true"
            >
              <circle cx="12" cy="12" r="10" /><line
                x1="12"
                y1="8"
                x2="12"
                y2="12"
              /><line x1="12" y1="16" x2="12.01" y2="16" />
            </svg>
            <span>{errorMsg}</span>
          </div>
        {/if}

        <!-- Success -->
        {#if successMsg}
          <div
            class="flex items-start gap-2.5 bg-emerald-50 border border-emerald-200 text-emerald-800 rounded-xl px-3.5 py-3 text-sm leading-snug"
            role="status"
          >
            <svg
              class="w-4 h-4 shrink-0 mt-0.5"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
              viewBox="0 0 24 24"
              aria-hidden="true"
            >
              <circle cx="12" cy="12" r="10" /><polyline
                points="9 12 11 14 15 10"
              />
            </svg>
            <span>{successMsg}</span>
          </div>
        {/if}

        <!-- Username -->
        <div class="flex flex-col gap-1.5">
          <label
            for="username"
            class="text-xs font-bold text-slate-600 tracking-wide uppercase"
          >
            Username
          </label>
          <input
            id="username"
            name="username"
            type="text"
            required
            bind:value={username}
            placeholder="Masukkan username Anda"
            class="w-full px-3.5 py-2.5 text-sm text-slate-900 bg-slate-50 border border-slate-200 rounded-xl outline-none placeholder:text-slate-400 focus:border-blue-400 focus:bg-white focus:ring-2 focus:ring-blue-400/50 transition-all"
          />
        </div>

        <!-- Password -->
        <div class="flex flex-col gap-1.5">
          <label
            for="password"
            class="text-xs font-bold text-slate-600 tracking-wide uppercase"
          >
            Password
          </label>
          <div class="relative">
            <input
              id="password"
              name="password"
              type={showPassword ? "text" : "password"}
              required
              bind:value={password}
              placeholder="••••••••"
              autocomplete="current-password"
              class="w-full px-3.5 py-2.5 pr-10 text-sm text-slate-900 bg-slate-50 border border-slate-200 rounded-xl outline-none placeholder:text-slate-400 focus:border-blue-400 focus:bg-white focus:ring-2 focus:ring-blue-400/50 transition-all"
            />
            <button
              type="button"
              class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 transition-colors"
              onclick={() => (showPassword = !showPassword)}
              aria-label={showPassword
                ? "Sembunyikan password"
                : "Tampilkan password"}
            >
              {#if !showPassword}
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="18"
                  height="18"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  class="icon icon-tabler icons-tabler-outline icon-tabler-eye"
                  ><path stroke="none" d="M0 0h24v24H0z" fill="none" /><path
                    d="M10 12a2 2 0 1 0 4 0a2 2 0 0 0 -4 0"
                  /><path
                    d="M21 12c-2.4 4 -5.4 6 -9 6c-3.6 0 -6.6 -2 -9 -6c2.4 -4 5.4 -6 9 -6c3.6 0 6.6 2 9 6"
                  /></svg
                >
              {:else}
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="18"
                  height="18"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  class="icon icon-tabler icons-tabler-outline icon-tabler-eye-off"
                  ><path stroke="none" d="M0 0h24v24H0z" fill="none" /><path
                    d="M10.585 10.587a2 2 0 0 0 2.829 2.828"
                  /><path
                    d="M16.681 16.673a8.717 8.717 0 0 1 -4.681 1.327c-3.6 0 -6.6 -2 -9 -6c1.272 -2.12 2.712 -3.678 4.32 -4.674m2.86 -1.146a9.055 9.055 0 0 1 1.82 -.18c3.6 0 6.6 2 9 6c-.666 1.11 -1.379 2.067 -2.138 2.87"
                  /><path d="M3 3l18 18" /></svg
                >
              {/if}
            </button>
          </div>
        </div>

        <!-- Submit -->
        <button
          type="submit"
          disabled={isLoading}
          class="mt-1 w-full flex items-center justify-center gap-2 px-4 py-2.5 bg-blue-600 text-white text-sm font-bold rounded-xl border border-blue-500 hover:bg-blue-500 hover:-translate-y-px hover:shadow-lg disabled:opacity-60 disabled:cursor-not-allowed disabled:hover:translate-y-0 disabled:hover:shadow-none transition-all duration-200"
        >
          {#if isLoading}
            <span
              class="w-3.5 h-3.5 border-2 border-white/30 border-t-white rounded-full animate-spin flex-shrink-0"
              aria-hidden="true"
            ></span>
            Memproses...
          {:else}
            Masuk
          {/if}
        </button>
      </form>
    </div>

    <!-- Back link -->
    <a
      href="/"
      class="group flex items-center gap-1.5 text-sm text-slate-500 hover:text-slate-800 font-medium transition-colors"
    >
      <svg
        class="w-3.5 h-3.5 group-hover:-translate-x-0.5 transition-transform"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
        viewBox="0 0 24 24"
        aria-hidden="true"
      >
        <line x1="19" y1="12" x2="5" y2="12" />
        <polyline points="12 19 5 12 12 5" />
      </svg>
      Kembali ke Beranda
    </a>
  </main>


</div>

<style>
  .w-13 {
    width: 3.25rem;
  }
  .h-13 {
    height: 3.25rem;
  }

  /* Override browser autofill styles */
  input:-webkit-autofill,
  input:-webkit-autofill:hover,
  input:-webkit-autofill:focus,
  input:-webkit-autofill:active {
    -webkit-box-shadow: 0 0 0 30px white inset !important;
    -webkit-text-fill-color: #0f172a !important; /* text-slate-900 */
    border-radius: 0.75rem; /* rounded-xl */
  }
</style>
