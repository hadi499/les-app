<script lang="ts">
  // import favicon from "$lib/assets/favicon.svg";
  import favicon from "$lib/assets/lb8.png";
  import "../app.css";
  import { page } from "$app/state";
  import Navbar from "$lib/components/Navbar.svelte";
  import Toast from "$lib/components/Toast.svelte";
  import { onMount, setContext } from "svelte";
  import type { Snippet } from "svelte";

  let { children }: { children: Snippet } = $props();

  let authState = $state({
    isAuthenticated: false,
    authChecked: false,
  });

  setContext("authState", authState);

  onMount(async () => {
    try {
      const res = await fetch(`/me`, { credentials: "include" });
      if (res.ok) {
        const data = await res.json();
        if (data.authenticated) {
          authState.isAuthenticated = true;
        }
      }
    } catch (e) {
      // ignore
    } finally {
      authState.authChecked = true;
    }
  });

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
  <link rel="icon" type="image/png" href={favicon} />
</svelte:head>

{#if showNavbar()}
  <div class="print:hidden">
    <Navbar />
  </div>
{/if}

<Toast />
{@render children()}
