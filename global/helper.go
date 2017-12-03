package global

import (
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
