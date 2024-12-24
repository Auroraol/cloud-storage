package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

const (
	appKey         = "app"
	stackKey       = "stack"
	sourceKey      = "source"
	shopIDKey      = "shop_id"
	orderIDKey     = "order_id"
	goodsIDKey     = "goods_id"
	platGoodsIDKey = "plat_goods_id"
	platUserKey    = "plat_user"
	companyKey     = "company"
	hostKey        = "host"
	timeKey        = "time"
	msgKey         = "msg"
)

var standardKeys = map[string]int8{
	logrus.FieldKeyTime:  1, //time
	logrus.FieldKeyLevel: 1, //level
	appKey:               1, //app
	logrus.FieldKeyMsg:   1, //msg
	FileKey:              1, // file
	stackKey:             1, //stack
	sourceKey:            1, //source
	shopIDKey:            1, //shop_id
	orderIDKey:           1, //order_id
	goodsIDKey:           1, //goods_id
	platGoodsIDKey:       1, //plat_goods_id
	platUserKey:          1, //plat_user_key
	companyKey:           1, //company
	logrus.ErrorKey:      1, //error
	TraceKey:             1, //trace
	FuncKey:              1, //func
	hostKey:              1, //host
}

type MergeHook struct{}

func NewMergeHook() *MergeHook {
	return &MergeHook{}
}

func (h *MergeHook) Fire(entry *logrus.Entry) error {
	customValue := map[string]interface{}{}

	if err, ok := entry.Data[logrus.ErrorKey]; ok && err == nil {
		delete(entry.Data, logrus.ErrorKey)
	}

	for key, v := range entry.Data {
		if _, ok := standardKeys[key]; ok || key == customKey {
			continue
		}

		customValue[key] = v
		delete(entry.Data, key)
	}

	if len(AppName) > 0 {
		entry.Data[appKey] = AppName
	}

	if len(customValue) > 0 {
		entry.Data[customKey] = customValue
	}

	entry.Data[hostKey] = os.Getenv(xdRealHost)

	return nil
}

func (h *MergeHook) Levels() []Level {
	return AllLevels
}
