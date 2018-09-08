package spriter

type DrawerExample struct {
	data       *Model
	path       string
}

func MakeSpriterDrawer(data *Model, path string) *DrawerExample {
	d := &DrawerExample{
		data: data,
		path: path,
	}
	// Initialize other data structures here
	return d
}

func (d *DrawerExample) LoadAssets() {
	//for _, file := range d.data.Files {
	//	Load texture file 'file.Name'
	//}
}

func (d *DrawerExample) Draw(p *EntityPlayer) {
	var file *File
	numObjectsToDraw := p.GetNumObjectsToDraw()
	for i :=0; i< numObjectsToDraw; i++ {
		o := p.GetKeyObjectToDraw(i)
		if o.ObjectType() == "object" {
			fileIndex := p.GetMappedFileIndexForKeyObject(o)
			file = d.data.Files[fileIndex]
			if file == nil {
				continue
			}
			// texture <-- Get the loaded texture for 'file.Name'
			// pivotX <-- o.Pivot.X() * float64(file.Width)
			// pivotY <--  o.Pivot.Y() * float64(file.Height)
			// scale <--  o.Scale
			// angle <-- o.Angle
			// position <-- o.Position
			// Draw!
		}
	}
}
