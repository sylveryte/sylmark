package server

type GraphStore struct {
	nodeStore map[string]Node
	idStore   map[int]string
}

func newGraphStore() *GraphStore {
	return &GraphStore{
		nodeStore: map[string]Node{},
		idStore:   map[int]string{},
	}
}

func (idStore *GraphStore) GetNodeFromName(name string) (node Node, found bool) {
	s := *idStore
	node, found = s.nodeStore[name]
	return node, found
}

func (idStore *GraphStore) GetNodeFromId(id int) (node Node, found bool) {
	s := *idStore
	name, found := s.idStore[id]
	if found {
		node, found = s.nodeStore[name]
	}
	return node, found
}

func (idStore *GraphStore) StoreAndGetId(name string, val int, kind NodeKind) Node {

	s := *idStore
	node, found := s.nodeStore[name]
	if !found {
		// add new node
		id := len(s.nodeStore)
		s.nodeStore[name] = Node{
			Id:   id,
			Name: name,
			Val:  val,
			Kind: kind,
		}
		s.idStore[id] = name
	}
	return node
}
