# isshues
# Starting design choices / questions

## Core
- Fixed Category/Priority (like v2), or defineable c/p like v1? 
  - Github has "labels". I think we can use that for categories, instead of fixed categories.
  - Priority can be a separate field. It would be defined PER PROJECT and a setting would let you set a priority order,
    useful for sorting.

  - Actually NO; let's do labels, with some having additional data (eg. color). 
  - Priorities are fixed, but in the server config you can set "thresholds" for integer values. 
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

### tea Messages, caching and other optimizations
a possibility for optimizing (not doing too many queries) is to do this.
Let's look at an example:
1. A user completes an issue #13
2. a `tea.Cmd` is issued to edit the issue in the database
3. a `tea.Msg` tells all clients that #13 was updated
4. now here we have to choose: 
    - do we send the new data together with the `tea.Msg`?
    - or do we just say "hey it updated, go fetch it yourself"
        - in this case a *cache* could be useful.
            - does the requested issue already exist in the cache? return it
            - does it not exist? query it and return it
            - did we get a tea.Msg telling us that a certain issue is updated? requery it
    - or we can also do it like this:
        - client A updates db entry
        - client A will fetch and broadcast_exclusively (everyone except A) the new db entry to all clients X
        - client A refreshes it's view(s) accordingly
        - all clients X will refresh their view(s) accordingly
    - or
        - client logs in, selects a project.
        - load all necessary data from that project (ie. all issues, labels, priorities....)
            - except that data is already "on" the server (keep in mind "clients" are actually on the server)
            - so what do we do in this case?
            - use the cache method?
            - fuck it we can just run a query. shouldn't be that slow... plus we're not scaling to millions of users...
        - when another client changes something about 


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
         // eg. a "heat" statistic which is the average of the priorities
```
By the way, in shorthand syntax then use:
+ for labels 
@ for assigning (plus special @NOBODY)
! for predefined priorities, or !<integer> for constant
> for dependencies

## IssueAssignees(IssueId, UserId)
## IssueLabels(IssueId, LabelId)

## Label(Id, Name, Color, ProjectId)
In v2 this would correspond to both category roles and tags

## IssueChatMessage(Id, CreatedAt, Content, UserId, IssueId)
NOTE: we could also even have a discord bot that JUST creates "chat rooms"

## IssueRelationship(FromIssueId, ToIssueId, Kind)
Kind = dependency, ...
#KERR-12 DEP #KERR-13
means 12 depends on 13

so inbound relationships, check ToIssueId
outbounds,                check FromIssueId
