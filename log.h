#ifndef LOG_H_GUARD
#define LOG_H_GUARD

#include <stdio.h>
#include <time.h>

#define log_print(format, ...) do { \
	time_t _now; \
	time(&_now); \
	char _buf[sizeof "2011-10-08T07:07:09Z"]; \
	strftime(_buf, sizeof _buf, "%FT%TZ", gmtime(&_now)); \
	fprintf(stderr, "%s\t" format "\n", _buf, ##__VA_ARGS__); \
} while(0)

#endif // LOG_H_GUARD
