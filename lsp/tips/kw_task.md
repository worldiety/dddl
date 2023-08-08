### Erläuterung: Definition einer Aufgabe

An dieser Stelle ist _%[1]s_ ein Schlüsselwort, dass eine Aufgabe innerhalb eines Arbeitsablaufs beschreibt.
Eine Aufgabe entspricht einem blauen PostIt-_Command_ im _Event Storming_ bzw. einer Aktivität in der BPMN.
Nach dem Schlüsselwort muss entweder der Bezeichner eines anderen Arbeitsablaufes aus demselben Kontext bzw. ein kontextfreier Arbeitsablauf folgen oder ein Literal, dass die Aufgabe atomar und für alle Stakeholder eindeutig beschreibt.
Die internationale Schreibweise ist `task`.

Beispiele:

```ddd
// minimal
Aufgabe ImmoFinanzScoringBerechnen

// mit Eingabe-Ereignis
Aufgabe ImmoFinanzScoringBerechnen(EingereichterAntrag)

// mit Ausgabe
Aufgabe ImmoFinanzScoringBerechnen(EingereichterAntrag) -> BerechneterAntrag | BerechnungsFehler

// mit bedingtem Programmfluss
Aufgabe ImmoFinanzScoringBerechnen(EingereichterAntrag) -> (BerechneterAntrag , Abgelehnt|BerechnungsFehler) {
    wenn "eingereichter Antrag Basisscore < 5" -> Abgelehnt
    sonst {
        wenn "Sonne scheint" {
            -> Berechnungsfehler
        }
    }
   
   -> BerechnterAntrag
}

```



