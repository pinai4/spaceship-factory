package notification

import (
	"embed"
	"text/template"

	"github.com/pinai4/spaceship-factory/notification/internal/client"
	def "github.com/pinai4/spaceship-factory/notification/internal/service"
)

var _ def.NotificationService = (*service)(nil)

//go:embed templates/*.tmpl
var templatesFS embed.FS

type service struct {
	telegramClient        client.TelegramClient
	userClient            client.UserClient
	orderPaidTemplate     *template.Template
	shipAssembledTemplate *template.Template
}

func NewService(telegramClient client.TelegramClient, userClient client.UserClient) *service {
	orderPaidTemplate := template.Must(template.ParseFS(templatesFS, "templates/order_paid_notification.tmpl"))
	shipAssembledTemplate := template.Must(template.ParseFS(templatesFS, "templates/ship_assembled_notification.tmpl"))

	return &service{
		telegramClient:        telegramClient,
		userClient:            userClient,
		orderPaidTemplate:     orderPaidTemplate,
		shipAssembledTemplate: shipAssembledTemplate,
	}
}
