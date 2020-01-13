module github.com/kyoh86/gordon

go 1.13

require (
	github.com/alecthomas/kingpin v2.2.6+incompatible
	github.com/alessio/shellescape v0.0.0-20190409004728-b115ca0f9053
	github.com/comail/colog v0.0.0-20160416085026-fba8e7b1f46c
	github.com/google/go-github/v24 v24.0.1
	github.com/joeshaw/envdecode v0.0.0-20190604014844-d6d9849fcc2c
	github.com/karrick/godirwalk v1.14.0 // indirect
	github.com/kyoh86/ask v0.0.7
	github.com/kyoh86/gogh v1.3.1
	github.com/kyoh86/xdg v1.2.0
	github.com/pkg/errors v0.8.1
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/stretchr/testify v1.4.0
	github.com/zalando/go-keyring v0.0.0-20200106095630-91fe8abcd771
	golang.org/x/crypto v0.0.0-20200109152110-61a87790db17 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sys v0.0.0-20200107162124-548cf772de50 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v3 v3.0.0-20191120175047-4206685974f2
)

replace github.com/joeshaw/envdecode => github.com/kyoh86/envdecode v0.0.1
