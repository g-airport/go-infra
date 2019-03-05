package validate

import (
	"regexp"

	"github.com/asaskevich/govalidator"
)

const (
	ChineseChar = "chinesechar"
	VisibleChar = "visiblechar"
)

var (
	RxChineseChar = regexp.MustCompile("^[\u4e00-\u9fa5\\w]+$")
	RxVisibleChar = regexp.MustCompile("^[\\w\\/\\.]+$")
)

func init() {
	govalidator.TagMap[ChineseChar] =
		govalidator.Validator(func(str string) bool {
			return RxChineseChar.MatchString(str)
		})

	govalidator.TagMap[VisibleChar] =
		govalidator.Validator(func(str string) bool {
			return RxVisibleChar.MatchString(str)
		})
}
