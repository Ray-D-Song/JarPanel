import http from '@/api'

/**
 * 获取 JAR 包状态
 */
export const getJarStatus = () => {
  return http.get<any>('/jar/status')
}

/**
 * 删除 JAR 包
 * @param id JAR 包 ID
 */
export const deleteJar = (id: string) => {
  return http.delete<void>(`/jar/delete?id=${id}`)
}

/**
 * 启动 JAR 包
 * @param id JAR 包 ID
 */
export const startJar = (id: string) => {
  return http.put<void>(`/jar/start?id=${id}`)
}

/**
 * 停止 JAR 包
 * @param id JAR 包 ID
 */
export const stopJar = (id: string) => {
  return http.put<void>(`/jar/stop?id=${id}`) 
}
