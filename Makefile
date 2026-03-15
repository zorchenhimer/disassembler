.PHONY: all clean run

all: bin/dasm
run: all
	bin/dasm testdata/sbx.cfg

clean:
	-rm bin/* testdata/*.asm

bin/dasm: cmd/main.go *.go config/*.go instr_6502/*.go types/*.go
	go build -o $@ $<
