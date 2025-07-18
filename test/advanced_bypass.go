package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/tebeka/selenium"
)

func main() {
	log.Println("Starting advanced Cloudflare bypass...")

	// Start ChromeDriver service
	service, err := selenium.NewChromeDriverService("/usr/local/bin/chromedriver", 4444)
	if err != nil {
		log.Fatalf("Failed to start ChromeDriver service: %v", err)
	}
	defer service.Stop()

	log.Println("ChromeDriver service started")

	// Configure Chrome capabilities with more stealth options
	caps := selenium.Capabilities{}
	caps["browserName"] = "chrome"
	caps["goog:chromeOptions"] = map[string]interface{}{
		"args": []string{
			"--headless",
			"--no-sandbox",
			"--disable-dev-shm-usage",
			"--disable-blink-features=AutomationControlled",
			"--disable-web-security",
			"--disable-features=VizDisplayCompositor",
			"--disable-background-timer-throttling",
			"--disable-backgrounding-occluded-windows",
			"--disable-renderer-backgrounding",
			"--disable-field-trial-config",
			"--disable-ipc-flooding-protection",
			"--disable-extensions",
			"--disable-plugins-discovery",
			"--disable-default-apps",
			"--disable-sync",
			"--disable-translate",
			"--hide-scrollbars",
			"--mute-audio",
			"--no-first-run",
			"--no-default-browser-check",
			"--disable-hang-monitor",
			"--disable-prompt-on-repost",
			"--disable-client-side-phishing-detection",
			"--disable-component-update",
			"--disable-domain-reliability",
			"--disable-features=TranslateUI",
			"--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		},
		"prefs": map[string]interface{}{
			"profile.default_content_setting_values.notifications": 2,
			"profile.managed_default_content_settings.images":      1,
		},
	}

	// Connect to WebDriver
	log.Println("Connecting to WebDriver...")
	wd, err := selenium.NewRemote(caps, "")
	if err != nil {
		log.Fatalf("Failed to connect to WebDriver: %v", err)
	}
	defer wd.Quit()

	// Set timeouts
	wd.SetImplicitWaitTimeout(30 * time.Second)
	wd.SetPageLoadTimeout(60 * time.Second)

	url := "https://novelbin.com/b/defiance-of-the-fall/prologuewelcome-to-the-multi-verse"
	log.Printf("Navigating to: %s", url)

	// Navigate to the URL
	if err := wd.Get(url); err != nil {
		log.Fatalf("Failed to navigate: %v", err)
	}

	// Wait for initial page load
	time.Sleep(5 * time.Second)

	// Advanced Cloudflare bypass logic
	log.Println("Starting Cloudflare bypass...")

	// Wait for Cloudflare challenge to complete
	for i := 0; i < 120; i++ { // Wait up to 2 minutes
		title, err := wd.Title()
		if err != nil {
			log.Printf("Error getting title: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		currentURL, err := wd.CurrentURL()
		if err != nil {
			log.Printf("Error getting current URL: %v", err)
		}

		log.Printf("Current title: %s", title)
		log.Printf("Current URL: %s", currentURL)

		// Check if we're past Cloudflare
		if !isCloudflarePage(title) && !isCloudflarePage(currentURL) {
			log.Println("✅ Cloudflare challenge appears to be complete!")
			break
		}

		log.Printf("⏳ Still on Cloudflare page (attempt %d/120)...", i+1)
		time.Sleep(2 * time.Second)
	}

	// Additional wait for content to load
	log.Println("Waiting for content to load...")
	time.Sleep(10 * time.Second)

	// Get the page source
	html, err := wd.PageSource()
	if err != nil {
		log.Fatalf("Failed to get page source: %v", err)
	}

	title, err := wd.Title()
	if err != nil {
		log.Printf("Error getting final title: %v", err)
		title = "Unknown"
	}

	log.Printf("Final title: %s", title)
	log.Printf("HTML length: %d characters", len(html))

	// Parse the content
	content := parseNovelContent(html, title, url)

	fmt.Printf("\n=== NOVEL CONTENT ===\n")
	fmt.Printf("Title: %s\n", content.Title)
	fmt.Printf("URL: %s\n", content.URL)
	fmt.Printf("Next URL: %s\n", content.NextURL)
	fmt.Printf("Content length: %d characters\n", len(content.Content))
	fmt.Printf("\n=== CONTENT PREVIEW ===\n")

	if len(content.Content) > 0 {
		preview := content.Content
		if len(preview) > 1000 {
			preview = preview[:1000] + "..."
		}
		fmt.Println(preview)
	} else {
		fmt.Println("No content found!")
		fmt.Printf("\n=== HTML PREVIEW ===\n")
		if len(html) > 500 {
			fmt.Println(html[:500] + "...")
		} else {
			fmt.Println(html)
		}
	}
}

func isCloudflarePage(text string) bool {
	text = strings.ToLower(text)
	indicators := []string{
		"just a moment",
		"checking your browser",
		"cloudflare",
		"ray id:",
		"verification",
		"challenge",
		"cf_chl_",
		"enable javascript",
	}

	for _, indicator := range indicators {
		if strings.Contains(text, indicator) {
			return true
		}
	}
	return false
}

type NovelContent struct {
	Title   string
	Content string
	URL     string
	NextURL string
}

func parseNovelContent(html, title, url string) *NovelContent {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Printf("Failed to parse HTML: %v", err)
		return &NovelContent{
			Title:   title,
			URL:     url,
			Content: "Failed to parse HTML",
			NextURL: "",
		}
	}

	content := extractContent(doc, "#chr-content")
	nextURL := extractLink(doc, "#next_chap")

	return &NovelContent{
		Title:   title,
		Content: content,
		URL:     url,
		NextURL: nextURL,
	}
}

func extractContent(doc *goquery.Document, selector string) string {
	selection := doc.Find(selector)
	if selection.Length() == 0 {
		return ""
	}

	// Remove unwanted elements
	selection.Find("script, style, nav, header, footer, .ads, .advertisement, .sidebar, .comments, .cloudflare").Remove()

	// Get text content
	text := selection.Text()

	// Clean up the text
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "\n\n\n", "\n\n")
	text = strings.ReplaceAll(text, "\t", " ")

	return text
}

func extractLink(doc *goquery.Document, selector string) string {
	selection := doc.Find(selector)
	if selection.Length() == 0 {
		return ""
	}

	// Remove unwanted elements
	selection.Find("script, style, nav, header, footer, .ads, .advertisement, .sidebar, .comments, .cloudflare").Remove()

	// Get text content
	link, exists := selection.Attr("href")
	if !exists {
		return ""
	}

	return link
}
