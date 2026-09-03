import { describe, it, expect } from 'vitest';
import {
  calculatorReducer,
  formatDisplayNumber,
} from '../hooks/useCalculator';
import { CalculatorState } from '../types/calculator';

const getInitialState = (): CalculatorState => ({
  currentInput: '0',
  previousInput: null,
  operation: null,
  displayValue: '0',
  formula: '',
  isEvaluating: false,
  error: null,
  lastEvaluatedOperandB: null,
  lastEvaluatedOperation: null,
});

describe('useCalculator Reducer & State Machine Unit Tests', () => {
  it('rendersInitialDisplayWithZero', () => {
    const state = getInitialState();
    expect(state.displayValue).toBe('0');
    expect(state.currentInput).toBe('0');
  });

  it('preventsMultipleConsecutiveDecimals: pressing . multiple times retains single decimal', () => {
    let state = getInitialState();

    state = calculatorReducer(state, { type: 'INPUT_DIGIT', digit: '5' });
    state = calculatorReducer(state, { type: 'INPUT_DECIMAL' });
    state = calculatorReducer(state, { type: 'INPUT_DECIMAL' });
    state = calculatorReducer(state, { type: 'INPUT_DIGIT', digit: '3' });

    expect(state.displayValue).toBe('5.3');
    expect(state.currentInput).toBe('5.3');
  });

  it('handlesLeadingZerosCorrectly: multiple zeros do not produce numbers like 0004', () => {
    let state = getInitialState();

    state = calculatorReducer(state, { type: 'INPUT_DIGIT', digit: '0' });
    state = calculatorReducer(state, { type: 'INPUT_DIGIT', digit: '0' });
    state = calculatorReducer(state, { type: 'INPUT_DIGIT', digit: '0' });
    state = calculatorReducer(state, { type: 'INPUT_DIGIT', digit: '4' });

    expect(state.displayValue).toBe('4');
    expect(state.currentInput).toBe('4');

    // Also verify leading zero before decimal
    let decimalState = getInitialState();
    decimalState = calculatorReducer(decimalState, { type: 'INPUT_DIGIT', digit: '0' });
    decimalState = calculatorReducer(decimalState, { type: 'INPUT_DECIMAL' });
    decimalState = calculatorReducer(decimalState, { type: 'INPUT_DIGIT', digit: '4' });
    expect(decimalState.displayValue).toBe('0.4');
  });

  it('preventsMultipleOperators: entering 5 + * replaces operator with *', () => {
    let state = getInitialState();

    state = calculatorReducer(state, { type: 'INPUT_DIGIT', digit: '5' });
    state = calculatorReducer(state, { type: 'CHOOSE_OPERATION', operation: 'add' });
    expect(state.operation).toBe('add');
    expect(state.formula).toBe('5 +');

    // Replace operator with multiply
    state = calculatorReducer(state, { type: 'CHOOSE_OPERATION', operation: 'multiply' });
    expect(state.operation).toBe('multiply');
    expect(state.formula).toBe('5 ×');
  });

  it('clearsStateOnAllClear_AC: resets current input, memory, and active operator', () => {
    let state = getInitialState();

    state = calculatorReducer(state, { type: 'INPUT_DIGIT', digit: '9' });
    state = calculatorReducer(state, { type: 'CHOOSE_OPERATION', operation: 'add' });
    state = calculatorReducer(state, { type: 'INPUT_DIGIT', digit: '4' });
    state = calculatorReducer(state, { type: 'CLEAR_ALL' });

    expect(state.displayValue).toBe('0');
    expect(state.currentInput).toBe('0');
    expect(state.previousInput).toBeNull();
    expect(state.operation).toBeNull();
    expect(state.formula).toBe('');
  });

  it('clearsCurrentEntryOnClearEntry_CE: resets only current buffer without clearing pending calculations', () => {
    let state = getInitialState();

    state = calculatorReducer(state, { type: 'INPUT_DIGIT', digit: '1' });
    state = calculatorReducer(state, { type: 'INPUT_DIGIT', digit: '5' });
    state = calculatorReducer(state, { type: 'CHOOSE_OPERATION', operation: 'add' });
    state = calculatorReducer(state, { type: 'INPUT_DIGIT', digit: '9' });
    state = calculatorReducer(state, { type: 'INPUT_DIGIT', digit: '9' });

    // Clear only 99
    state = calculatorReducer(state, { type: 'CLEAR_ENTRY' });
    expect(state.displayValue).toBe('0');
    expect(state.currentInput).toBe('0');
    expect(state.previousInput).toBe('15');
    expect(state.operation).toBe('add');
  });

  it('disablesOperatorsDuringErrorState: operators are ignored while error is active', () => {
    let state = getInitialState();
    state = calculatorReducer(state, {
      type: 'CALCULATION_ERROR',
      error: 'Cannot divide by zero',
    });

    expect(state.error).toBe('Cannot divide by zero');

    // Attempt to choose an operator
    const nextState = calculatorReducer(state, {
      type: 'CHOOSE_OPERATION',
      operation: 'add',
    });
    expect(nextState.operation).toBeNull();
    expect(nextState.error).toBe('Cannot divide by zero');

    // Entering a new digit clears the error
    const recoveredState = calculatorReducer(state, {
      type: 'INPUT_DIGIT',
      digit: '7',
    });
    expect(recoveredState.error).toBeNull();
    expect(recoveredState.displayValue).toBe('7');
  });

  it('formatsLargeNumbersForDisplay: converts large or long fractional numbers gracefully', () => {
    expect(formatDisplayNumber(10000000000000)).toContain('e');
    expect(formatDisplayNumber('0.3333333333333333').length).toBeLessThanOrEqual(12);
  });
});
