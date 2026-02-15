export interface BlockingStatus {
  enabled: boolean;
  disabledGroups?: string[];
  autoEnableInSec?: number;
}

export interface DnsQueryRequest {
  domain: string;
  type: string;
}

export interface DnsQueryResponse {
  reason: string;
  response: string;
  responseType: string;
  returnCode: string;
}

export interface QueryHistoryEntry {
  domain: string;
  type: string;
  response: DnsQueryResponse;
  timestamp: number;
}

export interface ParsedMetrics {
  cacheHits?: number;
  cacheMisses?: number;
  cacheEntryCount?: number;
  listEntries?: Record<string, number>;
  allowlistEntries?: Record<string, number>;
  prefetchHits?: number;
  prefetches?: number;
  errors?: number;
  blockingEnabled?: boolean;
  totalQueries?: number;
  queriesByClient?: Record<string, number>;
  queriesByType?: Record<string, number>;
  totalResponses?: number;
  responsesByReason?: Record<string, number>;
  responsesByType?: Record<string, number>;
  responsesByCode?: Record<string, number>;
  requestDurationSum?: number;
  requestDurationCount?: number;
  lastListRefresh?: number;
  failedDownloads?: number;
  buildInfo?: { version: string; buildTime: string };
}

// Sidecar types

export interface SidecarStatsResponse {
  period: { start: string; end: string; files_parsed: number };
  summary: {
    total_queries: number;
    blocked_queries: number;
    cached_queries: number;
    unique_domains: number;
    unique_clients: number;
    avg_duration_ms: number;
    p95_duration_ms: number;
  };
  hourly: { hour: number; total: number; blocked: number; cached: number }[];
  top_domains: { domain: string; count: number }[];
  top_blocked: { domain: string; count: number; reason: string }[];
  clients: { ip: string; name: string; total: number; blocked: number }[];
  query_types: Record<string, number>;
  response_categories: Record<string, number>;
  return_codes: Record<string, number>;
}

export interface SidecarTimelineBucket {
  timestamp: string;
  total: number;
  blocked: number;
  cached: number;
}

export interface SidecarLogEntry {
  timestamp: string;
  client_ip: string;
  client_name: string;
  domain: string;
  query_type: string;
  response_reason: string;
  return_code: string;
  duration_ms: number;
  response_answer: string;
}

export interface SidecarLogsResponse {
  total: number;
  offset: number;
  limit: number;
  entries: SidecarLogEntry[];
}

export interface SidecarServiceStatus {
  active: string;
  sub_state: string;
  pid?: string;
  memory?: string;
  uptime?: string;
  full_status: string;
}

export class ConnectionError extends Error {
  constructor(
    message: string,
    public cause?: unknown,
  ) {
    super(message);
    this.name = "ConnectionError";
  }
}

export class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
    public body?: string,
  ) {
    super(message);
    this.name = "ApiError";
  }
}
