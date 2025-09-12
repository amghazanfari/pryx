package auth

import (
	"crypto/subtle"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"
	"pryx/internal/models"
)

func Middleware(db *gorm.DB, requireScope string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authz := r.Header.Get("Authorization")
			if !strings.HasPrefix(authz, "Bearer ") {
				http.Error(w, `{"error":"missing bearer token"}`, http.StatusUnauthorized)
				return
			}
			plain := strings.TrimSpace(strings.TrimPrefix(authz, "Bearer "))
			if len(plain) < 16 {
				http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
				return
			}
			prefix := plain[:8]
			hash := HashAPIKey(plain)

			var key models.APIKey
			if err := db.WithContext(r.Context()).
				Where("prefix = ? AND hash = ? AND revoked = false", prefix, hash).
				First(&key).Error; err != nil {
				http.Error(w, `{"error":"invalid or revoked token"}`, http.StatusUnauthorized)
				return
			}

			if requireScope != "" && !HasScope(key.Scopes, requireScope) {
				http.Error(w, `{"error":"forbidden: missing scope"}`, http.StatusForbidden)
				return
			}

			now := time.Now()
			_ = db.Model(&key).Update("last_used_at", &now).Error

			ctx := WithAuth(r.Context(), AuthContext{
				UserID: key.UserID, APIKeyID: key.ID, Scopes: key.Scopes,
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func jsonErr(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// SharedSecretMiddleware protects admin routes with a single shared secret.
// Provide the secret via env (e.g., ADMIN_SECRET) and clients must send:
//
//	Authorization: Bearer <secret>
//
// or
//
//	X-Admin-Secret: <secret>
func SharedSecretMiddleware(secret string) func(http.Handler) http.Handler {
	secret = strings.TrimSpace(secret)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if secret == "" {
				jsonErr(w, http.StatusForbidden, "admin disabled: no secret configured")
				return
			}
			got := ""
			if h := r.Header.Get("Authorization"); strings.HasPrefix(h, "Bearer ") {
				got = strings.TrimSpace(strings.TrimPrefix(h, "Bearer "))
			}
			if got == "" {
				got = strings.TrimSpace(r.Header.Get("X-Admin-Secret"))
			}
			// Constant-time compare to avoid timing leaks
			if subtle.ConstantTimeCompare([]byte(got), []byte(secret)) != 1 {
				jsonErr(w, http.StatusUnauthorized, "unauthorized")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func IPAllowlistMiddlewareFromEnv() func(http.Handler) http.Handler {
	raw := strings.TrimSpace(os.Getenv("ADMIN_IP_ALLOWLIST"))
	if raw == "" {
		return func(next http.Handler) http.Handler { return next }
	}
	var nets []*net.IPNet
	for _, part := range strings.Split(raw, ",") {
		_, n, err := net.ParseCIDR(strings.TrimSpace(part))
		if err == nil && n != nil {
			nets = append(nets, n)
		}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host, _, _ := net.SplitHostPort(r.RemoteAddr)
			ip := net.ParseIP(host)
			allowed := false
			for _, n := range nets {
				if n.Contains(ip) {
					allowed = true
					break
				}
			}
			if !allowed {
				jsonErr(w, http.StatusForbidden, "forbidden: ip not allowed")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
