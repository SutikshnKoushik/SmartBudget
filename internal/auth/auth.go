package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
	store             *sessions.CookieStore
	sessionName       = "smart-budget-session"
)

type UserInfo struct {
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

func Init() {
	// Initialize OAuth config
	googleOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	// Initialize session store with a random key
	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		// Generate a random key if not provided
		key := make([]byte, 32)
		rand.Read(key)
		sessionKey = base64.StdEncoding.EncodeToString(key)
		fmt.Println("WARNING: Using randomly generated session key. Set SESSION_KEY environment variable for production.")
	}
	store = sessions.NewCookieStore([]byte(sessionKey))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}
}

func GenerateStateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func GetLoginURL(state string) string {
	return googleOauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func ExchangeCode(code string) (*oauth2.Token, error) {
	return googleOauthConfig.Exchange(context.Background(), code)
}

func GetUserInfo(token *oauth2.Token) (*UserInfo, error) {
	client := googleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func SaveSession(w http.ResponseWriter, r *http.Request, userInfo *UserInfo) error {
	session, _ := store.Get(r, sessionName)
	session.Values["email"] = userInfo.Email
	session.Values["name"] = userInfo.Name
	session.Values["picture"] = userInfo.Picture
	session.Values["authenticated"] = true
	session.Values["login_time"] = time.Now().Unix()
	return session.Save(r, w)
}

func GetSession(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, sessionName)
}

func IsAuthenticated(r *http.Request) bool {
	session, err := store.Get(r, sessionName)
	if err != nil {
		return false
	}
	auth, ok := session.Values["authenticated"].(bool)
	return ok && auth
}

func ClearSession(w http.ResponseWriter, r *http.Request) error {
	session, _ := store.Get(r, sessionName)
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1
	return session.Save(r, w)
}

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(r) {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next(w, r)
	}
}