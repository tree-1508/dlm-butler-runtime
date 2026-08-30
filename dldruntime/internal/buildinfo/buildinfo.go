package buildinfo

import (
	"fmt"
	"strconv"
	"time"
)

var (
	Version = "head"
	BuiltAt = ""
	Commit  = ""
)

func VersionString() string {
	timeString := "no build date"
	if BuiltAt != "" {
		if epoch, err := strconv.ParseInt(BuiltAt, 10, 64); err == nil {
			timeString = fmt.Sprintf("built on %s", time.Unix(epoch, 0).Format("Jan _2 2006 @ 15:04:05"))
		}
	}
	res := fmt.Sprintf("%s, %s", Version, timeString)
	if Commit != "" {
		res = fmt.Sprintf("%s, ref %s", res, Commit)
	}
	return res
}
