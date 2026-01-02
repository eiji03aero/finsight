package session

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

const (
	sessionName = "finsight_session"
	userIDKey   = "user_id"
	emailKey    = "email"
	workspaceIDKey = "workspace_id"
)

var store *sessions.CookieStore

// InitStore initializes the session store with the provided secret key
func InitStore(secretKey string) *sessions.CookieStore {
	store = sessions.NewCookieStore([]byte(secretKey))

	// Session options
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   604800, // 7 days in seconds
		HttpOnly: true,
		Secure:   false, // Set to true in production (HTTPS only)
		SameSite: http.SameSiteLaxMode,
	}

	return store
}

// SetSession sets user session data
func SetSession(c *gin.Context, userID int, email string, workspaceID int) error {
	session, err := store.Get(c.Request, sessionName)
	if err != nil {
		return err
	}

	session.Values[userIDKey] = userID
	session.Values[emailKey] = email
	session.Values[workspaceIDKey] = workspaceID

	return session.Save(c.Request, c.Writer)
}

// GetSession retrieves user session data
func GetSession(c *gin.Context) (userID int, email string, workspaceID int, err error) {
	session, err := store.Get(c.Request, sessionName)
	if err != nil {
		return 0, "", 0, err
	}

	userIDVal, ok := session.Values[userIDKey]
	if !ok {
		return 0, "", 0, http.ErrNoCookie
	}

	emailVal, ok := session.Values[emailKey]
	if !ok {
		return 0, "", 0, http.ErrNoCookie
	}

	workspaceIDVal, ok := session.Values[workspaceIDKey]
	if !ok {
		return 0, "", 0, http.ErrNoCookie
	}

	return userIDVal.(int), emailVal.(string), workspaceIDVal.(int), nil
}
