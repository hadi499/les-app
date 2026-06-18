<script lang="ts">
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";
  import { fetchAllProgress, fetchAllGameScores } from "$lib/stores/progress";

  let { children } = $props();

  let isAuthenticated = $state(false);

  onMount(async () => {
    try {
      const res = await fetch(`/me`, {
        credentials: "include",
      });
      if (res.ok) {
        const data = await res.json();
        if (data.authenticated) {
          await fetchAllProgress();
          await fetchAllGameScores();
          isAuthenticated = true;
          return;
        }
      }
    } catch (error) {
      console.error("Auth check error:", error);
    }

    // Jika tidak terautentikasi, alihkan ke halaman login
    goto("/login");
  });
</script>

{#if isAuthenticated}
  <div class="typing-app-wrapper pt-16">
    {@render children()}
  </div>
{:else}
  <div class="min-h-screen flex items-center justify-center bg-orange-100">
    <div class="animate-pulse flex flex-col items-center">
      <div
        class="w-12 h-12 border-4 border-indigo-500/30 border-t-indigo-500 rounded-full animate-spin"
      ></div>
      <p
        class="mt-4 text-orange-800 font-medium tracking-widest text-sm uppercase"
      >
        Memeriksa sesi...
      </p>
    </div>
  </div>
{/if}
