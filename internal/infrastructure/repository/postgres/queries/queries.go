package queries

import "embed"

//go:generate go run github.com/istovpets/sqlset/cmd/sqlset-gen --dir=. --out=queries-gen.go --pkg=queries

//go:embed *
var QueriesFS embed.FS
