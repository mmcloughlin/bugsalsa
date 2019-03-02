# bugsalsa
Investigating a 32-bit overflow bug in SUPERCOP-derived salsa20 implementations

## NaCL and libsodium

The following was compiled and run against NaCL and libsodium.

```c
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <inttypes.h>
#include <assert.h>

#include "crypto_stream_salsa20.h"

#define BLOCK_SIZE (64ul)
#define WRAP_POSITION (BLOCK_SIZE << 32)
#define COMPARE_SIZE (512ul)
#define CIPHER_LENGTH (WRAP_POSITION + (4*COMPARE_SIZE))

int main()
{
	// Produce keystream of length CIPHER_LENGTH.
	uint8_t nonce[crypto_stream_salsa20_NONCEBYTES] = {0};
	uint8_t key[crypto_stream_salsa20_KEYBYTES] = {0};

	uint8_t *buffer = (uint8_t *)calloc(CIPHER_LENGTH, sizeof(uint8_t));
	assert(buffer != NULL);
	fprintf(stderr, "buffer allocated\n");

	int status = crypto_stream_salsa20_xor(buffer, buffer, CIPHER_LENGTH, nonce, key);
	assert(status == 0);
	fprintf(stderr, "keystream written\n");

	// Dump it.
	fwrite(buffer, CIPHER_LENGTH, sizeof(uint8_t), stdout);
	fprintf(stderr, "output written\n");

	// Cleanup.
	free(buffer);
	fprintf(stderr, "finish\n");
}
```

Comparing the two files we get:

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
