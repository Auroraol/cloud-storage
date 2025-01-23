import { inject, provide } from "vue"

const key = Symbol("s-table")

export const createTableContext = (instance) => provide(key, instance)

export const useTableContext = () => inject(key)
