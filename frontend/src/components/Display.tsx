interface DisplayProps {
  expression: string;
  value: string;
  errorMessage: string | null;
}

// The ledger tape: a dim running-expression line above a large right-aligned
// current value, matching an accounting ledger's column convention. An API
// error replaces the value with a rotated stamp mark instead of a toast.
export function Display({ expression, value, errorMessage }: DisplayProps) {
  return (
    <div className="display">
      <div className="display-expression" aria-hidden={errorMessage ? 'true' : undefined}>
        {expression || ' '}
      </div>
      <div
        className="display-value-row"
        aria-live={errorMessage ? 'assertive' : 'polite'}
        aria-atomic="true"
      >
        {errorMessage ? (
          <span className="error-stamp" role="alert">
            {errorMessage}
          </span>
        ) : (
          <span className="display-value">{value}</span>
        )}
      </div>
    </div>
  );
}
