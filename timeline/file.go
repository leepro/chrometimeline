package timeline

import (
	"encoding/json"
	"io/ioutil"
)

func Render(fn string, output string) error {
	bs, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}

	var data map[string]interface{}
	err = json.Unmarshal(bs, &data)
	if err != nil {
		return err
	}

	var ret []*Item
	var domComplete float64
	var loadEventEnd float64

	items := data["timings"].([]interface{})
	for _, _d := range items {
		d := _d.(map[string]interface{})
		if d["entryType"].(string) == "navigation" {
			ret = append(ret, &Item{
				d["startTime"].(float64), d["responseEnd"].(float64), d["name"].(string),
				d["initiatorType"].(string),
			})

			domComplete = d["domComplete"].(float64)
			loadEventEnd = d["loadEventEnd"].(float64)
		} else if d["entryType"].(string) == "resource" {
			ret = append(ret, &Item{
				d["startTime"].(float64), d["responseEnd"].(float64), d["name"].(string),
				d["initiatorType"].(string),
			})
		}
	}

	var max float64
	category := make(map[string]int)
	for _, d := range ret {
		// fmt.Printf("[%d] %d %d %s\n", i, int(d.Start)/10, int(d.End)/10, d.Type)
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
