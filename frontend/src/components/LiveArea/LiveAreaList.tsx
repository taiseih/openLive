// Single Responsibility Principle: ライブエリア一覧の表示のみに責任を持つ
import React from 'react';
import { LiveArea } from '@/types/models';
import { LiveAreaCard } from './LiveAreaCard';
import { Loading } from '@/components/common/Loading';

interface LiveAreaListProps {
  liveAreas: LiveArea[];
  isLoading?: boolean;
  onAreaClick?: (liveArea: LiveArea) => void;
}

export const LiveAreaList: React.FC<LiveAreaListProps> = ({
  liveAreas,
  isLoading = false,
  onAreaClick,
}) => {
  if (isLoading) {
    return <Loading text="読み込み中..." />;
  }

  if (liveAreas.length === 0) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500">ライブエリアがありません</p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {liveAreas.map((area) => (
        <LiveAreaCard key={area.id} liveArea={area} onClick={onAreaClick} />
      ))}
    </div>
  );
};

