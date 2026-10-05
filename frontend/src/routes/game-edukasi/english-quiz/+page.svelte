<script>
  import { onMount } from "svelte";

  // Data vocabulary untuk kuis
  const VOCABULARY = [
    { emoji: "🍎", word: "Apple", indonesian: "Apel", options: ["Apple", "Orange", "Banana", "Grape"] },
    { emoji: "🐶", word: "Dog", indonesian: "Anjing", options: ["Cat", "Dog", "Rabbit", "Bird"] },
    { emoji: "🚗", word: "Car", indonesian: "Mobil", options: ["Bus", "Train", "Car", "Bicycle"] },
    { emoji: "🐱", word: "Cat", indonesian: "Kucing", options: ["Dog", "Mouse", "Cat", "Cow"] },
    { emoji: "🏠", word: "House", indonesian: "Rumah", options: ["School", "House", "Hospital", "Park"] },
    { emoji: "☀️", word: "Sun", indonesian: "Matahari", options: ["Moon", "Star", "Cloud", "Sun"] },
    { emoji: "🌙", word: "Moon", indonesian: "Bulan", options: ["Sun", "Moon", "Earth", "Sky"] },
    { emoji: "🌲", word: "Tree", indonesian: "Pohon", options: ["Flower", "Grass", "Tree", "Leaf"] },
    { emoji: "📘", word: "Book", indonesian: "Buku", options: ["Pencil", "Bag", "Ruler", "Book"] },
    { emoji: "⚽", word: "Ball", indonesian: "Bola", options: ["Doll", "Kite", "Ball", "Toy"] }
  ];

  /** @type {any[]} */
  let questions = $state([]);
  let currentQuestionIndex = $state(0);
  let score = $state(0);
  let isGameOver = $state(false);
  let showSuccess = $state(false);
  let errorShake = $state(-1);
  let isFlipped = $state(false); // Untuk animasi kartu berbalik

  /** @param {any[]} array */
  function shuffleArray(array) {
    return [...array].sort(() => Math.random() - 0.5);
  }

  function startGame() {
    score = 0;
    currentQuestionIndex = 0;
    isGameOver = false;
    isFlipped = false;
    
    // Acak urutan pertanyaan dan pilihan jawabannya
    questions = shuffleArray(VOCABULARY).slice(0, 10).map(q => ({
      ...q,
      options: shuffleArray(q.options)
    }));
  }

  /**
   * @param {string} selectedWord
   * @param {number} index
   */
  function handleAnswer(selectedWord, index) {
    if (showSuccess || isFlipped) return;

    const currentQuestion = questions[currentQuestionIndex];

    if (selectedWord === currentQuestion.word) {
      // Benar
      new Audio('/sounds/benar.mp3').play().catch(e => console.error("Audio play failed:", e));
      score++;
      isFlipped = true; // Balik kartu untuk menampilkan arti bahasa Indonesia
      
      setTimeout(() => {
        showSuccess = true;
      }, 600);

      setTimeout(() => {
        showSuccess = false;
        isFlipped = false;
        if (currentQuestionIndex < questions.length - 1) {
          currentQuestionIndex++;
        } else {
          isGameOver = true;
        }
      }, 2000);
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

<div class="min-h-screen bg-blue-50 flex flex-col items-center pt-24 pb-10 px-4 font-sans relative overflow-hidden">
  
  <div class="absolute top-[-10%] left-[-10%] w-96 h-96 bg-blue-300 rounded-full mix-blend-multiply filter blur-3xl opacity-40 animate-blob"></div>
  <div class="absolute bottom-[-10%] right-[-10%] w-96 h-96 bg-purple-300 rounded-full mix-blend-multiply filter blur-3xl opacity-40 animate-blob animation-delay-2000"></div>

  <div class="w-full max-w-md z-10">
    <!-- Header: Tombol kembali dan indikator progres -->
    <div class="flex items-center justify-between mb-8">
      <a href="/game-edukasi" aria-label="Kembali ke menu game edukasi" class="p-2 bg-white rounded-full shadow-sm text-blue-500 hover:text-blue-600 hover:shadow transition-all border border-blue-100">
        <svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><path d="m15 18-6-6 6-6"/></svg>
      </a>
      <div class="flex gap-1">
        {#if questions.length > 0}
          {#each questions as _, i}
            <div class="w-3 h-3 rounded-full {i <= currentQuestionIndex ? 'bg-blue-500' : 'bg-blue-200'}"></div>
          {/each}
        {/if}
      </div>
    </div>

    <!-- Area Permainan Utama -->
    {#if questions.length > 0 && !isGameOver}
      <div class="bg-white rounded-3xl p-6 shadow-xl border-4 border-blue-100 text-center mb-8 relative">
        <h2 class="text-2xl font-black text-blue-500 mb-2 tracking-wide">English Quiz!</h2>
        <p class="text-blue-800 font-bold text-sm mb-6 opacity-80">Apa bahasa Inggris dari gambar ini?</p>
        
        <!-- Kartu (Flashcard) dengan efek Flip -->
        <div class="flip-card-container mb-8 aspect-square max-w-50 mx-auto">
          <div class="flip-card {isFlipped ? 'flipped' : ''}">
            <!-- Bagian Depan Kartu -->
            <div class="flip-card-front bg-blue-50 border-4 border-blue-200 rounded-3xl flex flex-col items-center justify-center shadow-lg">
              <span class="text-8xl animate-pop-in">{questions[currentQuestionIndex].emoji}</span>
            </div>
            <!-- Bagian Belakang Kartu (Hasil Benar) -->
            <div class="flip-card-back bg-green-50 border-4 border-green-200 rounded-3xl flex flex-col items-center justify-center shadow-lg">
              <span class="text-6xl mb-2">🎉</span>
              <span class="text-2xl font-black text-green-600">{questions[currentQuestionIndex].word}</span>
              <span class="text-lg font-bold text-green-500 opacity-80">{questions[currentQuestionIndex].indonesian}</span>
            </div>
          </div>
        </div>

        <!-- Tombol Pilihan Jawaban -->
        {#if !isFlipped}
          <div class="grid grid-cols-2 gap-4">
            {#each questions[currentQuestionIndex].options as option, index}
              <button
                onclick={() => handleAnswer(option, index)}
                class="py-4 px-2 bg-white hover:bg-blue-50 border-4 border-blue-200 rounded-2xl shadow-[0_6px_0_0_rgba(191,219,254,1)] active:shadow-none active:translate-y-1.5 flex items-center justify-center transition-all {errorShake === index ? 'animate-shake bg-red-50 border-red-200' : ''}"
              >
                <span class="text-xl font-black text-blue-700">{option}</span>
              </button>
            {/each}
          </div>
        {/if}
        
        <!-- Pesan Sukses Mengambang -->
        {#if showSuccess}
          <div class="absolute inset-0 flex items-center justify-center z-20 pointer-events-none">
            <div class="text-8xl animate-bounce">✨</div>
          </div>
        {/if}
      </div>
    {/if}

    <!-- Modal Game Over -->
    {#if isGameOver}
      <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-blue-900/40 backdrop-blur-md animate-fade-in">
        <div class="bg-white p-8 rounded-4xl shadow-2xl max-w-sm w-full text-center border-4 border-blue-200 transform transition-all scale-100 animate-pop">
          <div class="text-7xl mb-4 animate-bounce">🏆</div>
          <h2 class="text-4xl font-black text-blue-500 mb-2">LUAR BIASA!</h2>
          <p class="text-slate-600 font-bold text-lg mb-8">
            Skor Kamu: <span class="text-blue-600 text-2xl">{score}</span> dari {questions.length}
          </p>
          
          <div class="flex flex-col gap-4">
            <button
              onclick={startGame}
              class="w-full py-4 bg-blue-500 hover:bg-blue-600 text-white text-xl font-black rounded-2xl shadow-[0_6px_0_0_rgba(59,130,246,1)] active:shadow-none active:translate-y-1.5 transition-all"
            >
              Main Lagi
            </button>
            <a
              href="/game-edukasi"
              class="w-full py-4 bg-yellow-400 hover:bg-yellow-500 text-yellow-900 text-xl font-black rounded-2xl shadow-[0_6px_0_0_rgba(202,138,4,1)] active:shadow-none active:translate-y-1.5 transition-all"
            >
              Menu Utama
            </a>
          </div>
        </div>
      </div>
    {/if}

  </div>
</div>

<style>
  /* Animasi Flip Card */
  .flip-card-container {
    perspective: 1000px;
  }
  .flip-card {
    width: 100%;
    height: 100%;
    position: relative;
    transition: transform 0.6s;
    transform-style: preserve-3d;
  }
  .flip-card.flipped {
    transform: rotateY(180deg);
  }
  .flip-card-front, .flip-card-back {
    width: 100%;
    height: 100%;
    position: absolute;
    backface-visibility: hidden;
  }
  .flip-card-back {
    transform: rotateY(180deg);
  }

  /* Animasi Lainnya */
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
