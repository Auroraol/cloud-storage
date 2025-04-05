<script lang="ts">
import { defineComponent } from "vue"
import { type ListItem } from "./data"

export default defineComponent({
  name: "NotifyList",
  props: {
    list: {
      type: Array as () => ListItem[],
      required: true
    }
  },
  setup(props) {
    return { props }
  }
})
</script>

<template>
  <el-empty v-if="props.list.length === 0" />
  <el-card v-else v-for="(item, index) in props.list" :key="index" shadow="never" class="card-container">
    <template #header>
      <div class="card-header">
        <div>
          <span>
            <span class="card-title">{{ item.title }}</span>
            <el-tag v-if="item.extra" :type="item.status" effect="plain" size="small">{{ item.extra }}</el-tag>
          </span>
          <div class="card-time">{{ item.datetime }}</div>
        </div>
        <div v-if="item.avatar" class="card-avatar">
          <img :src="item.avatar" width="34" />
        </div>
      </div>
    </template>
    <div class="card-body">
      {{ item.description ?? "No Data" }}
    </div>
  </el-card>
</template>

<style lang="scss" scoped>
.card-container {
  margin-bottom: 10px;
  transition:
    transform 0.3s,
    box-shadow 0.3s;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .card-title {
      font-weight: bold;
      margin-right: 10px;
      color: var(--el-text-color-primary);
    }

    .card-time {
      font-size: 12px;
      color: var(--el-text-color-secondary);
      margin-top: 4px;
    }

    .card-avatar {
      display: flex;
      align-items: center;

      img {
        border-radius: 50%;
        object-fit: cover;
      }
    }
  }

  .card-body {
    font-size: 13px;
    color: var(--el-text-color-regular);
    line-height: 1.6;
    word-break: break-word;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
  }
}
</style>
