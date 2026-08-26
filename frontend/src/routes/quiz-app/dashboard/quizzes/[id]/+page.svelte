<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { slide, fade } from "svelte/transition";
  import { page } from "$app/state";
  import { goto, beforeNavigate } from "$app/navigation";
  import Modal from "$lib/components/Modal.svelte";
  import katex from "katex";
  import "katex/dist/katex.min.css";

  let quizId = $derived(page.params.id);

  type Question = {
    id: number;
    question: string;
    image?: string;
    options: string[];
    answer: number;
  };

  type Quiz = {
    id: number;
    title: string;
    category: string;
    timeLimit: number; // in minutes
    questions: Question[];
  };

  let quiz = $state<Quiz | null>(null);
  let isLoading = $state(true);

  // State
  let currentQuestionIndex = $state(0);
  let userAnswers = $state<
    {
      questionId: string;
      question: string;
      image?: string;
      answer: string | null;
      answerIndex: number | null;
      correct: string;
      isCorrect: boolean;
    }[]
  >([]);
  let isFinished = $state(false);
  let timeLeft = $state(0);
  let timerInterval: ReturnType<typeof setInterval>;
  let isSubmitting = $state(false);
  let showLeaveModal = $state(false);
  let targetUrl = $state<string | null>(null);
  let hasConfirmedLeave = $state(false);
  let quizStartTime = $state(0);

  let finalScore = $state(0); // Score returned by KuisApp backend

  let showErrorModal = $state(false);
  let errorMessage = $state("");
  let errorRedirectTarget = $state("/quiz-app/dashboard");

  function closeErrorModal() {
    showErrorModal = false;
    goto(errorRedirectTarget);
  }

  const currentQuestion = $derived(
    quiz && quiz.questions ? quiz.questions[currentQuestionIndex] : null,
  );
  
  // Calculate score locally just for the UI review if needed, but we'll use finalScore from backend.
  const scoreLocal = $derived(userAnswers.filter((a) => a.isCorrect).length * 10);

  beforeNavigate((navigation) => {
    if (quiz && currentQuestion && !isFinished && !hasConfirmedLeave) {
      navigation.cancel();
      targetUrl = navigation.to?.url.pathname || "/quiz-app/dashboard";
      showLeaveModal = true;
    }
  });

  async function confirmLeave() {
    showLeaveModal = false;
    hasConfirmedLeave = true;

    // Submit score 0 saat user memilih keluar
    try {
      await fetch(`/api/kuisapp/quizzes/${quizId}/submit`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({
          answers: {},
          duration: Math.floor((Date.now() - quizStartTime) / 1000)
        }),
      });
    } catch (e) {
      console.error(e);
    }

    if (targetUrl) {
      goto(targetUrl);
    }
  }

  function cancelLeave() {
    showLeaveModal = false;
    targetUrl = null;
  }

  function handleBeforeUnload(e: BeforeUnloadEvent) {
    if (quiz && currentQuestion && !isFinished && !hasConfirmedLeave) {
      e.preventDefault();
      e.returnValue = "";
    }
  }

  function handlePageHide() {
    if (quiz && currentQuestion && !isFinished && !hasConfirmedLeave) {
      hasConfirmedLeave = true;
      fetch(`/api/kuisapp/quizzes/${quizId}/submit`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({
          answers: {},
          duration: Math.floor((Date.now() - quizStartTime) / 1000)
        }),
        keepalive: true,
      }).catch(console.error);
    }
  }

  function handleVisibilityChange() {
    if (document.visibilityState === "hidden") {
      handlePageHide();
    }
  }

  onMount(async () => {
    window.addEventListener("beforeunload", handleBeforeUnload);
    window.addEventListener("pagehide", handlePageHide);
    document.addEventListener("visibilitychange", handleVisibilityChange);
    
    try {
      // Fetch KuisApp user profile
      const authRes = await fetch(`/api/kuisapp/me?t=${Date.now()}`, { credentials: "include", cache: "no-store" });
      const authData = await authRes.json();
      if (!authRes.ok) {
        goto("/quiz-app/login");
        return;
      }

      // Fetch Quiz details
      const res = await fetch(`/api/kuisapp/quizzes/${quizId}`, {
        credentials: "include",
      });
      if (res.ok) {
        const json = await res.json();
        quiz = json.data;
        if (quiz && quiz.questions && quiz.questions.length > 0) {
          // The user requested: waktu kuis berdasarkan detik per soal
          timeLeft = quiz.timeLimit;
          quizStartTime = Date.now();
          startTimer();
        } else {
          errorMessage = "Kuis tidak memiliki pertanyaan.";
          showErrorModal = true;
        }
      } else {
        const errorData = await res.json().catch(() => ({}));
        errorMessage = errorData.error || "Kuis tidak ditemukan.";
        showErrorModal = true;
      }
    } catch (e) {
      console.error(e);
      errorMessage = "Gagal memuat kuis.";
      showErrorModal = true;
    } finally {
      isLoading = false;
    }
  });

  onDestroy(() => {
    clearInterval(timerInterval);
    if (typeof window !== "undefined") {
      window.removeEventListener("beforeunload", handleBeforeUnload);
      window.removeEventListener("pagehide", handlePageHide);
      document.removeEventListener("visibilitychange", handleVisibilityChange);
    }
  });

  function startTimer() {
    clearInterval(timerInterval);
    timerInterval = setInterval(() => {
      timeLeft--;
      if (timeLeft <= 0) {
        handleTimeout();
      }
    }, 1000);
  }

  function formatTime(seconds: number) {
    if (seconds < 0) return "0 dtk";
    if (seconds < 60) return `${seconds} dtk`;
    const m = Math.floor(seconds / 60);
    const s = seconds % 60;
    if (s === 0) return `${m} mnt`;
    return `${m}m ${s}dtk`;
  }

  function handleTimeout() {
    recordAnswer(null, -1);
  }

  function handleAnswer(optionText: string, optionIndex: number) {
    recordAnswer(optionText, optionIndex);
  }

  function recordAnswer(answerText: string | null, answerIndex: number) {
    if (!currentQuestion) return;

    userAnswers = [
      ...userAnswers,
      {
        questionId: currentQuestion.id.toString(),
        question: currentQuestion.question,
        image: currentQuestion.image,
        answer: answerText,
        answerIndex: answerIndex !== -1 ? answerIndex : null,
        correct: currentQuestion.options[currentQuestion.answer],
        isCorrect: answerIndex === currentQuestion.answer,
      },
    ];

    if (quiz && currentQuestionIndex < quiz.questions.length - 1) {
      currentQuestionIndex++;
      timeLeft = quiz.timeLimit; // Reset timer for the next question
    } else {
      clearInterval(timerInterval);
      isFinished = true;
      submitScore();
    }
  }

  async function submitScore() {
    isSubmitting = true;
    const duration = Math.floor((Date.now() - quizStartTime) / 1000);
    
    // Build answers map for KuisApp backend
    const answersPayload: Record<string, number> = {};
    userAnswers.forEach((ua) => {
      if (ua.answerIndex !== null) {
        answersPayload[ua.questionId] = ua.answerIndex;
      }
    });

    try {
      const res = await fetch(`/api/kuisapp/quizzes/${quizId}/submit`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({
          answers: answersPayload,
          duration: duration
        }),
      });

      if (res.status === 401) {
        errorMessage = "Sesi Anda telah berakhir. Silakan login kembali.";
        errorRedirectTarget = "/quiz-app/login";
        showErrorModal = true;
        return;
      }

      if (!res.ok) {
        let errorMsg = `Gagal submit score: ${res.status}`;
        try {
          const errorJson = await res.json();
          errorMsg = errorJson.error || errorMsg;
        } catch (err) {}
        errorMessage = errorMsg;
        showErrorModal = true;
      } else {
        const json = await res.json();
        const finalScore = json.score || 0;
        const pointsEarned = json.points_earned || 0;
        const pointsAlreadyClaimed = json.points_already_claimed || false;
        
        // Simpan hasil ke sessionStorage untuk ditampilkan di halaman result
        sessionStorage.setItem(`kuisapp_result_${quizId}`, JSON.stringify({
          finalScore,
          pointsEarned,
          pointsAlreadyClaimed,
          userAnswers,
          quizTitle: quiz?.title,
          quizQuestionsLength: quiz?.questions.length,
          duration
        }));
        
        // Redirect ke halaman result terpisah
        goto(`/quiz-app/dashboard/quizzes/${quizId}/result`);
      }
    } finally {
      isSubmitting = false;
    }
  }

  function restartQuiz() {
    currentQuestionIndex = 0;
    userAnswers = [];
    isFinished = false;
    finalScore = 0;
    if (quiz) {
      timeLeft = quiz.timeLimit * 60;
      startTimer();
    }
  }

  function renderText(text: string | null) {
    if (!text) return "";

    // Replace block math $$...$$
    let rendered = text.replace(/\$\$([\s\S]*?)\$\$/g, (match, math) => {
      try {
        return katex.renderToString(math, {
          displayMode: true,
          throwOnError: false,
        });
      } catch (e) {
        return match;
      }
    });

    // Replace inline math $...$
    rendered = rendered.replace(/\$([^$]*?)\$/g, (match, math) => {
      try {
        return katex.renderToString(math, {
          displayMode: false,
          throwOnError: false,
        });
      } catch (e) {
        return match;
      }
    });

    return rendered;
  }
</script>

<svelte:head>
  <title>{quiz ? quiz.title : 'Kuis'} | KuisApp</title>
</svelte:head>

<div class="w-full flex flex-col items-center py-6 sm:py-12 relative overflow-visible">

  <div class="relative z-10 w-full max-w-2xl mx-auto flex flex-col gap-6">
    {#if isLoading}
      <div class="flex justify-center p-12">
        <div class="w-10 h-10 border-4 border-slate-200 border-t-blue-600 rounded-full animate-spin"></div>
      </div>
    {:else if quiz && currentQuestion}


      {#if !isFinished}
        <!-- Header Info: Soal ke & Timer -->
        <div class="sticky top-0 sm:top-16 z-50 flex justify-between items-center py-4 px-2 sm:px-0 mb-6 sm:mb-8 bg-slate-50">
          <div class="text-sm font-bold tracking-[0.1em] text-slate-800 uppercase flex flex-col">
            <span>{quiz.title}</span>
            <div class="flex items-center gap-1.5 mt-1.5 flex-wrap max-w-[200px] sm:max-w-xs">
              {#each quiz.questions as _, idx}
                <div
                  class="w-2 h-2 rounded-full transition-colors duration-300 {idx === currentQuestionIndex ? 'bg-indigo-600 scale-125 shadow-sm shadow-indigo-200' : idx < currentQuestionIndex ? 'bg-indigo-300' : 'bg-slate-200'}"
                  aria-hidden="true"
                ></div>
              {/each}
            </div>
          </div>
          <div class="flex items-center gap-2">
            <div class="px-3 h-10 rounded-full bg-slate-100 border border-slate-300 flex items-center justify-center font-bold text-slate-800 shadow-sm {timeLeft <= 60 ? 'text-red-600 bg-red-50 border-red-200' : ''}">
              {formatTime(timeLeft)}
            </div>
          </div>
        </div>

        <!-- Question Card -->
        <div in:fade={{ duration: 300 }} class="bg-white rounded-3xl p-8 shadow-xl border border-slate-200 flex flex-col gap-8 relative overflow-hidden">
          <div class="absolute top-0 left-0 w-full h-2 bg-gradient-to-r from-indigo-500 to-blue-500"></div>
          {#if currentQuestion.image}
            <div class="w-[85%] md:w-[50%] lg:w-[40%] mx-auto flex justify-center mb-2">
              <img
                src={currentQuestion.image}
                alt="Gambar Pertanyaan"
                class="max-w-full h-auto rounded-xl shadow-sm object-contain max-h-[300px] md:max-h-[250px]"
              />
            </div>
          {/if}

          <h2 class="text-xl sm:text-2xl font-bold text-slate-900 leading-tight">
            {@html renderText(currentQuestion.question)}
          </h2>

          <div class="flex flex-col gap-3">
            {#each currentQuestion.options as option, optIndex}
              <button
                onclick={() => handleAnswer(option, optIndex)}
                class="w-full flex items-center gap-4 text-left px-6 py-4 rounded-xl border-2 border-slate-200 bg-slate-50 text-slate-800 font-semibold transition-all duration-200 cursor-pointer shadow-sm group"
              >
                <span class="flex items-center justify-center w-8 h-8 rounded-lg bg-white border border-slate-300 text-slate-600 font-bold flex-shrink-0 transition-colors">
                  {String.fromCharCode(65 + optIndex)}
                </span>
                <span>{@html renderText(option)}</span>
              </button>
            {/each}
          </div>
        </div>
      {/if}
    {/if}
  </div>

  <!-- Leave Confirmation Modal -->
  <Modal show={showLeaveModal} onclose={cancelLeave}>
    <div class="space-y-4 text-center">
      <div class="w-12 h-12 mx-auto rounded-full bg-red-100 flex items-center justify-center">
        <svg class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
      </div>
      <div>
        <h3 class="text-lg font-semibold text-slate-900">Konfirmasi Keluar</h3>
        <p class="text-sm text-slate-600 mt-1">
          Anda sedang mengerjakan kuis. Jika Anda keluar sekarang, skor Anda akan dihitung 0. Apakah Anda yakin ingin keluar?
        </p>
      </div>
      <div class="flex gap-2 justify-center pt-2">
        <button onclick={cancelLeave} class="px-4 py-2 text-sm rounded-lg border border-slate-300 hover:bg-slate-50 text-slate-900 cursor-pointer font-medium">
          Batal
        </button>
        <button onclick={confirmLeave} class="px-4 py-2 text-sm rounded-lg bg-red-500 text-white hover:bg-red-600 cursor-pointer font-medium shadow-sm">
          Ya, Keluar
        </button>
      </div>
    </div>
  </Modal>

  <!-- Error Modal -->
  <Modal show={showErrorModal} onclose={closeErrorModal}>
    <div class="text-center">
      <div class="w-12 h-12 rounded-full bg-red-100 mx-auto mb-4 flex items-center justify-center">
        <svg class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
      </div>
      <div>
        <h3 class="text-lg font-semibold text-slate-900">Akses Ditolak</h3>
        <p class="text-sm text-slate-600 mt-2">
          {errorMessage}
        </p>
      </div>
      <div class="flex gap-2 justify-center pt-5">
        <button onclick={closeErrorModal} class="px-5 py-2.5 text-sm rounded-xl bg-blue-600 text-white hover:bg-blue-700 cursor-pointer font-bold shadow-sm transition-colors">
          Kembali ke Daftar Kuis
        </button>
      </div>
    </div>
  </Modal>
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
