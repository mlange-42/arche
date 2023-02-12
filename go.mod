module github.com/mlange-42/arche

go 1.20

require github.com/stretchr/testify v1.8.1

require internal/base v0.0.0

replace internal/base => ./internal/base

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
