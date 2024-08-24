package charthtml

type Symbol struct {
	values            []string
	currentValueIndex int
	assignedSymbols   map[string]string
}

func NewSymbol() *Symbol {
	return &Symbol{
		values:            []string{"diamond", "circle", "triangle", "rect", "roundRect", "pin", "arrow"},
		assignedSymbols:   make(map[string]string),
		currentValueIndex: -1,
	}
}

func (s *Symbol) Next() string {
	s.currentValueIndex++
	if s.currentValueIndex >= len(s.values) {
		s.currentValueIndex = 0
	}

	return s.values[s.currentValueIndex]
}

func (s *Symbol) GetFor(language string) string {
	symbol, ok := s.assignedSymbols[language]
	if !ok {
		symbol = s.Next()
		s.assignedSymbols[language] = symbol
	}

	return symbol
}
