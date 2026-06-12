<script lang="ts">
  import { getProgress } from "$lib/stores/progress";
  import { lessons } from "$lib/data/lessons";
</script>

<svelte:head>
  <title>Ketik 10 Jari - Belajar Mengetik Seru!</title>
</svelte:head>

<div class="min-h-screen bg-[#0C134F] text-white font-sans selection:bg-zinc-800 selection:text-white">
  <main class="max-w-4xl mx-auto p-4 sm:p-8 relative z-10">
    <!-- Basic Practice Card -->
    <a href="/mengetik/practice/basic" class="flex flex-col sm:flex-row items-center gap-6 p-6 sm:p-8 bg-green-900/20 backdrop-blur-sm border border-green-700/50 hover:border-green-500 hover:shadow-xl hover:shadow-green-500/20 hover:-translate-y-1 transition-all rounded-3xl mb-8 group no-underline">
      <div class="text-6xl shrink-0 drop-shadow-md">⌨️</div>
      <div class="flex-1 text-center sm:text-left">
        <h2 class="text-2xl sm:text-3xl font-medium mb-2 text-white drop-shadow-sm">Latihan Dasar</h2>
        <p class="text-green-100 font-light m-0">Latihan posisi jari home row: A S D F J K L ;</p>
      </div>
      <div class="text-3xl text-green-400 transition-transform group-hover:translate-x-2 hidden sm:block">→</div>
    </a>

    <!-- Game Card -->
    <a href="/mengetik/game" class="flex flex-col sm:flex-row items-center gap-6 p-6 sm:p-8 bg-amber-900/20 backdrop-blur-sm border border-amber-700/50 hover:border-amber-500 hover:shadow-xl hover:shadow-amber-500/20 hover:-translate-y-1 transition-all rounded-3xl mb-8 group no-underline">
      <div class="text-6xl shrink-0 drop-shadow-md">🎮</div>
      <div class="flex-1 text-center sm:text-left">
        <h2 class="text-2xl sm:text-3xl font-medium mb-2 text-white drop-shadow-sm">Game Typing Seru!</h2>
        <p class="text-amber-100 font-light m-0">Latih kecepatan mengetikmu dengan game yang menyenangkan</p>
      </div>
      <div class="text-3xl text-amber-400 transition-transform group-hover:translate-x-2 hidden sm:block">→</div>
    </a>

    <!-- Lessons Intro -->
    <div class="mb-8 text-center sm:text-left">
      <h2 class="text-3xl font-medium text-white mb-2 drop-shadow-sm">📚 Pilih Pelajaran</h2>
      <p class="text-blue-200 font-light m-0">Mulai dari home row, lalu naik level satu per satu!</p>
    </div>

    <!-- Lessons Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-8">
      {#each lessons as lesson, i}
        {@const p = getProgress(lesson.id)}
        {@const prevProgress = i > 0 ? getProgress(lessons[i - 1].id) : null}
        {@const isLocked = lesson.unlockScore > 0 && (!prevProgress || prevProgress.bestAccuracy < 80)}
        <a
          href={isLocked ? undefined : `/mengetik/practice/${lesson.id}`}
          class="flex gap-4 p-5 rounded-2xl border backdrop-blur-sm transition-all duration-300 no-underline {isLocked ? 'opacity-60 cursor-not-allowed bg-zinc-900/30 border-zinc-800' : p.completed ? 'bg-green-900/20 border-green-700/50 hover:border-green-500 hover:shadow-lg hover:shadow-green-500/10 hover:-translate-y-1' : 'bg-zinc-900/50 border-zinc-700 hover:border-blue-400 hover:shadow-lg hover:shadow-blue-500/10 hover:-translate-y-1 hover:bg-blue-900/20'}"
        >
          <div class="text-4xl shrink-0 w-14 h-14 flex items-center justify-center bg-zinc-800/80 rounded-xl border border-zinc-700/50 drop-shadow-sm">
            {isLocked ? "🔒" : lesson.icon}
          </div>
          <div class="flex-1 min-w-0">
            <h3 class="text-lg font-medium text-white mb-1 drop-shadow-sm">{lesson.title}</h3>
            <p class="text-sm text-blue-200 mb-3 line-clamp-2 leading-snug font-light">{lesson.description}</p>
            {#if !isLocked}
              <div class="flex items-center gap-4">
                {#if p.bestWPM > 0}
                  <span class="text-sm font-bold text-amber-400">{p.bestWPM} WPM</span>
                {/if}
                {#if p.bestAccuracy > 0}
                  <span class="text-sm font-bold text-green-400">{p.bestAccuracy}%</span>
                {/if}
              </div>
            {:else}
              <p class="text-xs font-medium text-zinc-500 m-0">Butuh akurasi 80% di pelajaran sebelumnya</p>
            {/if}
          </div>
        </a>
      {/each}
    </div>

    <!-- Tips -->
    <div class="bg-indigo-900/30 backdrop-blur-sm border border-indigo-700/50 rounded-2xl p-6 sm:p-8">
      <h3 class="text-xl font-medium text-white mb-4 drop-shadow-sm">💡 Tips Mengetik 10 Jari</h3>
      <ul class="list-disc pl-5 space-y-2 text-indigo-100 font-light m-0">
        <li>Letakkan jari di <strong class="text-indigo-300 font-medium">home row</strong> (ASDF - JKL;)</li>
        <li>Jangan lihat keyboard, lihat layar!</li>
        <li>Mulai pelan, kecepatan datang dengan latihan</li>
        <li>Gunakan jari yang tepat untuk setiap tombol</li>
      </ul>
    </div>
  </main>
</div>
