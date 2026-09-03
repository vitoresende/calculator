import React, { useEffect } from 'react';
import { Display } from './Display';
import { Keypad } from './Keypad';
import { ErrorBanner } from './ErrorBanner';
import { useCalculator } from '../hooks/useCalculator';

export const Calculator: React.FC = () => {
  const {
    state,
    inputDigit,
    inputDecimal,
    chooseOperation,
    evaluate,
    clearAll,
    clearEntry,
  } = useCalculator();

  // Handle physical keyboard inputs
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      // Ignore key events when focusing inputs/textareas if any
      if (['INPUT', 'TEXTAREA'].includes((e.target as HTMLElement)?.tagName)) {
        return;
      }

      if (/^[0-9]$/.test(e.key)) {
        e.preventDefault();
        inputDigit(e.key);
      } else if (e.key === '.') {
        e.preventDefault();
        inputDecimal();
      } else if (e.key === '+') {
        e.preventDefault();
        chooseOperation('add');
      } else if (e.key === '-') {
        e.preventDefault();
        chooseOperation('subtract');
      } else if (e.key === '*') {
        e.preventDefault();
        chooseOperation('multiply');
      } else if (e.key === '/') {
        e.preventDefault();
        chooseOperation('divide');
      } else if (e.key === '^') {
        e.preventDefault();
        chooseOperation('pow');
      } else if (e.key === '%') {
        e.preventDefault();
        chooseOperation('percentage');
      } else if (e.key === 'Enter' || e.key === '=') {
        e.preventDefault();
        evaluate();
      } else if (e.key === 'Escape') {
        e.preventDefault();
        clearAll();
      } else if (e.key === 'Backspace') {
        e.preventDefault();
        clearEntry();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [inputDigit, inputDecimal, chooseOperation, evaluate, clearAll, clearEntry]);

  return (
    <div
      className="w-full max-w-sm sm:max-w-md mx-auto bg-slate-900 border border-slate-800 rounded-2xl shadow-2xl p-4 sm:p-6"
      role="region"
      aria-label="Scientific Calculator"
    >
      <div className="flex items-center justify-between mb-3 px-1">
        <div className="flex items-center gap-2">
          <span className="w-3 h-3 rounded-full bg-rose-500 inline-block" />
          <span className="w-3 h-3 rounded-full bg-amber-500 inline-block" />
          <span className="w-3 h-3 rounded-full bg-emerald-500 inline-block" />
        </div>
        <span className="text-xs font-mono font-medium text-slate-400 uppercase tracking-wider">
          Go + React Engine
        </span>
      </div>

      <ErrorBanner message={state.error} onDismiss={clearAll} />

      <Display
        value={state.displayValue}
        formula={state.formula}
        isEvaluating={state.isEvaluating}
      />

      <Keypad
        onDigit={inputDigit}
        onDecimal={inputDecimal}
        onOperation={chooseOperation}
        onEvaluate={evaluate}
        onClearAll={clearAll}
        onClearEntry={clearEntry}
        isEvaluating={state.isEvaluating}
        hasError={state.error !== null}
      />
    </div>
  );
};
