import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import App from '../App';

describe('App Component', () => {
  it('renders application header, calculator interface, and footer', () => {
    render(<App />);

    expect(
      screen.getByRole('heading', { name: /Full-Stack Cloud Calculator/i })
    ).toBeInTheDocument();
    expect(screen.getByTestId('calculator-display')).toBeInTheDocument();
    expect(
      screen.getByText(/Deployed to Google Cloud Run & Firebase Hosting/i)
    ).toBeInTheDocument();
  });
});
