<script lang="ts">
  import "katex/dist/katex.min.css";
  import { onMount } from "svelte";
  import { Editor, Extension } from "@tiptap/core";
  import StarterKit from "@tiptap/starter-kit";
  import TextAlign from "@tiptap/extension-text-align";
  import ImageResize from "tiptap-extension-resize-image";

  const TabHandler = Extension.create({
    name: "tabHandler",
    addKeyboardShortcuts() {
      return {
        Tab: ({ editor }) => {
          editor.commands.insertContent("\t");
          return true;
        },
      };
    },
  });

  let {
    value = "",
    minHeight = "min-h-[80px]",
    containerMinHeight = "min-h-[100px]",
    textSize = "prose-sm",
  }: {
    value?: string;
    minHeight?: string;
    containerMinHeight?: string;
    textSize?: string;
  } = $props();

  let editorEl = $state<HTMLDivElement>();
  let editor = $state<Editor>();

  // Reactive active states — diperbarui tiap TipTap selection/update
  let active = $state({
    bold: false,
    italic: false,
    underline: false,
    strike: false,
    heading: false,
    bulletList: false,
    orderedList: false,
    alignLeft: false,
    alignCenter: false,
    alignRight: false,
    alignJustify: false,
    isImageSelected: false,
  });

  function updateActive() {
    if (!editor) return;
    active.bold = editor.isActive("bold");
    active.italic = editor.isActive("italic");
    active.underline = editor.isActive("underline");
    active.strike = editor.isActive("strike");
    active.heading = editor.isActive("heading", { level: 2 });
    active.bulletList = editor.isActive("bulletList");
    active.orderedList = editor.isActive("orderedList");
    active.alignLeft = editor.isActive({ textAlign: "left" });
    active.alignCenter = editor.isActive({ textAlign: "center" });
    active.alignRight = editor.isActive({ textAlign: "right" });
    active.alignJustify = editor.isActive({ textAlign: "justify" });
    active.isImageSelected = editor.isActive("image") || editor.isActive("imageResize");
  }

  onMount(() => {
    if (!editorEl) return;

    editor = new Editor({
      element: editorEl,
      extensions: [
        StarterKit.configure({
          heading: { levels: [2] },
          code: false,
          codeBlock: false,
        }),
        TextAlign.configure({ types: ["heading", "paragraph"] }),
        ImageResize,
        TabHandler,
      ],
      content: value || "",
      parseOptions: {
        preserveWhitespace: "full",
      },
      editorProps: {
        attributes: {
          class: `prose ${textSize} max-w-none focus:outline-none ${minHeight}`,
        },
      },
    });

    // Sync active states ke Svelte reactivity
    editor.on("selectionUpdate", updateActive);
    editor.on("transaction", updateActive);
    updateActive();
  });

  function getHTML(): string {
    return editor?.getHTML() ?? "";
  }

  function setHTML(html: string) {
    editor?.commands.setContent(html || "", { parseOptions: { preserveWhitespace: "full" } });
  }

  let showMathModal = $state(false);
  let mathType = $state<"inline" | "display">("inline");
  let mathInput = $state("");

  function openMathModal(type: "inline" | "display") {
    mathType = type;
    mathInput = type === "inline" ? "x^2 + y^2 = r^2" : "\\frac{-b \\pm \\sqrt{b^2-4ac}}{2a}";
    showMathModal = true;
  }

  function submitMath() {
    if (!mathInput || !editor) return;
    if (mathType === "inline") {
      editor.chain().focus().insertContent(`$${mathInput}$`).run();
    } else {
      editor.chain().focus().insertContent(`<p>$$${mathInput}$$</p>`).run();
    }
    showMathModal = false;
  }

  function insertMathInline() {
    openMathModal("inline");
  }

  function insertMathDisplay() {
    openMathModal("display");
  }

  function handleAlignLeft() {
    if (!editor) return;
    const { state } = editor;
    const { tr, selection } = state;

    const deletions: {from: number, to: number}[] = [];

    state.doc.nodesBetween(selection.from, selection.to, (node, pos) => {
      if (node.isBlock) {
        const firstChild = node.firstChild;
        if (firstChild && firstChild.isText && firstChild.text) {
          const match = firstChild.text.match(/^[\s\t]+/);
          if (match) {
            deletions.push({
              from: pos + 1,
              to: pos + 1 + match[0].length
            });
          }
        }
      }
    });

    for (let i = deletions.length - 1; i >= 0; i--) {
      tr.delete(deletions[i].from, deletions[i].to);
    }

    if (deletions.length > 0) {
      editor.view.dispatch(tr);
    }
    
    editor.chain().focus().setTextAlign("left").run();
  }

  let fileInput = $state<HTMLInputElement>();

  async function handleImageUpload(e: Event) {
    const file = (e.target as HTMLInputElement).files?.[0];
    if (!file) return;

    const formData = new FormData();
    formData.append("image", file);

    try {
      const res = await fetch("/api/upload?type=materi", {
        method: "POST",
        body: formData,
        credentials: "include"
      });
      if (res.ok) {
        const data = await res.json();
        if (data.url && editor) {
          editor.chain().focus().setImage({ src: data.url }).run();
        }
      } else {
        console.error("Upload failed");
      }
    } catch (err) {
      console.error(err);
    }
    
    if (fileInput) {
      fileInput.value = "";
    }
  }

  export { getHTML, setHTML };
</script>

<div class="flex flex-col gap-1 relative">
  <div
    class="flex flex-wrap gap-0.5 p-1 border border-gray-300 rounded-t-lg bg-gray-50 sticky top-0 z-10"
  >
    <button
      type="button"
      onclick={() => editor?.chain().focus().toggleBold().run()}
      class="px-1.5 py-1 text-xs rounded cursor-pointer {active.bold
        ? 'bg-indigo-100 text-indigo-700'
        : 'text-gray-600 hover:bg-gray-100'}"
      title="Bold"><b>B</b></button
    >
    <button
      type="button"
      onclick={() => editor?.chain().focus().toggleItalic().run()}
      class="px-1.5 py-1 text-xs rounded cursor-pointer {active.italic
        ? 'bg-indigo-100 text-indigo-700'
        : 'text-gray-600 hover:bg-gray-100'}"
      title="Italic"><i>I</i></button
    >
    <button
      type="button"
      onclick={() => editor?.chain().focus().toggleUnderline().run()}
      class="px-1.5 py-1 text-xs rounded cursor-pointer {active.underline
        ? 'bg-indigo-100 text-indigo-700'
        : 'text-gray-600 hover:bg-gray-100'}"
      title="Underline"><u>U</u></button
    >
    <button
      type="button"
      onclick={() => editor?.chain().focus().toggleStrike().run()}
      class="px-1.5 py-1 text-xs rounded cursor-pointer {active.strike
        ? 'bg-indigo-100 text-indigo-700'
        : 'text-gray-600 hover:bg-gray-100'}"
      title="Strikethrough"><s>S</s></button
    >
    <span class="w-px bg-gray-300 mx-0.5"></span>
    <button
      type="button"
      onclick={() => editor?.chain().focus().toggleHeading({ level: 2 }).run()}
      class="px-1.5 py-1 text-xs rounded cursor-pointer {active.heading
        ? 'bg-indigo-100 text-indigo-700'
        : 'text-gray-600 hover:bg-gray-100'}"
      title="Heading">H</button
    >
    <span class="w-px bg-gray-300 mx-0.5"></span>
    <button
      type="button"
      onclick={() => editor?.chain().focus().toggleBulletList().run()}
      class="px-1.5 py-1 text-xs rounded cursor-pointer {active.bulletList
        ? 'bg-indigo-100 text-indigo-700'
        : 'text-gray-600 hover:bg-gray-100'}"
      title="Bullet List">&#8226;</button
    >
    <button
      type="button"
      onclick={() => editor?.chain().focus().toggleOrderedList().run()}
      class="px-1.5 py-1 text-xs rounded cursor-pointer {active.orderedList
        ? 'bg-indigo-100 text-indigo-700'
        : 'text-gray-600 hover:bg-gray-100'}"
      title="Numbered List">1.</button
    >
    <span class="w-px bg-gray-300 mx-0.5"></span>
    <button
      type="button"
      onclick={handleAlignLeft}
      class="px-1.5 py-1 text-xs rounded cursor-pointer {active.alignLeft
        ? 'bg-indigo-100 text-indigo-700'
        : 'text-gray-600 hover:bg-gray-100'}"
      title="Align Left"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="16"
        height="16"
        viewBox="0 0 24 24"
        stroke-width="2"
        stroke="currentColor"
        fill="none"
        stroke-linecap="round"
        stroke-linejoin="round"
        ><path stroke="none" d="M0 0h24v24H0z" fill="none" /><path
          d="M4 6l16 0"
        /><path d="M4 12l10 0" /><path d="M4 18l14 0" /></svg
      >
    </button>
    <button
      type="button"
      onclick={() => editor?.chain().focus().setTextAlign("center").run()}
      class="px-1.5 py-1 text-xs rounded cursor-pointer {active.alignCenter
        ? 'bg-indigo-100 text-indigo-700'
        : 'text-gray-600 hover:bg-gray-100'}"
      title="Align Center"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="16"
        height="16"
        viewBox="0 0 24 24"
        stroke-width="2"
        stroke="currentColor"
        fill="none"
        stroke-linecap="round"
        stroke-linejoin="round"
        ><path stroke="none" d="M0 0h24v24H0z" fill="none" /><path
          d="M4 6l16 0"
        /><path d="M8 12l8 0" /><path d="M6 18l12 0" /></svg
      >
    </button>
    <button
      type="button"
      onclick={() => editor?.chain().focus().setTextAlign("right").run()}
      class="px-1.5 py-1 text-xs rounded cursor-pointer {active.alignRight
        ? 'bg-indigo-100 text-indigo-700'
        : 'text-gray-600 hover:bg-gray-100'}"
      title="Align Right"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="16"
        height="16"
        viewBox="0 0 24 24"
        stroke-width="2"
        stroke="currentColor"
        fill="none"
        stroke-linecap="round"
        stroke-linejoin="round"
        ><path stroke="none" d="M0 0h24v24H0z" fill="none" /><path
          d="M4 6l16 0"
        /><path d="M10 12l10 0" /><path d="M6 18l14 0" /></svg
      >
    </button>
    <button
      type="button"
      onclick={() => editor?.chain().focus().setTextAlign("justify").run()}
      class="px-1.5 py-1 text-xs rounded cursor-pointer {active.alignJustify
        ? 'bg-indigo-100 text-indigo-700'
        : 'text-gray-600 hover:bg-gray-100'}"
      title="Align Justify"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="16"
        height="16"
        viewBox="0 0 24 24"
        stroke-width="2"
        stroke="currentColor"
        fill="none"
        stroke-linecap="round"
        stroke-linejoin="round"
        ><path stroke="none" d="M0 0h24v24H0z" fill="none" /><path
          d="M4 6l16 0"
        /><path d="M4 12l16 0" /><path d="M4 18l16 0" /></svg
      >
    </button>
    <span class="w-px bg-gray-300 mx-0.5"></span>
    <button
      type="button"
      onclick={() => fileInput?.click()}
      class="px-1.5 py-1 text-xs rounded cursor-pointer text-gray-600 hover:bg-gray-100"
      title="Insert Image"
    >
      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M15 8h.01"/><path d="M3 6a3 3 0 0 1 3 -3h12a3 3 0 0 1 3 3v12a3 3 0 0 1 -3 3h-12a3 3 0 0 1 -3 -3v-12z"/><path d="M3 16l5 -5c.928 -.893 2.072 -.893 3 0l5 5"/><path d="M14 14l1 -1c.928 -.893 2.072 -.893 3 0l3 3"/></svg>
    </button>
    {#if active.isImageSelected}
      <span class="w-px bg-gray-300 mx-0.5"></span>
      <button
        type="button"
        onclick={() => editor?.chain().focus().deleteSelection().run()}
        class="px-1.5 py-1 text-xs rounded cursor-pointer text-red-600 hover:bg-red-100 bg-red-50"
        title="Hapus Gambar"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M4 7l16 0"/><path d="M10 11l0 6"/><path d="M14 11l0 6"/><path d="M5 7l1 12a2 2 0 0 0 2 2h8a2 2 0 0 0 2 -2l1 -12"/><path d="M9 7v-3a1 1 0 0 1 1 -1h4a1 1 0 0 1 1 1v3"/></svg>
      </button>
    {/if}
    <span class="w-px bg-gray-300 mx-0.5"></span>
    <button
      type="button"
      onclick={insertMathInline}
      class="px-1.5 py-1 text-xs rounded cursor-pointer text-gray-600 hover:bg-gray-100"
      title="Inline Math">$</button
    >
    <button
      type="button"
      onclick={insertMathDisplay}
      class="px-1.5 py-1 text-xs rounded cursor-pointer text-gray-600 hover:bg-gray-100"
      title="Display Math">$$</button
    >
  </div>
  <div
    class={`border border-gray-300 rounded-b-lg focus-within:ring-2 focus-within:ring-indigo-500 focus-within:border-indigo-500 ${containerMinHeight} px-3 py-2 text-sm cursor-text`}
    onmousedown={(e) => {
      const selection = window.getSelection();
      // Hanya intercept jika ada teks yang sedang terseleksi
      if (!selection || selection.toString().length === 0) return;

      const isProseMirrorContent = (e.target as HTMLElement).closest(
        ".ProseMirror",
      );
      if (isProseMirrorContent) {
        window.getSelection()?.removeAllRanges();
        editor?.commands.focus("end");
      }
    }}
  >
    <div bind:this={editorEl}></div>
  </div>
</div>

<style>
  :global(.ProseMirror) {
    height: 100%;
    white-space: pre-wrap;
    tab-size: 4;
  }
  :global(.ProseMirror img) {
    max-width: 100%;
    height: auto;
    border-radius: 6px;
  }
  :global(.ProseMirror img.ProseMirror-selectednode) {
    outline: 3px solid #6366f1;
    box-shadow: 0 0 0 4px rgba(99, 102, 241, 0.2);
  }
</style>

{#if showMathModal}
  <div class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-slate-900/50 backdrop-blur-sm animate-in fade-in duration-200">
    <div class="bg-white rounded-2xl shadow-xl w-full max-w-lg overflow-hidden animate-in zoom-in-95 duration-200">
      <div class="p-6">
        <h3 class="text-xl font-bold text-slate-900 mb-4">
          Masukkan Rumus {mathType === "inline" ? "Baris (Inline)" : "Blok (Display)"}
        </h3>
        <textarea
          bind:value={mathInput}
          class="w-full px-4 py-3 border border-slate-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:outline-none min-h-[100px] font-mono text-sm resize-y"
          placeholder="Tuliskan LaTeX di sini..."
        ></textarea>
      </div>
      <div class="px-6 py-4 bg-slate-50 border-t border-slate-100 flex justify-end gap-3">
        <button
          onclick={() => showMathModal = false}
          class="px-4 py-2 text-sm font-semibold text-slate-700 hover:bg-slate-200 rounded-xl transition-colors cursor-pointer"
        >
          Batal
        </button>
        <button
          onclick={submitMath}
          class="px-4 py-2 text-sm font-semibold text-white bg-indigo-600 hover:bg-indigo-700 rounded-xl transition-colors cursor-pointer"
        >
          Sisipkan
        </button>
      </div>
    </div>
  </div>
{/if}

<input type="file" bind:this={fileInput} onchange={handleImageUpload} accept="image/*" class="hidden" />
