### Erläuterung: Definition eines Auswahl-Typs

%s definiert hier einen neuen Auswahltypen.
Algebraisch handelt es sich hierbei um einen Aufzählungstypen bzw. einen _Tagged Union_.
Die internationale Schreibweise wird mit einem `choice` eingeführt. 

Beispiele:

```ddd
// minimal
Auswahl Kunde

// erweitert
Auswahl Kunde {
    Firmenkunde
    oder Privatkunde
}
```