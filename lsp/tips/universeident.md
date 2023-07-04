### Erläuterung: Vordefinierter Datentyp

An dieser Stelle wird der vordefinierte Datentyp _%[1]s_ definiert.
Die folgenden Bezeichner sind als Basisdatentypen vordefiniert:

* _Text_ definiert eine UTF-8 Zeichenkette.
* _Zahl_ definiert entweder eine Ganz- oder Gleitkommazahl.
* _Ganzzahl_ definiert eine Zahl wie 5, 13 oder 42.
* _Gleitkommazahl_ definiert eine Zahl wie 2.31.
* _Menge_ definiert ein ungeordnetes Set, bei dem jeder Wert eindeutig ist. Eine Menge muss mit einem Typparameter versehen werden, also z.B. Menge[Zahl].
* _Liste_ definiert eine geordnete List bzw. ein Array oder Slice. Eine Liste kann die selben oder die gleichen Werte mehrfach enthalten und muss mit einem Typparameter versehen werden, also z.B. Liste[Text].
* _Zuordnung_ definiert eine ungeordnete Map, bei dem jedem Schlüssel eindeutig ein Wert zugeordnet ist. Eine Zuoordnung muss mit zwei Typparametern versehen werden, also z.B. Zuordnung[Zahl,Text].


_Tipp:_ Verwende niemals Gleitkommazahlen, wenn es um Geld oder Geldäquivalente geht, also z.B. auch wenn Zeiteinheiten oder Produktionswerte in monetäre Werte umgerechnet werden.
"Grund ist, dass sich bei Gleitkommazahlen (je nach Anwendungsfall) nicht tolerierbare Rundungs- und Darstellungsfehler ergeben.
Versuche die kleinste nicht mehr teilbare Einheit zu finden, bei Währungen ist dies z.B. der Eurocent-Wert oder bei Zeitwerten z.B. Minuten oder Sekunden (je nachdem was die Fachlichkeit erfordert).

_Tipp:_ Es gibt absichtlich keinen Wahrheitswert (bool). Modelliere stattdessen mit Hilfe eines Choice-Types (algebraischer oder-Typ bzw. tagged union).
Softwareentwickler können dies auch typsicher mittels polymorpher Mechaniken wie Vererbung oder die Verwendung von Interface-Marker-Methoden umsetzen.