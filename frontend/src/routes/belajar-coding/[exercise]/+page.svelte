<script lang="ts">
  import { page } from "$app/stores";
  import { onMount } from "svelte";
  import {
    exercises,
    type Language,
  } from "$lib/components/LiveEditor/starters";
  import Editor from "$lib/components/LiveEditor/Editor.svelte";
  import Preview from "$lib/components/LiveEditor/Preview.svelte";
  import PythonRunner from "$lib/components/LiveEditor/PythonRunner.svelte";

  let exerciseId = $derived($page.params.exercise || "");
  let exercise = $derived(exercises[exerciseId]);

  let htmlCode = $state("");
  let cssCode = $state("");
  let jsCode = $state("");
  let pythonCode = $state("");

  let previewHtml = $state("");
  let previewCss = $state("");
  let previewJs = $state("");

  function runWebCode() {
    previewHtml = htmlCode;
    previewCss = cssCode;
    previewJs = jsCode;
  }

  let activeTab: Language = $state("html");

  function getStorageKey(lang: string) {
    return `lesbalong_coding_${exerciseId}_${lang}`;
  }

  function loadCode() {
    if (!exercise) return;

    if (exercise.type === "web") {
      htmlCode =
        localStorage.getItem(getStorageKey("html")) ??
        exercise.starter.html ??
        "";
      cssCode =
        localStorage.getItem(getStorageKey("css")) ??
        exercise.starter.css ??
        "";
      jsCode =
        localStorage.getItem(getStorageKey("js")) ?? exercise.starter.js ?? "";
      activeTab = "html";
      runWebCode();
    } else if (exercise.type === "python") {
      pythonCode =
        localStorage.getItem(getStorageKey("python")) ??
        exercise.starter.python ??
        "";
      activeTab = "python";
    }
  }

  function resetCode() {
    if (
      !exercise ||
      !confirm(
        "Apakah kamu yakin ingin mengulang dari awal? Semua perubahan akan hilang.",
      )
    )
      return;

    if (exercise.type === "web") {
      localStorage.removeItem(getStorageKey("html"));
      localStorage.removeItem(getStorageKey("css"));
      localStorage.removeItem(getStorageKey("js"));
    } else if (exercise.type === "python") {
      localStorage.removeItem(getStorageKey("python"));
    }
    loadCode();
  }

  onMount(() => {
    loadCode();
  });

  // Save on change (debounced)
  let saveTimeout: ReturnType<typeof setTimeout>;

  function saveCode() {
    clearTimeout(saveTimeout);
    saveTimeout = setTimeout(() => {
      if (exercise?.type === "web") {
        localStorage.setItem(getStorageKey("html"), htmlCode);
        localStorage.setItem(getStorageKey("css"), cssCode);
        localStorage.setItem(getStorageKey("js"), jsCode);
      } else if (exercise?.type === "python") {
        localStorage.setItem(getStorageKey("python"), pythonCode);
      }
    }, 500);
  }

  $effect(() => {
    // Re-run saveCode whenever code changes
    htmlCode;
    cssCode;
    jsCode;
    pythonCode;
    saveCode();
  });
</script>

<svelte:head>
  <title
    >{exercise ? exercise.title : "Belajar Coding"} - Les Balong Garut</title
  >
</svelte:head>

<div class="h-screen w-full flex flex-col bg-gray-50 font-sans">
  <!-- Header -->
  <header
    class="bg-white border-b border-gray-200 px-4 py-3 flex justify-between items-center shrink-0"
  >
    <div>
      <h1 class="text-xl font-bold text-gray-800">
        {#if exercise}
          {exercise.title}
        {:else}
          Latihan Tidak Ditemukan
        {/if}
      </h1>
    </div>
    <div class="flex items-center gap-3">
      <button
        class="px-5 py-2.5 text-sm font-semibold text-slate-600 hover:text-slate-900 bg-slate-100 hover:bg-slate-200 rounded-xl transition-colors"
        onclick={resetCode}
      >
        Ulangi dari Awal
      </button>
      <a
        href="/belajar-coding"
        class="px-5 py-2.5 flex items-center gap-2 text-sm font-semibold text-slate-600 hover:text-slate-900 bg-white border border-slate-200 hover:border-slate-300 hover:bg-slate-50 hover:shadow-sm rounded-xl transition-all"
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
            d="M10 19l-7-7m0 0l7-7m-7 7h18"
          ></path></svg
        >
        Kembali
      </a>
    </div>
  </header>

  <!-- Workspace -->
  {#if exercise}
    <div class="flex-1 overflow-hidden p-4 flex gap-4 h-full">
      {#if exercise.type === "web"}
        <!-- Web Editor Panel -->
        <div
          class="flex-1 flex flex-col bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden min-w-0"
        >
          <div class="flex border-b border-gray-200 shrink-0">
            <button
              class="px-4 py-2 text-sm font-medium {activeTab === 'html'
                ? 'text-blue-600 border-b-2 border-blue-600'
                : 'text-gray-500 hover:text-gray-700'}"
              onclick={() => (activeTab = "html")}
            >
              HTML
            </button>
            <button
              class="px-4 py-2 text-sm font-medium {activeTab === 'css'
                ? 'text-blue-600 border-b-2 border-blue-600'
                : 'text-gray-500 hover:text-gray-700'}"
              onclick={() => (activeTab = "css")}
            >
              CSS
            </button>
            <button
              class="px-4 py-2 text-sm font-medium {activeTab === 'js'
                ? 'text-blue-600 border-b-2 border-blue-600'
                : 'text-gray-500 hover:text-gray-700'}"
              onclick={() => (activeTab = "js")}
            >
              JavaScript
            </button>
          </div>
          <div class="flex-1 overflow-hidden relative">
            <div class="absolute inset-0" class:hidden={activeTab !== "html"}>
              <Editor bind:value={htmlCode} language="html" />
            </div>
            <div class="absolute inset-0" class:hidden={activeTab !== "css"}>
              <Editor bind:value={cssCode} language="css" />
            </div>
            <div class="absolute inset-0" class:hidden={activeTab !== "js"}>
              <Editor bind:value={jsCode} language="js" />
            </div>
          </div>
        </div>

        <!-- Web Preview Panel -->
        <div
          class="flex-1 flex flex-col bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden min-w-0"
        >
          <div
            class="px-4 py-2 border-b border-gray-200 shrink-0 flex justify-between items-center"
          >
            <span class="text-sm font-medium text-gray-700">Preview</span>
            <button
              class="bg-blue-600 hover:bg-blue-700 text-white text-xs py-1 px-3 rounded shadow-sm flex items-center gap-1 transition-colors"
              onclick={runWebCode}
            >
              <svg
                class="w-3.5 h-3.5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"
                ></path><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                ></path></svg
              >
              Run Code
            </button>
          </div>
          <div class="flex-1 p-2 bg-gray-100">
            <Preview html={previewHtml} css={previewCss} js={previewJs} />
          </div>
        </div>
      {:else if exercise.type === "python"}
        <!-- Python Editor Panel -->
        <div
          class="flex-[1.5] flex flex-col bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden min-w-0"
        >
          <div class="px-4 py-2 border-b border-gray-200 shrink-0">
            <span class="text-sm font-medium text-gray-700">script.py</span>
          </div>
          <div class="flex-1 relative overflow-hidden">
            <Editor bind:value={pythonCode} language="python" />
          </div>
        </div>

        <!-- Python Console Panel -->
        <div
          class="flex-1 flex flex-col rounded-lg shadow-sm overflow-hidden min-w-0"
        >
          <PythonRunner code={pythonCode} />
        </div>
      {/if}
    </div>
  {:else}
    <div class="flex-1 flex items-center justify-center">
      <div class="text-center">
        <h2 class="text-2xl font-bold text-gray-700">
          Latihan tidak ditemukan
        </h2>
        <p class="text-gray-500 mt-2">Pastikan URL yang kamu masukkan benar.</p>
        <a
          href="/dashboard"
          class="inline-block mt-4 px-4 py-2 bg-blue-600 text-white rounded-md"
          >Kembali ke Dashboard</a
        >
      </div>
    </div>
  {/if}
</div>
