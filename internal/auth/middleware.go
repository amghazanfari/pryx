package auth

import (
	"net/http"
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
