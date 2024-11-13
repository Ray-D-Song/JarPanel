<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue';
import { ElMessage } from 'element-plus';
import { stopJar, startJar, deleteJar } from '@/api/modules/jar';
import { ServiceItem } from './types';
import NewServiceDialog from './NewServiceDialog.vue';
import DetailDialog from './DetailDialog.vue';
import useServiceStatus from './useServiceStatus';

const { serviceStatus } = useServiceStatus();

// 处理删除操作
const handleDelete = (row: ServiceItem) => {
  if (row.status === 'running') {
    return ElMessage.error('请先停止应用再删除');
  }
  deleteJar(row.id).then((res) => {
    const { code } = res;
    if (code === 200) {
      ElMessage.success('删除成功');
    }
  });
};

// 处理启动操作
const handleStart = (row: ServiceItem) => {
  startJar(row.id).then((res) => {
    const { code } = res;
    if (code === 200) {
      return ElMessage.success('启动成功');
    }
    return ElMessage.error(res.message || '启动失败');
  });
};

// 处理停止操作
const handleStop = (row: ServiceItem) => {
  stopJar(row.id).then((res) => {
    const { code } = res;
    if (code === 200) {
      ElMessage.success('停止成功');
    }
  });
};

const newServiceDialogVisible = ref(false);
const detailDialogRef = ref<any>(null);

function handleOpenDetail(service: ServiceItem) {
  detailDialogRef.value.handleOpen(service);
}
</script>

<template>
  <div class="jar-manage">
    <div class="header">
      <h2>JAR包管理</h2>
      <el-button type="primary" @click="newServiceDialogVisible = true">新建服务</el-button>
    </div>
    <DetailDialog ref="detailDialogRef" />

    <el-table :data="serviceStatus" border style="width: 100%; margin-top: 20px">
      <el-table-column prop="name" label="服务名称" />
      <el-table-column prop="currentVersion" label="Jar包" />
      <el-table-column prop="createTime" label="添加日期" />
      <el-table-column prop="status" label="当前状态" />
      <el-table-column label="操作" width="300">
        <template #default="{ row }">
          <el-button type="info" size="small" @click="handleOpenDetail(row)">查看</el-button>
          <el-button type="primary" size="small" @click="handleStart(row)" v-if="row.status === 'stopped'">
            启动
          </el-button>
          <el-button type="warning" size="small" @click="handleStop(row)" v-if="row.status === 'running'">
            停止
          </el-button>
          <el-button type="danger" size="small" @click="handleDelete(row)" v-if="row.status === 'stopped'">
            删除
          </el-button>
          <el-button type="success" size="small">实时日志</el-button>
        </template>
      </el-table-column>
    </el-table>

    <NewServiceDialog v-model="newServiceDialogVisible" />
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
</style>
