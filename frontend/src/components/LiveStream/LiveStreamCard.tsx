// Single Responsibility Principle: ライブ配信カードの表示のみに責任を持つ
import React from 'react';
import { LiveStream } from '@/types/models';
import { Card, CardBody } from '@/components/common/Card';

interface LiveStreamCardProps {
  stream: LiveStream;
  onClick?: (stream: LiveStream) => void;
}

export const LiveStreamCard: React.FC<LiveStreamCardProps> = ({ stream, onClick }) => {
  const getStatusBadge = () => {
    switch (stream.status) {
      case 'live':
        return (
          <span className="px-2 py-1 text-xs bg-red-500 text-white rounded-full animate-pulse">
            LIVE
          </span>
        );
      case 'scheduled':
        return (
          <span className="px-2 py-1 text-xs bg-blue-100 text-blue-700 rounded-full">
            予定
          </span>
        );
      case 'ended':
        return (
          <span className="px-2 py-1 text-xs bg-gray-100 text-gray-700 rounded-full">
            終了
          </span>
        );
    }
  };

  return (
    <Card onClick={() => onClick?.(stream)}>
      {stream.thumbnailUrl && (
        <div className="relative aspect-video bg-gray-200">
          <img
            src={stream.thumbnailUrl}
            alt={stream.title}
            className="w-full h-full object-cover"
          />
          <div className="absolute top-2 right-2">{getStatusBadge()}</div>
        </div>
      )}
      <CardBody>
        <h3 className="text-lg font-semibold text-gray-800 mb-1">{stream.title}</h3>
        {stream.description && (
          <p className="text-sm text-gray-600 line-clamp-2">{stream.description}</p>
        )}
        {stream.startedAt && (
          <p className="text-xs text-gray-500 mt-2">
            開始: {new Date(stream.startedAt).toLocaleString('ja-JP')}
          </p>
        )}
      </CardBody>
    </Card>
  );
};

