all: push
TAG = latest
PREFIX = operator-template
build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -o controller-template .
container: build
	docker build -t $(PREFIX):$(TAG) .
push: container
	docker push $(PREFIX):$(TAG)
clean:
	rm -f controller-template

