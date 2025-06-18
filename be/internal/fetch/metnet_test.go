package fetch

import "testing"

func TestBucketToMMPerHr(t *testing.T) {
	v := bucketToMMPerHr(3)
	if v <= 0 {
		t.Fatalf("expected non-zero for bucket 3 got %f", v)
	}
}
