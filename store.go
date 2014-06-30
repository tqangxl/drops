package drops

import (
	"fmt"
	"github.com/AduroIdea/dplrlib"
	"log"
	"reflect"
)

var stores map[string]*Store

type Store struct {
	Type *reflect.Type
	dao  dplrlib.DAO
}

func NewStore() {
	stores = make(map[string]*Store)
}

func AddStore(model interface{}, dao dplrlib.DAO) {
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
