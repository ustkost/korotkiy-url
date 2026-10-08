<script lang="ts">
	import { onMount } from 'svelte';
	import { push } from 'svelte-spa-router';
	import { getLink, updateOriginalUrl, updateShortCode } from '../api';

	export let params: { id: string };

	let id: number;
	let originalUrl = '';
	let shortCode = '';
	let originalUrlInitial = '';
	let shortCodeInitial = '';
	let loading = true;
	let submitting = false;
	let error: string | null = null;

	onMount(async () => {
		id = Number(params.id);
		try {
			const link = await getLink(id);
			originalUrl = link.original_url;
			shortCode = link.short_code;
			originalUrlInitial = link.original_url;
			shortCodeInitial = link.short_code;
		} catch (e) {
			error = e instanceof Error ? e.message : 'unknown error';
		} finally {
			loading = false;
		}
	});

	async function handleSubmit() {
		submitting = true;
		error = null;
		try {
			if (originalUrl !== originalUrlInitial) {
				await updateOriginalUrl(id, originalUrl);
			}
			if (shortCode !== shortCodeInitial) {
				await updateShortCode(id, shortCode);
			}
			push('/');
		} catch (e) {
			error = e instanceof Error ? e.message : 'unknown error';
		} finally {
			submitting = false;
		}
	}
</script>

<h1>Edit Link</h1>

{#if loading}
	<p>Loading...</p>
{:else}
	<form on:submit|preventDefault={handleSubmit}>
		<label>
			Original URL
			<input type="url" bind:value={originalUrl} required />
		</label>

		<label>
			Short code (leave empty for regenerated)
			<input type="text" bind:value={shortCode} />
		</label>

		{#if error}
			<p class="error">{error}</p>
		{/if}

		<button type="submit" disabled={submitting} class="primary">
			{submitting ? 'Saving...' : 'Save'}
		</button>
	</form>
{/if}
