package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewRepository(t *testing.T) {
	// テスト：Repositoryオブジェクトが正しく作成されることを確認する
	t.Run("Repositoryオブジェクトが正しく作成される", func(t *testing.T) {
		id := "12345"
		name := "test-repo"
		url := "https://github.com/example/test-repo"
		accessToken := "ghp_abcdefg"

		beforeCreate := time.Now()
		repo := NewRepository(id, name, url, accessToken)
		afterCreate := time.Now()

		assert.NotNil(t, repo)
		assert.Equal(t, id, repo.ID())
		assert.Equal(t, name, repo.Name())
		assert.Equal(t, url, repo.URL())
		assert.Equal(t, accessToken, repo.AccessToken())

		// 作成時間が現在時刻に近いことを確認
		assert.True(t, repo.CreatedAt().After(beforeCreate) || repo.CreatedAt().Equal(beforeCreate))
		assert.True(t, repo.CreatedAt().Before(afterCreate) || repo.CreatedAt().Equal(afterCreate))

		// 更新時間も同様に確認
		assert.True(t, repo.UpdatedAt().After(beforeCreate) || repo.UpdatedAt().Equal(beforeCreate))
		assert.True(t, repo.UpdatedAt().Before(afterCreate) || repo.UpdatedAt().Equal(afterCreate))
	})

	// テスト：アクセストークンが空の場合でも正しく作成される
	t.Run("アクセストークンが空の場合でも正しく作成される", func(t *testing.T) {
		id := "12345"
		name := "test-repo"
		url := "https://github.com/example/test-repo"
		accessToken := "" // 空のトークン

		repo := NewRepository(id, name, url, accessToken)

		assert.NotNil(t, repo)
		assert.Equal(t, accessToken, repo.AccessToken())
	})
}

func TestReconstructRepository(t *testing.T) {
	// テスト：永続化データからRepositoryオブジェクトが正しく再構築されることを確認する
	t.Run("永続化データからRepositoryオブジェクトが正しく再構築される", func(t *testing.T) {
		id := "12345"
		name := "test-repo"
		url := "https://github.com/example/test-repo"
		accessToken := "ghp_abcdefg"
		createdAt := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC)

		repo := ReconstructRepository(id, name, url, accessToken, createdAt, updatedAt)

		assert.NotNil(t, repo)
		assert.Equal(t, id, repo.ID())
		assert.Equal(t, name, repo.Name())
		assert.Equal(t, url, repo.URL())
		assert.Equal(t, accessToken, repo.AccessToken())
		assert.Equal(t, createdAt, repo.CreatedAt())
		assert.Equal(t, updatedAt, repo.UpdatedAt())
	})
}

func TestRepositoryGetters(t *testing.T) {
	// テスト：各ゲッターメソッドが正しい値を返すことを確認する
	t.Run("各ゲッターメソッドが正しい値を返す", func(t *testing.T) {
		id := "12345"
		name := "test-repo"
		url := "https://github.com/example/test-repo"
		accessToken := "ghp_abcdefg"
		createdAt := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC)

		repo := ReconstructRepository(id, name, url, accessToken, createdAt, updatedAt)

		assert.Equal(t, id, repo.ID())
		assert.Equal(t, name, repo.Name())
		assert.Equal(t, url, repo.URL())
		assert.Equal(t, accessToken, repo.AccessToken())
		assert.Equal(t, createdAt, repo.CreatedAt())
		assert.Equal(t, updatedAt, repo.UpdatedAt())
	})
}

func TestSetUpdatedAt(t *testing.T) {
	// テスト：SetUpdatedAtが現在時刻に更新時間を設定することを確認する
	t.Run("SetUpdatedAtが現在時刻に更新時間を設定する", func(t *testing.T) {
		id := "12345"
		name := "test-repo"
		url := "https://github.com/example/test-repo"
		accessToken := "ghp_abcdefg"
		createdAt := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC)

		repo := ReconstructRepository(id, name, url, accessToken, createdAt, updatedAt)

		// 元の更新時間を確認
		assert.Equal(t, updatedAt, repo.UpdatedAt())

		// 少し待ってから更新
		time.Sleep(time.Millisecond)
		beforeUpdate := time.Now()
		repo.SetUpdatedAt()
		afterUpdate := time.Now()

		// 更新時間が変更されていることを確認
		assert.NotEqual(t, updatedAt, repo.UpdatedAt())

		// 更新時間が現在時刻に近いことを確認
		assert.True(t, repo.UpdatedAt().After(beforeUpdate) || repo.UpdatedAt().Equal(beforeUpdate))
		assert.True(t, repo.UpdatedAt().Before(afterUpdate) || repo.UpdatedAt().Equal(afterUpdate))
	})
}

func TestSetAccessToken(t *testing.T) {
	// テスト：SetAccessTokenがトークンを更新し、更新時間も変更することを確認する
	t.Run("SetAccessTokenがトークンを更新し、更新時間も変更する", func(t *testing.T) {
		id := "12345"
		name := "test-repo"
		url := "https://github.com/example/test-repo"
		accessToken := "ghp_abcdefg"
		createdAt := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC)

		repo := ReconstructRepository(id, name, url, accessToken, createdAt, updatedAt)

		// 元のトークンと更新時間を確認
		assert.Equal(t, accessToken, repo.AccessToken())
		assert.Equal(t, updatedAt, repo.UpdatedAt())

		// 少し待ってからトークンを更新
		time.Sleep(time.Millisecond)
		newToken := "ghp_newtoken"
		beforeUpdate := time.Now()
		repo.SetAccessToken(newToken)
		afterUpdate := time.Now()

		// トークンが更新されていることを確認
		assert.Equal(t, newToken, repo.AccessToken())

		// 更新時間も変更されていることを確認
		assert.NotEqual(t, updatedAt, repo.UpdatedAt())
		assert.True(t, repo.UpdatedAt().After(beforeUpdate) || repo.UpdatedAt().Equal(beforeUpdate))
		assert.True(t, repo.UpdatedAt().Before(afterUpdate) || repo.UpdatedAt().Equal(afterUpdate))
	})

	// テスト：空のアクセストークンを設定できることを確認する
	t.Run("空のアクセストークンを設定できる", func(t *testing.T) {
		id := "12345"
		name := "test-repo"
		url := "https://github.com/example/test-repo"
		accessToken := "ghp_abcdefg"

		repo := NewRepository(id, name, url, accessToken)
		assert.Equal(t, accessToken, repo.AccessToken())

		// 空トークンを設定
		repo.SetAccessToken("")

		// トークンが空に更新されていることを確認
		assert.Equal(t, "", repo.AccessToken())
	})
}
