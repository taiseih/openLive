"use client";

import { Button } from "@/components/common/Button";
import { Loading } from "@/components/common/Loading";
import { LiveStreamCard } from "@/components/LiveStream/LiveStreamCard";
import { useAuth } from "@/hooks/useAuth";
import { useLiveStreams } from "@/hooks/useLiveStream";
import { useRouter } from "next/navigation";

export default function Home() {
  const router = useRouter();
  const { user, isAuthenticated, isLoading: authLoading, signOut } = useAuth();
  const { streams, isLoading: streamsLoading } = useLiveStreams();

  if (authLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <Loading text="読み込み中..." />
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* ヘッダー */}
      <header className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex justify-between items-center">
            <h1 className="text-2xl font-bold text-primary-600">openLive</h1>
            <div className="flex items-center space-x-4">
              {isAuthenticated ? (
                <>
                  <span className="text-gray-700">
                    こんにちは、{user?.displayName || "ゲスト"}さん
                  </span>
                  <Button onClick={() => router.push("/areas")}>
                    マイエリア
                  </Button>
                  <Button variant="secondary" onClick={signOut}>
                    ログアウト
                  </Button>
                </>
              ) : (
                <Button onClick={() => router.push("/login")}>ログイン</Button>
              )}
            </div>
          </div>
        </div>
      </header>

      {/* メインコンテンツ */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* ヒーローセクション */}
        <section className="mb-12 text-center">
          <h2 className="text-4xl font-bold text-gray-900 mb-4">
            自由にライブ配信を楽しもう
          </h2>
          <p className="text-xl text-gray-600 mb-8">
            あなただけの配信広場を作って、仲間とつながろう
          </p>
          {isAuthenticated && (
            <Button onClick={() => router.push("/areas/new")}>
              ライブエリアを作成
            </Button>
          )}
        </section>

        {/* 配信中の配信一覧 */}
        <section>
          <h3 className="text-2xl font-bold text-gray-900 mb-6">
            配信中のライブ
          </h3>
          {streamsLoading ? (
            <Loading />
          ) : streams.length === 0 ? (
            <div className="text-center py-12 bg-white rounded-lg shadow">
              <p className="text-gray-500">現在配信中のライブはありません</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {streams.map((stream) => (
                <LiveStreamCard
                  key={stream.id}
                  stream={stream}
                  onClick={(s) => router.push(`/streams/${s.id}`)}
                />
              ))}
            </div>
          )}
        </section>
      </main>

      {/* フッター */}
      <footer className="bg-white border-t border-gray-200 mt-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <p className="text-center text-gray-600">
            © 2024 openLive. All rights reserved.
          </p>
        </div>
      </footer>
    </div>
  );
}
