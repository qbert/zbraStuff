default: build

build:
	go fmt
	go vet
	go build 

test: build
	go test

clean:
	rm zbraStuff

deploy: build
	git add -u
	git commit -m "Deployment"
	git push openshift

github: build
	git add -u
	git commit -m "Deployment"
	git push github

coverage-test:
	go test -coverprofile=coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out