<script lang="ts">
	import { onMount } from 'svelte';
	import { getLink, listClicks, type Link, type Click } from '../api';

	export let params: { id: string };

	const PAGE_SIZE = 20;

	let id: number;
	let linkData: Link | null = null;
	let clicks: Click[] = [];
	let total = 0;
	let loading = true;
	let loadingMore = false;
	let error: string | null = null;

	onMount(async () => {
		try {
			id = Number(params.id);
			const [link, clickResult] = await Promise.all([
				getLink(id),
				listClicks(id, PAGE_SIZE, 0),
			]);
			linkData = link;
			clicks = clickResult.clicks ?? [];
			total = clickResult.total;
		} catch (e) {
			error = e instanceof Error ? e.message : 'unknown error';
		} finally {
			loading = false;
		}
	});

	async function loadMore() {
		loadingMore = true;
		error = null;
		try {
			const result = await listClicks(id, PAGE_SIZE, clicks.length);
			clicks = [...clicks, ...(result.clicks ?? [])];
			total = result.total;
		} catch (e) {
			error = e instanceof Error ? e.message : 'unknown error';
		} finally {
			loadingMore = false;
		}
	}
</script>

<h1>Link Stats</h1>

{#if loading}
	<p>Loading...</p>
{:else if !linkData}
	<p class="error">Error: {error ?? 'link not found'}</p>
{:else}
	<p><strong>{linkData.short_code}</strong> → {linkData.original_url}</p>
	<p>Total clicks: {total}</p>

	{#if error}
		<p class="error">Error: {error}</p>
	{/if}

	{#if clicks.length === 0}
		<p>No clicks yet.</p>
	{:else}
		<table>
			<thead>
				<tr>
					<th>Time</th>
					<th>Referrer</th>
				</tr>
			</thead>
			<tbody>
				{#each clicks as c (c.id)}
					<tr>
						<td>{new Date(c.timestamp).toLocaleString()}</td>
						<td>{c.referrer || 'direct'}</td>
					</tr>
				{/each}
			</tbody>
		</table>

		{#if clicks.length < total}
			<button on:click={loadMore} disabled={loadingMore}>
				{loadingMore ? 'Loading...' : 'Load more'}
			</button>
		{/if}
	{/if}
{/if}
