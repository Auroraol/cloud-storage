<template>
  <div class="table">
    {{ tableData.length }}
    <div class="file-table" v-show="true || tableData.length > 0">
      <el-table
        ref="tableElRef"
        v-loading="getLoading"
        v-bind="getBindValues"
        @row-click="handleRowClick"
        @selection-change="handleSelectionChange"
        :row-class-name="tableRowClassName"
      >
        <el-table-column v-if="props.selection" type="selection" width="48px" align="center" class-name="checkbox" />
        <template v-for="(column, index) in props.columns" :key="index">
          <el-table-column
            v-bind="
              column.scopedSlots
                ? {
                    prop: column.prop,
                    label: column.label,
                    align: column.align || 'left',
                    width: column.width,
                    minWidth: column.minWidth,
                    showOverflowTooltip: column.ellipsis
                  }
                : {
                    prop: column.prop,
                    label: column.label,
                    align: column.align || 'left',
                    width: column.width,
                    minWidth: column.minWidth,
                    showOverflowTooltip: column.ellipsis
                  }
            "
          >
            <template v-if="column.scopedSlots" #default="scope">
              <slot :name="column.scopedSlots" :index="scope.$index" :row="scope.row" />
            </template>
          </el-table-column>
        </template>
      </el-table>
    </div>
    <div class="empty" v-show="!getLoading && tableData.length === 0">
      <slot name="empty" />
      {{ getLoading }}
      {{ tableData.length }}
      {{ getProps.dataSource }}
    </div>
  </div>
</template>

<script lang="ts" setup>
import { defineProps, defineEmits } from "vue"
import { computed, nextTick, onMounted, ref, unref, watch } from "vue"
import { basicProps } from "./props"
import { useLoading } from "./hooks/useLoading"
import { usePagination } from "./hooks/usePagination"
import { useDataSource } from "./hooks/useDataSource"
import { getViewportOffset } from "@/utils/domUtils"
import { useWindowSizeFn } from "@/utils/useWindowSizeFn"

// 发送子参数
const emit = defineEmits(["fetch-success", "fetch-error", "select-change"])

// 接受父参数
const props = defineProps({ ...basicProps })

const tableHeight = ref(150)
const tableData = ref([])
const tableElRef = ref()
const innerPropsRef = ref()

const getProps = computed(() => {
  // 解包 innerPropsRef 并与 props 合并
  const innerProps = unref(innerPropsRef) || {} // 确保 innerPropsRef 解包后是一个对象
  return { ...props, ...innerProps } // 合并两个对象
})

console.log(getProps.value)

const { getLoading, setLoading } = useLoading(getProps)
console.log(getLoading.value) // 这里打印getLoading的值

const { getPaginationInfo, setPagination } = usePagination(getProps)
console.log(getPaginationInfo.value)
const { getDataSourceRef, getRowKey } = useDataSource(
  getProps,
  {
    getPaginationInfo,
    setPagination,
    // tableDat,
    tableData,
    setLoading
  },
  emit
)

tableData.value = unref(getProps).dataSource
console.log("初始化tableData的值:", tableData.value)

// 计算table的绑定值
const getBindValues = computed(() => {
  const tableData = unref(getDataSourceRef)
  console.log(tableData)
  const maxHeight = tableData.length ? `${unref(tableHeight)}px` : "auto"
  return {
    ...unref(getProps),
    data: tableData,
    rowKey: unref(getRowKey),
    maxHeight: maxHeight
  }
})

let numbers = []
const handleRowClick = (row) => {
  const isSelect = tableElRef.value.getSelectionRows().includes(row) && tableElRef.value.getSelectionRows().length == 1
  tableElRef.value.clearSelection()
  tableElRef.value.toggleRowSelection(row, !isSelect)
}
const handleSelectionChange = (selection) => {
  const ids = selection.map((item) => item.id)
  numbers = ids
  emit("select-change", selection)
}
const tableRowClassName = ({ row, _ }) => {
  let color = ""
  numbers.forEach((id) => {
    if (id === row.id) {
      color = "is-selected"
    }
  })
  return color
}

const computeTableHeight = async () => {
  const table = unref(tableElRef)
  if (!table) return
  const tableEl = table?.$el
  const headerEl = tableEl.querySelector(".el-table__header")
  const { bottomIncludeBody } = getViewportOffset(headerEl)
  tableHeight.value = bottomIncludeBody
}

useWindowSizeFn(computeTableHeight, 280, 1000)

watch(
  () => tableData.value.length,
  () => {
    computeTableHeight()
  }
)

watch(
  () => getProps.value.dataSource,
  (newData) => {
    tableData.value = newData
    console.log("更新后的tableData:", tableData.value)
  }
)

onMounted(() => {
  nextTick(() => {
    computeTableHeight()
  })
})
</script>
