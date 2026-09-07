<script lang="ts">
  import { getProgress } from "$lib/stores/progress";
  import { lessons } from "$lib/data/lessons";
</script>

<svelte:head>
  <title>Ketik 10 Jari - Belajar Mengetik Seru!</title>
  <meta name="description" content="Belajar mengetik 10 jari dengan cara yang menyenangkan di Les Balongarut. Latihan interaktif dari level dasar hingga mahir, cocok untuk semua usia." />
  <link rel="canonical" href="https://lesbalonggarut.my.id/mengetik" />

  <!-- Open Graph -->
  <meta property="og:type" content="website" />
  <meta property="og:url" content="https://lesbalonggarut.my.id/mengetik" />
  <meta property="og:title" content="Ketik 10 Jari - Belajar Mengetik Seru!" />
  <meta property="og:description" content="Latihan mengetik 10 jari dari level dasar hingga mahir. Interaktif dan menyenangkan!" />
  <meta property="og:site_name" content="Les Balongarut" />
  <meta property="og:locale" content="id_ID" />

  <!-- Twitter Card -->
  <meta name="twitter:card" content="summary" />
  <meta name="twitter:title" content="Ketik 10 Jari - Belajar Mengetik Seru!" />
  <meta name="twitter:description" content="Belajar mengetik 10 jari dengan cara yang menyenangkan di Les Balongarut." />
</svelte:head>

<div
  class="min-h-screen bg-slate-50 text-slate-900 font-sans selection:bg-blue-200 selection:text-blue-900"
>
  <main class="max-w-4xl mx-auto p-4 sm:p-8 relative z-10">
    <!-- Mobile Warning -->
    <div class="mb-8 p-4 bg-amber-100 border border-amber-300 rounded-2xl sm:hidden flex items-start gap-3 shadow-sm">
      <div class="text-2xl shrink-0 drop-shadow-sm">⚠️</div>
      <div>
        <h3 class="font-medium text-amber-900 mb-1 drop-shadow-sm">Akses Melalui Desktop</h3>
        <p class="text-sm text-amber-800 m-0 font-light">
          Aplikasi mengetik ini berfungsi optimal jika dibuka melalui Desktop atau Laptop dengan keyboard fisik.
        </p>
      </div>
    </div>

    <!-- Basic Practice Card -->
    <a
      href="/mengetik/practice/basic"
      class="flex flex-col sm:flex-row items-center gap-6 p-6 sm:p-8 bg-emerald-100/50 backdrop-blur-sm border border-emerald-300 hover:border-emerald-500 hover:shadow-xl hover:shadow-emerald-900/10 hover:-translate-y-1 transition-all rounded-3xl mb-8 group no-underline"
    >
      <div class="text-6xl shrink-0 drop-shadow-md">⌨️</div>
      <div class="flex-1 text-center sm:text-left">
        <h2
          class="text-2xl sm:text-3xl font-medium mb-2 text-slate-900 drop-shadow-sm"
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
          class="text-2xl sm:text-3xl font-medium mb-2 text-slate-900 drop-shadow-sm"
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
      <h2 class="text-3xl font-medium text-slate-900 mb-2 drop-shadow-sm">
        📚 Pilih Pelajaran
      </h2>
      <p class="text-slate-600 font-light m-0">
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
            ? 'opacity-60 cursor-not-allowed bg-white/40 border-slate-200'
            : p.completed
              ? 'bg-emerald-100/50 border-emerald-300 hover:border-emerald-500 hover:shadow-lg hover:shadow-green-500/10 hover:-translate-y-1'
              : 'bg-white/60 border-slate-300 hover:border-blue-400 hover:shadow-lg hover:shadow-blue-900/10 hover:-translate-y-1 hover:bg-slate-50'}"
        >
          <div
            class="text-4xl shrink-0 w-14 h-14 flex items-center justify-center bg-white/80 rounded-xl border border-slate-300 drop-shadow-sm"
          >
            {isLocked ? "🔒" : lesson.icon}
          </div>
          <div class="flex-1 min-w-0">
            <h3 class="text-lg font-medium text-slate-900 mb-1 drop-shadow-sm">
              {lesson.title}
            </h3>
            <p
              class="text-sm text-slate-600 mb-3 line-clamp-2 leading-snug font-light"
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
              <p class="text-xs font-medium text-slate-500 m-0">
                Butuh akurasi 80% di pelajaran sebelumnya
              </p>
            {/if}
          </div>
        </a>
      {/each}
    </div>

    <!-- Tips -->
    <div
      class="bg-white/60 backdrop-blur-sm border border-slate-200 rounded-2xl p-6 sm:p-8"
    >
      <h3 class="text-xl font-medium text-slate-900 mb-4 drop-shadow-sm">
        💡 Tips Mengetik 10 Jari
      </h3>
      <ul class="list-disc pl-5 space-y-2 text-slate-600 font-light m-0">
        <li>
          Letakkan jari di <strong class="text-slate-800 font-medium"
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
