import { describe, expect, it } from 'vitest';
import { ApiError, NetworkError } from '../api/client';
import { friendlyMessage } from './friendlyMessage';

describe('friendlyMessage', () => {
  it('maps DIVISION_BY_ZERO to a plain, direct message', () => {
    expect(friendlyMessage(new ApiError('DIVISION_BY_ZERO', 'division by zero'))).toBe(
      "Can't divide by zero",
    );
  });

  it('maps MATH_DOMAIN_ERROR to a plain, direct message', () => {
    expect(
      friendlyMessage(new ApiError('MATH_DOMAIN_ERROR', 'operation is undefined for the given operands')),
    ).toBe('Undefined for these numbers');
  });

  it('maps RESULT_OUT_OF_RANGE to a plain, direct message', () => {
    expect(friendlyMessage(new ApiError('RESULT_OUT_OF_RANGE', 'result is not a finite number'))).toBe(
      'Result too large to display',
    );
  });

  it('maps a NetworkError to a connectivity message', () => {
    expect(friendlyMessage(new NetworkError())).toBe("Can't reach the server");
  });

  it('falls back to a generic message for anything else', () => {
    expect(friendlyMessage(new Error('boom'))).toBe('Something went wrong. Try again.');
  });
});
