import { useForm } from "@tanstack/react-form"
import { useNavigate } from "@tanstack/react-router"
import { signupSchema, type SignupFormData } from "@/pages/Signup/model/signupSchema"
import { useSignup } from "@/features/auth/model/useSignup"
import { Button } from "@/shared/shadcn/ui/button"
import { Input } from "@/shared/shadcn/ui/input"
import { Label } from "@/shared/shadcn/ui/label"
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/shared/shadcn/ui/card"
import { Alert, AlertDescription } from "@/shared/shadcn/ui/alert"

export function SignupForm() {
  const navigate = useNavigate()
  const signup = useSignup()

  const form = useForm({
    defaultValues: {
      email: "",
      password: "",
      workspaceName: "",
    } as SignupFormData,
    onSubmit: async ({ value }) => {
      try {
        await signup.mutateAsync(value)
        // On success, navigate to home page
        navigate({ to: "/" })
      } catch (error) {
        // Error is handled in the form UI
        console.error("Signup failed:", error)
      }
    },
  })

  return (
    <div className="flex min-h-screen items-center justify-center p-4 bg-gradient-to-br from-slate-50 to-slate-100 dark:from-slate-950 dark:to-slate-900">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle className="text-2xl">Create an account</CardTitle>
          <CardDescription>
            Enter your information to get started with your new workspace
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form
            onSubmit={(e) => {
              e.preventDefault()
              e.stopPropagation()
              form.handleSubmit()
            }}
            className="space-y-4"
          >
            {/* Email Field */}
            <form.Field
              name="email"
              validators={{
                onChange: signupSchema.shape.email,
              }}
            >
              {(field) => (
                <div className="space-y-2">
                  <Label htmlFor="email">Email</Label>
                  <Input
                    id="email"
                    type="email"
                    placeholder="you@example.com"
                    value={field.state.value}
                    onBlur={field.handleBlur}
                    onChange={(e) => field.handleChange(e.target.value)}
                    className={
                      field.state.meta.errors.length > 0
                        ? "border-destructive"
                        : ""
                    }
                  />
                  {field.state.meta.errors.length > 0 && (
                    <p className="text-sm text-destructive">
                      {typeof field.state.meta.errors[0] === "string"
                        ? field.state.meta.errors[0]
                        : field.state.meta.errors[0]?.message}
                    </p>
                  )}
                </div>
              )}
            </form.Field>

            {/* Password Field */}
            <form.Field
              name="password"
              validators={{
                onChange: signupSchema.shape.password,
              }}
            >
              {(field) => (
                <div className="space-y-2">
                  <Label htmlFor="password">Password</Label>
                  <Input
                    id="password"
                    type="password"
                    placeholder="••••••••"
                    value={field.state.value}
                    onBlur={field.handleBlur}
                    onChange={(e) => field.handleChange(e.target.value)}
                    className={
                      field.state.meta.errors.length > 0
                        ? "border-destructive"
                        : ""
                    }
                  />
                  {field.state.meta.errors.length > 0 && (
                    <p className="text-sm text-destructive">
                      {typeof field.state.meta.errors[0] === "string"
                        ? field.state.meta.errors[0]
                        : field.state.meta.errors[0]?.message}
                    </p>
                  )}
                </div>
              )}
            </form.Field>

            {/* Workspace Name Field */}
            <form.Field
              name="workspaceName"
              validators={{
                onChange: signupSchema.shape.workspaceName,
              }}
            >
              {(field) => (
                <div className="space-y-2">
                  <Label htmlFor="workspaceName">Workspace Name</Label>
                  <Input
                    id="workspaceName"
                    type="text"
                    placeholder="My Workspace"
                    value={field.state.value}
                    onBlur={field.handleBlur}
                    onChange={(e) => field.handleChange(e.target.value)}
                    className={
                      field.state.meta.errors.length > 0
                        ? "border-destructive"
                        : ""
                    }
                  />
                  {field.state.meta.errors.length > 0 && (
                    <p className="text-sm text-destructive">
                      {typeof field.state.meta.errors[0] === "string"
                        ? field.state.meta.errors[0]
                        : field.state.meta.errors[0]?.message}
                    </p>
                  )}
                </div>
              )}
            </form.Field>

            {/* Error Message */}
            {signup.isError && (
              <Alert variant="destructive">
                <AlertDescription>
                  {signup.error instanceof Error
                    ? signup.error.message
                    : "An error occurred. Please try again."}
                </AlertDescription>
              </Alert>
            )}

            {/* Submit Button */}
            <Button
              type="submit"
              disabled={signup.isPending}
              className="w-full"
              size="lg"
            >
              {signup.isPending ? "Creating account..." : "Sign Up"}
            </Button>
          </form>
        </CardContent>
        <CardFooter className="flex flex-col space-y-2">
          <p className="text-sm text-muted-foreground text-center">
            Already have an account?{" "}
            <a
              href="/login"
              className="text-primary underline-offset-4 hover:underline font-medium"
            >
              Sign in
            </a>
          </p>
        </CardFooter>
      </Card>
    </div>
  )
}
