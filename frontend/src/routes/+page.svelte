<script lang="ts">
  import { onMount } from "svelte";

  let isAuthenticated = $state(false);
  let authChecked = $state(false);

  // Mengambil status dari server
  let isClassOpenPaud = $state(true);
  let isClassOpenSd = $state(true);
  let isSettingsLoaded = $state(false);

  type Quote = { id: number; quote: string; arti: string; author: string; created_at: string; updated_at: string };
  let quotes = $state<Quote[]>([]);
  let quotesLoaded = $state(false);

  onMount(async () => {
    // Fetch Quotes
    try {
      const resQuotes = await fetch("/api/quotes");
      if (resQuotes.ok) {
        quotes = await resQuotes.json();
      }
    } catch (e) {
      console.error("Gagal memuat quotes", e);
    } finally {
      quotesLoaded = true;
    }
    // Fetch Pengaturan Sistem
    try {
      const resSettings = await fetch("/api/settings");
      if (resSettings.ok) {
        const settings = await resSettings.json();
        if (settings.is_class_open_paud !== undefined) {
          isClassOpenPaud = settings.is_class_open_paud === "true";
        }
        if (settings.is_class_open_sd !== undefined) {
          isClassOpenSd = settings.is_class_open_sd === "true";
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
  <title>Les Balonggarut | Bimbingan Belajar Terbaik</title>
  <meta
    name="description"
    content="Les  Matematika, Bahasa Inggris, TKA dan Komputer  di Balong Garut Sidoarjo"
  />
  <meta
    name="keywords"
    content="Les Balong Garut, Balong Garut Sidoarjo, les SDN Balong Garut,  les TK Balong Garut, les privat, bimbingan belajar, les matematika, les komputer, les bahasa inggris, bimbel"
  />
  <link rel="canonical" href="https://lesbalonggarut.my.id/" />

  <!-- Open Graph -->
  <meta property="og:type" content="website" />
  <meta property="og:url" content="https://lesbalonggarut.my.id/" />
  <meta
    property="og:title"
    content="Les Balong Garut | Bimbingan Belajar Terbaik"
  />
  <meta
    property="og:description"
    content="Les Balong Garut menyediakan bimbingan belajar untuk Komputer, Matematika, Bahasa Inggris, dan TKA. Tingkatkan prestasimu sekarang!"
  />
  <meta property="og:locale" content="id_ID" />
  <meta property="og:site_name" content="Les Balong Garut" />

  <!-- Twitter Card -->
  <meta name="twitter:card" content="summary" />
  <meta
    name="twitter:title"
    content="Les Balong Garut | Bimbingan Belajar Terbaik"
  />
  <meta
    name="twitter:description"
    content="Bimbingan belajar Komputer, Matematika, Bahasa Inggris, dan TKA."
  />
</svelte:head>

<div
  class="min-h-screen bg-slate-50 font-sans selection:bg-blue-200 selection:text-blue-900 flex flex-col relative overflow-x-hidden text-slate-900"
>
  <!-- Background Ambient -->
  <div class="absolute inset-0 z-0 pointer-events-none fixed">
    <div
      class="absolute top-1/4 left-1/4 w-[250px] h-[250px] sm:w-[500px] sm:h-[500px] bg-blue-100/60 rounded-full blur-[80px] sm:blur-[120px]"
    ></div>
    <div
      class="absolute bottom-1/4 right-1/4 w-[200px] h-[200px] sm:w-[400px] sm:h-[400px] bg-slate-200/60 rounded-full blur-[80px] sm:blur-[120px]"
    ></div>
  </div>

  <!-- Hero Section -->
  <section
    class="relative z-10 w-full max-w-4xl mx-auto px-6 pt-28 pb-12 lg:pt-36 lg:pb-32 flex flex-col items-center text-center gap-16"
  >
    <!-- Text -->
    <div class="flex flex-col items-center gap-8">
      {#if isSettingsLoaded}
        <div class="flex flex-row flex-wrap justify-center gap-2 mt-4">
          <!-- Status PAUD/TK -->
          {#if isClassOpenPaud}
            <div
              class="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-emerald-50/80 backdrop-blur-sm border border-emerald-200/50 text-emerald-700 text-xs sm:text-sm font-bold tracking-widest uppercase animate-in fade-in zoom-in duration-500 shadow-sm shadow-emerald-500/10"
            >
              <span class="relative flex h-2.5 w-2.5">
                <span
                  class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"
                ></span>
                <span
                  class="relative inline-flex rounded-full h-2.5 w-2.5 bg-emerald-500"
                ></span>
              </span>
              PAUD/TK: Buka
            </div>
          {:else}
            <div
              class="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-rose-50/80 backdrop-blur-sm border border-rose-200/50 text-rose-700 text-xs sm:text-sm font-bold tracking-widest uppercase animate-in fade-in zoom-in duration-500 shadow-sm shadow-rose-500/10"
            >
              <svg
                class="w-4 h-4"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              PAUD/TK: Libur
            </div>
          {/if}

          <!-- Status SD/SMP -->
          {#if isClassOpenSd}
            <div
              class="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-emerald-50/80 backdrop-blur-sm border border-emerald-200/50 text-emerald-700 text-xs sm:text-sm font-bold tracking-widest uppercase animate-in fade-in zoom-in duration-500 shadow-sm shadow-emerald-500/10"
            >
              <span class="relative flex h-2.5 w-2.5">
                <span
                  class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"
                ></span>
                <span
                  class="relative inline-flex rounded-full h-2.5 w-2.5 bg-emerald-500"
                ></span>
              </span>
              SD/SMP: Buka
            </div>
          {:else}
            <div
              class="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-rose-50/80 backdrop-blur-sm border border-rose-200/50 text-rose-700 text-xs sm:text-sm font-bold tracking-widest uppercase animate-in fade-in zoom-in duration-500 shadow-sm shadow-rose-500/10"
            >
              <svg
                class="w-4 h-4"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              SD/SMP: Libur
            </div>
          {/if}
        </div>
      {:else}
        <!-- Placeholder agar layout tidak lompat -->
        <div class="h-9 sm:h-9 mt-4 opacity-0"></div>
      {/if}

      <div class="flex flex-col gap-6 items-center mt-2">

        <h1
          class="text-3xl sm:text-5xl lg:text-[4rem] font-bold tracking-[0.1em] text-slate-900 uppercase leading-tight text-center drop-shadow-sm"
        >
          Les Balonggarut
        </h1>

        <!-- Subject Badges -->
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-2 mt-1 w-full max-w-xs sm:max-w-none">
          <!-- Komputer -->
          <span class="inline-flex items-center justify-center gap-1.5 px-3 py-1.5 rounded-full bg-blue-50 border border-blue-200 text-blue-700 text-xs font-semibold tracking-wide">
            <svg class="w-3.5 h-3.5 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect x="2" y="3" width="20" height="14" rx="2"/>
              <path d="M8 21h8M12 17v4"/>
            </svg>
            Komputer
          </span>
          <!-- Matematika -->
          <span class="inline-flex items-center justify-center gap-1.5 px-3 py-1.5 rounded-full bg-violet-50 border border-violet-200 text-violet-700 text-xs font-semibold tracking-wide">
            <svg class="w-3.5 h-3.5 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M12 2v20M2 12h20"/>
            </svg>
            Matematika
          </span>
          <!-- Bahasa Inggris -->
          <span class="inline-flex items-center justify-center gap-1.5 px-3 py-1.5 rounded-full bg-emerald-50 border border-emerald-200 text-emerald-700 text-xs font-semibold tracking-wide">
            <svg class="w-3.5 h-3.5 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
            </svg>
            Bahasa Inggris
          </span>
          <!-- TKA -->
          <span class="inline-flex items-center justify-center gap-1.5 px-3 py-1.5 rounded-full bg-amber-50 border border-amber-200 text-amber-700 text-xs font-semibold tracking-wide">
            <svg class="w-3.5 h-3.5 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M9 11l3 3L22 4"/>
              <path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"/>
            </svg>
            TKA
          </span>
        </div>

        <div class="flex items-center gap-4 w-full max-w-xs opacity-30 mt-1">
          <div class="flex-1 h-px bg-slate-400"></div>
          <svg class="w-3 h-3 text-slate-400 shrink-0" viewBox="0 0 24 24" fill="currentColor">
            <circle cx="12" cy="12" r="3" />
          </svg>
          <div class="flex-1 h-px bg-slate-400"></div>
        </div>
        <p class="text-sm sm:text-base text-slate-500 font-light tracking-wide leading-7 text-center max-w-sm sm:max-w-md">
          Tempat belajar modern dengan metode
          <span class="font-semibold text-slate-700">Card Memory</span>
          <span class="whitespace-nowrap">dan <span class="font-semibold text-slate-700">latihan soal</span></span>
          terstruktur untuk hasil yang optimal.
        </p>
      </div>

      <div class="pt-6 flex flex-row items-center justify-center gap-3 flex-wrap">
        {#if authChecked}
          {#if !isAuthenticated}
            <!-- Primary: Masuk Portal -->
            <a
              href="/login"
              class="group inline-flex items-center gap-2.5 px-6 py-3 rounded-full bg-blue-600 hover:bg-blue-700 text-white text-xs font-bold tracking-[0.15em] uppercase shadow-lg shadow-blue-500/30 hover:shadow-blue-500/50 hover:scale-105 transition-all duration-300 no-underline"
            >
              <svg class="w-3.5 h-3.5 transition-transform duration-300 group-hover:translate-x-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1"/>
              </svg>
              Masuk Portal
            </a>
          {:else}
            <!-- Primary: Dashboard -->
            <a
              href="/dashboard"
              class="group inline-flex items-center gap-2.5 px-6 py-3 rounded-full bg-blue-600 hover:bg-blue-700 text-white text-xs font-bold tracking-[0.15em] uppercase shadow-lg shadow-blue-500/30 hover:shadow-blue-500/50 hover:scale-105 transition-all duration-300 no-underline"
            >
              <svg class="w-3.5 h-3.5 transition-transform duration-300 group-hover:translate-x-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"/>
              </svg>
              Dashboard
            </a>
          {/if}
        {/if}

        <!-- Secondary: Panduan -->
        <a
          href="/panduan"
          class="group inline-flex items-center gap-2.5 px-6 py-3 rounded-full bg-white/70 hover:bg-white border border-slate-300 hover:border-slate-400 text-slate-700 hover:text-slate-900 text-xs font-semibold tracking-[0.15em] uppercase backdrop-blur-sm shadow-sm hover:shadow-md transition-all duration-300 no-underline"
        >
          Panduan
          <svg class="w-3.5 h-3.5 transition-transform duration-300 group-hover:translate-x-0.5 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3"/>
          </svg>
        </a>
      </div>
    </div>
  </section>

  <!-- Quotes Section -->
  {#if quotesLoaded && quotes.length > 0}
  <section class="relative z-10 w-full max-w-5xl mx-auto px-6 pb-24 flex flex-col items-center gap-8 mt-4">
    <!-- Section Label -->
    <div class="flex items-center gap-3 text-slate-400">
      <div class="h-px w-8 bg-slate-300"></div>
      <span class="text-xs tracking-[0.25em] uppercase font-medium">Inspirasi</span>
      <div class="h-px w-8 bg-slate-300"></div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-5 w-full">
      {#each quotes as q, i}
        <div class="{
          i % 3 === 0 ? 'bg-blue-50/80 border-l-blue-500' :
          i % 3 === 1 ? 'bg-violet-50/80 border-l-violet-500' :
                        'bg-emerald-50/80 border-l-emerald-500'
        } backdrop-blur-md p-6 rounded-2xl shadow-sm border border-slate-200/50 border-l-4 flex flex-col justify-between hover:shadow-lg hover:-translate-y-0.5 transition-all duration-300 group">
          <div>
            <!-- Big quote icon -->
            <svg class="w-9 h-9 opacity-25 mb-4 group-hover:opacity-40 transition-opacity duration-300 {
              i % 3 === 0 ? 'text-blue-600' :
              i % 3 === 1 ? 'text-violet-600' :
                            'text-emerald-600'
            }" fill="currentColor" viewBox="0 0 24 24">
              <path d="M14.017 21v-7.391c0-5.704 3.731-9.57 8.983-10.609l.995 2.151c-2.432.917-3.995 3.638-3.995 5.849h4v10h-9.983zm-14.017 0v-7.391c0-5.704 3.748-9.57 9-10.609l.996 2.151c-2.433.917-3.996 3.638-3.996 5.849h3.983v10h-9.983z" />
            </svg>
            <!-- English quote -->
            <p class="text-slate-900 font-semibold text-sm leading-relaxed mb-3 italic">&ldquo;{q.quote}&rdquo;</p>
            <!-- Indonesian translation -->
            <p class="text-slate-500 text-xs leading-relaxed mb-6 border-t border-slate-200/70 pt-3">{q.arti}</p>
          </div>
          <!-- Author -->
          <div class="flex items-center gap-2">
            <div class="w-5 h-0.5 rounded-full {
              i % 3 === 0 ? 'bg-blue-500' :
              i % 3 === 1 ? 'bg-violet-500' :
                            'bg-emerald-500'
            }"></div>
            <p class="text-slate-700 font-bold text-xs tracking-wider uppercase">{q.author}</p>
          </div>
        </div>
      {/each}
    </div>
  </section>
  {/if}

  <!-- Footer -->
  <footer class="relative z-10 w-full border-t border-slate-200/70 bg-white/40 backdrop-blur-sm">
    <div class="max-w-5xl mx-auto px-6 py-8 flex flex-col sm:flex-row items-center justify-between gap-4">
      <!-- Brand -->
      <div class="flex flex-col items-center sm:items-start gap-1">
        <span class="text-sm font-bold tracking-[0.15em] uppercase text-slate-800">Les Balonggarut</span>
        <span class="text-xs text-slate-500 flex items-center gap-1.5">
          <svg class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/>
            <circle cx="12" cy="10" r="3"/>
          </svg>
          Balong Garut, Krembung, Sidoarjo
        </span>
      </div>

      <!-- Nav Links -->
      <nav class="flex items-center gap-5 text-xs text-slate-500 font-medium">
        <a href="/panduan" class="hover:text-slate-800 transition-colors duration-200 no-underline">Panduan</a>
        <a href="/quiz" class="hover:text-slate-800 transition-colors duration-200 no-underline">Kuis</a>
        <a href="/login" class="hover:text-slate-800 transition-colors duration-200 no-underline">Masuk</a>
      </nav>

      <!-- Copyright -->
      <span class="text-xs text-slate-400">&copy; {new Date().getFullYear()} Les Balonggarut</span>
    </div>
  </footer>

  <!-- Noise Overlay for texture -->
  <div
    class="absolute inset-0 opacity-[0.03] pointer-events-none mix-blend-screen"
    style="background-image: url('data:image/svg+xml,%3Csvg viewBox=%220 0 200 200%22 xmlns=%22http://www.w3.org/2000/svg%22%3E%3Cfilter id=%22noiseFilter%22%3E%3CfeTurbulence type=%22fractalNoise%22 baseFrequency=%220.8%22 numOctaves=%223%22 stitchTiles=%22stitch%22/%3E%3C/filter%3E%3Crect width=%22100%25%22 height=%22100%25%22 filter=%22url(%23noiseFilter)%22/%3E%3C/svg%3E');"
  ></div>
</div>
