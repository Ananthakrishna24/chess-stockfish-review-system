import React from 'react';
import { cn } from '@/lib/utils';
import { Clock } from 'lucide-react';

interface PlayerInfoBarProps {
  playerName: string;
  playerRating?: number;
  playerAvatarUrl?: string;
  isTurn?: boolean;
  className?: string;
}

const PlayerInfoBar = ({
  playerName,
  playerRating,
  playerAvatarUrl,
  isTurn,
  className,
}: PlayerInfoBarProps) => {
  return (
    <div className={cn("flex items-center justify-between", className)}>
      <div className="flex items-center gap-3">
        <div>
          <span className="font-semibold text-foreground">{playerName}</span>
          {playerRating && <span className="text-sm text-muted-foreground ml-2">({playerRating})</span>}
        </div>
      </div>
      <div className="flex items-center gap-2 bg-black/20 px-3 py-1.5 rounded">
        <Clock className="h-5 w-5 text-muted-foreground" />
        <span className="font-mono text-lg font-semibold">9:59</span>
      </div>
    </div>
  );
};

export default PlayerInfoBar; 