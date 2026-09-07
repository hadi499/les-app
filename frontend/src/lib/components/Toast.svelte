<script lang="ts">
  import { toast } from "$lib/stores/toast.svelte";
</script>

{#if toast.visible}
  <div
    class="fixed top-4 right-4 z-[100] px-4 py-3 rounded-lg shadow-lg text-sm font-medium transition-all {toast.type ===
    'success'
      ? 'bg-green-300 text-slate-900'
      : 'bg-red-500 text-slate-900'}"
  >
    {toast.message}
  </div>
{/if}
