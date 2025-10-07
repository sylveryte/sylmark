Note: Github is mirror, original repo on [codeberg](https://codeberg.org/sylveryte/sylmark)

# Sylmark

Personal Knowledge Mangement(PKM) Language Server (LSP) with markdown files in golang

## Build

- Using go
  `CGO_ENABLED=1 go build`

## Neovim Setup

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

## Roadmap

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
  - [x] Common dates links
- [x] Go To Definitions
  - [x] Wikilinks File
  - [x] Wikilinks with sub headings
- [x] Go to references
  - [x] Tags
  - [x] Wikilinks
  - [x] Headings
- [x] Dim unresolved wikilinks
- [x] Diagnostics
- [x] Code actions
  - [x] Created unresolved
    - [x] Update internal data
  - [x] Append heading
    - [x] Update internal data
- [x] Wikilinks within file
  - [x] Completions
  - [x] References
- [x] Symbols
  - [x] Dynamic workspace symbols
- [x] Switch to official Treesitter markdown parsers
- [/] Inline link (Markdown style links )
  - [x] Any file
  - [x] Images file link include !
  - [ ] Markdown files
    - [ ] Make proper links in store
    - [ ] References
    - [ ] Hover
    - [ ] Definitions
    - [ ] Interoperability with wikilinks
    - [ ] Headings
- [x] Footnotes
- [ ] Rename heading across workspace
- [ ] Rename file changes across workspace
- [ ] Better nested tag support
- [/] Sylgraph
  - [x] Graph view of all nodes
    - [x] Files
    - [x] Links
    - [x] Tags
  - [x] Click to open in editor
  - [ ] Local mode
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

### Store

```d2
vars:{
  d2-config{
  layout-engine:elk
  pad:5
  }
}
s: Store
gls: "map[GTarget]Glink"{

  style.double-border:true


}
l:Location{

    style.multiple:true}
t: "map[Tag][]Location"{
     style.double-border:true

   }
t->l
gl:gLink{
  style.multiple:true
    "Defs []Location"{
      style.multiple:true}
  "Refs []Location"{
    style.multiple:true}
}

ds:"map[lsp.DocumentURI]DocumentData"{
     style.double-border:true
   }
#gl."Defs []Location"->d
ds->dd
dd:"DocumentData"{
  style.multiple:true
    "Trees"{
      style.multiple:true}
  Content
    Headings
    FootNotes
}
fns: "map[string]FootNoteRef"{
       style.double-border:true
     }
fnr:"FootNoteRef"{
      style.multiple:true
        "Def *lsp.Range"
        "Refs [].Range"{
          style.multiple:true
        }
      Excert
    }
fns->fnr
hds:"map[string]Subheading"{
  style.double-border:true
}
shd:"Subheading"{
      style.multiple:true
        "Def lsp.Range"
        "Refs [].Range"{
          style.multiple:true
        }
    }
hds->shd
sm:"[]string"{
          style.multiple:true
  }

dd.FootNotes->fns
dd.HeadingsStore->hds
s:{
Tags
GLinkStore
DocStore
OtherFiles
}
s.Tags -> t
s.GLinkStore->gls
gls->gl
s.DocStore->ds
s.OtherFiles->sm
```
