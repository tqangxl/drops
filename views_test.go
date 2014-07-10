//Testing of component creation
// If i want live reload i can use chromix
//chromix with "file:///.*/component/.*.html" reload
package drops

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/mkasner/drops"
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

func TestHandlersSuite(t *testing.T) {
	testView(t)

}

//Creating view ad writing to tempfile
// tempfile can be livereload with command chromix with "file:///.*/component/.*.html" reload
func testView(t *testing.T) {
	fmt.Println("Testing views")
	var project ProjectModel
	typ := reflect.TypeOf(project)
	drops.NewStore()
	var node NodeModel
	drops.AddStore(node, &NodeDAO{})
	//Rules for generating fields
	rules := map[string]string{
		"Id":     `ignore:"true"`,
		"Origin": `ignore:"false" foreign:"NodeModel"`,
	}

	view := drops.New(typ, rules, "", "")
	r := view.Render()
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
	result := make([]drops.Model, 3)
	result[0] = drops.Model{"Id": 1, "Name": "localhost"}
	result[1] = drops.Model{"Id": 1, "Name": "sweden"}
	result[2] = drops.Model{"Id": 1, "Name": "litva"}

	return result
}
func (n *NodeDAO) DeleteModel(id int64) {
}
func (n *NodeDAO) SaveModel(data map[string]interface{}) {
}
