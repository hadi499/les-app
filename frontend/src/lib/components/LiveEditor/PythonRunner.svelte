<script lang="ts">
  import { onMount } from 'svelte';

  let { code = '' } = $props<{ code?: string }>();

  let output: string[] = $state([]);
  let pyodide: any = $state(null);
  let isLoading = $state(true);
  let runError: string | null = $state(null);

  onMount(async () => {
    try {
      // Load pyodide script dynamically
      if (!(window as any).loadPyodide) {
        const script = document.createElement('script');
        script.src = 'https://cdn.jsdelivr.net/pyodide/v0.25.0/full/pyodide.js';
        document.head.appendChild(script);
        
        await new Promise((resolve) => {
          script.onload = resolve;
        });
      }

      // Initialize Pyodide
      pyodide = await (window as any).loadPyodide();
      
      // Redirect stdout
      pyodide.setStdout({
        batched: (msg: string) => {
          output = [...output, msg];
        }
      });
      isLoading = false;
    } catch (e) {
      console.error("Failed to load pyodide", e);
      runError = "Gagal memuat Python. Coba muat ulang halaman.";
      isLoading = false;
    }
  });

  async function runCode() {
    if (!pyodide) return;
    
    output = [];
    runError = null;
    
    try {
      await pyodide.runPythonAsync(code);
    } catch (e: any) {
      // Simplify error message for beginners
      console.error(e);
      let errorMsg = e.message || e.toString();
      
      // Basic heuristic to simplify python stack trace
      if (errorMsg.includes('SyntaxError')) {
        runError = "Ada kesalahan penulisan kode (SyntaxError). Coba periksa lagi kode yang kamu tulis.";
      } else if (errorMsg.includes('NameError')) {
        runError = "Ada variabel atau perintah yang tidak dikenali (NameError). Coba periksa ejaannya.";
      } else {
        // Just show the last line of the stack trace which usually contains the actual error
        const lines = errorMsg.split('\\n').filter((l: string) => l.trim() !== '');
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
      disabled={isLoading}
      onclick={runCode}
    >
      {isLoading ? 'Memuat Python...' : 'Run Code'}
    </button>
  </div>
  
  <div class="flex-1 p-4 overflow-y-auto text-sm">
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
