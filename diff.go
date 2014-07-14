//Contains diff function that compares DOMs and produces difference
package drops

import (
	"fmt"

	"github.com/mkasner/drops/element"
)

type Patch struct {
	Element string `json:"element"`
	Payload string `json:"payload"`
}

const capacity = 1024

//Produces diff patch of two DOMs that can be used to update client side dom
func Diff(view1 *element.View, view2 *element.View) []Patch {
	// fmt.Printf("\n\n\nview1: %v\n", *view1)
	// fmt.Printf("view2: %v\n", *view2)
	cont := true
	patches := make([]Patch, 0)
	var newPatches []Patch
	fmt.Printf("Testing views: %s vs %s\n", view1.Template, view2.Template)
	fmt.Println("Testing length...")
	newPatches, cont = testLength(view1, view2)
	patches = append(patches, newPatches...)
	if cont {
		fmt.Println("Testing Equality...")

		newPatches, cont = testEquality(view1, view2)
		patches = append(patches, newPatches...)
	}
	if cont {
		fmt.Println("Testing ChildrenEquality...")
		newPatches, cont = testChildrenEquality(view1, view2)
		patches = append(patches, newPatches...)
	}

	return patches
}

//Rules
//test for length of children, if length different, generate patch
func testLength(view1 *element.View, view2 *element.View) ([]Patch, bool) {
	patches := make([]Patch, 0)
	cont := true //continue
	if len(view1.Children) != len(view2.Children) {
		buff := Render(view2)
		patchElement := view2.InjectInto

		patch := &Patch{Element: patchElement, Payload: buff.String()}
		patches = append(patches, *patch)
		cont = false
	}
	if len(patches) > 0 {
		fmt.Println("\nGenerated patches on testLength")
	}
	return patches, cont
}

//test to see if views are equal
func testEquality(view1 *element.View, view2 *element.View) ([]Patch, bool) {
	patches := make([]Patch, 0)
	cont := true //continue
	same := true
	if view1 != nil && view2 != nil {
		same = true
	}
	if same && view1.Model != view2.Model {
		same = false
	}
	if same && view1.Template != view2.Template {
		same = false
	}
	if !same {
		buff := Render(view2)
		patchElement := view2.InjectInto

		patch := &Patch{Element: patchElement, Payload: buff.String()}
		patches = append(patches, *patch)
		// fmt.Printf("\nview1: %+v\n", view1)
		// fmt.Printf("\nview2: %+v\n", view2)
		// fmt.Printf("\nequal: %+v\n", view1.Model == view2.Model)

		cont = false
	}
	if len(patches) > 0 {
		fmt.Println("\nGenerated patches on testEquality")
	}
	return patches, cont
}

//test to see if elements of children are equal for every element
func testChildrenEquality(view1 *element.View, view2 *element.View) ([]Patch, bool) {
	patches := make([]Patch, 0)
	cont := true //continue
	biggerChildren := len(view1.Children)
	if view2.Children != nil && len(view2.Children) > biggerChildren {
		biggerChildren = len(view2.Children)
	}
	for i := 0; i < biggerChildren; i++ {
		var v1 *element.View
		var v2 *element.View
		if i < len(view1.Children) {
			v1 = view1.Children[i]
		}
		if i < len(view2.Children) {
			v2 = view2.Children[i]
		}
		if v1 != v2 {
			buff := Render(v2)
			patchElement := v2.InjectInto

			patch := &Patch{Element: patchElement, Payload: buff.String()}
			patches = append(patches, *patch)

			cont = false
		} else {
			patches = append(patches, Diff(v1, v2)...)
		}
	}
	return patches, cont
}

// func Diff(view1 *element.View, view2 *element.View) []Patch {
// 	// fmt.Printf("\n\n\nview1: %v\n", *view1)
// 	// fmt.Printf("view2: %v\n", *view2)
// 	patches := make([]Patch, 0)

// 	var biggerChildren int
// 	//Get same level children pointer addresses
// 	var children1 []*element.View
// 	if view1.Children != nil {
// 		biggerChildren = len(children1)
// 		if view2.Children != nil && len(view2.Children) > biggerChildren {
// 			biggerChildren = len(view2.Children)
// 		}
// 		children1 = make([]*element.View, biggerChildren)
// 		for i, e := range view1.Children {
// 			// view := e.Value.(*element.View)
// 			children1[i] = e
// 		}
// 	}
// 	var children2 []*element.View
// 	if view2.Children != nil {
// 		children2 = make([]*element.View, biggerChildren)
// 		for i, e := range view2.Children {
// 			// view := e.Value.(*element.View)
// 			children2[i] = e
// 		}
// 	}

// 	// if children are the same go one level deeper
// 	// fmt.Printf("children1: %v\n", children1)
// 	// fmt.Printf("children2: %v\n", children2)
// 	notSame := false

// 	for i := 0; i < biggerChildren; i++ {
// 		if children1[i] != children2[i] {
// 			notSame = true
// 		}
// 	}
// 	if notSame {
// 		// fmt.Printf("same: %v\n", children1 == children2)
// 		for i := 0; i < capacity; i++ {
// 			if children1[i] != nil && children2[i] != nil {
// 				view1 := children1[i]
// 				view2 := children2[i]
// 				childPatches := Diff(view1, view2)
// 				patches = append(patches, childPatches...)
// 			}
// 		}
// 	} else {
// 		// fmt.Printf("Children not same: %v\n", children1 == children2)

// 		//Children are not the same - find out the differences and render them to patch
// 		differences := list.New()
// 		for i := 0; i < capacity; i++ {
// 			var view1 *element.View
// 			var view2 *element.View
// 			if len(children1)-1 >= i && children1[i] != nil {
// 				view1 = children1[i]
// 			}
// 			if len(children2)-1 >= i && children1[i] != nil {
// 				view2 = children2[i]
// 			}
// 			// if children1[i] != nil && children2[i] != nil {
// 			// view1 = children1[i]
// 			// view2 = children2[i]
// 			// //Check for modelStruct equality
// 			// modelStruct1 := view1.ModelStruct
// 			// modelStruct2 := view2.ModelStruct
// 			// modelStructEq := reflect.DeepEqual(modelStruct1, modelStruct2)
// 			// if modelStructEq {
// 			// 	//modelstructs are equal go level deeper
// 			// 	//They are the same, try deeper level
// 			// 	childPatches := Diff(view1, view2)
// 			// 	patches = append(patches, childPatches...)
// 			// } else {
// 			// 	fmt.Printf("different: %v\n", view2)
// 			// 	differences.PushBack(view2)
// 			// }
// 			// map1 := make(map[string]interface{})
// 			// //Copy maps and ignore content key
// 			// for k, v := range view1.Model {
// 			// 	if k != "Content" {
// 			// 		map1[k] = v
// 			// 	}
// 			// }
// 			// map2 := make(map[string]interface{})
// 			// for k, v := range view2.Model {
// 			// 	if k != "Content" {
// 			// 		map2[k] = v
// 			// 	}
// 			// }

// 			//Bypassing comparison
// 			modelEq := false
// 			if view1 != nil && view2 != nil {
// 				// fmt.Println("Comparing views...")
// 				// fmt.Printf("View1:%+v %+v\n", &*view1, view1)
// 				// fmt.Printf("View2:%+v %+v\n", &*view2, view2)
// 				modelEq = reflect.DeepEqual(view1.Model, view2.Model)
// 			}
// 			// if children1[i] != children2[i] {
// 			if !modelEq {
// 				if view2 != nil {
// 					// fmt.Printf("different: %+v\n", view2)
// 					differences.PushBack(view2)
// 				}
// 			} else {
// 				//They are the same, try deeper level
// 				childPatches := Diff(view1, view2)
// 				patches = append(patches, childPatches...)
// 			}

// 			// }
// 		}
// 		for e := differences.Front(); e != nil; e = e.Next() {
// 			view := e.Value.(*element.View)
// 			buff := Render(view)
// 			patchElement := view.InjectInto
// 			if view.Return != "" {
// 				patchElement = view.Provides
// 			}
// 			patch := &Patch{Element: patchElement, Payload: buff.String()}
// 			patches = append(patches, *patch)
// 		}

// 	}
// 	// buffer := view1.Render()
// 	// patch := &Patch{Element: "html", Payload: buffer.String()}
// 	// patches[0] = *patch
// 	// fmt.Printf("\nPatches size: %+v\n", len(patches))
// 	return patches
// }
