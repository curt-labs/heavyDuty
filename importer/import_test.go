package importer

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestImport(t *testing.T) {
	Convey("test matchVehicles", t, func() {
		getVehiclePartMap()
		var vps []VehiclePart
		vp := VehiclePart{
			Vehicle: Vehicle{
				Year:  1997,
				Make:  "Ford",
				Model: "F-150",
				Style: "All",
			},
		}

		vps = append(vps, vp)

		vv, err := MatchVehicles(vps)
		So(err, ShouldBeNil)
		t.Log(vv)
	})
}

func TestMaps(t *testing.T) {
	Convey("Test getVehiclePartMap", t, func() {
		vmap, _, err := getVehiclePartMap()
		So(err, ShouldBeNil)
		for i, v := range vmap {
			t.Log(i, "--", v)
		}
	})
}
