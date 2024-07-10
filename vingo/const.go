package vingo

const (
	Enable  = 1
	Disable = 2

	True  = 1
	False = 2

	DateYearFormat        = "2006"
	DateYearMonthFormat   = "2006-01"
	DateFormat            = "2006-01-02"
	DateFormatChinese     = "2006年01月02日"
	DatetimeFormat        = "2006-01-02 15:04:05"
	DatetimeFormatChinese = "2006年01月02日 15时04分05秒"

	Add    = true
	Remove = false

	Male   = "男"
	Female = "女"

	DaySec = 3600 * 24  // 1天秒数
	GB     = 1073741824 // 1GM

	VIDEO_TYPE = "video"
	AUDIO_TYPE = "audio"
	IMAGE_TYPE = "image"
	FILE_TYPE  = "file"

	ASC  = "ASC"
	DESC = "DESC"
)

// 身份证号系数
var idCardFactors = []uint{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}

// 身份证号效验码
var idCardCodes = map[uint]string{0: "1", 1: "0", 2: "X", 3: "9", 4: "8", 5: "7", 6: "6", 7: "5", 8: "4", 9: "3", 10: "2"}
