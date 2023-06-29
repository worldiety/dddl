package parser

type KeywordEvent struct {
	node
	Keyword string `("event" | "Ereignis")`
}

type KeywordEventSent struct {
	node
	Keyword string `("sent" | "Versendet")`
}

type KeywordActivity struct {
	node
	Keyword string `("step" | "Aktivität")`
}

type KeywordWorkflow struct {
	node
	Keyword string `("workflow" | "Arbeitsablauf")`
}

type KeywordDecision struct {
	node
	Keyword string `("decision" | "Entscheidung")`
}

type KeywordIf struct {
	node
	Keyword string `("if" | "wenn")`
}

type KeywordThen struct {
	node
	Keyword string `("then" | "dann")`
}

type KeywordElse struct {
	node
	Keyword string `("else" | "sonst")`
}

type KeywordActor struct {
	node
	Keyword string `("actor" | "Akteur")`
}

type KeywordData struct {
	node
	Keyword string `("data" | "Daten")`
}

type KeywordReturn struct {
	node
	Keyword string `("return" | "Endereignis")`
}

type KeywordReturnError struct {
	node
	Keyword string `("error" | "Fehler")`
}

type KeywordContext struct {
	node
	Keyword string `("context" | "Kontext")`
}

type KeywordTodo struct {
	node
	Keyword string `"TODO"`
}
