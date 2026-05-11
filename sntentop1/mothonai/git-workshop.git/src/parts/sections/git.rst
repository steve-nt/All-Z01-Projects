.. _git_init:

############################
How to make a git repository
############################

``git init``

It will:

-  create a “``.git``” directory under your current directory
-  be an empty repository regardless if you have files in there

Now that we have a git repository, how would we work on/with it?


*********************
Repository definition
*********************

It is a linked list, essentially, that tracks changes in a sequencial
way.

A repository is a directory with files and directories and, yes, it can
include symlinks [3]_ and permissions [4]_.

It has a *work tree* typically, unless it’s a bare repository.

********************
Repository structure
********************

While you probably already worked with ``git`` already, it’s possible
that the definition above might sound right. However, while this might
look right, it isn’t really. Let’s check what is in this very directory
with ``ls -a``:

::

   .
   ..
   .git
   README.md

As we see there is a ``.git`` directory… Let’s check what is in there
with:


That’s a lot of things! And yeah, that’s the repository!

Repository inner structure fs-wise:

::

   HEAD
   config
   description
   hooks
   info
   objects
   packed-refs
   refs


.. [3]
   Symbolic links (symlinks) are simple pointers to other places in the
   filesystem.

.. [4]
   Permissions are rights about an entity in the filesystem. *rwx* are
   the main ones and correspond to read, write and execute permissions.

