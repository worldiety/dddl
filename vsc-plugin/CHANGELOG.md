# Changelog

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