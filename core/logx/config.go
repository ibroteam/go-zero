package logx

type CustomLoggerFn func(title, content string)

type SlsConf struct {
	AccessKeyID     string
	AccessKeySecret string
	Endpoint        string
	Project         string

	InfoStore   string
	ErrorStore  string
	StatStore   string
	SlowStore   string
	StackStore  string
	SevereStore string

	// 纯粹用于占位传参，下面的配置不通过xml，而是直接在main.go中根据需要传入对应的实现
	customLoggerFn              CustomLoggerFn `yaml:"-"` // 自定义的日志操作实现，例如钉钉机器人告警
	customLoggerIntervalLimitMs int            `yaml:"-"` // 设置间隔期，用于控制自定义的日志数量防止太多
}

func (m *SlsConf) SetCustomLogger(customLoggerFn CustomLoggerFn, customLoggerIntervalLimitMs int) {
	m.customLoggerIntervalLimitMs = customLoggerIntervalLimitMs
	m.customLoggerFn = customLoggerFn
}

// A LogConf is a logging config.
type LogConf struct {
	ServiceName         string  `json:",optional"`
	Mode                string  `json:",default=console,options=console|file|volume|sls"`
	TimeFormat          string  `json:",optional"`
	Path                string  `json:",default=logs"`
	Level               string  `json:",default=info,options=info|error|severe"`
	Compress            bool    `json:",optional"`
	KeepDays            int     `json:",optional"`
	StackCooldownMillis int     `json:",default=100"`
	Sls                 SlsConf `json:"sls,optional"`
}
