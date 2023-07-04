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

// kontextfreier Arbeitsablauf

Arbeitsablauf AuftragPlatzieren {
    "AuftragPlatzieren ist der Kernprozess."

    Kontext Produktion {
        Akteur Mitarbeiter {
            Ereignis "Auftrag erhalten"
            Aufgabe "Auftrag prüfen" 
        }
    }


    Kontext Einkauf {
        Akteur HauptEinkäufer {
            Ereignis MaterialWurdeVerbraucht
            Aufgabe "Preise prüfen und neues Material bestellen"
            Entscheidung wenn "Material lieferbar" dann {
                Zwischenereignis ArtikelBestellungAusgelöst
            } sonst{
                // Bezeichner
                Aufgabe FehllisteErstellen {
                    Ansicht Dashboard
                    Eingabe "Kladde vom Mitarbeiter"
                    Ausgabe Exceldatei
                }
                Zwischenereignis PanikEMailAnChef
                Fehler KeinMaterialVerfügbar
            }
        }
    }

    Kontext Produktion {
        Akteur Mitarbeiter {
            Wiederholung solange "noch kein Material da"{
                Aufgabe "busy waiting"
            }
            Aufgabe "Beende Fertigung"
            Zwischenereignis "Fertigung beendet"
        }
    }

    Kontext Fakturierung {
        Akteur Rechnungssteller {
            Ereignis "Fertigung beendet"
            Aufgabe RechnungsErstellung 
                TODO "@Torben: Dieser Prozess muss noch modelliert werden"
            Endereignis "Auftrag fertig"
        }
    }
}

// ggf. in anderen Dateien

// Eine leere Kontext-Deklaration kann als Erinnerungsstütze nützlich sein
Kontext Rechnungsstellung

Kontext Zeiterfassung {
    TODO "
    Es bleibt noch viel zu tun:
        * @Torben: Rechtslage prüfen
        * @Olaf: Format der Stundenzettel zeigen lassen
    "
    
    "Die Zeiterfassung ist ein _spannender_ **Kontext**."
    
    Daten Vorgang
    
    Arbeitsablauf ZeitBuchen
}


```
