<script>
	import Router, { link } from 'svelte-spa-router';
	import routes from './routes.js';
</script>

<header>
	<a href="/" use:link class="brand">korotkiy-url</a>
</header>

<main>
	<Router {routes} />
</main>
