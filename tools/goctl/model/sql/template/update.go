package template

// Update defines a template for generating update codes
var Update = `
func (m *default{{.upperStartCamelObject}}Model) Update(ctx context.Context, data {{.upperStartCamelObject}}) error {
	{{if .withCache}}{{.keys}}
    _, err := m.Exec(ctx, func(ctx2 context.Context, conn sqlx.SqlConnCtx) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where {{.originalPrimaryKey}} = ?", m.table, {{.lowerStartCamelObject}}RowsWithPlaceHolder)
		return conn.Exec(ctx2, query, {{.expressionValues}})
	}, {{.keyValues}}){{else}}query := fmt.Sprintf("update %s set %s where {{.originalPrimaryKey}} = ?", m.table, {{.lowerStartCamelObject}}RowsWithPlaceHolder)
    _,err:=m.conn.Exec(ctx, query, {{.expressionValues}}){{end}}
	return err
}
`

// UpdateMethod defines an interface method template for generating update codes
var UpdateMethod = `Update(ctx context.Context, data {{.upperStartCamelObject}}) error`
