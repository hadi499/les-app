<script lang="ts">
  import { onMount } from "svelte";
  import { page } from "$app/state";

  let isAuthenticated = $state(false);
  let user: { username: string } | null = $state(null);
  let isDropdownOpen = $state(false);
  let authChecked = $state(false);

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
    } finally {
      authChecked = true;
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
  class="bg-[#D98C4A]/95 backdrop-blur-md shadow-sm shadow-orange-900/15 border-b border-orange-400 fixed top-0 w-full z-50 subpixel-antialiased"
>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
    <div class="flex justify-between h-16 items-center">
      <div class="flex-shrink-0 flex items-center gap-3">
        <div
          class="w-10 h-10 bg-white rounded-xl flex items-center justify-center text-orange-600 font-extrabold text-xl shadow-md"
        >
          LB
        </div>
        <a
          href="/"
          class="font-extrabold text-xl text-orange-950 tracking-tight drop-shadow-sm no-underline hover:text-orange-100 transition-colors"
          >Les Balongarut</a
        >
      </div>
      <div class="hidden sm:ml-6 sm:flex sm:space-x-8 items-center">
        <a
          href="/"
          class="inline-flex items-center px-1 pt-1 border-b-2 text-sm font-semibold drop-shadow-sm transition-colors no-underline {currentPath ===
          '/'
            ? 'border-white text-orange-950'
            : 'border-transparent text-orange-100 hover:border-white/50 hover:text-orange-950'}"
          >Beranda</a
        >
        <a
          href="/mengetik"
          class="inline-flex items-center px-1 pt-1 border-b-2 text-sm font-semibold drop-shadow-sm transition-colors no-underline {currentPath.startsWith(
            '/mengetik',
          )
            ? 'border-white text-orange-950'
            : 'border-transparent text-orange-100 hover:border-white/50 hover:text-orange-950'}"
          >Mengetik 10 Jari</a
        >
        <a
          href="/cetak-kode"
          class="inline-flex items-center px-1 pt-1 border-b-2 text-sm font-semibold drop-shadow-sm transition-colors no-underline {currentPath.startsWith(
            '/cetak-kode',
          )
            ? 'border-white text-orange-950'
            : 'border-transparent text-orange-100 hover:border-white/50 hover:text-orange-950'}"
          >Cetak Kode</a
        >

        {#if !authChecked}
          <div class="ml-4 w-10 h-10"></div>
        {:else if isAuthenticated}
          <!-- Avatar Dropdown -->
          <div class="relative ml-4">
            <button
              onclick={() => (isDropdownOpen = !isDropdownOpen)}
              class="flex items-center gap-3 focus:outline-none cursor-pointer border-none bg-transparent p-0 group"
            >
              <div
                class="w-10 h-10 rounded-full bg-white flex items-center justify-center text-orange-600 font-bold border border-transparent shadow-sm group-hover:ring-2 group-hover:ring-white transition-all"
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
                class="absolute right-0 mt-2 w-48 bg-white rounded-xl shadow-xl ring-1 ring-black/5 z-50 overflow-hidden animate-in fade-in slide-in-from-top-2 duration-200"
              >
                <div
                  class="px-4 py-3 border-b border-orange-100 bg-orange-50 text-left"
                >
                  <p class="text-xs text-orange-800 m-0">Masuk sebagai</p>
                  <p class="text-sm font-bold text-orange-950 truncate m-0">
                    {user?.username}
                  </p>
                </div>
                <div class="py-1">
                  <a
                    href="/dashboard"
                    class="flex items-center px-4 py-2 text-sm text-orange-900 hover:bg-orange-100 font-medium no-underline"
                  >
                    <svg
                      class="mr-3 w-4 h-4 text-orange-500"
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
                <div class="border-t border-orange-100 py-1 text-left">
                  <button
                    onclick={handleLogout}
                    class="flex w-full items-center px-4 py-2 text-sm text-red-600 hover:bg-red-50 hover:text-red-700 font-medium border-none bg-transparent cursor-pointer"
                  >
                    <svg
                      class="mr-3 w-4 h-4 text-red-500"
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
            class="ml-4 px-4 py-2 rounded-lg bg-white text-orange-900 hover:bg-orange-50 font-semibold text-sm transition-colors border border-transparent shadow-sm no-underline"
            >Masuk Portal</a
          >
        {/if}
      </div>
    </div>
  </div>
</nav>
