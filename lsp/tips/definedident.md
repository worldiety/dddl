### Was ist _%[1]s_?

An dieser Stelle wird die Verwendung des Typs mit dem Bezeichner _%[1]s_ definiert.

Eine Definition bezeichnet eine instantiierte Verwendung.
Befindet sich beispielsweise der Bezeichner auf der rechten Seite einer anderen Datentyp-Deklaration, ist es die Verwendung als Definition.

Wird der Begriff innerhalb eines Arbeitsablaufs verwendet, handelt es sich ebenfalls immer um eine Referenz auf die Deklaration.

Wird der Begriff innerhalb eines _Bounded Context_ benutzt, gilt die Verwendung grundsätzlich nur für die innerhalb des Kontextes deklarierten Typen.
Danach darf erst nach Deklarationen im anonymen Bereich (_Shared Kernel_) gesucht und aufgelöst werden.
Deklarationen aus anderen Kontexten (bzw. Subdomänen) dürfen notwendigerweise nicht aufgelöst werden, da es **niemals** zu einer direkten Kopplung zwischen Kontexten kommen darf.

_Tipp:_ Wenn es erforderlich ist, dass auf die Prozesse oder Datentypen anderer Kontexte zugegriffen werden muss, dann sind diese zu duplizieren und entsprechend über die Definition verbal zu verknüpfen. 
Der Entwickler wird später ein Interface (objektorientiertes Design) bzw. einen Funktionstyp (funktionales Design) sowie eine passende Adapter Implementierung (z.B. einen _Driver Adapter_ in der Hexagonal Architektur) erstellen.
Das strategische Design beschreibt hierzu auch strategische Muster (z.B. Customer/Supplier), mit der ein Projektmanager bewusst organisatorische Maßnahmen durchsetzen kann.
Das ist wichtig, um zu bestimmen welcher Bounded Context _zuerst_ implementiert wird und wie die API für nachfolgende Kontexte aussieht.

_Tipp:_ Wenn sehr viele Elemente _kopiert_ werden müssen, ist das ein Indiz dafür sein, dass die Kopplung zwischen zwei Kontexten zu groß ist und es sich eher um einen einzigen Kontext handelt und eine Trennung mehr Schaden anrichtet als es nützt.