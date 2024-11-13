<script setup lang="ts">
import { ref } from 'vue';
import { ServiceItem } from './types';
import FileListDialog from './FileListDialog.vue';

const visible = ref(false);
const form = ref<ServiceItem | null>(null);
function handleOpen(service: ServiceItem) {
  form.value = service;
  visible.value = true;
}

const fileListDialogRef = ref<any>(null);
function handleOpenFileList() {
  fileListDialogRef.value.handleOpen(form.value?.id, form.value?.currentVersion);
}

defineExpose({
  handleOpen,
});
</script>

<template>
  <el-dialog v-model="visible" title="服务详情">
    <el-form :model="form" label-width="100px">
      <el-form-item label="应用名称：">
        <span>{{ form?.name }}</span>
      </el-form-item>
      <el-form-item label="前缀参数：">
        <span>{{ form?.prefixArgs }}</span>
      </el-form-item>
      <el-form-item label="后缀参数：">
        <span>{{ form?.suffixArgs }}</span>
      </el-form-item>
      <el-form-item label="创建时间：">
        <span>{{ form?.createTime }}</span>
      </el-form-item>
      <el-form-item label="部署时间：">
        <span>{{ form?.deployTime }}</span>
      </el-form-item>
      <el-form-item label="当前状态：">
        <el-tag :type="form?.status === 'running' ? 'success' : 'danger'">{{
          form?.status
        }}</el-tag>
      </el-form-item>
      <el-form-item label="运行版本：">
        <span>{{ form?.currentVersion }}</span>
      </el-form-item>
      <el-form-item label="文件夹：">
        <el-button type="primary" size="small" @click="handleOpenFileList">查看文件</el-button>
      </el-form-item>
    </el-form>
    <FileListDialog ref="fileListDialogRef" />
  </el-dialog>
</template>
