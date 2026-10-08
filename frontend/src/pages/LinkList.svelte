<script lang="ts">
	import { onMount } from 'svelte';
	import { link } from 'svelte-spa-router';
	import { listLinks, deleteLink, type Link as LinkType } from '../api';

	const PAGE_SIZE = 20;

	let links: LinkType[] = [];
	let total = 0;
	let loading = true;
	let loadingMore = false;
	let error: string | null = null;

	async function load() {
		loading = true;
		error = null;
		try {
			const result = await listLinks(PAGE_SIZE, 0);
			links = result.links ?? [];
			total = result.total;
		} catch (e) {
			error = e instanceof Error ? e.message : 'unknown error';
		} finally {
			loading = false;
		}
	}

	async function loadMore() {
		loadingMore = true;
		error = null;
		try {
			const result = await listLinks(PAGE_SIZE, links.length);
			links = [...links, ...(result.links ?? [])];
			total = result.total;
		} catch (e) {
			error = e instanceof Error ? e.message : 'unknown error';
		} finally {
			loadingMore = false;
		}
	}

	async function handleDelete(id: number) {
		if (!confirm('Delete this link?')) return;
		try {
			await deleteLink(id);
			links = links.filter((l) => l.id !== id);
			total -= 1;
		} catch (e) {
			error = e instanceof Error ? e.message : 'unknown error';
		}
	}

  let copiedId: number | null = null;

  async function copyLink(l: LinkType) {
    const url = `${window.location.origin}/s/${l.short_code}`;
    try {
      await navigator.clipboard.writeText(url);
      copiedId = l.id;
      setTimeout(() => {
        if (copiedId === l.id) copiedId = null;
      }, 1500);
    } catch {
      error = 'could not copy to clipboard';
    }
  }

	onMount(load);
</script>

<h1>Links</h1>

<a href="/links/new" use:link>+ Create Link</a>

{#if loading}
	<p>Loading...</p>
{:else}
	{#if error}
		<p class="error">Error: {error}</p>
	{/if}

	<p>Total links: {total}</p>
	<ul>
		{#each links as l (l.id)}
      <li>
        <strong>{l.short_code}</strong> → {l.original_url}
        <button on:click={() => copyLink(l)}>
          {copiedId === l.id ? 'Copied!' : 'Copy'}
        </button>
        <a href="/links/{l.id}/edit" use:link>Edit</a>
        <a href="/links/{l.id}/stats" use:link>Stats</a>
        <button on:click={() => handleDelete(l.id)}>Delete</button>
      </li>
		{/each}
	</ul>

	{#if links.length < total}
		<button on:click={loadMore} disabled={loadingMore}>
			{loadingMore ? 'Loading...' : 'Load more'}
		</button>
	{/if}
{/if}
