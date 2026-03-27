package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"smartcommunity/internal/global"
	"smartcommunity/internal/model"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const communityGreenPointsLeaderboardKey = "community_green_points_total_earned_leaderboard"

type GreenPointService struct {
	aiService      AIService
	storageService StorageService
}

type GarbageRewardResponse struct {
	ImageURL    string `json:"image_url"`
	ObjectKey   string `json:"object_key"`
	Points      int    `json:"points"`
	Reason      string `json:"reason"`
	GreenPoints int    `json:"green_points"`
}

type GreenPointLeaderboardItem struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	RealName string `json:"real_name"`
	Avatar   string `json:"avatar"`
	Rank     int    `json:"rank"`
	Points   int    `json:"points"`
	Nickname string `json:"nickname"`
}

func (s *GreenPointService) UploadGarbage(userID int64, fileHeader *multipart.FileHeader) (*GarbageRewardResponse, error) {
	imageURL, objectKey, err := s.storageService.UploadMultipartFile(fileHeader, "green-points")
	if err != nil {
		// Development fallback: if MinIO is unreachable, use base64 data URL directly for vision model.
		imageURL, err = buildDataURLFromFile(fileHeader)
		if err != nil {
			log.Printf("upload garbage image failed and fallback failed: %v", err)
			return nil, err
		}
		objectKey = ""
		log.Printf("upload garbage image failed, fallback to base64: %v", err)
	}

	recognitionResult, err := s.aiService.RecognizeGarbage(imageURL)
	if err != nil {
		log.Printf("AI garbage recognition failed: %v", err)
		return nil, err
	}

	result := &GarbageRewardResponse{
		ImageURL:  imageURL,
		ObjectKey: objectKey,
		Points:    recognitionResult.Points,
		Reason:    recognitionResult.Reason,
	}

	err = global.DB.Transaction(func(tx *gorm.DB) error {
		var user model.SysUser
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, userID).Error; err != nil {
			return errors.New("user not found")
		}

		if err := tx.Model(&model.SysUser{}).
			Where("id = ?", userID).
			Update("green_points", gorm.Expr("green_points + ?", recognitionResult.Points)).Error; err != nil {
			return err
		}

		record := model.GreenPointRecord{
			UserID:    userID,
			Action:    "garbage_classification",
			Points:    recognitionResult.Points,
			CreatedAt: time.Now(),
		}
		if err := tx.Create(&record).Error; err != nil {
			return err
		}

		result.GreenPoints = user.GreenPoints + recognitionResult.Points
		return nil
	})
	if err != nil {
		log.Printf("persist green point reward failed: %v", err)
		return nil, err
	}

	if err := adjustLeaderboardScore(userID, recognitionResult.Points); err != nil {
		log.Printf("update leaderboard after garbage reward failed: %v", err)
	}

	return result, nil
}

func (s *GreenPointService) GetLeaderboard(limit int64) ([]GreenPointLeaderboardItem, error) {
	if limit <= 0 {
		limit = 10
	}

	ctx := context.Background()
	items, err := global.RDB.ZRevRangeWithScores(ctx, communityGreenPointsLeaderboardKey, 0, limit-1).Result()
	if err != nil {
		log.Printf("read leaderboard from redis failed: %v", err)
		return nil, err
	}

	if len(items) == 0 {
		return warmUpLeaderboardFromMySQL(limit)
	}

	return enrichLeaderboardItems(items)
}

func adjustLeaderboardScore(userID int64, delta int) error {
	if delta == 0 {
		return nil
	}
	ctx := context.Background()
	return global.RDB.ZIncrBy(ctx, communityGreenPointsLeaderboardKey, float64(delta), strconv.FormatInt(userID, 10)).Err()
}

func syncLeaderboardScore(userID int64, score int) error {
	ctx := context.Background()
	return global.RDB.ZAdd(ctx, communityGreenPointsLeaderboardKey, redis.Z{
		Score:  float64(score),
		Member: strconv.FormatInt(userID, 10),
	}).Err()
}

func warmUpLeaderboardFromMySQL(limit int64) ([]GreenPointLeaderboardItem, error) {
	type leaderboardAggregate struct {
		UserID int64 `gorm:"column:user_id"`
		Points int64 `gorm:"column:points"`
	}

	var aggregates []leaderboardAggregate
	if err := global.DB.Model(&model.GreenPointRecord{}).
		Select("user_id, COALESCE(SUM(points), 0) AS points").
		Where("points > 0").
		Group("user_id").
		Order("points desc").
		Limit(int(limit)).
		Scan(&aggregates).Error; err != nil {
		return nil, err
	}

	if len(aggregates) == 0 {
		return warmUpLeaderboardFromUserBalance(limit)
	}

	userIDs := make([]int64, 0, len(aggregates))
	for _, item := range aggregates {
		userIDs = append(userIDs, item.UserID)
	}

	var users []model.SysUser
	if err := global.DB.Select("id", "username", "real_name", "avatar").
		Where("id IN ?", userIDs).
		Find(&users).Error; err != nil {
		return nil, err
	}

	userMap := make(map[int64]model.SysUser, len(users))
	for _, user := range users {
		userMap[user.ID] = user
	}

	zs := make([]redis.Z, 0, len(aggregates))
	result := make([]GreenPointLeaderboardItem, 0, len(aggregates))
	for index, item := range aggregates {
		user := userMap[item.UserID]
		points := int(item.Points)
		zs = append(zs, redis.Z{
			Score:  float64(points),
			Member: strconv.FormatInt(item.UserID, 10),
		})
		result = append(result, GreenPointLeaderboardItem{
			UserID:   item.UserID,
			Username: user.Username,
			RealName: user.RealName,
			Avatar:   user.Avatar,
			Rank:     index + 1,
			Points:   points,
			Nickname: buildUserNickname(user),
		})
	}

	if len(zs) > 0 {
		if err := global.RDB.ZAdd(context.Background(), communityGreenPointsLeaderboardKey, zs...).Err(); err != nil {
			log.Printf("warm up leaderboard failed: %v", err)
		}
	}

	return result, nil
}

func warmUpLeaderboardFromUserBalance(limit int64) ([]GreenPointLeaderboardItem, error) {
	var users []model.SysUser
	if err := global.DB.Select("id", "username", "real_name", "avatar", "green_points").
		Where("green_points > 0").
		Order("green_points desc").
		Limit(int(limit)).
		Find(&users).Error; err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return []GreenPointLeaderboardItem{}, nil
	}

	zs := make([]redis.Z, 0, len(users))
	result := make([]GreenPointLeaderboardItem, 0, len(users))
	for index, user := range users {
		zs = append(zs, redis.Z{
			Score:  float64(user.GreenPoints),
			Member: strconv.FormatInt(user.ID, 10),
		})
		result = append(result, GreenPointLeaderboardItem{
			UserID:   user.ID,
			Username: user.Username,
			RealName: user.RealName,
			Avatar:   user.Avatar,
			Rank:     index + 1,
			Points:   user.GreenPoints,
			Nickname: buildUserNickname(user),
		})
	}

	if len(zs) > 0 {
		if err := global.RDB.ZAdd(context.Background(), communityGreenPointsLeaderboardKey, zs...).Err(); err != nil {
			log.Printf("warm up leaderboard failed: %v", err)
		}
	}

	return result, nil
}

func enrichLeaderboardItems(items []redis.Z) ([]GreenPointLeaderboardItem, error) {
	userIDs := make([]int64, 0, len(items))
	for _, item := range items {
		userID, err := memberToUserID(item.Member)
		if err != nil {
			continue
		}
		userIDs = append(userIDs, userID)
	}

	var users []model.SysUser
	if len(userIDs) > 0 {
		if err := global.DB.Select("id", "username", "real_name", "avatar").
			Where("id IN ?", userIDs).
			Find(&users).Error; err != nil {
			return nil, err
		}
	}

	userMap := make(map[int64]model.SysUser, len(users))
	for _, user := range users {
		userMap[user.ID] = user
	}

	result := make([]GreenPointLeaderboardItem, 0, len(items))
	for index, item := range items {
		userID, err := memberToUserID(item.Member)
		if err != nil {
			continue
		}
		user := userMap[userID]
		result = append(result, GreenPointLeaderboardItem{
			UserID:   userID,
			Username: user.Username,
			RealName: user.RealName,
			Avatar:   user.Avatar,
			Rank:     index + 1,
			Points:   int(item.Score),
			Nickname: buildUserNickname(user),
		})
	}

	return result, nil
}

func memberToUserID(member interface{}) (int64, error) {
	switch value := member.(type) {
	case string:
		return strconv.ParseInt(value, 10, 64)
	case []byte:
		return strconv.ParseInt(string(value), 10, 64)
	default:
		return 0, fmt.Errorf("unsupported redis member type: %T", member)
	}
}

func buildUserNickname(user model.SysUser) string {
	if strings.TrimSpace(user.RealName) != "" {
		return user.RealName
	}
	if strings.TrimSpace(user.Username) != "" {
		return user.Username
	}
	return fmt.Sprintf("User-%d", user.ID)
}

func buildDataURLFromFile(fileHeader *multipart.FileHeader) (string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return "", err
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if strings.TrimSpace(contentType) == "" {
		contentType = "image/jpeg"
	}

	return "data:" + contentType + ";base64," + base64.StdEncoding.EncodeToString(data), nil
}
