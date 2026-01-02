// Base API URL from environment variable
export const API_URL = import.meta.env.VITE_API_URL || "http://localhost:28080"

// Base fetch wrapper with credentials support
export async function apiClient<T>(
  endpoint: string,
  options?: RequestInit
): Promise<T> {
  const response = await fetch(`${API_URL}${endpoint}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options?.headers,
    },
    credentials: "include", // Include cookies for session management
  })

  if (!response.ok) {
    const error = await response.json().catch(() => ({
      error: "An error occurred",
    }))
    throw new Error(error.error || `HTTP ${response.status}`)
  }

  return response.json()
}
