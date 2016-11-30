package bioportal

type Ontology struct {
	AdministeredBy []string `json:"administeredBy"`
	Acronym        string   `json:"acronym"`
	Name           string   `json:"name"`
	SummaryOnly    bool     `json:"summaryOnly"`
	OntologyType   string   `json:"ontologyType"`
	ID             string   `json:"@id"`
	Type           string   `json:"@type"`
	Links          struct {
		Submissions      string `json:"submissions"`
		Properties       string `json:"properties"`
		Classes          string `json:"classes"`
		SingleClass      string `json:"single_class"`
		Roots            string `json:"roots"`
		Instances        string `json:"instances"`
		Metrics          string `json:"metrics"`
		Reviews          string `json:"reviews"`
		Notes            string `json:"notes"`
		Groups           string `json:"groups"`
		Categories       string `json:"categories"`
		LatestSubmission string `json:"latest_submission"`
		Projects         string `json:"projects"`
		Download         string `json:"download"`
		Views            string `json:"views"`
		Analytics        string `json:"analytics"`
		UI               string `json:"ui"`
		Context          struct {
			Submissions      string `json:"submissions"`
			Properties       string `json:"properties"`
			Classes          string `json:"classes"`
			SingleClass      string `json:"single_class"`
			Roots            string `json:"roots"`
			Instances        string `json:"instances"`
			Metrics          string `json:"metrics"`
			Reviews          string `json:"reviews"`
			Notes            string `json:"notes"`
			Groups           string `json:"groups"`
			Categories       string `json:"categories"`
			LatestSubmission string `json:"latest_submission"`
			Projects         string `json:"projects"`
			Download         string `json:"download"`
			Views            string `json:"views"`
			Analytics        string `json:"analytics"`
			UI               string `json:"ui"`
		} `json:"@context"`
	} `json:"links"`
	Context struct {
		Vocab          string `json:"@vocab"`
		Acronym        string `json:"acronym"`
		Name           string `json:"name"`
		AdministeredBy struct {
			ID   string `json:"@id"`
			Type string `json:"@type"`
		} `json:"administeredBy"`
		OntologyType struct {
			ID   string `json:"@id"`
			Type string `json:"@type"`
		} `json:"ontologyType"`
	} `json:"@context"`
}

type Class struct {
	PrefLabel    string        `json:"prefLabel"`
	Synonym      []string      `json:"synonym"`
	Definition   []string      `json:"definition"`
	Cui          []interface{} `json:"cui"`
	SemanticType []interface{} `json:"semanticType"`
	Obsolete     interface{}   `json:"obsolete"`
	ID           string        `json:"@id"`
	Type         string        `json:"@type"`
	Links        struct {
		Self        string `json:"self"`
		Ontology    string `json:"ontology"`
		Children    string `json:"children"`
		Parents     string `json:"parents"`
		Descendants string `json:"descendants"`
		Ancestors   string `json:"ancestors"`
		Instances   string `json:"instances"`
		Tree        string `json:"tree"`
		Notes       string `json:"notes"`
		Mappings    string `json:"mappings"`
		UI          string `json:"ui"`
		Context     struct {
			Self        string `json:"self"`
			Ontology    string `json:"ontology"`
			Children    string `json:"children"`
			Parents     string `json:"parents"`
			Descendants string `json:"descendants"`
			Ancestors   string `json:"ancestors"`
			Instances   string `json:"instances"`
			Tree        string `json:"tree"`
			Notes       string `json:"notes"`
			Mappings    string `json:"mappings"`
			UI          string `json:"ui"`
		} `json:"@context"`
	} `json:"links"`
	Context struct {
		Vocab        string `json:"@vocab"`
		PrefLabel    string `json:"prefLabel"`
		Synonym      string `json:"synonym"`
		Definition   string `json:"definition"`
		Obsolete     string `json:"obsolete"`
		SemanticType string `json:"semanticType"`
		Cui          string `json:"cui"`
	} `json:"@context"`
}

type ClassesPaginated struct {
	Page      int         `json:"page"`
	PageCount int         `json:"pageCount"`
	PrevPage  interface{} `json:"prevPage"`
	NextPage  interface{} `json:"nextPage"`
	Links     struct {
		NextPage interface{} `json:"nextPage"`
		PrevPage interface{} `json:"prevPage"`
	} `json:"links"`
	Collection Class `json:"collection"`
}

type Mapping struct {
	MappingID interface{} `json:"id"`
	Source    string      `json:"source"`
	Classes   []Class     `json:"classes"`
	Process   interface{} `json:"process"`
	ID        string      `json:"@id"`
	Type      string      `json:"@type"`
}

type Tree struct {
	PrefLabel   string `json:"prefLabel"`
	HasChildren bool   `json:"hasChildren"`
	Children    []Tree `json:"children"`
	Obsolete    bool   `json:"obsolete"`
	ID          string `json:"@id"`
	Type        string `json:"@type"`
	Links       struct {
		Self        string `json:"self"`
		Ontology    string `json:"ontology"`
		Children    string `json:"children"`
		Parents     string `json:"parents"`
		Descendants string `json:"descendants"`
		Ancestors   string `json:"ancestors"`
		Instances   string `json:"instances"`
		Tree        string `json:"tree"`
		Notes       string `json:"notes"`
		Mappings    string `json:"mappings"`
		UI          string `json:"ui"`
		Context     struct {
			Self        string `json:"self"`
			Ontology    string `json:"ontology"`
			Children    string `json:"children"`
			Parents     string `json:"parents"`
			Descendants string `json:"descendants"`
			Ancestors   string `json:"ancestors"`
			Instances   string `json:"instances"`
			Tree        string `json:"tree"`
			Notes       string `json:"notes"`
			Mappings    string `json:"mappings"`
			UI          string `json:"ui"`
		} `json:"@context"`
	} `json:"links"`
	Context struct {
		Vocab     string `json:"@vocab"`
		PrefLabel string `json:"prefLabel"`
		Obsolete  string `json:"obsolete"`
	} `json:"@context"`
}
