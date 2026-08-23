<script lang="ts">
  import { onMount } from "svelte";

  type User = {
    id: number;
    username: string;
    role: string;
    class?: string;
    created_at?: string;
    last_active_at?: string;
    points?: number;
    is_suspended?: boolean;
    parent_id?: number | null;
    parent?: User;
    children?: User[];
    teacher_id?: number | null;
    teacher?: User;
    students?: User[];
  };

  let parentsList: User[] = $state([]);
  let teachersList: User[] = $state([]);

  let users: User[] = $state([]);
  let isLoading = $state(true);
  let showLoadingSpinner = $state(false);
  let errorMsg = $state("");
  let totalUsers = $state(0);
  let totalParents = $state(0);
  let totalStudents = $state(0);

  // Paginasi
  let currentPage = $state(1);
  let itemsPerPage = $state(25);
  let totalPages = $state(1);

  let showDeleteModal = $state(false);
  let userToDelete: { id: number; username: string } | null = $state(null);
  let isDeleting = $state(false);

  let showAddModal = $state(false);
  let isAdding = $state(false);
  let newUser = $state({
    username: "",
    password: "",
    role: "student",
    class: "",
    parent_id: null as number | null,
    teacher_id: null as number | null,
  });
  let addErrorMsg = $state("");
  let showPassword = $state(false);

  let showEditModal = $state(false);
  let isEditing = $state(false);
  let editUser = $state({ id: 0, username: "", role: "student", class: "", parent_id: null as number | null, remove_parent: false, teacher_id: null as number | null, remove_teacher: false });
  let editErrorMsg = $state("");

  let showResetModal = $state(false);
  let userToReset: { id: number; username: string } | null = $state(null);
  let isResetting = $state(false);
  let resetPassword = $state("");
  let resetErrorMsg = $state("");
  let showResetPassword = $state(false);

  let showSuspendModal = $state(false);
  let userToSuspend: {
    id: number;
    username: string;
    is_suspended: boolean;
  } | null = $state(null);
  let isSuspending = $state(false);

  let flashMessage = $state("");
  let flashType = $state("success"); // "success" | "error"

  function showFlash(msg: string, type: "success" | "error" = "success") {
    flashMessage = msg;
    flashType = type;
    setTimeout(() => {
      flashMessage = "";
    }, 3000);
  }

  async function fetchUsers() {
    isLoading = true;
    showLoadingSpinner = false;
    setTimeout(() => {
      showLoadingSpinner = true;
    }, 150);
    errorMsg = "";
    try {
      const res = await fetch(`/api/users?role=parent&page=${currentPage}&limit=${itemsPerPage}`, {
        credentials: "include",
      });
      if (res.status === 403) {
        errorMsg = "Akses ditolak. Anda bukan Admin.";
        isLoading = false;
        return;
      }
      if (!res.ok) {
        throw new Error("Gagal mengambil data users");
      }
      const data = (await res.json()) as { 
        users: User[],
        total_users: number,
        total_parents: number,
        total_students: number,
        total_pages: number
      };
      users = data.users || [];
      totalUsers = data.total_users || 0;
      totalParents = data.total_parents || 0;
      totalStudents = data.total_students || 0;
      totalPages = data.total_pages || 1;
      
      // Pastikan halaman tidak melebihi total halaman (misal saat penghapusan data)
      if (currentPage > totalPages && totalPages > 0) {
        currentPage = totalPages;
        fetchUsers();
      }
    } catch (e) {
      errorMsg = e instanceof Error ? e.message : String(e);
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    fetchUsers();
    fetchParents();
    fetchTeachers();
  });

  async function fetchParents() {
    try {
      const res = await fetch(`/api/users?role=parent&limit=1000`, {
        credentials: "include",
      });
      if (res.ok) {
        const data = await res.json();
        parentsList = data.users || [];
      }
    } catch (e) {
      console.error("Gagal mengambil daftar parent:", e);
    }
  }

  async function fetchTeachers() {
    try {
      const res = await fetch(`/api/users?role=parent&limit=1000`, {
        credentials: "include",
      });
      if (res.ok) {
        const data = await res.json();
        teachersList = data.users || [];
      }
    } catch (e) {
      console.error("Gagal mengambil daftar guru:", e);
    }
  }

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
        showFlash(data.error || "Gagal menghapus user", "error");
        return;
      }

      showFlash("User berhasil dihapus", "success");
      showDeleteModal = false;
      userToDelete = null;
      fetchUsers();
      fetchParents();
    } catch (e) {
      showFlash(
        "Terjadi kesalahan: " + (e instanceof Error ? e.message : String(e)),
        "error",
      );
    } finally {
      isDeleting = false;
    }
  }

  function cancelDelete() {
    showDeleteModal = false;
    userToDelete = null;
  }

  function promptAdd(role: string) {
    newUser = { username: "", password: "", role: role, class: "", parent_id: null, teacher_id: null };
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
    if (newUser.password.length < 6) {
      addErrorMsg = "Password minimal 6 karakter.";
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
      fetchParents();
    } catch (e) {
      addErrorMsg =
        "Terjadi kesalahan: " + (e instanceof Error ? e.message : String(e));
    } finally {
      isAdding = false;
    }
  }

  function promptEdit(u: User) {
    editUser = {
      id: u.id,
      username: u.username,
      role: u.role,
      class: u.class || "",
      parent_id: u.parent_id || null,
      remove_parent: false,
      teacher_id: u.teacher_id || null,
      remove_teacher: false,
    };
    editErrorMsg = "";
    showEditModal = true;
  }

  function cancelEdit() {
    showEditModal = false;
  }

  async function confirmEdit() {
    if (!editUser.username) {
      editErrorMsg = "Username harus diisi.";
      return;
    }
    isEditing = true;
    editErrorMsg = "";
    try {
      const res = await fetch(`/api/users/${editUser.id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({
          username: editUser.username,
          role: editUser.role,
          class: editUser.role === 'student' ? editUser.class : undefined,
          parent_id: editUser.role === 'student' ? (editUser.parent_id || undefined) : undefined,
          remove_parent: editUser.role === 'student' && editUser.parent_id === null,
          teacher_id: editUser.role === 'student' ? (editUser.teacher_id || undefined) : undefined,
          remove_teacher: editUser.role === 'student' && editUser.teacher_id === null,
        }),
      });
      const data = await res.json();
      if (!res.ok) {
        editErrorMsg = data.error || "Gagal mengubah user";
        return;
      }
      showFlash("User berhasil diubah", "success");
      showEditModal = false;
      fetchUsers();
      fetchParents();
    } catch (e) {
      editErrorMsg =
        "Terjadi kesalahan: " + (e instanceof Error ? e.message : String(e));
    } finally {
      isEditing = false;
    }
  }

  function promptReset(userId: number, username: string) {
    userToReset = { id: userId, username };
    resetPassword = "";
    resetErrorMsg = "";
    showResetPassword = false;
    showResetModal = true;
  }

  function cancelReset() {
    showResetModal = false;
    userToReset = null;
  }

  async function confirmReset() {
    if (!userToReset) return;
    if (!resetPassword) {
      resetErrorMsg = "Password baru harus diisi.";
      return;
    }
    if (resetPassword.length < 6) {
      resetErrorMsg = "Password minimal 6 karakter.";
      return;
    }

    isResetting = true;
    resetErrorMsg = "";
    try {
      const res = await fetch(`/api/users/${userToReset.id}/reset-password`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ new_password: resetPassword }),
      });

      const data = await res.json();

      if (!res.ok) {
        resetErrorMsg = data.error || "Gagal mereset password";
        return;
      }

      showFlash("Password berhasil direset!", "success");
      showResetModal = false;
      userToReset = null;
    } catch (e) {
      resetErrorMsg =
        "Terjadi kesalahan: " + (e instanceof Error ? e.message : String(e));
    } finally {
      isResetting = false;
    }
  }

  function promptSuspend(u: User) {
    userToSuspend = {
      id: u.id,
      username: u.username,
      is_suspended: u.is_suspended || false,
    };
    showSuspendModal = true;
  }

  function cancelSuspend() {
    showSuspendModal = false;
    userToSuspend = null;
  }

  async function confirmSuspend() {
    if (!userToSuspend) return;
    isSuspending = true;
    try {
      const res = await fetch(`/api/users/${userToSuspend.id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ is_suspended: !userToSuspend.is_suspended }),
      });
      const data = await res.json();
      if (!res.ok) {
        showFlash(data.error || "Gagal mengubah status suspend", "error");
        return;
      }
      showFlash(
        userToSuspend.is_suspended
          ? "User berhasil diaktifkan kembali"
          : "User berhasil ditangguhkan",
        "success",
      );
      showSuspendModal = false;
      userToSuspend = null;
      fetchUsers();
    } catch (e) {
      showFlash(
        "Terjadi kesalahan: " + (e instanceof Error ? e.message : String(e)),
        "error",
      );
    } finally {
      isSuspending = false;
    }
  }
</script>

<svelte:head>
  <title>Manajemen Orang Tua - Portal Admin</title>
</svelte:head>

<div class="animate-in fade-in duration-500 relative">
  <!-- Flash Message -->
  {#if flashMessage}
    <div
      class="fixed top-6 right-6 z-[110] flex items-center gap-3 px-5 py-3 rounded-2xl shadow-xl shadow-slate-900/10 border animate-in slide-in-from-right-8 fade-in duration-300 {flashType ===
      'success'
        ? 'bg-green-50 border-green-200 text-green-800'
        : 'bg-red-50 border-red-200 text-red-800'}"
    >
      {#if flashType === "success"}
        <svg
          class="w-5 h-5 text-green-600"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          ><path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M5 13l4 4L19 7"
          ></path></svg
        >
      {:else}
        <svg
          class="w-5 h-5 text-red-600"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          ><path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M6 18L18 6M6 6l12 12"
          ></path></svg
        >
      {/if}
      <p class="text-sm font-semibold">{flashMessage}</p>
    </div>
  {/if}

  <div
    class="mb-8 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4"
  >
    <div>
      <h1
        class="text-2xl font-bold text-slate-900 sm:text-3xl tracking-tight drop-shadow-sm"
      >
        Manajemen Orang Tua
      </h1>
      <p
        class="mt-2 text-slate-600 text-sm sm:text-base font-light tracking-wide"
      >
        Lihat dan kelola seluruh guru di sistem ini.
      </p>
    </div>
    <div class="flex flex-wrap gap-2">
      <button
        onclick={() => promptAdd('parent')}
        class="inline-flex items-center gap-2 px-4 py-2 text-sm font-bold text-purple-700 bg-purple-50 hover:bg-purple-100 border border-purple-200 rounded-xl transition-all shadow-sm focus:outline-none focus:ring-2 focus:ring-purple-500/50 cursor-pointer"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"
          ><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
        Tambah Orang Tua
      </button>
    </div>
  </div>

  {#if isLoading}
    <div
      class="fixed inset-0 z-[100] flex flex-col items-center justify-center bg-slate-50/50 backdrop-blur-sm {showLoadingSpinner
        ? 'opacity-100'
        : 'opacity-0'} transition-opacity duration-300"
    >
      <div
        class="w-12 h-12 border-4 border-slate-200 border-t-indigo-600 rounded-full animate-spin shadow-sm"
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
    <div class="flex items-center gap-2 mb-4">
      <span class="px-3 py-1 bg-purple-100 text-purple-700 text-xs font-bold rounded-full border border-purple-200 shadow-sm" title="Total Orang Tua">
        {totalParents} Total Orang Tua
      </span>
    </div>
    <div
      class="bg-white/60 rounded-3xl border border-slate-200 shadow-lg shadow-slate-800/10 overflow-hidden"
    >
      <div class="overflow-x-auto">
        <table
          class="w-full text-left border-collapse whitespace-nowrap min-w-[600px]"
        >
          <thead>
            <tr class="bg-white/40 border-b border-slate-200">
              <th class="py-4 px-6 font-bold text-slate-900 text-sm">ID</th>
              <th class="py-4 px-6 font-bold text-slate-900 text-sm"
                >Username</th
              >
              <th class="py-4 px-6 font-bold text-slate-900 text-sm">Role</th>
              <th class="py-4 px-6 font-bold text-slate-900 text-sm">Terdaftar</th>
              <th class="py-4 px-6 font-bold text-slate-900 text-sm">Status</th>
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
                <td class="py-4 px-6 text-sm text-slate-700">
                  {u.created_at ? new Date(u.created_at).toLocaleString('id-ID', {timeZone: 'Asia/Jakarta', day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit'}) : "-"}
                </td>

                <td class="py-4 px-6 text-sm flex flex-col gap-1 items-start">
                  {#if u.is_suspended}
                    <span
                      class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-md text-xs font-medium bg-red-100 text-red-700 border border-red-300"
                    >
                      Suspended
                    </span>
                  {/if}
                  {#if u.last_active_at && Date.now() - new Date(u.last_active_at).getTime() < 5 * 60 * 1000}
                    <span
                      class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-md text-xs font-medium bg-green-100 text-green-700 border border-green-300 mt-1"
                    >
                      <span
                        class="w-1.5 h-1.5 rounded-full bg-green-500 animate-pulse"
                      ></span>
                      Online
                    </span>
                  {:else}
                    <span
                      class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-md text-xs font-medium bg-slate-100 text-slate-600 border border-slate-300 mt-1"
                    >
                      <span class="w-1.5 h-1.5 rounded-full bg-slate-400"
                      ></span>
                      Offline
                    </span>
                  {/if}
                </td>
                <td class="py-4 px-6 text-center">
                  <div class="flex items-center justify-center gap-2">
                    <a
                      href="/dashboard/parents/{u.id}"
                      class="p-2 bg-blue-50 text-blue-600 hover:bg-blue-100 rounded-xl transition-colors cursor-pointer flex items-center justify-center"
                      title="Detail User"
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
                          d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                        ></path><path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.543 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                        ></path></svg
                      >
                    </a>
                    <button
                      onclick={() => promptEdit(u)}
                      class="inline-flex items-center p-1.5 text-blue-700 bg-blue-50 hover:bg-blue-100 rounded-lg transition-colors border border-blue-200 cursor-pointer focus:outline-none focus:ring-2 focus:ring-blue-500/50"
                      title="Edit User"
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
                          d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                        ></path></svg
                      >
                    </button>
                    <button
                      onclick={() => promptReset(u.id, u.username)}
                      class="inline-flex items-center p-1.5 text-amber-700 bg-amber-50 hover:bg-amber-100 rounded-lg transition-colors border border-amber-200 cursor-pointer focus:outline-none focus:ring-2 focus:ring-amber-500/50"
                      title="Reset Password"
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
                          d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"
                        ></path></svg
                      >
                    </button>
                    <button
                      onclick={() => promptSuspend(u)}
                      class="inline-flex items-center p-1.5 {u.is_suspended
                        ? 'text-green-700 bg-green-50 hover:bg-green-100 border-green-200 focus:ring-green-500/50'
                        : 'text-orange-700 bg-orange-50 hover:bg-orange-100 border-orange-200 focus:ring-orange-500/50'} rounded-lg transition-colors border cursor-pointer focus:outline-none focus:ring-2"
                      title={u.is_suspended ? "Aktifkan User" : "Suspend User"}
                    >
                      {#if u.is_suspended}
                        <svg
                          class="w-4 h-4"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                          ><path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                          ></path></svg
                        >
                      {:else}
                        <svg
                          class="w-4 h-4"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                          ><path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"
                          ></path></svg
                        >
                      {/if}
                    </button>
                    <button
                      onclick={() => promptDelete(u.id, u.username)}
                      class="inline-flex items-center p-1.5 text-red-600 bg-red-100 hover:bg-red-200 rounded-lg transition-colors border border-red-300 cursor-pointer focus:outline-none focus:ring-2 focus:ring-red-500/50"
                      title="Hapus User"
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
                    </button>
                  </div>
                </td>
              </tr>
            {/each}
            {#if users.length === 0}
              <tr>
                <td
                  colspan="5"
                  class="py-8 text-center text-slate-500 font-light"
                  >Tidak ada user ditemukan.</td
                >
              </tr>
            {/if}
          </tbody>
        </table>
      </div>
      <!-- Pagination Controls -->
      {#if totalPages > 1}
        <div class="px-4 py-4 sm:px-6 border-t border-slate-200 flex items-center justify-between bg-slate-50/50">
          
          <!-- Mobile View -->
          <div class="flex flex-1 justify-between sm:hidden items-center">
            <button
              disabled={currentPage === 1}
              onclick={() => { currentPage -= 1; fetchUsers(); }}
              class="relative inline-flex items-center justify-center w-10 h-10 rounded-xl border border-slate-300 bg-white text-slate-700 hover:bg-slate-50 disabled:opacity-50 disabled:cursor-not-allowed shadow-sm transition-all"
            >
              <span class="sr-only">Sebelumnya</span>
              <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                <path fill-rule="evenodd" d="M12.79 5.23a.75.75 0 01-.02 1.06L8.832 10l3.938 3.71a.75.75 0 11-1.04 1.08l-4.5-4.25a.75.75 0 010-1.08l4.5-4.25a.75.75 0 011.06.02z" clip-rule="evenodd" />
              </svg>
            </button>
            <span class="text-sm font-medium text-slate-500">
              Hal {currentPage} / {totalPages}
            </span>
            <button
              disabled={currentPage === totalPages}
              onclick={() => { currentPage += 1; fetchUsers(); }}
              class="relative inline-flex items-center justify-center w-10 h-10 rounded-xl border border-slate-300 bg-white text-slate-700 hover:bg-slate-50 disabled:opacity-50 disabled:cursor-not-allowed shadow-sm transition-all"
            >
              <span class="sr-only">Selanjutnya</span>
              <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                <path fill-rule="evenodd" d="M7.21 14.77a.75.75 0 01.02-1.06L11.168 10 7.23 6.29a.75.75 0 111.04-1.08l4.5 4.25a.75.75 0 010 1.08l-4.5 4.25a.75.75 0 01-1.06-.02z" clip-rule="evenodd" />
              </svg>
            </button>
          </div>

          <!-- Desktop View -->
          <div class="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
            <div>
              <p class="text-sm text-slate-600 font-light">
                Menampilkan <span class="font-semibold text-slate-900">{(currentPage - 1) * itemsPerPage + 1}</span> hingga <span class="font-semibold text-slate-900">{Math.min(currentPage * itemsPerPage, totalUsers)}</span> dari <span class="font-semibold text-slate-900">{totalUsers}</span> user
              </p>
            </div>
            <div>
              <nav class="isolate inline-flex -space-x-px rounded-xl shadow-sm" aria-label="Pagination">
                <button
                  disabled={currentPage === 1}
                  onclick={() => { currentPage -= 1; fetchUsers(); }}
                  class="relative inline-flex items-center rounded-l-xl px-3 py-2 text-slate-400 ring-1 ring-inset ring-slate-300 hover:bg-slate-50 focus:z-20 focus:outline-offset-0 disabled:opacity-50 disabled:cursor-not-allowed transition-all bg-white"
                >
                  <span class="sr-only">Sebelumnya</span>
                  <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                    <path fill-rule="evenodd" d="M12.79 5.23a.75.75 0 01-.02 1.06L8.832 10l3.938 3.71a.75.75 0 11-1.04 1.08l-4.5-4.25a.75.75 0 010-1.08l4.5-4.25a.75.75 0 011.06.02z" clip-rule="evenodd" />
                  </svg>
                </button>
                
                {#each Array.from({ length: totalPages }) as _, i}
                  <button
                    onclick={() => { currentPage = i + 1; fetchUsers(); }}
                    class="relative inline-flex items-center px-4 py-2 text-sm font-semibold {currentPage === i + 1 ? 'z-10 bg-blue-600 text-white ring-blue-600 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600' : 'text-slate-900 ring-1 ring-inset ring-slate-300 hover:bg-slate-50 focus:z-20 focus:outline-offset-0 bg-white'} transition-all"
                  >
                    {i + 1}
                  </button>
                {/each}
                
                <button
                  disabled={currentPage === totalPages}
                  onclick={() => { currentPage += 1; fetchUsers(); }}
                  class="relative inline-flex items-center rounded-r-xl px-3 py-2 text-slate-400 ring-1 ring-inset ring-slate-300 hover:bg-slate-50 focus:z-20 focus:outline-offset-0 disabled:opacity-50 disabled:cursor-not-allowed transition-all bg-white"
                >
                  <span class="sr-only">Selanjutnya</span>
                  <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                    <path fill-rule="evenodd" d="M7.21 14.77a.75.75 0 01.02-1.06L11.168 10 7.23 6.29a.75.75 0 111.04-1.08l4.5 4.25a.75.75 0 010 1.08l-4.5 4.25a.75.75 0 01-1.06-.02z" clip-rule="evenodd" />
                  </svg>
                </button>
              </nav>
            </div>
          </div>
        </div>
      {/if}
    </div>
  {/if}

  <!-- Delete Confirmation Modal -->
  {#if showDeleteModal && userToDelete}
    <div
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-sm"
    >
      <!-- animate-in and fade-in handled here or by external css, but keeping simple inline classes -->
      <div
        class="bg-slate-50 rounded-3xl p-6 w-full max-w-md shadow-2xl border border-slate-200"
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
    <div
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-sm"
    >
      <div
        class="bg-slate-50 rounded-3xl p-6 w-full max-w-md shadow-2xl border border-slate-200"
      >
        <h3 class="text-xl font-bold text-slate-900 mb-4">
          Tambah {newUser.role === 'student' ? 'Murid' : newUser.role === 'teacher' ? 'Guru' : newUser.role === 'parent' ? 'Orang Tua' : 'User'} Baru
        </h3>

        {#if addErrorMsg}
          <div
            class="mb-4 p-3 bg-red-100 text-red-700 text-sm rounded-lg border border-red-200"
          >
            {addErrorMsg}
          </div>
        {/if}

        <div class="space-y-4 mb-6">
          <div>
            <label
              class="block text-sm font-medium text-slate-700 mb-1"
              for="username">Username</label
            >
            <input
              id="username"
              type="text"
              bind:value={newUser.username}
              class="w-full px-4 py-2 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-shadow bg-white text-slate-900"
              placeholder="Masukkan username"
            />
          </div>
          <div>
            <label
              class="block text-sm font-medium text-slate-700 mb-1"
              for="password">Password</label
            >
            <div class="relative">
              <input
                id="password"
                type={showPassword ? "text" : "password"}
                bind:value={newUser.password}
                class="w-full pl-4 pr-12 py-2 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-shadow bg-white text-slate-900"
                placeholder="Masukkan password"
              />
              <button
                type="button"
                class="absolute inset-y-0 right-0 px-3 flex items-center text-slate-400 hover:text-blue-500 focus:outline-none transition-colors cursor-pointer"
                onclick={() => (showPassword = !showPassword)}
                title={showPassword
                  ? "Sembunyikan password"
                  : "Tampilkan password"}
              >
                {#if showPassword}
                  <!-- Eye slash icon -->
                  <svg
                    class="h-5 w-5"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"
                    />
                  </svg>
                {:else}
                  <!-- Eye icon -->
                  <svg
                    class="h-5 w-5"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                    />
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.543 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                    />
                  </svg>
                {/if}
              </button>
            </div>
          </div>
          <div>
            <label
              class="block text-sm font-medium text-slate-700 mb-1"
              for="role">Role</label
            >
            <input
              id="role"
              type="text"
              value={newUser.role === 'student' ? 'Murid' : newUser.role === 'teacher' ? 'Guru' : newUser.role === 'parent' ? 'Orang Tua' : newUser.role}
              disabled
              class="w-full px-4 py-2 border border-slate-300 rounded-xl bg-slate-100 text-slate-500 cursor-not-allowed"
            />
          </div>
          {#if newUser.role === "student"}
            <div>
              <label
                class="block text-sm font-medium text-slate-700 mb-1"
                for="addParentId">Pilih Orang Tua (Opsional)</label
              >
              <select
                id="addParentId"
                bind:value={newUser.parent_id}
                class="w-full px-4 py-2 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-shadow bg-white text-slate-900"
              >
                <option value={null}>-- Tidak Ada --</option>
                {#each parentsList as p}
                  <option value={p.id}>{p.username}</option>
                {/each}
              </select>
            </div>
            <div>
              <label
                class="block text-sm font-medium text-slate-700 mb-1"
                for="addTeacherId">Pilih Guru (Opsional)</label
              >
              <select
                id="addTeacherId"
                bind:value={newUser.teacher_id}
                class="w-full px-4 py-2 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-shadow bg-white text-slate-900"
              >
                <option value={null}>-- Tidak Ada --</option>
                {#each teachersList as t}
                  <option value={t.id}>{t.username}</option>
                {/each}
              </select>
            </div>
            <div>
              <label
                class="block text-sm font-medium text-slate-700 mb-1"
                for="userclass">Kelas</label
              >
              <input
                id="userclass"
                type="text"
                bind:value={newUser.class}
                class="w-full px-4 py-2 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-shadow bg-white text-slate-900"
                placeholder="Misal: SD, SMP, PAUD"
              />
            </div>
          {/if}
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
              <div
                class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"
              ></div>
            {/if}
            Simpan
          </button>
        </div>
      </div>
    </div>
  {/if}

  <!-- Edit User Modal -->
  {#if showEditModal}
    <div
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-sm"
    >
      <div
        class="bg-slate-50 rounded-3xl p-6 w-full max-w-md shadow-2xl border border-slate-200"
      >
        <h3 class="text-xl font-bold text-slate-900 mb-4">Edit User</h3>

        {#if editErrorMsg}
          <div
            class="mb-4 p-3 bg-red-100 text-red-700 text-sm rounded-lg border border-red-200"
          >
            {editErrorMsg}
          </div>
        {/if}

        <div class="space-y-4 mb-6">
          <div>
            <label
              class="block text-sm font-medium text-slate-700 mb-1"
              for="editUsername">Username</label
            >
            <input
              id="editUsername"
              type="text"
              bind:value={editUser.username}
              class="w-full px-4 py-2 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-shadow bg-white text-slate-900"
            />
          </div>
          <div>
            <label
              class="block text-sm font-medium text-slate-700 mb-1"
              for="editRole">Role</label
            >
            <select
              id="editRole"
              bind:value={editUser.role}
              class="w-full px-4 py-2 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-shadow bg-white text-slate-900"
            >
              <option value="student">Student</option>
              <option value="teacher">Teacher</option>
              <option value="admin">Admin</option>
              <option value="parent">Orang Tua</option>
            </select>
          </div>
          {#if editUser.role === "student"}
            <div>
              <label
                class="block text-sm font-medium text-slate-700 mb-1"
                for="editParentId">Tautkan ke Orang Tua (Opsional)</label
              >
              <div class="flex items-center gap-2">
                <select
                  id="editParentId"
                  bind:value={editUser.parent_id}
                  disabled={editUser.remove_parent}
                  class="flex-1 px-4 py-2 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-shadow bg-white text-slate-900 disabled:opacity-50"
                >
                  <option value={null}>-- Pilih Orang Tua --</option>
                  {#each parentsList as p}
                    <option value={p.id}>{p.username}</option>
                  {/each}
                </select>
                <label class="flex items-center gap-2 text-sm text-slate-700">
                  <input type="checkbox" bind:checked={editUser.remove_parent} class="rounded border-slate-300 text-blue-600 focus:ring-blue-500" />
                  Hapus Tautan
                </label>
              </div>
            </div>
            <div>
              <label
                class="block text-sm font-medium text-slate-700 mb-1"
                for="editTeacherId">Tautkan ke Guru (Opsional)</label
              >
              <div class="flex items-center gap-2">
                <select
                  id="editTeacherId"
                  bind:value={editUser.teacher_id}
                  disabled={editUser.remove_teacher}
                  class="flex-1 px-4 py-2 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-shadow bg-white text-slate-900 disabled:opacity-50"
                >
                  <option value={null}>-- Pilih Guru --</option>
                  {#each teachersList as t}
                    <option value={t.id}>{t.username}</option>
                  {/each}
                </select>
                <label class="flex items-center gap-2 text-sm text-slate-700">
                  <input type="checkbox" bind:checked={editUser.remove_teacher} class="rounded border-slate-300 text-blue-600 focus:ring-blue-500" />
                  Hapus Tautan
                </label>
              </div>
            </div>
            <div>
              <label
                class="block text-sm font-medium text-slate-700 mb-1"
                for="editClass">Kelas</label
              >
              <input
                id="editClass"
                type="text"
                bind:value={editUser.class}
                class="w-full px-4 py-2 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-shadow bg-white text-slate-900"
                placeholder="Misal: SD, SMP, PAUD"
              />
            </div>
          {/if}
        </div>

        <div class="flex justify-end gap-3">
          <button
            onclick={cancelEdit}
            disabled={isEditing}
            class="px-4 py-2 text-sm font-medium text-slate-800 bg-white shadow-md border border-slate-200 hover:bg-slate-100 rounded-xl transition-colors disabled:opacity-50"
          >
            Batal
          </button>
          <button
            onclick={confirmEdit}
            disabled={isEditing}
            class="px-4 py-2 text-sm font-medium text-white bg-blue-600 shadow-md hover:bg-blue-700 rounded-xl transition-colors flex items-center gap-2 disabled:opacity-50"
          >
            {#if isEditing}
              <div
                class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"
              ></div>
            {/if}
            Simpan
          </button>
        </div>
      </div>
    </div>
  {/if}

  <!-- Reset Password Modal -->
  {#if showResetModal && userToReset}
    <div
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-sm"
    >
      <div
        class="bg-slate-50 rounded-3xl p-6 w-full max-w-md shadow-2xl border border-slate-200"
      >
        <h3 class="text-xl font-bold text-slate-900 mb-2">Reset Password</h3>
        <p class="text-slate-600 mb-4 text-sm">
          Masukkan password baru untuk user <span
            class="font-bold text-slate-900">"{userToReset.username}"</span
          >.
        </p>

        {#if resetErrorMsg}
          <div
            class="mb-4 p-3 bg-red-100 text-red-700 text-sm rounded-lg border border-red-200"
          >
            {resetErrorMsg}
          </div>
        {/if}

        <div class="space-y-4 mb-6">
          <div>
            <label
              class="block text-sm font-medium text-slate-700 mb-1"
              for="resetPassword">Password Baru</label
            >
            <div class="relative">
              <input
                id="resetPassword"
                type={showResetPassword ? "text" : "password"}
                bind:value={resetPassword}
                class="w-full pl-4 pr-12 py-2 border border-slate-300 rounded-xl focus:ring-2 focus:ring-amber-500 focus:border-amber-500 outline-none transition-shadow bg-white text-slate-900"
                placeholder="Masukkan password baru"
              />
              <button
                type="button"
                class="absolute inset-y-0 right-0 px-3 flex items-center text-slate-400 hover:text-amber-500 focus:outline-none transition-colors cursor-pointer"
                onclick={() => (showResetPassword = !showResetPassword)}
                title={showResetPassword
                  ? "Sembunyikan password"
                  : "Tampilkan password"}
              >
                {#if showResetPassword}
                  <!-- Eye slash icon -->
                  <svg
                    class="h-5 w-5"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"
                    />
                  </svg>
                {:else}
                  <!-- Eye icon -->
                  <svg
                    class="h-5 w-5"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                    />
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.543 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                    />
                  </svg>
                {/if}
              </button>
            </div>
          </div>
        </div>

        <div class="flex justify-end gap-3">
          <button
            onclick={cancelReset}
            disabled={isResetting}
            class="px-4 py-2 text-sm font-medium text-slate-800 bg-white shadow-md border border-slate-200 hover:bg-slate-100 rounded-xl transition-colors disabled:opacity-50"
          >
            Batal
          </button>
          <button
            onclick={confirmReset}
            disabled={isResetting}
            class="px-4 py-2 text-sm font-medium text-white bg-amber-600 shadow-md hover:bg-amber-700 rounded-xl transition-colors flex items-center gap-2 disabled:opacity-50"
          >
            {#if isResetting}
              <div
                class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"
              ></div>
            {/if}
            Reset Password
          </button>
        </div>
      </div>
    </div>
  {/if}

  <!-- Suspend Confirmation Modal -->
  {#if showSuspendModal && userToSuspend}
    <div
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-sm"
    >
      <div
        class="bg-slate-50 rounded-3xl p-6 w-full max-w-md shadow-2xl border border-slate-200"
      >
        <h3 class="text-xl font-bold text-slate-900 mb-2">
          Konfirmasi {userToSuspend.is_suspended ? "Aktivasi" : "Suspend"}
        </h3>
        <p class="text-slate-600 mb-6 text-sm">
          Apakah Anda yakin ingin {userToSuspend.is_suspended
            ? "mengaktifkan kembali"
            : "menangguhkan (suspend)"} user
          <span class="font-bold text-slate-900"
            >"{userToSuspend.username}"</span
          >?
          {#if !userToSuspend.is_suspended}
            <br /><span class="text-orange-600 font-medium block mt-1"
              >User tidak akan bisa login selama ditangguhkan.</span
            >
          {/if}
        </p>
        <div class="flex justify-end gap-3">
          <button
            onclick={cancelSuspend}
            disabled={isSuspending}
            class="px-4 py-2 text-sm font-medium text-slate-800 bg-white shadow-md border border-slate-200 hover:bg-slate-100 rounded-xl transition-colors disabled:opacity-50"
          >
            Batal
          </button>
          <button
            onclick={confirmSuspend}
            disabled={isSuspending}
            class="px-4 py-2 text-sm font-medium text-white {userToSuspend.is_suspended
              ? 'bg-green-600 hover:bg-green-700'
              : 'bg-orange-500 hover:bg-orange-600'} shadow-md rounded-xl transition-colors flex items-center gap-2 disabled:opacity-50"
          >
            {#if isSuspending}
              <div
                class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"
              ></div>
            {/if}
            Ya, {userToSuspend.is_suspended ? "Aktifkan" : "Suspend"}
          </button>
        </div>
      </div>
    </div>
  {/if}
</div>
