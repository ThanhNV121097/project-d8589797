// Data-access layer for "Display Hello Word from database".
// Talks to the real backend: GET /api/v1/content (services.md §3.1).
// The component maps any rejection, and an empty value, to its error state.

export interface ContentResponse {
  value: string;
}

const apiBase = process.env.NEXT_PUBLIC_API_URL ?? "/api";

export async function fetchContent(): Promise<ContentResponse> {
  const res = await fetch(`${apiBase}/v1/content`);
  if (!res.ok) {
    throw new Error(`content request failed: ${res.status}`);
  }
  return (await res.json()) as ContentResponse;
}
