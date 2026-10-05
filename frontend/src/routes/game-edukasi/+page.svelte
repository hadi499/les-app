<script lang="ts">
</script>

<svelte:head>
  <title>Game Edukasi - Les Balonggarut</title>
</svelte:head>

<div class="min-h-screen bg-slate-50 font-sans selection:bg-blue-200 selection:text-blue-900 pt-20 sm:pt-24 pb-16 px-4 sm:px-6 text-slate-900">
  <div class="max-w-4xl mx-auto">
    <!-- Header with back button -->
    <div class="flex items-center gap-4 mb-8 sm:mb-12">
      <a
        href="/"
        class="p-1.5 sm:p-2 bg-white rounded-full shadow-sm text-slate-500 hover:text-slate-800 hover:shadow transition-all border border-slate-200"
        aria-label="Kembali ke Beranda"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="w-5 h-5 sm:w-6 sm:h-6"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2.5"
          stroke-linecap="round"
          stroke-linejoin="round"><path d="m15 18-6-6 6-6" /></svg>
      </a>
      <div>
        <h1 class="text-2xl sm:text-3xl font-bold text-slate-800">Game Edukasi</h1>
        <p class="text-sm text-slate-500 mt-1">Kumpulan permainan interaktif untuk melatih memori dan logika.</p>
      </div>
    </div>

    <!-- The grid of games -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6 w-full text-left">
      <!-- Card Berhitung -->
      <a
        href="/game-edukasi/berhitung"
        class="group relative flex flex-col items-start p-6 bg-white/70 backdrop-blur-sm border border-slate-300 hover:border-pink-400 hover:bg-white rounded-3xl hover:shadow-xl hover:shadow-pink-500/10 transition-all duration-300 no-underline overflow-hidden"
      >
        <div class="absolute -right-6 -top-6 w-32 h-32 bg-pink-100 rounded-full blur-2xl opacity-40 group-hover:opacity-70 transition-opacity duration-500"></div>
        <div class="flex items-center gap-4 mb-4 relative z-10 w-full">
          <div class="w-12 h-12 shrink-0 rounded-2xl bg-pink-100 text-pink-600 flex items-center justify-center shadow-sm border border-pink-200/50">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z"/></svg>
          </div>
          <h3 class="text-lg font-bold text-slate-800 group-hover:text-pink-600 transition-colors">Latihan Berhitung</h3>
        </div>
        <p class="text-sm text-slate-500 leading-relaxed relative z-10 flex-1">Tingkatkan kecepatan dan ketepatan menghitung dengan modul aritmetika yang terstruktur untuk semua level.</p>
        <div class="mt-6 flex items-center text-sm font-bold text-pink-600 opacity-100 md:opacity-0 group-hover:opacity-100 transform translate-x-0 md:-translate-x-2.5 group-hover:translate-x-0 transition-all duration-300 relative z-10">
          Mulai Berlatih
          <svg class="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"></path></svg>
        </div>
      </a>

      <!-- Card Tebak 4 Gambar -->
      <a
        href="/game-edukasi/tebak-4-gambar"
        class="group relative flex flex-col items-start p-6 bg-white/70 backdrop-blur-sm border border-slate-300 hover:border-yellow-400 hover:bg-white rounded-3xl hover:shadow-xl hover:shadow-yellow-500/10 transition-all duration-300 no-underline overflow-hidden"
      >
        <div class="absolute -right-6 -top-6 w-32 h-32 bg-yellow-100 rounded-full blur-2xl opacity-40 group-hover:opacity-70 transition-opacity duration-500"></div>
        <div class="flex items-center gap-4 mb-4 relative z-10 w-full">
          <div class="w-12 h-12 shrink-0 rounded-2xl bg-yellow-100 text-yellow-600 flex items-center justify-center shadow-sm border border-yellow-200/50">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
          </div>
          <h3 class="text-lg font-bold text-slate-800 group-hover:text-yellow-600 transition-colors">Tebak 4 Gambar</h3>
        </div>
        <p class="text-sm text-slate-500 leading-relaxed relative z-10 flex-1">Latih daya ingat tingkat mudah dengan menebak posisi 4 gambar.</p>
        <div class="mt-6 flex items-center text-sm font-bold text-yellow-600 opacity-100 md:opacity-0 group-hover:opacity-100 transform translate-x-0 md:-translate-x-2.5 group-hover:translate-x-0 transition-all duration-300 relative z-10">
          Lihat Daftar Game
          <svg class="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"></path></svg>
        </div>
      </a>

      <!-- Card Tebak 6 Gambar -->
      <a
        href="/game-edukasi/tebak-6-gambar"
        class="group relative flex flex-col items-start p-6 bg-white/70 backdrop-blur-sm border border-slate-300 hover:border-emerald-400 hover:bg-white rounded-3xl hover:shadow-xl hover:shadow-emerald-500/10 transition-all duration-300 no-underline overflow-hidden"
      >
        <div class="absolute -right-6 -top-6 w-32 h-32 bg-emerald-100 rounded-full blur-2xl opacity-40 group-hover:opacity-70 transition-opacity duration-500"></div>
        <div class="flex items-center gap-4 mb-4 relative z-10 w-full">
          <div class="w-12 h-12 shrink-0 rounded-2xl bg-emerald-100 text-emerald-600 flex items-center justify-center shadow-sm border border-emerald-200/50">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
          </div>
          <h3 class="text-lg font-bold text-slate-800 group-hover:text-emerald-600 transition-colors">Tebak 6 Gambar</h3>
        </div>
        <p class="text-sm text-slate-500 leading-relaxed relative z-10 flex-1">Tantang daya ingat tingkat sedang dengan menebak posisi 6 gambar.</p>
        <div class="mt-6 flex items-center text-sm font-bold text-emerald-600 opacity-100 md:opacity-0 group-hover:opacity-100 transform translate-x-0 md:-translate-x-2.5 group-hover:translate-x-0 transition-all duration-300 relative z-10">
          Lihat Daftar Game
          <svg class="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"></path></svg>
        </div>
      </a>

      <!-- Card Memory Game -->
      <a
        href="/game-edukasi/memory-game"
        class="group relative flex flex-col items-start p-6 bg-white/70 backdrop-blur-sm border border-slate-300 hover:border-purple-400 hover:bg-white rounded-3xl hover:shadow-xl hover:shadow-purple-500/10 transition-all duration-300 no-underline overflow-hidden"
      >
        <div class="absolute -right-6 -top-6 w-32 h-32 bg-purple-100 rounded-full blur-2xl opacity-40 group-hover:opacity-70 transition-opacity duration-500"></div>
        <div class="flex items-center gap-4 mb-4 relative z-10 w-full">
          <div class="w-12 h-12 shrink-0 rounded-2xl bg-purple-100 text-purple-600 flex items-center justify-center shadow-sm border border-purple-200/50">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
          </div>
          <h3 class="text-lg font-bold text-slate-800 group-hover:text-purple-600 transition-colors">Memory Game</h3>
        </div>
        <p class="text-sm text-slate-500 leading-relaxed relative z-10 flex-1">Latih daya ingat dengan mencocokkan pasangan kartu emoji yang tersembunyi.</p>
        <div class="mt-6 flex items-center text-sm font-bold text-purple-600 opacity-100 md:opacity-0 group-hover:opacity-100 transform translate-x-0 md:-translate-x-2.5 group-hover:translate-x-0 transition-all duration-300 relative z-10">
          Mulai Bermain
          <svg class="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"></path></svg>
        </div>
      </a>

      <!-- Card Odd One Out -->
      <a
        href="/game-edukasi/odd-one-out"
        class="group relative flex flex-col items-start p-6 bg-white/70 backdrop-blur-sm border border-slate-300 hover:border-orange-400 hover:bg-white rounded-3xl hover:shadow-xl hover:shadow-orange-500/10 transition-all duration-300 no-underline overflow-hidden"
      >
        <div class="absolute -right-6 -top-6 w-32 h-32 bg-orange-100 rounded-full blur-2xl opacity-40 group-hover:opacity-70 transition-opacity duration-500"></div>
        <div class="flex items-center gap-4 mb-4 relative z-10 w-full">
          <div class="w-12 h-12 shrink-0 rounded-2xl bg-orange-100 text-orange-600 flex items-center justify-center shadow-sm border border-orange-200/50">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
          </div>
          <h3 class="text-lg font-bold text-slate-800 group-hover:text-orange-600 transition-colors">Cari yang Beda</h3>
        </div>
        <p class="text-sm text-slate-500 leading-relaxed relative z-10 flex-1">Cari satu gambar yang berbeda dari sekumpulan gambar yang mirip.</p>
        <div class="mt-6 flex items-center text-sm font-bold text-orange-600 opacity-100 md:opacity-0 group-hover:opacity-100 transform translate-x-0 md:-translate-x-2.5 group-hover:translate-x-0 transition-all duration-300 relative z-10">
          Mulai Bermain
          <svg class="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"></path></svg>
        </div>
      </a>

      <!-- Card Schulte Table -->
      <a
        href="/game-edukasi/sculte-game"
        class="group relative flex flex-col items-start p-6 bg-white/70 backdrop-blur-sm border border-slate-300 hover:border-sky-400 hover:bg-white rounded-3xl hover:shadow-xl hover:shadow-sky-500/10 transition-all duration-300 no-underline overflow-hidden"
      >
        <div class="absolute -right-6 -top-6 w-32 h-32 bg-sky-100 rounded-full blur-2xl opacity-40 group-hover:opacity-70 transition-opacity duration-500"></div>
        <div class="flex items-center gap-4 mb-4 relative z-10 w-full">
          <div class="w-12 h-12 shrink-0 rounded-2xl bg-sky-100 text-sky-600 flex items-center justify-center shadow-sm border border-sky-200/50">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
          </div>
          <h3 class="text-lg font-bold text-slate-800 group-hover:text-sky-600 transition-colors">Tabel Schulte</h3>
        </div>
        <p class="text-sm text-slate-500 leading-relaxed relative z-10 flex-1">Latih kecepatan membaca dan fokus dengan mencari angka berurutan di tabel.</p>
        <div class="mt-6 flex items-center text-sm font-bold text-sky-600 opacity-100 md:opacity-0 group-hover:opacity-100 transform translate-x-0 md:-translate-x-2.5 group-hover:translate-x-0 transition-all duration-300 relative z-10">
          Mulai Bermain
          <svg class="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"></path></svg>
        </div>
      </a>

      <!-- Card Schulte Table Prima -->
      <a
        href="/game-edukasi/sculte-game-prime"
        class="group relative flex flex-col items-start p-6 bg-white/70 backdrop-blur-sm border border-slate-300 hover:border-indigo-400 hover:bg-white rounded-3xl hover:shadow-xl hover:shadow-indigo-500/10 transition-all duration-300 no-underline overflow-hidden"
      >
        <div class="absolute -right-6 -top-6 w-32 h-32 bg-indigo-100 rounded-full blur-2xl opacity-40 group-hover:opacity-70 transition-opacity duration-500"></div>
        <div class="flex items-center gap-4 mb-4 relative z-10 w-full">
          <div class="w-12 h-12 shrink-0 rounded-2xl bg-indigo-100 text-indigo-600 flex items-center justify-center shadow-sm border border-indigo-200/50">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
          </div>
          <h3 class="text-lg font-bold text-slate-800 group-hover:text-indigo-600 transition-colors">Tabel Schulte Prima</h3>
        </div>
        <p class="text-sm text-slate-500 leading-relaxed relative z-10 flex-1">Latih kecepatan membaca dan fokus dengan mencari 25 bilangan prima pertama.</p>
        <div class="mt-6 flex items-center text-sm font-bold text-indigo-600 opacity-100 md:opacity-0 group-hover:opacity-100 transform translate-x-0 md:-translate-x-2.5 group-hover:translate-x-0 transition-all duration-300 relative z-10">
          Mulai Bermain
          <svg class="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"></path></svg>
        </div>
      </a>

      <!-- Card English Quiz -->
      <a
        href="/game-edukasi/english-quiz"
        class="group relative flex flex-col items-start p-6 bg-white/70 backdrop-blur-sm border border-slate-300 hover:border-blue-400 hover:bg-white rounded-3xl hover:shadow-xl hover:shadow-blue-500/10 transition-all duration-300 no-underline overflow-hidden"
      >
        <div class="absolute -right-6 -top-6 w-32 h-32 bg-blue-100 rounded-full blur-2xl opacity-40 group-hover:opacity-70 transition-opacity duration-500"></div>
        <div class="flex items-center gap-4 mb-4 relative z-10 w-full">
          <div class="w-12 h-12 shrink-0 rounded-2xl bg-blue-100 text-blue-600 flex items-center justify-center shadow-sm border border-blue-200/50">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
          </div>
          <h3 class="text-lg font-bold text-slate-800 group-hover:text-blue-600 transition-colors">Kuis Bahasa Inggris</h3>
        </div>
        <p class="text-sm text-slate-500 leading-relaxed relative z-10 flex-1">Bermain tebak kata bahasa Inggris dengan kartu flash bergambar.</p>
        <div class="mt-6 flex items-center text-sm font-bold text-blue-600 opacity-100 md:opacity-0 group-hover:opacity-100 transform translate-x-0 md:-translate-x-2.5 group-hover:translate-x-0 transition-all duration-300 relative z-10">
          Mulai Bermain
          <svg class="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"></path></svg>
        </div>
      </a>

      <!-- Card English Flashcard (SD & SMP) -->
      <a
        href="/game-edukasi/english-flashcard"
        class="group relative flex flex-col items-start p-6 bg-white/70 backdrop-blur-sm border border-slate-300 hover:border-teal-400 hover:bg-white rounded-3xl hover:shadow-xl hover:shadow-teal-500/10 transition-all duration-300 no-underline overflow-hidden"
      >
        <div class="absolute -right-6 -top-6 w-32 h-32 bg-teal-100 rounded-full blur-2xl opacity-40 group-hover:opacity-70 transition-opacity duration-500"></div>
        <div class="flex items-center gap-4 mb-4 relative z-10 w-full">
          <div class="w-12 h-12 shrink-0 rounded-2xl bg-teal-100 text-teal-600 flex items-center justify-center shadow-sm border border-teal-200/50">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/></svg>
          </div>
          <h3 class="text-lg font-bold text-slate-800 group-hover:text-teal-600 transition-colors">Flashcard Bahasa Inggris</h3>
        </div>
        <p class="text-sm text-slate-500 leading-relaxed relative z-10 flex-1">Latihan grammar & vocab SD - SMP. Tebak jawaban dengan 2x kesempatan.</p>
        <div class="mt-6 flex items-center text-sm font-bold text-teal-600 opacity-100 md:opacity-0 group-hover:opacity-100 transform translate-x-0 md:-translate-x-2.5 group-hover:translate-x-0 transition-all duration-300 relative z-10">
          Mulai Bermain
          <svg class="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"></path></svg>
        </div>
      </a>
    </div>
  </div>
</div>
