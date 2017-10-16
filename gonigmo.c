#include "gonigmo.h"

int GonigNewDefault(char *pattern, int pattern_length, int option, OnigRegex *regex,
                    OnigRegion **region, OnigErrorInfo **error_info, char **error_buffer)
{
  int result = ONIG_NORMAL;
  int error_msg_len = 0;

  OnigUChar *pattern_start = (OnigUChar *)pattern;
  OnigUChar *pattern_end = (OnigUChar *)(pattern + pattern_length);

  *error_info = (OnigErrorInfo *)malloc(sizeof(OnigErrorInfo));
  memset(*error_info, 0, sizeof(OnigErrorInfo));

  *error_buffer = (char *)malloc(ONIG_MAX_ERROR_MESSAGE_LEN * sizeof(char));
  memset(*error_buffer, 0, ONIG_MAX_ERROR_MESSAGE_LEN * sizeof(char));

  *region = onig_region_new();

  result =
      onig_new_default(regex, pattern_start, pattern_end, (OnigOptionType)(option), *error_info);

  if (result != ONIG_NORMAL) {
    error_msg_len = onig_error_code_to_str((unsigned char *)(*error_buffer), result, *error_info);
    if (error_msg_len >= ONIG_MAX_ERROR_MESSAGE_LEN) {
      error_msg_len = ONIG_MAX_ERROR_MESSAGE_LEN - 1;
    }
    (*error_buffer)[error_msg_len] = "\0";
  }

  return result;
}

int GonigSearch(void *str, int str_length, int offset, int option, OnigRegex regex,
                OnigRegion *region, OnigErrorInfo *error_info, char *error_buffer, int *captures,
                int *numCaptures)
{
  int result = ONIG_MISMATCH;
  int error_msg_len = 0;

  OnigUChar *str_start = (OnigUChar *)str;
  OnigUChar *str_end = (OnigUChar *)(str_start + str_length);
  OnigUChar *search_start = (OnigUChar *)(str_start + offset);
  OnigUChar *search_end = str_end;

  result = onig_search(regex, str_start, str_end, search_start, search_end, region, option);

  if (result < 0 && error_buffer != NULL) {
    error_msg_len = onig_error_code_to_str((unsigned char *)(error_buffer), result, error_info);
    if (error_msg_len >= ONIG_MAX_ERROR_MESSAGE_LEN) {
      error_msg_len = ONIG_MAX_ERROR_MESSAGE_LEN - 1;
    }
    error_buffer[error_msg_len] = "\0";

  } else if (captures != NULL) {
    int i;
    int count = 0;
    for (i = 0; i < region->num_regs; i++) {
      captures[2 * count] = region->beg[i];
      captures[2 * count + 1] = region->end[i];
      count++;
    }
    *numCaptures = count;
  }

  return result;
}

int GonigMatch(void *str, int str_length, int offset, int option, OnigRegex regex,
               OnigRegion *region)
{
  int result = ONIG_MISMATCH;
  int error_msg_len = 0;

  OnigUChar *str_start = (OnigUChar *)str;
  OnigUChar *str_end = (OnigUChar *)(str_start + str_length);
  OnigUChar *search_start = (OnigUChar *)(str_start + offset);

  result = onig_match(regex, str_start, str_end, search_start, region, option);

  return result;
}
