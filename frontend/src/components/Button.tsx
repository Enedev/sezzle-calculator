interface ButtonProps {
  label: string;
  onClick: () => void;
  variant?: 'digit' | 'operator' | 'equals' | 'function';
  disabled?: boolean;
  ariaLabel?: string;
}

export function Button({ label, onClick, variant = 'digit', disabled = false, ariaLabel }: ButtonProps) {
  return (
    <button
      type="button"
      className={`btn btn-${variant}`}
      onClick={onClick}
      disabled={disabled}
      aria-label={ariaLabel ?? label}
    >
      {label}
    </button>
  );
}
