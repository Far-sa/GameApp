package protobufencoder

import (
	"encoding/base64"
	"game-app/contract/golang/matching"
	"game-app/entity"
	"game-app/pkg/slice"

	"google.golang.org/protobuf/proto"
)

// ! prefer to use case by case due to loose type saftey if ..
func EncodeMatchedUsersEvent(mu entity.MatchUsers) string {

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

// * general way for support all types of events (fat function)
func EncodeEvent(event entity.Event, data any) string {
	var payload []byte
	switch event {
	case entity.MatchingUserEvent:
		mu, ok := data.(entity.MatchUsers)
		if !ok {
			// TODO: log error
			// TODO: update metrics
			return ""
		}

		pbMu := matching.MatchUsers{
			Category: string(mu.Category),
			UserIds:  slice.MapFromUintToUint64(mu.UserIDs),
		}

		var err error
		payload, err = proto.Marshal(&pbMu)
		if err != nil {
			// TODO: log error
			// TODO: update metrics
			return ""
		}
	}

	return base64.StdEncoding.EncodeToString(payload)

}

func DecoderEvent(event entity.Event, data string) any {

	payload, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		// TODO: log error
		// TODO: update metrics
		return nil
	}

	switch event {
	case entity.MatchingUserEvent:

		pbMu := matching.MatchUsers{}

		if err := proto.Unmarshal(payload, &pbMu); err != nil {
			// TODO: log error
			// TODO: update metrics
			return nil
		}

		return entity.MatchUsers{
			Category: entity.Category(pbMu.Category),
			UserIDs:  slice.MapFromUint64ToUint(pbMu.UserIds),
		}
	}

	return nil

}
