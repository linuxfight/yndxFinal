package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateToken(t *testing.T) {
	now := time.Now()
	validSecret := "super-secret-key"
	userID := "user-123"

	tests := []struct {
		name        string
		userID      string
		expires     time.Time
		secret      string
		wantErr     bool
		wantExpired bool
	}{
		{
			name:    "valid token",
			userID:  userID,
			expires: now.Add(1 * time.Hour),
			secret:  validSecret,
			wantErr: false,
		},
		{
			name:    "empty secret",
			userID:  userID,
			expires: now.Add(1 * time.Hour),
			secret:  "",
			wantErr: false,
		},
		{
			name:        "expired token",
			userID:      userID,
			expires:     now.Add(-1 * time.Hour),
			secret:      validSecret,
			wantErr:     false,
			wantExpired: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenString, err := generateToken(tt.userID, tt.expires, tt.secret)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, tokenString)

			// Parse and validate the token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(tt.secret), nil
			})

			if tt.wantExpired {
				assert.ErrorContains(t, err, "token is expired")
				return
			}

			require.NoError(t, err)
			assert.True(t, token.Valid)

			// Validate claims
			claims, ok := token.Claims.(jwt.MapClaims)
			require.True(t, ok, "invalid claims type")

			assert.Equal(t, tt.userID, claims["sub"])
			assert.IsType(t, float64(0), claims["iat"])
			assert.IsType(t, float64(0), claims["exp"])

			// Validate timestamps
			issuedAt := time.Unix(int64(claims["iat"].(float64)), 0)
			expiresAt := time.Unix(int64(claims["exp"].(float64)), 0)

			assert.WithinDuration(t, now, issuedAt, time.Second)
			assert.WithinDuration(t, tt.expires, expiresAt, time.Second)

			// Validate signing method
			assert.IsType(t, &jwt.SigningMethodHMAC{}, token.Method)
			assert.Equal(t, jwt.SigningMethodHS256, token.Method)
		})
	}
}
