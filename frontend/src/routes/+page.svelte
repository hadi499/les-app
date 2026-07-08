<script lang="ts">
  import { onMount } from "svelte";

  let isAuthenticated = $state(false);
  let authChecked = $state(false);
  
  // Mengambil status dari server
  let isClassOpen = $state(true);
  let isSettingsLoaded = $state(false);

  onMount(async () => {
    // Fetch Pengaturan Sistem
    try {
      const resSettings = await fetch("/api/settings");
      if (resSettings.ok) {
        const settings = await resSettings.json();
        if (settings.is_class_open !== undefined) {
          isClassOpen = settings.is_class_open === "true";
        }
      }
    } catch (e) {
      console.error("Gagal memuat pengaturan", e);
    } finally {
      isSettingsLoaded = true;
    }

    try {
      const res = await fetch(`/me`, {
        credentials: "include",
      });
      if (res.ok) {
        const data = (await res.json()) as { authenticated: boolean };
        if (data.authenticated) {
          isAuthenticated = true;
        }
      }
    } catch (e) {
      // Ignore error if no session
    } finally {
      authChecked = true;
    }
  });
</script>

<svelte:head>
  <title>Les Balongarut | Bimbingan Belajar Terbaik</title>
  <meta
    name="description"
    content="Les Balongarut menyediakan bimbingan belajar untuk Komputer, Matematika, Bahasa Inggris, dan TKA dengan pengajar profesional. Tingkatkan prestasimu sekarang!"
  />
  <meta name="keywords" content="Les Balongarut, les privat, bimbingan belajar, les matematika, les komputer, les bahasa inggris, bimbel" />
</svelte:head>

<div
  class="min-h-screen bg-slate-50 font-sans selection:bg-blue-200 selection:text-blue-900 flex flex-col relative overflow-x-hidden text-slate-900"
>
  <!-- Background Ambient -->
  <div class="absolute inset-0 z-0 pointer-events-none fixed">
    <div
      class="absolute top-1/4 left-1/4 w-[600px] h-[600px] bg-blue-100/60 rounded-full blur-[120px]"
    ></div>
    <div
      class="absolute bottom-1/4 right-1/4 w-[500px] h-[500px] bg-slate-200/60 rounded-full blur-[120px]"
    ></div>

  </div>

  <!-- Hero Section -->
  <section
    class="relative z-10 w-full max-w-4xl mx-auto px-6 py-20 lg:py-32 flex flex-col items-center text-center gap-16"
  >
    <!-- Text -->
    <div class="flex flex-col items-center gap-8">


      {#if isSettingsLoaded}
        {#if isClassOpen}
          <div class="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-emerald-50/80 backdrop-blur-sm border border-emerald-200/50 text-emerald-700 text-xs sm:text-sm font-bold tracking-widest uppercase animate-in fade-in zoom-in duration-500 shadow-sm shadow-emerald-500/10 mt-4">
            <span class="relative flex h-2.5 w-2.5">
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
              <span class="relative inline-flex rounded-full h-2.5 w-2.5 bg-emerald-500"></span>
            </span>
            Kelas Hari Ini Buka
          </div>
        {:else}
          <div class="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-rose-50/80 backdrop-blur-sm border border-rose-200/50 text-rose-700 text-xs sm:text-sm font-bold tracking-widest uppercase animate-in fade-in zoom-in duration-500 shadow-sm shadow-rose-500/10 mt-4">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            Kelas Sedang Libur
          </div>
        {/if}
      {:else}
        <!-- Placeholder agar layout tidak lompat -->
        <div class="h-9 sm:h-9 mt-4 opacity-0"></div>
      {/if}

      <div class="flex flex-col gap-4 items-center mt-2">
        <h1
          class="text-4xl sm:text-5xl lg:text-[4rem] font-bold tracking-[0.1em] text-slate-900 uppercase leading-tight text-center drop-shadow-sm"
        >
          Les Balongarut
        </h1>
        <div
          class="text-xs sm:text-sm tracking-[0.2em] text-slate-600 uppercase flex flex-col items-center gap-3 font-medium"
        >
          <span>Komputer</span>
          <span>Matematika</span>
          <span>Bahasa Inggris</span>
          <span>TKA</span>
        </div>
      </div>

      <div class="pt-6 flex flex-col items-center gap-6">
        <a
          href="/panduan"
          class="w-64 group relative inline-flex items-center justify-center px-8 py-3 text-xs tracking-[0.2em] font-bold uppercase text-slate-800 border border-slate-300 hover:text-slate-950 hover:border-blue-500 transition-all duration-700 bg-white/60 overflow-hidden cursor-pointer backdrop-blur-md no-underline shadow-sm hover:shadow-md hover:shadow-blue-500/10"
        >
          <span class="relative z-10 flex items-center gap-3">
            Panduan
            <svg
              class="w-3.5 h-3.5 transition-transform duration-500 group-hover:translate-x-1 text-blue-600"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M17 8l4 4m0 0l-4 4m4-4H3"
              ></path>
            </svg>
          </span>
          <div
            class="absolute inset-0 -translate-x-full bg-blue-100/50 group-hover:translate-x-0 transition-transform duration-700 ease-out z-0"
          ></div>
        </a>

        {#if authChecked}
          {#if !isAuthenticated}
            <a
              href="/login"
              class="w-64 group relative inline-flex items-center justify-center px-8 py-3 text-xs tracking-[0.2em] font-bold uppercase text-white border border-blue-600 hover:border-blue-700 transition-all duration-700 bg-blue-600 hover:bg-blue-700 overflow-hidden cursor-pointer backdrop-blur-md no-underline shadow-md hover:shadow-lg"
            >
              <span class="relative z-10 flex items-center gap-3">
                Masuk Portal
                <svg
                  class="w-3.5 h-3.5 transition-transform duration-500 group-hover:translate-x-1 text-white"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1"
                  ></path>
                </svg>
              </span>
              <div
                class="absolute inset-0 -translate-x-full bg-blue-700/50 group-hover:translate-x-0 transition-transform duration-700 ease-out z-0"
              ></div>
            </a>
          {:else}
            <a
              href="/dashboard"
              class="w-64 group relative inline-flex items-center justify-center px-8 py-3 text-xs tracking-[0.2em] font-bold uppercase text-white border border-blue-600 hover:border-blue-700 transition-all duration-700 bg-blue-600 hover:bg-blue-700 overflow-hidden cursor-pointer backdrop-blur-md no-underline shadow-md hover:shadow-lg"
            >
              <span class="relative z-10 flex items-center gap-3">
                Dashboard
                <svg
                  class="w-3.5 h-3.5 transition-transform duration-500 group-hover:translate-x-1 text-white"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
                  ></path>
                </svg>
              </span>
              <div
                class="absolute inset-0 -translate-x-full bg-blue-700/50 group-hover:translate-x-0 transition-transform duration-700 ease-out z-0"
              ></div>
            </a>
          {/if}
        {/if}
      </div>
    </div>
  </section>

  <!-- Spacer to push footer elements if any -->
  <div class="flex-1"></div>

  <!-- Noise Overlay for texture -->
  <div
    class="absolute inset-0 opacity-[0.03] pointer-events-none mix-blend-screen"
    style="background-image: url('data:image/svg+xml,%3Csvg viewBox=%220 0 200 200%22 xmlns=%22http://www.w3.org/2000/svg%22%3E%3Cfilter id=%22noiseFilter%22%3E%3CfeTurbulence type=%22fractalNoise%22 baseFrequency=%220.8%22 numOctaves=%223%22 stitchTiles=%22stitch%22/%3E%3C/filter%3E%3Crect width=%22100%25%22 height=%22100%25%22 filter=%22url(%23noiseFilter)%22/%3E%3C/svg%3E');"
  ></div>
</div>
