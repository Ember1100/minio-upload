package routers

type SysUploadTask struct {
	ID             int64  `gorm:"column:id"`
	UploadID       string `gorm:"column:upload_id"`
	FileIdentifier string `gorm:"column:file_identifier"`
	FileName       string `gorm:"column:file_name"`
	BucketName     string `gorm:"column:bucket_name"`
	ObjectKey      string `gorm:"column:object_key"`
	TotalSize      int64  `gorm:"column:total_size"`
	ChunkSize      int64  `gorm:"column:chunk_size"`
	ChunkNum       int    `gorm:"column:chunk_num"`
}

func (SysUploadTask) TableName() string {
	return "sys_upload_task"
}
