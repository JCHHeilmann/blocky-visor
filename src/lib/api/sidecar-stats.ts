import { sidecarRequest } from "./sidecar";
import { settingsStore } from "$lib/stores/settings.svelte";
import type {
  SidecarStatsResponse,
  SidecarTimelineBucket,
  SidecarLogsResponse,
} from "$lib/types/api";

export type StatsRange = "today" | "yesterday" | "7d" | "30d";

export async function fetchStats(
  range: StatsRange = "today",
): Promise<SidecarStatsResponse> {
  return sidecarRequest<SidecarStatsResponse>(`/api/stats?range=${range}`);
}

export async function fetchTimeline(
  range: StatsRange = "today",
  interval: "5m" | "15m" | "1h" | "1d" = "15m",
): Promise<SidecarTimelineBucket[]> {
  return sidecarRequest<SidecarTimelineBucket[]>(
    `/api/stats/timeline?range=${range}&interval=${interval}`,
  );
}

export async function fetchLogs(
  params: {
    range?: StatsRange;
    limit?: number;
    offset?: number;
    client?: string;
    domain?: string;
    type?: string;
  } = {},
): Promise<SidecarLogsResponse> {
  const searchParams = new URLSearchParams();
  if (params.range) searchParams.set("range", params.range);
  if (params.limit) searchParams.set("limit", String(params.limit));
  if (params.offset) searchParams.set("offset", String(params.offset));
  if (params.client) searchParams.set("client", params.client);
  if (params.domain) searchParams.set("domain", params.domain);
  if (params.type) searchParams.set("type", params.type);

  const qs = searchParams.toString();
  return sidecarRequest<SidecarLogsResponse>(`/api/logs${qs ? `?${qs}` : ""}`);
}

export function buildLogStreamUrl(filters?: {
  client?: string;
  domain?: string;
  type?: string;
}): string {
  const { sidecarUrl, sidecarApiKey } = settingsStore;
  const params = new URLSearchParams();
  params.set("key", sidecarApiKey);
  if (filters?.client) params.set("client", filters.client);
  if (filters?.domain) params.set("domain", filters.domain);
  if (filters?.type) params.set("type", filters.type);
  return `${sidecarUrl}/api/logs/stream?${params.toString()}`;
}
