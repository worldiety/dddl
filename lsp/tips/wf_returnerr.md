### Erläuterung: Fehlerdefinition mit Text

An dieser Stelle ist _%[1]s_ ein Text-Literal, dass den Arbeitsablauf beendet und ein Ausnahmeereignis beschreibt.
Durch die Literal-Schreibweise wird ausgedrückt, dass noch keine Arbeitsablauf-Definition im Glossar erfolgt ist.

Die BPMN bezeichnet diese Situation als _Fehler-Endereignis_.
Im _Event Storming_ werden solche Ereignisse mit einem roten PostIt repräsentiert.

_Achtung:_ Die Verwendung eines Literals an dieser Stelle ist ein Hinweis darauf, dass die Definition fehlt und der Entwickler später sich etwas ausdenken muss, was ggf. fachlich nicht richtig ist.

_Tipp:_ Definiere für jeden Fehlerfall einen eigenen Datentyp.
Der Entwickler wird dies als (ggf. polymorphe) Exception oder als Error-Value einer Funktion oder Methode modellieren. 

_Tipp:_ Es ist auch möglich einen (Transaktions-)Abbruch (vgl. BPMN) bzw. Prozesseskalation auszudrücken, indem weder Literal noch Bezeichner angegeben wird.
Es muss dann das letzte Element eines Blockes sein.
Technisch wird ein Entwickler später entweder eine _RuntimeException_ oder eine _panic_ einsetzen, um diese Eskalation auszudrücken.
Ursächlich für diese Modellierung können technische Randfälle sein, die durch manuelle Prüfung bearbeitet werden müssen.