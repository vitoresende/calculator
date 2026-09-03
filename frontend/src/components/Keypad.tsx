import React from 'react';
import { Button } from './Button';
import { OperationType } from '../types/calculator';

interface KeypadProps {
  onDigit: (digit: string) => void;
  onDecimal: () => void;
  onOperation: (op: OperationType) => void;
  onEvaluate: () => void;
  onClearAll: () => void;
  onClearEntry: () => void;
  isEvaluating?: boolean;
  hasError?: boolean;
}

export const Keypad: React.FC<KeypadProps> = ({
  onDigit,
  onDecimal,
  onOperation,
  onEvaluate,
  onClearAll,
  onClearEntry,
  isEvaluating = false,
  hasError = false,
}) => {
  return (
    <div className="grid grid-cols-4 gap-2.5 sm:gap-3 w-full" role="group" aria-label="Calculator keypad">
      {/* Row 1: Memory/Clear and Advanced Operations */}
      <Button
        label="AC"
        ariaLabel="All Clear"
        onClick={onClearAll}
        variant="action"
      />
      <Button
        label="CE"
        ariaLabel="Clear Entry"
        onClick={onClearEntry}
        variant="action"
      />
      <Button
        label="√"
        ariaLabel="Square Root"
        onClick={() => onOperation('sqrt')}
        variant="function"
        disabled={hasError || isEvaluating}
      />
      <Button
        label="^"
        ariaLabel="Exponentiation"
        onClick={() => onOperation('pow')}
        variant="function"
        disabled={hasError || isEvaluating}
      />

      {/* Row 2: 7, 8, 9, Divide */}
      <Button
        label="7"
        ariaLabel="7"
        onClick={() => onDigit('7')}
        variant="number"
      />
      <Button
        label="8"
        ariaLabel="8"
        onClick={() => onDigit('8')}
        variant="number"
      />
      <Button
        label="9"
        ariaLabel="9"
        onClick={() => onDigit('9')}
        variant="number"
      />
      <Button
        label="÷"
        ariaLabel="Divide"
        onClick={() => onOperation('divide')}
        variant="operator"
        disabled={hasError || isEvaluating}
      />

      {/* Row 3: 4, 5, 6, Multiply */}
      <Button
        label="4"
        ariaLabel="4"
        onClick={() => onDigit('4')}
        variant="number"
      />
      <Button
        label="5"
        ariaLabel="5"
        onClick={() => onDigit('5')}
        variant="number"
      />
      <Button
        label="6"
        ariaLabel="6"
        onClick={() => onDigit('6')}
        variant="number"
      />
      <Button
        label="×"
        ariaLabel="Multiply"
        onClick={() => onOperation('multiply')}
        variant="operator"
        disabled={hasError || isEvaluating}
      />

      {/* Row 4: 1, 2, 3, Subtract */}
      <Button
        label="1"
        ariaLabel="1"
        onClick={() => onDigit('1')}
        variant="number"
      />
      <Button
        label="2"
        ariaLabel="2"
        onClick={() => onDigit('2')}
        variant="number"
      />
      <Button
        label="3"
        ariaLabel="3"
        onClick={() => onDigit('3')}
        variant="number"
      />
      <Button
        label="−"
        ariaLabel="Subtract"
        onClick={() => onOperation('subtract')}
        variant="operator"
        disabled={hasError || isEvaluating}
      />

      {/* Row 5: 0, Decimal point, Percentage, Add */}
      <Button
        label="0"
        ariaLabel="0"
        onClick={() => onDigit('0')}
        variant="number"
      />
      <Button
        label="."
        ariaLabel="Decimal Point"
        onClick={onDecimal}
        variant="number"
      />
      <Button
        label="%"
        ariaLabel="Percentage"
        onClick={() => onOperation('percentage')}
        variant="function"
        disabled={hasError || isEvaluating}
      />
      <Button
        label="+"
        ariaLabel="Add"
        onClick={() => onOperation('add')}
        variant="operator"
        disabled={hasError || isEvaluating}
      />

      {/* Row 6: Equals button spanning all 4 columns */}
      <Button
        label="="
        ariaLabel="Calculate"
        onClick={onEvaluate}
        variant="equal"
        className="col-span-4"
        disabled={hasError || isEvaluating}
      />
    </div>
  );
};
