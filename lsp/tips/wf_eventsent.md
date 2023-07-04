### Erläuterung: Zwischenereignis versenden mit Text

An dieser Stelle ist _%[1]s_ ein Text-Literal, dass im Arbeitsablauf ein versendetes Zwischenereignis beschreibt und die Ausführung anschließend fortsetzt.
Durch die Literal-Schreibweise wird ausgedrückt, dass noch keine Definition im Glossar erfolgt ist.

Die BPMN bezeichnet diese Situation als _Zwischenereignis_.
Im _Event Storming_ werden solche Ereignisse nicht eigens definiert, daher werden sie hier auch mit einem orangen PostIt repräsentiert.

_Achtung:_ Die Verwendung eines Literals an dieser Stelle ist ein Hinweis darauf, dass die Definition fehlt und der Entwickler später sich etwas ausdenken muss, was ggf. fachlich nicht richtig ist.

_Tipp:_ Definiere für jeden Fehlerfall einen eigenen Datentyp.
Der Entwickler wird dies wahrscheinlich mit einem Listener oder Messagebus umsetzen.