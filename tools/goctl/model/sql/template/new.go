package template

// New defines an template for creating model instance
var New = `
func New{{.upperStartCamelObject}}Model({{if .withCache}}conn sqlx.SqlConnCtx, c cache.CacheConf{{else}}conn sqlx.SqlConnCtx{{end}}) {{.upperStartCamelObject}}Model {
	return &default{{.upperStartCamelObject}}Model{
		{{if .withCache}}CachedConnCtx: sqlc.NewConnCtx(conn, c, cache.WithExpiry(time.Minute), cache.WithNotFoundExpiry(time.Minute)){{else}}conn:conn{{end}},
		table:      "{{.table}}",
	}
}
`
