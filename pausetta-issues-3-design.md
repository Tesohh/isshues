# Questions

## Core
- Fixed Category/Priority (like v2), or defineable c/p like v1? 
  - Github has "labels". I think we can use that for categories, instead of fixed categories.
  - Priority can be a separate field. It would be defined PER PROJECT and a setting would let you set a priority order,
    useful for sorting.
- statuses?
    - Two options:
        - Have fixed statuses: Todo, Doing, Done, Cancelled
        - Have defineable statuses, and the 4 categories (todo, doing, done, cancelled) on all statuses
- Epics?
    - this can be done with just labels. keep it simple
    - when making a new issue, a field lets you pick 
- Due dates / Deadlines?
    - Absolutely. But this can be added later with no issues.
- Project prefixes, use 4 characters
- Do we want "orgs"? So the same server can be used by multiple companies?
    - Honestly no, just have a single org with projects, and people will be added to the right projects.

## Frontend
- We want crossterm events instead of raw input bytes  
- Perhaps, when we receive `data` we can send it to the crosstermbackend (by faking stdin) and then letting it do what it needs to do

## Auth
- Done through SSH keys, entirely.
- More SSH keys can be associated to accounts.
    - Users can freely accept "sign in" requests on their own:
        1. Say I signed up on my Desktop.
        2. I want to login on my laptop.
        3. I ssh and the server asks me what's my name?
        4. a. That name DOES NOT match an account. A new account is created. END
        4. b. That name DOES match an account. A random word is shown on the laptop screen 
        5. On my desktop (trusted device), i'll receive a login request. 
        6. If i type the random word, the new SSH key will be associated to my main account. END 
- Admins can set the behaviour from a cli, either:
    - users can freely sign up (aka assigning a new account to their SSH key), and then admins would set permissions
    - users can freely TRY to sign up, but then the admin must approve requests (and set permissions)

# Models

## User(Id, UNIQUE Username, IsAdmin)
Admins have ALL permissions on all projects.

## GlobalPermission(Id)
Global permissions are applied per user and are valid throughout the whole app, except for admin accounts which skip all permission checks.
Examples: `new_project`

### UserPermissions(UserId, GlobalPermissionId)

## ProjectPermission(Id)
Permissions are applied per project per user, except for admin accounts which skip all permission checks.
Examples: `member`, `read`, `write_issues`, `edit` (edit project)

## Project(Id, Title, Prefix)

### ProjectUserPermissions(ProjectId, UserId, ProjectPermissionId)
Users with `member` permissions are automatically considered members of a project.

## Issue
```go
id
title
description // long form description
code // serial number per project eg KERR-[[100]]
status // todo, progress, done, cancelled. keep it opinionated.

project_id
recruiter_user_id

priority // an integer. In the UI, will be shown as a name, if a "label" is associated to this specific value.
         // eg. LOW = 60, NORMAL = 100, HIGH = 150, CRITICAL = 999
         // with this we can do some crazy calcs
```
in shorthand syntax then use:
+ for labels 
@ for assigning (plus special @NOBODY)
! for predefined priorities, or !<integer> for constant

## IssueAssignees(IssueId, UserId)
## IssueLabels(IssueId, LabelId)

## Label(Id, Name, Color, ProjectId)
In v2 this would correspond to both category roles and tags. Fuck it we merge <C-D-.><C-D-.><C-D-.><C-D-.>

## IssueChatMessage(Id, CreatedAt, Content, UserId, IssueId)
NOTE: we could also even have a discord bot that JUST creates "chat rooms"

## IssueRelationship(FromIssueId, ToIssueId, Kind)
Kind = dependency, ...
#KERR-12 DEP #KERR-13
means 12 depends on 13

so inbound relationships, check ToIssueId
outbounds,                check FromIssueId
