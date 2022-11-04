#include <stdlib.h>

#include "spall.h"

#ifndef SPALLWRAP
#define SPALLWRAP

SpallProfile *NewSpallProfile(uintptr_t goProfileHandle, double timestampUnit);
SpallBuffer *NewSpallBuffer(size_t size);

void FreeSpallProfile(SpallProfile *p);
void FreeSpallBuffer(SpallBuffer *b);

#endif
