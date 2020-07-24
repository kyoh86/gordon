module github.com/kyoh86/gordon

go 1.13

require (
	github.com/alecthomas/kingpin v2.2.6+incompatible
	github.com/blang/semver v3.5.1+incompatible
	github.com/golang/mock v1.4.1
	github.com/google/go-github/v29 v29.0.3
	github.com/kyoh86/appenv v0.0.20
	github.com/kyoh86/ask v0.0.7
	github.com/kyoh86/gogh v1.7.0
	github.com/kyoh86/xdg v1.2.0
	github.com/morikuni/aec v1.0.0
	github.com/saracen/walker v0.1.1
	github.com/stoewer/go-strcase v1.2.0
	github.com/stretchr/testify v1.5.1
	github.com/zalando/go-keyring v0.1.0
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sys v0.0.0-20200722175500-76b94024e4b6 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
)

replace github.com/joeshaw/envdecode => github.com/kyoh86/envdecode v0.0.1
