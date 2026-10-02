import { act, renderHook } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { ApiError } from '../api/client';
import type { BinaryOperator } from './useCalculator';
import { useCalculator } from './useCalculator';

vi.mock('../api/client', async (importOriginal) => {
  const actual = await importOriginal<typeof import('../api/client')>();
  return { ...actual, calculate: vi.fn(), getOperations: vi.fn() };
});

import { calculate } from '../api/client';

const mockCalculate = vi.mocked(calculate);

beforeEach(() => {
  mockCalculate.mockReset();
});

afterEach(() => {
  vi.restoreAllMocks();
});

// pressOperator/pressEquals/pressSqrt are all async (they may call the API),
// so every call — even ones that only take the synchronous "just queue it"
// path — must be awaited inside act(), or React's test renderer can trip
// over an unresolved act() scope on the next synchronous update.
type Hook = ReturnType<typeof useCalculator>;

async function pressOperator(result: { current: Hook }, op: BinaryOperator) {
  await act(async () => {
    await result.current.pressOperator(op);
  });
}

async function pressEquals(result: { current: Hook }) {
  await act(async () => {
    await result.current.pressEquals();
  });
}

async function pressSqrt(result: { current: Hook }) {
  await act(async () => {
    await result.current.pressSqrt();
  });
}

describe('useCalculator', () => {
  it('starts idle showing 0', () => {
    const { result } = renderHook(() => useCalculator());
    expect(result.current.display).toBe('0');
    expect(result.current.expression).toBe('');
    expect(result.current.isLoading).toBe(false);
    expect(result.current.errorMessage).toBeNull();
  });

  it('builds the display as digits are pressed', () => {
    const { result } = renderHook(() => useCalculator());
    act(() => result.current.pressDigit('7'));
    act(() => result.current.pressDigit('7'));
    expect(result.current.display).toBe('77');
  });

  it('handles a decimal point without duplicating it', () => {
    const { result } = renderHook(() => useCalculator());
    act(() => result.current.pressDigit('7'));
    act(() => result.current.pressDecimal());
    act(() => result.current.pressDigit('5'));
    act(() => result.current.pressDecimal());
    expect(result.current.display).toBe('7.5');
  });

  it('backspace removes the last character', () => {
    const { result } = renderHook(() => useCalculator());
    act(() => {
      result.current.pressDigit('1');
      result.current.pressDigit('2');
      result.current.pressDigit('3');
    });
    act(() => result.current.pressBackspace());
    expect(result.current.display).toBe('12');
  });

  it('queues an operator without calling the API when there is no second operand yet', async () => {
    const { result } = renderHook(() => useCalculator());
    act(() => result.current.pressDigit('7'));
    await pressOperator(result, 'add');
    expect(result.current.expression).toBe('7 +');
    expect(result.current.display).toBe('7');
    expect(mockCalculate).not.toHaveBeenCalled();
  });

  it('replaces a queued operator if another operator is pressed before a second operand', async () => {
    const { result } = renderHook(() => useCalculator());
    act(() => result.current.pressDigit('2'));
    await pressOperator(result, 'add');
    await pressOperator(result, 'multiply');
    expect(result.current.expression).toBe('2 ×');
    expect(mockCalculate).not.toHaveBeenCalled();
  });

  it('resolves a binary operation on equals', async () => {
    mockCalculate.mockResolvedValue(5);
    const { result } = renderHook(() => useCalculator());
    act(() => result.current.pressDigit('2'));
    await pressOperator(result, 'add');
    act(() => result.current.pressDigit('3'));
    await pressEquals(result);
    expect(mockCalculate).toHaveBeenCalledWith('add', [2, 3]);
    expect(result.current.display).toBe('5');
    expect(result.current.expression).toBe('');
  });

  it('chains operators left-to-right like a physical calculator (2 + 3 x 4 = 20, not 14)', async () => {
    mockCalculate.mockResolvedValueOnce(5).mockResolvedValueOnce(20);
    const { result } = renderHook(() => useCalculator());
    act(() => result.current.pressDigit('2'));
    await pressOperator(result, 'add');
    act(() => result.current.pressDigit('3'));
    await pressOperator(result, 'multiply');
    expect(mockCalculate).toHaveBeenNthCalledWith(1, 'add', [2, 3]);
    expect(result.current.display).toBe('5');
    expect(result.current.expression).toBe('5 ×');

    act(() => result.current.pressDigit('4'));
    await pressEquals(result);
    expect(mockCalculate).toHaveBeenNthCalledWith(2, 'multiply', [5, 4]);
    expect(result.current.display).toBe('20');
  });

  it('applies sqrt immediately to the current entry without disturbing a pending operator', async () => {
    mockCalculate.mockResolvedValueOnce(4).mockResolvedValueOnce(13);
    const { result } = renderHook(() => useCalculator());
    act(() => result.current.pressDigit('9'));
    await pressOperator(result, 'add');
    act(() => result.current.pressDigit('1'));
    act(() => result.current.pressDigit('6'));

    await pressSqrt(result);
    expect(mockCalculate).toHaveBeenNthCalledWith(1, 'sqrt', [16]);
    expect(result.current.display).toBe('4');
    expect(result.current.expression).toBe('9 +');

    await pressEquals(result);
    expect(mockCalculate).toHaveBeenNthCalledWith(2, 'add', [9, 4]);
    expect(result.current.display).toBe('13');
  });

  it('treats percentage as a regular chain operator', async () => {
    mockCalculate.mockResolvedValue(10);
    const { result } = renderHook(() => useCalculator());
    act(() => result.current.pressDigit('2'));
    act(() => result.current.pressDigit('0'));
    await pressOperator(result, 'percentage');
    act(() => result.current.pressDigit('5'));
    act(() => result.current.pressDigit('0'));
    await pressEquals(result);
    expect(mockCalculate).toHaveBeenCalledWith('percentage', [20, 50]);
    expect(result.current.display).toBe('10');
  });

  it('sets isLoading while a request is in flight', async () => {
    let resolvePromise!: (value: number) => void;
    mockCalculate.mockReturnValue(
      new Promise<number>((resolve) => {
        resolvePromise = resolve;
      }),
    );
    const { result } = renderHook(() => useCalculator());
    act(() => result.current.pressDigit('2'));
    await pressOperator(result, 'add');
    act(() => result.current.pressDigit('3'));

    let pending!: Promise<void>;
    act(() => {
      pending = result.current.pressEquals();
    });
    expect(result.current.isLoading).toBe(true);

    await act(async () => {
      resolvePromise(5);
      await pending;
    });
    expect(result.current.isLoading).toBe(false);
  });

  it('resets on failure when resolving a chained operator (not just on equals)', async () => {
    mockCalculate.mockRejectedValue(new ApiError('RESULT_OUT_OF_RANGE', 'result is not a finite number'));
    const { result } = renderHook(() => useCalculator());
    act(() => result.current.pressDigit('5'));
    await pressOperator(result, 'add');
    act(() => result.current.pressDigit('3'));
    await pressOperator(result, 'multiply');

    expect(result.current.errorMessage).toBe('Result too large to display');
    expect(result.current.display).toBe('0');
  });

  it('resets on failure when sqrt itself fails', async () => {
    mockCalculate.mockRejectedValue(new ApiError('MATH_DOMAIN_ERROR', 'operation is undefined for the given operands'));
    const { result } = renderHook(() => useCalculator());
    act(() => result.current.pressDigit('4'));
    await pressSqrt(result);

    expect(result.current.errorMessage).toBe('Undefined for these numbers');
    expect(result.current.display).toBe('0');
  });

  it('on API error, shows a friendly message and resets the chain entirely', async () => {
    mockCalculate.mockRejectedValue(new ApiError('DIVISION_BY_ZERO', 'division by zero'));
    const { result } = renderHook(() => useCalculator());
    act(() => result.current.pressDigit('5'));
    await pressOperator(result, 'divide');
    act(() => result.current.pressDigit('0'));
    await pressEquals(result);

    expect(result.current.errorMessage).toBe("Can't divide by zero");
    expect(result.current.display).toBe('0');
    expect(result.current.expression).toBe('');

    act(() => result.current.pressDigit('3'));
    expect(result.current.errorMessage).toBeNull();
    expect(result.current.display).toBe('3');
  });

  it('equals with an empty second operand repeats the first operand (5 + = -> 10)', async () => {
    mockCalculate.mockResolvedValue(10);
    const { result } = renderHook(() => useCalculator());
    act(() => result.current.pressDigit('5'));
    await pressOperator(result, 'add');
    await pressEquals(result);
    expect(mockCalculate).toHaveBeenCalledWith('add', [5, 5]);
    expect(result.current.display).toBe('10');
  });

  it('clear resets everything', async () => {
    const { result } = renderHook(() => useCalculator());
    act(() => result.current.pressDigit('9'));
    await pressOperator(result, 'add');
    act(() => result.current.pressClear());
    expect(result.current.display).toBe('0');
    expect(result.current.expression).toBe('');
  });
});
