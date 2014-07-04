package drops

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
)

var stores map[string]*Store

type Store struct {
	Type *reflect.Type
	dao  DAO
}

func NewStore() {
	stores = make(map[string]*Store)
}

func AddStore(model interface{}, dao DAO) {
	typ := reflect.TypeOf(model)
	stores[typ.Name()] = &Store{Type: &typ, dao: dao}
	fmt.Printf("Created store for model %s\n", typ.Name())
}

func AddModel(model string, data map[string]interface{}) {
	if store, ok := stores[model]; ok {
		fmt.Printf("Found store for model %s\n", model)
		store.dao.SaveModel(data)
	} else {
		log.Printf("Store for model %s not initialized.\n", model)
	}
}

func GetAll(model string) []map[string]interface{} {
	if store, ok := stores[model]; ok {
		fmt.Printf("Found store for model %s\n", model)
		return store.dao.GetAll()
	} else {
		log.Printf("Store for model %s not initialized.\n", model)
	}
	return nil
}

func SaveModel(model string, data map[string]interface{}) {
	if store, ok := stores[model]; ok {
		fmt.Printf("Found store for model %s\n", model)
		store.dao.UpdateModel(data)
	} else {
		log.Printf("Store for model %s not initialized.\n", model)
	}
}

func DeleteModel(model string, data map[string]interface{}) {
	if store, ok := stores[model]; ok {
		fmt.Printf("Found store for model %s\n", model)
		id, err := strconv.ParseInt(data["id"].(string), 0, 64)
		if err != nil {
			log.Printf("Error parsing id %v", err)
		}
		store.dao.DeleteModel(id)
	} else {
		log.Printf("Store for model %s not initialized.\n", model)
	}
}
