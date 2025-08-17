# Sylmark

Personal Knowledge Mangement(PKM) Language Server (LSP) with markdown files in golang

## Entities

- Hash Tags `#sylmark #lsp`
- Wikilinks File `[[Example]]`
- Wikilinks with sub headings `[[Example#Objective]]`
- Wikilinks within file `[[#Work Items]]`
- Links Any File `[Go mod file](./go.mod)`

## Work Items

- [x] Minimal Treesitter parser
  - [tree_sitter_sylmark](https://github.com/sylveryte/tree-sitter-sylmark)
- [x] Hover
  - [x] Tag
  - [ ] Wikilinks
- [/] Completions
  - [x] Hash Tags
  - [/] Wikilinks File
  - [/] Wikilinks with sub headings
  - [ ] Wikilinks within file
  - [ ] ?Links Any File
- [ ] Go To Definitions
  - [ ] Wikilinks File
  - [ ] Wikilinks with sub headings
  - [ ] Wikilinks within file
  - [ ] ?Links Any File
- [/] Go to references
  - [x] Tags
  - [ ] Wikilinks

## Structure

- Package flow

  `main <- handle <- data <- lsp`
