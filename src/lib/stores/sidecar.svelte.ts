import { settingsStore } from "$lib/stores/settings.svelte";

function createSidecarStore() {
  let connected = $state(false);
  let loading = $state(false);
  let error = $state<string | null>(null);

  async function checkConnection() {
    if (!settingsStore.sidecarConfigured) {
      connected = false;
      return;
    }

    loading = true;
    error = null;
    try {
      const controller = new AbortController();
      const timer = setTimeout(() => controller.abort(), 5000);
      const response = await fetch(`${settingsStore.sidecarUrl}/api/health`, {
        signal: controller.signal,
      });
      clearTimeout(timer);
      connected = response.ok;
      if (!response.ok) {
        error = `Health check returned ${response.status}`;
      }
    } catch (err) {
      connected = false;
      error = err instanceof Error ? err.message : "Connection failed";
    } finally {
      loading = false;
    }
  }

  return {
    get connected() {
      return connected;
    },
    get loading() {
      return loading;
    },
    get error() {
      return error;
    },
    get configured() {
      return settingsStore.sidecarConfigured;
    },
    checkConnection,
    setDisconnected() {
      connected = false;
      error = "Connection lost";
    },
  };
}

export const sidecarStore = createSidecarStore();
