<script>
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  
  let { children } = $props();

  // Data dinamis dari backend
  let user = $state({ username: 'Loading...', role: '...' });
  let isMobileMenuOpen = $state(false);
  let isLoading = $state(true);

  onMount(async () => {
    try {
      const res = await fetch('http://localhost:8080/me', { credentials: 'include' });
      if (!res.ok) {
        goto('/login');
        return;
      }
      const data = await res.json();
      if (!data.authenticated) {
        goto('/login');
        return;
      }
      user = data.user;
    } catch(e) {
      goto('/login');
    } finally {
      isLoading = false;
    }
  });

  async function handleLogout() {
    try {
      await fetch('http://localhost:8080/api/auth/logout', {
        method: 'POST',
        credentials: 'include'
      });
      goto('/login');
    } catch(e) {
      console.error(e);
      goto('/login');
    }
  }
</script>

{#if isLoading}
  <!-- Layar Loading Sederhana -->
  <div class="min-h-screen bg-[#0C134F] flex items-center justify-center">
    <div class="flex flex-col items-center gap-4">
      <div class="w-12 h-12 border-4 border-indigo-900 border-t-indigo-500 rounded-full animate-spin"></div>
      <p class="text-blue-200 font-medium">Memuat portal...</p>
    </div>
  </div>
{:else}
<div class="min-h-screen bg-[#0C134F] flex selection:bg-zinc-800 selection:text-white font-sans text-white">
  
  <!-- Sidebar (Desktop) -->
  <aside class="w-64 bg-zinc-900/50 backdrop-blur-md border-r border-zinc-800 hidden md:flex flex-col flex-shrink-0 sticky top-0 h-screen print:hidden shadow-lg shadow-black/20">
    <!-- Logo -->
    <div class="h-16 flex items-center px-6 border-b border-zinc-800/50">
      <div class="flex items-center gap-3">
        <div class="w-8 h-8 bg-zinc-800 border border-zinc-600 rounded-lg flex items-center justify-center text-white font-extrabold text-sm shadow-md">
          L
        </div>
        <span class="font-bold text-lg text-white tracking-tight drop-shadow-sm">{user?.role === 'teacher' ? 'Portal Guru' : 'Portal Siswa'}</span>
      </div>
    </div>

    <!-- Nav Links -->
    <nav class="flex-1 overflow-y-auto py-6 px-4 space-y-1.5">
      <a href="/dashboard" class="flex items-center gap-3 px-3 py-2.5 rounded-xl font-medium transition-colors no-underline { $page.url.pathname === '/dashboard' ? 'bg-blue-900/30 text-blue-200 shadow-sm shadow-blue-900/10 border border-blue-800/50' : 'text-zinc-400 hover:bg-zinc-800/50 hover:text-white' }">
        <svg class="w-5 h-5 { $page.url.pathname === '/dashboard' ? 'text-blue-400' : 'text-zinc-500' }" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"></path></svg>
        Dashboard
      </a>

      
      {#if user?.role === 'teacher'}
        <a href="/dashboard/users" class="flex items-center gap-3 px-3 py-2.5 rounded-xl font-medium transition-colors no-underline { $page.url.pathname.includes('/users') ? 'bg-blue-900/30 text-blue-200 shadow-sm shadow-blue-900/10 border border-blue-800/50' : 'text-zinc-400 hover:bg-zinc-800/50 hover:text-white' }">
          <svg class="w-5 h-5 { $page.url.pathname.includes('/users') ? 'text-blue-400' : 'text-zinc-500' }" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"></path></svg>
          Manajemen Users
        </a>
        
        <a href="/dashboard/card-memory" class="flex items-center gap-3 px-3 py-2.5 rounded-xl font-medium transition-colors no-underline { $page.url.pathname.includes('/card-memory') ? 'bg-blue-900/30 text-blue-200 shadow-sm shadow-blue-900/10 border border-blue-800/50' : 'text-zinc-400 hover:bg-zinc-800/50 hover:text-white' }">
          <svg class="w-5 h-5 { $page.url.pathname.includes('/card-memory') ? 'text-blue-400' : 'text-zinc-500' }" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path></svg>
          Card Memory
        </a>
      {/if}
      
      <div class="pt-4 mt-2 border-t border-zinc-800/50">
        <a href="/" class="flex items-center gap-3 px-3 py-2.5 rounded-xl font-medium transition-colors text-zinc-400 hover:bg-zinc-800/50 hover:text-white no-underline">
          <svg class="w-5 h-5 text-zinc-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9"></path></svg>
          Halaman Utama
        </a>
      </div>
    </nav>

    <!-- User & Logout -->
    <div class="p-4 border-t border-zinc-800 bg-zinc-900/80">
      <div class="flex items-center gap-3 mb-4 px-2">
        <div class="w-10 h-10 rounded-full bg-zinc-800 flex items-center justify-center text-indigo-400 font-bold border border-zinc-600 shadow-sm flex-shrink-0">
          {user.username ? user.username.charAt(0).toUpperCase() : 'U'}
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-bold text-white truncate m-0">{user.username}</p>
          <p class="text-xs font-light text-blue-200 truncate capitalize m-0 mt-0.5">{user.role}</p>
        </div>
      </div>
      <button onclick={handleLogout} class="w-full flex items-center justify-center gap-2 px-4 py-2.5 text-sm font-medium text-red-400 bg-zinc-800 border border-zinc-700 hover:bg-zinc-700 hover:text-red-300 rounded-xl transition-all shadow-sm cursor-pointer">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"></path></svg>
        Keluar Akun
      </button>
    </div>
  </aside>

  <!-- Mobile Header -->
  <div class="md:hidden fixed top-0 inset-x-0 h-16 bg-zinc-900/80 backdrop-blur-md border-b border-zinc-800 z-50 flex items-center justify-between px-4 shadow-sm print:hidden">
     <div class="flex items-center gap-2">
        <div class="w-8 h-8 bg-zinc-800 border border-zinc-600 rounded-lg flex items-center justify-center text-white font-extrabold text-sm shadow-md">
          L
        </div>
        <span class="font-bold text-lg text-white tracking-tight drop-shadow-sm">{user?.role === 'teacher' ? 'Portal Guru' : 'Portal Siswa'}</span>
     </div>
     <button class="text-zinc-400 focus:outline-none p-2 bg-zinc-800/50 rounded-lg border border-zinc-700 hover:bg-zinc-800 transition-colors" onclick={() => isMobileMenuOpen = !isMobileMenuOpen}>
       <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path></svg>
     </button>
  </div>
  
  {#if isMobileMenuOpen}
  <div class="md:hidden fixed inset-0 z-40 bg-black/60 backdrop-blur-sm" onclick={() => isMobileMenuOpen = false}></div>
  <aside class="md:hidden fixed inset-y-0 left-0 w-64 bg-zinc-900 border-r border-zinc-800 shadow-2xl z-50 flex flex-col transform transition-transform {isMobileMenuOpen ? 'translate-x-0' : '-translate-x-full'} print:hidden">
     <!-- Copy of sidebar content for mobile -->
     <div class="h-16 flex items-center justify-between px-6 border-b border-zinc-800/50 bg-zinc-900">
        <span class="font-bold text-lg text-white tracking-tight drop-shadow-sm">Menu Utama</span>
        <button onclick={() => isMobileMenuOpen = false} class="text-zinc-500 hover:text-white border-none bg-transparent cursor-pointer"><svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg></button>
     </div>
     <nav class="flex-1 overflow-y-auto py-6 px-4 space-y-1.5 bg-zinc-900">
       <a href="/dashboard" class="flex items-center gap-3 px-3 py-2.5 rounded-xl font-medium transition-colors no-underline { $page.url.pathname === '/dashboard' ? 'bg-blue-900/30 text-blue-200 border border-blue-800/50' : 'text-zinc-400 hover:text-white hover:bg-zinc-800/50' }">
         <svg class="w-5 h-5 { $page.url.pathname === '/dashboard' ? 'text-blue-400' : 'text-zinc-500' }" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"></path></svg>
         Dashboard
       </a>

       {#if user?.role === 'teacher'}
         <a href="/dashboard/users" class="flex items-center gap-3 px-3 py-2.5 rounded-xl font-medium transition-colors no-underline { $page.url.pathname.includes('/users') ? 'bg-blue-900/30 text-blue-200 border border-blue-800/50' : 'text-zinc-400 hover:text-white hover:bg-zinc-800/50' }">
           <svg class="w-5 h-5 { $page.url.pathname.includes('/users') ? 'text-blue-400' : 'text-zinc-500' }" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"></path></svg>
           Manajemen Users
         </a>

         <a href="/dashboard/card-memory" class="flex items-center gap-3 px-3 py-2.5 rounded-xl font-medium transition-colors no-underline { $page.url.pathname.includes('/card-memory') ? 'bg-blue-900/30 text-blue-200 border border-blue-800/50' : 'text-zinc-400 hover:text-white hover:bg-zinc-800/50' }">
           <svg class="w-5 h-5 { $page.url.pathname.includes('/card-memory') ? 'text-blue-400' : 'text-zinc-500' }" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path></svg>
           Card Memory
         </a>
       {/if}

       <div class="pt-4 mt-2 border-t border-zinc-800/50">
         <a href="/" class="flex items-center gap-3 px-3 py-2.5 rounded-xl font-medium transition-colors text-zinc-400 hover:bg-zinc-800/50 hover:text-white no-underline">
           <svg class="w-5 h-5 text-zinc-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9"></path></svg>
           Halaman Utama
         </a>
       </div>
     </nav>
  </aside>
  {/if}

  <!-- Main Content Area -->
  <main class="flex-1 flex flex-col min-w-0 md:pt-0 pt-16 h-screen overflow-y-auto print:pt-0 print:h-auto print:overflow-visible print:block bg-transparent z-0 relative">
    <div class="flex-1 p-6 md:p-8 max-w-6xl mx-auto w-full relative z-10 print:p-0 print:m-0">
      {@render children()}
    </div>
  </main>
</div>
{/if}
