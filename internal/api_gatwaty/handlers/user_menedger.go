package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	authpb "github.com/mmrgreenteaa/user-management-service/internal/gen/proto/auth"
	user_mendgerpb "github.com/mmrgreenteaa/user-management-service/internal/gen/proto/user_manegement"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (api *ApiGatway) RegistrateUser(c *gin.Context) {

	type user struct {
		Login string `json:"login"`
		Pass  string `json:"password"`
	}
	var useReq user

	if err := c.ShouldBindJSON(&useReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req := user_mendgerpb.UserRegistrationReq{
		Login: useReq.Login,
		Pass:  useReq.Pass,
	}

	_, err := api.UserMenedger.RegistrationUser(context.Background(), &req)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Fatalf("неизвестная ошибка: %v", err)
		}
		switch st.Code() {
		case codes.InvalidArgument:
			log.Println("Ошибка валидации:", st.Message())
		case codes.NotFound:
			log.Println("Ресурс не найден:", st.Message())
		case codes.PermissionDenied:
			log.Println("Доступ запрещён:", st.Message())
		case codes.Unavailable:
			log.Println("Сервис временно недоступен:", st.Message())
		case codes.DeadlineExceeded:
			log.Println("Таймаут запроса:", st.Message())
		default:
			log.Printf("Неожиданная ошибка (%v): %v", st.Code(), st.Message())
		}
		return
	}

}

func (apiGtwy *ApiGatway) EditUserLogin(c *gin.Context) {

	type UserReq struct {
		OldLogin string `json:"old_login"`
		NewLogin string `json:"new_login"`
	}

	var userReq UserReq
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed json values"})
		return
	}

	req := user_mendgerpb.UserLoginEditReq{
		OldLogin: userReq.OldLogin,
		NewLogin: userReq.NewLogin,
	}
	_, err := apiGtwy.UserMenedger.EditLogin(c.Request.Context(), &req)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Printf("Operation failed: %v", err) // или реальная ошибка
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			return
		}
		switch st.Code() {
		case codes.Internal:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			return
		case codes.InvalidArgument:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "the old login is not valid",
			})
			return

		case codes.Unavailable:
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "The service is temporarily unavailable:",
			})
			return
		}
	}
}

func (apiGtwy *ApiGatway) DeleteAccount(c *gin.Context) {

	ref, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "refresh cookie not found",
		})
	}
	c.SetCookie("refresh_token", "", -1, "/", "yourdomain.com", true, true)
	reqlogout := authpb.RefreshRequst{
		RefreshToken: ref,
	}

	_, err = apiGtwy.AuthServis.LogOut(c.Request.Context(), &reqlogout)
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
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "The service is temporarily unavailable:",
			})
			return
		}
	}

	type UserReq struct {
		Login string `json:"login"`
	}

	var userReq UserReq
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed json values"})
		return
	}

	req := user_mendgerpb.LoginReq{
		Login: userReq.Login,
	}
	_, err = apiGtwy.UserMenedger.DeleteUser(c.Request.Context(), &req)
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
				"error": "the login is not valid",
			})
			return

		case codes.Unavailable:
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "The service is temporarily unavailable:",
			})
			return
		}
	}

}
