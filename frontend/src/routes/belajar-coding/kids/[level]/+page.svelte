<script lang="ts">
  import { page } from '$app/stores';
  import { levels } from '$lib/components/BlockPlayground/levels';
  import BlocklyWorkspace from '$lib/components/BlockPlayground/BlocklyWorkspace.svelte';
  import BlockRunner from '$lib/components/BlockPlayground/BlockRunner.svelte';
  import { goto } from '$app/navigation';

  let levelId = $derived(parseInt($page.params.level || '1', 10));
  let level = $derived(levels[levelId]);

  let currentBlocks: any[] = $state([]);

  function handleBlocksChanged(blocks: any[]) {
    currentBlocks = blocks;
  }

  function nextLevel() {
    if (levels[levelId + 1]) {
      goto(`/belajar-coding/kids/${levelId + 1}`);
    } else {
      alert("Selamat! Kamu telah menyelesaikan semua level!");
      goto('/belajar-coding');
    }
  }
</script>

<svelte:head>
  <title>Coding Anak - {level ? level.title : 'Level Tidak Ditemukan'}</title>
</svelte:head>

<div class="min-h-screen lg:h-screen w-full flex flex-col bg-slate-100 overflow-y-auto lg:overflow-hidden font-sans">
  <header class="bg-white border-b border-slate-200 px-6 py-4 flex justify-between items-center shrink-0 shadow-sm z-10">
    <div class="flex items-center gap-4">
      <div class="w-12 h-12 bg-green-100 text-green-600 rounded-2xl flex items-center justify-center shadow-inner">
        <svg class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
      </div>
      <div>
        <h1 class="text-2xl font-extrabold text-slate-800">
          Belajar Logika (Anak)
        </h1>
        <p class="text-slate-500 font-medium">Susun blok untuk mencapai bintang!</p>
      </div>
    </div>
    
    <a 
      href="/belajar-coding"
      class="px-5 py-2.5 text-sm font-bold text-slate-600 bg-slate-100 hover:bg-slate-200 rounded-xl transition"
    >
      Tutup
    </a>
  </header>

  {#if level}
    <div class="flex-1 p-4 sm:p-6 flex flex-col lg:flex-row gap-6 lg:overflow-hidden">
      
      <!-- Toolbox and Workspace (Left) -->
      <div class="min-h-[500px] lg:min-h-0 flex-[1.2] bg-white rounded-2xl shadow-sm border border-slate-200 overflow-hidden relative flex flex-col">
        <div class="bg-blue-50 text-blue-800 px-4 py-3 text-sm font-medium border-b border-blue-100 flex items-center gap-2 z-10 shrink-0 shadow-sm">
          <svg class="w-5 h-5 animate-pulse text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 5l7 7-7 7M5 5l7 7-7 7"></path></svg>
          Petunjuk: Seret blok dari kiri ke area putih ini, sambungkan, lalu klik "Jalankan"!
        </div>
        <div class="flex-1 relative">
          <!-- Re-keying BlocklyWorkspace so it completely re-initializes when level (and allowedBlocks) changes -->
          {#key level.id}
            <BlocklyWorkspace 
              allowedBlocks={level.allowedBlocks}
              onBlocksChanged={handleBlocksChanged}
            />
          {/key}
        </div>
      </div>

      <!-- Execution Canvas (Right) -->
      <div class="min-h-[500px] lg:min-h-0 flex-1 flex flex-col">
        <!-- Re-key BlockRunner as well to reset cleanly -->
        {#key level.id}
          <BlockRunner 
            level={level}
            blockSequence={currentBlocks}
            onNextLevel={nextLevel}
          />
        {/key}
      </div>
      
    </div>
  {:else}
    <div class="flex-1 flex flex-col items-center justify-center p-6 text-center">
      <h2 class="text-3xl font-bold text-slate-800 mb-4">Level tidak ditemukan!</h2>
      <a href="/belajar-coding" class="px-6 py-3 bg-blue-600 text-white font-bold rounded-xl shadow-lg">Kembali</a>
    </div>
  {/if}
</div>
