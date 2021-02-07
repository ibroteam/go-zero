package logx

type WaringRobotConf struct {
	NotifyUrl                 string
	Secret                    string
	ReportIntervalLimitMillis int `json:",default=2000"`
}

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
	WaringRobot WaringRobotConf `json:"WaringRobot,optional"`
}

type LogConf struct {
	ServiceName         string  `json:",optional"`
	Mode                string  `json:",default=console,options=console|file|volume|sls"`
	Path                string  `json:",default=logs"`
	Level               string  `json:",default=info,options=info|error|severe"`
	Compress            bool    `json:",optional"`
	KeepDays            int     `json:",optional"`
	StackCooldownMillis int     `json:",default=100"`
	Sls                 SlsConf `json:"sls,optional"`
}
