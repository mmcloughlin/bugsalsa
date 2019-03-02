#ifndef PARAMS_H_GUARD
#define PARAMS_H_GUARD

//
// Salsa20 block size.
//
#define BLOCK_SIZE UINT64_C(64)

//
// Number of bytes at which the 32-bit counter overflows.
//
#define OVERFLOW_POSITION (BLOCK_SIZE << 32)

//
// How many bytes after overflow does it take for the high 32-bits to zero.
//
#define WRAP_BYTES (((32 / 3) + 1) * 256)

//
// Length required to observe keystream cycle.
//
#define CYCLE_MIN_LENGTH (OVERFLOW_POSITION + WRAP_BYTES)

//
// Byte sizes for convenience.
//
#define KiB (UINT64_C(1) << 10)
#define MiB (KiB << 10)
#define GiB (MiB << 10)

#endif // PARAMS_H_GUARD
