package record

type Section3 struct {
	Length                       int
	GridSourceDefinition         int
	TotalPoints                  int
	OctetsForOptionalPointList   int
	InterpretationOfPointList    int
	GridDefinitionTemplateNumber int
	GridDefinitionTemplateData   []byte
	OptionalPointListData        []byte
}
