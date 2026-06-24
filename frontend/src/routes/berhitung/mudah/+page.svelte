<script lang="ts">
  import { onMount } from "svelte";

  const EMOJIS: string[] = [
    "🍎",
    "🐶",
    "🍩",
    "🚗",
    "🧸",
    "🐱",
    "🍓",
    "🦋",
    "🐮",
    "⚽",
  ];

  let audioSuccess: HTMLAudioElement | undefined = $state();
  let audioError: HTMLAudioElement | undefined = $state();

  let currentCount: number = $state(0);
  let currentEmoji: string = $state("");
  let options: number[] = $state([]);
  let score: number = $state(0);
  let questionNumber: number = $state(1);
  let isGameOver: boolean = $state(false);
  let availableCounts: number[] = [];

  function shuffleArray(array: number[]) {
    for (let i = array.length - 1; i > 0; i--) {
      const j = Math.floor(Math.random() * (i + 1));
      [array[i], array[j]] = [array[j], array[i]];
    }
    return array;
  }
  let showSuccess: boolean = $state(false);
  let errorShake: boolean = $state(false);

  function playSound(audioElement: HTMLAudioElement | undefined) {
    if (audioElement) {
      audioElement.pause();
      audioElement.currentTime = 0;
      let playPromise = audioElement.play();
      if (playPromise !== undefined) {
        playPromise.catch((error) =>
          console.log("Audio diblokir browser:", error),
        );
      }
    }
  }

  function startGame() {
    score = 0;
    questionNumber = 1;
    isGameOver = false;
    availableCounts = shuffleArray(Array.from({length: 10}, (_, i) => i + 1));
    nextQuestion();
  }

  function nextQuestion() {
    if (questionNumber > 10) {
      isGameOver = true;
      return;
    }
    currentCount = availableCounts[questionNumber - 1];
    currentEmoji = EMOJIS[Math.floor(Math.random() * EMOJIS.length)];

    let ops = new Set<number>([currentCount]);
    while (ops.size < 3) {
      let wrongAns = Math.floor(Math.random() * 10) + 1;
      ops.add(wrongAns);
    }
    options = Array.from(ops).sort(() => Math.random() - 0.5);
  }

  function handleAnswer(selected: number) {
    if (showSuccess) return;

    if (selected === currentCount) {
      playSound(audioSuccess);
      score++;
      showSuccess = true;

      setTimeout(() => {
        showSuccess = false;
        questionNumber++;
        nextQuestion();
      }, 1000);
    } else {
      playSound(audioError);
      errorShake = true;
      setTimeout(() => {
        errorShake = false;
      }, 400);
    }
  }

  onMount(() => {
    startGame();
  });
</script>

<svelte:head>
  <title>Level Mudah - Game Berhitung</title>
</svelte:head>

<div
  class="min-h-screen bg-slate-50 flex flex-col items-center pt-24 pb-10 px-4 font-sans relative overflow-hidden"
>
  <div
    class="absolute top-[-10%] left-[-10%] w-96 h-96 bg-yellow-200 rounded-full mix-blend-multiply filter blur-3xl opacity-60 animate-blob"
  ></div>
  <div
    class="absolute bottom-[-10%] right-[-10%] w-96 h-96 bg-pink-300 rounded-full mix-blend-multiply filter blur-3xl opacity-60 animate-blob animation-delay-2000"
  ></div>

  <div class="w-full max-w-md z-10">
    <!-- Header -->
    <div class="flex items-center justify-between mb-8">
      <a
        href="/berhitung"
        class="p-2 bg-white rounded-full shadow-sm text-pink-500 hover:text-pink-600 hover:shadow transition-all border border-pink-100"
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
      </a>
      <div class="flex gap-1">
        {#each Array(10) as _, i}
          <div
            class="w-3 h-3 rounded-full {i < questionNumber - 1
              ? 'bg-pink-500'
              : 'bg-pink-200'}"
          ></div>
        {/each}
      </div>
    </div>

    <!-- Instruksi Visual -->
    <div
      class="bg-white rounded-3xl p-6 shadow-xl border-4 border-pink-100 text-center mb-8 relative"
    >
      <h2
        class="text-2xl font-black text-pink-500 mb-6 tracking-wide select-none"
      >
        Berapa jumlahnya?
      </h2>

      <!-- Area Gambar -->
      <div
        class="flex flex-wrap justify-center gap-2 mb-6 min-h-[120px] items-center"
      >
        {#if showSuccess}
          <div class="text-7xl animate-bounce select-none">🌟</div>
        {:else}
          {#each Array(currentCount) as _}
            <div
              class="text-5xl animate-pop-in cursor-pointer hover:scale-110 transition-transform select-none"
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
          class="select-none [-webkit-tap-highlight-color:transparent] [-webkit-touch-callout:none] aspect-square bg-white hover:bg-pink-50 border-4 border-pink-200 rounded-3xl shadow-[0_8px_0_0_rgba(251,207,232,1)] active:shadow-[0_0px_0_0_rgba(251,207,232,1)] active:translate-y-2 flex items-center justify-center text-5xl font-black text-pink-500 transition-all {errorShake
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
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-pink-900/40 backdrop-blur-md animate-fade-in"
      >
        <div
          class="bg-white p-8 rounded-[2rem] shadow-2xl max-w-sm w-full text-center border-4 border-pink-200 transform transition-all scale-100 animate-pop"
        >
          <div class="text-7xl mb-4 animate-bounce">🏆</div>
          <h2 class="text-4xl font-black text-pink-500 mb-2">HOREEE!</h2>
          <p class="text-slate-600 font-bold text-lg mb-8">
            Adik sangat pintar! Betul <span class="text-pink-600 text-2xl"
              >{score}</span
            > dari 10.
          </p>

          <div class="flex flex-col gap-4">
            <button
              onclick={startGame}
              class="w-full py-4 bg-pink-500 hover:bg-pink-600 text-white text-xl font-black rounded-2xl shadow-[0_6px_0_0_rgba(219,39,119,1)] active:shadow-none active:translate-y-[6px] transition-all"
            >
              Main Lagi
            </button>
            <a
              href="/berhitung"
              class="w-full py-4 bg-yellow-400 hover:bg-yellow-500 text-yellow-900 text-xl font-black rounded-2xl shadow-[0_6px_0_0_rgba(202,138,4,1)] active:shadow-none active:translate-y-[6px] transition-all no-underline block"
            >
              Pulang
            </a>
          </div>
        </div>
      </div>
    {/if}
  </div>
</div>

<audio bind:this={audioSuccess} src="/sounds/benar.mp3" preload="auto"></audio>
<audio bind:this={audioError} src="/sounds/wrong.mp3" preload="auto"></audio>

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
