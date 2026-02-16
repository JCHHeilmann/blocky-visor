import { getBlockingStatus } from "$lib/api/blocking";
import type { BlockingStatus } from "$lib/types/api";

function createBlockingStore() {
  let status = $state<BlockingStatus | null>(null);
  let connected = $state(false);
  let loading = $state(false);
  let initialLoad = $state(true);
  let error = $state<string | null>(null);
  let lastUpdated = $state<number | null>(null);

  async function refresh() {
    loading = true;
    error = null;
    try {
      status = await getBlockingStatus();
      connected = true;
      lastUpdated = Date.now();
    } catch (err) {
      connected = false;
      error = err instanceof Error ? err.message : "Unknown error";
    } finally {
      loading = false;
      initialLoad = false;
    }
  }

  return {
    get status() {
      return status;
    },
    get connected() {
      return connected;
    },
    get loading() {
      return loading;
    },
    get initialLoad() {
      return initialLoad;
    },
    get error() {
      return error;
    },
    get lastUpdated() {
      return lastUpdated;
    },
    get enabled() {
      return status?.enabled ?? false;
    },
    get autoEnableInSec() {
      return status?.autoEnableInSec;
    },
    get disabledGroups() {
      return status?.disabledGroups;
    },
    refresh,
    setDisconnected() {
      connected = false;
      error = "Connection lost";
    },
  };
}

export const blockingStore = createBlockingStore();
