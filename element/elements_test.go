//Testing of component creation
// If i want live reload i can use chromix
//chromix with "file:///.*/component/.*.html" reload
package element

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/mkasner/drops/element"
)

type ProjectModel struct {
	Id      int64
	Name    string `sql:"type:varchar(100);unique"`
	GitName string `sql:"type:varchar(100);" json:"git_name"`
	Origin  string `sql:"type:varchar(100);"`
}
type NodeModel struct {
	Id       int64
	Name     string `sql:"type:varchar(100);unique"`
	User     string `sql:"type:varchar(100);"`
	Host     string `sql:"type:varchar(100);"`
	Port     string `sql:"type:varchar(100);"`
	RepoPath string `sql:"type:varchar(255);" json:"repo_path"`
	AppPath  string `sql:"type:varchar(255);" json:"app_path"`
}

func TestSuite(t *testing.T) {
	// testView(t)
	testAdd(t)
	testAddToView(t)
}

//Testing Add function
//Test if length of children is equal expected
func testAdd(t *testing.T) {
	fmt.Println("Testing Add(DOM, view)")
	dom1 := &element.DOM{View: element.View{Children: make([]*element.View, 0), Template: "views/base.html", Provides: "html", Model: &element.Model{MAP: make(map[string]interface{})}}, Id: "1"}
	dom1.IdMap = make(map[string]*element.View)
	dom1.IdMap[dom1.View.Provides] = &dom1.View
	head := &element.Head{View: element.View{Template: "views/head.html", Provides: "head", Return: "Head", Model: &element.Model{MAP: make(map[string]interface{})}}}
	body := &element.Body{View: element.View{Template: "views/body.html", Provides: "body", Return: "Body",
		Model: &element.Model{MAP: map[string]interface{}{
			"Content": "Body1",
		}}}}
	dom1 = Add(dom1, &head.View)
	dom1 = Add(dom1, &body.View)
	expected := 2
	if len(dom1.View.Children) != expected {
		t.Errorf("Children not added for first level: %v != %v\n", len(dom1.Children), expected)
		t.Fail()
	}

	//Test adding to deeper level, expected children should be the same
	header := &Header{View: &element.View{Template: "views/header.html", Provides: "#header", InjectInto: "body", Model: &element.Model{MAP: make(map[string]interface{})}}}
	dom1 = Add(dom1, header.View)
	if len(dom1.View.Children) != expected {
		t.Errorf("Children should be the same as first expected: %v != %v\n", len(dom1.Children), expected)
		t.Fail()
	}

	expected = 1
	if len(body.View.Children) != expected {
		t.Errorf("Children not added for body element: %v != %v\n", len(body.View.Children), expected)
		t.Fail()
	}

}

func testAddToView(t *testing.T) {
	fmt.Println("Testing AddToView(View, newVIew)")
	dom1 := &element.DOM{View: element.View{Children: make([]*element.View, 0), Template: "views/base.html", Provides: "html", Model: &element.Model{MAP: make(map[string]interface{})}}, Id: "1"}
	dom1.IdMap = make(map[string]*element.View)
	dom1.IdMap[dom1.View.Provides] = &dom1.View
	head := &element.Head{View: element.View{Template: "views/head.html", Provides: "head", Return: "Head", Model: &element.Model{MAP: make(map[string]interface{})}}}
	body := &element.Body{View: element.View{Template: "views/body.html", Provides: "body", Return: "Body",
		Model: &element.Model{MAP: map[string]interface{}{
			"Content": "Body1",
		}}}}
	dom1 = Add(dom1, &head.View)
	dom1 = Add(dom1, &body.View)
	expected := 2
	if len(dom1.View.Children) != expected {
		t.Errorf("Children not added: %v != %v\n", len(dom1.Children), expected)
		t.Fail()
	}
	header := &Header{View: &element.View{Template: "views/header.html", Provides: "#header", InjectInto: "body", Model: &element.Model{MAP: make(map[string]interface{})}}}
	AddToView(&body.View, header.View)
	expected = 1
	if len(body.View.Children) != expected {
		t.Errorf("Children not added for body element: %v != %v\n", len(body.View.Children), expected)
		t.Fail()
	}
}

func TestCopyDOM(t *testing.T) {
	fmt.Println("Testing CopyDOM(element.DOM)")
	dom1 := &element.DOM{View: element.View{Children: make([]*element.View, 0), Template: "views/base.html", Provides: "html", Model: &element.Model{MAP: make(map[string]interface{})}}, Id: "1"}
	dom1.IdMap = make(map[string]*element.View)
	dom1.IdMap[dom1.View.Provides] = &dom1.View
	head := &element.Head{View: element.View{Template: "views/head.html", Provides: "head", Return: "Head", Model: &element.Model{MAP: make(map[string]interface{})}}}
	body := &element.Body{View: element.View{Template: "views/body.html", Provides: "body", Return: "Body",
		Model: &element.Model{MAP: map[string]interface{}{
			"Content": "Body1",
		}}}}
	dom1 = Add(dom1, &head.View)
	dom1 = Add(dom1, &body.View)
	expected := 2
	if len(dom1.View.Children) != expected {
		t.Errorf("Children not added: %v != %v\n", len(dom1.Children), expected)
		t.Fail()
	}
	header := &Header{View: &element.View{Template: "views/header.html", Provides: "#header", InjectInto: "body", Model: &element.Model{MAP: make(map[string]interface{})}}}
	AddToView(&body.View, header.View)
	expected = 1
	if len(body.View.Children) != expected {
		t.Errorf("Children not added for body element: %v != %v\n", len(body.View.Children), expected)
		t.Fail()
	}
	dom2 := CopyDom(*dom1)
	if dom1.Id == dom2.Id {
		t.Errorf("Id not changed: %s == %s\n", dom1.Id, dom2.Id)
		t.Fail()
	}

	equals(t, len(dom1.Children), len(dom2.Children))
}

func TestCopyView(t *testing.T) {
	fmt.Println("Testing CopyView(element.View, map[string]*element.View)")
	//Not implemented yet
}

//Creating view ad writing to tempfile
// tempfile can be livereload with command chromix with "file:///.*/component/.*.html" reload
func testView(t *testing.T) {
	fmt.Println("Testing views")
	var project ProjectModel
	typ := reflect.TypeOf(project)
	NewStore()
	var node NodeModel
	AddStore(node, &NodeDAO{})
	//Rules for generating fields
	rules := map[string]string{
		"Id":     `ignore:"true"`,
		"Origin": `ignore:"false" foreign:"NodeModel"`,
	}

	view := New(typ, rules, "", "")
	r := Render(view)
	// open output file
	fo, err := os.Create("/tmp/component.html")
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	fmt.Printf("View %+v", r.String())
	if _, err := fo.Write(r.Bytes()); err != nil {
		panic(err)
	}
}

type NodeDAO struct{}

func (n *NodeDAO) GetAll() []map[string]interface{} {
	result := make([]map[string]interface{}, 3)
	result[0] = map[string]interface{}{"Id": 1, "Name": "localhost"}
	result[1] = map[string]interface{}{"Id": 1, "Name": "sweden"}
	result[2] = map[string]interface{}{"Id": 1, "Name": "litva"}

	return result
}
func (n *NodeDAO) DeleteModel(id int64) {
}
func (n *NodeDAO) SaveModel(data map[string]interface{}) {
}
func (n *NodeDAO) UpdateModel(data map[string]interface{}) {
}

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
