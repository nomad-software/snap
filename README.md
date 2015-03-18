#Snap
**Version control for database schemas.**

---

## Overview

Snap is a proof of concept tool to start exploring version control for database 
schemas. Usually when maintaining and updating database schemas over multiple 
environments things start to get confusing very quickly. Snap is a tool 
inspired by [git](http://git-scm.com/) allowing you to manage and interrogate 
snap managed databases.

## Notes

 * Currently only [MySql](http://www.mysql.com/) databases are supported.
 * Because of reliance on the external `diff` tool, only Posix environments are 
   currently supported.
