package logx

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
}

type LogConf struct {
	ServiceName         string  `json:",optional"`
	Mode                string  `json:",default=console,options=console|file|volume"`
	Path                string  `json:",default=logs"`
	Level               string  `json:",default=info,options=info|error|severe"`
	Compress            bool    `json:",optional"`
	KeepDays            int     `json:",optional"`
	StackCooldownMillis int     `json:",default=100"`
	Sls                 SlsConf `json:"sls,optional"`
}
