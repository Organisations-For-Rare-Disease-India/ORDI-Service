package notification

import (
	"ORDI/cmd/web"
	"ORDI/internal/handlers/token"
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
	notifications, err := n.notificationRepository.FindAllByField(ctx, "user_email", claims.Email)

	for _, notification := range notifications {
		notification.IsRead = true
		err = n.notificationRepository.Save(ctx, &notification)
	}
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	templ.Handler(web.NotificationsPage(notifications)).ServeHTTP(w, r)

}
