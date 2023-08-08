### Erläuterung: Definition eines Typen

%s definiert hier einen neuen eigenständigen Typ, der die gleiche Struktur aufweist wie der deklarierte Basistyp.
Es gibt weiterhin die Typen `Auswahl`, `Synonym`, `Daten` und `Aufgabe`.

Beispiele:

```ddd
// minimal, ohne etwas auszusagen
Typ Kunde

// Kundennummer hat die gleiche Struktur wie der Universe-Typ "Text"
Typ Kundennummer = Text
```