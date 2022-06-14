package error

type ErrorCodeType int

const (
	ErrorCodeBaseCommon       ErrorCodeType = 0
	ErrorCodeBaseUserService  ErrorCodeType = 100000
	ErrorCodeBaseGoodsService ErrorCodeType = 200000
	ErrorCodeBaseBiService    ErrorCodeType = 300000
	ErrorCodeBaseSdkCommon    ErrorCodeType = 400000
)

const (
	ErrorCodeOK         ErrorCodeType = 0
	ErrorCodeWrongParam ErrorCodeType = 1
	ErrorCodeSystem     ErrorCodeType = 2
	ErrorCodeUnknown    ErrorCodeType = 3
	ErrorCodeEnd        ErrorCodeType = 4
)
