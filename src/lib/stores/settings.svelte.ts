const STORAGE_KEY = "blocky-visor-settings";
const API_URL_KEY = "blocky-api-url";

const DEFAULTS: Settings = {
  apiUrl: "http://localhost:4000",
  refreshInterval: 5,
  metricsInterval: 30,
  sidecarUrl: "",
  sidecarApiKey: "",
};

interface Settings {
  apiUrl: string;
  refreshInterval: number;
  metricsInterval: number;
  sidecarUrl: string;
  sidecarApiKey: string;
}

function loadSettings(): Settings {
  if (typeof window === "undefined") return { ...DEFAULTS };
  try {
    const stored = localStorage.getItem(STORAGE_KEY);
    if (stored) return { ...DEFAULTS, ...JSON.parse(stored) };
  } catch {}
  return { ...DEFAULTS };
}

function persist(settings: Settings): void {
  if (typeof window === "undefined") return;
  localStorage.setItem(STORAGE_KEY, JSON.stringify(settings));
  localStorage.setItem(API_URL_KEY, settings.apiUrl);
}

function createSettingsStore() {
  let settings = $state<Settings>(loadSettings());

  if (typeof window !== "undefined") {
    localStorage.setItem(API_URL_KEY, settings.apiUrl);
  }

  function update<K extends keyof Settings>(key: K, value: Settings[K]): void {
    settings[key] = value;
    persist(settings);
  }

  return {
    get apiUrl() {
      return settings.apiUrl;
    },
    set apiUrl(value: string) {
      update("apiUrl", value);
    },
    get refreshInterval() {
      return settings.refreshInterval;
    },
    set refreshInterval(value: number) {
      update("refreshInterval", value);
    },
    get metricsInterval() {
      return settings.metricsInterval;
    },
    set metricsInterval(value: number) {
      update("metricsInterval", value);
    },
    get sidecarUrl() {
      return settings.sidecarUrl;
    },
    set sidecarUrl(value: string) {
      update("sidecarUrl", value);
    },
    get sidecarApiKey() {
      return settings.sidecarApiKey;
    },
    set sidecarApiKey(value: string) {
      update("sidecarApiKey", value);
    },
    get sidecarConfigured() {
      return settings.sidecarUrl !== "" && settings.sidecarApiKey !== "";
    },
    resetDefaults() {
      settings = { ...DEFAULTS };
      persist(settings);
    },
  };
}

export const settingsStore = createSettingsStore();
