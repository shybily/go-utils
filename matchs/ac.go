package matchs

import (
	"container/list"
)

type acNode struct {
	count int
	fail  *acNode
	child map[rune]*acNode
	index int
}

func newAcNode() *acNode {
	return &acNode{
		count: 0,
		fail:  nil,
		child: make(map[rune]*acNode),
		index: -1,
	}
}

type AcMatcher struct {
	root *acNode
	size int
	mark []bool
}

func NewAcMatcher() *AcMatcher {
	return &AcMatcher{
		root: newAcNode(),
		size: 0,
		mark: make([]bool, 0),
	}
}

// initialize the ahocorasick
func (m *AcMatcher) Build(dictionary []string) {
	for i, _ := range dictionary {
		m.insert(dictionary[i])
	}
	m.build()
	m.mark = make([]bool, m.size)
}

// string matchs search
// return all strings matched as indexes into the original dictionary
func (m *AcMatcher) Match(s string) []int {

	m.resetMark()

	var (
		curNode = m.root
		ret     []int
		n       *acNode
	)

	for _, v := range s {
		for curNode.child[v] == nil && curNode != m.root {
			curNode = curNode.fail
		}
		curNode = curNode.child[v]
		if curNode == nil {
			curNode = m.root
		}

		n = curNode
		for n != m.root && n.count > 0 && !m.mark[n.index] {
			m.mark[n.index] = true
			for i := 0; i < n.count; i++ {
				ret = append(ret, n.index)
			}
			n = n.fail
		}
	}

	return ret
}

// just return the number of len(Match(s))
func (m *AcMatcher) GetMatchResultSize(s string) int {

	m.resetMark()

	var (
		curNode = m.root
		n       *acNode
		num     = 0
	)

	for _, v := range s {
		for curNode.child[v] == nil && curNode != m.root {
			curNode = curNode.fail
		}
		curNode = curNode.child[v]
		if curNode == nil {
			curNode = m.root
		}

		n = curNode
		for n != m.root && n.count > 0 && !m.mark[n.index] {
			m.mark[n.index] = true
			num += n.count
			n = n.fail
		}
	}

	return num
}

func (m *AcMatcher) build() {
	ll := list.New()
	ll.PushBack(m.root)
	for ll.Len() > 0 {
		temp := ll.Remove(ll.Front()).(*acNode)
		var p *acNode = nil

		for i, v := range temp.child {
			if temp == m.root {
				v.fail = m.root
			} else {
				p = temp.fail
				for p != nil {
					if p.child[i] != nil {
						v.fail = p.child[i]
						break
					}
					p = p.fail
				}
				if p == nil {
					v.fail = m.root
				}
			}
			ll.PushBack(v)
		}
	}
}

func (m *AcMatcher) insert(s string) {
	curNode := m.root
	for _, v := range s {
		if curNode.child[v] == nil {
			curNode.child[v] = newAcNode()
		}
		curNode = curNode.child[v]
	}
	curNode.count++
	curNode.index = m.size
	m.size++
}

func (m *AcMatcher) resetMark() {
	for i := 0; i < m.size; i++ {
		m.mark[i] = false
	}
}
