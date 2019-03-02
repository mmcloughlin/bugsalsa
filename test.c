#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <inttypes.h>
#include <assert.h>

#include "crypto_stream_salsa20.h"

#include "params.h"
#include "log.h"

#define LENGTH CYCLE_MIN_LENGTH

int main()
{
	log_print("start");

	// Produce keystream of length LENGTH.
	uint8_t nonce[crypto_stream_salsa20_NONCEBYTES] = {0};
	uint8_t key[crypto_stream_salsa20_KEYBYTES] = {0};

	uint8_t *buffer = (uint8_t *)calloc(LENGTH, sizeof(uint8_t));
	assert(buffer != NULL);
	log_print("buffer allocated");

	int status = crypto_stream_salsa20_xor(buffer, buffer, LENGTH, nonce, key);
	assert(status == 0);
	log_print("keystream written");

	// Dump it.
	fwrite(buffer, LENGTH, sizeof(uint8_t), stdout);
	log_print("output written");

	// Cleanup.
	free(buffer);
	log_print("finish");
}
