import React from 'react';

export interface InputProps extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'size'> {
  label?: string;
  error?: string;
  helperText?: string;
  size?: 'sm' | 'md' | 'lg';
}

const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ className = '', label, error, helperText, size = 'md', ...props }, ref) => {
    const sizeClasses = {
      sm: 'px-3 py-1.5 text-sm',
      md: 'px-4 py-2.5 text-sm', 
      lg: 'px-5 py-3 text-base'
    };

    const inputStyles = `
      w-full border rounded-lg transition-all duration-200
      focus:outline-none focus:ring-2 focus:ring-opacity-50
      disabled:opacity-50 disabled:cursor-not-allowed
      placeholder:text-gray-400
      ${sizeClasses[size]}
      ${error ? 'border-red-300 focus:ring-red-500 focus:border-red-500' : 'focus:ring-[--chess-accent] focus:border-[--chess-accent]'}
      ${className}
    `;

    return (
      <div className="space-y-2">
        {label && (
          <label 
            className="block text-sm font-medium"
            style={{ color: 'var(--text-primary)' }}
          >
            {label}
          </label>
        )}
        <input
          ref={ref}
          className={inputStyles}
          style={{
            borderColor: error ? undefined : 'var(--border-medium)',
            backgroundColor: 'var(--surface)',
            color: 'var(--text-primary)'
          }}
          {...props}
        />
        {error && (
          <p className="text-sm font-medium" style={{ color: 'var(--chess-danger)' }}>
            {error}
          </p>
        )}
        {helperText && !error && (
          <p className="text-sm" style={{ color: 'var(--text-secondary)' }}>
            {helperText}
          </p>
        )}
      </div>
    );
  }
);

Input.displayName = 'Input';

export interface TextareaProps extends Omit<React.TextareaHTMLAttributes<HTMLTextAreaElement>, 'size'> {
  label?: string;
  error?: string;
  helperText?: string;
  size?: 'sm' | 'md' | 'lg';
}

export const Textarea = React.forwardRef<HTMLTextAreaElement, TextareaProps>(
  ({ className = '', label, error, helperText, size = 'md', ...props }, ref) => {
    const sizeClasses = {
      sm: 'px-3 py-1.5 text-sm',
      md: 'px-4 py-2.5 text-sm',
      lg: 'px-5 py-3 text-base'
    };

    const textareaStyles = `
      w-full border rounded-lg resize-none transition-all duration-200
      focus:outline-none focus:ring-2 focus:ring-opacity-50
      disabled:opacity-50 disabled:cursor-not-allowed
      placeholder:text-gray-400
      ${sizeClasses[size]}
      ${error ? 'border-red-300 focus:ring-red-500 focus:border-red-500' : 'focus:ring-[--chess-accent] focus:border-[--chess-accent]'}
      ${className}
    `;

    return (
      <div className="space-y-2">
        {label && (
          <label 
            className="block text-sm font-medium"
            style={{ color: 'var(--text-primary)' }}
          >
            {label}
          </label>
        )}
        <textarea
          ref={ref}
          className={textareaStyles}
          style={{
            borderColor: error ? undefined : 'var(--border-medium)',
            backgroundColor: 'var(--surface)',
            color: 'var(--text-primary)'
          }}
          {...props}
        />
        {error && (
          <p className="text-sm font-medium" style={{ color: 'var(--chess-danger)' }}>
            {error}
          </p>
        )}
        {helperText && !error && (
          <p className="text-sm" style={{ color: 'var(--text-secondary)' }}>
            {helperText}
          </p>
        )}
      </div>
    );
  }
);

Textarea.displayName = 'Textarea';

export default Input; 