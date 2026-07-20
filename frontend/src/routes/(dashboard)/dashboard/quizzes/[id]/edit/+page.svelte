<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import { onMount } from "svelte";

  let quizId = $derived(page.params.id);

  let title = $state("");
  let category = $state("Matematika");
  let timeLimit = $state(15);
  let isPublished = $state(true);
  
  let questions: { question: string; options: string[]; answer: number }[] = $state([]);
  let lastResetAt = $state<string | null>(null);
  let isResetThisMonth = $derived.by(() => {
    if (!lastResetAt) return false;
    const resetDate = new Date(lastResetAt);
    const now = new Date();
    return resetDate.getFullYear() === now.getFullYear() && resetDate.getMonth() === now.getMonth();
  });
  let isLoading = $state(true);
  let showLoadingSpinner = $state(false);
  let isSubmitting = $state(false);

  let showResetModal = $state(false);
  let isResetting = $state(false);

  async function resetScores() {
    isResetting = true;
    try {
      const res = await fetch(`/api/quizzes/${quizId}/scores`, {
        method: "DELETE",
        credentials: "include",
      });
      if (res.ok) {
        alert("Riwayat nilai berhasil direset!");
        lastResetAt = new Date().toISOString();
        showResetModal = false;
      } else {
        const err = await res.json();
        alert("Gagal mereset nilai: " + (err.error || ""));
      }
    } catch (e) {
      alert("Terjadi kesalahan jaringan saat mereset nilai.");
    } finally {
      isResetting = false;
    }
  }

  async function fetchQuiz() {
    try {
      const res = await fetch(`/api/quizzes/${quizId}`, {
        credentials: "include",
      });
      if (res.ok) {
        const json = await res.json();
        const data = json.data;
        title = data.title;
        category = data.category;
        timeLimit = data.timeLimit;
        if (data.is_published !== undefined) {
          isPublished = data.is_published;
        }
        if (data.last_reset_at) {
          lastResetAt = data.last_reset_at;
        }
        if (data.questions && data.questions.length > 0) {
          questions = data.questions.map((q: any) => ({
            question: q.question,
            options: q.options || ["", "", "", ""],
            answer: q.answer || 0
          }));
        }
      } else {
        alert("Gagal memuat kuis");
        goto("/dashboard/quizzes");
      }
    } catch (e) {
      console.error(e);
      alert("Terjadi kesalahan");
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    fetchQuiz();
  });

  function addQuestion() {
    questions = [...questions, { question: "", options: ["", "", "", ""], answer: 0 }];
  }

  function removeQuestion(index: number) {
    questions = questions.filter((_, i) => i !== index);
  }

  async function handleSubmit(e?: Event) {
    if (e) e.preventDefault();
    isSubmitting = true;

    try {
      const res = await fetch(`/api/quizzes/${quizId}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({
          title,
          category,
          timeLimit: Number(timeLimit),
          is_published: isPublished,
          questions: questions.map(q => ({
            question: q.question,
            options: q.options,
            answer: Number(q.answer)
          }))
        })
      });

      if (res.ok) {
        goto("/dashboard/quizzes");
      } else {
        const err = await res.json();
        alert("Gagal: " + (err.error || "Terjadi kesalahan"));
      }
    } catch (err) {
      console.error(err);
      alert("Error menghubungi server");
    } finally {
      isSubmitting = false;
    }
  }
</script>

<svelte:head>
  <title>Edit Kuis | Les Balongarut</title>
</svelte:head>

<div class="sticky top-0 z-40 bg-slate-50/95 backdrop-blur-md py-3 sm:py-4 mb-6 sm:mb-8 border-b border-slate-200/60 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 -mx-4 px-4 sm:-mx-8 sm:px-8">
  <div class="flex items-center gap-4">
    <a
      href="/dashboard/quizzes"
      class="p-2 text-slate-600 bg-white/50 hover:bg-white rounded-xl border border-slate-200 shadow-sm transition-all"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
      </svg>
    </a>
    <div>
      <h1 class="text-2xl sm:text-3xl font-bold text-slate-900 mb-1">Edit Kuis</h1>
      <p class="text-sm text-slate-600 m-0 hidden sm:block">Perbarui informasi dan daftar pertanyaan.</p>
    </div>
  </div>

  <div class="flex items-center gap-3 w-full sm:w-auto">
    <button
      type="button"
      onclick={addQuestion}
      class="flex-1 sm:flex-none inline-flex items-center justify-center gap-2 px-4 py-2 bg-white text-blue-600 border border-blue-200 rounded-xl font-bold hover:bg-blue-50 transition-colors shadow-sm cursor-pointer whitespace-nowrap"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
      </svg>
      Tambah Soal
    </button>
    <button
      type="button"
      onclick={handleSubmit}
      disabled={isSubmitting}
      class="flex-1 sm:flex-none inline-flex items-center justify-center gap-2 px-6 py-2 bg-blue-600 text-white font-bold rounded-xl shadow-md hover:bg-blue-700 hover:shadow-lg transition-all disabled:opacity-50 disabled:cursor-not-allowed border-none cursor-pointer whitespace-nowrap"
    >
      {isSubmitting ? 'Menyimpan...' : 'Perbarui Kuis'}
    </button>
  </div>
</div>

  {#if isLoading}
    <div class="fixed inset-0 z-[100] flex flex-col items-center justify-center bg-slate-50/50 backdrop-blur-sm {showLoadingSpinner ? 'opacity-100' : 'opacity-0'} transition-opacity duration-300">
      <div class="w-12 h-12 border-4 border-slate-200 border-t-indigo-600 rounded-full animate-spin shadow-sm"></div>
    </div>
  {:else}
  <form onsubmit={handleSubmit} class="space-y-8 pb-12">
    <!-- Detail Kuis -->
    <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-200 flex flex-col gap-5">
      <h2 class="text-xl font-bold text-slate-800 border-b border-slate-100 pb-3">Informasi Utama</h2>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
        <div class="flex flex-col gap-2">
          <label for="title" class="text-sm font-semibold text-slate-800">Judul Kuis</label>
          <input
            type="text"
            id="title"
            bind:value={title}
            required
            class="px-4 py-2.5 rounded-xl border border-slate-200 bg-slate-50 focus:bg-white focus:ring-2 focus:ring-slate-400 focus:border-slate-400 outline-none transition-all"
          />
        </div>
        
        <div class="flex flex-col gap-2">
          <label for="category" class="text-sm font-semibold text-slate-800">Kategori</label>
          <input
            type="text"
            id="category"
            bind:value={category}
            required
            class="px-4 py-2.5 rounded-xl border border-slate-200 bg-slate-50 focus:bg-white focus:ring-2 focus:ring-slate-400 focus:border-slate-400 outline-none transition-all"
            placeholder="Misal: Matematika Dasar"
          />
        </div>

        <div class="flex flex-col gap-2 md:col-span-2">
          <label for="timeLimit" class="text-sm font-semibold text-slate-800">Batas Waktu per Soal (Detik)</label>
          <input
            type="number"
            id="timeLimit"
            bind:value={timeLimit}
            min="5"
            required
            class="px-4 py-2.5 rounded-xl border border-slate-200 bg-slate-50 focus:bg-white focus:ring-2 focus:ring-slate-400 focus:border-slate-400 outline-none transition-all w-full md:w-1/3"
          />
        </div>

        <div class="flex items-center justify-between p-4 bg-slate-50 border border-slate-200 rounded-xl md:col-span-2">
          <div>
            <h3 class="text-sm font-bold text-slate-900">Publikasikan Kuis</h3>
            <p class="text-xs text-slate-500 mt-0.5">Jika dimatikan, kuis akan disimpan sebagai Draf dan tidak bisa diakses murid.</p>
          </div>
          <button 
            type="button"
            role="switch"
            aria-checked={isPublished}
            onclick={() => isPublished = !isPublished}
            class="relative inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-blue-600 focus:ring-offset-2 {isPublished ? 'bg-blue-600' : 'bg-slate-300'}"
          >
            <span class="pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out {isPublished ? 'translate-x-5' : 'translate-x-0'}"></span>
          </button>
        </div>

        <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 p-4 bg-red-50/50 border border-red-200 rounded-xl md:col-span-2">
          <div>
            <h3 class="text-sm font-bold text-red-900">Pertandingkan Ulang (Reset Riwayat)</h3>
            <p class="text-xs text-red-700 mt-0.5">Hapus seluruh data nilai murid di kuis ini secara permanen agar mereka bisa mendapat poin lagi.</p>
          </div>
          <button
            type="button"
            onclick={() => showResetModal = true}
            disabled={isResetThisMonth}
            class="px-4 py-2 bg-white text-red-600 font-bold rounded-lg shadow-sm hover:bg-red-50 transition-colors border border-red-200 cursor-pointer text-sm whitespace-nowrap flex items-center gap-2 shrink-0 disabled:opacity-75 disabled:cursor-not-allowed disabled:bg-slate-100 disabled:text-slate-500 disabled:border-slate-300"
          >
            {#if isResetThisMonth}
              <svg class="w-4 h-4 text-emerald-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
              Sudah Direset Bulan Ini
            {:else}
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path></svg>
              Reset Nilai
            {/if}
          </button>
        </div>
      </div>
    </div>

    <!-- Daftar Soal -->
    <div class="space-y-4">
      <div class="flex justify-between items-center">
        <h2 class="text-xl font-bold text-slate-800">Daftar Pertanyaan</h2>
      </div>

      {#each questions as q, index}
        <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-200 relative group transition-all hover:border-slate-400">
            <button
              type="button"
              onclick={() => removeQuestion(index)}
              class="absolute top-4 right-4 p-2 text-red-500 bg-red-50 hover:bg-red-100 rounded-lg opacity-0 group-hover:opacity-100 transition-opacity cursor-pointer border-none"
              title="Hapus Soal"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </button>

          <div class="flex items-center gap-3 mb-4">
            <div class="w-8 h-8 rounded-full bg-slate-100 text-slate-600 font-bold flex items-center justify-center text-sm">
              {index + 1}
            </div>
            <h3 class="font-semibold text-slate-900 m-0">Pertanyaan</h3>
          </div>

          <div class="flex flex-col gap-4 pl-11">
            <input
              type="text"
              bind:value={q.question}
              required
              placeholder="Tulis pertanyaan di sini..."
              class="px-4 py-2.5 rounded-xl border border-slate-200 bg-slate-50 focus:bg-white focus:ring-2 focus:ring-slate-400 focus:border-slate-400 outline-none w-full"
            />

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mt-2">
              {#each [0, 1, 2, 3] as optIndex}
                <div class="flex items-center gap-3">
                  <input
                    type="radio"
                    name={`answer-${index}`}
                    value={optIndex}
                    bind:group={q.answer}
                    class="w-5 h-5 accent-blue-600 bg-transparent cursor-pointer"
                    required
                  />
                  <input
                    type="text"
                    bind:value={q.options[optIndex]}
                    required
                    placeholder={`Pilihan ${String.fromCharCode(65 + optIndex)}`}
                    class="flex-1 px-4 py-2 rounded-lg border border-slate-200 bg-white focus:ring-2 focus:ring-slate-400 outline-none text-sm"
                  />
                </div>
              {/each}
            </div>
          </div>
        </div>
      {/each}
    </div>

    <!-- Tombol aksi telah dipindah ke sticky header di atas -->
  </form>
{/if}

<!-- Reset Scores Modal -->
{#if showResetModal}
  <div class="fixed inset-0 z-[100] flex items-center justify-center p-4 sm:p-0">
    <div class="absolute inset-0 bg-slate-900/40 backdrop-blur-sm" onclick={() => showResetModal = false}></div>
    <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-sm p-6 overflow-hidden animate-in zoom-in-95 duration-200">
      <div class="w-12 h-12 rounded-full bg-red-100 flex items-center justify-center mb-4 text-red-600">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path></svg>
      </div>
      <h3 class="text-lg font-bold text-slate-900 mb-2">Pertandingkan Ulang?</h3>
      <p class="text-sm text-slate-500 mb-6 leading-relaxed">
        Apakah Anda yakin ingin mereset riwayat kuis ini? Semua data nilai dan status klaim poin murid di kuis ini akan <strong>dihapus permanen</strong>, sehingga mereka bisa mengerjakannya dan mendapat poin lagi dari awal.
      </p>
      <div class="flex flex-col sm:flex-row gap-3">
        <button
          onclick={() => showResetModal = false}
          class="flex-1 px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 font-bold text-sm rounded-xl transition-colors text-center cursor-pointer"
        >
          Batal
        </button>
        <button
          onclick={resetScores}
          disabled={isResetting}
          class="flex-1 px-4 py-2 bg-red-600 hover:bg-red-700 text-white font-bold text-sm rounded-xl transition-colors text-center flex items-center justify-center gap-2 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {#if isResetting}
            <div class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
            Mereset...
          {:else}
            Ya, Reset
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}
