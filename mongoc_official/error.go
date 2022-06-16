package mongoc_official

import (
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

// IsDup returns whether err informs of a duplicate key error because
// a primary key index or a secondary unique index already has an entry
// with the given value.
func isDupCode(code int, msg string) bool {
	return code == 11000 || code == 11001 || code == 12582 || code == 16460 && strings.Contains(msg, " E11000 ")
}

func IsDup(err error) bool {
	switch er := err.(type) {
	case mongo.WriteException:
		for _, e := range er.WriteErrors {
			if isDupCode(e.Code, e.Message) {
				return true
			}
		}
	case mongo.BulkWriteException:
		for _, e := range er.WriteErrors {
			if isDupCode(e.Code, e.Message) {
				return true
			}
		}
	}
	return false
}
