import React from 'react';

interface CardProps {
  children: React.ReactNode;
  className?: string;
  variant?: 'default' | 'elevated' | 'outlined' | 'minimal';
}

export function Card({ children, className = '', variant = 'default' }: CardProps) {
  const baseClasses = 'rounded-xl transition-all duration-200';
  
  const variantClasses = {
    default: 'bg-white border border-[--border-light] shadow-[--shadow-sm] hover:shadow-[--shadow-md]',
    elevated: 'bg-white border-0 shadow-[--shadow-lg] hover:shadow-[--shadow-xl]',
    outlined: 'bg-transparent border-2 border-[--border-medium] hover:border-[--chess-accent] hover:bg-white/50',
    minimal: 'bg-white/60 backdrop-blur-sm border border-[--border-light]/50'
  };

  return (
    <div 
      className={`${baseClasses} ${variantClasses[variant]} ${className}`}
      style={{
        backgroundColor: variant === 'default' || variant === 'elevated' ? 'var(--surface)' : undefined,
        borderColor: variant === 'default' ? 'var(--border-light)' : undefined,
        boxShadow: variant === 'default' ? 'var(--shadow-sm)' : variant === 'elevated' ? 'var(--shadow-lg)' : undefined
      }}
    >
      {children}
    </div>
  );
}

interface CardHeaderProps {
  children: React.ReactNode;
  className?: string;
  noBorder?: boolean;
}

export function CardHeader({ children, className = '', noBorder = false }: CardHeaderProps) {
  return (
    <div 
      className={`px-6 py-5 ${!noBorder ? 'border-b' : ''} ${className}`}
      style={{ 
        borderColor: !noBorder ? 'var(--border-light)' : undefined 
      }}
    >
      {children}
    </div>
  );
}

interface CardContentProps {
  children: React.ReactNode;
  className?: string;
  size?: 'sm' | 'md' | 'lg';
}

export function CardContent({ children, className = '', size = 'md' }: CardContentProps) {
  const sizeClasses = {
    sm: 'px-4 py-3',
    md: 'px-6 py-5',
    lg: 'px-8 py-6'
  };

  return (
    <div className={`${sizeClasses[size]} ${className}`}>
      {children}
    </div>
  );
}

interface CardTitleProps {
  children: React.ReactNode;
  className?: string;
  size?: 'sm' | 'md' | 'lg';
}

export function CardTitle({ children, className = '', size = 'md' }: CardTitleProps) {
  const sizeClasses = {
    sm: 'text-base',
    md: 'text-lg',
    lg: 'text-xl'
  };

  return (
    <h3 
      className={`font-semibold leading-tight ${sizeClasses[size]} ${className}`}
      style={{ color: 'var(--text-primary)' }}
    >
      {children}
    </h3>
  );
}

interface CardDescriptionProps {
  children: React.ReactNode;
  className?: string;
}

export function CardDescription({ children, className = '' }: CardDescriptionProps) {
  return (
    <p 
      className={`text-sm leading-relaxed mt-1.5 ${className}`}
      style={{ color: 'var(--text-secondary)' }}
    >
      {children}
    </p>
  );
}

interface CardFooterProps {
  children: React.ReactNode;
  className?: string;
}

export function CardFooter({ children, className = '' }: CardFooterProps) {
  return (
    <div 
      className={`px-6 py-4 border-t bg-[--surface-secondary] rounded-b-xl ${className}`}
      style={{ 
        borderColor: 'var(--border-light)',
        backgroundColor: 'var(--surface-secondary)'
      }}
    >
      {children}
    </div>
  );
} 