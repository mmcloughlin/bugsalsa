CC=gcc
CFLAGS=-Wall -O2 -g

all: test.nacl test.sodium

clean:
	$(RM) *.nacl *.sodium

%.nacl: %.c
	$(CC) -o $@ $(CFLAGS) -I/usr/include/nacl $^ -lnacl

%.sodium: %.c
	$(CC) -o $@ $(CFLAGS) -I/usr/include/sodium $^ -lsodium

key.%: test.%
	./$< > $@
