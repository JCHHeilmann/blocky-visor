import type { ParsedMetrics } from "$lib/types/api";

interface MetricLine {
  name: string;
  labels: Record<string, string>;
  value: number;
}

function parseLine(line: string): MetricLine | null {
  if (line.startsWith("#") || line.trim() === "") return null;

  const match = line.match(/^(\w+)(?:\{(.+?)\})?\s+(.+)$/);
  if (!match) return null;

  const [, name, labelStr, valueStr] = match;
  const value = parseFloat(valueStr);
  if (isNaN(value)) return null;

  const labels: Record<string, string> = {};
  if (labelStr) {
    for (const pair of labelStr.split(",")) {
      const m = pair.match(/(\w+)="([^"]*)"/);
      if (m) labels[m[1]] = m[2];
    }
  }

  return { name, labels, value };
}

function addToRecord(
  record: Record<string, number>,
  key: string,
  value: number,
) {
  record[key] = (record[key] ?? 0) + value;
}

export function parsePrometheusMetrics(text: string): ParsedMetrics {
  const lines = text.split("\n");
  const metrics: ParsedMetrics = {};

  let cacheHits = 0;
  let cacheMisses = 0;
  let totalQueries = 0;
  let totalResponses = 0;
  let durationSum = 0;
  let durationCount = 0;
  const listEntries: Record<string, number> = {};
  const allowlistEntries: Record<string, number> = {};
  const queriesByClient: Record<string, number> = {};
  const queriesByType: Record<string, number> = {};
  const responsesByReason: Record<string, number> = {};
  const responsesByType: Record<string, number> = {};
  const responsesByCode: Record<string, number> = {};

  for (const line of lines) {
    const parsed = parseLine(line);
    if (!parsed) continue;

    switch (parsed.name) {
      case "blocky_cache_hits_total":
        cacheHits += parsed.value;
        break;
      case "blocky_cache_misses_total":
        cacheMisses += parsed.value;
        break;
      case "blocky_cache_entries":
        metrics.cacheEntryCount = parsed.value;
        break;
      case "blocky_denylist_cache_entries":
        if (parsed.labels.group) {
          listEntries[parsed.labels.group] = parsed.value;
        }
        break;
      case "blocky_allowlist_cache_entries":
        if (parsed.labels.group) {
          allowlistEntries[parsed.labels.group] = parsed.value;
        }
        break;
      case "blocky_prefetch_hits_total":
        metrics.prefetchHits = parsed.value;
        break;
      case "blocky_prefetches_total":
        metrics.prefetches = parsed.value;
        break;
      case "blocky_error_total":
        metrics.errors = parsed.value;
        break;
      case "blocky_blocking_enabled":
        metrics.blockingEnabled = parsed.value === 1;
        break;
      case "blocky_query_total":
        totalQueries += parsed.value;
        if (parsed.labels.client)
          addToRecord(queriesByClient, parsed.labels.client, parsed.value);
        if (parsed.labels.type)
          addToRecord(queriesByType, parsed.labels.type, parsed.value);
        break;
      case "blocky_response_total":
        totalResponses += parsed.value;
        if (parsed.labels.reason)
          addToRecord(responsesByReason, parsed.labels.reason, parsed.value);
        if (parsed.labels.response_type)
          addToRecord(
            responsesByType,
            parsed.labels.response_type,
            parsed.value,
          );
        if (parsed.labels.response_code)
          addToRecord(
            responsesByCode,
            parsed.labels.response_code,
            parsed.value,
          );
        break;
      case "blocky_request_duration_seconds_sum":
        durationSum += parsed.value;
        break;
      case "blocky_request_duration_seconds_count":
        durationCount += parsed.value;
        break;
      case "blocky_last_list_group_refresh_timestamp_seconds":
        metrics.lastListRefresh = parsed.value;
        break;
      case "blocky_failed_downloads_total":
        metrics.failedDownloads = parsed.value;
        break;
      case "blocky_build_info":
        if (parsed.labels.version) {
          metrics.buildInfo = {
            version: parsed.labels.version,
            buildTime: parsed.labels.build_time ?? "",
          };
        }
        break;
    }
  }

  if (cacheHits > 0) metrics.cacheHits = cacheHits;
  if (cacheMisses > 0) metrics.cacheMisses = cacheMisses;
  if (Object.keys(listEntries).length) metrics.listEntries = listEntries;
  if (Object.keys(allowlistEntries).length)
    metrics.allowlistEntries = allowlistEntries;
  if (totalQueries > 0) {
    metrics.totalQueries = totalQueries;
    metrics.queriesByClient = queriesByClient;
    metrics.queriesByType = queriesByType;
  }
  if (totalResponses > 0) {
    metrics.totalResponses = totalResponses;
    metrics.responsesByReason = responsesByReason;
    metrics.responsesByType = responsesByType;
    metrics.responsesByCode = responsesByCode;
  }
  if (durationCount > 0) {
    metrics.requestDurationSum = durationSum;
    metrics.requestDurationCount = durationCount;
  }

  return metrics;
}
