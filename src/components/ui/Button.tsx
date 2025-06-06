import React from 'react';

export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'outline' | 'ghost' | 'success' | 'warning' | 'danger';
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl';
  children: React.ReactNode;
  isLoading?: boolean;
  leftIcon?: React.ReactNode;
  rightIcon?: React.ReactNode;
  fullWidth?: boolean;
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ 
    className = '', 
    variant = 'primary', 
    size = 'md', 
    children, 
    isLoading, 
    disabled, 
    leftIcon, 
    rightIcon, 
    fullWidth = false,
    ...props 
  }, ref) => {
    const baseStyles = `
      inline-flex items-center justify-center font-medium transition-all duration-200 
      focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-opacity-50
      disabled:opacity-50 disabled:pointer-events-none disabled:cursor-not-allowed
      active:scale-[0.98] rounded-lg relative overflow-hidden
      ${fullWidth ? 'w-full' : ''}
    `;
    
    const variants = {
      primary: `
        text-white shadow-sm hover:shadow-md active:shadow-sm
        focus:ring-[--chess-accent]
      `,
      secondary: `
        text-white shadow-sm hover:shadow-md active:shadow-sm
        focus:ring-[--chess-secondary]
      `,
      outline: `
        bg-transparent border-2 hover:shadow-sm
        focus:ring-[--chess-accent]
      `,
      ghost: `
        bg-transparent hover:shadow-sm
        focus:ring-[--chess-accent]
      `,
      success: `
        text-white shadow-sm hover:shadow-md active:shadow-sm
        focus:ring-[--chess-success]
      `,
      warning: `
        text-white shadow-sm hover:shadow-md active:shadow-sm
        focus:ring-[--chess-warning]
      `,
      danger: `
        text-white shadow-sm hover:shadow-md active:shadow-sm
        focus:ring-[--chess-danger]
      `
    };

    const sizes = {
      xs: 'h-7 px-2.5 text-xs gap-1.5',
      sm: 'h-8 px-3 text-sm gap-1.5',
      md: 'h-10 px-4 text-sm gap-2',
      lg: 'h-12 px-6 text-base gap-2',
      xl: 'h-14 px-8 text-lg gap-2.5'
    };

    const getVariantStyles = () => {
      const baseVariant = variants[variant];
      
      switch (variant) {
        case 'primary':
          return `${baseVariant} bg-[--chess-accent] hover:bg-[--chess-accent]/90 border-transparent`;
        case 'secondary':
          return `${baseVariant} bg-[--chess-secondary] hover:bg-[--chess-secondary]/90 border-transparent`;
        case 'outline':
          return `${baseVariant} border-[--border-medium] text-[--text-primary] hover:bg-[--surface-secondary] hover:border-[--chess-accent] hover:text-[--chess-accent]`;
        case 'ghost':
          return `${baseVariant} text-[--text-secondary] hover:bg-[--surface-secondary] hover:text-[--text-primary]`;
        case 'success':
          return `${baseVariant} bg-[--chess-success] hover:bg-[--chess-success]/90 border-transparent`;
        case 'warning':
          return `${baseVariant} bg-[--chess-warning] hover:bg-[--chess-warning]/90 border-transparent`;
        case 'danger':
          return `${baseVariant} bg-[--chess-danger] hover:bg-[--chess-danger]/90 border-transparent`;
        default:
          return baseVariant;
      }
    };

    const sizeStyles = sizes[size];
    const variantStyles = getVariantStyles();
    
    return (
      <button
        ref={ref}
        className={`${baseStyles} ${variantStyles} ${sizeStyles} ${className}`}
        disabled={disabled || isLoading}
        style={{
          backgroundColor: variant === 'primary' ? 'var(--chess-accent)' : 
                          variant === 'secondary' ? 'var(--chess-secondary)' :
                          variant === 'success' ? 'var(--chess-success)' :
                          variant === 'warning' ? 'var(--chess-warning)' :
                          variant === 'danger' ? 'var(--chess-danger)' : undefined,
          borderColor: variant === 'outline' ? 'var(--border-medium)' : undefined,
          color: variant === 'outline' || variant === 'ghost' ? 'var(--text-primary)' : undefined
        }}
        {...props}
      >
        {isLoading && (
          <svg
            className={`animate-spin ${size === 'xs' ? 'h-3 w-3' : size === 'sm' ? 'h-3.5 w-3.5' : 'h-4 w-4'}`}
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              className="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              strokeWidth="4"
            />
            <path
              className="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            />
          </svg>
        )}
        
        {!isLoading && leftIcon && (
          <span className={`${size === 'xs' ? '-ml-0.5' : '-ml-1'}`}>
            {leftIcon}
          </span>
        )}
        
        <span className={isLoading ? 'ml-2' : ''}>{children}</span>
        
        {!isLoading && rightIcon && (
          <span className={`${size === 'xs' ? '-mr-0.5' : '-mr-1'}`}>
            {rightIcon}
          </span>
        )}
      </button>
    );
  }
);

Button.displayName = 'Button';

export default Button; 