<script>
  import { onMount } from "svelte";


  // Pasangan emoji yang bentuk/warnanya mirip untuk melatih ketelitian
  const EMOJI_PAIRS = [
    ["😃", "😀"], ["🍎", "🍅"], ["🚗", "🚕"], ["🐶", "🦊"], 
    ["🐻", "🐼"], ["🌲", "🌳"], ["🌕", "🌝"], ["🍔", "🥪"],
    ["⚽", "🏀"], ["🍉", "🍓"], ["🦋", "🐝"], ["🌻", "🌼"]
  ];

  /** @type {string[]} */
  let grid = $state([]);
  let oddIndex = $state(-1);
  let score = $state(0);
  let questionNumber = $state(1);
  let isGameOver = $state(false);
  let showSuccess = $state(false);
  let errorShake = $state(-1);

  function startGame() {
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

    // Pilih pasangan emoji acak
    const pairIndex = Math.floor(Math.random() * EMOJI_PAIRS.length);
    const pair = EMOJI_PAIRS[pairIndex];
    
    // Tentukan mana yang mayoritas, mana yang beda (odd)
    const isFirstOdd = Math.random() > 0.5;
    const majorityEmoji = isFirstOdd ? pair[1] : pair[0];
    const oddEmoji = isFirstOdd ? pair[0] : pair[1];

    // Buat grid 6 kotak (mayoritas 5, beda 1)
    const newGrid = Array(6).fill(majorityEmoji);
    oddIndex = Math.floor(Math.random() * 6);
    newGrid[oddIndex] = oddEmoji;
    
    grid = newGrid;
  }

  /** @param {number} index */
  function handleAnswer(index) {
    if (showSuccess) return;

    if (index === oddIndex) {
      // Benar
      new Audio('/sounds/benar.mp3').play().catch(e => console.error("Audio play failed:", e));
      score++;
      showSuccess = true;
      
      setTimeout(() => {
        showSuccess = false;
        questionNumber++;
        nextQuestion();
      }, 1200);
    } else {
      // Salah
      errorShake = index;
      setTimeout(() => {
        errorShake = -1;
      }, 400);
    }
  }

  onMount(() => {
    startGame();
  });
</script>

<div class="min-h-screen bg-orange-50 flex flex-col items-center pt-24 pb-10 px-4 font-sans relative overflow-hidden">
  
  <div class="absolute top-[-10%] left-[-10%] w-96 h-96 bg-yellow-300 rounded-full mix-blend-multiply filter blur-3xl opacity-40 animate-blob"></div>
  <div class="absolute bottom-[-10%] right-[-10%] w-96 h-96 bg-orange-300 rounded-full mix-blend-multiply filter blur-3xl opacity-40 animate-blob animation-delay-2000"></div>

  <div class="w-full max-w-md z-10">
    <div class="flex items-center justify-between mb-8">
      <a href="/game-edukasi" class="p-2 bg-white rounded-full shadow-sm text-orange-500 hover:text-orange-600 hover:shadow transition-all border border-orange-100">
        <svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><path d="m15 18-6-6 6-6"/></svg>
      </a>
      <div class="flex gap-1">
        {#each Array(10) as _, i}
          <div class="w-3 h-3 rounded-full {i < questionNumber - 1 ? 'bg-orange-500' : 'bg-orange-200'}"></div>
        {/each}
      </div>
    </div>

    <div class="bg-white rounded-3xl p-6 shadow-xl border-4 border-orange-100 text-center mb-8">
      <h2 class="text-2xl font-black text-orange-500 mb-2 tracking-wide">Cari yang Beda!</h2>
      <p class="text-orange-800 font-bold text-sm mb-6 opacity-80">Gunakan mata elangmu untuk mencari 1 gambar yang berbeda dari teman-temannya.</p>
      
      {#if showSuccess}
        <div class="h-64 flex items-center justify-center">
          <div class="text-8xl animate-bounce">✨</div>
        </div>
      {:else}
        <div class="grid grid-cols-2 sm:grid-cols-3 gap-4">
          {#each grid as emoji, index}
            <button
              onclick={() => handleAnswer(index)}
              class="aspect-square bg-white hover:bg-orange-50 border-4 border-orange-200 rounded-3xl shadow-[0_6px_0_0_rgba(253,186,116,1)] active:shadow-none active:translate-y-[6px] flex items-center justify-center text-6xl transition-all {errorShake === index ? 'animate-shake bg-red-50 border-red-200' : ''}"
            >
              <span class="animate-pop-in" style="animation-delay: {index * 50}ms">{emoji}</span>
            </button>
          {/each}
        </div>
      {/if}
    </div>

    {#if isGameOver}
      <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-orange-900/40 backdrop-blur-md animate-fade-in">
        <div class="bg-white p-8 rounded-[2rem] shadow-2xl max-w-sm w-full text-center border-4 border-orange-200 transform transition-all scale-100 animate-pop">
          <div class="text-7xl mb-4 animate-bounce">🦅</div>
          <h2 class="text-4xl font-black text-orange-500 mb-2">MATA ELANG!</h2>
          <p class="text-slate-600 font-bold text-lg mb-8">
            Kamu berhasil menemukan <span class="text-orange-600 text-2xl">{score}</span> gambar yang bersembunyi.
          </p>
          
          <div class="flex flex-col gap-4">
            <button
              onclick={startGame}
              class="w-full py-4 bg-orange-500 hover:bg-orange-600 text-white text-xl font-black rounded-2xl shadow-[0_6px_0_0_rgba(249,115,22,1)] active:shadow-none active:translate-y-[6px] transition-all"
            >
              Main Lagi
            </button>
            <a
              href="/game-edukasi"
              class="w-full py-4 bg-yellow-400 hover:bg-yellow-500 text-yellow-900 text-xl font-black rounded-2xl shadow-[0_6px_0_0_rgba(202,138,4,1)] active:shadow-none active:translate-y-[6px] transition-all"
            >
              Pulang
            </a>
          </div>
        </div>
      </div>
    {/if}

  </div>
</div>

<style>
  @keyframes popIn {
    0% { transform: scale(0); opacity: 0; }
    80% { transform: scale(1.1); opacity: 1; }
    100% { transform: scale(1); opacity: 1; }
  }
  .animate-pop-in {
    animation: popIn 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275) backwards;
  }

  @keyframes shake {
    0%, 100% { transform: translateX(0); }
    25% { transform: translateX(-8px); }
    50% { transform: translateX(8px); }
    75% { transform: translateX(-8px); }
  }
  .animate-shake {
    animation: shake 0.3s cubic-bezier(.36,.07,.19,.97) both;
  }

  @keyframes blob {
    0% { transform: translate(0px, 0px) scale(1); }
    33% { transform: translate(30px, -50px) scale(1.1); }
    66% { transform: translate(-20px, 20px) scale(0.9); }
    100% { transform: translate(0px, 0px) scale(1); }
  }
  .animate-blob {
    animation: blob 7s infinite;
  }
  .animation-delay-2000 {
    animation-delay: 2s;
  }

  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }
  @keyframes pop {
    0% { transform: scale(0.8); opacity: 0; }
    100% { transform: scale(1); opacity: 1; }
  }
  .animate-fade-in {
    animation: fadeIn 0.3s ease-out forwards;
  }
  .animate-pop {
    animation: pop 0.5s cubic-bezier(0.175, 0.885, 0.32, 1.275) forwards;
  }
</style>
