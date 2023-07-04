### Erläuterung: Definition einer Eingabe

An dieser Stelle ist _%[1]s_ ein Schlüsselwort, dass eine Eingabe (die ein anderer Akteur zuliefert) für die Erfüllung einer Aufgabe innerhalb eines Arbeitsablaufs beschreibt.
Eine (Datei-)Eingabe hat kein richtiges Äquivalent im _Event Storming_. 
Man kann es als ein Ereignis definieren, dass aus einem externen System zuliefert wird, weshalb es mit einem pinken PostIt dargestellt wird.
In der BPMN wird dies mit dem entsprechenden Dokumentensymbol dargestellt.

Eine Eingabe sollte daher auch im BPMN-Sinne verwendet werden, um z.B. einen Dateiupload einer Exceldatei durch den Nutzer zu repräsentieren.

Nach dem Schlüsselwort muss entweder der Bezeichner einer Daten-Deklaration oder ein Literal, dass die Eingabe atomar und für alle Stakeholder eindeutig beschreibt, folgen.
Ein Prozessbezeichner ist nicht zulässig.

Die Datei-Eingabe eines Arbeitsablaufes ist technisch als Parameter einer freien Funktion oder Methode zu verstehen.
