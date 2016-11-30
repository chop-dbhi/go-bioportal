package bioportal

type SearchOptions struct {
	BaseOptions

	Query                string   `url:"q"`
	Ontologies           []string `url:"ontologies,omitempty"`
	RequireExactMatch    bool     `url:"require_exact_match"`
	Suggest              bool     `url:"suggest"`
	RequireDefinitions   bool     `url:"require_definitions"`
	AlsoSearchProperties bool     `url:"also_search_properties"`
	AlsoSearchViews      bool     `url:"also_search_views"`
	AlsoSearchObsolete   bool     `url:"also_search_obsolete"`
	CUI                  []string `url:"cui,omitempty"`
	SemanticTypes        []string `url:"semantic_types,omitempty"`
}

func DefaultSearchOptions() *SearchOptions {
	return &SearchOptions{
		BaseOptions:          *DefaultBaseOptions(),
		RequireExactMatch:    false,
		Suggest:              false,
		RequireDefinitions:   false,
		AlsoSearchProperties: false,
		AlsoSearchViews:      false,
		AlsoSearchObsolete:   false,
	}
}

type SearchResult struct {
	Page      int         `json:"page"`
	PageCount int         `json:"pageCount"`
	PrevPage  interface{} `json:"prevPage"`
	NextPage  int         `json:"nextPage"`
	Links     struct {
		NextPage string      `json:"nextPage"`
		PrevPage interface{} `json:"prevPage"`
	} `json:"links"`
	Collection []struct {
		PrefLabel    string   `json:"prefLabel"`
		Synonym      []string `json:"synonym,omitempty"`
		Cui          []string `json:"cui,omitempty"`
		SemanticType []string `json:"semanticType,omitempty"`
		Obsolete     bool     `json:"obsolete"`
		MatchType    string   `json:"matchType"`
		OntologyType string   `json:"ontologyType"`
		Provisional  bool     `json:"provisional"`
		ID           string   `json:"@id"`
		Type         string   `json:"@type"`
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
			Obsolete     string `json:"obsolete"`
			SemanticType string `json:"semanticType"`
			Cui          string `json:"cui"`
		} `json:"@context"`
		Definition []string `json:"definition,omitempty"`
	} `json:"collection"`
}
