<script lang="ts">
  import { onMount, getContext } from "svelte";

  let authState = getContext<{
    isAuthenticated: boolean;
    authChecked: boolean;
  }>("authState");
  let isAuthenticated = $derived(authState.isAuthenticated);
  let authChecked = $derived(authState.authChecked);

  type Quote = {
    id: number;
    quote: string;
    arti: string;
    author: string;
    created_at: string;
    updated_at: string;
  };
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
  <div class="absolute inset-0 z-0 pointer-events-none overflow-hidden">
    <div
      class="absolute top-1/4 left-1/4 w-[250px] h-[250px] sm:w-[500px] sm:h-[500px] bg-blue-100/60 rounded-full blur-[80px] sm:blur-[120px]"
    ></div>
    <div
      class="absolute bottom-1/4 right-1/4 w-[200px] h-[200px] sm:w-[400px] sm:h-[400px] bg-slate-200/60 rounded-full blur-[80px] sm:blur-[120px]"
    ></div>
  </div>

  <!-- Hero Section -->
  <section
    class="relative z-10 w-full max-w-4xl mx-auto px-6 pt-20 pb-12 lg:pt-28 lg:pb-32 flex flex-col items-center text-center gap-16"
  >
    <!-- Text -->
    <div class="flex flex-col items-center gap-8">
      <div class="flex flex-col gap-6 items-center mt-8 sm:mt-2">
        <h1
          class="text-3xl font-['Concert_One'] sm:text-5xl lg:text-[4.5rem] font-bold tracking-normal sm:tracking-[0.1em] whitespace-nowrap bg-gradient-to-r from-blue-600 to-slate-900 text-transparent bg-clip-text uppercase leading-tight text-center drop-shadow-sm"
        >
          Les Balonggarut
        </h1>

        <!-- Subject Badges -->
        <div
          class="flex flex-col sm:grid sm:grid-cols-4 gap-3 sm:gap-2 mt-1 w-fit min-w-50 sm:w-full sm:max-w-none"
        >
          <!-- Komputer -->
          <span
            class="inline-flex items-center justify-center gap-1.5 px-3 py-1.5 rounded-full bg-blue-50 border border-blue-200 text-blue-700 text-xs font-semibold tracking-wide"
          >
            <svg
              class="w-3.5 h-3.5 shrink-0"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <rect x="2" y="3" width="20" height="14" rx="2" />
              <path d="M8 21h8M12 17v4" />
            </svg>
            Komputer
          </span>
          <!-- Matematika -->
          <span
            class="inline-flex items-center justify-center gap-1.5 px-3 py-1.5 rounded-full bg-violet-50 border border-violet-200 text-violet-700 text-xs font-semibold tracking-wide"
          >
            <svg
              class="w-3.5 h-3.5 shrink-0"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M12 2v20M2 12h20" />
            </svg>
            Matematika
          </span>
          <!-- Bahasa Inggris -->
          <span
            class="inline-flex items-center justify-center gap-1.5 px-3 py-1.5 rounded-full bg-emerald-50 border border-emerald-200 text-emerald-700 text-xs font-semibold tracking-wide"
          >
            <svg
              class="w-3.5 h-3.5 shrink-0"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path
                d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"
              />
            </svg>
            Bahasa Inggris
          </span>
          <!-- TKA -->
          <span
            class="inline-flex items-center justify-center gap-1.5 px-3 py-1.5 rounded-full bg-amber-50 border border-amber-200 text-amber-700 text-xs font-semibold tracking-wide"
          >
            <svg
              class="w-3.5 h-3.5 shrink-0"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M9 11l3 3L22 4" />
              <path
                d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"
              />
            </svg>
            TKA
          </span>
        </div>

        <div class="flex items-center gap-4 w-full max-w-xs opacity-30 mt-1">
          <div class="flex-1 h-px bg-slate-400"></div>
          <svg
            class="w-3 h-3 text-slate-400 shrink-0"
            viewBox="0 0 24 24"
            fill="currentColor"
          >
            <circle cx="12" cy="12" r="3" />
          </svg>
          <div class="flex-1 h-px bg-slate-400"></div>
        </div>
        <p
          class="text-sm sm:text-base text-slate-500 font-light tracking-wide leading-7 text-center max-w-sm sm:max-w-md"
        >
          Tempat belajar modern dengan metode
          <span class="font-semibold text-slate-700">Card Memory</span>
          <span class="whitespace-nowrap"
            >dan <span class="font-semibold text-slate-700">latihan soal</span
            ></span
          >
          terstruktur untuk hasil yang optimal.
        </p>
      </div>

      <div class="pt-6 flex flex-col items-center justify-center gap-8 w-full">
        <!-- Main Actions -->
        <div class="flex flex-row items-center justify-center gap-5 flex-wrap">
          {#if authChecked}
            {#if !isAuthenticated}
              <!-- Secondary: Masuk Portal -->
              <a
                href="/login"
                class="group inline-flex items-center justify-center w-48 gap-2.5 px-6 py-3 rounded-full bg-white/70 hover:bg-white border border-indigo-400 hover:border-indigo-500 text-indigo-600 hover:text-indigo-700 text-xs font-bold tracking-[0.15em] uppercase backdrop-blur-sm shadow-md hover:shadow-lg hover:scale-105 transition-all duration-300 no-underline"
              >
                <svg
                  class="w-3.5 h-3.5 transition-transform duration-300 group-hover:translate-x-0.5"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M12 6.042A8.967 8.967 0 006 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 016 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 016-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0018 18a8.967 8.967 0 00-6 2.292m0-14.25v14.25"
                  />
                </svg>
                Les App
              </a>
            {:else}
              <!-- Primary: Dashboard -->
              <a
                href="/dashboard"
                class="group inline-flex items-center justify-center w-48 gap-2.5 px-6 py-3 rounded-full bg-white/70 hover:bg-white border border-indigo-400 hover:border-indigo-500 text-indigo-600 hover:text-indigo-700 text-xs font-bold tracking-[0.15em] uppercase backdrop-blur-sm shadow-md hover:shadow-lg hover:scale-105 transition-all duration-300 no-underline"
              >
                <svg
                  class="w-3.5 h-3.5 transition-transform duration-300 group-hover:translate-x-0.5"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
                  />
                </svg>
                Dashboard
              </a>
            {/if}

            <!-- Kuis App -->
            <a
              href="/quiz-app"
              class="group inline-flex items-center justify-center w-48 gap-2.5 px-6 py-3 rounded-full bg-white/70 hover:bg-white border border-indigo-400 hover:border-indigo-500 text-indigo-600 hover:text-indigo-700 text-xs font-bold tracking-[0.15em] uppercase backdrop-blur-sm shadow-md hover:shadow-lg hover:scale-105 transition-all duration-300 no-underline"
            >
              Kuis App
              <svg
                class="w-3.5 h-3.5 transition-transform duration-300 group-hover:translate-x-0.5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M13 10V3L4 14h7v7l9-11h-7z"
                />
              </svg>
            </a>
          {:else}
            <!-- Placeholders to prevent layout shift while checking auth -->
            <div
              class="inline-flex items-center justify-center w-48 gap-2.5 px-6 py-3 rounded-full bg-white/50 border border-indigo-200 backdrop-blur-sm shadow-md animate-pulse"
            >
              <div class="w-3.5 h-3.5 bg-indigo-200 rounded-full"></div>
              <div class="h-3 w-20 bg-indigo-200 rounded"></div>
            </div>
            
            <div
              class="inline-flex items-center justify-center w-48 gap-2.5 px-6 py-3 rounded-full bg-white/50 border border-indigo-200 backdrop-blur-sm shadow-md animate-pulse"
            >
              <div class="h-3 w-16 bg-indigo-200 rounded"></div>
              <div class="w-3.5 h-3.5 bg-indigo-200 rounded-full"></div>
            </div>
          {/if}
        </div>

        <!-- Featured New Feature: Belajar Coding & Berhitung -->
        <div
          class="flex flex-col items-center gap-6 w-full animate-in fade-in slide-in-from-bottom-2 duration-500 delay-150 mt-4"
        >
          <div
            class="flex items-center gap-4 w-full max-w-[300px] opacity-60 mb-2"
          >
            <div class="flex-1 h-px bg-slate-300"></div>
            <span
              class="text-[10px] uppercase tracking-[0.2em] text-slate-500 font-bold whitespace-nowrap"
              >Fitur Pendukung</span
            >
            <div class="flex-1 h-px bg-slate-300"></div>
          </div>

          <div
            class="grid grid-cols-1 md:grid-cols-2 gap-6 w-full max-w-3xl text-left"
          >
            <!-- Card Belajar Coding -->
            <a
              href="/belajar-coding"
              class="group relative flex flex-col items-start p-6 bg-white/70 backdrop-blur-sm border border-slate-300 hover:border-blue-400 hover:bg-white rounded-3xl hover:shadow-xl hover:shadow-blue-500/10 transition-all duration-300 no-underline overflow-hidden"
            >
              <div
                class="absolute -right-6 -top-6 w-32 h-32 bg-blue-100 rounded-full blur-2xl opacity-40 group-hover:opacity-70 transition-opacity duration-500"
              ></div>

              <div class="flex items-center gap-4 mb-4 relative z-10 w-full">
                <div
                  class="w-12 h-12 shrink-0 rounded-2xl bg-blue-100 text-blue-600 flex items-center justify-center shadow-sm border border-blue-200/50"
                >
                  <svg
                    class="w-6 h-6"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    ><path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"
                    /></svg
                  >
                </div>
                <h3
                  class="text-lg font-bold text-slate-800 group-hover:text-blue-600 transition-colors"
                >
                  Belajar Coding
                </h3>
              </div>
              <p
                class="text-sm text-slate-500 leading-relaxed relative z-10 flex-1"
              >
                Latih logika komputasional anak hingga remaja dengan merakit
                blok visual maupun menulis kode secara interaktif.
              </p>

              <div
                class="mt-6 flex items-center text-sm font-bold text-blue-600 opacity-100 md:opacity-0 group-hover:opacity-100 transform translate-x-0 md:-translate-x-2.5 group-hover:translate-x-0 transition-all duration-300 relative z-10"
              >
                Coba Sekarang
                <svg
                  class="w-4 h-4 ml-1"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  ><path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M14 5l7 7m0 0l-7 7m7-7H3"
                  ></path></svg
                >
              </div>
            </a>

            <!-- Card Berhitung -->
            <a
              href="/berhitung"
              class="group relative flex flex-col items-start p-6 bg-white/70 backdrop-blur-sm border border-slate-300 hover:border-pink-400 hover:bg-white rounded-3xl hover:shadow-xl hover:shadow-pink-500/10 transition-all duration-300 no-underline overflow-hidden"
            >
              <div
                class="absolute -right-6 -top-6 w-32 h-32 bg-pink-100 rounded-full blur-2xl opacity-40 group-hover:opacity-70 transition-opacity duration-500"
              ></div>

              <div class="flex items-center gap-4 mb-4 relative z-10 w-full">
                <div
                  class="w-12 h-12 shrink-0 rounded-2xl bg-pink-100 text-pink-600 flex items-center justify-center shadow-sm border border-pink-200/50"
                >
                  <svg
                    class="w-6 h-6"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    ><path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z"
                    /></svg
                  >
                </div>
                <h3
                  class="text-lg font-bold text-slate-800 group-hover:text-pink-600 transition-colors"
                >
                  Latihan Berhitung
                </h3>
              </div>
              <p
                class="text-sm text-slate-500 leading-relaxed relative z-10 flex-1"
              >
                Tingkatkan kecepatan dan ketepatan menghitung dengan modul
                aritmetika yang terstruktur untuk semua level.
              </p>

              <div
                class="mt-6 flex items-center text-sm font-bold text-pink-600 opacity-100 md:opacity-0 group-hover:opacity-100 transform translate-x-0 md:-translate-x-2.5 group-hover:translate-x-0 transition-all duration-300 relative z-10"
              >
                Mulai Berlatih
                <svg
                  class="w-4 h-4 ml-1"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  ><path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M14 5l7 7m0 0l-7 7m7-7H3"
                  ></path></svg
                >
              </div>
            </a>

            <!-- Card Ketik 10 Jari -->
            <a
              href="/mengetik"
              class="group relative flex flex-col items-start p-6 bg-white/70 backdrop-blur-sm border border-slate-300 hover:border-teal-400 hover:bg-white rounded-3xl hover:shadow-xl hover:shadow-teal-500/10 transition-all duration-300 no-underline overflow-hidden"
            >
              <div
                class="absolute -right-6 -top-6 w-32 h-32 bg-teal-100 rounded-full blur-2xl opacity-40 group-hover:opacity-70 transition-opacity duration-500"
              ></div>

              <div class="flex items-center gap-4 mb-4 relative z-10 w-full">
                <div
                  class="w-12 h-12 shrink-0 rounded-2xl bg-teal-100 text-teal-600 flex items-center justify-center shadow-sm border border-teal-200/50"
                >
                  <svg
                    class="w-6 h-6"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    ><path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z"
                    ></path></svg
                  >
                </div>
                <h3
                  class="text-lg font-bold text-slate-800 group-hover:text-teal-600 transition-colors"
                >
                  Ketik 10 Jari
                </h3>
              </div>
              <p
                class="text-sm text-slate-500 leading-relaxed relative z-10 flex-1"
              >
                Latih kecepatan dan akurasi mengetik dengan metode 10 jari yang
                efisien dan interaktif.
              </p>

              <div
                class="mt-6 flex items-center text-sm font-bold text-teal-600 opacity-100 md:opacity-0 group-hover:opacity-100 transform translate-x-0 md:-translate-x-2.5 group-hover:translate-x-0 transition-all duration-300 relative z-10"
              >
                Mulai Berlatih
                <svg
                  class="w-4 h-4 ml-1"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  ><path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M14 5l7 7m0 0l-7 7m7-7H3"
                  ></path></svg
                >
              </div>
            </a>

            <!-- Card Compress Image -->
            <a
              href="/compress-image"
              class="group relative flex flex-col items-start p-6 bg-white/70 backdrop-blur-sm border border-slate-300 hover:border-orange-400 hover:bg-white rounded-3xl hover:shadow-xl hover:shadow-orange-500/10 transition-all duration-300 no-underline overflow-hidden"
            >
              <div
                class="absolute -right-6 -top-6 w-32 h-32 bg-orange-100 rounded-full blur-2xl opacity-40 group-hover:opacity-70 transition-opacity duration-500"
              ></div>

              <div class="flex items-center gap-4 mb-4 relative z-10 w-full">
                <div
                  class="w-12 h-12 shrink-0 rounded-2xl bg-orange-100 text-orange-600 flex items-center justify-center shadow-sm border border-orange-200/50"
                >
                  <svg
                    class="w-6 h-6"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    ><path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
                    ></path></svg
                  >
                </div>
                <h3
                  class="text-lg font-bold text-slate-800 group-hover:text-orange-600 transition-colors"
                >
                  Compress Image
                </h3>
              </div>
              <p
                class="text-sm text-slate-500 leading-relaxed relative z-10 flex-1"
              >
                Perkecil ukuran foto atau gambar Anda dengan cepat tanpa
                mengurangi kualitas secara drastis.
              </p>

              <div
                class="mt-6 flex items-center text-sm font-bold text-orange-600 opacity-100 md:opacity-0 group-hover:opacity-100 transform translate-x-0 md:-translate-x-2.5 group-hover:translate-x-0 transition-all duration-300 relative z-10"
              >
                Buka Alat
                <svg
                  class="w-4 h-4 ml-1"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  ><path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M14 5l7 7m0 0l-7 7m7-7H3"
                  ></path></svg
                >
              </div>
            </a>
          </div>
        </div>
      </div>
    </div>
  </section>

  <!-- Quotes Section -->
  {#if quotesLoaded && quotes.length > 0}
    <section
      class="relative z-10 w-full max-w-5xl mx-auto px-6 pb-24 flex flex-col items-center gap-8 mt-4"
    >
      <!-- Section Label -->
      <div class="flex items-center gap-3 text-slate-400">
        <div class="h-px w-8 bg-slate-300"></div>
        <span class="text-xs tracking-[0.25em] uppercase font-medium"
          >Inspirasi</span
        >
        <div class="h-px w-8 bg-slate-300"></div>
      </div>

      <div class="grid grid-cols-1 max-w-2xl mx-auto w-full">
        {#each quotes.slice(0, 1) as q, i}
          <div
            class="flex flex-col items-center text-center px-4 py-8 md:py-10"
          >
            <!-- English quote -->
            <p
              class="text-slate-800 font-semibold text-lg md:text-xl leading-relaxed mb-4 italic max-w-xl"
            >
              &ldquo;{q.quote}&rdquo;
            </p>
            <!-- Indonesian translation -->
            <p
              class="text-slate-500 text-sm md:text-base leading-relaxed mb-8 max-w-lg"
            >
              {q.arti}
            </p>
            <!-- Author -->
            <div class="flex flex-col items-center gap-2">
              <div class="w-8 h-0.5 rounded-full bg-blue-500/50"></div>
              <p
                class="text-slate-700 font-bold text-xs tracking-[0.2em] uppercase mt-2"
              >
                {q.author}
              </p>
            </div>
          </div>
        {/each}
      </div>
    </section>
  {/if}

  <!-- Footer -->
  <footer
    class="relative z-10 w-full border-t border-slate-200/70 bg-white/40 backdrop-blur-sm"
  >
    <div
      class="max-w-5xl mx-auto px-6 py-8 flex flex-col sm:flex-row items-center justify-between gap-4"
    >
      <!-- Brand -->
      <div class="flex flex-col items-center sm:items-start gap-1">
        <span
          class="text-sm font-bold tracking-[0.15em] uppercase text-slate-800"
          >Les Balonggarut</span
        >
        <span class="text-xs text-slate-500 flex items-center gap-1.5">
          <svg
            class="w-3 h-3"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z" />
            <circle cx="12" cy="10" r="3" />
          </svg>
          Balong Garut, Krembung, Sidoarjo
        </span>
      </div>

      <!-- Nav Links -->
      <nav class="flex items-center gap-5 text-xs text-slate-500 font-medium">
        <a
          href="/panduan"
          class="hover:text-slate-800 transition-colors duration-200 no-underline"
          >Panduan</a
        >
        <a
          href="/belajar-coding"
          class="hover:text-slate-800 transition-colors duration-200 no-underline"
          >Belajar Coding</a
        >

        <a
          href="/login"
          class="hover:text-slate-800 transition-colors duration-200 no-underline"
          >Masuk</a
        >
      </nav>

      <!-- Copyright -->
      <span class="text-xs text-slate-400"
        >&copy; {new Date().getFullYear()} Les Balonggarut</span
      >
    </div>
  </footer>

  <!-- Noise Overlay for texture removed for better mobile performance -->
</div>
