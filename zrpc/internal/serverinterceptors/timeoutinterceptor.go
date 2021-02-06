package serverinterceptors

import (
	"context"
	"fmt"
	"github.com/tal-tech/go-zero/core/dingtalk"
	"github.com/tal-tech/go-zero/core/dingtalk/message"
	"github.com/tal-tech/go-zero/core/logx"
	"os"
	"reflect"
	"time"

	"github.com/tal-tech/go-zero/core/contextx"
	"google.golang.org/grpc"
)

func UnaryTimeoutInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx, cancel := contextx.ShrinkDeadline(ctx, timeout)
		defer cancel()
		// @CRACK begin
		// 2021-02-06 by hujiachao
		// 当RPC运行时遇到异常，将其报出来
		result, err := handler(ctx, req)
		if err != nil {
			hostname, _ := os.Hostname()
			title := "RPC ERROR@" + hostname
			content := fmt.Sprintf("%+v\n\nREQ:\n%s\n%+v\n", err, reflect.TypeOf(req).String(), req)

			env := os.Getenv("Env")
			if env == "test" {
				secret := "SECf01e230ef9a6b902b41490dabc3786cd53e5d2fe2570c318a8fcbd40bc104085"
				url := "https://oapi.dingtalk.com/robot/send?access_token=46c57006994630f925fb8caf9bf655ee24a63be791dcf87b1eeec4760a3d0ffa"
				go dingtalk.SendRobotMessage(url, secret, message.NewMarkdownMessageGeneral(title, content))
			} else if env == "prod" {
				secret := "SEC682cf73cd54ec0eb77cb79219ba6e5c97eb308f2131bd260f1b7764db266cd17"
				url := "https://oapi.dingtalk.com/robot/send?access_token=71b8a94373e7c46a061c7eb8305de9d4f1eb940cc35098a26029ba758ba0bac9"
				go dingtalk.SendRobotMessage(url, secret, message.NewMarkdownMessageGeneral(title, content))
			}
			logx.Error(content)
		}
		return result, err
		// return handler(ctx, req)
		// @CRACK end
	}
}
