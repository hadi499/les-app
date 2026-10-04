<script lang="ts">
  let gameState = $state("start");
  let countdown = $state(0);
  let message = $state("");
  let targetCard = $state<any>(null);
  let cards = $state<any[]>([]);
  let selectedTime = $state(10);
  let isCorrect = $state<boolean | null>(null);

  // 1. TAMBAHKAN BINDING AUDIO
  let audioCorrect: HTMLAudioElement;
  let audioWrong: HTMLAudioElement;

  const images = [
    "/images/doraemon/d1.png",
    "/images/doraemon/d2.png",
    "/images/doraemon/d3.png",
    "/images/doraemon/d4.png",
  ];

  // 2. FUNGSI PEMUTAR SUARA YANG AMAN
  function playSound(audioElement: HTMLAudioElement | undefined) {
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
    let shuffled = [...images];
    for (let i = shuffled.length - 1; i > 0; i--) {
      const j = Math.floor(Math.random() * (i + 1));
      [shuffled[i], shuffled[j]] = [shuffled[j], shuffled[i]];
    }

    cards = shuffled.map((img, index) => ({
      id: index + 1,
      number: index + 1,
      imageUrl: img,
    }));

    targetCard = cards[Math.floor(Math.random() * cards.length)];
    gameState = "memorize";
    countdown = selectedTime;
    message = "Hafalkan posisi gambar-gambar ini!";
    isCorrect = null;
  }

  $effect(() => {
    if (gameState === "memorize" && countdown > 0) {
      const timer = setTimeout(() => countdown--, 1000);
      return () => clearTimeout(timer);
    } else if (gameState === "memorize" && countdown === 0) {
      gameState = "guessing";
      message = "Di nomor berapakah gambar ini berada?";
    }
  });

  function guess(card: any) {
    if (gameState !== "guessing") return;

    if (card.id === targetCard.id) {
      message = `Benar sekali, gambar itu ada di nomor ${card.number}.`;
      isCorrect = true;
      playSound(audioCorrect); // 3. MAINKAN SUARA BENAR
    } else {
      message = `Salah, Gambar yang benar ada di nomor ${targetCard.number}.`;
      isCorrect = false;
      playSound(audioWrong); // 4. MAINKAN SUARA SALAH
    }
    gameState = "result";
  }
</script>

<div class="max-w-xl mx-auto pt-20 sm:pt-24 pb-12 px-4 sm:px-6 min-h-screen">
  <div class="flex items-center justify-between mb-4 sm:mb-6">
    <a
      href="/game-edukasi/tebak-4-gambar"
      class="p-1.5 sm:p-2 bg-white rounded-full shadow-sm text-blue-500 hover:text-blue-600 hover:shadow transition-all border border-slate-200"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="w-5 h-5 sm:w-6 sm:h-6"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2.5"
        stroke-linecap="round"
        stroke-linejoin="round"><path d="m15 18-6-6 6-6" /></svg
      >
    </a>
  </div>
  <div
    class="bg-gray-50 rounded-2xl shadow-lg text-center select-none p-4 sm:p-6 border border-slate-200"
  >
    <h1 class="text-xl sm:text-2xl font-bold">Game 4 Gambar</h1>
    <h2 class="text-base sm:text-lg text-orange-600 font-semibold">Doraemon</h2>

    <div class="min-h-[220px] flex flex-col items-center justify-center">
      {#if gameState === "start"}
        <div class="flex flex-col items-center gap-3 sm:gap-5 w-full max-w-sm mx-auto">
          <p class="text-sm sm:text-base font-bold text-gray-700">
            Sudah siap bermain?
          </p>

          <!-- Time selector -->
          <div class="bg-blue-50/80 p-3 sm:p-4 rounded-2xl w-full border border-blue-100">
            <p class="mb-2.5 text-[11px] sm:text-xs font-bold text-blue-700 uppercase tracking-widest">
              Waktu Menghafal
            </p>

            <div class="flex justify-center gap-1.5 sm:gap-2">
              {#each [10, 7, 5] as time}
                <button
                  class={`px-3 py-1.5 sm:px-5 sm:py-2 text-xs sm:text-sm rounded-full font-bold border-2 transition-all
                  ${
                    selectedTime === time
                      ? "bg-blue-600 text-white border-blue-600 scale-105 shadow-md"
                      : "bg-white text-blue-600 border-blue-200 hover:bg-blue-100"
                  }`}
                  onclick={() => (selectedTime = time)}
                >
                  {time} Detik
                </button>
              {/each}
            </div>
          </div>

          <button
            class="bg-orange-500 hover:bg-orange-600 active:translate-y-1 text-white px-6 py-2 sm:px-8 sm:py-3 text-base sm:text-lg font-bold rounded-full shadow-md"
            onclick={startGame}
          >
            ▶ MULAI MAIN
          </button>
        </div>
      {:else if gameState === "memorize"}
        <p class="text-base sm:text-lg font-bold text-gray-700">
          {message} <br />
          <span class="text-red-500 text-xl sm:text-2xl">
            Waktu: {countdown} detik
          </span>
        </p>
      {:else}
        <p
          class={`text-base sm:text-lg font-bold mb-2 ${
            gameState === "result"
              ? isCorrect
                ? "text-blue-600"
                : "text-red-600"
              : "text-gray-700"
          }`}
        >
          {message}
        </p>

        {#if targetCard}
          <div class="my-3">
            <img
              src={targetCard.imageUrl}
              class="w-20 h-20 sm:w-28 sm:h-28 mx-auto bg-white border-4 border-dashed border-yellow-400 rounded-xl p-2"
              draggable="false"
              alt=""
            />
          </div>
        {/if}

        {#if gameState === "result"}
          <button
            class="bg-green-600 hover:bg-green-700 text-white px-5 py-2 rounded-full font-bold mb-4 transition"
            onclick={startGame}
          >
            Main Lagi
          </button>
        {/if}
      {/if}
    </div>

    {#if gameState !== "start"}
      <div class="grid grid-cols-2 gap-6 sm:gap-8 justify-items-center mt-6 sm:mt-8 mx-auto w-full max-w-xs sm:max-w-md">
        {#each cards as card}
          <button
            class="w-full aspect-square max-w-[130px] sm:max-w-[160px] [perspective:1000px]"
            onclick={() => guess(card)}
            disabled={gameState !== "guessing"}
          >
            <div
              class={`relative w-full h-full transition-transform duration-500 [transform-style:preserve-3d]
              ${gameState === "guessing" ? "rotate-y-180" : ""}`}
            >
              <!-- Front -->
              <div
                class="absolute w-full h-full flex items-center justify-center bg-white border-[3px] sm:border-4 border-gray-200 rounded-xl shadow-md backface-hidden"
              >
                <img
                  src={card.imageUrl}
                  class="w-[70%] h-[70%] object-contain"
                  draggable="false"
                  alt=""
                />
              </div>

              <!-- Back -->
              <div
                class="absolute w-full h-full flex items-center justify-center bg-yellow-400 border-[3px] sm:border-4 border-yellow-500 rounded-xl shadow-md text-white text-4xl sm:text-5xl font-bold rotate-y-180 backface-hidden"
              >
                {card.number}
              </div>
            </div>
          </button>
        {/each}
      </div>
    {/if}
  </div>
</div>

<!-- 5. ELEMEN AUDIO DARI FOLDER LOKAL -->
<audio bind:this={audioCorrect} src="/sounds/benar.mp3" preload="auto"></audio>
<audio bind:this={audioWrong} src="/sounds/wrong.mp3" preload="auto"></audio>

<style>
  /* Tambahan kecil karena Tailwind belum support full 3D */
  .backface-hidden {
    backface-visibility: hidden;
  }
</style>
