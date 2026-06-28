<script lang="ts">
  import { fade, scale } from "svelte/transition";
  import { onMount } from "svelte";

  type TodoList = {
    id: number;
    title: string;
    created_at: string;
  };

  let lists = $state<TodoList[]>([]);
  let newListTitle = $state("");
  let isLoading = $state(true);
  let errorMsg = $state("");

  let showDeleteListModal = $state(false);
  let listToDelete = $state<TodoList | null>(null);
  let showEditListModal = $state(false);
  let listToEdit = $state<TodoList | null>(null);
  let editListTitle = $state("");
  let openMenuId = $state<number | null>(null);

  onMount(async () => {
    await fetchLists();
  });

  async function fetchLists() {
    isLoading = true;
    errorMsg = "";
    try {
      const res = await fetch("/api/todolists", { credentials: "include" });
      if (!res.ok) {
        if (res.status === 403) {
          throw new Error("Akses ditolak. Anda bukan Guru.");
        }
        throw new Error("Gagal mengambil data todolist");
      }
      lists = (await res.json()) || [];
    } catch (e: any) {
      errorMsg = e.message;
    } finally {
      isLoading = false;
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
        body: JSON.stringify({ title: newListTitle.trim() }),
      });
      if (res.ok) {
        const newList = await res.json();
        lists = [...lists, newList];
        newListTitle = "";
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
        lists = lists.filter((l) => l.id !== listToDelete!.id);
        showDeleteListModal = false;
        listToDelete = null;
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
    showEditListModal = true;
  }

  async function confirmEditList() {
    if (!listToEdit || !editListTitle.trim()) return;

    try {
      const res = await fetch(`/api/todolists/${listToEdit.id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ title: editListTitle.trim() }),
      });
      if (res.ok) {
        const updatedList = await res.json();
        lists = lists.map((l) => (l.id === updatedList.id ? updatedList : l));
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

<svelte:head>
  <title>Todolist - Portal Guru</title>
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
    <form onsubmit={addList} class="w-full md:w-auto flex gap-2">
      <input
        type="text"
        bind:value={newListTitle}
        placeholder="Nama list baru..."
        class="w-full md:w-64 px-4 py-2 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all bg-white text-slate-900 shadow-sm"
      />
      <button
        type="submit"
        disabled={!newListTitle.trim()}
        class="bg-blue-600 hover:bg-blue-700 disabled:bg-slate-300 text-white font-medium px-4 py-2 rounded-xl transition-all shadow-md shrink-0 cursor-pointer"
      >
        Buat List
      </button>
    </form>
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
    <div class="flex flex-col gap-4">
      {#each lists as list (list.id)}
        <a
          href={`/dashboard/todos/${list.id}`}
          class="group relative flex items-center justify-between bg-white border border-slate-200 rounded-2xl px-4 py-2 shadow-sm hover:shadow-md hover:border-blue-300 transition-all duration-300 {openMenuId ===
          list.id
            ? 'z-50'
            : 'z-0'}"
        >
          <div class="flex items-center gap-4 pr-4">
            <div class=" text-blue-500 rounded-xl sm:block shrink-0">
              <svg
                class="w-6 h-6"
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
            <div class="flex flex-col md:flex-row md:items-center gap-1.5 md:gap-3 mt-1 md:mt-0">
              <h3
                class="text-lg md:text-xl font-semibold text-slate-800 group-hover:text-blue-600 transition-colors break-words leading-tight"
              >
                {list.title}
              </h3>
              <div class="flex items-center text-[11px] md:text-xs text-slate-500 font-medium bg-slate-100 px-2.5 py-1 rounded-md w-max">
                <svg class="w-3.5 h-3.5 mr-1 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg>
                {new Date(list.created_at).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })}
              </div>
            </div>
          </div>

          <div class="relative shrink-0 ml-auto">
            <button
              onclick={(e) => {
                e.preventDefault();
                e.stopPropagation();
                openMenuId = openMenuId === list.id ? null : list.id;
              }}
              class="md:hidden p-2 text-slate-400 hover:bg-slate-100 rounded-xl transition-colors cursor-pointer"
            >
              <svg class="w-6 h-6" fill="currentColor" viewBox="0 0 24 24"
                ><path
                  d="M12 8a2 2 0 110-4 2 2 0 010 4zm0 6a2 2 0 110-4 2 2 0 010 4zm0 6a2 2 0 110-4 2 2 0 010 4z"
                /></svg
              >
            </button>

            <div
              class="absolute md:static right-0 top-12 md:top-auto md:right-auto {openMenuId ===
              list.id
                ? 'flex flex-col'
                : 'hidden md:flex'} md:flex-row p-2 md:p-0 bg-white md:bg-transparent border md:border-none border-slate-100 shadow-xl md:shadow-none rounded-xl md:opacity-0 md:group-hover:opacity-100 items-center gap-1 transition-all z-10 w-max"
            >
              <button
                onclick={(e) => {
                  openMenuId = null;
                  promptEditList(list, e);
                }}
                class="w-full flex items-center gap-2 p-3 md:p-2 bg-white md:bg-transparent text-slate-600 md:text-slate-400 hover:text-blue-600 hover:bg-blue-50 rounded-xl transition-all cursor-pointer"
                title="Edit List"
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
                <span class="md:hidden text-sm font-medium pr-2">Edit</span>
              </button>

              <button
                onclick={(e) => {
                  openMenuId = null;
                  promptDeleteList(list, e);
                }}
                class="w-full flex items-center gap-2 p-3 md:p-2 bg-white md:bg-transparent text-red-600 md:text-slate-400 hover:text-red-500 hover:bg-red-50 rounded-xl transition-all cursor-pointer"
                title="Hapus List"
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
                    d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                  ></path></svg
                >
                <span class="md:hidden text-sm font-medium pr-2">Hapus</span>
              </button>
            </div>
          </div>
        </a>
      {/each}
    </div>
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
        class="w-full px-4 py-2 mb-6 border border-slate-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all bg-white text-slate-900 shadow-sm"
        onkeydown={(e) => e.key === "Enter" && confirmEditList()}
      />
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
