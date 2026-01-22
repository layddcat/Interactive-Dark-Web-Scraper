# Interactive Dark Web Scraper Threat Intelligence Dashboard
Bu proje, Dark Web (.onion) kaynaklarından toplanan içeriklerin işlenerek görsel bir dashboard üzerinde gösterilmesini amaçlamaktadır. Sistem, toplanan verileri analiz eder, her içerik için otomatik başlık üretir ve kayıtları kategori ile kritiklik derecesine göre sınıflandırır.

## 1. Proje Amacı
Projenin amacı, Dark Web ortamından elde edilen karmaşık ve dağınık içerikleri yapılandırılmış bir veritabanında saklamak ve bu verileri web tabanlı bir arayüz üzerinden analistlerin kolayca inceleyebileceği hale getirmektir.

## 1.1 Başlıklandırma Mekanizması
Toplanan içeriklerde öncelikle HTML `<title>` etiketi kontrol edilir. Geçerli bir başlık bulunamazsa, temizlenmiş metin içerisinden ilk kelimeler kullanılarak içerikle uyumlu bir başlık otomatik olarak oluşturulur. Bu yöntem sayesinde dashboard üzerinde her kayıt okunabilir ve anlamlı bir şekilde listelenir.

## 2. Kullanılan Teknolojiler
- Go (Golang)
- Gin-Gonic
- GORM
- PostgreSQL
- Tor SOCKS5 Proxy
- Docker & Docker Compose
- Tailwind CSS

## 3. Kurulum ve Çalıştırma
Proje Docker ortamında çalışacak şekilde hazırlanmıştır. Gerekli tüm bileşenler container içinde yer almaktadır.

docker-compose up --build

## 4. Uygulama Arayüzü: 
URL: http://localhost:8080/login

Yetkili kullanıcı bilgileri:
- Kullanıcı Adı: admin
- Parola: cti123