### Erläuterung: Definition einer Ausgabe

An dieser Stelle ist _%[1]s_ ein Schlüsselwort, dass eine Ausgabe (die ein anderer Akteur benötigt) für die Erfüllung einer Aufgabe innerhalb eines Arbeitsablaufs beschreibt.
Eine (Datei-)Ausgabe hat kein richtiges Äquivalent im _Event Storming_. 
Man kann es als gewöhnliches Ereignis eines _Commands_ definieren, weshalb es als orangefarbenes PostIt dargestellt wird.
In der BPMN wird dies mit dem entsprechenden Dokumentensymbol dargestellt.

Eine Ausgabe sollte daher auch im BPMN-Sinne verwendet werden, um z.B. einen Dateidownload einer Exceldatei durch den Nutzer zu repräsentieren.

Nach dem Schlüsselwort muss entweder der Bezeichner einer Daten-Deklaration oder ein Literal, dass die Ausgabe atomar und für alle Stakeholder eindeutig beschreibt, folgen.
Ein Prozessbezeichner ist nicht zulässig.


Die Datei-Ausgabe eines Arbeitsablaufes ist technisch als Rückgabewert einer freien Funktion oder Methode zu verstehen oder kann aber auch über einen Listener oder einen Messagebus erfolgen.