<template>
  <div class="transfer-records">
    <h1>传输记录</h1>
    <div class="records">
      <div class="record-section">
        <h2>上传记录</h2>
        <ul>
          <li v-for="record in uploadRecords" :key="record.id" class="record-item">
            <span class="file-name">{{ record.fileName }}</span>
            <span class="file-size">{{ record.size }}</span>
            <span class="file-date">{{ record.date }}</span>
          </li>
        </ul>
      </div>
      <div class="record-section">
        <h2>下载记录</h2>
        <ul>
          <li v-for="record in downloadRecords" :key="record.id" class="record-item">
            <span class="file-name">{{ record.fileName }}</span>
            <span class="file-size">{{ record.size }}</span>
            <span class="file-date">{{ record.date }}</span>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { ref } from "vue"

export default {
  name: "TransferRecords",
  setup() {
    const uploadRecords = ref([
      { id: 1, fileName: "file1.txt", size: "1 MB", date: "2023-01-01" },
      { id: 2, fileName: "file2.txt", size: "2 MB", date: "2023-01-02" }
    ])

    const downloadRecords = ref([
      { id: 1, fileName: "file3.txt", size: "3 MB", date: "2023-01-03" },
      { id: 2, fileName: "file4.txt", size: "4 MB", date: "2023-01-04" }
    ])

    return { uploadRecords, downloadRecords }
  }
}
</script>

<style scoped>
.transfer-records {
  padding: 20px;
  background-color: #f9f9f9;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.records {
  margin-top: 20px;
}

.record-section {
  margin-bottom: 30px;
}

h2 {
  margin-top: 20px;
  color: #333;
}

.record-item {
  display: flex;
  justify-content: space-between;
  padding: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  margin: 5px 0;
  background-color: #fff;
  transition: background-color 0.3s;
}

.record-item:hover {
  background-color: #f1f1f1;
}

.file-name {
  flex: 2;
}

.file-size {
  flex: 1;
  text-align: center;
}

.file-date {
  flex: 1;
  text-align: right;
  color: #666;
}
</style>
