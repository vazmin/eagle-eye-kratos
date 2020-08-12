package service

import (
	pb "github.com/vazmin/eagle-eye-kratos/service/licensing/api"
	"testing"
)

func TestMerge(t *testing.T) {
	dest := &pb.License{LicenseId: "abc", ProduceName: "foo", LicenseMax: 10}
	src := &pb.License{LicenseId: "abc", ProduceName: "bar", LicenseAllocated:  100}
	dest.XXX_Merge(src)
	if dest.ProduceName != "bar" {
		t.Error("")
	}
}
