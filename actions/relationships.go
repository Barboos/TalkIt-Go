package actions

import (
	"s_app/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

// RelationshipsCreate registers a new relationship with the application.
func RelationshipsCreate(c buffalo.Context) error {

	uid := c.Session().Get("current_user_id")
	current_user := &models.User{}
	tx := c.Value("tx").(*pop.Connection)
	err := tx.Find(current_user, uid)
  if err != nil {
    return errors.WithStack(err)
  }

  user := &models.User{}
  err = tx.Find(user, c.Param("user_id"))
  if err != nil {
    return errors.WithStack(err)
  }

  relationship := &models.Relationship{
    FollowerID: current_user.ID,
    FollowedID: user.ID,
  }

	// Validate the data from the html form
	_, err = tx.ValidateAndCreate(relationship)
	if err != nil {
		return err
	}

	return c.Redirect(302, "/")
}

// RelationshipsDestroy deletes a Relationship from the DB. This function is mapped
// to the path DELETE /microposts/{micropost_id}
func RelationshipsDestroy(c buffalo.Context) error {
  uid := c.Session().Get("current_user_id")
	current_user := &models.User{}
	tx := c.Value("tx").(*pop.Connection)
	err := tx.Find(current_user, uid)
  if err != nil {
    return errors.WithStack(err)
  }

  user := &models.User{}
  err = tx.Find(user, c.Param("user_id"))
  if err != nil {
    return errors.WithStack(err)
  }

	relationship := &models.Relationships{}

	if err := tx.Where("followed_id in (?)", user.ID).Where("follower_id in (?)", current_user.ID).All(relationship); err != nil {
		return c.Error(404, err)
	}
  
	if err := tx.Destroy(relationship); err != nil {
		return err
	}

	return c.Redirect(302, "/")
}

func RelationshipsList(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	relationships := &models.Relationships{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all relationships from the DB
	if err := q.All(relationships); err != nil {
		return err
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)
	c.Set("relationships", relationships)
	return c.Render(200, r.HTML("/relationships.html"))
}
