<script lang="ts">
  import type { Card } from "$lib/types";
  import { cardsStore } from "$lib/stores/cards.svelte";
  import { cardFolders, fetchCardFolders, createCardFolder } from "$lib/stores/cardFolders";
  import { onMount } from "svelte";

  let {
    onsave,
    oncancel,
    edit,
    initialFolderId,
  }: {
    onsave: (card: Omit<Card, "id">) => void;
    oncancel?: () => void;
    edit?: Card | null;
    initialFolderId?: number | null;
  } = $props();

  let card_folder_id = $state<number | null>(initialFolderId ?? null);
  let newFolderName = $state("");
  let isCreatingFolder = $state(false);

  let image = $state("");
  let selectedFile = $state<File | null>(null);
  let previewUrl = $state("");
  let title = $state("");
  let uploading = $state(false);
  let uploadMsg = $state("");
  let dragging = $state(false);

  onMount(() => {
    fetchCardFolders();
  });

  $effect(() => {
    card_folder_id = edit?.card_folder_id ?? initialFolderId ?? null;
    image = edit?.image ?? "";
    title = edit?.title ?? "";
  });

  async function handleSubmit(e: Event) {
    e.preventDefault();
    if (!title.trim() || (!image.trim() && !selectedFile)) return;
    if (!card_folder_id) {
        alert("Pilih folder terlebih dahulu");
        return;
    }
    
    let finalImageUrl = image;

    if (selectedFile) {
      uploading = true;
      try {
        finalImageUrl = await cardsStore.uploadImage(selectedFile);
      } catch (err: any) {
        uploadMsg = "Upload gagal: " + (err.message || "unknown");
        uploading = false;
        return;
      }
      uploading = false;
    }

    onsave({
      card_folder_id,
      image: finalImageUrl.trim(),
      title: title.trim(),
      content: "",
      size: "6",
      cardType: "image",
    });

    if (!edit) {
      card_folder_id = null;
      image = "";
      title = "";
      selectedFile = null;
      if (previewUrl) {
        URL.revokeObjectURL(previewUrl);
        previewUrl = "";
      }
    }
  }

  async function handleCreateFolder() {
    if (!newFolderName.trim()) return;
    await createCardFolder(newFolderName.trim());
    newFolderName = "";
    isCreatingFolder = false;
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
    if (file) processFile(file);
  }

  function handleFileSelect(e: Event) {
    const target = e.target as HTMLInputElement;
    const file = target.files?.[0];
    if (file) processFile(file);
    target.value = "";
  }

  function processFile(file: File) {
    if (!file.type.startsWith("image/")) {
      uploadMsg = "File harus berupa gambar";
      return;
    }
    if (file.size > 1024 * 1024) {
      uploadMsg = "Ukuran file maksimal 1MB";
      return;
    }
    
    selectedFile = file;
    if (previewUrl) URL.revokeObjectURL(previewUrl);
    previewUrl = URL.createObjectURL(file);
    uploadMsg = "";
    image = "";
  }
</script>

<form onsubmit={handleSubmit} class="flex flex-col gap-4">
  <!-- Upload / Pilih Gambar -->
  {#if image || previewUrl}
    <div class="relative border border-slate-300 rounded-lg p-4">
      <img
        src={previewUrl || image}
        alt="Preview"
        class="max-h-56 mx-auto rounded object-contain"
      />
      <button
        type="button"
        onclick={() => {
          image = "";
          selectedFile = null;
          if (previewUrl) URL.revokeObjectURL(previewUrl);
          previewUrl = "";
        }}
        class="absolute top-2 right-2 w-7 h-7 flex items-center justify-center rounded-full bg-red-900/300 text-slate-900 text-sm hover:bg-red-600 cursor-pointer"
        >&times;</button
      >
    </div>
  {:else}
    <!-- Drag & Drop Area -->
    <div
      class="relative flex flex-col items-center justify-center gap-3 border-2 border-dashed rounded-lg p-8 text-center transition-colors cursor-pointer {dragging
        ? 'border-indigo-500 bg-white/60'
        : 'border-slate-300 hover:border-indigo-400'}"
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
          <span class="text-sm text-slate-700">Mengupload...</span>
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
              : 'text-slate-700'}"
          >
            {dragging ? "Lepaskan di sini..." : "Seret & lepas gambar di sini"}
          </p>
          <p class="text-xs text-slate-500 mt-1">atau</p>
        </div>
        <label
          class="px-4 py-2 text-sm rounded-lg bg-blue-600 text-white hover:bg-blue-700 cursor-pointer"
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
      class="text-xs {uploadMsg === 'Upload berhasil'
        ? 'text-green-600'
        : 'text-red-500'}">{uploadMsg}</span
    >
  {/if}

  <!-- Kategori -->
  <label class="flex flex-col gap-1 text-sm font-medium text-slate-700">
    Folder Kartu
    {#if isCreatingFolder}
      <div class="flex gap-2">
          <input
              type="text"
              bind:value={newFolderName}
              placeholder="Nama folder baru..."
              class="flex-1 px-3 py-2 border border-slate-300 rounded-lg text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none"
          />
          <button type="button" onclick={handleCreateFolder} class="px-3 py-2 bg-indigo-600 text-white rounded-lg text-sm">Buat</button>
          <button type="button" onclick={() => (isCreatingFolder = false)} class="px-3 py-2 bg-slate-200 text-slate-700 rounded-lg text-sm">Batal</button>
      </div>
    {:else}
      <div class="flex gap-2">
          <select
              bind:value={card_folder_id}
              class="flex-1 px-3 py-2 border border-slate-300 rounded-lg text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none bg-white"
          >
              <option value={null} disabled>-- Pilih Folder --</option>
              {#each $cardFolders as folder}
                  <option value={folder.id}>{folder.name}</option>
              {/each}
          </select>
          <button type="button" onclick={() => (isCreatingFolder = true)} class="px-3 py-2 bg-slate-100 text-slate-600 border border-slate-300 rounded-lg text-sm">+</button>
      </div>
    {/if}
  </label>

  <!-- Judul -->
  <label class="flex flex-col gap-1 text-sm font-medium text-slate-700">
    Judul
    <input
      type="text"
      bind:value={title}
      placeholder="Judul kartu gambar"
      required
      class="px-3 py-2 border border-slate-300 rounded-lg text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none"
    />
  </label>

  <div class="flex gap-2 justify-end pt-2">
    <button
      type="button"
      onclick={() => {
        card_folder_id = null;
        image = "";
        selectedFile = null;
        if (previewUrl) URL.revokeObjectURL(previewUrl);
        previewUrl = "";
        title = "";
        uploadMsg = "";
        if (oncancel) oncancel();
      }}
      class="px-4 py-2 text-sm rounded-lg border border-slate-300 hover:bg-slate-50 text-slate-700 cursor-pointer"
    >
      Batal
    </button>
    <button
      type="submit"
      disabled={uploading}
      class="px-4 py-2 text-sm rounded-lg bg-blue-600 text-white hover:bg-blue-700 font-medium cursor-pointer disabled:opacity-50"
    >
      {uploading ? "Menyimpan..." : (edit ? "Simpan" : "Tambah Kartu Gambar")}
    </button>
  </div>
</form>
