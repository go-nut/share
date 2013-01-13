
build:
	go build .

install:
	go install -p 3 .

gofmt:
	gofmt -w *.go

loc:
	find ./ -name '*.go' -print |sort |xargs wc -l

tags:
	find ./ -name '*.go' -print0 |xargs -0 gotags > TAGS

push:
	git push origin master
