<script lang="ts">
  import { onMount } from 'svelte';
  
  let { children } = $props();
  let isAllowed = $state(false);

  onMount(async () => {
    try {
      const res = await fetch(`/me`, { credentials: 'include' });
      if (!res.ok) {
        window.location.href = '/dashboard';
        return;
      }
      const data = await res.json();
      if (!data.authenticated || data.user?.role !== 'teacher') {
        window.location.href = '/dashboard';
        return;
      }
      isAllowed = true;
    } catch (e) {
      window.location.href = '/dashboard';
    }
  });
</script>

{#if isAllowed}
  {@render children()}
{/if}
