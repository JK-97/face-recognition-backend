package model

import (
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/config"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
)

func isOutside(box [4]float32) bool {
	cfg := config.Config()
	fence := cfg.Get("fence-position").(schema.FencePos)

	center := [2]float32{(box[0] + box[2]) / 2, (box[1] + box[3]) / 2}
	if center[0] > fence.Xmin && center[0] < fence.Xmax && center[1] > fence.Ymin && center[1] < fence.Ymax {
		return false
	}
	return true
}

// SetFence sets fence position
func SetFence(fence schema.FencePos) {
	cfg := config.Config()
	cfg.Set("fence-position", fence)
}

// GetFence gets fence position
func GetFence() schema.FencePos {
	cfg := config.Config()
	return cfg.Get("fence-position").(schema.FencePos)
}
