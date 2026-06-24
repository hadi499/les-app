<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { slide, fade } from "svelte/transition";
  import { page } from "$app/state";
  import { goto, beforeNavigate } from "$app/navigation";

  let quizId = $derived(page.params.id);

  type Question = {
    id: number;
    question: string;
    options: string[];
    answer: number;
  };

  type Quiz = {
    id: number;
    title: string;
    category: string;
    timeLimit: number;
    questions: Question[];
  };

  let quiz = $state<Quiz | null>(null);
  let isLoading = $state(true);

  // State
  let currentQuestionIndex = $state(0);
  let userAnswers = $state<
    {
      question: string;
      answer: string | null;
      correct: string;
      isCorrect: boolean;
    }[]
  >([]);
  let isFinished = $state(false);
  let timeLeft = $state(15);
  let timerInterval: ReturnType<typeof setInterval>;
  let isSubmitting = $state(false);

  const currentQuestion = $derived(
    quiz && quiz.questions ? quiz.questions[currentQuestionIndex] : null,
  );
  const score = $derived(userAnswers.filter((a) => a.isCorrect).length * 10); // 10 points per correct answer

  beforeNavigate((navigation) => {
    if (quiz && currentQuestion && !isFinished) {
      const confirmLeave = confirm(
        "Anda sedang mengerjakan kuis. Jika Anda keluar sekarang, skor Anda akan dihitung 0. Apakah Anda yakin ingin keluar?",
      );
      if (!confirmLeave) {
        navigation.cancel();
      } else {
        // Submit score 0 saat user memilih keluar
        fetch(`/api/scores/quizzes`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
          keepalive: true,
          body: JSON.stringify({
            quiz_id: Number(quizId),
            score: 0,
          }),
        }).catch(console.error);
      }
    }
  });

  function handleBeforeUnload(e: BeforeUnloadEvent) {
    if (quiz && currentQuestion && !isFinished) {
      e.preventDefault();
      e.returnValue = "";
    }
  }

  onMount(async () => {
    window.addEventListener("beforeunload", handleBeforeUnload);
    try {
      // Cek apakah user sudah login
      const authRes = await fetch(`/me`, { credentials: "include" });
      const authData = await authRes.json();
      if (!authData.authenticated) {
        goto("/login");
        return;
      }

      const res = await fetch(`/api/quizzes/${quizId}`, {
        credentials: "include",
      });
      if (res.ok) {
        const json = await res.json();
        quiz = json.data;
        if (quiz && quiz.questions && quiz.questions.length > 0) {
          timeLeft = quiz.timeLimit;
          startTimer();
        } else {
          alert("Kuis tidak memiliki pertanyaan.");
          goto("/quiz");
        }
      } else {
        alert("Kuis tidak ditemukan.");
        goto("/quiz");
      }
    } catch (e) {
      console.error(e);
      alert("Gagal memuat kuis.");
      goto("/quiz");
    } finally {
      isLoading = false;
    }
  });

  onDestroy(() => {
    clearInterval(timerInterval);
    if (typeof window !== "undefined") {
      window.removeEventListener("beforeunload", handleBeforeUnload);
    }
  });

  function startTimer() {
    clearInterval(timerInterval);
    if (quiz) {
      timeLeft = quiz.timeLimit;
    }
    timerInterval = setInterval(() => {
      timeLeft--;
      if (timeLeft <= 0) {
        handleTimeout();
      }
    }, 1000);
  }

  function handleTimeout() {
    recordAnswer(null, -1);
  }

  function handleAnswer(optionText: string, optionIndex: number) {
    recordAnswer(optionText, optionIndex);
  }

  function recordAnswer(answerText: string | null, answerIndex: number) {
    clearInterval(timerInterval);
    if (!currentQuestion) return;

    userAnswers = [
      ...userAnswers,
      {
        question: currentQuestion.question,
        answer: answerText,
        correct: currentQuestion.options[currentQuestion.answer],
        isCorrect: answerIndex === currentQuestion.answer,
      },
    ];

    if (quiz && currentQuestionIndex < quiz.questions.length - 1) {
      currentQuestionIndex++;
      startTimer();
    } else {
      isFinished = true;
      submitScore();
    }
  }

  async function submitScore() {
    isSubmitting = true;
    console.log("Mencoba submit score:", {
      quiz_id: Number(quizId),
      score: score,
    });
    try {
      const res = await fetch(`/api/scores/quizzes`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({
          quiz_id: Number(quizId),
          score: score,
        }),
      });
      console.log("Response status:", res.status);
      if (!res.ok) {
        const errorText = await res.text();
        console.error("Gagal submit score, status:", res.status, errorText);
        alert(`Gagal submit score: ${res.status} - ${errorText}`);
      } else {
        console.log("Berhasil submit score!");
      }
    } catch (e) {
      console.error("Fetch error:", e);
      alert(
        "Error saat koneksi ke server untuk submit score: " +
          (e instanceof Error ? e.message : String(e)),
      );
    } finally {
      isSubmitting = false;
    }
  }

  function restartQuiz() {
    currentQuestionIndex = 0;
    userAnswers = [];
    isFinished = false;
    startTimer();
  }
</script>

<svelte:head>
  <title>Kuis | Les Balongarut</title>
</svelte:head>

<div
  class="min-h-screen bg-slate-50 font-sans selection:bg-blue-200 selection:text-blue-900 flex flex-col items-center py-12 px-4 relative overflow-x-hidden pt-24"
>
  <!-- Background Ambient -->
  <div class="absolute inset-0 z-0 pointer-events-none fixed">
    <div
      class="absolute top-1/4 left-1/4 w-[600px] h-[600px] bg-white/40 rounded-full blur-[120px]"
    ></div>
    <div
      class="absolute bottom-1/4 right-1/4 w-[500px] h-[500px] bg-blue-100/50 rounded-full blur-[120px]"
    ></div>

  </div>

  <div class="relative z-10 w-full max-w-2xl mx-auto flex flex-col gap-6">
    {#if isLoading}
      <div class="flex justify-center p-12">
        <div
          class="w-10 h-10 border-4 border-slate-200 border-t-blue-600 rounded-full animate-spin"
        ></div>
      </div>
    {:else if quiz && currentQuestion}
      {#if !isFinished}
        <div class="flex justify-between items-center mb-2">
          <a
            href="/quiz"
            class="text-sm font-bold tracking-[0.1em] text-slate-600 hover:text-slate-900 uppercase no-underline flex items-center gap-1"
          >
            <svg
              class="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M15 19l-7-7 7-7"
              />
            </svg>
            Kembali
          </a>
        </div>

        <!-- Header Info: Soal ke & Timer -->
        <div
          class="flex justify-between items-center bg-white/60 backdrop-blur-md rounded-2xl p-4 shadow-sm border border-slate-200"
        >
          <div
            class="text-sm font-bold tracking-[0.1em] text-slate-800 uppercase flex flex-col"
          >
            <span>{quiz.title}</span>
            <div class="flex items-center gap-1.5 mt-1.5">
              {#each quiz.questions as _, idx}
                <div
                  class="w-2 h-2 rounded-full transition-colors duration-300 {idx === currentQuestionIndex ? 'bg-blue-600 scale-125' : idx < currentQuestionIndex ? 'bg-blue-300' : 'bg-slate-200'}"
                  aria-hidden="true"
                ></div>
              {/each}
            </div>
          </div>
          <div class="flex items-center gap-2">
            <div
              class="text-xs font-semibold uppercase tracking-wider text-slate-600"
            >
              Waktu:
            </div>
            <div
              class="w-10 h-10 rounded-full bg-slate-100 border border-slate-300 flex items-center justify-center font-bold text-slate-800 shadow-sm {timeLeft <=
              5
                ? 'text-red-600 bg-red-50 border-red-200 animate-pulse'
                : ''}"
            >
              {timeLeft}
            </div>
          </div>
        </div>

        <!-- Question Card -->
        <div
          in:fade={{ duration: 300 }}
          class="bg-white/80 backdrop-blur-xl rounded-3xl p-8 shadow-md border border-slate-200 flex flex-col gap-8"
        >
          <h2
            class="text-xl sm:text-2xl font-bold text-slate-900 leading-tight"
          >
            {currentQuestion.question}
          </h2>

          <div class="flex flex-col gap-3">
            {#each currentQuestion.options as option, optIndex}
              <button
                onclick={() => handleAnswer(option, optIndex)}
                class="w-full flex items-center gap-4 text-left px-6 py-4 rounded-xl border-2 border-slate-200 hover:border-slate-300 bg-white hover:bg-slate-50 text-slate-800 font-semibold transition-colors duration-200 cursor-pointer shadow-sm"
              >
                <span
                  class="flex items-center justify-center w-8 h-8 rounded-lg bg-slate-200 text-slate-700 font-bold flex-shrink-0"
                >
                  {String.fromCharCode(65 + optIndex)}
                </span>
                <span>{option}</span>
              </button>
            {/each}
          </div>
        </div>
      {:else}
        <!-- Result Modal / Card -->
        <div
          in:slide
          class="bg-white/90 backdrop-blur-xl rounded-3xl p-8 shadow-xl border border-slate-300 flex flex-col gap-8 relative overflow-hidden"
        >
          <!-- Score Section -->
          <div
            class="flex flex-col items-center gap-2 text-center relative z-10 border-b border-slate-200 pb-8"
          >
            <div
              class="text-sm font-bold tracking-[0.2em] uppercase text-slate-600 mb-2"
            >
              Skor Akhir
            </div>
            <div
              class="text-[5rem] leading-none font-black text-blue-600 drop-shadow-sm"
            >
              {score}
            </div>
            <p class="text-slate-600 font-medium text-sm">
              Total Benar: {userAnswers.filter((a) => a.isCorrect).length} dari {quiz
                .questions.length} Soal
            </p>
            {#if isSubmitting}
              <p class="text-xs text-blue-500 font-medium animate-pulse">
                Menyimpan skor...
              </p>
            {/if}
          </div>

          <!-- Review Section -->
          <div class="flex flex-col gap-6 relative z-10">
            <h3
              class="text-lg font-bold tracking-[0.1em] text-slate-800 uppercase border-l-4 border-blue-400 pl-3"
            >
              Ulasan Jawaban
            </h3>

            <div
              class="flex flex-col gap-4 max-h-[50vh] overflow-y-auto pr-2 custom-scrollbar"
            >
              {#each userAnswers as ans, index}
                <div
                  class="p-4 rounded-xl border border-slate-200 bg-white/50 space-y-2"
                >
                  <div class="flex justify-between items-start gap-4">
                    <p class="font-semibold text-slate-800 text-sm m-0">
                      {index + 1}. {ans.question}
                    </p>
                    {#if ans.isCorrect}
                      <span
                        class="inline-flex items-center justify-center px-2 py-1 rounded text-xs font-bold bg-green-100 text-green-700"
                        >Benar</span
                      >
                    {:else}
                      <span
                        class="inline-flex items-center justify-center px-2 py-1 rounded text-xs font-bold bg-red-100 text-red-700"
                        >Salah</span
                      >
                    {/if}
                  </div>

                  <div class="text-sm flex flex-col gap-1 mt-2">
                    <div class="flex gap-2 items-center">
                      <span class="text-slate-500 font-medium"
                        >Jawaban Anda:</span
                      >
                      <span
                        class={ans.isCorrect
                          ? "text-green-700 font-semibold"
                          : "text-red-600 font-semibold"}
                      >
                        {ans.answer === null ? "Tidak Menjawab" : ans.answer}
                      </span>
                    </div>
                    {#if !ans.isCorrect}
                      <div class="flex gap-2 items-center">
                        <span class="text-slate-500 font-medium"
                          >Jawaban Benar:</span
                        >
                        <span class="text-green-700 font-semibold"
                          >{ans.correct}</span
                        >
                      </div>
                    {/if}
                  </div>
                </div>
              {/each}
            </div>
          </div>

          <div class="flex justify-center pt-4 relative z-10 gap-4">
            <a
              href="/quiz"
              class="inline-flex items-center justify-center px-6 py-3 text-xs tracking-[0.1em] font-bold uppercase text-slate-700 bg-white border border-slate-300 hover:border-slate-500 hover:text-slate-900 transition-all duration-300 rounded-lg shadow-sm cursor-pointer no-underline"
            >
              Kembali
            </a>
            <button
              onclick={restartQuiz}
              class="inline-flex items-center justify-center px-8 py-3 text-xs tracking-[0.2em] font-bold uppercase text-white bg-blue-600 border border-transparent hover:bg-blue-700 transition-all duration-300 rounded-lg overflow-hidden cursor-pointer shadow-md hover:shadow-lg"
            >
              Ulangi Kuis
            </button>
          </div>
        </div>
      {/if}
    {/if}
  </div>
</div>

<style>
  .custom-scrollbar::-webkit-scrollbar {
    width: 6px;
  }
  .custom-scrollbar::-webkit-scrollbar-track {
    background: rgba(241, 245, 249, 0.5); /* slate-100 */
    border-radius: 4px;
  }
  .custom-scrollbar::-webkit-scrollbar-thumb {
    background: rgba(148, 163, 184, 0.5); /* slate-400 */
    border-radius: 4px;
  }
  .custom-scrollbar::-webkit-scrollbar-thumb:hover {
    background: rgba(71, 85, 105, 0.8); /* slate-600 */
  }
</style>
