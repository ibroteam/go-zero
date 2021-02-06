package serverinterceptors

import (
	"context"
	"fmt"
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
			title := "error@" + hostname
			content := fmt.Sprintf("%+v\n\nREQ:\n%s\n%+v\n", err, reflect.TypeOf(req).String(), req)

			fmt.Println(title)
			fmt.Println(content)
		}
		return result, err
		// return handler(ctx, req)
		// @CRACK end
	}
}
