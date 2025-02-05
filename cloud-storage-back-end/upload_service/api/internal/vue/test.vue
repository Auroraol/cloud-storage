<!-- FileUploader.vue -->
<template>
    <div class="file-uploader">
      <!-- 上传区域 -->
      <div 
        class="upload-area" 
        @drop.prevent="handleDrop"
        @dragover.prevent
        @dragenter.prevent="isDragging = true"
        @dragleave.prevent="isDragging = false"
        :class="{ 'dragging': isDragging }"
      >
        <div v-if="!currentFile" class="upload-placeholder">
          <i class="el-icon-upload"></i>
          <div class="text">将文件拖到此处，或<em @click="triggerFileInput">点击上传</em></div>
          <div class="hint">支持任意类型文件，大小不超过2GB</div>
        </div>
  
        <!-- 上传进度 -->
        <div v-else class="upload-progress">
          <div class="file-info">
            <i class="el-icon-document"></i>
            <span class="filename">{{ currentFile.name }}</span>
            <span class="filesize">({{ formatSize(currentFile.size) }})</span>
          </div>
          <el-progress 
            :percentage="uploadProgress" 
            :status="uploadStatus === 'error' ? 'exception' : undefined"
          />
          <div class="upload-status">{{ statusText }}</div>
          <div class="operations">
            <el-button 
              v-if="uploadStatus === 'paused'" 
              type="primary" 
              size="small" 
              @click="resumeUpload"
            >继续上传</el-button>
            <el-button 
              v-if="uploadStatus === 'uploading'" 
              type="warning" 
              size="small" 
              @click="pauseUpload"
            >暂停</el-button>
            <el-button 
              v-if="uploadStatus !== 'success'" 
              type="danger" 
              size="small" 
              @click="cancelUpload"
            >取消</el-button>
          </div>
        </div>
      </div>
  
      <!-- 上传历史 -->
      <div class="upload-history" v-if="uploadHistory.length > 0">
        <div class="history-title">上传历史</div>
        <div class="history-list">
          <div v-for="item in uploadHistory" :key="item.key" class="history-item">
            <i class="el-icon-document"></i>
            <div class="file-info">
              <div class="filename">{{ item.filename }}</div>
              <div class="filesize">{{ formatSize(item.size) }}</div>
            </div>
            <div class="file-status">
              <el-tag :type="item.status === 'success' ? 'success' : 'danger'">
                {{ item.status === 'success' ? '上传成功' : '上传失败' }}
              </el-tag>
            </div>
            <div class="operations">
              <el-button 
                type="text" 
                v-if="item.status === 'success'"
                @click="copyUrl(item.url)"
              >复制链接</el-button>
            </div>
          </div>
        </div>
      </div>
  
      <!-- 隐藏的文件输入框 -->
      <input 
        type="file" 
        ref="fileInput" 
        style="display: none" 
        @change="handleFileSelect"
      />
    </div>
  </template>
  
  <script>
  import axios from 'axios'
  import { ElMessage } from 'element-plus'
  
  export default {
    name: 'FileUploader',
    data() {
      return {
        isDragging: false,
        currentFile: null,
        uploadProgress: 0,
        uploadStatus: '', // uploading, paused, error, success
        statusText: '',
        uploadHistory: [],
        // 分片上传相关
        chunkSize: 5 * 1024 * 1024, // 5MB
        chunks: [],
        currentChunkIndex: 0,
        uploadId: '',
        uploadedETags: [],
        // 断点续传相关
        pauseUpload: false,
        uploadController: null
      }
    },
  
    methods: {
      // 触发文件选择
      triggerFileInput() {
        this.$refs.fileInput.click()
      },
  
      // 处理文件选择
      async handleFileSelect(event) {
        const file = event.target.files[0]
        if (file) {
          await this.startUpload(file)
        }
      },
  
      // 处理拖拽
      async handleDrop(event) {
        const file = event.dataTransfer.files[0]
        if (file) {
          this.isDragging = false
          await this.startUpload(file)
        }
      },
  
      // 开始上传
      async startUpload(file) {
        this.currentFile = file
        this.uploadProgress = 0
        this.uploadStatus = 'uploading'
        this.statusText = '准备上传...'
  
        try {
          if (file.size > 20 * 1024 * 1024) {
            await this.startChunkUpload(file)
          } else {
            await this.startNormalUpload(file)
          }
        } catch (error) {
          this.handleUploadError(error)
        }
      },
  
      // 普通上传
      async startNormalUpload(file) {
        const formData = new FormData()
        formData.append('file', file)
        formData.append('metadata', JSON.stringify({
          filename: file.name,
          type: file.type
        }))
  
        try {
          const response = await axios.post('/upload_service/v1/file/upload', formData, {
            headers: {
              'Content-Type': 'multipart/form-data'
            },
            onUploadProgress: this.handleProgress
          })
  
          this.handleUploadSuccess(response.data)
        } catch (error) {
          this.handleUploadError(error)
        }
      },
  
      // 分片上传
      async startChunkUpload(file) {
        try {
          // 1. 初始化分片上传
          const initResponse = await axios.post('/upload_service/v1/file/multipart/init', {
            fileName: file.name,
            fileSize: file.size,
            metadata: JSON.stringify({
              type: file.type
            })
          })
  
          this.uploadId = initResponse.data.uploadId
          const key = initResponse.data.key
  
          // 2. 准备分片
          const chunkCount = Math.ceil(file.size / this.chunkSize)
          this.chunks = Array.from({ length: chunkCount }, (_, index) => {
            const start = index * this.chunkSize
            const end = Math.min(file.size, start + this.chunkSize)
            return file.slice(start, end)
          })
  
          // 3. 上传分片
          this.uploadedETags = []
          for (let i = 0; i < this.chunks.length && !this.pauseUpload; i++) {
            this.currentChunkIndex = i
            const formData = new FormData()
            formData.append('file', this.chunks[i])
            formData.append('uploadId', this.uploadId)
            formData.append('key', key)
            formData.append('chunkIndex', i + 1)
  
            const response = await axios.post('/upload_service/v1/file/multipart/upload', formData)
            this.uploadedETags.push(response.data.etag)
            this.updateChunkProgress(i)
          }
  
          if (this.pauseUpload) {
            this.uploadStatus = 'paused'
            this.statusText = '上传已暂停'
            return
          }
  
          // 4. 完成上传
          const completeResponse = await axios.post('/upload_service/v1/file/multipart/complete', {
            uploadId: this.uploadId,
            key: key,
            etags: this.uploadedETags
          })
  
          this.handleUploadSuccess(completeResponse.data)
        } catch (error) {
          this.handleUploadError(error)
        }
      },
  
      // 更新分片上传进度
      updateChunkProgress(chunkIndex) {
        const progress = ((chunkIndex + 1) / this.chunks.length) * 100
        this.uploadProgress = Math.round(progress)
        this.statusText = `正在上传第 ${chunkIndex + 1}/${this.chunks.length} 个分片`
      },
  
      // 处理上传进度
      handleProgress(progressEvent) {
        if (this.currentFile.size <= 20 * 1024 * 1024) {
          const progress = (progressEvent.loaded / progressEvent.total) * 100
          this.uploadProgress = Math.round(progress)
          this.statusText = '正在上传...'
        }
      },
  
      // 处理上传成功
      handleUploadSuccess(data) {
        this.uploadStatus = 'success'
        this.uploadProgress = 100
        this.statusText = '上传成功'
        
        this.uploadHistory.unshift({
          filename: this.currentFile.name,
          size: this.currentFile.size,
          status: 'success',
          url: data.url,
          key: data.key
        })
  
        setTimeout(() => {
          this.resetUpload()
        }, 2000)
      },
  
      // 处理上传错误
      handleUploadError(error) {
        this.uploadStatus = 'error'
        this.statusText = `上传失败: ${error.message}`
        
        this.uploadHistory.unshift({
          filename: this.currentFile.name,
          size: this.currentFile.size,
          status: 'error'
        })
      },
  
      // 重置上传状态
      resetUpload() {
        this.currentFile = null
        this.uploadProgress = 0
        this.uploadStatus = ''
        this.statusText = ''
        this.uploadId = ''
        this.uploadedETags = []
        this.chunks = []
        this.currentChunkIndex = 0
        this.pauseUpload = false
      },
  
      // 暂停上传
      pauseUpload() {
        this.pauseUpload = true
        this.uploadStatus = 'paused'
        this.statusText = '正在暂停...'
      },
  
      // 继续上传
      async resumeUpload() {
        this.pauseUpload = false
        this.uploadStatus = 'uploading'
        await this.startChunkUpload(this.currentFile)
      },
  
      // 取消上传
      cancelUpload() {
        this.resetUpload()
      },
  
      // 格式化文件大小
      formatSize(bytes) {
        if (bytes === 0) return '0 B'
        const k = 1024
        const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
        const i = Math.floor(Math.log(bytes) / Math.log(k))
        return `${(bytes / Math.pow(k, i)).toFixed(2)} ${sizes[i]}`
      },
  
      // 复制URL
      copyUrl(url) {
        navigator.clipboard.writeText(url).then(() => {
          ElMessage.success('链接已复制到剪贴板')
        }).catch(() => {
          ElMessage.error('复制失败')
        })
      }
    }
  }
  </script>
  
  <style scoped>
  .file-uploader {
    width: 100%;
    max-width: 800px;
    margin: 0 auto;
  }
  
  .upload-area {
    border: 2px dashed #dcdfe6;
    border-radius: 6px;
    padding: 20px;
    text-align: center;
    cursor: pointer;
    transition: all 0.3s;
  }
  
  .upload-area.dragging {
    border-color: #409eff;
    background-color: rgba(64, 158, 255, 0.1);
  }
  
  .upload-placeholder {
    color: #909399;
  }
  
  .upload-placeholder i {
    font-size: 48px;
    margin-bottom: 10px;
  }
  
  .upload-placeholder .text {
    font-size: 16px;
    margin-bottom: 5px;
  }
  
  .upload-placeholder .text em {
    color: #409eff;
    font-style: normal;
    cursor: pointer;
  }
  
  .upload-placeholder .hint {
    font-size: 12px;
  }
  
  .upload-progress {
    padding: 20px;
  }
  
  .file-info {
    display: flex;
    align-items: center;
    margin-bottom: 15px;
  }
  
  .file-info i {
    font-size: 24px;
    margin-right: 10px;
  }
  
  .file-info .filename {
    font-weight: bold;
    margin-right: 5px;
  }
  
  .file-info .filesize {
    color: #909399;
  }
  
  .upload-status {
    margin-top: 10px;
    color: #909399;
  }
  
  .operations {
    margin-top: 15px;
  }
  
  .upload-history {
    margin-top: 30px;
  }
  
  .history-title {
    font-size: 16px;
    font-weight: bold;
    margin-bottom: 15px;
  }
  
  .history-item {
    display: flex;
    align-items: center;
    padding: 10px;
    border-bottom: 1px solid #ebeef5;
  }
  
  .history-item i {
    font-size: 20px;
    margin-right: 10px;
  }
  
  .history-item .file-info {
    flex: 1;
    margin-bottom: 0;
  }
  
  .history-item .file-status {
    margin: 0 15px;
  }
  
  .history-item .operations {
    margin-top: 0;
  }
  </style>