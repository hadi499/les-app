<script lang="ts">
  import favicon from "$lib/assets/favicon.svg";
  import "../app.css";
  import { page } from "$app/state";
  import Navbar from "$lib/components/Navbar.svelte";
  import type { Snippet } from "svelte";

  let { children }: { children: Snippet } = $props();

  const showNavbar = $derived(() => {
    const path = page.url?.pathname || "/";
    // Hide navbar in dashboard, login, register, and select-user
    if (
      path.startsWith("/dashboard") ||
      path.startsWith("/login") ||
      path.startsWith("/register") ||
      path.includes("select-user")
    ) {
      return false;
    }
    return true;
  });
</script>

<svelte:head>
  <link rel="icon" href={favicon} />
</svelte:head>

{#if showNavbar()}
  <div class="print:hidden">
    <Navbar />
  </div>
{/if}

{@render children()}
