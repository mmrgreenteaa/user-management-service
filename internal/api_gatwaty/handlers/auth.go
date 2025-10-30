package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	authpb "github.com/mmrgreenteaa/user-management-service/internal/gen/proto/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (gatway *ApiGatway) LogIn(c *gin.Context) {

	if gatway == nil || gatway.AuthServis == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "service unavailable"})
		return
	}
	type user struct {
		Login string `json:"login"`
		Pass  string `json:"password"`
	}
	var useReq user

	if err := c.ShouldBindJSON(&useReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//log.Println(useReq)
	req := authpb.UserInfoRequest{
		Login:    useReq.Login,
		Password: useReq.Pass,
	}
	res, err := gatway.AuthServis.Login(context.Background(), &req)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Fatalf("неизвестная ошибка: %v", err)
		}
		switch st.Code() {
		case codes.Internal:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			log.Println("Здесь лшибка ")
			c.Abort()
		case codes.InvalidArgument:
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "invalid data"},
			)
		case codes.Unavailable:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "The service is temporarily unavailable:",
			})
			log.Println("Сервис временно недоступен:", st.Message())
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error"})
		}

	}
	c.Header("Authorization", res.AccessToken)
	c.SetCookie("refresh_token", res.RefreshToken,
		172800,      // maxAge: 2 дня (в секундах)
		"/",         // путь
		"localhost", // домен (или "")
		false,       // Secure — только по HTTPS
		false,       // HttpOnly — для защиты от XSS
	)
	c.JSON(http.StatusOK, gin.H{
		"message":      "login successful",
		"access_token": res.AccessToken,
	})
}

func (apiGtwy ApiGatway) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Требуется токен авторизации",
			})
			c.Abort()
			return
		}
		req := authpb.AccessRequest{
			Token: token,
		}
		_, err := apiGtwy.AuthServis.VerificatAccess(context.Background(), &req)
		if err != nil {
			log.Println(err)
			st, ok := status.FromError(err)
			if !ok {
				log.Fatalf("неизвестная ошибка: %v", err)
			}
			switch st.Code() {
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
				c.Abort()
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "The token is not valid",
				})
				c.Abort()
				log.Println("error valid:", st.Message())

			case codes.Unauthenticated:
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid or expired token",
				})

			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "The service is temporarily unavailable:",
				})
				log.Println("Сервис временно недоступен:", st.Message())
			}
		}

	}
}

func (apiGtwy *ApiGatway) RefreshToken(c *gin.Context) {

	access := c.GetHeader("Authorization")
	if access == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Требуется токен авторизации",
		})
		c.Abort()
		return
	}
	ref, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "refresh cookie not found",
		})
	}
	req := authpb.RefreshRequst{
		RefreshToken: ref,
	}
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", access)
	res, err := apiGtwy.AuthServis.RefreshToken(ctx, &req)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Fatalf("unknown error: %v", err)
		}
		switch st.Code() {
		case codes.Internal:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			c.Abort()
		case codes.InvalidArgument:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "The tokens is not valid",
			})
			c.Abort()
			log.Println("error valid:", st.Message())

		case codes.Unauthenticated:
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired tokens",
			})

		case codes.Unavailable:
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "The service is temporarily unavailable:",
			})
			log.Println("Сервис временно недоступен:", st.Message())
		}

	}
	c.Header("Authorization", res.AccessToken)
	c.SetCookie("refresh_token", res.RefreshToken,
		172800,      // maxAge: 2 дня (в секундах)
		"/",         // путь
		"localhost", // домен (или "")
		false,       // Secure — только по HTTPS
		false,       // HttpOnly — для защиты от XSS
	)

}

func (apiGtwy *ApiGatway) LogOut(c *gin.Context) {

	ref, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "refresh cookie not found",
		})
		return
	}
	if ref == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "empty refresh token",
		})
		return
	}
	c.SetCookie("refresh_token", "", -1, "/", "yourdomain.com", true, true)
	req := authpb.RefreshRequst{
		RefreshToken: ref,
	}
	_, err = apiGtwy.AuthServis.LogOut(c.Request.Context(), &req)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			return
		}
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

}
