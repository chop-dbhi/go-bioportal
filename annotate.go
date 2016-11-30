package bioportal

type AnnotateOptions struct {
	BaseOptions

	Text                         string   `url:"text"`
	Ontologies                   []string `url:"ontologies,omitempty"`
	SemanticTypes                []string `url:"semantic_types,omitempty"`
	ExpandSemanticTypesHierarchy bool     `url:"expand_semantic_types_hierarchy"`
	ExpandClassHierarchy         bool     `url:"expand_class_hierarchy"`
	ClassHierarchyMaxLevel       uint     `url:"class_hierarchy_max"`
	ExpandMappings               bool     `url:"expand_mappings"`
	StopWords                    []string `url:"stop_words,omitempty"`
	MinimumMatchLength           uint     `url:"minimum_match_length"`
	ExcludeNumbers               bool     `url:"exclude_numbers"`
	WholeWordOnly                bool     `url:"whole_word_only"`
	ExcludeSynonyms              bool     `url:"exclude_synonyms"`
	LongestOnly                  bool     `url:"longest_only"`
}

func DefaultAnnotateOptions() *AnnotateOptions {
	return &AnnotateOptions{
		BaseOptions:   *DefaultBaseOptions(),
		WholeWordOnly: true,
	}
}

type AnnotationResult struct {
	AnnotatedClass struct {
		ID    string `json:"@id"`
		Type  string `json:"@type"`
		Links struct {
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
			Vocab string `json:"@vocab"`
		} `json:"@context"`
	} `json:"annotatedClass"`
	Hierarchy   []interface{} `json:"hierarchy"`
	Annotations []struct {
		From      int    `json:"from"`
		To        int    `json:"to"`
		MatchType string `json:"matchType"`
		Text      string `json:"text"`
	} `json:"annotations"`
	Mappings []interface{} `json:"mappings"`
}
