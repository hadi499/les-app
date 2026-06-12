<script lang="ts">
  import { goto } from "$app/navigation";
  import { getCurrentUserId, getGameHighScores } from "$lib/stores/users";

  let highScores = $state({ left: 0, right: 0, both: 0, letters: 0, all: 0 });

  function selectMode(mode: "left" | "right" | "both" | "letters" | "all") {
    goto(`/mengetik/game/${mode}`);
  }

  $effect(() => {
    const userId = getCurrentUserId();
    if (userId) {
      highScores = getGameHighScores(userId);
    }
  });
</script>

<svelte:head>
  <title>Game Typing - Ketik 10 Jari</title>
</svelte:head>

<div class="min-h-screen bg-[#0C134F] p-4 sm:p-8 text-white font-sans selection:bg-zinc-800 selection:text-white">
  <div class="max-w-5xl mx-auto relative z-10">
    <header class="text-center mb-10 sm:mb-16">
      <h1 class="text-3xl sm:text-4xl md:text-5xl font-medium mb-4 bg-gradient-to-br from-amber-400 to-amber-500 bg-clip-text text-transparent drop-shadow-md">🎮 Game Typing Seru!</h1>
      <p class="text-lg sm:text-xl text-blue-200 font-light tracking-wide">Pilih mode latihan yang ingin kamu mainkan</p>
    </header>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-6 sm:gap-8 mb-12">
      <!-- Left -->
      <button class="bg-zinc-900/50 backdrop-blur-sm border-2 border-zinc-700 rounded-3xl p-6 sm:p-8 text-center cursor-pointer transition-all hover:-translate-y-2 hover:shadow-xl hover:shadow-blue-500/20 hover:border-blue-500 hover:bg-blue-900/30 group" onclick={() => selectMode("left")}>
        <div class="text-5xl sm:text-6xl mb-4 group-hover:scale-110 transition-transform drop-shadow-sm">👈</div>
        <h2 class="text-xl sm:text-2xl font-medium mb-2 text-white drop-shadow-sm">Home Row Kiri</h2>
        <p class="text-sm sm:text-base text-blue-200 mb-6 font-light">Latih jari tangan kiri</p>
        <div class="font-mono text-xl sm:text-2xl font-bold tracking-widest text-amber-400 mb-4 drop-shadow-sm">A S D F</div>
        <div class="inline-block px-4 py-2 bg-amber-900/30 border border-amber-700/50 text-amber-300 rounded-full text-xs sm:text-sm font-medium">Mudah</div>
        {#if highScores.left > 0}
          <div class="mt-4 text-sm sm:text-lg font-bold text-green-400">🏆 {highScores.left}</div>
        {/if}
      </button>

      <!-- Right -->
      <button class="bg-zinc-900/50 backdrop-blur-sm border-2 border-zinc-700 rounded-3xl p-6 sm:p-8 text-center cursor-pointer transition-all hover:-translate-y-2 hover:shadow-xl hover:shadow-pink-500/20 hover:border-pink-500 hover:bg-pink-900/30 group" onclick={() => selectMode("right")}>
        <div class="text-5xl sm:text-6xl mb-4 group-hover:scale-110 transition-transform drop-shadow-sm">👉</div>
        <h2 class="text-xl sm:text-2xl font-medium mb-2 text-white drop-shadow-sm">Home Row Kanan</h2>
        <p class="text-sm sm:text-base text-blue-200 mb-6 font-light">Latih jari tangan kanan</p>
        <div class="font-mono text-xl sm:text-2xl font-bold tracking-widest text-amber-400 mb-4 drop-shadow-sm">J K L ;</div>
        <div class="inline-block px-4 py-2 bg-amber-900/30 border border-amber-700/50 text-amber-300 rounded-full text-xs sm:text-sm font-medium">Mudah</div>
        {#if highScores.right > 0}
          <div class="mt-4 text-sm sm:text-lg font-bold text-green-400">🏆 {highScores.right}</div>
        {/if}
      </button>

      <!-- Both -->
      <button class="bg-zinc-900/50 backdrop-blur-sm border-2 border-zinc-700 rounded-3xl p-6 sm:p-8 text-center cursor-pointer transition-all hover:-translate-y-2 hover:shadow-xl hover:shadow-green-500/20 hover:border-green-500 hover:bg-green-900/30 group" onclick={() => selectMode("both")}>
        <div class="text-5xl sm:text-6xl mb-4 group-hover:scale-110 transition-transform drop-shadow-sm">🎯</div>
        <h2 class="text-xl sm:text-2xl font-medium mb-2 text-white drop-shadow-sm">Home Row Lengkap</h2>
        <p class="text-sm sm:text-base text-blue-200 mb-6 font-light">Latih kedua tangan sekaligus</p>
        <div class="font-mono text-xl sm:text-2xl font-bold tracking-widest text-amber-400 mb-4 drop-shadow-sm">A S D F J K L ;</div>
        <div class="inline-block px-4 py-2 bg-amber-900/30 border border-amber-700/50 text-amber-300 rounded-full text-xs sm:text-sm font-medium">Sedang</div>
        {#if highScores.both > 0}
          <div class="mt-4 text-sm sm:text-lg font-bold text-green-400">🏆 {highScores.both}</div>
        {/if}
      </button>

      <!-- Letters -->
      <button class="bg-zinc-900/50 backdrop-blur-sm border-2 border-zinc-700 rounded-3xl p-6 sm:p-8 text-center cursor-pointer transition-all hover:-translate-y-2 hover:shadow-xl hover:shadow-amber-500/20 hover:border-amber-500 hover:bg-amber-900/30 group" onclick={() => selectMode("letters")}>
        <div class="text-5xl sm:text-6xl mb-4 group-hover:scale-110 transition-transform drop-shadow-sm">🔤</div>
        <h2 class="text-xl sm:text-2xl font-medium mb-2 text-white drop-shadow-sm">Semua Huruf</h2>
        <p class="text-sm sm:text-base text-blue-200 mb-6 font-light">Latih semua huruf A-Z</p>
        <div class="font-mono text-xl sm:text-2xl font-bold tracking-widest text-amber-400 mb-4 drop-shadow-sm">A-Z</div>
        <div class="inline-block px-4 py-2 bg-amber-900/30 border border-amber-700/50 text-amber-300 rounded-full text-xs sm:text-sm font-medium">Sedang</div>
        {#if highScores.letters > 0}
          <div class="mt-4 text-sm sm:text-lg font-bold text-green-400">🏆 {highScores.letters}</div>
        {/if}
      </button>

      <!-- All -->
      <button class="md:col-span-2 bg-zinc-900/50 backdrop-blur-sm border-2 border-zinc-700 rounded-3xl p-6 sm:p-8 text-center cursor-pointer transition-all hover:-translate-y-2 hover:shadow-xl hover:shadow-purple-500/20 hover:border-purple-500 hover:bg-purple-900/30 group" onclick={() => selectMode("all")}>
        <div class="text-5xl sm:text-6xl mb-4 group-hover:scale-110 transition-transform drop-shadow-sm">🚀</div>
        <h2 class="text-xl sm:text-2xl font-medium mb-2 text-white drop-shadow-sm">Semua Huruf & Karakter</h2>
        <p class="text-sm sm:text-base text-blue-200 mb-6 font-light">Tantang dirimu dengan semua tombol</p>
        <div class="font-mono text-xl sm:text-2xl font-bold tracking-widest text-amber-400 mb-4 drop-shadow-sm">A-Z `-=[]\;',./</div>
        <div class="inline-block px-4 py-2 bg-amber-900/30 border border-amber-700/50 text-amber-300 rounded-full text-xs sm:text-sm font-medium">Sulit</div>
        {#if highScores.all > 0}
          <div class="mt-4 text-sm sm:text-lg font-bold text-green-400">🏆 {highScores.all}</div>
        {/if}
      </button>
    </div>

    <div class="bg-indigo-900/30 backdrop-blur-sm border border-indigo-700/50 rounded-2xl p-6 sm:p-8">
      <h3 class="text-xl font-medium text-amber-400 mb-4 drop-shadow-sm">Cara Bermain:</h3>
      <ul class="list-disc pl-6 space-y-2 text-blue-100 font-light m-0">
        <li>Huruf akan melayang dari bawah ke atas</li>
        <li>Ketik huruf sebelum menyentuh garis atas</li>
        <li>Setiap 100 poin, kecepatan meningkat</li>
        <li>Kamu punya 5 nyawa, jangan sampai habis!</li>
      </ul>
    </div>
  </div>
</div>
