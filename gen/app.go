package gen

type AppTemplate map[interface{}]interface{}

func (a AppTemplate) Merge(target AppTemplate) AppTemplate {
	for ti, tv := range target {
		if _, ok := a[ti]; ok == false {
			a[ti] = tv
		} else {
			switch tv.(type) {
			case AppTemplate:
				a[ti] = a[ti].(AppTemplate).Merge(tv.(AppTemplate))
			default:
				a[ti] = tv
			}
		}
	}
	return a
}
