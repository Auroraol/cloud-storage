<template>
  <div class="editable-cell">
    <div v-show="!isEdit" class="editable-cell-content">
      {{ currentValue }}
      <!-- 显示当前值 -->
      <div class="flex editable-cell-content" v-show="isEdit" v-click-outside="onClickOutside">
        <!-- 当处于编辑状态时显示输入框 -->
        <div class="editable-cell-content">
          <ElInput ref="elRef" v-model="currentValue" />
          <!-- 输入框，绑定当前值 -->
        </div>
        <div class="editable-cell-action">
          <n-icon class="mx-2 cursor-pointer">
            <CheckOutlined @click="submitEdit" />
            <!-- 提交编辑 -->
          </n-icon>
          <n-icon class="mx-2 cursor-pointer">
            <CloseOutlined @click="cancelEdit" />
            <!-- 取消编辑 -->
          </n-icon>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "EditableCell" // 组件名称
}
</script>

<script setup>
import { isBoolean } from "@/utils/is" // 导入工具函数
import { ElInput } from "element-plus" // 导入Element Plus的输入框组件
import { nextTick, ref, unref, watchEffect } from "vue" // 导入Vue的响应式API
import { useTableContext } from "../../hooks/useTableContext" // 导入表格上下文

const props = {
  value: {
    type: String,
    default: "" // 默认值为空字符串
  },
  record: {
    type: Object // 记录对象
  },
  column: {
    type: Object,
    default: () => ({}) // 默认列对象
  },
  index: {
    type: Number // 行索引
  }
}

const isEdit = ref(false) // 编辑状态
const elRef = ref() // 输入框引用
const table = useTableContext() // 获取表格上下文
const currentValue = ref(props.value) // 当前值
const defaultValue = ref(props.value) // 默认值

watchEffect(() => {
  const { editable } = props.column // 获取列的可编辑属性
  isEdit.value = isBoolean(editable) ? !!editable : false // 设置编辑状态
})

const handleEdit = () => {
  isEdit.value = true // 开始编辑
  nextTick(() => {
    const el = unref(elRef) // 获取输入框引用
    el?.focus?.() // 聚焦输入框
  })
}

const submitEdit = (needEmit = true, valid = true) => {
  if (valid && !true) return false // 校验

  const { column, index, record } = props // 获取列、索引和记录
  if (!record) return false // 如果没有记录则返回
  const { prop } = column // 获取列的属性
  if (!prop) return // 如果没有属性则返回

  needEmit && table.emit?.("edit-end", { record, index, prop, value: unref(currentValue) }) // 触发编辑结束事件
  isEdit.value = false // 结束编辑
}

const handleEnter = async () => {
  submitEdit() // 提交编辑
}

const cancelEdit = () => {
  isEdit.value = false // 取消编辑
  currentValue.value = defaultValue.value // 恢复默认值
  const { column, index, record } = props // 获取列、索引和记录
  const { prop } = column // 获取列的属性
  table.emit?.("edit-cancel", {
    // 触发取消编辑事件
    record,
    index,
    prop,
    value: unref(currentValue)
  })
}

if (props.record) {
  props.record.cancelCbs = cancelEdit // 记录取消回调
  props.record.submitCbs = submitEdit // 记录提交回调
  props.record.onCancelEdit = () => props?.cancelCbs() // 取消编辑的回调
  props.record.onSubmitEdit = () => {
    props.record?.submitCbs() // 提交编辑的回调
    return true
  }
}
</script>
<style lang="scss" scoped></style>
