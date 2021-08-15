package utils

import (
	"testing"
)

func Test_Check_Disc_Space(t *testing.T) {
	ds, err := Disk_Space("D:/")
	if err != nil {
		panic(err)
	}
	t.Log("DiskUsage: ", SpaceUsed(ds))
	t.Log("DiskUsage Percent: ", SpaceUsedPercent((ds)))
}
