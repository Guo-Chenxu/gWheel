package limiter

type ILimiter interface {
	// 尝试是否限流，成功获取返回true，被限流返回False
	Try(count int32) bool
}
