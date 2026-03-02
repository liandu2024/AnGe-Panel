package models

type File struct {
	BaseModel
	Src      string `json:"src"`
	UserId   uint   `json:"userId"`
	FileName string `json:"fileName" gorm:"varchar(255)"`           // 文件名
	Method   int    `gorm:"int(5)" json:"method"`                   // 上传方式
	Ext      string `gorm:"varchar(255)" json:"ext"`                // 扩展名
	FileType string `gorm:"varchar(32);default:''" json:"fileType"` // icon | wallpaper | ''(unknown)
}

// 添加一个文件记录
func (m *File) AddFile(userId uint, fileName, ext, src, fileType string) (File, error) {
	file := File{
		UserId:   userId,
		FileName: fileName,
		Src:      src,
		Ext:      ext,
		FileType: fileType,
	}
	err := Db.Create(&file).Error

	return file, err
}
