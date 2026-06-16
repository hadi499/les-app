<script lang="ts">
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";

  let username = $state("");
  let password = $state("");
  let confirmPassword = $state("");
  let errorMsg = $state("");
  let successMsg = $state("");
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
      // Abaikan error jaringan
    }
  });

  async function handleRegister(e: SubmitEvent) {
    e.preventDefault();
    isLoading = true;
    errorMsg = "";
    successMsg = "";

    if (password !== confirmPassword) {
      errorMsg = "Password dan konfirmasi password tidak cocok";
      isLoading = false;
      return;
    }

    try {
      const res = await fetch("http://localhost:8080/api/auth/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
        // Role and avatar use default in backend
      });

      const data = await res.json();

      if (!res.ok) {
        throw new Error(data.error || "Pendaftaran gagal");
      }

      successMsg = "Pendaftaran berhasil! Mengalihkan ke halaman login...";
      setTimeout(() => {
        goto("/login");
      }, 2000);
    } catch (err) {
      errorMsg = err instanceof Error ? err.message : String(err);
    } finally {
      isLoading = false;
    }
  }
</script>

<svelte:head>
  <title>Daftar — Portal Les</title>
  <meta
    name="description"
    content="Buat akun baru untuk mengakses portal belajar."
  />
</svelte:head>

<!-- Page wrapper -->
<div
  class="relative min-h-screen bg-orange-50 flex items-center justify-center px-4 py-12 overflow-hidden font-sans selection:bg-orange-200 selection:text-orange-900"
  style="color-scheme: light;"
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

  <!-- Content wrapper -->
  <main class="relative z-10 w-full max-w-sm flex flex-col items-center gap-6">
    <!-- Header -->
    <div class="flex flex-col items-center gap-2 text-center">
      <div
        class="w-13 h-13 bg-white border border-orange-200 text-orange-600 rounded-2xl flex items-center justify-center shadow-sm mb-1"
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
        class="text-2xl font-bold tracking-tight text-orange-950 drop-shadow-sm"
      >
        Buat akun baru
      </h1>
      <p class="text-sm text-orange-800 font-medium">
        Daftar untuk mulai belajar di portal kami
      </p>
    </div>

    <!-- Card -->
    <div
      class="w-full bg-white/60 backdrop-blur-md rounded-2xl border border-orange-200 shadow-lg p-8"
    >
      <form onsubmit={handleRegister} novalidate class="flex flex-col gap-4">
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
              class="w-4 h-4 flex-shrink-0 mt-0.5"
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
            class="text-xs font-bold text-orange-800 tracking-wide uppercase"
          >
            Username
          </label>
          <input
            id="username"
            name="username"
            type="text"
            required
            bind:value={username}
            placeholder="Pilih username"
            autocomplete="username"
            class="w-full px-3.5 py-2.5 text-sm text-orange-950 bg-white/80 border border-orange-200 rounded-xl outline-none placeholder:text-orange-400 focus:border-orange-400 focus:bg-white focus:ring-2 focus:ring-orange-400/50 transition-all"
          />
        </div>

        <!-- Password -->
        <div class="flex flex-col gap-1.5">
          <label
            for="password"
            class="text-xs font-bold text-orange-800 tracking-wide uppercase"
          >
            Password
          </label>
          <input
            id="password"
            name="password"
            type="password"
            required
            bind:value={password}
            placeholder="••••••••"
            autocomplete="new-password"
            class="w-full px-3.5 py-2.5 text-sm text-orange-950 bg-white/80 border border-orange-200 rounded-xl outline-none placeholder:text-orange-400 focus:border-orange-400 focus:bg-white focus:ring-2 focus:ring-orange-400/50 transition-all"
          />
        </div>

        <!-- Confirm Password -->
        <div class="flex flex-col gap-1.5">
          <label
            for="confirmPassword"
            class="text-xs font-bold text-orange-800 tracking-wide uppercase"
          >
            Konfirmasi Password
          </label>
          <input
            id="confirmPassword"
            name="confirmPassword"
            type="password"
            required
            bind:value={confirmPassword}
            placeholder="••••••••"
            autocomplete="new-password"
            class="w-full px-3.5 py-2.5 text-sm text-orange-950 bg-white/80 border border-orange-200 rounded-xl outline-none placeholder:text-orange-400 focus:border-orange-400 focus:bg-white focus:ring-2 focus:ring-orange-400/50 transition-all"
          />
        </div>

        <!-- Submit -->
        <button
          type="submit"
          disabled={isLoading}
          class="mt-1 w-full flex items-center justify-center gap-2 px-4 py-2.5 bg-orange-600 text-white text-sm font-bold rounded-xl border border-orange-500 hover:bg-orange-500 hover:-translate-y-px hover:shadow-lg disabled:opacity-60 disabled:cursor-not-allowed disabled:hover:translate-y-0 disabled:hover:shadow-none transition-all duration-200"
        >
          {#if isLoading}
            <span
              class="w-3.5 h-3.5 border-2 border-white/30 border-t-white rounded-full animate-spin flex-shrink-0"
              aria-hidden="true"
            ></span>
            Memproses...
          {:else}
            Daftar
          {/if}
        </button>
      </form>

      <!-- Login link -->
      <p class="mt-5 text-center text-xs text-orange-800 font-medium">
        Sudah punya akun?
        <a
          href="/login"
          class="text-orange-600 font-bold hover:text-orange-500 hover:underline transition-colors"
        >
          Masuk di sini
        </a>
      </p>
    </div>

    <!-- Back link -->
    <a
      href="/"
      class="group flex items-center gap-1.5 text-sm text-orange-700 hover:text-orange-900 font-medium transition-colors"
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

  <!-- Noise Overlay for texture -->
  <div
    class="absolute inset-0 opacity-[0.02] pointer-events-none mix-blend-screen"
    style="background-image: url('data:image/svg+xml,%3Csvg viewBox=%220 0 200 200%22 xmlns=%22http://www.w3.org/2000/svg%22%3E%3Cfilter id=%22noiseFilter%22%3E%3CfeTurbulence type=%22fractalNoise%22 baseFrequency=%220.8%22 numOctaves=%223%22 stitchTiles=%22stitch%22/%3E%3C/filter%3E%3Crect width=%22100%25%22 height=%22100%25%22 filter=%22url(%23noiseFilter)%22/%3E%3C/svg%3E');"
  ></div>
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
    -webkit-text-fill-color: #431407 !important; /* text-orange-950 */
    border-radius: 0.75rem; /* rounded-xl */
  }
</style>
