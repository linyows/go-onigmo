#ifndef GONIGMO_H
#define GONIGMO_H

#include <onigmo.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

extern int GonigNewDefault(char *pattern, int pattern_length, int option, OnigRegex *regex,
                           OnigRegion **region, OnigErrorInfo **error_info, char **error_buffer);

extern int GonigSearch(void *str, int str_length, int offset, int option, OnigRegex regex,
                       OnigRegion *region, OnigErrorInfo *error_info, char *error_buffer,
                       int *captures, int *numCaptures);

extern int GonigMatch(void *str, int str_length, int offset, int option, OnigRegex regex,
                      OnigRegion *region);

#endif /* GONIGMO_H */
