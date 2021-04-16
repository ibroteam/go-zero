package generator

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/core/collection"
	conf "github.com/tal-tech/go-zero/tools/goctl/config"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/parser"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/format"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
)

const (
	serverTemplate = `{{.head}}

package server

import (
	"context"

	{{.imports}}
)

type {{.server}}Server struct {
	svcCtx *svc.ServiceContext
}

func New{{.server}}Server(svcCtx *svc.ServiceContext) *{{.server}}Server {
	return &{{.server}}Server{
		svcCtx: svcCtx,
	}
}

{{.funcs}}
`
	functionTemplate = `
{{if .hasComment}}{{.comment}}{{end}}
func (s *{{.server}}Server) {{.method}} (ctx context.Context, in {{.request}}) ({{.response}}, error) {
	l := logic.New{{.logicName}}(ctx, s.svcCtx)
	return l.{{.method}}(in)
}
`

	functionStreamTemplate = `
{{if .hasComment}}{{.comment}}{{end}}
func (s *{{.server}}Server) {{.method}} (in {{.request}}) error {
	l := logic.New{{.logicName}}(in.Context(), s.svcCtx)
	req, err := in.Recv()
	if err != nil {
		return err
	}

	resp, err := l.SignIn(req)
	if err != nil {
		return err
	}

	in.Send(resp)
	return nil
}
`
)

// GenServer generates rpc server file, which is an implementation of rpc server
func (g *DefaultGenerator) GenServer(ctx DirContext, proto parser.Proto, cfg *conf.Config) error {
	dir := ctx.GetServer()
	logicImport := fmt.Sprintf(`"%v"`, ctx.GetLogic().Package)
	svcImport := fmt.Sprintf(`"%v"`, ctx.GetSvc().Package)
	pbImport := fmt.Sprintf(`"%v"`, ctx.GetPb().Package)

	imports := collection.NewSet()
	imports.AddStr(logicImport, svcImport, pbImport)

	head := util.GetHead(proto.Name)
	service := proto.Service
	serverFilename, err := format.FileNamingFormat(cfg.NamingFormat, service.Name+"_server")
	if err != nil {
		return err
	}

	serverFile := filepath.Join(dir.Filename, serverFilename+".go")
	funcList, err := g.genFunctions(proto.PbPackage, service)
	if err != nil {
		return err
	}

	text, err := util.LoadTemplate(category, serverTemplateFile, serverTemplate)
	if err != nil {
		return err
	}

	err = util.With("server").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"head":    head,
		"server":  stringx.From(service.Name).ToCamel(),
		"imports": strings.Join(imports.KeysStr(), util.NL),
		"funcs":   strings.Join(funcList, util.NL),
	}, serverFile, true)
	return err
}

func (g *DefaultGenerator) genFunctions(goPackage string, service parser.Service) ([]string, error) {
	var (
		err          error
		text         string
		request      string
		functionList []string
	)

	for _, rpc := range service.RPC {
		if rpc.StreamsReturns {
			text, err = util.LoadTemplate(category, serverFuncTemplateFile, functionStreamTemplate)
			request = fmt.Sprintf("%s.%s_%sServer", goPackage, stringx.From(service.Name).ToCamel(), parser.CamelCase(rpc.Name))
		} else {
			text, err = util.LoadTemplate(category, serverFuncTemplateFile, functionTemplate)
			request = fmt.Sprintf("*%s.%s", goPackage, parser.CamelCase(rpc.RequestType))
		}

		if err != nil {
			return nil, err
		}

		comment := parser.GetComment(rpc.Doc())
		buffer, err := util.With("func").Parse(text).Execute(map[string]interface{}{
			"server":     stringx.From(service.Name).ToCamel(),
			"logicName":  fmt.Sprintf("%sLogic", stringx.From(rpc.Name).ToCamel()),
			"method":     parser.CamelCase(rpc.Name),
			"request":    request,
			"response":   fmt.Sprintf("*%s.%s", goPackage, parser.CamelCase(rpc.ReturnsType)),
			"hasComment": len(comment) > 0,
			"comment":    comment,
		})
		if err != nil {
			return nil, err
		}

		functionList = append(functionList, buffer.String())
	}
	return functionList, nil
}
