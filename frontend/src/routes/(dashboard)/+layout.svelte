<script lang="ts">
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import { onMount, onDestroy } from "svelte";
  import { chatStore } from "$lib/stores/chatStore";

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
      
      // Setup chat
      chatStore.setMyUserId(user.id);
      try {
        const unreadRes = await fetch("/api/chat/unread-count", { credentials: "include" });
        if (unreadRes.ok) {
          const unreadData = await unreadRes.json();
          chatStore.setUnreadCount(unreadData.count || 0);
        }
      } catch (e) {
        console.error("Failed to fetch unread count", e);
      }
      chatStore.connect();
      
    } catch (e) {
      goto("/login");
    } finally {
      isLoading = false;
    }
  });

  onDestroy(() => {
    chatStore.disconnect();
  });

  $effect(() => {
    if (!isLoading && user.role !== "teacher") {
      const path = page.url.pathname;
      const teacherOnlyRoutes = [
        "/dashboard/notes",
        "/dashboard/card-memory",
        "/dashboard/subjects",
        "/dashboard/typing-monitoring",
        "/dashboard/users",
      ];
      if (teacherOnlyRoutes.some((r) => path.startsWith(r))) {
        goto("/dashboard", { replaceState: true });
      }
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

  let showPasswordModal = $state(false);
  let oldPassword = $state("");
  let newPassword = $state("");
  let confirmPassword = $state("");
  let passwordError = $state("");
  let passwordSuccess = $state("");
  let isSubmittingPassword = $state(false);

  let showOldPassword = $state(false);
  let showNewPassword = $state(false);
  let showConfirmPassword = $state(false);

  let showUserDropdown = $state(false);
  let showMobileUserDropdown = $state(false);

  function resetPasswordForm() {
    oldPassword = "";
    newPassword = "";
    confirmPassword = "";
    passwordError = "";
    passwordSuccess = "";
    showOldPassword = false;
    showNewPassword = false;
    showConfirmPassword = false;
  }

  async function handlePasswordChange(e: Event) {
    e.preventDefault();
    if (!oldPassword || !newPassword || !confirmPassword) {
      passwordError = "Semua kolom password harus diisi.";
      return;
    }
    if (newPassword !== confirmPassword) {
      passwordError = "Password baru dan konfirmasi tidak cocok.";
      return;
    }
    if (newPassword.length < 6) {
      passwordError = "Password minimal 6 karakter.";
      return;
    }

    isSubmittingPassword = true;
    passwordError = "";
    passwordSuccess = "";

    try {
      const res = await fetch("/api/auth/change-password", {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          old_password: oldPassword,
          new_password: newPassword,
        }),
      });

      const data = await res.json();
      if (!res.ok) {
        passwordError = data.error || "Gagal mengubah password";
      } else {
        passwordSuccess = "Password berhasil diubah!";
        setTimeout(() => {
          showPasswordModal = false;
          resetPasswordForm();
        }, 1500);
      }
    } catch (err) {
      passwordError = "Terjadi kesalahan koneksi";
    } finally {
      isSubmittingPassword = false;
    }
  }
</script>

<svelte:window
  onclick={() => {
    showUserDropdown = false;
    showMobileUserDropdown = false;
  }}
/>

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
    class="min-h-screen w-full overflow-x-hidden bg-slate-50 flex selection:bg-slate-200 selection:text-slate-800 font-sans text-slate-900"
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
          href="/dashboard/chat"
          class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline {page.url.pathname.includes(
            '/chat',
          )
            ? 'bg-white/80 text-blue-700 font-medium shadow-sm shadow-slate-800/5 border border-slate-300'
            : 'text-slate-700 hover:bg-white/50 hover:text-slate-900 border border-transparent'}"
        >
          <div
            class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
              '/chat',
            )
              ? 'bg-blue-600 text-white shadow-md shadow-blue-500/20'
              : 'bg-slate-100 text-slate-500 group-hover:bg-slate-200 group-hover:text-slate-700'}"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"></path></svg>
          </div>
          <span class="flex-1">Pesan</span>
          {#if $chatStore.unreadCount > 0}
            <span class="bg-red-500 text-white text-[10px] font-bold px-2 py-0.5 rounded-full min-w-[20px] text-center">
              {$chatStore.unreadCount}
            </span>
          {/if}
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

        <a
          href="/dashboard/writing-progress"
          class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline {page.url.pathname.includes(
            '/writing-progress',
          )
            ? 'bg-white/80 text-blue-700 font-medium shadow-sm shadow-slate-800/5 border border-slate-300'
            : 'text-slate-700 hover:bg-white/50 hover:text-slate-900 border border-transparent'}"
        >
          <div
            class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
              '/writing-progress',
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
                d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"
              ></path></svg
            >
          </div>
          Menulis TK/PAUD
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

        {#if user?.role === "teacher"}
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
      <div class="p-4 border-t border-slate-300 relative">
        <button
          onclick={(e) => {
            e.stopPropagation();
            showUserDropdown = !showUserDropdown;
          }}
          class="w-full group flex items-center justify-between gap-3 px-2 py-2 rounded-xl hover:bg-white/50 transition-colors cursor-pointer border border-transparent"
        >
          <div class="flex items-center gap-3 min-w-0">
            <div
              class="w-10 h-10 rounded-full bg-white flex items-center justify-center text-blue-600 font-bold border border-slate-200 shadow-sm shrink-0"
            >
              {user.username ? user.username.charAt(0).toUpperCase() : "U"}
            </div>
            <div class="flex-1 min-w-0 text-left">
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
          <svg
            class="w-5 h-5 text-slate-400 group-hover:text-slate-600 transition-transform {showUserDropdown
              ? 'rotate-180'
              : ''}"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            ><path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M5 15l7-7 7 7"
            ></path></svg
          >
        </button>

        {#if showUserDropdown}
          <!-- svelte-ignore a11y_click_events_have_key_events -->
          <!-- svelte-ignore a11y_no_static_element_interactions -->
          <div
            class="absolute bottom-[calc(100%-1rem)] left-4 right-4 mb-2 bg-white rounded-xl shadow-lg border border-slate-200 overflow-hidden animate-in fade-in slide-in-from-bottom-2 duration-200 z-50"
            onclick={(e) => e.stopPropagation()}
          >
            <button
              onclick={() => {
                showUserDropdown = false;
                showPasswordModal = true;
                resetPasswordForm();
              }}
              class="w-full flex items-center gap-3 px-4 py-3 text-sm font-medium text-slate-700 hover:bg-slate-50 hover:text-slate-900 transition-colors cursor-pointer text-left border-b border-slate-100"
            >
              <svg
                class="w-4 h-4 text-slate-400"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"
                ></path></svg
              >
              Ganti Password
            </button>
            <button
              onclick={handleLogout}
              class="w-full flex items-center gap-3 px-4 py-3 text-sm font-medium text-red-600 hover:bg-red-50 hover:text-red-700 transition-colors cursor-pointer text-left"
            >
              <svg
                class="w-4 h-4 text-red-500"
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
        {/if}
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
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors no-underline {page
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
            href="/dashboard/chat"
            onclick={() => (isMobileMenuOpen = false)}
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors no-underline {page.url.pathname.includes(
              '/chat',
            )
              ? 'bg-white/80 text-blue-700 font-medium border border-slate-300'
              : 'text-slate-700 hover:text-slate-900 hover:bg-white/50 border border-transparent'}"
          >
            <div
              class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
                '/chat',
              )
                ? 'bg-blue-600 text-white shadow-md shadow-blue-500/20'
                : 'bg-slate-100 text-slate-500 group-hover:bg-slate-200 group-hover:text-slate-700'}"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"></path></svg>
            </div>
            <span class="flex-1">Pesan</span>
            {#if $chatStore.unreadCount > 0}
              <span class="bg-red-500 text-white text-[10px] font-bold px-2 py-0.5 rounded-full min-w-[20px] text-center">
                {$chatStore.unreadCount}
              </span>
            {/if}
          </a>

          <a
            href="/dashboard/todos"
            onclick={() => (isMobileMenuOpen = false)}
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors no-underline {page.url.pathname.includes(
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

          <a
            href="/dashboard/writing-progress"
            onclick={() => (isMobileMenuOpen = false)}
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors no-underline {page.url.pathname.includes(
              '/writing-progress',
            )
              ? 'bg-white/80 text-blue-700 font-medium border border-slate-300'
              : 'text-slate-700 hover:text-slate-900 hover:bg-white/50 border border-transparent'}"
          >
            <div
              class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
                '/writing-progress',
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
                  d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"
                ></path></svg
              >
            </div>
            Menulis TK/PAUD
          </a>

          {#if user?.role === "teacher"}
            <a
              href="/dashboard/notes"
              onclick={() => (isMobileMenuOpen = false)}
              class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors no-underline {page.url.pathname.includes(
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
              class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors no-underline {page.url.pathname.includes(
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
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors no-underline {page.url.pathname.includes(
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
              class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors no-underline {page.url.pathname.includes(
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
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors no-underline {page.url.pathname.includes(
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


          <a
            href="/dashboard/absen"
            onclick={() => (isMobileMenuOpen = false)}
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors no-underline {page.url.pathname.includes(
              '/absen',
            )
              ? 'bg-white/80 text-blue-700 font-medium border border-slate-300'
              : 'text-slate-700 hover:text-slate-900 hover:bg-white/50 border border-transparent'}"
          >
            <div
              class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
                '/absen',
              )
                ? 'bg-blue-600 text-white shadow-md'
                : 'bg-slate-100 text-slate-500 group-hover:bg-slate-200 group-hover:text-slate-700'}"
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
                  d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                ></path></svg
              >
            </div>
            Absensi
          </a>

          {#if user?.role === "teacher"}
            <a
              href="/dashboard/typing-monitoring"
              onclick={() => (isMobileMenuOpen = false)}
              class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors no-underline {page.url.pathname.includes(
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
              class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors no-underline {page.url.pathname.includes(
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
              class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors text-slate-700 hover:bg-white/50 hover:text-slate-900 no-underline"
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
        <div class="p-4 border-t border-slate-200 bg-slate-100 relative">
          <button
            onclick={(e) => {
              e.stopPropagation();
              showMobileUserDropdown = !showMobileUserDropdown;
            }}
            class="w-full group flex items-center justify-between gap-3 px-2 py-2 rounded-xl hover:bg-white/50 transition-colors cursor-pointer border border-transparent"
          >
            <div class="flex items-center gap-3 min-w-0">
              <div
                class="w-10 h-10 rounded-full bg-white flex items-center justify-center text-blue-600 font-bold border border-slate-200 shadow-sm shrink-0"
              >
                {user.username ? user.username.charAt(0).toUpperCase() : "U"}
              </div>
              <div class="flex-1 min-w-0 text-left">
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
            <svg
              class="w-5 h-5 text-slate-400 group-hover:text-slate-600 transition-transform {showMobileUserDropdown
                ? 'rotate-180'
                : ''}"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              ><path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M5 15l7-7 7 7"
              ></path></svg
            >
          </button>

          {#if showMobileUserDropdown}
            <!-- svelte-ignore a11y_click_events_have_key_events -->
            <!-- svelte-ignore a11y_no_static_element_interactions -->
            <div
              class="absolute bottom-[calc(100%-1rem)] left-4 right-4 mb-2 bg-white rounded-xl shadow-lg border border-slate-200 overflow-hidden animate-in fade-in slide-in-from-bottom-2 duration-200 z-50"
              onclick={(e) => e.stopPropagation()}
            >
              <button
                onclick={() => {
                  showMobileUserDropdown = false;
                  showPasswordModal = true;
                  resetPasswordForm();
                  isMobileMenuOpen = false;
                }}
                class="w-full flex items-center gap-3 px-4 py-3 text-sm font-medium text-slate-700 hover:bg-slate-50 hover:text-slate-900 transition-colors cursor-pointer text-left border-b border-slate-100"
              >
                <svg
                  class="w-4 h-4 text-slate-400"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  ><path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"
                  ></path></svg
                >
                Ganti Password
              </button>
              <button
                onclick={handleLogout}
                class="w-full flex items-center gap-3 px-4 py-3 text-sm font-medium text-red-600 hover:bg-red-50 hover:text-red-700 transition-colors cursor-pointer text-left"
              >
                <svg
                  class="w-4 h-4 text-red-500"
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
          {/if}
        </div>
      </aside>
    {/if}

    <!-- Main Content Area -->
    <main
      class="flex-1 flex flex-col min-w-0 md:pt-0 pt-16 h-[100dvh] overflow-y-auto print:pt-0 print:h-auto print:overflow-visible print:block bg-transparent"
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

  <!-- Password Modal -->
  {#if showPasswordModal}
    <div
      class="fixed inset-0 z-[100] flex items-center justify-center bg-black/50 backdrop-blur-sm p-4"
    >
      <div
        class="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden animate-in fade-in zoom-in-95 duration-200"
      >
        <div class="p-6">
          <h3 class="text-xl font-bold text-slate-900 mb-1">Ganti Password</h3>
          <p class="text-sm text-slate-500 mb-6">
            Masukkan password lama dan password baru Anda.
          </p>

          {#if passwordError}
            <div
              class="mb-4 p-3 bg-red-100 text-red-700 text-sm rounded-lg border border-red-200"
            >
              {passwordError}
            </div>
          {/if}
          {#if passwordSuccess}
            <div
              class="mb-4 p-3 bg-green-100 text-green-700 text-sm rounded-lg border border-green-200"
            >
              {passwordSuccess}
            </div>
          {/if}

          <form onsubmit={handlePasswordChange} class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-slate-700 mb-1"
                >Password Lama</label
              >
              <div class="relative">
                <input
                  type={showOldPassword ? "text" : "password"}
                  bind:value={oldPassword}
                  class="w-full px-4 py-2 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors pr-10"
                />
                <button
                  type="button"
                  aria-label="Toggle password visibility"
                  class="absolute inset-y-0 right-0 px-3 flex items-center text-slate-400 hover:text-slate-600 cursor-pointer"
                  onclick={() => (showOldPassword = !showOldPassword)}
                >
                  {#if showOldPassword}
                    <svg
                      class="w-5 h-5"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                      ><path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                      ></path><path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                      ></path></svg
                    >
                  {:else}
                    <svg
                      class="w-5 h-5"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                      ><path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"
                      ></path></svg
                    >
                  {/if}
                </button>
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-700 mb-1"
                >Password Baru</label
              >
              <div class="relative">
                <input
                  type={showNewPassword ? "text" : "password"}
                  bind:value={newPassword}
                  class="w-full px-4 py-2 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors pr-10"
                />
                <button
                  type="button"
                  aria-label="Toggle password visibility"
                  class="absolute inset-y-0 right-0 px-3 flex items-center text-slate-400 hover:text-slate-600 cursor-pointer"
                  onclick={() => (showNewPassword = !showNewPassword)}
                >
                  {#if showNewPassword}
                    <svg
                      class="w-5 h-5"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                      ><path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                      ></path><path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                      ></path></svg
                    >
                  {:else}
                    <svg
                      class="w-5 h-5"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                      ><path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"
                      ></path></svg
                    >
                  {/if}
                </button>
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-700 mb-1"
                >Konfirmasi Password Baru</label
              >
              <div class="relative">
                <input
                  type={showConfirmPassword ? "text" : "password"}
                  bind:value={confirmPassword}
                  class="w-full px-4 py-2 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors pr-10"
                />
                <button
                  type="button"
                  aria-label="Toggle password visibility"
                  class="absolute inset-y-0 right-0 px-3 flex items-center text-slate-400 hover:text-slate-600 cursor-pointer"
                  onclick={() => (showConfirmPassword = !showConfirmPassword)}
                >
                  {#if showConfirmPassword}
                    <svg
                      class="w-5 h-5"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                      ><path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                      ></path><path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                      ></path></svg
                    >
                  {:else}
                    <svg
                      class="w-5 h-5"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                      ><path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"
                      ></path></svg
                    >
                  {/if}
                </button>
              </div>
            </div>

            <div class="pt-4 flex gap-3">
              <button
                type="button"
                onclick={() => (showPasswordModal = false)}
                class="flex-1 px-4 py-2 border border-slate-300 text-slate-700 rounded-xl hover:bg-slate-50 font-medium transition-colors cursor-pointer"
              >
                Batal
              </button>
              <button
                type="submit"
                disabled={isSubmittingPassword}
                class="flex-1 px-4 py-2 bg-blue-600 text-white rounded-xl hover:bg-blue-700 font-medium transition-colors disabled:opacity-50 flex justify-center items-center cursor-pointer"
              >
                {#if isSubmittingPassword}
                  <svg
                    class="animate-spin h-5 w-5 text-white"
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    ><circle
                      class="opacity-25"
                      cx="12"
                      cy="12"
                      r="10"
                      stroke="currentColor"
                      stroke-width="4"
                    ></circle><path
                      class="opacity-75"
                      fill="currentColor"
                      d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                    ></path></svg
                  >
                {:else}
                  Simpan
                {/if}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  {/if}
{/if}
