### Erläuterung: Ausgabe mit Text

An dieser Stelle ist _%[1]s_ ein Text-Literal, dass eine Ausgabe-Datei bzw. ein abstraktes Ausgabe-Modell für eine Aufgabe innerhalb eines Arbeitsablaufs bezeichnet.
Durch die Literal-Schreibweise wird ausgedrückt, dass noch keine Definition im Glossar erfolgt ist.

Eine Ausgabe ist unabhängig von jeglicher Oberfläche zu sehen und kann je nach Umsetzung unterschiedliche erfolgen, wie z.B. mit einer Nutzeroberfläche, einer Datenbank oder einer Excel-Datei.
Technisch spielt das hier keine Rolle, da es sich in der Schichtenarchitektur ohnehin "nur" um einen (Driven-)Adapter in der Infrastruktur- oder Präsentationsschicht handelt.
Der Entwickler wird dies als (ggf. polymorphen) Ausgabewert einer Funktion oder Methode modellieren. 

_Achtung:_ Die Verwendung eines Literals an dieser Stelle ist ein Hinweis darauf, dass die Definition fehlt und der Entwickler später sich etwas ausdenken muss, was ggf. fachlich nicht richtig ist.

_Tipp:_ Es spricht nichts dagegen einen anderen Prozess bzw. Arbeitsablauf als Ausgabe zu definieren.
Die BPMN erlaubt so etwas auch explizit, was auch technisch durch den Entwickler problemlos umgesetzt werden kann.