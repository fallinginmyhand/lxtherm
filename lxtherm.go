package lxtherm

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// FilterDir to filter list of dir
func filterDir(fp string, filter string) ([]os.FileInfo, error) {
	fi, fierr := listDir(fp)
	if fierr != nil {
		return nil, fierr
	}
	ffi := []os.FileInfo{}
TOPLP:
	for _, f := range fi {
		switch strings.ToLower(filter) {
		case "f":
			if f.IsDir() == false {
				ffi = append(ffi, f)
			}
			break
		case "d":
			if f.IsDir() == true {
				ffi = append(ffi, f)
			}
			break

		}
		ffi = append(ffi, fi...)
		break TOPLP

	}
	return ffi, nil
}

// Thermals fetch thermal sensors' value
func Thermals() map[string]float64 {
	sn := map[string]float64{}
	thermals, errthermals := filterDir("/sys/class/thermal/", "d")
	if errthermals == nil {
		for _, th := range thermals {
			if strings.Index(th.Name(), "therm") == 0 {
				sensorfp := "/sys/class/thermal/" + th.Name() + "/temp"
				f, e := os.Open(sensorfp)
				if e != nil {
					continue
				}
				rb := make([]byte, 10)
				rbl, errrbl := f.Read(rb)
				if errrbl != nil || rbl == 0 {
					continue
				}
				thv, errthv := strconv.ParseFloat(string(rb[0:rbl-1]), 64)
				if errthv != nil {
					continue
				}
				sn[th.Name()] = thv / 1000
			}
		}

	}
	return sn
}

// ListDir get list of dir
func listDir(fp string) ([]os.FileInfo, error) {
	fd, err := ioutil.ReadDir(fp)
	if err != nil {
		return nil, err
	}
	return fd, nil
}
