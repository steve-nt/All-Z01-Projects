.. _game:

#########################
Git Turn-Based Story Game
#########################

********
Overview
********

A collaborative Git workshop for 2–3 participants.

Each player takes turns adding a line to a shared story file. The workflow
ensures everyone practices fetching, merging, branching, committing, and pushing
to their own remote.

Note that for this game, you don't need to open any kind of editor for the most
parts of it. A terminal is just enough.

**************
Workshop Setup
**************

Each participant must:

- :ref:`set_up_environment`
- :ref:`create_remote_repo`
- :ref:`add_collaborators_read`

Create personal forks/remotes of the main repository called ``git-game`` and grant each other read access.

Remember, when you clone your repository locally, add your teammates remotes.

For example


:: 

    git remote add <username> https://platform.zone01.gr/git/<username>/git-game

Reference: :ref:`add_remotes_locally`



Step 1: Player 1 Starts
=======================

Player 1 initializes the repository and adds the first line:

::

   git clone https://platform.zone01.gr/git/<player1>/git-game
   cd git-game
   echo "Player 1 starts the story." > game.txt
   git add game.txt
   git commit -m "Player 1 starts the story"
   git push origin main

Step 2: Player 2’s Turn
=======================

1. Player 2 should clone locally their fork

::

   git clone https://platform.zone01.gr/git/<player2>/git-game
   cd git-game

2. Add player's 1 remote

::
   
   git remote add player1 https://platform.zone01.gr/git/<player1>/git-game

3. Fetch all remote information

::

   git fetch --all

4. Pull the changes from player1

::

   git pull player1 main

5. Create your branch with name ``turn-1``

::

   git checkout -b turn-1

6. Player 2 should now add their line:

::

   echo "Player 2 adds a twist!" >> game.txt
   git add game.txt

7. Commit the changes with a appropriate message

::

   git commit -m "turn-1 played"

8. Push explicitly the new turn to *your* remote

::

   git push origin turn-1

9. Make available the ``main`` branch in your remote as well

::

   git push origin main

Step 3: Player 3’s Turn
=======================

1. Player 3 clones their fork as well

::

   git clone https://platform.zone01.gr/git/<player3>/git-game
   cd git-game

2. Player 3 adds the remotes of player1 and player2 and fetches all information
   from the remotes

::

   git remote add player1 https://platform.zone01.gr/<player1>/git-game
   git remote add player2 https://platform.zone01.gr/<player2>/git-game
   git fetch --all

3. Make sure you are on ``main`` branch

::

   git checkout main

4. Before playing, you will need to add the ``turn-1`` branch onto ``main``. In
   other words, you are going to do a merge!

::

   git merge player2/turn-1

5. After you made your changes onto ``main``, you will have to push it on your
   repository (remote).

::

   git push origin main

Now, you essentially approved the changes of player2.

6. Create your branch for your turn (``turn-2``)

::

   git checkout -b turn-2

7. You can now add your line:

::

   echo "Player 3 introduces a surprise!" >> game.txt
   git add game.txt
   git commit -m "turn-2 played"

8. Now, you can push your branch to your remote

::

   git push origin turn-2

Step 4: Next Round
==================

1. Player 1 fetches everything

::

   git fetch --all

2. Make sure you are on ``main`` branch

::

   git checkout main

3. Approve player's 3 turn by merging it

::

   git merge player3/turn-2

4. Push the updates on your remote for next player to see

::

   git push origin main

5. Create your branch for your turn (``turn-3``)

::

   git checkout -b turn-3

6. Append your new line and add it to repo like before

::

   echo "Player 1 continues the story" >> game.txt
   git add game.txt
   git commit -m "Player 1 continues the story"

7. Push your branch onto your remote

::

   git push origin turn-3

Continue the game
=================

Repeat the process for 6 rounds, each player creating a new branch each turn,
merging from the previous player’s branch, adding a line, and pushing to their
own remote.

After the forth round, the player that is about to play, has to just accept the
previous player's turn by merging it as before.

We are going to make an asynchronous continuation for a while in order to see
simple cases of rebasing and conflict resolving.

Create conflict
===============

Players now should play out of order. Keep your sequence but instead of waiting
your turn, make a branch locally that its name, reflects your next turn.

Since we started our first turn with 0, by the end of the 6th round, the turn
number should be 17.

Given that, the next player should make a branch ``turn-18``. Everything should
be done as before, but now, before the passing the turn to the next player, the
current one should make another branch called ``turn-21``, add a line as before
and push it as well on their remote.

Next player, should merge the right turn (``turn-18``) and add their turn play
(``turn-19``). After this is done, they should checkout again on ``main`` and 
make a new branch +3, so ``turn-22``, add a new line, commit, push etc.

A typical procedure would look like (assume player1 is playing):

::

    git checkout main
    git fetch --all
    git pull player3 main
    git merge player3/turn-17
    git checkout -b turn-18
    echo "A new sentence" >> game.txt
    git add game.txt
    git commit -m "turn-18 played"
    git checkout main
    git checkout -b turn-21
    echo "Another new sentence" >> game.txt
    git add game.txt
    git commit -m "turn-21 played (early)"
    git push origin main
    git push origin turn-18
    git push origin turn-21

player2 should fetch as always from player's 1 remote and do the following:

::

    git checkout main
    git fetch --all
    git pull player1 main
    git merge player1/turn-18
    git checkout -b turn-19
    echo "A new sentence by player2" >> game.txt
    git add game.txt
    git commit -m "turn-19 played"
    git checkout main
    git checkout -b turn-22
    echo "Another new sentence by player2" >> game.txt
    git add game.txt
    git commit -m "turn-22 played (early)"
    git push origin main
    git push origin turn-19
    git push origin turn-22
    
The same applies for player3, adjust numbers accordingly. After player3 has
finished approving the right turn and added theirs as well as their early turn,
player1 should encounter a conflict.

::

    git checkout main
    git fetch --all
    git pull player3 main
    git merge player3/turn-20

Now, it's time to do a rebase for turn-21 so we can put it in the right position

::

    git checkout turn-21
    git rebase main

At this point, a conflict message should appear. It would look like this:

::

    ❯ git rebase main
    Auto-merging game.txt
    CONFLICT (content): Merge conflict in game.txt
    error: could not apply b4daeac... early turn-12 test
    hint: Resolve all conflicts manually, mark them as resolved with
    hint: "git add/rm <conflicted_files>", then run "git rebase --continue".
    hint: You can instead skip this commit: run "git rebase --skip".
    hint: To abort and get back to the state before "git rebase", run "git rebase --abort".
    hint: Disable this message with "git config set advice.mergeConflict false"
    Could not apply b4daeac... # early turn-12 test


Open an editor to edit the ``game.txt`` file.

Adjust the file by removing the ``>>>>>>``, ``========`` and ``<<<<<<<`` parts
(whole lines) and make sure only the actual sentences are there and in the right
order.

Then proceed by adding the edits:

::

    git add game.txt

And proceed with the rebase to finish the process:

::

    git rebase --continue

An editor would open with the commit message in place. You may remove the
"(early)" part of it, save it and exit.

After this is done, you can proceed pushing the branch on your remote:

::

    git push origin turn-21 --force

The ``--force`` flag is needed here since there is difference with the remote
version of it.

Don't forget to also push the ``main`` branch as well after it, since we did
approved/merged previous player's turn branch.

::

    git push origin main

The next player should be able to proceed by fetching and ultimately repeating
the same steps, only difference would be the turns numbers for their own case.

Make sure that all players at the end of this round have the same stuff.

Hints
=====
- use ``git blame game.txt`` to see the sequencial additions,
- use ``git log`` to check if the commits are in the right order or in case you
  feel lost,
- use ``git status`` to check the status of your workdir.

--------------
Workshop Goals
--------------

Participants will practice:

- Fetching from another player’s remote
- Merging into a local branch
- Creating a new branch for a turn
- Adding and committing changes
- Pushing branches to their own remote
- Resolving merge conflicts collaboratively
- Understanding distributed workflows in Git

----------
Game Notes
----------

- Everyone should only push to their own remote; no designated merger.
- Conflicts may occur naturally — resolving them is part of the game.
- The story evolves collaboratively while each player practices real-world Git
  commands.

