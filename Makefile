export APP=mossgow
export LDFLAGS="-w -s"

detect:
	go run -ldflags $(LDFLAGS) ./cmd/mossgow detect
	./moss -d uploads/*/*/* uploads/*/*/*

build:
	go build -ldflags $(LDFLAGS) ./cmd/mossgow

install:
	go install -ldflags $(LDFLAGS) ./cmd/mossgow
