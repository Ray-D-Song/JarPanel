<script setup lang="ts">
import { getServiceFileList } from '@/api/modules/jar';
import { ref } from 'vue';

const fileList = ref<{ name: string }[]>([]);
const currentVersion = ref('');
const visible = ref(false);
const serviceId = ref('');
function handleOpen(id: string, version: string) {
  serviceId.value = id;
  getServiceFileList(id).then((res) => {
    if (res.code === 200) {
      fileList.value = res.data.map((file) => ({
        name: file,
      }));
    }
  });
  currentVersion.value = version;
  visible.value = true;
}

function handleDownload(name: string) {
  fetch(`/api/v1/jar/download?id=${serviceId.value}&name=${name}`).then(async (res) => {
    if (!res.ok) {
      return;
    }

    const blob = await res.blob();
    // 设置文件名
    const fileName = name;
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = fileName;
    a.click();
    // 释放 URL 对象
    URL.revokeObjectURL(url);
    // 删除 a 元素
    document.body.removeChild(a);
  });
}

const file = ref<File | null>(null);
function handleUpload() {
  if (!file.value) {
    ElMessage.error('请选择文件');
    return;
  }
  const formData = new FormData();
  formData.append('file', file.value);
  fetch(`/api/v1/jar/upload?id=${serviceId.value}&name=${currentVersion.value}`, {
    method: 'POST',
    body: formData,
  });
}

defineExpose({
  handleOpen,
});
</script>

<template>
  <el-dialog v-model="visible" width="80%">
    <template #title>
      <section class="flex items-center gap-4">
        <h4 class="text-lg font-bold">文件列表</h4>
        <el-upload :before-upload="handleUpload">
          <el-button type="primary" size="small">上传文件</el-button>
        </el-upload>
      </section>
    </template>
    <div class="p-6 pt-0">


    <el-table :data="fileList" border>
      <el-table-column prop="name" label="文件名">
              <template #default="{ row }">
                <div class="flex items-center gap-4">
                  <span>{{ row.name }}</span>
                  <el-tag type="primary" size="small" v-if="row.name === currentVersion">运行版本</el-tag>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="操作">
              <template #default="{ row }">
                <el-button type="primary" size="small" @click="handleDownload(row.name)">下载</el-button>
                <el-button type="danger" size="small">删除</el-button>
                <el-button v-if="row.name !== currentVersion" type="warning" size="small">设为运行版本</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
  </el-dialog>
</template>
