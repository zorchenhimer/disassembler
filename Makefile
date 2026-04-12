.PHONY: all clean run win

all: bin/dasm
win: bin/dasm.exe
run: all
	bin/dasm testdata/sbx.cfg

clean:
	-rm bin/* testdata/*.asm

bin/dasm: cmd/main.go *.go config/*.go instr_6502/*.go instr_sbx/*.go types/*.go
	go build -o $@ $<

bin/dasm.exe: cmd/main.go *.go config/*.go instr_6502/*.go instr_sbx/*.go types/*.go
	GOOS=windows GOARCH=amd64 go build -o $@ $<
