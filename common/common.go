package common

// Constants related to Prometheus metrics.
const (
	ConsumeCPUAddress       = "/ConsumeCPU"
	ConsumeMemAddress       = "/ConsumeMem"
	BumpMetricAddress       = "/BumpMetric"
	GetCurrentStatusAddress = "/GetCurrentStatus"
	MetricsAddress          = "/metrics"

	MillicoresQuery              = "millicores"
	MegabytesQuery               = "megabytes"
	MetricNameQuery              = "metric"
	DeltaQuery                   = "delta"
	DurationSecQuery             = "durationSec"
	RequestSizeInMillicoresQuery = "requestSizeMillicores"
	RequestSizeInMegabytesQuery  = "requestSizeMegabytes"
	RequestSizeCustomMetricQuery = "requestSizeMetrics"

	BadRequest                = "Bad request. Not a POST request"
	UnknownFunction           = "unknown function"
	IncorrectFunctionArgument = "incorrect function argument"
	NotGivenFunctionArgument  = "not given function argument"
)
