package utils

import (
	"fmt"
	"time"
)

var dayNames = []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
var monthNames = []string{"Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}
var numberWords = []string{"Nol", "Satu", "Dua", "Tiga", "Empat", "Lima", "Enam", "Tujuh", "Delapan", "Sembilan", "Sepuluh", "Sebelas", "Dua Belas", "Tiga Belas", "Empat Belas", "Lima Belas", "Enam Belas", "Tujuh Belas", "Delapan Belas", "Sembilan Belas"}

func FormatTime(t time.Time) *string {
	if t.IsZero() {
		return nil
	}
	formattedTime := t.Format("2006-01-02 15:04:05")
	return &formattedTime
}

func ConvertNumberToWords(n int) string {
	switch {
	case n < 10:
		return numberWords[n]
	case n < 20:
		return numberWords[n-10] + " Belas"
	case n < 100:
		tens := n / 10
		units := n % 10
		if units == 0 {
			return numberWords[tens] + " Puluh"
		}
		return numberWords[tens] + " Puluh " + numberWords[units]
	case n < 1000:
		hundreds := n / 100
		remainder := n % 100
		if hundreds == 1 {
			if remainder == 0 {
				return "Seratus"
			}
			return "Seratus " + ConvertNumberToWords(remainder)
		}
		return numberWords[hundreds] + " Ratus " + ConvertNumberToWords(remainder)
	case n < 10000:
		thousands := n / 1000
		remainder := n % 1000
		if thousands == 1 {
			if remainder == 0 {
				return "Seribu"
			}
			return "Seribu " + ConvertNumberToWords(remainder)
		}
		return numberWords[thousands] + " Ribu " + ConvertNumberToWords(remainder)
	default:
		return "Invalid"
	}
}

func FormatTimeToText(t time.Time) string {
	dayName := dayNames[int(t.Weekday())]
	dateInWords := ConvertNumberToWords(t.Day())
	monthName := monthNames[int(t.Month())-1]
	yearInWords := ConvertNumberToWords(t.Year())

	formattedText := fmt.Sprintf("%s, tanggal %s bulan %s tahun %s", dayName, dateInWords, monthName, yearInWords)
	return formattedText
}
