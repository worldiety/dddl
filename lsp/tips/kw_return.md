### Erläuterung: Definition des Prozessendes

An dieser Stelle ist _%[1]s_ ein Schlüsselwort, dass das reguläre Ende eines Arbeitsablaufes beschreibt.

Ein Prozessende wird im _Event Storming_ nicht explizit modelliert, weshalb es hier als übliches orangefarbiges PostIt modelliert wird.
In der BPMN wird dies eindeutig als _Endereignis_ bezeichnet und für die bessere Dokumentation auch so übernommen.

Die folgenden Situationen sind gültig:
* folgt nach dem Schlüsselwort ein Bezeichner, muss es sich um einen Bezeichner eines Datentyps handeln
* folgt nach dem Schlüsselwort ein Literal, bedeutet dass, das der Rückgabewert noch nicht eindeutig definiert wurde.
* folgt nach dem Schlüsselwort weder Literal noch Bezeichner, handelt es sich folglich um einen Prozess, der nur Seiteneffekte erzeugen soll.
Eine solche Modellierung kann es zwar geben, sollte aber möglichst vermieden werden, da er sich nur sehr schlecht testen und in anderen Prozessen komponieren lässt.
Logischerweise kann diese Schreibweise nur der letzte Ausdruck innerhalb eines Blocks sein, da ansonsten das nächste lexikalische Token zum Schlüsselwort gezählt wird und _merkwürdige_ Parserfehler verursachen kann.

Für den Entwickler bedeutet ein _%[1]s_ die Definition konventioneller Rückgabewerte einer freien Funktion oder einer Methode.
Die Menge aller möglichen Rückgabewerte kann technisch z.B. als polymorpher Typ, Interface oder _Tagged Union_ erfolgen.