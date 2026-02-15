import { apiRequest } from './client';

export async function refreshLists(): Promise<void> {
	await apiRequest('/api/lists/refresh', { method: 'POST' });
}
