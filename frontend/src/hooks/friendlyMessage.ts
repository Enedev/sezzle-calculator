import { ApiError, NetworkError } from '../api/client';
import type { ErrorCode } from '../types/api';

const GENERIC_MESSAGE = 'Something went wrong. Try again.';

// Only DIVISION_BY_ZERO, MATH_DOMAIN_ERROR, and RESULT_OUT_OF_RANGE are
// reachable from a well-formed request this client ever sends. The rest of
// the backend's error codes are mapped defensively in case that ever
// changes, rather than left to fall through to a confusing message.
const MESSAGES: Record<ErrorCode, string> = {
  DIVISION_BY_ZERO: "Can't divide by zero",
  MATH_DOMAIN_ERROR: 'Undefined for these numbers',
  RESULT_OUT_OF_RANGE: 'Result too large to display',
  INVALID_OPERAND_COUNT: GENERIC_MESSAGE,
  UNKNOWN_OPERATION: GENERIC_MESSAGE,
  MISSING_FIELD: GENERIC_MESSAGE,
  INVALID_JSON: GENERIC_MESSAGE,
  UNKNOWN_FIELD: GENERIC_MESSAGE,
  UNSUPPORTED_MEDIA_TYPE: GENERIC_MESSAGE,
  REQUEST_TOO_LARGE: GENERIC_MESSAGE,
  NOT_FOUND: GENERIC_MESSAGE,
  METHOD_NOT_ALLOWED: GENERIC_MESSAGE,
  INTERNAL_ERROR: GENERIC_MESSAGE,
};

export function friendlyMessage(err: unknown): string {
  if (err instanceof ApiError) {
    return MESSAGES[err.code] ?? GENERIC_MESSAGE;
  }
  if (err instanceof NetworkError) {
    return "Can't reach the server";
  }
  return GENERIC_MESSAGE;
}
