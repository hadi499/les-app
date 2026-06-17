<script>
  import { onMount } from "svelte";
  import { Link } from "svelte-routing";

  const EMOJIS = ["🍎", "🐶", "🎈", "🚗", "🧸", "🐱", "🍓", "🦋", "⭐", "⚽"];

  // 1. TAMBAHKAN BINDING AUDIO
  let audioCorrect = $state();
  let audioWrong = $state();
  let audioApplause = $state();

  let currentCount = $state(0);
  let currentEmoji = $state("");
  let options = $state([]);
  let score = $state(0);
  let questionNumber = $state(1);
  let isGameOver = $state(false);
  let showSuccess = $state(false);
  let showError = $state(false);
  let errorShake = $state(false);

  // 2. FUNGSI PEMUTAR SUARA
  function playSound(audioElement) {
    if (audioElement) {
      audioElement.pause();
      audioElement.currentTime = 0;
      let playPromise = audioElement.play();
      if (playPromise !== undefined) {
        playPromise.catch((error) => {
          console.log("Audio diblokir browser:", error);
        });
      }
    }
  }

  function startGame() {
    // 3. MATIKAN SUARA TEPUK TANGAN SAAT RESTART
    if (audioApplause) {
      audioApplause.pause();
      audioApplause.currentTime = 0;
    }

    score = 0;
    questionNumber = 1;
    isGameOver = false;
    nextQuestion();
  }

  function nextQuestion() {
    if (questionNumber > 10) {
      isGameOver = true;
      return;
    }

    currentCount = Math.floor(Math.random() * 10) + 11;
    currentEmoji = EMOJIS[Math.floor(Math.random() * EMOJIS.length)];

    let ops = new Set([currentCount]);
    while (ops.size < 3) {
      let wrongAns = Math.floor(Math.random() * 10) + 11;
      ops.add(wrongAns);
    }

    options = Array.from(ops).sort(() => Math.random() - 0.5);
  }

  function handleAnswer(selected) {
    if (showSuccess || showError) return;

    // 4. DETEKSI SOAL TERAKHIR
    const isLastQuestion = questionNumber >= 10;

    if (selected === currentCount) {
      score++;
      showSuccess = true;

      // 5. MAINKAN SUARA BENAR / TEPUK TANGAN
      if (isLastQuestion) {
        playSound(audioApplause);
      } else {
        playSound(audioCorrect);
      }
    } else {
      showError = true;
      errorShake = true;

      // 6. MAINKAN SUARA SALAH / TEPUK TANGAN
      if (isLastQuestion) {
        playSound(audioApplause);
      } else {
        playSound(audioWrong);
      }
    }

    setTimeout(() => {
      showSuccess = false;
      showError = false;
      errorShake = false;

      if (isLastQuestion) {
        isGameOver = true;
      } else {
        questionNumber++;
        nextQuestion();
      }
    }, 1000);
  }

  onMount(() => {
    startGame();
  });
</script>

<div
  class="min-h-screen bg-indigo-50 flex flex-col items-center py-10 px-4 font-sans relative overflow-hidden"
>
  <div
    class="absolute top-[-10%] left-[-10%] w-96 h-96 bg-blue-200 rounded-full mix-blend-multiply filter blur-3xl opacity-60 animate-blob"
  ></div>
  <div
    class="absolute bottom-[-10%] right-[-10%] w-96 h-96 bg-indigo-300 rounded-full mix-blend-multiply filter blur-3xl opacity-60 animate-blob animation-delay-2000"
  ></div>

  <div class="w-full max-w-md z-10">
    <!-- Header -->
    <div class="flex items-center justify-between mb-8">
      <Link
        to="/"
        class="p-2 bg-white rounded-full shadow-sm text-indigo-500 hover:text-indigo-600 hover:shadow transition-all border border-indigo-100"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="28"
          height="28"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="3"
          stroke-linecap="round"
          stroke-linejoin="round"><path d="m15 18-6-6 6-6" /></svg
        >
      </Link>
      <div class="flex gap-1">
        {#each Array(10) as _, i}
          <div
            class="w-3 h-3 rounded-full {i < questionNumber - 1
              ? 'bg-indigo-500'
              : 'bg-indigo-200'}"
          ></div>
        {/each}
      </div>
    </div>

    <!-- Instruksi Visual -->
    <div
      class="bg-white rounded-3xl p-6 shadow-xl border-4 border-indigo-100 text-center mb-8 relative"
    >
      <h2
        class="text-2xl font-black text-indigo-500 mb-6 tracking-wide select-none"
      >
        Berapa jumlahnya?
      </h2>

      <!-- Area Gambar -->
      <div
        class="flex flex-wrap justify-center gap-1.5 mb-6 min-h-40 items-center"
      >
        {#if showSuccess}
          <div class="text-7xl animate-bounce select-none">🌟</div>
        {:else if showError}
          <div class="text-7xl animate-shake select-none">😢</div>
        {:else}
          {#each Array(currentCount) as _}
            <div
              class="text-4xl animate-pop-in cursor-pointer hover:scale-110 transition-transform select-none"
            >
              {currentEmoji}
            </div>
          {/each}
        {/if}
      </div>
    </div>

    <!-- Pilihan Jawaban -->
    <div class="grid grid-cols-3 gap-4">
      {#each options as opt}
        <button
          onclick={() => handleAnswer(opt)}
          class="select-none [-webkit-tap-highlight-color:transparent] [-webkit-touch-callout:none] aspect-square bg-white border-4 border-indigo-200 rounded-3xl flex items-center justify-center text-5xl font-black text-indigo-500 transition-all {showSuccess ||
          showError
            ? 'cursor-not-allowed opacity-80'
            : 'hover:bg-indigo-50 shadow-[0_8px_0_0_rgba(199,210,254,1)] active:shadow-[0_0px_0_0_rgba(199,210,254,1)] active:translate-y-2'} {errorShake
            ? 'animate-shake'
            : ''}"
        >
          {opt}
        </button>
      {/each}
    </div>

    <!-- Layar Hasil Akhir -->
    {#if isGameOver}
      <div
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-indigo-900/40 backdrop-blur-md animate-fade-in"
      >
        <div
          class="bg-white p-8 rounded-4xl shadow-2xl max-w-sm w-full text-center border-4 border-indigo-200 transform transition-all scale-100 animate-pop"
        >
          {#if score === 10}
            <div class="text-7xl mb-2 animate-bounce">🏆</div>
            <h2 class="text-4xl font-black text-indigo-500 mb-4 tracking-wider">
              HOREEE!
            </h2>
            <p class="text-slate-600 font-bold mb-6">
              Level 2 Lulus! <br /> Kamu super teliti!
            </p>
          {:else}
            <div class="text-7xl mb-2 animate-bounce">✨</div>
            <h2 class="text-3xl font-black text-yellow-500 mb-4">Bagus!</h2>
            <p class="text-slate-600 font-bold mb-6">
              Ayo coba lagi, <br /> sedikit lagi dapat nilai sempurna!
            </p>
          {/if}

          <div
            class="bg-indigo-50 rounded-2xl py-6 mb-8 border-2 border-indigo-100"
          >
            <p class="text-slate-500 font-bold text-lg mb-1">
              Skor Akhir Kamu:
            </p>
            <div
              class="text-indigo-600 text-[100px] font-black drop-shadow-md flex items-baseline justify-center gap-1 leading-none py-2"
            >
              {score}<span class="text-4xl text-indigo-300 font-bold">/10</span>
            </div>
          </div>

          <div class="flex flex-col gap-4">
            <button
              onclick={startGame}
              class="w-full py-4 bg-indigo-500 hover:bg-indigo-600 text-white text-xl font-black rounded-2xl shadow-[0_6px_0_0_rgba(99,102,241,1)] active:shadow-none active:translate-y-1.5 transition-all"
            >
              Main Lagi
            </button>
            <Link
              to="/"
              class="w-full py-4 bg-yellow-400 hover:bg-yellow-500 text-yellow-900 text-xl font-black rounded-2xl shadow-[0_6px_0_0_rgba(202,138,4,1)] active:shadow-none active:translate-y-1.5 transition-all"
            >
              Pulang
            </Link>
          </div>
        </div>
      </div>
    {/if}
  </div>
</div>

<!-- 7. ELEMEN AUDIO DARI FOLDER LOKAL -->
<audio bind:this={audioCorrect} src="/sounds/benar.mp3" preload="auto"></audio>
<audio bind:this={audioWrong} src="/sounds/wrong.mp3" preload="auto"></audio>
<audio bind:this={audioApplause} src="/sounds/tepuk-tangan.mp3" preload="auto"
></audio>

<style>
  @keyframes popIn {
    0% {
      transform: scale(0);
      opacity: 0;
    }
    80% {
      transform: scale(1.2);
      opacity: 1;
    }
    100% {
      transform: scale(1);
      opacity: 1;
    }
  }
  .animate-pop-in {
    animation: popIn 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275) backwards;
  }

  @keyframes shake {
    0%,
    100% {
      transform: translateX(0);
    }
    25% {
      transform: translateX(-8px);
    }
    50% {
      transform: translateX(8px);
    }
    75% {
      transform: translateX(-8px);
    }
  }
  .animate-shake {
    animation: shake 0.3s cubic-bezier(0.36, 0.07, 0.19, 0.97) both;
  }

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

  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }
  @keyframes pop {
    0% {
      transform: scale(0.8);
      opacity: 0;
    }
    100% {
      transform: scale(1);
      opacity: 1;
    }
  }
  .animate-fade-in {
    animation: fadeIn 0.3s ease-out forwards;
  }
  .animate-pop {
    animation: pop 0.5s cubic-bezier(0.175, 0.885, 0.32, 1.275) forwards;
  }
</style>
