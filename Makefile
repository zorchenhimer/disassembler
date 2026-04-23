.PHONY: all clean run win

SRC = dasm.go \
	  formatter.go \
	  labelmanager.go \
	  config/*.go \
	  instr_6502/*.go \
	  instr_sbx/*.go \
	  types/*.go \
	  stats/*.go

all: bin/dasm
win: bin/dasm.exe
run: all
	bin/dasm testdata/sbx.cfg

test: all
	bin/dasm testdata/sbx_script.cfg

clean:
	-rm bin/* testdata/*.asm

bin/dasm: cmd/main.go *.go $(SRC)
	go build -o $@ $<

bin/dasm.exe: cmd/main.go *.go $(SRC)
	GOOS=windows GOARCH=amd64 go build -o $@ $<
