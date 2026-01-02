import type { Workspace } from "@/entities/workspace/model/types"

// User entity type
export interface User {
  id: number
  email: string
  createdAt: string // ISO 8601 timestamp
  updatedAt: string
}

// Session response payload
export interface SessionResponse {
  user: User
  workspace: Workspace
  authenticated: boolean
}
