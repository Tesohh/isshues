# isshues
TUI/CLI Project manager for small teams over SSH

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
