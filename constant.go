package utils

import "strings"

const VESSEL_CODE = "bbo//cd/jenis-kapal/vessel"
const SPB_CODE = "bbo//cd/jenis-kapal/spb"
const TONGKANG_CODE = "bbo//cd/jenis-kapal/tugboat"

const ConstFailedDeleteApprovedData = "Data already been approved, cannot be deleted"
const ConstFailedInitCaptcha = "Failed to init Captcha, please try again later"
const ConstInternalServerError = "Internal server error, please try again later"
const ConstNoQrCode = "Data QRCode tidak ditemukan dalam dokumen"
const ConstNotFound = "Data tidak ditemukan"
const ConstForeignKeyViolated = "Data tidak dapat diubah karena masih digunakan sebagai referensi oleh data lain"

func CapitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

type BboResponseError struct {
	Msg string
}

func (e BboResponseError) Error() string {
	return e.Msg
}
