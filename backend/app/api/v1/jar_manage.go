package v1

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/1Panel-dev/1Panel/backend/app/api/v1/helper"
	"github.com/1Panel-dev/1Panel/backend/app/service"
	"github.com/1Panel-dev/1Panel/backend/constant"
	"github.com/1Panel-dev/1Panel/backend/global"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var jarDir = filepath.Join(global.CONF.System.BaseDir, "uploads", "jars")

func init() {
	// 确保 uploads 目录存在
	if _, err := os.Stat(filepath.Join(global.CONF.System.BaseDir, "uploads")); os.IsNotExist(err) {
		os.MkdirAll(filepath.Join(global.CONF.System.BaseDir, "uploads"), 0755)
	}
	// 确保 jars 目录存在
	if _, err := os.Stat(jarDir); os.IsNotExist(err) {
		os.MkdirAll(jarDir, 0755)
	}
	// 确保 record.json 文件存在
	if _, err := os.Stat(filepath.Join(jarDir, "record.json")); os.IsNotExist(err) {
		os.Create(filepath.Join(jarDir, "record.json"))
		os.WriteFile(filepath.Join(jarDir, "record.json"), []byte("[]"), 0644)
	}
}

type JarRecord struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	FileName   string `json:"fileName"`
	CreateTime string `json:"createTime"`
}

// @Tags Jar
// @Summary Upload jar
// @Description 上传jar包
// @Accept multipart/form-data
// @Param file formData file true "jar包文件"
// @Param name formData string true "应用名称"
// @Success 200
// @Security ApiKeyAuth
// @Router /jars/upload [post]
// @x-panel-log {"bodyKeys":["name"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"上传jar包 [name]","formatEN":"Upload jar [name]"}
func (b *BaseApi) UploadJar(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, err)
		return
	}
	name := c.PostForm("name")
	if name == "" {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, nil)
		return
	}

	if filepath.Ext(file.Filename) != ".jar" {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, nil)
		return
	}

	jarService := service.NewJarService()
	if err := jarService.Upload(file); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	// 读取 record.json 文件
	record, err := os.ReadFile(filepath.Join(jarDir, "record.json"))
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	// 解析 record.json 文件
	var recordData []JarRecord
	if err := json.Unmarshal(record, &recordData); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	// 添加新记录
	recordData = append(recordData, JarRecord{
		Id:         uuid.New().String(),
		Name:       name,
		FileName:   file.Filename,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	})

	// 写入 record.json 文件
	recordBytes, err := json.Marshal(recordData)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	if err := os.WriteFile(filepath.Join(jarDir, "record.json"), recordBytes, 0644); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	helper.SuccessWithData(c, nil)
}

// @Tags Jar
// @Summary Start jar
// @Description 启动jar包
// @Accept json
// @Param request body request.JarOperate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /jars/start [post]
// @x-panel-log {"bodyKeys":["name"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"启动jar包 [name]","formatEN":"Start jar [name]"}
func (b *BaseApi) StartJar(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, nil)
		return
	}
	// 读取 record.json 文件
	records, err := os.ReadFile(filepath.Join(jarDir, "record.json"))
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	var recordData []JarRecord
	if err := json.Unmarshal(records, &recordData); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	// 找到对应的jar包
	for _, record := range recordData {
		if record.Id == id {
			// 启动jar包
			cmd := exec.Command("java", "-jar", filepath.Join(jarDir, record.FileName))
			if err := cmd.Run(); err != nil {
				helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
				return
			}
		}
	}

	helper.SuccessWithData(c, nil)
}

// @Tags Jar
// @Summary Stop jar
// @Description 停止jar包
// @Accept json
// @Param request body request.JarOperate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /jars/stop [post]
// @x-panel-log {"bodyKeys":["name"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"停止jar包 [name]","formatEN":"Stop jar [name]"}
func (b *BaseApi) StopJar(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, nil)
		return
	}

	// 读取 record.json 文件
	records, err := os.ReadFile(filepath.Join(jarDir, "record.json"))
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	var recordData []JarRecord
	if err := json.Unmarshal(records, &recordData); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	// 找到对应的jar包
	for _, record := range recordData {
		if record.Id == id {
			// 停止jar包
			cmd := exec.Command("kill", "-9", record.Name)
			if err := cmd.Run(); err != nil {
				helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
				return
			}
		}
	}

	helper.SuccessWithData(c, nil)
}

// @Tags Jar
// @Summary Delete jar
// @Description 删除jar包
// @Accept json
// @Param request body request.JarDelete true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /jars/delete [post]
// @x-panel-log {"bodyKeys":["name"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"删除jar包 [name]","formatEN":"Delete jar [name]"}
func (b *BaseApi) DeleteJar(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, nil)
		return
	}

	// 读取 record.json 文件
	records, err := os.ReadFile(filepath.Join(jarDir, "record.json"))
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	var recordData []JarRecord
	if err := json.Unmarshal(records, &recordData); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	// 删除对应的jar包
	for _, record := range recordData {
		if record.Id == id {
			os.Remove(filepath.Join(jarDir, record.FileName))
		}
	}

	// 写入 record.json 文件
	if _, err := json.Marshal(recordData); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	helper.SuccessWithData(c, nil)
}

type JarStatus struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	FileName   string `json:"fileName"`
	CreateTime string `json:"createTime"`
	Status     string `json:"status"`
}

// @Tags Jar
// @Summary Get jar status
// @Description 获取jar包运行状态
// @Accept json
// @Success 200 {object} JarStatus
// @Security ApiKeyAuth
// @Router /jars/status [post]
// @x-panel-log {"bodyKeys":[],"paramKeys":[],"BeforeFunctions":[],"formatZH":"获取jar包状态","formatEN":"Get jar status"}
func (b *BaseApi) GetJarStatus(c *gin.Context) {
	// 从record.json中获取jar包信息
	records, err := os.ReadFile(filepath.Join(jarDir, "record.json"))
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	var recordData []JarRecord
	if err := json.Unmarshal(records, &recordData); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	cmd := exec.Command("ps", "-ef")
	output, err := cmd.Output()
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	// 遍历recordData，检查每个jar包是否在运行, 返回Record + Running状态
	var statuses []JarStatus
	for _, record := range recordData {
		if strings.Contains(string(output), record.Name) {
			statuses = append(statuses, JarStatus{Id: record.Id, Name: record.Name, FileName: record.FileName, CreateTime: record.CreateTime, Status: "running"})
		} else {
			statuses = append(statuses, JarStatus{Id: record.Id, Name: record.Name, FileName: record.FileName, CreateTime: record.CreateTime, Status: "stopped"})
		}
	}

	helper.SuccessWithData(c, statuses)
}
