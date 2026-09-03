export type OperationType =
  | 'add'
  | 'subtract'
  | 'multiply'
  | 'divide'
  | 'pow'
  | 'sqrt'
  | 'percentage';

export interface CalculationRequest {
  operation: OperationType;
  a: number;
  b?: number;
}

export interface CalculationResponse {
  result: number;
  operation: string;
  operands: number[];
}

export interface ApiErrorDetail {
  code: string;
  message: string;
}

export interface ApiErrorResponse {
  error: ApiErrorDetail;
}

export interface CalculatorState {
  currentInput: string;
  previousInput: string | null;
  operation: OperationType | null;
  displayValue: string;
  formula: string;
  isEvaluating: boolean;
  error: string | null;
  lastEvaluatedOperandB: number | null;
  lastEvaluatedOperation: OperationType | null;
}

export type CalculatorAction =
  | { type: 'INPUT_DIGIT'; digit: string }
  | { type: 'INPUT_DECIMAL' }
  | { type: 'CHOOSE_OPERATION'; operation: OperationType }
  | { type: 'CALCULATION_START' }
  | { type: 'CALCULATION_SUCCESS'; result: number; formula: string; isUnary?: boolean }
  | { type: 'CALCULATION_ERROR'; error: string }
  | { type: 'CLEAR_ALL' }
  | { type: 'CLEAR_ENTRY' }
  | { type: 'SET_ERROR'; error: string };
