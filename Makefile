PACKAGES=$(shell go list ./... | grep -v docloco/vendor)

all: clean embed-data
	@go build -o bin/docloco

clean:
	@rm -fr bin/*
	@find . -name rice-box.go -delete

embed-data:
	@go get github.com/GeertJohan/go.rice/rice
	@rice embed-go -i ./docloco

test: checklist
	@go vet $(PACKAGES)
	@go test $(PACKAGES)

deploy:
	gcloud app deploy

gae-logs:
	@gcloud app logs tail -s default

gae-browse:
	@gcloud app browse

save-deps:
	@godep save -t $(PACKAGES)

checklist:
	@go get honnef.co/go/tools/cmd/staticcheck
	@go get honnef.co/go/tools/cmd/unused
	@go get honnef.co/go/tools/cmd/gosimple
	@staticcheck $(PACKAGES)
	@unused $(PACKAGES)
	@gosimple $(PACKAGES)

.PHONY: test deploy gae-logs gae-browse clean all save-deps checklist