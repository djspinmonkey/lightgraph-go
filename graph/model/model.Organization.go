package model

type Organization struct {
	ID   string `json:"data.attributes.id"`
	Name string `json:"data.attributes.name"`
}

func (o *Organization) Project(id string) *Project {
	project, err := FetchProject(o.ID, id)
	if err != nil {
		panic(err)
	}

	return project
}
