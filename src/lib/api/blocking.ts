import { apiRequest } from "./client";
import type { BlockingStatus } from "$lib/types/api";

export async function getBlockingStatus(): Promise<BlockingStatus> {
  return apiRequest<BlockingStatus>("/api/blocking/status");
}

export async function enableBlocking(): Promise<void> {
  await apiRequest("/api/blocking/enable", { method: "GET" });
}

export async function disableBlocking(
  duration?: string,
  groups?: string[],
): Promise<void> {
  const params = new URLSearchParams();
  if (duration) params.set("duration", duration);
  if (groups?.length) params.set("groups", groups.join(","));

  const query = params.toString();
  await apiRequest(`/api/blocking/disable${query ? `?${query}` : ""}`, {
    method: "GET",
  });
}
