module github.com/lubiedo/yav

go 1.16

replace github.com/lubiedo/yav/src/models => ./src/models

replace github.com/lubiedo/yav/src/utils => ./src/utils

require (
	github.com/Masterminds/sprig/v3 v3.2.2
	github.com/gomarkdown/markdown v0.0.0-20210408062403-ad838ccf8cdd
	github.com/google/uuid v1.2.0 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/mitchellh/copystructure v1.1.2 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2 // indirect
	gopkg.in/yaml.v2 v2.4.0
)
