#include "wrap.h"
#include "spall.h"

extern bool gowrite(uintptr_t goProfileHandle, const void *p, size_t n);
extern bool goflush(uintptr_t goProfileHandle);
extern void goclose(uintptr_t goProfileHandle);

bool goSpall_Write(SpallProfile *ctx, const void *p, size_t n) {
    return gowrite((uintptr_t)ctx->userdata, p, n);
}

bool goSpall_Flush(SpallProfile *ctx) {
    return goflush((uintptr_t)ctx->userdata);
}

void goSpall_Close(SpallProfile *ctx) {
    goclose((uintptr_t)ctx->userdata);
}

SpallProfile *NewSpallProfile(uintptr_t goProfileHandle, double timestampUnit) {
    SpallProfile *p = malloc(sizeof(SpallProfile));
    *p = (SpallProfile) {
		.timestamp_unit = timestampUnit,

        .write = goSpall_Write,
        .flush = goSpall_Flush,
        .close = goSpall_Close,

        .userdata = (void*)goProfileHandle,
	};

    // Copy-pasted from spall.h. Very sad.
    SpallHeader header;
    header.magic_header = 0x0BADF00D;
    header.version = 1;
    header.timestamp_unit = p->timestamp_unit;
    header.must_be_0 = 0;
    if (!p->write(p, &header, sizeof(header))) {
        SpallQuit(p);
        return p;
    }

    return p;
}

void FreeSpallProfile(SpallProfile *p) {
    free(p);
}

SpallBuffer *NewSpallBuffer(size_t size) {
    SpallBuffer *b = malloc(sizeof(SpallBuffer));
    *b = (SpallBuffer) {
        .data = malloc(size),
        .length = size,
    };
    return b;
}

void FreeSpallBuffer(SpallBuffer *b) {
    free(b->data);
    free(b);
}
