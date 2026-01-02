import { useMutation } from "@tanstack/react-query"
import { apiClient } from "@/shared/api/client"
import type { User } from "@/entities/user/model/types"
import type { Workspace } from "@/entities/workspace/model/types"

// Signup request payload
export interface SignupRequest {
  email: string
  password: string
  workspaceName: string
}

// Signup response payload
export interface SignupResponse {
  user: User
  workspace: Workspace
  message: string
}

export function useSignup() {
  return useMutation({
    mutationFn: async (data: SignupRequest): Promise<SignupResponse> => {
      return apiClient<SignupResponse>("/api/auth/signup", {
        method: "POST",
        body: JSON.stringify(data),
      })
    },
  })
}
