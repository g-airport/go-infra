package avatarbuilder

import "image/color"

type BuilderOption func(*AvatarBuilder)

func SetFontSize(size float64) BuilderOption {
	return func(builder *AvatarBuilder) {
		builder.fontSize = size
	}
}

func SetAvatarSize(w int, h int) BuilderOption {
	return func(builder *AvatarBuilder) {
		builder.W = w
		builder.H = h
	}
}

func SetFrontGroundColor(c color.Color) BuilderOption {
	return func(builder *AvatarBuilder) {
		builder.fg = c
	}
}

func SetBackGroundColor(c color.Color) BuilderOption {
	return func(builder *AvatarBuilder) {
		builder.bg = c
	}
}

func SetFrontGroundColorHex(hex uint32) BuilderOption {
	return func(builder *AvatarBuilder) {
		builder.fg = builder.hexToRGBA(hex)
	}
}

func SetBackGroundColorHex(hex uint32) BuilderOption {
	return func(builder *AvatarBuilder) {
		builder.bg = builder.hexToRGBA(hex)
	}
}
