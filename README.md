# Sylmark

Personal Knowledge Mangement(PKM) Language Server (LSP) with markdown files in golang

## Entities

- Hash Tags `#sylmark #lsp`
- Wikilinks File `[[Example]]`
- Wikilinks with sub headings `[[Example#Objective]]`
- Wikilinks within file `[[#Work Items]]`
- Links Any File `[Go mod file](./go.mod)`

## Work Items

### v0.1 (current)

- [x] Minimal Treesitter parser
  - [tree_sitter_sylmark](https://codeberg.org/sylveryte/tree-sitter-sylmark)
- [x] Hover
  - [x] Tag
  - [x] Wikilinks
  - [x] Headings (references)
- [/] Completions
  - [x] Hash Tags
  - [x] Wikilinks File
  - [x] Wikilinks with sub headings
  - [ ] Wikilinks within file
- [/] Go To Definitions
  - [x] Wikilinks File
  - [x] Wikilinks with sub headings
  - [ ] Wikilinks within file
- [x] Go to references
  - [x] Tags
  - [x] Wikilinks
  - [x] Headings

## Roadmap

### v0.2 (next)

- [ ] Sub tag support
- [ ] Dim nonexisting wikilinks
- [ ] Rename heading across workspace
- [ ] Code actions
- [ ] Diagnostics
- [ ] Links Any File

### v0.3

- [ ] Sylgraph - Graph view of all nodes

## Structure

- Package flow

  `main <- handle <- data <- lsp`
