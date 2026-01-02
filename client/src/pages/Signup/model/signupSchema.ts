import { z } from "zod"

// Signup form validation schema
export const signupSchema = z.object({
  email: z.string().email("Invalid email format"),
  password: z.string().min(8, "Password must be at least 8 characters"),
  workspaceName: z.string().min(1, "Workspace name is required"),
})

export type SignupFormData = z.infer<typeof signupSchema>
