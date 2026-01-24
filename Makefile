.PHONY: all clean run

all: bin/dasm
run: all
	bin/dasm testdata/sbx.cfg

clean:
	-rm bin/*

bin/dasm: cmd/main.go *.go config/*.go instructions/*.go types/*.go
	go build -o $@ $<
