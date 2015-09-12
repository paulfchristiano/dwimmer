package synonyms

type UF struct {
	parent map[int]int
	rank   map[int]int
}

func NewUF() *UF {
	return &UF{
		parent: make(map[int]int),
		rank:   make(map[int]int),
	}
}

func (uf *UF) Find(n int) (m int) {
	m, ok := uf.parent[n]
	if !ok {
		uf.parent[n] = n
		return n
	}
	if m == n {
		return m
	}
	m = uf.Find(m)
	uf.parent[n] = m
	return
}

func (uf *UF) Union(n, m int) {
	n = uf.Find(n)
	m = uf.Find(m)
	a := uf.rank[n]
	b := uf.rank[m]
	if a < b {
		uf.parent[n] = m
	} else {
		uf.parent[m] = n
		if a == b {
			uf.rank[n]++
		}
	}
}
