package main

import (
	downloader "github.com/HikariKnight/quickpassthrough/internal/ls_iommu_downloader"
	"github.com/HikariKnight/quickpassthrough/internal/tuimode"
	"github.com/HikariKnight/quickpassthrough/pkg/params"
)

func main() {
	// Get all our arguments in 1 neat struct
	pArg := params.NewParams()

	if !pArg.Flag["gui"] {
		downloader.GetLsIOMMU()
		tuimode.App()
	}
}
