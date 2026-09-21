package http

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"strings"

	"github.com/JerryJeager/ohara-be/internal/models"
	"github.com/JerryJeager/ohara-be/internal/service/documents"
	"github.com/JerryJeager/ohara-be/internal/service/websites"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/gocolly/colly/v2"
	"github.com/google/uuid"
	"golang.org/x/net/html"
)

type WebsiteController struct {
	serv websites.WebsiteSv
}

func NewWebsiteController(serv websites.WebsiteSv) *WebsiteController {
	return &WebsiteController{serv: serv}
}

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_3_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.3 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:123.0) Gecko/20100101 Firefox/123.0",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/122.0.6261.89 Mobile/15E148 Safari/604.1",
}

var blockElements = map[string]bool{
	"p": true, "div": true, "br": true, "li": true,
	"h1": true, "h2": true, "h3": true, "h4": true, "h5": true, "h6": true,
	"tr": true, "table": true, "section": true, "article": true,
	"blockquote": true, "pre": true, "ul": true, "ol": true, "header": true,
	"footer": true, "main": true,
}

var tagsToSkip = map[string]bool{
	"script": true, "style": true, "noscript": true, "svg": true,
	"iframe": true, "form": true, "button": true, "select": true,
	"nav": true, "footer": true, "head": true,
}

// skipExtensions are link targets we never want to treat as crawlable pages.
var skipExtensions = []string{
	".pdf", ".jpg", ".jpeg", ".png", ".gif", ".svg", ".webp", ".ico",
	".zip", ".css", ".js", ".mp4", ".mp3", ".woff", ".woff2", ".xml",
}

const (
	maxPagesPerSite = 20
	maxCrawlDepth   = 2
	crawlDelay      = 500 * time.Millisecond
	spaThreshold    = 200
)

func (c *WebsiteController) CreateWebsite(ctx *gin.Context) {
	var website models.Website
	if err := ctx.ShouldBindJSON(&website); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "invalid request body",
		})
		return
	}

	websiteID, err := c.serv.CreateWebsite(ctx, &website)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	go saveContentEmbedding(websiteID, website.Url)

	ctx.JSON(http.StatusAccepted, gin.H{
		"website_id": websiteID,
	})
}

// ---------- fetching ----------

func fetchWithColly(target string) (string, error) {
	c1 := colly.NewCollector()
	c1.SetRequestTimeout(15 * time.Second)

	var body string
	var visitErr error

	c1.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])
	})

	c1.OnResponse(func(r *colly.Response) {
		body = string(r.Body)
	})

	c1.OnError(func(r *colly.Response, err error) {
		visitErr = err
	})

	if err := c1.Visit(target); err != nil {
		return "", err
	}
	c1.Wait()

	if visitErr != nil {
		return "", visitErr
	}
	return body, nil
}

// rodSession wraps a single headless-Chrome instance so a whole crawl can
// share one browser process instead of launching a new one per page.
type rodSession struct {
	browser *rod.Browser
}

func newRodSession() (*rodSession, error) {
	launcherURL, err := launcher.New().
		Headless(true).
		Set("disable-gpu", "").
		Set("no-sandbox", "").
		Launch()
	if err != nil {
		return nil, err
	}

	browser := rod.New().ControlURL(launcherURL)
	if err := browser.Connect(); err != nil {
		return nil, err
	}
	return &rodSession{browser: browser}, nil
}

func (s *rodSession) fetch(target string) (string, error) {
	page, err := s.browser.Timeout(20 * time.Second).Page(proto.TargetCreateTarget{URL: target})
	if err != nil {
		return "", err
	}
	defer page.Close()

	if err := page.WaitStable(1 * time.Second); err != nil {
		return "", fmt.Errorf("page did not stabilize: %w", err)
	}

	return page.HTML()
}

func (s *rodSession) Close() {
	if s == nil || s.browser == nil {
		return
	}
	_ = s.browser.Close()
}

// fetchRawWithRetry tries colly first, falls back to rod (lazily starting a
// shared session on first need) if the result looks like an unrendered SPA
// shell, and retries the whole sequence up to maxAttempts times. Returns the
// raw HTML (not extracted text) so callers can both extract content AND
// discover outgoing links from it.
func fetchRawWithRetry(target string, maxAttempts int, session **rodSession) (string, error) {
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		rawHTML, collyErr := fetchWithColly(target)
		content := extractCleanText(rawHTML)

		if collyErr == nil && len(strings.TrimSpace(content)) >= spaThreshold {
			return rawHTML, nil
		}

		log.Printf("attempt %d/%d: colly insufficient for %s (err: %v), trying rod",
			attempt, maxAttempts, target, collyErr)

		if *session == nil {
			s, err := newRodSession()
			if err != nil {
				lastErr = fmt.Errorf("failed to start rod session: %w", err)
				if attempt < maxAttempts {
					time.Sleep(time.Duration(attempt) * 3 * time.Second)
				}
				continue
			}
			*session = s
		}

		renderedHTML, rodErr := (*session).fetch(target)
		if rodErr == nil {
			if strings.TrimSpace(extractCleanText(renderedHTML)) != "" {
				return renderedHTML, nil
			}
			rodErr = fmt.Errorf("rod returned no extractable content")
		}

		lastErr = fmt.Errorf("attempt %d: colly err: %v, rod err: %v", attempt, collyErr, rodErr)

		if attempt < maxAttempts {
			backoff := time.Duration(attempt) * 3 * time.Second
			log.Printf("attempt %d failed for %s, retrying in %v", attempt, target, backoff)
			time.Sleep(backoff)
		}
	}

	return "", lastErr
}

// ---------- crawling ----------

type pageResult struct {
	url     string
	title   string
	content string
}

// normalizeURL strips fragments and trailing slashes so "/about" and
// "/about/" (and "/about#team") are treated as the same page.
func normalizeURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	u.Fragment = ""
	u.Path = strings.TrimSuffix(u.Path, "/")
	return u.String()
}

func shouldSkipLink(u *url.URL) bool {
	lower := strings.ToLower(u.Path)
	for _, ext := range skipExtensions {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}

// discoverLinks pulls same-host <a href> links out of raw HTML (works for
// both colly's static HTML and rod's rendered HTML).
func discoverLinks(rawHTML, baseURL string) []string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rawHTML))
	if err != nil {
		return nil
	}
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil
	}

	seen := map[string]bool{}
	var links []string

	doc.Find("a[href]").Each(func(_ int, sel *goquery.Selection) {
		href, exists := sel.Attr("href")
		if !exists {
			return
		}
		href = strings.TrimSpace(href)
		if href == "" || strings.HasPrefix(href, "#") ||
			strings.HasPrefix(href, "mailto:") || strings.HasPrefix(href, "tel:") ||
			strings.HasPrefix(href, "javascript:") {
			return
		}

		resolved, err := base.Parse(href)
		if err != nil || resolved.Host != base.Host || shouldSkipLink(resolved) {
			return
		}
		resolved.Fragment = ""

		normalized := normalizeURL(resolved.String())
		if !seen[normalized] {
			seen[normalized] = true
			links = append(links, normalized)
		}
	})

	return links
}

func extractTitle(rawHTML string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rawHTML))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(doc.Find("title").First().Text())
}

// crawlSite does a breadth-first crawl starting from startURL, bounded by
// maxPagesPerSite and maxCrawlDepth, sharing one rod session (started lazily,
// only if some page actually needs JS rendering) across the whole crawl.
func crawlSite(startURL string) []pageResult {
	var rodSess *rodSession
	defer rodSess.Close()

	type queueItem struct {
		url   string
		depth int
	}

	visited := map[string]bool{}
	queue := []queueItem{{url: startURL, depth: 0}}
	var results []pageResult

	for len(queue) > 0 && len(results) < maxPagesPerSite {
		item := queue[0]
		queue = queue[1:]

		norm := normalizeURL(item.url)
		if visited[norm] {
			continue
		}
		visited[norm] = true

		rawHTML, err := fetchRawWithRetry(item.url, 2, &rodSess)
		if err != nil {
			log.Printf("failed to fetch %s: %v", item.url, err)
			continue
		}

		content := extractCleanText(rawHTML)
		if strings.TrimSpace(content) == "" {
			continue
		}

		results = append(results, pageResult{
			url:     item.url,
			title:   extractTitle(rawHTML),
			content: content,
		})

		if item.depth < maxCrawlDepth {
			for _, link := range discoverLinks(rawHTML, item.url) {
				if !visited[normalizeURL(link)] {
					queue = append(queue, queueItem{url: link, depth: item.depth + 1})
				}
			}
		}

		time.Sleep(crawlDelay)
	}

	return results
}

// ---------- text extraction ----------

func extractCleanText(rawHTML string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rawHTML))
	if err != nil {
		return ""
	}

	doc.Find("script, style, noscript, svg, iframe, nav, footer, header, form, [aria-hidden='true']").Remove()

	selectors := []string{"main", "article", "[role='main']", "#content", ".content", "body"}
	for _, sel := range selectors {
		if selection := doc.Find(sel).First(); selection.Length() > 0 {
			if node := selection.Get(0); node != nil {
				text := walkAndExtract(node)
				if strings.TrimSpace(text) != "" {
					return normalizeWhitespace(text)
				}
			}
		}
	}
	return ""
}

func walkAndExtract(n *html.Node) string {
	var sb strings.Builder
	var walk func(*html.Node)

	walk = func(node *html.Node) {
		if node.Type == html.ElementNode && tagsToSkip[node.Data] {
			return
		}
		if node.Type == html.TextNode {
			text := strings.TrimSpace(node.Data)
			if text != "" {
				sb.WriteString(text)
				sb.WriteString(" ")
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
		if node.Type == html.ElementNode && blockElements[node.Data] {
			sb.WriteString("\n")
		}
	}

	walk(n)
	return sb.String()
}

var (
	reMultiSpace   = regexp.MustCompile(`[ \t]+`)
	reMultiNewline = regexp.MustCompile(`\n{3,}`)
)

func normalizeWhitespace(s string) string {
	s = reMultiSpace.ReplaceAllString(s, " ")
	s = reMultiNewline.ReplaceAllString(s, "\n\n")
	return strings.TrimSpace(s)
}

// ---------- orchestration ----------

func saveContentEmbedding(websiteID uuid.UUID, websiteUrl string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic in saveContentEmbedding for %s: %v", websiteID, r)
			_ = websites.UpdateWebsiteStatus(websiteID, "failed")
		}
	}()

	pages := crawlSite(websiteUrl)
	if len(pages) == 0 {
		log.Printf("no pages successfully crawled for website %s", websiteID)
		_ = websites.UpdateWebsiteStatus(websiteID, "failed")
		return
	}

	var embedErrCount int
	for _, page := range pages {
		if err := documents.ChunkAndEmbedDocument(page.content, page.title, page.url, websiteID); err != nil {
			log.Printf("failed to embed page %s for website %s: %v", page.url, websiteID, err)
			embedErrCount++
			continue
		}
	}

	if embedErrCount == len(pages) {
		_ = websites.UpdateWebsiteStatus(websiteID, "failed")
		log.Printf("all %d pages failed to embed for website %s", len(pages), websiteID)
		return
	}

	_ = websites.UpdateWebsiteStatus(websiteID, "success")
	log.Printf("successfully crawled, chunked and embedded %d/%d pages for website %s",
		len(pages)-embedErrCount, len(pages), websiteID)
}

func (c *WebsiteController) GetWebsite(ctx *gin.Context) {
	userId, err := GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	userID, err := uuid.Parse(userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	website, err := c.serv.GetWebsite(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, website)
}

func (c *WebsiteController) GetIndexedPages(ctx *gin.Context) {
	var websiteID WebsiteIDPP
	if err := ctx.ShouldBindUri(&websiteID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	indexedPages, err := c.serv.GetIndexedPages(ctx, uuid.MustParse(websiteID.WebsiteID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"indexed_pages": indexedPages,
	})

}

func (c *WebsiteController) UpdateLocalDev(ctx *gin.Context) {
	var websiteID WebsiteIDPP
	if err := ctx.ShouldBindUri(&websiteID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	var enabled models.IsLocalDevEnabled
	if err := ctx.ShouldBindJSON(&enabled); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if err := c.serv.UpdateLocalDev(ctx, &enabled, uuid.MustParse(websiteID.WebsiteID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *WebsiteController) DeleteWebsite(ctx *gin.Context) {
	var websiteID WebsiteIDPP
	if err := ctx.ShouldBindUri(&websiteID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	err := c.serv.DeleteWebsite(ctx, uuid.MustParse(websiteID.WebsiteID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}
