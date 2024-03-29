package global

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

// GlowImage takes a asset path, load it, and add a glow effect around it
func GlowImage(spriteassetpath string, color uint32) *sdl.Surface {
	log.Debug("GlowImage(", spriteassetpath, ",", color, ")")
	red := byte((color & 0xff0000) >> 16)
	green := byte((color & 0x00ff00) >> 8)
	blue := byte(color & 0x0000ff)

	if image, err := LoadImageAsset(spriteassetpath); err == nil {
		if image2, err := LoadImageAsset(spriteassetpath); err == nil {
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

// FlashImage will flash/brighten a surface depending on the 'level'
// - level=8 -> completely white
// - level=0 -> the normal image
func FlashImage(surface *sdl.Surface, level uint32) {
	log.Debug("FlashImage(", surface, ",", level, ")")
	if level > 8 {
		level = 8
	}
	if level > 0 {
		threshold := byte(level*32 - 1)
		surface.Lock()
		pixels := surface.Pixels()
		lenpixels := len(pixels)
		for i := 0; i < lenpixels; i += 4 {
			if pixels[i+3] != 0 {
				if pixels[i] < threshold {
					pixels[i] = threshold
				}
				if pixels[i+1] < threshold {
					pixels[i+1] = threshold
				}
				if pixels[i+2] < threshold {
					pixels[i+2] = threshold
				}
			}
		}
		surface.Unlock()
	}
}

// ParseMega parse a string like 100M and translate it into
// a int32 in Megabytes
func ParseMega(str string) int32 {
	log.Debug("ParseMega(", str, ")")
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
		if value < 2048*1024 && value > 0 {
			value = value * 1024
		} else {
			value = 2147483647
		}
	}
	if values[2] == "T" {
		if value < 2048 && value > 0 {
			value = value * 1024 * 1024
		} else {
			value = 2147483647
		}
	}
	if int32(value) < 0 {
		return 2147483647
	}
	return int32(value)
}

// AdjustMega get a number in megabyte, reduce it in GB, TB and return the string result
func AdjustMega(mega int32) string {
	if mega >= 2000000 {
		return strconv.Itoa(int(mega/1000000)) + " TB"
	}
	if mega >= 2000 {
		return strconv.Itoa(int(mega/1000)) + " GB"
	}
	return strconv.Itoa(int(mega)) + " MB"
}

// LoadImageAsset load an SDL (PNG) image directly from assets
func LoadImageAsset(filename string) (*sdl.Surface, error) {
	data, err := Asset(filename)
	if err != nil {
		return nil, err
	}
	src, err := sdl.RWFromMem(data)
	if err != nil {
		return nil, fmt.Errorf("Unable to read " + filename + ": " + err.Error())
	}
	imagetype := strings.ToUpper(filename[len(filename)-3:])

	return img.LoadTypedRW(src, false, imagetype)
}

// AdjustImage will resize an image to fit into w x h surface image
func AdjustImage(image *sdl.Surface, w, h int32) (*sdl.Surface, error) {
	dst, err := sdl.CreateRGBSurface(0, w, h, 32, 0x00ff0000, 0x0000ff00, 0x000000ff, 0xff000000)
	if err != nil {
		return nil, err
	}
	dstw := w
	dsth := h
	if image.W*h > image.H*w {
		dsth = image.H * w / image.W
	} else {
		dstw = image.W * h / image.H
	}
	xshift := (w - dstw) / 2
	yshift := (h - dsth) / 2
	image.BlitScaled(&sdl.Rect{X: 0, Y: 0, W: image.W, H: image.H}, dst, &sdl.Rect{X: xshift, Y: yshift, W: dstw, H: dsth})
	return dst, nil
}

// helper function, to know which pixel is in (x.y)
//
// It is mainly used to know if we are on a transparent pixel
func GetSurfacePixel(surface *sdl.Surface, x, y int32) (red, green, blue, alpha uint8) {
	if x < 0 || x >= surface.W || y < 0 || y >= surface.H {
		return 0, 0, 0, 0
	}
	err := surface.Lock()
	if err != nil {
		panic(err)
	}
	bpp := surface.Format.BytesPerPixel
	bytes := surface.Pixels()
	red = bytes[int(y)*int(surface.Pitch)+int(x)*int(bpp)]
	green = bytes[int(y)*int(surface.Pitch)+int(x)*int(bpp)+1]
	blue = bytes[int(y)*int(surface.Pitch)+int(x)*int(bpp)+2]
	alpha = bytes[int(y)*int(surface.Pitch)+int(x)*int(bpp)+3]

	surface.Unlock()
	return
}
