import { ConnectionError, ApiError } from "$lib/types/api";

const DEFAULT_TIMEOUT = 10000;

function getBaseUrl(): string {
  if (typeof window === "undefined") return "";
  return localStorage.getItem("blocky-api-url") || "http://localhost:4000";
}

export async function apiRequest<T>(
  path: string,
  options: RequestInit & { timeout?: number } = {},
): Promise<T> {
  const { timeout = DEFAULT_TIMEOUT, ...fetchOptions } = options;
  const url = `${getBaseUrl()}${path}`;

  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), timeout);

  try {
    const response = await fetch(url, {
      ...fetchOptions,
      signal: controller.signal,
      headers: {
        "Content-Type": "application/json",
        ...fetchOptions.headers,
      },
    });

    if (!response.ok) {
      const body = await response.text().catch(() => "");
      throw new ApiError(
        `API error: ${response.status}`,
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
      throw new ConnectionError("Request timed out");
    }
    throw new ConnectionError("Failed to connect to Blocky", err);
  } finally {
    clearTimeout(timer);
  }
}

export async function apiRequestRaw(
  path: string,
  options: RequestInit & { timeout?: number } = {},
): Promise<Response> {
  const { timeout = DEFAULT_TIMEOUT, ...fetchOptions } = options;
  const url = `${getBaseUrl()}${path}`;

  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), timeout);

  try {
    const response = await fetch(url, {
      ...fetchOptions,
      signal: controller.signal,
    });
    return response;
  } catch (err) {
    if (err instanceof DOMException && err.name === "AbortError") {
      throw new ConnectionError("Request timed out");
    }
    throw new ConnectionError("Failed to connect to Blocky", err);
  } finally {
    clearTimeout(timer);
  }
}

export async function testConnection(baseUrl: string): Promise<boolean> {
  try {
    const controller = new AbortController();
    const timer = setTimeout(() => controller.abort(), 5000);
    const response = await fetch(`${baseUrl}/api/blocking/status`, {
      signal: controller.signal,
    });
    clearTimeout(timer);
    return response.ok;
  } catch {
    return false;
  }
}
