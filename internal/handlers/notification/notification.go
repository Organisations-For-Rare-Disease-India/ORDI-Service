package notification

import (
	"ORDI/cmd/web"
	"ORDI/internal/handlers/token"
	"ORDI/internal/models"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

func (n *notificationHandler) ShowNotifications(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, err := token.ValidateJWT(w, r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	fmt.Println(claims.UserId)
	notifications, err := n.notificationRepository.FindAllByField(ctx, "user_id", claims.UserId)

	var viewNotifications []models.ViewNotification
	for _, notification := range notifications {
		if !notification.IsRead {
			noti := models.ViewNotification{
				SentTime: notification.SentTime.Format("02 Jan 2006, 03:04 PM"),
				Message:  notification.Message,
			}
			viewNotifications = append(viewNotifications, noti)
		}

	}
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	templ.Handler(web.NotificationsPage(viewNotifications)).ServeHTTP(w, r)
}
