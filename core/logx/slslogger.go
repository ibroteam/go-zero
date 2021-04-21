package logx

import (
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"github.com/gogo/protobuf/proto"
	"github.com/tal-tech/go-zero/core/netx"
	"github.com/tal-tech/go-zero/core/sysx"
	"time"
)

type slsWriter struct {
	*limitedExecutor
	project        string
	logStore       string
	source         string
	topic          string
	producer       *producer.Producer
	customLoggerFn CustomLoggerFn
}

func newSlsWriter(AppName, Endpoint, Project, AccessKeyID, AccessKeySecret, LogStore string, customLoggerFn CustomLoggerFn, customLoggerIntervalLimitMs int) *slsWriter {
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
	if customLoggerFn != nil {
		l.customLoggerFn = customLoggerFn
		l.limitedExecutor = newLimitedExecutor(customLoggerIntervalLimitMs)
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

	if l.customLoggerFn != nil {
		l.logOrDiscard(func() {
			title := "error@" + sysx.Hostname()
			content := string(data)
			l.customLoggerFn(title, content)
		})
	}

	err := l.producer.SendLog(l.project, l.logStore, l.topic, l.source, log)
	return len(data), err
}
