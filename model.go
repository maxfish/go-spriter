package spriter

type Model struct {
	SconVersion      string    `xml:"scon_version,attr"`
	Generator        string    `xml:"generator,attr"`
	GeneratorVersion string    `xml:"generator_version,attr"`
	Entities         []*Entity `xml:"entity"`
	Folders          []*Folder `xml:"folder"`
	nameToEntity     map[string]*Entity
	Files            map[int]*File
}

type Folder struct {
	Id    int     `xml:"id,attr"`
	Name  string  `xml:"name,attr"`
	Files []*File `xml:"file"`
}

func (m *Model) GetEntityIndex(name string) int {
	for i := range m.Entities {
		e := m.Entities[i]
		if e.Name == name {
			return e.Id
		}
	}
	return -1
}

func (m *Model) GetEntityByName(name string) *Entity {
	if m.nameToEntity == nil {
		m.nameToEntity = make(map[string]*Entity)
		for i := range m.Entities {
			t := m.Entities[i]
			m.nameToEntity[t.Name] = t
		}
	}
	return m.nameToEntity[name]
}

func (m *Model) GetFile(fileIndex int) *File {
	if fileIndex < 0 {
		return nil
	} else {
		return m.Files[fileIndex]
	}
}
