.. _set_up_environment:

=======================
Set up your environment
=======================

-  ``git config [--global] user.name "<your name>"``
-  ``git config [--global] user.email "<your email>"``

``git`` is a collaboration tool, hense it's essential to declare who are you so
others can know who introduces what changes into the repository.

===============================
Set up credentials helper store
===============================

Enable the credentials helper so you don't need to input your credentials each
time you do something with remotes:

::

    git config --global credential.helper store


