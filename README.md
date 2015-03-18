#Snap
**Version control for database schemas.**

---

## Overview

Snap is a proof of concept tool to start exploring version control for database 
schemas. Usually when maintaining and updating database schemas over multiple 
environments things start to get confusing very quickly. Snap is a tool 
inspired by [git](http://git-scm.com/) allowing you to manage and interrogate 
snap managed databases.

## Installation

Make sure you have Go installed and correctly configured then issue the 
following command: `go install github.com/nomad-software/snap`.

## Usage

Snap is invoked on the command line by using the program name followed by a 
command and any required arguments. Some commands have optional arguments too.
```
snap command <arguments...> [optional]
```
The following is an overview of each command. To read the entire command 
documentation see snap's built-in help.

| Command | Description |
| :------ | :---------- |
| commit  | Commit changes to a schema. |
| copy    | Copy a database from a specified revision. |
| diff    | Show differences between schema revisions. |
| dump    | Dump the entire schema at a specified revision. |
| help    | View the help. |
| init    | Initialise a database for use with snap. |
| list    | List all managed databases. |
| log     | Show a log of changes to a database schema. |
| show    | Show the changes made at a specified schema revision. |
| update  | Update a database schema to any previously commit change. |
| version | Show version information. |

## Built-in help

Full help is available from within the program, viewable after issuing the 
command `snap help`. Further command specific help is available by specifing 
the command too, like this: `snap help init`.

## Supported environments

 * Currently only [MySql](http://www.mysql.com/) databases are supported.
 * Because of reliance on the external `diff` tool, only Posix environments are 
   currently supported.
