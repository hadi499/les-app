<script>
  import { onMount, onDestroy } from "svelte";


  const emojis = ["🚀", "🍎", "🐶", "🍕", "🎮", "⚽", "🚗", "🎸"];
  
  let cards = $state([]);
  let flippedIndices = $state([]);
  let matchedIndices = $state([]);
  let moves = $state(0);
  let time = $state(60); 
  let isPlaying = $state(false);
  let isGameOver = $state(false);
  let isWin = $state(false);
  let timerInterval;

  // Memulai atau mengulang permainan
  function startGame() {
    // Gandakan array emoji lalu acak (shuffle)
    const shuffled = [...emojis, ...emojis]
      .sort(() => Math.random() - 0.5)
      .map((emoji, id) => ({ id, emoji }));
      
    cards = shuffled;
    flippedIndices = [];
    matchedIndices = [];
    moves = 0;
    time = 60;
    isPlaying = true;
    isGameOver = false;
    isWin = false;

    if (timerInterval) clearInterval(timerInterval);
    
    timerInterval = setInterval(() => {
      time--;
      if (time <= 0) {
        clearInterval(timerInterval);
        isGameOver = true;
        isPlaying = false;
      }
    }, 1000);
  }

  // Menangani klik pada kartu
  function handleCardClick(index) {
    // Abaikan klik jika game over, belum mulai, kartu sudah dicocokkan, kartu sudah dibalik, atau 2 kartu sedang dibalik
    if (
      !isPlaying || 
      isGameOver || 
      matchedIndices.includes(index) || 
      flippedIndices.includes(index) || 
      flippedIndices.length >= 2
    ) {
      return;
    }

    // Balikkan kartu
    flippedIndices = [...flippedIndices, index];

    // Jika 2 kartu terbuka, cek kecocokan
    if (flippedIndices.length === 2) {
      moves++;
      const [firstIndex, secondIndex] = flippedIndices;
      
      if (cards[firstIndex].emoji === cards[secondIndex].emoji) {
        // Cocok!
        matchedIndices = [...matchedIndices, firstIndex, secondIndex];
        flippedIndices = [];
        
        // Cek apakah menang
        if (matchedIndices.length === cards.length) {
          clearInterval(timerInterval);
          isWin = true;
          isGameOver = true;
          isPlaying = false;
        }
      } else {
        // Tidak cocok, tutup kembali setelah 1 detik
        setTimeout(() => {
          flippedIndices = [];
        }, 1000);
      }
    }
  }

  // Format detik menjadi MM:SS
  const formatTime = (seconds) => {
    const m = Math.floor(seconds / 60).toString().padStart(2, "0");
    const s = (seconds % 60).toString().padStart(2, "0");
    return `${m}:${s}`;
  };

  onMount(() => {
    // Tunggu sebentar sebelum auto-start agar UI termuat sepenuhnya
    setTimeout(() => {
      startGame();
    }, 300);
  });

  onDestroy(() => {
    if (timerInterval) clearInterval(timerInterval);
  });
</script>

<div class="min-h-screen bg-slate-50 flex flex-col items-center pt-24 pb-10 px-4 font-sans relative overflow-hidden">
  
  <!-- Hiasan Background Background -->
  <div class="absolute top-[-10%] left-[-10%] w-96 h-96 bg-sky-200 rounded-full mix-blend-multiply filter blur-3xl opacity-50 animate-blob"></div>
  <div class="absolute bottom-[-10%] right-[-10%] w-96 h-96 bg-emerald-200 rounded-full mix-blend-multiply filter blur-3xl opacity-50 animate-blob animation-delay-2000"></div>

  <div class="w-full max-w-md z-10">
    <!-- Header -->
    <div class="flex items-center justify-between mb-8">
      <a href="/game-edukasi" class="p-2 bg-white rounded-full shadow-sm text-slate-500 hover:text-sky-500 hover:shadow transition-all border border-slate-200">
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="m15 18-6-6 6-6"/></svg>
      </a>
      <h1 class="text-2xl font-extrabold text-slate-800 tracking-tight">Latihan Memori</h1>
      <button onclick={startGame} class="p-2 bg-white rounded-full shadow-sm text-slate-500 hover:text-emerald-500 hover:shadow transition-all border border-slate-200" title="Ulangi Permainan">
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"/><path d="M3 3v5h5"/></svg>
      </button>
    </div>

    <!-- HUD (Heads Up Display) -->
    <div class="bg-white/80 backdrop-blur-md p-4 rounded-2xl shadow-sm border border-slate-200/60 flex justify-between items-center mb-8">
      <div class="text-center px-4">
        <div class="text-xs font-bold text-slate-400 uppercase tracking-widest mb-1">Waktu</div>
        <div class="text-2xl font-black {time <= 10 ? 'text-red-500 animate-pulse' : 'text-slate-700'}">
          {formatTime(time)}
        </div>
      </div>
      <div class="w-px h-10 bg-slate-200"></div>
      <div class="text-center px-4">
        <div class="text-xs font-bold text-slate-400 uppercase tracking-widest mb-1">Langkah</div>
        <div class="text-2xl font-black text-sky-600">
          {moves}
        </div>
      </div>
    </div>

    <!-- Papan Permainan -->
    <div class="grid grid-cols-4 gap-3 md:gap-4 mb-8 perspective-1000">
      {#each cards as card, index}
        {@const isFlipped = flippedIndices.includes(index) || matchedIndices.includes(index)}
        {@const isMatched = matchedIndices.includes(index)}
        
        <button
          class="aspect-square relative focus:outline-none focus-visible:ring-4 focus-visible:ring-sky-300 rounded-2xl group"
          onclick={() => handleCardClick(index)}
          disabled={!isPlaying || isFlipped}
        >
          <div class="w-full h-full transition-transform duration-500 transform-style-3d {isFlipped ? 'rotate-y-180' : ''}">
            
            <!-- Bagian Belakang Kartu (yang menghadap atas saat tertutup) -->
            <div class="absolute inset-0 w-full h-full bg-gradient-to-br from-sky-400 to-indigo-500 rounded-2xl shadow-md border-2 border-white/20 flex items-center justify-center backface-hidden group-hover:shadow-lg transition-shadow">
              <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" class="text-white/50"><path d="M12 2v20"/><path d="M2 12h20"/></svg>
            </div>

            <!-- Bagian Depan Kartu (yang ada Emoji) -->
            <div class="absolute inset-0 w-full h-full bg-white rounded-2xl shadow-md border-2 {isMatched ? 'border-emerald-400 bg-emerald-50' : 'border-slate-200'} flex items-center justify-center backface-hidden rotate-y-180">
              <span class="text-4xl md:text-5xl filter {isMatched ? 'drop-shadow-md' : ''} transition-all">
                {card.emoji}
              </span>
            </div>

          </div>
        </button>
      {/each}
    </div>

    <!-- Layar Hasil Akhir -->
    {#if isGameOver}
      <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-sm animate-fade-in">
        <div class="bg-white p-8 rounded-3xl shadow-2xl max-w-sm w-full text-center border border-slate-100 transform transition-all scale-100 animate-pop">
          <div class="text-6xl mb-4">
            {isWin ? '🎉' : '⏱️'}
          </div>
          <h2 class="text-2xl font-black text-slate-900 mb-2">
            {isWin ? 'Luar Biasa!' : 'Waktu Habis!'}
          </h2>
          <p class="text-slate-500 font-medium mb-6">
            {isWin 
              ? `Kamu menyelesaikan permainan dalam ${60 - time} detik dengan ${moves} langkah.` 
              : `Kamu hampir saja menyelesaikannya! Terus berlatih ya.`}
          </p>
          
          <div class="flex flex-col gap-3">
            <button
              onclick={startGame}
              class="w-full py-3.5 bg-sky-500 hover:bg-sky-600 text-white font-bold rounded-xl shadow-lg shadow-sky-500/30 transition-all hover:-translate-y-0.5"
            >
              Main Lagi
            </button>
            <a
              href="/game-edukasi"
              class="w-full py-3.5 bg-slate-100 hover:bg-slate-200 text-slate-700 font-bold rounded-xl transition-all"
            >
              Kembali ke Beranda
            </a>
          </div>
        </div>
      </div>
    {/if}

  </div>
</div>

<style>
  /* Kustomisasi CSS untuk Animasi 3D Card Flip */
  .perspective-1000 {
    perspective: 1000px;
  }
  .transform-style-3d {
    transform-style: preserve-3d;
  }
  .backface-hidden {
    backface-visibility: hidden;
  }
  .rotate-y-180 {
    transform: rotateY(180deg);
  }

  /* Animasi Background Blob */
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

  /* Animasi Popup Result */
  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }
  @keyframes pop {
    0% { transform: scale(0.9); opacity: 0; }
    100% { transform: scale(1); opacity: 1; }
  }
  .animate-fade-in {
    animation: fadeIn 0.3s ease-out forwards;
  }
  .animate-pop {
    animation: pop 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275) forwards;
  }
</style>

