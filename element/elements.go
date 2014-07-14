package element

type DOM struct {
	View
	IdTree *ViewTrie
	Id     string
}
type Head struct {
	View
}
type Body struct {
	View
}

type Model struct {
	MAP  map[string]interface{}
	TYPE interface{}
}
type View struct {
	Children   []*View
	Template   string
	Content    string
	InjectInto string //ID to inject this view into
	Provides   string //ID that this view provides, into which we can inject other views
	*Model
	Return string //key which view returns and will be used in model for templating
	Parent *View
}
