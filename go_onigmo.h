#ifndef GONIGMO_H
#define GONIGMO_H

#include <onigmo.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

extern int go_onigmo();
extern int example();
extern int gonig_new(UChar* pattern, regex_t* reg);
extern int goning_search(UChar* str, regex_t* reg);

#endif /* GONIGMO_H */
