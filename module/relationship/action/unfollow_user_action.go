package action

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/HotPotatoC/twitter-clone/module"
	"github.com/HotPotatoC/twitter-clone/module/relationship/entity"
	"github.com/HotPotatoC/twitter-clone/module/relationship/service"
	userEntity "github.com/HotPotatoC/twitter-clone/module/user/entity"
	"github.com/gofiber/fiber/v2"
)

type unfollowUserAction struct {
	service service.FollowUserService
}

func NewUnfollowUserAction(service service.FollowUserService) module.Action {
	return unfollowUserAction{service: service}
}

func (a unfollowUserAction) Execute(c *fiber.Ctx) error {
	followerID := c.Locals("userID").(float64)
	followedID, err := strconv.ParseInt(c.Params("userID"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request",
		})
	}

	if int64(followerID) == followedID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "You cannot unfollow yourself",
		})
	}

	username, err := a.service.Execute(int64(followerID), followedID)
	if err != nil {
		switch {
		case errors.Is(err, userEntity.ErrUserDoesNotExist):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User not found",
			})
		case errors.Is(err, entity.ErrUserIsNotFollowing):
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "You are not following that user",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "There was a problem on our side",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": fmt.Sprintf("Unfollowed %s", username),
	})
}