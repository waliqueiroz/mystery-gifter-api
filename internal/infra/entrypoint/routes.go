package entrypoint

import (
	"github.com/gofiber/fiber/v2"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/entrypoint/rest"
)

func CreateRoutes(router fiber.Router, authMiddleware fiber.Handler, userController *rest.UserController, authController *rest.AuthController, groupController *rest.GroupController) {
	api := router.Group("/api/v1")

	// swagger:operation POST /api/v1/login Login
	//
	// Authenticate user and get access token
	//
	// This endpoint authenticates a user with email and password and returns a JWT token.
	//
	// ---
	// tags:
	// - auth
	// produces:
	// - application/json
	// consumes:
	// - application/json
	// parameters:
	// - name: CredentialsDTO
	//   in: body
	//   description: User credentials for authentication
	//   required: true
	//   schema:
	//     "$ref": '#/definitions/CredentialsDTO'
	// responses:
	//   '200':
	//     description: Authentication successful
	//     schema:
	//       "$ref": '#/definitions/AuthSessionDTO'
	//   '400':
	//     description: Invalid credentials
	//   '401':
	//     description: Authentication failed
	//   '422':
	//     description: Invalid request body
	api.Post("/login", authController.Login)

	// swagger:operation POST /api/v1/users CreateUser
	//
	// Create a new user
	//
	// This endpoint creates a new user account with the provided information.
	//
	// ---
	// tags:
	// - users
	// produces:
	// - application/json
	// consumes:
	// - application/json
	// parameters:
	// - name: CreateUserDTO
	//   in: body
	//   description: User information for account creation
	//   required: true
	//   schema:
	//     "$ref": '#/definitions/CreateUserDTO'
	// responses:
	//   '201':
	//     description: User created successfully
	//     schema:
	//       "$ref": '#/definitions/UserDTO'
	//   '400':
	//     description: Invalid user data
	//   '409':
	//     description: User already exists
	//   '422':
	//     description: Invalid request body
	api.Post("/users", userController.Create)

	api.Use(authMiddleware) // from now on, all routes will require authentication

	// swagger:operation GET /api/v1/users SearchUsers
	//
	// Search users with filters and pagination
	//
	// This endpoint searches for users based on filters and returns paginated results.
	// Requires authentication.
	//
	// ---
	// tags:
	// - users
	// produces:
	// - application/json
	// security:
	// - Bearer: []
	// parameters:
	// - name: name
	//   in: query
	//   description: Filter by user name
	//   required: false
	//   type: string
	// - name: email
	//   in: query
	//   description: Filter by user email
	//   required: false
	//   type: string
	// - name: limit
	//   in: query
	//   description: Number of results per page
	//   required: false
	//   type: integer
	// - name: offset
	//   in: query
	//   description: Number of results to skip
	//   required: false
	//   type: integer
	// - name: sort
	//   in: query
	//   description: Sort direction (asc, desc)
	//   required: false
	//   type: string
	// - name: sort_field
	//   in: query
	//   description: Field to sort by
	//   required: false
	//   type: string
	// responses:
	//   '200':
	//     description: Search completed successfully
	//     schema:
	//       "$ref": '#/definitions/UserSearchResultDTO'
	//   '400':
	//     description: Invalid search parameters
	//   '401':
	//     description: Authentication required
	//   '422':
	//     description: Invalid query parameters
	api.Get("/users", userController.Search)

	// swagger:operation GET /api/v1/users/{userID} GetUserByID
	//
	// Get user by ID
	//
	// This endpoint retrieves a specific user by their ID.
	// Requires authentication.
	//
	// ---
	// tags:
	// - users
	// produces:
	// - application/json
	// security:
	// - Bearer: []
	// parameters:
	// - name: userID
	//   in: path
	//   description: Unique user identifier
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//     description: User found successfully
	//     schema:
	//       "$ref": '#/definitions/UserDTO'
	//   '401':
	//     description: Authentication required
	//   '404':
	//     description: User not found
	api.Get("/users/:userID", userController.GetByID)

	// swagger:operation GET /api/v1/groups SearchGroups
	//
	// Search groups with filters and pagination
	//
	// This endpoint searches for groups based on filters and returns paginated results.
	// Requires authentication.
	//
	// ---
	// tags:
	// - groups
	// produces:
	// - application/json
	// security:
	// - Bearer: []
	// parameters:
	// - name: name
	//   in: query
	//   description: Filter by group name
	//   required: false
	//   type: string
	// - name: owner_id
	//   in: query
	//   description: Filter by group owner ID
	//   required: false
	//   type: string
	// - name: status
	//   in: query
	//   description: Filter by group status
	//   required: false
	//   type: string
	// - name: limit
	//   in: query
	//   description: Number of results per page
	//   required: false
	//   type: integer
	// - name: offset
	//   in: query
	//   description: Number of results to skip
	//   required: false
	//   type: integer
	// - name: sort
	//   in: query
	//   description: Sort direction (asc, desc)
	//   required: false
	//   type: string
	// - name: sort_field
	//   in: query
	//   description: Field to sort by
	//   required: false
	//   type: string
	// responses:
	//   '200':
	//     description: Search completed successfully
	//     schema:
	//       "$ref": '#/definitions/GroupSearchResultDTO'
	//   '400':
	//     description: Invalid search parameters
	//   '401':
	//     description: Authentication required
	//   '422':
	//     description: Invalid query parameters
	api.Get("/groups", groupController.Search)

	// swagger:operation POST /api/v1/groups CreateGroup
	//
	// Create a new group
	//
	// This endpoint creates a new secret santa group.
	// The authenticated user becomes the group owner.
	//
	// ---
	// tags:
	// - groups
	// produces:
	// - application/json
	// consumes:
	// - application/json
	// security:
	// - Bearer: []
	// parameters:
	// - name: CreateGroupDTO
	//   in: body
	//   description: Group information for creation
	//   required: true
	//   schema:
	//     "$ref": '#/definitions/CreateGroupDTO'
	// responses:
	//   '201':
	//     description: Group created successfully
	//     schema:
	//       "$ref": '#/definitions/GroupDTO'
	//   '400':
	//     description: Invalid group data
	//   '401':
	//     description: Authentication required
	//   '422':
	//     description: Invalid request body
	api.Post("/groups", groupController.Create)

	// swagger:operation GET /api/v1/groups/{groupID} GetGroupByID
	//
	// Get group by ID
	//
	// This endpoint retrieves a specific group by its ID.
	// Requires authentication.
	//
	// ---
	// tags:
	// - groups
	// produces:
	// - application/json
	// security:
	// - Bearer: []
	// parameters:
	// - name: groupID
	//   in: path
	//   description: Unique group identifier
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//     description: Group found successfully
	//     schema:
	//       "$ref": '#/definitions/GroupDTO'
	//   '401':
	//     description: Authentication required
	//   '404':
	//     description: Group not found
	api.Get("/groups/:groupID", groupController.GetByID)

	// swagger:operation POST /api/v1/groups/{groupID}/users AddUserToGroup
	//
	// Add user to group
	//
	// This endpoint adds a user to an existing group.
	// Only the group owner can add users.
	//
	// ---
	// tags:
	// - groups
	// produces:
	// - application/json
	// consumes:
	// - application/json
	// security:
	// - Bearer: []
	// parameters:
	// - name: groupID
	//   in: path
	//   description: Unique group identifier
	//   required: true
	//   type: string
	// - name: AddUserDTO
	//   in: body
	//   description: User information to add to group
	//   required: true
	//   schema:
	//     "$ref": '#/definitions/AddUserDTO'
	// responses:
	//   '200':
	//     description: User added successfully
	//     schema:
	//       "$ref": '#/definitions/GroupDTO'
	//   '400':
	//     description: Invalid request data
	//   '401':
	//     description: Authentication required
	//   '403':
	//     description: Insufficient permissions
	//   '404':
	//     description: Group not found
	//   '409':
	//     description: User already in group
	//   '422':
	//     description: Invalid request body
	api.Post("/groups/:groupID/users", groupController.AddUser)

	// swagger:operation DELETE /api/v1/groups/{groupID}/users/{userID} RemoveUserFromGroup
	//
	// Remove user from group
	//
	// This endpoint removes a user from a group.
	// Only the group owner can remove users.
	//
	// ---
	// tags:
	// - groups
	// produces:
	// - application/json
	// security:
	// - Bearer: []
	// parameters:
	// - name: groupID
	//   in: path
	//   description: Unique group identifier
	//   required: true
	//   type: string
	// - name: userID
	//   in: path
	//   description: Unique user identifier
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//     description: User removed successfully
	//     schema:
	//       "$ref": '#/definitions/GroupDTO'
	//   '401':
	//     description: Authentication required
	//   '403':
	//     description: Insufficient permissions
	//   '404':
	//     description: Group or user not found
	api.Delete("/groups/:groupID/users/:userID", groupController.RemoveUser)

	// swagger:operation POST /api/v1/groups/{groupID}/matches GenerateMatches
	//
	// Generate matches for the group
	//
	// This endpoint generates random matches between users in the group.
	// Only the group owner can generate matches.
	//
	// ---
	// tags:
	// - groups
	// produces:
	// - application/json
	// security:
	// - Bearer: []
	// parameters:
	// - name: groupID
	//   in: path
	//   description: Unique group identifier
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//     description: Matches generated successfully
	//     schema:
	//       "$ref": '#/definitions/GroupDTO'
	//   '401':
	//     description: Authentication required
	//   '403':
	//     description: Insufficient permissions
	//   '404':
	//     description: Group not found
	//   '409':
	//     description: Cannot generate matches (insufficient users)
	api.Post("/groups/:groupID/matches", groupController.GenerateMatches)

	// swagger:operation POST /api/v1/groups/{groupID}/reopen ReopenGroup
	//
	// Reopen an archived group
	//
	// This endpoint reopens an archived group.
	// Only the group owner can reopen groups.
	//
	// ---
	// tags:
	// - groups
	// produces:
	// - application/json
	// security:
	// - Bearer: []
	// parameters:
	// - name: groupID
	//   in: path
	//   description: Unique group identifier
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//     description: Group reopened successfully
	//     schema:
	//       "$ref": '#/definitions/GroupDTO'
	//   '401':
	//     description: Authentication required
	//   '403':
	//     description: Insufficient permissions
	//   '404':
	//     description: Group not found
	//   '409':
	//     description: Group cannot be reopened
	api.Post("/groups/:groupID/reopen", groupController.Reopen)

	// swagger:operation POST /api/v1/groups/{groupID}/archive ArchiveGroup
	//
	// Archive a group
	//
	// This endpoint archives a group.
	// Only the group owner can archive groups.
	//
	// ---
	// tags:
	// - groups
	// produces:
	// - application/json
	// security:
	// - Bearer: []
	// parameters:
	// - name: groupID
	//   in: path
	//   description: Unique group identifier
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//     description: Group archived successfully
	//     schema:
	//       "$ref": '#/definitions/GroupDTO'
	//   '401':
	//     description: Authentication required
	//   '403':
	//     description: Insufficient permissions
	//   '404':
	//     description: Group not found
	//   '409':
	//     description: Group cannot be archived
	api.Post("/groups/:groupID/archive", groupController.Archive)

	// swagger:operation GET /api/v1/groups/{groupID}/matches/user GetUserMatch
	//
	// Get user's match in the group
	//
	// This endpoint returns the user that the authenticated user should give a gift to.
	// Requires authentication and the user must be a member of the group.
	//
	// ---
	// tags:
	// - groups
	// produces:
	// - application/json
	// security:
	// - Bearer: []
	// parameters:
	// - name: groupID
	//   in: path
	//   description: Unique group identifier
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//     description: User match found successfully
	//     schema:
	//       "$ref": '#/definitions/UserDTO'
	//   '401':
	//     description: Authentication required
	//   '403':
	//     description: User not a member of group
	//   '404':
	//     description: Group not found or no match available
	api.Get("/groups/:groupID/matches/user", groupController.GetUserMatch)
}
