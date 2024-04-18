package router

import (
	"dy/api"
	"dy/middleware"
	"dy/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/basic/docs"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())
	//r.Static(utils.CONF)
	r.Static("/assets", "./assets")
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/index", service.GetIndex)

	r.GET("/user/getUserList", middleware.JWTHandler(), api.GetUserList)
	r.GET("/user/getUser", middleware.JWTHandler(), api.GetUser)
	r.GET("/user/getFriends", middleware.JWTHandler(), api.GetFriends)
	r.GET("/user/getGroups", middleware.JWTHandler(), api.GetGroups)
	r.GET("/user/getListMessage", middleware.JWTHandler(), api.GetListMessage)
	r.POST("/user/createUser", api.CreateUser)
	r.POST("/user/deleteUser", middleware.JWTHandler(), api.DeleteUser)
	r.POST("/user/updateUser", middleware.JWTHandler(), api.UpdateUser)
	r.GET("/forgotPassword", api.ForgotPassword)
	r.GET("/getUsers", middleware.JWTHandler(), api.GetUsers)
	r.POST("/login", api.Login)
	//验证码
	r.GET("/getCode", api.SendVerificationEmail)
	r.GET("/verifyCode", api.VerifyCode)

	//发送消息
	r.GET("/message/sendMsg", api.SendMsg)
	r.GET("/message/sendUserMsg", api.SendUserMsg)

	r.POST("/message/getMessages", api.GetP2PMessages)
	r.POST("/message/getGroupMessages", api.GetGroupMessages)
	//文件
	r.POST("/file/uploadFile", middleware.JWTHandler(), api.UploadFile)
	r.GET("/file/deleteFile", middleware.JWTHandler(), api.DeleteFile)

	//视频
	//获取关注视频
	r.POST("/fallow", middleware.JWTHandler(), api.GetFallowVideosByOwnerId)
	//对视频的操作
	r.GET("/actionForVideo", middleware.JWTHandler(), api.ActionForVideo)
	r.GET("/getVideosByActionType", middleware.JWTHandler(), api.GetVideosByActionType)
	r.POST("/saveVideo", middleware.JWTHandler(), api.SaveVideo)
	r.GET("/getVideoByState", middleware.JWTHandler(), api.GetVideosByState)
	r.GET("/saveDisplayVideo", middleware.JWTHandler(), api.SaveDisplayVideoHistory)
	r.GET("/getVideosByVisible", middleware.JWTHandler(), api.GetVideosByVisible)

	return r
}
