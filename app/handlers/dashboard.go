package handlers

import (
	"bufio"
	"interactive-scraper/database"
	"net/http"
	"os"
	"strings"
	"strconv" 

	"github.com/gin-gonic/gin"
)

func GetDashboard(c *gin.Context) {
	var items []database.IntelligenceData
	database.DB.Order("criticality_score desc, id desc").Find(&items)

	sourceCount := 0
	file, err := os.Open("targets.yaml")
	if err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "- http") || strings.HasPrefix(line, "http") {
				sourceCount++
			}
		}
		file.Close()
	}

	var totalCount int64 = int64(len(items))
	var highCriticalityCount int64
	var leakCount int64
	var offlineCount int64

	for _, item := range items {
		if item.CriticalityScore >= 7 {
			highCriticalityCount++
		}
		if item.Category == "Data Breach" {
			leakCount++
		}
		if !item.IsActive {
			offlineCount++
		}
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"items":                items,
		"totalCount":           totalCount,
		"sourceCount":          sourceCount,
		"highCriticalityCount": highCriticalityCount,
		"leakCount":            leakCount,
		"offlineCount":         offlineCount,
	})
}

func GetDataDetail(c *gin.Context) {
	id := c.Param("id")
	var item database.IntelligenceData
	
	if err := database.DB.First(&item, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/dashboard")
		return
	}

	var analysisComment string
	var recommendedAction string
	var impactLevel string

	switch item.Category {
	case "Exploit & Vulnerability":
		analysisComment = "Tespit edilen exploit veya zafiyet paylaşımı, kurumun dışa açık servislerinde yetkisiz erişim riskini artırmaktadır. Tehdit aktörlerinin aktif sömürü ihtimali göz önünde bulundurulmalıdır."
		recommendedAction = "Etkilenen servislerin versiyonlarını doğrulayın, ilgili CVE kayıtlarını inceleyin ve WAF/IPS imzalarını güncelleyin."
		if item.CriticalityScore >= 8 {
			impactLevel = "KRİTİK - Aktif Sömürü Riski"
		} else {
			impactLevel = "YÜKSEK - Zafiyet Maruziyeti"
		}

	case "Data Breach":
		analysisComment = "Paylaşılan veri seti, kuruma ait kimlik, erişim veya müşteri verilerini içerebilir. İkincil saldırılar ve kimlik avı riski artmaktadır."
		recommendedAction = "Sızan verileri IOC olarak ele alın, etkilenen hesaplar için parola sıfırlama ve MFA zorunluluğu uygulayın."
		impactLevel = "YÜKSEK - Veri Gizliliği İhlali"

	case "Financial Fraud":
		analysisComment = "Finansal dolandırıcılık faaliyeti, doğrudan maddi kayıp ve itibar zedelenmesine yol açabilir."
		recommendedAction = "Şüpheli işlem paternlerini izleyin, finansal sistemlerde ek doğrulama kontrolleri uygulayın."
		impactLevel = "ORTA/YÜKSEK - Maddi ve Hukuki Risk"

	case "System Alert":
		analysisComment = "Kaynak geçici veya kalıcı olarak çevrimdışı duruma geçmiş olabilir. Bu durum, takedown operasyonu veya altyapı değişikliğine işaret edebilir."
		recommendedAction = "Kaynağı alternatif kanallardan doğrulayın ve izleme kapsamından düşürmeden önce manuel kontrol sağlayın."
		impactLevel = "DÜŞÜK/ORTA - Görünürlük Kaybı"

	default:
		analysisComment = "Tehdit doğrudan saldırı içermemekle birlikte, tehdit aktörü veya platform davranışları açısından izlenmelidir."
		recommendedAction = "Aktivite yoğunluğunu ve tekrar sıklığını takip edin."
		impactLevel = "DÜŞÜK - İstihbarat Amaçlı"
	}

	c.HTML(http.StatusOK, "detail.html", gin.H{
		"item":              item,
		"analysisComment":   analysisComment,
		"recommendedAction": recommendedAction,
		"impactLevel":       impactLevel,
	})
}


func UpdateCriticality(c *gin.Context) {
	id := c.PostForm("id")
	scoreStr := c.PostForm("score")

	score, err := strconv.Atoi(scoreStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz skor formatı"})
		return
	}

	database.DB.Model(&database.IntelligenceData{}).Where("id = ?", id).Update("criticality_score", score)
	
	c.Redirect(http.StatusFound, "/dashboard")
}