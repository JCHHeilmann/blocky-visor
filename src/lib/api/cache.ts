import { apiRequest } from './client';

export async function flushCache(): Promise<void> {
	await apiRequest('/api/cache/flush', { method: 'PUT' });
}
