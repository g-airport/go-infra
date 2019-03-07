package auth

const RightsNone = 0

//const SessionExpireTime = 1800 * time.Second

const (
	RightsUser = 1 << iota // in hex:0x1 - 1
)

// for example
var RightsMap = map[string]int32{
	"RightsUser": RightsUser,
}
