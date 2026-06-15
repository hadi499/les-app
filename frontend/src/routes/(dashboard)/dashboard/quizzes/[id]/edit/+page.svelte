<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import { onMount } from "svelte";

  let quizId = $derived(page.params.id);

  let title = $state("");
  let category = $state("Matematika");
  let timeLimit = $state(15);
  
  let questions: { question: string; options: string[]; answer: number }[] = $state([]);

  let isLoading = $state(true);
  let isSubmitting = $state(false);

  async function fetchQuiz() {
    try {
      const res = await fetch(`http://localhost:8080/api/quizzes/${quizId}`, {
        credentials: "include",
      });
      if (res.ok) {
        const json = await res.json();
        const data = json.data;
        title = data.title;
        category = data.category;
        timeLimit = data.timeLimit;
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

  async function handleSubmit(e: Event) {
    e.preventDefault();
    isSubmitting = true;

    try {
      const res = await fetch(`http://localhost:8080/api/quizzes/${quizId}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({
          title,
          category,
          timeLimit: Number(timeLimit),
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

<div class="mb-8 flex items-center gap-4">
  <a
    href="/dashboard/quizzes"
    class="p-2 text-orange-800 bg-white/50 hover:bg-white rounded-xl border border-orange-200 shadow-sm transition-all"
  >
    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
    </svg>
  </a>
  <div>
    <h1 class="text-3xl font-bold text-orange-950 mb-1">Edit Kuis</h1>
    <p class="text-orange-800 m-0">Perbarui informasi dan daftar pertanyaan.</p>
  </div>
</div>

{#if isLoading}
  <div class="flex justify-center p-12">
    <div class="w-8 h-8 border-4 border-orange-200 border-t-orange-600 rounded-full animate-spin"></div>
  </div>
{:else}
  <form onsubmit={handleSubmit} class="space-y-8 pb-12">
    <!-- Detail Kuis -->
    <div class="bg-white p-6 rounded-2xl shadow-sm border border-orange-200 flex flex-col gap-5">
      <h2 class="text-xl font-bold text-orange-900 border-b border-orange-100 pb-3">Informasi Utama</h2>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
        <div class="flex flex-col gap-2">
          <label for="title" class="text-sm font-semibold text-orange-900">Judul Kuis</label>
          <input
            type="text"
            id="title"
            bind:value={title}
            required
            class="px-4 py-2.5 rounded-xl border border-orange-200 bg-orange-50 focus:bg-white focus:ring-2 focus:ring-orange-400 focus:border-orange-400 outline-none transition-all"
          />
        </div>
        
        <div class="flex flex-col gap-2">
          <label for="category" class="text-sm font-semibold text-orange-900">Kategori</label>
          <input
            type="text"
            id="category"
            bind:value={category}
            required
            class="px-4 py-2.5 rounded-xl border border-orange-200 bg-orange-50 focus:bg-white focus:ring-2 focus:ring-orange-400 focus:border-orange-400 outline-none transition-all"
            placeholder="Misal: Matematika Dasar"
          />
        </div>

        <div class="flex flex-col gap-2 md:col-span-2">
          <label for="timeLimit" class="text-sm font-semibold text-orange-900">Batas Waktu per Soal (Detik)</label>
          <input
            type="number"
            id="timeLimit"
            bind:value={timeLimit}
            min="5"
            required
            class="px-4 py-2.5 rounded-xl border border-orange-200 bg-orange-50 focus:bg-white focus:ring-2 focus:ring-orange-400 focus:border-orange-400 outline-none transition-all w-full md:w-1/3"
          />
        </div>
      </div>
    </div>

    <!-- Daftar Soal -->
    <div class="space-y-4">
      <div class="flex justify-between items-center">
        <h2 class="text-xl font-bold text-orange-900">Daftar Pertanyaan</h2>
        <button
          type="button"
          onclick={addQuestion}
          class="inline-flex items-center gap-2 px-4 py-2 bg-white text-orange-600 border border-orange-300 rounded-lg font-medium hover:bg-orange-50 transition-colors shadow-sm cursor-pointer"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Tambah Soal
        </button>
      </div>

      {#each questions as q, index}
        <div class="bg-white p-6 rounded-2xl shadow-sm border border-orange-200 relative group transition-all hover:border-orange-400">
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
            <div class="w-8 h-8 rounded-full bg-orange-100 text-orange-800 font-bold flex items-center justify-center text-sm">
              {index + 1}
            </div>
            <h3 class="font-semibold text-orange-950 m-0">Pertanyaan</h3>
          </div>

          <div class="flex flex-col gap-4 pl-11">
            <input
              type="text"
              bind:value={q.question}
              required
              placeholder="Tulis pertanyaan di sini..."
              class="px-4 py-2.5 rounded-xl border border-orange-200 bg-orange-50 focus:bg-white focus:ring-2 focus:ring-orange-400 focus:border-orange-400 outline-none w-full"
            />

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mt-2">
              {#each [0, 1, 2, 3] as optIndex}
                <div class="flex items-center gap-3">
                  <input
                    type="radio"
                    name={`answer-${index}`}
                    value={optIndex}
                    bind:group={q.answer}
                    class="w-5 h-5 accent-orange-600 bg-transparent cursor-pointer"
                    required
                  />
                  <input
                    type="text"
                    bind:value={q.options[optIndex]}
                    required
                    placeholder={`Pilihan ${String.fromCharCode(65 + optIndex)}`}
                    class="flex-1 px-4 py-2 rounded-lg border border-orange-200 bg-white focus:ring-2 focus:ring-orange-400 outline-none text-sm"
                  />
                </div>
              {/each}
            </div>
          </div>
        </div>
      {/each}
    </div>

    <div class="flex justify-end pt-6 border-t border-orange-200">
      <button
        type="submit"
        disabled={isSubmitting}
        class="px-8 py-3 bg-orange-600 text-white font-bold rounded-xl shadow-md hover:bg-orange-700 hover:shadow-lg transition-all disabled:opacity-50 disabled:cursor-not-allowed border-none cursor-pointer"
      >
        {isSubmitting ? 'Menyimpan...' : 'Perbarui Kuis'}
      </button>
    </div>
  </form>
{/if}
