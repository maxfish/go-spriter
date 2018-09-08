package spriter

import "fmt"

type Mainline struct {
	Keys []*MainlineKey `xml:"key"`
}

func (m *Mainline) getKeyBeforeTime(time int) *MainlineKey {
	found := m.Keys[0]
	for i := range m.Keys {
		k := m.Keys[i]
		if k.Time <= time {
			found = k
		} else {
			break
		}
	}
	return found
}

func (m *Mainline) String() string {
	toReturn := "Mainline ["
	for i := range m.Keys {
		toReturn += "\n" + m.Keys[i].String()
	}
	toReturn += "]"
	return toReturn
}

type MainlineKey struct {
	Id         int          `xml:"id,attr"`
	Time       int          `xml:"time,attr"`
	BoneRefs   []*ObjectRef `xml:"bone_ref"`
	ObjectRefs []*ObjectRef `xml:"object_ref"`
	CurveType  *string      `xml:"curve_type,attr"`
	curve      *Curve
}

func (k *MainlineKey) String() string {
	toReturn := fmt.Sprintf("MainlineKey [id:%d, time:%d, %s\n", k.Id, k.Time, k.curve.String())
	for i := range k.BoneRefs {
		toReturn += "\t" + k.BoneRefs[i].String() + "\n"
	}
	for i := range k.ObjectRefs {
		toReturn += "\t" + k.ObjectRefs[i].String() + "\n"
	}
	toReturn += "]"
	return toReturn
}

func (k *MainlineKey) Curve() *Curve {
	if k.curve == nil {
		k.curve = MakeCurve()
	}
	return k.curve
}

type ObjectRef struct {
	Id        int    `xml:"id,attr"`
	Key       int    `xml:"key,attr"`
	Parent    *int   `xml:"parent,attr"`
	Timeline  int    `xml:"timeline,attr"`
	ZIndex    string `xml:"z_index,attr"`
	ParentRef *ObjectRef
}

func (r *ObjectRef) String() string {
	return fmt.Sprintf("ObjectRef [id:%d, key:%d, parent:%d, timeline:%d, zIndex:%s]", r.Id, r.Key, r.Parent, r.Timeline, r.ZIndex)
}
