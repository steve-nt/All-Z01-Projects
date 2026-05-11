.. _use-cases:

#########
Use cases
#########

**************
Back-up system
**************

We all used ``git`` already. ``<- That's an assumption``

What a simple workflow would be like
====================================

-  ``git add <file>`` - to add files
-  ``git commit -m <message>`` - to “commit” the added files to the
   repository

Okay, but what is this committing stuff to repository? See :ref:`commits`

******************
Collaboration tool
******************

You can use ``git`` to collaborate with others!!! Hooray!

It’s one of the main things happening in here!

Advice:

-  learn it,
-  practice it,
-  practice it good,
-  **git gud** at it ``>.<``

But really:

-  learn how to contribute,
-  learn how to maintain.

***************
Version control
***************

    A quick detour

Assume git-workshop v1.0.0 [1]_ is out!! Then imagine git-workshop
v1.1.0 is an update to v1.0.0. What are these numbers?!

-  Major v1
-  Minor 0 or 1
-  Patch 0

How to version?
===============

   We can use the tag command

   ``git tag <tag> [<commit>]``

-  It will *tag* a commit with the ``<tag>`` we give it.
-  It can be anything but for versioning, we mostly going to the concept
   described briefly, before.

How to list tags?
=================

   ``git [-P] tag``

-  will list all tags available,

What ``-P`` flag does?
----------------------

-  it’s a ``git`` flag that can be inserted in most of the commands we
   issue
-  disables the use of pager [2]_ on the output
-  useful for scripting

By default, ``git``\ ’s output is piped on your pager which is set by
the ``$PAGER`` or ``$GIT_PAGER`` environment variables.

.. [1]
   This is the most basic versioning scheme. For more information:
   `wikipedia <https://en.wikipedia.org/wiki/Software_versioning>`__.

.. [2]
   Pager is a program that takes as an input a text file (or other) and
   provides us the scrolling abilitiy. Typical pager programs are *more*
   and *less*.
