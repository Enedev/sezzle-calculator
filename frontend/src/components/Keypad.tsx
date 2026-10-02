import type { BinaryOperator } from '../hooks/useCalculator';
import { Button } from './Button';

interface KeypadProps {
  onDigit: (digit: string) => void;
  onDecimal: () => void;
  onOperator: (operator: BinaryOperator) => void;
  onEquals: () => void;
  onSqrt: () => void;
  onClear: () => void;
  onBackspace: () => void;
  disabled: boolean;
}

export function Keypad({
  onDigit,
  onDecimal,
  onOperator,
  onEquals,
  onSqrt,
  onClear,
  onBackspace,
  disabled,
}: KeypadProps) {
  return (
    <div className="keypad">
      <Button label="C" onClick={onClear} variant="function" disabled={disabled} ariaLabel="Clear" />
      <Button
        label="⌫"
        onClick={onBackspace}
        variant="function"
        disabled={disabled}
        ariaLabel="Backspace"
      />
      <Button label="√" onClick={onSqrt} variant="function" disabled={disabled} ariaLabel="Square root" />
      <Button
        label="÷"
        onClick={() => onOperator('divide')}
        variant="operator"
        disabled={disabled}
        ariaLabel="Divide"
      />

      <Button label="7" onClick={() => onDigit('7')} disabled={disabled} />
      <Button label="8" onClick={() => onDigit('8')} disabled={disabled} />
      <Button label="9" onClick={() => onDigit('9')} disabled={disabled} />
      <Button
        label="×"
        onClick={() => onOperator('multiply')}
        variant="operator"
        disabled={disabled}
        ariaLabel="Multiply"
      />

      <Button label="4" onClick={() => onDigit('4')} disabled={disabled} />
      <Button label="5" onClick={() => onDigit('5')} disabled={disabled} />
      <Button label="6" onClick={() => onDigit('6')} disabled={disabled} />
      <Button
        label="−"
        onClick={() => onOperator('subtract')}
        variant="operator"
        disabled={disabled}
        ariaLabel="Subtract"
      />

      <Button label="1" onClick={() => onDigit('1')} disabled={disabled} />
      <Button label="2" onClick={() => onDigit('2')} disabled={disabled} />
      <Button label="3" onClick={() => onDigit('3')} disabled={disabled} />
      <Button
        label="+"
        onClick={() => onOperator('add')}
        variant="operator"
        disabled={disabled}
        ariaLabel="Add"
      />

      <Button
        label="%"
        onClick={() => onOperator('percentage')}
        variant="operator"
        disabled={disabled}
        ariaLabel="Percentage"
      />
      <Button label="0" onClick={() => onDigit('0')} disabled={disabled} />
      <Button label="." onClick={onDecimal} disabled={disabled} ariaLabel="Decimal point" />
      <Button
        label="^"
        onClick={() => onOperator('power')}
        variant="operator"
        disabled={disabled}
        ariaLabel="Power"
      />

      <Button label="=" onClick={onEquals} variant="equals" disabled={disabled} ariaLabel="Equals" />
    </div>
  );
}
