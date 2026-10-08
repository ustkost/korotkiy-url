import LinkList from './pages/LinkList.svelte';
import CreateLink from './pages/CreateLink.svelte';
import EditLink from './pages/EditLink.svelte';
import LinkStats from './pages/LinkStats.svelte';

export default {
	'/': LinkList,
	'/links/new': CreateLink,
	'/links/:id/edit': EditLink,
	'/links/:id/stats': LinkStats,
};
