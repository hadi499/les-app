<script lang="ts">
  import type { Card } from "$lib/types";
  import { cardsStore } from "$lib/stores/cards.svelte";

  let {
    onsave,
    oncancel,
    edit,
  }: {
    onsave: (card: Omit<Card, "id">) => void;
    oncancel?: () => void;
    edit?: Card | null;
  } = $props();

  let image = $state("");
  let title = $state("");
  let uploading = $state(false);
  let uploadMsg = $state("");
  let dragging = $state(false);

  $effect(() => {
    image = edit?.image ?? "";
    title = edit?.title ?? "";
  });

  function handleSubmit(e: Event) {
    e.preventDefault();
    if (!title.trim() || !image.trim()) return;
    onsave({
      category: "",
      image: image.trim(),
      title: title.trim(),
      content: "",
      size: "6",
      cardType: "image",
    });
    if (!edit) {
      image = "";
      title = "";
    }
  }

  function handleDragOver(e: DragEvent) {
    e.preventDefault();
    dragging = true;
  }

  function handleDragLeave() {
    dragging = false;
  }

  function handleDrop(e: DragEvent) {
    e.preventDefault();
    dragging = false;
    const file = e.dataTransfer?.files?.[0];
    if (file) uploadFile(file);
  }

  async function handleFileSelect(e: Event) {
    const target = e.target as HTMLInputElement;
    const file = target.files?.[0];
    if (file) await uploadFile(file);
    target.value = "";
  }

  async function uploadFile(file: File) {
    if (!file.type.startsWith("image/")) {
      uploadMsg = "File harus berupa gambar";
      return;
    }
    uploading = true;
    uploadMsg = "";
    try {
      const url = await cardsStore.uploadImage(file);
      image = url;
      uploadMsg = "Upload berhasil";
    } catch (err: any) {
      uploadMsg = "Upload gagal: " + (err.message || "unknown");
    } finally {
      uploading = false;
    }
  }
</script>

<form onsubmit={handleSubmit} class="flex flex-col gap-4">
  <!-- Upload / Pilih Gambar -->
  {#if image}
    <div class="relative border border-zinc-700 rounded-lg p-4">
      <img
        src={image}
        alt="Preview"
        class="max-h-56 mx-auto rounded object-contain"
      />
      <button
        type="button"
        onclick={() => (image = "")}
        class="absolute top-2 right-2 w-7 h-7 flex items-center justify-center rounded-full bg-red-900/300 text-white text-sm hover:bg-red-600 cursor-pointer"
        >&times;</button
      >
    </div>
  {:else}
    <!-- Drag & Drop Area -->
    <div
      class="relative flex flex-col items-center justify-center gap-3 border-2 border-dashed rounded-lg p-8 text-center transition-colors cursor-pointer {dragging
        ? 'border-indigo-500 bg-indigo-900/30'
        : 'border-zinc-600 hover:border-indigo-400'}"
      ondragover={handleDragOver}
      ondragleave={handleDragLeave}
      ondrop={handleDrop}
      onclick={() => {}}
    >
      {#if uploading}
        <div class="flex flex-col items-center gap-2">
          <svg
            class="w-10 h-10 text-indigo-500 animate-spin"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              class="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              stroke-width="4"
            ></circle>
            <path
              class="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
            ></path>
          </svg>
          <span class="text-sm text-zinc-400">Mengupload...</span>
        </div>
      {:else}
        <svg
          class="w-14 h-14 {dragging ? 'text-indigo-400' : 'text-gray-300'}"
          fill="none"
          stroke="currentColor"
          stroke-width="1.5"
          viewBox="0 0 24 24"
        >
          <rect x="3" y="3" width="18" height="18" rx="2" ry="2" /><circle
            cx="8.5"
            cy="8.5"
            r="1.5"
          /><polyline points="21 15 16 10 5 21" />
        </svg>
        <div>
          <p
            class="text-sm font-medium {dragging
              ? 'text-indigo-600'
              : 'text-zinc-400'}"
          >
            {dragging ? "Lepaskan di sini..." : "Seret & lepas gambar di sini"}
          </p>
          <p class="text-xs text-zinc-500 mt-1">atau</p>
        </div>
        <label
          class="px-4 py-2 text-sm rounded-lg bg-indigo-600 text-white hover:bg-indigo-700 cursor-pointer"
        >
          Pilih File
          <input
            type="file"
            accept="image/*"
            onchange={handleFileSelect}
            class="hidden"
          />
        </label>
      {/if}
    </div>
  {/if}

  {#if uploadMsg}
    <span
      class="text-xs {uploadMsg.includes('gagal')
        ? 'text-red-500'
        : 'text-green-600'}">{uploadMsg}</span
    >
  {/if}

  <!-- Judul -->
  <label class="flex flex-col gap-1 text-sm font-medium text-blue-200">
    Judul
    <input
      type="text"
      bind:value={title}
      placeholder="Judul kartu gambar"
      required
      class="px-3 py-2 border border-zinc-600 rounded-lg text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none"
    />
  </label>

  <div class="flex gap-2 justify-end pt-2">
    <button
      type="button"
      onclick={oncancel}
      class="px-4 py-2 text-sm rounded-lg border border-zinc-600 text-blue-200 hover:bg-zinc-900/50 cursor-pointer"
    >
      Batal
    </button>
    <button
      type="submit"
      class="px-4 py-2 text-sm rounded-lg bg-indigo-600 text-white hover:bg-indigo-700 font-medium cursor-pointer"
    >
      {edit ? "Simpan" : "Tambah Kartu Gambar"}
    </button>
  </div>
</form>
