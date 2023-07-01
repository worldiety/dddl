package lsp

import (
	"fmt"
	"github.com/worldiety/dddl/lsp/protocol"
	"github.com/worldiety/dddl/parser"
	"log"
	"strings"
)

// Handle a hover event.
func (s *Server) Hover(params *protocol.HoverParams) protocol.Hover {
	file := s.files[params.TextDocument.URI]
	doc, err := parser.ParseText(string(file.Uri), file.Content)
	if err != nil {
		log.Println("cannot parse", err)
		return protocol.Hover{
			Contents: protocol.MarkupContent{
				Kind:  "markdown",
				Value: "## Syntaxfehler\nPrüfe deinen Text und die Fehlermeldung: " + err.Error(),
			},
		}
	}

	tokens := IntoTokens(doc)
	token := tokens.FindBy(params.Position)
	if token == nil {
		log.Println("token not found")
		return protocol.Hover{
			Contents: protocol.MarkupContent{
				Kind:  "markdown",
				Value: "",
			},
		}
	}

	return protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  "markdown",
			Value: s.hoverText(token),
		},
		Range: protocol.Range{
			Start: protocol.Position{
				Line:      uint32(token.Node.Position().Line - 1),
				Character: uint32(token.Node.Position().Column - 1),
			},
			End: protocol.Position{
				Line:      uint32(token.Node.EndPosition().Line - 1),
				Character: uint32(token.Node.EndPosition().Column - 1),
			},
		},
	}
}

func (s *Server) hoverText(token *VSCToken) string {
	switch n := token.Node.(type) {
	case *parser.KeywordContext:
		return fmt.Sprintf("### Was ist _%s_?\nAn dieser Stelle ist _%s_ ein Schlüsselwort, dass einen _Bounded Context_ deklariert. "+
			"Hierbei handelt es sich um eine Subdomäne der gesamten Domäne (Fachlichkeit), die nur eine schwache Kopplung an andere Subdomänen aufweist. "+
			"Indizien für die Grenzen eines _Bounded Context_ sind gleiche Begriffe mit anderen Definitionen sowie die Zuständigkeit anderer "+
			"Domänenexperten.\n\n"+
			"Ein _Bounded Context_ kann in mehreren Dateien und mehrfach innerhalb einer Datei vorkommen. Allerdings darf nur einer davon ein TODO und eine Definition enthalten.", n.Keyword, n.Keyword)
	case *parser.Ident:
		return fmt.Sprintf("### Was ist _%s_?\n"+
			"An dieser Stelle ist _%s_ ein Bezeichner. Bezeichner treten entweder als Deklaration oder Definition auf. "+
			"Eine Deklaration muss eindeutig sein und darf innerhalb eines _Bounded Contexts_ nur einmal erfolgen. "+
			"Eine Definition bezeichnet eine instantiierte Verwendung. "+
			"Befindet sich beispielsweise der Bezeichner auf der rechten Seite einer anderen Datentyp-Deklaration, ist es die Verwendung als Definition. "+
			"Wird der Begriff innerhalb eines Arbeitsablaufs verwendet, handelt es sich ebenfalls immer um eine Referenz auf die Deklaration. "+
			"\n\n"+
			"Die folgenden Bezeichner wurden als Basisdatentypen vordefiniert:\n"+
			"* _Text_ definiert eine UTF-8 Zeichenkette.\n"+
			"* _Zahl_ definiert entweder eine Ganz- oder Gleitkommazahl.\n"+
			"* _Ganzzahl_ definiert eine Zahl wie 5, 13 oder 42.\n"+
			"* _Gleitkommazahl_ definiert eine Zahl wie 2.31.\n"+
			"* _Menge_ definiert ein ungeordnetes Set, bei dem jeder Wert eindeutig ist. Eine Menge muss mit einem Typparameter versehen werden, also z.B. Menge[Zahl].\n"+
			"* _Liste_ definiert eine geordnete List bzw. ein Array oder Slice. Eine Liste kann die selben oder die gleichen Werte mehrfach enthalten und muss mit einem Typparameter versehen werden, also z.B. Liste[Text].\n"+
			"* _Zuordnung_ definiert eine ungeordnete Map, bei dem jedem Schlüssel eindeutig ein Wert zugeordnet ist. Eine Zuoordnung muss mit zwei Typparametern versehen werden, also z.B. Zuordnung[Zahl,Text].\n"+
			"\n\n"+
			"_Tipp:_ Verwende niemals Gleitkommazahlen, wenn es um Geld oder Geldäquivalente geht, also z.B. auch wenn Zeiteinheiten oder Produktionswerte in monetäre Werte umgerechnet werden. "+
			"Grund ist, dass sich bei Gleitkommazahlen (je nach Anwendungsfall) nicht tolerierbare Rundungs- und Darstellungsfehler ergeben. "+
			"Versuche die kleinste nicht mehr teilbare Einheit zu finden, bei Währungen ist dies z.B. der Eurocent-Wert oder bei Zeitwerten z.B. Minuten oder Sekunden (je nachdem was die Fachlichkeit erfordert). \n"+
			"\n"+
			"_Tipp:_ Es gibt absichtlich keinen Wahrheitswert (bool). Modelliere stattdessen mit Hilfe eines Choice-Types (algebraischer oder-Typ bzw. tagged union). "+
			"Softwareentwickler können dies auch typsicher mittels polymorpher Mechaniken wie Vererbung oder die Verwendung von Interface-Marker-Methoden umsetzen.", n.Value, n.Value)

	case *parser.Definition:
		s := shortStringLit(n.Text)
		return fmt.Sprintf("### Was ist _%s_?\n"+
			"Dieser Text ist ein String-Literal und definiert einen Kontext, einen Datentyp oder einen Arbeitsablauf eindeutig. "+
			"Diese Definition ist ein sehr wichtiger Bestandteil der gemeinsamen Sprache (Ubiquituous Language) von allen beteiligten Stakeholdern. "+
			"Stakeholder sind beispielsweise die Domänenexperten, Endkunden, Geschäftsführer, Entwickler, Tester oder Projektleiter, also alle die ein Interesse an dem Projekt haben. "+
			"Dieser Text stellt die Referenzdokumentation für das ganze Projekt dar und muss von allen Beteiligten akzeptiert und verstanden worden sein. \n\n"+
			"_Tipp:_ Die Verwendung von Markdown ist möglich und zur Strukturierung von längeren Texten empfohlen. \n\n"+
			"_Tipp:_ Wird ??? irgendwo im Text verwendet, wird der Text auch den offenen Aufgaben hinzugefügt.", s)

	case *parser.ToDoText:
		s := shortStringLit(n.Text)
		return fmt.Sprintf("### Was ist _%s_?\n"+
			"Dieser Text ist ein String-Literal und annotiert ungeklärte Fragen zu einem Kontext, einem Datentyp, zu einem Arbeitsablauf oder auch zu einem Element innerhalb eines Arbeitsablaufes. "+
			"Bei der Exploration der Domäne treten häufig noch ungeklärte Fragen zu Details auf, die explizit ausformuliert werden sollten. "+
			"\n\n"+
			"_Tipp:_ Die Verwendung von Markdown ist möglich und zur Strukturierung von längeren Texten empfohlen. ", s)
	case *parser.KeywordData:
		return fmt.Sprintf("### Was ist _%s_?\n"+
			"An dieser Stelle ist _%s_ ein Schlüsselwort, dass einen (algebraischen) Datentyp einführt bzw. eindeutig deklariert. "+
			"Bei der Exploration der Domäne ist es insbesondere bei frühen Workshops sinnvoll zunächst nur Namen und Konzepte (d.h. kurze Definitionen) für Datenelemente zu notieren. "+
			"Erst in darauf folgenden Workshops bzw. gezielten Einzelinterviews mit den zuständigen Domänenexperten sollten die genauen Spezifika ermittelt werden. "+
			"\n\n"+
			"Es werden entweder Produktypen (d.h. Tupel, Structs, Records, Classes o.ä.) oder Ko-Produkttypen (choice-type, tagged unions, disjoint unions, variant types o.ä.) unterstützt. \n"+
			"Notationsbeispiele:\n\n"+
			"* Auswahl-Typ: `Daten Kunde { Firmenkunde oder Privatkunde }`\n"+
			"* Produkt-Typ: `Daten Firmenkunde { Name und Rechtsform und Registergericht und Adresse }`\n"+
			"\n\n"+
			"_Achtung:_ Die genaue Definition erfolgt bei der Implementierung immer. "+
			"Entweder interpretiert und ergänzt der Entwickler die fehlenden Sachverhalte stillschweigend und ohne Prüfmöglichkeit oder diese Informationen werden vorher festgelegt. "+
			"Die Informationen werden spätestens beim Grooming und dem Anhängen von _Tasks_ oder _technischen Stories_ an die erstellten _User Stories_ definiert.",
			n.Keyword, n.Keyword)

	default:
		return fmt.Sprintf("%T", token.Node)
	}

}

func shortStringLit(s string) string {
	const limit = 30
	s = strings.Split(s, ".")[0]
	s = strings.Split(s, "\n")[0]
	if len(s) < limit {
		return s
	}

	return s[:limit] + "..."
}
