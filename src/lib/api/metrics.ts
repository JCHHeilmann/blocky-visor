import { apiRequestRaw } from './client';
import { parsePrometheusMetrics } from '$lib/utils/metrics-parser';
import type { ParsedMetrics } from '$lib/types/api';

export async function fetchMetrics(): Promise<ParsedMetrics | null> {
	try {
		const response = await apiRequestRaw('/metrics', { timeout: 5000 });
		if (!response.ok) return null;
		const text = await response.text();
		return parsePrometheusMetrics(text);
	} catch {
		return null;
	}
}
