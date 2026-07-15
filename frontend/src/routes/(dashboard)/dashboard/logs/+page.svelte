<script lang="ts">
  import { onMount } from "svelte";

  type UserLog = {
    id: number;
    user_id: number;
    action: string;
    ip_address: string;
    user_agent: string;
    created_at: string;
    user: {
      username: string;
    };
  };

  let logs: UserLog[] = $state([]);
  let isLoading = $state(true);
  let errorMsg = $state("");
  let isTeacher = $state(false);

  // Pagination
  let currentPage = $state(1);
  let totalItems = $state(0);
  let limit = $state(50);
  let totalPages = $derived(Math.ceil(totalItems / limit) || 1);

  function handleLimitChange() {
    currentPage = 1;
    fetchLogs(1);
  }

  async function checkRole() {
    try {
      const res = await fetch(`/me`, { credentials: "include" });
      const data = await res.json();
      if (data.authenticated && data.user && data.user.role === "teacher") {
        isTeacher = true;
      } else {
        errorMsg = "Akses ditolak. Hanya guru yang dapat melihat halaman ini.";
      }
    } catch (e) {
      errorMsg = "Gagal memverifikasi akses.";
    }
  }

  async function fetchLogs(page = 1) {
    if (!isTeacher) return;
    isLoading = true;
    errorMsg = "";
    try {
      const res = await fetch(`/api/logs?page=${page}&limit=${limit}`, {
        credentials: "include",
      });
      if (!res.ok) throw new Error("Gagal mengambil data log");
      const data = await res.json();
      logs = data.data || [];
      totalItems = data.meta.total;
      currentPage = data.meta.page;
    } catch (e) {
      errorMsg = e instanceof Error ? e.message : String(e);
    } finally {
      isLoading = false;
    }
  }

  onMount(async () => {
    await checkRole();
    if (isTeacher) {
      await fetchLogs(1);
    } else {
      isLoading = false;
    }
  });

  function nextPage() {
    if (currentPage < totalPages) {
      fetchLogs(currentPage + 1);
    }
  }

  function prevPage() {
    if (currentPage > 1) {
      fetchLogs(currentPage - 1);
    }
  }

  function formatDate(dateStr: string) {
    const d = new Date(dateStr);
    return new Intl.DateTimeFormat("id-ID", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
    }).format(d);
  }
</script>

<svelte:head>
  <title>User Logs | Les Balong</title>
</svelte:head>

<div class="mb-6 flex flex-col items-start justify-between gap-4 md:flex-row md:items-center">
  <div>
    <h1 class="text-3xl font-extrabold text-gray-900 dark:text-white">
      User Logs
    </h1>
    <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
      Pantau riwayat login dan logout pengguna
    </p>
  </div>
</div>

{#if isLoading}
  <div class="flex items-center justify-center p-12">
    <div class="h-12 w-12 animate-spin rounded-full border-4 border-indigo-600 border-t-transparent"></div>
  </div>
{:else if errorMsg}
  <div class="rounded-xl border border-red-200 bg-red-50 p-6 text-center dark:border-red-900/50 dark:bg-red-900/20">
    <p class="font-medium text-red-600 dark:text-red-400">{errorMsg}</p>
  </div>
{:else}
  <div class="overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm dark:border-gray-800 dark:bg-gray-900">
    <div class="overflow-x-auto">
      <table class="w-full text-left text-sm text-gray-500 dark:text-gray-400">
        <thead class="bg-gray-50/50 text-xs uppercase text-gray-700 dark:bg-gray-800/50 dark:text-gray-300">
          <tr>
            <th class="px-6 py-4 font-semibold">User</th>
            <th class="px-6 py-4 font-semibold">Aksi</th>
            <th class="px-6 py-4 font-semibold">Waktu</th>
            <th class="px-6 py-4 font-semibold">IP Address</th>
            <th class="px-6 py-4 font-semibold">User Agent</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200 dark:divide-gray-800">
          {#each logs as log}
            <tr class="transition-colors hover:bg-gray-50 dark:hover:bg-gray-800/50">
              <td class="whitespace-nowrap px-6 py-4 font-medium text-gray-900 dark:text-white">
                {log.user?.username || `User ID: ${log.user_id}`}
              </td>
              <td class="whitespace-nowrap px-6 py-4">
                <span class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium {log.action === 'login' ? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400' : 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400'}">
                  {log.action.toUpperCase()}
                </span>
              </td>
              <td class="whitespace-nowrap px-6 py-4">
                {formatDate(log.created_at)}
              </td>
              <td class="whitespace-nowrap px-6 py-4 text-xs font-mono">
                {log.ip_address}
              </td>
              <td class="px-6 py-4 text-xs max-w-xs truncate" title={log.user_agent}>
                {log.user_agent}
              </td>
            </tr>
          {:else}
            <tr>
              <td colspan="5" class="px-6 py-12 text-center text-gray-500 dark:text-gray-400">
                Belum ada data log aktivitas.
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>

    {#if totalPages > 1 || totalItems > 0}
      <div class="flex flex-col gap-4 border-t border-gray-200 px-4 py-4 sm:flex-row sm:items-center sm:justify-between sm:px-6 dark:border-gray-800">
        
        <div class="flex items-center justify-between sm:justify-start sm:gap-4">
          <p class="text-sm text-gray-500 dark:text-gray-400">
            Hal <span class="font-medium text-gray-900 dark:text-white">{currentPage}</span>/<span class="font-medium text-gray-900 dark:text-white">{totalPages}</span>
          </p>
          
          <div class="flex items-center gap-2 text-sm text-gray-500 sm:border-l sm:border-gray-300 sm:pl-4 dark:text-gray-400 dark:sm:border-gray-700">
            <label for="per-page" class="hidden sm:inline">Per halaman:</label>
            <label for="per-page" class="sm:hidden">Tampil:</label>
            <select
              id="per-page"
              bind:value={limit}
              onchange={handleLimitChange}
              class="rounded-md border border-gray-300 bg-white py-1 pl-2 pr-6 text-sm shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500 dark:border-gray-600 dark:bg-gray-800 dark:text-white dark:focus:border-indigo-500"
            >
              <option value={10}>10</option>
              <option value={20}>20</option>
              <option value={50}>50</option>
            </select>
          </div>
        </div>

        <div class="flex w-full justify-between gap-2 sm:w-auto sm:justify-end">
          <button
            onclick={prevPage}
            disabled={currentPage === 1}
            class="flex-1 rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-700 shadow-sm transition-colors hover:bg-gray-50 disabled:opacity-50 disabled:hover:bg-white sm:flex-none dark:border-gray-700 dark:bg-gray-800 dark:text-gray-200 dark:hover:bg-gray-700 dark:disabled:hover:bg-gray-800"
          >
            Sebelumnya
          </button>
          <button
            onclick={nextPage}
            disabled={currentPage === totalPages}
            class="flex-1 rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-700 shadow-sm transition-colors hover:bg-gray-50 disabled:opacity-50 disabled:hover:bg-white sm:flex-none dark:border-gray-700 dark:bg-gray-800 dark:text-gray-200 dark:hover:bg-gray-700 dark:disabled:hover:bg-gray-800"
          >
            Selanjutnya
          </button>
        </div>

      </div>
    {/if}
  </div>
{/if}
