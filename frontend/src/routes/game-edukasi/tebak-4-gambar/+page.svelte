<script>
  import Navbar from "$lib/components/Navbar.svelte";
</script>

<svelte:head>
  <title>Tebak 4 Gambar - Game Edukasi</title>
</svelte:head>

<div class="min-h-screen bg-slate-50 pb-12 font-sans selection:bg-blue-200">
  <Navbar />

  <!-- Hero Section -->
  <div class="relative pt-24 pb-16 overflow-hidden">
    <div class="absolute inset-0 bg-linear-to-b from-blue-100/50 to-transparent"></div>
    <div class="absolute -right-20 -top-20 w-96 h-96 bg-blue-200/40 rounded-full blur-3xl"></div>
    <div class="absolute -left-20 top-20 w-72 h-72 bg-indigo-200/40 rounded-full blur-3xl"></div>

    <div class="max-w-6xl mx-auto px-6 sm:px-8 relative z-10 text-center">
      <h1 class="text-4xl md:text-5xl lg:text-6xl font-extrabold text-slate-800 tracking-tight mb-4">
        Tebak <span class="text-blue-600">4 Gambar</span>
      </h1>
      <p class="text-lg text-slate-600 max-w-2xl mx-auto leading-relaxed">
        Pilih tema permainan di bawah ini dan latih daya ingat tingkat mudah dengan menebak posisi 4 gambar.
      </p>
      
      <!-- Back button to Game Edukasi -->
      <a href="/game-edukasi" class="inline-flex items-center justify-center mt-6 px-6 py-2.5 bg-white border border-slate-200 text-slate-600 rounded-full hover:bg-slate-50 hover:text-blue-600 transition-colors shadow-sm font-medium">
        <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path></svg>
        Kembali ke Kumpulan Game
      </a>
    </div>
  </div>

  <!-- Cards Section -->
  <div class="max-w-6xl mx-auto px-6 sm:px-8">
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 sm:gap-8">
      
      <!-- Card Doraemon -->
      <a
        href="/game-edukasi/tebak-4-gambar/doraemon"
        class="group relative flex flex-col items-start p-6 bg-white/70 backdrop-blur-sm border border-slate-300 hover:border-yellow-400 hover:bg-white rounded-3xl hover:shadow-xl hover:shadow-yellow-500/10 transition-all duration-300 no-underline overflow-hidden"
      >
        <div class="absolute -right-6 -top-6 w-32 h-32 bg-yellow-100 rounded-full blur-2xl opacity-40 group-hover:opacity-70 transition-opacity duration-500"></div>
        <div class="flex items-center gap-4 mb-4 relative z-10 w-full">
          <div class="w-12 h-12 shrink-0 rounded-2xl bg-yellow-100 text-yellow-600 flex items-center justify-center shadow-sm border border-yellow-200/50">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
          </div>
          <h3 class="text-lg font-bold text-slate-800 group-hover:text-yellow-600 transition-colors">Doraemon</h3>
        </div>
        <p class="text-sm text-slate-500 leading-relaxed relative z-10 flex-1">Bermain tebak posisi gambar karakter favorit (Mudah).</p>
        <div class="mt-6 flex items-center text-sm font-bold text-yellow-600 opacity-100 md:opacity-0 group-hover:opacity-100 transform translate-x-0 md:-translate-x-2.5 group-hover:translate-x-0 transition-all duration-300 relative z-10">
          Mulai Bermain
          <svg class="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"></path></svg>
        </div>
      </a>

      <!-- Card Tayo -->
      <a
        href="/game-edukasi/tebak-4-gambar/tayo"
        class="group relative flex flex-col items-start p-6 bg-white/70 backdrop-blur-sm border border-slate-300 hover:border-blue-400 hover:bg-white rounded-3xl hover:shadow-xl hover:shadow-blue-500/10 transition-all duration-300 no-underline overflow-hidden"
      >
        <div class="absolute -right-6 -top-6 w-32 h-32 bg-blue-100 rounded-full blur-2xl opacity-40 group-hover:opacity-70 transition-opacity duration-500"></div>
        <div class="flex items-center gap-4 mb-4 relative z-10 w-full">
          <div class="w-12 h-12 shrink-0 rounded-2xl bg-blue-100 text-blue-600 flex items-center justify-center shadow-sm border border-blue-200/50">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
          </div>
          <h3 class="text-lg font-bold text-slate-800 group-hover:text-blue-600 transition-colors">Tayo</h3>
        </div>
        <p class="text-sm text-slate-500 leading-relaxed relative z-10 flex-1">Menebak posisi gambar bus Tayo dan teman-teman (Mudah).</p>
        <div class="mt-6 flex items-center text-sm font-bold text-blue-600 opacity-100 md:opacity-0 group-hover:opacity-100 transform translate-x-0 md:-translate-x-2.5 group-hover:translate-x-0 transition-all duration-300 relative z-10">
          Mulai Bermain
          <svg class="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"></path></svg>
        </div>
      </a>

      <!-- Card Buah -->
      <a
        href="/game-edukasi/tebak-4-gambar/buah"
        class="group relative flex flex-col items-start p-6 bg-white/70 backdrop-blur-sm border border-slate-300 hover:border-emerald-400 hover:bg-white rounded-3xl hover:shadow-xl hover:shadow-emerald-500/10 transition-all duration-300 no-underline overflow-hidden"
      >
        <div class="absolute -right-6 -top-6 w-32 h-32 bg-emerald-100 rounded-full blur-2xl opacity-40 group-hover:opacity-70 transition-opacity duration-500"></div>
        <div class="flex items-center gap-4 mb-4 relative z-10 w-full">
          <div class="w-12 h-12 shrink-0 rounded-2xl bg-emerald-100 text-emerald-600 flex items-center justify-center shadow-sm border border-emerald-200/50">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
          </div>
          <h3 class="text-lg font-bold text-slate-800 group-hover:text-emerald-600 transition-colors">Buah</h3>
        </div>
        <p class="text-sm text-slate-500 leading-relaxed relative z-10 flex-1">Menebak letak posisi gambar buah segar (Mudah).</p>
        <div class="mt-6 flex items-center text-sm font-bold text-emerald-600 opacity-100 md:opacity-0 group-hover:opacity-100 transform translate-x-0 md:-translate-x-2.5 group-hover:translate-x-0 transition-all duration-300 relative z-10">
          Mulai Bermain
          <svg class="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"></path></svg>
        </div>
      </a>
    </div>
  </div>
</div>
