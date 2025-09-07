package lsp

type ShowDocumentFx func(uri DocumentURI,external bool, rng Range)error
