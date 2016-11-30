package bioportal

type RecommendOptions struct {
	BaseOptions

	Terms                []string `url:"input"`
	Ontologies           []string `url:"ontologies,omitempty"`
	InputType            int      `url:"input_type,omitempty"`
	OutputType           int      `url:"output_type,omitempty"`
	MaxElementsSet       int      `url:"max_elements_set,omitempty"`
	CoverageWeight       float32  `url:"wc,omitempty"`
	SpecializationWeight float32  `url:"ws,omitempty"`
	AcceptanceWeight     float32  `url:"wa,omitempty"`
	DetailWeight         float32  `url:"wd,omitempty"`
}

func DefaultRecommendOptions() *RecommendOptions {
	return &RecommendOptions{
		BaseOptions:    *DefaultBaseOptions(),
		InputType:      1,
		OutputType:     1,
		MaxElementsSet: 3,
	}
}

type RecommendResult struct {
	EvaluationScore float64 `json:"evaluationScore"`
	Ontologies      []struct {
		Acronym string `json:"acronym"`
		ID      string `json:"@id"`
		Type    string `json:"@type"`
	} `json:"ontologies"`
	CoverageResult struct {
		Score              int `json:"score"`
		NormalizedScore    int `json:"normalizedScore"`
		NumberTermsCovered int `json:"numberTermsCovered"`
		NumberWordsCovered int `json:"numberWordsCovered"`
		Annotations        []struct {
			From           int    `json:"from"`
			To             int    `json:"to"`
			MatchType      string `json:"matchType"`
			Text           string `json:"text"`
			AnnotatedClass struct {
				ID   string `json:"@id"`
				Type string `json:"@type"`
			} `json:"annotatedClass"`
			HierarchySize int `json:"hierarchySize"`
		} `json:"annotations"`
	} `json:"coverageResult"`
	SpecializationResult struct {
		Score           float64 `json:"score"`
		NormalizedScore int     `json:"normalizedScore"`
	} `json:"specializationResult"`
	AcceptanceResult struct {
		NormalizedScore float64 `json:"normalizedScore"`
		BioportalScore  float64 `json:"bioportalScore"`
		UmlsScore       int     `json:"umlsScore"`
	} `json:"acceptanceResult"`
	DetailResult struct {
		NormalizedScore  float64 `json:"normalizedScore"`
		DefinitionsScore int     `json:"definitionsScore"`
		SynonymsScore    int     `json:"synonymsScore"`
		PropertiesScore  float64 `json:"propertiesScore"`
	} `json:"detailResult"`
}
