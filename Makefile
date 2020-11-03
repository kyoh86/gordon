FORCE:
.PHONY: FORCE

VERSION := `git vertag get`
COMMIT  := `git rev-parse HEAD`

generate-clear: gen-clear
.PHONY: generate-clear
gen-clear:
	rm internal/**/*_gen.go
.PHONY: gen-clear

generate: gen
.PHONY: generate
gen: gen-clear internal/env/go_dist_gen.go
	go generate -x ./...
	go run "github.com/rjeczalik/interfaces/cmd/interfacer" -for github.com/kyoh86/gordon/internal/env.Access -as command.Env -o internal/command/env_gen.go
	go run "github.com/golang/mock/mockgen" -source internal/command/env_gen.go -destination internal/command/env_mock_test.go -package command_test
	
	go run "github.com/rjeczalik/interfaces/cmd/interfacer" -for github.com/kyoh86/gordon/internal/env.Access -as gordon.Env -o internal/gordon/env_gen.go
	go run "github.com/golang/mock/mockgen" -source internal/gordon/env_gen.go -destination internal/gordon/env_mock_test.go -package gordon_test
	
	go run "github.com/rjeczalik/interfaces/cmd/interfacer" -for github.com/kyoh86/gordon/internal/hub.Client -as command.HubClient -o internal/command/hub_gen.go
	go run "github.com/golang/mock/mockgen" -source internal/command/hub_gen.go -destination internal/command/hub_mock_test.go -package command_test
.PHONY: gen

internal/env/go_dist_gen.go: FORCE
	@echo "// Code generated by 'make gen'; DO NOT EDIT." > $@
	@echo "" >> $@
	@echo "package env" >> $@
	@echo "" >> $@
	@echo -n "var goDists = map[string]map[string]byte" >> $@
	@go tool dist list -json | jq -rcM 'reduce .[]as$$o({};.*{($$o.GOOS):{($$o.GOARCH):1}})' >> $@
	@echo -n "var goArchs = map[string]byte" >> $@
	@go tool dist list -json | jq -rcM 'reduce .[]as$$o({};.*{($$o.GOARCH):1})' >> $@
	@gofmt -s -w $@

lint: gen
	golangci-lint run
.PHONY: lint

test: lint
	go test -v --race ./...
.PHONY: test

install: test
	go install -a -ldflags "-X=main.version=$(VERSION) -X=main.commit=$(COMMIT)" ./...
.PHONY: install

man:
	go run . --help-man > gordon.1
.PHONY: man
