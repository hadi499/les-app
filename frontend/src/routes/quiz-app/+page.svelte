<script lang="ts">
  import { getContext } from "svelte";

  const auth = getContext<{current: any, loading: boolean}>("quizUser");
</script>

<svelte:head>
  <title>LB Quiz | Les Balonggarut</title>
</svelte:head>

<div
  class="min-h-screen bg-slate-50 relative overflow-hidden font-sans selection:bg-indigo-200 selection:text-indigo-900"
>
  <!-- Dynamic Backgrounds -->
  <div
    class="absolute top-[-10%] left-[-10%] w-[50%] h-[50%] bg-indigo-200/50 rounded-full blur-[120px] mix-blend-multiply animate-blob"
  ></div>
  <div
    class="absolute top-[20%] right-[-10%] w-[40%] h-[40%] bg-purple-200/50 rounded-full blur-[120px] mix-blend-multiply animate-blob animation-delay-2000"
  ></div>
  <div
    class="absolute bottom-[-20%] left-[20%] w-[60%] h-[60%] bg-blue-200/50 rounded-full blur-[120px] mix-blend-multiply animate-blob animation-delay-4000"
  ></div>

  <!-- Content Container -->
  <div
    class="relative z-10 container mx-auto px-6 py-12 min-h-screen flex flex-col justify-center items-center"
  >
      <div class="text-center max-w-4xl mx-auto mb-16">
        <div
          class="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-white/50 backdrop-blur-md border border-white/60 shadow-sm mb-6 text-sm font-semibold text-indigo-700 uppercase tracking-widest"
        >
          <span class="relative flex h-3 w-3">
            <span
              class="animate-ping absolute inline-flex h-full w-full rounded-full bg-indigo-400 opacity-75"
            ></span>
            <span
              class="relative inline-flex rounded-full h-3 w-3 bg-indigo-500"
            ></span>
          </span>
          Portal Integrasi
        </div>

        <h1
          class="text-3xl py-3 md:text-7xl font-extrabold text-transparent bg-clip-text bg-linear-to-r from-indigo-800 via-purple-700 to-blue-800 mb-6 drop-shadow-sm"
        >
          Les Balonggarut Quiz
        </h1>
        <p
          class="text-lg md:text-xl text-slate-600 font-medium max-w-2xl mx-auto leading-relaxed"
        >
          Platform KuisApp terintegrasi dengan ekosistem Les Balongarut. Uji
          kemampuanmu, kumpulkan poin, dan raih peringkat teratas dengan cara
          yang menyenangkan.
        </p>
      </div>

      <div
        class="grid grid-cols-1 md:grid-cols-2 gap-8 w-full max-w-4xl mx-auto perspective-1000"
      >
        <!-- Les App Card -->
        <div class="group relative">
          <div
            class="absolute inset-0 bg-linear-to-br from-blue-400 to-cyan-300 rounded-3xl blur opacity-25 group-hover:opacity-50 transition duration-500"
          ></div>
          <div
            class="relative h-full bg-white/60 backdrop-blur-xl border border-white/80 rounded-3xl p-10 shadow-xl hover:shadow-2xl hover:-translate-y-2 transition-all duration-500 flex flex-col items-center text-center"
          >
            <div
              class="w-20 h-20 bg-linear-to-br from-blue-500 to-cyan-400 rounded-2xl flex items-center justify-center shadow-lg mb-6 group-hover:scale-110 transition-transform duration-500 text-white"
            >
              <svg
                class="w-10 h-10"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"
                ></path></svg
              >
            </div>
            <h2 class="text-2xl font-bold text-slate-800 mb-3">
              Les Balonggarut
            </h2>
            <p class="text-slate-600 mb-8 flex-grow">
              Akses kembali ke modul pembelajaran utama, latihan materi, dan
              manajemen kelas Anda.
            </p>
            <a
              href="/"
              class="w-full py-4 rounded-xl font-bold text-white bg-slate-800 hover:bg-slate-900 shadow-md hover:shadow-xl transition-all duration-300 flex items-center justify-center gap-2 relative overflow-hidden group/btn"
            >
              <span
                class="absolute inset-0 bg-white/20 translate-x-[-100%] group-hover/btn:translate-x-[100%] transition-transform duration-500"
              ></span>
              Kembali
              <svg
                class="w-5 h-5 group-hover/btn:translate-x-1 transition-transform"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M17 8l4 4m0 0l-4 4m4-4H3"
                ></path></svg
              >
            </a>
          </div>
        </div>

        <!-- Quiz App Card -->
        <div class="group relative">
          <div
            class="absolute inset-0 bg-linear-to-br from-indigo-500 to-purple-500 rounded-3xl blur opacity-25 group-hover:opacity-50 transition duration-500"
          ></div>
          <div
            class="relative h-full bg-white/60 backdrop-blur-xl border border-white/80 rounded-3xl p-10 shadow-xl hover:shadow-2xl hover:-translate-y-2 transition-all duration-500 flex flex-col items-center text-center"
          >
            <div
              class="w-20 h-20 bg-linear-to-br from-indigo-600 to-purple-500 rounded-2xl flex items-center justify-center shadow-lg mb-6 group-hover:scale-110 transition-transform duration-500 text-white"
            >
              <svg
                class="w-10 h-10"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z"
                ></path></svg
              >
            </div>
            <h2 class="text-2xl font-bold text-slate-800 mb-3">
              LB Quiz Portal
            </h2>
            <p class="text-slate-600 mb-8 flex-grow">
              Masuk ke arena kuis, kerjakan soal menantang, dan pantau statistik
              nilai Anda secara real-time.
            </p>
            <a
              href={auth?.current ? (auth.current.role === 'admin' ? '/quiz-app/admin' : '/quiz-app/dashboard') : "/quiz-app/login"}
              class="w-full py-4 rounded-xl font-bold text-white bg-linear-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 shadow-md hover:shadow-indigo-500/30 transition-all duration-300 flex items-center justify-center gap-2 relative overflow-hidden group/btn"
            >
              <span
                class="absolute inset-0 bg-white/20 translate-x-[-100%] group-hover/btn:translate-x-[100%] transition-transform duration-500"
              ></span>
              {auth?.current ? 'Dashboard' : 'Masuk'}
              <svg
                class="w-5 h-5 group-hover/btn:translate-x-1 transition-transform"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M13 10V3L4 14h7v7l9-11h-7z"
                ></path></svg
              >
            </a>
          </div>
        </div>
      </div>
  </div>
</div>

<style>
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
  .animation-delay-4000 {
    animation-delay: 4s;
  }
  .perspective-1000 {
    perspective: 1000px;
  }
</style>
