package logic

import (
	"context"
	"github.com/Auroraol/cloud-storage/common/cache"
	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/upload_service/rpc/uploadservice"
	"github.com/Auroraol/cloud-storage/user_center/rpc/client/userrepository"
	jsoniter "github.com/json-iterator/go"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"strconv"

	"github.com/Auroraol/cloud-storage/share_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/share_service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShareBasicDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取资源详情
func NewShareBasicDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareBasicDetailLogic {
	return &ShareBasicDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareBasicDetailLogic) ShareBasicDetail(req *types.DetailRequest) (resp *types.DetailResponse, err error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	//先查询缓存中有没有该数据
	redisQueryKey := cache.CacheShareKey + strconv.FormatInt(req.Id, 10)
	ifExists, err := l.svcCtx.RedisClient.Exists(redisQueryKey)
	if err != nil {
		return nil, err
	}
	if ifExists == true {
		jsonStr, err := l.svcCtx.RedisClient.Get(redisQueryKey)
		if err != nil {
			return nil, err
		}
		//判断数据是否为空
		if jsonStr == "" {
			return nil, response.NewErrCodeMsg(response.SUCCESS, "查无此分享信息")
		}
		var shareInfo types.DetailResponse
		err = json.UnmarshalFromString(jsonStr, &shareInfo)
		if err != nil {
			return nil, err
		}
		//增加点击数
		err = l.svcCtx.ShareBasicModel.AddOneClick(l.ctx, req.Id)
		if err != nil {
			return nil, response.NewErrMsg("增加点击数失败！")
		}
		return &shareInfo, nil
	}
	//从数据库查询数据
	//申请分布式锁，获取repositoryId和返回user库、repositoryPool库的相应值
	redisLockKey := redisQueryKey
	redisLock := redis.NewRedisLock(l.svcCtx.RedisClient, redisLockKey)
	redisLock.SetExpire(cache.RedisLockExpireSeconds)
	if ok, err := redisLock.Acquire(); !ok || err != nil {
		return nil, response.NewErrCodeMsg(response.SUCCESS, "当前有其他用户正在进行操作，请稍后重试")
	}
	defer func() {
		recover()
		redisLock.Release()
	}()
	shareInfo, err := l.svcCtx.ShareBasicModel.FindOne(l.ctx, uint64(req.Id))
	switch err {
	case nil:
		break
	case sqlc.ErrNotFound:
		//缓存空数据
		err = l.svcCtx.RedisClient.Setex(redisQueryKey, "", cache.RedisLockExpireSeconds)
		if err != nil {
			return nil, err
		}
		return nil, response.NewErrCodeMsg(response.SUCCESS, "查无此分享信息")
	default:
		return nil, err
	}
	userRepositoryName, err := l.svcCtx.UserCenterRepositoryRpc.GetUserRepositoryNameByRepositoryId(l.ctx, &userrepository.RepositoryIdReq{
		RepositoryId: int64(shareInfo.RepositoryId),
	})
	if err != nil {
		return nil, response.NewErrMsg("无法获得用户储存库的信息！")
	}
	RepositoryPool, err := l.svcCtx.UploadServiceRpc.GetRepositoryPoolByRepositoryId(l.ctx, &uploadservice.RepositoryReq{
		RepositoryId: int64(shareInfo.RepositoryId),
	})
	if err != nil {
		return nil, response.NewErrMsg("无法获得储存池的信息！")
	}
	//把数据存储到缓存中
	DetailInfo := types.DetailResponse{
		RepositoryId: int64(shareInfo.RepositoryId),
		Name:         userRepositoryName.RepositoryName,
		Ext:          RepositoryPool.Ext,
		Size:         RepositoryPool.Size,
		Path:         RepositoryPool.Path,
	}
	jsonStr, err := json.MarshalToString(DetailInfo)
	if err != nil {
		return nil, err
	}
	l.svcCtx.RedisClient.Setex(redisQueryKey, jsonStr, cache.RedisLockExpireSeconds)
	//增加点击数
	err = l.svcCtx.ShareBasicModel.AddOneClick(l.ctx, req.Id)
	if err != nil {
		return nil, response.NewErrMsg("增加点击数失败！")
	}
	return &DetailInfo, nil
}
