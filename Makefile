export APP=mossgow
export LDFLAGS="-w -s"

detect:
	go run -ldflags $(LDFLAGS) ./cmd/mossgow detect
	./moss -d uploads/*/*/* uploads/*/*/*

build:
	CGO_ENABLED=1 go build -ldflags $(LDFLAGS) ./cmd/mossgow

install:
	CGO_ENABLED=1 go install -ldflags $(LDFLAGS) ./cmd/mossgow
