<template>
  <div>
    <div class="action-text" v-if="mode === 'row'">
      <p :class="['text', value.showEdit ? 'show' : 'hover-hidden']">{{ value[field] }}</p>
      <div :class="['actions', value.showEdit ? '' : 'hover-show']" :style="{ left: `${offset}px` }">
        <div class="flex">
          <div class="action-group">
            <template v-for="(action, index) in fileTypeActions.slice(0, 4)" :key="index">
              <div
                class="action-item"
                @click.stop="action.onClick(value)"
                v-show="action.isShow ? action.isShow(value) : true"
              >
                <i :class="['iconfont', 'icon-' + action.icon]" />
                <template v-if="action.showLabel">{{ action.title }}</template>
              </div>
            </template>
          </div>
          <div class="action-more" v-if="fileTypeActions.length > 4">
            <el-popover
              placement="right"
              trigger="hover"
              :offset="2"
              :show-arrow="false"
              :hide-after="0"
              width="auto"
              :teleported="false"
              popper-style="padding: 4px 10px; min-width: auto"
            >
              <template #reference>
                <div class="action-item">
                  <i class="iconfont icon-more" />
                </div>
              </template>
              <template #default>
                <div class="pop-action-group">
                  <template v-for="(action, index) in fileTypeActions.slice(4)" :key="index">
                    <div class="pop-action-item" @click.stop="action.onClick(value)">
                      <el-tooltip :content="action.title" placement="top">
                        <SvgIcon :name="action.icon" />
                      </el-tooltip>
                    </div>
                  </template>
                </div>
              </template>
            </el-popover>
          </div>
        </div>
      </div>
    </div>
    <div class="action-btn" v-else-if="mode === 'btns'">
      <el-button-group>
        <template v-for="(action, index) in fileTypeActions" :key="index">
          <el-button
            text
            round
            type="primary"
            bg
            v-if="action.isShow ? action.isShow(value) : true"
            @click="action.onClick(value)"
          >
            <SvgIcon :name="action.icon" />
            {{ action.title }}
          </el-button>
        </template>
      </el-button-group>
    </div>
  </div>
</template>
<script setup lang="ts">
import { onMounted, ref } from "vue"
import SvgIcon from "@/components/SvgIcon/index.vue"
const props = defineProps({
  value: {
    type: [Array, Object],
    default: () => []
  },
  mode: {
    type: String,
    default: "btns"
  },
  field: {
    type: String,
    default: ""
  },
  actions: {
    type: Array,
    default: () => []
  },
  offset: {
    type: [String, Number],
    default: 0
  }
})

// Start of Selection
const fileTypeActions = computed(() => props.actions)

// 处理操作
function filterActions() {
  // 目前没有下载选项
  if (props.value.folderType === 1) {
    fileTypeActions.value = props.actions.filter((item) => item.label !== "download")
  }
}

onMounted(() => {
  filterActions()
})
</script>
<style lang="scss">
.action-btn {
  .iconfont {
    font-size: 16px;
    margin-right: 6px;
  }
}

.action-text {
  position: relative;
  height: 50px;
  line-height: 50px;

  .actions {
    display: none;
    z-index: 100;
    position: absolute;
    bottom: 0;
    .flex {
      display: flex;
      align-items: center;
    }

    .action-group {
      display: flex;

      .action-item {
        display: flex;
        justify-content: center;
        background-color: transparent;
        padding: 0;
        margin-right: 12px;
        color: #06a7ff;

        &:last-child {
          margin-right: 0;
        }

        .iconfont {
          margin-right: 2px;
          font-size: 14px;
        }
      }
    }
  }
}

.action-more {
  margin-left: 12px;
  color: #06a7ff;

  .pop-action-group {
    display: flex;
    flex-direction: row;

    .pop-action-item {
      display: inline-flex;
      align-items: center;
      font-size: 12px;
      color: #03081a;
      cursor: pointer;
      padding: 4px;
      margin-left: 12px;

      &:first-child {
        margin-left: 0;
      }

      .iconfont {
        font-size: 14px;
      }
    }
  }
}
</style>
