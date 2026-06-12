<script lang="ts">
  import { onMount } from "svelte";
  import { page } from "$app/state";

  let isAuthenticated = $state(false);
  let user: { username: string } | null = $state(null);
  let isDropdownOpen = $state(false);

  const currentPath = $derived(page.url?.pathname || "/");

  onMount(async () => {
    try {
      const res = await fetch("http://localhost:8080/me", {
        credentials: "include",
      });
      if (res.ok) {
        const data = await res.json();
        if (data.authenticated) {
          isAuthenticated = true;
          user = data.user;
        }
      }
    } catch (e) {
      // Abaikan error jika tidak ada sesi
    }
  });

  async function handleLogout() {
    try {
      await fetch("http://localhost:8080/api/auth/logout", {
        method: "POST",
        credentials: "include",
      });
      isAuthenticated = false;
      user = null;
      isDropdownOpen = false;
      window.location.href = "/"; // redirect to home after logout
    } catch (e) {
      console.error(e);
    }
  }
</script>

<nav
  class="bg-[#0C134F]/80 backdrop-blur-md shadow-sm shadow-black/20 border-b border-indigo-400/30 fixed top-0 w-full z-50"
>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
    <div class="flex justify-between h-16 items-center">
      <div class="flex-shrink-0 flex items-center gap-3">
        <div
          class="w-10 h-10 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-xl flex items-center justify-center text-white font-extrabold text-xl shadow-md"
        >
          L
        </div>
        <a
          href="/"
          class="font-bold text-xl text-blue-100 tracking-tight no-underline hover:text-indigo-400 transition-colors"
          >Les Balongarut</a
        >
      </div>
      <div class="hidden sm:ml-6 sm:flex sm:space-x-8 items-center">
        <a
          href="/"
          class="inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium transition-colors no-underline {currentPath ===
          '/'
            ? 'border-indigo-400 text-white'
            : 'border-transparent text-blue-200 hover:border-zinc-500 hover:text-white'}"
          >Beranda</a
        >
        <a
          href="/mengetik"
          class="inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium transition-colors no-underline {currentPath.startsWith(
            '/mengetik',
          )
            ? 'border-indigo-400 text-white'
            : 'border-transparent text-blue-200 hover:border-zinc-500 hover:text-white'}"
          >Mengetik 10 Jari</a
        >
        <a
          href="/cetak-kode"
          class="inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium transition-colors no-underline {currentPath.startsWith(
            '/cetak-kode',
          )
            ? 'border-indigo-400 text-white'
            : 'border-transparent text-blue-200 hover:border-zinc-500 hover:text-white'}"
          >Cetak Kode</a
        >

        {#if isAuthenticated}
          <!-- Avatar Dropdown -->
          <div class="relative ml-4">
            <button
              onclick={() => (isDropdownOpen = !isDropdownOpen)}
              class="flex items-center gap-3 focus:outline-none cursor-pointer border-none bg-transparent p-0 group"
            >
              <div
                class="w-10 h-10 rounded-full bg-zinc-800 flex items-center justify-center text-indigo-400 font-bold border border-zinc-600 shadow-sm group-hover:ring-2 group-hover:ring-indigo-500 transition-all"
              >
                {user?.username ? user.username.charAt(0).toUpperCase() : "U"}
              </div>
            </button>

            {#if isDropdownOpen}
              <!-- Backdrop for closing dropdown -->
              <button
                class="fixed inset-0 z-40 w-full h-full cursor-default focus:outline-none border-none bg-transparent"
                onclick={() => (isDropdownOpen = false)}
                aria-label="Tutup dropdown"
              ></button>

              <!-- Dropdown Menu -->
              <div
                class="absolute right-0 mt-2 w-48 bg-zinc-900 rounded-xl shadow-lg ring-1 ring-white/10 z-50 overflow-hidden animate-in fade-in slide-in-from-top-2 duration-200"
              >
                <div
                  class="px-4 py-3 border-b border-zinc-800 bg-zinc-800/50 text-left"
                >
                  <p class="text-xs text-blue-200 m-0">Masuk sebagai</p>
                  <p class="text-sm font-bold text-white truncate m-0">
                    {user?.username}
                  </p>
                </div>
                <div class="py-1">
                  <a
                    href="/dashboard"
                    class="flex items-center px-4 py-2 text-sm text-blue-100 hover:bg-zinc-800 hover:text-white font-medium no-underline"
                  >
                    <svg
                      class="mr-3 w-4 h-4 text-zinc-400"
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
                    Dashboard
                  </a>
                </div>
                <div class="border-t border-zinc-800 py-1 text-left">
                  <button
                    onclick={handleLogout}
                    class="flex w-full items-center px-4 py-2 text-sm text-red-400 hover:bg-zinc-800 hover:text-red-300 font-medium border-none bg-transparent cursor-pointer"
                  >
                    <svg
                      class="mr-3 w-4 h-4 text-red-400"
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
                    Keluar Akun
                  </button>
                </div>
              </div>
            {/if}
          </div>
        {:else}
          <a
            href="/login"
            class="ml-4 px-4 py-2 rounded-lg bg-zinc-800 text-slate-200 hover:bg-zinc-700 font-medium text-sm transition-colors border border-zinc-600 no-underline"
            >Masuk Portal</a
          >
        {/if}
      </div>
    </div>
  </div>
</nav>
