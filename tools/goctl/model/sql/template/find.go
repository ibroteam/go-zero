package template

// FindOne defines find row by id.
var FindOne = `
func (m *default{{.upperStartCamelObject}}Model) FindOne(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) (*{{.upperStartCamelObject}}, error) {
	{{if .withCache}}{{.cacheKey}}
	var resp {{.upperStartCamelObject}}
	err := m.QueryRow(ctx, &resp, {{.cacheKeyVariable}}, func(ctx2 context.Context, conn sqlx.SqlConnCtx, v interface{}) error {
		query :=  fmt.Sprintf("select %s from %s where {{.originalPrimaryKey}} = ? limit 1", {{.lowerStartCamelObject}}Rows, m.table)
		return conn.QueryRow(ctx2, v, query, {{.lowerStartCamelPrimaryKey}})
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}{{else}}query := fmt.Sprintf("select %s from %s where {{.originalPrimaryKey}} = ? limit 1", {{.lowerStartCamelObject}}Rows, m.table)
	var resp {{.upperStartCamelObject}}
	err := m.conn.QueryRow(ctx, &resp, query, {{.lowerStartCamelPrimaryKey}})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}{{end}}
}
`

// FindOneByField defines find row by field.
var FindOneByField = `
func (m *default{{.upperStartCamelObject}}Model) FindOneBy{{.upperField}}(ctx context.Context, {{.in}}) (*{{.upperStartCamelObject}}, error) {
    var resp {{.upperStartCamelObject}}
	{{if .withCache}}{{.cacheKey}}
	err := m.QueryRowIndex(ctx, &resp, {{.cacheKeyVariable}}, m.formatPrimary, func(ctx2 context.Context, conn sqlx.SqlConnCtx, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where {{.originalField}} limit 1", {{.lowerStartCamelObject}}Rows, m.table)
		if err := conn.QueryRow(ctx2, &resp, query, {{.lowerStartCamelField}}); err != nil {
			return nil, err
		}
		return resp.{{.upperStartCamelPrimaryKey}}, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}{{else}}var resp {{.upperStartCamelObject}}
	query := fmt.Sprintf("select %s from %s where {{.originalField}} limit 1", {{.lowerStartCamelObject}}Rows, m.table )
	err := m.conn.QueryRow(ctx, &resp, query, {{.lowerStartCamelField}})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}{{end}}
`

// FindOneByFieldExtraMethod defines find row by field with extras.
var FindOneByFieldExtraMethod = `
func (m *default{{.upperStartCamelObject}}Model) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", {{.primaryKeyLeft}}, primary)
}

func (m *default{{.upperStartCamelObject}}Model) queryPrimary(ctx context.Context, conn sqlx.SqlConnCtx, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where {{.originalPrimaryField}} = ? limit 1", {{.lowerStartCamelObject}}Rows, m.table )
	return conn.QueryRow(ctx, v, query, primary)
}
`

// FindOneMethod defines find row method.
var FindOneMethod = `FindOne(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) (*{{.upperStartCamelObject}}, error)`

// FindOneByFieldMethod defines find row by field method.
var FindOneByFieldMethod = `FindOneBy{{.upperField}}(ctx context.Context, {{.in}}) (*{{.upperStartCamelObject}}, error) `
