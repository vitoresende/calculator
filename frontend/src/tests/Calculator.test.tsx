import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { Calculator } from '../components/Calculator';
import * as api from '../services/api';

vi.mock('../services/api', () => ({
  calculate: vi.fn(),
  CalculatorApiError: class extends Error {
    public code: string;
    public status: number;
    constructor(message: string, code: string, status: number) {
      super(message);
      this.code = code;
      this.status = status;
    }
  },
}));

describe('Calculator Component Integration & User Interactions', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('rendersInitialDisplayWithZero and has accessible display', () => {
    render(<Calculator />);
    const display = screen.getByTestId('calculator-display');
    expect(display).toHaveTextContent('0');
    expect(display).toHaveAttribute('aria-live', 'polite');
  });

  it('handles user digit clicks and decimal entry', async () => {
    const user = userEvent.setup();
    render(<Calculator />);

    await user.click(screen.getByRole('button', { name: /^5$/ }));
    await user.click(screen.getByRole('button', { name: /decimal point/i }));
    await user.click(screen.getByRole('button', { name: /^2$/ }));

    const display = screen.getByTestId('calculator-display');
    expect(display).toHaveTextContent('5.2');
  });

  it('preventsMultipleConsecutiveDecimals', async () => {
    const user = userEvent.setup();
    render(<Calculator />);

    await user.click(screen.getByRole('button', { name: /^5$/ }));
    await user.click(screen.getByRole('button', { name: /decimal point/i }));
    await user.click(screen.getByRole('button', { name: /decimal point/i }));
    await user.click(screen.getByRole('button', { name: /^3$/ }));

    const display = screen.getByTestId('calculator-display');
    expect(display).toHaveTextContent('5.3');
  });

  it('executes addition calculation successfully', async () => {
    const user = userEvent.setup();
    vi.mocked(api.calculate).mockResolvedValueOnce({
      result: 7,
      operation: 'add',
      operands: [3, 4],
    });

    render(<Calculator />);

    await user.click(screen.getByRole('button', { name: /^3$/ }));
    await user.click(screen.getByRole('button', { name: /^Add$/i }));
    await user.click(screen.getByRole('button', { name: /^4$/ }));
    await user.click(screen.getByRole('button', { name: /^Calculate$/i }));

    await waitFor(() => {
      expect(api.calculate).toHaveBeenCalledWith({
        operation: 'add',
        a: 3,
        b: 4,
      });
      const display = screen.getByTestId('calculator-display');
      expect(display).toHaveTextContent('7');
    });
  });

  it('displaysFriendlyDivisionByZeroMessage when backend returns division by zero error', async () => {
    const user = userEvent.setup();
    const errorInstance = new api.CalculatorApiError(
      'Cannot divide by zero',
      'DIVISION_BY_ZERO',
      422
    );
    vi.mocked(api.calculate).mockRejectedValueOnce(errorInstance);

    render(<Calculator />);

    await user.click(screen.getByRole('button', { name: /^9$/ }));
    await user.click(screen.getByRole('button', { name: /^Divide$/i }));
    await user.click(screen.getByRole('button', { name: /^0$/ }));
    await user.click(screen.getByRole('button', { name: /^Calculate$/i }));

    await waitFor(() => {
      const display = screen.getByTestId('calculator-display');
      expect(display).toHaveTextContent('0');
      expect(screen.getByRole('alert')).toHaveTextContent('Cannot divide by zero');
    });
  });

  it('clears state on AC button press', async () => {
    const user = userEvent.setup();
    render(<Calculator />);

    await user.click(screen.getByRole('button', { name: /^8$/ }));
    await user.click(screen.getByRole('button', { name: /^9$/ }));
    expect(screen.getByTestId('calculator-display')).toHaveTextContent('89');

    await user.click(screen.getByRole('button', { name: /All Clear/i }));
    expect(screen.getByTestId('calculator-display')).toHaveTextContent('0');
  });

  it('supports physical keyboard input', async () => {
    const user = userEvent.setup();
    render(<Calculator />);

    await user.keyboard('4');
    await user.keyboard('.');
    await user.keyboard('7');

    expect(screen.getByTestId('calculator-display')).toHaveTextContent('4.7');

    await user.keyboard('{Escape}');
    expect(screen.getByTestId('calculator-display')).toHaveTextContent('0');
  });

  it('supports all operator and editing keyboard shortcuts', async () => {
    const user = userEvent.setup();
    render(<Calculator />);

    // Test backspace (clear entry)
    await user.keyboard('12');
    expect(screen.getByTestId('calculator-display')).toHaveTextContent('12');
    await user.keyboard('{Backspace}');
    expect(screen.getByTestId('calculator-display')).toHaveTextContent('0');

    // Test operator keys: +, -, *, /, ^, %
    await user.keyboard('5');
    await user.keyboard('+');
    expect(screen.getByText('5 +')).toBeInTheDocument();

    await user.keyboard('-');
    expect(screen.getByText('5 −')).toBeInTheDocument();

    await user.keyboard('*');
    expect(screen.getByText('5 ×')).toBeInTheDocument();

    await user.keyboard('/');
    expect(screen.getByText('5 ÷')).toBeInTheDocument();

    await user.keyboard('^');
    expect(screen.getByText('5 ^')).toBeInTheDocument();

    // Type 2 and calculate with Enter
    vi.mocked(api.calculate).mockResolvedValueOnce({
      result: 25,
      operation: 'pow',
      operands: [5, 2],
    });

    await user.keyboard('2');
    await user.keyboard('{Enter}');

    await waitFor(() => {
      expect(screen.getByTestId('calculator-display')).toHaveTextContent('25');
    });
  });

  it('executes unary square root operation on button click', async () => {
    const user = userEvent.setup();
    vi.mocked(api.calculate).mockResolvedValueOnce({
      result: 9,
      operation: 'sqrt',
      operands: [81],
    });

    render(<Calculator />);

    await user.click(screen.getByRole('button', { name: /^8$/ }));
    await user.click(screen.getByRole('button', { name: /^1$/ }));
    await user.click(screen.getByRole('button', { name: /Square Root/i }));

    await waitFor(() => {
      expect(api.calculate).toHaveBeenCalledWith({
        operation: 'sqrt',
        a: 81,
      });
      expect(screen.getByTestId('calculator-display')).toHaveTextContent('9');
    });
  });

  it('executes direct percentage operation on button click', async () => {
    const user = userEvent.setup();
    vi.mocked(api.calculate).mockResolvedValueOnce({
      result: 0.5,
      operation: 'percentage',
      operands: [50],
    });

    render(<Calculator />);

    await user.click(screen.getByRole('button', { name: /^5$/ }));
    await user.click(screen.getByRole('button', { name: /^0$/ }));
    await user.click(screen.getByRole('button', { name: /Percentage/i }));

    await waitFor(() => {
      expect(api.calculate).toHaveBeenCalledWith({
        operation: 'percentage',
        a: 50,
      });
      expect(screen.getByTestId('calculator-display')).toHaveTextContent('0.5');
    });
  });
});
