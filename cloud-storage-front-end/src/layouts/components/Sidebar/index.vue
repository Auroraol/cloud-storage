<script lang="ts" setup>
import { computed } from "vue"
import { useRoute } from "vue-router"
import { useAppStore } from "@/store/modules/app"
import { usePermissionStore } from "@/store/modules/permission"
import { useSettingsStore } from "@/store/modules/settings"
import SidebarItem from "./SidebarItem.vue"
import Logo from "../Logo/index.vue"
import { useDevice } from "@/hooks/useDevice"
import { useLayoutMode } from "@/hooks/useLayoutMode"
import { getCssVar } from "@/utils/css"

const v3SidebarMenuBgColor = getCssVar("--v3-sidebar-menu-bg-color")
const v3SidebarMenuTextColor = getCssVar("--v3-sidebar-menu-text-color")
const v3SidebarMenuActiveTextColor = getCssVar("--v3-sidebar-menu-active-text-color")

const { isMobile } = useDevice()
const { isLeft, isTop } = useLayoutMode()
const route = useRoute()
const appStore = useAppStore()
const permissionStore = usePermissionStore()
const settingsStore = useSettingsStore()

const activeMenu = computed(() => {
  const {
    meta: { activeMenu },
    path
  } = route
  return activeMenu ? activeMenu : path
})
const noHiddenRoutes = computed(() => permissionStore.routes.filter((item) => !item.meta?.hidden))
const isCollapse = computed(() => !appStore.sidebar.opened)
const isLogo = computed(() => isLeft.value && settingsStore.showLogo)
const backgroundColor = computed(() => (isLeft.value ? v3SidebarMenuBgColor : undefined))
const textColor = computed(() => (isLeft.value ? v3SidebarMenuTextColor : undefined))
const activeTextColor = computed(() => (isLeft.value ? v3SidebarMenuActiveTextColor : undefined))
const sidebarMenuItemHeight = computed(() => {
  return !isTop.value ? "var(--v3-sidebar-menu-item-height)" : "var(--v3-navigationbar-height)"
})
const sidebarMenuHoverBgColor = computed(() => {
  return !isTop.value ? "var(--v3-sidebar-menu-hover-bg-color)" : "transparent"
})
const tipLineWidth = computed(() => {
  return !isTop.value ? "2px" : "0px"
})

import { useUserStore } from "@/store/modules/user"
import { formatFileSize } from "@/utils/format/formatFileSize"

const userStore = useUserStore()
const capacity = computed(() => userStore.capacity)
const capacityPercentage = computed(() => {
  const percentage = (capacity.value.now_volume / capacity.value.total_volume) * 100
  return percentage.toFixed(2)
})
</script>

<template>
  <div :class="{ 'has-logo': isLogo }">
    <Logo v-if="isLogo" :collapse="isCollapse" />
    <el-scrollbar wrap-class="scrollbar-wrapper">
      <el-menu
        :default-active="activeMenu"
        :collapse="isCollapse && !isTop"
        :background-color="backgroundColor"
        :text-color="textColor"
        :active-text-color="activeTextColor"
        :unique-opened="true"
        :collapse-transition="false"
        :mode="isTop && !isMobile ? 'horizontal' : 'vertical'"
      >
        <SidebarItem v-for="route in noHiddenRoutes" :key="route.path" :item="route" :base-path="route.path" />
      </el-menu>
      <!-- 存储占用 -->
      <div class="storage-occupancy">
        <el-progress type="circle" class="storage-occupancy-progress" :percentage="capacityPercentage" />
        <div class="storage-occupancy-tip">
          <el-tooltip effect="dark" content="存储占用：上传文件后会消耗存储" placement="top">
            <div style="color: #696969">存储占用情况</div>
          </el-tooltip>
          <div style="color: #696969; margin-top: 5px">
            <span style="white-space: nowrap">
              {{ formatFileSize(capacity.now_volume) }}/{{ formatFileSize(capacity.total_volume) }}
            </span>
            {{ capacityPercentage }}%
          </div>
        </div>
      </div>
    </el-scrollbar>
  </div>
</template>

<style lang="scss" scoped>
%tip-line {
  &::before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    width: v-bind(tipLineWidth);
    height: 100%;
    background-color: var(--v3-sidebar-menu-tip-line-bg-color);
  }
}

.has-logo {
  .el-scrollbar {
    height: calc(100% - var(--v3-header-height));
  }
}

.el-scrollbar {
  height: 100%;
  :deep(.scrollbar-wrapper) {
    // 限制水平宽度
    overflow-x: hidden !important;
  }
  // 滚动条
  :deep(.el-scrollbar__bar) {
    &.is-horizontal {
      // 隐藏水平滚动条
      display: none;
    }
  }
}

.el-menu {
  border: none;
  width: 100% !important;
}

.el-menu--horizontal {
  height: v-bind(sidebarMenuItemHeight);
}

:deep(.el-menu-item),
:deep(.el-sub-menu__title),
:deep(.el-sub-menu .el-menu-item),
:deep(.el-menu--horizontal .el-menu-item) {
  height: v-bind(sidebarMenuItemHeight);
  line-height: v-bind(sidebarMenuItemHeight);
  &.is-active,
  &:hover {
    background-color: v-bind(sidebarMenuHoverBgColor);
  }
}

:deep(.el-sub-menu) {
  &.is-active {
    > .el-sub-menu__title {
      color: v-bind(activeTextColor) !important;
    }
  }
}

:deep(.el-menu-item.is-active) {
  @extend %tip-line;
}

.el-menu--collapse {
  :deep(.el-sub-menu.is-active) {
    .el-sub-menu__title {
      @extend %tip-line;
    }
  }
}

.storage-occupancy {
  // 内容居中
  .storage-occupancy-progress {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, 50%);
    margin-top: 100px; // 修改为0以便于调整位置
    color: #fcbc52;
  }
  .storage-occupancy-tip {
    position: absolute;
    top: 61%;
    left: 50%;
    transform: translate(-50%, 50%);
    text-align: center; // 使文字居中
    margin-top: 190px; // 添加上边距以便于与进度条分开
  }
}
</style>
