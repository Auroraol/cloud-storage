<template>
  <div class="sidebar-item">
    <el-menu :default-active="resolvePath(props.item.path)" mode="vertical" :unique-opened="true" class="main-menu">
      <template v-for="(item, index) in menuItems" :key="index">
        <!-- 没有子菜单的菜单项 -->
        <el-menu-item
          v-if="!item.children || item.children.length === 0"
          :index="resolvePath(item.path)"
          @click="handleMenuClick(item)"
        >
          <SvgIcon v-if="item.meta?.svgIcon" :name="item.meta.svgIcon" />
          <component v-else-if="item.meta?.elIcon" :is="item.meta.elIcon" class="el-icon" />
          <span v-if="item.meta?.title">{{ item.meta.title }}</span>
        </el-menu-item>

        <!-- 有子菜单的菜单项 -->
        <div v-else class="parent-menu-item" @click="toggleSubmenu(item)" :class="{ 'is-active': activeMenu === item }">
          <div class="menu-title">
            <SvgIcon v-if="item.meta?.svgIcon" :name="item.meta.svgIcon" />
            <component v-else-if="item.meta?.elIcon" :is="item.meta.elIcon" class="el-icon" />
            <span v-if="item.meta?.title">{{ item.meta.title }}</span>
            <span class="arrow" :class="{ 'is-active': activeMenu === item }">›</span>
          </div>
        </div>
      </template>
    </el-menu>

    <!-- 子菜单容器 -->
    <transition name="submenu-fade">
      <div v-if="activeMenu" class="submenu-container">
        <div class="submenu">
          <el-menu :default-active="resolvePath(activeMenu.path)" mode="vertical">
            <el-menu-item
              v-for="(child, childIndex) in activeMenu.children"
              :key="childIndex"
              :index="resolvePath(child.path)"
              @click="handleSubMenuClick(child)"
            >
              <SvgIcon v-if="child.meta?.svgIcon" :name="child.meta.svgIcon" />
              <component v-else-if="child.meta?.elIcon" :is="child.meta.elIcon" class="el-icon" />
              <span v-if="child.meta?.title">{{ child.meta.title }}</span>
            </el-menu-item>
          </el-menu>
        </div>
      </div>
    </transition>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref } from "vue"
import { useRouter } from "vue-router"
import { type RouteRecordRaw } from "vue-router"
import SvgIcon from "@/components/SvgIcon/index.vue"
import { isExternal } from "@/utils/validate"
import path from "path-browserify"

interface Props {
  item: RouteRecordRaw
  basePath?: string
}

const props = withDefaults(defineProps<Props>(), {
  basePath: ""
})

const router = useRouter()
const activeMenu = ref<RouteRecordRaw | null>(null)

const toggleSubmenu = (item: RouteRecordRaw) => {
  if (activeMenu.value === item) {
    activeMenu.value = null
  } else {
    activeMenu.value = item
  }
}

const handleMenuClick = (item: RouteRecordRaw) => {
  try {
    if (item.path) {
      router.push(resolvePath(item.path))
    }
  } catch (error) {
    console.error("Navigation error:", error)
  }
}

const handleSubMenuClick = (item: RouteRecordRaw) => {
  try {
    if (item.path) {
      const fullPath = resolvePath(item.path)
      router.push(fullPath)
      // 点击子菜单项后关闭子菜单
      activeMenu.value = null
    }
  } catch (error) {
    console.error("Navigation error:", error)
  }
}

/** 解析路径 */
const resolvePath = (routePath: string) => {
  if (!routePath) return ""
  switch (true) {
    case isExternal(routePath):
      return routePath
    case isExternal(props.basePath):
      return props.basePath
    default:
      return path.resolve(props.basePath, routePath)
  }
}

/** 菜单项列表 */
const menuItems = computed(() => {
  return props.item.children ? props.item.children.filter((child) => !child.meta?.hidden) : []
})
</script>

<style scoped>
.sidebar-item {
  position: relative;
  width: 200px;
}

.main-menu {
  background-color: var(--v3-sidebar-menu-bg-color);
  border-right: none;
}

.parent-menu-item {
  position: relative;
  cursor: pointer;
}

.menu-title {
  height: 56px;
  line-height: 56px;
  padding: 0 20px;
  color: var(--el-menu-text-color);
  display: flex;
  align-items: center;
  background-color: var(--v3-sidebar-menu-bg-color);
  transition: background-color 0.3s;
}

.menu-title:hover {
  background-color: var(--el-menu-hover-bg-color);
}

.parent-menu-item.is-active .menu-title {
  background-color: var(--el-menu-hover-bg-color);
}

.arrow {
  margin-left: auto;
  font-size: 20px;
  transition: transform 0.3s;
}

.arrow.is-active {
  transform: rotate(90deg);
}

.submenu-container {
  position: fixed;
  left: 200px;
  top: 0;
  height: 100vh;
  z-index: 100;
}

.submenu {
  background-color: var(--v3-sidebar-menu-bg-color);
  min-width: 200px;
  box-shadow: var(--el-box-shadow-light);
  height: 100%;
  border-left: 1px solid var(--el-border-color-light);
}

.submenu-fade-enter-active,
.submenu-fade-leave-active {
  transition:
    opacity 0.3s,
    transform 0.3s;
}

.submenu-fade-enter-from,
.submenu-fade-leave-to {
  opacity: 0;
  transform: translateX(-10px);
}

:deep(.el-menu-item) {
  color: var(--el-menu-text-color) !important;
  height: 56px;
  line-height: 56px;
}

:deep(.el-menu-item:hover) {
  background-color: var(--el-menu-hover-bg-color) !important;
}

:deep(.el-menu-item.is-active) {
  background-color: var(--el-menu-hover-bg-color) !important;
  color: var(--el-menu-active-color) !important;
  border-right: 2px solid var(--el-menu-active-color);
}

.submenu :deep(.el-menu) {
  border-right: none;
  background-color: transparent;
}

.submenu :deep(.el-menu-item) {
  padding: 0 20px !important;
  margin: 0;
  min-width: 200px;
}

.svg-icon {
  margin-right: 12px;
  font-size: 18px;
  color: inherit;
}

.el-icon {
  margin-right: 12px !important;
  font-size: 18px;
  color: inherit;
}

/* 移除 el-menu 的默认边框 */
:deep(.el-menu) {
  border-right: none !important;
}

/* 确保菜单项内容垂直居中 */
:deep(.el-menu-item),
.menu-title {
  display: flex;
  align-items: center;
}
</style>
