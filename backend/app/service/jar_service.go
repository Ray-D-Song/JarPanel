package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/1Panel-dev/1Panel/backend/global"
	"github.com/google/uuid"
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

func (s *JarService) Upload(file *multipart.FileHeader) (string, error) {
	// 将 jar 包重命名为 uuid.jar
	fileName := fmt.Sprintf("%s.jar", uuid.New().String())
	dst := filepath.Join(jarDir, fileName)
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return fileName, err
}
