# DDDL

Die worldiety _Domain-Driven Design Language_ kann zur Unterstützung einer strukturierten Anforderungserhebung dienen.
Dazu wird eine Interpretation verschiedener Elemente basierend auf den strategischen Mustern des Domain-Driven Design mit Hilfe einer domänenspezifischen Sprache bereitgestellt.  

Grundsätzlich tragen die TxT-Dateien die Endung _*.ddd_ und können nach Belieben arrangiert und verschachtelt werden.

## Installation

Damit die Vorschaugenerierung richtig funktioniert, ist plantuml und der BPMN-Font erforderlich:

```bash
brew install plantuml
```
Download und Installation der [bpmn.ttf](https://github.com/bpmn-io/bpmn-font/blob/master/dist/font/bpmn.ttf) aus dem bpmn-font Projekt.

## Kommandos

Mit shift+cmd+p öffnet sich die Kommandopalette von VSC.
Die folgenden Kommandos stehen zur Verfügung

* `Vorschau`: Die ist die empfohlene build-in Vorschau.
* `als HTML exportieren`: Dies ist eine statisches Ausgabe zur Weitergabe oder Drucken. VSC hat keine eingebaute Preview für HTML-Dokumente. 
Es gibt dafür die _Live Preview_ als entsprechende Extension von Microsoft.
* `generiere Go Code`: Hiermit wird im aktuellen Workspace Go Code generiert.
Dazu wird das übergeordnete Go-Module gesucht und dort im internal Package für jeden Bounded Context ein eigenes Package generiert.
Gibt es kein Go-Module, wird der Code in den Workspace-Folder geschrieben.

  
## DevOps Tools

Das Projekt enthält neben dem VSC-Plugin und dem LSP-Server auch die folgenden Tools.

### dddc

Der DDD Compiler für das automatisierte Erzeugen von verschiedenen Ausgabeformaten kann wie folgt installiert werden:

```bash
go install github.com/worldiety/dddl/cmd/dddc@latest

# within the working directory of the  *.ddd files project
# e.g. to generate a standalone html file
dddc -format=html -out=index.html

# or generate go code (by default places always in mod-root/internal/<context name>
dddc -format=go
```



## Werkzeug für Workshops

Dieses Werkzeug eignet sich als Ergänzung für die Techniken des _Event Stormings_, des _Domain Story Tellings_ sowie der _BMPN_.
Es erscheint sinnvoll den ersten Workshop mit allen Stakeholdern zu machen und alle Prozesse mit ihren Beteiligten aufzunehmen. 
In den dann folgenden Workshop sollten in Einzelinterviews mit den jeweiligen Domänenexperten:innen die exakte Daten und Ablaufdokumentation entstehen.
Erst hiernach können sinnvolle _User Stories_ mitsamt Akzeptanzkritieren formuliert werden.
Neue Erkenntnisse und Fragen aus dem _Refinement_ müssen in dem gemeinsamen ddd-Modell wieder ergänzt werden, um das gemeinsame mentale Modell fortwährend aktuell zu halten.

#### Event Storming

Diese Methode definiert ein PostIt-basiertes Workshop-Verfahren, das extra zur Erkundung der Domäne (Fachlichkeit) im
Domain-driven Design entworfen wurde.
Kernidee ist, dass man alle Stakeholder in einem Workshop zusammenbringt und von den Arbeitsergebnissen her denkt.
Dazu wird zuerst gezielt nach Domänen-Ereignissen gefragt (Domain Events).
Alle Stakeholder zu beteiligen ist wichtig, da auch diese üblicherweise in ihren eigenen _Bounded Contexten_
professionalisiert sind und daher auch nur jeweils Teile von Geschäftsprozessen wissen.
Domänen-Ereignisse sind in der Vergangenheitsform zu formulieren, als etwas, dass eingetreten ist.
Nachdem alle wesentlichen Ereignisse identifiziert wurden, wozu auch Fehlerzustände gehören, werden in folgenden
Workshop-Iterationen die Akteure (Actors), deren Kommandos (Commands) und Ansichten (Views) bestimmt.
Es können auch noch Elemente für Aggregate-Roots und Externe Systeme festgelegt werden.
Im letzten Schritt kann eine erste Gruppierung in Subdomänen bzw. _Bounded Contexte_ erfolgen.

Das Verfahren ist weder gedacht noch dazu geeignet, um komplexe Prozesse im Detail zu dokumentieren, beispielweise mit
Schleifen und Sonderfällen.
Schwierig erscheint auch die Bedeutung in der Kommunikation und im späteren Entwicklungsprozess.
Mit diesem Verfahren erkundete Domänen, fordern eine ereignisbasierte gemeinsame Sprache (ubiquitous language) und damit
auch eine _CQRS_-Architektur (Command-Query-Responsibility-Segregation) mit _Event Sourcing_.
Ansonsten ist ein kompletter Bruch zu erwarten, der ein Auseinanderdriften der Dokumentation und Implementierung stark
begünstigt.

#### Domain Story Telling

Diese Methode definiert ebenfalls ein Workshop-Verfahren zur Identifizierung der relevanten Prozesse und beschreibt
diese durch das Festlegen von Akteuren, Aktivitäten und Arbeitsobjekten.
Ein Akteur nimmt dazu üblicherweise die Rolle des Subjektes ein, die Aktivität die des Verbs und das Arbeitsobjekt
stellt das Objekt dar, sodass einfach zu lesende gerichte Flussgraphen entstehen.
Damit diese später besser nachvollzogen werden können, muss die Ausführungsreihenfolge durch eine Nummerierung
dokumentiert werden.
Abschließend werden die _Bounded Contexte_ bestimmt und durch das Zeichnen von entsprechenden Begrenzungsrechtecken
dargestellt.

Das Verfahren ist weder gedacht noch dazu geeignet, um komplexe Prozesse für Dritte zu dokumentieren.
Stattdessen ist es nur als Gedankenstütze für die Workshopteilnehmer gedacht.
Insbesondere fehlen Möglichkeiten Sonderfälle und Varianten abzubilden.
Dazu wird empfohlen, entweder mit Kommentaren zu arbeiten oder Prozesse zu kopieren und jede Variante in einer eigenen
Grafik darzustellen.
Hierbei ist kritisch zu bedenken, dass der kombinatorische Aufwand eine nicht mehr beherrschbare Komplexität verursacht,
sodass sich das Verfahren wirklich nur für stark vereinfachte Prozessdarstellungen eignet.

#### BPMN

Die _Business Process Model and Notation_ ist eine grafische Notation zur Dokumentation von Geschäftsprozessen und wie
die UML durch die OMG standardisiert.
Diese Notation ist durch ihren Umfang und Komplexität nur bedingt workshoptauglich, dafür aber besonders geeignet um in
Einzelinterviews mit Domänenexperten einen Arbeitsablauf mit allen notwendigen Facetten zu dokumentieren.
Ziel ist die voll umfassende Prozessdokumentation, sodass auch Dritte in der Lage sind, einen Prozess nachzuvollziehen
und ggf. selbst ausführen zu können.

#### Fazit

Die hier unterstützte Notation für Arbeitsabläufe soll ebenso in Domänen-Erkundungs-Workshops verwendbar sein, wie das
_Event Storming_ und _Domain Story Telling_ aber auch ausdrucksstark genug sein, um im Nachgang alle fachlichen
Sonderfälle für die eigentliche Entwicklung ähnlich der BPMN zu dokumentieren.
Dabei beträgt der Umfang allerdings im Vergleich zur BPMN nur einen Bruchteil, ist aber sprachlich wie technisch
ausreichend vollständig.
Durch die Verwandtschaft mit dem _Event Storming_ dominiert auch hier eine ereignisorientierte Modellierung, die zu
einem Bruch im gemeinsamen mentalen Modell führen kann.
Allerdings kann hier explizit in der fachlichen Modellierung zwischen der Referenzierung mittels Bezeichnern und
beschreibenden Fließtexten unterschieden werden.
Wenn also wirklich ein _CQRS_-System fachlich modelliert werden soll, muss ein Ereignis als _Daten_ Deklaration
ausgeführt werden und eben genau dieser Bezeichner auch im Arbeitsablauf verwendet werden.
Hierbei ist ganz genau zu klären, welche nicht-technischen aber fachlichen Anforderungen an die Daten- und
Prozesskonsistenz zu garantieren sind.

Es gibt viele weitere Techniken zur Abbildung von Anforderungen.
Dazu zählen die nutzerzentrierten _Szenarien_ sowie _User Stories_.
Diese beschreiben einen Prozess vollständig aus der Sicht des Nutzers.
Der hier definierte Arbeitsablauf hat diese Beschränkung nicht, hinterfragt aber auch nicht die Notwendigkeit eines
Geschäftsprozesses, der als gegeben hingenommen wird.
Daher ist es im Workshop sehr sinnvoll, den Sinn und die Bedeutung eines Prozesses zu prüfen und statt des
_Ist-Zustands_ besser den zukünftig erwünschten und ggf. neu digitalisierten Ablauf zu dokumentieren.

## Beispiel

```ddd

"TODO: Format noch zu klären."
Typ Kundennummer = Text

"
@Torben: Heißt das wirklich Kunde oder wird eigentlich ein anderer Name verwendet?
* Verbraucher?
* Gewerk?
"
Auswahl Kunde {
    Firmenkunde
    oder Privatkunde
}

"Diesen _Bounded Kontext_ konnten wir bereits identifizieren."
Kontext Scoring {

    Daten Kunde {
        Kundennummer
    }
    
    Synonym Verbraucher = Kunde
    
    "Ganz kompliziert, geht aber auch ohne Body und ohne Parameter, um die Erkundung der Domäne zu vereinfachen."
    Aufgabe ImmoFinanzScoringBerechnen(Speichern, EingereichterAntrag) -> (BerechneterAntrag , Abgelehnt|BerechnungsFehler) {
        wenn "eingereichter Antrag Basisscore < 5" -> Abgelehnt
        sonst {
            wenn "Sonne scheint" {
                -> Berechnungsfehler
            }
        }
        
        solange "x > y"{
            PrüfeErneut
        }
        
        Speichern(EingereichterAntrag)
       
       -> BerechneterAntrag
    } 
    
    Aufgabe Speichern
    
}


 


```

## Strukturierung nach DDD

Im Domain-Driven Design werden strategische und taktische Mustern unterschieden.

### Kontext (engl. context)

Die Zerlegung in verschiedene _Bounded Contexte_ gehört zu den strategischen, also organisatorischen Mustern.
Dabei steht häufig im Vordergrund, dass die gesamte Domäne (z.B. ein Unternehmen) in verschiedene und möglichst entkoppelte Kontexte aufgeteilt werden kann.
Um diese Zerteilung vorzunehmen, gibt es kein Regelwerk, da jede Domäne individuell sein kann.
Heuristiken zur Erkennung von verschiedenen _Bounded Contexts_ sind:
* derselbe Begriff taucht mit unterschiedlichen Definitionen auf
* Ein Domänenexperte verweist auf einen anderen Domänenexperten
* Es sind _Pivot Events_ erkennbar
* Es ist eine natürliche und schwache Kopplung mit einer (natürlichen) eventuellen Konsistenz erkennbar.

Der Projektmanager bestimmt danach die Priorität innerhalb der Projektumsetzung.
D.h. es wird mit dem Kunden zusammen bestimmt, welches der wichtigste _Bounded Context_ ist.
Von der erfolgreichen Umsetzung dieses Kontextes hängen alle anderen Kontexte ab, was als Upstream-Kontext bezeichnet wird.
Upstream und Downstream bezeichnen in dieser Mustersprache **nicht** den Datenfluss, sondern die Team- bzw. Umsetzungsprioritäten.

Beispiel

```ddd
Kontext Buchausleihe {
  Typ Buch
  
  Aufgabe Ausleihen(Buch)->BuchWurdeAusgeliehen
}

Kontext Buchsuche
```

### Aggregat (engl. aggregate)

Ein Aggregat bezeichnet ein Muster, das im Wesentlichen die Grenze einer fachlichen Transaktion beschreibt.
Hierbei muss mit den fachlichen Experten geklärt werden, welche Konsistenzgarantien wirklich erforderlich sind.
Als Daumenregel gilt, dass je größer und umfassender die Anforderungen an die atomare und unmittelbar sichtbare Konsistenz sind, desto unflexibler wird die technische Architektur und desto schlechter lässt sich die Entwicklung mit mehreren unabhängigen Teams organisieren.
Kennzeichnend eines _Aggregats_ ist das _aggregate root_. 
Hierbei handelt es sich um eine Entität, die über eine eindeutige ID referenziert werden kann.
Diese Entität stellt ihre eigene Konsistenz (d.h. die Einhaltung der fachlichen Invarianten) per _Information Hiding_ sicher.
Normalerweise ist der Zustand einer Entität auf Seiteneffekte ausgelegt, d.h. mind. das persistente Lesen, Speichern, Bearbeiten und Löschen, was in der Persistenz bzw. Infrastrukturschicht implementiert wird.
Aggregate dürfen untereinander nur über ihre ID aufeinander Bezug nehmen.
Durch die Schichtentrennung und dem _Information Hiding_ ist ein perfektes Design zudem eigentlich unmöglich technisch auszudrücken: entweder muss das Domänen-Objekt für die Persistenzschicht das _information hiding_ Prinzip verletzten (Deserialisierung und Konstruktion) oder die Domänenschicht erhält Persistenzanteile.

Zudem sind Aggregate die Quelle und die Empfänger von Domänenereignissen.

```ddd
Kontext Buchausleihe {

  Aggregat Ausleihe {
    Daten Nutzer
    Daten Buch
    Daten BuchWurdeAusgeliehen
    
    Aufgabe Ausleihen(Buch)->BuchWurdeAusgeliehen
    Aufgabe Speichern(BuchWurdeAusgeliehen)
  }
}
```

Entwickler:innen können ein Aggregat auf verschiedene Weisen technisch ausdrücken.
Dies kann im objektorientierten Stil (_Klasse_) erfolgen oder auch über eine Reihe von freien (_pure_) Funktionen möglich sein.
Je nach Umfang erscheint häufig ein eigenes _Package_ praktisch, da hierbei das Prinzip des _Information Hiding_ mittels _package private_ Sichtbarkeiten klarer durchgesetzt werden kann.

Ein Aggregat ist die kleinstmögliche Einheit für einen Microservice.

## Typen

Ein Typ definiert eine Menge von Werten und damit verbundene Eigenschaften und erlaubte Operationen.

### Boolesche Typen

Es ist Absicht, dass es keinen booleschen Typ gibt,
da diese immer als _Auswahltyp_ domänenspezifisch definiert werden sollen.
Entwickler:innen haben später die Möglichkeit diese je nach Sprache als Konstanten oder eigenständige Typen abzubilden.

Schlechtes Beispiel:

```ddd
Daten Kreis {
  // verwende keine booleschen Definition, da diese
  // fast immer fachlich unklar bis falsch sind.
  bool als Farbig 
}
```

Besseres Beispiel:

```ddd
Daten Kreis {
  RenderOption
}

Auswahl RenderOption {
  Outline
  oder Farbig
}
```

### Zahl (engl. Number)

Zahl definiert einen unspezifischen Zahlentyp, der sowohl eine Ganzzahl oder eine Fließkommazahl sein kann.
Entwickler:innen müssen daraus einen Auswahltyp machen, um möglichst exakt zu sein.

```ddd

// dies ist die vordefinierte Semantik
Auswahl Zahl {
  Ganzzahl
  oder Gleitkommazahl
}
```

### Ganzzahl (engl. Integer)

Eine Ganzzahl definiert eine positive oder negative natürliche Zahl.
Sofern der Wertebereich nicht eingeschränkt ist, sollten Entwickler:innen einen _int64_ bzw _long_ als Basisdatentyp wählen.

### Gleitkommazahl (engl. Float)

Eine Gleitkommazahl definiert einen IEEE Float.
Grundsätzlich sollten Entwickler:innen einen _float64_ bzw. _double_ als Basisdatentyp wählen.

### Text (engl. String)

Ein Text sollte immer eine Byte-Sequenz im UTF-8 Format sein.
Andere Formate kommen heute in der Verarbeitung intern eigentlich nicht mehr vor.
Die Anbindung von Legacy-Systemen sollte entsprechend klar modelliert werden, sodass offensichtlich ist, wenn es sich nicht um UTF-8 handelt.
Je nach Programmiersprache bzw. Plattform müssen Entwickler:innen dann auch abweichende Basistypen bzw. Decoder wählen.


### Liste (engl. List)

Eine Liste ist eine generische Datenstruktur zum Halten einer geordnete Menge von gleichartigen Werten.

Beispiel

```ddd
Typ Einkaufsliste = Liste[Artikel]
```

### Zuordnung (engl. Map)

Eine Zuordnung ist eine generische Datenstruktur zum Halten von eindeutigen Schlüssel-Wert-Beziehungen.
Entwickler:innen sollten die Gleichheit von Schlüsseln anhand der Werte prüfen und nicht an der technischen Eindeutigkeit von Zeigern oder Referenzen.

Beispiel

```ddd
Typ Telefonbuch = Zuordnung[Name, Telefonnummer]
```

### Menge (engl. Set)

Eine Menge ist eine generische Datenstruktur zum Halten von eindeutigen Werten.
Die Reihenfolge ist undefiniert.
Entwickler:innen sollten die Gleichheit von Werten prüfen und nicht an der technischen Eindeutigkeit von Zeigern oder Referenzen.

Beispiel

```ddd
Typ GewählteOptionen = Menge[BonbonSorte]
```

### Auswahl (engl. choice)

Ein Auswahltyp ist ein algebraischer Koprodukt-Typ, bei dem nur jeweils genau einer der definierten Typen angenommen werden kann (disjunkt).
Entwickler:innen kennen dies aus anderen Sprachen z.B. als _Sealed Type_ (Java), _enum_ (Rust) oder per |-Notation (Typescript).
Im Grunde lässt sich dies in jeder Sprache mit Polymorphismus abbilden, entweder über klassische Vererbung oder mittels Interfaces.

Beispiel

```ddd
Auswahl Kreditablehnung {
  EigenanteilZuNiedrig
  oder KeinKunde
  oder ExternesScoringRot
}
```

### Daten (engl. data)

Ein Datentyp ist ein algebraischer Produkttyp, bei dem alle Elemente gleichzeitig vorhanden sind.
Siehe auch _Auswahltyp_.
Entwickler:innen kennen dies als technisches _Struct_, _Record_ oder _Class_.

Beispiel

```ddd
Auswahl Kunde {
  Kundennummer
  und Vorname
  und Nachname
  und Adresse
}
```

Häufig kommt es in Geschäftsprozessen zu einem mehrstufigen Arbeitsablauf, der Daten in mehreren Schritten transformiert.
Diese sollten wann immer möglich so exakt wie möglich notiert werden.
Optionale Felder sind dann meist unnötig und ein Zeichen schwacher oder gar falscher Domänenmodellierung.

Schlechtes Beispiel:

```ddd

// Bei dieser Modellierung ist alles optional.
// Eine fachliche Reihenfolge der Prozesse ist nicht erkennbar.
// Dementsprechend kann weder der Kunde noch der
// Entwickler später nachvollziehen, wann und welche
// Zustände im Arbeitsablauf zulässig sind.
Daten Warenkorb {
  Artikelliste?
  Zahlmethode?
  Versandadresse?
  Rechnungsadresse?
  AGBAkzeptiert?
}
```

Besseres Beispiel:

```ddd
Daten UngeprüfterAuftrag {
  Text? als UngeprüfterVorname
  Text? als UngeprüfterNachname
  Text? als UngeprüfterGewünschteZahlmethode
  Text? als UngeprüfteVersandadresse
  Text? als UngeprüfteKundennummer
  Text? als UngeprüfteAbweichendeRechnungsaddresse
}

Typ Artikelliste = Liste[Artikel]
Daten ArtikelMitZahlmethode {
  Artikelliste
  und Zahlmethode
}

Daten WarenkorbMitRechnungsadresse {
  ArtikelMitZahlmethode
  und RechnungsAdresse
} 

Daten BestellbarerWarenkorb {
  Artikelliste
  Zahlmethode
  Versandadresse
  Rechnungsadresse
}

Aufgabe ArtikelPrüfen(UngeprüfterAuftrag) -> Artikelliste
Aufgabe ZahlmethodeWählen(Artikelliste) -> ArtikelMitZahlmethode
Aufgabe RechnungsadresseWählen(ArtikelMitZahlmethode) -> WarenkorbMitRechnungsadresse
Aufgabe Warenkorb(WarenkorbMitRechnungsadresse)-> BestellbarerWarenkorb

Aufgabe Bestellprozess()->(Bestellung,Bestellfehler){
  ArtikelPrüfen
  ZahlmethodeWählen
  RechnungsadresseWählen
  Warenkorb
}

```

## Annotationen

Annotationen verwenden die Bezeichnergrammatik und erlauben eine flexible Anzahl an Schlüssel-Wert Paaren zu definieren.
Dadurch sind keine weiteren belegten Schlüsselwörter in der Sprache erforderlich und das System lässt sich einfacher erweitern.
Derzeit sind die folgenden Annotationen spezifiziert.

### @Ereignis (engl. @event)

Domänenereignisse dienen insbesondere auch der Kommunikationen zwischen _Bounded Contexten_. 
Innerhalb eines _Bounded Context_ ist das Erzeugen und Konsumieren von eigenen Domänenereignissen je nach Architektur mitunter nur unnötiger Aufwand und die derzeitige Empfehlung ist es, diese eher zu vermeiden.
Für eine starke Entkopplung von Aggregaten oder ganzen _Bounded Contexten_ eignet sich auf Kosten von zusätzlicher Duplikation und mit der Einführung von eventueller Konsistenz allerdings sehr gut.
Um noch spezifischer zu formulieren, sollte daher jedes Ereignis noch um Eigenschaften ergänzt werden, die anzeigen, ob es sich um ein- oder ausgehende Ereignisse handeln.


Beispiel

```ddd
@Ereignis(eingehend,ausgehend)
Daten WarenkorbBestellt 
```

Entwickler:innen wissen dann, dass für die zu implementierenden Datenstrukturen ein technischer Adapter an eine _Event Sourcing_ Architektur erstellt werden muss.
Typischerweise sollte ein _Bounded Context_ oder mind. ein Microservice mit einem einzelnen Aggregat, seinen eigenen _Event Store_ als Audit-Log halten sowie Event-unabhängige Persistenz, das die _Single Source of Truth_ enthält.
Hintergrund sind die folgenden Argumente:
* Event Sourcing in seiner Reinform ist in den meisten Szenarien absolut ineffizient
* Bei einer Kopplung an einen zentralen Event Store, gibt es eine ungewollte Kopplung an eine zentrale Datenbank, die einer unabhängigen Microservice-Struktur entgegensteht.
* Anforderungen der DSGVO erlauben regelhaft keine unbeschränkte Aufzeichnung von personenbezogenen Daten.
* Migrationscode für die allerersten Events müsste über die komplette Lebensdauer mitgeführt werden.  

### @Fehler (engl. @error)

Fehlerzustände bzw. fehlerartige Randfälle sollten explizit als Fehler gekennzeichnet werden.
Entwickler:innen modellieren dies in der Domäne später ggf. auf sprachspezifische Weise beispielsweise mit _errors_ oder _Exceptions_.

Beispiel

```ddd
@Fehler
Typ PersonNichtGefunden = PersonID 

@Fehler
Auswahl PersonenFehler {
  PersonNichtGefunden
  oder AusweisAbgelaufen
  oder UnerlaubteAktion
}
```


### @Fremdsystem (engl. @error)

Ein _Bounded Context_ interagiert zwecks Seiteneffekte eigentlich immer mit externen Fremdsystemen.
Dazu können neben klassischen Systemen wie E-Mails oder Notifications auch Infrastruktur-Themen gezählt werden, wie das Persistieren in einer Datenbank oder in S3.
Wird also ein Fremdsystem aus der Sicht der Domänenmodellierung identifiziert, so sollte dafür eine _Aufgabe_ ohne Block definiert werden.
Entwickler:innen wissen dann, dass innerhalb der Domänenimplementierung nur ein Funktionstyp bzw. ein Interface erstellt werden muss und die eigentliche Implementierung später per Injection reingereicht wird.

Beispiel

```ddd
@Fremdsystem
Aufgabe PersonSpeichern(Person)->PersonenFehler
```



## Grammatik

![dddl grammar](/ebnf/grammar.png)


## Roadmap

* Aufgaben returnen ist etwas anderes als Ereignisse zurückliefern: Bedeutung dokumentieren
* Externe Ereignisse?
* Event senden?
* Annotationen
  * Fehler
  * Ereignis