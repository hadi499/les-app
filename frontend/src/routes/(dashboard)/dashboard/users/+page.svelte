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
    <h1 class="text-2xl font-bold text-orange-950 sm:text-3xl tracking-tight drop-shadow-sm">Manajemen Users</h1>
    <p class="mt-2 text-orange-800 text-sm sm:text-base font-light tracking-wide">Lihat dan kelola seluruh pengguna di sistem ini.</p>
  </div>

  {#if isLoading}
    <div class="flex justify-center p-12">
      <div class="w-10 h-10 border-4 border-orange-200 border-t-orange-500 rounded-full animate-spin"></div>
    </div>
  {:else if errorMsg}
    <div class="bg-red-100 text-red-800 p-6 rounded-2xl border border-red-300 font-medium flex items-center gap-3">
      <svg class="w-6 h-6 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
      {errorMsg}
    </div>
  {:else}
    <div class="bg-white/60 backdrop-blur-md rounded-3xl border border-orange-200 shadow-lg shadow-orange-900/10 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-white/40 border-b border-orange-200">
              <th class="py-4 px-6 font-bold text-orange-950 text-sm">ID</th>
              <th class="py-4 px-6 font-bold text-orange-950 text-sm">Username</th>
              <th class="py-4 px-6 font-bold text-orange-950 text-sm">Role</th>
              <th class="py-4 px-6 font-bold text-orange-950 text-sm text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-orange-200">
            {#each users as u}
              <tr class="hover:bg-white/40 transition-colors">
                <td class="py-4 px-6 text-sm text-orange-800">{u.id}</td>
                <td class="py-4 px-6 text-sm font-medium text-orange-950 drop-shadow-sm">{u.username}</td>
                <td class="py-4 px-6 text-sm">
                  <span class={`px-2.5 py-1 rounded-md text-xs font-medium capitalize border ${u.role === 'teacher' ? 'bg-purple-100 text-purple-700 border-purple-300' : 'bg-blue-100 text-orange-800 border-blue-300'}`}>
                    {u.role}
                  </span>
                </td>
                <td class="py-4 px-6 text-right">
                  <button 
                    onclick={() => handleDelete(u.id, u.username)}
                    class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-red-600 bg-red-100 hover:bg-red-200 rounded-lg transition-colors border border-red-300 cursor-pointer focus:outline-none focus:ring-2 focus:ring-red-500/50"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
                    Hapus
                  </button>
                </td>
              </tr>
            {/each}
            {#if users.length === 0}
              <tr>
                <td colspan="4" class="py-8 text-center text-orange-700 font-light">Tidak ada user ditemukan.</td>
              </tr>
            {/if}
          </tbody>
        </table>
      </div>
    </div>
  {/if}
</div>
