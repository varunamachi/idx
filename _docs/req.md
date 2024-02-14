# Idx
Simple identity provider


## Initial Ideas

### App
- An 'App' defines an application that needs user management and might or might not have permission based access control


### User
- User is an entity that interacts with one or more apps
- User can represent a real user or a user account used for automation purposes (Service account)
- Normal user will have to give an email, service account might not need this
- User is independent of an app
  - But user access can be restricted to one or more apps
  - User's permission depends on the app being accessed
  - An app has the authority to allow or disallow users
- User email will be verified with auto-generated mail
  

### Group
- Groups are associated with permissions
- Groups are specific to apps
- An user can have one or more groups associated
- When an user has more than one group, the permissions are the union of perms associated with the groups (i.e logical OR)


### Permissions

### Authentication
- At present only key and secret authentication is supported
- User authenticates with idx with user-name and password
- Service authenticates with its name and a secret token (this can be roatated automatically)

### Authorization
- Authorization happens per service
- Service asks idx about permissions a authenticated user has WRT that service
- 
