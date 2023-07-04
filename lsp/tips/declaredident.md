### Was ist _%[1]s_?

An dieser Stelle ist _%[1]s_ ein deklarierter bzw. neu eingeführter Bezeichner.
Bezeichner treten entweder als Deklaration oder Definition auf.
Eine Deklaration muss eindeutig sein und darf innerhalb eines _Bounded Contexts_ nur einmal erfolgen.

Es dürfen auch kontextlose _freie_ Deklarationen erfolgen, die dann einem gemeinsamen _Shared Kernel_ zugeordnet werden.

_Tipp:_ Grundsätzlich sollten _Shared Kernel_ immer vermieden werden, da sie tendenziell immer zum _Big Ball of Mud_ Muster degenerieren.

_Tipp:_ Freie Deklarationen von Arbeitsabläufen sind perfekt geeignet, um kontextübergreifende Gesamtprozesse zu dokumentieren.
Dazu können Teilbereiche des Ablaufs auch immer mit einer _Kontext_-Annotation ergänzt werden.
Die Details dieses Prozesses sollten dann in einem späteren Workshop im jeweiligen Kontext noch detaillierter definiert werden.