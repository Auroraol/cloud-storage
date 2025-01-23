<template>
  <div class="share-container">
    <div class="share-header">
      <div class="top" v-if="selectFiles.length === 0">
        <span class="title">链接分享</span>
        <span class="content">(分享失败超过1年以上的链接记录将被自动清理)</span>
      </div>
      <div class="actions" v-else>
        <!-- <FileAction :value="selectFiles" mode="btns" :actions="actions" /> -->
      </div>
    </div>
    <div class="share-body">
      <div class="left-container">
        <div class="nav">
          <span class="txt">全部文件</span>
          <span class="count">已加载{{ tableDataNum }}个</span>
        </div>
        <div class="file-list">
          <vue-good-table
            ref="shareTableRef"
            :columns="columns"
            :rows="tableRows"
            :select-options="{
              enabled: true,
              selectOnCheckboxOnly: true,
              selectionText: '选中的行数',
              clearSelectionText: ''
            }"
            :search-options="{
              enabled: true,
              placeholder: '搜索分享',
              skipDiacritics: true
            }"
            :pagination-options="{
              enabled: true,
              mode: 'pages',
              perPage: 10,
              perPageDropdown: [10, 20, 50],
              dropdownAllowAll: false,
              nextLabel: '下一页',
              prevLabel: '上一页',
              rowsPerPageLabel: '每页显示',
              pageLabel: '页',
              ofLabel: '/'
            }"
            v-on:row-click="handleTableRowSelect"
            styleClass="vgt-table"
          >
            <template #table-row="props">
              <span v-if="props.column.field === 'filename'">
                <div class="file-item">
                  <template v-if="(props.row.fileType === 3 || props.row.fileType === 1) && props.row.status === 2">
                    <Icon :cover="props.row.fileCover" :width="32" />
                  </template>
                  <template v-else>
                    <Icon v-if="props.row.folderType === 0" :file-type="props.row.fileType" />
                    <Icon v-else icon-name="folder_2" :width="40" />
                  </template>
                  <span class="filename">{{ props.row.filename }}</span>
                </div>
              </span>
            </template>
            <template #selected-row-actions>
              <!-- <FileAction :value="selectFiles" mode="btns" :actions="actions" /> -->
            </template>
            <template #emptystate>
              <el-empty description="空空如也" />
            </template>
          </vue-good-table>
        </div>
      </div>
      <div class="right-detail">
        <el-scrollbar>
          <ShareDetail :selectFileProps="selectFiles" />
        </el-scrollbar>
      </div>
      <style lang="scss" scoped />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, computed } from "vue"
import { VueGoodTable } from "vue-good-table-next"
import { columns } from "./columns"
import Icon from "@/components/FileIcon/Icon.vue"
import ShareDetail from "./detail.vue"
import useClipboard from "vue-clipboard3"
import { ElMessage } from "element-plus"
import { isArray } from "@/utils/is"
// import { FileAction } from "@/components/FileAction"

console.log("当前数据:")

interface ShareItem {
  filename: string
  fileType: number
  status: number
  fileId: string // 文件id
  userId: string // 用户id
  validType: boolean // 有效期类型 0:1天 1:7天 2:30天 3:永久有效
  code: string // 提取码
  browseCount: number // 浏览次数
  saveCount: number // 收藏次数
  downloadCount: number // 下载次数
  expireTime: Date // 过期时间
  createTime: Date // 更新时间
}

const shareTableRef = ref()
const selectFiles = ref<ShareItem[]>([])
const shareIds = ref<string>("")
const tableDataNum = ref<number>(0)
const { toClipboard } = useClipboard()
const shareUrl = ref<string>(document.location.origin + "/share/")

const actions = [
  {
    title: "复制链接",
    label: "copyLink",
    showLabel: true,
    icon: "link",
    isShow: (data: ShareItem | ShareItem[]) => {
      return !isArray(data) || (isArray(data) && data.length === 1)
    },
    onClick: (data: ShareItem | ShareItem[]) => {
      copyLink2Clipboard(data)
    }
  },
  {
    title: "取消分享",
    label: "cancelShare",
    showLabel: true,
    icon: "stop",
    onClick: (data: ShareItem | ShareItem[]) => {
      cancelShare(data)
    }
  }
]

const dialogConfig = reactive({
  show: false,
  title: "确认取消分享",
  showCancel: true,
  buttons: [
    {
      type: "primary",
      text: "确认",
      click: () => {
        deleteShare()
      }
    }
  ]
})

const loadDataTable = (res: any) => {
  return loadShareList(res)
}

const fetchSuccess = ({ resultTotal }: { resultTotal: number }) => {
  tableDataNum.value = resultTotal
}

const reloadTable = () => {
  shareTableRef.value?.reload()
}

const handleTableRowSelect = (params) => {
  // 确保 selectFiles 是一个数组
  selectFiles.value = [params.row] // 将选中的行作为数组传递
}

const handleRowClick = (params) => {
  console.log("Row clicked:", params)
}

const cancelShare = (data: ShareItem | ShareItem[]) => {
  if (isArray(data)) {
    shareIds.value = data.map((item) => item.id).join(",")
  } else {
    shareIds.value = data.id.toString()
  }
  dialogConfig.show = true
}

const deleteShare = () => {
  ElMessage.success("取消外链分享成功")
  reloadTable()
  dialogConfig.show = false
}

const copyLink2Clipboard = async (data: ShareItem | ShareItem[]) => {
  const shareData = isArray(data) ? data[0] : data
  await toClipboard(`链接: ${shareUrl.value}${shareData.shareId} 提取码: ${shareData.code}`)
  ElMessage.success("复制成功")
}

// 测试文件列表数据
const fileList = ref<ShareItem[]>([
  {
    filename: "文件1",
    fileType: 0,
    status: 2,
    fileId: "1",
    userId: "2",
    validType: true,
    code: "ABC123",
    browseCount: 10,
    saveCount: 5,
    downloadCount: 2,
    expireTime: new Date("2024-12-31"),
    createTime: new Date("2023-01-01")
  },
  {
    filename: "文件2",
    fileType: 9,
    status: 1,
    fileId: "2",
    userId: "1",
    validType: false,
    code: "XYZ789",
    browseCount: 20,
    saveCount: 10,
    downloadCount: 5,
    expireTime: new Date("2023-06-30"),
    createTime: new Date("2023-02-01")
  }
])

// 格式化日期的函数
const formatDate = (date: string): string => {
  if (!date) return "-"
  try {
    const dateObj = new Date(date)
    const year = dateObj.getFullYear()
    const month = String(dateObj.getMonth() + 1).padStart(2, "0")
    const day = String(dateObj.getDate()).padStart(2, "0")
    const hours = String(dateObj.getHours()).padStart(2, "0")
    const minutes = String(dateObj.getMinutes()).padStart(2, "0")
    const seconds = String(dateObj.getSeconds()).padStart(2, "0")

    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
  } catch {
    return "-"
  }
}

// 格式化表格数据
const tableRows = computed(() => {
  return fileList.value.map((file) => ({
    ...file,
    createTime: formatDate(file.createTime.toISOString()),
    expireTime: formatDate(file.expireTime.toISOString())
  }))
})
</script>

<style lang="scss">
.share-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 10px 4px 0;
}

.share-header {
  height: 50px;
  padding-left: 4px;
  font-size: 12px;
  display: flex;
  justify-content: space-between;
  align-items: center;

  .top {
    display: inline-block;
    line-height: 36px;
  }
}

.share-body {
  display: flex;
  height: calc(100% - 50px); // 确保填充剩余空间

  .left-container {
    flex: 1; // 使左侧容器占据剩余空间
    display: flex;
    flex-direction: column;

    .nav {
      font-size: 12px;
      color: #03081a;
      padding-left: 4px;
      display: flex;
      justify-content: space-between;

      .txt {
        color: #03081a;
        font-weight: 600;
        font-size: 18px;
        margin-bottom: 15px;
      }
      .count {
        font-size: 18px;
        color: #afb3bf;
      }
    }

    .file-list {
      padding-left: 4px;
      flex: 1; // 使文件列表占据剩余空间
      overflow-y: auto; // 添加滚动条
      min-height: 272px;
    }
  }

  .right-detail {
    width: 272px; // 固定宽度
    overflow: hidden; // 确保内容不溢出
    display: inline-block;
    position: relative;
    height: 100%;
    margin-left: 16px;
    // background-color: #f5f6fa;
    /* background-color: #f5f6f8; */
    background-color: #f9f9f9;
    border-radius: 8px;
    padding: 16px;
    font-size: 13px;
    margin-right: 4px; // 添加右侧边距
  }
}

.vgt-table {
  min-height: 450px;
  td {
    vertical-align: middle !important;
  }

  .file-item {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .filename {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}
</style>
