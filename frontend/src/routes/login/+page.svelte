<script lang="ts">
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";

  let username = $state("");
  let password = $state("");
  let errorMsg = $state("");
  let isLoading = $state(false);

  onMount(async () => {
    try {
      const res = await fetch("http://localhost:8080/me", {
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
      const res = await fetch("http://localhost:8080/api/auth/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
        credentials: "include",
      });

      const data = await res.json();

      if (!res.ok) {
        throw new Error(data.error || "Login gagal");
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
</svelte:head>

<!-- Page wrapper -->
<div class="relative min-h-screen bg-[#0C134F] flex items-center justify-center px-4 py-12 overflow-hidden font-sans selection:bg-zinc-800 selection:text-white">

  <!-- Decorative blobs (dark theme adapted) -->
  <div class="pointer-events-none absolute -top-24 -right-20 w-96 h-96 rounded-full bg-indigo-500 opacity-20 blur-[100px] animate-blob-1" aria-hidden="true"></div>
  <div class="pointer-events-none absolute -bottom-16 -left-16 w-72 h-72 rounded-full bg-blue-500 opacity-20 blur-[100px] animate-blob-2" aria-hidden="true"></div>

  <!-- Content wrapper -->
  <main class="relative z-10 w-full max-w-sm flex flex-col items-center gap-6">

    <!-- Header -->
    <div class="flex flex-col items-center gap-2 text-center">
      <div class="w-13 h-13 bg-zinc-900 border border-zinc-700 text-white rounded-2xl flex items-center justify-center shadow-lg mb-1" aria-hidden="true">
        <svg width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M12 2L2 7l10 5 10-5-10-5z"/>
          <path d="M2 17l10 5 10-5"/>
          <path d="M2 12l10 5 10-5"/>
        </svg>
      </div>
      <h1 class="text-2xl font-semibold tracking-tight text-white drop-shadow-sm">Selamat datang</h1>
      <p class="text-sm text-blue-200 font-light">Masuk untuk melanjutkan ke portal belajar</p>
    </div>

    <!-- Card -->
    <div class="w-full bg-zinc-900/50 backdrop-blur-md rounded-2xl border border-zinc-700 shadow-xl shadow-black/50 p-8">
      <form onsubmit={handleLogin} novalidate class="flex flex-col gap-4">

        <!-- Error -->
        {#if errorMsg}
          <div class="flex items-start gap-2.5 bg-red-900/50 border border-red-700 text-red-200 rounded-xl px-3.5 py-3 text-sm leading-snug" role="alert">
            <svg class="w-4 h-4 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" viewBox="0 0 24 24" aria-hidden="true">
              <circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>
            </svg>
            <span>{errorMsg}</span>
          </div>
        {/if}

        <!-- Username -->
        <div class="flex flex-col gap-1.5">
          <label for="username" class="text-xs font-medium text-blue-100 tracking-wide uppercase">
            Username
          </label>
          <input
            id="username"
            name="username"
            type="text"
            required
            bind:value={username}
            placeholder="Masukkan username Anda"
            autocomplete="username"
            class="w-full px-3.5 py-2.5 text-sm text-white bg-zinc-800/80 border border-zinc-600 rounded-xl outline-none placeholder:text-zinc-500 focus:border-indigo-400 focus:bg-zinc-800 focus:ring-2 focus:ring-indigo-500/50 transition-all"
          />
        </div>

        <!-- Password -->
        <div class="flex flex-col gap-1.5">
          <label for="password" class="text-xs font-medium text-blue-100 tracking-wide uppercase">
            Password
          </label>
          <input
            id="password"
            name="password"
            type="password"
            required
            bind:value={password}
            placeholder="••••••••"
            autocomplete="current-password"
            class="w-full px-3.5 py-2.5 text-sm text-white bg-zinc-800/80 border border-zinc-600 rounded-xl outline-none placeholder:text-zinc-500 focus:border-indigo-400 focus:bg-zinc-800 focus:ring-2 focus:ring-indigo-500/50 transition-all"
          />
        </div>

        <!-- Remember me -->
        <div class="flex items-center gap-2">
          <input
            id="remember-me"
            name="remember-me"
            type="checkbox"
            class="w-3.5 h-3.5 accent-indigo-500 cursor-pointer border-zinc-600 bg-zinc-800"
          />
          <label for="remember-me" class="text-xs text-blue-200 cursor-pointer select-none">
            Ingat saya
          </label>
        </div>

        <!-- Submit -->
        <button
          type="submit"
          disabled={isLoading}
          class="mt-1 w-full flex items-center justify-center gap-2 px-4 py-2.5 bg-indigo-600 text-white text-sm font-medium rounded-xl border border-indigo-500 hover:bg-indigo-500 hover:-translate-y-px hover:shadow-lg disabled:opacity-60 disabled:cursor-not-allowed disabled:hover:translate-y-0 disabled:hover:shadow-none transition-all duration-200"
        >
          {#if isLoading}
            <span class="w-3.5 h-3.5 border-2 border-white/30 border-t-white rounded-full animate-spin flex-shrink-0" aria-hidden="true"></span>
            Memproses...
          {:else}
            Masuk
          {/if}
        </button>

      </form>

      <!-- Register link -->
      <p class="mt-5 text-center text-xs text-blue-300 font-light">
        Belum punya akun?
        <a href="/register" class="text-indigo-400 font-medium hover:text-indigo-300 hover:underline transition-colors">
          Daftar sekarang
        </a>
      </p>
    </div>

    <!-- Back link -->
    <a href="/" class="group flex items-center gap-1.5 text-xs text-zinc-400 hover:text-white transition-colors">
      <svg class="w-3.5 h-3.5 group-hover:-translate-x-0.5 transition-transform" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" viewBox="0 0 24 24" aria-hidden="true">
        <line x1="19" y1="12" x2="5" y2="12"/>
        <polyline points="12 19 5 12 12 5"/>
      </svg>
      Kembali ke Beranda
    </a>

  </main>
</div>

<style>
  @keyframes blob-drift-1 {
    from { transform: translate(0, 0) scale(1); }
    to   { transform: translate(24px, 18px) scale(1.06); }
  }
  @keyframes blob-drift-2 {
    from { transform: translate(0, 0) scale(1); }
    to   { transform: translate(-20px, -14px) scale(1.04); }
  }

  .animate-blob-1 {
    animation: blob-drift-1 12s ease-in-out infinite alternate;
  }
  .animate-blob-2 {
    animation: blob-drift-2 16s ease-in-out infinite alternate;
  }

  /* Custom w-13 h-13 (52px) not in default Tailwind scale */
  .w-13 { width: 3.25rem; }
  .h-13 { height: 3.25rem; }
</style>
