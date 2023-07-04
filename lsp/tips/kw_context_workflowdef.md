### Erläuterung: Kontext-Definition im Arbeitsablauf

An dieser Stelle ist _%[1]s_ ein Schlüsselwort, dass einen _Bounded Context_ innerhalb eines Arbeitsablaufs definiert.
Diese Notation ist nur zulässig, wenn der Arbeitsablauf nicht innerhalb eines Kontexts definiert wurde.
Es ist unklar, was die Benennung eines fremden _Bounded Contexts_ innerhalb eines anderen Kontexts überhaupt bedeuten soll.
Der Entwickler darf derartige Kopplung nicht im Domänen-Modell implementieren, sondern nur als Adapter in der Infrastrukturschicht -  sei es Persistenz oder Präsentation.

Beispiel:

```ddd
Arbeitsablauf ZeitenBuchen {

    Kontext Produktion {
        Ereignis "Auftrag erhalten"
        Aufgabe "Auftrag prüfen"        
    }
}
    
```