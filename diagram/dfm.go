package diagram

import (
	"github.com/awalterschulze/gographviz"

	"strings"
	"log"
	"os"
)

/*
Struct Dfm defines the star schema and includes a pointer to 
a graph.
*/
type Dfm struct {
	Graph *gographviz.Graph	 //graph
}

/*
Struct Fact defines a fact and includes a name 
and atrributes.
*/
type Fact struct {
	name string
	attributes []string
}

/*
Attributes are maps of attributes for modifying nodes and edges 
in the dfm schema.
*/
var  (
	DFM_nodeAtt = map[string]string{"shape":"circle", "label":"\"\""}
	DFM_edgeAtt = map[string]string{"arrowhead":"none", "len":"0.5"}
	DFM_factAtt =  map[string]string{"shape":"plain", "root":"true"}
	DFM_descriptiveAtt = map[string]string{"shape":"underline"}
	DFM_optionalAtt = map[string]string{"arrowhead":"icurve"}
	DFM_hierarchyAtt = map[string]string{"arrowhead":"none"}
)

/*
Function NewDFM creates a new Dfm object for creating the dfm schema.
*/
func NewDFM() *Dfm {

	graphAst, _ := gographviz.ParseString(`digraph G { 
		layout=twopi;
		overlap=prism;
		overlap_scaling=4.5;
	}`)
	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		panic(err)
	}

	return &Dfm{
		graph,
	}
}

/*
Function CreateFact creates a new Fact object and renders
the fact in the schema
*/
func (schema *Dfm) CreateFact(title string, attributes string) *Fact {

	t_attributes := strings.Split(attributes, " ")

	fact := &Fact {
		title,
		t_attributes,
	}


	schema.RenderFact(fact)
	return fact
}

/*
Function RenderFact renders a given fact in the schema.
*/
func (schema *Dfm) RenderFact(f *Fact) {

	label := `<<table border="0" cellborder="1" cellspacing="0" cellpadding="20"> <tr> <td bgcolor="lightblue">`+f.name+`</td> </tr>`
	for _, att := range(f.attributes) {
		label += `<tr> <td>`+att+`</td> </tr>`
	}
	label += `</table>>`

	fact_att := DFM_factAtt
	fact_att["label"] = label

	schema.Graph.AddNode("G", f.name, fact_att)
}

/*
Function AddDimension adds a dimension in the schema
*/
func (schema *Dfm) AddDimension(label string, attach string) {
	
	//refactor this
	node_att := DFM_nodeAtt
	node_att["xlabel"] = label
	node_att["fixedsize"] = "true"

	schema.Graph.AddNode("G", label, node_att)
	schema.Graph.AddEdge(attach, label, true, DFM_edgeAtt)
}


/*
Function AddSequenceDimension adds a sequence of dimensions
in the schema
*/
func (schema *Dfm) AddSequenceDimension(labels string, startAttach string) {

	arr_labels := strings.Split(labels, " ")

	schema.AddDimension(arr_labels[0], startAttach)
	for i := 0; i < len(arr_labels)-1; i++ {
		schema.AddDimension(arr_labels[i+1], arr_labels[i])
	}

}

/*
Function AddConvergence adds a convergence node in the schema
*/
func (schema *Dfm) AddConvergence(label string, attach string) {
	node_att := DFM_nodeAtt
	node_att["xlabel"] = label

	schema.Graph.AddNode("G", label, node_att)
	schema.Graph.AddEdge(attach, label, true, nil)
}

/*
Function AddHierarchy adds a hyerarchy node in the schema
*/
func (schema *Dfm) AddHierarchy(labels string, from string, to string) {

	node_att := DFM_nodeAtt
	node_att["label"] = to

	schema.Graph.AddNode("G", to, node_att)

	for _, label := range strings.Split(labels, " ") {
		tmpAtt := DFM_hierarchyAtt
		tmpAtt["xlabel"] = label
		schema.Graph.AddEdge(from, to, true, tmpAtt)
	}

}

/*
Function AddOptional adds an optional node in the schema
*/
func (schema *Dfm) AddOptional(label string, attach string) {

	node_att := DFM_nodeAtt
	node_att["xlabel"] = label

	schema.Graph.AddNode("G", label, node_att)
	schema.Graph.AddEdge(attach, label, true, DFM_optionalAtt)
}

/*
Function AddDescriptive adds a descriptive node in the schema
*/
func (schema *Dfm) AddDescriptive(label string, to string) {

	schema.Graph.AddNode("G", label, DFM_descriptiveAtt)
	schema.Graph.AddEdge(to, label, true, DFM_edgeAtt)
}

/*
Function AddSequenceDescriptive adds a sequence of descreptive nodes
in the schema
*/
func (schema *Dfm) AddSequenceDescriptive(labels string, to string) {

	for _, label := range strings.Split(labels, " ") {
		schema.AddDescriptive(label, to)
	}
}

/*
Function RenderDiagram renders the entire schema to a dot file
*/
func (schema *Dfm) RenderDiagram() {
	output := schema.Graph.String()

	f, err := os.Create("dot.dot")
	if err != nil {
        log.Fatal(err)
    }

	defer f.Close()
	f.WriteString(output)
}