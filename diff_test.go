package drops

import (
	"fmt"
	"testing"
	"text/template"

	"github.com/mkasner/drops/element"
)

type Header struct {
	*element.View
}

type DeploymentView struct {
	*element.View
}
type ProjectView struct {
	*element.View
}

type Project struct {
	Id          int64
	Name        string `sql:"type:varchar(100);unique"`
	GitName     string `sql:"type:varchar(100);" json:"git_name"`
	Origin      string `sql:"type:varchar(100);"`
	Deployments []Deployment
}
type Deployment struct {
	Id             int64
	Name           string `sql:"type:varchar(100);"`
	Environment    string `sql:"type:varchar(100);" json:"from_path"`
	FromPath       string `sql:"type:varchar(100);"`
	Target         string `sql:"type:varchar(100);"`
	Checkout       string `sql:"type:varchar(100);"`
	Node           Node
	Protocol       Protocol
	NodeProtocolId int64
	Project        Project
	ProjectId      int64
}

type Protocol struct {
	Id   int64
	Name string `sql:"type:varchar(100);"`
}

type Node struct {
	Id       int64
	Name     string `sql:"type:varchar(100);unique"`
	User     string `sql:"type:varchar(100);"`
	Host     string `sql:"type:varchar(100);"`
	Port     string `sql:"type:varchar(100);"`
	RepoPath string `sql:"type:varchar(255);" json:"repo_path"`
	AppPath  string `sql:"type:varchar(255);" json:"app_path"`
}

func testDiffSuite(t *testing.T) {
	// testEchoHandler(t)
	// testDeploymentHandler(t)
	// testProjectAllHandler(t)
	// testDiff(t)
	// testDiffSimple(t)

}

//Simple diff test with only one view
func TestDiffSame(t *testing.T) {
	loadTemplates()
	fmt.Println("TestDiffSame...")
	dom1 := &element.DOM{View: element.View{Children: make([]*element.View, 0), Template: "base.html", Provides: "html", Model: &element.Model{MAP: make(map[string]interface{})}}, Id: "1"}
	dom1.IdTree = &element.ViewTrie{}
	dom1.IdTree.Put(dom1.View.Provides, &dom1.View)
	dom2 := copyDom(*dom1)
	patches := Diff(&dom1.View, &dom2.View)
	if len(patches) != 0 {
		t.Error("Patches generated for same DOM")
		t.Fail()
	}
	// fmt.Printf("\ndom1: %+v\n", dom1.IdTree)
	// fmt.Printf("\ndom2: %+v\n", dom1.IdTree)
}

func TestDiffSimple(t *testing.T) {
	loadTemplates()
	fmt.Println("TestDiffSimple...")
	dom1 := &element.DOM{View: element.View{Children: make([]*element.View, 0), Template: "Base", Provides: "html", Model: &element.Model{MAP: make(map[string]interface{})}}, Id: "1"}
	dom1.IdTree = &element.ViewTrie{}
	dom1.IdTree.Put(dom1.View.Provides, &dom1.View)
	dom2 := copyDom(*dom1)
	dom2.Id = "2"
	dom2.View = element.View{Children: make([]*element.View, 0), Template: "Body", Provides: "html", Model: &element.Model{MAP: make(map[string]interface{})}}
	patches := Diff(&dom1.View, &dom2.View)
	if len(patches) != 1 {
		t.Errorf("Patches not generated %v\n", len(patches))
		t.Fail()
	} else {
		// fmt.Printf("\nPatches: %+v\n", patches)
	}
}

func TestDiffThreeLevel(t *testing.T) {
	loadTemplates()
	fmt.Println("TestDiffThreeLevel...")
	dom1 := &element.DOM{View: element.View{Children: make([]*element.View, 0), Template: "base.html", Provides: "html", Model: &element.Model{MAP: make(map[string]interface{})}}, Id: "1"}
	dom1.IdTree = &element.ViewTrie{}
	dom1.IdTree.Put(dom1.View.Provides, &dom1.View)
	head := &element.Head{View: element.View{Template: "head.html", Provides: "head", Return: "Head", InjectInto: "html", Model: &element.Model{MAP: make(map[string]interface{})}}}
	body := &element.Body{View: element.View{Template: "body.html", Provides: "body", Return: "Body", InjectInto: "html",
		Model: &element.Model{MAP: map[string]interface{}{
			"Content": "Body1",
		}}}}
	dom1 = Add(dom1, &head.View)
	dom1 = Add(dom1, &body.View)
	dom2 := copyDom(*dom1)
	dom2.Id = "2"
	dom2.View = element.View{Children: make([]*element.View, 0), Template: "base.html", Provides: "html", Model: &element.Model{MAP: make(map[string]interface{})}}
	body2 := &element.Body{View: element.View{Template: "body.html", Provides: "body", Return: "Body", InjectInto: "html",
		Model: &element.Model{MAP: map[string]interface{}{
			"Content": "Body2",
		}}}}
	dom2 = Add(dom2, &body2.View)
	patches := Diff(&dom1.View, &dom2.View)
	if len(patches) != 1 {
		t.Errorf("Patches not generated %v\n", len(patches))
		t.Fail()
	} else {
		fmt.Printf("\nPatches: %+v\n", patches)
		PrintDOM(dom1, dom1.Id)
		PrintDOM(dom2, dom2.Id)
	}
}

func loadTemplates() {
	//Load templates
	tmpl := template.Must(template.ParseGlob("views/*.html"))
	for _, t := range tmpl.Templates() {
		// log.Printf("Template imported: %s", t.Name())
		Templates[t.Name()] = t
	}
}

func TestDiff(t *testing.T) {
	loadTemplates()
	dom1 := &element.DOM{View: element.View{Children: make([]*element.View, 0), Template: "base.html", Provides: "html", Model: &element.Model{MAP: make(map[string]interface{})}}, Id: "1"}
	dom1.IdTree = &element.ViewTrie{}
	dom1.IdTree.Put(dom1.View.Provides, &dom1.View)
	head := &element.Head{View: element.View{Template: "head.html", Provides: "head", Return: "Head", InjectInto: "html", Model: &element.Model{MAP: make(map[string]interface{})}}}
	body := &element.Body{View: element.View{Template: "body.html", Provides: "body", Return: "Body", InjectInto: "html",
		Model: &element.Model{MAP: make(map[string]interface{})}}}
	dom1 = Add(dom1, &head.View)
	dom1 = Add(dom1, &body.View)
	header := &Header{View: &element.View{Template: "header.html", Provides: "#header", InjectInto: "body", Model: &element.Model{MAP: make(map[string]interface{})}}}
	dom1 = Add(dom1, header.View)

	result := []Deployment{Deployment{Name: "aqd-testv1", Environment: "dev"}}
	rowValues := make([]map[string]interface{}, len(result))

	for i, d := range result {
		row := make(map[string]interface{})
		row["ColumnValues"] = []string{d.Name, d.Environment}
		rowValues[i] = row
	}

	depl := &DeploymentView{View: &element.View{Template: "list.html", InjectInto: "#header", Model: &element.Model{MAP: map[string]interface{}{
		"Title":       "Deployments",
		"ColumnNames": []string{"Name", "Environment"},
		"RowValues":   rowValues,
	}}}}
	dom1 = Add(dom1, depl.View)

	dom2 := copyDom(*dom1)
	if dom1.Id == dom2.Id {
		t.Errorf("Id not changed: %s == %s\n", dom1.Id, dom2.Id)
		t.Fail()
	}
	result2 := []Project{Project{Name: "AQD", GitName: "AQD.git", Origin: "aqd@git.aduro.hr:AQD.git"}}
	rowValues2 := make([]map[string]interface{}, len(result2))

	for i, d := range result2 {
		row := make(map[string]interface{})
		row["ColumnValues"] = []string{d.Name, d.GitName, d.Origin}
		rowValues2[i] = row
	}

	proj := &ProjectView{View: &element.View{Template: "list.html", InjectInto: "#header", Model: &element.Model{MAP: map[string]interface{}{
		"Title":       "Projects",
		"ColumnNames": []string{"Name", "GitName", "Origin"},
		"RowValues":   rowValues2,
	}}}}
	// buf := Render(header.View)
	fmt.Printf("Header before %+v\n", len(header.View.Children))
	dom2 = Add(dom2, proj.View)
	// buf = Render(header.View)
	fmt.Printf("Header after %+v\n", len(header.View.Children))
	patches := Diff(&dom1.View, &dom2.View)
	// message, err := json.Marshal(patches)
	// if err != nil {
	// 	log.Println("Error marshaling patch")
	// }
	// fmt.Printf("result %+v\n", string(message))
	expected := 1
	if len(patches) != expected {
		t.Errorf("Patches not generated: expected: %d  actual: %d", expected, len(patches))
		// PrintDOM(dom1, dom1.Id)
		// PrintDOM(dom2, dom2.Id)
		t.Fail()
	}
}
