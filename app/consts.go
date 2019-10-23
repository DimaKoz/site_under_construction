package app

// I use these constants inside this project
const (
	LogPath = "log.txt"

	PathPatternRoot         = "/"
	PathPatternNotFound     = "/404"
	PathPatternUnknownError = "/500"
	PathPatternFavicon      = "/favicon.ico"
	PathPatternWoff2        = "/assets/woff2/"
	PathPatternCss          = "/assets/css/"
	PathPatternJs           = "/assets/script/"
	PathPatternImage        = "/assets/image/"
	Html500                 = "./html/error_500_page.html"
	Html404                 = "./html/error_404_page.html"
	HtmlUnderConstruction   = "./html/under_construction.html"

	HeaderKeyEtag            = "Etag"
	HeaderKeyCacheControl    = "Cache-Control"
	HeaderKeyIfNoneMatch     = "If-None-Match"
	HeaderValueCacheControl  = "max-age=3600"
	HeaderKeyContentType     = "Content-Type"
	HeaderKeyAcceptEncoding  = "Accept-Encoding"
	HeaderValueEncoding      = "gzip"
	HeaderKeyContentLength   = "Content-Length"
	HeaderKeyContentEncoding = "Content-Encoding"

	FaviconData = "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"
)
