<script lang="ts">
  import { onMount } from "svelte";
  import { fetchFlashcards, fetchFlashcardCategories } from "$lib/api_english_flashcard";
  import type { EnglishFlashcard, EnglishFlashcardCategory } from "$lib/api_english_flashcard";
  
  let questions = $state<EnglishFlashcard[]>([]);
  let categories = $state<EnglishFlashcardCategory[]>([]);
  
  let selectedCategory = $state<EnglishFlashcardCategory | null>(null);
  let hasStarted = $state(false);
  let isLoading = $state(true);
  
  let currentQuestionIndex = $state(0);
  let score = $state(0);
  let attempts = $state(0); // maksimal 1
  let wrongSelections = $state<number[]>([]); // Indeks jawaban salah agar tombol disabled
  
  let isGameOver = $state(false);
  let isFlipped = $state(false); 
  let isFlippingBack = $state(false);
  let isSuccess = $state(false); // Untuk warna kartu saat di-flip
  
  let errorShake = $state(-1);

  function shuffleArray<T>(array: T[]): T[] {
    return [...array].sort(() => Math.random() - 0.5);
  }

  async function loadCategories() {
    isLoading = true;
    try {
      categories = await fetchFlashcardCategories();
    } catch(e) {
      console.error(e);
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    loadCategories();
  });

  async function selectCategoryAndStart(category: EnglishFlashcardCategory) {
    selectedCategory = category;
    isLoading = true;
    
    try {
      const data = await fetchFlashcards(category.id);
      if (data.length === 0) {
        alert("Belum ada soal untuk kategori ini.");
        selectedCategory = null;
        isLoading = false;
        return;
      }
      
      // Acak urutan pertanyaan dan batasi 10 soal
      questions = shuffleArray(data)
        .slice(0, 10)
        .map((q) => ({
          ...q,
          options: shuffleArray(q.options),
        }));

      score = 0;
      currentQuestionIndex = 0;
      isGameOver = false;
      hasStarted = true;
      resetCard();
    } catch (e) {
      console.error(e);
      alert("Gagal memuat soal dari server.");
      selectedCategory = null;
    } finally {
      isLoading = false;
    }
  }
  
  function backToMenu() {
    hasStarted = false;
    selectedCategory = null;
    questions = [];
  }

  function resetCard() {
    attempts = 0;
    wrongSelections = [];
    isFlipped = false;
    isFlippingBack = false;
    isSuccess = false;
  }

  function handleAnswer(selectedOption: string, index: number) {
    if (isFlipped || wrongSelections.includes(index)) return;

    const currentQuestion = questions[currentQuestionIndex];

    if (selectedOption === currentQuestion.answer) {
      // Benar
      new Audio("/sounds/benar.mp3")
        .play()
        .catch((e) => console.error("Audio play failed:", e));
      isSuccess = true;
      isFlipped = true;

      // Skor: 1 jika tebakan benar
      score += 1;
    } else {
      // Salah
      attempts++;
      wrongSelections = [...wrongSelections, index];
      errorShake = index;

      setTimeout(() => {
        errorShake = -1;
      }, 400);

      if (attempts >= 1) {
        // Kesempatan habis
        isSuccess = false;
        isFlipped = true;
      }
    }
  }

  function nextQuestion() {
    if (currentQuestionIndex < questions.length - 1) {
      isFlippingBack = true;
      isFlipped = false;
      
      setTimeout(() => {
        currentQuestionIndex++;
        resetCard();
      }, 300);
    } else {
      isGameOver = true;
    }
  }
</script>

<div
  class="min-h-screen bg-teal-50 flex flex-col items-center pt-24 pb-10 px-4 font-sans relative overflow-hidden"
>
  <div
    class="absolute top-[-10%] left-[-10%] w-96 h-96 bg-teal-300 rounded-full mix-blend-multiply filter blur-3xl opacity-40 animate-blob"
  ></div>
  <div
    class="absolute bottom-[-10%] right-[-10%] w-96 h-96 bg-green-300 rounded-full mix-blend-multiply filter blur-3xl opacity-40 animate-blob animation-delay-2000"
  ></div>

  <div class="w-full max-w-4xl z-10 mt-6 md:mt-12">
    {#if !hasStarted}
      <!-- Layar Pemilihan Kategori (List Kartu) -->
      <div class="text-center animate-fade-in relative z-20">
        <h1 class="text-4xl font-black text-teal-800 mb-4 drop-shadow-sm">
          English Flashcard
        </h1>
        <p class="text-slate-600 mb-8 font-medium text-lg">
          Pilih kategori untuk memulai permainan. Kamu punya 2 kesempatan
          menebak untuk setiap kartu!
        </p>

        {#if isLoading}
          <div class="py-8 text-teal-600 font-bold animate-pulse text-xl">
            Memuat daftar kategori...
          </div>
        {:else if categories.length === 0}
          <div class="py-12 bg-white/60 rounded-3xl backdrop-blur-sm border border-teal-100 shadow-sm max-w-xl mx-auto">
            <div class="text-6xl mb-4 opacity-50">📂</div>
            <p class="text-slate-600 font-bold">Belum ada kategori yang ditambahkan.</p>
          </div>
        {:else}
          <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
            {#each categories as cat}
              <button
                onclick={() => selectCategoryAndStart(cat)}
                class="group text-left bg-white p-6 rounded-3xl shadow-lg border-b-8 border-r-4 border-teal-200 hover:border-teal-400 hover:shadow-xl hover:-translate-y-2 transition-all duration-300 flex flex-col h-full"
              >
                <div class="text-4xl mb-4 group-hover:scale-110 transition-transform origin-left">📚</div>
                <h3 class="text-2xl font-black text-slate-800 mb-2 group-hover:text-teal-700 transition-colors">
                  {cat.name}
                </h3>
                <p class="text-slate-500 font-medium text-sm line-clamp-2">
                  {cat.description || "Mulai belajar kosakata di kategori ini"}
                </p>
                <div class="mt-auto pt-6 flex justify-end w-full">
                  <span class="bg-teal-100 text-teal-700 font-bold px-4 py-1.5 rounded-full text-sm group-hover:bg-teal-500 group-hover:text-white transition-colors flex items-center gap-1">
                    Mainkan
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M14 5l7 7m0 0l-7 7m7-7H3" />
                    </svg>
                  </span>
                </div>
              </button>
            {/each}
          </div>
        {/if}

        <div class="mt-12">
          <a
            href="/game-edukasi"
            class="inline-flex items-center gap-2 text-teal-700 hover:text-teal-900 font-black text-lg bg-white/50 px-6 py-3 rounded-full shadow-sm hover:bg-white transition-all backdrop-blur-sm"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
              <path stroke-linecap="round" stroke-linejoin="round" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
            Kembali ke Menu Utama
          </a>
        </div>
      </div>
    {:else}
      <!-- Header -->
      <div class="flex items-center justify-between mb-8 max-w-xl mx-auto">
        <button
          onclick={backToMenu}
          aria-label="Kembali ke menu kategori"
          class="p-2 bg-white rounded-full shadow-sm text-teal-500 hover:text-teal-600 hover:shadow transition-all border border-teal-100"
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
        </button>

        <div
          class="bg-white px-4 py-2 rounded-2xl shadow-sm border border-teal-100 flex gap-4 font-bold text-teal-700"
        >
          <div>Soal: {currentQuestionIndex + 1}/{questions.length}</div>
          <div>Skor: {score}</div>
        </div>
      </div>

      {#if questions.length > 0 && !isGameOver}
        <div class="max-w-xl mx-auto">
          <div class="text-center mb-6">
            <h2 class="text-2xl font-black text-slate-800 mb-2">{selectedCategory?.name}</h2>
            {#if selectedCategory?.description}
              <p class="text-slate-600 font-medium">{selectedCategory.description}</p>
            {/if}
          </div>

          <!-- Area Flashcard (Berbentuk persegi panjang mirip kartu kredit di desktop, lebih tinggi di mobile) -->
          <div
            class="flip-card-container mb-8 w-full max-w-125 mx-auto h-[380px] sm:h-auto sm:aspect-[1.586/1]"
          >
            <div class="flip-card {isFlipped ? 'flipped' : ''}">
              <!-- Bagian Depan (Pertanyaan) -->
              <div
                class="flip-card-front bg-white border-4 border-teal-200 rounded-3xl p-6 sm:p-8 flex flex-col items-center justify-center shadow-xl"
              >
                <h2
                  class="text-xl sm:text-2xl font-black text-slate-800 text-center leading-relaxed"
                >
                  {questions[currentQuestionIndex].question}
                </h2>
              </div>

              <!-- Bagian Belakang (Jawaban & Penjelasan) -->
              <div
                class="flip-card-back border-4 rounded-3xl p-4 sm:p-6 flex flex-col items-center justify-center shadow-xl {isSuccess
                  ? 'bg-green-50 border-green-300'
                  : 'bg-red-50 border-red-300'}"
              >
                <div class="text-4xl sm:text-5xl mb-2 animate-pop-in shrink-0">
                  {isSuccess ? "🎉" : "💡"}
                </div>
                <h3
                  class="text-lg sm:text-xl font-black mb-2 shrink-0 {isSuccess
                    ? 'text-green-600'
                    : 'text-red-600'}"
                >
                  {isSuccess ? "Hebat! Jawaban Benar" : "Jawaban Benar:"}
                </h3>
                <div
                  class="text-xl sm:text-2xl font-black text-slate-800 mb-2 sm:mb-3 bg-white px-5 py-1.5 rounded-2xl shadow-sm text-center shrink-0"
                >
                  {questions[currentQuestionIndex].answer}
                </div>
                <div class="w-full flex-1 min-h-0 overflow-y-auto mt-1 sm:mt-2">
                  <p
                    class="text-slate-600 font-bold text-center text-xs sm:text-sm bg-white/60 px-3 py-1.5 sm:p-3 rounded-xl m-0"
                  >
                    {questions[currentQuestionIndex].explanation}
                  </p>
                </div>
              </div>
            </div>
          </div>

          <!-- Pilihan Jawaban / Tombol Lanjut -->
          {#if !isFlipped && !isFlippingBack}
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              {#each questions[currentQuestionIndex].options as option, index}
                {@const isWrong = wrongSelections.includes(index)}
                <button
                  onclick={() => handleAnswer(option, index)}
                  disabled={isWrong}
                  class="py-4 px-6 bg-white hover:bg-teal-50 border-4 rounded-2xl flex items-center justify-center transition-all
                  {isWrong
                    ? 'border-slate-200 bg-slate-50 opacity-50 shadow-none'
                    : 'border-teal-200 shadow-[0_6px_0_0_rgba(153,246,228,1)] active:shadow-none active:translate-y-1.5'} 
                  {errorShake === index
                    ? 'animate-shake border-red-300 bg-red-50'
                    : ''}"
                >
                  <span
                    class="text-xl font-black {isWrong
                      ? 'text-slate-400 line-through'
                      : 'text-teal-700'}">{option}</span
                  >
                </button>
              {/each}
            </div>
          {:else if isFlipped}
            <div class="animate-fade-in mt-2 flex justify-center">
              <button
                onclick={nextQuestion}
                class="w-full sm:w-auto px-12 py-4 bg-teal-500 hover:bg-teal-600 text-white text-xl font-black rounded-2xl shadow-[0_6px_0_0_rgba(20,184,166,1)] active:shadow-none active:translate-y-1.5 transition-all"
              >
                Lanjut
              </button>
            </div>
          {/if}
        </div>
      {/if}

      <!-- Modal Game Over -->
      {#if isGameOver}
        <div
          class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-teal-900/40 backdrop-blur-md animate-fade-in"
        >
          <div
            class="bg-white p-8 rounded-4xl shadow-2xl max-w-sm w-full text-center border-4 border-teal-200 transform transition-all scale-100 animate-pop"
          >
            <div class="text-7xl mb-4 animate-bounce">🎓</div>
            <h2 class="text-4xl font-black text-teal-600 mb-2">SELESAI!</h2>
            <p class="text-slate-600 font-bold text-lg mb-8">
              Total Skormu: <span class="text-teal-600 text-3xl block mt-2"
                >{score}</span
              >
            </p>

            <div class="flex flex-col gap-4">
              <button
                onclick={() => selectedCategory && selectCategoryAndStart(selectedCategory)}
                class="w-full py-4 bg-teal-500 hover:bg-teal-600 text-white text-xl font-black rounded-2xl shadow-[0_6px_0_0_rgba(20,184,166,1)] active:shadow-none active:translate-y-1.5 transition-all"
              >
                Main Lagi Kategori Ini
              </button>
              <button
                onclick={backToMenu}
                class="w-full py-4 bg-amber-400 hover:bg-amber-500 text-amber-900 text-xl font-black rounded-2xl shadow-[0_6px_0_0_rgba(251,191,36,1)] active:shadow-none active:translate-y-1.5 transition-all"
              >
                Pilih Kategori Lain
              </button>
            </div>
          </div>
        </div>
      {/if}
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
    transition: transform 0.6s cubic-bezier(0.4, 0.2, 0.2, 1);
    transform-style: preserve-3d;
  }
  .flip-card.flipped {
    transform: rotateY(180deg);
  }
  .flip-card-front,
  .flip-card-back {
    width: 100%;
    height: 100%;
    position: absolute;
    backface-visibility: hidden;
    left: 0;
    top: 0;
  }
  .flip-card-back {
    transform: rotateY(180deg);
  }

  /* Animasi Lainnya */
  @keyframes popIn {
    0% {
      transform: scale(0);
      opacity: 0;
    }
    80% {
      transform: scale(1.1);
      opacity: 1;
    }
    100% {
      transform: scale(1);
      opacity: 1;
    }
  }
  .animate-pop-in {
    animation: popIn 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275) backwards;
    animation-delay: 0.3s;
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
