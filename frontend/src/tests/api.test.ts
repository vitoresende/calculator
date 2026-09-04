import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { calculate, CalculatorApiError } from '../services/api';

describe('API Client (calculate & CalculatorApiError)', () => {
  const originalFetch = global.fetch;

  beforeEach(() => {
    vi.restoreAllMocks();
  });

  afterEach(() => {
    global.fetch = originalFetch;
  });

  it('executes successful calculation request', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({
        result: 42,
        operation: 'add',
        operands: [20, 22],
      }),
    } as unknown as Response);

    const res = await calculate({ operation: 'add', a: 20, b: 22 });
    expect(res.result).toBe(42);
    expect(res.operation).toBe('add');
    expect(res.operands).toEqual([20, 22]);
  });

  it('maps DIVISION_BY_ZERO code to friendly user message', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      ok: false,
      status: 422,
      json: async () => ({
        error: { code: 'DIVISION_BY_ZERO', message: 'division by zero is not allowed' },
      }),
    } as unknown as Response);

    await expect(calculate({ operation: 'divide', a: 10, b: 0 })).rejects.toThrow(
      'Cannot divide by zero'
    );

    try {
      await calculate({ operation: 'divide', a: 10, b: 0 });
    } catch (err) {
      expect(err).toBeInstanceOf(CalculatorApiError);
      if (err instanceof CalculatorApiError) {
        expect(err.code).toBe('DIVISION_BY_ZERO');
        expect(err.status).toBe(422);
      }
    }
  });

  it('maps NEGATIVE_SQUARE_ROOT code to friendly user message', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      ok: false,
      status: 422,
      json: async () => ({
        error: { code: 'NEGATIVE_SQUARE_ROOT', message: 'square root of negative number' },
      }),
    } as unknown as Response);

    await expect(calculate({ operation: 'sqrt', a: -16 })).rejects.toThrow(
      'Invalid input for square root'
    );
  });

  it('maps INVALID_DOMAIN code to friendly user message', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      ok: false,
      status: 422,
      json: async () => ({
        error: { code: 'INVALID_DOMAIN', message: 'undefined domain' },
      }),
    } as unknown as Response);

    await expect(calculate({ operation: 'pow', a: -4, b: 0.5 })).rejects.toThrow('Domain error');
  });

  it('maps ARITHMETIC_OVERFLOW code to friendly user message', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      ok: false,
      status: 422,
      json: async () => ({
        error: { code: 'ARITHMETIC_OVERFLOW', message: 'overflow' },
      }),
    } as unknown as Response);

    await expect(calculate({ operation: 'pow', a: 1e200, b: 2 })).rejects.toThrow('Overflow error');
  });

  it('handles non-JSON error response from backend gracefully', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      ok: false,
      status: 502,
      json: async () => {
        throw new Error('Bad Gateway');
      },
    } as unknown as Response);

    await expect(calculate({ operation: 'add', a: 1, b: 2 })).rejects.toThrow(
      'Server returned status 502'
    );
  });
});
