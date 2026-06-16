<script lang="ts">
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import {
    getLesson,
    lessons,
    KEY_MAP,
    FINGER_COLORS,
  } from "$lib/data/lessons";
  import { saveLessonProgress, calculateStars } from "$lib/stores/progress";
  import type { TypingStats } from "$lib/types";
  import Keyboard from "$lib/components/Keyboard.svelte";
  import TypingArea from "$lib/components/TypingArea.svelte";
  import ResultModal from "$lib/components/ResultModal.svelte";
  import Navbar from "$lib/components/Navbar.svelte";

  // Safe parameter extraction with validation
  const lessonId = $derived.by(() => {
    const param = page.params?.lessonId;
    const id = Number(param);
    return isNaN(id) ? 1 : id; // default to lesson 1 if invalid
  });
  const lesson = $derived(getLesson(lessonId));

  let mode = $state<"words" | "sentences">("words");
  let text = $state("");
  let currentIndex = $state(0);
  let errors = $state(new Set<number>());
  let startTime = $state<number | null>(null);
  let elapsed = $state(0);
  let finished = $state(false);
  let pressedKey = $state("");
  let timerInterval: ReturnType<typeof setInterval> | null = null;

  const activeKey = $derived(
    currentIndex < text.length ? text[currentIndex].toLowerCase() : "",
  );

  let stats: TypingStats = $state({
    wpm: 0,
    accuracy: 100,
    correctChars: 0,
    totalChars: 0,
    errors: 0,
    time: 0,
  });

  let resultStars = $state(0);
  let hasNext = $derived(lessons.some((l) => l.id === lessonId + 1));

  function generateText(): string {
    if (!lesson) return "";
    if (mode === "words") {
      const words: string[] = [];
      for (let i = 0; i < 15; i++) {
        words.push(
          lesson.words[Math.floor(Math.random() * lesson.words.length)],
        );
      }
      return words.join(" ");
    } else {
      const sentences: string[] = [];
      for (let i = 0; i < 3; i++) {
        sentences.push(
          lesson.sentences[Math.floor(Math.random() * lesson.sentences.length)],
        );
      }
      return sentences.join(" ");
    }
  }

  function startLesson() {
    text = generateText();
    currentIndex = 0;
    errors = new Set<number>();
    startTime = null;
    elapsed = 0;
    finished = false;
    resultStars = 0;
    stats = {
      wpm: 0,
      accuracy: 100,
      correctChars: 0,
      totalChars: 0,
      errors: 0,
      time: 0,
    };
    if (timerInterval) clearInterval(timerInterval);
  }

  function updateStats() {
    if (!startTime) return;
    const now = Date.now();
    elapsed = Math.floor((now - startTime) / 1000);
    const minutes = elapsed / 60;
    const wordsTyped = currentIndex / 5;
    const wpm = minutes > 0 ? Math.round(wordsTyped / minutes) : 0;
    const accuracy =
      stats.totalChars > 0
        ? Math.round((stats.correctChars / stats.totalChars) * 100)
        : 100;
    stats = { ...stats, wpm, accuracy, time: elapsed };
  }

  function handleKeyDown(e: KeyboardEvent) {
    if (finished) return;

    const key = e.key;
    pressedKey = key.toLowerCase();
    setTimeout(() => {
      pressedKey = "";
    }, 150);

    if (key === "Escape") {
      startLesson();
      return;
    }

    if (key === "Backspace") {
      e.preventDefault();
      if (currentIndex > 0) {
        currentIndex--;
        errors.delete(currentIndex);
      }
      return;
    }

    if (key.length !== 1) return;
    e.preventDefault();

    if (!startTime) {
      startTime = Date.now();
      timerInterval = setInterval(updateStats, 500);
    }

    const expected = text[currentIndex];
    const isCorrect = key.toLowerCase() === expected.toLowerCase();

    stats.totalChars++;
    if (isCorrect) {
      stats.correctChars++;
    } else {
      errors = new Set([...errors, currentIndex]);
      stats.errors++;
    }

    currentIndex++;

    updateStats();

    if (currentIndex >= text.length) {
      setTimeout(() => {
        finishLesson();
      }, 400);
    }
  }

  function finishLesson() {
    finished = true;
    if (timerInterval) clearInterval(timerInterval);
    updateStats();

    resultStars = calculateStars(stats.wpm, stats.accuracy);

    saveLessonProgress({
      lessonId,
      bestWPM: stats.wpm,
      bestAccuracy: stats.accuracy,
      completed: true,
      stars: resultStars,
      attempts: 0,
    });
  }

  function handleRetry() {
    startLesson();
  }

  function handleNext() {
    const nextLesson = lessons.find((l) => l.id === lessonId + 1);
    if (nextLesson) {
      goto(`/mengetik/practice/${nextLesson.id}`);
    }
  }

  function switchMode(newMode: "words" | "sentences") {
    mode = newMode;
    startLesson();
  }

  $effect(() => {
    if (lesson) {
      startLesson();
    }
    return () => {
      if (timerInterval) clearInterval(timerInterval);
    };
  });

  const highlightKeys = $derived(lesson?.keys ?? []);
</script>

<svelte:head>
  <title>{lesson?.title ?? "Latihan"} - Ketik 10 Jari</title>
</svelte:head>

<svelte:window onkeydown={handleKeyDown} />

{#if lesson}
  <div
    class="flex flex-col min-h-[calc(100vh-64px)] text-orange-950 bg-orange-50"
  >
    <header
      class="flex flex-col sm:flex-row items-center justify-between p-4 sm:px-8 gap-3 sm:gap-0"
    >
      <div class="flex items-center gap-3">
        <span class="text-3xl">{lesson.icon}</span>
        <h1 class="text-xl font-bold m-0 text-orange-950">{lesson.title}</h1>
      </div>
      <div
        class="flex gap-1 p-1 bg-white/60 border border-orange-300 rounded-xl shadow-sm"
      >
        <button
          class="px-4 py-2 border-none bg-transparent rounded-lg cursor-pointer font-semibold text-sm transition-all {mode ===
          'words'
            ? 'bg-amber-100 text-amber-600'
            : 'text-orange-800 hover:bg-white'}"
          onclick={() => switchMode("words")}
        >
          Kata
        </button>
        <button
          class="px-4 py-2 border-none bg-transparent rounded-lg cursor-pointer font-semibold text-sm transition-all {mode ===
          'sentences'
            ? 'bg-amber-100 text-amber-600'
            : 'text-orange-800 hover:bg-white'}"
          onclick={() => switchMode("sentences")}
        >
          Kalimat
        </button>
      </div>
    </header>

    <main
      class="flex-1 w-full max-w-4xl mx-auto p-4 sm:p-8 flex flex-col gap-8"
    >
      <TypingArea {text} {currentIndex} {errors} />

      <div class="flex flex-col items-center gap-6">
        <Keyboard {activeKey} {pressedKey} {highlightKeys} />

        {#if lesson.keys.length > 0}
          <div class="text-center">
            <p class="text-sm text-orange-800 m-0 mb-3 font-light">
              Jari untuk latihan ini:
            </p>
            <div class="flex justify-center gap-2">
              {#each lesson.keys as key}
                {@const info = KEY_MAP[key]}
                {#if info}
                  <span
                    class="w-9 h-9 flex items-center justify-center rounded-lg font-bold text-sm bg-white border-2 text-orange-950 shadow-sm"
                    style="border-color: {FINGER_COLORS[info.finger]}"
                  >
                    {key.toUpperCase()}
                  </span>
                {/if}
              {/each}
            </div>
          </div>
        {/if}
      </div>

      <div
        class="flex justify-center gap-6 text-orange-800 text-sm mt-4 font-light"
      >
        <span
          ><kbd
            class="bg-white px-2 py-0.5 rounded border border-orange-300 text-xs shadow-sm mr-1 text-orange-950"
            >Esc</kbd
          > Ulangi</span
        >
        <span
          ><kbd
            class="bg-white px-2 py-0.5 rounded border border-orange-300 text-xs shadow-sm mr-1 text-orange-950"
            >Backspace</kbd
          > Hapus</span
        >
      </div>
    </main>

    <ResultModal
      show={finished}
      wpm={stats.wpm}
      accuracy={stats.accuracy}
      stars={resultStars}
      onRetry={handleRetry}
      onNext={handleNext}
      {hasNext}
    />
  </div>
{:else}
  <div
    class="min-h-screen flex flex-col items-center justify-center gap-4 bg-gradient-to-br from-[#EAE4BD] to-[#EAD5B8] text-orange-950"
  >
    <h1 class="text-2xl font-bold">Lesson tidak ditemukan</h1>
    <a href="/mengetik" class="text-amber-500 hover:underline"
      >Kembali ke Menu</a
    >
  </div>
{/if}
