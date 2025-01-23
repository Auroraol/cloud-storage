import { h } from "vue"
import EditableCell from "./EditableCell.vue"

interface Column {
  prop: string
  // 其他属性可以根据需要添加
}

interface Record {
  [key: string]: any
  onEdit?: (edit: boolean, submit?: boolean) => boolean
  editable?: boolean
  onSubmitEdit?: () => boolean
  onCancelEdit?: () => void
}

/**
 * 渲染为可编辑单元格，并挂载属性
 * @param {Column} column
 */
export const renderEditCell = (column: Column) => (record: Record, index: number) => {
  const _prop = column.prop
  const value = record[_prop]
  record.onEdit = (edit: boolean, submit: boolean = false): boolean => {
    if (!submit) {
      record.editable = edit
    }
    // 提交编辑
    if (!edit && submit) {
      const res = record.onSubmitEdit?.()
      if (res) {
        record.editable = false
        return true
      }
      return false
    }
    // 取消编辑
    if (!edit && !submit) {
      record.onCancelEdit?.()
    }
    return true
  }
  return h(EditableCell, {
    value,
    record,
    column,
    index
  })
}
