# Doc

## Structure

- Package flow

  `main <- handle <- server <- data <- lsp`

### Store

```d2
vars: {
  d2-config: {
    layout-engine: elk
    pad: 6
  }
}
s: Store
ls: "map[id]Link" {
  style.double-border: true
}
ts: "map[Target][]Id" {
  style.double-border: true
}
t: "map[Tag][]IdLocation" {
  style.double-border: true
}
lk: link {
  style.multiple: true
  Def
  Refs
}
lnr: "map[SubTarget][]IdLocation" {
  style.double-border: true
}
lnd: "map[SubTarget]Range" {
  style.double-border: true
}
ds: "map[Id]DocumentData" {
  style.double-border: true
}
# gl."Defs []Location"->d
ds -> dd
dd: "DocumentData" {
  FootNotes
  HeadingsStore
  style.multiple: true
  "Trees": {
    style.multiple: true
  }
  Content
  FootNotes
}
fns: "map[string]FootNoteRef" {
  style.double-border: true
}
fnr: "FootNoteRef" {
  style.multiple: true
  "Def *Range"
  "Refs []Range": {
    style.multiple: true
  }
  Excert
}
fns -> fnr
hds: "map[string]Subheading" {
  style.double-border: true
}
shd: "Subheading" {
  style.multiple: true
  "Def Range"
  "Refs []Range": {
    style.multiple: true
  }
}
hds -> shd
sm: "[]string" {
  style.multiple: true
}

dd.FootNotes -> fns
dd.HeadingsStore -> hds
lk.Def -> lnd
lk.Refs -> lnr
s: {
  Tags
  LinkStore
  TargetStore
  DocStore
  IdStore
  OtherFiles
}
s.TargetStore -> ts
s.Tags -> t
s.LinkStore -> ls
ls -> lk
s.DocStore -> ds
s.OtherFiles -> sm
```

## Graph Store

```d2


vars: {
  d2-config: {
    layout-engine: elk
    pad: 6
  }
}
n:Node{
  NodeId
  Name
  Val
  Kind
  uri
  style.multiple: true
}

ns: "map[NodeId]Node" {
  style.double-border: true
}
ls: "map[NodeId]map[NodeId]int" {
  style.double-border: true
}
gs:GraphStore {
  NodeStore
  LinkStore
}
gs.NodeStore -> ns
gs.LinkStore -> ls
ns->n
g:Graph{
  Nodes
  Links
  }
n->g.Nodes
ls->g.Links
s:Store{
  style:{
    stroke-dash:4
    }
}
s->gs
```
