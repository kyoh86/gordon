module github.com/kyoh86/gordon

go 1.13

require (
	github.com/99designs/keyring v1.1.6
	github.com/alecthomas/kingpin v2.2.6+incompatible
	github.com/blang/semver v3.5.1+incompatible
	github.com/dvsekhvalnov/jose2go v0.0.0-20201001154944-b09cfaf05951 // indirect
	github.com/fatih/color v1.10.0 // indirect
	github.com/goccy/go-yaml v1.8.3
	github.com/golang/mock v1.4.4
	github.com/google/go-github/v29 v29.0.3
	github.com/keybase/go-keychain v0.0.0-20200502122510-cda31fe0c86d // indirect
	github.com/kyoh86/appenv v0.1.0
	github.com/kyoh86/ask v0.0.7
	github.com/kyoh86/gogh v1.7.1
	github.com/kyoh86/xdg v1.2.0
	github.com/morikuni/aec v1.0.0
	github.com/saracen/walker v0.1.1
	github.com/stoewer/go-strcase v1.2.0
	github.com/stretchr/testify v1.6.1
	github.com/ulikunitz/xz v0.5.8
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
	golang.org/x/sys v0.0.0-20201101102859-da207088b7d1 // indirect
)

replace github.com/rjeczalik/interfaces => github.com/kyoh86/interfaces v0.1.2-0.20201103060818-43fd8be1be4d
