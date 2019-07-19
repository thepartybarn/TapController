ImageName = thepartybarn/production:tapcontroller
ImageNameWithVersion = $(ImageName)-$(GitVersion)
BuildDate := $(shell date -u +%Y%m%d.%H%M%S)
GitVersion = $(shell git describe --always --long --dirty=-test)

all: format test checkin build push

format:
	go get ./...
	go fmt *.go

test: format
	go test -v ./...

checkin: format test
	-rm main
	-git pull
	-git commit -a
	-git push

build: format
	env GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "-X main._buildDate=$(BuildDate) -X main._buildVersion=$(GitVersion)" -o main *.go

push: format test checkin build

	@echo Image Name: $(ImageName)
	@echo Image Name: $(ImageNameWithVersion)

	sudo docker build -t $(ImageName) .
	sudo docker tag $(ImageName) $(ImageNameWithVersion)
	sudo docker login

	sudo docker push $(ImageName)
	sudo docker push $(ImageNameWithVersion)
local: 
	env GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "-X main._buildDate=$(BuildDate) -X main._buildVersion=$(GitVersion)" -o main *.go
	sudo docker build -t $(ImageName) .
