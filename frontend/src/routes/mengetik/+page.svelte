<script lang="ts">
  import { getProgress } from "$lib/stores/progress";
  import { lessons } from "$lib/data/lessons";
</script>

<svelte:head>
  <title>Ketik 10 Jari - Belajar Mengetik Seru!</title>
</svelte:head>

<div
  class="min-h-screen bg-orange-50 text-orange-950 font-sans selection:bg-white selection:text-orange-900"
>
  <main class="max-w-4xl mx-auto p-4 sm:p-8 relative z-10">
    <!-- Basic Practice Card -->
    <a
      href="/mengetik/practice/basic"
      class="flex flex-col sm:flex-row items-center gap-6 p-6 sm:p-8 bg-emerald-100/50 backdrop-blur-sm border border-emerald-300 hover:border-emerald-500 hover:shadow-xl hover:shadow-emerald-900/10 hover:-translate-y-1 transition-all rounded-3xl mb-8 group no-underline"
    >
      <div class="text-6xl shrink-0 drop-shadow-md">⌨️</div>
      <div class="flex-1 text-center sm:text-left">
        <h2
          class="text-2xl sm:text-3xl font-medium mb-2 text-orange-950 drop-shadow-sm"
        >
          Latihan Dasar
        </h2>
        <p class="text-emerald-800 font-light m-0">
          Latihan posisi jari home row: A S D F J K L ;
        </p>
      </div>
      <div
        class="text-3xl text-emerald-700 transition-transform group-hover:translate-x-2 hidden sm:block"
      >
        →
      </div>
    </a>

    <!-- Game Card -->
    <a
      href="/mengetik/game"
      class="flex flex-col sm:flex-row items-center gap-6 p-6 sm:p-8 bg-amber-100/50 backdrop-blur-sm border border-amber-300 hover:border-amber-500 hover:shadow-xl hover:shadow-amber-900/10 hover:-translate-y-1 transition-all rounded-3xl mb-8 group no-underline"
    >
      <div class="text-6xl shrink-0 drop-shadow-md">🎮</div>
      <div class="flex-1 text-center sm:text-left">
        <h2
          class="text-2xl sm:text-3xl font-medium mb-2 text-orange-950 drop-shadow-sm"
        >
          Game Typing Seru!
        </h2>
        <p class="text-amber-800 font-light m-0">
          Latih kecepatan mengetikmu dengan game yang menyenangkan
        </p>
      </div>
      <div
        class="text-3xl text-amber-600 transition-transform group-hover:translate-x-2 hidden sm:block"
      >
        →
      </div>
    </a>

    <!-- Lessons Intro -->
    <div class="mb-8 text-center sm:text-left">
      <h2 class="text-3xl font-medium text-orange-950 mb-2 drop-shadow-sm">
        📚 Pilih Pelajaran
      </h2>
      <p class="text-orange-800 font-light m-0">
        Mulai dari home row, lalu naik level satu per satu!
      </p>
    </div>

    <!-- Lessons Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-8">
      {#each lessons as lesson, i}
        {@const p = getProgress(lesson.id)}
        {@const prevProgress = i > 0 ? getProgress(lessons[i - 1].id) : null}
        {@const isLocked =
          lesson.unlockScore > 0 &&
          (!prevProgress || prevProgress.bestAccuracy < 80)}
        <a
          href={isLocked ? undefined : `/mengetik/practice/${lesson.id}`}
          class="flex gap-4 p-5 rounded-2xl border backdrop-blur-sm transition-all duration-300 no-underline {isLocked
            ? 'opacity-60 cursor-not-allowed bg-white/40 border-orange-200'
            : p.completed
              ? 'bg-emerald-100/50 border-emerald-300 hover:border-emerald-500 hover:shadow-lg hover:shadow-green-500/10 hover:-translate-y-1'
              : 'bg-white/60 border-orange-300 hover:border-orange-400 hover:shadow-lg hover:shadow-orange-900/10 hover:-translate-y-1 hover:bg-orange-50'}"
        >
          <div
            class="text-4xl shrink-0 w-14 h-14 flex items-center justify-center bg-white/80 rounded-xl border border-orange-300 drop-shadow-sm"
          >
            {isLocked ? "🔒" : lesson.icon}
          </div>
          <div class="flex-1 min-w-0">
            <h3 class="text-lg font-medium text-orange-950 mb-1 drop-shadow-sm">
              {lesson.title}
            </h3>
            <p
              class="text-sm text-orange-800 mb-3 line-clamp-2 leading-snug font-light"
            >
              {lesson.description}
            </p>
            {#if !isLocked}
              <div class="flex items-center gap-4">
                {#if p.bestWPM > 0}
                  <span class="text-sm font-bold text-amber-600"
                    >{p.bestWPM} WPM</span
                  >
                {/if}
                {#if p.bestAccuracy > 0}
                  <span class="text-sm font-bold text-emerald-700"
                    >{p.bestAccuracy}%</span
                  >
                {/if}
              </div>
            {:else}
              <p class="text-xs font-medium text-orange-700 m-0">
                Butuh akurasi 80% di pelajaran sebelumnya
              </p>
            {/if}
          </div>
        </a>
      {/each}
    </div>

    <!-- Tips -->
    <div
      class="bg-white/60 backdrop-blur-sm border border-orange-200 rounded-2xl p-6 sm:p-8"
    >
      <h3 class="text-xl font-medium text-orange-950 mb-4 drop-shadow-sm">
        💡 Tips Mengetik 10 Jari
      </h3>
      <ul class="list-disc pl-5 space-y-2 text-orange-800 font-light m-0">
        <li>
          Letakkan jari di <strong class="text-orange-900 font-medium"
            >home row</strong
          > (ASDF - JKL;)
        </li>
        <li>Jangan lihat keyboard, lihat layar!</li>
        <li>Mulai pelan, kecepatan datang dengan latihan</li>
        <li>Gunakan jari yang tepat untuk setiap tombol</li>
      </ul>
    </div>
  </main>
</div>
