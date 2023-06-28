# DDDL

Die worldiety _Domain-Driven Design Language_ kann zur Unterstützung einer strukturierten Anforderungserhebung dienen.
Dazu wird eine Interpretation verschiedener Elemente basierend auf den strategischen Mustern des Domain-Driven Design mit Hilfe einer domänenspezifischen Sprache bereitgestellt.  

Grundsätzlich tragen die TxT-Dateien die Endung _*.ddd_ und können nach Belieben arrangiert und verschachtelt werden.

## Bounded Context

Es können beliebig viele _Bounded Kontexte_ definiert werden, wobei bei Wiederholungen in mehreren Dateien nur eine Definition und ein TODO zulässig ist.

```ddd
Kontext Produktbestellung
    
    TODO: "Die Definition muss noch überarbeitet werden."
    
    "Die Produktbestellung ist die zentrale und wertschöpfende Domäne..."

```

## Daten

```ddd
// Daten gehören immer zu einem übergeordneten Kontext
Daten Kunde = Privatkunde oder Firmenkunde
    TODO: "Besser beschreiben"
    
    "Ein Kunde ist ein zu überladener Begriff."
    
Daten Firmenkunde =
    Vorname
    und Nachname
    und Adresse
    
```