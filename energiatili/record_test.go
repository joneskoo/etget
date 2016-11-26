package energiatili_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/joneskoo/etget/energiatili"
)

func TestRecordBidirectional(t *testing.T) {
	var in, out energiatili.Record

	in = energiatili.Record{
		Timestamp: time.Date(2016, 11, 26, 6, 0, 0, 0, time.UTC),
		Value:     1.5,
	}
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("JSON marshal of Record: %s", err)
	}
	err = json.Unmarshal(b, &out)
	if err != nil {
		t.Fatalf("JSON unmarshal of marshal result: %s", err)
	}
	if out.Value != in.Value || !out.Timestamp.Equal(in.Timestamp) {
		t.Fatalf("want unmarshal(marshal(...)) to restore time %s, got %s", in.Timestamp.UTC(), out.Timestamp.UTC())
	}
}
