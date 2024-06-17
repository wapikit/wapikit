package contact_list_service

import (
	"github.com/sarthakjdev/wapikit/api/services"
	"github.com/sarthakjdev/wapikit/internal/interfaces"
)

type ContactListService struct {
	services.BaseService `json:"-,inline"`
}

func NewContactListService() *ContactListService {
	return &ContactListService{
		BaseService: services.BaseService{
			Name:        "Contact List Service",
			RestApiPath: "/api/contact-list",
			Routes: []interfaces.Route{
				{},
			},
		},
	}
}
