import { useCalculator } from '../hooks/useCalculator';
import { useKeyboardInput } from '../hooks/useKeyboardInput';
import { Display } from './Display';
import { Keypad } from './Keypad';

export function Calculator() {
  const calc = useCalculator();
  useKeyboardInput(calc);

  return (
    <div className="calculator" data-loading={calc.isLoading}>
      <Display expression={calc.expression} value={calc.display} errorMessage={calc.errorMessage} />
      <Keypad
        onDigit={calc.pressDigit}
        onDecimal={calc.pressDecimal}
        onOperator={calc.pressOperator}
        onEquals={calc.pressEquals}
        onSqrt={calc.pressSqrt}
        onClear={calc.pressClear}
        onBackspace={calc.pressBackspace}
        disabled={calc.isLoading}
      />
    </div>
  );
}
