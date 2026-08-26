<script lang="ts">
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import { onMount, setContext } from "svelte";

  import type { Snippet } from "svelte";
  let { children }: { children: Snippet } = $props();

  let user: any = $state(null);
  let isLoading = $state(true);
  let isMobileMenuOpen = $state(false);

  setContext("quizUser", {
    get current() {
      return user;
    },
    get loading() {
      return isLoading;
    },
  });

  const isLandingPage = $derived(
    page.url.pathname === "/quiz-app" || page.url.pathname === "/quiz-app/",
  );

  const isQuizAreaPage = $derived(
    page.url.pathname.startsWith("/quiz-app/dashboard/quizzes/"),
  );

  onMount(async () => {
    try {
      const res = await fetch(`/api/kuisapp/me?t=${Date.now()}`, {
        credentials: "include",
        cache: "no-store",
      });
      if (res.ok) {
        const data = await res.json();
        user = data.user;
      }
    } catch (e) {
      console.error("Auth check error:", e);
    } finally {
      isLoading = false;
    }
  });

  $effect(() => {
    if (!isLoading) {
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
    {#if user && !isLandingPage}
      <nav
        class="bg-white border-b border-slate-200 sticky top-0 z-50 shadow-sm {user.role !==
        'admin'
          ? 'hidden sm:block'
          : ''}"
      >
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div class="flex justify-between h-16">
            <div class="flex">
              {#if !isQuizAreaPage}
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
                      class="{page.url.pathname ===
                      '/quiz-app/dashboard/history'
                        ? 'border-indigo-500 text-slate-900'
                        : 'border-transparent text-slate-500 hover:border-slate-300 hover:text-slate-700'} inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium transition-colors"
                    >
                      Riwayat Kuis
                    </a>
                  {/if}
                </div>
              {/if}
            </div>

            <div class="flex items-center gap-2">
              {#if !isQuizAreaPage}
                <button
                  onclick={handleLogout}
                  class="inline-flex items-center justify-center gap-2 px-3 sm:px-4 py-2 bg-red-50 text-red-600 hover:bg-red-100 hover:text-red-700 font-bold rounded-xl transition-colors text-sm"
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
              {/if}
            </div>
          </div>
        </div>
      </nav>
    {/if}

    <!-- Main Content -->
    <main
      class="flex-1 w-full {isLandingPage
        ? ''
        : 'max-w-7xl mx-auto p-4 sm:p-6 lg:p-8'} {user && user.role !== 'admin'
        ? 'pb-24 sm:pb-8'
        : ''}"
    >
      {@render children()}
    </main>

    <!-- Bottom Navigation Bar for Mobile (User Only) -->
    {#if user && user.role !== "admin" && !isLandingPage && !isQuizAreaPage}
      <nav
        class="sm:hidden fixed bottom-0 w-full bg-white border-t border-slate-200 z-50 pb-[env(safe-area-inset-bottom)] shadow-[0_-2px_10px_rgba(0,0,0,0.05)]"
      >
        <div class="flex justify-around items-center h-16">
          <a
            href="/quiz-app/dashboard"
            class="flex flex-col items-center justify-center w-full h-full {page
              .url.pathname === '/quiz-app/dashboard'
              ? 'text-indigo-600'
              : 'text-slate-500 hover:text-slate-800 hover:bg-slate-50'} transition-colors"
          >
            <svg
              class="w-6 h-6 mb-1 {page.url.pathname === '/quiz-app/dashboard'
                ? 'fill-indigo-100/50'
                : ''}"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              ><path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
              ></path></svg
            >
            <span class="text-[10px] font-bold tracking-wide">Beranda</span>
          </a>
          <a
            href="/quiz-app/dashboard/history"
            class="flex flex-col items-center justify-center w-full h-full {page
              .url.pathname === '/quiz-app/dashboard/history'
              ? 'text-indigo-600'
              : 'text-slate-500 hover:text-slate-800 hover:bg-slate-50'} transition-colors"
          >
            <svg
              class="w-6 h-6 mb-1 {page.url.pathname ===
              '/quiz-app/dashboard/history'
                ? 'fill-indigo-100/50'
                : ''}"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              ><path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
              ></path></svg
            >
            <span class="text-[10px] font-bold tracking-wide">Riwayat</span>
          </a>
          <button
            onclick={handleLogout}
            class="flex flex-col items-center justify-center w-full h-full text-slate-500 hover:text-red-600 hover:bg-red-50 transition-colors"
          >
            <svg
              class="w-6 h-6 mb-1"
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
            <span class="text-[10px] font-bold tracking-wide">Keluar</span>
          </button>
        </div>
      </nav>
    {/if}
  </div>
{/if}
