<script lang="ts">
  import { fade, scale } from "svelte/transition";
  import { onMount } from "svelte";

  type TodoList = {
    id: number;
    title: string;
  };

  let lists = $state<TodoList[]>([]);
  let newListTitle = $state("");
  let isLoading = $state(true);
  let errorMsg = $state("");

  let showDeleteListModal = $state(false);
  let listToDelete = $state<TodoList | null>(null);
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
</script>

<svelte:head>
  <title>Todolist - Portal Guru</title>
</svelte:head>

<div
  class="max-w-7xl mx-auto p-4 md:p-6 lg:p-8 animate-in fade-in duration-500"
>
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
    <div
      class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6"
    >
      {#each lists as list (list.id)}
        <a
          href={`/dashboard/todos/${list.id}`}
          class="group block bg-white border border-slate-200 rounded-3xl p-6 shadow-sm hover:shadow-md hover:border-blue-300 transition-all duration-300 relative overflow-hidden"
        >
          <div class="flex items-start justify-between gap-4">
            <h3
              class="text-xl font-bold text-slate-800 group-hover:text-blue-600 transition-colors line-clamp-2"
            >
              {list.title}
            </h3>

            <div class="relative shrink-0">
              <!-- Mobile 3-dots button -->
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

              <button
                onclick={(e) => {
                  openMenuId = null;
                  promptDeleteList(list, e);
                }}
                class="absolute md:static right-0 top-12 md:top-auto md:right-auto {openMenuId ===
                list.id
                  ? 'flex'
                  : 'hidden md:flex'} md:opacity-0 md:group-hover:opacity-100 items-center gap-2 p-3 md:p-2 bg-white md:bg-transparent border md:border-transparent border-slate-100 shadow-xl md:shadow-none text-red-600 md:text-slate-400 hover:text-red-500 hover:bg-red-50 rounded-xl transition-all cursor-pointer z-10"
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
          <div
            class="mt-8 flex items-center justify-between text-sm text-slate-500 font-medium"
          >
            <span>Buka List</span>
            <svg
              class="w-4 h-4 transform group-hover:translate-x-1 transition-transform"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              ><path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M9 5l7 7-7 7"
              ></path></svg
            >
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
          Ya, Hapus
        </button>
      </div>
    </div>
  </div>
{/if}
