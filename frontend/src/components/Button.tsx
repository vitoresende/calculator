import React from 'react';

type ButtonVariant = 'number' | 'operator' | 'action' | 'equal' | 'function';

interface ButtonProps {
  label: string;
  ariaLabel: string;
  onClick: () => void;
  variant?: ButtonVariant;
  className?: string;
  disabled?: boolean;
}

export const Button: React.FC<ButtonProps> = ({
  label,
  ariaLabel,
  onClick,
  variant = 'number',
  className = '',
  disabled = false,
}) => {
  const baseStyles =
    'relative inline-flex items-center justify-center font-mono font-medium rounded-xl text-lg sm:text-xl transition-all duration-150 select-none shadow-sm focus:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-900 active:scale-95 disabled:opacity-40 disabled:cursor-not-allowed disabled:active:scale-100 min-h-[56px]';

  const variantStyles: Record<ButtonVariant, string> = {
    number:
      'bg-slate-800 text-slate-100 hover:bg-slate-700 active:bg-slate-600 focus-visible:ring-slate-400 border border-slate-700/60',
    operator:
      'bg-amber-600/90 text-amber-50 hover:bg-amber-500 active:bg-amber-600 focus-visible:ring-amber-400 border border-amber-500/40',
    function:
      'bg-slate-700 text-slate-200 hover:bg-slate-600 active:bg-slate-500 focus-visible:ring-slate-300 border border-slate-600/50',
    action:
      'bg-rose-950/80 text-rose-200 hover:bg-rose-900 active:bg-rose-800 focus-visible:ring-rose-400 border border-rose-800/60',
    equal:
      'bg-blue-600 text-white hover:bg-blue-500 active:bg-blue-600 focus-visible:ring-blue-400 border border-blue-500 shadow-blue-900/30',
  };

  return (
    <button
      type="button"
      aria-label={ariaLabel}
      onClick={onClick}
      disabled={disabled}
      className={`${baseStyles} ${variantStyles[variant]} ${className}`}
    >
      {label}
    </button>
  );
};
