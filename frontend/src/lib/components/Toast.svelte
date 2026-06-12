<script lang="ts">
  import { toast } from '$lib/stores/toast.svelte';
</script>

{#if toast.visible}
  <div class="fixed top-4 right-4 z-[100] px-4 py-3 rounded-lg shadow-lg text-sm font-medium transition-all {toast.type === 'success' ? 'bg-green-600 text-white' : 'bg-red-600 text-white'}">
    {toast.message}
  </div>
{/if}
