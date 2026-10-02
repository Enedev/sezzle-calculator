import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import { Keypad } from './Keypad';

function renderKeypad(disabled = false) {
  const handlers = {
    onDigit: vi.fn(),
    onDecimal: vi.fn(),
    onOperator: vi.fn(),
    onEquals: vi.fn(),
    onSqrt: vi.fn(),
    onClear: vi.fn(),
    onBackspace: vi.fn(),
  };
  render(
    <Keypad
      onDigit={handlers.onDigit}
      onDecimal={handlers.onDecimal}
      onOperator={handlers.onOperator}
      onEquals={handlers.onEquals}
      onSqrt={handlers.onSqrt}
      onClear={handlers.onClear}
      onBackspace={handlers.onBackspace}
      disabled={disabled}
    />,
  );
  return handlers;
}

describe('Keypad', () => {
  it('invokes onDigit for every digit button', async () => {
    const handlers = renderKeypad();
    const user = userEvent.setup();
    for (const digit of ['0', '1', '2', '3', '4', '5', '6', '7', '8', '9']) {
      await user.click(screen.getByRole('button', { name: digit }));
    }
    expect(handlers.onDigit.mock.calls.map((c) => c[0])).toEqual([
      '0',
      '1',
      '2',
      '3',
      '4',
      '5',
      '6',
      '7',
      '8',
      '9',
    ]);
  });

  it('invokes onOperator with the right operator for each operator button', async () => {
    const handlers = renderKeypad();
    const user = userEvent.setup();
    const cases: [string, string][] = [
      ['Add', 'add'],
      ['Subtract', 'subtract'],
      ['Multiply', 'multiply'],
      ['Divide', 'divide'],
      ['Power', 'power'],
      ['Percentage', 'percentage'],
    ];
    for (const [label, operator] of cases) {
      await user.click(screen.getByRole('button', { name: label }));
      expect(handlers.onOperator).toHaveBeenLastCalledWith(operator);
    }
  });

  it('invokes onDecimal, onSqrt, onClear, onBackspace, and onEquals', async () => {
    const handlers = renderKeypad();
    const user = userEvent.setup();
    await user.click(screen.getByRole('button', { name: 'Decimal point' }));
    await user.click(screen.getByRole('button', { name: 'Square root' }));
    await user.click(screen.getByRole('button', { name: 'Clear' }));
    await user.click(screen.getByRole('button', { name: 'Backspace' }));
    await user.click(screen.getByRole('button', { name: 'Equals' }));
    expect(handlers.onDecimal).toHaveBeenCalledTimes(1);
    expect(handlers.onSqrt).toHaveBeenCalledTimes(1);
    expect(handlers.onClear).toHaveBeenCalledTimes(1);
    expect(handlers.onBackspace).toHaveBeenCalledTimes(1);
    expect(handlers.onEquals).toHaveBeenCalledTimes(1);
  });

  it('disables every button when disabled is true', () => {
    renderKeypad(true);
    for (const button of screen.getAllByRole('button')) {
      expect(button).toBeDisabled();
    }
  });
});
