.. _branching:

#########
Branching
#########

**************
Quick overview
**************

Definition
==========

A **branch** is a series of changes/commits. In a repository there can be more
than one branch serving various purposes.

Default is “master” or, more recently we see the upsising of, “main”.

Cheatsheet
==========

The following sequence of commands is showing a typical workflow:

- ``git checkout [main-branch-name]`` - checkout to the branch you want to start from
- ``git checkout -b [new-branch-name]`` - makes a new branch and checkouts there
- ``git add ...`` - add stuff
- ``git commit ...`` - commit your added changes
- ``git push [remote] [new-branch-name]`` - pushes your branch to the remote

Example usage
=============

There is a project that lacks a feature. A developer (or more) should create a
branch out of the current mainline so they can work on the feature. This is
helpful for a number of reasons, but mostly because partial work commited during
development, won't interfere with the mainline of the project.

After work is complete, the developer should make a request to the maintainer of
the project to pull their fresh feature into the mainline of the project. This
part is called "merge". Git web services call the request for merge as "pull 
request" (PR), as it is also requiring them to pull the branch.

***************
Create a branch
***************

The right conditions
====================

Local only repo
---------------

.. _git-branch:

Check your current branch
^^^^^^^^^^^^^^^^^^^^^^^^^

Be sure to be on the branch you want to branch-off from. You can check this with
::

    $ git branch
    * 01-patch
      dev
      fixing-a-bit
      main
      more-fixes

Since the asterisk is next to "01-patch", but for a new feature we may want to
start from the "main" branch. We can quickly checkout into that branch with

.. _checkout-to-right-branch:

Checkout to the right branch
""""""""""""""""""""""""""""

::

    $ git checkout main

If you encounter an error saying you have changes that might be overwritten, use
:ref:`git-stash` to "stash" them. Then, retry the command.

Now, proceed to make the new branch out of main.

.. _make-new-branch:

Make new branch
"""""""""""""""

::

    $ git checkout -b other-branch

Fresh clone from remote
^^^^^^^^^^^^^^^^^^^^^^^
In a fresh clone, there are no changes. Just :ref:`make-new-branch`.

Working in sync with remote
^^^^^^^^^^^^^^^^^^^^^^^^^^^
Always fetch fresh with:

.. _git_fetch:

git fetch
"""""""""""""

::

    git fetch --all

This will ensure that our local copy has all the information about all branches
and remotes.

*************************
Consilidation of branches
*************************

Multiple branches be like
=========================

::

   -o-o-o-o-o    (main)
     `-o-o-o-o   (other-branch)

At this point, main and other-branch differ.

This can mean that if we would want to 

``git rebase``
--------------

::

   # rebasing will be like
   -o-o-o-o-o          (main)
             `-o-o-o-o (other branch)

``git merge``
-------------

while merging would be like

::

   -o-o-o-o-o          (main)
             `-o-o-o-o (other branch)

which would make it look like

::

   -o-o-o-o-o-o-o-o-o-o (main)

**Last “o” is the merge commit**

This is known as **fast-forward** as well.

