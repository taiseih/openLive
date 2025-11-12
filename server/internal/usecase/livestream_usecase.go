package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/taiseihayashizaki/openLive/server/internal/domain"
)

// LiveStreamUseCase はライブ配信関連のユースケースを提供します
type LiveStreamUseCase interface {
	CreateStream(ctx context.Context, liveAreaID, streamerID uuid.UUID, title, description string) (*domain.LiveStream, error)
	GetStream(ctx context.Context, id uuid.UUID) (*domain.LiveStream, error)
	GetStreamsByLiveArea(ctx context.Context, liveAreaID uuid.UUID) ([]*domain.LiveStream, error)
	GetLiveStreams(ctx context.Context, limit, offset int) ([]*domain.LiveStream, error)
	StartStream(ctx context.Context, id, streamerID uuid.UUID) (*domain.LiveStream, error)
	EndStream(ctx context.Context, id, streamerID uuid.UUID) (*domain.LiveStream, error)
	JoinStream(ctx context.Context, streamID uuid.UUID, userID *uuid.UUID) (*domain.LiveStreamViewer, error)
	LeaveStream(ctx context.Context, viewerID uuid.UUID) error
	GetViewers(ctx context.Context, streamID uuid.UUID) ([]*domain.LiveStreamViewer, error)
	GetViewerCount(ctx context.Context, streamID uuid.UUID) (int, error)
}

type liveStreamUseCase struct {
	streamRepo       domain.LiveStreamRepository
	viewerRepo       domain.LiveStreamViewerRepository
	liveAreaRepo     domain.LiveAreaRepository
	liveAreaMemberRepo domain.LiveAreaMemberRepository
}

// NewLiveStreamUseCase は新しいLiveStreamUseCaseを作成します
func NewLiveStreamUseCase(
	streamRepo domain.LiveStreamRepository,
	viewerRepo domain.LiveStreamViewerRepository,
	liveAreaRepo domain.LiveAreaRepository,
	liveAreaMemberRepo domain.LiveAreaMemberRepository,
) LiveStreamUseCase {
	return &liveStreamUseCase{
		streamRepo:       streamRepo,
		viewerRepo:       viewerRepo,
		liveAreaRepo:     liveAreaRepo,
		liveAreaMemberRepo: liveAreaMemberRepo,
	}
}

func (u *liveStreamUseCase) CreateStream(ctx context.Context, liveAreaID, streamerID uuid.UUID, title, description string) (*domain.LiveStream, error) {
	// ライブエリアが存在するか確認
	liveArea, err := u.liveAreaRepo.FindByID(ctx, liveAreaID)
	if err != nil {
		return nil, fmt.Errorf("live area not found: %w", err)
	}

	// 配信者がメンバーかオーナーか確認
	member, err := u.liveAreaMemberRepo.FindByLiveAreaAndUser(ctx, liveAreaID, streamerID)
	if err != nil || (member == nil && liveArea.OwnerID != streamerID) {
		return nil, fmt.Errorf("unauthorized: not a member of this live area")
	}

	stream := domain.NewLiveStream(liveAreaID, streamerID, title, description)
	if err := u.streamRepo.Create(ctx, stream); err != nil {
		return nil, fmt.Errorf("failed to create stream: %w", err)
	}

	return stream, nil
}

func (u *liveStreamUseCase) GetStream(ctx context.Context, id uuid.UUID) (*domain.LiveStream, error) {
	stream, err := u.streamRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get stream: %w", err)
	}
	return stream, nil
}

func (u *liveStreamUseCase) GetStreamsByLiveArea(ctx context.Context, liveAreaID uuid.UUID) ([]*domain.LiveStream, error) {
	streams, err := u.streamRepo.FindByLiveAreaID(ctx, liveAreaID)
	if err != nil {
		return nil, fmt.Errorf("failed to get streams: %w", err)
	}
	return streams, nil
}

func (u *liveStreamUseCase) GetLiveStreams(ctx context.Context, limit, offset int) ([]*domain.LiveStream, error) {
	streams, err := u.streamRepo.FindLiveStreams(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get live streams: %w", err)
	}
	return streams, nil
}

func (u *liveStreamUseCase) StartStream(ctx context.Context, id, streamerID uuid.UUID) (*domain.LiveStream, error) {
	stream, err := u.streamRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("stream not found: %w", err)
	}

	// 配信者チェック
	if stream.StreamerID != streamerID {
		return nil, fmt.Errorf("unauthorized: not the streamer")
	}

	// 既に配信中の場合はエラー
	if stream.IsLive() {
		return nil, fmt.Errorf("stream is already live")
	}

	stream.Start()
	if err := u.streamRepo.Update(ctx, stream); err != nil {
		return nil, fmt.Errorf("failed to start stream: %w", err)
	}

	return stream, nil
}

func (u *liveStreamUseCase) EndStream(ctx context.Context, id, streamerID uuid.UUID) (*domain.LiveStream, error) {
	stream, err := u.streamRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("stream not found: %w", err)
	}

	// 配信者チェック
	if stream.StreamerID != streamerID {
		return nil, fmt.Errorf("unauthorized: not the streamer")
	}

	stream.End()
	if err := u.streamRepo.Update(ctx, stream); err != nil {
		return nil, fmt.Errorf("failed to end stream: %w", err)
	}

	return stream, nil
}

func (u *liveStreamUseCase) JoinStream(ctx context.Context, streamID uuid.UUID, userID *uuid.UUID) (*domain.LiveStreamViewer, error) {
	// 配信が存在するか確認
	stream, err := u.streamRepo.FindByID(ctx, streamID)
	if err != nil {
		return nil, fmt.Errorf("stream not found: %w", err)
	}

	// 配信中でない場合はエラー
	if !stream.IsLive() {
		return nil, fmt.Errorf("stream is not live")
	}

	viewer := domain.NewLiveStreamViewer(streamID, userID)
	if err := u.viewerRepo.Create(ctx, viewer); err != nil {
		return nil, fmt.Errorf("failed to join stream: %w", err)
	}

	return viewer, nil
}

func (u *liveStreamUseCase) LeaveStream(ctx context.Context, viewerID uuid.UUID) error {
	// 視聴者情報を取得（実装はシンプル化のため省略）
	// 実際にはviewerRepoにFindByIDメソッドを追加する必要があります
	return fmt.Errorf("not implemented")
}

func (u *liveStreamUseCase) GetViewers(ctx context.Context, streamID uuid.UUID) ([]*domain.LiveStreamViewer, error) {
	viewers, err := u.viewerRepo.FindActiveByStreamID(ctx, streamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get viewers: %w", err)
	}
	return viewers, nil
}

func (u *liveStreamUseCase) GetViewerCount(ctx context.Context, streamID uuid.UUID) (int, error) {
	count, err := u.viewerRepo.CountActiveByStreamID(ctx, streamID)
	if err != nil {
		return 0, fmt.Errorf("failed to get viewer count: %w", err)
	}
	return count, nil
}

