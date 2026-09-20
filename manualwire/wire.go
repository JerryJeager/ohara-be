package manualwire

import (
	"github.com/JerryJeager/ohara-be/config"
	"github.com/JerryJeager/ohara-be/internal/http"
	"github.com/JerryJeager/ohara-be/internal/service/documents"
	"github.com/JerryJeager/ohara-be/internal/service/users"
	"github.com/JerryJeager/ohara-be/internal/service/websites"
)

func GetUserRepository() *users.UserRepo {
	repo := config.GetSession()
	return users.NewUserRepo(repo)
}

func GetUserService(repo users.UserStore) *users.UserServ {
	return users.NewUserService(repo)
}

func GetUserController() *http.UserController {
	repo := GetUserRepository()
	service := GetUserService(repo)
	return http.NewUserController(service)
}

func GetWebsiteRepository() *websites.WebsiteRepo {
	repo := config.GetSession()
	return websites.NewWebsiteRepo(repo)
}

func GetWebsiteService(repo websites.WebsiteStore) *websites.WebsiteServ {
	return websites.NewWebsiteService(repo)
}

func GetWebsiteController() *http.WebsiteController {
	repo := GetWebsiteRepository()
	service := GetWebsiteService(repo)
	return http.NewWebsiteController(service)
}

func GetDocumentRepository() *documents.DocumentRepo {
	repo := config.GetSession()
	return documents.NewDocumentRepo(repo)
}

func GetDocumentService(repo documents.DocumentStore) *documents.DocumentServ {
	return documents.NewDocumentService(repo)
}

func GetDocumentController() *http.DocumentController {
	repo := GetDocumentRepository()
	service := GetDocumentService(repo)
	return http.NewDocumentController(service)
}
