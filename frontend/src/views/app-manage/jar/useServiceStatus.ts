import { onBeforeUnmount } from 'vue';

import { onMounted } from 'vue';

import { ref } from 'vue';
import { ServiceItem } from './types';
import { getJarStatus } from '@/api/modules/jar';

function useServiceStatus() {
  const serviceStatus = ref<ServiceItem[]>([]);

  let timer: NodeJS.Timeout;
  onMounted(() => {
    timer = setInterval(() => {
      getJarStatus().then((res) => {
        const { code, data } = res;
        if (code === 200) {
          serviceStatus.value = data ? data : [];
        }
      });
    }, import.meta.env.PROD ? 1000 : 3000);
  });

  onBeforeUnmount(() => {
    clearInterval(timer);
  });

  return {
    serviceStatus,
  };
}

export default useServiceStatus;
