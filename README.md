## Setup config

https://www.strava.com/oauth/authorize?client_id={clientId}&response_type=code&redirect_uri=http://localhost/exchange_token&approval_prompt=force&scope=activity:read_all

## Help with golang

```
# Linter
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.53.3
export PATH="$PATH:$HOME/go/bin"
golangci-lint run
# Formatter
gofmt -w *.go # installed by default
```