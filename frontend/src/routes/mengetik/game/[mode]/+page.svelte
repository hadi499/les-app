<script lang="ts">
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";
  import {
    getCurrentUserId,
    saveGameHighScore,
    getGameHighScores,
  } from "$lib/stores/users";

  type GameMode = "left" | "right" | "both" | "letters" | "all";

  const validModes: GameMode[] = ["left", "right", "both", "letters", "all"];

  // Safe parameter extraction - params sudah di-decode oleh SvelteKit
  const mode = $derived.by(() => {
    const param = page.params?.mode;
    if (!param || typeof param !== "string") {
      return "left";
    }
    // Validasi tanpa decode lagi (sudah di-decode oleh SvelteKit)
    if (!validModes.includes(param as GameMode)) {
      return "left";
    }
    return param as GameMode;
  });

  let gameState = $state<"ready" | "playing" | "paused" | "gameover">("ready");
  let score = $state(0);
  let lives = $state(5);
  let level = $state(1);
  let letters = $state<
    { id: number; char: string; x: number; y: number; speed: number }[]
  >([]);
  let nextLetterId = $state(0);
  let spawnInterval: ReturnType<typeof setInterval> | null = null;
  let animationFrame: number | null = null;
  let lastTime = $state(0);
  let keysPressed = new Set<string>();

  const keyMap = {
    left: ["a", "s", "d", "f"],
    right: ["j", "k", "l", ";"],
    both: ["a", "s", "d", "f", "j", "k", "l", ";"],
    letters: [
      "a",
      "b",
      "c",
      "d",
      "e",
      "f",
      "g",
      "h",
      "i",
      "j",
      "k",
      "l",
      "m",
      "n",
      "o",
      "p",
      "q",
      "r",
      "s",
      "t",
      "u",
      "v",
      "w",
      "x",
      "y",
      "z",
    ],
    all: [
      "a",
      "b",
      "c",
      "d",
      "e",
      "f",
      "g",
      "h",
      "i",
      "j",
      "k",
      "l",
      "m",
      "n",
      "o",
      "p",
      "q",
      "r",
      "s",
      "t",
      "u",
      "v",
      "w",
      "x",
      "y",
      "z",
      "`",
      "-",
      "=",
      "[",
      "]",
      "\\",
      ";",
      "'",
      ",",
      ".",
      "/",
    ],
  };

  const allKeysDisplay = "A-Z `-=[]\\;',./";

  function getModeTitle(m: GameMode): string {
    const titles = {
      left: "Home Row Kiri",
      right: "Home Row Kanan",
      both: "Home Row Lengkap",
      letters: "Semua Huruf",
      all: "Semua Huruf & Karakter",
    };
    return titles[m];
  }

  function getSpawnInterval(): number {
    const baseInterval = 2000;
    const reduction = Math.floor(score / 100) * 150;
    return Math.max(600, baseInterval - reduction);
  }

  function getSpeedMultiplier(): number {
    return 1 + Math.floor(score / 100) * 0.15;
  }

  function spawnLetter() {
    const keys = keyMap[mode];
    const spawnCount = level >= 3 ? 2 : 1;
    const positions: number[] = [];
    const minDistance = 15; // Jarak minimal 15% lebar layar

    for (let i = 0; i < spawnCount; i++) {
      const char = keys[Math.floor(Math.random() * keys.length)];
      const baseSpeed = 0.3 + Math.random() * 0.2;
      const speed = baseSpeed * getSpeedMultiplier();

      // Generate posisi x yang tidak bertumpuk
      let x: number;
      let attempts = 0;
      do {
        x = 10 + Math.random() * 80;
        attempts++;
      } while (
        positions.some((pos) => Math.abs(pos - x) < minDistance) &&
        attempts < 10
      );

      positions.push(x);
      letters = [...letters, { id: nextLetterId++, char, x, y: 100, speed }];
    }
  }

  function startGame() {
    score = 0;
    lives = 5;
    level = 1;
    letters = [];
    nextLetterId = 0;
    gameState = "playing";
    lastTime = performance.now();
    keysPressed.clear();

    spawnInterval = setInterval(spawnLetter, getSpawnInterval());
    animationFrame = requestAnimationFrame(gameLoop);
  }

  function gameLoop(currentTime: number) {
    if (gameState !== "playing") return;

    const deltaTime = currentTime - lastTime;
    lastTime = currentTime;

    const newLevel = Math.floor(score / 100) + 1;
    if (newLevel !== level) {
      level = newLevel;
      if (spawnInterval) clearInterval(spawnInterval);
      spawnInterval = setInterval(spawnLetter, getSpawnInterval());
    }

    letters = letters
      .map((l) => ({ ...l, y: l.y - l.speed * (deltaTime / 16) }))
      .filter((l) => {
        if (l.y < -5) {
          lives--;
          if (lives <= 0) {
            endGame();
          }
          return false;
        }
        return true;
      });

    animationFrame = requestAnimationFrame(gameLoop);
  }

  function endGame() {
    gameState = "gameover";
    if (spawnInterval) {
      clearInterval(spawnInterval);
      spawnInterval = null;
    }
    if (animationFrame) {
      cancelAnimationFrame(animationFrame);
      animationFrame = null;
    }

    // Simpan high score untuk user saat ini
    const userId = getCurrentUserId();
    if (userId) {
      saveGameHighScore(userId, mode, score);
    }
  }

  function handleKeydown(e: KeyboardEvent) {
    if (gameState !== "playing") return;

    const key = e.key;
    const validKeys = keyMap[mode];

    const keyToCheck =
      key.length === 1 && /[a-zA-Z]/.test(key) ? key.toLowerCase() : key;

    if (!validKeys.includes(keyToCheck)) return;

    if (keysPressed.has(keyToCheck)) return;

    keysPressed.add(keyToCheck);

    if (keysPressed.size > 2) {
      endGame();
      return;
    }

    const index = letters.findIndex((l) => l.char === keyToCheck);
    if (index !== -1) {
      letters = letters.filter((_, i) => i !== index);
      score += 10;
    }
  }

  function handleKeyup(e: KeyboardEvent) {
    const key = e.key;
    const keyToCheck =
      key.length === 1 && /[a-zA-Z]/.test(key) ? key.toLowerCase() : key;
    keysPressed.delete(keyToCheck);
  }

  function cleanup() {
    if (spawnInterval) {
      clearInterval(spawnInterval);
      spawnInterval = null;
    }
    if (animationFrame) {
      cancelAnimationFrame(animationFrame);
      animationFrame = null;
    }
  }

  onMount(() => {
    if (!validModes.includes(mode)) {
      goto("/mengetik/game");
      return;
    }

    window.addEventListener("keydown", handleKeydown);
    window.addEventListener("keyup", handleKeyup);
    return () => {
      window.removeEventListener("keydown", handleKeydown);
      window.removeEventListener("keyup", handleKeyup);
      cleanup();
    };
  });
</script>

<svelte:head>
  <title>{getModeTitle(mode)} - Game Typing</title>
</svelte:head>

<div
  class="min-h-screen flex flex-col bg-slate-50 text-slate-900 overflow-hidden font-sans"
>
  <header
    class="flex flex-col md:flex-row items-center justify-between px-4 md:px-8 py-4 md:py-6 gap-2 md:gap-4 relative z-10"
  >
    <h1
      class="text-base md:text-xl font-bold m-0 flex-1 text-center drop-shadow-sm text-slate-700"
    >
      {getModeTitle(mode)}
    </h1>
    <div class="flex gap-4 md:gap-8">
      <div class="flex flex-col items-center">
        <span
          class="text-[0.625rem] md:text-xs text-slate-600 uppercase tracking-wider font-light"
          >Skor</span
        >
        <span
          class="text-base md:text-2xl font-bold text-amber-600 drop-shadow-sm"
          >{score}</span
        >
      </div>
      <div class="flex flex-col items-center">
        <span
          class="text-[0.625rem] md:text-xs text-slate-600 uppercase tracking-wider font-light"
          >Level</span
        >
        <span
          class="text-base md:text-2xl font-bold text-amber-600 drop-shadow-sm"
          >{level}</span
        >
      </div>
      <div class="flex flex-col items-center">
        <span
          class="text-[0.625rem] md:text-xs text-slate-600 uppercase tracking-wider font-light"
          >Nyawa</span
        >
        <span
          class="text-base md:text-2xl font-bold text-amber-600 drop-shadow-sm"
          >{"❤️".repeat(lives)}</span
        >
      </div>
    </div>
  </header>

  <main
    class="flex-1 relative w-full max-w-[900px] mx-auto px-4 md:px-8 my-4 md:my-8 z-10"
  >
    <div
      class="relative w-full h-[55vh] md:h-[calc(60vh+50px)] bg-slate-100/50 border border-slate-300 rounded-2xl overflow-hidden shadow-2xl shadow-slate-300/50 backdrop-blur-sm"
    >
      {#if gameState === "ready"}
        <div
          class="absolute inset-0 bg-slate-900/10 text-slate-900 flex items-center justify-center z-10 rounded-2xl backdrop-blur-sm"
        >
          <div class="text-center max-w-[400px] px-4">
            <h2
              class="text-xl md:text-3xl font-bold mb-2 text-slate-800 drop-shadow-md"
            >
              Siap Bermain?
            </h2>
            <p class="text-slate-600 text-xs md:text-sm mb-6">
              Ketik huruf yang melayang sebelum menyentuh garis atas!
            </p>
            <div
              class="bg-white/80 border border-slate-200 rounded-xl p-4 mb-6 shadow-inner"
            >
              <p class="text-[0.625rem] md:text-xs mb-2 text-slate-600">
                Huruf yang akan muncul:
              </p>
              {#if mode === "all"}
                <div
                  class="font-mono text-xs md:text-sm font-bold tracking-[0.15rem] leading-relaxed text-amber-600 drop-shadow-sm"
                >
                  {allKeysDisplay}
                </div>
              {:else}
                <div
                  class="font-mono text-base md:text-xl font-bold tracking-[0.5rem] text-amber-600 drop-shadow-sm"
                >
                  {keyMap[mode].map((k) => k.toUpperCase()).join(" ")}
                </div>
              {/if}
            </div>
            <button
              class="bg-blue-500 hover:bg-blue-400 text-white border border-blue-600 rounded-xl px-8 py-3 text-sm font-medium cursor-pointer transition-all hover:scale-105 hover:shadow-lg shadow-blue-900/5"
              onclick={startGame}>Mulai Game</button
            >
          </div>
        </div>
      {:else if gameState === "gameover"}
        <div
          class="absolute inset-0 bg-slate-900/10 text-slate-900 flex items-center justify-center z-10 rounded-2xl backdrop-blur-md"
        >
          <div class="text-center max-w-[400px] px-4">
            <h2
              class="text-xl md:text-3xl font-bold mb-2 bg-gradient-to-br from-red-400 to-amber-500 bg-clip-text text-transparent drop-shadow-sm"
            >
              Game Over!
            </h2>
            <div
              class="bg-white/80 border border-slate-200 rounded-xl p-5 mb-6 shadow-lg shadow-slate-900/10"
            >
              <div
                class="text-xs text-slate-600 uppercase tracking-widest mb-2 font-light"
              >
                Skor Akhir
              </div>
              <div
                class="text-5xl font-black text-amber-600 leading-none mb-2 drop-shadow-md"
              >
                {score}
              </div>
              <div class="text-sm text-slate-600 mb-4 font-light">
                Level {level}
              </div>
              {#if getCurrentUserId()}
                {@const highScores = getGameHighScores(getCurrentUserId()!)}
                <div
                  class="pt-3 border-t border-slate-200 flex justify-center items-center gap-2"
                >
                  <span class="text-xs text-slate-600 font-light"
                    >High Score:</span
                  >
                  <span class="text-lg font-bold text-emerald-700"
                    >{highScores[mode]}</span
                  >
                </div>
              {/if}
            </div>
            <div class="flex flex-col md:flex-row gap-3 justify-center">
              <button
                class="bg-blue-500 hover:bg-blue-400 text-white border border-blue-600 rounded-xl px-6 py-2.5 text-sm font-medium cursor-pointer transition-all hover:scale-105 hover:shadow-lg shadow-blue-900/5"
                onclick={startGame}>Main Lagi</button
              >
              <a
                href="/mengetik/game"
                class="bg-white hover:bg-slate-50 text-slate-800 border border-slate-200 rounded-xl px-6 py-2.5 text-sm font-medium no-underline transition-all hover:shadow-lg shadow-slate-900/5 inline-flex items-center justify-center"
                >Pilih Mode Lain</a
              >
            </div>
          </div>
        </div>
      {/if}

      <div
        class="absolute top-[5%] left-0 right-0 h-1 bg-gradient-to-r from-red-500 via-orange-500 to-yellow-500 shadow-[0_0_20px_rgba(239,68,68,0.5)] z-0"
      ></div>
      {#each letters as letter (letter.id)}
        <div
          class="absolute text-2xl md:text-4xl font-black text-blue-600 drop-shadow-[0_0_15px_rgba(251,191,36,0.6)] -translate-x-1/2 -translate-y-1/2 select-none pointer-events-none z-0"
          style="left: {letter.x}%; top: {letter.y}%"
        >
          {letter.char.toUpperCase()}
        </div>
      {/each}
    </div>
  </main>
</div>
