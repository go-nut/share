
build:
	go build *.go

install:
	go install - 3 .

gofmt:
	gofmt -w *.go

loc:
	find ./ -name '*.go' -print |sort |xargs wc -l

tags:
	find ./ -name '*.go' -print0 |xargs -0 gotags > TAGS

push:
	git push origin master
