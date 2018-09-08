package spriter

import "fmt"

type Entity struct {
	Id               int             `xml:"id,attr"`
	Name             string          `xml:"name,attr"`
	Animations       []*Animation    `xml:"animation"`
	CharacterMaps    []*CharacterMap `xml:"character_map"`
	ObjectInfos      []*ObjectInfo   `xml:"obj_info"`
	MaxNumTimelines  int
	animationPointer int
	namedAnimations  map[string]*Animation
}

func (e *Entity) String() string {
	return fmt.Sprintf("id: %d, name:%s", e.Id, e.Name)
}

func (e *Entity) getAnimationByIndex(index int) *Animation {
	return e.Animations[index]
}

func (e *Entity) getAnimationByName(name string) *Animation {
	if e.namedAnimations == nil {
		e.namedAnimations = make(map[string]*Animation)
		for i := range e.Animations {
			a := e.Animations[i]
			e.namedAnimations[a.Name] = a
		}
	}

	return e.namedAnimations[name]
}

func (e *Entity) AnimationsCount() int {
	return len(e.Animations)
}

func (e *Entity) containsAnimation(anim *Animation) bool {
	for i := range e.Animations {
		if anim == e.Animations[i] {
			return true
		}
	}
	return false
}

func (e *Entity) getCharacterMap(name string) *CharacterMap {
	for i := range e.CharacterMaps {
		cp := e.CharacterMaps[i]
		if cp.Name == name {
			return cp
		}
	}
	return nil
}

func (e *Entity) getInfoByIndex(index int) *ObjectInfo {
	return e.ObjectInfos[index]
}

func (e *Entity) getInfoByName(name string) *ObjectInfo {
	for i := range e.ObjectInfos {
		oi := e.ObjectInfos[i]
		if oi.Name == name {
			return oi
		}
	}
	return nil
}

type ObjectType string

const (
	TypeBone   ObjectType = "bone"
	TypeSkin   ObjectType = "skin"
	TypeBox    ObjectType = "box"
	TypePoint  ObjectType = "point"
	TypeSprite ObjectType = "sprite"
)

type ObjectInfo struct {
	Name   string     `xml:"name,attr"`
	Type   ObjectType `xml:"type,attr"`
	Width  float64    `xml:"w,attr"`
	Height float64    `xml:"h,attr"`
}

func MakeObjectInfo(name string, t ObjectType, w float64, h float64) *ObjectInfo {
	o := &ObjectInfo{
		Name:   name,
		Type:   t,
		Width:  w,
		Height: h,
	}
	return o
}

func (oi *ObjectInfo) String() string {
	return fmt.Sprintf("[name: %s, type: %s, size: %dx%d", oi.Name, oi.Type, oi.Width, oi.Height)
}

type CharacterMap struct {
	Id           int                  `xml:"id,attr"`
	Name         string               `xml:"name"`
	Maps         []mapInstructionData `xml:"map"`
	FilesMapping map[int]int
}

type mapInstructionData struct {
	File         int  `xml:"file,attr"`
	Folder       int  `xml:"folder,attr"`
	TargetFile   *int `xml:"target_file,attr"`
	TargetFolder *int `xml:"target_folder,attr"`
}
