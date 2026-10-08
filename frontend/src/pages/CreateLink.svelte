<script lang="ts">
	import { push } from 'svelte-spa-router';
	import { createLink } from '../api';

	let originalUrl = '';
	let customCode = '';
	let submitting = false;
	let error: string | null = null;

	async function handleSubmit() {
		submitting = true;
		error = null;
		try {
			await createLink(originalUrl, customCode);
			push('/');
		} catch (e) {
			error = e instanceof Error ? e.message : 'unknown error';
		} finally {
			submitting = false;
		}
	}
</script>

<h1>Create Link</h1>

<form on:submit|preventDefault={handleSubmit}>
	<label>
		Original URL
		<input type="url" bind:value={originalUrl} required />
	</label>

	<label>
		Custom code (optional)
		<input type="text" bind:value={customCode} />
	</label>

	{#if error}
		<p class="error">{error}</p>
	{/if}

	<button type="submit" disabled={submitting}>
		{submitting ? 'Creating...' : 'Create'}
	</button>
</form>
