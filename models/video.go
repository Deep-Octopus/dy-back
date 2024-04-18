package models

import (
	"dy/utils"
	"errors"
	"gorm.io/gorm"
	"strconv"
)

type Video struct {
	gorm.Model
	UserID   uint   `json:"userID"`
	Caption  string `json:"caption"`
	Cover    string `json:"cover"` //封面路径
	FilePath string `json:"filePath"`
	State    uint   `json:"state"`
	Visible  uint   `json:"visible"` //是否公开
	Desc     string `json:"desc"`
	// 其他视频信息字段
}
type VideoDto struct {
	Video
	Plays     int64 `json:"plays"`     //播放量
	Likes     int64 `json:"likes"`     //点赞数
	Comments  int64 `json:"comments"`  //评论量
	Collects  int64 `json:"collects"`  //收藏量
	Transmits int64 `json:"transmits"` //转发量
}
type UserAndVideo struct {
	gorm.Model
	UserId     uint   `json:"userId"`
	VideoId    uint   `json:"videoId"`
	ActionType string `json:"actionType"`
	Desc       string `json:"desc"`
}

func GetAllCountsForVideo(video Video) VideoDto {
	vd := VideoDto{
		Video:     video,
		Plays:     GetCountForVideoByActionType(video.ID, "play"),
		Likes:     GetCountForVideoByActionType(video.ID, "like"),
		Comments:  GetCountForVideoByActionType(video.ID, "comment"),
		Collects:  GetCountForVideoByActionType(video.ID, "collect"),
		Transmits: GetCountForVideoByActionType(video.ID, "transmit"),
	}
	return vd
}
func GetCountForVideoByActionType(videoId uint, actionType string) int64 {
	var cnt int64
	utils.DB.Model(&UserAndVideo{}).
		Where("video_id = ? AND action_type = ?", videoId, actionType).
		Count(&cnt)
	return cnt
}

func GetVideosByUserIdAndState(userId, state uint) []Video {
	var videos = make([]Video, 0)
	utils.DB.Where("user_id = ? and state = ?", userId, state).Find(&videos)
	return videos
}
func CreateVideo(video Video) error {
	if len(video.Caption) == 0 {
		video.Caption = "抖音视频@" + strconv.Itoa(int(video.ID))
	}
	if err := utils.DB.Create(&video).Error; err != nil {
		return errors.New("新增用户失败")
	}
	return nil
}
func GetVideosByUserId(userId uint) []Video {
	var videos = make([]Video, 0)
	utils.DB.Where("user_id = ?", userId).Find(&videos)
	return videos
}

func GetVideosByVideoId(videoId uint) []Video {
	var videos = make([]Video, 0)
	utils.DB.Where("id = ? and visible = 1", videoId).Find(&videos)
	return videos
}

func AddActionsForVideo(uv UserAndVideo) error {
	return utils.DB.Create(&uv).Error
}
func SearchVideosByActionTypes(userId uint, actions []string) []Video {
	uvs := make([]UserAndVideo, 0)
	utils.DB.Where("user_id = ? and action_type in ?", userId, actions).Find(&uvs)

	videos := make([]Video, 0)
	for _, uv := range uvs {
		videos = append(videos, GetVideosByVideoId(uv.VideoId)...)
	}
	return videos
}
