<script lang="ts">
  import { page } from "$app/state";
  import { onMount } from "svelte";

  type Child = {
    id: number;
    username: string;
  };

  type User = {
    id: number;
    username: string;
    role: string;
    is_suspended: boolean;
    last_active_at: string;
    created_at: string;
    children: Child[];
  };

  let user: User | null = $state(null);
  let isLoading = $state(true);
  let errorMsg = $state("");

  const id = page.params.id;

  async function fetchUser() {
    isLoading = true;
    errorMsg = "";
    try {
      const res = await fetch(`/api/users/${id}`, {
        credentials: "include",
      });
      if (!res.ok) {
        throw new Error("Gagal mengambil data pengguna");
      }
      const data = await res.json();
      user = data.user;
    } catch (e: any) {
      console.error(e);
      errorMsg = e.message;
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    fetchUser();
  });
</script>

<svelte:head>
  <title>Detail Orang Tua - Portal Admin</title>
</svelte:head>

<div class="animate-in fade-in duration-500 relative">
  <!-- Header -->
  <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-8">
    <div>
      <div class="flex items-center gap-2 mb-2">
        <a href="/dashboard/parents" class="text-blue-600 hover:text-blue-800 flex items-center gap-1 text-sm font-medium transition-colors">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path></svg>
          Kembali
        </a>
      </div>
      <h1 class="text-2xl font-bold text-slate-900 sm:text-3xl tracking-tight drop-shadow-sm">
        Detail Orang Tua
      </h1>
      <p class="mt-2 text-slate-600 text-sm sm:text-base font-light tracking-wide">
        Informasi detail mengenai akun orang tua ini.
      </p>
    </div>
  </div>

  {#if isLoading}
    <div class="flex flex-col items-center justify-center py-20 bg-white/60 backdrop-blur-xl rounded-3xl border border-white/80 shadow-xl">
      <div class="w-12 h-12 border-4 border-blue-100 border-t-blue-600 rounded-full animate-spin mb-4 shadow-md"></div>
      <p class="text-slate-500 font-medium">Memuat data orang tua...</p>
    </div>
  {:else if errorMsg}
    <div class="p-6 bg-red-50 text-red-700 rounded-3xl border border-red-200 shadow-sm flex items-start gap-4">
      <div class="p-2 bg-red-100 rounded-xl shrink-0">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
      </div>
      <div>
        <h3 class="font-bold text-lg mb-1">Gagal Memuat Data</h3>
        <p class="text-sm opacity-90">{errorMsg}</p>
        <button onclick={fetchUser} class="mt-4 px-4 py-2 bg-red-600 text-white text-sm font-bold rounded-xl hover:bg-red-700 transition-colors shadow-sm cursor-pointer">
          Coba Lagi
        </button>
      </div>
    </div>
  {:else if user}
    <div class="bg-white/80 backdrop-blur-xl rounded-3xl shadow-xl border border-white p-6 sm:p-8 space-y-8">
      
      <!-- Profile Header -->
      <div class="flex items-center gap-6 pb-8 border-b border-slate-100">
        <div class="w-24 h-24 rounded-2xl bg-indigo-100 text-indigo-600 flex items-center justify-center font-bold text-4xl shadow-inner border border-indigo-50 shrink-0">
          {user.username.charAt(0).toUpperCase()}
        </div>
        <div>
          <h2 class="text-2xl font-bold text-slate-900 leading-tight mb-1">{user.username}</h2>
          <div class="flex items-center gap-3">
            <span class="px-3 py-1 bg-purple-100 text-purple-700 text-xs font-bold rounded-full capitalize shadow-sm border border-purple-200">{user.role}</span>
            <span class="text-slate-400 text-sm font-medium">ID: {user.id}</span>
          </div>
        </div>
      </div>

      <!-- Detail Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div class="bg-slate-50 rounded-2xl p-5 border border-slate-100">
          <p class="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-1">Status Akun</p>
          <p class="font-medium text-lg">
            {#if user.is_suspended}
              <span class="text-red-600 flex items-center gap-2">
                <span class="w-2 h-2 rounded-full bg-red-600 shrink-0"></span> Suspended
              </span>
            {:else if user.last_active_at && Date.now() - new Date(user.last_active_at).getTime() < 5 * 60 * 1000}
              <span class="text-green-600 flex items-center gap-2">
                <span class="w-2 h-2 rounded-full bg-green-600 shrink-0"></span> Online
              </span>
            {:else}
              <span class="text-slate-500 flex items-center gap-2">
                <span class="w-2 h-2 rounded-full bg-slate-400 shrink-0"></span> Offline
              </span>
            {/if}
          </p>
        </div>
        
        <div class="bg-slate-50 rounded-2xl p-5 border border-slate-100">
          <p class="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-1">Terdaftar Pada</p>
          <p class="font-medium text-slate-800 text-lg">
            {user.created_at ? new Date(user.created_at).toLocaleString('id-ID', {timeZone: 'Asia/Jakarta', day: 'numeric', month: 'short', year: 'numeric'}) : "-"}
          </p>
        </div>
      </div>

      <!-- Anak-Anak -->
      <div class="bg-slate-50 p-6 rounded-3xl border border-slate-100">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-bold text-slate-900 flex items-center gap-2">
            <svg class="w-5 h-5 text-blue-500 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"></path></svg>
            Daftar Anak
          </h3>
          <span class="px-3 py-1 bg-blue-100 text-blue-700 text-xs font-bold rounded-full">{user.children ? user.children.length : 0} Anak</span>
        </div>
        
        {#if user.children && user.children.length > 0}
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            {#each user.children as child}
              <div class="flex items-center gap-3 bg-white p-4 rounded-2xl shadow-sm border border-slate-100 hover:border-blue-200 transition-colors">
                <div class="w-10 h-10 rounded-full bg-blue-100 text-blue-600 flex items-center justify-center font-bold text-sm shrink-0">
                  {child.username.charAt(0).toUpperCase()}
                </div>
                <div class="flex-1 min-w-0">
                  <p class="font-bold text-slate-900 truncate">{child.username}</p>
                  <p class="text-xs text-slate-500 truncate">ID: {child.id}</p>
                </div>
              </div>
            {/each}
          </div>
        {:else}
          <div class="text-center py-8 bg-white rounded-2xl border border-slate-100 border-dashed">
            <div class="w-12 h-12 bg-slate-100 rounded-full flex items-center justify-center mx-auto mb-3 text-slate-400">
              <svg class="w-6 h-6 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"></path></svg>
            </div>
            <p class="text-slate-500 font-medium">Belum ada data anak.</p>
          </div>
        {/if}
      </div>

    </div>
  {/if}
</div>
