# Changelog

## 0.0.22
Fixed crashes in various context assignee conditions.


## 0.0.21
Official support for the following annotations:
* @event(incoming,outgoing)|@Ereignis(eingehend,ausgehend)
* @external|@Fremdsystem
* @error|@Fehler

Various bugfixes, more spec and tutorial content.

## 0.0.20
Using functions as parameters renders as dependency instead of event.
Added while-loop support.
Added annotations with key and key-value notation.

## 0.0.19
New Preview Layout, allowing larger bounded contexts without hiding the entire view.
Fixed SVG distorted images.
New features:
 * Supporting rendering of white space in type names 
 * Added support for optional type declarations, like `Text?` which shall be interpreted as `choice OptText {Text or None}`
 * Supporting `aggregate` grouping
 * added optional field name support (`data UnparsedAddress { Text as UnparsedZip }`)


## 0.0.18
Rendering UML trees only with a depth of 2, inserting more comments for base types.

## 0.0.17
Various bug fixes: scroll stable preview, replace <> with [] in UML preview, render base types as Class, invert extension between types.

## 0.0.16
Redesigned grammar based on experience reports.

## 0.0.15
Added basic code generator support for Go.

## 0.0.14
Fixed rendering and resolving of multiple context declarations.

## 0.0.13
Added html export and dddc commandline tool.

## 0.0.12
Removed : from TODO declarations, to be consistent.
Improved and added more hover documentations, introduced "Wiederholung" prefix keyword for consistency and allowing Subprocess identifiers.
Added Auto-Link support over type declaration resolver.

## 0.0.11
Refactored Linting-Engine to Workspace modell and Resolver-Engine and added support for @-Notations and fixed a lot of linter render and linking bugs. Added a lot of tooltips.

## 0.0.10

Replaced (broken) semantic tokens with static textmate grammar tokens.
Started rewriting of hover texts with more insights in their meaning and usage.

## 0.0.9

More Bugfixes and changed a lot of grammar to have less ambiguous parsing behavior in edge cases. This causes Contexts, Data and Workflows to require { } brackets. Supporting more BPMN-Style notations.

## 0.0.8

More Bugfixes, partial AST DE/EN-Keyword support.
Changed Grammar for Data and Workflow.
Introduced hybrid notation for Event Storming / Domain Story Telling / BPMN.