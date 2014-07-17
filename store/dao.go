package store

type DAO interface {
	SaveModel(data map[string]interface{})
	UpdateModel(data map[string]interface{})
	DeleteModel(id int64)
	GetAll() []map[string]interface{}
}

//if struct can convert itself to interface
//not using at the moment
type Mappable interface {
	StructToMap() map[string]interface{}
}
