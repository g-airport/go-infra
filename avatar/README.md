# avatar builder

Using [go free_type](https://github.com/golang/freetype) to build default avatar with string

## Usage

You can reference `./example`

Some snippet is as blow

```
  // init avatarbuilder, you need to tell builder ttf file and how to alignment text
	ab := avatarbuilder.NewAvatarBuilder("./SourceHanSansSC-Medium.ttf", &calc.SourceHansSansSCMedium{})
	ab.SetBackgroundColorHex(colors[1])
	ab.SetFrontgroundColor(color.White)
	ab.SetFontSize(80)
	ab.SetAvatarSize(200, 200)
	if err := ab.GenerateImageAndSave("哈哈哈", "./out.png"); err != nil {
		fmt.Println(err)
		return
	}
```

## Extend Other Font

Because element of width of each font is different, so you need tell avatar how to align the content.
Avatar already implement a free font(made by google and adobe)'s center algorithm in `./calc`,
If you need other font, feel free to PR or issue.
