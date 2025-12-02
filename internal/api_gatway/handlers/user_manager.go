package handlers

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	user_mendgerpb "github.com/mmrgreenteaa/user-management-service/internal/gen/proto/user_manegement"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type UserReg struct {
	Login string `json:"login"  binding:"required" example:"my_user"`
	Pass  string `json:"password"  binding:"required" example:"123"`
	Email string `json:"email" binding:"required" example:"mrland@mail.com"`
}

// RegistrateUser godoc
//
//	@Summary        Registration new user in user_management service
//	@Tags           user_management
//	@Accept         json
//	@Produce        json
//	@Param          input body UserReg true "user's info"
//
// @Success 204 	"User added"
//
//	@Failure 	 500 {object} map[string]string "Internal Server Error"
//
// @Failure      400 {object} map[string]string "Invalid input"
// @Failure      504 {object} map[string]string "it is not possible to connect to the service"
// @Router         /users [post]
func (apiGtwy *ApiGatway) RegistrateUser(c *gin.Context) {

	var useReq UserReg

	if err := c.ShouldBindJSON(&useReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json requst incorrect"})
		return
	}
	if useReq.Email == "" || useReq.Login == "" || useReq.Pass == "" {

		c.JSON(http.StatusBadRequest, gin.H{"error": "email or login or email incorrect"})
		return

	}

	req := user_mendgerpb.UserRegistrationReq{
		Login:    useReq.Login,
		Password: useReq.Pass,
		Email:    useReq.Email,
	}

	ctx := metadata.NewOutgoingContext(c.Request.Context(), nil)
	_, err := apiGtwy.UserMenedger.RegistrationUser(ctx, &req)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			apiGtwy.logger.Error("non-gRPC error during registrate user", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			return
		}

		apiGtwy.logger.Warn("registrate user failed",
			slog.String("grpc_code", st.Code().String()),
			slog.String("message", st.Message()))
		switch st.Code() {
		case codes.Internal:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			return
		case codes.InvalidArgument:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid values",
			})
			return
		case codes.Unavailable:
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"error": "the service is temporarily unavailable:",
			})
			return
		}
	}

}

type LoginEdit struct {
	NewLogin string `json:"new_login" binding:"required" example:"my_usenewr"`
}

// EditUserLogin godoc
//
//	@Summary       Changes the user's login
//	@Tags           user_management
//	@Accept         json
//	@Produce        json
// 	@Security     ApiKeyAuth
//	@Param          input body LoginEdit true "new login"
// 	@Param accessToken header string true "jwt access token"
// 	@Success 204 	"User edit"
//	@Failure 	 500 {object} map[string]string "Internal Server Error"
// 	@Failure      400 {object} map[string]string "Invalid input"
// 	@Failure      504 {object} map[string]string "it is not possible to connect to the service"
// 	@Router         /users [patch]
func (apiGtwy *ApiGatway) EditUserLogin(c *gin.Context) {

	var userReq LoginEdit
	if err := c.ShouldBindJSON(&userReq); err != nil {
		apiGtwy.logger.Error("failed edit user login json login not correct")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed json values"})
		return
	}
	log.Println("контекст", c)

	value, ok := c.Get("user_id")
	if !ok {
		apiGtwy.logger.Error("failed get user_id in context")
		c.JSON(http.StatusInternalServerError, "failed get user_id")
		return
	}
	uid, ok := value.(string)
	if !ok {
		apiGtwy.logger.Error("failed convert user_id to string")
		c.JSON(http.StatusInternalServerError, "failed get user_id")
		return
	}
	if uid == "" || userReq.NewLogin == "" {
		apiGtwy.logger.Error("failed  userid or login is empty")
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	req := user_mendgerpb.UserLoginEditReq{
		UserId:   uid,
		NewLogin: userReq.NewLogin,
	}
	ctx := metadata.NewOutgoingContext(c.Request.Context(), nil)
	_, err := apiGtwy.UserMenedger.EditLogin(ctx, &req)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			apiGtwy.logger.Error("non-gRPC error during auth", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			return
		}

		apiGtwy.logger.Warn("auth failed",
			slog.String("grpc_code", st.Code().String()),
			slog.String("message", st.Message()))
		switch st.Code() {
		case codes.Internal:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			return
		case codes.InvalidArgument:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "user id invalid",
			})
			return
		case codes.Unavailable:
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"error": "the service is temporarily unavailable:",
			})
			return
		}
	}
}

// DeleteAccount godoc
//
//	@Summary       delete user info 
//	@Tags           user_management
//	@Accept         json
//	@Produce        json
// @Security     ApiKeyAuth
// @Param accessToken header string true "jwt access token"
// @Success 204 	"User delete"
//	@Failure 	 500 {object} map[string]string "Internal Server Error"
// @Failure      400 {object} map[string]string "Invalid input"
// @Failure      504 {object} map[string]string "it is not possible to connect to the service"
// @Router         /users [delete]
func (apiGtwy *ApiGatway) DeleteAccount(c *gin.Context) {

	value, ok := c.Get("user_id")
	if !ok {
		apiGtwy.logger.Error("failed get user_id in context")
		c.JSON(http.StatusInternalServerError, "failed get user_id")
		return
	}
	uid, ok := value.(string)
	if !ok || uid == "" {
		apiGtwy.logger.Error("failed convert user_id to string")
		c.JSON(http.StatusBadRequest, "failed get user_id")
		return
	}

	req := user_mendgerpb.DeleteReq{
		UserId: uid,
	}
	ctx := metadata.NewOutgoingContext(c.Request.Context(), nil)
	_, err := apiGtwy.UserMenedger.DeleteUser(ctx, &req)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			apiGtwy.logger.Error("non-gRPC error during delete user", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			return
		}

		apiGtwy.logger.Warn("delete user failed",
			slog.String("grpc_code", st.Code().String()),
			slog.String("message", st.Message()))
		switch st.Code() {
		case codes.Internal:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			return
		case codes.InvalidArgument:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "the user_id is not valid",
			})
			return

		case codes.Unavailable:
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"error": "The service is temporarily unavailable:",
			})
			return
		}
	}

}
