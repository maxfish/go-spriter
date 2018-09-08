package spriter

import (
	"fmt"
	"math"
)

type Animation struct {
	Id         int    `xml:"id,attr"`
	Name       string `xml:"name,attr"`
	Length     int    `xml:"length,attr"`
	Interval   int    `xml:"interval,attr"`
	XMLLooping *bool  `xml:"looping,attr"`
	Looping    bool
	Mainline   *Mainline   `xml:"mainline"`
	Timelines  []*Timeline `xml:"timeline"`

	nameToTimeline           map[string]*Timeline
	currentKey               *MainlineKey
	unmappedInterpolatedKeys []*TimelineKey
	interpolatedKeys         []*TimelineKey
}

func (a *Animation) String() string {
	toReturn := fmt.Sprintf("Anim [id:%d, name:%s, interval:%d, duration:%d, looping:'%t'", a.Id, a.Name, a.Interval, a.Length, a.Looping)
	toReturn += "\n" + a.Mainline.String()
	toReturn += "\nTimelines:\n"
	for i := range a.Timelines {
		toReturn += a.Timelines[i].String() + "\n"
	}
	toReturn += "]"
	return toReturn
}

func (a *Animation) initialize() {
	a.interpolatedKeys = make([]*TimelineKey, len(a.Timelines))
	a.unmappedInterpolatedKeys = make([]*TimelineKey, len(a.Timelines))

	for i := 0; i < len(a.interpolatedKeys); i++ {
		a.interpolatedKeys[i] = MakeTimelineKey(i)
		a.interpolatedKeys[i].setObject(MakeTimelineKeyObject())
		a.unmappedInterpolatedKeys[i] = MakeTimelineKey(i)
		a.unmappedInterpolatedKeys[i].setObject(MakeTimelineKeyObject())
	}
	if len(a.Mainline.Keys) > 0 {
		a.currentKey = a.Mainline.Keys[0]
	}
}

func (a *Animation) getTimelineByName(name string) *Timeline {
	if a.nameToTimeline == nil {
		a.nameToTimeline = make(map[string]*Timeline)
		for i := range a.Timelines {
			t := a.Timelines[i]
			a.nameToTimeline[t.Name] = t
		}
	}
	return a.nameToTimeline[name]
}

func (a *Animation) timelines() int {
	return len(a.Timelines)
}

func (a *Animation) update(time int, root *TimelineKeyObject) {
	if root == nil {
		fmt.Println("Error: The root can not be nil! Set a root bone to apply this animation relative to the root bone.")
		return
	}
	a.currentKey = a.Mainline.getKeyBeforeTime(time)

	for i := range a.unmappedInterpolatedKeys {
		a.unmappedInterpolatedKeys[i].active = false
	}
	for i := range a.currentKey.BoneRefs {
		a.updateObjectRef(a.currentKey.BoneRefs[i], root, time)
	}
	for i := range a.currentKey.ObjectRefs {
		a.updateObjectRef(a.currentKey.ObjectRefs[i], root, time)
	}
}

func (a *Animation) updateObjectRef(ref *ObjectRef, root *TimelineKeyObject, time int) {
	// Get the timelines, the refs pointing to
	timeline := a.Timelines[ref.Timeline]
	key := timeline.Keys[ref.Key]
	nextKey := timeline.Keys[(ref.Key+1)%len(timeline.Keys)]
	currentTime := key.Time
	nextTime := nextKey.Time

	// This happens when time is lower than the time of the first key in the timeline
	// e.g. |--x---█--------█----█|
	if time < currentTime {
		nextTime += a.Length
		time += a.Length
	}

	// This happens when time is greater than the time of the last key in the timeline
	// e.g. |█-----█--------█--x--|
	if nextTime < currentTime {
		if !a.Looping {
			nextKey = key
		} else {
			nextTime += a.Length
		}
	}

	// Normalize the time
	t := float64(time-currentTime) / float64(nextTime-currentTime)
	if math.IsNaN(t) || math.IsInf(t, 0) {
		t = 1
	}
	if a.currentKey.Time > currentTime {
		tMid := float64(a.currentKey.Time-currentTime) / float64(nextTime-currentTime)
		if math.IsNaN(tMid) || math.IsInf(tMid, 0) {
			tMid = 0
		}
		t = float64(time-a.currentKey.Time) / float64(nextTime-a.currentKey.Time)
		if math.IsNaN(t) || math.IsInf(t, 0) {
			t = 1
		}
		t = a.currentKey.Curve().interpolate(tMid, 1, t)
	} else {
		t = a.currentKey.Curve().interpolate(0, 1, t)
	}

	bone1 := key.object
	bone2 := nextKey.object
	tweenTarget := a.interpolatedKeys[ref.Timeline].object
	tweenTarget.objectType = bone1.objectType
	a.interpolateObject(bone1, bone2, tweenTarget, t, key.Curve, key.Spin)
	a.unmappedInterpolatedKeys[ref.Timeline].active = true
	refParent := root
	if ref.ParentRef != nil {
		refParent = a.unmappedInterpolatedKeys[ref.ParentRef.Timeline].object
	}
	a.unmapTimelineObject(ref.Timeline, refParent)
}

func (a *Animation) unmapTimelineObject(timeline int, root *TimelineKeyObject) {
	mapTarget := a.unmappedInterpolatedKeys[timeline].object
	mapTarget.setWithBone(a.interpolatedKeys[timeline].object)
	mapTarget.unmapCoordinates(root)
}

func (a *Animation) interpolateBone(bone1 *TimelineKeyObject, bone2 *TimelineKeyObject, target *TimelineKeyObject, t float64, curve *Curve, spin int) {
	target.Angle = curve.interpolateAngleWithSpin(bone1.Angle, bone2.Angle, t, spin)
	curve.interpolatePoints(bone1.Position, bone2.Position, t, target.Position)
	curve.interpolatePoints(bone1.Scale, bone2.Scale, t, target.Scale)
	curve.interpolatePoints(bone1.Pivot, bone2.Pivot, t, target.Pivot)
}

func (a *Animation) interpolateObject(object1 *TimelineKeyObject, object2 *TimelineKeyObject, target *TimelineKeyObject, t float64, curve *Curve, spin int) {
	a.interpolateBone(object1, object2, target, t, curve, spin)
	target.Alpha = curve.interpolateAngle(object1.Alpha, object2.Alpha, t)
	target.fileIndex = object1.fileIndex
}
