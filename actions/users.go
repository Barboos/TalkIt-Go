package actions

import (
	"s_app/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
	//"github.com/gobuffalo/uuid"
)

func UsersNew(c buffalo.Context) error {
	u := models.User{}
	c.Set("user", u)
	return c.Render(200, r.HTML("users/new.html"))
}

func UsersShow(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}
	// Allocate an empty User
	user := &models.User{}
	micropost := models.Micropost{}
	c.Set("Micropost", micropost)
	// To find the User the parameter user_id is used.
	if err := tx.Eager("Microposts").Find(user, c.Param("user_id")); err != nil {
		return c.Error(404, err)
	}
	c.Set("current_user", user)
	q := tx.PaginateFromParams(c.Params())
	microposts := &models.Microposts{}
	q.Order("created_at desc").BelongsTo(user).All(microposts)
	c.Set("microposts", microposts)
	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	following := &models.Relationships{}
	tx.Where("follower_id in (?)", user.ID).All(following)
	c.Set("following", following)

	followers := &models.Relationships{}
	tx.Where("followed_id in (?)", user.ID).All(followers)
	c.Set("followers", followers)

	relationship := models.Relationship{}
	c.Set("relationship", relationship)

	c.Set("c", c)
	root := &models.User{}
	uid := c.Session().Get("current_user_id")
	tx.Eager("Microposts").Find(root, uid)
	c.Set("root_user", root)

	return c.Render(200, r.HTML("users/show.html"))
}

// UsersCreate registers a new user with the application.
func UsersCreate(c buffalo.Context) error {
	u := &models.User{}
	if err := c.Bind(u); err != nil {
		return errors.WithStack(err)
	}
	tx := c.Value("tx").(*pop.Connection)
	verrs, err := u.Create(tx)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		c.Set("user", u)
		c.Set("errors", verrs)
		return c.Render(200, r.HTML("users/new.html"))
	}
	c.Session().Set("current_user_id", u.ID)
	c.Flash().Add("success", "Welcome to Buffalo!")
	return c.Redirect(302, "/")
}

func UsersList(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}
	users := &models.Users{}
	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())
	// Retrieve all Users from the DB
	if err := q.All(users); err != nil {
		return err
	}
	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)
	c.Set("users", users)
	return c.Render(200, r.HTML("users/index.html"))
}

// Edit renders a edit form for a User. This function is
// mapped to the path GET /users/{user_id}/edit
func UsersEdit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	// Allocate an empty User
	user := &models.User{}

	if err := tx.Find(user, c.Param("user_id")); err != nil {
		return c.Error(404, err)
	}
	c.Set("user", user)

	return c.Render(200, r.HTML("users/edit.html"))
}

// Update changes a User in the DB. This function is mapped to
// the path PUT /users/{user_id}
func UsersUpdate(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}
	// Allocate an empty User
	user := &models.User{}
	if err := tx.Find(user, c.Param("user_id")); err != nil {
		return c.Error(404, err)
	}
	// Bind User to the html form elements
	if err := c.Bind(user); err != nil {
		return err
	}
	verrs, err := tx.ValidateAndUpdate(user)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)
		// Render again the edit.html template that the user can
		// correct the input.
		c.Set("user", user)
		return c.Render(422, r.HTML("/users/edit.html"))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", T.Translate(c, "user.updated.success"))
	// and redirect to the user page
	c.Set("user", user)
	return c.Redirect(302, "/")
}

// SetCurrentUser attempts to find a user based on the current_user_id
// in the session. If one is found it is set on the context.
func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid != nil {
			u := &models.User{}
			micropost := models.Micropost{}
			c.Set("Micropost", micropost)
			tx := c.Value("tx").(*pop.Connection)
			err := tx.Eager("Microposts").Find(u, uid)
			if err != nil {
				return errors.WithStack(err)
			}
			c.Set("current_user", u)

			following := &models.Relationships{}
			tx.Where("follower_id in (?)", u.ID).All(following)
			c.Set("following", following)

			followers := &models.Relationships{}
			tx.Where("followed_id in (?)", u.ID).All(followers)
			c.Set("followers", followers)

			// Paginate results. Params "page" and "per_page" control pagination.
			// Default values are "page=1" and "per_page=20".
			q := tx.PaginateFromParams(c.Params())
			microposts := &models.Microposts{}
			//ids := &models.Relationships{}
			//tx.Select("followed_id").Where("follower_id in (?)", u.ID).All(ids)
			q.Order("created_at desc").Where("user_id in (SELECT followed_id FROM relationships WHERE  follower_id = (?)) OR user_id = (?)", u.ID, u.ID).Eager("User").All(microposts)
			c.Set("microposts", microposts)
			// Add the paginator to the context so it can be used in the template.
			c.Set("pagination", q.Paginator)
		}
		return next(c)
	}
}

// Authorize require a user be logged in before accessing a route
func Authorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid == nil {
			c.Session().Set("redirectURL", c.Request().URL.String())

			err := c.Session().Save()
			if err != nil {
				return errors.WithStack(err)
			}

			c.Flash().Add("danger", "You must be authorized to see that page")
			return c.Redirect(302, "/")
		}
		return next(c)
	}
}
