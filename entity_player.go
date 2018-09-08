package spriter

import (
	"fmt"
)

type PlayerListenerInterface interface {
	beforeUpdate(player *EntityPlayer)
	afterUpdate(player *EntityPlayer)
	mainlineKeyChanged(prevKey *MainlineKey, newKey *MainlineKey)
	animationFinished(animation *Animation)
	animationChanged(oldAnim *Animation, newAnim *Animation)
}

type EntityPlayer struct {
	position  *Point
	pivot     *Point
	angle     float64
	entity    *Entity
	animation *Animation
	time      int

	interpolatedKeys             []*TimelineKey
	unmappedInterpolatedKeys     []*TimelineKey
	tempInterpolatedKeys         []*TimelineKey
	tempUnmappedInterpolatedKeys []*TimelineKey

	root        *TimelineKeyObject
	rootIsDirty bool

	currentKey           *MainlineKey
	previousKey          *MainlineKey
	listeners            []PlayerListenerInterface
	objToTimeline        map[*TimelineKeyObject]*TimelineKey
	enabledCharacterMaps map[string]*CharacterMap
}

func (p *EntityPlayer) String() string {
	return fmt.Sprintf("entity: %s, key:%d, time:%d", p.entity.Name, p.currentKey.Id, p.time)
}

func MakeEntityPlayer(entity *Entity) *EntityPlayer {
	p := &EntityPlayer{}
	p.root = MakeTimelineKeyBone()
	p.position = MakePoint(0, 0)
	p.pivot = MakePoint(0, 0)
	p.listeners = make([]PlayerListenerInterface, 0)
	p.objToTimeline = make(map[*TimelineKeyObject]*TimelineKey)
	p.enabledCharacterMaps = make(map[string]*CharacterMap)
	p.setEntity(entity)
	return p
}

func (p *EntityPlayer) Update(timeDeltaMs int) {
	for i := range p.listeners {
		p.listeners[i].beforeUpdate(p)
	}
	if p.rootIsDirty {
		p.updateRoot()
	}
	p.animation.update(p.time, p.root)
	p.currentKey = p.animation.currentKey
	if p.previousKey != p.currentKey {
		for i := range p.listeners {
			p.listeners[i].mainlineKeyChanged(p.previousKey, p.currentKey)
		}
		p.previousKey = p.currentKey
	}

	p.interpolatedKeys = p.tempInterpolatedKeys
	p.unmappedInterpolatedKeys = p.tempUnmappedInterpolatedKeys
	for i := range p.animation.interpolatedKeys {
		p.interpolatedKeys[i].active = p.animation.interpolatedKeys[i].active
		p.unmappedInterpolatedKeys[i].active = p.animation.unmappedInterpolatedKeys[i].active
		p.interpolatedKeys[i].object.setWithBone(p.animation.interpolatedKeys[i].object)
		p.unmappedInterpolatedKeys[i].object.setWithBone(p.animation.unmappedInterpolatedKeys[i].object)
	}

	for i := range p.listeners {
		p.listeners[i].afterUpdate(p)
	}
	p.increaseTime(timeDeltaMs)
}

func (p *EntityPlayer) increaseTime(millisecs int) {
	p.time += millisecs
	if p.time > p.animation.Length {
		p.time = p.time - p.animation.Length
		for i := range p.listeners {
			p.listeners[i].animationFinished(p.animation)
		}
	}
	if p.time < 0 {
		for i := range p.listeners {
			p.listeners[i].animationFinished(p.animation)
		}
		p.time += p.animation.Length
	}
}

func (p *EntityPlayer) updateRoot() {
	p.root.Angle = p.angle
	position := MakePoint(p.pivot.X(), p.pivot.Y())
	position.Rotate(p.angle)
	position.Add(p.position)
	p.root.Position.Set(position)
	p.rootIsDirty = false
}

func (p *EntityPlayer) getBone(index int) *TimelineKeyObject {
	return p.unmappedInterpolatedKeys[p.currentKey.BoneRefs[index].Timeline].object
}

func (p *EntityPlayer) getObject(index int) *TimelineKeyObject {
	return p.unmappedInterpolatedKeys[p.currentKey.ObjectRefs[index].Timeline].object
}

func (p *EntityPlayer) getBoneIndex(name string) int {
	for i := range p.currentKey.BoneRefs {
		bone := p.currentKey.BoneRefs[i]
		if p.animation.Timelines[bone.Timeline].Name == name {
			return bone.Id
		}
	}
	return -1
}

func (p *EntityPlayer) getBoneByName(name string) *TimelineKeyObject {
	return p.unmappedInterpolatedKeys[p.animation.getTimelineByName(name).Id].object
}

func (p *EntityPlayer) getObjectIndex(name string) int {
	for i := range p.currentKey.ObjectRefs {
		obj := p.currentKey.ObjectRefs[i]
		if p.animation.Timelines[obj.Timeline].Name == name {
			return obj.Id
		}
	}
	return -1
}

func (p *EntityPlayer) getObjectByName(name string) *TimelineKeyObject {
	return p.unmappedInterpolatedKeys[p.animation.getTimelineByName(name).Id].object
}

func (p *EntityPlayer) getObjectInfoFor(boneOrObject *TimelineKeyObject) *ObjectInfo {
	return p.animation.Timelines[p.objToTimeline[boneOrObject].Id].objectInfo
}

func (p *EntityPlayer) getKeyFor(boneOrObject *TimelineKeyObject) *TimelineKey {
	return p.objToTimeline[boneOrObject]
}

func (p *EntityPlayer) unmapObjects(base *ObjectRef) {
	start := -1
	if base != nil {
		start = base.Id - 1
	}
	for i := start + 1; i < len(p.currentKey.BoneRefs); i++ {
		ref := p.currentKey.BoneRefs[i]
		if ref.ParentRef != base && base != nil {
			continue
		}
		parent := p.unmappedInterpolatedKeys[ref.ParentRef.Timeline].object
		if ref.ParentRef == nil {
			parent = p.root
		}
		p.unmappedInterpolatedKeys[ref.Timeline].object.setWithBone(p.interpolatedKeys[ref.Timeline].object)
		p.unmappedInterpolatedKeys[ref.Timeline].object.unmapCoordinates(parent)
		p.unmapObjects(ref)
	}
	for i := range p.currentKey.ObjectRefs {
		ref := p.currentKey.ObjectRefs[i]
		if ref.ParentRef != base && base != nil {
			continue
		}
		parent := p.unmappedInterpolatedKeys[ref.ParentRef.Timeline].object
		if ref.ParentRef == nil {
			parent = p.root
		}
		p.unmappedInterpolatedKeys[ref.Timeline].object.setWithBone(p.interpolatedKeys[ref.Timeline].object)
		p.unmappedInterpolatedKeys[ref.Timeline].object.unmapCoordinates(parent)
	}
}

func (p *EntityPlayer) setEntity(entity *Entity) {
	if entity == nil {
		fmt.Println("Error: Entity cannot be nil")
		return
	}
	p.entity = entity
	maxTimelineKeys := entity.MaxNumTimelines
	p.tempInterpolatedKeys = make([]*TimelineKey, maxTimelineKeys)
	p.tempUnmappedInterpolatedKeys = make([]*TimelineKey, maxTimelineKeys)
	p.interpolatedKeys = make([]*TimelineKey, maxTimelineKeys)
	p.unmappedInterpolatedKeys = make([]*TimelineKey, maxTimelineKeys)

	for i := 0; i < maxTimelineKeys; i++ {
		key := MakeTimelineKey(i)
		keyU := MakeTimelineKey(i)
		key.setObject(MakeTimelineKeyBone())
		keyU.setObject(MakeTimelineKeyBone())
		p.interpolatedKeys[i] = key
		p.unmappedInterpolatedKeys[i] = keyU
		p.objToTimeline[keyU.object] = keyU
	}
	p.tempInterpolatedKeys = p.interpolatedKeys
	p.tempUnmappedInterpolatedKeys = p.unmappedInterpolatedKeys
	p.setAnimation(entity.getAnimationByIndex(0))
}

func (p *EntityPlayer) getEntity() *Entity {
	return p.entity
}

func (p *EntityPlayer) setAnimation(animation *Animation) {
	prevAnim := p.animation
	if animation == p.animation {
		return
	}
	if animation == nil {
		fmt.Println("Error: Animation cannot be nil")
		return
	}
	if !p.entity.containsAnimation(animation) && animation.Id != -1 {
		fmt.Println("Error: Animation should be in the same entity as the current one")
		return
	}
	if animation != p.animation {
		p.time = 0
	}
	p.animation = animation
	tempTime := p.time
	p.time = 0
	p.Update(0)
	p.time = tempTime
	for i := range p.listeners {
		p.listeners[i].animationChanged(prevAnim, p.animation)
	}
}

func (p *EntityPlayer) SetAnimationByName(name string) {
	p.setAnimation(p.entity.getAnimationByName(name))
}

func (p *EntityPlayer) SetAnimationByIndex(index int) {
	p.setAnimation(p.entity.getAnimationByIndex(index))
}

func (p *EntityPlayer) GetAnimation() *Animation {
	return p.animation
}

func (p *EntityPlayer) GetCurrentKey() *MainlineKey {
	return p.currentKey
}

func (p *EntityPlayer) getTime() int {
	return p.time
}

func (p *EntityPlayer) setTime(time int) *EntityPlayer {
	p.time = time
	p.increaseTime(0)
	return p
}

func (p *EntityPlayer) setScale(scale float64) *EntityPlayer {
	p.root.Scale[0] = scale * float64(p.flippedX())
	p.root.Scale[1] = scale * float64(p.flippedY())
	return p
}

func (p *EntityPlayer) SetFlip(x bool, y bool) *EntityPlayer {
	if x {
		p.SetFlipX()
	}
	if y {
		p.SetFlipY()
	}
	return p
}

func (p *EntityPlayer) SetFlipX() *EntityPlayer {
	p.root.Scale[0] *= -1
	return p
}

func (p *EntityPlayer) SetFlipY() *EntityPlayer {
	p.root.Scale[1] *= -1
	return p
}

func (p *EntityPlayer) flippedX() int {
	return int(signum(p.root.Scale.X()))
}

func (p *EntityPlayer) flippedY() int {
	return int(signum(p.root.Scale.Y()))
}

func (p *EntityPlayer) setPosition(x float64, y float64) *EntityPlayer {
	p.rootIsDirty = true
	p.position.SetCoords(x, y)
	return p
}

func (p *EntityPlayer) getX() float64 {
	return p.position.X()
}

func (p *EntityPlayer) getY() float64 {
	return p.position.Y()
}

func (p *EntityPlayer) SetAngle(angle float64) *EntityPlayer {
	p.rootIsDirty = true
	p.angle = angle
	return p
}

func (p *EntityPlayer) getAngle() float64 {
	return p.angle
}

func (p *EntityPlayer) setPivot(x float64, y float64) *EntityPlayer {
	p.rootIsDirty = true
	p.pivot.SetCoords(x, y)
	return p
}

func (p *EntityPlayer) EnableCharacterMap(mapName string) {
	p.enabledCharacterMaps[mapName] = p.entity.getCharacterMap(mapName)
}

func (p *EntityPlayer) DisableCharacterMap(mapName string) {
	delete(p.enabledCharacterMaps, mapName)
}

// Helper function for drawing the sprite

func (p *EntityPlayer) GetNumObjectsToDraw() int {
	return len(p.currentKey.ObjectRefs)
}

func (p *EntityPlayer) GetKeyObjectToDraw(index int) *TimelineKeyObject {
	objRef := p.currentKey.ObjectRefs[index]
	return p.unmappedInterpolatedKeys[objRef.Timeline].object
}

func (p *EntityPlayer) GetMappedFileIndexForKeyObject(object *TimelineKeyObject) int {
	fileIndex := object.fileIndex
	for _, v := range p.enabledCharacterMaps {
		if _, ok := v.FilesMapping[fileIndex]; ok {
			fileIndex = v.FilesMapping[fileIndex]
		}
	}
	return fileIndex
}
