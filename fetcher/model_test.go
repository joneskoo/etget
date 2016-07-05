package fetcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestModel(t *testing.T) {
	filename := "/Users/joneskoo/go/src/github.com/joneskoo/etget/energiatili.json"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var report ConsumptionReport
	err = json.Unmarshal(data, &report)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("%#v\n", report)
	err = report.Update()
	if err != nil {
		t.Errorf("Update(): %v", err)
	}
	fmt.Printf("%#v\n", report.Hours.Consumptions[0].Series.Data)
}
