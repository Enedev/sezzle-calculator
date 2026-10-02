import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { ApiError } from '../api/client';
import { Calculator } from './Calculator';

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

describe('Calculator', () => {
  it('computes 2 + 3 = 5 via keypad clicks', async () => {
    mockCalculate.mockResolvedValue(5);
    const user = userEvent.setup();
    render(<Calculator />);

    await user.click(screen.getByRole('button', { name: '2' }));
    await user.click(screen.getByRole('button', { name: 'Add' }));
    await user.click(screen.getByRole('button', { name: '3' }));
    await user.click(screen.getByRole('button', { name: 'Equals' }));

    expect(mockCalculate).toHaveBeenCalledWith('add', [2, 3]);
    expect(await screen.findByText('5', { selector: '.display-value' })).toBeInTheDocument();
  });

  it('supports physical keyboard input end to end', async () => {
    mockCalculate.mockResolvedValue(10);
    const user = userEvent.setup();
    render(<Calculator />);

    await user.keyboard('5+5{Enter}');

    expect(mockCalculate).toHaveBeenCalledWith('add', [5, 5]);
    expect(await screen.findByText('10')).toBeInTheDocument();
  });

  it('shows a friendly message on a backend error and recovers on the next keypress', async () => {
    mockCalculate.mockRejectedValue(new ApiError('DIVISION_BY_ZERO', 'division by zero'));
    const user = userEvent.setup();
    render(<Calculator />);

    await user.click(screen.getByRole('button', { name: '5' }));
    await user.click(screen.getByRole('button', { name: 'Divide' }));
    await user.click(screen.getByRole('button', { name: '0' }));
    await user.click(screen.getByRole('button', { name: 'Equals' }));

    expect(await screen.findByRole('alert')).toHaveTextContent("Can't divide by zero");

    await user.click(screen.getByRole('button', { name: '3' }));
    expect(screen.queryByRole('alert')).not.toBeInTheDocument();
    expect(screen.getByText('3', { selector: '.display-value' })).toBeInTheDocument();
  });

  it('disables the keypad while a request is in flight', async () => {
    let resolvePromise!: (value: number) => void;
    mockCalculate.mockReturnValue(
      new Promise((resolve) => {
        resolvePromise = resolve;
      }),
    );
    const user = userEvent.setup();
    render(<Calculator />);

    await user.click(screen.getByRole('button', { name: '2' }));
    await user.click(screen.getByRole('button', { name: 'Add' }));
    await user.click(screen.getByRole('button', { name: '3' }));
    await user.click(screen.getByRole('button', { name: 'Equals' }));

    expect(screen.getByRole('button', { name: '2' })).toBeDisabled();

    resolvePromise(5);
    expect(await screen.findByText('5', { selector: '.display-value' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: '2' })).not.toBeDisabled();
  });
});
