<script lang="ts">
  import { exercises } from '$lib/components/LiveEditor/starters';
  
  // Convert object to array for easier iteration
  const exerciseList = Object.values(exercises);
</script>

<svelte:head>
  <title>Pilih Bahasa Pemrograman - Belajar Coding</title>
</svelte:head>

<div class="min-h-screen bg-slate-50 font-sans p-6 sm:p-12">
  <div class="max-w-4xl mx-auto">
    <!-- Header -->
    <header class="mb-12 text-center">
      <a href="/" class="inline-flex items-center gap-2 px-4 py-2 text-sm font-semibold text-slate-600 hover:text-slate-900 bg-white border border-slate-200 hover:border-slate-300 hover:bg-slate-50 hover:shadow-sm rounded-xl transition-all mb-8">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path></svg>
        Kembali ke Beranda
      </a>
      <h1 class="text-4xl sm:text-5xl font-extrabold text-slate-900 tracking-tight mb-4">
        Belajar Coding
      </h1>
      <p class="text-lg text-slate-600 max-w-2xl mx-auto">
        Pilih bahasa pemrograman yang ingin kamu pelajari. Semua modul dapat dijalankan langsung di peramban tanpa perlu instalasi aplikasi!
      </p>
    </header>

    <!-- Cards Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      {#each exerciseList as exercise}
        <a 
          href="/belajar-coding/{exercise.id}" 
          class="group relative flex flex-col bg-white border border-slate-200 rounded-2xl p-6 sm:p-8 hover:border-blue-400 hover:shadow-xl hover:shadow-blue-500/10 transition-all duration-300 no-underline overflow-hidden"
        >
          <!-- Decorative Background element -->
          <div class="absolute -right-12 -top-12 w-40 h-40 bg-gradient-to-br {exercise.type === 'web' ? 'from-blue-50 to-indigo-100' : 'from-yellow-50 to-yellow-100'} rounded-full blur-2xl opacity-50 group-hover:opacity-100 transition-opacity"></div>
          
          <div class="relative z-10 flex items-center gap-4 mb-4">
            <div class="flex items-center justify-center w-12 h-12 rounded-xl {exercise.type === 'web' ? 'bg-blue-100 text-blue-600' : 'bg-yellow-100 text-yellow-600'}">
              {#if exercise.type === 'web'}
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"></path></svg>
              {:else}
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg>
              {/if}
            </div>
            <h2 class="text-xl font-bold text-slate-800 group-hover:text-blue-600 transition-colors">
              {exercise.title}
            </h2>
          </div>
          
          <p class="relative z-10 text-slate-600 text-sm leading-relaxed flex-1">
            {exercise.description}
          </p>
          
          <div class="relative z-10 mt-6 flex items-center text-sm font-semibold text-blue-600 opacity-100 md:opacity-0 group-hover:opacity-100 transform translate-x-0 md:-translate-x-2.5 group-hover:translate-x-0 transition-all duration-300">
            Mulai Belajar
            <svg class="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"></path></svg>
          </div>
        </a>
      {/each}

      <!-- Coding Anak Card -->
      <a 
        href="/belajar-coding/kids/1" 
        class="group relative flex flex-col bg-white border border-slate-200 rounded-2xl p-6 sm:p-8 hover:border-green-400 hover:shadow-xl hover:shadow-green-500/10 transition-all duration-300 no-underline overflow-hidden md:col-span-2 lg:col-span-1"
      >
        <div class="absolute -right-12 -top-12 w-40 h-40 bg-gradient-to-br from-green-50 to-emerald-100 rounded-full blur-2xl opacity-50 group-hover:opacity-100 transition-opacity"></div>
        
        <div class="relative z-10 flex items-center gap-4 mb-4">
          <div class="flex items-center justify-center w-12 h-12 rounded-xl bg-green-100 text-green-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 10l-2 1m0 0l-2-1m2 1v2.5M20 7l-2 1m2-1l-2-1m2 1v2.5M14 4l-2-1-2 1M4 7l2-1M4 7l2 1M4 7v2.5M12 21l-2-1m2 1l2-1m-2 1v-2.5M6 18l-2-1v-2.5M18 18l2-1v-2.5"></path></svg>
          </div>
          <h2 class="text-xl font-bold text-slate-800 group-hover:text-green-600 transition-colors">
            Coding Anak (Blockly)
          </h2>
        </div>
        
        <p class="relative z-10 text-slate-600 text-sm leading-relaxed flex-1">
          Belajar logika pemrograman tanpa harus mengetik! Susun blok-blok bergambar layaknya bermain puzzle untuk menggerakkan karakter. (Cocok untuk usia 5-10 tahun)
        </p>
        
        <div class="relative z-10 mt-6 flex items-center text-sm font-semibold text-green-600 opacity-100 md:opacity-0 group-hover:opacity-100 transform translate-x-0 md:-translate-x-2.5 group-hover:translate-x-0 transition-all duration-300">
          Mulai Bermain
          <svg class="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"></path></svg>
        </div>
      </a>
    </div>
  </div>
</div>
