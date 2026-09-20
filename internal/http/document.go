package http

import (
	"net/http"

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
