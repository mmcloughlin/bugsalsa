# bugsalsa
Investigating a 32-bit overflow bug in SUPERCOP-derived salsa20 implementations

## Results

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

Note this offset is `(1<<38) + 129`.
