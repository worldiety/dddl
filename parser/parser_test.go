package parser

import (
	"encoding/json"
	"fmt"
	"testing"
)

const test = `
Kontext Rechnungsgenerator {

	Arbeitsablauf Zeitloggen {
	
		"Super beschreibung"
	TODO: "schrott"
	
			Akteur Mitarbeiter {
				Ereignis "Aufgabe erledigt"
				Aktivität "Auf Abschicken klicken"
			}
		
			Kontext Rechnungsgenerator{
				Akteur Rechnungssteller{
						Ereignis "Monatslog abgeschlossen"
						Aktivität "schließt Monat ab"
		
						Entscheidung wenn "Monat == 30 Tage" dann{
							
						}
		
						Kontext WasAnderes {
							Fehler
						}
					}
			}
		   
		
	
	}

}

Kontext X

Kontext Y {}

Kontext Z{
	Daten Z
}

`

func TestParse(t *testing.T) {
	v, err := ParseText("abc", test)
	if err != nil {
		t.Fatal(err)
	}

	buf, err := json.MarshalIndent(v, " ", " ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buf))
}
