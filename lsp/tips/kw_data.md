### Erläuterung: Definition eines Datentyps

An dieser Stelle ist _%[1]s_ ein Schlüsselwort, dass einen (algebraischen Summen) Datentyp einführt bzw. eindeutig deklariert.
Bei der Exploration der Domäne ist es insbesondere bei frühen Workshops sinnvoll zunächst nur Namen und Konzepte (d.h. kurze Definitionen) für Datenelemente zu notieren.
Erst in darauf folgenden Workshops bzw. gezielten Einzelinterviews mit den zuständigen Domänenexperten sollten die genauen Spezifika ermittelt werden.
Die internationale Variante lautet `data`.
Es gibt weiterhin die Typen `Auswahl`, `Synonym`, `Typ` und `Aufgabe`.

Beispiele:

```ddd
// Schnell mal eine leere Daten-Deklaration beim Explorieren
Daten Zeitbuchender 

"
So sieht ein allgemeines TODO aus.
Vorgang ist soetwas wie ein Ticket. @Torben: Bitte einmal nachfragen, ob ein Ticket wirklich ein Synonym ist.
"
Daten Vorgang {
    Vorgangsnummer
    und Titel
    und Beschreibung
}

```