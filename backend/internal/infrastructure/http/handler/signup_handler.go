package handler

import (
	"net/http"

	"backend/internal/application/usecase"
	"backend/internal/infrastructure/session"

	"github.com/gin-gonic/gin"
)

type SignupHandler struct {
	signupUseCase *usecase.SignupUseCase
}

func NewSignupHandler(signupUseCase *usecase.SignupUseCase) *SignupHandler {
	return &SignupHandler{signupUseCase: signupUseCase}
}

type SignupRequest struct {
	Email         string `json:"email" binding:"required"`
	Password      string `json:"password" binding:"required"`
	WorkspaceName string `json:"workspaceName" binding:"required"`
}

type UserResponse struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type WorkspaceResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type SignupResponse struct {
	User      UserResponse      `json:"user"`
	Workspace WorkspaceResponse `json:"workspace"`
	Message   string            `json:"message"`
}

type ErrorResponse struct {
	Error   string                 `json:"error"`
	Code    string                 `json:"code"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Signup handles user signup requests
func (h *SignupHandler) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
			Code:  "VALIDATION_ERROR",
			Details: map[string]interface{}{
				"message": err.Error(),
			},
		})
		return
	}

	// Execute signup usecase
	user, workspace, err := h.signupUseCase.Execute(
		c.Request.Context(),
		req.Email,
		req.Password,
		req.WorkspaceName,
	)
	if err != nil {
		// Check for duplicate email error
		if err.Error() == "email already registered" {
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: "Email already registered",
				Code:  "EMAIL_EXISTS",
				Details: map[string]interface{}{
					"field":   "email",
					"message": "An account with this email already exists",
				},
			})
			return
		}

		// Other validation errors
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
			Code:  "VALIDATION_ERROR",
		})
		return
	}

	// Set session
	if err := session.SetSession(c, user.ID, user.Email, workspace.ID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to create session",
			Code:  "INTERNAL_ERROR",
		})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, SignupResponse{
		User: UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		},
		Workspace: WorkspaceResponse{
			ID:        workspace.ID,
			Name:      workspace.Name,
			CreatedAt: workspace.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: workspace.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		},
		Message: "Account created successfully",
	})
}

// GetSession returns the current session information
func (h *SignupHandler) GetSession(c *gin.Context) {
	userID, email, workspaceID, err := session.GetSession(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Not authenticated",
			Code:  "UNAUTHENTICATED",
			Details: map[string]interface{}{
				"message": "No active session found",
			},
		})
		return
	}

	// In a real implementation, you would fetch user and workspace from the database
	// For now, return basic session info
	c.JSON(http.StatusOK, gin.H{
		"authenticated": true,
		"user": gin.H{
			"id":    userID,
			"email": email,
		},
		"workspace": gin.H{
			"id": workspaceID,
		},
	})
}
