package spriter

import (
	"os"
	"fmt"
	"io/ioutil"
	"math"
	"encoding/xml"
)

func NewSpriterModelFromFile(fileName string) *Model {
	fmt.Println("=== Reading SCML ===")

	// Open the xml file
	xmlFile, err := os.Open(fileName)
	defer xmlFile.Close()
	if err != nil {
		fmt.Println(err)
	}

	model := &Model{}

	// Read the data into the model
	byteData, _ := ioutil.ReadAll(xmlFile)
	err = xml.Unmarshal(byteData, &model)
	if err != nil {
		fmt.Println(err)
	}

	initializeData(model)
	return model
}

func initializeData(data *Model) {
	for i := range data.Entities {
		// File mappings
		data.Files = make(map[int]*File)
		for i := range data.Folders {
			folder := data.Folders[i]
			for j := range folder.Files {
				file := folder.Files[j]
				data.Files[FolderAndFileToFileIndex(i, j)] = file
			}
		}

		// Entities
		entity := data.Entities[i]
		for j := range entity.CharacterMaps {
			m := entity.CharacterMaps[j]
			m.FilesMapping = make(map[int]int)
			if m.Name == "" {
				m.Name = fmt.Sprintf("charMap%d", j)
			}
			for k := range m.Maps {
				mapping := m.Maps[k]
				if mapping.TargetFile == nil || mapping.TargetFolder == nil {
					m.FilesMapping[FolderAndFileToFileIndex(mapping.Folder, mapping.File)] = -1
				} else {
					m.FilesMapping[FolderAndFileToFileIndex(mapping.Folder, mapping.File)] = FolderAndFileToFileIndex(*mapping.TargetFolder, *mapping.TargetFile)
				}
			}
		}

		// Animations
		for j := range entity.Animations {
			a := entity.Animations[j]
			a.initialize()
			a.Looping = optionalBool(a.XMLLooping, true)

			// Updates the max number of timelines for the entity
			if len(a.Timelines) > entity.MaxNumTimelines {
				entity.MaxNumTimelines = len(a.Timelines)
			}

			// Mainline
			for k := range a.Mainline.Keys {
				key := a.Mainline.Keys[k]
				curve := MakeCurve()
				key.curve = curve
				if key.CurveType == nil {
					key.curve.curveType = TypeLinear
				} else {
					key.curve.curveType = getCurveTypeFromName(*key.CurveType)
				}
				for z := range key.BoneRefs {
					ref := key.BoneRefs[z]
					if ref.Parent != nil && *ref.Parent >= 0 && *ref.Parent < len(key.BoneRefs) {
						ref.ParentRef = key.BoneRefs[*ref.Parent]
					}
				}
				for z := range key.ObjectRefs {
					ref := key.ObjectRefs[z]
					if ref.Parent != nil && *ref.Parent >= 0 && *ref.Parent < len(key.BoneRefs) {
						ref.ParentRef = key.BoneRefs[*ref.Parent]
					}
				}
			}
			// Timelines
			for k := range a.Timelines {
				timeline := a.Timelines[k]
				o := entity.getInfoByName(timeline.Name)
				if o == nil {
					o = MakeObjectInfo(timeline.Name, TypeSprite, 0, 0)
				}
				timeline.objectInfo = o
				if timeline.ObjectType == "" {
					timeline.ObjectType = TypeSprite
				}
				timeline.objectInfo.Type = timeline.ObjectType

				for z := range timeline.Keys {
					key := timeline.Keys[z]
					curve := MakeCurveWithType(getCurveTypeFromName(key.CurveType))
					curve.constraints = [4]float64{key.C1, key.C2, key.C3, key.C4}
					key.Curve = curve
					key.Spin = optionalInt(key.XMLSpin, 1)

					if key.XMLDataBone != nil {
						// This is a bone
						key.object = key.XMLDataBone
						key.object.objectType = TypeBone
						key.object.Pivot = MakePoint(
							optionalFloat(key.object.XMLPivotX, 0),
							optionalFloat(key.object.XMLPivotY, 0.5),
						)
						key.XMLDataBone = nil
					} else {
						// This could be a Sprite, a Point or a Box
						key.object = key.XMLDataObject
						key.object.objectType = timeline.ObjectType
						o := key.object
						o.fileIndex = FolderAndFileToFileIndex(o.Folder, o.File)
						f := data.Files[o.fileIndex]
						if timeline.objectInfo != nil {
							timeline.objectInfo.Width = float64(f.Width)
							timeline.objectInfo.Height = float64(f.Height)
						}
						key.object.Pivot = MakePoint(
							optionalFloat(key.object.XMLPivotX, f.PivotX),
							optionalFloat(key.object.XMLPivotY, f.PivotY),
						)
						key.XMLDataObject = nil
					}
					key.object.Position = MakePoint(key.object.XMLX, key.object.XMLY)
					key.object.Scale = MakePoint(
						optionalFloat(key.object.XMLScaleX, 1),
						optionalFloat(key.object.XMLScaleY, 1),
					)
					// Convert degrees to radians
					key.object.Angle = (math.Pi * key.object.Angle) / 180.0
				}
			}
		}
	}
}

func optionalFloat(value *float64, defValue float64) float64 {
	if value != nil {
		return *value
	} else {
		return defValue
	}
}

func optionalInt(value *int, defValue int) int {
	if value != nil {
		return *value
	} else {
		return defValue
	}
}

func optionalBool(value *bool, defValue bool) bool {
	if value != nil {
		return *value
	} else {
		return defValue
	}
}
