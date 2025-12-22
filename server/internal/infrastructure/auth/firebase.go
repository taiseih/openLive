package auth

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// FirebaseAuth はFirebase認証クライアントをラップします
type FirebaseAuth struct {
	client *auth.Client
}

// NewFirebaseAuth は新しいFirebase認証クライアントを作成します
func NewFirebaseAuth(ctx context.Context, projectID, emulatorHost, credentialsPath string) (*FirebaseAuth, error) {
	var opts []option.ClientOption

	// 認証情報ファイルが指定されている場合
	if credentialsPath != "" {
		opts = append(opts, option.WithCredentialsFile(credentialsPath))
	}

	config := &firebase.Config{
		ProjectID: projectID,
	}

	app, err := firebase.NewApp(ctx, config, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize firebase app: %w", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get auth client: %w", err)
	}

	log.Println("Successfully initialized Firebase Auth")
	return &FirebaseAuth{client: client}, nil
}

// VerifyIDToken はFirebase IDトークンを検証します
func (f *FirebaseAuth) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	token, err := f.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}
	return token, nil
}

// GetUser はユーザーIDからユーザー情報を取得します
func (f *FirebaseAuth) GetUser(ctx context.Context, uid string) (*auth.UserRecord, error) {
	user, err := f.client.GetUser(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// GetClient はFirebase Authクライアントを返します
func (f *FirebaseAuth) GetClient() *auth.Client {
	return f.client
}
