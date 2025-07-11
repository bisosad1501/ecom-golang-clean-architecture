package database

import (
	"context"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userSessionRepository struct {
	db *gorm.DB
}

// NewUserSessionRepository creates a new user session repository
func NewUserSessionRepository(db *gorm.DB) repositories.UserSessionRepository {
	return &userSessionRepository{db: db}
}

// Create creates a new user session
func (r *userSessionRepository) Create(ctx context.Context, session *entities.UserSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

// GetByID retrieves a user session by ID
func (r *userSessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.UserSession, error) {
	var session entities.UserSession
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&session).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrUserNotFound
		}
		return nil, err
	}
	return &session, nil
}

// GetByToken retrieves a user session by token
func (r *userSessionRepository) GetByToken(ctx context.Context, token string) (*entities.UserSession, error) {
	var session entities.UserSession
	err := r.db.WithContext(ctx).Where("session_token = ? AND is_active = ? AND expires_at > ?", token, true, time.Now()).First(&session).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrUserNotFound
		}
		return nil, err
	}
	return &session, nil
}

// Update updates an existing user session
func (r *userSessionRepository) Update(ctx context.Context, session *entities.UserSession) error {
	return r.db.WithContext(ctx).Save(session).Error
}

// Delete deletes a user session by ID
func (r *userSessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&entities.UserSession{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entities.ErrUserNotFound
	}
	return nil
}

// GetActiveSessionsByUserID retrieves active sessions for a user
func (r *userSessionRepository) GetActiveSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.UserSession, error) {
	var sessions []*entities.UserSession
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ? AND expires_at > ?", userID, true, time.Now()).
		Order("last_activity DESC").
		Find(&sessions).Error
	return sessions, err
}

// GetSessionsByUserID retrieves sessions for a user with pagination
func (r *userSessionRepository) GetSessionsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.UserSession, error) {
	var sessions []*entities.UserSession
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&sessions).Error
	return sessions, err
}

// InvalidateUserSessions invalidates all sessions for a user
func (r *userSessionRepository) InvalidateUserSessions(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entities.UserSession{}).
		Where("user_id = ?", userID).
		Update("is_active", false).Error
}

// InvalidateSessionByToken invalidates a session by token
func (r *userSessionRepository) InvalidateSessionByToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).
		Model(&entities.UserSession{}).
		Where("session_token = ?", token).
		Update("is_active", false).Error
}

// DeleteExpiredSessions deletes expired sessions
func (r *userSessionRepository) DeleteExpiredSessions(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&entities.UserSession{}).Error
}

// DeleteInactiveSessions deletes inactive sessions
func (r *userSessionRepository) DeleteInactiveSessions(ctx context.Context, inactiveThreshold time.Duration) error {
	cutoff := time.Now().Add(-inactiveThreshold)
	return r.db.WithContext(ctx).
		Where("last_activity < ? OR is_active = ?", cutoff, false).
		Delete(&entities.UserSession{}).Error
}

// CountActiveSessionsByUserID counts active sessions for a user
func (r *userSessionRepository) CountActiveSessionsByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.UserSession{}).
		Where("user_id = ? AND is_active = ? AND expires_at > ?", userID, true, time.Now()).
		Count(&count).Error
	return count, err
}

type userLoginHistoryRepository struct {
	db *gorm.DB
}

// NewUserLoginHistoryRepository creates a new user login history repository
func NewUserLoginHistoryRepository(db *gorm.DB) repositories.UserLoginHistoryRepository {
	return &userLoginHistoryRepository{db: db}
}

// Create creates a new login history record
func (r *userLoginHistoryRepository) Create(ctx context.Context, history *entities.UserLoginHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

// GetByUserID retrieves login history for a user
func (r *userLoginHistoryRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.UserLoginHistory, error) {
	var history []*entities.UserLoginHistory
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&history).Error
	return history, err
}

// GetFailedLoginAttempts retrieves failed login attempts since a specific time
func (r *userLoginHistoryRepository) GetFailedLoginAttempts(ctx context.Context, userID uuid.UUID, since time.Time) ([]*entities.UserLoginHistory, error) {
	var history []*entities.UserLoginHistory
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND success = ? AND created_at > ?", userID, false, since).
		Order("created_at DESC").
		Find(&history).Error
	return history, err
}

// CountLoginAttempts counts login attempts since a specific time
func (r *userLoginHistoryRepository) CountLoginAttempts(ctx context.Context, userID uuid.UUID, since time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.UserLoginHistory{}).
		Where("user_id = ? AND created_at > ?", userID, since).
		Count(&count).Error
	return count, err
}

// CountFailedAttempts counts failed login attempts since a specific time
func (r *userLoginHistoryRepository) CountFailedAttempts(ctx context.Context, userID uuid.UUID, since time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.UserLoginHistory{}).
		Where("user_id = ? AND success = ? AND created_at > ?", userID, false, since).
		Count(&count).Error
	return count, err
}

// DeleteOldHistory deletes old login history
func (r *userLoginHistoryRepository) DeleteOldHistory(ctx context.Context, olderThan time.Time) error {
	return r.db.WithContext(ctx).
		Where("created_at < ?", olderThan).
		Delete(&entities.UserLoginHistory{}).Error
}

type userActivityRepository struct {
	db *gorm.DB
}

// NewUserActivityRepository creates a new user activity repository
func NewUserActivityRepository(db *gorm.DB) repositories.UserActivityRepository {
	return &userActivityRepository{db: db}
}

// Create creates a new user activity
func (r *userActivityRepository) Create(ctx context.Context, activity *entities.UserActivity) error {
	return r.db.WithContext(ctx).Create(activity).Error
}

// GetByID retrieves a user activity by ID
func (r *userActivityRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.UserActivity, error) {
	var activity entities.UserActivity
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&activity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrUserNotFound
		}
		return nil, err
	}
	return &activity, nil
}

// GetByUserID retrieves activities for a user
func (r *userActivityRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.UserActivity, error) {
	var activities []*entities.UserActivity
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&activities).Error
	return activities, err
}

// GetByUserIDAndType retrieves activities for a user by type
func (r *userActivityRepository) GetByUserIDAndType(ctx context.Context, userID uuid.UUID, activityType entities.ActivityType, limit, offset int) ([]*entities.UserActivity, error) {
	var activities []*entities.UserActivity
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ?", userID, activityType).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&activities).Error
	return activities, err
}

// GetRecentActivity retrieves recent activities for a user
func (r *userActivityRepository) GetRecentActivity(ctx context.Context, userID uuid.UUID, since time.Time) ([]*entities.UserActivity, error) {
	var activities []*entities.UserActivity
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND created_at > ?", userID, since).
		Order("created_at DESC").
		Find(&activities).Error
	return activities, err
}

// GetActivityStats retrieves activity statistics for a user
func (r *userActivityRepository) GetActivityStats(ctx context.Context, userID uuid.UUID, dateFrom, dateTo time.Time) (map[entities.ActivityType]int64, error) {
	var results []struct {
		Type  entities.ActivityType `json:"type"`
		Count int64                 `json:"count"`
	}

	err := r.db.WithContext(ctx).
		Model(&entities.UserActivity{}).
		Select("type, COUNT(*) as count").
		Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, dateFrom, dateTo).
		Group("type").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	stats := make(map[entities.ActivityType]int64)
	for _, result := range results {
		stats[result.Type] = result.Count
	}

	return stats, nil
}

// GetMostActiveUsers retrieves most active users
func (r *userActivityRepository) GetMostActiveUsers(ctx context.Context, limit int, dateFrom, dateTo time.Time) ([]*entities.User, error) {
	var users []*entities.User
	err := r.db.WithContext(ctx).
		Table("users").
		Select("users.*").
		Joins("JOIN user_activities ON users.id = user_activities.user_id").
		Where("user_activities.created_at BETWEEN ? AND ?", dateFrom, dateTo).
		Group("users.id").
		Order("COUNT(user_activities.id) DESC").
		Limit(limit).
		Find(&users).Error
	return users, err
}

// DeleteOldActivities deletes old activities
func (r *userActivityRepository) DeleteOldActivities(ctx context.Context, olderThan time.Time) error {
	return r.db.WithContext(ctx).
		Where("created_at < ?", olderThan).
		Delete(&entities.UserActivity{}).Error
}

// DeleteByUserID deletes all activities for a user
func (r *userActivityRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&entities.UserActivity{}).Error
}

type userPreferencesRepository struct {
	db *gorm.DB
}

// NewUserPreferencesRepository creates a new user preferences repository
func NewUserPreferencesRepository(db *gorm.DB) repositories.UserPreferencesRepository {
	return &userPreferencesRepository{db: db}
}

// Create creates new user preferences
func (r *userPreferencesRepository) Create(ctx context.Context, preferences *entities.UserPreferences) error {
	return r.db.WithContext(ctx).Create(preferences).Error
}

// GetByUserID retrieves user preferences by user ID
func (r *userPreferencesRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entities.UserPreferences, error) {
	var preferences entities.UserPreferences
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&preferences).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrUserNotFound
		}
		return nil, err
	}
	return &preferences, nil
}

// Update updates user preferences
func (r *userPreferencesRepository) Update(ctx context.Context, preferences *entities.UserPreferences) error {
	return r.db.WithContext(ctx).Save(preferences).Error
}

// Delete deletes user preferences
func (r *userPreferencesRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&entities.UserPreferences{}).Error
}

// UpdateTheme updates user theme preference
func (r *userPreferencesRepository) UpdateTheme(ctx context.Context, userID uuid.UUID, theme string) error {
	return r.db.WithContext(ctx).Model(&entities.UserPreferences{}).
		Where("user_id = ?", userID).
		Update("theme", theme).Error
}

// UpdateLanguage updates user language preference
func (r *userPreferencesRepository) UpdateLanguage(ctx context.Context, userID uuid.UUID, language string) error {
	return r.db.WithContext(ctx).Model(&entities.UserPreferences{}).
		Where("user_id = ?", userID).
		Update("language", language).Error
}

// UpdateCurrency updates user currency preference
func (r *userPreferencesRepository) UpdateCurrency(ctx context.Context, userID uuid.UUID, currency string) error {
	return r.db.WithContext(ctx).Model(&entities.UserPreferences{}).
		Where("user_id = ?", userID).
		Update("currency", currency).Error
}

// UpdateTimezone updates user timezone preference
func (r *userPreferencesRepository) UpdateTimezone(ctx context.Context, userID uuid.UUID, timezone string) error {
	return r.db.WithContext(ctx).Model(&entities.UserPreferences{}).
		Where("user_id = ?", userID).
		Update("timezone", timezone).Error
}

// UpdateNotificationSettings updates notification settings
func (r *userPreferencesRepository) UpdateNotificationSettings(ctx context.Context, userID uuid.UUID, settings map[string]bool) error {
	updates := make(map[string]interface{})
	for key, value := range settings {
		updates[key] = value
	}
	updates["updated_at"] = time.Now()

	return r.db.WithContext(ctx).Model(&entities.UserPreferences{}).
		Where("user_id = ?", userID).
		Updates(updates).Error
}

// UpdatePrivacySettings updates privacy settings
func (r *userPreferencesRepository) UpdatePrivacySettings(ctx context.Context, userID uuid.UUID, settings map[string]interface{}) error {
	settings["updated_at"] = time.Now()
	return r.db.WithContext(ctx).Model(&entities.UserPreferences{}).
		Where("user_id = ?", userID).
		Updates(settings).Error
}

// UpdateShoppingSettings updates shopping settings
func (r *userPreferencesRepository) UpdateShoppingSettings(ctx context.Context, userID uuid.UUID, settings map[string]interface{}) error {
	settings["updated_at"] = time.Now()
	return r.db.WithContext(ctx).Model(&entities.UserPreferences{}).
		Where("user_id = ?", userID).
		Updates(settings).Error
}

// GetPreferencesByTheme gets preferences by theme
func (r *userPreferencesRepository) GetPreferencesByTheme(ctx context.Context, theme string, limit, offset int) ([]*entities.UserPreferences, error) {
	var preferences []*entities.UserPreferences
	err := r.db.WithContext(ctx).
		Where("theme = ?", theme).
		Limit(limit).
		Offset(offset).
		Find(&preferences).Error
	return preferences, err
}

// GetPreferencesByLanguage gets preferences by language
func (r *userPreferencesRepository) GetPreferencesByLanguage(ctx context.Context, language string, limit, offset int) ([]*entities.UserPreferences, error) {
	var preferences []*entities.UserPreferences
	err := r.db.WithContext(ctx).
		Where("language = ?", language).
		Limit(limit).
		Offset(offset).
		Find(&preferences).Error
	return preferences, err
}

type userVerificationRepository struct {
	db *gorm.DB
}

// NewUserVerificationRepository creates a new user verification repository
func NewUserVerificationRepository(db *gorm.DB) repositories.UserVerificationRepository {
	return &userVerificationRepository{db: db}
}

// Create creates a new user verification
func (r *userVerificationRepository) Create(ctx context.Context, verification *entities.UserVerification) error {
	return r.db.WithContext(ctx).Create(verification).Error
}

// GetByID retrieves a user verification by ID
func (r *userVerificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.UserVerification, error) {
	var verification entities.UserVerification
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&verification).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrUserNotFound
		}
		return nil, err
	}
	return &verification, nil
}

// GetByToken retrieves a user verification by token
func (r *userVerificationRepository) GetByToken(ctx context.Context, token string) (*entities.UserVerification, error) {
	var verification entities.UserVerification
	err := r.db.WithContext(ctx).Where("token = ? AND is_verified = ? AND expires_at > ?", token, false, time.Now()).First(&verification).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrUserNotFound
		}
		return nil, err
	}
	return &verification, nil
}

// GetByUserIDAndType retrieves a user verification by user ID and type
func (r *userVerificationRepository) GetByUserIDAndType(ctx context.Context, userID uuid.UUID, verificationType string) (*entities.UserVerification, error) {
	var verification entities.UserVerification
	err := r.db.WithContext(ctx).Where("user_id = ? AND type = ? AND is_verified = ? AND expires_at > ?", userID, verificationType, false, time.Now()).First(&verification).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrUserNotFound
		}
		return nil, err
	}
	return &verification, nil
}

// Update updates a user verification
func (r *userVerificationRepository) Update(ctx context.Context, verification *entities.UserVerification) error {
	return r.db.WithContext(ctx).Save(verification).Error
}

// Delete deletes a user verification
func (r *userVerificationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.UserVerification{}, id).Error
}

// GetActiveVerifications retrieves active verifications for a user
func (r *userVerificationRepository) GetActiveVerifications(ctx context.Context, userID uuid.UUID) ([]*entities.UserVerification, error) {
	var verifications []*entities.UserVerification
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_verified = ? AND expires_at > ?", userID, false, time.Now()).
		Order("created_at DESC").
		Find(&verifications).Error
	return verifications, err
}

// GetByCode retrieves a verification by code and type
func (r *userVerificationRepository) GetByCode(ctx context.Context, code string, verificationType string) (*entities.UserVerification, error) {
	var verification entities.UserVerification
	err := r.db.WithContext(ctx).Where("code = ? AND type = ? AND is_verified = ? AND expires_at > ?", code, verificationType, false, time.Now()).First(&verification).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrUserNotFound
		}
		return nil, err
	}
	return &verification, nil
}

// MarkAsVerified marks a verification as verified
func (r *userVerificationRepository) MarkAsVerified(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&entities.UserVerification{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_verified": true,
			"verified_at": now,
			"updated_at":  now,
		}).Error
}

// IncrementAttempt increments the attempt count
func (r *userVerificationRepository) IncrementAttempt(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entities.UserVerification{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"attempt_count": gorm.Expr("attempt_count + 1"),
			"updated_at":    time.Now(),
		}).Error
}

// DeleteExpiredVerifications deletes expired verifications
func (r *userVerificationRepository) DeleteExpiredVerifications(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&entities.UserVerification{}).Error
}

// DeleteByUserID deletes all verifications for a user
func (r *userVerificationRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&entities.UserVerification{}).Error
}

// CountVerificationsByType counts verifications by type in date range
func (r *userVerificationRepository) CountVerificationsByType(ctx context.Context, verificationType string, dateFrom, dateTo time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.UserVerification{}).
		Where("type = ? AND created_at BETWEEN ? AND ?", verificationType, dateFrom, dateTo).
		Count(&count).Error
	return count, err
}

// GetFailedVerifications retrieves failed verifications for a user
func (r *userVerificationRepository) GetFailedVerifications(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.UserVerification, error) {
	var verifications []*entities.UserVerification
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND attempt_count >= max_attempts", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&verifications).Error
	return verifications, err
}
