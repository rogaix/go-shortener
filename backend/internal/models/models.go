package models

type ShortenedURL struct {
	ID         string `gorm:"primaryKey;type:varchar(191)"`
	LongURL    string `gorm:"type:varchar(2048);index:idx_long_url,length:768"`
	ShortURL   string `gorm:"type:varchar(255)"`
	ClickCount uint   `gorm:"default:0"`
	QRCode     []byte `gorm:"type:longblob"`
}
