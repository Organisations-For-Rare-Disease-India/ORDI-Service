package notification

import (
	"ORDI/internal/models"
	"ORDI/internal/repositories"
	"net/http"
)

type Notification interface {
	ShowNotifications(http.ResponseWriter, *http.Request)
}

type notificationHandler struct {
	notificationRepository repositories.Repository[models.Notification]
}
type NotificationHandlerConfig struct {
	NotificationRepository repositories.Repository[models.Notification]
}

func NewNotificationHandler(config NotificationHandlerConfig) Notification {
	return &notificationHandler{
		notificationRepository: config.NotificationRepository,
	}
}
