package logic

//
//type FileUploadLogic struct {
//	logx.Logger
//	ctx    context.Context
//	svcCtx *svc.ServiceContext
//}
//
//// 文件上传
//func NewFileUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadLogic {
//	return &FileUploadLogic{
//		Logger: logx.WithContext(ctx),
//		ctx:    ctx,
//		svcCtx: svcCtx,
//	}
//}
//
//// 保存上传的文件到临时目录
//func saveUploadedFile(file *multipart.FileHeader) (string, error) {
//	src, err := file.Open()
//	if err != nil {
//		return "", err
//	}
//	defer src.Close()
//
//	// 创建临时文件
//	tempFile, err := os.CreateTemp("", "upload-*"+filepath.Ext(file.Filename))
//	if err != nil {
//		return "", err
//	}
//	defer tempFile.Close()
//
//	// 复制文件内容
//	_, err = io.Copy(tempFile, src)
//	if err != nil {
//		os.Remove(tempFile.Name()) // 清理临时文件
//		return "", err
//	}
//
//	return tempFile.Name(), nil
//}
//
//func (l *FileUploadLogic) FileUpload(req *types.FileUploadRequest) (resp *types.FileUploadResponse, err error) {
//	// 获取上传的文件
//	file, header, err := l.svcCtx.Request.FormFile("file")
//	if err != nil {
//		return nil, fmt.Errorf("获取上传文件失败: %v", err)
//	}
//	defer file.Close()
//
//	// 检查文件大小
//	if header.Size > 20*1024*1024 { // 20MB
//		return nil, fmt.Errorf("文件大小超过限制，请使用分片上传")
//	}
//
//	// 保存文件到临时目录
//	tempFile, err := saveUploadedFile(header)
//	if err != nil {
//		return nil, fmt.Errorf("保存临时文件失败: %v", err)
//	}
//	defer os.Remove(tempFile) // 确保清理临时文件
//
//	// 解析元数据
//	var metadata map[string]string
//	if req.Metadata != "" {
//		if err := json.Unmarshal([]byte(req.Metadata), &metadata); err != nil {
//			return nil, fmt.Errorf("解析元数据失败: %v", err)
//		}
//	}
//
//	// 构建OSS对象键（使用用户ID作为前缀）
//	userId := l.ctx.Value("userId").(string)
//	objectKey := fmt.Sprintf("users/%s/files/%s", userId, filepath.Base(header.Filename))
//
//	// 准备上传选项
//	options := []oss.Option{
//		oss.ContentType(header.Header.Get("Content-Type")),
//	}
//	for k, v := range metadata {
//		options = append(options, oss.Meta(k, v))
//	}
//
//	// 获取OSS bucket
//	bucket := oss.Bucket()
//	if bucket == nil {
//		return nil, fmt.Errorf("获取OSS Bucket失败")
//	}
//
//	// 上传文件到OSS
//	err = bucket.PutObjectFromFile(objectKey, tempFile, options...)
//	if err != nil {
//		return nil, fmt.Errorf("上传文件到OSS失败: %v", err)
//	}
//
//	// 生成文件URL
//	url := fmt.Sprintf("https://%s.%s/%s", oss.Config.BucketName, oss.Config.Endpoint, objectKey)
//
//	return &types.FileUploadResponse{
//		URL:  url,
//		Key:  objectKey,
//		Size: header.Size,
//	}, nil
//}

//
//// 文件上传
//func CosUpload(fileHeader *multipart.FileHeader, newId int64, b []byte) (string, string, error) {
//	u, _ := url.Parse(CosUrl)
//	bs := &cos.BaseURL{BucketURL: u}
//	c := cos.NewClient(bs, &http.Client{
//		Transport: &cos.AuthorizationTransport{
//			SecretID:  SecretID,
//			SecretKey: SecretKey,
//		},
//	})
//	baseName := path.Base(fileHeader.Filename)
//	name := "butane-netdisk/" + strconv.FormatInt(newId, 10) + baseName
//	_, err := c.Object.Put(context.Background(), name, bytes.NewReader(b), nil)
//	if err != nil {
//		return "", "", err
//	}
//	filePath := CosUrl + "/" + name
//	return filePath, baseName, nil
//}
