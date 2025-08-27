package templates

type Template interface {
	GridDefinitionEnd(template int, section3Bytes []byte) (int, bool)    // Table 3.1
	ProductDefinitionEnd(template int, section4Bytes []byte) (int, bool) // Table 4.0
}

type RetrieveEndOctet func([]byte) int

func constantEnd(end int) RetrieveEndOctet {
	return func(bytes []byte) int {
		return end
	}
}

type template struct {
	gridDefinition    map[int]RetrieveEndOctet
	productDefinition map[int]RetrieveEndOctet
}

func (t *template) GridDefinitionEnd(template int, section3Bytes []byte) (int, bool) {
	retriever, ok := t.gridDefinition[template]
	if !ok {
		return 0, false
	}
	return retriever(section3Bytes), true
}

func (t *template) ProductDefinitionEnd(template int, section4Bytes []byte) (int, bool) {
	retriever, ok := t.productDefinition[template]
	if !ok {
		return 0, false
	}
	return retriever(section4Bytes), true
}
