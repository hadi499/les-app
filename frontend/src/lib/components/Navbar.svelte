<script lang="ts">
  import { onMount } from "svelte";
  import { page } from "$app/state";

  let isAuthenticated = $state(false);
  let user: { username: string } | null = $state(null);
  let isDropdownOpen = $state(false);
  let isMobileMenuOpen = $state(false);
  let authChecked = $state(false);
  let dropdownRef: HTMLDivElement | undefined = $state();

  const currentPath = $derived(page.url?.pathname || "/");

  onMount(async () => {
    try {
      const res = await fetch(`/me`, {
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
      isDropdownOpen = false;
      authChecked = false; // Sembunyikan UI auth untuk mencegah efek berkedip
      await fetch(`/api/auth/logout`, {
        method: "POST",
        credentials: "include",
      });
      isAuthenticated = false;
      user = null;
      window.location.href = "/"; // redirect to home after logout
    } catch (e) {
      console.error(e);
      authChecked = true;
    }
  }

  function toggleMobileMenu() {
    isMobileMenuOpen = !isMobileMenuOpen;
  }

  function closeMobileMenu() {
    isMobileMenuOpen = false;
  }
</script>

<svelte:window
  onclick={(e) => {
    if (
      isDropdownOpen &&
      dropdownRef &&
      !dropdownRef.contains(e.target as Node)
    ) {
      isDropdownOpen = false;
    }
  }}
/>

<nav
  class="bg-white/80 backdrop-blur-md shadow-md border-b border-slate-100 fixed top-0 w-full z-50 subpixel-antialiased {isMobileMenuOpen
    ? 'bottom-0'
    : ''}"
>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
    <div class="flex justify-between h-16 items-center">
      <div class="shrink-0 flex items-center gap-3">
        <a
          href="/"
          class="w-10 h-10 bg-white border border-blue-300 rounded-xl shadow-sm flex items-center justify-center text-blue-600 font-extrabold text-xl no-underline hover:scale-105 transition-transform"
        >
          LB
        </a>
        <a
          href="/"
          class="hidden sm:block font-['Concert_One'] text-2xl bg-gradient-to-r from-blue-600 to-slate-900 text-transparent bg-clip-text tracking-wide no-underline hover:opacity-80 transition-opacity"
          >Les Balonggarut</a
        >
      </div>
      <div class="hidden lg:ml-6 lg:flex lg:space-x-8 items-center">
        <a
          href="/"
          class="inline-flex items-center px-1 pt-1 border-b-2 text-sm font-semibold transition-colors no-underline {currentPath ===
          '/'
            ? 'border-blue-500 text-blue-700'
            : 'border-transparent text-slate-700 hover:text-blue-600'}"
          >Beranda</a
        >
        <a
          href="/mengetik"
          class="inline-flex items-center px-1 pt-1 border-b-2 text-sm font-semibold transition-colors no-underline {currentPath.startsWith(
            '/mengetik',
          )
            ? 'border-blue-500 text-blue-700'
            : 'border-transparent text-slate-700 hover:text-blue-600'}"
          >Mengetik 10 Jari</a
        >
        <a
          href="/cetak-kode"
          class="inline-flex items-center px-1 pt-1 border-b-2 text-sm font-semibold transition-colors no-underline {currentPath.startsWith(
            '/cetak-kode',
          )
            ? 'border-blue-500 text-blue-700'
            : 'border-transparent text-slate-700 hover:text-blue-600'}"
          >Cetak Kode</a
        >
        <a
          href="/quiz"
          class="inline-flex items-center px-1 pt-1 border-b-2 text-sm font-semibold transition-colors no-underline {currentPath.startsWith(
            '/quiz',
          )
            ? 'border-blue-500 text-blue-700'
            : 'border-transparent text-slate-700 hover:text-blue-600'}">Kuis</a
        >
        <a
          href="/berhitung"
          class="inline-flex items-center px-1 pt-1 border-b-2 text-sm font-semibold transition-colors no-underline {currentPath.startsWith(
            '/berhitung',
          )
            ? 'border-blue-500 text-blue-700'
            : 'border-transparent text-slate-700 hover:text-blue-600'}"
          >Berhitung</a
        >
        <a
          href="/compress-image"
          class="inline-flex items-center px-1 pt-1 border-b-2 text-sm font-semibold transition-colors no-underline {currentPath.startsWith(
            '/compress-image',
          )
            ? 'border-blue-500 text-blue-700'
            : 'border-transparent text-slate-700 hover:text-blue-600'}"
          >Compress Image</a
        >

        {#if !authChecked}
          <div class="ml-4 w-10 h-10"></div>
        {:else if isAuthenticated}
          <!-- Avatar Dropdown -->
          <div class="relative ml-4" bind:this={dropdownRef}>
            <button
              onclick={() => (isDropdownOpen = !isDropdownOpen)}
              class="flex items-center gap-3 focus:outline-none cursor-pointer border-none bg-transparent p-0 group"
            >
              <div
                class="w-10 h-10 rounded-full bg-white flex items-center justify-center text-blue-600 font-bold border border-slate-200 shadow-sm group-hover:ring-2 group-hover:ring-blue-100 transition-all"
              >
                {user?.username ? user.username.charAt(0).toUpperCase() : "U"}
              </div>
            </button>

            {#if isDropdownOpen}
              <!-- Dropdown Menu -->
              <div
                class="absolute right-0 mt-2 w-48 bg-white rounded-xl shadow-xl ring-1 ring-black/5 z-50 overflow-hidden animate-in fade-in slide-in-from-top-2 duration-200"
              >
                <div
                  class="px-4 py-3 border-b border-slate-100 bg-slate-50 text-left"
                >
                  <p class="text-sm font-bold text-slate-900 truncate m-0">
                    {user?.username}
                  </p>
                </div>
                <div class="py-1">
                  <a
                    href="/dashboard"
                    class="flex items-center px-4 py-2 text-sm text-slate-700 hover:bg-slate-50 font-medium no-underline"
                  >
                    <svg
                      class="mr-3 w-4 h-4 text-slate-500"
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
                <div class="border-t border-slate-100 mt-4 pt-2 pb-1 text-left">
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
        {/if}
      </div>

      <!-- Mobile menu button -->
      <div class="flex items-center lg:hidden gap-4">
        <button
          onclick={toggleMobileMenu}
          class="inline-flex items-center justify-center p-2 rounded-md text-slate-800 hover:text-blue-600 hover:bg-blue-50 focus:outline-none border-none bg-transparent cursor-pointer transition-colors"
          aria-expanded="false"
        >
          <span class="sr-only">Buka menu utama</span>
          {#if !isMobileMenuOpen}
            <svg
              class="block h-6 w-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"
              />
            </svg>
          {:else}
            <svg
              class="block h-6 w-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          {/if}
        </button>
      </div>
    </div>
  </div>

  <!-- Mobile menu dropdown -->
  {#if isMobileMenuOpen}
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div
      class="lg:hidden absolute inset-0 top-16 bg-black/10 z-40 animate-in fade-in duration-200"
      onclick={closeMobileMenu}
    ></div>
    <div
      class="lg:hidden relative z-50 bg-white/95 shadow-lg border-t border-slate-100 animate-in slide-in-from-top-2 duration-200 max-h-[calc(100vh-4rem)] overflow-y-auto"
    >
      {#if isAuthenticated}
        <div class="pt-4 pb-3 border-b border-slate-200">
          <div class="flex items-center px-5">
            <div class="shrink-0">
              <div
                class="h-10 w-10 rounded-full bg-blue-50 flex items-center justify-center text-blue-600 font-bold border border-blue-200 shadow-sm"
              >
                {user?.username ? user.username.charAt(0).toUpperCase() : "U"}
              </div>
            </div>
            <div class="ml-3">
              <div class="text-base font-medium leading-none text-slate-900">
                {user?.username}
              </div>
            </div>
          </div>
          <div class="mt-6 px-2 space-y-1">
            <a
              href="/dashboard"
              onclick={closeMobileMenu}
              class="block px-3 py-2 rounded-md text-base font-medium text-slate-700 hover:text-slate-900 hover:bg-slate-50 no-underline"
              >Dashboard</a
            >
          </div>
        </div>
      {/if}
      <div class="px-2 pt-3 pb-3 space-y-1">
        <a
          href="/"
          onclick={closeMobileMenu}
          class="block px-3 py-2 rounded-md text-base font-medium no-underline {currentPath ===
          '/'
            ? 'bg-blue-50 text-blue-700'
            : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'}"
          >Beranda</a
        >
        <a
          href="/mengetik"
          onclick={closeMobileMenu}
          class="block px-3 py-2 rounded-md text-base font-medium no-underline {currentPath.startsWith(
            '/mengetik',
          )
            ? 'bg-blue-50 text-blue-700'
            : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'}"
          >Mengetik 10 Jari</a
        >
        <a
          href="/cetak-kode"
          onclick={closeMobileMenu}
          class="block px-3 py-2 rounded-md text-base font-medium no-underline {currentPath.startsWith(
            '/cetak-kode',
          )
            ? 'bg-blue-50 text-blue-700'
            : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'}"
          >Cetak Kode</a
        >
        <a
          href="/quiz"
          onclick={closeMobileMenu}
          class="block px-3 py-2 rounded-md text-base font-medium no-underline {currentPath.startsWith(
            '/quiz',
          )
            ? 'bg-blue-50 text-blue-700'
            : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'}">Kuis</a
        >
        <a
          href="/berhitung"
          onclick={closeMobileMenu}
          class="block px-3 py-2 rounded-md text-base font-medium no-underline {currentPath.startsWith(
            '/berhitung',
          )
            ? 'bg-blue-50 text-blue-700'
            : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'}"
          >Berhitung</a
        >
        <a
          href="/compress-image"
          onclick={closeMobileMenu}
          class="block px-3 py-2 rounded-md text-base font-medium no-underline {currentPath.startsWith(
            '/compress-image',
          )
            ? 'bg-blue-50 text-blue-700'
            : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'}"
          >Compress Image</a
        >
      </div>
      {#if isAuthenticated}
        <div class="px-2 pt-2 pb-3 border-t border-slate-200 mt-6">
          <button
            onclick={handleLogout}
            class="block w-full text-left px-3 py-2 rounded-md text-base font-medium text-red-600 hover:text-red-700 hover:bg-red-50 border-none bg-transparent cursor-pointer"
            >Keluar Akun</button
          >
        </div>
      {/if}
    </div>
  {/if}
</nav>
