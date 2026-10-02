import { useCallback, useReducer } from 'react';
import { calculate } from '../api/client';
import { friendlyMessage } from './friendlyMessage';

// Binary operators are regular chain operators under the immediate-execution
// model (classic pocket-calculator semantics: 2 + 3 x 4 = 20, not 14).
// percentage is deliberately included here, not treated as a special unary
// "/100" key, because the backend's percentage contract is fixed-binary
// ("a percent of b") — folding it into the same chaining logic keeps one
// code path for every operator instead of a second, inconsistent one.
export type BinaryOperator = 'add' | 'subtract' | 'multiply' | 'divide' | 'power' | 'percentage';

const OPERATOR_SYMBOLS: Record<BinaryOperator, string> = {
  add: '+',
  subtract: '−',
  multiply: '×',
  divide: '÷',
  power: '^',
  percentage: '%',
};

interface State {
  accumulator: number | null;
  pendingOperator: BinaryOperator | null;
  entry: string;
  status: 'idle' | 'loading' | 'error';
  errorMessage: string | null;
}

const initialState: State = {
  accumulator: null,
  pendingOperator: null,
  entry: '',
  status: 'idle',
  errorMessage: null,
};

type Action =
  | { type: 'DIGIT'; digit: string }
  | { type: 'DECIMAL' }
  | { type: 'BACKSPACE' }
  | { type: 'CLEAR' }
  | { type: 'QUEUE_OPERATOR'; operator: BinaryOperator }
  | { type: 'REQUEST_START' }
  | { type: 'BINARY_SUCCESS'; result: number; pendingOperator: BinaryOperator | null }
  | { type: 'UNARY_SUCCESS'; result: number; replaceEntry: boolean }
  | { type: 'FAILURE'; message: string };

function reducer(state: State, action: Action): State {
  switch (action.type) {
    case 'DIGIT': {
      if (state.status === 'loading') return state;
      const entry = state.status === 'error' ? action.digit : state.entry + action.digit;
      return { ...state, entry, status: 'idle', errorMessage: null };
    }
    case 'DECIMAL': {
      if (state.status === 'loading') return state;
      const base = state.status === 'error' ? '' : state.entry;
      if (base.includes('.')) return { ...state, status: 'idle', errorMessage: null };
      return { ...state, entry: base === '' ? '0.' : base + '.', status: 'idle', errorMessage: null };
    }
    case 'BACKSPACE': {
      if (state.status === 'loading') return state;
      if (state.status === 'error') return { ...state, status: 'idle', errorMessage: null };
      return { ...state, entry: state.entry.slice(0, -1) };
    }
    case 'CLEAR':
      return initialState;
    case 'QUEUE_OPERATOR':
      return {
        ...state,
        accumulator: state.accumulator ?? parseFloat(state.entry || '0'),
        pendingOperator: action.operator,
        entry: '',
        status: 'idle',
        errorMessage: null,
      };
    case 'REQUEST_START':
      return { ...state, status: 'loading', errorMessage: null };
    case 'BINARY_SUCCESS':
      return {
        ...state,
        accumulator: action.result,
        pendingOperator: action.pendingOperator,
        entry: '',
        status: 'idle',
        errorMessage: null,
      };
    case 'UNARY_SUCCESS':
      return action.replaceEntry
        ? { ...state, entry: String(action.result), status: 'idle', errorMessage: null }
        : { ...state, accumulator: action.result, status: 'idle', errorMessage: null };
    case 'FAILURE':
      return {
        accumulator: null,
        pendingOperator: null,
        entry: '',
        status: 'error',
        errorMessage: action.message,
      };
    default:
      return state;
  }
}

export function useCalculator() {
  const [state, dispatch] = useReducer(reducer, initialState);

  const pressDigit = useCallback((digit: string) => dispatch({ type: 'DIGIT', digit }), []);
  const pressDecimal = useCallback(() => dispatch({ type: 'DECIMAL' }), []);
  const pressBackspace = useCallback(() => dispatch({ type: 'BACKSPACE' }), []);
  const pressClear = useCallback(() => dispatch({ type: 'CLEAR' }), []);

  const pressOperator = useCallback(
    async (operator: BinaryOperator) => {
      if (state.status === 'loading') return;

      if (state.pendingOperator && state.entry !== '') {
        const a = state.accumulator ?? 0;
        const b = parseFloat(state.entry);
        dispatch({ type: 'REQUEST_START' });
        try {
          const result = await calculate(state.pendingOperator, [a, b]);
          dispatch({ type: 'BINARY_SUCCESS', result, pendingOperator: operator });
        } catch (err) {
          dispatch({ type: 'FAILURE', message: friendlyMessage(err) });
        }
        return;
      }

      dispatch({ type: 'QUEUE_OPERATOR', operator });
    },
    [state.status, state.pendingOperator, state.entry, state.accumulator],
  );

  const pressEquals = useCallback(async () => {
    if (state.status === 'loading') return;
    if (!state.pendingOperator) return;

    const a = state.accumulator ?? 0;
    const b = parseFloat(state.entry || String(state.accumulator ?? 0));
    dispatch({ type: 'REQUEST_START' });
    try {
      const result = await calculate(state.pendingOperator, [a, b]);
      dispatch({ type: 'BINARY_SUCCESS', result, pendingOperator: null });
    } catch (err) {
      dispatch({ type: 'FAILURE', message: friendlyMessage(err) });
    }
  }, [state.status, state.pendingOperator, state.entry, state.accumulator]);

  const pressSqrt = useCallback(async () => {
    if (state.status === 'loading') return;

    const usingEntry = state.entry !== '';
    const value = usingEntry ? parseFloat(state.entry) : (state.accumulator ?? 0);
    dispatch({ type: 'REQUEST_START' });
    try {
      const result = await calculate('sqrt', [value]);
      dispatch({ type: 'UNARY_SUCCESS', result, replaceEntry: usingEntry });
    } catch (err) {
      dispatch({ type: 'FAILURE', message: friendlyMessage(err) });
    }
  }, [state.status, state.entry, state.accumulator]);

  const display = state.entry !== '' ? state.entry : String(state.accumulator ?? 0);
  const expression = state.pendingOperator
    ? `${state.accumulator ?? 0} ${OPERATOR_SYMBOLS[state.pendingOperator]}`
    : '';

  return {
    display,
    expression,
    isLoading: state.status === 'loading',
    errorMessage: state.status === 'error' ? state.errorMessage : null,
    pressDigit,
    pressDecimal,
    pressBackspace,
    pressClear,
    pressOperator,
    pressEquals,
    pressSqrt,
  };
}
