<script lang="ts">
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";

  import type { Snippet } from "svelte";
  let { children }: { children: Snippet } = $props();

  let user: any = $state(null);
  let isLoading = $state(true);
  let isMobileMenuOpen = $state(false);

  const isLandingPage = $derived(
    page.url.pathname === "/quiz-app" || page.url.pathname === "/quiz-app/",
  );

  onMount(async () => {
    try {
      const res = await fetch(`/api/kuisapp/me`, { credentials: "include" });
      if (res.ok) {
        const data = await res.json();
        user = data.user;
      }
    } catch (e) {
      console.error("Auth check error:", e);
    } finally {
      isLoading = false;

      const path = page.url.pathname as string;
      const isAuthPage =
        path === "/quiz-app/login" || path === "/quiz-app/register";

      if (!user && !isAuthPage && !isLandingPage) {
        goto("/quiz-app/login");
      } else if (user && isAuthPage) {
        if (user.role === "admin") {
          goto("/quiz-app/admin");
        } else {
          goto("/quiz-app/dashboard");
        }
      } else if (
        user &&
        user.role === "admin" &&
        path.startsWith("/quiz-app/dashboard")
      ) {
        goto("/quiz-app/admin");
      } else if (
        user &&
        user.role !== "admin" &&
        path.startsWith("/quiz-app/admin")
      ) {
        goto("/quiz-app/dashboard");
      }
    }
  });

  async function handleLogout() {
    await fetch("/api/kuisapp/logout", {
      method: "POST",
      credentials: "include",
    });
    user = null;
    goto("/quiz-app");
  }
</script>

{#if isLoading}
  <div class="min-h-screen flex items-center justify-center bg-slate-50">
    <div
      class="w-10 h-10 border-4 border-slate-200 border-t-indigo-600 rounded-full animate-spin"
    ></div>
  </div>
{:else}
  <div class="min-h-screen bg-slate-50 flex flex-col font-sans">
    <!-- Navbar -->
    {#if user}
      <nav
        class="bg-white border-b border-slate-200 sticky top-0 z-50 shadow-sm"
      >
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div class="flex justify-between h-16">
            <div class="flex">
              <div class="flex-shrink-0 flex items-center">
                <span
                  class="text-xl font-black bg-gradient-to-r from-indigo-600 to-violet-600 bg-clip-text text-transparent"
                >
                  LB Quiz
                </span>
              </div>
              <div class="hidden sm:ml-6 sm:flex sm:space-x-8">
                {#if user.role !== "admin"}
                  <a
                    href="/quiz-app/dashboard"
                    class="{page.url.pathname === '/quiz-app/dashboard'
                      ? 'border-indigo-500 text-slate-900'
                      : 'border-transparent text-slate-500 hover:border-slate-300 hover:text-slate-700'} inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium transition-colors"
                  >
                    Dashboard
                  </a>
                  <a
                    href="/quiz-app/dashboard/history"
                    class="{page.url.pathname === '/quiz-app/dashboard/history'
                      ? 'border-indigo-500 text-slate-900'
                      : 'border-transparent text-slate-500 hover:border-slate-300 hover:text-slate-700'} inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium transition-colors"
                  >
                    Riwayat Kuis
                  </a>
                {/if}
                {#if user.role === "admin"}
                  <a
                    href="/quiz-app/admin"
                    class="{page.url.pathname.startsWith('/quiz-app/admin')
                      ? 'border-indigo-500 text-slate-900'
                      : 'border-transparent text-slate-500 hover:border-slate-300 hover:text-slate-700'} inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium transition-colors"
                  >
                    Admin Panel
                  </a>
                {/if}
              </div>
            </div>
            <div class="flex items-center gap-2">
              <button
                onclick={handleLogout}
                class="hidden sm:inline-flex items-center justify-center gap-2 px-4 py-2 bg-red-50 text-red-600 hover:bg-red-100 hover:text-red-700 font-bold rounded-xl transition-colors text-sm"
                title="Logout"
              >
                <svg
                  class="w-4 h-4"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  ><path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
                  ></path></svg
                >
                Keluar
              </button>

              <!-- Mobile menu button -->
              <button
                onclick={() => (isMobileMenuOpen = !isMobileMenuOpen)}
                class="sm:hidden inline-flex items-center justify-center p-2 rounded-md text-slate-400 hover:text-slate-500 hover:bg-slate-100 transition-colors"
              >
                <svg
                  class="h-6 w-6"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  {#if isMobileMenuOpen}
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M6 18L18 6M6 6l12 12"
                    />
                  {:else}
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M4 6h16M4 12h16M4 18h16"
                    />
                  {/if}
                </svg>
              </button>
            </div>
          </div>
        </div>

        <!-- Mobile Menu Panel -->
        {#if isMobileMenuOpen}
          <div class="sm:hidden border-t border-slate-200 bg-white">
            <div class="pt-2 pb-3 space-y-1">
              {#if user.role !== "admin"}
                <a
                  href="/quiz-app/dashboard"
                  onclick={() => (isMobileMenuOpen = false)}
                  class="{page.url.pathname === '/quiz-app/dashboard'
                    ? 'bg-indigo-50 border-indigo-500 text-indigo-700'
                    : 'border-transparent text-slate-600 hover:bg-slate-50 hover:border-slate-300 hover:text-slate-800'} block pl-3 pr-4 py-2 border-l-4 text-base font-medium"
                >
                  Dashboard
                </a>
                <a
                  href="/quiz-app/dashboard/history"
                  onclick={() => (isMobileMenuOpen = false)}
                  class="{page.url.pathname === '/quiz-app/dashboard/history'
                    ? 'bg-indigo-50 border-indigo-500 text-indigo-700'
                    : 'border-transparent text-slate-600 hover:bg-slate-50 hover:border-slate-300 hover:text-slate-800'} block pl-3 pr-4 py-2 border-l-4 text-base font-medium"
                >
                  Riwayat Kuis
                </a>
              {/if}
              {#if user.role === "admin"}
                <a
                  href="/quiz-app/admin"
                  onclick={() => (isMobileMenuOpen = false)}
                  class="{page.url.pathname.startsWith('/quiz-app/admin')
                    ? 'bg-indigo-50 border-indigo-500 text-indigo-700'
                    : 'border-transparent text-slate-600 hover:bg-slate-50 hover:border-slate-300 hover:text-slate-800'} block pl-3 pr-4 py-2 border-l-4 text-base font-medium"
                >
                  Admin Panel
                </a>
              {/if}
              <button
                onclick={() => {
                  isMobileMenuOpen = false;
                  handleLogout();
                }}
                class="w-full text-left border-transparent text-red-600 hover:bg-red-50 hover:border-red-300 block pl-3 pr-4 py-2 border-l-4 text-base font-medium"
              >
                Keluar
              </button>
            </div>
          </div>
        {/if}
      </nav>
    {/if}

    <!-- Main Content -->
    <main
      class="flex-1 w-full {isLandingPage
        ? ''
        : 'max-w-7xl mx-auto p-4 sm:p-6 lg:p-8'}"
    >
      {@render children()}
    </main>
  </div>
{/if}
