package reX

type T map[string]interface{}
type StringMap map[string]string
type O interface {
	interface{}
}

func Map[SRC O, DEST O, RESULT []DEST](os []SRC, convertFn func(SRC) DEST) RESULT {
	lst := make([]DEST, 0, len(os))
	for _, o := range os {
		lst = append(lst, convertFn(o))
	}
	return lst
}
