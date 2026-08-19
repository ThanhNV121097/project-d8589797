// Mock data-access layer for "Display Hello Word from database".
// Success shape matches docs/architecture/services.md §3.1 (GET /api/v1/content).
// A non-2xx response uses the §2.3 error envelope; the component maps any
// rejection to its error state. This file is deleted when the real API replaces it.

export interface ContentResponse {
  value: string;
}

const MOCK_LATENCY_MS = 700;

export async function fetchContent(): Promise<ContentResponse> {
  // Simulate a round-trip so the loading state is observable.
  await new Promise((resolve) => setTimeout(resolve, MOCK_LATENCY_MS));
  // Simulated value from the seeded DB row, served by the API — not UI copy.
  return { value: "Hello Word" };
}
