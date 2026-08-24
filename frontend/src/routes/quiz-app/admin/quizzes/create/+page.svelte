<script lang="ts">
  import { goto } from "$app/navigation";
  import { compressImageFile } from "$lib/utils";
  import { onMount, onDestroy } from "svelte";

  let title = $state("");
  let category_id = $state<number | null>(null);
  let timeLimit = $state(15);
  let isPublished = $state(false);
  let categories: any[] = $state([]);

  onMount(async () => {
    const res = await fetch("/api/kuisapp/categories", { credentials: "include" });
    if (res.ok) {
      const d = await res.json();
      categories = d.data || [];
    }
  });
  
  let questions: { question: string; image?: string; options: string[]; answer: number }[] = $state([]);
  let isUploadingImage: Record<number, boolean> = $state({});

  let isSubmitting = $state(false);
  let isSuccess = false;

  let newlyUploadedImages: string[] = [];

  onDestroy(() => {
    if (!isSuccess && newlyUploadedImages.length > 0) {
      newlyUploadedImages.forEach((url) => deleteImageFromServer(url));
    }
  });

  async function deleteImageFromServer(url?: string) {
    if (!url) return;
    try {
      await fetch(`/api/upload?url=${encodeURIComponent(url)}`, {
        method: "DELETE",
        credentials: "include",
      });
    } catch (err) {
      console.error("Gagal menghapus gambar:", err);
    }
  }

  function addQuestion() {
    questions = [...questions, { question: "", image: "", options: ["", "", "", ""], answer: 0 }];
  }

  function removeQuestion(index: number) {
    if (questions[index].image) {
      deleteImageFromServer(questions[index].image);
    }
    questions = questions.filter((_, i) => i !== index);
  }

  async function handleImageUpload(e: Event, index: number) {
    const target = e.target as HTMLInputElement;
    if (!target.files || target.files.length === 0) return;

    isUploadingImage[index] = true;
    try {
      let file = target.files[0];
      file = await compressImageFile(file);

      if (file.size > 1 * 1024 * 1024) {
        alert("Ukuran file setelah kompresi masih lebih dari 1MB. Silakan pilih gambar lain.");
        return;
      }

      const formData = new FormData();
      formData.append("image", file);

      const res = await fetch("/api/kuisapp/upload", {
        method: "POST",
        body: formData,
        credentials: "include",
      });
      if (!res.ok) {
        const err = await res.json();
        throw new Error(err.error || "Gagal mengunggah gambar");
      }
      const data = await res.json();
      questions[index].image = data.url;
      newlyUploadedImages.push(data.url);
    } catch (err: any) {
      alert(err.message || "Terjadi kesalahan");
    } finally {
      isUploadingImage[index] = false;
    }
  }

  async function handleSubmit(e?: Event) {
    if (e) e.preventDefault();
    isSubmitting = true;

    try {
      const res = await fetch(`/api/kuisapp/quizzes`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({
          title,
          category_id: category_id ? Number(category_id) : null,
          timeLimit: Number(timeLimit),
          is_published: isPublished,
          questions: questions.map(q => ({
            question: q.question,
            image: q.image,
            options: q.options,
            answer: Number(q.answer)
          }))
        })
      });

      if (res.ok) {
        isSuccess = true;
        goto("/quiz-app/admin");
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
  <title>Tambah Kuis | Les Balongarut</title>
</svelte:head>

<div class="sticky top-16 z-40 bg-slate-50/95 backdrop-blur-md py-3 sm:py-4 mb-6 sm:mb-8 border-b border-slate-200/60 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 -mx-4 px-4 sm:-mx-8 sm:px-8">
  <div class="flex items-center gap-4">
    <a
      href="/quiz-app/admin"
      class="p-2 text-slate-600 bg-white/50 hover:bg-white rounded-xl border border-slate-200 shadow-sm transition-all"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
      </svg>
    </a>
    <div>
      <h1 class="text-2xl sm:text-3xl font-bold text-slate-900 mb-1">Tambah Kuis Baru</h1>
      <p class="text-sm text-slate-600 m-0 hidden sm:block">Buat set kuis dan daftar pertanyaannya.</p>
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
      {isSubmitting ? 'Menyimpan...' : 'Simpan Kuis'}
    </button>
  </div>
</div>

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
          placeholder="Misal: Kuis Matematika Dasar"
        />
      </div>
      
      <div class="flex flex-col gap-2">
        <label for="category" class="text-sm font-semibold text-slate-800">Kategori</label>
        <select
          id="category"
          bind:value={category_id}
          class="px-4 py-2.5 rounded-xl border border-slate-200 bg-slate-50 focus:bg-white focus:ring-2 focus:ring-slate-400 focus:border-slate-400 outline-none transition-all"
        >
          <option value={null}>-- Tanpa Kategori --</option>
          {#each categories as cat}
            <option value={cat.id}>{cat.name}</option>
          {/each}
        </select>
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

        <div class="flex items-center justify-between p-4 bg-slate-50 border border-slate-200 rounded-xl">
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

        <div class="flex flex-col gap-4">
          <div class="flex flex-col gap-2">
            <label class="text-sm font-semibold text-slate-800">Gambar Pendukung (Opsional)</label>
            <div class="flex items-start gap-4">
              {#if q.image}
                <div class="relative w-32 h-32 rounded-lg border border-slate-200 overflow-hidden bg-slate-50 shrink-0">
                  <img src={q.image} alt="Preview" class="w-full h-full object-cover" />
                  <button
                    type="button"
                    onclick={() => {
                      deleteImageFromServer(questions[index].image);
                      questions[index].image = "";
                    }}
                    class="absolute top-1 right-1 p-1 bg-white/90 rounded-md text-red-500 hover:text-red-700 hover:bg-white shadow-sm border border-slate-200 cursor-pointer"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
                  </button>
                </div>
              {/if}
              <div class="flex-1">
                <input
                  type="file"
                  accept="image/*"
                  onchange={(e) => handleImageUpload(e, index)}
                  disabled={isUploadingImage[index]}
                  class="block w-full text-sm text-slate-500 file:mr-4 file:py-2 file:px-4 file:rounded-xl file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100 disabled:opacity-50"
                />
                {#if isUploadingImage[index]}
                  <p class="text-xs text-blue-600 mt-2 animate-pulse">Mengunggah gambar...</p>
                {/if}
              </div>
            </div>
          </div>

          <div class="flex items-center gap-3">
            <div class="w-8 h-8 shrink-0 rounded-full bg-slate-100 text-slate-600 font-bold flex items-center justify-center text-sm">
              {index + 1}
            </div>
            <input
              type="text"
              bind:value={q.question}
              required
              placeholder="Tulis pertanyaan di sini..."
              class="px-4 py-2.5 rounded-xl border border-slate-200 bg-slate-50 focus:bg-white focus:ring-2 focus:ring-slate-400 focus:border-slate-400 outline-none w-full"
            />
          </div>

          <div class="pl-11 w-full">
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
            <p class="text-xs text-blue-600/70 italic m-0 mt-3">* Pilih radio button pada jawaban yang benar</p>
          </div>
        </div>
      </div>
    {/each}
  </div>

    <!-- Tombol aksi telah dipindah ke sticky header di atas -->
</form>
