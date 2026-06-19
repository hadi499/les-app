<script lang="ts">
  import { onMount } from "svelte";

  let isAuthenticated = $state(false);
  let authChecked = $state(false);

  onMount(async () => {
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
  <title>Les Balongarut | Portal Akademik Eksklusif</title>
  <meta
    name="description"
    content="Sistem pendukung akademik terintegrasi untuk Les Balongarut."
  />
</svelte:head>

<div
  class="min-h-screen bg-orange-50 font-sans selection:bg-orange-200 selection:text-orange-900 flex flex-col relative overflow-x-hidden"
>
  <!-- Background Ambient -->
  <div class="absolute inset-0 z-0 pointer-events-none fixed">
    <div
      class="absolute top-1/4 left-1/4 w-[600px] h-[600px] bg-white/40 rounded-full blur-[120px]"
    ></div>
    <div
      class="absolute bottom-1/4 right-1/4 w-[500px] h-[500px] bg-amber-100/50 rounded-full blur-[120px]"
    ></div>

  </div>

  <!-- Hero Section -->
  <section
    class="relative z-10 w-full max-w-4xl mx-auto px-6 py-20 lg:py-32 flex flex-col items-center text-center gap-16"
  >
    <!-- Text -->
    <div class="flex flex-col items-center gap-8">
      <div class="flex items-center gap-3">
        <div
          class="w-8 h-8 border border-orange-200 bg-white flex items-center justify-center rotate-45 shadow-sm"
        >
          <div class="w-1.5 h-1.5 bg-orange-400"></div>
        </div>
        <span
          class="text-xs font-medium tracking-[0.3em] uppercase text-orange-800"
        >
          Sistem Internal
        </span>
      </div>

      <div class="flex flex-col gap-4 items-center">
        <h1
          class="text-4xl sm:text-5xl lg:text-[4rem] font-bold tracking-[0.1em] text-orange-950 uppercase leading-tight text-center drop-shadow-sm"
        >
          Les Balongarut
        </h1>
        <p
          class="text-[10px] sm:text-xs tracking-[0.2em] text-orange-800 uppercase flex flex-wrap justify-center gap-4 items-center font-medium"
        >
          <span>Komputer</span>
          <span class="w-1.5 h-1.5 bg-orange-300 rounded-full"></span>
          <span>Matematika</span>
          <span class="w-1.5 h-1.5 bg-orange-300 rounded-full"></span>
          <span>Bahasa Inggris</span>
        </p>
      </div>

      <div class="pt-4 flex justify-center">
        {#if authChecked && !isAuthenticated}
          <a
            href="/login"
            class="group relative inline-flex items-center justify-center px-8 py-3 text-xs tracking-[0.2em] font-bold uppercase text-orange-900 border border-orange-300 hover:text-orange-950 hover:border-orange-500 transition-all duration-700 bg-white/50 overflow-hidden cursor-pointer backdrop-blur-sm no-underline shadow-md hover:shadow-lg"
          >
            <span class="relative z-10 flex items-center gap-3">
              Inisiasi Sesi
              <svg
                class="w-3.5 h-3.5 transition-transform duration-500 group-hover:translate-x-1"
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
              class="absolute inset-0 -translate-x-full bg-orange-200/40 group-hover:translate-x-0 transition-transform duration-700 ease-out z-0"
            ></div>
          </a>
        {/if}
      </div>
    </div>
  </section>

  <!-- Spacer to push footer elements if any -->
  <div class="flex-1"></div>

  <!-- Noise Overlay for texture -->
  <div
    class="absolute inset-0 opacity-[0.02] pointer-events-none mix-blend-screen"
    style="background-image: url('data:image/svg+xml,%3Csvg viewBox=%220 0 200 200%22 xmlns=%22http://www.w3.org/2000/svg%22%3E%3Cfilter id=%22noiseFilter%22%3E%3CfeTurbulence type=%22fractalNoise%22 baseFrequency=%220.8%22 numOctaves=%223%22 stitchTiles=%22stitch%22/%3E%3C/filter%3E%3Crect width=%22100%25%22 height=%22100%25%22 filter=%22url(%23noiseFilter)%22/%3E%3C/svg%3E');"
  ></div>
</div>
