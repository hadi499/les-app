<script lang="ts">
  import hljs from "highlight.js";
  import "highlight.js/styles/atom-one-dark.css";

  let codeInput = $state("");

  let highlightedCode = $derived.by(() => {
    if (!codeInput) return "";
    const rawCode = codeInput;
    const html = hljs.highlightAuto(rawCode).value;

    const rawLines = rawCode.split("\n");
    const htmlLines = html.split("\n");

    let openTags: string[] = [];

    const processedLines = htmlLines.map((lineHtml, i) => {
      const leadingWhitespace = rawLines[i].match(/^\s*/)?.[0] || "";
      const spaces = leadingWhitespace.replace(/\t/g, "  ").length;

      let prefix = openTags.join("");

      const tagRegex = /<\/?span[^>]*>/g;
      let match;
      while ((match = tagRegex.exec(lineHtml)) !== null) {
        const tag = match[0];
        if (tag.startsWith("</")) {
          openTags.pop();
        } else {
          openTags.push(tag);
        }
      }

      let suffix = openTags.map(() => "</span>").join("");

      let lineContent = prefix + lineHtml + suffix;
      if (!lineHtml && !lineContent) {
        lineContent = " ";
      } else if (!lineHtml) {
        lineContent = prefix + " " + suffix;
      }

      return `<div class="code-line" style="padding-left: ${spaces}ch; text-indent: -${spaces}ch;">${lineContent}</div>`;
    });

    return processedLines.join("");
  });

  function handleKeydown(event: KeyboardEvent) {
    if (event.key === "Tab") {
      event.preventDefault();
      const target = event.target as HTMLTextAreaElement;
      const start = target.selectionStart;
      const end = target.selectionEnd;

      // Menggunakan 2 spasi untuk tab
      codeInput =
        codeInput.substring(0, start) + "  " + codeInput.substring(end);

      setTimeout(() => {
        target.selectionStart = target.selectionEnd = start + 2;
      }, 0);
    }
  }

  function printCode() {
    window.print();
  }

  function copyCode() {
    navigator.clipboard.writeText(codeInput);
    alert("Kode berhasil disalin!");
  }
</script>

<svelte:head>
  <title>Cetak Kode Programming - Les Balongarut</title>
  <meta name="description" content="Cetak kode pemrograman dengan tampilan yang rapi dan indah langsung dari browser. Gratis dan mudah digunakan." />
  <link rel="canonical" href="https://lesbalonggarut.my.id/cetak-kode" />

  <!-- Open Graph -->
  <meta property="og:type" content="website" />
  <meta property="og:url" content="https://lesbalonggarut.my.id/cetak-kode" />
  <meta property="og:title" content="Cetak Kode Programming - Les Balongarut" />
  <meta property="og:description" content="Cetak kode pemrograman dengan tampilan rapi dan indah langsung dari browser, gratis." />
  <meta property="og:site_name" content="Les Balongarut" />
  <meta property="og:locale" content="id_ID" />

  <!-- Twitter Card -->
  <meta name="twitter:card" content="summary" />
  <meta name="twitter:title" content="Cetak Kode Programming - Les Balongarut" />
  <meta name="twitter:description" content="Cetak kode pemrograman dengan tampilan rapi langsung dari browser." />
</svelte:head>

<div
  class="min-h-screen bg-slate-50 font-sans selection:bg-blue-200 selection:text-blue-900 flex flex-col relative overflow-x-hidden print:bg-white print:text-black print:p-0"
>
  <!-- Background Ambient -->
  <div class="absolute inset-0 z-0 pointer-events-none fixed print:hidden">
    <div
      class="absolute top-1/4 left-1/4 w-[600px] h-[600px] bg-white/40 rounded-full blur-[120px]"
    ></div>
    <div
      class="absolute bottom-1/4 right-1/4 w-[500px] h-[500px] bg-blue-100/50 rounded-full blur-[120px]"
    ></div>

  </div>

  <div
    class="relative z-10 w-full max-w-5xl mx-auto p-4 pt-24 md:p-8 md:pt-32 space-y-8 print:space-y-0 print:max-w-none print:w-full print:p-0"
  >
    <div class="text-center print:hidden flex flex-col items-center gap-4">
      <div class="flex items-center gap-3">
        <div
          class="w-6 h-6 border border-slate-200 bg-white flex items-center justify-center rotate-45 shadow-sm"
        >
          <div class="w-1.5 h-1.5 bg-blue-500"></div>
        </div>
        <span
          class="text-xs font-medium tracking-[0.2em] uppercase text-slate-500"
        >
          Utility System
        </span>
      </div>
      <h1
        class="text-3xl font-bold tracking-[0.1em] text-slate-900 uppercase drop-shadow-sm"
      >
        Cetak Kode Programming
      </h1>
      <p class="text-xs tracking-[0.1em] text-slate-600 font-medium uppercase">
        Paste script code Anda di bawah, lalu klik "Cetak" untuk mem-print
        dokumen.
      </p>
    </div>

    <!-- Input Section (Hidden when printing) -->
    <div class="flex flex-col space-y-4 print:hidden">
      <div
        class="flex items-center justify-between border-b border-slate-300 pb-3"
      >
        <label
          for="code-input"
          class="text-sm font-bold tracking-[0.1em] uppercase text-slate-800"
        >
          Source Code
        </label>
        <div class="flex space-x-3">
          <button
            onclick={() => (codeInput = "")}
            class="text-xs px-4 py-2 bg-white hover:bg-slate-50 text-slate-700 rounded-md border border-slate-300 transition-colors disabled:bg-white/50 disabled:text-slate-400 disabled:border-slate-200 tracking-wider uppercase font-extrabold backdrop-blur-sm shadow-sm hover:shadow-md"
            disabled={!codeInput}
          >
            Clear
          </button>
          <button
            onclick={copyCode}
            class="text-xs px-4 py-2 bg-white hover:bg-slate-50 text-slate-700 rounded-md border border-slate-300 transition-colors disabled:bg-white/50 disabled:text-slate-400 disabled:border-slate-200 tracking-wider uppercase font-extrabold backdrop-blur-sm shadow-sm hover:shadow-md"
            disabled={!codeInput}
          >
            Copy
          </button>
          <button
            onclick={printCode}
            class="group relative inline-flex items-center justify-center px-6 py-2 text-xs tracking-[0.1em] font-extrabold uppercase text-white border border-blue-600 hover:border-blue-500 bg-blue-600 hover:bg-blue-500 rounded-md transition-all duration-300 disabled:bg-blue-300 disabled:text-white disabled:border-blue-300 backdrop-blur-sm shadow-md overflow-hidden"
            disabled={!codeInput}
          >
            <span class="relative z-10 flex items-center gap-2">
              <svg
                class="w-4 h-4"
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
              Cetak
            </span>
          </button>
        </div>
      </div>
      <div class="relative group">
        <div
          class="absolute -inset-0.5 bg-gradient-to-r from-blue-400/30 to-indigo-500/30 rounded-xl blur opacity-0 group-hover:opacity-100 transition duration-500"
        ></div>
        <textarea
          id="code-input"
          bind:value={codeInput}
          onkeydown={handleKeydown}
          placeholder="Paste script code di sini..."
          class="relative w-full p-6 rounded-xl bg-white/60 backdrop-blur-md border border-slate-200 text-slate-900 font-mono text-base focus:ring-2 focus:ring-blue-400 focus:border-blue-400 outline-none resize-y min-h-[650px] shadow-xl whitespace-pre overflow-x-auto placeholder-slate-400/80"
          spellcheck="false"
        ></textarea>
      </div>
    </div>

    <!-- Output Section (Hidden on screen, ONLY visible when printing) -->
    <div class="hidden print:block print:w-full">
      {#if codeInput}
        <pre
          class="m-0 text-sm font-mono leading-relaxed text-black whitespace-pre-wrap break-words"
          style="tab-size: 2;"><code
            class="hljs !bg-transparent !text-black p-0"
            >{@html highlightedCode}</code
          ></pre>
      {/if}
    </div>
  </div>

  <!-- Noise Overlay for texture -->
  <div
    class="absolute inset-0 opacity-[0.02] pointer-events-none mix-blend-screen print:hidden"
    style="background-image: url('data:image/svg+xml,%3Csvg viewBox=%220 0 200 200%22 xmlns=%22http://www.w3.org/2000/svg%22%3E%3Cfilter id=%22noiseFilter%22%3E%3CfeTurbulence type=%22fractalNoise%22 baseFrequency=%220.8%22 numOctaves=%223%22 stitchTiles=%22stitch%22/%3E%3C/filter%3E%3Crect width=%22100%25%22 height=%22100%25%22 filter=%22url(%23noiseFilter)%22/%3E%3C/svg%3E');"
  ></div>
</div>

<style>
  /* Global overrides for highlight.js in this specific component to keep it transparent */
  :global(.hljs) {
    background: transparent !important;
    padding: 0 !important;
  }

  :global(.code-line) {
    /* Mencegah line break di tengah-tengah baris saat print perpindahan halaman jika memungkinkan */
    break-inside: avoid-page;
    /* Pastikan word-break diaplikasikan ke div level baris juga */
    word-break: break-word;
    overflow-wrap: break-word;
  }

  /* Print specific overrides to make code look like text editor on white paper */
  @media print {
    :global(body) {
      background: white !important;
    }
    :global(.hljs),
    :global(.hljs *) {
      color: black !important;
    }
    :global(.hljs-keyword),
    :global(.hljs-built_in) {
      font-weight: bold;
    }
    :global(.hljs-comment) {
      font-style: italic;
    }
  }
</style>
