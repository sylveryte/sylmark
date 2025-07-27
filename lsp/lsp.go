package lsp

type DocumentURI string

type InitializeParams struct {
	ProcessID             int                `json:"processId,omitempty"`
	RootURI               DocumentURI        `json:"rootUri,omitempty"`
	InitializationOptions *InitializeOptions `json:"initializationOptions,omitempty"`
	Capabilities          ClientCapabilities `json:"capabilities,omitempty"`
	Trace                 string             `json:"trace,omitempty"`
}

type InitializeOptions struct {
	DocumentFormatting bool `json:"documentFormatting"`
	RangeFormatting    bool `json:"documentRangeFormatting"`
	Hover              bool `json:"hover"`
	DocumentSymbol     bool `json:"documentSymbol"`
	CodeAction         bool `json:"codeAction"`
	Completion         bool `json:"completion"`
}

type ClientCapabilities struct{}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities,omitempty"`
}

type TextDocumentSyncKind int

const (
	TDSKNone TextDocumentSyncKind = iota
	TDSKFull
	TDSKIncremental
)

type CompletionProvider struct {
	ResolveProvider   bool     `json:"resolveProvider,omitempty"`
	TriggerCharacters []string `json:"triggerCharacters"`
}

type WorkspaceFoldersServerCapabilities struct {
	Supported           bool `json:"supported"`
	ChangeNotifications bool `json:"changeNotifications"`
}

type ServerCapabilitiesWorkspace struct {
	WorkspaceFolders WorkspaceFoldersServerCapabilities `json:"workspaceFolders"`
}

type ServerCapabilities struct {
	TextDocumentSync           TextDocumentSyncKind         `json:"textDocumentSync,omitempty"`
	DocumentSymbolProvider     bool                         `json:"documentSymbolProvider,omitempty"`
	CompletionProvider         *CompletionProvider          `json:"completionProvider,omitempty"`
	DefinitionProvider         bool                         `json:"definitionProvider,omitempty"`
	DocumentFormattingProvider bool                         `json:"documentFormattingProvider,omitempty"`
	RangeFormattingProvider    bool                         `json:"documentRangeFormattingProvider,omitempty"`
	HoverProvider              bool                         `json:"hoverProvider,omitempty"`
	CodeActionProvider         bool                         `json:"codeActionProvider,omitempty"`
	Workspace                  *ServerCapabilitiesWorkspace `json:"workspace,omitempty"`
}

type TextDocumentIdentifier struct {
	URI DocumentURI `json:"uri"`
}

type TextDocumentItem struct {
	URI        DocumentURI `json:"uri"`
	LanguageID string      `json:"languageId"`
	Version    int         `json:"version"`
	Text       string      `json:"text"`
}

type DidOpenTextDocumentParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}

type DidCloseTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type TextDocumentContentChangeEvent struct {
	Range       Range  `json:"range"`
	RangeLength int    `json:"rangeLength"`
	Text        string `json:"text"`
}

type VersionedTextDocumentIdentifier struct {
	TextDocumentIdentifier
	Version int `json:"version"`
}

type DidChangeTextDocumentParams struct {
	TextDocument   VersionedTextDocumentIdentifier  `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

type DidSaveTextDocumentParams struct {
	Text         *string                `json:"text"`
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type TextDocumentPositionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}

type CompletionParams struct {
	TextDocumentPositionParams
	CompletionContext CompletionContext `json:"contentChanges"`
}

type CompletionContext struct {
	TriggerKind      int     `json:"triggerKind"`
	TriggerCharacter *string `json:"triggerCharacter"`
}

type HoverParams struct {
	TextDocumentPositionParams
}

type Location struct {
	URI   DocumentURI `json:"uri"`
	Range Range       `json:"range"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

// Hover is
type Hover struct {
	Contents any    `json:"contents"`
	Range    *Range `json:"range"`
}
