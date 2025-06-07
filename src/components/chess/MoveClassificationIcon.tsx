import React from 'react';
import { Star, ThumbsUp, AlertTriangle, X, Zap, Award, CheckCircle } from 'lucide-react';
import { MoveClassification } from '@/types/analysis';
import { cn } from '@/lib/utils';

interface MoveClassificationIconProps {
  classification: MoveClassification;
  size?: 'sm' | 'md' | 'lg';
  showLabel?: boolean;
  className?: string;
}

export function MoveClassificationIcon({ 
  classification, 
  size = 'md', 
  showLabel = false,
  className 
}: MoveClassificationIconProps) {
  const sizeClasses = {
    sm: 'w-3 h-3',
    md: 'w-4 h-4',
    lg: 'w-5 h-5'
  };

  const iconSize = sizeClasses[size];

  const getIcon = () => {
    switch (classification) {
      case 'brilliant':
        return <Zap className={cn(iconSize, 'text-orange-400')} />;
      case 'great':
        return <Award className={cn(iconSize, 'text-blue-400')} />;
      case 'best':
        return <Star className={cn(iconSize, 'text-green-400')} />;
      case 'excellent':
        return <CheckCircle className={cn(iconSize, 'text-green-300')} />;
      case 'good':
        return <ThumbsUp className={cn(iconSize, 'text-lime-400')} />;
      case 'inaccuracy':
        return <span className={cn('text-yellow-400 font-bold', size === 'sm' ? 'text-xs' : size === 'lg' ? 'text-lg' : 'text-sm')}>?</span>;
      case 'mistake':
        return <AlertTriangle className={cn(iconSize, 'text-orange-500')} />;
      case 'blunder':
        return <span className={cn('text-red-500 font-bold', size === 'sm' ? 'text-xs' : size === 'lg' ? 'text-lg' : 'text-sm')}>??</span>;
      case 'miss':
        return <X className={cn(iconSize, 'text-red-400')} />;
      default:
        return null;
    }
  };

  const getLabel = () => {
    switch (classification) {
      case 'brilliant': return 'Brilliant';
      case 'great': return 'Great';
      case 'best': return 'Best';
      case 'excellent': return 'Excellent';
      case 'good': return 'Good';
      case 'inaccuracy': return 'Inaccuracy';
      case 'mistake': return 'Mistake';
      case 'blunder': return 'Blunder';
      case 'miss': return 'Miss';
      default: return '';
    }
  };

  const getColor = () => {
    switch (classification) {
      case 'brilliant': return 'text-orange-400';
      case 'great': return 'text-blue-400';
      case 'best': return 'text-green-400';
      case 'excellent': return 'text-green-300';
      case 'good': return 'text-lime-400';
      case 'inaccuracy': return 'text-yellow-400';
      case 'mistake': return 'text-orange-500';
      case 'blunder': return 'text-red-500';
      case 'miss': return 'text-red-400';
      default: return 'text-muted-foreground';
    }
  };

  return (
    <div className={cn('flex items-center gap-1', className)}>
      {getIcon()}
      {showLabel && (
        <span className={cn('text-xs font-medium', getColor())}>
          {getLabel()}
        </span>
      )}
    </div>
  );
} 