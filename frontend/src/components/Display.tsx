import React from 'react';

interface DisplayProps {
  value: string;
  formula: string;
  isEvaluating?: boolean;
}

export const Display: React.FC<DisplayProps> = ({ value, formula, isEvaluating }) => {
  return (
    <div
      className="w-full bg-slate-950 border border-slate-800 rounded-xl p-5 mb-4 shadow-inner flex flex-col justify-end text-right min-h-[110px]"
      aria-label="Calculator display"
    >
      {/* Sub-display for expression history and active operator */}
      <div
        className="text-slate-400 text-sm font-mono h-6 overflow-hidden text-ellipsis whitespace-nowrap"
        aria-hidden="true"
      >
        {formula}
      </div>

      {/* Main Display value */}
      <div
        data-testid="calculator-display"
        aria-live="polite"
        className={`text-3xl sm:text-4xl font-mono font-semibold tracking-tight text-white overflow-x-auto whitespace-nowrap scrollbar-none transition-opacity duration-150 ${
          isEvaluating ? 'opacity-50' : 'opacity-100'
        }`}
      >
        {value}
      </div>
    </div>
  );
};
