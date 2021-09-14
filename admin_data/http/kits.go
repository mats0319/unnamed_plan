package http

import (
    "fmt"
    "github.com/mats9693/unnamed_plan/admin_data/db/model"
    "github.com/pkg/errors"
)

func sortUsersByUserID(users []*model.User, order []string) ([]*model.User, error) {
    if len(users) != len(order) {
        return nil, errors.New(fmt.Sprintf("unmatched users amount, users %d, orders %d", len(users), len(order)))
    }

    length := len(users)
    for i := 0; i < length; i++ {
        for j := i; j < length; j++ {
            if order[j] == users[i].UserID {
                users[i], users[j] = users[j], users[i]
                break
            }
        }
    }

    unmatchedIndex := -1
    for i := 0; i < length; i++ {
        if users[i].UserID != order[i] {
            unmatchedIndex = i
            break
        }
    }

    if unmatchedIndex >= 0 {
        return nil, errors.New(fmt.Sprintf("unmatched user id: %s", users[unmatchedIndex].UserID))
    }

    return users, nil
}
