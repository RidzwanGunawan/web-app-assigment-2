package middleware

import (
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Auth() gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		
		sessionToken, tokenErr := ctx.Cookie("session_token")
		requestContentType := ctx.GetHeader("Content-Type")

		if tokenErr != nil {
			if requestContentType == "application/json" {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse("unauthorized"))
			} else {
				ctx.AbortWithStatusJSON(http.StatusSeeOther, model.NewErrorResponse("token is missing"))
			}
			return
		}

		parsedToken, parseErr := jwt.Parse(sessionToken, func(token *jwt.Token) (interface{}, error) {
			if _, validMethod := token.Method.(*jwt.SigningMethodHMAC); !validMethod {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return model.JwtKey, nil
		})

		if parseErr != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse("invalid token"))
			return
		}

		claims, validClaims := parsedToken.Claims.(jwt.MapClaims)
		if !validClaims || !parsedToken.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse("unauthorized"))
			return
		}

		// Extract user ID from token claims
		claimData, _ := json.Marshal(claims)
		var userClaims model.Claims
		json.Unmarshal(claimData, &userClaims)

		ctx.Set("userID", userClaims.UserID)
		ctx.Next()
		// TODO: answer here
	})
}
