package spriter

import "fmt"

const (
	NumBitsPerFolder = 10
)

type File struct {
	Id     int     `xml:"id,attr"`
	Name   string  `xml:"name,attr"`
	Width  int     `xml:"width,attr"`
	Height int     `xml:"height,attr"`
	PivotX float64 `xml:"pivot_x,attr"`
	PivotY float64 `xml:"pivot_y,attr"`
	// Atlas - Currently not supported
	//AtlasX        float64 `xml:"ax,attr"`
	//AtlasY        float64 `xml:"ay,attr"`
	//AtlasXOffset  float64 `xml:"axoff,attr"`
	//AtlasYOffset  float64 `xml:"ayoff,attr"`
	//AtlasWidth    float64 `xml:"aw,attr"`
	//AtlasHeight   float64 `xml:"ah,attr"`
	//AtlasRotation string  `xml:"aror,attr"`
}

func (f *File) String() string {
	return fmt.Sprintf("[id: %d, name: %s, size: %dx%d, pivot: %f,%f", f.Id, f.Name, f.Width, f.Height, f.PivotX, f.PivotY)
}

func FolderAndFileToFileIndex(folder int, file int) int {
	return folder<<NumBitsPerFolder + file
}
