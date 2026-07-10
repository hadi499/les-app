<script lang="ts">
  import { fade, scale } from "svelte/transition";
  import { onMount } from "svelte";

  type TodoItem = {
    id: number;
    todolist_id: number;
    text: string;
    completed: boolean;
  };

  type TodoList = {
    id: number;
    title: string;
    student_username?: string;
    created_at: string;
    items?: TodoItem[];
  };

  let lists = $state<TodoList[]>([]);
  let newListTitle = $state("");
  let newListStudentUsername = $state("");
  let isLoading = $state(true);
  let errorMsg = $state("");

  let showCreateListModal = $state(false);
  let showDeleteListModal = $state(false);
  let listToDelete = $state<TodoList | null>(null);
  let showEditListModal = $state(false);
  let listToEdit = $state<TodoList | null>(null);
  let editListTitle = $state("");
  let editListStudentUsername = $state("");
  let openMenuId = $state<number | null>(null);
  let userRole = $state("");
  let students = $state<{id: number, username: string, role: string}[]>([]);

  // Paginasi state
  let currentPage = $state(1);
  let itemsPerPage = 24;
  let totalPages = $state(1);
  let totalItems = $state(0);
  let isInitialLoad = $state(true);

  // Ambil ulang data ketika currentPage berubah (tapi bukan saat load pertama)
  $effect(() => {
    if (!isInitialLoad) {
      fetchLists();
    }
  });

  onMount(async () => {
    try {
      const meRes = await fetch("/me");
      if (meRes.ok) {
        const data = await meRes.json();
        userRole = data.user?.role || "";
      }

      if (userRole === "teacher") {
        const usersRes = await fetch("/api/users", { credentials: "include" });
        if (usersRes.ok) {
          const data = await usersRes.json();
          students = data.users?.filter((u: any) => u.role === "student") || [];
        }
      }
    } catch (e) {
      console.error(e);
    }
    await fetchLists();
  });

  async function fetchLists() {
    isLoading = true;
    errorMsg = "";
    try {
      const res = await fetch(`/api/todolists?page=${currentPage}&limit=${itemsPerPage}`, { credentials: "include" });
      if (!res.ok) {
        if (res.status === 403) {
          throw new Error("Akses ditolak. Anda bukan Guru.");
        }
        throw new Error("Gagal mengambil data todolist");
      }
      const data = await res.json();
      if (data && data.data) {
        lists = data.data;
        totalPages = data.total_pages;
        currentPage = data.current_page;
        totalItems = data.total_items;
      } else {
        lists = [];
        totalPages = 1;
        totalItems = 0;
      }
    } catch (e: any) {
      errorMsg = e.message;
    } finally {
      isLoading = false;
      isInitialLoad = false;
    }
  }

  async function addList(e?: Event) {
    if (e) e.preventDefault();
    if (!newListTitle.trim()) return;

    try {
      const res = await fetch("/api/todolists", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ 
          title: newListTitle.trim(),
          student_username: newListStudentUsername.trim()
        }),
      });
      if (res.ok) {
        newListTitle = "";
        newListStudentUsername = "";
        showCreateListModal = false;
        // Pindah ke halaman pertama dan fetch ulang untuk melihat list baru
        if (currentPage === 1) {
          fetchLists();
        } else {
          currentPage = 1;
        }
      } else {
        const err = await res.json();
        alert(err.error || "Gagal membuat list");
      }
    } catch (err) {
      console.error(err);
      alert("Terjadi kesalahan sistem");
    }
  }

  function promptDeleteList(list: TodoList, e: Event) {
    e.preventDefault();
    e.stopPropagation();
    listToDelete = list;
    showDeleteListModal = true;
  }

  async function confirmDeleteList() {
    if (!listToDelete) return;

    try {
      const res = await fetch(`/api/todolists/${listToDelete.id}`, {
        method: "DELETE",
        credentials: "include",
      });
      if (res.ok) {
        showDeleteListModal = false;
        listToDelete = null;

        // Jika data di halaman ini habis, mundur 1 halaman, jika tidak cukup fetch ulang
        if (lists.length === 1 && currentPage > 1) {
          currentPage--;
        } else {
          fetchLists();
        }
      } else {
        alert("Gagal menghapus list");
      }
    } catch (err) {
      console.error(err);
      alert("Terjadi kesalahan");
    }
  }

  function promptEditList(list: TodoList, e: Event) {
    e.preventDefault();
    e.stopPropagation();
    listToEdit = list;
    editListTitle = list.title;
    editListStudentUsername = list.student_username || "";
    showEditListModal = true;
  }

  async function confirmEditList() {
    if (!listToEdit || !editListTitle.trim()) return;

    try {
      const res = await fetch(`/api/todolists/${listToEdit.id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ 
          title: editListTitle.trim(),
          student_username: editListStudentUsername.trim()
        }),
      });
      if (res.ok) {
        const updatedList = await res.json();
        lists = lists.map((l) => (l.id === updatedList.id ? { ...updatedList, items: l.items } : l));
        showEditListModal = false;
        listToEdit = null;
      } else {
        const err = await res.json();
        alert(err.error || "Gagal mengedit list");
      }
    } catch (err) {
      console.error(err);
      alert("Terjadi kesalahan sistem");
    }
  }
</script>

<svelte:window onclick={() => openMenuId = null} />

<svelte:head>
  <title>Todolist - Les App</title>
</svelte:head>

<div class="max-w-4xl mx-auto md:p-4 lg:p-8 animate-in fade-in duration-500">
  <div
    class="mb-8 flex flex-col md:flex-row justify-between items-start md:items-end gap-4"
  >
    <div>
      <h1
        class="text-2xl font-bold text-slate-900 sm:text-3xl tracking-tight drop-shadow-sm"
      >
        Todolist
      </h1>
      <p
        class="mt-2 text-slate-600 text-sm sm:text-base font-light tracking-wide"
      >
        Kelola dan kategorikan daftar tugas Anda.
      </p>
    </div>
    {#if userRole === "teacher"}
      <button
        onclick={() => (showCreateListModal = true)}
        class="text-blue-700 bg-blue-50 hover:bg-blue-100 border border-blue-200 font-bold px-4 md:px-5 py-2.5 rounded-xl transition-all shadow-sm flex items-center gap-2 cursor-pointer w-max justify-center"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
        <span>Buat Todolist</span>
      </button>
    {/if}
  </div>

  {#if isLoading}
    <div class="flex items-center justify-center p-12">
      <div
        class="w-10 h-10 border-4 border-slate-200 border-t-blue-600 rounded-full animate-spin"
      ></div>
    </div>
  {:else if errorMsg}
    <div
      class="p-6 bg-red-50 text-red-600 rounded-xl border border-red-200 flex items-center gap-3"
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
      <p class="font-semibold">{errorMsg}</p>
    </div>
  {:else if lists.length === 0}
    <div
      class="bg-white/60 backdrop-blur-md rounded-3xl border border-slate-200 shadow-sm p-16 text-center"
    >
      <svg
        class="w-20 h-20 mx-auto text-slate-300 mb-6"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
        ><path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="1.5"
          d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 002-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
        ></path></svg
      >
      <h3 class="text-2xl font-bold text-slate-700 mb-2">Belum Ada Todolist</h3>
      <p class="text-slate-500">
        Buat list pertama Anda dengan mengetikkan nama di pojok kanan atas.
      </p>
    </div>
  {:else}
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-6">
      {#each lists as list (list.id)}
        <div class="bg-white border border-slate-200 rounded-2xl p-5 shadow-sm hover:shadow-md transition-shadow relative flex flex-col h-full {openMenuId === list.id ? 'z-50' : 'z-0'}">
          <div class="flex justify-between items-start mb-2">
            <div>
              <span class="text-xs font-bold text-slate-400">#{list.id}</span>
              <h3 class="text-lg font-bold text-slate-900 mt-1 leading-tight">{list.title}</h3>
            </div>
            {#if list.student_username}
              <span class="px-2 py-1 bg-blue-50 text-blue-600 rounded-md text-[10px] font-bold border border-blue-100 whitespace-nowrap uppercase tracking-wide">
                @{list.student_username}
              </span>
            {/if}
          </div>
          
          <div class="flex items-center text-xs font-medium text-slate-500 mb-4 bg-slate-50 w-max px-2.5 py-1.5 rounded-md border border-slate-100">
            <svg class="w-4 h-4 mr-1.5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg>
            {new Date(list.created_at).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })}
          </div>

          {#if list.items && list.items.length > 0 && list.items.every(i => i.completed)}
            <div class="mb-4">
              <span class="px-2 py-1 text-xs font-bold bg-emerald-50 text-emerald-600 rounded-md border border-emerald-100 flex items-center gap-1.5 w-max">
                <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7"></path></svg>
                Selesai
              </span>
            </div>
          {/if}

          <div class="mt-auto flex gap-2 pt-4 border-t border-slate-100">
            <a
              href={`/dashboard/todos/${list.id}`}
              class="flex-1 inline-flex items-center justify-center p-2.5 text-sm font-semibold text-blue-700 bg-blue-50 rounded-xl border border-blue-100 hover:bg-blue-100 transition-colors no-underline"
              title="Buka Todo"
            >
              <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4"></path></svg>
              Buka
            </a>
            
            {#if userRole === "teacher"}
              <div class="relative shrink-0">
                <button
                  onclick={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                    openMenuId = openMenuId === list.id ? null : list.id;
                  }}
                  class="p-2.5 text-slate-500 bg-slate-50 border border-slate-200 hover:bg-slate-100 rounded-xl transition-colors cursor-pointer flex items-center justify-center h-full"
                  title="Menu"
                >
                  <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24"
                    ><path
                      d="M12 8a2 2 0 110-4 2 2 0 010 4zm0 6a2 2 0 110-4 2 2 0 010 4zm0 6a2 2 0 110-4 2 2 0 010 4z"
                    /></svg
                  >
                </button>

                <div
                  class="absolute right-0 bottom-14 {openMenuId === list.id
                    ? 'flex flex-col'
                    : 'hidden'} p-2 bg-white border border-slate-100 shadow-xl rounded-xl items-center gap-1 transition-all z-50 w-max"
                >
                  <button
                    onclick={(e) => {
                      openMenuId = null;
                      promptEditList(list, e);
                    }}
                    class="w-full flex items-center gap-2 p-2 text-slate-600 hover:text-blue-600 hover:bg-blue-50 rounded-xl transition-all cursor-pointer"
                    title="Edit List"
                  >
                    <svg class="w-5 h-5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path></svg>
                    <span class="text-sm font-medium pr-2">Edit</span>
                  </button>

                  <button
                    onclick={(e) => {
                      openMenuId = null;
                      promptDeleteList(list, e);
                    }}
                    class="w-full flex items-center gap-2 p-2 text-red-600 hover:text-red-500 hover:bg-red-50 rounded-xl transition-all cursor-pointer"
                    title="Hapus List"
                  >
                    <svg class="w-5 h-5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
                    <span class="text-sm font-medium pr-2">Hapus</span>
                  </button>
                </div>
              </div>
            {/if}
          </div>
        </div>
      {/each}
    </div>

    <!-- Paginasi (Biasa) -->
    {#if totalPages > 1}
      <div class="flex items-center justify-between px-4 py-3 bg-white border border-slate-200 rounded-2xl shadow-sm mt-4 mb-4">
        <div class="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
          <div>
            <p class="text-sm text-slate-700 m-0">
              Menampilkan <span class="font-medium">{totalItems === 0 ? 0 : (currentPage - 1) * itemsPerPage + 1}</span> 
              hingga <span class="font-medium">{Math.min(currentPage * itemsPerPage, totalItems)}</span> 
              dari <span class="font-medium">{totalItems}</span> todolist
            </p>
          </div>
          <div>
            <nav class="isolate inline-flex -space-x-px rounded-md shadow-sm" aria-label="Pagination">
              <button
                onclick={() => currentPage = Math.max(1, currentPage - 1)}
                disabled={currentPage === 1}
                class="relative inline-flex items-center rounded-l-md px-2 py-2 text-slate-400 ring-1 ring-inset ring-slate-300 hover:bg-slate-50 focus:z-20 focus:outline-offset-0 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
              >
                <span class="sr-only">Previous</span>
                <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true"><path fill-rule="evenodd" d="M12.79 5.23a.75.75 0 01-.02 1.06L8.832 10l3.938 3.71a.75.75 0 11-1.04 1.08l-4.5-4.25a.75.75 0 010-1.08l4.5-4.25a.75.75 0 011.06.02z" clip-rule="evenodd" /></svg>
              </button>
              {#each Array(totalPages) as _, i}
                <button
                  onclick={() => currentPage = i + 1}
                  class="relative inline-flex items-center px-4 py-2 text-sm font-semibold {currentPage === i + 1 ? 'z-10 bg-blue-600 text-white focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600' : 'text-slate-900 ring-1 ring-inset ring-slate-300 hover:bg-slate-50 focus:z-20 focus:outline-offset-0'} cursor-pointer"
                >
                  {i + 1}
                </button>
              {/each}
              <button
                onclick={() => currentPage = Math.min(totalPages, currentPage + 1)}
                disabled={currentPage === totalPages}
                class="relative inline-flex items-center rounded-r-md px-2 py-2 text-slate-400 ring-1 ring-inset ring-slate-300 hover:bg-slate-50 focus:z-20 focus:outline-offset-0 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
              >
                <span class="sr-only">Next</span>
                <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true"><path fill-rule="evenodd" d="M7.21 14.77a.75.75 0 01.02-1.06L11.168 10 7.23 6.29a.75.75 0 111.04-1.08l4.5 4.25a.75.75 0 010 1.08l-4.5 4.25a.75.75 0 01-1.06-.02z" clip-rule="evenodd" /></svg>
              </button>
            </nav>
          </div>
        </div>
        <!-- Mobile pagination view -->
        <div class="flex flex-1 justify-center gap-4 sm:hidden">
          <button
            onclick={() => currentPage = Math.max(1, currentPage - 1)}
            disabled={currentPage === 1}
            class="relative inline-flex items-center justify-center w-10 h-10 rounded-full bg-white border border-slate-200 text-slate-500 hover:bg-slate-50 hover:text-blue-600 disabled:opacity-50 disabled:pointer-events-none transition-colors shadow-sm cursor-pointer"
          >
            <span class="sr-only">Previous</span>
            <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true"><path fill-rule="evenodd" d="M12.79 5.23a.75.75 0 01-.02 1.06L8.832 10l3.938 3.71a.75.75 0 11-1.04 1.08l-4.5-4.25a.75.75 0 010-1.08l4.5-4.25a.75.75 0 011.06.02z" clip-rule="evenodd" /></svg>
          </button>
          <span class="flex items-center text-sm font-semibold text-slate-700">
            {currentPage} <span class="text-slate-400 mx-1.5">/</span> {totalPages}
          </span>
          <button
            onclick={() => currentPage = Math.min(totalPages, currentPage + 1)}
            disabled={currentPage === totalPages}
            class="relative inline-flex items-center justify-center w-10 h-10 rounded-full bg-white border border-slate-200 text-slate-500 hover:bg-slate-50 hover:text-blue-600 disabled:opacity-50 disabled:pointer-events-none transition-colors shadow-sm cursor-pointer"
          >
            <span class="sr-only">Next</span>
            <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true"><path fill-rule="evenodd" d="M7.21 14.77a.75.75 0 01.02-1.06L11.168 10 7.23 6.29a.75.75 0 111.04-1.08l4.5 4.25a.75.75 0 010 1.08l-4.5 4.25a.75.75 0 01-1.06-.02z" clip-rule="evenodd" /></svg>
          </button>
        </div>
      </div>
    {/if}
  {/if}
</div>

{#if showDeleteListModal && listToDelete}
  <div
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-sm"
    transition:fade={{ duration: 150 }}
  >
    <div
      class="bg-white rounded-3xl shadow-xl max-w-sm w-full p-6 border border-slate-100"
      transition:scale={{ duration: 150, start: 0.95 }}
    >
      <h3 class="text-xl font-bold text-slate-900 mb-2">Hapus List?</h3>
      <p class="text-slate-600 text-sm mb-6">
        Anda yakin ingin menghapus list <span class="font-bold text-slate-900"
          >"{listToDelete.title}"</span
        > beserta semua tugas di dalamnya?
      </p>
      <div class="flex justify-end gap-3">
        <button
          onclick={() => (showDeleteListModal = false)}
          class="px-4 py-2 text-sm font-medium text-slate-800 bg-white shadow-sm border border-slate-200 hover:bg-slate-50 rounded-xl transition-colors cursor-pointer"
        >
          Batal
        </button>
        <button
          onclick={confirmDeleteList}
          class="px-4 py-2 text-sm font-medium text-white bg-red-500 shadow-sm hover:bg-red-600 rounded-xl transition-colors cursor-pointer"
        >
          Hapus
        </button>
      </div>
    </div>
  </div>
{/if}

{#if showEditListModal && listToEdit}
  <div
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-sm"
    transition:fade={{ duration: 150 }}
  >
    <div
      class="bg-white rounded-3xl shadow-xl max-w-sm w-full p-6 border border-slate-100"
      transition:scale={{ duration: 150, start: 0.95 }}
    >
      <h3 class="text-xl font-bold text-slate-900 mb-4">Edit Todolist</h3>
      <input
        type="text"
        bind:value={editListTitle}
        placeholder="Nama Todolist"
        class="w-full px-4 py-2 mb-3 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all bg-white text-slate-900 shadow-sm"
        onkeydown={(e) => e.key === "Enter" && confirmEditList()}
      />
      <select
        bind:value={editListStudentUsername}
        class="w-full px-4 pr-10 py-2 mb-6 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all bg-white text-slate-900 shadow-sm appearance-none bg-[url('data:image/svg+xml;charset=US-ASCII,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20width%3D%22292.4%22%20height%3D%22292.4%22%3E%3Cpath%20fill%3D%22%2364748b%22%20d%3D%22M287%2069.4a17.6%2017.6%200%200%200-13-5.4H18.4c-5%200-9.3%201.8-12.9%205.4A17.6%2017.6%200%200%200%200%2082.2c0%205%201.8%209.3%205.4%2012.9l128%20127.9c3.6%203.6%207.8%205.4%2012.8%205.4s9.2-1.8%2012.8-5.4L287%2095c3.5-3.5%205.4-7.8%205.4-12.8%200-5-1.9-9.2-5.5-12.8z%22%2F%3E%3C%2Fsvg%3E')] bg-[length:0.7rem_auto] bg-[right_1rem_center] bg-no-repeat"
      >
        <option value="">Pilih Murid (opsional)</option>
        {#each students as student}
          <option value={student.username}>@{student.username}</option>
        {/each}
      </select>
      <div class="flex justify-end gap-3">
        <button
          onclick={() => {
            showEditListModal = false;
            listToEdit = null;
          }}
          class="px-5 py-2.5 text-sm font-medium text-slate-600 bg-white shadow-sm border border-slate-200 hover:bg-slate-50 hover:text-slate-900 rounded-xl transition-all cursor-pointer"
        >
          Batal
        </button>
        <button
          onclick={confirmEditList}
          disabled={!editListTitle.trim()}
          class="px-5 py-2.5 text-sm font-medium text-white bg-blue-600 shadow-sm shadow-blue-500/30 hover:bg-blue-700 hover:shadow-blue-500/50 rounded-xl transition-all cursor-pointer disabled:opacity-50"
        >
          Simpan
        </button>
      </div>
    </div>
  </div>
{/if}

{#if showCreateListModal}
  <div
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-sm"
    transition:fade={{ duration: 150 }}
  >
    <div
      class="bg-white rounded-3xl shadow-xl max-w-sm w-full p-6 border border-slate-100"
      transition:scale={{ duration: 150, start: 0.95 }}
    >
      <h3 class="text-xl font-bold text-slate-900 mb-4">Buat Todolist Baru</h3>
      
      <div class="mb-4">
        <label class="block text-sm font-medium text-slate-700 mb-1" for="new_title">Nama Todolist</label>
        <input
          id="new_title"
          type="text"
          bind:value={newListTitle}
          placeholder="Contoh: Tugas Matematika"
          class="w-full px-4 py-2 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all bg-white text-slate-900 shadow-sm"
          onkeydown={(e) => e.key === "Enter" && addList()}
        />
      </div>

      <div class="mb-6">
        <label class="block text-sm font-medium text-slate-700 mb-1" for="new_student">Pilih Murid (Opsional)</label>
        <select
          id="new_student"
          bind:value={newListStudentUsername}
          class="w-full px-4 pr-10 py-2 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all bg-white text-slate-900 shadow-sm appearance-none bg-[url('data:image/svg+xml;charset=US-ASCII,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20width%3D%22292.4%22%20height%3D%22292.4%22%3E%3Cpath%20fill%3D%22%2364748b%22%20d%3D%22M287%2069.4a17.6%2017.6%200%200%200-13-5.4H18.4c-5%200-9.3%201.8-12.9%205.4A17.6%2017.6%200%200%200%200%2082.2c0%205%201.8%209.3%205.4%2012.9l128%20127.9c3.6%203.6%207.8%205.4%2012.8%205.4s9.2-1.8%2012.8-5.4L287%2095c3.5-3.5%205.4-7.8%205.4-12.8%200-5-1.9-9.2-5.5-12.8z%22%2F%3E%3C%2Fsvg%3E')] bg-[length:0.7rem_auto] bg-[right_1rem_center] bg-no-repeat"
        >
          <option value="">- Tidak ada (Semua Murid) -</option>
          {#each students as student}
            <option value={student.username}>@{student.username}</option>
          {/each}
        </select>
      </div>

      <div class="flex justify-end gap-3">
        <button
          onclick={() => {
            showCreateListModal = false;
            newListTitle = "";
            newListStudentUsername = "";
          }}
          class="px-5 py-2.5 text-sm font-medium text-slate-600 bg-white shadow-sm border border-slate-200 hover:bg-slate-50 hover:text-slate-900 rounded-xl transition-all cursor-pointer"
        >
          Batal
        </button>
        <button
          onclick={addList}
          disabled={!newListTitle.trim()}
          class="px-5 py-2.5 text-sm font-medium text-white bg-blue-600 shadow-sm shadow-blue-500/30 hover:bg-blue-700 hover:shadow-blue-500/50 rounded-xl transition-all cursor-pointer disabled:opacity-50"
        >
          Buat
        </button>
      </div>
    </div>
  </div>
{/if}
