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

func (m *default{{.upperStartCamelObject}}Model) UpdateSomeByPrimaryId(ctx context.Context, primaryId int64, updateFields map[string]interface{}) error {
	sb := strings.Builder{}
	sz := len(updateFields)
	params := make([]interface{}, 0, sz+1)
	sb.WriteString("update " + m.table + " set ")
	for k, v := range updateFields {
		sb.WriteString("," + k + "=?")
		params = append(params, v)
	}
	params = append(params, primaryId)
	sb.WriteString(" where {{.originalPrimaryKey}} = ?")
	{{if .withCache}}{{.keyValues}} := fmt.Sprintf("%s%v", {{.originalPrimaryKeyPrefix}}, primaryId)
    _, err := m.Exec(ctx, func(ctx2 context.Context, conn sqlx.SqlConnCtx) (result sql.Result, err error) {
		return conn.Exec(ctx2, sb.String()[1:], params...)
	}, {{.keyValues}}){{else}}
    _,err:=m.conn.Exec(ctx, sb.String()[1:], params...){{end}}
	return err
}
`

// UpdateMethod defines an interface method template for generating update codes
var UpdateMethod = `Update(ctx context.Context, data {{.upperStartCamelObject}}) error`
