const API_BASE = '/api';

export interface Link {
	id: number;
	short_code: string;
	original_url: string;
	created_at: string;
}

export interface LinkListResponse {
	links: Link[];
	total: number;
	limit: number;
	offset: number;
}

export async function listLinks(limit = 20, offset = 0): Promise<LinkListResponse> {
	const res = await fetch(`${API_BASE}/links?limit=${limit}&offset=${offset}`);
	if (!res.ok) {
		throw new Error(`failed to list links: ${res.status}`);
	}
	return res.json();
}

export async function createLink(originalUrl: string, customCode: string): Promise<Link> {
	const res = await fetch(`${API_BASE}/links`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({
			original_url: originalUrl,
			custom_code: customCode,
		}),
	});
	if (!res.ok) {
		const text = await res.text();
		throw new Error(text || `failed to create link: ${res.status}`);
	}
	return res.json();
}

export async function getLink(id: number): Promise<Link> {
	const res = await fetch(`${API_BASE}/links/${id}`);
	if (!res.ok) {
		throw new Error(`failed to get link: ${res.status}`);
	}
	return res.json();
}

export async function updateOriginalUrl(id: number, originalUrl: string): Promise<Link> {
	const res = await fetch(`${API_BASE}/links/${id}/url`, {
		method: 'PATCH',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ original_url: originalUrl }),
	});
	if (!res.ok) {
		const text = await res.text();
		throw new Error(text || `failed to update URL: ${res.status}`);
	}
	return res.json();
}

export async function updateShortCode(id: number, customCode: string): Promise<Link> {
	const res = await fetch(`${API_BASE}/links/${id}/code`, {
		method: 'PATCH',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ custom_code: customCode }),
	});
	if (!res.ok) {
		const text = await res.text();
		throw new Error(text || `failed to update short code: ${res.status}`);
	}
	return res.json();
}

export async function deleteLink(id: number): Promise<void> {
	const res = await fetch(`${API_BASE}/links/${id}`, { method: 'DELETE' });
	if (!res.ok) {
		const text = await res.text();
		throw new Error(text || `failed to delete link: ${res.status}`);
	}
}

export interface Click {
	id: number;
	link_id: number;
	timestamp: string;
	referrer: string;
}

export interface ClickListResponse {
	clicks: Click[];
	total: number;
	limit: number;
	offset: number;
}

export async function listClicks(linkId: number, limit = 20, offset = 0): Promise<ClickListResponse> {
	const res = await fetch(`${API_BASE}/links/${linkId}/clicks?limit=${limit}&offset=${offset}`);
	if (!res.ok) {
		throw new Error(`failed to list clicks: ${res.status}`);
	}
	return res.json();
}
