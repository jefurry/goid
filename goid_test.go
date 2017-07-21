package goid

import (
	"testing"
	"time"
)

func TestIdString(t *testing.T) {
	data := []struct {
		serverRoomID, clusterID, machineID, workID, opID uint32
	}{
		{0, 0, 0, 0, 0},
		{1, 2, 3, 4, 5},
		{2, 3, 4, 5, 6},
		{5, 2, 1, 3, 53},
	}

	id := New()
	id.SetSalt("HDf6%n2#lfY(76Ve)")

	for _, v := range data {
		s, err := id.Encode(v.serverRoomID, v.clusterID, v.machineID, v.workID, v.opID)
		if err != nil {
			t.Fatal(err.Error())
		}

		err = id.FromString(s)
		if err != nil {
			t.Fatal(err.Error())
		}

		if v.serverRoomID != id.ServerRoomID() ||
			v.clusterID != id.ClusterID() ||
			v.machineID != id.MachineID() ||
			v.workID != id.WorkID() ||
			v.opID != id.OpID() {
			t.Fatal("id mismatching.")
		}

		timestamp := id.Timestamp()
		if timestamp <= 0 || timestamp > time.Now().UnixNano() {
			t.Fatal("timestamp invalid.")
		}
	}
}
