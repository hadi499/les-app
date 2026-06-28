<script lang="ts">
  import { onMount } from "svelte";

  type User = { id: number; username: string; role: string };

  let users: User[] = $state([]);
  let isLoading = $state(true);
  let errorMsg = $state("");

  let showDeleteModal = $state(false);
  let userToDelete: { id: number; username: string } | null = $state(null);
  let isDeleting = $state(false);

  let showAddModal = $state(false);
  let isAdding = $state(false);
  let newUser = $state({ username: "", password: "" });
  let addErrorMsg = $state("");
  let showPassword = $state(false);

  async function fetchUsers() {
    isLoading = true;
    errorMsg = "";
    try {
      const res = await fetch(`/api/users`, {
        credentials: "include",
      });
      if (res.status === 403) {
        errorMsg = "Akses ditolak. Anda bukan Guru.";
        isLoading = false;
        return;
      }
      if (!res.ok) {
        throw new Error("Gagal mengambil data users");
      }
      const data = (await res.json()) as { users: User[] };
      users = data.users || [];
    } catch (e) {
      errorMsg = e instanceof Error ? e.message : String(e);
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    fetchUsers();
  });

  function promptDelete(userId: number, username: string) {
    userToDelete = { id: userId, username };
    showDeleteModal = true;
  }

  async function confirmDelete() {
    if (!userToDelete) return;

    isDeleting = true;
    try {
      const res = await fetch(`/api/users/${userToDelete.id}`, {
        method: "DELETE",
        credentials: "include",
      });

      const data = await res.json();

      if (!res.ok) {
        alert(data.error || "Gagal menghapus user");
        return;
      }

      showDeleteModal = false;
      userToDelete = null;
      fetchUsers();
    } catch (e) {
      alert(
        "Terjadi kesalahan: " + (e instanceof Error ? e.message : String(e)),
      );
    } finally {
      isDeleting = false;
    }
  }

  function cancelDelete() {
    showDeleteModal = false;
    userToDelete = null;
  }

  function promptAdd() {
    newUser = { username: "", password: "" };
    addErrorMsg = "";
    showPassword = false;
    showAddModal = true;
  }

  function cancelAdd() {
    showAddModal = false;
  }

  async function confirmAdd() {
    if (!newUser.username || !newUser.password) {
      addErrorMsg = "Username dan password harus diisi.";
      return;
    }
    isAdding = true;
    addErrorMsg = "";
    try {
      const res = await fetch(`/api/auth/register`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(newUser),
      });
      const data = await res.json();
      if (!res.ok) {
        addErrorMsg = data.error || "Gagal menambahkan user";
        return;
      }
      showAddModal = false;
      fetchUsers();
    } catch (e) {
      addErrorMsg = "Terjadi kesalahan: " + (e instanceof Error ? e.message : String(e));
    } finally {
      isAdding = false;
    }
  }
</script>

<svelte:head>
  <title>Manajemen Users - Portal Guru</title>
</svelte:head>

<div class="animate-in fade-in duration-500">
  <div class="mb-8 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
    <div>
      <h1
        class="text-2xl font-bold text-slate-900 sm:text-3xl tracking-tight drop-shadow-sm"
      >
        Manajemen Users
      </h1>
      <p
        class="mt-2 text-slate-600 text-sm sm:text-base font-light tracking-wide"
      >
        Lihat dan kelola seluruh pengguna di sistem ini.
      </p>
    </div>
    <button
      onclick={promptAdd}
      class="inline-flex items-center gap-2 px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-xl hover:bg-blue-700 transition-colors shadow-md shadow-blue-500/20 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
      Tambah User
    </button>
  </div>

  {#if isLoading}
    <div class="flex justify-center p-12">
      <div
        class="w-10 h-10 border-4 border-slate-200 border-t-blue-500 rounded-full animate-spin"
      ></div>
    </div>
  {:else if errorMsg}
    <div
      class="bg-red-100 text-red-800 p-6 rounded-2xl border border-red-300 font-medium flex items-center gap-3"
    >
      <svg
        class="w-6 h-6 flex-shrink-0"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
        ><path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
        ></path></svg
      >
      {errorMsg}
    </div>
  {:else}
    <div
      class="bg-white/60 backdrop-blur-md rounded-3xl border border-slate-200 shadow-lg shadow-slate-800/10 overflow-hidden"
    >
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse whitespace-nowrap min-w-[600px]">
          <thead>
            <tr class="bg-white/40 border-b border-slate-200">
              <th class="py-4 px-6 font-bold text-slate-900 text-sm">ID</th>
              <th class="py-4 px-6 font-bold text-slate-900 text-sm"
                >Username</th
              >
              <th class="py-4 px-6 font-bold text-slate-900 text-sm">Role</th>
              <th class="py-4 px-6 font-bold text-slate-900 text-sm text-center"
                >Aksi</th
              >
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-200">
            {#each users as u}
              <tr class="hover:bg-white/40 transition-colors">
                <td class="py-4 px-6 text-sm text-slate-600">{u.id}</td>
                <td
                  class="py-4 px-6 text-sm font-medium text-slate-900 drop-shadow-sm"
                  >{u.username}</td
                >
                <td class="py-4 px-6 text-sm">
                  <span
                    class={`px-2.5 py-1 rounded-md text-xs font-medium capitalize border ${u.role === "teacher" ? "bg-purple-100 text-purple-700 border-purple-300" : "bg-blue-100 text-slate-600 border-blue-300"}`}
                  >
                    {u.role}
                  </span>
                </td>
                <td class="py-4 px-6 text-center">
                  <button
                    onclick={() => promptDelete(u.id, u.username)}
                    class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-red-600 bg-red-100 hover:bg-red-200 rounded-lg transition-colors border border-red-300 cursor-pointer focus:outline-none focus:ring-2 focus:ring-red-500/50"
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
                        d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                      ></path></svg
                    >
                    Hapus
                  </button>
                </td>
              </tr>
            {/each}
            {#if users.length === 0}
              <tr>
                <td
                  colspan="4"
                  class="py-8 text-center text-slate-500 font-light"
                  >Tidak ada user ditemukan.</td
                >
              </tr>
            {/if}
          </tbody>
        </table>
      </div>
    </div>
  {/if}

  <!-- Delete Confirmation Modal -->
  {#if showDeleteModal && userToDelete}
    <div
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-sm"
    >
      <!-- animate-in and fade-in handled here or by external css, but keeping simple inline classes -->
      <div
        class="bg-slate-50 backdrop-blur-md rounded-3xl p-6 w-full max-w-md shadow-2xl border border-slate-200 transform transition-all"
      >
        <h3 class="text-xl font-bold text-slate-900 mb-2">Konfirmasi Hapus</h3>
        <p class="text-slate-600 mb-6 text-sm">
          Hapus permanen user <span class="font-bold text-slate-900"
            >"{userToDelete.username}"</span
          >
          beserta seluruh riwayat belajarnya?<br />
          <span class="text-red-500 font-medium block mt-1"
            >Tindakan ini tidak dapat dibatalkan.</span
          >
        </p>
        <div class="flex justify-end gap-3">
          <button
            onclick={cancelDelete}
            disabled={isDeleting}
            class="px-4 py-2 text-sm font-medium text-slate-800 bg-white shadow-md border border-slate-200 hover:bg-slate-100 rounded-xl transition-colors disabled:opacity-50"
          >
            Batal
          </button>
          <button
            onclick={confirmDelete}
            disabled={isDeleting}
            class="px-4 py-2 text-sm font-medium text-white bg-red-500 shadow-md hover:bg-red-600 rounded-xl transition-colors flex items-center gap-2 disabled:opacity-50"
          >
            {#if isDeleting}
              <div
                class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"
              ></div>
            {/if}
            Ya, Hapus
          </button>
        </div>
      </div>
    </div>
  {/if}

  <!-- Add User Modal -->
  {#if showAddModal}
    <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-sm">
      <div class="bg-slate-50 backdrop-blur-md rounded-3xl p-6 w-full max-w-md shadow-2xl border border-slate-200 transform transition-all">
        <h3 class="text-xl font-bold text-slate-900 mb-4">Tambah User Baru</h3>
        
        {#if addErrorMsg}
          <div class="mb-4 p-3 bg-red-100 text-red-700 text-sm rounded-lg border border-red-200">
            {addErrorMsg}
          </div>
        {/if}

        <div class="space-y-4 mb-6">
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1" for="username">Username</label>
            <input
              id="username"
              type="text"
              bind:value={newUser.username}
              class="w-full px-4 py-2 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all bg-white text-slate-900"
              placeholder="Masukkan username"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1" for="password">Password</label>
            <div class="relative">
              <input
                id="password"
                type={showPassword ? "text" : "password"}
                bind:value={newUser.password}
                class="w-full pl-4 pr-12 py-2 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all bg-white text-slate-900"
                placeholder="Masukkan password"
              />
              <button
                type="button"
                class="absolute inset-y-0 right-0 px-3 flex items-center text-slate-400 hover:text-blue-500 focus:outline-none transition-colors cursor-pointer"
                onclick={() => (showPassword = !showPassword)}
                title={showPassword ? "Sembunyikan password" : "Tampilkan password"}
              >
                {#if showPassword}
                  <!-- Eye slash icon -->
                  <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
                  </svg>
                {:else}
                  <!-- Eye icon -->
                  <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.543 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                  </svg>
                {/if}
              </button>
            </div>
          </div>
        </div>

        <div class="flex justify-end gap-3">
          <button
            onclick={cancelAdd}
            disabled={isAdding}
            class="px-4 py-2 text-sm font-medium text-slate-800 bg-white shadow-md border border-slate-200 hover:bg-slate-100 rounded-xl transition-colors disabled:opacity-50"
          >
            Batal
          </button>
          <button
            onclick={confirmAdd}
            disabled={isAdding}
            class="px-4 py-2 text-sm font-medium text-white bg-blue-600 shadow-md hover:bg-blue-700 rounded-xl transition-colors flex items-center gap-2 disabled:opacity-50"
          >
            {#if isAdding}
              <div class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
            {/if}
            Simpan
          </button>
        </div>
      </div>
    </div>
  {/if}
</div>
