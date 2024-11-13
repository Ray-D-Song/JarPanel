package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/1Panel-dev/1Panel/backend/app/api/v1/helper"
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

type ServiceRecord struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	PrefixArgs      string `json:"prefixArgs"`
	SuffixArgs      string `json:"suffixArgs"`
	CreateTime      string `json:"createTime"`
	DeployTime      string `json:"deployTime"`
	CurrentVersion  string `json:"currentVersion"`
	PreviousVersion string `json:"previousVersion"`
}

func (b *BaseApi) NewService(c *gin.Context) {
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

	// 创建新的服务文件夹
	newDirName := uuid.New().String()
	newDirPath := filepath.Join(jarDir, newDirName)
	if err := os.MkdirAll(newDirPath, 0755); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	// jar 包命名为类似于 2024-01-01-12-00.jar
	jarFileName := time.Now().Format("2006-01-02-15-00") + ".jar"
	// 将 jar 包写入新的服务文件夹
	jarFilePath := filepath.Join(newDirPath, jarFileName)
	if err := c.SaveUploadedFile(file, jarFilePath); err != nil {
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
	var recordData []ServiceRecord
	if err := json.Unmarshal(record, &recordData); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	// 添加新记录
	recordData = append(recordData, ServiceRecord{
		Id:              newDirName,
		Name:            name,
		PrefixArgs:      "",
		SuffixArgs:      "",
		CreateTime:      time.Now().Format("2006-01-02 15:04:05"),
		DeployTime:      "",
		CurrentVersion:  jarFileName,
		PreviousVersion: "",
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
func (b *BaseApi) StartService(c *gin.Context) {
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

	var recordData []ServiceRecord
	if err := json.Unmarshal(records, &recordData); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	// 找到对应的jar包
	for _, record := range recordData {
		if record.Id == id {
			// 进入jar包目录
			jarDirPath := filepath.Join(jarDir, record.Id)
			// 组合启动命令
			command := fmt.Sprintf("cd %s && java -jar %s %s %s", jarDirPath, record.CurrentVersion, record.PrefixArgs, record.SuffixArgs)
			cmd := exec.Command("sh", "-c", command)
			cmd.Stderr = cmd.Stdout

			// 以后台方式启动进程
			if err := cmd.Start(); err != nil {
				helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
				return
			}

			// 将进程解绑，使其在后台独立运行
			if err := cmd.Process.Release(); err != nil {
				helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
				return
			}
			break
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
func (b *BaseApi) StopService(c *gin.Context) {
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

	var recordData []ServiceRecord
	if err := json.Unmarshal(records, &recordData); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	for _, record := range recordData {
		if record.Id == id {
			// 使用 jps -l获取所有 java 进程
			cmd := exec.Command("jps", "-l")
			output, err := cmd.Output()
			if err != nil {
				helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
				return
			}

			// 在输出中查找对应的 jar 包进程
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.Contains(line, record.CurrentVersion) {
					// 提取 PID
					pid := strings.Split(line, " ")[0]
					// 先尝试优雅停止
					killCmd := exec.Command("kill", pid)
					if err := killCmd.Run(); err != nil {
						// 如果优雅停止失败,再强制终止
						killCmd = exec.Command("kill", "-9", pid)
						if err := killCmd.Run(); err != nil {
							helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
							return
						}
					}
					// 等待进程终止
					maxRetries := 5
					for i := 0; i < maxRetries; i++ {
						// 检查进程是否还在运行
						checkCmd := exec.Command("kill", "-0", pid)
						if err := checkCmd.Run(); err != nil {
							// 进程已终止
							break
						}
						time.Sleep(time.Second)
					}
					break
				}
			}
			break
		}
	}

	helper.SuccessWithData(c, nil)
}

type ServiceStatus struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	PrefixArgs      string `json:"prefixArgs"`
	SuffixArgs      string `json:"suffixArgs"`
	CreateTime      string `json:"createTime"`
	DeployTime      string `json:"deployTime"`
	Status          string `json:"status"`
	CurrentVersion  string `json:"currentVersion"`
	PreviousVersion string `json:"previousVersion"`
}

// @Tags Jar
// @Summary Get jar status
// @Description 获取jar包运行状态
// @Accept json
// @Success 200 {object} ServiceStatus
// @Security ApiKeyAuth
// @Router /jars/status [post]
// @x-panel-log {"bodyKeys":[],"paramKeys":[],"BeforeFunctions":[],"formatZH":"获取jar包状态","formatEN":"Get jar status"}
func (b *BaseApi) GetServiceStatus(c *gin.Context) {
	// 从record.json中获取jar包信息
	records, err := os.ReadFile(filepath.Join(jarDir, "record.json"))
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	var recordData []ServiceRecord
	if err := json.Unmarshal(records, &recordData); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	// 使用 jps -l 获取所有java进程
	cmd := exec.Command("jps", "-l")
	output, err := cmd.Output()
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	os.WriteFile(filepath.Join(jarDir, "dev.log"), output, 0644)

	// 遍历recordData，检查每个jar包是否在运行, 返回Record + Running状态
	var statuses []ServiceStatus
	for _, record := range recordData {
		if strings.Contains(string(output), record.CurrentVersion) {
			statuses = append(statuses, ServiceStatus{Id: record.Id, Name: record.Name, PrefixArgs: record.PrefixArgs, SuffixArgs: record.SuffixArgs, CreateTime: record.CreateTime, DeployTime: record.DeployTime, Status: "running", CurrentVersion: record.CurrentVersion, PreviousVersion: record.PreviousVersion})
		} else {
			statuses = append(statuses, ServiceStatus{Id: record.Id, Name: record.Name, PrefixArgs: record.PrefixArgs, SuffixArgs: record.SuffixArgs, CreateTime: record.CreateTime, DeployTime: record.DeployTime, Status: "stopped", CurrentVersion: record.CurrentVersion, PreviousVersion: record.PreviousVersion})
		}
	}

	helper.SuccessWithData(c, statuses)
}

// @Tags Jar
// @Summary Get jar files
// @Description 获取jar包文件列表
// @Accept json
// @Success 200 {array} string
// @Security ApiKeyAuth
// @Router /jars/files [get]
func (b *BaseApi) GetServiceFileList(c *gin.Context) {
	id := c.Query("id")
	serviceDirPath := filepath.Join(jarDir, id)
	files, err := os.ReadDir(serviceDirPath)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
	}
	var fileList []string
	for _, file := range files {
		fileList = append(fileList, file.Name())
	}
	helper.SuccessWithData(c, fileList)
}

// @Tags Jar
// @Summary Download jar file
// @Description 下载jar包文件
// @Accept json
// @Success 200
// @Security ApiKeyAuth
// @Router /jars/download [get]
func (b *BaseApi) DownloadFile(c *gin.Context) {
	id := c.Query("id")
	name := c.Query("name")
	filePath := filepath.Join(jarDir, id, name)
	// 获取文件
	file, err := os.ReadFile(filePath)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	c.Data(http.StatusOK, "application/octet-stream", file)
}

func (b *BaseApi) UploadFile(c *gin.Context) {
	id := c.Query("id")
	name := c.Query("name")
	if id == "" || name == "" {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, nil)
		return
	}
	// 如果是 .jar 文件，则按照日期命名
	if filepath.Ext(name) == ".jar" {
		name = time.Now().Format("2006-01-02-15-00") + ".jar"
	}
	filePath := filepath.Join(jarDir, id, name)
	file, err := c.FormFile("file")
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, err)
		return
	}

	dst, err := os.Create(filePath)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	defer dst.Close()

	src, err := file.Open()
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	defer src.Close()

	if _, err := io.Copy(dst, src); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}
