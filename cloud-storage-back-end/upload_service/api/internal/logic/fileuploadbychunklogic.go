package logic

//type FileUploadByChunkLogic struct {
//	logx.Logger
//	ctx    context.Context
//	svcCtx *svc.ServiceContext
//}

//// 文件分片上传
//func NewFileUploadByChunkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadByChunkLogic {
//	return &FileUploadByChunkLogic{
//		Logger: logx.WithContext(ctx),
//		ctx:    ctx,
//		svcCtx: svcCtx,
//	}
//}
//
//func (l *FileUploadByChunkLogic) FileUploadByChunk(req *types.FileUploadByChunkRequest, file multipart.File, fileHeader *multipart.FileHeader) (resp *types.FileUploadByChunkResponse, err error) {
//	// 判断是否已达用户容量上限
//	userId := token.GetUidFromCtx(l.ctx)
//	if userId == 0 {
//		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
//	}
//	volumeInfo, err := l.svcCtx.UserCenterRpc.FindVolumeById(l.ctx, &user.FindVolumeReq{Id: userId})
//	if err != nil {
//		return nil, err
//	}
//	if volumeInfo.NowVolume+fileHeader.Size > volumeInfo.TotalVolume {
//		return nil, response.NewErrCode(response.FILE_TOO_LARGE_ERROR)
//	}
//	// 增加用户当前已存储容量
//	_, err = l.svcCtx.UserCenterRpc.AddVolume(l.ctx, &user.AddVolumeReq{
//		Id:   userId,
//		Size: fileHeader.Size,
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	//判断文件是否已存在，若已存在则为秒传成功
//	b := make([]byte, fileHeader.Size)
//	_, err = file.Read(b)
//	if err != nil {
//		return nil, err
//	}
//	md5Str := utils.Md5ByBytes(b)
//	//
//	count, err := l.svcCtx.RepositoryPoolModel.CountByHash(l.ctx, md5Str)
//	if count > 0 {
//		repositoryInfo, err := l.svcCtx.RepositoryPoolModel.FindRepositoryPoolByHash(l.ctx, md5Str)
//		if err != nil {
//			return nil, err
//		}
//		return &types.FileUploadByChunkResponse{Id: repositoryInfo.Identity}, err
//	}
//
//	newId := uuid.NewV4().String()
//	//// 上传文件到cos，并得到filepath
//	//// 先将文件分块
//	//err = GenerateChunk(file, fileHeader, md5Str)
//	//if err != nil {
//	//	return nil, err
//	//}
//	//// 将分块后的文件进行上传
//	//parseInt, _ := strconv.ParseInt(newId, 10, 64)
//	//filePath, baseName, err := CosUploadByPart(fileHeader, md5Str, parseInt)
//	//if err != nil {
//	//	return nil, err
//	//}
//
//	var baseName string
//	var filePath string
//
//	// 插入数据
//	_, err = l.svcCtx.RepositoryPoolModel.InsertWithId(l.ctx, &model.RepositoryPool{
//		Identity: newId,
//		Hash:     md5Str,
//		Ext:      path.Ext(baseName),
//		Size:     fileHeader.Size,
//		Path:     filePath,
//		Name:     fileHeader.Filename,
//	})
//
//	if err != nil {
//		return nil, response.NewErrCode(response.FILE_UPLOAD_ERROR)
//	}
//
//	return &types.FileUploadByChunkResponse{Id: newId}, err
//}

//// 文件分块
//func GenerateChunk(file multipart.File, fileHeader *multipart.FileHeader, md5Str string) error {
//	ChunkSize := ChunkSize
//	chunkNum := math.Ceil(float64(fileHeader.Size) / float64(ChunkSize))
//	for i := 0; i < int(chunkNum); i++ {
//		//新建块，初始化大小
//		nowBlo := make([]byte, ChunkSize)
//		file.Seek(int64(i*ChunkSize), 0)
//		if int64(ChunkSize) > fileHeader.Size-int64(i*ChunkSize) {
//			nowBlo = make([]byte, int64(ChunkSize)-(fileHeader.Size-int64(i*ChunkSize)))
//		}
//		//读入块数据，向nowBlow中读入file的数据
//		file.Read(nowBlo)
//		f, err := os.OpenFile("service/repository/filePath/"+md5Str+strconv.FormatInt(int64(i), 10)+".chunk", os.O_CREATE, 0666)
//		if err != nil {
//			return err
//		}
//		//输出块
//		f.Write(nowBlo)
//		f.Close()
//	}
//	file.Close()
//	return nil
//}

//// 文件分块上传
//func CosUploadByPart(fileHeader *multipart.FileHeader, md5Str string, newId int64) (string, string, error) {
//	//获得上传的Upload，表示现在将上传的文件
//	ChunkSize := ChunkSize
//	chunkNum := math.Ceil(float64(fileHeader.Size) / float64(ChunkSize))
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
//	v, _, err := c.Object.InitiateMultipartUpload(context.Background(), name, nil)
//	if err != nil {
//		return "", "", err
//	}
//	UploadID := v.UploadID
//	opt := &cos.CompleteMultipartUploadOptions{}
//	for i := 0; i < int(chunkNum); i++ {
//		//获得该块的md5码值PartETag
//		f, err := os.ReadFile("service/repository/filePath/" + md5Str + strconv.FormatInt(int64(i), 10) + ".chunk")
//		if err != nil {
//			return "", "", err
//		}
//		resp, err := c.Object.UploadPart(
//			context.Background(), name, UploadID, i+1, bytes.NewReader(f), nil,
//		)
//		if err != nil {
//			return "", "", err
//		}
//		PartETag := resp.Header.Get("ETag")
//		//将该块塞入块数组opt
//		opt.Parts = append(opt.Parts, cos.Object{PartNumber: i + 1, ETag: PartETag})
//	}
//	//将所有块上传
//	_, _, err = c.Object.CompleteMultipartUpload(
//		context.Background(), name, UploadID, opt,
//	)
//	if err != nil {
//		return "", "", err
//	}
//	filePath := CosUrl + "/" + name
//	return filePath, baseName, nil
//}
