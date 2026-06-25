<script lang="ts">
  import Modal from "$lib/components/Modal.svelte";
  import { onMount, tick } from "svelte";
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import RichEditor from "$lib/components/RichEditor.svelte";
  import { renderMathContent } from "$lib/extensions/math";
  import { toast } from "$lib/stores/toast.svelte";

  interface RichEditorRef {
    getHTML(): string;
    setHTML(html: string): void;
  }

  type Folder = {
    id: number;
    name: string;
  };

  type Note = {
    id: number;
    title: string;
    folder_id: number | null;
    folder: Folder | null;
    content: string;
    created_at: string;
  };

  let folders = $state<Folder[]>([]);
  let notes = $state<Note[]>([]);
  let searchQuery = $state("");
  let currentFolder = $state<Folder | null>(null);
  let viewingNote = $state<Note | null>(null);

  onMount(async () => {
    try {
      // Fetch Folders
      const resFolders = await fetch("/api/folders", {
        credentials: "include",
      });
      if (resFolders.status === 401 || resFolders.status === 403) {
        goto("/dashboard");
        return;
      }
      if (resFolders.ok) {
        folders = await resFolders.json();
      }

      // Fetch Notes
      const resNotes = await fetch("/api/notes", { credentials: "include" });
      if (resNotes.ok) {
        notes = await resNotes.json();
      }
    } catch (e) {
      console.error(e);
    }
  });

  $effect(() => {
    const folderId = page.url.searchParams.get("folder");
    if (folderId && folders.length > 0) {
      const foundFolder = folders.find((f) => f.id === Number(folderId));
      if (foundFolder && currentFolder?.id !== foundFolder.id)
        currentFolder = foundFolder;
    } else if (folders.length > 0 && !folderId) {
      if (currentFolder) currentFolder = null;
    }

    const noteId = page.url.searchParams.get("note");
    if (noteId && notes.length > 0) {
      const foundNote = notes.find((n) => n.id === Number(noteId));
      if (foundNote && viewingNote?.id !== foundNote.id)
        viewingNote = foundNote;
    } else if (notes.length > 0 && !noteId) {
      if (viewingNote) viewingNote = null;
    }
  });

  function openFolder(folder: Folder) {
    currentFolder = folder;
    const url = new URL(page.url);
    url.searchParams.set("folder", String(folder.id));
    goto(url, { keepFocus: true, noScroll: true });
  }

  function closeFolder() {
    currentFolder = null;
    const url = new URL(page.url);
    url.searchParams.delete("folder");
    goto(url, { keepFocus: true, noScroll: true });
  }

  function openNote(note: Note) {
    viewingNote = note;
    const url = new URL(page.url);
    url.searchParams.set("note", String(note.id));
    goto(url, { keepFocus: true, noScroll: true });
  }

  function closeNote() {
    viewingNote = null;
    const url = new URL(page.url);
    url.searchParams.delete("note");
    goto(url, { keepFocus: true, noScroll: true });
  }

  let filteredNotes = $derived(
    notes.filter((n) => {
      if (searchQuery.trim() !== "") {
        return (
          n.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
          n.content.toLowerCase().includes(searchQuery.toLowerCase())
        );
      }
      return currentFolder
        ? n.folder_id === currentFolder.id
        : n.folder_id === null;
    }),
  );

  // Modals state
  let showNoteForm = $state(false);
  let showFolderForm = $state(false);
  let showDeleteModal = $state(false);
  let noteToDelete = $state<Note | null>(null);
  let editingNote = $state<Note | null>(null);
  let editingFolder = $state<Folder | null>(null);
  let showMoveModal = $state(false);
  let noteToMove = $state<Note | null>(null);
  let moveTargetFolderId = $state<number | "">("");

  // Forms state
  let formTitle = $state("");
  let formFolderId = $state<number | "">("");
  let formContent = $state("");
  let newFolderName = $state("");
  let isFolderDropdownOpen = $state(false);
  let editorRef = $state<RichEditorRef>();

  function openNewFolder() {
    editingFolder = null;
    newFolderName = "";
    showFolderForm = true;
  }

  function openEditFolder(folder: Folder) {
    editingFolder = folder;
    newFolderName = folder.name;
    showFolderForm = true;
  }

  async function saveFolder(e: Event) {
    e.preventDefault();
    const name = newFolderName.trim();
    if (!name) return;

    try {
      if (editingFolder) {
        const res = await fetch(`/api/folders/${editingFolder.id}`, {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
          body: JSON.stringify({ name }),
        });
        if (res.ok) {
          const updatedFolder = await res.json();
          folders = folders.map((f) => (f.id === updatedFolder.id ? updatedFolder : f));
          if (currentFolder?.id === updatedFolder.id) currentFolder = updatedFolder;
        }
      } else {
        const res = await fetch("/api/folders", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
          body: JSON.stringify({ name }),
        });
        if (res.ok) {
          const newFolder = await res.json();
          folders = [...folders, newFolder];
          currentFolder = newFolder; // Auto-enter the newly created folder
        }
      }
    } catch (e) {
      console.error(e);
    }
    showFolderForm = false;
  }

  async function deleteFolder(folder: Folder) {
    if (confirm(`Hapus folder "${folder.name}" dan semua isinya?`)) {
      try {
        const res = await fetch(`/api/folders/${folder.id}`, {
          method: "DELETE",
          credentials: "include",
        });
        if (res.ok) {
          folders = folders.filter((f) => f.id !== folder.id);
          notes = notes.filter((n) => n.folder_id !== folder.id);
          if (currentFolder?.id === folder.id) currentFolder = null;
        }
      } catch (e) {
        console.error(e);
      }
    }
  }

  async function openEdit(note?: Note) {
    if (note) {
      editingNote = note;
      formTitle = note.title;
      formFolderId = note.folder_id === null ? "" : note.folder_id;
      formContent = note.content;
    } else {
      editingNote = null;
      formTitle = "";
      formFolderId =
        currentFolder && currentFolder.id !== 0 ? currentFolder.id : "";
      formContent = "";
    }
    showNoteForm = true;
    await tick();
    if (editorRef) {
      editorRef.setHTML(note ? note.content : "");
    }
  }

  async function saveNote(e: Event) {
    e.preventDefault();
    const html = editorRef?.getHTML() ?? "";
    if (!formTitle.trim() || !(html.trim() || html.includes("<img"))) return;

    const payload = {
      title: formTitle,
      folder_id: formFolderId === "" ? null : Number(formFolderId),
      content: html,
    };

    try {
      if (editingNote) {
        const res = await fetch(`/api/notes/${editingNote.id}`, {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
          body: JSON.stringify(payload),
        });
        if (res.ok) {
          const updatedNote = await res.json();
          notes = notes.map((n) => (n.id === updatedNote.id ? updatedNote : n));
          if (viewingNote?.id === updatedNote.id) viewingNote = updatedNote;
        }
      } else {
        const res = await fetch("/api/notes", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
          body: JSON.stringify(payload),
        });
        if (res.ok) {
          const newNote = await res.json();
          notes = [newNote, ...notes];
        }
      }
    } catch (e) {
      console.error(e);
    }
    showNoteForm = false;
  }

  function openMoveModal(note: Note) {
    noteToMove = note;
    moveTargetFolderId = note.folder_id === null ? "" : note.folder_id;
    showMoveModal = true;
  }

  async function moveNoteToFolder() {
    if (!noteToMove) return;
    const payload = {
      title: noteToMove.title,
      folder_id: moveTargetFolderId === "" ? null : Number(moveTargetFolderId),
      content: noteToMove.content,
    };
    try {
      const res = await fetch(`/api/notes/${noteToMove.id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify(payload),
      });
      if (res.ok) {
        const updatedNote = await res.json();
        notes = notes.map((n) => (n.id === updatedNote.id ? updatedNote : n));
        if (viewingNote?.id === updatedNote.id) viewingNote = updatedNote;
        toast.success("Catatan dipindahkan");
      }
    } catch (e) {
      console.error(e);
      toast.error("Gagal memindahkan catatan");
    }
    showMoveModal = false;
    noteToMove = null;
  }

  async function copyNote(note: Note) {
    const payload = {
      title: `${note.title} (Copy)`,
      folder_id: note.folder_id,
      content: note.content,
    };
    try {
      const res = await fetch("/api/notes", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify(payload),
      });
      if (res.ok) {
        const newNote = await res.json();
        notes = [newNote, ...notes];
        toast.success("Catatan berhasil diduplikat");
      }
    } catch (e) {
      console.error(e);
      toast.error("Gagal menduplikat catatan");
    }
  }

  function deleteNote(note: Note) {
    noteToDelete = note;
    showDeleteModal = true;
  }

  async function confirmDelete() {
    if (!noteToDelete) return;
    const id = noteToDelete.id;
    showDeleteModal = false;
    try {
      const res = await fetch(`/api/notes/${id}`, {
        method: "DELETE",
        credentials: "include",
      });
      if (res.ok) {
        notes = notes.filter((n) => n.id !== id);
        if (viewingNote?.id === id) closeNote();
        if (editingNote?.id === id) showNoteForm = false;
        toast.success("Catatan berhasil dihapus");
      } else {
        toast.error("Gagal menghapus catatan");
      }
    } catch (e) {
      console.error(e);
      toast.error("Terjadi kesalahan saat menghapus");
    } finally {
      noteToDelete = null;
    }
  }

  function printNote() {
    window.print();
  }

  let printFontSize = $state(16);
</script>

<svelte:head>
  <title>Catatan - Les Balongarut</title>
  <style>
    @media print {
      @page {
        size: A4 portrait;
        margin: 20mm;
      }
      body {
        background: white !important;
      }
    }
  </style>
</svelte:head>

<!-- VIEW: DETAIL CATATAN -->
{#if viewingNote}
  <div
    class="max-w-[210mm] mx-auto bg-white p-6 md:p-12 rounded-3xl shadow-sm border border-slate-200 print:border-none print:shadow-none print:p-0 print:m-0 print:min-h-0 animate-in fade-in zoom-in-95 duration-200 relative min-h-[297mm]"
  >
    <div
      class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 print:hidden mb-8 border-b border-slate-100 pb-5"
    >
      <button
        onclick={() => closeNote()}
        class="flex items-center gap-2 text-slate-500 hover:text-slate-900 font-medium transition-colors cursor-pointer"
      >
        <svg
          class="w-5 h-5"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          ><path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M10 19l-7-7m0 0l7-7m-7 7h18"
          ></path></svg
        >
        Kembali
      </button>
      <div class="flex items-center gap-2 w-full sm:w-auto">

        <button
          onclick={() => deleteNote(viewingNote!)}
          class="p-2.5 text-red-600 bg-red-50 hover:bg-red-100 border border-red-200 rounded-xl transition-colors cursor-pointer flex items-center justify-center shrink-0"
          title="Hapus"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
        </button>
        <button
          onclick={() => openEdit(viewingNote!)}
          class="p-2.5 text-blue-600 bg-blue-50 hover:bg-blue-100 border border-blue-200 rounded-xl transition-colors cursor-pointer flex items-center justify-center shrink-0"
          title="Edit Catatan"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path></svg>
        </button>
        <div class="flex items-center bg-slate-100 rounded-xl overflow-hidden print:hidden border border-slate-200 shrink-0">
          <button
            onclick={() => printFontSize = Math.max(10, printFontSize - 1)}
            class="px-3 py-2 text-sm font-bold text-slate-600 hover:bg-slate-200 hover:text-slate-900 transition-colors cursor-pointer"
            title="Perkecil Font"
          >A-</button>
          <div class="px-2 py-2 text-xs font-medium text-slate-600 bg-white border-x border-slate-200 text-center min-w-12">
            {printFontSize}px
          </div>
          <button
            onclick={() => printFontSize = Math.min(36, printFontSize + 1)}
            class="px-3 py-2 text-base font-bold text-slate-600 hover:bg-slate-200 hover:text-slate-900 transition-colors cursor-pointer"
            title="Perbesar Font"
          >A+</button>
        </div>
        <button
          onclick={printNote}
          class="p-2.5 text-slate-700 bg-slate-100 hover:bg-slate-200 border border-slate-200 rounded-xl transition-colors cursor-pointer flex items-center justify-center shrink-0"
          title="Cetak A4"
        >
          <svg
            class="w-5 h-5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            ><path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z"
            ></path></svg
          >
        </button>
      </div>
    </div>

    <!-- Print Area -->
    <div class="print:block">
      <h1
        class="text-3xl sm:text-4xl font-semibold text-slate-900 mb-4 tracking-tight leading-tight print:hidden"
      >
        {viewingNote.title}
      </h1>
      <div
        class="flex items-center gap-3 text-sm text-slate-500 mb-8 pb-8 border-b border-slate-100 print:hidden"
      >
        <span
          class="px-3 py-1 rounded-lg bg-slate-100 text-slate-700 font-medium flex items-center gap-1.5"
        >
          <svg
            class="w-4 h-4 text-slate-400"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            ><path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"
            ></path></svg
          >
          {viewingNote.folder ? viewingNote.folder.name : "Tanpa Folder"}
        </span>
        <span class="flex items-center gap-1.5">
          <svg
            class="w-4 h-4 text-slate-400"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            ><path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
            ></path></svg
          >
          {new Date(viewingNote.created_at).toLocaleDateString("id-ID", {
            year: "numeric",
            month: "short",
            day: "numeric",
          })}
        </span>
      </div>
      <div
        class="prose prose-slate max-w-none text-slate-800 leading-loose whitespace-pre-wrap"
        style="font-size: {printFontSize}px;"
      >
        {@html renderMathContent(viewingNote.content)}
      </div>
    </div>
  </div>

  <!-- VIEW: FOLDERS OR NOTES LIST -->
{:else}
  <div class="space-y-6 animate-in fade-in duration-300 print:hidden">
    <!-- Header & Actions -->
    <div
      class="flex flex-col sm:flex-row sm:items-center justify-between gap-4"
    >
      <div>
        {#if currentFolder && !searchQuery}
          <div class="flex items-center gap-2 mb-1">
            <button
              onclick={() => closeFolder()}
              class="text-slate-400 hover:text-slate-700 transition-colors cursor-pointer"
            >
              Catatan Saya
            </button>
            <span class="text-slate-300">/</span>
            <span class="text-slate-700 font-medium">{currentFolder.name}</span>
          </div>
          <h1
            class="text-2xl font-bold text-slate-900 tracking-tight flex items-center gap-2"
          >
            <svg
              class="w-6 h-6 text-blue-500"
              fill="currentColor"
              viewBox="0 0 20 20"
              ><path
                d="M2 6a2 2 0 012-2h5l2 2h5a2 2 0 012 2v6a2 2 0 01-2 2H4a2 2 0 01-2-2V6z"
              ></path></svg
            >
            {currentFolder.name}
          </h1>
        {:else if searchQuery}
          <h1 class="text-2xl font-bold text-slate-900 tracking-tight">
            Hasil Pencarian
          </h1>
          <p class="text-sm text-slate-500 mt-1">Mencari: "{searchQuery}"</p>
        {:else}
          <h1 class="text-2xl font-bold text-slate-900 tracking-tight">
            Catatan Saya
          </h1>
          <p class="text-sm text-slate-500 mt-1">
            Kelola catatan dalam folder.
          </p>
        {/if}
      </div>

      <div class="flex items-center gap-3">
        {#if !currentFolder && !searchQuery}
          <button
            onclick={openNewFolder}
            class="inline-flex items-center gap-1.5 px-3 py-2 text-sm font-medium text-slate-500 hover:text-slate-900 bg-transparent border-none transition-colors cursor-pointer"
          >
            <svg
              class="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              ><path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M9 13h6m-3-3v6m-9 1V7a2 2 0 012-2h6l2 2h6a2 2 0 012 2v8a2 2 0 01-2 2H5a2 2 0 01-2-2z"
              ></path></svg
            >
            Folder Baru
          </button>
        {/if}
        {#if currentFolder}
          <button
            onclick={() => openEdit()}
            class="inline-flex items-center gap-2 px-5 py-2.5 text-sm font-medium rounded-full bg-blue-600 text-white hover:bg-blue-700 transition-colors shadow-sm cursor-pointer"
          >
            <svg
              class="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              ><path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M12 4v16m8-8H4"
              ></path></svg
            >
            Catatan Baru
          </button>
        {/if}
      </div>
    </div>

    <!-- Search Bar -->
    <div class="relative max-w-xl">
      <div
        class="absolute inset-y-0 left-0 pl-5 flex items-center pointer-events-none"
      >
        <svg
          class="w-4 h-4 text-slate-400"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          ><path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
          ></path></svg
        >
      </div>
      <input
        type="text"
        bind:value={searchQuery}
        placeholder="Cari catatan..."
        class="block w-full pl-14 pr-12 py-2.5 border border-slate-200 rounded-full bg-white/60 focus:bg-white text-sm text-slate-800 placeholder-slate-400 focus:outline-none focus:border-slate-300 focus:ring-4 focus:ring-slate-100 transition-all"
      />
      {#if searchQuery}
        <button
          onclick={() => (searchQuery = "")}
          aria-label="Bersihkan pencarian"
          class="absolute right-4 top-1/2 -translate-y-1/2 p-1 text-slate-400 hover:text-slate-600 rounded-full hover:bg-slate-100 transition-colors cursor-pointer"
        >
          <svg
            class="w-4 h-4"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            ><path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M6 18L18 6M6 6l12 12"
            ></path></svg
          >
        </button>
      {/if}
    </div>

    <!-- Content Grid -->
    {#if !currentFolder && !searchQuery}
      <!-- FOLDERS VIEW -->
      <h2
        class="text-sm font-bold text-slate-400 uppercase tracking-wider mb-2 mt-8"
      >
        Folder Kategori
      </h2>
      <div
        class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4"
      >
        {#each folders as folder}
          <div
            class="group relative bg-white rounded-2xl border border-slate-200 p-4 shadow-sm hover:shadow-md hover:border-blue-300 transition-all cursor-pointer overflow-hidden flex items-center gap-4"
            onclick={() => openFolder(folder)}
            role="button"
            tabindex="0"
            onkeydown={(e) => e.key === "Enter" && openFolder(folder)}
          >
            <div
              class="w-12 h-12 rounded-xl bg-blue-50 text-blue-500 flex items-center justify-center shrink-0 group-hover:bg-blue-500 group-hover:text-white transition-colors"
            >
              <svg class="w-6 h-6" fill="currentColor" viewBox="0 0 20 20"
                ><path
                  d="M2 6a2 2 0 012-2h5l2 2h5a2 2 0 012 2v6a2 2 0 01-2 2H4a2 2 0 01-2-2V6z"
                ></path></svg
              >
            </div>
            <div class="flex-1 min-w-0">
              <h3 class="font-bold text-slate-800 truncate">{folder.name}</h3>
              <p class="text-xs text-slate-500">
                {notes.filter((n) => n.folder_id === folder.id).length} catatan
              </p>
            </div>
            <div class="absolute top-2 right-2 flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-all">
              <button
                onclick={(e) => {
                  e.stopPropagation();
                  openEditFolder(folder);
                }}
                class="p-1.5 text-slate-500 hover:text-blue-600 cursor-pointer rounded-lg hover:bg-blue-50"
                title="Edit folder"
              >
                <svg
                  class="w-4 h-4"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  ><path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                  ></path></svg
                >
              </button>
              <button
                onclick={(e) => {
                  e.stopPropagation();
                  deleteFolder(folder);
                }}
                class="p-1.5 text-slate-500 hover:text-red-600 cursor-pointer rounded-lg hover:bg-red-50"
                title="Hapus folder"
              >
                <svg
                  class="w-4 h-4"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  ><path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                  ></path></svg
                >
              </button>
            </div>
          </div>
        {/each}
      </div>
    {:else}
      <!-- NOTES VIEW (In Folder or Search) -->
      <div class="flex items-center justify-between mb-2 mt-8">
        <h2 class="text-sm font-bold text-slate-400 uppercase tracking-wider">
          Daftar Catatan
        </h2>
        <span
          class="text-xs font-medium text-slate-400 bg-slate-100 px-2 py-1 rounded-md"
          >{filteredNotes.length} item</span
        >
      </div>

      {#if filteredNotes.length === 0}
        <div
          class="text-center py-16 bg-white/60 backdrop-blur-sm border border-slate-200 rounded-3xl shadow-sm"
        >
          <div
            class="w-16 h-16 bg-slate-100 rounded-2xl flex items-center justify-center mx-auto mb-4 text-slate-400"
          >
            <svg
              class="w-8 h-8"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              ><path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
              ></path></svg
            >
          </div>
          <h3 class="text-base font-semibold text-slate-900">Kosong</h3>
          <p class="mt-1 text-sm text-slate-500">Belum ada catatan di sini.</p>
        </div>
      {:else}
        <div
          class="bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden divide-y divide-slate-100"
        >
          {#each filteredNotes as note (note.id)}
            <div
              class="group relative flex items-center justify-between p-4 hover:bg-slate-50 transition-colors cursor-pointer"
              onclick={() => openNote(note)}
              role="button"
              tabindex="0"
              onkeydown={(e) => e.key === "Enter" && openNote(note)}
            >
              <div class="flex items-center gap-3 min-w-0">
                <svg
                  class="w-5 h-5 text-slate-400 shrink-0 group-hover:text-blue-500 transition-colors"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  ><path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                  ></path></svg
                >
                <div class="min-w-0">
                  <h3
                    class="text-[15px] font-semibold text-slate-800 truncate group-hover:text-blue-600 transition-colors"
                  >
                    {note.title}
                  </h3>
                  {#if searchQuery}
                    <p class="text-xs text-slate-500 mt-0.5 truncate">
                      {note.folder ? note.folder.name : "Tanpa Folder"}
                    </p>
                  {/if}
                </div>
              </div>

              <div class="flex items-center gap-4 shrink-0 pl-4">
                <span class="text-xs text-slate-400 font-medium hidden sm:block"
                  >{new Date(note.created_at).toLocaleDateString("id-ID", {
                    year: "numeric",
                    month: "short",
                    day: "numeric",
                  })}</span
                >
                <div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-all">
                  <button onclick={(e) => { e.stopPropagation(); copyNote(note); }} class="p-1.5 text-slate-500 hover:text-emerald-600 rounded-lg hover:bg-emerald-50 cursor-pointer" title="Duplikat">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" viewBox="0 0 24 24"><rect x="9" y="9" width="13" height="13" rx="2" ry="2" /><path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1" /></svg>
                  </button>
                  <button onclick={(e) => { e.stopPropagation(); openMoveModal(note); }} class="p-1.5 text-slate-500 hover:text-indigo-600 rounded-lg hover:bg-indigo-50 cursor-pointer" title="Pindah Folder">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" /></svg>
                  </button>
                  <button
                    onclick={(e) => {
                      e.stopPropagation();
                      deleteNote(note);
                    }}
                    class="p-1.5 text-slate-500 hover:text-red-600 rounded-lg hover:bg-red-50 cursor-pointer"
                    title="Hapus"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>
                  </button>
                </div>
              </div>
            </div>
          {/each}
        </div>
      {/if}
    {/if}
  </div>
{/if}

<!-- MODAL: ADD/EDIT NOTE -->
<Modal
  show={showNoteForm}
  onclose={() => {
    showNoteForm = false;
    isFolderDropdownOpen = false;
  }}
  title={editingNote ? "Edit Catatan" : "Catatan Baru"}
>
  <form onsubmit={saveNote} class="space-y-5 relative">
    <div>
      <label for="title" class="block text-sm font-medium text-slate-700 mb-1.5"
        >Judul Catatan</label
      >
      <input
        type="text"
        id="title"
        bind:value={formTitle}
        required
        class="w-full bg-white border border-slate-300 rounded-xl px-4 py-2.5 text-slate-900 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all"
        placeholder="Masukkan judul..."
      />
    </div>

    <div class={`relative ${isFolderDropdownOpen ? 'z-50' : 'z-10'}`}>
      <label
        for="folder"
        class="block text-sm font-medium text-slate-700 mb-1.5"
        >Folder Penempatan</label
      >

      <!-- Transparent backdrop for click-outside -->
      {#if isFolderDropdownOpen}
        <div
          class="fixed inset-0 z-40"
          onclick={() => (isFolderDropdownOpen = false)}
          role="button"
          tabindex="0"
          onkeydown={(e) =>
            e.key === "Escape" && (isFolderDropdownOpen = false)}
        ></div>
      {/if}

      <!-- Custom Select Trigger -->
      <button
        type="button"
        onclick={() => {
          if (!editingNote && currentFolder !== null) return;
          isFolderDropdownOpen = !isFolderDropdownOpen;
        }}
        class={`w-full text-left bg-white border border-slate-300 rounded-xl px-4 py-2.5 text-slate-900 outline-none transition-all flex items-center justify-between shadow-sm cursor-pointer relative ${!editingNote && currentFolder !== null ? "bg-slate-50 cursor-not-allowed text-slate-500" : "hover:border-blue-400 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"}`}
      >
        <span class="flex items-center gap-2 truncate">
          {#if formFolderId === ""}
            <div
              class="w-6 h-6 rounded-md bg-slate-100 flex items-center justify-center shrink-0"
            >
              <svg
                class="w-4 h-4 text-slate-500"
                fill="currentColor"
                viewBox="0 0 20 20"
                ><path
                  fill-rule="evenodd"
                  d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4zm2 6a1 1 0 011-1h6a1 1 0 110 2H7a1 1 0 01-1-1zm1 3a1 1 0 100 2h6a1 1 0 100-2H7z"
                  clip-rule="evenodd"
                ></path></svg
              >
            </div>
            <span class="font-medium text-slate-700">Tanpa Folder</span>
          {:else}
            {@const selectedFolder = folders.find((f) => f.id === formFolderId)}
            {#if selectedFolder}
              <div
                class="w-6 h-6 rounded-md bg-blue-50 text-blue-500 flex items-center justify-center shrink-0"
              >
                <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20"
                  ><path
                    d="M2 6a2 2 0 012-2h5l2 2h5a2 2 0 012 2v6a2 2 0 01-2 2H4a2 2 0 01-2-2V6z"
                  ></path></svg
                >
              </div>
              <span class="font-medium">{selectedFolder.name}</span>
            {:else if currentFolder && currentFolder.id === formFolderId}
              <div
                class="w-6 h-6 rounded-md bg-blue-50 text-blue-500 flex items-center justify-center shrink-0"
              >
                <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20"
                  ><path
                    d="M2 6a2 2 0 012-2h5l2 2h5a2 2 0 012 2v6a2 2 0 01-2 2H4a2 2 0 01-2-2V6z"
                  ></path></svg
                >
              </div>
              <span class="font-medium">{currentFolder.name}</span>
            {:else}
              <span class="text-slate-500">Pilih Folder...</span>
            {/if}
          {/if}
        </span>
        <svg
          class={`w-5 h-5 transition-transform ${isFolderDropdownOpen ? "rotate-180 text-blue-500" : "text-slate-400"}`}
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          ><path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M19 9l-7 7-7-7"
          ></path></svg
        >
      </button>

      <!-- Dropdown Menu -->
      {#if isFolderDropdownOpen}
        <div
          class="absolute z-50 mt-2 w-full bg-white border border-slate-200 rounded-2xl shadow-xl py-2 overflow-y-auto max-h-56 animate-in fade-in slide-in-from-top-2 duration-200"
        >
          {#each folders as folder}
            <button
              type="button"
              class="w-full text-left px-4 py-3 text-sm text-slate-700 hover:bg-blue-50 hover:text-blue-700 transition-colors flex items-center gap-3 cursor-pointer"
              onclick={() => {
                formFolderId = folder.id;
                isFolderDropdownOpen = false;
              }}
            >
              <div
                class="w-8 h-8 rounded-lg bg-blue-50 text-blue-500 flex items-center justify-center shrink-0"
              >
                <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20"
                  ><path
                    d="M2 6a2 2 0 012-2h5l2 2h5a2 2 0 012 2v6a2 2 0 01-2 2H4a2 2 0 01-2-2V6z"
                  ></path></svg
                >
              </div>
              <span class="font-medium">{folder.name}</span>
            </button>
          {/each}
        </div>
      {/if}
    </div>

    <div class="relative z-20">
      <label
        for="content"
        class="block text-sm font-medium text-slate-700 mb-1.5"
        >Isi Catatan</label
      >
      <RichEditor bind:this={editorRef} value={editingNote?.content ?? ""} />
    </div>

    <div class="flex justify-end gap-3 pt-5 border-t border-slate-100">
      <button
        type="button"
        onclick={() => {
          showNoteForm = false;
          isFolderDropdownOpen = false;
        }}
        class="px-5 py-2.5 text-sm font-medium rounded-xl border border-slate-300 text-slate-700 hover:bg-slate-50 transition-colors cursor-pointer"
      >
        Batal
      </button>
      <button
        type="submit"
        class="px-5 py-2.5 text-sm font-medium rounded-xl bg-blue-600 text-white hover:bg-blue-700 hover:shadow-lg transition-all cursor-pointer"
      >
        {editingNote ? "Simpan Perubahan" : "Simpan Catatan"}
      </button>
    </div>
  </form>
</Modal>

<!-- MODAL: ADD FOLDER -->
<Modal
  show={showFolderForm}
  onclose={() => (showFolderForm = false)}
  title={editingFolder ? "Edit Folder" : "Folder Baru"}
>
  <form onsubmit={saveFolder} class="space-y-5">
    <div>
      <label
        for="folderName"
        class="block text-sm font-medium text-slate-700 mb-1.5"
        >Nama Folder</label
      >
      <input
        type="text"
        id="folderName"
        bind:value={newFolderName}
        required
        class="w-full bg-white border border-slate-300 rounded-xl px-4 py-2.5 text-slate-900 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all"
        placeholder="Contoh: Logika Algoritma"
      />
    </div>

    <div class="flex justify-end gap-3 pt-5 border-t border-slate-100">
      <button
        type="button"
        onclick={() => (showFolderForm = false)}
        class="px-5 py-2.5 text-sm font-medium rounded-xl border border-slate-300 text-slate-700 hover:bg-slate-50 transition-colors cursor-pointer"
      >
        Batal
      </button>
      <button
        type="submit"
        class="px-5 py-2.5 text-sm font-medium rounded-xl bg-blue-600 text-white hover:bg-blue-700 hover:shadow-lg transition-all cursor-pointer"
      >
        {editingFolder ? "Simpan Perubahan" : "Buat Folder"}
      </button>
    </div>
  </form>
</Modal>

<!-- MODAL: DELETE CONFIRMATION -->
<Modal
  show={showDeleteModal}
  onclose={() => (showDeleteModal = false)}
  title="Hapus Catatan"
>
  <div class="space-y-5">
    <p class="text-sm text-slate-600">
      Apakah Anda yakin ingin menghapus catatan <strong>"{noteToDelete?.title}"</strong>? Tindakan ini tidak dapat
      dibatalkan.
    </p>
    <div class="flex justify-end gap-3 pt-5 border-t border-slate-100">
      <button
        type="button"
        onclick={() => (showDeleteModal = false)}
        class="px-5 py-2.5 text-sm font-medium rounded-xl border border-slate-300 text-slate-700 hover:bg-slate-50 transition-colors cursor-pointer"
      >
        Batal
      </button>
      <button
        type="button"
        onclick={confirmDelete}
        class="px-5 py-2.5 text-sm font-medium rounded-xl bg-red-600 text-white hover:bg-red-700 transition-colors shadow-sm cursor-pointer"
      >
        Hapus
      </button>
    </div>
  </div>
</Modal>

<!-- MODAL: MOVE NOTE -->
<Modal
  show={showMoveModal}
  onclose={() => (showMoveModal = false)}
  title="Pindah Folder"
>
  <div class="space-y-5">
    <div>
      <label class="block text-sm font-medium text-slate-700 mb-1.5">Pilih Folder Tujuan</label>
      <div class="relative">
        <select
          bind:value={moveTargetFolderId}
          class="w-full bg-white border border-slate-300 rounded-xl px-4 py-2.5 text-slate-900 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all appearance-none cursor-pointer"
        >
          <option value="">Tanpa Folder</option>
          {#each folders as folder}
            <option value={folder.id}>{folder.name}</option>
          {/each}
        </select>
        <div class="absolute inset-y-0 right-0 flex items-center pr-4 pointer-events-none">
          <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
        </div>
      </div>
    </div>
    <div class="flex justify-end gap-3 pt-5 border-t border-slate-100">
      <button
        type="button"
        onclick={() => (showMoveModal = false)}
        class="px-5 py-2.5 text-sm font-medium rounded-xl border border-slate-300 text-slate-700 hover:bg-slate-50 transition-colors cursor-pointer"
      >
        Batal
      </button>
      <button
        type="button"
        onclick={moveNoteToFolder}
        class="px-5 py-2.5 text-sm font-medium rounded-xl bg-indigo-600 text-white hover:bg-indigo-700 hover:shadow-lg transition-all cursor-pointer"
      >
        Pindahkan
      </button>
    </div>
  </div>
</Modal>
