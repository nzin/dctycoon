package global

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func DigAsset(t *testing.T, dirname string) {
	dir, err := AssetDir(dirname)
	assert.Empty(t, err, "Asset "+dirname+" is a directory")
	for _, d := range dir {
		path := dirname + "/" + d
		if _, err := AssetDir(path); err == nil {
			// this is a directory
			DigAsset(t, path)
		} else {
			// this is a file
			data, err := Asset(path)
			assert.NotEmpty(t, data, "Asset "+path+" can be loaded")
			assert.Empty(t, err, "Asset "+path+" can be loaded")
		}
	}
}

func TestAssets(t *testing.T) {
	dir, err := AssetDir("assets")

	assert.NotEmpty(t, dir, "We have an asset directory")
	assert.Empty(t, err, "We can read the root directory")

	DigAsset(t, "assets")
}
