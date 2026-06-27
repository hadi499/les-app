<script lang="ts">
  import type { Card } from "$lib/types";
  import { cardsStore } from "$lib/stores/cards.svelte";
  import { cardFolders, fetchCardFolders, createCardFolder } from "$lib/stores/cardFolders";
  import { onMount } from "svelte";
  import RichEditor from "./RichEditor.svelte";

  interface RichEditorRef {
    getHTML(): string;
    setHTML(html: string): void;
  }

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
  let title = $state("");
  let size = $state("6");
  let uploading = $state(false);
  let uploadMsg = $state("");

  let editorRef = $state<RichEditorRef>();

  onMount(() => {
    fetchCardFolders();
  });

  $effect(() => {
    card_folder_id = edit?.card_folder_id ?? initialFolderId ?? null;
    image = edit?.image ?? "";
    title = edit?.title ?? "";
    size = edit?.size ?? "6";
    if (edit?.content) {
      editorRef?.setHTML(edit.content);
    }
  });

  function handleSubmit(e: Event) {
    e.preventDefault();
    const html = editorRef?.getHTML() ?? "";
    if (!title.trim() || !(html.trim() || html.includes("<img"))) return;
    if (!card_folder_id) {
        alert("Pilih folder terlebih dahulu");
        return;
    }
    onsave({
      card_folder_id,
      image: image.trim(),
      title: title.trim(),
      content: html,
      size,
      cardType: "regular",
    });
    if (!edit) {
      card_folder_id = null;
      image = "";
      title = "";
      size = "6";
      editorRef?.setHTML("");
    }
  }

  async function handleCreateFolder() {
    if (!newFolderName.trim()) return;
    await createCardFolder(newFolderName.trim());
    newFolderName = "";
    isCreatingFolder = false;
  }

  async function handleFileUpload(e: Event) {
    const target = e.target as HTMLInputElement;
    const file = target.files?.[0];
    if (!file) return;
    
    if (!file.type.startsWith("image/")) {
      uploadMsg = "File harus berupa gambar";
      target.value = "";
      return;
    }
    if (file.size > 1024 * 1024) {
      uploadMsg = "Ukuran file maksimal 1MB";
      target.value = "";
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
      target.value = "";
    }
  }
</script>

<form onsubmit={handleSubmit} class="flex flex-col gap-3">
  <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
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

    <div class="flex flex-col gap-1">
      <span class="text-sm font-medium text-slate-700">Gambar</span>
      <div class="flex gap-2 items-start">
        <input
          type="text"
          bind:value={image}
          placeholder="URL gambar atau upload..."
          class="flex-1 px-3 py-2 border border-slate-300 rounded-lg text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none"
        />
        {#if image}
          <button
            type="button"
            onclick={() => {
              image = "";
              uploadMsg = "";
            }}
            class="px-3 py-2 text-sm rounded-lg bg-red-50 border border-red-500 text-red-600 hover:bg-red-100 cursor-pointer whitespace-nowrap"
          >
            Batal
          </button>
        {:else}
          <label
            class="px-3 py-2 text-sm rounded-lg border border-blue-500 text-blue-600 hover:bg-blue-50 cursor-pointer whitespace-nowrap {uploading
              ? 'opacity-50 pointer-events-none'
              : ''}"
          >
            {uploading ? "Uploading..." : "Upload"}
            <input
              type="file"
              accept="image/*"
              onchange={handleFileUpload}
              class="hidden"
              disabled={uploading}
            />
          </label>
        {/if}
      </div>
      {#if uploadMsg}
        <span
          class="text-xs {uploadMsg === 'Upload berhasil'
            ? 'text-green-600'
            : 'text-red-500'}">{uploadMsg}</span
        >
      {/if}
    </div>
  </div>

  <label class="flex flex-col gap-1 text-sm font-medium text-slate-700">
    Ukuran Kartu
    <div class="flex border border-slate-300 rounded-lg overflow-hidden w-fit">
      <button
        type="button"
        onclick={() => (size = "6")}
        class="px-4 py-2 text-sm cursor-pointer {size === '6'
          ? 'bg-blue-600 text-white'
          : 'bg-white text-slate-700 hover:bg-slate-50'}">Kecil (6/A4)</button
      >
      <button
        type="button"
        onclick={() => (size = "4")}
        class="px-4 py-2 text-sm cursor-pointer border-x border-slate-300 {size ===
        '4'
          ? 'bg-blue-600 text-white'
          : 'bg-white text-slate-700 hover:bg-slate-50'}">Sedang (4/A4)</button
      >
      <button
        type="button"
        onclick={() => (size = "2")}
        class="px-4 py-2 text-sm cursor-pointer {size === '2'
          ? 'bg-blue-600 text-white'
          : 'bg-white text-slate-700 hover:bg-slate-50'}">Besar (2/A4)</button
      >
    </div>
  </label>

  <label class="flex flex-col gap-1 text-sm font-medium text-slate-700">
    Judul
    <input
      type="text"
      bind:value={title}
      placeholder="Judul materi pelajaran"
      required
      class="px-3 py-2 border border-slate-300 rounded-lg text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none"
    />
  </label>

  <div class="flex flex-col gap-1">
    <label class="text-sm font-medium text-slate-700">Konten / Ringkasan</label>
    <div class="font-normal">
      <RichEditor bind:this={editorRef} value={edit?.content ?? ""} />
    </div>
  </div>

  <div class="flex gap-2 justify-end">
    <button
      type="button"
      onclick={oncancel}
      class="px-4 py-2 text-sm rounded-lg border border-blue-500 text-blue-600 hover:bg-blue-50 cursor-pointer"
    >
      Batal
    </button>
    <button
      type="submit"
      class="px-4 py-2 text-sm rounded-lg bg-blue-600 text-white hover:bg-blue-700 font-medium cursor-pointer"
    >
      {edit ? "Simpan" : "Tambah Kartu"}
    </button>
  </div>
</form>
