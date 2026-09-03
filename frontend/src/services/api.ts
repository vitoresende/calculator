import { CalculationRequest, CalculationResponse, ApiErrorResponse } from '../types/calculator';

const API_BASE_URL = (import.meta.env.VITE_API_BASE_URL || '/api/v1').replace(/\/$/, '');

export class CalculatorApiError extends Error {
  public readonly code: string;
  public readonly status: number;

  constructor(message: string, code: string, status: number) {
    super(message);
    this.name = 'CalculatorApiError';
    this.code = code;
    this.status = status;
  }
}

export async function calculate(
  request: CalculationRequest,
  signal?: AbortSignal
): Promise<CalculationResponse> {
  const response = await fetch(`${API_BASE_URL}/calculate`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(request),
    signal,
  });

  if (!response.ok) {
    let errorCode = 'UNKNOWN_ERROR';
    let errorMessage = `Server returned status ${response.status}`;

    try {
      const errorData = (await response.json()) as ApiErrorResponse;
      if (errorData?.error) {
        errorCode = errorData.error.code;
        errorMessage = errorData.error.message;
      }
    } catch {
      // JSON decode failed, use default response text or status message
    }

    // Map common error codes to friendly UI messages
    if (errorCode === 'DIVISION_BY_ZERO') {
      errorMessage = 'Cannot divide by zero';
    } else if (errorCode === 'NEGATIVE_SQUARE_ROOT') {
      errorMessage = 'Invalid input for square root';
    } else if (errorCode === 'INVALID_DOMAIN') {
      errorMessage = 'Domain error';
    } else if (errorCode === 'ARITHMETIC_OVERFLOW') {
      errorMessage = 'Overflow error';
    }

    throw new CalculatorApiError(errorMessage, errorCode, response.status);
  }

  return (await response.json()) as CalculationResponse;
}
