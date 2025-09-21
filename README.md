# Sylmark

Personal Knowledge Mangement(PKM) Language Server (LSP) with markdown files in golang

## Installation

- on linux
  `CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build`

## Setup

```lua
    vim.lsp.config.sylmark = {
      cmd = { "path/to/binary" },
      root_markers = { '.sylroot' },
      filetypes = { 'markdown' },
      on_attach = function(client, bufnr)
        vim.api.nvim_create_user_command(
          "Daily",
          function(args)
            local input = args.args

            client:exec_cmd({
              title = "Show",
              command = "show",
              arguments = { input }, -- Also works with `vim.NIL`
            }, { bufnr = bufnr })
          end,
          { desc = 'Open daily note', nargs = "*" }
        )
        vim.api.nvim_create_user_command(
          "Graph",
          function(args)
            local input = args.args

            client:exec_cmd({
              title = "Open Graph",
              command = "graph",
              arguments = { input }, -- Also works with `vim.NIL`
            }, { bufnr = bufnr })
          end,
          { desc = 'Start graph server and open', nargs = "*" }
        )
      end
    }

    vim.lsp.enable({
      "sylmark",
    })
```

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
- [x] Sylgraph
  - [x] Move into lsp
  * [x] Graph view of all nodes
    - [x] Files
    - [x] Links
    - [x] Tags
  * [x] Click to open in editor
- [x] Code actions
  - [x] Created unresolved
    - [x] Update internal data
  - [x] Append heading
    - [x] Update internal data

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
