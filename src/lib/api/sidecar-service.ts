import { sidecarRequest } from './sidecar';
import type { SidecarServiceStatus } from '$lib/types/api';

export async function fetchServiceStatus(): Promise<SidecarServiceStatus> {
	return sidecarRequest<SidecarServiceStatus>('/api/service/status');
}

export async function restartService(): Promise<{ status: string }> {
	return sidecarRequest<{ status: string }>('/api/service/restart', { method: 'POST' });
}
