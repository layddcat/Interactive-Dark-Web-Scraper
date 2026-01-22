package services

import (
	"interactive-scraper/database"
	"regexp"
	"strings"
)

func ProcessRawData(sourceName, url, content, date string, isActive bool) {
	cleanText := cleanHTMLContent(content)
	title := extractHTMLTitle(content)
	
	if title == "" {
		words := strings.Fields(cleanText)
		if len(words) > 6 {
			title = strings.Join(words[:6], " ") + "..."
		} else {
			title = url
		}
	}

	score, category := AnalyzeIntelligence(cleanText, isActive)

	data := database.IntelligenceData{
		Title:            title,
		SourceName:       sourceName,
		SourceURL:        url,
		RawContent:       content,
		Category:         category,
		CriticalityScore: score,
		PublishedAt:      date,
		IsActive:         isActive && score > 0, 
	}
	database.DB.Create(&data)
}

func cleanHTMLContent(html string) string {
	reStyle := regexp.MustCompile(`(?is)<style.*?>.*?</style>`)
	html = reStyle.ReplaceAllString(html, " ")
	reScript := regexp.MustCompile(`(?is)<script.*?>.*?</script>`)
	html = reScript.ReplaceAllString(html, " ")
	reTags := regexp.MustCompile("<[^>]*>")
	html = reTags.ReplaceAllString(html, " ")
	return strings.Join(strings.Fields(html), " ")
}

func extractHTMLTitle(html string) string {
	re := regexp.MustCompile(`(?i)<title>(.*?)</title>`)
	match := re.FindStringSubmatch(html)
	if len(match) > 1 {
		res := strings.TrimSpace(match[1])
		if strings.Contains(res, "{") || strings.Contains(res, "<") { return "" }
		return res
	}
	return ""
}

func AnalyzeIntelligence(content string, isActive bool) (int, string) {
	lower := strings.ToLower(content)

	if !isActive || len(content) < 150 ||
		strings.Contains(lower, "connection refused") ||
		strings.Contains(lower, "timeout") {
		return 0, "Inactive / Unreachable"
	}

	score := 1 
	categoryScores := map[string]int{}

	if strings.Contains(lower, "ransomware") { 
		score += 4
		categoryScores["Ransomware & Malware"] += 4
	}
	if strings.Contains(lower, "lockbit") || strings.Contains(lower, "payload") {
		score += 3
		categoryScores["Ransomware & Malware"] += 3
	}

	if strings.Contains(lower, "leak") || strings.Contains(lower, "dump") {
		score += 2
		categoryScores["Data Breach"] += 2
	}
	if strings.Contains(lower, "database") || strings.Contains(lower, "sql") {
		score += 2
		categoryScores["Data Breach"] += 2
	}

	if strings.Contains(lower, "cvv") || strings.Contains(lower, "carding") {
		score += 3
		categoryScores["Financial Fraud"] += 3
	}
	if strings.Contains(lower, "btc") || strings.Contains(lower, "monero") {
		score += 1
		categoryScores["Financial Fraud"] += 1
	}

	if strings.Contains(lower, "market") || strings.Contains(lower, "vendor") {
		score += 2
		categoryScores["Dark Web Marketplace"] += 2
	}
	if strings.Contains(lower, "shop") || strings.Contains(lower, "order") {
		score += 1
		categoryScores["Dark Web Marketplace"] += 1
	}

	if len(content) > 1000 {
		score += 1
	}
	if len(content) > 3000 {
		score += 1
	}

	if score > 10 {
		score = 10
	}

	finalCategory := "General Intelligence"
	maxCatScore := 0
	for cat, val := range categoryScores {
		if val > maxCatScore {
			maxCatScore = val
			finalCategory = cat
		}
	}

	return score, finalCategory
}
