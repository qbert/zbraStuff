default: build

build:
	go fmt
	go vet
	go build --ldflags="-X github.com/qbert/heartbeat-golang.CommitHash=`git rev-parse HEAD`"

test: build
	go test

clean:
	rm zbraStuff
	rm coverage.out

deploy: build
	git add -u
	git commit -m "Deployment"
	git push openshift master

github: build
	git add -u
	git commit -m "Deployment"
	git push github master

coverage-test:
	go test -coverprofile=coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out