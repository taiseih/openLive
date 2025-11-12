'use client';

import { useRouter } from 'next/navigation';
import { useAuth } from '@/hooks/useAuth';
import { useLiveAreas } from '@/hooks/useLiveAreas';
import { Button } from '@/components/common/Button';
import { Loading } from '@/components/common/Loading';
import { LiveAreaList } from '@/components/LiveArea/LiveAreaList';

export default function AreasPage() {
  const router = useRouter();
  const { isAuthenticated, isLoading: authLoading } = useAuth();
  const { liveAreas, isLoading, error } = useLiveAreas();

  if (authLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <Loading text="読み込み中..." />
      </div>
    );
  }

  if (!isAuthenticated) {
    router.push('/');
    return null;
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* ヘッダー */}
      <header className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex justify-between items-center">
            <h1 className="text-2xl font-bold text-primary-600 cursor-pointer" onClick={() => router.push('/')}>
              openLive
            </h1>
            <Button onClick={() => router.push('/areas/new')}>
              新しいエリアを作成
            </Button>
          </div>
        </div>
      </header>

      {/* メインコンテンツ */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <h2 className="text-3xl font-bold text-gray-900 mb-8">ライブエリア一覧</h2>

        {error ? (
          <div className="bg-red-50 border border-red-200 rounded-lg p-4">
            <p className="text-red-600">エラーが発生しました: {error.message}</p>
          </div>
        ) : (
          <LiveAreaList
            liveAreas={liveAreas}
            isLoading={isLoading}
            onAreaClick={(area) => router.push(`/areas/${area.id}`)}
          />
        )}
      </main>
    </div>
  );
}

