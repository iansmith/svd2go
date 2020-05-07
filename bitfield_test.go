package svd

import (
	"github.com/iansmith/svd/runtime/volatile"
	"testing"
)

func TestBasics(t *testing.T) {
	var mypDep MypDef

	r:=(volatile.Register32)(mypDep.r1.Reg)
	if r.Get() !=0 {
		t.Errorf("first get should be 0 (%d)",mypDep.r1.Reg.Get())
	}

}