package model

import (
	"gitlab.jiangxingai.com/luyor/tf-fence-backend/config"
)

func isOutside(box [4]float32) bool {
	cfg := config.Config()
	fence := cfg.Get("fence-position").([4]float32)

	center := [2]float32{(box[0] + box[2]) / 2, (box[1] + box[3]) / 2}
	if center[0] > fence[0] && center[0] < fence[2] && center[1] > fence[1] && center[1] < fence[3] {
		return false
	}
	return true
}

// SetFence sets fence position
func SetFence(fence [4]float32) {
	cfg := config.Config()
	cfg.Set("fence-position", fence)
}

// GetFence gets fence position
func GetFence() [4]float32 {
	cfg := config.Config()
	return cfg.Get("fence-position").([4]float32)
}
