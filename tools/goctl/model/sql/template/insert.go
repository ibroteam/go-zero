package template

// Insert defines a template for insert code in model
var Insert = `
func (m *default{{.upperStartCamelObject}}Model) Insert(ctx context.Context, data {{.upperStartCamelObject}}) (sql.Result,error) {
	{{if .withCache}}{{if .containsIndexCache}}{{.keys}}
    ret, err := m.Exec(ctx, func(ctx2 context.Context, conn sqlx.SqlConnCtx) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values ({{.expression}})", m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet)
		return conn.Exec(ctx2, query, {{.expressionValues}})
	}, {{.keyValues}}){{else}}query := fmt.Sprintf("insert into %s (%s) values ({{.expression}})", m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet)
    ret,err:=m.ExecNoCache(ctx, query, {{.expressionValues}})
	{{end}}{{else}}query := fmt.Sprintf("insert into %s (%s) values ({{.expression}})", m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet)
    ret,err:=m.conn.Exec(ctx, query, {{.expressionValues}}){{end}}
	return ret,err
}
`

// InsertMethod defines a interface method template for insert code in model
var InsertMethod = `Insert(ctx context.Context, data {{.upperStartCamelObject}}) (sql.Result,error)`
