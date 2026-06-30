<script lang="ts">
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";

  import type { Snippet } from "svelte";
  let { children }: { children: Snippet } = $props();

  // Data dinamis dari backend
  type User = { username: string; role: string };
  let user: User = $state({ username: "Loading...", role: "..." });
  let isMobileMenuOpen = $state(false);
  let isDesktopSidebarOpen = $state(true);
  let isLoading = $state(true);

  onMount(async () => {
    try {
      const res = await fetch(`/me`, {
        credentials: "include",
      });
      if (!res.ok) {
        goto("/login");
        return;
      }
      const data = await res.json();
      if (!data.authenticated) {
        goto("/login");
        return;
      }
      user = data.user;
    } catch (e) {
      goto("/login");
    } finally {
      isLoading = false;
    }
  });

  async function handleLogout() {
    try {
      await fetch(`/api/auth/logout`, {
        method: "POST",
        credentials: "include",
      });
      goto("/login");
    } catch (e) {
      console.error(e);
      goto("/login");
    }
  }
</script>

{#if isLoading}
  <!-- Layar Loading Sederhana -->
  <div class="min-h-screen bg-slate-100 flex items-center justify-center">
    <div class="flex flex-col items-center gap-4">
      <div
        class="w-12 h-12 border-4 border-slate-200 border-t-blue-500 rounded-full animate-spin"
      ></div>
      <p class="text-slate-600 font-semibold">Memuat portal...</p>
    </div>
  </div>
{:else}
  <div
    class="min-h-screen bg-slate-50 flex selection:bg-slate-200 selection:text-slate-800 font-sans text-slate-900"
    style="--sidebar-width: {isDesktopSidebarOpen ? '16rem' : '0px'};"
  >
    <!-- Sidebar (Desktop) -->
    <aside
      class="bg-slate-400/10 backdrop-blur-md border-r border-slate-300 hidden md:flex flex-col shrink-0 sticky top-0 h-screen print:hidden shadow-lg shadow-slate-800/5 transition-all duration-300 z-30 {isDesktopSidebarOpen
        ? 'w-64 opacity-100'
        : 'w-0 opacity-0 overflow-hidden border-transparent'}"
    >
      <!-- Logo -->
      <div class="h-16 flex items-center px-6 border-b border-slate-300">
        <div class="group flex items-center gap-3">
          <div
            class="w-8 h-8 bg-white border border-transparent rounded-lg flex items-center justify-center text-blue-600 font-extrabold text-sm shadow-md"
          >
            LB
          </div>
          <span
            class="font-extrabold text-lg text-slate-900 tracking-tight drop-shadow-sm"
            >{user?.role === "teacher" ? "Portal Guru" : "Portal Siswa"}</span
          >
        </div>
      </div>

      <!-- Nav Links -->
      <nav class="flex-1 overflow-y-auto py-6 px-4 space-y-1">
        <a
          href="/dashboard"
          class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline {page
            .url.pathname === '/dashboard'
            ? 'bg-white/80 text-blue-700 font-medium shadow-sm shadow-slate-800/5 border border-slate-300'
            : 'text-slate-700 hover:bg-white/50 hover:text-slate-900 border border-transparent'}"
        >
          <div
            class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page
              .url.pathname === '/dashboard'
              ? 'bg-blue-600 text-white shadow-md shadow-blue-500/20'
              : 'bg-slate-100 text-slate-500 group-hover:bg-slate-200 group-hover:text-slate-700'}"
          >
            <svg
              class="w-5 h-5"
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
          </div>
          Dashboard
        </a>

        <a
          href="/dashboard/todos"
          class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline {page.url.pathname.includes(
            '/todos',
          )
            ? 'bg-white/80 text-blue-700 font-medium shadow-sm shadow-slate-800/5 border border-slate-300'
            : 'text-slate-700 hover:bg-white/50 hover:text-slate-900 border border-transparent'}"
        >
          <div
            class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
              '/todos',
            )
              ? 'bg-blue-600 text-white shadow-md shadow-blue-500/20'
              : 'bg-slate-100 text-slate-500 group-hover:bg-slate-200 group-hover:text-slate-700'}"
          >
            <svg
              class="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              ><path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4"
              ></path></svg
            >
          </div>
          Todolist
        </a>

        {#if user?.role === "teacher"}
          <a
            href="/dashboard/notes"
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline {page.url.pathname.includes(
              '/notes',
            )
              ? 'bg-white/80 text-blue-700 font-medium shadow-sm shadow-slate-800/5 border border-slate-300'
              : 'text-slate-700 hover:bg-white/50 hover:text-slate-900 border border-transparent'}"
          >
            <div
              class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
                '/notes',
              )
                ? 'bg-blue-600 text-white shadow-md shadow-blue-500/20'
                : 'bg-slate-100 text-slate-500 group-hover:bg-slate-200 group-hover:text-slate-700'}"
            >
              <svg
                class="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                ></path></svg
              >
            </div>
            Catatan
          </a>
          <a
            href="/dashboard/card-memory"
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline {page.url.pathname.includes(
              '/card-memory',
            )
              ? 'bg-white/80 text-blue-700 font-medium shadow-sm shadow-slate-800/5 border border-slate-300'
              : 'text-slate-700 hover:bg-white/50 hover:text-slate-900 border border-transparent'}"
          >
            <div
              class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
                '/card-memory',
              )
                ? 'bg-blue-600 text-white shadow-md shadow-blue-500/20'
                : 'bg-slate-100 text-slate-500 group-hover:bg-slate-200 group-hover:text-slate-700'}"
            >
              <svg
                class="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
                ></path></svg
              >
            </div>
            Card Memory
          </a>
        {/if}

        <a
          href="/dashboard/exams"
          class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline {page.url.pathname.includes(
            '/exams',
          )
            ? 'bg-white/80 text-blue-700 font-medium shadow-sm shadow-slate-800/5 border border-slate-300'
            : 'text-slate-700 hover:bg-white/50 hover:text-slate-900 border border-transparent'}"
        >
          <div
            class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
              '/exams',
            )
              ? 'bg-blue-600 text-white shadow-md shadow-blue-500/20'
              : 'bg-slate-100 text-slate-500 group-hover:bg-slate-200 group-hover:text-slate-700'}"
          >
            <svg
              class="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              ><path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
              ></path></svg
            >
          </div>
          Nilai Harian
        </a>

        {#if user?.role === "teacher"}
          <a
            href="/dashboard/subjects"
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline {page.url.pathname.includes(
              '/subjects',
            )
              ? 'bg-white/80 text-blue-700 font-medium shadow-sm shadow-slate-800/5 border border-slate-300'
              : 'text-slate-700 hover:bg-white/50 hover:text-slate-900 border border-transparent'}"
          >
            <div
              class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
                '/subjects',
              )
                ? 'bg-blue-600 text-white shadow-md shadow-blue-500/20'
                : 'bg-slate-100 text-slate-500 group-hover:bg-slate-200 group-hover:text-slate-700'}"
            >
              <svg
                class="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"
                ></path></svg
              >
            </div>
            Mata Pelajaran
          </a>
        {/if}

        <a
          href="/dashboard/quizzes"
          class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline {page.url.pathname.includes(
            '/quizzes',
          )
            ? 'bg-white/80 text-blue-700 font-medium shadow-sm shadow-slate-800/5 border border-slate-300'
            : 'text-slate-700 hover:bg-white/50 hover:text-slate-900 border border-transparent'}"
        >
          <div
            class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
              '/quizzes',
            )
              ? 'bg-blue-600 text-white shadow-md shadow-blue-500/20'
              : 'bg-slate-100 text-slate-500 group-hover:bg-slate-200 group-hover:text-slate-700'}"
          >
            <svg
              class="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              ><path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              ></path></svg
            >
          </div>
          Kuis & Nilai
        </a>

        {#if user?.role === "teacher"}
          <a
            href="/dashboard/absen"
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline {page.url.pathname.includes(
              '/absen',
            )
              ? 'bg-white/80 text-blue-700 font-medium shadow-sm shadow-slate-800/5 border border-slate-300'
              : 'text-slate-700 hover:bg-white/50 hover:text-slate-900 border border-transparent'}"
          >
            <div
              class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
                '/absen',
              )
                ? 'bg-blue-600 text-white shadow-md shadow-blue-500/20'
                : 'bg-slate-100 text-slate-500 group-hover:bg-slate-200 group-hover:text-slate-700'}"
            >
              <svg
                class="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                ></path></svg
              >
            </div>
            Absensi
          </a>
          <a
            href="/dashboard/typing-monitoring"
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline {page.url.pathname.includes(
              '/typing-monitoring',
            )
              ? 'bg-white/80 text-blue-700 font-medium shadow-sm shadow-slate-800/5 border border-slate-300'
              : 'text-slate-700 hover:bg-white/50 hover:text-slate-900 border border-transparent'}"
          >
            <div
              class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
                '/typing-monitoring',
              )
                ? 'bg-blue-600 text-white shadow-md shadow-blue-500/20'
                : 'bg-slate-100 text-slate-500 group-hover:bg-slate-200 group-hover:text-slate-700'}"
            >
              <svg
                class="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"
                ></path></svg
              >
            </div>
            Pantau Mengetik
          </a>
          <a
            href="/dashboard/users"
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline {page.url.pathname.includes(
              '/users',
            )
              ? 'bg-white/80 text-blue-700 font-medium shadow-sm shadow-slate-800/5 border border-slate-300'
              : 'text-slate-700 hover:bg-white/50 hover:text-slate-900 border border-transparent'}"
          >
            <div
              class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
                '/users',
              )
                ? 'bg-blue-600 text-white shadow-md shadow-blue-500/20'
                : 'bg-slate-100 text-slate-500 group-hover:bg-slate-200 group-hover:text-slate-700'}"
            >
              <svg
                class="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"
                ></path></svg
              >
            </div>
            Manajemen Users
          </a>
        {/if}

        <div class="pt-4 mt-2 border-t border-slate-200">
          <a
            href="/"
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors text-slate-700 hover:bg-white/50 hover:text-slate-900 no-underline"
          >
            <div
              class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all bg-slate-100 text-slate-500 group-hover:bg-slate-200 group-hover:text-slate-700"
            >
              <svg
                class="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9"
                ></path></svg
              >
            </div>
            Halaman Utama
          </a>
        </div>
      </nav>

      <!-- User & Logout -->
      <div class="p-4 border-t border-slate-300">
        <div class="group flex items-center gap-3 mb-4 px-2">
          <div
            class="w-10 h-10 rounded-full bg-white flex items-center justify-center text-blue-600 font-bold border border-transparent shadow-sm flex-shrink-0"
          >
            {user.username ? user.username.charAt(0).toUpperCase() : "U"}
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-bold text-slate-900 truncate m-0">
              {user.username}
            </p>
            <p
              class="text-xs font-medium text-slate-600 truncate capitalize m-0 mt-0.5"
            >
              {user.role}
            </p>
          </div>
        </div>
        <button
          onclick={handleLogout}
          class="w-full flex items-center justify-center gap-2 px-4 py-2.5 text-sm font-medium text-red-600 bg-white border border-transparent hover:bg-red-50 hover:text-red-700 rounded-xl transition-all shadow-sm cursor-pointer"
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
          Keluar Akun
        </button>
      </div>
    </aside>

    <!-- Desktop Sidebar Toggle Button -->
    <button
      onclick={() => (isDesktopSidebarOpen = !isDesktopSidebarOpen)}
      class="hidden md:flex fixed top-4 z-50 p-1.5 bg-white border border-slate-200 shadow-sm rounded-lg text-slate-500 hover:bg-slate-50 hover:text-blue-600 transition-all duration-300 cursor-pointer"
      style="left: {isDesktopSidebarOpen ? '15.25rem' : '1rem'};"
      title={isDesktopSidebarOpen ? "Sembunyikan Sidebar" : "Tampilkan Sidebar"}
    >
      <svg
        class="w-5 h-5"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        {#if isDesktopSidebarOpen}
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M11 19l-7-7 7-7m8 14l-7-7 7-7"
          />
        {:else}
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M13 5l7 7-7 7M5 5l7 7-7 7"
          />
        {/if}
      </svg>
    </button>

    <!-- Mobile Header -->
    <div
      class="md:hidden fixed top-0 inset-x-0 h-16 bg-slate-400/10 backdrop-blur-md border-b border-slate-200 z-50 flex items-center justify-between px-4 shadow-sm print:hidden"
    >
      <div class="flex items-center gap-2">
        <div
          class="w-8 h-8 bg-white border border-transparent rounded-lg flex items-center justify-center text-blue-600 font-extrabold text-sm shadow-md"
        >
          LB
        </div>
        <span
          class="font-extrabold text-lg text-slate-900 tracking-tight drop-shadow-sm"
          >{user?.role === "teacher" ? "Portal Guru" : "Portal Siswa"}</span
        >
      </div>
      <button
        aria-label="Buka menu navigasi"
        class="text-slate-600 focus:outline-none p-2 bg-white/50 rounded-lg border border-transparent hover:bg-white/80 transition-colors {isMobileMenuOpen
          ? 'hidden'
          : 'block'}"
        onclick={() => (isMobileMenuOpen = !isMobileMenuOpen)}
      >
        <svg
          class="w-5 h-5"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          ><path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M4 6h16M4 12h16M4 18h16"
          ></path></svg
        >
      </button>
    </div>

    {#if isMobileMenuOpen}
      <div
        class="md:hidden fixed inset-0 z-40 bg-black/60 backdrop-blur-sm"
        role="button"
        tabindex="-1"
        aria-label="Tutup menu"
        onclick={() => (isMobileMenuOpen = false)}
        onkeydown={(e) => {
          if (e.key === "Enter" || e.key === " ") isMobileMenuOpen = false;
        }}
      ></div>
      <aside
        class="md:hidden fixed inset-y-0 left-0 w-64 bg-slate-100 border-r border-slate-200 shadow-2xl z-50 flex flex-col transform transition-transform {isMobileMenuOpen
          ? 'translate-x-0'
          : '-translate-x-full'} print:hidden"
      >
        <!-- Copy of sidebar content for mobile -->
        <div
          class="h-16 flex items-center justify-between px-6 border-b border-slate-200 bg-slate-100"
        >
          <span
            class="font-extrabold text-lg text-slate-900 tracking-tight drop-shadow-sm"
            >Menu Utama</span
          >
          <button
            aria-label="Tutup menu"
            onclick={() => (isMobileMenuOpen = false)}
            class="text-slate-400 hover:text-white border-none bg-transparent cursor-pointer"
            ><svg
              class="w-6 h-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              ><path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M6 18L18 6M6 6l12 12"
              ></path></svg
            ></button
          >
        </div>
        <nav class="flex-1 overflow-y-auto py-6 px-4 space-y-1 bg-slate-100">
          <a
            href="/dashboard"
            onclick={() => (isMobileMenuOpen = false)}
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-sm transition-colors no-underline {page
              .url.pathname === '/dashboard'
              ? 'bg-white/80 text-blue-700 font-medium border border-slate-300'
              : 'text-slate-700 hover:text-slate-900 hover:bg-white/50 border border-transparent'}"
          >
            <div
              class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page
                .url.pathname === '/dashboard'
                ? 'bg-blue-600 text-white shadow-md shadow-blue-500/20'
                : 'bg-slate-100 text-slate-500 group-hover:bg-slate-200 group-hover:text-slate-700'}"
            >
              <svg
                class="w-5 h-5"
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
            </div>
            Dashboard
          </a>

          <a
            href="/dashboard/todos"
            onclick={() => (isMobileMenuOpen = false)}
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-sm transition-colors no-underline {page.url.pathname.includes(
              '/todos',
            )
              ? 'bg-white/80 text-blue-700 font-medium border border-slate-300'
              : 'text-slate-700 hover:text-slate-900 hover:bg-white/50 border border-transparent'}"
          >
            <div
              class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
                '/todos',
              )
                ? 'bg-blue-600 text-white shadow-md shadow-blue-500/20'
                : 'bg-slate-100 text-slate-500 group-hover:bg-slate-200 group-hover:text-slate-700'}"
            >
              <svg
                class="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4"
                ></path></svg
              >
            </div>
            Todolist
          </a>

          {#if user?.role === "teacher"}
            <a
              href="/dashboard/notes"
              onclick={() => (isMobileMenuOpen = false)}
              class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-sm transition-colors no-underline {page.url.pathname.includes(
                '/notes',
              )
                ? 'bg-white/80 text-blue-700 font-medium border border-slate-300'
                : 'text-slate-700 hover:text-slate-900 hover:bg-white/50 border border-transparent'}"
            >
              <div
                class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
                  '/notes',
                )
                  ? 'bg-blue-600 text-white shadow-md shadow-blue-500/20'
                  : 'bg-slate-100 text-slate-500 group-hover:bg-slate-200 group-hover:text-slate-700'}"
              >
                <svg
                  class="w-5 h-5"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  ><path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                  ></path></svg
                >
              </div>
              Catatan
            </a>
            <a
              href="/dashboard/card-memory"
              onclick={() => (isMobileMenuOpen = false)}
              class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-sm transition-colors no-underline {page.url.pathname.includes(
                '/card-memory',
              )
                ? 'bg-white/80 text-blue-700 font-medium border border-slate-300'
                : 'text-slate-700 hover:text-slate-900 hover:bg-white/50 border border-transparent'}"
            >
              <div
                class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
                  '/card-memory',
                )
                  ? 'bg-blue-600 text-white shadow-md shadow-blue-500/20'
                  : 'bg-slate-100 text-slate-500 group-hover:bg-slate-200 group-hover:text-slate-700'}"
              >
                <svg
                  class="w-5 h-5"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  ><path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
                  ></path></svg
                >
              </div>
              Card Memory
            </a>
          {/if}

          <a
            href="/dashboard/exams"
            onclick={() => (isMobileMenuOpen = false)}
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-sm transition-colors no-underline {page.url.pathname.includes(
              '/exams',
            )
              ? 'bg-white/80 text-blue-700 font-medium border border-slate-300'
              : 'text-slate-700 hover:text-slate-900 hover:bg-white/50 border border-transparent'}"
          >
            <div
              class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
                '/exams',
              )
                ? 'bg-blue-600 text-white shadow-md shadow-blue-500/20'
                : 'bg-slate-100 text-slate-500 group-hover:bg-slate-200 group-hover:text-slate-700'}"
            >
              <svg
                class="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                ></path></svg
              >
            </div>
            Nilai Harian
          </a>

          {#if user?.role === "teacher"}
            <a
              href="/dashboard/subjects"
              onclick={() => (isMobileMenuOpen = false)}
              class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-sm transition-colors no-underline {page.url.pathname.includes(
                '/subjects',
              )
                ? 'bg-white/80 text-blue-700 font-medium border border-slate-300'
                : 'text-slate-700 hover:text-slate-900 hover:bg-white/50 border border-transparent'}"
            >
              <div
                class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
                  '/subjects',
                )
                  ? 'bg-blue-600 text-white shadow-md shadow-blue-500/20'
                  : 'bg-slate-100 text-slate-500 group-hover:bg-slate-200 group-hover:text-slate-700'}"
              >
                <svg
                  class="w-5 h-5"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  ><path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"
                  ></path></svg
                >
              </div>
              Mata Pelajaran
            </a>
          {/if}

          <a
            href="/dashboard/quizzes"
            onclick={() => (isMobileMenuOpen = false)}
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-sm transition-colors no-underline {page.url.pathname.includes(
              '/quizzes',
            )
              ? 'bg-white/80 text-blue-700 font-medium border border-slate-300'
              : 'text-slate-700 hover:text-slate-900 hover:bg-white/50 border border-transparent'}"
          >
            <div
              class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
                '/quizzes',
              )
                ? 'bg-blue-600 text-white shadow-md shadow-blue-500/20'
                : 'bg-slate-100 text-slate-500 group-hover:bg-slate-200 group-hover:text-slate-700'}"
            >
              <svg
                class="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                ></path></svg
              >
            </div>
            Kuis & Nilai
          </a>

          {#if user?.role === "teacher"}
            <a
              href="/dashboard/absen"
              onclick={() => (isMobileMenuOpen = false)}
              class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-sm transition-colors no-underline {page.url.pathname.includes(
                '/absen',
              )
                ? 'bg-white/80 text-blue-700 font-medium border border-slate-300'
                : 'text-slate-700 hover:text-slate-900 hover:bg-white/50 border border-transparent'}"
            >
              <div
                class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
                  '/absen',
                )
                  ? 'bg-blue-600 text-white shadow-md shadow-blue-500/20'
                  : 'bg-slate-100 text-slate-500 group-hover:bg-slate-200 group-hover:text-slate-700'}"
              >
                <svg
                  class="w-5 h-5"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  ><path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                  ></path></svg
                >
              </div>
              Absensi
            </a>
            <a
              href="/dashboard/typing-monitoring"
              onclick={() => (isMobileMenuOpen = false)}
              class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-sm transition-colors no-underline {page.url.pathname.includes(
                '/typing-monitoring',
              )
                ? 'bg-white/80 text-blue-700 font-medium border border-slate-300'
                : 'text-slate-700 hover:text-slate-900 hover:bg-white/50 border border-transparent'}"
            >
              <div
                class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
                  '/typing-monitoring',
                )
                  ? 'bg-blue-600 text-white shadow-md shadow-blue-500/20'
                  : 'bg-slate-100 text-slate-500 group-hover:bg-slate-200 group-hover:text-slate-700'}"
              >
                <svg
                  class="w-5 h-5"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  ><path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"
                  ></path></svg
                >
              </div>
              Pantau Mengetik
            </a>
            <a
              href="/dashboard/users"
              onclick={() => (isMobileMenuOpen = false)}
              class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-sm transition-colors no-underline {page.url.pathname.includes(
                '/users',
              )
                ? 'bg-white/80 text-blue-700 font-medium border border-slate-300'
                : 'text-slate-700 hover:text-slate-900 hover:bg-white/50 border border-transparent'}"
            >
              <div
                class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
                  '/users',
                )
                  ? 'bg-blue-600 text-white shadow-md shadow-blue-500/20'
                  : 'bg-slate-100 text-slate-500 group-hover:bg-slate-200 group-hover:text-slate-700'}"
              >
                <svg
                  class="w-5 h-5"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  ><path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"
                  ></path></svg
                >
              </div>
              Manajemen Users
            </a>
          {/if}

          <div class="pt-4 mt-2 border-t border-slate-200">
            <a
              href="/"
              onclick={() => (isMobileMenuOpen = false)}
              class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-sm transition-colors text-slate-700 hover:bg-white/50 hover:text-slate-900 no-underline"
            >
              <svg
                class="w-5 h-5 text-slate-400"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9"
                ></path></svg
              >
              Halaman Utama
            </a>
          </div>
        </nav>

        <!-- User & Logout (Mobile) -->
        <div class="p-4 border-t border-slate-200 bg-slate-100">
          <div class="group flex items-center gap-3 mb-4 px-2">
            <div
              class="w-10 h-10 rounded-full bg-white flex items-center justify-center text-blue-600 font-bold border border-transparent shadow-sm shrink-0"
            >
              {user.username ? user.username.charAt(0).toUpperCase() : "U"}
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-bold text-slate-900 truncate m-0">
                {user.username}
              </p>
              <p
                class="text-xs font-medium text-slate-600 truncate capitalize m-0 mt-0.5"
              >
                {user.role}
              </p>
            </div>
          </div>
          <button
            onclick={handleLogout}
            class="w-full flex items-center justify-center gap-2 px-4 py-2.5 text-sm font-medium text-red-600 bg-white border border-transparent hover:bg-red-50 hover:text-red-700 rounded-xl transition-all shadow-sm cursor-pointer"
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
            Keluar Akun
          </button>
        </div>
      </aside>
    {/if}

    <!-- Main Content Area -->
    <main
      class="flex-1 flex flex-col min-w-0 md:pt-0 pt-16 h-screen overflow-y-auto print:pt-0 print:h-auto print:overflow-visible print:block bg-transparent"
    >
      <div
        class="flex-1 p-6 md:p-8 {isDesktopSidebarOpen
          ? ''
          : 'md:pl-16'} max-w-6xl mx-auto w-full print:p-0 print:m-0"
      >
        {@render children()}
      </div>
    </main>
  </div>
{/if}
