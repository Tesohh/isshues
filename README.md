# isshues
TUI/CLI Project manager for small teams over SSH

## Shorthand syntax
Definitions:
- Component: something preceded by a special symbol that will apply a label, dependency etc.

### Example
Say we want to create some issues in the project `OCVV` (Open ClasseViva).

This shorthand
```
- +feat !crit +frontend add Nuke
```
Will create an issue with:
- the labels `+feat` and `+frontend`
- the `!crit` priority
- since no assignees were specified, it will be automatically assigned to the creator (you)

---

This shorthand
```
- +feat +gfx add Nuke 3D graphics @lallos >1 !low
```
Will create an issue with:
- the labels `+feat` and `+gfx`
- the `!low` priority
- assigned to user `@lallos`
- depends on the issue `#OCVV-1` (which is "add Nuke").

Notice how the order of the components does not matter.

**User**names are parsed in a lenient fashion, so to tag `@lallos`, `@lal`, `@llos`, `@lall` would suffice.
this is to favor extremely short typing (especially in the CLI with no autocomplete).
Specifying a shortname is also allowed, and it takes precedence over the username.

---

This shorthand
```
- +idea consider adding nuke 4D graphics @quantum-team >2
```
Will create an issue with:
- the labels `+idea`
- assigned to the group `@quantum-team`
- depends on the issue `#OCVV-2` (which is "add 3d Nuke...").

Group names are not lenient, and they must be typed out precisely.

---

Finally, this shorthand
```
- +bug fix exploding phone bug @nobody 
```
Will create an issue with:
- the labels `+fix`
- assigned to nobody, not even the creator

## Permissions
`isshues` has two kinds of permissions:

### Global permissions
Apply to anything that isn't _inside_ a project.
For now the only global permission will be `create-projects`.

### Project group permissions
Groups are a feature similar to discord roles.

Multiple users can be members of a group, which is applied on a per project basis (ie. all projects have separate groups).

Project permissions, eg. `write-issues`, `read-issues`, `edit-project`, `delete-project` are given to the group.

Then, all members of the group will have the permission.

To apply permissions on a per-user per-project basis, a "anonymous" group is created automatically when a user joins a project. 
This new "anonymous" group will only contain that user and cannot have other users inside of it.

> [!NOTE]  
> users are members of a project if they are members of _any_ group inside of a project. 
> which means, kicking a user involves removing them from ALL groups (including the "anonymous" group which should be deleted.)
