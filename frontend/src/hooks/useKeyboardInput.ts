import { useEffect } from 'react';
import type { BinaryOperator } from './useCalculator';

interface KeyboardActions {
  pressDigit: (digit: string) => void;
  pressDecimal: () => void;
  pressOperator: (operator: BinaryOperator) => void | Promise<void>;
  pressEquals: () => void | Promise<void>;
  pressSqrt: () => void | Promise<void>;
  pressClear: () => void;
  pressBackspace: () => void;
}

const KEY_TO_OPERATOR: Record<string, BinaryOperator> = {
  '+': 'add',
  '-': 'subtract',
  '*': 'multiply',
  '/': 'divide',
  '^': 'power',
  '%': 'percentage',
};

// Keyboard support for the calculator: digits, the four arithmetic symbols
// plus ^ and %, Enter/= for equals, Escape to clear, Backspace, and r for
// square root (there's no standard physical key for √).
export function useKeyboardInput(actions: KeyboardActions): void {
  useEffect(() => {
    function handleKeyDown(event: KeyboardEvent) {
      if (event.key >= '0' && event.key <= '9') {
        actions.pressDigit(event.key);
        return;
      }
      if (event.key === '.') {
        actions.pressDecimal();
        return;
      }
      if (event.key in KEY_TO_OPERATOR) {
        void actions.pressOperator(KEY_TO_OPERATOR[event.key]);
        return;
      }
      if (event.key === 'Enter' || event.key === '=') {
        event.preventDefault();
        void actions.pressEquals();
        return;
      }
      if (event.key === 'Escape') {
        actions.pressClear();
        return;
      }
      if (event.key === 'Backspace') {
        actions.pressBackspace();
        return;
      }
      if (event.key === 'r' || event.key === 'R') {
        void actions.pressSqrt();
      }
    }

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [actions]);
}
