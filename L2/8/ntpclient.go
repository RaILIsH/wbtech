package ntpclient

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

// GetExactTime возвращает текущее время, полученное через NTP
func GetExactTime() (time.Time, error) {
	return ntp.Time("pool.ntp.org")
}

// PrintCurrentTime печатает текущее время, полученное через NTP
func PrintCurrentTime() {
	time1, err := GetExactTime()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка получения времени: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(time1.Format(time.RFC3339))
}
