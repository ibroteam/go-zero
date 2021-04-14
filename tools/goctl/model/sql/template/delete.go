package template

// Delete defines a delete template
var Delete = `
func (m *default{{.upperStartCamelObject}}Model) Delete(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) error {
	{{if .withCache}}{{if .containsIndexCache}}data, err:=m.FindOne({{.lowerStartCamelPrimaryKey}})
	if err!=nil{
		return err
	}{{end}}

	{{.keys}}
    _, err {{if .containsIndexCache}}={{else}}:={{end}} m.Exec(ctx, func(ctx2 context.Context, conn sqlx.SqlConnCtx) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where {{.originalPrimaryKey}} = ?", m.table)
		return conn.Exec(ctx2, query, {{.lowerStartCamelPrimaryKey}})
	}, {{.keyValues}}){{else}}query := fmt.Sprintf("delete from %s where {{.originalPrimaryKey}} = ?", m.table)
		_,err:=m.conn.Exec(ctx, query, {{.lowerStartCamelPrimaryKey}}){{end}}
	return err
}
`

// DeleteMethod defines a delete template for interface method
var DeleteMethod = `Delete(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) error`
