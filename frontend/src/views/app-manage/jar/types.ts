export interface ServiceItem {
  id: string;
  // 应用名称
  name: string;
  // 前缀参数
  prefixArgs: string;
  // 后缀参数
  suffixArgs: string;
  // 创建时间
  createTime: string;
  // 部署时间
  deployTime: string;
  // 状态
  status: 'running' | 'stopped';

  // 当前版本
  currentVersion: string;
  // 上一个版本
  previousVersion: string;
}
