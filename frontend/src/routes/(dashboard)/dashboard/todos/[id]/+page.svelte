<script lang="ts">
  import { page } from "$app/stores";
  import { fade, slide } from "svelte/transition";
  import { flip } from "svelte/animate";
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";

  const listId = $page.params.id;

  type TodoItem = {
    id: number;
    todolist_id: number;
    text: string;
    completed: boolean;
  };

  type TodoList = {
    id: number;
    title: string;
    items: TodoItem[];
  };

  let list = $state<TodoList | null>(null);
  let newItemText = $state("");
  let isLoading = $state(true);
  let errorMsg = $state("");
  let openMenuId = $state<number | null>(null);

  onMount(async () => {
    await fetchList();
  });

  async function fetchList() {
    isLoading = true;
    errorMsg = "";
    try {
      const res = await fetch(`/api/todolists/${listId}`, { credentials: "include" });
      if (!res.ok) {
        throw new Error("Gagal mengambil data todolist");
      }
      list = await res.json();
      if (!list?.items) {
          list!.items = [];
      }
    } catch (e: any) {
      errorMsg = e.message;
    } finally {
      isLoading = false;
    }
  }

  async function addItem(e?: Event) {
    if (e) e.preventDefault();
    if (!newItemText.trim() || !list) return;

    try {
      const res = await fetch(`/api/todolists/${list.id}/items`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ text: newItemText.trim() }),
      });
      if (res.ok) {
        const newItem = await res.json();
        list.items = [...list.items, newItem];
        newItemText = "";
      } else {
        alert("Gagal menambahkan tugas");
      }
    } catch (err) {
      console.error(err);
      alert("Terjadi kesalahan");
    }
  }

  async function toggleItem(itemId: number, currentStatus: boolean) {
    if (!list) return;
    
    // Optimistic update
    list.items = list.items.map((i) => i.id === itemId ? { ...i, completed: !currentStatus } : i);

    try {
      const res = await fetch(`/api/todolists/${list.id}/items/${itemId}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ completed: !currentStatus }),
      });
      if (!res.ok) {
        throw new Error("Failed to update");
      }
    } catch (err) {
      console.error(err);
      // Revert on fail
      list.items = list.items.map((i) => i.id === itemId ? { ...i, completed: currentStatus } : i);
    }
  }

  async function deleteItem(itemId: number) {
    if (!list) return;

    // Optimistic update
    const originalItems = [...list.items];
    list.items = list.items.filter((i) => i.id !== itemId);

    try {
      const res = await fetch(`/api/todolists/${list.id}/items/${itemId}`, {
        method: "DELETE",
        credentials: "include",
      });
      if (!res.ok) {
        throw new Error("Failed to delete");
      }
    } catch (err) {
      console.error(err);
      list.items = originalItems;
    }
  }

  function goBack() {
    goto("/dashboard/todos");
  }
</script>

<svelte:head>
  <title>{list?.title || "Todolist"} - Portal Guru</title>
</svelte:head>

<div class="max-w-4xl mx-auto p-4 md:p-6 lg:p-8 animate-in fade-in duration-500">
  
  <button
    onclick={goBack}
    class="inline-flex items-center gap-2 text-slate-500 hover:text-blue-600 mb-6 font-medium transition-colors cursor-pointer"
  >
    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path></svg>
    Kembali ke Daftar List
  </button>

  {#if isLoading}
    <div class="flex items-center justify-center p-12">
      <div class="w-10 h-10 border-4 border-slate-200 border-t-blue-600 rounded-full animate-spin"></div>
    </div>
  {:else if errorMsg}
    <div class="p-6 bg-red-50 text-red-600 rounded-xl border border-red-200 flex items-center gap-3">
      <svg class="w-6 h-6 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
      <p class="font-semibold">{errorMsg}</p>
    </div>
  {:else if list}
    <div class="bg-white/80 backdrop-blur-xl border border-slate-200/60 shadow-lg shadow-slate-200/40 rounded-3xl p-6 md:p-10 transition-all">
      <h1 class="text-2xl md:text-4xl font-bold text-slate-900 mb-6 md:mb-8 pb-4 md:pb-6 border-b border-slate-200">
        {list.title}
      </h1>

      <form onsubmit={addItem} class="flex items-center gap-3 mb-10">
        <input
          type="text"
          bind:value={newItemText}
          placeholder="Tambahkan tugas baru..."
          class="w-full px-4 py-3 md:px-5 md:py-4 border-2 border-slate-200 rounded-2xl focus:ring-0 focus:border-blue-500 outline-none transition-all bg-slate-50 hover:bg-white text-slate-800 text-base md:text-lg shadow-sm"
        />
        <button
          type="submit"
          disabled={!newItemText.trim()}
          class="bg-blue-600 hover:bg-blue-700 disabled:bg-slate-300 text-white font-semibold p-3 md:p-4 rounded-2xl transition-all shadow-lg shadow-blue-500/30 hover:shadow-blue-500/50 shrink-0 cursor-pointer"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 4v16m8-8H4"></path></svg>
        </button>
      </form>

      <div class="flex flex-col gap-3 min-h-[400px]">
        {#if list.items.length === 0}
          <div class="p-16 text-center text-slate-500 flex flex-col items-center justify-center h-full opacity-70">
            <svg class="w-20 h-20 mb-6 text-slate-300" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4"></path></svg>
            <p class="font-medium text-xl text-slate-600 mb-2">Belum ada tugas.</p>
            <p class="text-base">Ketik di atas lalu tekan enter untuk menambah tugas ke list ini.</p>
          </div>
        {:else}
          {#each list.items as item (item.id)}
            <div
              animate:flip={{ duration: 300 }}
              in:slide|local
              out:slide|local
              class="group relative flex items-center justify-between p-4 md:p-5 bg-white border {item.completed ? 'border-slate-100 bg-slate-50' : 'border-slate-200 shadow-sm'} rounded-2xl transition-all hover:border-blue-300 hover:shadow-md {openMenuId === item.id ? 'z-50' : 'z-0'}"
            >
              <div
                role="button"
                tabindex="0"
                class="flex items-center gap-5 flex-1 cursor-pointer outline-none"
                onclick={() => toggleItem(item.id, item.completed)}
                onkeydown={(e) => {
                  if (e.key === "Enter" || e.key === " ") toggleItem(item.id, item.completed);
                }}
              >
                <div class="relative flex items-center justify-center w-7 h-7 md:w-8 md:h-8 shrink-0">
                  <input
                    type="checkbox"
                    checked={item.completed}
                    tabindex="-1"
                    class="peer appearance-none w-6 h-6 md:w-7 md:h-7 border-2 border-slate-300 rounded-xl checked:bg-blue-500 checked:border-blue-500 transition-all cursor-pointer pointer-events-none"
                  />
                  <svg class="absolute w-4 h-4 md:w-5 md:h-5 text-white pointer-events-none opacity-0 peer-checked:opacity-100 transition-opacity" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7"></path></svg>
                </div>
                <span class="text-base md:text-lg transition-all duration-300 {item.completed ? 'text-slate-400 line-through decoration-slate-400 decoration-2' : 'text-slate-700 font-medium'}">
                  {item.text}
                </span>
              </div>

              <div class="relative shrink-0 flex items-center gap-1 z-10">
                <!-- Mobile 3-dots button -->
                <button
                  onclick={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                    openMenuId = openMenuId === item.id ? null : item.id;
                  }}
                  class="md:hidden p-2 text-slate-400 hover:bg-slate-100 rounded-xl transition-colors cursor-pointer"
                >
                  <svg class="w-6 h-6" fill="currentColor" viewBox="0 0 24 24"><path d="M12 8a2 2 0 110-4 2 2 0 010 4zm0 6a2 2 0 110-4 2 2 0 010 4zm0 6a2 2 0 110-4 2 2 0 010 4z"/></svg>
                </button>

                <button
                  onclick={() => { openMenuId = null; deleteItem(item.id); }}
                  class="absolute md:static right-0 top-12 md:top-auto md:right-auto {openMenuId === item.id ? 'flex' : 'hidden md:flex'} md:opacity-0 md:group-hover:opacity-100 items-center gap-2 p-3 md:p-2 bg-white md:bg-transparent border md:border-transparent border-slate-100 shadow-xl md:shadow-none text-red-600 md:text-slate-400 hover:text-red-500 hover:bg-red-50 rounded-xl transition-all cursor-pointer w-max"
                  title="Hapus tugas"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
                  <span class="md:hidden text-sm font-medium pr-2">Hapus</span>
                </button>
              </div>
            </div>
          {/each}
        {/if}
      </div>
    </div>
  {/if}
</div>
