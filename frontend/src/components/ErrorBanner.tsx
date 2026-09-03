import React from 'react';

interface ErrorBannerProps {
  message: string | null;
  onDismiss?: () => void;
}

export const ErrorBanner: React.FC<ErrorBannerProps> = ({ message, onDismiss }) => {
  if (!message) return null;

  return (
    <div
      role="alert"
      className="w-full mb-3 px-4 py-2.5 bg-rose-950/90 border border-rose-700/80 text-rose-200 rounded-lg text-sm flex items-center justify-between shadow-sm animate-in fade-in duration-200"
    >
      <div className="flex items-center gap-2">
        <svg
          className="w-4 h-4 text-rose-400 shrink-0"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          aria-hidden="true"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth="2"
            d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
          />
        </svg>
        <span className="font-medium">{message}</span>
      </div>
      {onDismiss && (
        <button
          type="button"
          aria-label="Dismiss error"
          onClick={onDismiss}
          className="text-rose-400 hover:text-rose-200 text-xs font-mono px-2 py-0.5 rounded hover:bg-rose-900/50"
        >
          ✕
        </button>
      )}
    </div>
  );
};
