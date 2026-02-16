import { settingsStore } from "$lib/stores/settings.svelte";
import { ConnectionError, ApiError } from "$lib/types/api";

const DEFAULT_TIMEOUT = 15000;

export async function sidecarRequest<T>(
  path: string,
  options: RequestInit & { timeout?: number } = {},
): Promise<T> {
  const { sidecarUrl, sidecarApiKey } = settingsStore;
  if (!sidecarUrl || !sidecarApiKey) {
    throw new ConnectionError("Sidecar not configured");
  }

  const { timeout = DEFAULT_TIMEOUT, ...fetchOptions } = options;
  const url = `${sidecarUrl}${path}`;

  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), timeout);

  try {
    const response = await fetch(url, {
      ...fetchOptions,
      signal: controller.signal,
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
  } catch (err) {
    if (err instanceof ApiError) throw err;
    if (err instanceof DOMException && err.name === "AbortError") {
      throw new ConnectionError("Sidecar request timed out");
    }
    throw new ConnectionError("Failed to connect to sidecar", err);
  } finally {
    clearTimeout(timer);
  }
}

export async function testSidecarConnection(
  baseUrl: string,
  apiKey: string,
): Promise<boolean> {
  try {
    const controller = new AbortController();
    const timer = setTimeout(() => controller.abort(), 5000);
    const response = await fetch(`${baseUrl}/api/health`, {
      signal: controller.signal,
    });
    clearTimeout(timer);
    return response.ok;
  } catch {
    return false;
  }
}
