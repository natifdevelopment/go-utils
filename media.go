package utils

import (
	"io"
	"mime/multipart"
	"net/http"
)

var allowedMediaTypes = []string{
	"application/pdf",
}

func IsValidMediaType(file multipart.File) (bool, string) {
	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return false, ""
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return false, ""
	}

	detectedType := http.DetectContentType(buffer)

	for _, mediaType := range allowedMediaTypes {
		if detectedType == mediaType {
			return true, detectedType
		}
	}

	return false, detectedType
}

func IsValidMediaTypeDownload(file io.ReadCloser) (bool, string) {
	defer file.Close()

	buffer := make([]byte, 512)
	bytesRead, err := io.ReadFull(file, buffer)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return false, ""
	}
	buffer = buffer[:bytesRead]

	detectedType := http.DetectContentType(buffer)

	for _, mediaType := range allowedMediaTypes {
		if detectedType == mediaType {
			return true, detectedType
		}
	}

	return false, detectedType
}

var dokIdToNameMap = map[string]string{
	"dokBast": "BAST",

	"dokDraughtSurvey": "Dokumen Draught Survey",

	"dokBaDk": "BA Denda Keterlambatan",

	"dokRincianBaDk":    "Dokumen Rincian Denda Keterlambatan",
	"fileIzin":          "File Izin Sandar",
	"fileNor":           "File NOR",
	"dokBlManifestSkab": "Dokumen Bill ok Lading (B/L), Manifest & SKAB",
	"dokNotulen":        "Dokumen Notulen",
	"dokTimesheetSof":   "Dokumen Timesheet / SOF",
	"fileCoa":           "File COA",

	"fileRoa": "File ROA",

	"fileNorLoadingSof": "File NOR Loading & SOF",
	"fileManifestBl":    "File Manifest, Bill of Lading (B/L)",

	"dokBaDs":               "BA Denda Susut",
	"fileTpb":               "File Tempat Penimbunan Berikat (TPB)",
	"fileNorUnloading":      "File NOR Unloading",
	"fileIzinSandarBongkar": "File Izin Sandar & Bongkar",

	"dokBaDm":      "BA Denda Muat",
	"docTimesheet": "Timesheet",

	"dokBaph":          "BAPH",
	"dokBaphTermin":    "BAPH",
	"dokInvoice":       "Dokumen Invoice",
	"dokBuktiRoyalti":  "Dokumen Bukti Royalti",
	"dokKwitansi":      "Dokumen Kwitansi",
	"dokRefPembayaran": "Nomor Referensi Pembayaran",
	"dokFakturPajak":   "Dokumen Faktur Pajak",
	"dokPenagihan":     "Dokumen Penagihan",

	"dokSpt":       "Dokumen SPT",
	"dokSptUpload": "Dokumen SPT Upload",

	"dokNotaDinas":     "Dokumen Nota Dinas",
	"dokRpt":           "Lembar Verifikasi",
	"dokRptVendor":     "Lampiran Verifikasi",
	"lembarVerifikasi": "Lembar Ceklist",

	"fileCowSurveyor": "File COW Surveyor",
	"fileCoaSurveyor": "File COA Surveyor",

	"dokPjbb":                 "Dokumen Kontrak PJBB",
	"dokPjbbCIF":              "CIF - Dokumen Kontrak PJBB",
	"dokPjbbFOB":              "FOB - Dokumen Kontrak PJBB",
	"dokPjbbCFR":              "CFR - Dokumen Kontrak PJBB",
	"dokPjbbBankGaransi":      "Dokumen Bank Garansi",
	"dokPjbbBankGaransiCIF":   "CIF - Dokumen Bank Garansi",
	"dokPjbbBankGaransiFOB":   "FOB - Dokumen Bank Garansi",
	"dokPjbbBankGaransiCFR":   "CFR - Dokumen Bank Garansi",
	"dokPjbbAmd":              "Dokumen Kontrak PJBB Amandemdn",
	"dokPjbbAmdCIF":           "CIF - Dokumen Kontrak PJBB Amandemdn",
	"dokPjbbAmdFOB":           "FOB - Dokumen Kontrak PJBB Amandemdn",
	"dokPjbbAmdCFR":           "CFR - Dokumen Kontrak PJBB Amandemdn",
	"dokPjab":                 "Dokumen Kontrak PJAB",
	"dokPjabAmd":              "Dokumen Kontrak PJAB Amandemen",
	"baTriwulanBaBankGaransi": "BA Triwulan & BA Bank Garansi",

	"dokNpwp": "Dokumen NPWP",
	"dokPkp":  "Dokumen PKP",

	"fileStokOpname": "File Stok Opname",

	"dokRincianDenda": "Dokumen Rincian Denda",

	"dokForceMajeure": "Dokumen Force Majeure",
	"dokNarasi":       "Dokumen Narasi",
}

func GetDokNameByDokId(dokId string) string {
	return dokIdToNameMap[dokId]
}
