import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { Display } from './Display';

describe('Display', () => {
  it('shows the current value and expression', () => {
    render(<Display expression="9 +" value="16" errorMessage={null} />);
    expect(screen.getByText('16')).toBeInTheDocument();
    expect(screen.getByText('9 +')).toBeInTheDocument();
  });

  it('uses a polite live region when there is no error', () => {
    render(<Display expression="" value="0" errorMessage={null} />);
    const row = screen.getByText('0').closest('[aria-live]');
    expect(row).toHaveAttribute('aria-live', 'polite');
  });

  it('shows the error stamp with an assertive live region instead of the value', () => {
    render(<Display expression="" value="0" errorMessage="Can't divide by zero" />);
    expect(screen.getByRole('alert')).toHaveTextContent("Can't divide by zero");
    const row = screen.getByRole('alert').closest('[aria-live]');
    expect(row).toHaveAttribute('aria-live', 'assertive');
    expect(screen.queryByText('0')).not.toBeInTheDocument();
  });
});
