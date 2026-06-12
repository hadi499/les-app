<script>
  import { onMount } from 'svelte';
  
  let users = $state([]);
  let isLoading = $state(true);
  let errorMsg = $state('');

  async function fetchUsers() {
    isLoading = true;
    errorMsg = '';
    try {
      const res = await fetch('http://localhost:8080/api/users', {
        credentials: 'include'
      });
      if (res.status === 403) {
        errorMsg = 'Akses ditolak. Anda bukan Guru.';
        isLoading = false;
        return;
      }
      if (!res.ok) {
        throw new Error('Gagal mengambil data users');
      }
      const data = await res.json();
      users = data.users || [];
    } catch(e) {
      errorMsg = e.message;
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    fetchUsers();
  });

  async function handleDelete(userId, username) {
    if (!confirm(`Hapus permanen user "${username}" beserta seluruh riwayat belajarnya? Tindakan ini tidak dapat dibatalkan.`)) {
      return;
    }

    try {
      const res = await fetch(`http://localhost:8080/api/users/${userId}`, {
        method: 'DELETE',
        credentials: 'include'
      });
      
      const data = await res.json();
      
      if (!res.ok) {
        alert(data.error || 'Gagal menghapus user');
        return;
      }
      
      alert('User berhasil dihapus');
      fetchUsers();
    } catch(e) {
      alert('Terjadi kesalahan: ' + e.message);
    }
  }
</script>

<svelte:head>
  <title>Manajemen Users - Portal Guru</title>
</svelte:head>

<div class="animate-in fade-in duration-500">
  <div class="mb-8">
    <h1 class="text-2xl font-bold text-white sm:text-3xl tracking-tight drop-shadow-sm">Manajemen Users</h1>
    <p class="mt-2 text-blue-200 text-sm sm:text-base font-light tracking-wide">Lihat dan kelola seluruh pengguna di sistem ini.</p>
  </div>

  {#if isLoading}
    <div class="flex justify-center p-12">
      <div class="w-10 h-10 border-4 border-indigo-900 border-t-indigo-500 rounded-full animate-spin"></div>
    </div>
  {:else if errorMsg}
    <div class="bg-red-900/30 text-red-200 p-6 rounded-2xl border border-red-800/50 font-medium flex items-center gap-3">
      <svg class="w-6 h-6 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
      {errorMsg}
    </div>
  {:else}
    <div class="bg-zinc-900/50 backdrop-blur-md rounded-3xl border border-zinc-800 shadow-lg shadow-black/30 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-zinc-800/50 border-b border-zinc-800">
              <th class="py-4 px-6 font-bold text-white text-sm">ID</th>
              <th class="py-4 px-6 font-bold text-white text-sm">Username</th>
              <th class="py-4 px-6 font-bold text-white text-sm">Role</th>
              <th class="py-4 px-6 font-bold text-white text-sm text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-800/50">
            {#each users as u}
              <tr class="hover:bg-zinc-800/50 transition-colors">
                <td class="py-4 px-6 text-sm text-zinc-400">{u.id}</td>
                <td class="py-4 px-6 text-sm font-medium text-white drop-shadow-sm">{u.username}</td>
                <td class="py-4 px-6 text-sm">
                  <span class={`px-2.5 py-1 rounded-md text-xs font-medium capitalize border ${u.role === 'teacher' ? 'bg-purple-900/30 text-purple-300 border-purple-800/50' : 'bg-blue-900/30 text-blue-300 border-blue-800/50'}`}>
                    {u.role}
                  </span>
                </td>
                <td class="py-4 px-6 text-right">
                  <button 
                    onclick={() => handleDelete(u.id, u.username)}
                    class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-red-400 bg-red-900/30 hover:bg-red-800/50 rounded-lg transition-colors border border-red-800/50 cursor-pointer focus:outline-none focus:ring-2 focus:ring-red-500/50"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
                    Hapus
                  </button>
                </td>
              </tr>
            {/each}
            {#if users.length === 0}
              <tr>
                <td colspan="4" class="py-8 text-center text-zinc-500 font-light">Tidak ada user ditemukan.</td>
              </tr>
            {/if}
          </tbody>
        </table>
      </div>
    </div>
  {/if}
</div>
