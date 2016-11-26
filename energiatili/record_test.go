package energiatili_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/joneskoo/etget/energiatili"
)

func TestRecordBidirectional(t *testing.T) {
	dp := energiatili.Record{Time: time.Date(2016, 11, 26, 6, 35, 0, 0, time.Local), Value: 1.5}
	b, err := json.Marshal(dp)
	if err != nil {
		t.Fatalf("JSON marshal of Record: %s", err)
	}
	t.Logf("%#v %q", dp, b)
	var dp2 energiatili.Record
	err = json.Unmarshal(b, &dp2)
	if err != nil {
		t.Fatalf("JSON unmarshal of marshal result: %s", err)
	}
	if dp2.Value != dp.Value || !dp2.Time.Equal(dp.Time) {
		t.Fatalf("want unmarshal(marshal(...)) to restore time %s, got %s", dp.Time.UTC(), dp2.Time.UTC())
	}
}
