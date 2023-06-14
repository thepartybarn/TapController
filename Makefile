ImageName = thepartybarn/production:tapcontroller
ImageNameWithVersion = $(ImageName)-$(GitVersion)
BuildDate := $(shell date -u +%Y%m%d.%H%M%S)
GitVersion = $(shell git describe --always --long --dirty=-test)

push: buildgo builddocker
	docker login

	docker push $(ImageName)
	docker push $(ImageNameWithVersion)

builddocker:
	docker build -t $(ImageName) .
	docker tag $(ImageName) $(ImageNameWithVersion)

buildgo:
	env GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "-X main._buildDate=$(BuildDate) -X main._buildVersion=$(GitVersion)" -o main *.go

run: build
	docker build -t $(ImageName) .

local: buildgo builddocker
	docker save $(ImageName) > tapcontroller.tar
