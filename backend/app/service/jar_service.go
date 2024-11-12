package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/1Panel-dev/1Panel/backend/global"
)

type JarService struct {
	id         string
	name       string
	fileName   string
	createTime time.Time
}

func NewJarService() *JarService {
	return &JarService{}
}

var jarDir = filepath.Join(global.CONF.System.BaseDir, "uploads", "jars")

func (s *JarService) Upload(file *multipart.FileHeader) error {
	// 使用 os.Stat 检查文件在服务器上是否存在
	if _, err := os.Stat(filepath.Join(jarDir, file.Filename)); err == nil {
		return fmt.Errorf("jar包 %s 已存在", file.Filename)
	}

	// 创建上传目录
	if err := os.MkdirAll(jarDir, 0755); err != nil {
		return err
	}

	// 保存文件
	dst := filepath.Join(jarDir, file.Filename)
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
