all: compile docker push clean

compile:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o goregtrain

docker:
	docker build --force-rm=true -t payload0911/goregtrain:single .

push:
	docker push payload0911/goregtrain:single

clean:
	rm goregtrain
