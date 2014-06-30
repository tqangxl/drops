//Contains diff function that compares DOMs and produces difference
package drops

import (
	"container/list"
	"fmt"
	"reflect"
)

type Patch struct {
	Element string `json:"element"`
	Payload string `json:"payload"`
}

const capacity = 1024

//Produces diff patch of two DOMs that can be used to update client side dom
func Diff(view1 *View, view2 *View) []Patch {
	// fmt.Printf("view1: %v\n", view1)
	// fmt.Printf("view2: %v\n", view2)
	patches := make([]Patch, 0)
	//Get same level children pointer addresses
	var children1 [capacity]interface{}
	if view1.Children != nil {
		i := 0
		for e := view1.Children.Front(); e != nil; e = e.Next() {
			// view := e.Value.(*View)
			children1[i] = e.Value
			i++
		}
	} else {

	}
	var children2 [capacity]interface{}
	if view2.Children != nil {
		i := 0
		for e := view2.Children.Front(); e != nil; e = e.Next() {
			// view := e.Value.(*View)
			children2[i] = e.Value
			i++
		}
	}
	// if children are the same go one level deeper
	// fmt.Printf("children1: %v\n", children1)
	// fmt.Printf("children2: %v\n", children2)
	if children1 == children2 {
		fmt.Printf("same: %b\n", children1 == children2)
		for i := 0; i < capacity; i++ {
			if children1[i] != nil && children2[i] != nil {
				view1 = children1[i].(*View)
				view2 = children2[i].(*View)
				childPatches := Diff(view1, view2)
				patches = append(patches, childPatches...)
			}
		}
	} else {
		//Children are not the same - find out the differences and render them to patch
		differences := list.New()
		for i := 0; i < capacity; i++ {
			var view1 *View
			var view2 *View
			if children1[i] != nil && children2[i] != nil {
				view1 = children1[i].(*View)
				view2 = children2[i].(*View)
				map1 := make(map[string]interface{})
				//Copy maps and ignore content key
				for k, v := range view1.Model {
					if k != "Content" {
						map1[k] = v
					}
				}
				map2 := make(map[string]interface{})
				for k, v := range view2.Model {
					if k != "Content" {
						map2[k] = v
					}
				}
				modelEq := reflect.DeepEqual(map1, map2)
				// if children1[i] != children2[i] {
				if !modelEq {
					fmt.Printf("different: %v\n", view2)
					differences.PushBack(view2)
				} else {
					//They are the same, try deeper level
					childPatches := Diff(view1, view2)
					patches = append(patches, childPatches...)
				}
			}
		}
		for e := differences.Front(); e != nil; e = e.Next() {
			view := e.Value.(*View)
			buff := view.Render()
			patch := &Patch{Element: view.InjectInto, Payload: buff.String()}
			patches = append(patches, *patch)
		}

	}
	// buffer := view1.Render()
	// patch := &Patch{Element: "html", Payload: buffer.String()}
	// patches[0] = *patch
	return patches
}
