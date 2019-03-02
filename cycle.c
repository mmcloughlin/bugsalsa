#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <inttypes.h>
#include <assert.h>

#include "crypto_stream_salsa20.h"

#include "params.h"
#include "log.h"

//
// How many repeated bytes do we want to see to confirm keystream cycle.
//
#define COMPARE_SIZE (32*MiB)

//
// How far beyond the overflow do we go?
//
#define EXTRA (WRAP_BYTES + COMPARE_SIZE)

//
// How much keystream do we need total?
//
#define LENGTH (OVERFLOW_POSITION + EXTRA)

int main()
{
	log_print("compare_size=%" PRId64, COMPARE_SIZE);
	log_print("length=%" PRId64, LENGTH);

	//
	// Produce keystream of length LENGTH.
	//
	uint8_t nonce[crypto_stream_salsa20_NONCEBYTES] = {0};
	uint8_t key[crypto_stream_salsa20_KEYBYTES] = {0};

	uint8_t *buffer = (uint8_t *)calloc(LENGTH, sizeof(uint8_t));
	assert(buffer != NULL);
	log_print("buffer allocated");

	int status = crypto_stream_salsa20_xor(buffer, buffer, LENGTH, nonce, key);
	assert(status == 0);
	log_print("keystream written");

	//
	// Look for keystream cycle.
	//
	for (size_t i = 0; i < EXTRA; i++)
	{
		for (size_t j = OVERFLOW_POSITION; j + COMPARE_SIZE <= LENGTH; j++)
		{
			if (memcmp(buffer + i, buffer + j, COMPARE_SIZE) == 0)
			{
				log_print("key repeat: [%" PRId64 ", %" PRId64 ") == [%" PRId64 ", %" PRId64 ")", i, i+COMPARE_SIZE, j, j+COMPARE_SIZE);
				goto finish;
			}
		}
	}

	//
	// Cleanup.
	//
finish:
	free(buffer);
	log_print("finish");
}
