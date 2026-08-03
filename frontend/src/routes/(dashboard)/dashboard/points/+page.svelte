<script lang="ts">
  import { onMount } from "svelte";
  import { toast } from "$lib/stores/toast.svelte";

  type User = { id: number; username: string; role: string; points?: number };

  let users: User[] = $state([]);
	let isLoading = $state(true);
	let showLoadingSpinner = $state(false);
	let errorMsg = $state("");
	let isTeacher = $state(false);
	let isResetting = $state(false);
	let showConfirmModal = $state(false);

	async function checkRole() {
		try {
			const res = await fetch("/me", { credentials: "include" });
			if (res.ok) {
				const data = await res.json();
				if (data.authenticated && data.user && data.user.role === "teacher") {
					isTeacher = true;
				}
			}
		} catch (e) {
			console.error("Gagal memeriksa role:", e);
		}
	}

	function resetPoints() {
		showConfirmModal = true;
	}

	async function executeResetPoints() {
		isResetting = true;
		errorMsg = "";
		try {
			const res = await fetch("/api/users/reset-points", {
				method: "POST",
				credentials: "include",
			});
			if (!res.ok) {
				throw new Error("Gagal mereset poin");
			}
			await fetchUsers(); // Refresh the list
			toast.success("Berhasil mereset semua poin dan riwayat nilai kuis.");
			showConfirmModal = false;
		} catch (e) {
			errorMsg = e instanceof Error ? e.message : String(e);
		} finally {
			isResetting = false;
		}
	}

  async function fetchUsers() {
    isLoading = true;
    showLoadingSpinner = false;
    setTimeout(() => { showLoadingSpinner = true; }, 150);
    errorMsg = "";
    try {
      const res = await fetch(`/api/users`, {
        credentials: "include",
      });
      if (res.status === 403) {
        errorMsg = "Akses ditolak.";
        isLoading = false;
        return;
      }
      if (!res.ok) {
        throw new Error("Gagal mengambil data users");
      }
      const data = (await res.json()) as { users: User[] };
      // Filter out teachers and sort users by points descending
      users = (data.users || [])
        .filter((u) => u.role !== "teacher")
        .sort((a, b) => (b.points || 0) - (a.points || 0));
    } catch (e) {
      errorMsg = e instanceof Error ? e.message : String(e);
    } finally {
      isLoading = false;
    }
  }

  onMount(async () => {
    await checkRole();
    await fetchUsers();
  });
</script>

<svelte:head>
  <title>Poin Users - Portal Guru</title>
</svelte:head>

<div class="animate-in fade-in duration-500 relative">
  <div class="mb-8 flex flex-col justify-start items-start gap-2">
    <h1
      class="text-2xl font-bold text-slate-900 sm:text-3xl tracking-tight drop-shadow-sm"
    >
      Poin Users
    </h1>
    <p
      class="text-slate-600 text-sm sm:text-base font-light tracking-wide"
    >
      Lihat peringkat poin dari seluruh pengguna (Leaderboard).
    </p>
  </div>

  {#if isTeacher}
    <div class="mb-10 flex justify-end">
      <button
        onclick={resetPoints}
        disabled={isResetting || isLoading}
        class="inline-flex items-center gap-2 text-red-600 hover:text-red-700 bg-red-50 hover:bg-red-100 border border-red-200 hover:border-red-300 font-semibold py-2 px-4 rounded-xl transition-all duration-200 shadow-sm disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-rotate-ccw"><path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"/><path d="M3 3v5h5"/></svg>
        {isResetting ? 'Mereset...' : 'Reset Semua Poin'}
      </button>
    </div>
  {/if}

  {#if isLoading}
    <div class="fixed inset-0 z-[100] flex flex-col items-center justify-center bg-slate-50/50 backdrop-blur-sm {showLoadingSpinner ? 'opacity-100' : 'opacity-0'} transition-opacity duration-300">
      <div class="w-12 h-12 border-4 border-slate-200 border-t-indigo-600 rounded-full animate-spin shadow-sm"></div>
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
    {#if users.length === 0}
      <div class="w-full bg-white/60 backdrop-blur-md rounded-3xl border border-slate-200 shadow-sm p-12 text-center text-slate-500 font-light">
        Tidak ada user ditemukan.
      </div>
    {:else}
      <div class="w-full grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
        {#each users as u}
          <div class="relative bg-white rounded-xl p-5 sm:p-6 shadow-sm border border-slate-200/60 overflow-hidden flex items-center justify-between">
            <!-- Decorative Background Element -->
            <div class="absolute -right-8 -top-8 w-32 h-32 bg-gradient-to-br from-indigo-100/40 to-blue-50/20 rounded-full blur-2xl pointer-events-none"></div>
            
            <div class="flex flex-col relative z-10 overflow-hidden pr-3">
              <span class="font-bold text-slate-800 text-lg sm:text-xl truncate">{u.username}</span>
              <div class="flex items-center gap-2 mt-1.5">
                <span class="w-1.5 h-1.5 rounded-full bg-indigo-500/80 shadow-[0_0_8px_rgba(99,102,241,0.5)]"></span>
                <span class="text-[13px] sm:text-sm font-semibold text-slate-500 capitalize tracking-wide">
                  {u.role}
                </span>
              </div>
            </div>
            
            <div class="relative z-10 flex flex-col items-center justify-center shrink-0 min-w-[3.5rem]">
              <span class="text-4xl sm:text-3xl font-black text-transparent bg-clip-text bg-gradient-to-r from-blue-600 to-indigo-600 leading-none text-center">
                {u.points || 0}
              </span>
              <span class="text-xs sm:text-[10px] font-bold text-slate-400 uppercase tracking-widest mt-1 text-center">
                Poin
              </span>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  {/if}
</div>

{#if showConfirmModal}
  <div class="fixed inset-0 z-150 flex items-center justify-center bg-slate-900/50 backdrop-blur-sm p-4">
    <div class="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden animate-in fade-in zoom-in-95 duration-200">
      <div class="p-6">
        <h3 class="text-xl font-bold text-slate-900 mb-2">Konfirmasi Reset Poin</h3>
        <p class="text-slate-600 mb-6">Apakah Anda yakin ingin mereset semua poin dan riwayat nilai kuis? Tindakan ini tidak dapat dibatalkan.</p>
        <div class="flex justify-end gap-3">
          <button onclick={() => showConfirmModal = false} class="px-4 py-2 font-medium text-slate-700 bg-slate-100 hover:bg-slate-200 rounded-xl transition-colors cursor-pointer">Batal</button>
          <button onclick={executeResetPoints} disabled={isResetting} class="px-4 py-2 font-medium text-white bg-red-600 hover:bg-red-700 rounded-xl transition-colors disabled:opacity-50 cursor-pointer">
            {isResetting ? 'Mereset...' : 'Ya, Reset Poin'}
          </button>
        </div>
      </div>
    </div>
  </div>
{/if}
