<script lang="ts">
    import Card from "$lib/components/ui/Card.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import Input from "$lib/components/ui/Input.svelte";
    import ApiUrlForm from "$lib/components/settings/ApiUrlForm.svelte";
    import SidecarForm from "$lib/components/settings/SidecarForm.svelte";
    import { settingsStore } from "$lib/stores/settings.svelte";
    import { toastStore } from "$lib/stores/toasts.svelte";
    import { themeStore } from "$lib/stores/theme.svelte";
    import { fetchMetrics } from "$lib/api/metrics";
    import { formatDate } from "$lib/utils/formatters";
    import type { ParsedMetrics } from "$lib/types/api";

    let blockyInfo = $state<ParsedMetrics | null>(null);

    $effect(() => {
        fetchMetrics().then((m) => {
            if (m) blockyInfo = m;
        });
    });

    let refreshInterval = $state(String(settingsStore.refreshInterval));
    let metricsInterval = $state(String(settingsStore.metricsInterval));

    const themeOptions = [
        {
            value: "system",
            label: "System",
            icon: "M9 17.25v1.007a3 3 0 0 1-.879 2.122L7.5 21h9l-.621-.621A3 3 0 0 1 15 18.257V17.25m6-12V15a2.25 2.25 0 0 1-2.25 2.25h-13.5A2.25 2.25 0 0 1 3 15V5.25m18 0A2.25 2.25 0 0 0 18.75 3H5.25A2.25 2.25 0 0 0 3 5.25m18 0V12a2.25 2.25 0 0 1-2.25 2.25h-13.5A2.25 2.25 0 0 1 3 12V5.25",
        },
        {
            value: "light",
            label: "Light",
            icon: "M12 3v2.25m6.364.386-1.591 1.591M21 12h-2.25m-.386 6.364-1.591-1.591M12 18.75V21m-4.773-4.227-1.591 1.591M5.25 12H3m4.227-4.773L5.636 5.636M15.75 12a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0Z",
        },
        {
            value: "dark",
            label: "Dark",
            icon: "M21.752 15.002A9.72 9.72 0 0 1 18 15.75c-5.385 0-9.75-4.365-9.75-9.75 0-1.33.266-2.597.748-3.752A9.753 9.753 0 0 0 3 11.25C3 16.635 7.365 21 12.75 21a9.753 9.753 0 0 0 9.002-5.998Z",
        },
    ] as const;

    function saveIntervals() {
        const refresh = parseInt(refreshInterval);
        const metrics = parseInt(metricsInterval);
        if (refresh >= 1 && refresh <= 300) {
            settingsStore.refreshInterval = refresh;
        }
        if (metrics >= 5 && metrics <= 600) {
            settingsStore.metricsInterval = metrics;
        }
        toastStore.success("Settings saved");
    }

    function resetDefaults() {
        settingsStore.resetDefaults();
        refreshInterval = String(settingsStore.refreshInterval);
        metricsInterval = String(settingsStore.metricsInterval);
        themeStore.setPreference("system");
        toastStore.info("Settings reset to defaults");
    }
</script>

<div class="max-w-2xl space-y-6">
    <Card>
        <h2 class="mb-4 text-lg font-semibold text-text-primary">
            API Connection
        </h2>
        <ApiUrlForm />
    </Card>

    <Card>
        <h2 class="mb-4 text-lg font-semibold text-text-primary">
            Sidecar Connection
        </h2>
        <SidecarForm />
    </Card>

    {#if blockyInfo?.buildInfo || blockyInfo?.lastListRefresh !== undefined}
        <Card>
            <h2 class="mb-4 text-lg font-semibold text-text-primary">
                Blocky Instance
            </h2>
            <div class="space-y-3 text-sm">
                {#if blockyInfo.buildInfo}
                    <div class="flex justify-between">
                        <span class="text-text-muted">Version</span>
                        <span class="font-mono text-text-primary"
                            >{blockyInfo.buildInfo.version}</span
                        >
                    </div>
                    {#if blockyInfo.buildInfo.buildTime}
                        <div class="flex justify-between">
                            <span class="text-text-muted">Build Time</span>
                            <span class="font-mono text-text-primary"
                                >{blockyInfo.buildInfo.buildTime}</span
                            >
                        </div>
                    {/if}
                {/if}
                {#if blockyInfo.lastListRefresh}
                    <div class="flex justify-between">
                        <span class="text-text-muted">Last List Refresh</span>
                        <span class="font-mono text-text-primary"
                            >{formatDate(
                                blockyInfo.lastListRefresh * 1000,
                            )}</span
                        >
                    </div>
                {/if}
                {#if blockyInfo.failedDownloads !== undefined && blockyInfo.failedDownloads > 0}
                    <div class="flex justify-between">
                        <span class="text-text-muted">Failed Downloads</span>
                        <span class="font-mono text-error"
                            >{blockyInfo.failedDownloads}</span
                        >
                    </div>
                {/if}
            </div>
            <p class="mt-4 text-xs text-text-faint">
                Blocky configuration is managed via YAML config files. Changes
                require editing the config file and restarting the Blocky
                service.
            </p>
        </Card>
    {/if}

    <Card>
        <h2 class="mb-4 text-lg font-semibold text-text-primary">Appearance</h2>
        <p class="mb-3 text-sm text-text-muted">
            Choose your preferred color theme
        </p>
        <div class="flex gap-2">
            {#each themeOptions as option}
                <button
                    onclick={() => themeStore.setPreference(option.value)}
                    class="flex items-center gap-2 rounded-lg border px-4 py-2.5 text-sm font-medium transition-colors cursor-pointer
						{themeStore.preference === option.value
                        ? 'border-accent-500 bg-accent-600/15 text-accent-600 dark:text-accent-400'
                        : 'border-surface-border bg-surface-secondary text-text-secondary hover:border-text-muted'}"
                >
                    <svg
                        class="h-4 w-4"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke-width="1.5"
                        stroke="currentColor"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            d={option.icon}
                        />
                    </svg>
                    {option.label}
                </button>
            {/each}
        </div>
    </Card>

    <Card>
        <h2 class="mb-4 text-lg font-semibold text-text-primary">
            Refresh Intervals
        </h2>
        <div class="space-y-4">
            <div>
                <Input
                    bind:value={refreshInterval}
                    label="Blocking status poll interval (seconds)"
                    type="number"
                />
                <p class="mt-1 text-xs text-text-faint">Min: 1, Max: 300</p>
            </div>
            <div>
                <Input
                    bind:value={metricsInterval}
                    label="Metrics refresh interval (seconds)"
                    type="number"
                />
                <p class="mt-1 text-xs text-text-faint">Min: 5, Max: 600</p>
            </div>
            <Button onclick={saveIntervals}>Save Intervals</Button>
        </div>
    </Card>

    <Card>
        <h2 class="mb-4 text-lg font-semibold text-text-primary">Reset</h2>
        <p class="mb-4 text-sm text-text-secondary">
            Reset all settings to their default values.
        </p>
        <Button variant="danger" onclick={resetDefaults}>Reset Defaults</Button>
    </Card>
</div>
