<template>
  <div class="app-container">
    <div class="share-header">
      <div class="top" v-if="selectFiles.length === 0">
        <span class="title">链接分享</span>
        <span class="content">(分享失败超过1年以上的链接记录将被自动清理)</span>
      </div>
      <div class="actions" v-else>
        <FileAction :value="selectFiles" mode="btns" :actions="actions" />
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
    </div>
    <FileDialog
      :show="dialogConfig.show"
      :title="dialogConfig.title"
      :buttons="dialogConfig.buttons"
      :show-cancel="true"
      width="400px"
      @close="dialogConfig.show = false"
    >
      <div class="delete-info">
        <div class="info">
          <span class="text">取消分享后，该条分享记录将被删除，好友将无法再访问此分享链接。 您确认要取消分享吗？</span>
        </div>
      </div>
    </FileDialog>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, computed } from "vue"
import { VueGoodTable } from "vue-good-table-next"
import useClipboard from "vue-clipboard3"
import { columns } from "./columns"
import Icon from "@/components/FileIcon/Icon.vue"
import ShareDetail from "./detail.vue"
import FileDialog from "@/components/FileDialog/index.vue"
import { ElMessage } from "element-plus"
import { isArray } from "@/utils/is"
import { FileAction } from "@/components/FileAction/index"

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
    icon: "cancel-share",
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

// 计算属性来处理选中的文件
const selectedFileIds = computed(() => {
  return selectFiles.value.map((file) => file.fileId).join(",")
})

const fetchSuccess = ({ resultTotal }: { resultTotal: number }) => {
  tableDataNum.value = resultTotal
}

const reloadTable = () => {
  shareTableRef.value?.reload()
}

const handleTableRowSelect = (params) => {
  selectFiles.value = [params.row] // 将选中的行作为数组传递
}

const handleRowClick = (params) => {
  console.log("Row clicked:", params)
}

const cancelShare = (data: ShareItem | ShareItem[]) => {
  if (isArray(data)) {
    shareIds.value = selectedFileIds.value // 使用计算属性
  } else {
    shareIds.value = data.fileId.toString()
  }
  dialogConfig.show = true
}

const deleteShare = () => {
  ElMessage.success("取消外链分享成功")
  dialogConfig.show = false
  reloadTable()
}

const copyLink2Clipboard = async (data: ShareItem | ShareItem[]) => {
  const shareData = isArray(data) ? data[0] : data
  await toClipboard(`链接: ${shareUrl.value}${shareData.fileId} 提取码: ${shareData.code}`)
  ElMessage.success("复制成功")
}

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

const loadDataTable = (res: any) => {
  return loadShareList(res)
}

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

// 计算属性来处理表格行数据
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
      flex: 1; // 使文件列表占据剩余空间
      overflow-y: auto; // 添加滚动条
      min-height: 272px;
    }
  }

  .right-detail {
    width: 220px; // 固定宽度
    overflow: hidden; // 确保内容不溢出
    display: inline-block;
    position: relative;
    height: 100%;
    margin-left: 16px;
    background-color: #f9f9f9;
    border-radius: 8px;
    padding: 1px;
    font-size: 13px;
  }
}

.vgt-table {
  min-height: 450px;

  tr {
    height: 50px !important; // 强制设置行高
  }

  td {
    padding: 10px !important; // 强制设置单元格内边距
    vertical-align: middle !important;
    text-align: center; // 设置水平对齐方式为居中
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
