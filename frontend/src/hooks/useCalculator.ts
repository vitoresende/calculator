import { useReducer, useCallback } from 'react';
import {
  CalculatorState,
  CalculatorAction,
  OperationType,
} from '../types/calculator';
import { calculate, CalculatorApiError } from '../services/api';

const MAX_DISPLAY_DIGITS = 12;

export function formatDisplayNumber(value: number | string): string {
  if (typeof value === 'string') {
    if (value === 'Error' || value.includes('Error') || value === 'Cannot divide by zero') {
      return value;
    }
    const num = Number(value);
    if (isNaN(num)) return value;
    if (value.endsWith('.')) return value;
    if (value.includes('.') && value.split('.')[1].endsWith('0')) return value;
    value = num;
  }

  if (Math.abs(value) >= 1e12 || (Math.abs(value) > 0 && Math.abs(value) < 1e-6)) {
    return value.toExponential(6).replace(/\.?0+e/, 'e');
  }

  const str = value.toString();
  if (str.length > MAX_DISPLAY_DIGITS) {
    return Number(value.toPrecision(8)).toString();
  }

  return str;
}

const initialState: CalculatorState = {
  currentInput: '0',
  previousInput: null,
  operation: null,
  displayValue: '0',
  formula: '',
  isEvaluating: false,
  error: null,
  lastEvaluatedOperandB: null,
  lastEvaluatedOperation: null,
};

export function calculatorReducer(
  state: CalculatorState,
  action: CalculatorAction
): CalculatorState {
  switch (action.type) {
    case 'INPUT_DIGIT': {
      // If recovering from error, reset
      if (state.error) {
        return {
          ...initialState,
          currentInput: action.digit,
          displayValue: action.digit,
        };
      }

      // Handle leading zeros: if current is "0", replace with digit unless digit is "0"
      let nextInput: string;
      if (state.currentInput === '0') {
        nextInput = action.digit;
      } else {
        if (state.currentInput.replace('-', '').length >= MAX_DISPLAY_DIGITS) {
          return state;
        }
        nextInput = state.currentInput + action.digit;
      }

      return {
        ...state,
        currentInput: nextInput,
        displayValue: formatDisplayNumber(nextInput),
      };
    }

    case 'INPUT_DECIMAL': {
      if (state.error) {
        return {
          ...initialState,
          currentInput: '0.',
          displayValue: '0.',
        };
      }

      // Prevent multiple consecutive decimals
      if (state.currentInput.includes('.')) {
        return state;
      }

      const nextInput = state.currentInput + '.';
      return {
        ...state,
        currentInput: nextInput,
        displayValue: nextInput,
      };
    }

    case 'CHOOSE_OPERATION': {
      if (state.error) {
        // Operators disabled during error state
        return state;
      }

      const opSymbol = getOperatorSymbol(action.operation);

      // If operator clicked without entering a new number, replace operator
      if (state.currentInput === '0' && state.previousInput !== null && state.operation !== null) {
        return {
          ...state,
          operation: action.operation,
          formula: `${state.previousInput} ${opSymbol}`,
        };
      }

      return {
        ...state,
        previousInput: state.currentInput,
        currentInput: '0',
        operation: action.operation,
        formula: `${state.currentInput} ${opSymbol}`,
        lastEvaluatedOperandB: null,
        lastEvaluatedOperation: null,
      };
    }

    case 'CALCULATION_START': {
      return {
        ...state,
        isEvaluating: true,
        error: null,
      };
    }

    case 'CALCULATION_SUCCESS': {
      const formattedResult = formatDisplayNumber(action.result);
      return {
        ...state,
        isEvaluating: false,
        displayValue: formattedResult,
        currentInput: action.isUnary ? formattedResult : '0',
        previousInput: action.isUnary ? state.previousInput : formattedResult,
        formula: action.formula,
        error: null,
      };
    }

    case 'CALCULATION_ERROR': {
      return {
        ...state,
        isEvaluating: false,
        error: action.error,
        displayValue: action.error,
        currentInput: '0',
        previousInput: null,
        operation: null,
      };
    }

    case 'CLEAR_ALL': {
      return initialState;
    }

    case 'CLEAR_ENTRY': {
      if (state.error) {
        return initialState;
      }
      return {
        ...state,
        currentInput: '0',
        displayValue: '0',
      };
    }

    case 'SET_ERROR': {
      return {
        ...state,
        isEvaluating: false,
        error: action.error,
        displayValue: action.error,
      };
    }

    default:
      return state;
  }
}

export function getOperatorSymbol(op: OperationType): string {
  switch (op) {
    case 'add':
      return '+';
    case 'subtract':
      return '−';
    case 'multiply':
      return '×';
    case 'divide':
      return '÷';
    case 'pow':
      return '^';
    case 'sqrt':
      return '√';
    case 'percentage':
      return '%';
    default:
      return '';
  }
}

export function useCalculator() {
  const [state, dispatch] = useReducer(calculatorReducer, initialState);

  const inputDigit = useCallback((digit: string) => {
    dispatch({ type: 'INPUT_DIGIT', digit });
  }, []);

  const inputDecimal = useCallback(() => {
    dispatch({ type: 'INPUT_DECIMAL' });
  }, []);

  const clearAll = useCallback(() => {
    dispatch({ type: 'CLEAR_ALL' });
  }, []);

  const clearEntry = useCallback(() => {
    dispatch({ type: 'CLEAR_ENTRY' });
  }, []);

  const chooseOperation = useCallback(
    async (op: OperationType) => {
      if (state.error) return;

      // Unary operation: sqrt or percentage
      if (op === 'sqrt' || op === 'percentage') {
        const val = Number(state.currentInput !== '0' ? state.currentInput : state.previousInput || '0');
        dispatch({ type: 'CALCULATION_START' });

        try {
          const resp = await calculate({ operation: op, a: val });
          const symbol = getOperatorSymbol(op);
          const formulaStr = op === 'sqrt' ? `${symbol}(${val})` : `${val}${symbol}`;
          dispatch({
            type: 'CALCULATION_SUCCESS',
            result: resp.result,
            formula: formulaStr,
            isUnary: true,
          });
        } catch (err: unknown) {
          const errorMsg =
            err instanceof CalculatorApiError ? err.message : 'Calculation error';
          dispatch({ type: 'CALCULATION_ERROR', error: errorMsg });
        }
        return;
      }

      // If there is already a pending binary calculation and user enters a new number, chain compute
      if (state.previousInput !== null && state.operation !== null && state.currentInput !== '0') {
        dispatch({ type: 'CALCULATION_START' });
        try {
          const aVal = Number(state.previousInput);
          const bVal = Number(state.currentInput);
          const resp = await calculate({ operation: state.operation, a: aVal, b: bVal });
          
          dispatch({
            type: 'CALCULATION_SUCCESS',
            result: resp.result,
            formula: `${formatDisplayNumber(resp.result)} ${getOperatorSymbol(op)}`,
          });
          dispatch({ type: 'CHOOSE_OPERATION', operation: op });
        } catch (err: unknown) {
          const errorMsg =
            err instanceof CalculatorApiError ? err.message : 'Calculation error';
          dispatch({ type: 'CALCULATION_ERROR', error: errorMsg });
        }
        return;
      }

      dispatch({ type: 'CHOOSE_OPERATION', operation: op });
    },
    [state]
  );

  const evaluate = useCallback(async () => {
    if (state.error) return;

    let op = state.operation;
    let aVal: number;
    let bVal: number;

    // Normal calculation: previousInput [op] currentInput
    if (state.previousInput !== null && op !== null) {
      aVal = Number(state.previousInput);
      bVal = Number(state.currentInput);
    } else if (state.lastEvaluatedOperation && state.lastEvaluatedOperandB !== null) {
      // Repeated equals: displayValue [lastOp] lastOperandB
      op = state.lastEvaluatedOperation;
      aVal = Number(state.displayValue);
      bVal = state.lastEvaluatedOperandB;
    } else {
      // Nothing to evaluate
      return;
    }

    dispatch({ type: 'CALCULATION_START' });

    try {
      const resp = await calculate({ operation: op, a: aVal, b: bVal });
      const symbol = getOperatorSymbol(op);
      const formulaStr = `${aVal} ${symbol} ${bVal} =`;

      dispatch({
        type: 'CALCULATION_SUCCESS',
        result: resp.result,
        formula: formulaStr,
      });

      // Update for repeated equals
      state.lastEvaluatedOperation = op;
      state.lastEvaluatedOperandB = bVal;
    } catch (err: unknown) {
      const errorMsg =
        err instanceof CalculatorApiError ? err.message : 'Calculation error';
      dispatch({ type: 'CALCULATION_ERROR', error: errorMsg });
    }
  }, [state]);

  return {
    state,
    inputDigit,
    inputDecimal,
    chooseOperation,
    evaluate,
    clearAll,
    clearEntry,
  };
}
