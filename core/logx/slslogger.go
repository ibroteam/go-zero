package logx

import (
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"github.com/gogo/protobuf/proto"
	"github.com/tal-tech/go-zero/core/dingtalk"
	"github.com/tal-tech/go-zero/core/dingtalk/message"
	"github.com/tal-tech/go-zero/core/netx"
	"github.com/tal-tech/go-zero/core/sysx"
	"time"
)

type slsWriter struct {
	*limitedExecutor
	project            string
	logStore           string
	source             string
	topic              string
	producer           *producer.Producer
	hasRobotWarning    bool
	warningRobotUrl    string
	warningRobotSecret string
}

func newSlsWriter(AppName, Endpoint, Project, AccessKeyID, AccessKeySecret, LogStore string, warn *WaringRobotConf) *slsWriter {
	localIp := netx.InternalIp()

	producerConfig := producer.GetDefaultProducerConfig()
	producerConfig.Endpoint = Endpoint
	producerConfig.AccessKeyID = AccessKeyID
	producerConfig.AccessKeySecret = AccessKeySecret
	producerInstance := producer.InitProducer(producerConfig)

	l := &slsWriter{
		project:  Project,
		logStore: LogStore,
		source:   localIp,
		topic:    AppName,
		producer: producerInstance,
	}

	if warn != nil && len(warn.NotifyUrl) > 0 && len(warn.Secret) > 0 {
		l.hasRobotWarning = true
		l.warningRobotUrl = warn.NotifyUrl
		l.warningRobotSecret = warn.Secret
		l.limitedExecutor = newLimitedExecutor(warn.ReportIntervalLimitMillis)
	}

	producerInstance.Start()
	return l
}

func (l *slsWriter) Close() error {
	l.producer.SafeClose()
	return nil
}

func (l *slsWriter) Write(data []byte) (int, error) {
	log := &sls.Log{
		Time: proto.Uint32(uint32(time.Now().Unix())),
		Contents: []*sls.LogContent{
			{
				Key:   proto.String("raw"),
				Value: proto.String(string(data)),
			},
		},
	}

	if l.hasRobotWarning {
		l.logOrDiscard(func() {
			title := "error@" + sysx.Hostname()
			content := string(data)
			go dingtalk.SendRobotMessage(l.warningRobotUrl, l.warningRobotSecret, message.NewMarkdownMessageGeneral(title, content))
		})
	}

	err := l.producer.SendLog(l.project, l.logStore, l.topic, l.source, log)
	return len(data), err
}
