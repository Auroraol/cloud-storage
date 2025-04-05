export interface ListItem {
  avatar?: string
  title: string
  datetime?: string
  description?: string
  status?: "primary" | "success" | "info" | "warning" | "danger"
  extra?: string
}

export const notifyData: ListItem[] = [
  {
    avatar: "https://gw.alipayobjects.com/zos/rmsportal/OKJXDXrmkNshAMvwtvhu.png",
    title: "云存储空间升级通知",
    datetime: "2小时前",
    description: "您的云存储空间已成功升级，现在您可以享受更大的存储容量和更快的访问速度。"
  },
  {
    avatar: "https://gw.alipayobjects.com/zos/rmsportal/OKJXDXrmkNshAMvwtvhu.png",
    title: "数据备份完成",
    datetime: "1天前",
    description: "您的文件数据已成功备份，可以在备份历史中查看详情。"
  }
]

export const messageData: ListItem[] = [
  {
    avatar: "https://gw.alipayobjects.com/zos/rmsportal/ThXAXghbEsBCCSDihZxY.png",
    title: "文件分享消息",
    description: "用户 admin 与您分享了文件：项目计划书.docx",
    datetime: "10分钟前"
  },
  {
    avatar: "https://gw.alipayobjects.com/zos/rmsportal/ThXAXghbEsBCCSDihZxY.png",
    title: "协作邀请",
    description: "用户 guest 邀请您加入协作文件夹：团队项目",
    datetime: "昨天"
  },
  {
    avatar: "https://gw.alipayobjects.com/zos/rmsportal/ThXAXghbEsBCCSDihZxY.png",
    title: "评论提醒",
    description: "用户 user1 评论了您的文件：会议记录.pdf",
    datetime: "3天前"
  }
]

export const todoData: ListItem[] = [
  {
    title: "存储空间清理",
    description: "您的存储空间即将用完，建议清理不必要的文件",
    extra: "紧急",
    status: "danger"
  },
  {
    title: "文件整理",
    description: "根据新的分类方案对文档进行整理",
    extra: "进行中",
    status: "primary"
  },
  {
    title: "安全检查",
    description: "对共享文件的权限进行安全检查",
    extra: "计划中",
    status: "info"
  }
]
