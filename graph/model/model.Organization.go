package model

type Organization struct {
	ID   string `json:"data.attributes.id"`
	Name string `json:"data.attributes.name"`
}

func (o *Organization) Project(id string) *Project {
	// The commented code below totally works, but since Project only has an ID and a Name, and they're always the same,
	// we can just make a Project struct with the provided ID/Name and return it to save ourselves a network call. If
	// we wanted to verify that the Project actually exists, we could uncomment the code below.
	//
	//project, err := FetchProject(o, id)
	//if err != nil {
	//	panic(err)
	//}
	//
	//return project

	return &Project{
		ID:           id,
		Name:         id,
		Organization: o,
	}
}
