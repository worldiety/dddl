### Erläuterung: Deklaration eines Datentyps

An dieser Stelle ist _%[1]s_ ein Schlüsselwort, dass einen (algebraischen) Datentyp einführt bzw. eindeutig deklariert.
Bei der Exploration der Domäne ist es insbesondere bei frühen Workshops sinnvoll zunächst nur Namen und Konzepte (d.h. kurze Definitionen) für Datenelemente zu notieren.
Erst in darauf folgenden Workshops bzw. gezielten Einzelinterviews mit den zuständigen Domänenexperten sollten die genauen Spezifika ermittelt werden.

Es werden entweder Produktypen (d.h. Tupel, Structs, Records, Classes o.ä.) oder Ko-Produkttypen (choice-type, tagged unions, disjoint unions, variant types o.ä.) unterstützt.
Notationsbeispiele:

* Auswahl-Typ: `Daten Kunde { Firmenkunde oder Privatkunde }`
* Produkt-Typ: `Daten Firmenkunde { Name und Rechtsform und Registergericht und Adresse }`

_Achtung:_ Die genaue Definition erfolgt bei der Implementierung immer. 
Entweder interpretiert und ergänzt der Entwickler die fehlenden Sachverhalte stillschweigend und ohne Prüfmöglichkeit oder diese Informationen werden vorher festgelegt.
Die Informationen werden spätestens beim Grooming und dem Anhängen von _Tasks_ oder _technischen Stories_ an die erstellten _User Stories_ definiert.

Beispiel:

```ddd
// Schnell mal eine leere Daten-Deklaration beim Explorieren
Daten Zeitbuchender 

Daten Vorgang {
    TODO "So sieht ein allgemeines TODO aus"
    
    "Vorgang ist soetwas wie ein Ticket. @Torben: Bitte einmal nachfragen, ob ein Ticket wirklich ein Synonym ist."
    
    Vorgangsnummer
    und Titel
    und Beschreibung
}

Daten Beschreibung {
    Klartext
    oder HTML
    oder Markdown
}
```