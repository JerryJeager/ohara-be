package http

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/JerryJeager/ohara-be/internal/models"
	"github.com/JerryJeager/ohara-be/internal/service/documents"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DocumentController struct {
	serv documents.DocumentSv
}

func NewDocumentController(serv documents.DocumentSv) *DocumentController {
	return &DocumentController{serv: serv}
}

func (c *DocumentController) EmbedDocument(ctx *gin.Context) {
	var content models.Document
	if err := ctx.ShouldBindJSON(&content); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := c.serv.EmbedDocument(ctx, &content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *DocumentController) QueryDocument(ctx *gin.Context) {
	var websiteID WebsiteIDPP
	if err := ctx.ShouldBindUri(&websiteID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	websiteId := uuid.MustParse(websiteID.WebsiteID)
	var query models.Query
	if err := ctx.ShouldBindJSON(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	website, err := c.serv.GetWebsite(ctx, websiteId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	allowedHost := getHostname(website.Url)

	origin := ctx.GetHeader("Origin")

	// Fallback to Referer
	if origin == "" {
		origin = ctx.GetHeader("Referer")
	}

	if origin == "" {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "missing origin",
		})
		return
	}

	requestHost := getHostname(origin)

	isLocalhost := requestHost == "localhost" ||
		requestHost == "127.0.0.1"

	if isLocalhost {
		if !website.IsLocalDevEnabled {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "local development not allowed",
			})
			return
		}
	} else if requestHost != allowedHost {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "unauthorized domain",
		})
		return
	}

	response, err := c.serv.QueryWebsiteDocument(ctx, websiteId, &query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"response": response,
	})
}

func (c *DocumentController) ChunkDocument(ctx *gin.Context) {
	var chunkReq models.ChunkReq
	if err := ctx.ShouldBindJSON(&chunkReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	chunks := c.serv.ChunkDocument(ctx, &chunkReq)

	ctx.JSON(http.StatusOK, gin.H{
		"chunks": chunks,
	})
}

func getHostname(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}

	return strings.ToLower(u.Hostname())
}
