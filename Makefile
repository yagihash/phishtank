JOB=build

build: test
	@ go build

test:
	@ go vet ./...
	@ richgo test -v -cover -race ./...

coverage:
	@ richgo test -v -race -coverprofile=/tmp/profile -covermode=atomic ./...
	@ go tool cover -html=/tmp/profile

validate-ci-config:
	@ circleci config validate -c .circleci/config.yml

local-ci:
	@ circleci local execute --job $(JOB)

