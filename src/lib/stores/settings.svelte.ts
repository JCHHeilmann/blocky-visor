const STORAGE_KEY = 'blocky-visor-settings';

const DEFAULTS = {
	apiUrl: 'http://localhost:4000',
	refreshInterval: 5,
	metricsInterval: 30,
	sidecarUrl: '',
	sidecarApiKey: ''
};

interface Settings {
	apiUrl: string;
	refreshInterval: number;
	metricsInterval: number;
	sidecarUrl: string;
	sidecarApiKey: string;
}

function loadSettings(): Settings {
	if (typeof window === 'undefined') return { ...DEFAULTS };
	try {
		const stored = localStorage.getItem(STORAGE_KEY);
		if (stored) return { ...DEFAULTS, ...JSON.parse(stored) };
	} catch {}
	return { ...DEFAULTS };
}

function saveSettings(settings: Settings) {
	if (typeof window === 'undefined') return;
	localStorage.setItem(STORAGE_KEY, JSON.stringify(settings));
	// Also keep the api URL in the key the client reads
	localStorage.setItem('blocky-api-url', settings.apiUrl);
}

function createSettingsStore() {
	let settings = $state<Settings>(loadSettings());

	// Sync api URL on init
	if (typeof window !== 'undefined') {
		const initial = loadSettings();
		localStorage.setItem('blocky-api-url', initial.apiUrl);
	}

	return {
		get apiUrl() { return settings.apiUrl; },
		set apiUrl(value: string) {
			settings.apiUrl = value;
			saveSettings(settings);
		},
		get refreshInterval() { return settings.refreshInterval; },
		set refreshInterval(value: number) {
			settings.refreshInterval = value;
			saveSettings(settings);
		},
		get metricsInterval() { return settings.metricsInterval; },
		set metricsInterval(value: number) {
			settings.metricsInterval = value;
			saveSettings(settings);
		},
		get sidecarUrl() { return settings.sidecarUrl; },
		set sidecarUrl(value: string) {
			settings.sidecarUrl = value;
			saveSettings(settings);
		},
		get sidecarApiKey() { return settings.sidecarApiKey; },
		set sidecarApiKey(value: string) {
			settings.sidecarApiKey = value;
			saveSettings(settings);
		},
		get sidecarConfigured() {
			return settings.sidecarUrl !== '' && settings.sidecarApiKey !== '';
		},
		resetDefaults() {
			settings = { ...DEFAULTS };
			saveSettings(settings);
		}
	};
}

export const settingsStore = createSettingsStore();
