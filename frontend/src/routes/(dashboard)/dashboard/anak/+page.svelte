<script lang="ts">
  import { onMount } from "svelte";

  type Child = {
    id: number;
    username: string;
    class: string;
    points: number;
    last_active_at?: string;
  };

  let children: Child[] = $state([]);
  let isLoading = $state(true);
  let errorMsg = $state("");

  async function fetchChildren() {
    isLoading = true;
    errorMsg = "";
    try {
      const res = await fetch(`/api/parents/my-children`, {
        credentials: "include",
      });
      if (!res.ok) {
        if (res.status === 403) {
          throw new Error("Akses ditolak. Anda bukan akun Orang Tua.");
        }
        throw new Error("Gagal mengambil data anak");
      }
      const data = await res.json();
      children = data.children || [];
    } catch (e) {
      errorMsg = e instanceof Error ? e.message : String(e);
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    fetchChildren();
  });
</script>

<svelte:head>
  <title>Anak Saya - Portal Orang Tua</title>
</svelte:head>

<div class="animate-in fade-in duration-500 relative max-w-5xl mx-auto">
  <div class="mb-8">
    <h1 class="text-2xl font-bold text-slate-900 sm:text-3xl tracking-tight drop-shadow-sm">
      Anak Saya
    </h1>
    <p class="mt-2 text-slate-600 text-sm sm:text-base font-light tracking-wide">
      Pantau daftar anak Anda yang terdaftar di aplikasi ini.
    </p>
  </div>

  {#if isLoading}
    <div class="flex justify-center p-12">
      <div class="w-10 h-10 border-4 border-slate-200 border-t-indigo-600 rounded-full animate-spin"></div>
    </div>
  {:else if errorMsg}
    <div class="bg-red-100 text-red-800 p-6 rounded-2xl border border-red-300 font-medium">
      {errorMsg}
    </div>
  {:else if children.length === 0}
    <div class="bg-white rounded-3xl border border-slate-200 shadow-sm p-12 text-center">
      <svg class="mx-auto h-12 w-12 text-slate-400 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
      </svg>
      <h3 class="text-lg font-medium text-slate-900">Belum Ada Anak yang Ditautkan</h3>
      <p class="mt-2 text-slate-500 text-sm">
        Akun Anda saat ini belum ditautkan ke akun murid manapun. Silakan hubungi Guru atau Admin untuk menautkan akun anak Anda.
      </p>
    </div>
  {:else}
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
      {#each children as child}
        <div class="bg-white rounded-3xl p-6 shadow-sm border border-slate-200 hover:shadow-md transition-shadow relative overflow-hidden group">
          <div class="absolute top-0 right-0 w-32 h-32 bg-indigo-50 rounded-bl-full -z-10 group-hover:scale-110 transition-transform duration-500"></div>
          
          <div class="flex items-start justify-between mb-4">
            <div class="w-12 h-12 rounded-full bg-indigo-100 text-indigo-600 flex items-center justify-center font-bold text-xl">
              {child.username.charAt(0).toUpperCase()}
            </div>
            
            <div class="flex flex-col items-end">
              <span class="px-2.5 py-1 bg-amber-100 text-amber-700 text-xs font-bold rounded-full">
                {child.points} Poin
              </span>
              {#if child.last_active_at && Date.now() - new Date(child.last_active_at).getTime() < 5 * 60 * 1000}
                <span class="inline-flex items-center gap-1 mt-2 text-[10px] font-medium text-green-600">
                  <span class="w-1.5 h-1.5 rounded-full bg-green-500 animate-pulse"></span>
                  Online
                </span>
              {/if}
            </div>
          </div>
          
          <h3 class="text-xl font-bold text-slate-900 mb-1">{child.username}</h3>
          <p class="text-sm text-slate-500 mb-6">Kelas: {child.class || "Tidak Ada"}</p>
          

        </div>
      {/each}
    </div>
  {/if}
</div>
