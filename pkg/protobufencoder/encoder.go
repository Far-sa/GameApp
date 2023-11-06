package protobufencoder

import (
	"encoding/base64"
	"game-app/contract/golang/matching"
	"game-app/entity"
	"game-app/pkg/slice"

	"google.golang.org/protobuf/proto"
)

func EncodeMatchedUsers(mu entity.MatchUsers) string {

	pbMu := matching.MatchUsers{
		Category: string(mu.Category),
		UserIds:  slice.MapFromUintToUint64(mu.UserIDs),
	}

	payload, err := proto.Marshal(&pbMu)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(payload)

}
