# Sylmark

Personal Knowledge Mangement(PKM) Language Server (LSP) with markdown files in golang

## Work Items

### v0.1 (current)

- [x] Minimal Treesitter parser
  - [tree_sitter_sylmark](https://codeberg.org/sylveryte/tree-sitter-sylmark)
- [x] Hover
  - [x] Tag
  - [x] Wikilinks
  - [x] Headings (references)
- [x] Completions
  - [x] Hash Tags
  - [x] Wikilinks File
  - [x] Wikilinks with sub headings
- [x] Go To Definitions
  - [x] Wikilinks File
  - [x] Wikilinks with sub headings
- [x] Go to references
  - [x] Tags
  - [x] Wikilinks
  - [x] Headings
- [x] Dim nonexisting wikilinks
- [x] Diagnostics
- [/] Sylgraph
  - [ ] Move into lsp
  * [/] Graph view of all nodes
    - [x] Files
    - [x] Links
    - [ ] Tags
  * [x] Click to open in editor
- [x] Code actions
  - [x] Created unresolved
    - [x] Update internal data
  - [x] Append heading
    - [x] Update internal data

### Bugs

- [ ] Spaced completions fetch results in crash

## Roadmap

### v0.2 (next)

- [ ] Sub tag support
- [ ] Rename heading across workspace
- [ ] Links Any File
- [ ] Wikilinks within file
  - [ ] Completions
  - [ ] References
- [ ] Sylgraph
  - [ ] Graph filters
  - [ ] Color based on groups

## Entities

- Hash Tags `#sylmark #lsp`
- Wikilinks File `[[Example]]`
- Wikilinks with sub headings `[[Example#Objective]]`
- Wikilinks within file `[[#Work Items]]`
- Links Any File `[Go mod file](./go.mod)`

## Structure

- Package flow

  `main <- handle <- server <- data <- lsp`
