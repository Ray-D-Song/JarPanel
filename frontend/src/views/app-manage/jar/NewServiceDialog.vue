<script setup lang="ts">
import { ElUpload, UploadFile, UploadProps, UploadRawFile } from 'element-plus';
import { watch, ref } from 'vue';

const modelValue = defineModel<boolean>({
  required: true,
});

const form = ref({
  name: '',
  file: null,
});
const upload = ref<InstanceType<typeof ElUpload>>();
watch(() => modelValue.value, () => {
  form.value = {
    name: '',
    file: null,
  };
  upload.value?.clearFiles();
});

const handleChange: UploadProps['onChange'] = (uploadFile) => {
  const file = uploadFile.raw as UploadRawFile;
  form.value.file = file;
};

const onCancel = () => {
  modelValue.value = false;
};
const onSubmit = () => {
  if (form.value.name === '') {
    ElMessage.error('服务名称不能为空');
    return;
  }
  if (form.value.file === null) {
    ElMessage.error('Jar包不能为空');
    return;
  }
  const formData = new FormData();
  formData.append('name', form.value.name);
  formData.append('file', form.value.file);
  fetch('/api/v1/jar/new', {
    method: 'POST',
    body: formData,
  }).then(async (res) => {
    if (!res.ok) {
      ElMessage.error('创建失败');
      return;
    }
    const { code } = await res.json()
    if (code === 200) {
      ElMessage.success('创建成功');
      setTimeout(() => {
        modelValue.value = false;
      }, 600);
    } else {
      ElMessage.error('创建失败');
    }
  });
};
</script>

<template>
  <el-dialog v-model="modelValue" title="新建服务" width="500px">
    <el-form label-width="80px">
      <el-form-item label="服务名称">
        <el-input v-model="form.name" placeholder="请输入服务名称" />
      </el-form-item>
      <el-form-item label="Jar包">
        <el-upload
          :auto-upload="false"
          :show-file-list="true"
          accept=".jar"
          :limit="1"
          :on-change="handleChange"
          ref="upload"
        >
          <template #trigger>
            <el-button type="primary">点击选择文件</el-button>
          </template>
        </el-upload>
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="onCancel">取消</el-button>
        <el-button type="primary" @click="onSubmit">确认</el-button>
      </span>
    </template>
  </el-dialog>
</template>
