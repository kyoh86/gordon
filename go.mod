module github.com/kyoh86/gordon

go 1.13

require (
	github.com/alecthomas/kingpin v2.2.6+incompatible
	github.com/blang/semver v3.5.1+incompatible
	github.com/golang/mock v1.4.1
	github.com/google/go-github/v29 v29.0.3
	github.com/kyoh86/appenv v0.0.20
	github.com/kyoh86/ask v0.0.7
	github.com/kyoh86/gogh v1.5.4
	github.com/kyoh86/xdg v1.2.0
	github.com/morikuni/aec v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/saracen/walker v0.0.0-20191201085201-324a081bae7e
	github.com/stoewer/go-strcase v1.2.0
	github.com/stretchr/testify v1.5.1
	github.com/zalando/go-keyring v0.0.0-20200121091418-667557018717
	golang.org/x/crypto v0.0.0-20200403201458-baeed622b8d8 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sys v0.0.0-20200331124033-c3d80250170d // indirect
	golang.org/x/tools v0.0.0-20200316194252-fafb6e2e8a4a // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
)

replace github.com/joeshaw/envdecode => github.com/kyoh86/envdecode v0.0.1
