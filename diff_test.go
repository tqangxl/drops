package drops

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/AduroIdea/dplrlib"
	"github.com/cyfdecyf/tst"
	"log"
	"testing"
)

type Header struct {
	*View
}

type Deployment struct {
	*View
}
type Project struct {
	*View
}

func TestDiffSuite(t *testing.T) {
	// testEchoHandler(t)
	// testDeploymentHandler(t)
	// testProjectAllHandler(t)
	testDiff(t)

}

func testDiff(t *testing.T) {
	l1 := list.New()
	dom1 := &DOM{View: &View{Children: l1, Template: "../../views/base.html", Provides: "body"}}
	dom1.IdTree = &tst.Trie{}
	dom1.IdTree.Put(dom1.View.Provides, dom1.View)
	head := &Head{View: &View{Template: "../../views/head.html", Provides: "head", Return: "Head"}}
	body := &Body{View: &View{Template: "../../views/body.html", Provides: "body", Return: "Body", InjectInto: "body",
		Model: map[string]interface{}{
			"Content": "Body1",
		}}}
	dom1 = Add(dom1, head.View)
	dom1 = Add(dom1, body.View)
	header := &Header{View: &View{Template: "../../views/header.html", Provides: "#header", InjectInto: "body"}}
	dom1 = Add(dom1, header.View)

	result := []dplr.Deployment{dplr.Deployment{Name: "aqd-testv1", Environment: "dev"}}
	rowValues := make([]map[string]interface{}, len(result))

	for i, d := range result {
		row := make(map[string]interface{})
		row["ColumnValues"] = []string{d.Name, d.Environment}
		rowValues[i] = row
	}

	depl := &Deployment{View: &View{Template: "../../views/index.html", InjectInto: "#header", Model: map[string]interface{}{
		"ColumnNames": []string{"Name", "Environment"},
		"RowValues":   rowValues,
	}}}
	dom1 = Add(dom1, depl.View)

	l2 := list.New()
	dom2 := &DOM{View: &View{Children: l2, Template: "../../views/base.html", Provides: "body"}}
	dom2.IdTree = &tst.Trie{}
	dom2.IdTree.Put(dom2.View.Provides, dom2.View)

	// head = &Head{View: &View{Template: "../../views/head.html", Provides: "head", Return: "Head"}}
	body = &Body{View: &View{Template: "../../views/body.html", Provides: "body", Return: "Body", InjectInto: "body",
		Model: map[string]interface{}{
			"Content": "Body1",
		}}}
	dom2 = Add(dom2, head.View)
	dom2 = Add(dom2, body.View)
	header = &Header{View: &View{Template: "../../views/header.html", Provides: "#header", InjectInto: "body"}}
	dom2 = Add(dom2, header.View)

	result2 := []dplr.Project{dplr.Project{Name: "AQD", GitName: "AQD.git", Origin: "aqd@git.aduro.hr:AQD.git"}}
	rowValues2 := make([]map[string]interface{}, len(result2))

	for i, d := range result2 {
		row := make(map[string]interface{})
		row["ColumnValues"] = []string{d.Name, d.GitName, d.Origin}
		rowValues2[i] = row
	}

	proj := &Project{View: &View{Template: "../../views/index.html", InjectInto: "#header", Model: map[string]interface{}{
		"ColumnNames": []string{"Name", "GitName", "Origin"},
		"RowValues":   rowValues2,
	}}}
	dom2 = Add(dom2, proj.View)

	patches := Diff(dom1.View, dom2.View)
	message, err := json.Marshal(patches)
	if err != nil {
		log.Println("Error marshaling patch")
	}
	fmt.Printf("result %+v\n", string(message))

}
