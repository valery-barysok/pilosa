package index

import (
	"log"
	"testing"

	//	"io/ioutil"
	//   "time"
	"pilosa/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFragment(t *testing.T) {

	//id := util.Id()
	id := util.Hex_to_SUUID("1")
	id2 := util.Hex_to_SUUID("2")
	dummy := NewFragmentContainer()
	dummy.AddFragment("general", "25", 0, id)
	dummy.AddFragment("Brand", "25", 0, id2)

	Convey("Get ", t, func() {
		bh, _ := dummy.Get(id, 1234)
		So(bh, ShouldNotEqual, 0)
	})

	Convey("SetBit/Count 1 1", t, func() {
		//	bh, _ := dummy.Get(id, 1234)
		bi1 := uint64(1234)
		changed, _ := dummy.SetBit(id, bi1, 1)
		So(changed, ShouldEqual, true)
		changed, _ = dummy.SetBit(id, bi1, 1)
		So(changed, ShouldEqual, false)
		bh, _ := dummy.Get(id, bi1)
		num, _ := dummy.Count(id, bh)
		So(num, ShouldEqual, 1)
	})

	Convey("Union/Intersect", t, func() {
		bi1 := uint64(1234)
		bi2 := uint64(4321)

		dummy.SetBit(id, bi2, 2) //set_bit creates the bitmap

		bh1, _ := dummy.Get(id, bi1)
		bh2, _ := dummy.Get(id, bi2)

		handles := []BitmapHandle{bh1, bh2}
		result, _ := dummy.Union(id, handles)

		num, _ := dummy.Count(id, result)
		So(num, ShouldEqual, 2)
		result, _ = dummy.Intersect(id, handles)

		num, _ = dummy.Count(id, result)
		So(num, ShouldEqual, 0)
	})
	Convey("Bytes", t, func() {
		bi1 := uint64(1234)
		bh1, _ := dummy.Get(id, bi1)
		before, _ := dummy.Count(id, bh1)

		bytes, _ := dummy.GetBytes(id, bh1)
		bh2, _ := dummy.FromBytes(id, bytes)

		after, _ := dummy.Count(id, bh2)
		So(before, ShouldEqual, after)
	})
	Convey("Empty ", t, func() {
		bh, _ := dummy.Empty(id)
		before, _ := dummy.Count(id, bh)
		So(before, ShouldEqual, 0)
	})

	Convey("GetList ", t, func() {
		bhs, _ := dummy.GetList(id, []uint64{1234, 4321, 789})
		result, _ := dummy.Union(id, bhs)
		num, _ := dummy.Count(id, result)
		So(num, ShouldNotEqual, 2)
	})

	Convey("Brand SetBit", t, func() {
		bi1 := uint64(1231)
		bi2 := uint64(1232)
		bi3 := uint64(1233)
		bi4 := uint64(1234)
		for x := uint64(0); x < 1000; x++ {
			if x < 100 {
				dummy.SetBit(id2, bi1, x)
				dummy.SetBit(id2, bi4, x)
			}
			if x < 500 {
				dummy.SetBit(id2, bi2, x)
			}
			if x%3 == 0 {
				dummy.SetBit(id2, bi3, x)
			}
			if x > 700 {
				dummy.SetBit(id2, bi4, x)
			}
		}
		bh1, _ := dummy.Get(id2, bi1)
		log.Println(dummy.TopN(id2, bh1, 4))
		//		dummy.Rank()
		So(1, ShouldEqual, 1)
	})
}
