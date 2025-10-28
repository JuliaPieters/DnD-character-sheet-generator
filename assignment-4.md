# Modern Programming Practices - Exam

| Metadata    | Value              |
|-------------|--------------------|
| Course code | CU75045V1          |
| Lecturer    | Loek van der Linde |
| Date        | 27-10-2025         |
| Deadline    | 02-11-2025 23:59   |
| Opportunity | 1st                |

## Do this first

If you haven't already version controlled your code, do so **now**. Initialize a git repository and commit the current status. You have to hand in a git diff comparing your exam with your current code at the end of this exam. If you can't hand in a git diff from your current code, your exam will not be valid and will not be graded.

You don't have to put your repository on GitHub. Just create a local repo and commit the current status.

If you need help with initializing, committing and diffing then look at your exam from CIN where you had to do all of that. 

## Assignment

Read the grading criteria on the last page of this document. If your app can't possibly pass the *Architecture* criterion, you're allowed to rewrite it so that it does[^footnote1]. By now you've hopefully figured out there are only two sane architectural decisions: either a standard three-tier app, or Onion architecture (/Clean/Hegaxonal architecture). 

If you're choosing to rewrite, you can keep using CodeGrade to check whether your code still passes all the tests. When you're happy with your rewrite, make another commit and save the git diff comparing it against your first code. If you rewrite, you have to submit 2 git diffs! One of your first code comparing with your new code, and then another one comparing your exam implementation with your new code.

[^footnote1]: Note that this goes for about half of you: I didn't teach you 11 lectures + had you do loads of reading on good programming just so you can write the equivalent of a 1000 line if-statement and expect to pass. Don't get me wrong, I'm not salty or angry. I just don't understand why I say it's an 80-100 hour assignment and your submission basically amounts to the Gilded Rose starter code pumped up on five barrels of steroids, mostly written by [insert your favorite chatbot]. I hope it was fun writing it at least..?

When you're happy with your (new) code, your exam assignment is below.

### Exam - Weapon damage (simple)

There a suprising amount of depth to calculating weapon damage in D&D, and you're going to implement a simple version of it: `damage die + STR/DEX modifier`. You'll need the https://www.dnd5eapi.co/api/2014 API to retrieve the damage die and Finesse property. Look up what the Finesse property is in the SRD (p64).

Note that you should pre-fetch this data again, i.e. by extending the `enrich` command that you probably wrote during the regular assignments. Make sure the weapon damage is visible in your character sheets!

Below are three test cases:

| Test case                          | Damage |
| ---------------------------------- | ---------- |
| Level 1 Half-Orc Acolyte Barbarian with 17 STR and 14 DEX, equipped with a Greataxe | 1d12 + 3 |
| Level 1 Tiefling Acolyte Ranger with 13 STR and 15 DEX, equipped with a Shortsword | 1d6 + 1 |
| Level 1 Dwarf Acolyte Rogue with 8 STR and 15 DEX, equipped with a Rapier | 1d8 + 2 |

## What do you submit

1. A git diff of your exam code compared to the code that you started with before the exam;
2. In case you opted to rewrite your starting code: a git diff comparing your first starting code with your new starting code;
3. A .zip of your code;
4. A short report containing explanation for the Maintainability and Testing criteria. Take as much or as little words/graphics as you need.

## Grading criteria

You can earn 100 points for the exam. Your grade is your amount of points divided by 10.

There are four entry requirements for this exam:

1. You have finished the CodeGrade assignment on 100%;
2. You successfully implemented the requested feature;
3. Your code compiles;
4. You have submitted all the necessary things.

When you don't meet these criteria, your submission will not be graded and you will not receive feedback.

| Criterion | Points | 0% | 50% | 100% |
|---|---|---|---|---|
| Architecture (measured through extensibility) | 60 | Responsiblities are all over the place. Implementing the feature forces cascading changes throughout the codebase, thereby violating the open/closed principle | One or two functionalities exist in the wrong layer, but otherwise well-written. Implementing the feature forces cascading changes, but barely: in only one or two places, thereby almost adhering to the open/closed principle | Good architecture where all functionalities exist in the right layers. Implementing the features does not force cascading changes, fully adhering to the open/closed principle |
| Maintainability | 20 | No (good enough) effort has been put in explaining why the code is good | Can explain why the code is good but does not show concrete proof, i.e. metrics, (graphic) models or other reputable sources | Shows why the code is good with concrete proof |
| Testing | 20 | No testing has been done | Tests show the functionality passes the happy path | Tests show a combination of manual and automated tests OR show that important edge cases have been handled |
