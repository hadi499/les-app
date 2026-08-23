<script lang="ts">
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import { onMount, onDestroy } from "svelte";
  import { chatStore } from "$lib/stores/chatStore";

  import type { Snippet } from "svelte";
  let { children }: { children: Snippet } = $props();

  // Data dinamis dari backend
  type User = { id?: number; username: string; role: string; class?: string };
  let user: User = $state({ username: "Loading...", role: "...", class: "-" });
  let isMobileMenuOpen = $state(false);
  let isDesktopSidebarOpen = $state(true);
  let isLoading = $state(true);
  let isNetworkError = $state(false);
  let showReloadButton = $state(false);

  onMount(async () => {
    let timeoutTimer = setTimeout(() => {
      if (isLoading) showReloadButton = true;
    }, 10000);

    const controller = new AbortController();
    const fetchTimeout = setTimeout(() => controller.abort(), 15000);

    try {
      const res = await fetch(`/me`, {
        credentials: "include",
        signal: controller.signal
      });
      clearTimeout(fetchTimeout);

      if (!res.ok) {
        if (res.status === 401 || res.status === 403) {
          window.location.href = "/login";
          return;
        }
        throw new Error("Server error");
      }
      const data = await res.json();
      if (!data.authenticated) {
        window.location.href = "/login";
        return;
      }
      user = data.user;

      // Setup chat
      if (user.id) {
        chatStore.setMyUserId(user.id);
      }
      try {
        const unreadRes = await fetch("/api/chat/unread-count", {
          credentials: "include",
        });
        if (unreadRes.ok) {
          const unreadData = await unreadRes.json();
          chatStore.setUnreadCount(unreadData.count || 0);
        }
      } catch (e) {
        console.error("Failed to fetch unread count", e);
      }
      chatStore.connect();
    } catch (e: any) {
      console.error("Layout init error:", e);
      if (e.name === 'AbortError' || !navigator.onLine || e.message === "Server error" || e.message.includes("Failed to fetch")) {
        isNetworkError = true;
      } else {
        window.location.href = "/login";
      }
    } finally {
      clearTimeout(timeoutTimer);
      isLoading = false;
    }
  });

  onDestroy(() => {
    chatStore.disconnect();
  });

  $effect(() => {
    if (!isLoading) {
      const path = page.url.pathname;
      const teacherOnlyRoutes = [
        "/dashboard/notes",
        "/dashboard/card-memory",
        "/dashboard/subjects",
        "/dashboard/typing-monitoring",
        "/dashboard/parents",
      ];
      const adminOnlyRoutes = [
        "/dashboard/users",
        "/dashboard/teachers",
        "/dashboard/logs",
        "/dashboard/quotes",
      ];

      const currentRole = user?.role?.toLowerCase() || "";
      if (currentRole !== "teacher" && currentRole !== "admin" && teacherOnlyRoutes.some((r) => path.startsWith(r))) {
        goto("/dashboard", { replaceState: true });
      }

      if (currentRole !== "admin" && adminOnlyRoutes.some((r) => path.startsWith(r))) {
        goto("/dashboard", { replaceState: true });
      }

      if (currentRole === "admin") {
        const adminAllowedRoutes = [
          "/dashboard/chat",
          "/dashboard/quotes",
          "/dashboard/users",
          "/dashboard/teachers",
          "/dashboard/parents",
          "/dashboard/logs",
          "/dashboard/card-memory"
        ];
        if (path !== "/dashboard" && path !== "/dashboard/" && !adminAllowedRoutes.some((r) => path.startsWith(r))) {
          goto("/dashboard", { replaceState: true });
        }
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
  let showProfileModal = $state(false);
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

<svelte:head>
  <meta name="robots" content="noindex, nofollow" />
</svelte:head>

<svelte:window
  onclick={() => {
    showUserDropdown = false;
    showMobileUserDropdown = false;
  }}
/>

{#if isLoading}
  <!-- Layar Loading Premium -->
  <div
    class="min-h-screen bg-slate-50 flex items-center justify-center relative overflow-hidden"
  >
    <!-- Ambient Background Orbs -->
    <div
      class="absolute top-[-10%] left-[-10%] w-[80vw] sm:w-[40vw] h-[80vw] sm:h-[40vw] rounded-full bg-indigo-200/60 blur-[80px] sm:blur-[100px] animate-pulse"
      style="animation-duration: 4s;"
    ></div>
    <div
      class="absolute bottom-[-10%] right-[-10%] w-[80vw] sm:w-[40vw] h-[80vw] sm:h-[40vw] rounded-full bg-blue-200/60 blur-[80px] sm:blur-[100px] animate-pulse"
      style="animation-duration: 5s; animation-delay: 1s;"
    ></div>

    <!-- Animation Container -->
    <div
      class="relative z-10 flex flex-col items-center gap-8 text-center"
    >
      <!-- 3 Bouncing/Floating Circles -->
      <div class="flex items-center justify-center gap-3 sm:gap-4 h-20 sm:h-24">
        <div class="w-3 h-3 sm:w-4 sm:h-4 rounded-full bg-blue-600 animate-bounce shadow-lg shadow-blue-600/40" style="animation-delay: 0s;"></div>
        <div class="w-3 h-3 sm:w-4 sm:h-4 rounded-full bg-blue-600 animate-bounce shadow-lg shadow-blue-600/40" style="animation-delay: 0.15s;"></div>
        <div class="w-3 h-3 sm:w-4 sm:h-4 rounded-full bg-blue-600 animate-bounce shadow-lg shadow-blue-600/40" style="animation-delay: 0.3s;"></div>
      </div>

      {#if showReloadButton}
        <div class="mt-4 animate-in fade-in slide-in-from-bottom-2 duration-500">
          <button
            onclick={() => window.location.reload()}
            class="px-6 py-2.5 bg-gradient-to-r from-indigo-500 to-blue-600 hover:from-indigo-600 hover:to-blue-700 text-white rounded-full font-bold text-sm shadow-lg shadow-indigo-500/30 transition-all hover:scale-105 hover:shadow-indigo-500/50 flex items-center gap-2"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path></svg>
            Muat Ulang Aplikasi
          </button>
        </div>
      {/if}
    </div>
  </div>
{:else if isNetworkError}
  <div class="min-h-screen bg-slate-50 flex items-center justify-center relative overflow-hidden px-4">
    <div class="relative z-10 flex flex-col items-center gap-6 bg-white/60 backdrop-blur-xl p-8 sm:p-12 rounded-[2rem] shadow-xl border border-white max-w-md text-center">
      <div class="w-20 h-20 bg-rose-100 rounded-full flex items-center justify-center text-rose-500 mb-2">
        <svg class="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path></svg>
      </div>
      <h2 class="text-2xl font-black text-slate-800">Koneksi Bermasalah</h2>
      <p class="text-slate-600 text-sm leading-relaxed">Gagal memuat data dari server. Hal ini biasanya terjadi jika koneksi internet terputus atau aplikasi kembali dari background.</p>
      <button
        onclick={() => window.location.reload()}
        class="mt-2 w-full px-6 py-3.5 bg-blue-600 hover:bg-blue-700 text-white rounded-xl font-bold text-sm shadow-lg shadow-blue-500/30 transition-all hover:-translate-y-0.5 flex items-center justify-center gap-2"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path></svg>
        Coba Muat Ulang
      </button>
    </div>
  </div>
{:else}
  <div
    class="min-h-screen w-full overflow-x-hidden bg-slate-50 flex selection:bg-slate-200 selection:text-slate-800 font-sans text-slate-900"
    style="--sidebar-width: {isDesktopSidebarOpen ? '16rem' : '0px'};"
  >
    <!-- Sidebar (Desktop) -->
    <aside
      class="bg-slate-400/10 backdrop-blur-md border-r border-slate-300 hidden lg:flex flex-col shrink-0 sticky top-0 h-screen print:hidden shadow-lg shadow-slate-800/5 transition-all duration-300 z-30 {isDesktopSidebarOpen
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
          class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline whitespace-nowrap {page
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
          class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
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
            <svg
              class="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              ><path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
              ></path></svg
            >
          </div>
          <span class="flex-1">Pesan</span>
          {#if $chatStore.unreadCount > 0}
            <span
              class="bg-red-500 text-white text-[10px] font-bold px-2 py-0.5 rounded-full min-w-[20px] text-center"
            >
              {$chatStore.unreadCount}
            </span>
          {/if}
        </a>

        {#if user?.role === "parent"}
          <a
            href="/dashboard/anak"
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
              '/anak',
            )
              ? 'bg-white/80 text-blue-700 font-medium shadow-sm shadow-slate-800/5 border border-slate-300'
              : 'text-slate-700 hover:bg-white/50 hover:text-slate-900 border border-transparent'}"
          >
            <div
              class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
                '/anak',
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
                ></path></svg>
            </div>
            Anak Saya
          </a>
        {/if}

        {#if user?.role !== "admin"}
        <a
          href="/dashboard/materi"
          class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
            '/materi',
          )
            ? 'bg-white/80 text-blue-700 font-medium shadow-sm shadow-slate-800/5 border border-slate-300'
            : 'text-slate-700 hover:bg-white/50 hover:text-slate-900 border border-transparent'}"
        >
          <div
            class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
              '/materi',
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
          Materi Pelajaran
        </a>

        {#if user?.role === "admin" || user?.role === "teacher"}
          <a
            href="/dashboard/subjects"
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
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
                  d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"
                ></path></svg
              >
            </div>
            Mata Pelajaran
          </a>
        {/if}

        <a
          href="/dashboard/todos"
          class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
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
          class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
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

        {#if user?.role === "teacher" || user?.role === "admin"}
          <a
            href="/dashboard/notes"
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
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
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
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
          class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
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

        <a
          href="/dashboard/quizzes"
          class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
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
          class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
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
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
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
        {/if}
        {#if user?.role !== "parent"}
        <a
          href="/dashboard/points"
          class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
            '/points',
          )
            ? 'bg-white/80 text-blue-700 font-medium shadow-sm shadow-slate-800/5 border border-slate-300'
            : 'text-slate-700 hover:bg-white/50 hover:text-slate-900 border border-transparent'}"
        >
          <div
            class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
              '/points',
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
                d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              ></path></svg
            >
          </div>
          Poin Murid
        </a>
        {/if}
        {/if}
        {#if user?.role === "admin"}
          <a
            href="/dashboard/card-memory"
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
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
          <a
            href="/dashboard/quotes"
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
              '/quotes',
            )
              ? 'bg-white/80 text-blue-700 font-medium shadow-sm shadow-slate-800/5 border border-slate-300'
              : 'text-slate-700 hover:bg-white/50 hover:text-slate-900 border border-transparent'}"
          >
            <div
              class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
                '/quotes',
              )
                ? 'bg-blue-600 text-white shadow-md shadow-blue-500/20'
                : 'bg-slate-100 text-slate-500 group-hover:bg-slate-200 group-hover:text-slate-700'}"
            >
              <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24"
                ><path
                  d="M14.017 21v-7.391c0-5.704 3.731-9.57 8.983-10.609l.995 2.151c-2.432.917-3.995 3.638-3.995 5.849h4v10h-9.983zm-14.017 0v-7.391c0-5.704 3.748-9.57 9-10.609l.996 2.151c-2.433.917-3.996 3.638-3.996 5.849h3.983v10h-9.983z"
                ></path></svg
              >
            </div>
            Quotes Inspirasi
          </a>
          <a
            href="/dashboard/users"
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
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
                ></path></svg>
            </div>
            Manajemen Users
          </a>
          <a
            href="/dashboard/teachers"
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
              '/teachers',
            )
              ? 'bg-white/80 text-blue-700 font-medium shadow-sm shadow-slate-800/5 border border-slate-300'
              : 'text-slate-700 hover:bg-white/50 hover:text-slate-900 border border-transparent'}"
          >
            <div
              class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
                '/teachers',
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
                  d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
                ></path></svg>
            </div>
            Manajemen Guru
          </a>
        {/if}
        {#if user?.role === "admin" || user?.role === "teacher"}
          <a
            href="/dashboard/parents"
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
              '/parents',
            )
              ? 'bg-white/80 text-blue-700 font-medium shadow-sm shadow-slate-800/5 border border-slate-300'
              : 'text-slate-700 hover:bg-white/50 hover:text-slate-900 border border-transparent'}"
          >
            <div
              class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
                '/parents',
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
                  d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
                ></path></svg>
            </div>
            Orang tua murid
          </a>
        {/if}
        {#if user?.role === "admin"}
          <a
            href="/dashboard/logs"
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
              '/logs',
            )
              ? 'bg-white/80 text-blue-700 font-medium shadow-sm shadow-slate-800/5 border border-slate-300'
              : 'text-slate-700 hover:bg-white/50 hover:text-slate-900 border border-transparent'}"
          >
            <div
              class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
                '/logs',
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
                ></path></svg>
            </div>
            User Logs
          </a>
        {/if}

        <div class="pt-4 mt-2 border-t border-slate-200">
          <a
            href="/"
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-[15px] transition-colors text-slate-700 hover:bg-white/50 hover:text-slate-900 no-underline whitespace-nowrap"
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
                showProfileModal = true;
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
                  d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
                ></path></svg
              >
              Profil Saya
            </button>
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
      class="hidden lg:flex fixed top-4 z-50 p-1.5 bg-white border border-slate-200 shadow-sm rounded-lg text-slate-500 hover:bg-slate-50 hover:text-blue-600 transition-all duration-300 cursor-pointer"
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
      class="lg:hidden fixed top-0 inset-x-0 h-16 bg-slate-400/10 backdrop-blur-md border-b border-slate-200 z-50 flex items-center justify-between px-4 shadow-sm print:hidden"
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
        class="lg:hidden fixed inset-0 z-40 bg-black/60 backdrop-blur-sm"
        role="button"
        tabindex="-1"
        aria-label="Tutup menu"
        onclick={() => (isMobileMenuOpen = false)}
        onkeydown={(e) => {
          if (e.key === "Enter" || e.key === " ") isMobileMenuOpen = false;
        }}
      ></div>
      <aside
        class="lg:hidden fixed inset-y-0 left-0 w-64 bg-slate-100 border-r border-slate-200 shadow-2xl z-50 flex flex-col transform transition-transform {isMobileMenuOpen
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
            class="text-slate-400 hover:text-red-500 border-none bg-transparent cursor-pointer"
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
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors no-underline whitespace-nowrap {page
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
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
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
              <svg
                class="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
                ></path></svg
              >
            </div>
            <span class="flex-1">Pesan</span>
            {#if $chatStore.unreadCount > 0}
              <span
                class="bg-red-500 text-white text-[10px] font-bold px-2 py-0.5 rounded-full min-w-[20px] text-center"
              >
                {$chatStore.unreadCount}
              </span>
            {/if}
          </a>

          {#if user?.role === "parent"}
            <a
              href="/dashboard/anak"
              onclick={() => (isMobileMenuOpen = false)}
              class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
                '/anak',
              )
                ? 'bg-white/80 text-blue-700 font-medium border border-slate-300'
                : 'text-slate-700 hover:text-slate-900 hover:bg-white/50 border border-transparent'}"
            >
              <div
                class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
                  '/anak',
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
                  ></path></svg>
              </div>
              Anak Saya
            </a>
          {/if}

          {#if user?.role !== "admin"}
          <a
            href="/dashboard/materi"
            onclick={() => (isMobileMenuOpen = false)}
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
              '/materi',
            )
              ? 'bg-white/80 text-blue-700 font-medium border border-slate-300'
              : 'text-slate-700 hover:text-slate-900 hover:bg-white/50 border border-transparent'}"
          >
            <div
              class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
                '/materi',
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
            Materi Pelajaran
          </a>

          {#if user?.role === "admin" || user?.role === "teacher"}
            <a
              href="/dashboard/subjects"
              onclick={() => (isMobileMenuOpen = false)}
              class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
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
                    d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"
                  ></path></svg
                >
              </div>
              Mata Pelajaran
            </a>
          {/if}

          <a
            href="/dashboard/todos"
            onclick={() => (isMobileMenuOpen = false)}
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
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
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
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

          {#if user?.role === "teacher" || user?.role === "admin"}
            <a
              href="/dashboard/notes"
              onclick={() => (isMobileMenuOpen = false)}
              class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
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
              class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
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
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
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

          <a
            href="/dashboard/quizzes"
            onclick={() => (isMobileMenuOpen = false)}
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
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
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
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
              class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
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
          {/if}
          {#if user?.role !== "parent"}
          <a
            href="/dashboard/points"
            onclick={() => (isMobileMenuOpen = false)}
            class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
              '/points',
            )
              ? 'bg-white/80 text-blue-700 font-medium border border-slate-300'
              : 'text-slate-700 hover:text-slate-900 hover:bg-white/50 border border-transparent'}"
          >
            <div
              class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
                '/points',
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
                  d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                ></path></svg
              >
            </div>
            Poin Murid
          </a>
          {/if}
          {/if}
          {#if user?.role === "admin"}
            <a
              href="/dashboard/card-memory"
              onclick={() => (isMobileMenuOpen = false)}
              class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
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
            <a
              href="/dashboard/quotes"
              onclick={() => (isMobileMenuOpen = false)}
              class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
                '/quotes',
              )
                ? 'bg-white/80 text-blue-700 font-medium border border-slate-300'
                : 'text-slate-700 hover:text-slate-900 hover:bg-white/50 border border-transparent'}"
            >
              <div
                class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
                  '/quotes',
                )
                  ? 'bg-blue-600 text-white shadow-md shadow-blue-500/20'
                  : 'bg-slate-100 text-slate-500 group-hover:bg-slate-200 group-hover:text-slate-700'}"
              >
                <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24"
                  ><path
                    d="M14.017 21v-7.391c0-5.704 3.731-9.57 8.983-10.609l.995 2.151c-2.432.917-3.995 3.638-3.995 5.849h4v10h-9.983zm-14.017 0v-7.391c0-5.704 3.748-9.57 9-10.609l.996 2.151c-2.433.917-3.996 3.638-3.996 5.849h3.983v10h-9.983z"
                  ></path></svg
                >
              </div>
              Quotes Inspirasi
            </a>
            <a
              href="/dashboard/users"
              onclick={() => (isMobileMenuOpen = false)}
              class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
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
                  ></path></svg>
              </div>
              Manajemen Users
            </a>
            <a
              href="/dashboard/teachers"
              onclick={() => (isMobileMenuOpen = false)}
              class="group flex items-center gap-3 px-4 py-3 rounded-2xl font-medium text-[15px] transition-all {page.url.pathname.includes(
                '/teachers',
              )
                ? 'bg-blue-50 text-blue-700 shadow-sm shadow-blue-500/10'
                : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'}"
            >
              <div
                class="flex items-center justify-center w-10 h-10 rounded-xl shrink-0 transition-all {page.url.pathname.includes(
                  '/teachers',
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
                    d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
                  ></path></svg>
              </div>
              Manajemen Guru
            </a>
          {/if}
          {#if user?.role === "admin" || user?.role === "teacher"}
            <a
              href="/dashboard/parents"
              onclick={() => (isMobileMenuOpen = false)}
              class="group flex items-center gap-3 px-4 py-3 rounded-2xl font-medium text-[15px] transition-all {page.url.pathname.includes(
                '/parents',
              )
                ? 'bg-blue-50 text-blue-700 shadow-sm shadow-blue-500/10'
                : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'}"
            >
              <div
                class="flex items-center justify-center w-10 h-10 rounded-xl shrink-0 transition-all {page.url.pathname.includes(
                  '/parents',
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
                    d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
                  ></path></svg>
              </div>
              Orang tua murid
            </a>
          {/if}
          {#if user?.role === "admin"}
            <a
              href="/dashboard/logs"
              onclick={() => (isMobileMenuOpen = false)}
              class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors no-underline whitespace-nowrap {page.url.pathname.includes(
                '/logs',
              )
                ? 'bg-white/80 text-blue-700 font-medium border border-slate-300'
                : 'text-slate-700 hover:text-slate-900 hover:bg-white/50 border border-transparent'}"
            >
              <div
                class="flex items-center justify-center w-8 h-8 rounded-lg shrink-0 transition-all {page.url.pathname.includes(
                  '/logs',
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
                  ></path></svg>
              </div>
              User Logs
            </a>
          {/if}

          <div class="pt-4 mt-2 border-t border-slate-200">
            <a
              href="/"
              onclick={() => (isMobileMenuOpen = false)}
              class="group flex items-center gap-2 px-3 py-2 rounded-xl font-normal text-base transition-colors text-slate-700 hover:bg-white/50 hover:text-slate-900 no-underline whitespace-nowrap"
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
                  showProfileModal = true;
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
                    d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
                  ></path></svg
                >
                Profil Saya
              </button>
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
      class="flex-1 flex flex-col min-w-0 lg:pt-0 pt-16 h-[100dvh] overflow-y-auto print:pt-0 print:h-auto print:overflow-visible print:block bg-transparent"
    >
      <div
        class="flex-1 p-6 lg:p-8 {isDesktopSidebarOpen
          ? ''
          : 'lg:pl-16'} max-w-6xl mx-auto w-full print:p-0 print:m-0"
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
                        d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"
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
                        d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                      ></path><path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
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
                        d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"
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
                        d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                      ></path><path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
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
                        d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"
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
                        d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                      ></path><path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
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

  {#if showProfileModal}
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div
      class="fixed inset-0 bg-slate-900/50 backdrop-blur-sm z-100 flex items-center justify-center p-4 animate-in fade-in duration-200"
      onclick={() => (showProfileModal = false)}
    >
      <div
        class="bg-white rounded-3xl shadow-2xl w-full max-w-md overflow-hidden animate-in zoom-in-95 duration-200"
        onclick={(e) => e.stopPropagation()}
      >
        <div class="p-6 relative">
          <button
            aria-label="Tutup Profil"
            onclick={() => (showProfileModal = false)}
            class="absolute top-4 right-4 text-slate-400 hover:text-slate-600 bg-white/50 rounded-full p-1 transition-colors cursor-pointer"
          >
            <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
          
          <h3 class="text-xl font-bold text-slate-900 mb-6 flex items-center gap-2">
            <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V8a2 2 0 00-2-2h-5m-4 0V5a2 2 0 114 0v1m-4 0a2 2 0 104 0m-5 8a2 2 0 100-4 2 2 0 000 4zm0 0c1.306 0 2.417.835 2.83 2M9 14a3.001 3.001 0 00-2.83 2M15 11h3m-3 4h2" />
            </svg>
            Profil Pengguna
          </h3>

          <!-- ATM Card UI -->
          <div class="relative bg-linear-to-br from-blue-700 via-blue-800 to-indigo-900 rounded-2xl p-4 sm:p-6 text-white shadow-xl overflow-hidden aspect-[1.586/1] flex flex-col justify-between">
            <!-- decorative shapes -->
            <div class="absolute top-0 right-0 -mr-8 -mt-8 w-32 h-32 rounded-full bg-white/10 blur-2xl pointer-events-none"></div>
            <div class="absolute bottom-0 left-0 -ml-8 -mb-8 w-24 h-24 rounded-full bg-blue-400/20 blur-xl pointer-events-none"></div>
            
            <div class="relative z-10 flex justify-between items-start">
              <div class="flex flex-col">
                <span class="text-[10px] sm:text-xs text-blue-200 uppercase tracking-wider font-semibold mb-0.5 sm:mb-1">ID Number</span>
                <span class="font-mono text-base sm:text-lg tracking-widest leading-none">{user.id ? user.id.toString().padStart(6, '0') : '000000'}</span>
              </div>
              <svg class="w-6 h-6 sm:w-8 sm:h-8 text-blue-200/80" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
              </svg>
            </div>

            <div class="relative z-10 my-auto py-2">
              <div class="text-[10px] sm:text-xs text-blue-200 uppercase tracking-wider font-semibold mb-0.5 sm:mb-1">Username</div>
              <div class="text-lg sm:text-xl font-bold tracking-wide truncate leading-none">{user.username}</div>
            </div>

            <div class="relative z-10 flex justify-between items-end mt-auto pt-2">
              <div>
                <div class="text-sm sm:text-base font-medium capitalize leading-none text-blue-100">{user.role}</div>
              </div>
              <div class="text-right">
                <span class="text-sm sm:text-base font-medium leading-none">Kelas {user.class && user.class !== "-" ? user.class : "-"}</span>
              </div>
            </div>
          </div>
          
        </div>
      </div>
    </div>
  {/if}
{/if}
