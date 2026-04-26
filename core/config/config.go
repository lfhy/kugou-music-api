package config

const (
	AppID         = "1005"
	ClientVer     = "20489"
	LiteAppID     = "3116"
	LiteClientVer = "11440"
)

func PlatformConfig(isLite bool) (appid, clientVer string) {
	if isLite {
		return LiteAppID, LiteClientVer
	}
	return AppID, ClientVer
}
