<template>
  <div class="app-container">
    <div class="detail-main">
      <div class="detail-title">
        <span>分享详情</span>
      </div>
      <!-- 多文件 -->
      <div class="muti-file" v-if="!selectFile && props.selectFileProps.length > 1">
        <!-- <Icon icon-name="folder_big" :width="128" /> -->
      </div>
      <!-- 单文件 -->
      <div class="detail-content" v-if="selectFile">
        <div class="filename">
          <template v-if="selectFile.ext">
            <Icon :cover="selectFile.path" />
          </template>
          <template v-else>
            <Icon :file-type="0" :width="32" />
          </template>
          <span class="text">{{ selectFile.name }}</span>
        </div>
        <template v-for="(item, index) in details" :key="index">
          <div class="content-item" v-if="item.label !== 'divider'">
            <div class="label">{{ item.label }}</div>
            <div class="value">
              <template v-if="item.key === 'update_time'">
                {{ formatTime(selectFile[item.key]) }}
              </template>
              <template v-else-if="item.key === 'expired_time'">
                {{ formatExpiredTime(selectFile[item.key]) }}
              </template>
              <template v-else>
                {{ selectFile[item.key] !== undefined ? selectFile[item.key] : 0 }}
              </template>
            </div>
          </div>
          <div class="divider" v-else />
        </template>
      </div>
      <!-- 空文件 -->
      <div class="detail-empty" v-if="props.selectFileProps.length === 0">
        <img src="https://nd-static.bdstatic.com/m-static/v20-main/home/img/empty-folder.55c81ea2.png" alt="空文件夹" />
        <p>选中文件，查看详情</p>
      </div>
    </div>
  </div>
</template>
<script>
export default {
  name: "detail"
}
</script>
<script setup>
import { computed } from "vue"
import dayjs from "dayjs"
import Icon from "@/components/FileIcon/Icon.vue"
import { formatTime } from "@/utils/format/formatTime"

const props = defineProps({
  selectFileProps: {
    type: Array,
    default: () => []
  }
})

const selectFile = computed(() => {
  // console.log("selectFileProps", props.selectFileProps[0])
  return props.selectFileProps.length === 1 ? props.selectFileProps[0] : null
})

const details = [
  { label: "分享时间", key: "update_time" },
  { label: "有效期", key: "expired_time" },
  { label: "提取码", key: "code" },
  { label: "浏览", key: "click_num" },
  { label: "保存", key: "saveCount" },
  { label: "下载", key: "downloadCount" }
]

// 计算过期时间
const formatExpiredTime = (expiredTime) => {
  if (!expiredTime) return "-"
  const expirationDate = dayjs(selectFile.value.update_time).add(expiredTime, "second") // 将当前时间加上过期时间（秒）
  return expirationDate.format("YYYY-MM-DD HH:mm:ss")
}
</script>

<style lang="scss" scoped>
.detail-main {
  // padding: 24px 16px;
  // background-color: #f9f9f9;
  // border-radius: 8px;
  // box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}
.detail-title {
  color: #03081a;
  font-weight: 600;
  font-size: 18px;
  margin-bottom: 15px;
}
.detail-empty {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 100%;
  text-align: center;
  transform: translate(-50%, -50%);
  img {
    width: 120px;
    height: 120px;
    margin-bottom: 8px;
  }
  p {
    color: #818999;
    font-weight: 400;
    font-size: 14px;
    line-height: 20px;
  }
}
.muti-file {
  background-color: rgba(214, 220, 224, 0.15);
  width: 100%;
  min-height: 134px;
  border-radius: 13px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;
}
.detail-content {
  .filename {
    width: 100%;
    font-size: 16px;
    font-weight: 500;
    margin-bottom: 12px;
    display: flex;
    align-items: center; // 图标和文字对齐
  }

  .text {
    margin-left: 8px;
  }
  .content-item {
    padding: 10px 0;
    line-height: 24px;
    display: flex;
    align-items: center; // 图标和文字对齐
    justify-content: space-between;
    .label {
      color: #878c9c;
      font-weight: 500;
      flex: 1;
    }
    .value {
      color: #03081a;
      font-weight: 600;
    }
  }
  .divider {
    margin: 10px 0;
    border-bottom: 1px solid #d4d7de;
  }
}
</style>
