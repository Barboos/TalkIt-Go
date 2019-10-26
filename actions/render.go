package actions

import (
	timeago "github.com/ararog/timeago"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr/v2"
	"github.com/zoonman/gravatar"
	"strings"
	"time"
	"github.com/gobuffalo/uuid"
	"s_app/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
)

var r *render.Engine
var assetsBox = packr.New("app:assets", "../public")

func init() {
	r = render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: "application.html",

		// Box containing all of the templates:
		TemplatesBox: packr.New("app:templates", "../templates"),
		AssetsBox:    assetsBox,

		// Add template helpers here:
		Helpers: render.Helpers{
			// for non-bootstrap form helpers uncomment the lines
			// below and import "github.com/gobuffalo/helpers/forms"
			// forms.FormKey:     forms.Form,
			// forms.FormForKey:  forms.FormFor,

			"gravatar_for": func(user string, size int) string {
				return gravatar.Avatar(user, uint(size))
			},

			"time_ago_in_words": func(created_at time.Time) string {
				got, _ := timeago.TimeAgoWithTime(time.Now(), created_at)
				return strings.ToLower(got)
			},

			"isFollowing": func(user uuid.UUID, other_user uuid.UUID, c buffalo.Context) bool{
				relationship := &models.Relationship{}
				tx, _ := c.Value("tx").(*pop.Connection)
				exist, _ := tx.Where("followed_id in (?)", other_user).Where("follower_id in (?)", user).Exists(relationship)
				return exist
			},
		},
	})
}
