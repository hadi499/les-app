<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import PyodideWorker from '$lib/workers/pyodide.worker.ts?worker';

  let { code = '' } = $props<{ code?: string }>();

  let output: string[] = $state([]);
  let isLoading = $state(true); // Loading on initial mount in background
  let isExecuting = $state(false); // When code is actually running
  let runError: string | null = $state(null);

  let worker: Worker;
  let messageIdCounter = 0;
  let resolvePromises: Record<number, { resolve: any, reject: any }> = {};

  onMount(() => {
    worker = new PyodideWorker();
    
    worker.onmessage = (event) => {
      const { id, type, msg, error } = event.data;
      
      if (type === 'STDOUT') {
        output = [...output, msg];
      } else if (type === 'INIT_DONE') {
        isLoading = false;
        if (resolvePromises[id]) {
          resolvePromises[id].resolve(null);
          delete resolvePromises[id];
        }
      } else if (type === 'RUN_DONE') {
        isExecuting = false;
        if (resolvePromises[id]) {
          resolvePromises[id].resolve(null);
          delete resolvePromises[id];
        }
      } else if (type === 'ERROR') {
        isLoading = false;
        isExecuting = false;
        if (resolvePromises[id]) {
          resolvePromises[id].reject(new Error(error));
          delete resolvePromises[id];
        }
      }
    };

    // Preload Pyodide in the background immediately
    const id = messageIdCounter++;
    resolvePromises[id] = {
      resolve: () => {},
      reject: (err: any) => { runError = "Gagal memuat Python: " + err.message; }
    };
    worker.postMessage({ id, type: 'INIT' });
  });

  onDestroy(() => {
    if (worker) worker.terminate();
  });

  async function runCode() {
    if (isLoading || isExecuting) return;
    
    output = [];
    runError = null;
    isExecuting = true;
    
    const id = messageIdCounter++;
    const runPromise = new Promise((resolve, reject) => {
      resolvePromises[id] = { resolve, reject };
    });
    
    worker.postMessage({ id, type: 'RUN', code });
    
    try {
      await runPromise;
    } catch (e: any) {
      let errorMsg = e.message || e.toString();
      
      if (errorMsg.includes('SyntaxError')) {
        runError = "Ada kesalahan penulisan kode (SyntaxError). Coba periksa lagi kode yang kamu tulis.";
      } else if (errorMsg.includes('NameError')) {
        runError = "Ada variabel atau perintah yang tidak dikenali (NameError). Coba periksa ejaannya.";
      } else {
        const lines = errorMsg.split('\n').filter((l: string) => l.trim() !== '');
        runError = "Error: " + lines[lines.length - 1];
      }
    }
  }
  
  $effect(() => {
    if (code === '') {
      output = [];
      runError = null;
    }
  });
</script>

<div class="flex flex-col h-full bg-gray-900 rounded overflow-hidden text-gray-100 font-mono">
  <div class="flex justify-between items-center p-2 bg-gray-800 border-b border-gray-700">
    <span class="text-sm font-semibold">Console (Python)</span>
    <button 
      class="bg-blue-600 hover:bg-blue-700 text-white text-xs py-1 px-3 rounded disabled:opacity-50"
      disabled={isLoading || isExecuting}
      onclick={runCode}
    >
      {isLoading ? 'Menyiapkan Mesin...' : (isExecuting ? 'Menjalankan...' : 'Run Code')}
    </button>
  </div>
  
  <div class="flex-1 p-4 overflow-y-auto text-sm relative">
    {#if isLoading}
      <div class="absolute inset-0 z-10 flex flex-col items-center justify-center bg-gray-900/90 text-white backdrop-blur-sm">
        <svg class="animate-spin h-8 w-8 text-blue-500 mb-3" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <span class="text-sm font-medium animate-pulse text-blue-200">Sedang memuat sistem Python...</span>
        <span class="text-xs text-blue-300 mt-1">(Tunggu sebentar, tidak akan macet)</span>
      </div>
    {/if}

    {#if isExecuting}
      <div class="absolute inset-0 z-10 flex items-center justify-center bg-gray-900/50 text-white">
        <span class="animate-pulse">Menjalankan kode...</span>
      </div>
    {/if}

    {#if runError}
      <div class="text-red-400 mb-2 whitespace-pre-wrap">{runError}</div>
    {/if}
    
    {#each output as line}
      <div class="whitespace-pre-wrap">{line}</div>
    {/each}
    
    {#if output.length === 0 && !runError}
      <div class="text-gray-500 italic">Output akan muncul di sini...</div>
    {/if}
  </div>
</div>
