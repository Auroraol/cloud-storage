import { computed, ref, unref, watch } from "vue"

/**
 * 表格加载动画钩子函数
 * @param {Prop} props
 */
export const useLoading = (props) => {
  const loadingRef = ref(unref(props).loading)

  // 监听loading的变化
  watch(
    () => unref(props).loading,
    (loading) => {
      loadingRef.value = loading
      console.log(loadingRef.value)
    }
  )

  const getLoading = computed(() => unref(loadingRef))

  const setLoading = (loading) => {
    loadingRef.value = loading
  }
  console.log(loadingRef.value)
  return { getLoading, setLoading }
}
