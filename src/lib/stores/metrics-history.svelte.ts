import type { ParsedMetrics } from "$lib/types/api";

export interface MetricsSnapshot {
  timestamp: number;
  metrics: ParsedMetrics;
}

export interface ActivityPoint {
  timestamp: number;
  total: number;
  blocked: number;
}

const MAX_SNAPSHOTS = 30;

function createMetricsHistoryStore() {
  let snapshots = $state<MetricsSnapshot[]>([]);

  function latestDelta(
    extract: (m: ParsedMetrics) => number,
  ): number | undefined {
    if (snapshots.length < 2) return undefined;
    const prev = snapshots[snapshots.length - 2].metrics;
    const curr = snapshots[snapshots.length - 1].metrics;
    const delta = extract(curr) - extract(prev);
    return delta >= 0 ? delta : undefined;
  }

  function push(metrics: ParsedMetrics) {
    snapshots = [...snapshots, { timestamp: Date.now(), metrics }].slice(
      -MAX_SNAPSHOTS,
    );
  }

  function clear() {
    snapshots = [];
  }

  return {
    get snapshots() {
      return snapshots;
    },

    get activity(): ActivityPoint[] {
      if (snapshots.length < 2) return [];

      const points: ActivityPoint[] = [];
      for (let i = 1; i < snapshots.length; i++) {
        const prev = snapshots[i - 1].metrics;
        const curr = snapshots[i].metrics;

        const totalDelta = (curr.totalQueries ?? 0) - (prev.totalQueries ?? 0);
        const currBlocked = sumBlockedResponses(curr.responsesByReason);
        const prevBlocked = sumBlockedResponses(prev.responsesByReason);
        const blockedDelta = currBlocked - prevBlocked;

        // Only include if counters moved forward (skip resets)
        if (totalDelta >= 0 && blockedDelta >= 0) {
          points.push({
            timestamp: snapshots[i].timestamp,
            total: totalDelta,
            blocked: blockedDelta,
          });
        }
      }
      return points;
    },

    get queriesPerInterval(): number | undefined {
      return latestDelta((m) => m.totalQueries ?? 0);
    },

    get blockedPerInterval(): number | undefined {
      return latestDelta((m) => sumBlockedResponses(m.responsesByReason));
    },

    push,
    clear,
  };
}

function sumBlockedResponses(
  responsesByReason?: Record<string, number>,
): number {
  if (!responsesByReason) return 0;
  let total = 0;
  for (const [reason, count] of Object.entries(responsesByReason)) {
    if (reason.toUpperCase().includes("BLOCKED")) {
      total += count;
    }
  }
  return total;
}

export const metricsHistoryStore = createMetricsHistoryStore();
