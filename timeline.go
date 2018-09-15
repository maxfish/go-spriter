package spriter

import "fmt"

type Timeline struct {
	Id         int            `xml:"id,attr"`
	Name       string         `xml:"name,attr"`
	Keys       []*TimelineKey `xml:"key"`
	objectInfo *ObjectInfo
}

func (b *Timeline) String() string {
	toReturn := fmt.Sprintf("Timeline [id:%d, name:%s, object_info:%s\n", b.Id, b.Name, b.objectInfo)
	for i := range b.Keys {
		toReturn += "\t" + b.Keys[i].String() + "\n"
	}
	toReturn += "]"
	return toReturn
}

type TimelineKey struct {
	Id        int `xml:"id,attr"`
	Spin      int
	active    bool
	object    *TimelineKeyObject
	Curve     *Curve
	Time      int     `xml:"time,attr"`
	CurveType string  `xml:"curve_type,attr"`
	C1        float64 `xml:"c1,attr"`
	C2        float64 `xml:"c2,attr"`
	C3        float64 `xml:"c3,attr"`
	C4        float64 `xml:"c4,attr"`

	// These fields are used only to read the data from the XML
	XMLDataBone   *TimelineKeyObject `xml:"bone"`
	XMLDataObject *TimelineKeyObject `xml:"object"`
	XMLSpin       *int               `xml:"spin,attr"`
}

func MakeTimelineKey(id int) *TimelineKey {
	return &TimelineKey{
		Id: id,
	}
}

func (b *TimelineKey) String() string {
	toReturn := fmt.Sprintf("TimelineKey [id:%d, time:%d, spin:%d, %s, object:%s]", b.Id, b.Time, b.Spin, b.Curve.String(), b.object.String())
	return toReturn
}

func (b *TimelineKey) setObject(obj *TimelineKeyObject) {
	b.object = obj
}

// This can represents either a Bone or an Object
type TimelineKeyObject struct {
	File       int     `xml:"file,attr"`
	Folder     int     `xml:"folder,attr"`
	Angle      float64 `xml:"angle,attr"`
	Position   *Point
	Pivot      *Point
	Scale      *Point
	Alpha      float64
	fileIndex  int
	objectType ObjectType

	// These fields are used only to read the data from the XML
	XMLX      float64  `xml:"x,attr"`
	XMLY      float64  `xml:"y,attr"`
	XMLPivotX *float64 `xml:"pivot_x,attr"`
	XMLPivotY *float64 `xml:"pivot_y,attr"`
	XMLScaleX *float64 `xml:"scale_x,attr"`
	XMLScaleY *float64 `xml:"scale_y,attr"`
}

func (b *TimelineKeyObject) String() string {
	return fmt.Sprintf(
		"TimelineKeyObject: [pos:%f,%f, pivot:%f,%f, Scale:%f,%f, angle:%f, type:%s, file:%d/%d]",
		b.Position.X(), b.Position.Y(), b.Pivot.X(), b.Pivot.Y(), b.Scale.X(), b.Scale.Y(), b.Angle, b.objectType, b.Folder, b.File,
	)
}

func MakeTimelineKeyObject() *TimelineKeyObject {
	return &TimelineKeyObject{
		Position:   MakePoint(0, 0),
		Pivot:      MakePoint(0, 1),
		Scale:      MakePoint(1, 1),
		objectType: "object",
	}
}

func MakeTimelineKeyBone() *TimelineKeyObject {
	return &TimelineKeyObject{
		Position:   MakePoint(0, 0),
		Pivot:      MakePoint(0, 1),
		Scale:      MakePoint(1, 1),
		objectType: "bone",
	}
}

func (b *TimelineKeyObject) ObjectType() ObjectType {
	return b.objectType
}

func (b *TimelineKeyObject) setWithBone(bone *TimelineKeyObject) {
	b.Position.Set(bone.Position)
	b.Scale.Set(bone.Scale)
	b.Angle = bone.Angle
	b.Pivot.Set(bone.Pivot)
	b.objectType = bone.objectType
	b.File = bone.File
	b.Folder = bone.Folder
	b.fileIndex = bone.fileIndex
}

func (b *TimelineKeyObject) set(x float64, y float64, angle float64, scaleX float64, scaleY float64) {
	b.Position.Set(&Point{x, y})
	b.Scale.Set(&Point{scaleX, scaleY})
	b.Angle = angle
}

func (b *TimelineKeyObject) unmapCoordinates(parent *TimelineKeyObject) {
	signScaleX := signum(parent.Scale.X())
	signScaleY := signum(parent.Scale.Y())
	b.Angle *= signScaleX * signScaleY
	b.Angle += parent.Angle

	scale := b.Scale.MakeCopy()
	scale.Scale(parent.Scale)
	b.Scale.Set(scale)

	position := b.Position.MakeCopy()
	position.Scale(parent.Scale)
	position.Rotate(parent.Angle)
	position.Add(parent.Position)
	b.Position.Set(position)
}

func (b *TimelineKeyObject) mapCoordinates(parent *TimelineKeyObject) {
	position := b.Position.MakeCopy()
	scale := b.Scale.MakeCopy()

	position.Sub(parent.Position)
	position.Rotate(-parent.Angle)
	position.ScaleCoords(1/parent.Scale.X(), 1/parent.Scale.Y())
	scale.ScaleCoords(1/parent.Scale.X(), 1/parent.Scale.Y())
	b.Angle -= parent.Angle
	b.Angle *= signum(parent.Scale.X()) * signum(parent.Scale.Y())

	b.Position.Set(position)
	b.Scale.Set(scale)
}
