%.o: %.asm
	ca65 -g -t nes -o $@ $<

%.bin: %.cfg %.o
	ld65 -o $@ -C $^

%.dasm.asm: %.dasm.cfg %.bin
	../bin/dasm $<
