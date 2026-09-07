<svelte:head>
  <title>Game Berhitung Seru - Les Balongarut</title>
  <meta name="description" content="Latih kemampuan berhitung dengan game matematika seru dari Les Balongarut. Tersedia level mudah, sedang, dan sulit untuk semua usia." />
  <link rel="canonical" href="https://lesbalonggarut.my.id/berhitung" />

  <!-- Open Graph -->
  <meta property="og:type" content="website" />
  <meta property="og:url" content="https://lesbalonggarut.my.id/berhitung" />
  <meta property="og:title" content="Game Berhitung Seru - Les Balongarut" />
  <meta property="og:description" content="Latih kemampuan berhitung dengan game matematika seru. Level mudah, sedang, dan sulit." />
  <meta property="og:site_name" content="Les Balongarut" />
  <meta property="og:locale" content="id_ID" />

  <!-- Twitter Card -->
  <meta name="twitter:card" content="summary" />
  <meta name="twitter:title" content="Game Berhitung Seru - Les Balongarut" />
  <meta name="twitter:description" content="Latih kemampuan berhitung dengan game matematika seru dari Les Balongarut." />
</svelte:head>

<div
  class="min-h-screen bg-slate-50 flex flex-col items-center pt-16 pb-10 px-4 font-sans relative overflow-hidden"
>
  <!-- Background Decorations -->
  <div
    class="absolute top-[-10%] left-[-10%] w-96 h-96 bg-yellow-200 rounded-full mix-blend-multiply filter blur-3xl opacity-60 animate-blob"
  ></div>
  <div
    class="absolute bottom-[-10%] right-[-10%] w-96 h-96 bg-pink-300 rounded-full mix-blend-multiply filter blur-3xl opacity-60 animate-blob animation-delay-2000"
  ></div>

  <div class="w-full max-w-4xl z-10 relative mt-10">
    <div class="text-center mb-12">
      <h1
        class="text-5xl md:text-6xl font-black text-pink-600 mb-4 drop-shadow-sm tracking-wide"
      >
        Game Berhitung
      </h1>
      <p class="text-xl text-pink-800 font-medium max-w-xl mx-auto">
        Ayo bermain sambil belajar menghitung! Pilih level kesulitanmu di bawah
        ini:
      </p>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-4xl mx-auto">
      <!-- Mudah -->
      <a
        href="/berhitung/mudah"
        class="group flex flex-col items-center bg-white rounded-[2.5rem] p-8 shadow-[0_12px_0_0_rgba(251,207,232,1)] hover:shadow-[0_4px_0_0_rgba(251,207,232,1)] hover:translate-y-2 transition-all border-4 border-pink-200 no-underline cursor-pointer"
      >
        <div
          class="text-6xl mb-6 group-hover:scale-110 transition-transform animate-bounce"
          style="animation-duration: 2s;"
        >
          🍎
        </div>
        <h2 class="text-3xl font-black text-pink-500 mb-2">Level 1</h2>
        <h3 class="text-xl font-bold text-slate-600 mb-4">Mudah (1 - 10)</h3>
        <p class="text-center text-slate-500 font-medium">
          Belajar menghitung benda dari satu sampai sepuluh.
        </p>
        <div
          class="mt-6 w-full py-3 bg-pink-100 text-pink-600 rounded-2xl font-bold text-center group-hover:bg-pink-500 group-hover:text-white transition-colors"
        >
          Mulai Bermain
        </div>
      </a>

      <!-- Sedang -->
      <a
        href="/berhitung/sedang"
        class="group flex flex-col items-center bg-white rounded-[2.5rem] p-8 shadow-[0_12px_0_0_rgba(167,243,208,1)] hover:shadow-[0_4px_0_0_rgba(167,243,208,1)] hover:translate-y-2 transition-all border-4 border-emerald-200 no-underline cursor-pointer"
      >
        <div
          class="text-6xl mb-6 group-hover:scale-110 transition-transform animate-bounce"
          style="animation-duration: 2.5s;"
        >
          ⭐
        </div>
        <h2 class="text-3xl font-black text-emerald-500 mb-2">Level 2</h2>
        <h3 class="text-xl font-bold text-slate-600 mb-4">Sedang (1 - 10)</h3>
        <p class="text-center text-slate-500 font-medium">
          Latihan berhitung dengan penilaian bintang jika benar semua!
        </p>
        <div
          class="mt-6 w-full py-3 bg-emerald-100 text-emerald-600 rounded-2xl font-bold text-center group-hover:bg-emerald-400 group-hover:text-emerald-900 transition-colors"
        >
          Mulai Bermain
        </div>
      </a>

      <!-- Sulit -->
      <a
        href="/berhitung/sulit"
        class="group flex flex-col items-center bg-white rounded-[2.5rem] p-8 shadow-[0_12px_0_0_rgba(199,210,254,1)] hover:shadow-[0_4px_0_0_rgba(199,210,254,1)] hover:translate-y-2 transition-all border-4 border-indigo-200 no-underline cursor-pointer"
      >
        <div
          class="text-6xl mb-6 group-hover:scale-110 transition-transform animate-bounce"
          style="animation-duration: 3s;"
        >
          🚀
        </div>
        <h2 class="text-3xl font-black text-indigo-500 mb-2">Level 3</h2>
        <h3 class="text-xl font-bold text-slate-600 mb-4">Sulit (11 - 20)</h3>
        <p class="text-center text-slate-500 font-medium">
          Tantangan untuk yang sudah jago! Hitung angka belasan.
        </p>
        <div
          class="mt-6 w-full py-3 bg-indigo-100 text-indigo-600 rounded-2xl font-bold text-center group-hover:bg-indigo-500 group-hover:text-white transition-colors"
        >
          Mulai Bermain
        </div>
      </a>
    </div>

    <div class="mt-16 text-center">
      <a
        href="/"
        class="inline-flex items-center gap-2 px-6 py-3 bg-white text-pink-600 font-bold rounded-2xl shadow-sm border border-pink-200 hover:bg-pink-50 hover:shadow-md transition-all no-underline"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-5 w-5"
          viewBox="0 0 20 20"
          fill="currentColor"
        >
          <path
            fill-rule="evenodd"
            d="M9.707 16.707a1 1 0 01-1.414 0l-6-6a1 1 0 010-1.414l6-6a1 1 0 011.414 1.414L5.414 9H17a1 1 0 110 2H5.414l4.293 4.293a1 1 0 010 1.414z"
            clip-rule="evenodd"
          />
        </svg>
        Kembali ke Beranda
      </a>
    </div>
  </div>
</div>

<style>
  @keyframes blob {
    0% {
      transform: translate(0px, 0px) scale(1);
    }
    33% {
      transform: translate(30px, -50px) scale(1.1);
    }
    66% {
      transform: translate(-20px, 20px) scale(0.9);
    }
    100% {
      transform: translate(0px, 0px) scale(1);
    }
  }
  .animate-blob {
    animation: blob 7s infinite;
  }
  .animation-delay-2000 {
    animation-delay: 2s;
  }
</style>
