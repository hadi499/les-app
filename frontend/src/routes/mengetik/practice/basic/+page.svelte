<script lang="ts">
  import QuickPractice from "$lib/components/QuickPractice.svelte";
</script>

<svelte:head>
  <title>Latihan Dasar - Ketik 10 Jari</title>
</svelte:head>

<div
  class="min-h-screen bg-slate-50 text-slate-900 selection:bg-blue-200 selection:text-blue-900"
>
  <main class="max-w-4xl mx-auto p-4 sm:p-8">
    <QuickPractice />
    <a
      href="/mengetik"
      class="inline-block text-slate-600 hover:text-slate-900 text-sm mt-4 transition-colors no-underline font-medium tracking-wide"
      >← Kembali ke Pelajaran</a
    >
  </main>
</div>
