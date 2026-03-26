# bouncr-client-go

A Go client library for the [Bouncr](https://github.com/kawasima/bouncr) identity and access management API.

## Installation

```sh
go get github.com/kawasima/bouncr-client-go
```

## Usage

Create a client with the `clientId` and `clientSecret` of an OIDC application registered in Bouncr. The client authenticates using the OAuth2 `client_credentials` grant.

```go
package main

import (
    "context"
    "fmt"
    "log"

    bouncr "github.com/kawasima/bouncr-client-go"
)

func main() {
    client, err := bouncr.NewClientWithOptions(
        "your-client-id",
        "your-client-secret",
        "http://localhost:3000", // Bouncr base URL
        false,                  // verbose logging
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // List users
    users, err := client.ListUsers(ctx, &bouncr.UserSearchParams{})
    if err != nil {
        log.Fatal(err)
    }
    for _, u := range users {
        fmt.Printf("ID=%d Account=%s\n", u.ID, u.Account)
    }
}
```

### Authentication

The client obtains an access token automatically via `POST /oauth2/token` before the first API call. Tokens are cached and refreshed transparently when they expire.

The OIDC application must have `client_credentials` included in its grant types and the required permissions assigned as scopes.

### Token endpoint override

If the OAuth2 token endpoint is not reachable at `{BaseURL}/oauth2/token`, you can override it:

```go
tokenURL, _ := url.Parse("http://localhost:3005/oauth2/token")
client.TokenURL = tokenURL
```

## API

All methods take a `context.Context` as the first argument.

### Users

| Method | Description |
| ------ | ----------- |
| `FindUser(ctx, account)` | Find a user by account name |
| `ListUsers(ctx, params)` | List users with optional search/pagination |
| `CreateUser(ctx, request)` | Create a user |
| `UpdateUser(ctx, account, request)` | Update a user |
| `DeleteUser(ctx, account)` | Delete a user |

### Groups

| Method | Description |
| ------ | ----------- |
| `FindGroup(ctx, name)` | Find a group by name |
| `ListGroups(ctx, params)` | List groups |
| `CreateGroup(ctx, request)` | Create a group |
| `UpdateGroup(ctx, name, request)` | Update a group |
| `DeleteGroup(ctx, name)` | Delete a group |
| `FindUsersInGroup(ctx, name)` | List users in a group |
| `AddUsersToGroup(ctx, group, accounts)` | Add users to a group |
| `RemoveUsersFromGroup(ctx, group, accounts)` | Remove users from a group |

### Roles

| Method | Description |
| ------ | ----------- |
| `FindRole(ctx, name)` | Find a role by name |
| `ListRoles(ctx, params)` | List roles |
| `CreateRole(ctx, request)` | Create a role |
| `UpdateRole(ctx, name, request)` | Update a role |
| `DeleteRole(ctx, name)` | Delete a role |
| `FindPermissionsInRole(ctx, name)` | List permissions in a role |
| `AddPermissionsToRole(ctx, role, permissions)` | Add permissions to a role |
| `RemovePermissionsFromRole(ctx, role, permissions)` | Remove permissions from a role |

### Permissions

| Method | Description |
| ------ | ----------- |
| `FindPermission(ctx, name)` | Find a permission by name |
| `ListPermissions(ctx, params)` | List permissions |
| `CreatePermission(ctx, request)` | Create a permission |
| `UpdatePermission(ctx, name, request)` | Update a permission |
| `DeletePermission(ctx, name)` | Delete a permission |

### Applications

| Method | Description |
| ------ | ----------- |
| `FindApplication(ctx, name)` | Find an application by name |
| `ListApplications(ctx, params)` | List applications |
| `CreateApplication(ctx, request)` | Create an application |
| `UpdateApplication(ctx, name, request)` | Update an application |
| `DeleteApplication(ctx, name)` | Delete an application |

### Realms

| Method | Description |
| ------ | ----------- |
| `FindRealm(ctx, appName, name)` | Find a realm within an application |
| `ListRealms(ctx, appName, params)` | List realms within an application |
| `CreateRealm(ctx, appName, request)` | Create a realm |
| `UpdateRealm(ctx, appName, name, request)` | Update a realm |
| `DeleteRealm(ctx, appName, name)` | Delete a realm |

### Assignments

| Method | Description |
| ------ | ----------- |
| `FindAssignment(ctx, request)` | Find an assignment |
| `CreateAssignments(ctx, request)` | Create assignments |
| `DeleteAssignments(ctx, request)` | Delete assignments |

### OIDC Applications

| Method | Description |
| ------ | ----------- |
| `FindOidcApplication(ctx, name)` | Find an OIDC application by name |
| `ListOidcApplications(ctx, params)` | List OIDC applications |
| `CreateOidcApplication(ctx, request)` | Create an OIDC application |
| `UpdateOidcApplication(ctx, name, request)` | Update an OIDC application |
| `DeleteOidcApplication(ctx, name)` | Delete an OIDC application |
| `RegenerateOidcApplicationSecret(ctx, name)` | Regenerate the client secret |

### OIDC Providers

| Method | Description |
| ------ | ----------- |
| `FindOidcProvider(ctx, name)` | Find an OIDC provider by name |
| `ListOidcProviders(ctx, params)` | List OIDC providers |
| `CreateOidcProvider(ctx, request)` | Create an OIDC provider |
| `UpdateOidcProvider(ctx, name, request)` | Update an OIDC provider |
| `DeleteOidcProvider(ctx, name)` | Delete an OIDC provider |

### Password Credentials

| Method | Description |
| ------ | ----------- |
| `CreatePasswordCredential(ctx, request)` | Create a password credential |
| `UpdatePasswordCredential(ctx, request)` | Update a password credential |
| `DeletePasswordCredential(ctx)` | Delete a password credential |

## License

Eclipse Public License 2.0
