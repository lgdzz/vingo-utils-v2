package oss

import (
	"github.com/lgdzz/vingo-utils-v2/db/mysql"
	"github.com/lgdzz/vingo-utils-v2/db/page"
	"github.com/lgdzz/vingo-utils-v2/vingo"
	"gorm.io/gorm"
)

type Resource struct {
	Id        int              `gorm:"primaryKey;column:id" json:"id"`
	OrgId     int              `gorm:"column:org_id" json:"orgId"`
	Type      string           `gorm:"column:type" json:"type"`                 // video-视频|audio-音频|image-图片|file-文件
	TakeNum   int64            `gorm:"column:take_num" json:"takeNum"`          // 引用次数，为0时可删除
	FileName  string           `gorm:"column:file_name" json:"fileName"`        // 文件名称
	FilePath  string           `gorm:"column:file_path" json:"filePath"`        // 源文件路径
	FileSize  int64            `gorm:"column:file_size" json:"fileSize"`        // 源文件大小
	FileType  string           `gorm:"column:file_type" json:"fileType"`        // 源文件类型
	Attr      any              `gorm:"column:attr;serializer:json" json:"attr"` // 文件属性
	CreatedAt *vingo.LocalTime `gorm:"column:created_at" json:"createdAt"`      // 上传时间
	UpdatedAt *vingo.LocalTime `gorm:"column:updated_at" json:"updatedAt"`      // 修改时间
	DeletedAt gorm.DeletedAt   `gorm:"column:deleted_at" json:"deletedAt"`      // 删除时间
}

func (s *Resource) CheckOrgId(orgId int) {
	if s.OrgId != orgId {
		panic("不允许跨组织操作素材")
	}
}

func (s *VideoAttr) CheckVideoLock() {
	if s.Lock {
		panic("视频正在处理中，禁止删除")
	}
}

func (s *Resource) TableName() string {
	return "resource"
}

type VideoAttr struct {
	Lock     bool            `json:"lock"`     // true-锁定|false-解锁（锁定时禁止删除）
	Cover    string          `json:"cover"`    // 封面图
	Duration int64           `json:"duration"` // 时长
	Ratio    string          `json:"ratio"`    // 分辨率
	Resource []VideoResource `json:"resource"`
}

type AudioAttr struct {
	Duration int64 `json:"duration"` // 时长
}

type ImageAttr struct {
}

type FileAttr struct {
}

type VideoResource struct {
	Name  string `json:"name"`
	Ratio string `json:"ratio"`
	Path  string `json:"path"`
}

type ResourceQuery struct {
	page.Limit
	Type     string `form:"type"`
	FileName string `form:"fileName"`
}

func GetResource(tx *gorm.DB, id uint) Resource {
	return mysql.FetchById[Resource](tx, id)
}

func (s *Resource) GetVideoAttr() (attr VideoAttr) {
	vingo.CustomOutput(s.Attr, &attr)
	return
}

func (s *Resource) GetAudioAttr() (attr AudioAttr) {
	vingo.CustomOutput(s.Attr, &attr)
	return
}

func (s *Resource) GetDuration() int64 {
	if s.Type == vingo.VIDEO_TYPE || s.Type == vingo.AUDIO_TYPE {
		return int64(s.Attr.(map[string]any)["duration"].(float64))
	} else {
		return 0
	}
}

type ResourceImage struct {
	Type     string `json:"type"`
	FileName string `json:"fileName"`
	FilePath string `json:"filePath"`
}

type ResourceFile struct {
	Type     string `json:"type"`
	FileName string `json:"fileName"`
	FilePath string `json:"filePath"`
}

type ResourceVideo struct {
	Type     string    `json:"type"`
	FileName string    `json:"fileName"`
	FilePath string    `json:"filePath"`
	Attr     VideoAttr `json:"attr"`
}

type ResourceAudio struct {
	Type     string    `json:"type"`
	FileName string    `json:"fileName"`
	FilePath string    `json:"filePath"`
	Attr     AudioAttr `json:"attr"`
}
