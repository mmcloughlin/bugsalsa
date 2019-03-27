# bugsalsa
Investigating a 32-bit overflow bug in SUPERCOP-derived salsa20 implementations.

## Background

[cryptofuzz](https://github.com/mmcloughlin/cryptofuzz) revealed a vulnerability in the Go `x/crypto` implementation of salsa20. A misimplementation of the counter increment in assembly causes the counter to cycle shortly after `2^32` blocks.

Since the asssembly version was ported directly from SUPERCOP, it seemed probable that the bug was also present in the upstream and other derived libraries. This repository contains some adhoc code used to test whether the bug is also present in NaCL and libsodium.

The bug was [fixed by the Golang team](https://groups.google.com/forum/#!topic/golang-dev/1X7VG7FDw2A) on March 20, 2019. It is now clear that the bug was known by others some years ago; it was patched in libsodium, but an empty file called `warning-256gb` was deemed sufficient for SUPERCOP and no action was taken to NaCL.

## Environment

Tested on Ubuntu 18.04 with

```
sudo apt-get install libnacl-dev libsodium-dev
```

Note that these tests will allocate over 256 GiB of memory!

## NaCL and libsodium Comparison

The `test.c` program was compiled against NaCL and libsodium and run to produce
two keyfiles. Comparing them we see a difference shortly after 256 GiB:

```
$ cmp key.nacl key.sodium
key.nacl key.sodium differ: byte 274877907073, line 1073727201
```

Note this offset is `(1<<38) + 129`. This is to be expected, since even with
the incorrect high 32-bit counter increment, the first two blocks after the
overflow will still agree. Therefore the first bad block should be the third,
which starts at byte 128 after the overflow. (It's 129 in the `cmp` output
because it reports offsets starting at 1.)

## Keystream Cycle

Example output from `cycle.c`

```
$ ./cycle.nacl
2019-03-02T15:25:34Z	compare_size=33554432
2019-03-02T15:25:34Z	length=274911464192
2019-03-02T15:25:34Z	buffer allocated
2019-03-02T15:34:22Z	keystream written
2019-03-02T15:34:22Z	key repeat: [2752, 33557184) == [274877909696, 274911464128)
2019-03-02T15:34:31Z	finish
```
