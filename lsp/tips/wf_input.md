### Erläuterung: Eingabe mit Text

An dieser Stelle ist _%[1]s_ ein Text-Literal, dass eine Eingabe-Datei bzw. ein abstraktes Eingabe-Modell für eine Aufgabe innerhalb eines Arbeitsablaufs bezeichnet.
Durch die Literal-Schreibweise wird ausgedrückt, dass noch keine Definition im Glossar erfolgt ist.

Eine Eingabe ist unabhängig von jeglicher Oberfläche zu sehen und kann je nach Umsetzung aus den verschiedensten Quellen stammen wie z.B. Nutzeroberfläche, einer Datenbank oder einer Excel-Datei.
Technisch spielt das hier keine Rolle, da es sich in der Schichtenarchitektur ohnehin "nur" um einen (Driver-)Adapter in der Präsentationsschicht handelt.
Der Entwickler wird dies als Parameter einer Funktion oder als Parameter einer Methode modellieren. 

_Achtung:_ Die Verwendung eines Literals an dieser Stelle ist ein Hinweis darauf, dass die Definition fehlt und der Entwickler später sich etwas ausdenken muss, was ggf. fachlich nicht richtig ist.

_Tipp:_ Es spricht nichts dagegen einen anderen Prozess bzw. Arbeitsablauf als Eingabe zu definieren.
Die BPMN erlaubt so etwas auch explizit, was auch technisch durch den Entwickler problemlos umgesetzt werden kann.