package handlers

import (
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	authpb "github.com/mmrgreenteaa/user-management-service/internal/gen/proto/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	refreshTokenMaxAge = 48 * time.Hour
	cookiePath         = "/"
)

func (apgt *ApiGatway) LogIn(c *gin.Context) {

	type user struct {
		Login string `json:"login"`
		Pass  string `json:"password"`
	}
	var useReq user

	if err := c.ShouldBindJSON(&useReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect data"})
		return
	}

	if useReq.Login == "" || useReq.Pass == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect data"})
		return
	}
	req := authpb.UserInfoRequest{
		Login:     useReq.Login,
		Password:  useReq.Pass,
		UserAgent: c.Request.UserAgent(),
		Ip:        c.ClientIP(),
	}
	ctx := metadata.NewOutgoingContext(c.Request.Context(), nil)
	res, err := apgt.AuthServis.Login(ctx, &req)
	apgt.logger.Info("login attempt", slog.String("login", req.Login))
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			apgt.logger.Error("non-gRPC error during login", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			return
		}

		apgt.logger.Warn("login failed",
			slog.String("grpc_code", st.Code().String()),
			slog.String("message", st.Message()))
		switch st.Code() {
		case codes.Internal:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			return
		case codes.InvalidArgument:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "login or password not specific",
			})
			return

		case codes.Unauthenticated:
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid login or password",
			})
			return
		case codes.Unavailable:
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"error": "the service is temporarily unavailable:",
			})
			return
		}

	}
	c.Header("Authorization", res.AccessToken)
	c.SetCookie("refresh_token", res.RefreshToken,
		int(refreshTokenMaxAge.Seconds()),
		cookiePath,
		"localhost",
		false, // Secure — только по HTTPS
		false, // HttpOnly — для защиты от XSS
	)
	c.JSON(http.StatusOK, nil)
}

func (apgt ApiGatway) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			apgt.logger.Warn("there is no access token")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization token is required",
			})
			return
		}
		tokenStr := strings.TrimPrefix(token, "Bearer ")
		req := authpb.AccessRequest{
			Token: tokenStr,
		}
		ctx := metadata.NewOutgoingContext(c.Request.Context(), nil)
		res, err := apgt.AuthServis.VerifyAccess(ctx, &req)
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				apgt.logger.Error("non-gRPC error during auth", slog.String("error", err.Error()))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
				return
			}

			apgt.logger.Warn("auth failed",
				slog.String("grpc_code", st.Code().String()),
				slog.String("message", st.Message()))

			switch st.Code() {
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
				return

			case codes.Unauthenticated:
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "invalid or expired token",
				})
				return
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "the service is temporarily unavailable:",
				})
				return
			}
		}
		log.Println(res.UserId)
		c.Set("user_id", res.UserId)
		c.Next()

	}
}

func (apgt *ApiGatway) RefreshToken(c *gin.Context) {

	access := c.GetHeader("Authorization")
	if access == "" {
		apgt.logger.Warn("access token not found")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "authorization token is required",
		})
		return
	}
	ref, err := c.Cookie("refresh_token")
	if err != nil {
		apgt.logger.Error("refresh token cookie not found")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "refresh cookie not found",
		})
		return
	}
	req := authpb.RefreshRequest{
		RefreshToken: ref,
		UserAgent:    c.Request.UserAgent(),
		Ip:           c.ClientIP(),
	}
	ctx := metadata.AppendToOutgoingContext(c.Request.Context(), "Authorization", access)
	res, err := apgt.AuthServis.RefreshToken(ctx, &req)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			apgt.logger.Error("non-gRPC error during auth", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			return
		}

		apgt.logger.Warn("refresh failed",
			slog.String("grpc_code", st.Code().String()),
			slog.String("message", st.Message()))

		switch st.Code() {
		case codes.Internal:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			return
		case codes.InvalidArgument:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "The tokens is not valid",
			})
			return

		case codes.Unauthenticated:
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired tokens, tokens delete",
			})
			deleteCookieHandler(c)
			apgt.logger.Info("refresh token delete cookie")
			return

		case codes.Unavailable:
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "the service is temporarily unavailable:",
			})
			return
		}

	}
	c.Header("Authorization", res.AccessToken)
	c.SetCookie("refresh_token", res.RefreshToken,
		int(refreshTokenMaxAge), 
		cookiePath,              
		"localhost",             
		false,                 
		false,                   
	)

}

func (apgt *ApiGatway) LogOut(c *gin.Context) {

	access := c.GetHeader("Authorization")
	if access == "" {
		apgt.logger.Warn("access token not found")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "authorization token is required",
		})
		return
	}

	ref, err := c.Cookie("refresh_token")
	if err != nil {
		apgt.logger.Error("refresh token cookie not found")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "refresh cookie not found",
		})
		return
	}
	if ref == "" {
		apgt.logger.Error("empty refresh token ")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "empty refresh token",
		})
		return
	}

	req := authpb.LogoutRequest{
		RefreshToken: ref,
		UserAgent:    c.Request.UserAgent(),
		Ip:           c.ClientIP(),
	}
	ctx := metadata.AppendToOutgoingContext(c.Request.Context(), "Authorization", access)
	_, err = apgt.AuthServis.Logout(ctx, &req)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			apgt.logger.Error("non-gRPC error during auth", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			return
		}

		apgt.logger.Warn("logout failed",
			slog.String("grpc_code", st.Code().String()),
			slog.String("message", st.Message()))
		switch st.Code() {
		case codes.Internal:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			return
		case codes.InvalidArgument:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "the tokens is not valid",
			})
			return

		case codes.Unauthenticated:
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired tokens",
			})
			return
		case codes.Unavailable:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "The service is temporarily unavailable:",
			})
			return
		}

	}
	deleteCookieHandler(c)
}

func deleteCookieHandler(c *gin.Context) {
	c.SetCookie("refresh_token", "refresh_token",
		-1,          
		cookiePath,  
		"localhost", 
		false,      
		false,       
	)

}
