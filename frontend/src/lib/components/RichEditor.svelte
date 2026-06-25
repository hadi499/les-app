<script lang="ts">
  import "katex/dist/katex.min.css";
  import { onMount } from "svelte";
  import { Editor } from "@tiptap/core";
  import StarterKit from "@tiptap/starter-kit";
  import TextAlign from "@tiptap/extension-text-align";

  let { value = "" }: { value?: string } = $props();

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
  }

  onMount(() => {
    if (!editorEl) return;

    editor = new Editor({
      element: editorEl,
      extensions: [
        StarterKit.configure({ heading: { levels: [2] } }),
        TextAlign.configure({ types: ["heading", "paragraph"] }),
      ],
      content: value || "",
      editorProps: {
        attributes: {
          class: "prose prose-sm max-w-none focus:outline-none min-h-[80px]",
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
    editor?.commands.setContent(html || "");
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

  export { getHTML, setHTML };
</script>

<div class="flex flex-col gap-1">
  <div
    class="flex flex-wrap gap-0.5 p-1 border border-gray-300 rounded-t-lg bg-gray-50"
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
      onclick={() => editor?.chain().focus().setTextAlign("left").run()}
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
    class="border border-gray-300 rounded-b-lg focus-within:ring-2 focus-within:ring-indigo-500 focus-within:border-indigo-500 min-h-[100px] px-3 py-2 text-sm cursor-text"
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
    min-height: 80px;
    height: 100%;
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
