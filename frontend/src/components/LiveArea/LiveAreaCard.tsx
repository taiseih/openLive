// Single Responsibility Principle: ライブエリアカードの表示のみに責任を持つ
import React from 'react';
import { LiveArea } from '@/types/models';
import { Card, CardBody } from '@/components/common/Card';

interface LiveAreaCardProps {
  liveArea: LiveArea;
  onClick?: (liveArea: LiveArea) => void;
}

export const LiveAreaCard: React.FC<LiveAreaCardProps> = ({ liveArea, onClick }) => {
  return (
    <Card onClick={() => onClick?.(liveArea)}>
      {liveArea.thumbnailUrl && (
        <div className="aspect-video bg-gray-200">
          <img
            src={liveArea.thumbnailUrl}
            alt={liveArea.name}
            className="w-full h-full object-cover"
          />
        </div>
      )}
      <CardBody>
        <div className="flex items-start justify-between">
          <div className="flex-1">
            <h3 className="text-lg font-semibold text-gray-800 mb-1">{liveArea.name}</h3>
            {liveArea.description && (
              <p className="text-sm text-gray-600 line-clamp-2">{liveArea.description}</p>
            )}
          </div>
          {liveArea.isPublic && (
            <span className="ml-2 px-2 py-1 text-xs bg-green-100 text-green-700 rounded-full">
              公開
            </span>
          )}
        </div>
      </CardBody>
    </Card>
  );
};

