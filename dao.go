package drops

type DAO interface {
	SaveModel(data map[string]interface{})
	UpdateModel(data map[string]interface{})
	DeleteModel(id int64)
	GetAll() []map[string]interface{}
}

//type for holding result as map
//not using at the moment
type Model map[string]interface{}

//if struct can convert itself to interface
//not using at the moment
type Mappable interface {
	StructToMap() map[string]interface{}
}
