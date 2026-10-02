import { renderHook } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { useKeyboardInput } from './useKeyboardInput';

function fireKey(key: string) {
  window.dispatchEvent(new KeyboardEvent('keydown', { key, cancelable: true }));
}

describe('useKeyboardInput', () => {
  const actions = {
    pressDigit: vi.fn(),
    pressDecimal: vi.fn(),
    pressOperator: vi.fn(),
    pressEquals: vi.fn(),
    pressSqrt: vi.fn(),
    pressClear: vi.fn(),
    pressBackspace: vi.fn(),
  };

  let unmount: () => void;

  beforeEach(() => {
    Object.values(actions).forEach((fn) => fn.mockReset());
    ({ unmount } = renderHook(() => useKeyboardInput(actions)));
  });

  afterEach(() => {
    unmount();
    vi.restoreAllMocks();
  });

  it('maps digit keys', () => {
    fireKey('7');
    expect(actions.pressDigit).toHaveBeenCalledWith('7');
  });

  it('maps the decimal point', () => {
    fireKey('.');
    expect(actions.pressDecimal).toHaveBeenCalled();
  });

  it.each([
    ['+', 'add'],
    ['-', 'subtract'],
    ['*', 'multiply'],
    ['/', 'divide'],
    ['^', 'power'],
    ['%', 'percentage'],
  ])('maps %s to the %s operator', (key, operator) => {
    fireKey(key);
    expect(actions.pressOperator).toHaveBeenCalledWith(operator);
  });

  it('maps Enter and = to equals', () => {
    fireKey('Enter');
    fireKey('=');
    expect(actions.pressEquals).toHaveBeenCalledTimes(2);
  });

  it('maps Escape to clear', () => {
    fireKey('Escape');
    expect(actions.pressClear).toHaveBeenCalled();
  });

  it('maps Backspace', () => {
    fireKey('Backspace');
    expect(actions.pressBackspace).toHaveBeenCalled();
  });

  it('maps r and R to square root', () => {
    fireKey('r');
    fireKey('R');
    expect(actions.pressSqrt).toHaveBeenCalledTimes(2);
  });

  it('ignores unrelated keys', () => {
    fireKey('F5');
    expect(Object.values(actions).some((fn) => fn.mock.calls.length > 0)).toBe(false);
  });
});
