package fetch

import (
	"encoding/json"
	"testing"
)

func TestBucketToMMPerHr(t *testing.T) {
	v := bucketToMMPerHr(3)
	if v <= 0 {
		t.Fatalf("expected non-zero for bucket 3 got %f", v)
	}
}

func TestNowcastRespCategories(t *testing.T) {
	data := `[{"Obj_id":"1","Date":"2024-06-20","toi":"1200","vupto":"1500","color":"2","cat1":"1","cat2":"0","message":"ok"}]`
	var arr []districtNowcastResp
	if err := json.Unmarshal([]byte(data), &arr); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(arr) != 1 {
		t.Fatalf("expected one element")
	}
	if arr[0].Cat1 != "1" || arr[0].Message != "ok" {
		t.Fatalf("unexpected fields: %+v", arr[0])
	}
}
