import { apiRequest } from './client';
import type { DnsQueryResponse } from '$lib/types/api';

export async function dnsQuery(domain: string, type: string): Promise<DnsQueryResponse> {
	return apiRequest<DnsQueryResponse>('/api/query', {
		method: 'POST',
		body: JSON.stringify({ query: domain, type })
	});
}
