<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { ElMessage } from 'element-plus'
import { getJarStatus, stopJar, startJar, deleteJar } from '@/api/modules/jar'

interface JarItem {
  id: string
  name: string
  fileName: string
  createTime: string
  status: 'running' | 'stopped'
}

const jarStatus = ref<JarStatus[]>([])

let timer: NodeJS.Timeout
onMounted(() => {
    timer = setInterval(() => {
        getJarStatus().then((res) => {
            const {code, data} = res
            if (code === 200) {
                jarStatus.value = data ? data : []
            }
        })
    }, 1000)
})
onBeforeUnmount(() => {
    clearInterval(timer)
})


// 处理删除操作
const handleDelete = (row: JarItem) => {
    if (row.status === 'running') {
        return ElMessage.error('请先停止应用再删除')
    }
    deleteJar(row.id).then((res) => {
        const {code} = res
        if (code === 200) {
            ElMessage.success('删除成功')
        }
    })
}

// 处理启动操作
const handleStart = (row: JarItem) => {
    startJar(row.id).then((res) => {
        const {code} = res
        if (code === 200) {
            return ElMessage.success('启动成功')
        }
        return ElMessage.error(res.message || '启动失败')
    })
}

// 处理停止操作
const handleStop = (row: JarItem) => {
    stopJar(row.id).then((res) => {
        const {code} = res
        if (code === 200) {
            ElMessage.success('停止成功')
        }
    })
}

// 上传对话框显示状态
const uploadDialogVisible = ref(false)
// 上传文件对象
const uploadFile = ref<File | null>(null)
// 应用名称
const appName = ref('')

const upload = ref<UploadInstance>()
// 处理文件选择
const handleFileChange = (file: any) => {
  if (file.raw) {
    uploadFile.value = file.raw
  }
}

// 处理上传确认
const handleUploadConfirm = async () => {
  if (!uploadFile.value || !appName.value) {
    ElMessage.warning('请填写完整信息')
    return
  }

  const formData = new FormData()
  formData.append('file', uploadFile.value)
  formData.append('name', appName.value)

  try {
    const res = await fetch('/api/v1/jar/upload', {
      method: 'POST',
      body: formData
    })

    if (!res.ok) {
        ElMessage.error('网络错误')
        return
    }
    const rsp = await res.json()
    if (rsp.code !== 200) {
        ElMessage.error(rsp.message || '上传失败')
        return
    }
    ElMessage.success('上传成功')
    uploadDialogVisible.value = false
    // 重置表单
    uploadFile.value = null
    appName.value = ''
    // 刷新列表
    getJarStatus()
  } catch (err) {
    console.error('Upload failed:', err)
    ElMessage.error('上传失败')
  }
}
</script>

<template>
  <div class="jar-manage">
    <div class="header">
      <h2>JAR包管理</h2>
      <el-button type="primary" @click="uploadDialogVisible = true">上传JAR包</el-button>
    </div>
    
    <el-table :data="jarStatus" border style="width: 100%; margin-top: 20px">
      <el-table-column prop="name" label="名称" />
      <el-table-column prop="fileName" label="文件名" />
      <el-table-column prop="createTime" label="添加日期" />
      <el-table-column prop="status" label="当前状态" />
      <el-table-column label="操作" width="300">
        <template #default="{ row }">
          <el-button 
            type="info" 
            size="small" 
            @click="handleStart(row)"
          >
            查看
          </el-button>
          <el-button 
            type="primary" 
            size="small" 
            @click="handleStart(row)"
            v-if="row.status === 'stopped'"
          >
            启动
          </el-button>
          <el-button 
            type="warning" 
            size="small" 
            @click="handleStop(row)"
            v-if="row.status === 'running'"
          >
            停止
          </el-button>
          <el-button 
            type="danger" 
            size="small" 
            @click="handleDelete(row)"
            v-if="row.status === 'stopped'"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 上传对话框 -->
    <el-dialog
      v-model="uploadDialogVisible"
      title="上传JAR包"
      width="500px"
    >
      <el-form label-width="80px">
        <el-form-item label="应用名称">
          <el-input v-model="appName" placeholder="请输入应用名称" />
        </el-form-item>
        <el-form-item label="选择文件">
            <el-upload
                :auto-upload="false"
                :show-file-list="true"
                accept=".jar"
                :limit="1"
                @change="handleFileChange"
                ref="upload"
            >
                <el-button type="primary">点击选择文件</el-button>
            </el-upload>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="uploadDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleUploadConfirm">确认</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.jar-manage {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.upload-area {
  border: 2px dashed #dcdfe6;
  border-radius: 6px;
  padding: 20px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
}

.upload-area:hover,
.upload-area.is-dragging {
  border-color: #409eff;
  background-color: #f5f7fa;
}

.upload-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.upload-icon {
  font-size: 48px;
  color: #909399;
}

.upload-text {
  color: #909399;
  margin-bottom: 8px;
}

.file-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px;
}

.file-name {
  color: #606266;
  font-size: 14px;
}
</style>
