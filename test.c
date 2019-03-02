#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <inttypes.h>
#include <assert.h>
#include <time.h>

#include "crypto_stream_salsa20.h"

#define BLOCK_SIZE (64ul)
#define WRAP_POSITION (BLOCK_SIZE << 32)
#define COMPARE_SIZE (512ul)
#define CIPHER_LENGTH (WRAP_POSITION + (4*COMPARE_SIZE))

static void log_print(const char* message)
{
	time_t now;
	time(&now);
	char buf[sizeof "2011-10-08T07:07:09Z"];
	strftime(buf, sizeof buf, "%FT%TZ", gmtime(&now));
	fprintf(stderr, "%s\t%s\n", buf, message);
}

int main()
{
	log_print("start");

	// Produce keystream of length CIPHER_LENGTH.
	uint8_t nonce[crypto_stream_salsa20_NONCEBYTES] = {0};
	uint8_t key[crypto_stream_salsa20_KEYBYTES] = {0};

	uint8_t *buffer = (uint8_t *)calloc(CIPHER_LENGTH, sizeof(uint8_t));
	assert(buffer != NULL);
	log_print("buffer allocated");

	int status = crypto_stream_salsa20_xor(buffer, buffer, CIPHER_LENGTH, nonce, key);
	assert(status == 0);
	log_print("keystream written");

	// Dump it.
	fwrite(buffer, CIPHER_LENGTH, sizeof(uint8_t), stdout);
	log_print("output written");

	// Cleanup.
	free(buffer);
	log_print("finish");
}
