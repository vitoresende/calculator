import React from 'react';
import { Calculator } from './components/Calculator';

export const App: React.FC = () => {
  return (
    <main className="min-h-screen flex flex-col justify-between py-8 px-4 sm:px-6 lg:px-8 bg-slate-950 text-slate-100">
      {/* Header */}
      <header className="text-center max-w-xl mx-auto mb-6">
        <h1 className="text-2xl sm:text-3xl font-extrabold tracking-tight bg-gradient-to-r from-blue-400 via-indigo-300 to-amber-300 bg-clip-text text-transparent">
          Full-Stack Cloud Calculator
        </h1>
        <p className="mt-2 text-sm text-slate-400">
          Clean Architecture Go REST API backed by React, TypeScript, and Tailwind CSS
        </p>
      </header>

      {/* Main Content */}
      <section className="flex-1 flex items-center justify-center">
        <Calculator />
      </section>

      {/* Footer */}
      <footer className="mt-8 text-center text-xs text-slate-500 font-mono">
        Deployed to Google Cloud Run &amp; Firebase Hosting
      </footer>
    </main>
  );
};

export default App;
