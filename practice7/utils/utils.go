package utils

import (
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret []byte

func init() {
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		jwtSecret = []byte("super-secret")
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

func GenerateJWT(userID uuid.UUID, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
			return
		}
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		claims, _ := token.Claims.(jwt.MapClaims)
		c.Set("userID", claims["user_id"].(string))
		c.Set("role", claims["role"].(string))
		c.Next()
	}
}

// RoleMiddleware checks if the user has a specific role required
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}
		c.Next()
	}
}

// Rate Limiter
type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*Visitor
	limit    int
}

type Visitor struct {
	count    int
	lastSeen time.Time
}

func NewRateLimiter(limit int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		limit:    limit,
	}
	go rl.cleanupVisitors()
	return rl
}

func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(1 * time.Minute)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 1*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) LimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var key string

		// If JWT Auth already ran, we can use userID
		userID, exists := c.Get("userID")
		if exists {
			key = userID.(string)
		} else {
			// Try to parse JWT if Authorization header is provided
			tokenStr := c.GetHeader("Authorization")
			if tokenStr != "" {
				tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
				token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				})
				if token != nil && token.Valid {
					if claims, ok := token.Claims.(jwt.MapClaims); ok {
						if id, ok := claims["user_id"].(string); ok {
							key = id
							// Optionally set context
							c.Set("userID", id)
							c.Set("role", claims["role"].(string))
						}
					}
				}
			}
		}

		if key == "" {
			key = c.ClientIP()
		}

		rl.mu.Lock()
		visitor, exists := rl.visitors[key]
		if !exists {
			rl.visitors[key] = &Visitor{count: 1, lastSeen: time.Now()}
			rl.mu.Unlock()
			c.Next()
			return
		}

		if time.Since(visitor.lastSeen) > 1*time.Minute {
			visitor.count = 1
			visitor.lastSeen = time.Now()
			rl.mu.Unlock()
			c.Next()
			return
		}

		if visitor.count >= rl.limit {
			rl.mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests"})
			return
		}

		visitor.count++
		visitor.lastSeen = time.Now()
		rl.mu.Unlock()

		c.Next()
	}
}
