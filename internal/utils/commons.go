package utils

import (
	"ORDI/internal/models"
	"ORDI/internal/repositories"
	"context"
	"strconv"
)

func GetNotificationCount(ctx context.Context, notificationRepository repositories.Repository[models.Notification], email string) (string, error) {

	messages, err := notificationRepository.FindAllByField(ctx, "user_email", email)
	if err != nil {
		return "", err
	}
	notificationCount := 0
	for _, message := range messages {
		if !message.IsRead {
			notificationCount++
		}
	}
	return strconv.Itoa(notificationCount), nil
}
