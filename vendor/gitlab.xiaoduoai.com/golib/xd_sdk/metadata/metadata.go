package metadata

import (
	"context"
	"net/http"
	"strings"
	"sync"
)

type metadataCtxKey struct{}

const metadataStrKey = "xd-metadata"

const (
	flowColor  = "flow_color"  // 流量标志的key
	mockColor  = "mock_color"  // Mock标志的Key
	printColor = "print_color" // 打印标志的key
)

const (
	mdParisSeparator    = "||" // k-v pairs 之间的间隔符
	mdKeyValueSeparator = "="  // k-v pair, k和v的间隔符

	flowColorTest = "pt" //全链路压测
	//flowColorProd = "prod"

	printColorAtLeastInfo = "print_at_least_info"

	//mockColorDefaultResp = "default_resp" //mock，默认状态，客户端没有提供返回数据,服务端返回默认值
	mockColorRandResp = "rand_resp" //mock, 某些字段随机数值
	mockColorHasResp  = "has_resp"  //mock，客户端提供返回数据
)

type FlowColor int

type Metadata interface {
	Get(key string) (string, bool)
	Set(key, val string)
	GetD() D
	isTestFlow() bool
	isPrintAtLeastInfo() bool
	isHasResp() bool
	isRandResp() bool
}

type metadata struct {
	D  map[string]string
	mu *sync.RWMutex
}

type D map[string]string

func NewMetadata(d D) Metadata {
	return &metadata{
		D:  d,
		mu: &sync.RWMutex{},
	}
}

func (m *metadata) Get(key string) (string, bool) {
	if m == nil || m.D == nil {
		return "", false
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.D[key]
	return s, ok
}

func (m *metadata) Set(key, val string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.D == nil {
		m.D = make(map[string]string)
	}
	m.D[key] = val
}

func (m *metadata) GetD() D {
	m.mu.RLock()
	defer m.mu.RUnlock()
	data := make(D, len(m.D))
	for k, v := range m.D {
		data[k] = v
	}
	return data
}

func (m *metadata) isTestFlow() bool {
	b, ok := m.Get(flowColor)
	if !ok {
		return false
	}

	return b == flowColorTest
}

func (m *metadata) isPrintAtLeastInfo() bool {
	b, ok := m.Get(printColor)
	if !ok {
		return false
	}

	return b == printColorAtLeastInfo
}

func (m *metadata) isHasResp() bool {
	b, ok := m.Get(mockColor)
	if !ok {
		return false
	}

	return b == mockColorHasResp
}

func (m *metadata) isRandResp() bool {
	b, ok := m.Get(mockColor)
	if !ok {
		return false
	}

	return b == mockColorRandResp
}

func NewMetadataWithFlowColorTest() Metadata {
	return NewMetadata(D{flowColor: flowColorTest})
}

// 添加日志等级info提升标记
func MarkPrintColorAtLeastInfo(ctx context.Context) context.Context {
	md := FromContext(ctx)
	md.Set(printColor, printColorAtLeastInfo)
	return WithMetadata(ctx, md)
}

// 添加Tag,
func AddTag(ctx context.Context, key, val string) context.Context {
	md := FromContext(ctx)
	md.Set(key, val)
	return WithMetadata(ctx, md)
}

// 添加Tag,
func AddTags(ctx context.Context, tags map[string]string) context.Context {
	md := FromContext(ctx)
	for k, v := range tags {
		md.Set(k, v)
	}
	return WithMetadata(ctx, md)
}

func IsTestFlow(ctx context.Context) bool {
	m := FromContext(ctx)
	return m.isTestFlow()
}

func IsPrintAtLeastInfo(ctx context.Context) bool {
	m := FromContext(ctx)
	return m.isPrintAtLeastInfo()
}

func IsMockHasResp(ctx context.Context) bool {
	m := FromContext(ctx)
	return m.isHasResp()
}

func IsMockRandResp(ctx context.Context) bool {
	m := FromContext(ctx)
	return m.isRandResp()
}

func FromContext(ctx context.Context) Metadata {
	//if ctx == nil {
	//	return make(Metadata)
	//}

	m, ok := ctx.Value(metadataCtxKey{}).(Metadata)
	if !ok {
		m = NewMetadata(nil)
		return m
	}

	return m
}

func WithMetadata(ctx context.Context, metadata Metadata) context.Context {
	return context.WithValue(ctx, metadataCtxKey{}, metadata)
}

func Encode(ctx context.Context) (key string, value string) {
	md := FromContext(ctx)

	//data, err := json.Marshal(md)
	//if err != nil {
	//	return metadataStrKey, ""
	//}
	//return metadataStrKey, string(data)
	return metadataStrKey, md2String(md)
}

// Header value format: k1=v1||k2=v2||k3=v3||...
func FromHTTPRequest(ctx context.Context, req *http.Request) context.Context {
	header := req.Header.Get(metadataStrKey)
	if header == "" {
		if c, err := req.Cookie(metadataStrKey); err == nil {
			header = c.Value
		}
	}

	if header == "" {
		return WithMetadata(ctx, NewMetadata(nil))
	}

	//ret := make(Metadata)
	//_ = json.Unmarshal([]byte(header), &ret)
	ret := string2Md(header)
	return WithMetadata(ctx, ret)
}

func FromStringMap(ctx context.Context, m map[string]string) context.Context {
	s := m[metadataStrKey]
	if s == "" {
		return WithMetadata(ctx, NewMetadata(nil))
	}

	//ret := make(Metadata)
	//_ = json.Unmarshal([]byte(s), &ret)
	ret := string2Md(s)
	return WithMetadata(ctx, ret)
}

// FromString: Decode 功能，与Encode对应
func FromString(ctx context.Context, ss string) context.Context {
	md := NewMetadata(nil)
	if ss == "" {
		return WithMetadata(ctx, md)
	}
	//_ = json.Unmarshal([]byte(ss), &ret)
	md = string2Md(ss)
	return WithMetadata(ctx, md)
}

// string s format: k1=v1||k2=v2||k3=v3...
func string2Md(s string) Metadata {
	items := strings.Split(s, mdParisSeparator)
	n := len(items)
	d := make(D, n)
	for _, item := range items {
		ss := strings.Split(item, mdKeyValueSeparator)
		if len(ss) != 2 {
			continue
		}
		d[ss[0]] = ss[1]
	}

	return NewMetadata(d)
}

// return string format: k1=v1||k2=v2||k3=v3...
func md2String(md Metadata) string {
	ret := ""
	for k, v := range md.GetD() {
		ret += k + mdKeyValueSeparator + v + mdParisSeparator
	}
	retLength := len(ret)
	if retLength > len(mdParisSeparator) {
		return ret[:len(ret)-len(mdParisSeparator)]
	}
	return ret
}

// 通用的生成影子key 函数, 返回shadow_XXX，确保前缀最多只有一个"shadow_"
func GenShadowString(ss string) string {
	if !strings.HasPrefix(ss, "shadow_") {
		ss = "shadow_" + ss
	}
	return ss
}
