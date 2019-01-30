package timeline

import (
	"io/ioutil"

	"github.com/tidwall/gjson"
)

func Render(fn string, output string) error {
	bs, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}

	var ret []*Item
	var domComplete float64
	var loadEventEnd float64
	var iType string

	value := gjson.Get(string(bs), "timings")
	for _, k := range value.Array() {
		v := k.Map()
		switch v["entryType"].Str {
		case "navigation":
			domComplete = v["domComplete"].Num
			loadEventEnd = v["loadEventEnd"].Num
		case "resource":
		default:
			continue
		}
		if v["initiatorType"].Str == "xmlhttprequest" {
			iType = "ajax"
		} else {
			iType = v["initiatorType"].Str
		}

		ret = append(ret, &Item{
			v["startTime"].Num,
			v["responseEnd"].Num,
			v["name"].Str,
			iType,
		})
	}

	var max float64
	category := make(map[string]int)
	for _, d := range ret {
		if max < d.End {
			max = d.End
		}
		if _, ok := category[d.Type]; ok {
			category[d.Type]++
		} else {
			category[d.Type] = 1
		}
	}

	DrawTimeline(max, domComplete, loadEventEnd, category, ret, output)

	return nil
}
