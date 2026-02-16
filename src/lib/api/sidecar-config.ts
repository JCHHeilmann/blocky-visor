import { sidecarRequest } from "./sidecar";

export async function fetchConfig(): Promise<string> {
  return sidecarRequest<string>("/api/config");
}

export async function saveConfig(
  yaml: string,
): Promise<{ status: string; backup: string }> {
  return sidecarRequest<{ status: string; backup: string }>("/api/config", {
    method: "PUT",
    body: yaml,
    headers: { "Content-Type": "text/plain" },
  });
}
