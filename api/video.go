package api

import (
	resp "dy/middleware"
	"dy/models"
	"dy/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func GetVideosByActionType(c *gin.Context) {
	actionTypes := c.Query("actionTypes")
	userId, _ := strconv.Atoi(c.Request.Header.Get("Id"))
	actions := strings.Split(actionTypes, "|")
	videos := models.SearchVideosByActionTypes(uint(userId), actions)
	vds := make([]models.VideoDto, len(videos))
	for i, video := range videos {
		vds[i] = models.GetAllCountsForVideo(video)
	}
	c.JSON(http.StatusOK, resp.OK.WithData(vds))
	//c.JSON(http.StatusOK, resp.OK.WithData(videos))
}

func GetFallowVideosByOwnerId(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.Header.Get("Id"))
	// 根据当前用户id找到关注的用户
	fallows := models.GetUsersByOwnerIdAndType(uint(userId), 3)

	videos := make([]models.Video, 0)
	// 获取关注的人下面的视频
	for _, fallow := range fallows {
		videos = append(videos, models.GetVideosByUserId(fallow.ID)...)
	}
	vds := make([]models.VideoDto, len(videos))
	for i, video := range videos {
		vds[i] = models.GetAllCountsForVideo(video)
	}
	c.JSON(http.StatusOK, resp.OK.WithData(vds))
}

// ActionForVideo 用户对视频进行操作，点赞这些
func ActionForVideo(c *gin.Context) {
	videoId, err := strconv.Atoi(c.Query("id"))
	actionType := c.Query("actionType")
	userId, _ := strconv.Atoi(c.Request.Header.Get("Id"))
	if err != nil || len(actionType) == 0 {
		c.JSON(http.StatusOK, resp.ErrParam.WithMsg("参数错误"))
		return
	}
	var uv = models.UserAndVideo{
		VideoId:    uint(videoId),
		UserId:     uint(userId),
		ActionType: actionType,
	}
	if er := models.AddActionsForVideo(uv); er != nil {
		c.JSON(http.StatusOK, resp.Err)
		return
	}
	c.JSON(http.StatusOK, resp.OK)
}

func SaveVideo(c *gin.Context) {
	video := models.Video{}
	if err := c.ShouldBind(&video); err != nil {
		c.JSON(http.StatusOK, resp.ErrParam)
		return
	}

	if err := models.CreateVideo(video); err != nil {
		c.JSON(http.StatusOK, resp.Err.WithMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, resp.OK.WithMsg("保存成功"))
}
func GetVideosByState(c *gin.Context) {
	state, err := strconv.Atoi(c.Query("state"))
	id, err := strconv.Atoi(c.Request.Header.Get("Id"))
	if err != nil {
		c.JSON(http.StatusOK, resp.ErrParam)
	}
	videos := models.GetVideosByUserIdAndState(uint(id), uint(state))
	vds := make([]models.VideoDto, len(videos))
	for i, video := range videos {
		vds[i] = models.GetAllCountsForVideo(video)
	}
	c.JSON(http.StatusOK, resp.OK.WithData(vds))
}
func getVideoDto(videos []models.Video) []models.VideoDto {
	vds := make([]models.VideoDto, len(videos))
	for i, video := range videos {
		vds[i] = models.GetAllCountsForVideo(video)
	}
	return vds
}
func SaveDisplayVideoHistory(c *gin.Context) {
	videoId, _ := strconv.Atoi(c.Query("videoId"))
	userId, _ := strconv.Atoi(c.Request.Header.Get("Id"))
	uv := models.UserAndVideo{
		VideoId:    uint(videoId),
		UserId:     uint(userId),
		ActionType: "play",
	}
	utils.DB.Create(uv)
	c.JSON(http.StatusOK, resp.OK)
}
func GetVideosByVisible(c *gin.Context) {
	id, _ := strconv.Atoi(c.Request.Header.Get("Id"))
	visible, _ := strconv.Atoi(c.Query("visible"))
	videos := make([]models.Video, 0)
	utils.DB.Where("user_id = ? and visible = ?", id, visible).Find(&videos)
	c.JSON(http.StatusOK, resp.OK.WithData(getVideoDto(videos)))
}
