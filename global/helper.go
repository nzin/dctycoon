package global

import (
	"regexp"
	"strconv"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

func GlowImage(spritepath string, color uint32) *sdl.Surface {
	red := byte((color & 0xff0000) >> 16)
	green := byte((color & 0x00ff00) >> 8)
	blue := byte(color & 0x0000ff)

	if image, err := img.Load(spritepath); err == nil {
		if image2, err := img.Load(spritepath); err == nil {
			if image2.Format.BytesPerPixel == 4 {
				pixels := image.Pixels()
				pixels2 := image2.Pixels()
				image2.Lock()
				for x := int32(1); x < image2.W-1; x++ {
					for y := int32(1); y < image2.H-1; y++ {
						if (pixels[(y*image2.W+x)*4+3] == 0) &&
							((pixels[((y+1)*image2.W+x)*4+3] != 0) ||
								(pixels[((y-1)*image2.W+x)*4+3] != 0) ||
								(pixels[(y*image2.W+(x+1))*4+3] != 0) ||
								(pixels[(y*image2.W+(x-1))*4+3] != 0)) {
							pixels2[(y*image2.W+x)*4+3] = 0xff
							pixels2[(y*image2.W+x)*4+0] = red
							pixels2[(y*image2.W+x)*4+1] = green
							pixels2[(y*image2.W+x)*4+2] = blue
						}
					}
				}
				image2.Unlock()
			}
			return image2
		}
	}
	return nil
}

//
// parse a string like 100M and translate it into
// a int32 int Megabytes
func ParseMega(str string) int32 {
	re := regexp.MustCompile("([0-9]+) *([MGT]?)")
	values := re.FindStringSubmatch(str)
	if len(values) < 2 {
		return 0
	}
	value, err := strconv.Atoi(values[1])
	if err != nil {
		return 0
	}
	if values[2] == "G" {
		if value < 2048*1024 {
			value = value * 1024 * 1024
		} else {
			value = value * 1024
		}
	}
	if values[2] == "T" {
		if value < 2048 {
			value = value * 1024 * 1024
		} else {
			value = 2147483647
		}
	}
	return int32(value)
}
