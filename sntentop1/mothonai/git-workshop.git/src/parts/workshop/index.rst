========
Workshop
========

-------
Prework
-------

.. toctree::

   ./environment.rst
   ./create-repo.rst
   ./remotes.rst

   
:doc:`game`

----
Sync
----

.. _git_fetch_all:

git fetch --all
^^^^^^^^^^^^^^^

::

    git fetch --all


--------------------
Create an new branch
--------------------

::

    git checkout -b <username>


-----------------------------
Add a README.md file, locally
-----------------------------

::

   # gwr-00

   git workshop repository - 00

-----------------
Add it to ``git``
-----------------

::

    git add README.md

---------
Commit it
---------

::

    git commit -m "Ready, set, go"

---------------------
Push it to the remote
---------------------

::

    git push origin <username>

After everyone has pushed to their remotes :ref:`git_fetch_all` again.


