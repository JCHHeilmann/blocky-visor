import { settingsStore } from "$lib/stores/settings.svelte";
import { ConnectionError, ApiError } from "$lib/types/api";

const DEFAULT_TIMEOUT = 15000;

async function fetchWithTimeout(
  url: string,
  options: RequestInit & { timeout?: number } = {},
): Promise<Response> {
  const { timeout = DEFAULT_TIMEOUT, ...fetchOptions } = options;
  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), timeout);

  try {
    return await fetch(url, { ...fetchOptions, signal: controller.signal });
  } catch (err) {
    if (err instanceof DOMException && err.name === "AbortError") {
      throw new ConnectionError("Sidecar request timed out");
    }
    throw new ConnectionError("Failed to connect to sidecar", { cause: err });
  } finally {
    clearTimeout(timer);
  }
}

export async function sidecarRequest<T>(
  path: string,
  options: RequestInit & { timeout?: number } = {},
): Promise<T> {
  const { sidecarUrl, sidecarApiKey } = settingsStore;
  if (!sidecarUrl || !sidecarApiKey) {
    throw new ConnectionError("Sidecar not configured");
  }

  const { timeout, ...fetchOptions } = options;
  const response = await fetchWithTimeout(`${sidecarUrl}${path}`, {
    timeout,
    ...fetchOptions,
    headers: {
      "Content-Type": "application/json",
      "X-API-Key": sidecarApiKey,
      ...fetchOptions.headers,
    },
  });

  if (!response.ok) {
    const body = await response.text().catch(() => "");
    throw new ApiError(
      `Sidecar error: ${response.status}`,
      response.status,
      body,
    );
  }

  const contentType = response.headers.get("content-type");
  if (contentType?.includes("application/json")) {
    return (await response.json()) as T;
  }
  return (await response.text()) as unknown as T;
}

export async function testSidecarConnection(baseUrl: string): Promise<boolean> {
  try {
    const response = await fetchWithTimeout(`${baseUrl}/api/health`, {
      timeout: 5000,
    });
    return response.ok;
  } catch {
    return false;
  }
}
