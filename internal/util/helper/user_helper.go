package helper

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

//Constants
const UserId = "claims:user_id"
const UserName = "claims:user_name"
const UserEmail = "claims:user_email"

func GetUserId(context *gin.Context) int64 {
	userId, _ := strconv.ParseInt(context.GetString(UserId), 10, 64)
	return userId
}

func GetUserName(context *gin.Context) string {
	return context.GetString(UserName)
}

func GetUserEmail(context *gin.Context) string {
	return context.GetString(UserEmail)
}
